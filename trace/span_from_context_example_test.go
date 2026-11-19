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
