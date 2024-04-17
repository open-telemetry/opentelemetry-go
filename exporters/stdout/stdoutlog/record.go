// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

import (
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// recordJSON is a JSON-serializable representation of a Record.
type recordJSON struct {
	Timestamp                 time.Time
	ObservedTimestamp         time.Time
	Severity                  log.Severity
	SeverityText              string
	Body                      log.Value
	Attributes                []log.KeyValue
	TraceID                   trace.TraceID
	SpanID                    trace.SpanID
	TraceFlags                trace.TraceFlags
	Resource                  resource.Resource
	Scope                     instrumentation.Scope
	AttributeValueLengthLimit int
	AttributeCountLimit       int
}

func (e *Exporter) newRecordJSON(r sdklog.Record) recordJSON {
	newRecord := recordJSON{
		Severity:     r.Severity(),
		SeverityText: r.SeverityText(),
		Body:         r.Body(),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]log.KeyValue, 0, r.AttributesLen()),

		Resource: r.Resource(),
		Scope:    r.InstrumentationScope(),
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, kv)
		return true
	})

	if e.timestamps {
		newRecord.Timestamp = r.Timestamp()
		newRecord.ObservedTimestamp = r.ObservedTimestamp()
	}

	return newRecord
}
