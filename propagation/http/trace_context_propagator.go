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

package http

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	apipropagation "go.opentelemetry.io/otel/api/propagation/http"
)

const (
	supportedVersion         = 0
	maxVersion               = 254
	TraceparentHeader        = "Traceparent"
	CorrelationContextHeader = "Correlation-Context"
)

// TraceContextPropagator propagates SpanContext in W3C TraceContext format.
type (
	TraceContextPropagator       struct{}
	CorrelationContextPropagator struct{}
	// TODO: BaggageContextPropagator, likely similar to
	// CorrelationContextPropagator
	//
	// TODO: ContextPropagator, grouping all three
)

var (
	_              apipropagation.SpanContextPropagator  = TraceContextPropagator{}
	_              apipropagation.CorrelationsPropagator = CorrelationContextPropagator{}
	traceCtxRegExp                                       = regexp.MustCompile("^[0-9a-f]{2}-[a-f0-9]{32}-[a-f0-9]{16}-[a-f0-9]{2}-?")
)

func (TraceContextPropagator) Inject(sc core.SpanContext, supplier apipropagation.Supplier) {
	if sc.IsValid() {
		h := fmt.Sprintf("%.2x-%s-%.16x-%.2x",
			supportedVersion,
			sc.TraceIDString(),
			sc.SpanID,
			sc.TraceFlags&core.TraceFlagsSampled)
		supplier.Set(TraceparentHeader, h)
	}
}

func (TraceContextPropagator) Extract(supplier apipropagation.Supplier) core.SpanContext {
	h := supplier.Get(TraceparentHeader)
	if h == "" {
		return core.EmptySpanContext()
	}

	h = strings.Trim(h, "-")
	if !traceCtxRegExp.MatchString(h) {
		return core.EmptySpanContext()
	}

	sections := strings.Split(h, "-")
	if len(sections) < 4 {
		return core.EmptySpanContext()
	}

	if len(sections[0]) != 2 {
		return core.EmptySpanContext()
	}
	ver, err := hex.DecodeString(sections[0])
	if err != nil {
		return core.EmptySpanContext()
	}
	version := int(ver[0])
	if version > maxVersion {
		return core.EmptySpanContext()
	}

	if version == 0 && len(sections) != 4 {
		return core.EmptySpanContext()
	}

	if len(sections[1]) != 32 {
		return core.EmptySpanContext()
	}

	var sc core.SpanContext

	sc.TraceID, err = core.TraceIDFromHex(sections[1][:32])
	if err != nil {
		return core.EmptySpanContext()
	}

	if len(sections[2]) != 16 {
		return core.EmptySpanContext()
	}
	sc.SpanID, err = core.SpanIDFromHex(sections[2][:])
	if err != nil {
		return core.EmptySpanContext()
	}

	if len(sections[3]) != 2 {
		return core.EmptySpanContext()
	}
	opts, err := hex.DecodeString(sections[3])
	if err != nil || len(opts) < 1 || (version == 0 && opts[0] > 2) {
		return core.EmptySpanContext()
	}
	sc.TraceFlags = opts[0] &^ core.TraceFlagsUnused

	if !sc.IsValid() {
		return core.EmptySpanContext()
	}

	return sc
}

func (CorrelationContextPropagator) Inject(correlations dctx.Correlations, supplier apipropagation.Supplier) {
	firstIter := true
	var headerValueBuilder strings.Builder
	correlations.Foreach(func(kv dctx.Correlation) bool {
		if kv.HopLimit == dctx.NoPropagation {
			return true
		}
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

func (CorrelationContextPropagator) Extract(supplier apipropagation.Supplier) dctx.Correlations {
	correlationContext := supplier.Get(CorrelationContextHeader)
	if correlationContext == "" {
		return dctx.NewCorrelations()
	}

	contextValues := strings.Split(correlationContext, ",")
	keyValues := make([]dctx.Correlation, 0, len(contextValues))
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

		// TODO (skaris): properties defined https://w3c.github.io/correlation-context/, are currently
		// just put as part of the value.
		var trimmedValueWithProps strings.Builder
		trimmedValueWithProps.WriteString(trimmedValue)
		for _, prop := range valueAndProps[1:] {
			trimmedValueWithProps.WriteRune(';')
			trimmedValueWithProps.WriteString(prop)
		}

		keyValues = append(keyValues, dctx.Correlation{
			KeyValue: key.String(trimmedName, trimmedValueWithProps.String()),
			HopLimit: dctx.UnlimitedPropagation,
		})
	}
	return dctx.NewCorrelations(dctx.CorrelationsUpdate{
		MultiKV: keyValues,
	})
}
