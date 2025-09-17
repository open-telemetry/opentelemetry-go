// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/sdk/trace/internal/observ"

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/trace/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

var measureAttrsPool = sync.Pool{
	New: func() any {
		// "component.name" + "component.type" + "error.type"
		const n = 1 + 1 + 1
		s := make([]attribute.KeyValue, 0, n)
		// Return a pointer to a slice instead of a slice itself
		// to avoid allocations on every call.
		return &s
	},
}

type SSP struct {
	componentNameAttr     attribute.KeyValue
	spansProcessedCounter otelconv.SDKProcessorSpanProcessed
}

// SSPComponentName returns the component name attribute for a
// SimpleSpanProcessor with the given ID.
func SSPComponentName(id int64) attribute.KeyValue {
	t := otelconv.ComponentTypeSimpleSpanProcessor
	name := fmt.Sprintf("%s/%d", t, id)
	return semconv.OTelComponentName(name)
}

func NewSSP(id int64) (*SSP, error) {
	if !x.Observability.Enabled() {
		return nil, nil
	}

	meter := otel.GetMeterProvider().Meter(
		ScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(SchemaURL),
	)
	componentName := SSPComponentName(id)
	spansProcessedCounter, err := otelconv.NewSDKProcessorSpanProcessed(meter)
	if err != nil {
		err = fmt.Errorf("failed to create SSP processed spans metric: %w", err)
	}
	return &SSP{
		componentNameAttr:     componentName,
		spansProcessedCounter: spansProcessedCounter,
	}, err
}

func (ssp *SSP) Record(ctx context.Context, count int64, err error) {
	attrs := measureAttrsPool.Get().(*[]attribute.KeyValue)
	defer func() {
		*attrs = (*attrs)[:0] // reset the slice for reuse
		measureAttrsPool.Put(attrs)
	}()
	if err != nil {
		*attrs = append(*attrs, semconv.ErrorType(err))
	}
	*attrs = append(*attrs,
		ssp.componentNameAttr,
		ssp.spansProcessedCounter.AttrComponentType(otelconv.ComponentTypeSimpleSpanProcessor))
	ssp.spansProcessedCounter.Add(ctx, count, *attrs...)
}
