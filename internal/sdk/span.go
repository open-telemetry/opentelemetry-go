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

package sdk

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
)

// Event is a Span Event.
type Event struct {
	// Name is the event name.
	Name string

	// Timestamp is the time of recording.
	Timestamp time.Time

	// Attributes are the identifying attributes.
	Attributes []kv.KeyValue
}

// Span is a mock span used in association with Tracer for testing purpose only.
type Span struct {
	sc           trace.SpanContext
	tracer       *Tracer
	Name         string
	Status       codes.Code
	StatusMsg    string
	attributeMap map[kv.Key]int
	Attributes   []kv.KeyValue
	Events       []Event
	Errors       []error
}

var _ trace.Span = (*Span)(nil)

// SpanContext returns associated kv.SpanContext. If the receiver is nil it returns
// an empty kv.SpanContext
func (s *Span) SpanContext() trace.SpanContext {
	if s == nil {
		return trace.EmptySpanContext()
	}
	return s.sc
}

// IsRecording always returns false for Span.
func (s *Span) IsRecording() bool {
	return false
}

// SetStatus sets the span Status and StatusMsg.
func (s *Span) SetStatus(status codes.Code, msg string) {
	s.Status = status
	s.StatusMsg = msg
}

// SetError does nothing.
func (s *Span) SetError(v bool) {
}

// SetAttributes sets attributes.
func (s *Span) SetAttributes(attrs ...kv.KeyValue) {
	for _, attr := range attrs {
		if attr.Value.Type() == kv.INVALID {
			continue
		}
		if idx, ok := s.attributeMap[attr.Key]; !ok {
			s.Attributes = append(s.Attributes, attr)
			s.attributeMap[attr.Key] = len(s.Attributes) - 1
		} else {
			s.Attributes[idx] = attr
		}
	}
}

// SetAttribute sets attribute k to value v.
func (s *Span) SetAttribute(k string, v interface{}) {
	attr := kv.Infer(k, v)
	if attr.Value.Type() != kv.INVALID {
		s.SetAttributes(attr)
	}
}

// End ends the span by calling the configured SpanRecorder.OnEnd.
func (s *Span) End(options ...trace.EndOption) {
	if s.tracer.Config.SpanRecorder != nil {
		s.tracer.Config.SpanRecorder.OnEnd(s)
	}
}

// RecordError records err
func (s *Span) RecordError(_ context.Context, err error, _ ...trace.ErrorOption) {
	s.Errors = append(s.Errors, err)
}

// SetName sets the span name.
func (s *Span) SetName(name string) {
	s.Name = name
}

// Tracer returns Tracer that created this Span.
func (s *Span) Tracer() trace.Tracer {
	return s.tracer
}

// AddEvent records an event.
func (s *Span) AddEvent(ctx context.Context, name string, attrs ...kv.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), name, attrs...)
}

// AddEvent records an event occuring at timestamp.
func (s *Span) AddEventWithTimestamp(_ context.Context, timestamp time.Time, name string, attrs ...kv.KeyValue) {
	s.Events = append(s.Events, Event{
		Name:       name,
		Timestamp:  timestamp,
		Attributes: attrs,
	})
}
