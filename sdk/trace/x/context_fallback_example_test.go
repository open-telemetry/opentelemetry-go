// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/sdk/trace/x"
	"go.opentelemetry.io/otel/trace"
)

// The provider applies the startup snapshot to every eligible root, including
// roots started after earlier work has ended. WithNewRoot opts out explicitly.
func ExampleWithContextFallback() {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	carrier := propagation.MapCarrier{
		"traceparent": "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01",
		"baggage":     "job=build",
	}
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		x.WithContextFallback(propagator, carrier),
		sdktrace.WithSpanProcessor(recorder),
	)
	defer func() { _ = tp.Shutdown(context.Background()) }()
	tracer := tp.Tracer("example.com/worker")

	ctx, job := tracer.Start(context.Background(), "job")
	fmt.Println("job baggage:", baggage.FromContext(ctx).Member("job").Value())
	job.End()

	_, later := tracer.Start(context.Background(), "later")
	later.End()
	_, independent := tracer.Start(context.Background(), "independent", trace.WithNewRoot())
	independent.End()

	for _, span := range recorder.Ended() {
		fmt.Printf("%s parent: %s\n", span.Name(), span.Parent().SpanID())
	}
	// Output:
	// job baggage: build
	// job parent: 0102030405060708
	// later parent: 0102030405060708
	// independent parent: 0000000000000000
}

// Ordinary scoped extraction associates only work that receives the extracted
// context with the invocation. A separate context starts independent work.
func ExampleWithContextFallback_scopedAlternative() {
	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier{
		"traceparent": "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01",
	}
	invocationCtx := propagator.Extract(context.Background(), carrier)
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	defer func() { _ = tp.Shutdown(context.Background()) }()
	tracer := tp.Tracer("example.com/worker")

	_, job := tracer.Start(invocationCtx, "job")
	job.End()
	_, request := tracer.Start(context.Background(), "unrelated request")
	request.End()

	for _, span := range recorder.Ended() {
		fmt.Printf("%s parent: %s\n", span.Name(), span.Parent().SpanID())
	}
	// Output:
	// job parent: 0102030405060708
	// unrelated request parent: 0000000000000000
}
