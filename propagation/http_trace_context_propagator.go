// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package propagation

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	apipropagation "go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

const (
	supportedVersion         = 0
	maxVersion               = 254
	TraceparentHeader        = "Traceparent"
	CorrelationContextHeader = "Correlation-Context"
)

// HTTPTraceContextPropagator propagates SpanContext in W3C TraceContext format.
type HTTPTraceContextPropagator struct{}

var _ apipropagation.TextFormatPropagator = HTTPTraceContextPropagator{}
var traceCtxRegExp = regexp.MustCompile("^[0-9a-f]{2}-[a-f0-9]{32}-[a-f0-9]{16}-[a-f0-9]{2}-?")

func (hp HTTPTraceContextPropagator) Inject(ctx context.Context, supplier apipropagation.Supplier) {
	sc := trace.CurrentSpan(ctx).SpanContext()
	if sc.IsValid() {
		h := fmt.Sprintf("%.2x-%s-%.16x-%.2x",
			supportedVersion,
			sc.TraceIDString(),
			sc.SpanID,
			sc.TraceFlags&otel.TraceFlagsSampled)
		supplier.Set(TraceparentHeader, h)
	}

	correlationCtx := dctx.FromContext(ctx)
	firstIter := true
	var headerValueBuilder strings.Builder
	correlationCtx.Foreach(func(kv otel.KeyValue) bool {
		if !firstIter {
			headerValueBuilder.WriteRune(',')
		}
		firstIter = false
		headerValueBuilder.WriteString(url.QueryEscape(strings.TrimSpace((string)(kv.Key))))
		headerValueBuilder.WriteRune('=')
		headerValueBuilder.WriteString(url.QueryEscape(strings.TrimSpace(kv.Value.Emit())))
		return true
	})
	if headerValueBuilder.Len() > 0 {
		headerString := headerValueBuilder.String()
		supplier.Set(CorrelationContextHeader, headerString)
	}
}

func (hp HTTPTraceContextPropagator) Extract(
	ctx context.Context, supplier apipropagation.Supplier,
) (otel.SpanContext, dctx.Map) {
	return hp.extractSpanContext(ctx, supplier), hp.extractCorrelationCtx(ctx, supplier)
}

func (hp HTTPTraceContextPropagator) extractSpanContext(
	ctx context.Context, supplier apipropagation.Supplier,
) otel.SpanContext {
	h := supplier.Get(TraceparentHeader)
	if h == "" {
		return otel.EmptySpanContext()
	}

	h = strings.Trim(h, "-")
	if !traceCtxRegExp.MatchString(h) {
		return otel.EmptySpanContext()
	}

	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return otel.EmptySpanContext()
	}

	if len(sections[0]) != 2 {
		return otel.EmptySpanContext()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return otel.EmptySpanContext()
	}
	version := int(ver[0])
	if version > maxVersion {
		return otel.EmptySpanContext()
	}

	if version == 0 && len(sections) != 4 {
		return otel.EmptySpanContext()
	}

	if len(sections[1]) != 32 {
		return otel.EmptySpanContext()
	}

	var sc otel.SpanContext

	sc.TraceID, err = otel.TraceIDFromHex(sections[1][:32])
	if err != nil {
		return otel.EmptySpanContext()
	}

	if len(sections[2]) != 16 {
		return otel.EmptySpanContext()
	}
	sc.SpanID, err = otel.SpanIDFromHex(sections[2][:])
	if err != nil {
		return otel.EmptySpanContext()
	}

	if len(sections[3]) != 2 {
		return otel.EmptySpanContext()
	}
	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return otel.EmptySpanContext()
	}
	sc.TraceFlags = opts[0] &^ otel.TraceFlagsUnused

	if !sc.IsValid() {
		return otel.EmptySpanContext()
	}

	return sc
}

func (hp HTTPTraceContextPropagator) extractCorrelationCtx(ctx context.Context, supplier apipropagation.Supplier) dctx.Map {
	correlationContext := supplier.Get(CorrelationContextHeader)
	if correlationContext == "" {
		return dctx.NewEmptyMap()
	}

	contextValues := strings.Split(correlationContext, ",")
	keyValues := make([]otel.KeyValue, 0, len(contextValues))
	for _, contextValue := range contextValues {
		valueAndProps := strings.Split(contextValue, ";")
		if len(valueAndProps) < 1 {
			continue
		}
		nameValue := strings.Split(valueAndProps[0], "=")
		if len(nameValue) < 2 {
			continue
		}
		name, err := url.QueryUnescape(nameValue[0])
		if err != nil {
			continue
		}
		trimmedName := strings.TrimSpace(name)
		value, err := url.QueryUnescape(nameValue[1])
		if err != nil {
			continue
		}
		trimmedValue := strings.TrimSpace(value)

		// TODO (skaris): properties defiend https://w3c.github.io/correlation-context/, are currently
		// just put as part of the value.
		var trimmedValueWithProps strings.Builder
		trimmedValueWithProps.WriteString(trimmedValue)
		for _, prop := range valueAndProps[1:] {
			trimmedValueWithProps.WriteRune(';')
			trimmedValueWithProps.WriteString(prop)
		}

		keyValues = append(keyValues, key.New(trimmedName).String(trimmedValueWithProps.String()))
	}
	return dctx.NewMap(dctx.MapUpdate{
		MultiKV: keyValues,
	})
}

func (hp HTTPTraceContextPropagator) GetAllKeys() []string {
	return []string{TraceparentHeader, CorrelationContextHeader}
}
