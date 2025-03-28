// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	testpb "github.com/opentracing-contrib/go-grpc/test/otgrpc_testing"
	ot "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/attribute"
	ototel "go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/bridge/opentracing/migration"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"
)

type testGRPCServer struct{}

var (
	statusCodeKey    = attribute.Key("status.code")
	statusMessageKey = attribute.Key("status.message")
	errorKey         = attribute.Key("error")
	nameKey          = attribute.Key("name")
)

type mockContextKeyValue struct {
	Key   any
	Value any
}

type mockTracer struct {
	embedded.Tracer

	FinishedSpans         []*mockSpan
	SpareTraceIDs         []trace.TraceID
	SpareSpanIDs          []trace.SpanID
	SpareContextKeyValues []mockContextKeyValue
	TraceFlags            trace.TraceFlags

	randLock sync.Mutex
	rand     *rand.Rand
}

func newMockTracer() *mockTracer {
	return &mockTracer{
		FinishedSpans:         nil,
		SpareTraceIDs:         nil,
		SpareSpanIDs:          nil,
		SpareContextKeyValues: nil,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (t *mockTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
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

func (*testGRPCServer) UnaryCall(ctx context.Context, r *testpb.SimpleRequest) (*testpb.SimpleResponse, error) {
	return &testpb.SimpleResponse{Payload: r.Payload * 2}, nil
}

func (*testGRPCServer) StreamingOutputCall(*testpb.SimpleRequest, testpb.TestService_StreamingOutputCallServer) error {
	return nil
}

func (*testGRPCServer) StreamingInputCall(testpb.TestService_StreamingInputCallServer) error {
	return nil
}

func (*testGRPCServer) StreamingBidirectionalCall(testpb.TestService_StreamingBidirectionalCallServer) error {
	return nil
}

func startTestGRPCServer(t *testing.T, tracer ot.Tracer) (*grpc.Server, net.Addr) {
	lis, _ := net.Listen("tcp", ":0")
	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)),
	)
	testpb.RegisterTestServiceServer(server, &testGRPCServer{})

	go func() {
		err := server.Serve(lis)
		require.NoError(t, err)
	}()

	return server, lis.Addr()
}

func TestBridgeTracer_ExtractAndInject_gRPC(t *testing.T) {
	tracer := newMockTracer()
	bridge := ototel.NewBridgeTracer()
	bridge.SetOpenTelemetryTracer(tracer)
	bridge.SetTextMapPropagator(propagation.TraceContext{})

	srv, addr := startTestGRPCServer(t, bridge)
	defer srv.Stop()

	conn, err := grpc.NewClient(
		addr.String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(bridge)),
	)
	require.NoError(t, err)
	cli := testpb.NewTestServiceClient(conn)

	ctx, cx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cx()
	res, err := cli.UnaryCall(ctx, &testpb.SimpleRequest{Payload: 42})
	require.NoError(t, err)
	assert.EqualValues(t, 84, res.Payload)

	checkSpans := func() bool {
		return len(tracer.FinishedSpans) == 2
	}
	require.Eventuallyf(t, checkSpans, 5*time.Second, 5*time.Millisecond, "expecting two spans")
	assert.Equal(t,
		tracer.FinishedSpans[0].SpanContext().TraceID(),
		tracer.FinishedSpans[1].SpanContext().TraceID(),
		"expecting same trace ID",
	)
}
