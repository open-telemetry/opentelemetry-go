// Copyright 2019, OpenTelemetry Authors
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
	"sync"
	"time"

	"google.golang.org/grpc/codes"

	otelcore "go.opentelemetry.io/otel/api/core"
	oteldctx "go.opentelemetry.io/otel/api/distributedcontext"
	otelkey "go.opentelemetry.io/otel/api/key"
	oteltrace "go.opentelemetry.io/otel/api/trace"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

var (
	ComponentKey = otelkey.New("component")
	ServiceKey   = otelkey.New("service")
	StatusKey    = otelkey.New("status")
	ErrorKey     = otelkey.New("error")
	NameKey      = otelkey.New("name")
)

type MockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type MockTracer struct {
	Resources             oteldctx.Map
	FinishedSpans         []*MockSpan
	SpareTraceIDs         []otelcore.TraceID
	SpareSpanIDs          []otelcore.SpanID
	SpareContextKeyValues []MockContextKeyValue

	randLock sync.Mutex
	rand     *rand.Rand
}

var _ oteltrace.Tracer = &MockTracer{}
var _ migration.DeferredContextSetupTracerExtension = &MockTracer{}

func NewMockTracer() *MockTracer {
	return &MockTracer{
		Resources:             oteldctx.NewEmptyMap(),
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *MockTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	ctx, span := t.Start(ctx, name)
	defer span.End()
	return body(ctx)
}

func (t *MockTracer) Start(ctx context.Context, name string, opts ...oteltrace.SpanOption) (context.Context, oteltrace.Span) {
	spanOpts := oteltrace.SpanOptions{}
	for _, opt := range opts {
		opt(&spanOpts)
	}
	startTime := spanOpts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanKind := spanOpts.SpanKind
	if spanKind == "" {
		spanKind = oteltrace.SpanKindInternal
	}
	spanContext := otelcore.SpanContext{
		TraceID:    t.getTraceID(ctx, &spanOpts),
		SpanID:     t.getSpanID(),
		TraceFlags: 0,
	}
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		recording:      spanOpts.Record,
		Attributes: oteldctx.NewMap(oteldctx.MapUpdate{
			MultiKV: spanOpts.Attributes,
		}),
		StartTime:    startTime,
		EndTime:      time.Time{},
		ParentSpanID: t.getParentSpanID(ctx, &spanOpts),
		Events:       nil,
		SpanKind:     spanKind,
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = oteltrace.SetCurrentSpan(ctx, span)
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

func (t *MockTracer) getTraceID(ctx context.Context, spanOpts *oteltrace.SpanOptions) otelcore.TraceID {
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

func (t *MockTracer) getParentSpanID(ctx context.Context, spanOpts *oteltrace.SpanOptions) otelcore.SpanID {
	if parent := t.getParentSpanContext(ctx, spanOpts); parent.IsValid() {
		return parent.SpanID
	}
	return otelcore.SpanID{}
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, spanOpts *oteltrace.SpanOptions) otelcore.SpanContext {
	if spanOpts.Relation.RelationshipType == oteltrace.ChildOfRelationship &&
		spanOpts.Relation.SpanContext.IsValid() {
		return spanOpts.Relation.SpanContext
	}
	if parentSpanContext := oteltrace.CurrentSpan(ctx).SpanContext(); parentSpanContext.IsValid() {
		return parentSpanContext
	}
	return otelcore.EmptySpanContext()
}

func (t *MockTracer) getSpanID() otelcore.SpanID {
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

func (t *MockTracer) getRandSpanID() otelcore.SpanID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	sid := otelcore.SpanID{}
	t.rand.Read(sid[:])

	return sid
}

func (t *MockTracer) getRandTraceID() otelcore.TraceID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	tid := otelcore.TraceID{}
	t.rand.Read(tid[:])

	return tid
}

func (t *MockTracer) DeferredContextSetupHook(ctx context.Context, span oteltrace.Span) context.Context {
	return t.addSpareContextValue(ctx)
}

type MockEvent struct {
	CtxAttributes oteldctx.Map
	Timestamp     time.Time
	Msg           string
	Attributes    oteldctx.Map
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer oteltrace.Tracer
	spanContext    otelcore.SpanContext
	SpanKind       oteltrace.SpanKind
	recording      bool

	Attributes   oteldctx.Map
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID otelcore.SpanID
	Events       []MockEvent
}

var _ oteltrace.Span = &MockSpan{}
var _ migration.OverrideTracerSpanExtension = &MockSpan{}

func (s *MockSpan) SpanContext() otelcore.SpanContext {
	return s.spanContext
}

func (s *MockSpan) IsRecording() bool {
	return s.recording
}

func (s *MockSpan) SetStatus(status codes.Code) {
	s.SetAttribute(NameKey.Uint32(uint32(status)))
}

func (s *MockSpan) SetName(name string) {
	s.SetAttribute(NameKey.String(name))
}

func (s *MockSpan) SetError(v bool) {
	s.SetAttribute(ErrorKey.Bool(v))
}

func (s *MockSpan) SetAttribute(attribute otelcore.KeyValue) {
	s.applyUpdate(oteldctx.MapUpdate{
		SingleKV: attribute,
	})
}

func (s *MockSpan) SetAttributes(attributes ...otelcore.KeyValue) {
	s.applyUpdate(oteldctx.MapUpdate{
		MultiKV: attributes,
	})
}

func (s *MockSpan) applyUpdate(update oteldctx.MapUpdate) {
	s.Attributes = s.Attributes.Apply(update)
}

func (s *MockSpan) End(options ...oteltrace.EndOption) {
	if !s.EndTime.IsZero() {
		return // already finished
	}
	endOpts := oteltrace.EndOptions{}

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

func (s *MockSpan) Tracer() oteltrace.Tracer {
	return s.officialTracer
}

func (s *MockSpan) AddEvent(ctx context.Context, msg string, attrs ...otelcore.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), msg, attrs...)
}

func (s *MockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...otelcore.KeyValue) {
	s.Events = append(s.Events, MockEvent{
		CtxAttributes: oteldctx.FromContext(ctx),
		Timestamp:     timestamp,
		Msg:           msg,
		Attributes: oteldctx.NewMap(oteldctx.MapUpdate{
			MultiKV: attrs,
		}),
	})
}

func (s *MockSpan) AddLink(link oteltrace.Link) {
	// TODO
}

func (s *MockSpan) Link(sc otelcore.SpanContext, attrs ...otelcore.KeyValue) {
	// TODO
}

func (s *MockSpan) OverrideTracer(tracer oteltrace.Tracer) {
	s.officialTracer = tracer
}
