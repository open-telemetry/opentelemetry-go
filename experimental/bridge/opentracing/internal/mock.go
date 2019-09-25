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
	otelkey "go.opentelemetry.io/api/key"
	oteltag "go.opentelemetry.io/api/tag"
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
	Resources             oteltag.Map
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
		Resources:             oteltag.NewEmptyMap(),
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
	defer span.Finish()
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
		TraceID:      t.getTraceID(ctx, &spanOpts),
		SpanID:       t.getSpanID(),
		TraceOptions: 0,
	}
	span := &MockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		recording:      spanOpts.RecordEvent,
		Attributes:     oteltag.NewMap(upsertMultiMapUpdate(spanOpts.Attributes...)),
		StartTime:      startTime,
		FinishTime:     time.Time{},
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
	CtxAttributes oteltag.Map
	Timestamp     time.Time
	Msg           string
	Attributes    oteltag.Map
}

type MockSpan struct {
	mockTracer     *MockTracer
	officialTracer oteltrace.Tracer
	spanContext    otelcore.SpanContext
	recording      bool

	Attributes   oteltag.Map
	StartTime    time.Time
	FinishTime   time.Time
	ParentSpanID uint64
	Events       []MockEvent
}

var _ oteltrace.Span = &MockSpan{}
var _ migration.OverrideTracerSpanExtension = &MockSpan{}

func (s *MockSpan) SpanContext() otelcore.SpanContext {
	return s.spanContext
}

func (s *MockSpan) IsRecordingEvents() bool {
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

func (s *MockSpan) ModifyAttribute(mutator oteltag.Mutator) {
	s.applyUpdate(oteltag.MapUpdate{
		SingleMutator: mutator,
	})
}

func (s *MockSpan) ModifyAttributes(mutators ...oteltag.Mutator) {
	s.applyUpdate(oteltag.MapUpdate{
		MultiMutator: mutators,
	})
}

func (s *MockSpan) applyUpdate(update oteltag.MapUpdate) {
	s.Attributes = s.Attributes.Apply(update)
}

func (s *MockSpan) Finish(options ...oteltrace.FinishOption) {
	if !s.FinishTime.IsZero() {
		return // already finished
	}
	finishOpts := oteltrace.FinishOptions{}

	for _, opt := range options {
		opt(&finishOpts)
	}

	finishTime := finishOpts.FinishTime
	if finishTime.IsZero() {
		finishTime = time.Now()
	}
	s.FinishTime = finishTime
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
		CtxAttributes: oteltag.FromContext(ctx),
		Timestamp:     timestamp,
		Msg:           msg,
		Attributes:    oteltag.NewMap(upsertMultiMapUpdate(attrs...)),
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

func upsertMapUpdate(kv otelcore.KeyValue) oteltag.MapUpdate {
	return singleMutatorMapUpdate(oteltag.UPSERT, kv)
}

func upsertMultiMapUpdate(kvs ...otelcore.KeyValue) oteltag.MapUpdate {
	return multiMutatorMapUpdate(oteltag.UPSERT, kvs...)
}

func singleMutatorMapUpdate(op oteltag.MutatorOp, kv otelcore.KeyValue) oteltag.MapUpdate {
	return oteltag.MapUpdate{
		SingleMutator: keyValueToMutator(op, kv),
	}
}

func multiMutatorMapUpdate(op oteltag.MutatorOp, kvs ...otelcore.KeyValue) oteltag.MapUpdate {
	return oteltag.MapUpdate{
		MultiMutator: keyValuesToMutators(op, kvs...),
	}
}

func keyValuesToMutators(op oteltag.MutatorOp, kvs ...otelcore.KeyValue) []oteltag.Mutator {
	var mutators []oteltag.Mutator
	for _, kv := range kvs {
		mutators = append(mutators, keyValueToMutator(op, kv))
	}
	return mutators
}

func keyValueToMutator(op oteltag.MutatorOp, kv otelcore.KeyValue) oteltag.Mutator {
	return oteltag.Mutator{
		MutatorOp: op,
		KeyValue:  kv,
	}
}
