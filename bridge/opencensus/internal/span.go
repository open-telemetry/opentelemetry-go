// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal // import "go.opentelemetry.io/otel/bridge/opencensus/internal"

import (
	"fmt"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	// MessageSendEvent is the name of the message send event.
	MessageSendEvent = "message send"
	// MessageReceiveEvent is the name of the message receive event.
	MessageReceiveEvent = "message receive"
)

var (
	// UncompressedKey is used for the uncompressed byte size attribute.
	UncompressedKey = attribute.Key("uncompressed byte size")
	// CompressedKey is used for the compressed byte size attribute.
	CompressedKey = attribute.Key("compressed byte size")
)

// Span is an OpenCensus SpanInterface wrapper for an OpenTelemetry Span.
type Span struct {
	otelSpan trace.Span
}

// NewSpan returns an OpenCensus Span wrapping an OpenTelemetry Span.
func NewSpan(s trace.Span) *octrace.Span {
	return octrace.NewSpan(&Span{otelSpan: s})
}

// IsRecordingEvents returns true if events are being recorded for this span.
func (s *Span) IsRecordingEvents() bool {
	return s.otelSpan.IsRecording()
}

// End ends this span.
func (s *Span) End() {
	s.otelSpan.End()
}

// SpanContext returns the SpanContext of this span.
func (s *Span) SpanContext() octrace.SpanContext {
	return otel2oc.SpanContext(s.otelSpan.SpanContext())
}

// SetName sets the name of this span, if it is recording events.
func (s *Span) SetName(name string) {
	s.otelSpan.SetName(name)
}

// SetStatus sets the status of this span, if it is recording events.
func (s *Span) SetStatus(status octrace.Status) {
	s.otelSpan.SetStatus(codes.Code(status.Code), status.Message)
}

// AddAttributes sets attributes in this span.
func (s *Span) AddAttributes(attributes ...octrace.Attribute) {
	s.otelSpan.SetAttributes(oc2otel.Attributes(attributes)...)
}

// Annotate adds an annotation with attributes to this span.
func (s *Span) Annotate(attributes []octrace.Attribute, str string) {
	s.otelSpan.AddEvent(str, trace.WithAttributes(oc2otel.Attributes(attributes)...))
}

// Annotatef adds a formatted annotation with attributes to this span.
func (s *Span) Annotatef(attributes []octrace.Attribute, format string, a ...interface{}) {
	s.Annotate(attributes, fmt.Sprintf(format, a...))
}

// AddMessageSendEvent adds a message send event to this span.
func (s *Span) AddMessageSendEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otelSpan.AddEvent(MessageSendEvent,
		trace.WithAttributes(
			attribute.KeyValue{
				Key:   UncompressedKey,
				Value: attribute.Int64Value(uncompressedByteSize),
			},
			attribute.KeyValue{
				Key:   CompressedKey,
				Value: attribute.Int64Value(compressedByteSize),
			}),
	)
}

// AddMessageReceiveEvent adds a message receive event to this span.
func (s *Span) AddMessageReceiveEvent(messageID, uncompressedByteSize, compressedByteSize int64) {
	s.otelSpan.AddEvent(MessageReceiveEvent,
		trace.WithAttributes(
			attribute.KeyValue{
				Key:   UncompressedKey,
				Value: attribute.Int64Value(uncompressedByteSize),
			},
			attribute.KeyValue{
				Key:   CompressedKey,
				Value: attribute.Int64Value(compressedByteSize),
			}),
	)
}

// AddLink adds a link to this span.
func (s *Span) AddLink(l octrace.Link) {
	Handle(fmt.Errorf("ignoring OpenCensus link %+v for span %q because OpenTelemetry doesn't support setting links after creation", l, s.String()))
}

// String prints a string representation of this span.
func (s *Span) String() string {
	return fmt.Sprintf("span %s", s.otelSpan.SpanContext().SpanID().String())
}
