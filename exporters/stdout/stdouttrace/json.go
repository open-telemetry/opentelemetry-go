// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace

import (
	"math"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// These wrappers preserve the typed debug representation while allowing
// non-finite attribute values that encoding/json cannot encode as numbers.
type jsonValue struct {
	Type  string
	Value any
}

type jsonAttribute struct {
	Key   attribute.Key
	Value jsonValue
}

func valueJSON(v attribute.Value) jsonValue {
	var value any
	switch v.Type() {
	case attribute.FLOAT64:
		value = floatJSON(v.AsFloat64())
	case attribute.FLOAT64SLICE:
		floats := v.AsFloat64Slice()
		values := make([]any, len(floats))
		for i, f := range floats {
			values[i] = floatJSON(f)
		}
		value = values
	case attribute.SLICE:
		slice := v.AsSlice()
		values := make([]jsonValue, len(slice))
		for i, item := range slice {
			values[i] = valueJSON(item)
		}
		value = values
	case attribute.MAP:
		value = attributesJSON(v.AsMap())
	default:
		value = v.AsInterface()
	}
	return jsonValue{Type: v.Type().String(), Value: value}
}

func floatJSON(f float64) any {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	default:
		return f
	}
}

func attributesJSON(attrs []attribute.KeyValue) []jsonAttribute {
	if attrs == nil {
		return nil
	}
	result := make([]jsonAttribute, len(attrs))
	for i, attr := range attrs {
		result[i] = jsonAttribute{Key: attr.Key, Value: valueJSON(attr.Value)}
	}
	return result
}

type jsonScope struct {
	instrumentation.Scope
	Attributes any
}

func scopeJSON(scope instrumentation.Scope) jsonScope {
	if scope.Attributes.Len() == 0 {
		// Let Set preserve the difference between zero and initialized empty sets.
		return jsonScope{Scope: scope, Attributes: &scope.Attributes}
	}
	return jsonScope{Scope: scope, Attributes: attributesJSON(scope.Attributes.ToSlice())}
}

type jsonEvent struct {
	trace.Event
	Attributes []jsonAttribute
}

type jsonLink struct {
	trace.Link
	Attributes []jsonAttribute
}

type jsonSpan struct {
	*tracetest.SpanStub
	Attributes             []jsonAttribute
	Events                 []jsonEvent
	Links                  []jsonLink
	Resource               []jsonAttribute
	InstrumentationScope   jsonScope
	InstrumentationLibrary jsonScope
}

func spanJSON(stub *tracetest.SpanStub) jsonSpan {
	result := jsonSpan{
		SpanStub:               stub,
		Attributes:             attributesJSON(stub.Attributes),
		Resource:               attributesJSON(stub.Resource.Attributes()),
		InstrumentationScope:   scopeJSON(stub.InstrumentationScope),
		InstrumentationLibrary: scopeJSON(stub.InstrumentationLibrary),
	}
	if stub.Events != nil {
		result.Events = make([]jsonEvent, len(stub.Events))
		for i, event := range stub.Events {
			result.Events[i] = jsonEvent{Event: event, Attributes: attributesJSON(event.Attributes)}
		}
	}
	if stub.Links != nil {
		result.Links = make([]jsonLink, len(stub.Links))
		for i, link := range stub.Links {
			result.Links[i] = jsonLink{Link: link, Attributes: attributesJSON(link.Attributes)}
		}
	}
	return result
}
