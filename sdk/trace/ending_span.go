// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/trace"
)

// endingSpan wraps a recordingSpan that has stopped recording (endTime set)
// and is handed to span processor OnEnding callbacks. It is the only way to
// mutate a span once the span has stopped recording; its mutating methods skip
// the isRecording check so that processors can still modify the span before it
// becomes truly immutable after all OnEnding callbacks have returned.
type endingSpan struct {
	*recordingSpan
}

var _ ReadWriteSpan = endingSpan{}

// IsRecording returns true because the span is still mutable during OnEnding.
func (endingSpan) IsRecording() bool { return true }

// End is a no-op: the span is already in the process of ending.
func (endingSpan) End(...trace.SpanEndOption) {}

// SetStatus sets the status of the span.
func (s endingSpan) SetStatus(code codes.Code, description string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setStatus(code, description)
}

// SetAttributes sets attributes on the span.
func (s endingSpan) SetAttributes(attributes ...attribute.KeyValue) {
	if len(attributes) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setAttributes(attributes...)
}

// SetName sets the name of the span.
func (s endingSpan) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setName(name)
}

// AddEvent adds an event to the span.
func (s endingSpan) AddEvent(name string, o ...trace.EventOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addEvent(name, o...)
}

// AddLink adds a link to the span.
func (s endingSpan) AddLink(link trace.Link) {
	if !link.SpanContext.IsValid() && len(link.Attributes) == 0 &&
		link.SpanContext.TraceState().Len() == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addLink(link)
}

// RecordError records an error event on the span.
func (s endingSpan) RecordError(err error, opts ...trace.EventOption) {
	if err == nil {
		return
	}
	// err.Error() is caller-controlled and may call back into this span.
	// Build the event options before acquiring s.mu to avoid deadlock.
	o := errorEventOptions(err, opts)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.addEvent(semconv.ExceptionEventName, o...)
}
