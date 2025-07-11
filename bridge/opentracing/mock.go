// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing // import "go.opentelemetry.io/otel/bridge/opentracing"

import (
	"context"
	"math/rand/v2"
	"reflect"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opentracing/migration"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"
)

//nolint:revive // ignoring missing comments for unexported global variables in an internal package.
var (
	statusCodeKey    = attribute.Key("status.code")
	statusMessageKey = attribute.Key("status.message")
	errorKey         = attribute.Key("error")
	nameKey          = attribute.Key("name")
)

type mockContextKeyValue struct {
	Key   interface{}
	Value interface{}
}

type mockTracer struct {
	embedded.Tracer

	FinishedSpans         []*mockSpan
	SpareTraceIDs         []trace.TraceID
	SpareSpanIDs          []trace.SpanID
	SpareContextKeyValues []mockContextKeyValue
	TraceFlags            trace.TraceFlags

	randLock sync.Mutex
	rand     *rand.ChaCha8
}

var (
	_ trace.Tracer                                  = &mockTracer{}
	_ migration.DeferredContextSetupTracerExtension = &mockTracer{}
)

func newMockTracer() *mockTracer {
	u := rand.Uint32()
	seed := [32]byte{byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24)}
	return &mockTracer{
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.NewChaCha8(seed),
	}
}

// Start returns a new trace span with the given name and options.
func (t *mockTracer) Start(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	config := trace.NewSpanStartConfig(opts...)
	startTime := config.Timestamp()
	if startTime.IsZero() {
		startTime = time.Now()
	}
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    t.getTraceID(ctx, &config),
		SpanID:     t.getSpanID(),
		TraceFlags: t.TraceFlags,
	})
	span := &mockSpan{
		mockTracer:     t,
		officialTracer: t,
		spanContext:    spanContext,
		Attributes:     config.Attributes(),
		StartTime:      startTime,
		EndTime:        time.Time{},
		ParentSpanID:   t.getParentSpanID(ctx, &config),
		Events:         nil,
		SpanKind:       trace.ValidateSpanKind(config.SpanKind()),
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = trace.ContextWithSpan(ctx, span)
		ctx = t.addSpareContextValue(ctx)
	}
	return ctx, span
}

func (t *mockTracer) addSpareContextValue(ctx context.Context) context.Context {
	if len(t.SpareContextKeyValues) > 0 {
		pair := t.SpareContextKeyValues[0]
		t.SpareContextKeyValues[0] = mockContextKeyValue{}
		t.SpareContextKeyValues = t.SpareContextKeyValues[1:]
		if len(t.SpareContextKeyValues) == 0 {
			t.SpareContextKeyValues = nil
		}
		ctx = context.WithValue(ctx, pair.Key, pair.Value)
	}
	return ctx
}

func (t *mockTracer) getTraceID(ctx context.Context, config *trace.SpanConfig) trace.TraceID {
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

func (t *mockTracer) getParentSpanID(ctx context.Context, config *trace.SpanConfig) trace.SpanID {
	if parent := t.getParentSpanContext(ctx, config); parent.IsValid() {
		return parent.SpanID()
	}
	return trace.SpanID{}
}

func (t *mockTracer) getParentSpanContext(ctx context.Context, config *trace.SpanConfig) trace.SpanContext {
	if !config.NewRoot() {
		return trace.SpanContextFromContext(ctx)
	}
	return trace.SpanContext{}
}

func (t *mockTracer) getSpanID() trace.SpanID {
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

func (t *mockTracer) getRandSpanID() trace.SpanID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	sid := trace.SpanID{}
	_, _ = t.rand.Read(sid[:])

	return sid
}

func (t *mockTracer) getRandTraceID() trace.TraceID {
	t.randLock.Lock()
	defer t.randLock.Unlock()

	tid := trace.TraceID{}
	_, _ = t.rand.Read(tid[:])

	return tid
}

// DeferredContextSetupHook implements the DeferredContextSetupTracerExtension interface.
func (t *mockTracer) DeferredContextSetupHook(ctx context.Context, span trace.Span) context.Context {
	return t.addSpareContextValue(ctx)
}

type mockEvent struct {
	Timestamp  time.Time
	Name       string
	Attributes []attribute.KeyValue
}

type mockLink struct {
	SpanContext trace.SpanContext
	Attributes  []attribute.KeyValue
}

type mockSpan struct {
	embedded.Span

	mockTracer     *mockTracer
	officialTracer trace.Tracer
	spanContext    trace.SpanContext
	SpanKind       trace.SpanKind
	recording      bool

	Attributes   []attribute.KeyValue
	StartTime    time.Time
	EndTime      time.Time
	ParentSpanID trace.SpanID
	Events       []mockEvent
	Links        []mockLink
}

var (
	_ trace.Span                            = &mockSpan{}
	_ migration.OverrideTracerSpanExtension = &mockSpan{}
)

func (s *mockSpan) SpanContext() trace.SpanContext {
	return s.spanContext
}

func (s *mockSpan) IsRecording() bool {
	return s.recording
}

func (s *mockSpan) SetStatus(code codes.Code, msg string) {
	s.SetAttributes(statusCodeKey.Int(int(code)), statusMessageKey.String(msg))
}

func (s *mockSpan) SetName(name string) {
	s.SetAttributes(nameKey.String(name))
}

func (s *mockSpan) SetError(v bool) {
	s.SetAttributes(errorKey.Bool(v))
}

func (s *mockSpan) SetAttributes(attributes ...attribute.KeyValue) {
	s.applyUpdate(attributes)
}

func (s *mockSpan) applyUpdate(update []attribute.KeyValue) {
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

func (s *mockSpan) End(options ...trace.SpanEndOption) {
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

func (s *mockSpan) RecordError(err error, opts ...trace.EventOption) {
	if err == nil {
		return // no-op on nil error
	}

	if !s.EndTime.IsZero() {
		return // already finished
	}

	s.SetStatus(codes.Error, "")
	opts = append(opts, trace.WithAttributes(
		semconv.ExceptionType(reflect.TypeOf(err).String()),
		semconv.ExceptionMessage(err.Error()),
	))
	s.AddEvent(semconv.ExceptionEventName, opts...)
}

func (s *mockSpan) Tracer() trace.Tracer {
	return s.officialTracer
}

func (s *mockSpan) AddEvent(name string, o ...trace.EventOption) {
	c := trace.NewEventConfig(o...)
	s.Events = append(s.Events, mockEvent{
		Timestamp:  c.Timestamp(),
		Name:       name,
		Attributes: c.Attributes(),
	})
}

func (s *mockSpan) AddLink(link trace.Link) {
	s.Links = append(s.Links, mockLink{
		SpanContext: link.SpanContext,
		Attributes:  link.Attributes,
	})
}

func (s *mockSpan) OverrideTracer(tracer trace.Tracer) {
	s.officialTracer = tracer
}

func (s *mockSpan) TracerProvider() trace.TracerProvider { return noop.NewTracerProvider() }
