package tracetest_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func doSomething(ctx context.Context) {
	tr := otel.Tracer("example.com/test")
	_, span := tr.Start(ctx, "doSomething")
	defer span.End()
	// simulate some logic
}

func TestDoSomething(t *testing.T) {
	// Create a SpanRecorder
	sr := tracetest.NewSpanRecorder()

	// Create a TracerProvider and register the recorder
	tp := trace.NewTracerProvider()
	tp.RegisterSpanProcessor(sr)

	// Set the provider globally
	otel.SetTracerProvider(tp)

	// Run the function
	ctx := context.Background()
	doSomething(ctx)

	// Get spans recorded
	spans := sr.Ended()

	if len(spans) != 1 {
		t.Fatalf("Expected 1 span, got %d", len(spans))
	}

	if spans[0].Name() != "doSomething" {
		t.Errorf("Expected span name 'doSomething', got %s", spans[0].Name())
	}
}
