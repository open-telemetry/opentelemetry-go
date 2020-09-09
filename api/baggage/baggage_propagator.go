// Copyright The OpenTelemetry Authors
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

package baggage

import (
	"context"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/label"
)

// Temporary header name until W3C finalizes format.
// https://github.com/open-telemetry/opentelemetry-specification/blob/18b2752ebe6c7f0cdd8c7b2bcbdceb0ae3f5ad95/specification/correlationcontext/api.md#header-name
const baggageHeader = "otcorrelations"

// Baggage propagates Key:Values in W3C CorrelationContext
// format.
// nolint:golint
type Baggage struct{}

var _ propagation.HTTPPropagator = Baggage{}

// DefaultHTTPPropagator returns the default context correlation HTTP
// propagator.
func DefaultHTTPPropagator() propagation.HTTPPropagator {
	return Baggage{}
}

// Inject implements HTTPInjector.
func (b Baggage) Inject(ctx context.Context, supplier propagation.HTTPSupplier) {
	baggageMap := MapFromContext(ctx)
	firstIter := true
	var headerValueBuilder strings.Builder
	baggageMap.Foreach(func(kv label.KeyValue) bool {
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
		supplier.Set(baggageHeader, headerString)
	}
}

// Extract implements HTTPExtractor.
func (b Baggage) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	baggage := supplier.Get(baggageHeader)
	if baggage == "" {
		return ctx
	}

	baggageValues := strings.Split(baggage, ",")
	keyValues := make([]label.KeyValue, 0, len(baggageValues))
	for _, baggageValue := range baggageValues {
		valueAndProps := strings.Split(baggageValue, ";")
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

		keyValues = append(keyValues, label.String(trimmedName, trimmedValueWithProps.String()))
	}

	if len(keyValues) > 0 {
		// Only update the context if valid values were found
		return ContextWithMap(ctx, NewMap(MapUpdate{
			MultiKV: keyValues,
		}))
	}

	return ctx
}

// GetAllKeys implements HTTPPropagator.
func (b Baggage) GetAllKeys() []string {
	return []string{baggageHeader}
}
