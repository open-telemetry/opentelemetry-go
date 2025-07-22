// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"errors"
	"fmt"

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

// Example_nestedSpans demonstrates parent and child span recording.
func Example_nestedSpans() {
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

	spans := recorder.Ended()
	fmt.Printf("Total ended spans: %d\n", len(spans))

	for _, s := range spans {
		fmt.Printf("Span: %s | Attributes: %d\n", s.Name(), len(s.Attributes()))
	}

	// Output:
	// Total ended spans: 2
	// Span: step-1 | Attributes: 1
	// Span: workflow | Attributes: 1
}

// Example_spanWithError simulates a span that records an error.
func Example_spanWithError() {
	recorder := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(recorder))
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := tp.Tracer("example/error")
	_, span := tracer.Start(context.Background(), "run-task")

	err := errors.New("disk not found")
	span.RecordError(err)
	span.SetStatus(1, "error during processing")
	span.End()

	spans := recorder.Ended()
	if len(spans) > 0 {
		fmt.Printf("Span had error? %t\n", spans[0].Status().Code != 0)
	}

	// Output:
	// Span had error? true
}

// Example_spanWithMultipleEvents demonstrates event logging inside a span.
func Example_spanWithMultipleEvents() {
	recorder := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(recorder))
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := tp.Tracer("example/events")
	_, span := tracer.Start(context.Background(), "eventful-operation")
	span.AddEvent("started")
	span.AddEvent("halfway")
	span.AddEvent("done")
	span.End()

	spans := recorder.Ended()
	fmt.Printf("Number of events: %d\n", len(spans[0].Events()))

	// Output:
	// Number of events: 3
}

// Example_spanAttributes demonstrates setting key-value attributes.
func Example_spanAttributes() {
	recorder := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(recorder))
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracer := tp.Tracer("example/attrs")
	_, span := tracer.Start(context.Background(), "my-span")
	span.SetAttributes(
		attribute.String("env", "staging"),
		attribute.Int("version", 2),
		attribute.Bool("feature.enabled", true),
	)
	span.End()

	spanAttrs := recorder.Ended()[0].Attributes()

	for _, attr := range spanAttrs {
		switch attr.Value.Type() {
		case attribute.STRING:
			fmt.Printf("%s: %s\n", attr.Key, attr.Value.AsString())
		case attribute.INT64:
			fmt.Printf("%s: %d\n", attr.Key, attr.Value.AsInt64())
		case attribute.BOOL:
			fmt.Printf("%s: %t\n", attr.Key, attr.Value.AsBool())
		}
	}

	// Output:
	// env: staging
	// version: 2
	// feature.enabled: true
}
