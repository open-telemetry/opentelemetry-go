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

	"google.golang.org/grpc/codes"

	otelcorrelation "go.opentelemetry.io/otel/api/correlation"
	otelcore "go.opentelemetry.io/otel/api/kv"
	oteltrace "go.opentelemetry.io/otel/api/trace"
	otelparent "go.opentelemetry.io/otel/internal/trace/parent"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

var (
	ComponentKey     = otelcore.Key("component")
	ServiceKey       = otelcore.Key("service")
	StatusCodeKey    = otelcore.Key("status.code")
	StatusMessageKey = otelcore.Key("status.message")
	ErrorKey         = otelcore.Key("error")
	NameKey          = otelcore.Key("name")
)

type MockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type MockTracer struct {
	Resources             otelcorrelation.Map
	FinishedSpans         []*MockSpan
	SpareTraceIDs         []oteltrace.ID
	SpareSpanIDs          []oteltrace.SpanID
	SpareContextKeyValues []MockContextKeyValue

	randLock sync.Mutex
	rand     *rand.Rand
}

var _ oteltrace.Tracer = &MockTracer{}
var _ migration.DeferredContextSetupTracerExtension = &MockTracer{}

func NewMockTracer() *MockTracer {
	return &MockTracer{
		Resources:             otelcorrelation.NewEmptyMap(),
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *MockTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error, opts ...oteltrace.StartOption) error {
	ctx, span := t.Start(ctx, name, opts...)
	defer span.End()
	return body(ctx)
}

func (t *MockTracer) Start(ctx context.Context, name string, opts ...oteltrace.StartOption) (context.Context, oteltrace.Span) {
	spanOpts := oteltrace.StartConfig{}
	for _, opt := range opts {
		opt(&spanOpts)
	}
	startTime := spanOpts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanContext := oteltrace.SpanContext{
		TraceID:    t.getTraceID(ctx, &spanOpts),
		SpanID:     t.getSpanID(),
		TraceFlags: 0,
	}
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		recording:      spanOpts.Record,
		Attributes: otelcorrelation.NewMap(otelcorrelation.MapUpdate{
			MultiKV: spanOpts.Attributes,
		}),
		StartTime:    startTime,
		EndTime:      time.Time{},
		ParentSpanID: t.getParentSpanID(ctx, &spanOpts),
		Events:       nil,
		SpanKind:     oteltrace.ValidateSpanKind(spanOpts.SpanKind),
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = oteltrace.ContextWithSpan(ctx, span)
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

func (t *MockTracer) getTraceID(ctx context.Context, spanOpts *oteltrace.StartConfig) oteltrace.ID {
	if parent := t.getParentSpanContext(ctx, spanOpts); parent.IsValid() {
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

func (t *MockTracer) getParentSpanID(ctx context.Context, spanOpts *oteltrace.StartConfig) oteltrace.SpanID {
	if parent := t.getParentSpanContext(ctx, spanOpts); parent.IsValid() {
		return parent.SpanID
	}
	return oteltrace.SpanID{}
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, spanOpts *oteltrace.StartConfig) oteltrace.SpanContext {
	spanCtx, _, _ := otelparent.GetSpanContextAndLinks(ctx, spanOpts.NewRoot)
	return spanCtx
}

func (t *MockTracer) getSpanID() oteltrace.SpanID {
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

func (t *MockTracer) getRandSpanID() oteltrace.SpanID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	sid := oteltrace.SpanID{}
	t.rand.Read(sid[:])

	return sid
}

func (t *MockTracer) getRandTraceID() oteltrace.ID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	tid := oteltrace.ID{}
	t.rand.Read(tid[:])

	return tid
}

func (t *MockTracer) DeferredContextSetupHook(ctx context.Context, span oteltrace.Span) context.Context {
	return t.addSpareContextValue(ctx)
}

type MockEvent struct {
	CtxAttributes otelcorrelation.Map
	Timestamp     time.Time
	Name          string
	Attributes    otelcorrelation.Map
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer oteltrace.Tracer
	spanContext    oteltrace.SpanContext
	SpanKind       oteltrace.SpanKind
	recording      bool

	Attributes   otelcorrelation.Map
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID oteltrace.SpanID
	Events       []MockEvent
}

var _ oteltrace.Span = &MockSpan{}
var _ migration.OverrideTracerSpanExtension = &MockSpan{}

func (s *MockSpan) SpanContext() oteltrace.SpanContext {
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

func (s *MockSpan) SetAttributes(attributes ...otelcore.KeyValue) {
	s.applyUpdate(otelcorrelation.MapUpdate{
		MultiKV: attributes,
	})
}

func (s *MockSpan) SetAttribute(k string, v interface{}) {
	s.SetAttributes(otelcore.Infer(k, v))
}

func (s *MockSpan) applyUpdate(update otelcorrelation.MapUpdate) {
	s.Attributes = s.Attributes.Apply(update)
}

func (s *MockSpan) End(options ...oteltrace.EndOption) {
	if !s.EndTime.IsZero() {
		return // already finished
	}
	endOpts := oteltrace.EndConfig{}

	for _, opt := range options {
		opt(&endOpts)
	}

	endTime := endOpts.EndTime
	if endTime.IsZero() {
		endTime = time.Now()
	}
	s.EndTime = endTime
	s.mockTracer.FinishedSpans = append(s.mockTracer.FinishedSpans, s)
}

func (s *MockSpan) RecordError(ctx context.Context, err error, opts ...oteltrace.ErrorOption) {
	if err == nil {
		return // no-op on nil error
	}

	if !s.EndTime.IsZero() {
		return // already finished
	}

	cfg := oteltrace.ErrorConfig{}

	for _, o := range opts {
		o(&cfg)
	}

	if cfg.Timestamp.IsZero() {
		cfg.Timestamp = time.Now()
	}

	if cfg.StatusCode != codes.OK {
		s.SetStatus(cfg.StatusCode, "")
	}

	s.AddEventWithTimestamp(ctx, cfg.Timestamp, "error",
		otelcore.String("error.type", reflect.TypeOf(err).String()),
		otelcore.String("error.message", err.Error()),
	)
}

func (s *MockSpan) Tracer() oteltrace.Tracer {
	return s.officialTracer
}

func (s *MockSpan) AddEvent(ctx context.Context, name string, attrs ...otelcore.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), name, attrs...)
}

func (s *MockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...otelcore.KeyValue) {
	s.Events = append(s.Events, MockEvent{
		CtxAttributes: otelcorrelation.MapFromContext(ctx),
		Timestamp:     timestamp,
		Name:          name,
		Attributes: otelcorrelation.NewMap(otelcorrelation.MapUpdate{
			MultiKV: attrs,
		}),
	})
}

func (s *MockSpan) OverrideTracer(tracer oteltrace.Tracer) {
	s.officialTracer = tracer
}
