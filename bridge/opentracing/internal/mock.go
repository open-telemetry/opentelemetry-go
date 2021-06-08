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

package internal

import (
	"context"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

//nolint:revive // ignoring missing comments for unexported global variables in an internal package.
var (
	ComponentKey     = attribute.Key("component")
	ServiceKey       = attribute.Key("service")
	StatusCodeKey    = attribute.Key("status.code")
	StatusMessageKey = attribute.Key("status.message")
	ErrorKey         = attribute.Key("error")
	NameKey          = attribute.Key("name")
)

type MockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type MockTracer struct {
	FinishedSpans         []*MockSpan
	SpareTraceIDs         []trace.TraceID
	SpareSpanIDs          []trace.SpanID
	SpareContextKeyValues []MockContextKeyValue

	randLock sync.Mutex
	rand     *rand.Rand
}

var _ trace.Tracer = &MockTracer{}
var _ migration.DeferredContextSetupTracerExtension = &MockTracer{}

func NewMockTracer() *MockTracer {
	return &MockTracer{
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *MockTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	config := trace.NewSpanStartConfig(opts...)
	startTime := config.Timestamp()
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    t.getTraceID(ctx, config),
		SpanID:     t.getSpanID(),
		TraceFlags: 0,
	})
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		Attributes:     config.Attributes(),
		StartTime:      startTime,
		EndTime:        time.Time{},
		ParentSpanID:   t.getParentSpanID(ctx, config),
		Events:         nil,
		SpanKind:       trace.ValidateSpanKind(config.SpanKind()),
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = trace.ContextWithSpan(ctx, span)
		ctx = t.addSpareContextValue(ctx)
	}
	return ctx, span
}

func (t *MockTracer) addSpareContextValue(ctx context.Context) context.Context {
	if len(t.SpareContextKeyValues) > 0 {
		pair := t.SpareContextKeyValues[0]
		t.SpareContextKeyValues[0] = MockContextKeyValue{}
		t.SpareContextKeyValues = t.SpareContextKeyValues[1:]
		if len(t.SpareContextKeyValues) == 0 {
			t.SpareContextKeyValues = nil
		}
		ctx = context.WithValue(ctx, pair.Key, pair.Value)
	}
	return ctx
}

func (t *MockTracer) getTraceID(ctx context.Context, config *trace.SpanConfig) trace.TraceID {
	if parent := t.getParentSpanContext(ctx, config); parent.IsValid() {
		return parent.TraceID()
	}
	if len(t.SpareTraceIDs) > 0 {
		traceID := t.SpareTraceIDs[0]
		t.SpareTraceIDs = t.SpareTraceIDs[1:]
		if len(t.SpareTraceIDs) == 0 {
			t.SpareTraceIDs = nil
		}
		return traceID
	}
	return t.getRandTraceID()
}

func (t *MockTracer) getParentSpanID(ctx context.Context, config *trace.SpanConfig) trace.SpanID {
	if parent := t.getParentSpanContext(ctx, config); parent.IsValid() {
		return parent.SpanID()
	}
	return trace.SpanID{}
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, config *trace.SpanConfig) trace.SpanContext {
	if !config.NewRoot() {
		return trace.SpanContextFromContext(ctx)
	}
	return trace.SpanContext{}
}

func (t *MockTracer) getSpanID() trace.SpanID {
	if len(t.SpareSpanIDs) > 0 {
		spanID := t.SpareSpanIDs[0]
		t.SpareSpanIDs = t.SpareSpanIDs[1:]
		if len(t.SpareSpanIDs) == 0 {
			t.SpareSpanIDs = nil
		}
		return spanID
	}
	return t.getRandSpanID()
}

func (t *MockTracer) getRandSpanID() trace.SpanID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	sid := trace.SpanID{}
	t.rand.Read(sid[:])

	return sid
}

func (t *MockTracer) getRandTraceID() trace.TraceID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	tid := trace.TraceID{}
	t.rand.Read(tid[:])

	return tid
}

func (t *MockTracer) DeferredContextSetupHook(ctx context.Context, span trace.Span) context.Context {
	return t.addSpareContextValue(ctx)
}

type MockEvent struct {
	Timestamp  time.Time
	Name       string
	Attributes []attribute.KeyValue
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer trace.Tracer
	spanContext    trace.SpanContext
	SpanKind       trace.SpanKind
	recording      bool

	Attributes   []attribute.KeyValue
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID trace.SpanID
	Events       []MockEvent
}

var _ trace.Span = &MockSpan{}
var _ migration.OverrideTracerSpanExtension = &MockSpan{}

func (s *MockSpan) SpanContext() trace.SpanContext {
	return s.spanContext
}

func (s *MockSpan) IsRecording() bool {
	return s.recording
}

func (s *MockSpan) SetStatus(code codes.Code, msg string) {
	s.SetAttributes(StatusCodeKey.Int(int(code)), StatusMessageKey.String(msg))
}

func (s *MockSpan) SetName(name string) {
	s.SetAttributes(NameKey.String(name))
}

func (s *MockSpan) SetError(v bool) {
	s.SetAttributes(ErrorKey.Bool(v))
}

func (s *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	s.applyUpdate(attributes)
}

func (s *MockSpan) applyUpdate(update []attribute.KeyValue) {
	updateM := make(map[attribute.Key]attribute.Value, len(update))
	for _, kv := range update {
		updateM[kv.Key] = kv.Value
	}

	seen := make(map[attribute.Key]struct{})
	for i, kv := range s.Attributes {
		if v, ok := updateM[kv.Key]; ok {
			s.Attributes[i].Value = v
			seen[kv.Key] = struct{}{}
		}
	}

	for k, v := range updateM {
		if _, ok := seen[k]; ok {
			continue
		}
		s.Attributes = append(s.Attributes, attribute.KeyValue{Key: k, Value: v})
	}
}

func (s *MockSpan) End(options ...trace.SpanEndOption) {
	if !s.EndTime.IsZero() {
		return // already finished
	}
	config := trace.NewSpanEndConfig(options...)
	endTime := config.Timestamp()
	if endTime.IsZero() {
		endTime = time.Now()
	}
	s.EndTime = endTime
	s.mockTracer.FinishedSpans = append(s.mockTracer.FinishedSpans, s)
}

func (s *MockSpan) RecordError(err error, opts ...trace.EventOption) {
	if err == nil {
		return // no-op on nil error
	}

	if !s.EndTime.IsZero() {
		return // already finished
	}

	s.SetStatus(codes.Error, "")
	opts = append(opts, trace.WithAttributes(
		semconv.ExceptionTypeKey.String(reflect.TypeOf(err).String()),
		semconv.ExceptionMessageKey.String(err.Error()),
	))
	s.AddEvent(semconv.ExceptionEventName, opts...)
}

func (s *MockSpan) Tracer() trace.Tracer {
	return s.officialTracer
}

func (s *MockSpan) AddEvent(name string, o ...trace.EventOption) {
	c := trace.NewEventConfig(o...)
	s.Events = append(s.Events, MockEvent{
		Timestamp:  c.Timestamp(),
		Name:       name,
		Attributes: c.Attributes(),
	})
}

func (s *MockSpan) OverrideTracer(tracer trace.Tracer) {
	s.officialTracer = tracer
}
