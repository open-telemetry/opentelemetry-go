// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal_test

import (
	"context"
	"testing"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel/bridge/opencensus/internal"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/oc2otel"
	"go.opentelemetry.io/otel/bridge/opencensus/internal/otel2oc"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"
)

type handler struct{ err error }

func (h *handler) Handle(e error) { h.err = e }

func withHandler() (*handler, func()) {
	h := new(handler)
	original := internal.Handle
	internal.Handle = h.Handle
	return h, func() { internal.Handle = original }
}

type tracer struct {
	embedded.Tracer

	ctx  context.Context
	name string
	opts []trace.SpanStartOption
}

func (t *tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	t.ctx, t.name, t.opts = ctx, name, opts
	sub := noop.NewTracerProvider().Tracer("testing")
	return sub.Start(ctx, name, opts...)
}

type ctxKey string

func TestTracerStartSpan(t *testing.T) {
	h, restore := withHandler()
	defer restore()

	otelTracer := &tracer{}
	ocTracer := internal.NewTracer(otelTracer)

	ctx := context.WithValue(context.Background(), ctxKey("key"), "value")
	name := "testing span"
	ocTracer.StartSpan(ctx, name, octrace.WithSpanKind(octrace.SpanKindClient))
	if h.err != nil {
		t.Fatalf("OC tracer.StartSpan errored: %v", h.err)
	}

	if otelTracer.ctx != ctx {
		t.Error("OTel tracer.Start called with wrong context")
	}
	if otelTracer.name != name {
		t.Error("OTel tracer.Start called with wrong name")
	}
	sk := trace.SpanKindClient
	c := trace.NewSpanStartConfig(otelTracer.opts...)
	if c.SpanKind() != sk {
		t.Errorf("OTel tracer.Start called with wrong options: %#v", c)
	}
}

func TestTracerStartSpanReportsErrors(t *testing.T) {
	h, restore := withHandler()
	defer restore()

	ocTracer := internal.NewTracer(&tracer{})
	ocTracer.StartSpan(context.Background(), "", octrace.WithSampler(octrace.AlwaysSample()))
	if h.err == nil {
		t.Error("OC tracer.StartSpan no error when converting Sampler")
	}
}

func TestTracerStartSpanWithRemoteParent(t *testing.T) {
	otelTracer := new(tracer)
	ocTracer := internal.NewTracer(otelTracer)
	sc := octrace.SpanContext{TraceID: [16]byte{1}, SpanID: [8]byte{1}}
	converted := oc2otel.SpanContext(sc).WithRemote(true)

	ocTracer.StartSpanWithRemoteParent(context.Background(), "", sc)

	got := trace.SpanContextFromContext(otelTracer.ctx)
	if !got.Equal(converted) {
		t.Error("tracer.StartSpanWithRemoteParent failed to set remote parent")
	}
}

func TestTracerFromContext(t *testing.T) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	tracer := noop.NewTracerProvider().Tracer("TestTracerFromContext")
	// Test using the fact that the No-Op span will propagate a span context .
	ctx, _ = tracer.Start(ctx, "test")

	got := internal.NewTracer(tracer).FromContext(ctx).SpanContext()
	// Do not test the conversion, only the propagation.
	want := otel2oc.SpanContext(sc)
	if got != want {
		t.Errorf("tracer.FromContext returned wrong context: %#v", got)
	}
}

func TestTracerNewContext(t *testing.T) {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: [16]byte{1},
		SpanID:  [8]byte{1},
	})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)

	tracer := noop.NewTracerProvider().Tracer("TestTracerNewContext")
	// Test using the fact that the No-Op span will propagate a span context .
	_, s := tracer.Start(ctx, "test")

	ocTracer := internal.NewTracer(tracer)
	ctx = ocTracer.NewContext(context.Background(), internal.NewSpan(s))
	got := trace.SpanContextFromContext(ctx)

	if !got.Equal(sc) {
		t.Error("tracer.NewContext did not attach Span to context")
	}
}

type differentSpan struct {
	octrace.SpanInterface
}

func (s *differentSpan) String() string { return "testing span" }

func TestTracerNewContextErrors(t *testing.T) {
	h, restore := withHandler()
	defer restore()

	ocTracer := internal.NewTracer(&tracer{})
	ocSpan := octrace.NewSpan(&differentSpan{})
	ocTracer.NewContext(context.Background(), ocSpan)
	if h.err == nil {
		t.Error("tracer.NewContext did not error for unrecognized span")
	}
}
