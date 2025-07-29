// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"fmt"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func Example_simulateWorkflow() {
	ctx := context.Background()

	// Set up an in-memory span recorder and tracer provider.
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(recorder),
	)

	defer func() {
		_ = tp.Shutdown(ctx)
	}()

	tracer := tp.Tracer("example/workflow")

	// Parent span "workflow"
	ctx, workflowSpan := tracer.Start(ctx, "workflow")

	// Child span "step-1"
	_, stepSpan := tracer.Start(ctx, "step-1")

	// End spans in reverse order
	stepSpan.End()
	workflowSpan.End()

	// Print span names in the order they ended.
	for _, s := range recorder.Ended() {
		fmt.Println(s.Name())
	}

	// Output:
	// step-1
	// workflow
}