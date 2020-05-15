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

package testtrace

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/kv/value"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

const (
	errorTypeKey    = kv.Key("error.type")
	errorMessageKey = kv.Key("error.message")
	errorEventName  = "error"
)

var _ trace.Span = (*Span)(nil)

type Span struct {
	lock          *sync.RWMutex
	tracer        *Tracer
	spanContext   trace.SpanContext
	parentSpanID  trace.SpanID
	ended         bool
	name          string
	startTime     time.Time
	endTime       time.Time
	statusCode    codes.Code
	statusMessage string
	attributes    map[kv.Key]value.Value
	events        []Event
	links         map[trace.SpanContext][]kv.KeyValue
}

func (s *Span) Tracer() trace.Tracer {
	return s.tracer
}

func (s *Span) End(opts ...trace.EndOption) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	var c trace.EndConfig

	for _, opt := range opts {
		opt(&c)
	}

	s.endTime = time.Now()

	if endTime := c.EndTime; !endTime.IsZero() {
		s.endTime = endTime
	}

	s.ended = true
}

func (s *Span) RecordError(ctx context.Context, err error, opts ...trace.ErrorOption) {
	if err == nil || s.ended {
		return
	}

	cfg := trace.ErrorConfig{}
	for _, o := range opts {
		o(&cfg)
	}

	if cfg.Timestamp.IsZero() {
		cfg.Timestamp = time.Now()
	}

	if cfg.StatusCode != codes.OK {
		s.SetStatus(cfg.StatusCode, "")
	}

	errType := reflect.TypeOf(err)
	errTypeString := fmt.Sprintf("%s.%s", errType.PkgPath(), errType.Name())
	if errTypeString == "." {
		errTypeString = errType.String()
	}

	s.AddEventWithTimestamp(ctx, cfg.Timestamp, errorEventName,
		errorTypeKey.String(errTypeString),
		errorMessageKey.String(err.Error()),
	)
}

func (s *Span) AddEvent(ctx context.Context, name string, attrs ...kv.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), name, attrs...)
}

func (s *Span) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...kv.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	attributes := make(map[kv.Key]value.Value)

	for _, attr := range attrs {
		attributes[attr.Key] = attr.Value
	}

	s.events = append(s.events, Event{
		Timestamp:  timestamp,
		Name:       name,
		Attributes: attributes,
	})
}

func (s *Span) IsRecording() bool {
	return true
}

func (s *Span) SpanContext() trace.SpanContext {
	return s.spanContext
}

func (s *Span) SetStatus(code codes.Code, msg string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.statusCode = code
	s.statusMessage = msg
}

func (s *Span) SetName(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	s.name = name
}

func (s *Span) SetAttributes(attrs ...kv.KeyValue) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.ended {
		return
	}

	for _, attr := range attrs {
		s.attributes[attr.Key] = attr.Value
	}
}

func (s *Span) SetAttribute(k string, v interface{}) {
	s.SetAttributes(kv.Infer(k, v))
}

// Name returns the name most recently set on the Span, either at or after creation time.
// It cannot be change after End has been called on the Span.
func (s *Span) Name() string {
	return s.name
}

// ParentSpanID returns the SpanID of the parent Span.
// If the Span is a root Span and therefore does not have a parent, the returned SpanID will be invalid
// (i.e., it will contain all zeroes).
func (s *Span) ParentSpanID() trace.SpanID {
	return s.parentSpanID
}

// Attributes returns the attributes set on the Span, either at or after creation time.
// If the same attribute key was set multiple times, the last call will be used.
// Attributes cannot be changed after End has been called on the Span.
func (s *Span) Attributes() map[kv.Key]value.Value {
	s.lock.RLock()
	defer s.lock.RUnlock()

	attributes := make(map[kv.Key]value.Value)

	for k, v := range s.attributes {
		attributes[k] = v
	}

	return attributes
}

// Events returns the events set on the Span.
// Events cannot be changed after End has been called on the Span.
func (s *Span) Events() []Event {
	return s.events
}

// Links returns the links set on the Span at creation time.
// If multiple links for the same SpanContext were set, the last link will be used.
func (s *Span) Links() map[trace.SpanContext][]kv.KeyValue {
	links := make(map[trace.SpanContext][]kv.KeyValue)

	for sc, attributes := range s.links {
		links[sc] = append([]kv.KeyValue{}, attributes...)
	}

	return links
}

// StartTime returns the time at which the Span was started.
// This will be the wall-clock time unless a specific start time was provided.
func (s *Span) StartTime() time.Time {
	return s.startTime
}

// EndTime returns the time at which the Span was ended if at has been ended,
// or false otherwise.
// If the span has been ended, the returned time will be the wall-clock time
// unless a specific end time was provided.
func (s *Span) EndTime() (time.Time, bool) {
	return s.endTime, s.ended
}

// Ended returns whether the Span has been ended,
// i.e., whether End has been called at least once on the Span.
func (s *Span) Ended() bool {
	return s.ended
}

// Status returns the status most recently set on the Span,
// or codes.OK if no status has been explicitly set.
// It cannot be changed after End has been called on the Span.
func (s *Span) StatusCode() codes.Code {
	return s.statusCode
}

// StatusMessage returns the status message most recently set on the
// Span or the empty string if no status mesaage was set.
func (s *Span) StatusMessage() string {
	return s.statusMessage
}
