// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// simulateWorkflow starts a “workflow” span, then a nested “step-1” span.
// It attaches attributes and an event to the child span.
func simulateWorkflow(ctx context.Context) {
	tracer := otel.Tracer("example/workflow")

	// Start parent span “workflow”
	ctx, workflowSpan := tracer.Start(ctx, "workflow")
	defer workflowSpan.End()

	workflowSpan.SetAttributes(attribute.String("workflow.phase", "start"))

	// Start child span “step-1”
	_, stepSpan := tracer.Start(ctx, "step-1")
	defer stepSpan.End()

	stepSpan.SetAttributes(attribute.Int("step.order", 1))
	stepSpan.AddEvent("Step 1 processing started")
}

func TestSimulateWorkflowCreatesSpans(t *testing.T) {
	// Prepare an in-memory recorder to capture completed spans.
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(recorder),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			t.Fatalf("failed to shut down tracer provider: %v", err)
		}
	}()

	// Swap in our test TracerProvider
	originalProvider := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(originalProvider)

	// Run the workflow
	simulateWorkflow(context.Background())

	endedSpans := recorder.Ended()
	const wantSpanCount = 2
	if len(endedSpans) != wantSpanCount {
		t.Fatalf("expected %d spans, got %d", wantSpanCount, len(endedSpans))
	}

	// The recorder yields spans in the order they ended: step-1, then workflow.

	// Check child span “step-1”
	stepSpan := endedSpans[0]
	if got, want := stepSpan.Name(), "step-1"; got != want {
		t.Errorf("child span name: got %q, want %q", got, want)
	}
	if len(stepSpan.Attributes()) != 1 {
		t.Errorf("child span attribute count: got %d, want 1", len(stepSpan.Attributes()))
	}

	// Check parent span “workflow”
	workflowSpan := endedSpans[1]
	if got, want := workflowSpan.Name(), "workflow"; got != want {
		t.Errorf("parent span name: got %q, want %q", got, want)
	}
}
