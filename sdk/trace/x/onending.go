// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// OnEndingSpanProcessor is implemented by [sdktrace.SpanProcessor]
// implementations that want to mutate a span during its
// [go.opentelemetry.io/otel/trace.Span.End] call, after the end timestamp has
// been set but before the span becomes immutable.
//
// OnEndingSpanProcessor is an optional interface. The [sdktrace.TracerProvider]
// detects it by type assertion at registration time, so no configuration
// beyond registering a processor that implements this interface is required.
//
// The [OnEndingSpanProcessor.OnEnding] method is called synchronously within
// [go.opentelemetry.io/otel/trace.Span.End], therefore it must not block or
// panic. If multiple [sdktrace.SpanProcessor] implementations are registered,
// their OnEnding callbacks are invoked in registration order, before any
// OnEnd callback is invoked.
//
// OnEnding may be called concurrently for different spans and may overlap
// other [sdktrace.SpanProcessor] methods on the same processor, so
// implementations must be safe for concurrent use.
//
// OnEnding may be called concurrently with or after [sdktrace.SpanProcessor.Shutdown]
// during processor unregistration. Implementations must handle this gracefully.
//
// Note: OnEnding alone does not suppress other processors, change
// TraceFlags.Sampled, or drop a span. Tail-based filtering that needs those
// effects requires a custom buffering/export pipeline on top of this hook.
//
// The ending span passed to [OnEndingSpanProcessor.OnEnding] must not be
// retained or mutated after the method returns.
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
