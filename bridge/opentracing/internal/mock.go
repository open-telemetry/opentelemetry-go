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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/internal/baggage"
	otelparent "go.opentelemetry.io/otel/internal/trace/parent"
	"go.opentelemetry.io/otel/label"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

var (
	ComponentKey     = label.Key("component")
	ServiceKey       = label.Key("service")
	StatusCodeKey    = label.Key("status.code")
	StatusMessageKey = label.Key("status.message")
	ErrorKey         = label.Key("error")
	NameKey          = label.Key("name")
)

type MockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type MockTracer struct {
	Resources             baggage.Map
	FinishedSpans         []*MockSpan
	SpareTraceIDs         []otel.TraceID
	SpareSpanIDs          []otel.SpanID
	SpareContextKeyValues []MockContextKeyValue

	randLock sync.Mutex
	rand     *rand.Rand
}

var _ otel.Tracer = &MockTracer{}
var _ migration.DeferredContextSetupTracerExtension = &MockTracer{}

func NewMockTracer() *MockTracer {
	return &MockTracer{
		Resources:             baggage.NewEmptyMap(),
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *MockTracer) Start(ctx context.Context, name string, opts ...otel.SpanOption) (context.Context, otel.Span) {
	config := otel.NewSpanConfig(opts...)
	startTime := config.Timestamp
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanContext := otel.SpanContext{
		TraceID:    t.getTraceID(ctx, config),
		SpanID:     t.getSpanID(),
		TraceFlags: 0,
	}
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		recording:      config.Record,
		Attributes: baggage.NewMap(baggage.MapUpdate{
			MultiKV: config.Attributes,
		}),
		StartTime:    startTime,
		EndTime:      time.Time{},
		ParentSpanID: t.getParentSpanID(ctx, config),
		Events:       nil,
		SpanKind:     otel.ValidateSpanKind(config.SpanKind),
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = otel.ContextWithSpan(ctx, span)
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

func (t *MockTracer) getTraceID(ctx context.Context, config *otel.SpanConfig) otel.TraceID {
	if parent := t.getParentSpanContext(ctx, config); parent.IsValid() {
		return parent.TraceID
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

func (t *MockTracer) getParentSpanID(ctx context.Context, config *otel.SpanConfig) otel.SpanID {
	if parent := t.getParentSpanContext(ctx, config); parent.IsValid() {
		return parent.SpanID
	}
	return otel.SpanID{}
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, config *otel.SpanConfig) otel.SpanContext {
	spanCtx, _, _ := otelparent.GetSpanContextAndLinks(ctx, config.NewRoot)
	return spanCtx
}

func (t *MockTracer) getSpanID() otel.SpanID {
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

func (t *MockTracer) getRandSpanID() otel.SpanID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	sid := otel.SpanID{}
	t.rand.Read(sid[:])

	return sid
}

func (t *MockTracer) getRandTraceID() otel.TraceID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	tid := otel.TraceID{}
	t.rand.Read(tid[:])

	return tid
}

func (t *MockTracer) DeferredContextSetupHook(ctx context.Context, span otel.Span) context.Context {
	return t.addSpareContextValue(ctx)
}

type MockEvent struct {
	CtxAttributes baggage.Map
	Timestamp     time.Time
	Name          string
	Attributes    baggage.Map
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer otel.Tracer
	spanContext    otel.SpanContext
	SpanKind       otel.SpanKind
	recording      bool

	Attributes   baggage.Map
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID otel.SpanID
	Events       []MockEvent
}

var _ otel.Span = &MockSpan{}
var _ migration.OverrideTracerSpanExtension = &MockSpan{}

func (s *MockSpan) SpanContext() otel.SpanContext {
	return s.spanContext
}

func (s *MockSpan) IsRecording() bool {
	return s.recording
}

func (s *MockSpan) SetStatus(code codes.Code, msg string) {
	s.SetAttributes(StatusCodeKey.Uint32(uint32(code)), StatusMessageKey.String(msg))
}

func (s *MockSpan) SetName(name string) {
	s.SetAttributes(NameKey.String(name))
}

func (s *MockSpan) SetError(v bool) {
	s.SetAttributes(ErrorKey.Bool(v))
}

func (s *MockSpan) SetAttributes(attributes ...label.KeyValue) {
	s.applyUpdate(baggage.MapUpdate{
		MultiKV: attributes,
	})
}

func (s *MockSpan) applyUpdate(update baggage.MapUpdate) {
	s.Attributes = s.Attributes.Apply(update)
}

func (s *MockSpan) End(options ...otel.SpanOption) {
	if !s.EndTime.IsZero() {
		return // already finished
	}
	config := otel.NewSpanConfig(options...)
	endTime := config.Timestamp
	if endTime.IsZero() {
		endTime = time.Now()
	}
	s.EndTime = endTime
	s.mockTracer.FinishedSpans = append(s.mockTracer.FinishedSpans, s)
}

func (s *MockSpan) RecordError(ctx context.Context, err error, opts ...otel.ErrorOption) {
	if err == nil {
		return // no-op on nil error
	}

	if !s.EndTime.IsZero() {
		return // already finished
	}

	cfg := otel.NewErrorConfig(opts...)

	if cfg.Timestamp.IsZero() {
		cfg.Timestamp = time.Now()
	}

	if cfg.StatusCode != codes.Ok {
		s.SetStatus(cfg.StatusCode, "")
	}

	s.AddEventWithTimestamp(ctx, cfg.Timestamp, "error",
		label.String("error.type", reflect.TypeOf(err).String()),
		label.String("error.message", err.Error()),
	)
}

func (s *MockSpan) Tracer() otel.Tracer {
	return s.officialTracer
}

func (s *MockSpan) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), name, attrs...)
}

func (s *MockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
	s.Events = append(s.Events, MockEvent{
		CtxAttributes: baggage.MapFromContext(ctx),
		Timestamp:     timestamp,
		Name:          name,
		Attributes: baggage.NewMap(baggage.MapUpdate{
			MultiKV: attrs,
		}),
	})
}

func (s *MockSpan) OverrideTracer(tracer otel.Tracer) {
	s.officialTracer = tracer
}
