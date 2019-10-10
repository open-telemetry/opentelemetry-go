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

	otelcore "go.opentelemetry.io/api/core"
	oteldctx "go.opentelemetry.io/api/distributedcontext"
	otelkey "go.opentelemetry.io/api/key"
	oteltrace "go.opentelemetry.io/api/trace"

	migration "go.opentelemetry.io/experimental/bridge/opentracing/migration"
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
	SpareSpanIDs          []uint64
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

func (t *MockTracer) WithResources(attributes ...otelcore.KeyValue) oteltrace.Tracer {
	t.Resources = t.Resources.Apply(upsertMultiMapUpdate(attributes...))
	return t
}

func (t *MockTracer) WithComponent(name string) oteltrace.Tracer {
	return t.WithResources(otelkey.New("component").String(name))
}

func (t *MockTracer) WithService(name string) oteltrace.Tracer {
	return t.WithResources(otelkey.New("service").String(name))
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
		Attributes:     oteldctx.NewMap(upsertMultiMapUpdate(spanOpts.Attributes...)),
		StartTime:      startTime,
		EndTime:        time.Time{},
		ParentSpanID:   t.getParentSpanID(ctx, &spanOpts),
		Events:         nil,
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
	uints := t.getNRandUint64(2)
	return otelcore.TraceID{
		High: uints[0],
		Low:  uints[1],
	}
}

func (t *MockTracer) getParentSpanID(ctx context.Context, spanOpts *oteltrace.SpanOptions) uint64 {
	if parent := t.getParentSpanContext(ctx, spanOpts); parent.IsValid() {
		return parent.SpanID
	}
	return 0
}

func (t *MockTracer) getParentSpanContext(ctx context.Context, spanOpts *oteltrace.SpanOptions) otelcore.SpanContext {
	if spanOpts.Reference.RelationshipType == oteltrace.ChildOfRelationship &&
		spanOpts.Reference.SpanContext.IsValid() {
		return spanOpts.Reference.SpanContext
	}
	if parentSpanContext := oteltrace.CurrentSpan(ctx).SpanContext(); parentSpanContext.IsValid() {
		return parentSpanContext
	}
	return otelcore.EmptySpanContext()
}

func (t *MockTracer) getSpanID() uint64 {
	if len(t.SpareSpanIDs) > 0 {
		spanID := t.SpareSpanIDs[0]
		t.SpareSpanIDs = t.SpareSpanIDs[1:]
		if len(t.SpareSpanIDs) == 0 {
			t.SpareSpanIDs = nil
		}
		return spanID
	}
	return t.getRandUint64()
}

func (t *MockTracer) getRandUint64() uint64 {
	return t.getNRandUint64(1)[0]
}

func (t *MockTracer) getNRandUint64(n int) []uint64 {
	uints := make([]uint64, n)
	t.randLock.Lock()
	defer t.randLock.Unlock()
	for i := 0; i < n; i++ {
		uints[i] = t.rand.Uint64()
	}
	return uints
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
	recording      bool

	Attributes   oteldctx.Map
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID uint64
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
	s.applyUpdate(upsertMapUpdate(attribute))
}

func (s *MockSpan) SetAttributes(attributes ...otelcore.KeyValue) {
	s.applyUpdate(upsertMultiMapUpdate(attributes...))
}

func (s *MockSpan) ModifyAttribute(mutator oteldctx.Mutator) {
	s.applyUpdate(oteldctx.MapUpdate{
		SingleMutator: mutator,
	})
}

func (s *MockSpan) ModifyAttributes(mutators ...oteldctx.Mutator) {
	s.applyUpdate(oteldctx.MapUpdate{
		MultiMutator: mutators,
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
		Attributes:    oteldctx.NewMap(upsertMultiMapUpdate(attrs...)),
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

func upsertMapUpdate(kv otelcore.KeyValue) oteldctx.MapUpdate {
	return singleMutatorMapUpdate(oteldctx.UPSERT, kv)
}

func upsertMultiMapUpdate(kvs ...otelcore.KeyValue) oteldctx.MapUpdate {
	return multiMutatorMapUpdate(oteldctx.UPSERT, kvs...)
}

func singleMutatorMapUpdate(op oteldctx.MutatorOp, kv otelcore.KeyValue) oteldctx.MapUpdate {
	return oteldctx.MapUpdate{
		SingleMutator: keyValueToMutator(op, kv),
	}
}

func multiMutatorMapUpdate(op oteldctx.MutatorOp, kvs ...otelcore.KeyValue) oteldctx.MapUpdate {
	return oteldctx.MapUpdate{
		MultiMutator: keyValuesToMutators(op, kvs...),
	}
}

func keyValuesToMutators(op oteldctx.MutatorOp, kvs ...otelcore.KeyValue) []oteldctx.Mutator {
	var mutators []oteldctx.Mutator
	for _, kv := range kvs {
		mutators = append(mutators, keyValueToMutator(op, kv))
	}
	return mutators
}

func keyValueToMutator(op oteldctx.MutatorOp, kv otelcore.KeyValue) oteldctx.Mutator {
	return oteldctx.Mutator{
		MutatorOp: op,
		KeyValue:  kv,
	}
}
