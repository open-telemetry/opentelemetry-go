// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetest_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func ExampleSpanRecorder() {
	ctx := context.Background()

	// Set up an in-memory span recorder and tracer provider.
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(sr),
	)
	defer tp.Shutdown(ctx) //nolint:errcheck // Example code, error handling omitted.

	tracer := tp.Tracer("example/simple")

	// Start and end a span.
	_, span := tracer.Start(ctx, "test-span")
	span.End()

	// Print the recorded span name.
	for _, s := range sr.Ended() {
		fmt.Println(s.Name())
	}

	// Output:
	// test-span
}
