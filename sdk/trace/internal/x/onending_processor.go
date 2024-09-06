// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/trace/internal/x"

import "go.opentelemetry.io/otel/trace"

// OnEndingSpanProcessor represents span processors that allow mutating spans
// just before they are ended and made immutable.
//
// This is useful for custom processor implementations that want to mutate
// spans when they are finished, and before they are made immutable, such as
// implementing tail-based sampling.
type OnEndingSpanProcessor interface {
	// OnEnding is called while the span is finished, and spans are still
	// mutable.
	//
	// This method is called synchronously during the span's End operation,
	// therefore it should not block or throw an exception.
	// If multiple [SpanProcessor] are registered, their `OnEnding` callbacks are
	// invoked in the order they have been registered.
	//
	// [SpanProcessor]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#SpanProcessor
	OnEnding(trace.Span)
}
