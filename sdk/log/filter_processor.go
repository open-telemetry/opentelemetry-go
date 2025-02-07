// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"reflect"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

// filterProcessor uses reflect to support [go.opentelemetry.io/otel/sdk/log/xlog.FilterProcessor]
// via duck typing.
type filterProcessor struct {
	enabledFn reflect.Value
	paramType reflect.Type
}

func asFilterProccessor(p any) (filterProcessor, bool) {
	fp := reflect.ValueOf(p)
	m := fp.MethodByName("Enabled")
	if m == (reflect.Value{}) {
		// No Enabled method.
		return filterProcessor{}, false
	}
	mty := m.Type()
	if mty.NumOut() != 1 {
		// Should return one output parameter.
		return filterProcessor{}, false
	}
	if reflect.Bool != mty.Out(0).Kind() {
		// Should return bool.
		return filterProcessor{}, false
	}
	if mty.NumIn() != 2 {
		// Should have two input parameters.
		return filterProcessor{}, false
	}
	if mty.In(0) != reflect.TypeFor[context.Context]() {
		// Should have context.Context as first input paramater.
		return filterProcessor{}, false
	}
	// Duck typing of EnabledParameters
	pt := mty.In(1)
	if res, ok := pt.FieldByName("Resource"); !ok || res.Type != reflect.TypeFor[resource.Resource]() {
		// The second paramater should have Resource resource.Resource field.
		return filterProcessor{}, false
	}
	if res, ok := pt.FieldByName("InstrumentationScope"); !ok || res.Type != reflect.TypeFor[instrumentation.Scope]() {
		// The second paramater should have InstrumentationScope instrumentation.Scope field.
		return filterProcessor{}, false
	}
	if res, ok := pt.FieldByName("Severity"); !ok || res.Type != reflect.TypeFor[log.Severity]() {
		// The second paramater should have Severity log.Severity field.
		return filterProcessor{}, false
	}

	return filterProcessor{
		enabledFn: m,
		paramType: pt,
	}, true
}

func (f filterProcessor) Enabled(ctx context.Context, param enabledParameters) bool {
	p := reflect.New(f.paramType).Elem()
	p.FieldByName("Resource").Set(reflect.ValueOf(param.Resource))
	p.FieldByName("InstrumentationScope").Set(reflect.ValueOf(param.InstrumentationScope))
	p.FieldByName("Severity").Set(reflect.ValueOf(param.Severity))

	ctxV := reflect.ValueOf(ctx)
	if ctxV == (reflect.Value{}) {
		// In order to not get panic: reflect: Call using zero Value argument.
		ctxV = reflect.Zero(reflect.TypeFor[context.Context]())
	}

	ret := f.enabledFn.Call([]reflect.Value{ctxV, p})
	return ret[0].Bool()
}

// enabledParameters represents payload for Enabled method.
type enabledParameters struct {
	Resource             resource.Resource
	InstrumentationScope instrumentation.Scope
	Severity             log.Severity
}
