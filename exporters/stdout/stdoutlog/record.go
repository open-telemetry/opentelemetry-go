// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"encoding/json"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

// recordJSON is a JSON-serializable representation of a Record.
type recordJSON struct {
	Timestamp         *time.Time `json:",omitempty"`
	ObservedTimestamp *time.Time `json:",omitempty"`
	EventName         string     `json:",omitempty"`
	Severity          log.Severity
	SeverityText      string
	Body              json.RawMessage
	Attributes        []json.RawMessage
	TraceID           trace.TraceID
	SpanID            trace.SpanID
	TraceFlags        trace.TraceFlags
	Resource          json.RawMessage
	Scope             scopeJSON
	DroppedAttributes int
}

func (e *Exporter) newRecordJSON(r sdklog.Record) recordJSON {
	res := r.Resource()
	scope := r.InstrumentationScope()
	newRecord := recordJSON{
		EventName:    r.EventName(),
		Severity:     r.Severity(),
		SeverityText: r.SeverityText(),
		Body:         newValueJSON(r.Body()),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]json.RawMessage, 0, r.AttributesLen()),

		Scope: scopeJSON{
			Scope:      scope,
			Attributes: json.RawMessage(scope.Attributes.String()),
		},

		DroppedAttributes: r.DroppedAttributes(),
	}

	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, json.RawMessage(kv.String()))
		return true
	})

	if res != nil {
		newRecord.Resource = json.RawMessage(res.Set().String())
	}

	if e.timestamps {
		timestamp := r.Timestamp()
		newRecord.Timestamp = &timestamp

		observedTimestamp := r.ObservedTimestamp()
		newRecord.ObservedTimestamp = &observedTimestamp
	}

	return newRecord
}

type scopeJSON struct {
	instrumentation.Scope
	Attributes json.RawMessage
}
