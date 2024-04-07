package stdoutlog

import (
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"time"
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

func newRecordJSON(r sdklog.Record) recordJSON {
	newRecord := recordJSON{
		Timestamp:         r.Timestamp(),
		ObservedTimestamp: r.ObservedTimestamp(),
		Severity:          r.Severity(),
		SeverityText:      r.SeverityText(),
		Body:              r.Body(),

		TraceID:    r.TraceID(),
		SpanID:     r.SpanID(),
		TraceFlags: r.TraceFlags(),

		Attributes: make([]log.KeyValue, 0, r.AttributesLen()),

		Resource:                  r.Resource(),
		Scope:                     r.InstrumentationScope(),
		AttributeValueLengthLimit: r.AttributeValueLengthLimit(),
		AttributeCountLimit:       r.AttributeCountLimit(),
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		newRecord.Attributes = append(newRecord.Attributes, kv)
		return true
	})

	return newRecord
}
