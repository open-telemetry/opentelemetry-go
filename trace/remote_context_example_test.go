// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

func ExampleContextWithRemoteSpanContext() {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x001},
		SpanID:  trace.SpanID{0x001},
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)

	fmt.Println("original remote:", sc.IsRemote())
	fmt.Println("extracted remote:", trace.SpanContextFromContext(ctx).IsRemote())
	// Output:
	// original remote: false
	// extracted remote: true
}
