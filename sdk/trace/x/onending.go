// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// OnEndingSpanProcessor is implemented by [sdktrace.SpanProcessor]
// implementations that want to mutate a span during its [sdktrace.Span.End]
// call, after the end timestamp has been set but before the span becomes
// immutable.
//
// This is useful for processors that need to read span data to make a
// sampling or enrichment decision and then modify the span accordingly, such
// as tail-based sampling.
//
// OnEndingSpanProcessor is an optional interface. The [sdktrace.TracerProvider]
// detects it by type assertion at registration time, so no configuration
// beyond registering a processor that implements this interface is required.
//
// The [OnEndingSpanProcessor.OnEnding] method is called synchronously within
// [sdktrace.Span.End], therefore it must not block or panic. If multiple
// [sdktrace.SpanProcessor] implementations are registered, their OnEnding
// callbacks are invoked in registration order, before any OnEnd callback is
// invoked.
//
// The span passed to [OnEndingSpanProcessor.OnEnding] must not be retained or
// mutated after the method returns.
type OnEndingSpanProcessor interface {
	// OnEnding is called during Span.End, after the end timestamp is set and
	// while the span is still mutable. Modifications to the span made here
	// will be visible in subsequent OnEnd callbacks.
	//
	// ending is the mutable span. It must not be retained or mutated after
	// the method returns.
	//
	// original is the exact [trace.Span] instance returned by
	// [go.opentelemetry.io/otel/trace.Tracer.Start] and passed to
	// [sdktrace.SpanProcessor.OnStart]. It has already ended, so mutations
	// through it are no-ops; its primary use is correlating state recorded
	// in OnStart with the ending span.
	OnEnding(ending sdktrace.ReadWriteSpan, original trace.Span)
}
