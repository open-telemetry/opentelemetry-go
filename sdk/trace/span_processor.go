// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

// SpanProcessor is a processing pipeline for spans in the trace signal.
// SpanProcessors registered with a TracerProvider and are called at the start
// and end of a Span's lifecycle, and are called in the order they are
// registered.
type SpanProcessor interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// OnStart is called when a span is started. It is called synchronously
	// and should not block.
	OnStart(parent context.Context, s ReadWriteSpan)
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// OnEnd is called when span is finished. It is called synchronously and
	// hence not block.
	OnEnd(s ReadOnlySpan)
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// Shutdown is called when the SDK shuts down. Any cleanup or release of
	// resources held by the processor should be done in this call.
	//
	// Calls to OnStart, OnEnd, or ForceFlush after this has been called
	// should be ignored.
	//
	// All timeouts and cancellations contained in ctx must be honored, this
	// should not block indefinitely.
	Shutdown(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.

	// ForceFlush exports all ended spans to the configured Exporter that have not yet
	// been exported.  It should only be called when absolutely necessary, such as when
	// using a FaaS provider that may suspend the process after an invocation, but before
	// the Processor can export the completed spans.
	ForceFlush(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and
	// must never be done outside of a new major release.
}

// onEndingSpanProcessor is the unexported interface that mirrors
// [go.opentelemetry.io/otel/sdk/trace/x.OnEndingSpanProcessor]. It is
// detected by structural type assertion at registration time.
type onEndingSpanProcessor interface {
	OnEnding(ReadWriteSpan, trace.Span)
}

type spanProcessorState struct {
	sp       SpanProcessor
	onEnding onEndingSpanProcessor
	state    sync.Once
}

func newSpanProcessorState(sp SpanProcessor) *spanProcessorState {
	s := &spanProcessorState{sp: sp}
	if oesp, ok := sp.(onEndingSpanProcessor); ok {
		s.onEnding = oesp
	}
	return s
}

type spanProcessorStates []*spanProcessorState

// hasOnEnding reports whether any processor in the slice implements the
// OnEnding callback.
func (s spanProcessorStates) hasOnEnding() bool {
	for _, sps := range s {
		if sps.onEnding != nil {
			return true
		}
	}
	return false
}
