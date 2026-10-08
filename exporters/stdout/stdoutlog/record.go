// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
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
	Body              valueJSON
	Attributes        []keyValueJSON
	TraceID           trace.TraceID
	SpanID            trace.SpanID
	TraceFlags        trace.TraceFlags
	Resource          any
	Scope             any
	DroppedAttributes int
}

func (e *Exporter) newRecordJSON(r sdklog.Record) recordJSON {
	res := r.Resource()
	newRecord := recordJSON{
		EventName:    r.EventName(),
		Severity:     r.Severity(),
		SeverityText: r.SeverityText(),
		Body:         newValueJSON(r.Body()),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]keyValueJSON, 0, r.AttributesLen()),

		Resource: res,
		Scope:    r.InstrumentationScope(),

		DroppedAttributes: r.DroppedAttributes(),
	}

	r.WalkAttributes(func(kv attribute.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, newKeyValueJSON(kv))
		return true
	})

	// Preserve the existing nil/empty resource and empty scope representations.
	if res.Len() > 0 {
		newRecord.Resource = newAttributeSetJSON(res.Iter())
	}
	scope := r.InstrumentationScope()
	if scope.Attributes.Len() > 0 {
		newRecord.Scope = scopeJSON{
			Scope:      scope,
			Attributes: newAttributeSetJSON(scope.Attributes.Iter()),
		}
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
	Attributes []keyValueJSON
}

func newAttributeSetJSON(iter attribute.Iterator) []keyValueJSON {
	attrs := make([]keyValueJSON, 0, iter.Len())
	for iter.Next() {
		attrs = append(attrs, newKeyValueJSON(iter.Attribute()))
	}
	return attrs
}
