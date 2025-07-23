// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// simulateWorkflow starts a “workflow” span and a nested “step-1” span.
func simulateWorkflow(ctx context.Context) {
	tracer := otel.Tracer("example/workflow")

	// Parent span “workflow”
	ctx, workflowSpan := tracer.Start(ctx, "workflow")
	defer workflowSpan.End()
	workflowSpan.SetAttributes(attribute.String("workflow.phase", "start"))

	// Child span “step-1”
	_, stepSpan := tracer.Start(ctx, "step-1")
	defer stepSpan.End()
	stepSpan.SetAttributes(attribute.Int("step.order", 1))
	stepSpan.AddEvent("Step 1 processing started")
}

// Example_simulateWorkflow is a runnable example.
// The // Output: comment makes it an executable doc test [2].
func Example_simulateWorkflow() {
	ctx := context.Background()

	// In-memory span recorder.
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(recorder),
	)
	defer tp.Shutdown(ctx)

	// Swap in the test TracerProvider and restore afterwards.
	orig := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(orig)

	// Run the workflow.
	simulateWorkflow(ctx)

	// Print span names in the order they ended.
	for _, s := range recorder.Ended() {
		fmt.Println(s.Name())
	}

	// Output:
	// step-1
	// workflow
}
