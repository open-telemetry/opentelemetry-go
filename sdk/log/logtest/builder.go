// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import (
	"context"
	"slices"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// RecordBuilder is used to facilitate unit testing implementations of
// [go.opentelemetry.io/otel/sdk/log.Exporter]
// and [go.opentelemetry.io/otel/sdk/log.Processor].
//
// Do not use RecordBuilder to create records in production code.
type RecordBuilder struct {
	timestamp         time.Time
	observedTimestamp time.Time
	severity          log.Severity
	severityText      string
	body              log.Value
	attrs             []log.KeyValue
	traceID           trace.TraceID
	spanID            trace.SpanID
	traceFlags        trace.TraceFlags

	resource *resource.Resource
	scope    instrumentation.Scope

	dropped int
}

// Record returns the accumulated log record.
func (b RecordBuilder) Record() sdklog.Record {
	var record sdklog.Record
	p := processor(func(r sdklog.Record) {
		r.SetTimestamp(b.timestamp)
		r.SetObservedTimestamp(b.observedTimestamp)
		r.SetSeverity(b.severity)
		r.SetSeverityText(b.severityText)
		r.SetBody(b.body)
		r.SetAttributes(b.attrs...)

		// Generate dropped attributes.
		for i := 0; i < b.dropped; i++ {
			r.AddAttributes(log.KeyValue{})
		}

		r.SetTraceID(b.traceID)
		r.SetSpanID(b.spanID)
		r.SetTraceFlags(b.traceFlags)

		record = r
	})

	attributeCountLimit := -1
	if b.dropped > 0 {
		// Make sure that we can generate dropped attributes.
		attributeCountLimit = len(b.attrs)
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(b.resource),
		sdklog.WithAttributeCountLimit(attributeCountLimit),
		sdklog.WithAttributeValueLengthLimit(-1),
		sdklog.WithProcessor(p),
	)

	l := provider.Logger(b.scope.Name,
		log.WithInstrumentationVersion(b.scope.Version),
		log.WithSchemaURL(b.scope.SchemaURL),
	)
	l.Emit(context.Background(), log.Record{}) // This executes the processor function.
	return record
}

// SetTimestamp sets the timestamp
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetTimestamp(t time.Time) RecordBuilder {
	b.timestamp = t
	return b
}

// SetObservedTimestamp sets the observed timestamp
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetObservedTimestamp(t time.Time) RecordBuilder {
	b.observedTimestamp = t
	return b
}

// SetSeverity sets the severity
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetSeverity(severity log.Severity) RecordBuilder {
	b.severity = severity
	return b
}

// SetSeverityText sets the severity text
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetSeverityText(text string) RecordBuilder {
	b.severityText = text
	return b
}

// SetBody sets the body
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetBody(v log.Value) RecordBuilder {
	b.body = v
	return b
}

// AddAttributes adds attributes
// to the record that is going to be returned by the builder.
func (b RecordBuilder) AddAttributes(attrs ...log.KeyValue) RecordBuilder {
	b.attrs = slices.Clone(b.attrs)
	b.attrs = append(b.attrs, attrs...)
	return b
}

// SetAttributes sets the attributes
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetAttributes(attrs ...log.KeyValue) RecordBuilder {
	b.attrs = slices.Clone(attrs)
	return b
}

// SetTraceID sets the trace ID
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetTraceID(traceID trace.TraceID) RecordBuilder {
	b.traceID = traceID
	return b
}

// SetSpanID sets the span ID
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetSpanID(spanID trace.SpanID) RecordBuilder {
	b.spanID = spanID
	return b
}

// SetTraceFlags sets the trace flags
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetTraceFlags(flags trace.TraceFlags) RecordBuilder {
	b.traceFlags = flags
	return b
}

// SetResource sets the resource
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetResource(r *resource.Resource) RecordBuilder {
	b.resource = r
	return b
}

// SetInstrumentationScope sets the instrumentation scope
// of the record that is going to be returned by the builder.
func (b RecordBuilder) SetInstrumentationScope(scope instrumentation.Scope) RecordBuilder {
	b.scope = scope
	return b
}

// SetDroppedAttributes sets the dropped attributes
// of the record that is going to be returned by the builder.
//
// Notice: The returned record is going to have an attribute count limit.
// Therefore, it will not be possible to add additional attributes on the record
// returned by the builder that has dropped attributes set to a value greater than 0
// (they will be dropped).
func (b RecordBuilder) SetDroppedAttributes(n int) RecordBuilder {
	b.dropped = n
	return b
}

type processor func(r sdklog.Record)

func (p processor) OnEmit(ctx context.Context, r sdklog.Record) error {
	p(r)
	return nil
}

func (processor) Enabled(context.Context, sdklog.Record) bool {
	return true
}

func (processor) Shutdown(ctx context.Context) error {
	return nil
}

func (processor) ForceFlush(context.Context) error {
	return nil
}
