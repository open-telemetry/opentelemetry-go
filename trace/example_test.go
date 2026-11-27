// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

func ExampleSpanFromContext() {
	// No Span is set in the context, so a no-op Span is returned.
	span := trace.SpanFromContext(context.Background())

	fmt.Println("recording:", span.IsRecording())
	fmt.Println("valid span context:", span.SpanContext().IsValid())
	// Output:
	// recording: false
	// valid span context: false
}

func ExampleContextWithRemoteSpanContext() {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x01},
		SpanID:  trace.SpanID{0x01},
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)

	fmt.Println("original remote:", sc.IsRemote())
	fmt.Println("extracted remote:", trace.SpanContextFromContext(ctx).IsRemote())
	// Output:
	// original remote: false
	// extracted remote: true
}
