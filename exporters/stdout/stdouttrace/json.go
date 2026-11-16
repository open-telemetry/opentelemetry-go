// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// attributesJSON retains list ordering while encoding values with the
// OpenTelemetry representation for non-OTLP protocols.
func attributesJSON(attrs []attribute.KeyValue) []json.RawMessage {
	if attrs == nil {
		return nil
	}
	result := make([]json.RawMessage, len(attrs))
	for i, attr := range attrs {
		result[i] = json.RawMessage(attr.String())
	}
	return result
}

type jsonScope struct {
	instrumentation.Scope
	Attributes json.RawMessage
}

func scopeJSON(scope instrumentation.Scope) jsonScope {
	return jsonScope{Scope: scope, Attributes: json.RawMessage(scope.Attributes.String())}
}

type jsonEvent struct {
	trace.Event
	Attributes []json.RawMessage
}

type jsonLink struct {
	trace.Link
	Attributes []json.RawMessage
}

type jsonSpan struct {
	*tracetest.SpanStub
	Attributes             []json.RawMessage
	Events                 []jsonEvent
	Links                  []jsonLink
	Resource               json.RawMessage
	InstrumentationScope   jsonScope
	InstrumentationLibrary jsonScope
}

func spanJSON(stub *tracetest.SpanStub) jsonSpan {
	scope := scopeJSON(stub.InstrumentationScope)
	result := jsonSpan{
		SpanStub:               stub,
		Attributes:             attributesJSON(stub.Attributes),
		InstrumentationScope:   scope,
		InstrumentationLibrary: scope,
	}
	if stub.Resource != nil {
		result.Resource = json.RawMessage(stub.Resource.Set().String())
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
