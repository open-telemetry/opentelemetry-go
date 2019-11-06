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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opentracing/migration"
)

var (
	ComponentKey = otel.NewKey("component")
	ServiceKey   = otel.NewKey("service")
	StatusKey    = otel.NewKey("status")
	ErrorKey     = otel.NewKey("error")
	NameKey      = otel.NewKey("name")
)

type MockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type MockTracer struct {
	Resources             otel.Map
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
		Resources:             otel.NewEmptyMap(),
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

func (t *MockTracer) Start(ctx context.Context, name string, opts ...otel.SpanOption) (context.Context, otel.Span) {
	spanOpts := otel.SpanOptions{}
	for _, opt := range opts {
		opt(&spanOpts)
	}
	startTime := spanOpts.StartTime
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanContext := otel.SpanContext{
		TraceID:    t.getTraceID(ctx, &spanOpts),
		SpanID:     t.getSpanID(),
		TraceFlags: 0,
	}
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		recording:      spanOpts.Record,
		Attributes: otel.NewMap(otel.MapUpdate{
			MultiKV: spanOpts.Attributes,
		}),
		StartTime:    startTime,
		EndTime:      time.Time{},
		ParentSpanID: t.getParentSpanID(ctx, &spanOpts),
		Events:       nil,
		SpanKind:     otel.ValidateSpanKind(spanOpts.SpanKind),
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = otel.SetCurrentSpan(ctx, span)
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

func (t *MockTracer) getTraceID(ctx context.Context, spanOpts *otel.SpanOptions) otel.TraceID {
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

func (t *MockTracer) getParentSpanID(ctx context.Context, spanOpts *otel.SpanOptions) otel.SpanID {
	if parent := t.getParentSpanContext(ctx, spanOpts); parent.IsValid() {
		return parent.SpanID
	}
	return otel.SpanID{}
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, spanOpts *otel.SpanOptions) otel.SpanContext {
	if spanOpts.Relation.RelationshipType == otel.ChildOfRelationship &&
		spanOpts.Relation.SpanContext.IsValid() {
		return spanOpts.Relation.SpanContext
	}
	if parentSpanContext := otel.CurrentSpan(ctx).SpanContext(); parentSpanContext.IsValid() {
		return parentSpanContext
	}
	return otel.EmptySpanContext()
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
	CtxAttributes otel.Map
	Timestamp     time.Time
	Msg           string
	Attributes    otel.Map
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer otel.Tracer
	spanContext    otel.SpanContext
	SpanKind       otel.SpanKind
	recording      bool

	Attributes   otel.Map
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

func (s *MockSpan) SetStatus(status codes.Code) {
	s.SetAttribute(NameKey.Uint32(uint32(status)))
}

func (s *MockSpan) SetName(name string) {
	s.SetAttribute(NameKey.String(name))
}

func (s *MockSpan) SetError(v bool) {
	s.SetAttribute(ErrorKey.Bool(v))
}

func (s *MockSpan) SetAttribute(attribute otel.KeyValue) {
	s.applyUpdate(otel.MapUpdate{
		SingleKV: attribute,
	})
}

func (s *MockSpan) SetAttributes(attributes ...otel.KeyValue) {
	s.applyUpdate(otel.MapUpdate{
		MultiKV: attributes,
	})
}

func (s *MockSpan) applyUpdate(update otel.MapUpdate) {
	s.Attributes = s.Attributes.Apply(update)
}

func (s *MockSpan) End(options ...otel.EndOption) {
	if !s.EndTime.IsZero() {
		return // already finished
	}
	endOpts := otel.EndOptions{}

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

func (s *MockSpan) Tracer() otel.Tracer {
	return s.officialTracer
}

func (s *MockSpan) AddEvent(ctx context.Context, msg string, attrs ...otel.KeyValue) {
	s.AddEventWithTimestamp(ctx, time.Now(), msg, attrs...)
}

func (s *MockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...otel.KeyValue) {
	s.Events = append(s.Events, MockEvent{
		CtxAttributes: otel.FromContext(ctx),
		Timestamp:     timestamp,
		Msg:           msg,
		Attributes: otel.NewMap(otel.MapUpdate{
			MultiKV: attrs,
		}),
	})
}

func (s *MockSpan) AddLink(link otel.Link) {
	// TODO
}

func (s *MockSpan) Link(sc otel.SpanContext, attrs ...otel.KeyValue) {
	// TODO
}

func (s *MockSpan) OverrideTracer(tracer otel.Tracer) {
	s.officialTracer = tracer
}
