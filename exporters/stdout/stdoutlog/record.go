// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// recordJSON is a JSON-serializable representation of a Record.
type recordJSON struct {
	Timestamp         *time.Time `json:",omitempty"`
	ObservedTimestamp *time.Time `json:",omitempty"`
	EventName         string     `json:",omitempty"`
	Severity          log.Severity
	SeverityText      string
	Body              attribute.Value
	Attributes        []attribute.KeyValue
	TraceID           trace.TraceID
	SpanID            trace.SpanID
	TraceFlags        trace.TraceFlags
	Resource          *resource.Resource
	Scope             scopeJSON
	DroppedAttributes int
}

type scopeJSON struct {
	Name       string
	Version    string
	SchemaURL  string
	Attributes json.RawMessage
}

func (e *Exporter) newRecordJSON(r sdklog.Record) recordJSON {
	res := r.Resource()
	scope := r.InstrumentationScope()
	newRecord := recordJSON{
		EventName:    r.EventName(),
		Severity:     r.Severity(),
		SeverityText: r.SeverityText(),
		Body:         r.Body(),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]attribute.KeyValue, 0, r.AttributesLen()),

		Resource: res,
		Scope: scopeJSON{
			Name:      scope.Name,
			Version:   scope.Version,
			SchemaURL: scope.SchemaURL,
			// Set.String retains nested values, including non-finite floats,
			// without relying on the Set's pointer-receiver JSON marshaler.
			Attributes: json.RawMessage(scope.Attributes.String()),
		},

		DroppedAttributes: r.DroppedAttributes(),
	}

	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, kv)
		return true
	})

	if e.timestamps {
		timestamp := r.Timestamp()
		newRecord.Timestamp = &timestamp

		observedTimestamp := r.ObservedTimestamp()
		newRecord.ObservedTimestamp = &observedTimestamp
	}

	return newRecord
}
