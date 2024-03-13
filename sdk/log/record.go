// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// Record is a log record emitted by the Logger.
type Record struct{}

// Timestamp returns the time when the log record occurred.
func (r *Record) Timestamp() time.Time {
	// TODO (#5064): Implement.
	return time.Time{}
}

// SetTimestamp sets the time when the log record occurred.
func (r *Record) SetTimestamp(t time.Time) {
	// TODO (#5064): Implement.
}

// ObservedTimestamp returns the time when the log record was observed.
func (r *Record) ObservedTimestamp() time.Time {
	// TODO (#5064): Implement.
	return time.Time{}
}

// SetObservedTimestamp sets the time when the log record was observed.
func (r *Record) SetObservedTimestamp(t time.Time) {
	// TODO (#5064): Implement.
}

// Severity returns the severity of the log record.
func (r *Record) Severity() log.Severity {
	// TODO (#5064): Implement.
	return log.Severity(0)
}

// SetSeverity sets the severity level of the log record.
func (r *Record) SetSeverity(level log.Severity) {
	// TODO (#5064): Implement.
}

// SeverityText returns severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Record) SeverityText() string {
	// TODO (#5064): Implement.
	return ""
}

// SetSeverityText sets severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Record) SetSeverityText(text string) {
	// TODO (#5064): Implement.
}

// Body returns the body of the log record.
func (r *Record) Body() log.Value {
	// TODO (#5064): Implement.
	return log.Value{}
}

// SetBody sets the body of the log record.
func (r *Record) SetBody(v log.Value) {
	// TODO (#5064): Implement.
}

// WalkAttributes walks all attributes the log record holds by calling f for
// each on each [log.KeyValue] in the [Record]. Iteration stops if f returns false.
func (r *Record) WalkAttributes(f func(log.KeyValue) bool) {
	// TODO (#5064): Implement.
}

// AddAttributes adds attributes to the log record.
func (r *Record) AddAttributes(attrs ...log.KeyValue) {
	// TODO (#5064): Implement.
}

// SetAttributes sets (and overrides) attributes to the log record.
func (r *Record) SetAttributes(attrs ...log.KeyValue) {
	// TODO (#5064): Implement.
}

// AttributesLen returns the number of attributes in the log record.
func (r *Record) AttributesLen() int {
	// TODO (#5064): Implement.
	return 0
}

// TraceID returns the trace ID or empty array.
func (r *Record) TraceID() trace.TraceID {
	// TODO (#5064): Implement.
	return trace.TraceID{}
}

// SetTraceID sets the trace ID.
func (r *Record) SetTraceID(id trace.TraceID) {
	// TODO (#5064): Implement.
}

// SpanID returns the span ID or empty array.
func (r *Record) SpanID() trace.SpanID {
	// TODO (#5064): Implement.
	return trace.SpanID{}
}

// SetSpanID sets the span ID.
func (r *Record) SetSpanID(id trace.SpanID) {
	// TODO (#5064): Implement.
}

// TraceFlags returns the trace flags.
func (r *Record) TraceFlags() trace.TraceFlags {
	return 0
}

// SetTraceFlags sets the trace flags.
func (r *Record) SetTraceFlags(flags trace.TraceFlags) {
	// TODO (#5064): Implement.
}

// Resource returns the entity that collected the log.
func (r *Record) Resource() resource.Resource {
	// TODO (#5064): Implement.
	return resource.Resource{}
}

// InstrumentationScope returns the scope that the Logger was created with.
func (r *Record) InstrumentationScope() instrumentation.Scope {
	// TODO (#5064): Implement.
	return instrumentation.Scope{}
}

// AttributeValueLengthLimit is the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value should be truncated to this length.
//
// Negative value means no limit should be applied.
func (r *Record) AttributeValueLengthLimit() int {
	// TODO (#5064): Implement.
	return 0
}

// AttributeCountLimit is the maximum allowed log record attribute count. Any
// attribute added to a log record once this limit is reached should be dropped.
//
// Zero means no attributes should be recorded.
//
// Negative value means no limit should be applied.
func (r *Record) AttributeCountLimit() int {
	// TODO (#5064): Implement.
	return 0
}

// Clone returns a copy of the record with no shared state. The original record
// and the clone can both be modified without interfering with each other.
func (r *Record) Clone() Record {
	// TODO (#5064): Implement.
	return *r
}
