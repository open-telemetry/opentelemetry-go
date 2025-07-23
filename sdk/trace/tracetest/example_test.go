// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
    "context"
    "testing"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// simulateWorkflow creates a parent span with a nested child span.
func simulateWorkflow(ctx context.Context) {
    tracer := otel.Tracer("example/workflow")
    ctx, parent := tracer.Start(ctx, "workflow")
    defer parent.End()

    parent.SetAttributes(attribute.String("workflow.phase", "start"))

    _, child := tracer.Start(ctx, "step-1")
    child.SetAttributes(attribute.Int("step.order", 1))
    child.AddEvent("Step 1 processing started")
    child.End()
}

// Example demonstrates parent and child span recording in unit tests.
func Example() {
    t := &testing.T{} // Provided by testing framework.
    
    recorder := tracetest.NewSpanRecorder()
    tp := trace.NewTracerProvider(trace.WithSpanProcessor(recorder))
    defer func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            panic(err)
        }
    }()

    original := otel.GetTracerProvider()
    otel.SetTracerProvider(tp)
    defer otel.SetTracerProvider(original)

    simulateWorkflow(context.Background())

    // Verify expected spans were created
    spans := recorder.Ended()
    expectedSpanCount := 2
    if len(spans) != expectedSpanCount {
        t.Errorf("Expected %d spans, got %d", expectedSpanCount, len(spans))
        return
    }

    // Verify first span (step-1)
    stepSpan := spans[0]
    expectedName := "step-1"
    if stepSpan.Name() != expectedName {
        t.Errorf("Expected span name %s, got %s", expectedName, stepSpan.Name())
        return
    }

    expectedAttrCount := 1
    if len(stepSpan.Attributes()) != expectedAttrCount {
        t.Errorf("Expected %d attributes, got %d", expectedAttrCount, len(stepSpan.Attributes()))
        return
    }

    // Verify second span (workflow)
    workflowSpan := spans[1]
    expectedWorkflowName := "workflow"
    if workflowSpan.Name() != expectedWorkflowName {
        t.Errorf("Expected span name %s, got %s", expectedWorkflowName, workflowSpan.Name())
        return
    }

    // Output:
    //
}
