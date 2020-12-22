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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

const (
	errorTypeKey    = label.Key("error.type")
	errorMessageKey = label.Key("error.message")
	errorEventName  = "error"
)

var _ trace.Span = (*Span)(nil)

// Span is an OpenTelemetry Span used for testing.
type Span struct {
	lock          sync.RWMutex
	tracer        *Tracer
	spanContext   trace.SpanContext
	parentSpanID  trace.SpanID
	ended         bool
	name          string
	startTime     time.Time
	endTime       time.Time
	statusCode    codes.Code
	statusMessage string
	attributes    map[label.Key]label.Value
	events        []Event
	links         []trace.Link
	spanKind      trace.SpanKind
}

// Tracer returns the Tracer that created s.
func (s *Span) Tracer() trace.Tracer {
	return s.tracer
}

// End ends s. If the Tracer that created s was configured with a
// SpanRecorder, that recorder's OnEnd method is called as the final part of
// this method.
func (s *Span) End(opts ...trace.SpanOption) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	c := trace.NewSpanConfig(opts...)
	s.endTime = time.Now()
	if endTime := c.Timestamp; !endTime.IsZero() {
		s.endTime = endTime
	}

	s.ended = true
	if s.tracer.config.SpanRecorder != nil {
		s.tracer.config.SpanRecorder.OnEnd(s)
	}
}

// RecordError records an error as a Span event.
func (s *Span) RecordError(err error, opts ...trace.EventOption) {
	if err == nil || s.ended {
		return
	}

	errType := reflect.TypeOf(err)
	errTypeString := fmt.Sprintf("%s.%s", errType.PkgPath(), errType.Name())
	if errTypeString == "." {
		errTypeString = errType.String()
	}

	s.SetStatus(codes.Error, "")
	opts = append(opts, trace.WithAttributes(
		errorTypeKey.String(errTypeString),
		errorMessageKey.String(err.Error()),
	))

	s.AddEvent(errorEventName, opts...)
}

// AddEvent adds an event to s.
func (s *Span) AddEvent(name string, o ...trace.EventOption) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	c := trace.NewEventConfig(o...)

	var attributes map[label.Key]label.Value
	if l := len(c.Attributes); l > 0 {
		attributes = make(map[label.Key]label.Value, l)
		for _, attr := range c.Attributes {
			attributes[attr.Key] = attr.Value
		}
	}

	s.events = append(s.events, Event{
		Timestamp:  c.Timestamp,
		Name:       name,
		Attributes: attributes,
	})
}

// IsRecording returns the recording state of s.
func (s *Span) IsRecording() bool {
	return true
}

// SpanContext returns the SpanContext of s.
func (s *Span) SpanContext() trace.SpanContext {
	return s.spanContext
}

// SetStatus sets the status of s in the form of a code and a message.
func (s *Span) SetStatus(code codes.Code, msg string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.statusCode = code
	s.statusMessage = msg
}

// SetName sets the name of s.
func (s *Span) SetName(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.name = name
}

// SetAttributes sets attrs as attributes of s.
func (s *Span) SetAttributes(attrs ...label.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	for _, attr := range attrs {
		s.attributes[attr.Key] = attr.Value
	}
}

// Name returns the name most recently set on s, either at or after creation
// time. It cannot be change after End has been called on s.
func (s *Span) Name() string { return s.name }

// ParentSpanID returns the SpanID of the parent Span. If s is a root Span,
// and therefore does not have a parent, the returned SpanID will be invalid
// (i.e., it will contain all zeroes).
func (s *Span) ParentSpanID() trace.SpanID { return s.parentSpanID }

// Attributes returns the attributes set on s, either at or after creation
// time. If the same attribute key was set multiple times, the last call will
// be used. Attributes cannot be changed after End has been called on s.
func (s *Span) Attributes() map[label.Key]label.Value {
	s.lock.RLock()
	defer s.lock.RUnlock()

	attributes := make(map[label.Key]label.Value)

	for k, v := range s.attributes {
		attributes[k] = v
	}

	return attributes
}

// Events returns the events set on s. Events cannot be changed after End has
// been called on s.
func (s *Span) Events() []Event { return s.events }

// Links returns the links set on s at creation time. If multiple links for
// the same SpanContext were set, the last link will be used.
func (s *Span) Links() []trace.Link { return s.links }

// StartTime returns the time at which s was started. This will be the
// wall-clock time unless a specific start time was provided.
func (s *Span) StartTime() time.Time { return s.startTime }

// EndTime returns the time at which s was ended if at has been ended, or
// false otherwise. If the span has been ended, the returned time will be the
// wall-clock time unless a specific end time was provided.
func (s *Span) EndTime() (time.Time, bool) { return s.endTime, s.ended }

// Ended returns whether s has been ended, i.e. whether End has been called at
// least once on s.
func (s *Span) Ended() bool { return s.ended }

// StatusCode returns the code of the status most recently set on s, or
// codes.OK if no status has been explicitly set. It cannot be changed after
// End has been called on s.
func (s *Span) StatusCode() codes.Code { return s.statusCode }

// StatusMessage returns the status message most recently set on s or the
// empty string if no status message was set.
func (s *Span) StatusMessage() string { return s.statusMessage }

// SpanKind returns the span kind of s.
func (s *Span) SpanKind() trace.SpanKind { return s.spanKind }
