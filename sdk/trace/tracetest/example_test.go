// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)


func processItem(ctx context.Context, itemID string) {
	tracer := otel.Tracer("example/processor")
	ctx, span := tracer.Start(ctx, "processItem")
	defer span.End()
	
	span.SetAttributes(attribute.String("item.id", itemID))
	// Simulate some processing work
	span.AddEvent("processing started")
}


func Example() {
	// Create a span recorder to capture spans
	recorder := tracetest.NewSpanRecorder()
	
	// Set up a tracer provider with our recorder
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(recorder),
	)
	defer tp.Shutdown(context.Background())
	
	// Set the tracer provider globally
	originalTP := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(originalTP)
	
	// Execute the instrumented function
	processItem(context.Background(), "item-123")
	
	// Retrieve and inspect the recorded spans
	spans := recorder.Ended()
	
	if len(spans) > 0 {
		span := spans[0]
		fmt.Printf("Span name: %s\n", span.Name())
		
		// Check for specific attributes
		attrs := span.Attributes()
		for _, attr := range attrs {
			if attr.Key == "item.id" {
				fmt.Printf("Item ID: %s\n", attr.Value.AsString())
			}
		}
	}
	
	// Output:
	// Span name: processItem
	// Item ID: item-123
}

func ExampleSpanRecorder() {
	// Create a new span recorder
	recorder := tracetest.NewSpanRecorder()
	
	// Configure tracer provider to use the recorder
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(recorder),
	)
	defer tp.Shutdown(context.Background())
	
	// Create and end a span
	tracer := tp.Tracer("example")

	_, span := tracer.Start(context.Background(), "example-operation")
	span.SetAttributes(attribute.String("example.key", "example-value"))
	span.End()
	
	// Verify spans were recorded
	spans := recorder.Ended()
	fmt.Printf("Recorded %d span(s)\n", len(spans))
	
	if len(spans) > 0 {
		fmt.Printf("First span: %s\n", spans[0].Name())
	}
	
	// Output:
	// Recorded 1 span(s)
	// First span: example-operation
}

// ExampleSpanRecorder_started demonstrates accessing spans that have been started
// but not necessarily ended.
func ExampleSpanRecorder_started() {
	recorder := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(recorder),
	)
	defer tp.Shutdown(context.Background())
	
	otel.SetTracerProvider(tp)
	
	tracer := otel.Tracer("example")
	_, span := tracer.Start(context.Background(), "long-running-operation")
	
	// Check started spans (including those not yet ended)
	started := recorder.Started()
	fmt.Printf("Started spans: %d\n", len(started))
	
	span.End()
	
	// Now check ended spans
	ended := recorder.Ended()
	fmt.Printf("Ended spans: %d\n", len(ended))

}