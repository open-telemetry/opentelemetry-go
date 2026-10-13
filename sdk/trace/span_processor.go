// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"sync"
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

type spanProcessorState struct {
	sp    SpanProcessor
	state sync.Once

	// admissionMu protects the fields below. They track the operations
	// (OnStart, OnEnd, ForceFlush) currently in flight on sp so that sp's
	// Shutdown is not invoked until they complete and no new ones start.
	admissionMu sync.Mutex
	active      int
	stopped     bool
	drained     chan struct{}
}

func newSpanProcessorState(sp SpanProcessor) *spanProcessorState {
	return &spanProcessorState{sp: sp}
}

type spanProcessorStates []*spanProcessorState

// begin reports whether a new operation on sp is admitted. If it returns true,
// the caller must call end when the operation completes.
func (s *spanProcessorState) begin() bool {
	s.admissionMu.Lock()
	defer s.admissionMu.Unlock()
	if s.stopped {
		return false
	}
	s.active++
	return true
}

func (s *spanProcessorState) end() {
	s.admissionMu.Lock()
	defer s.admissionMu.Unlock()
	s.active--
	if s.active == 0 && s.stopped {
		close(s.drained)
	}
}

// stop closes admission of new operations. It does not wait for admitted
// operations to complete; use wait for that. It is safe to call repeatedly.
func (s *spanProcessorState) stop() {
	s.admissionMu.Lock()
	defer s.admissionMu.Unlock()
	if s.stopped {
		return
	}
	s.stopped = true
	s.drained = make(chan struct{})
	if s.active == 0 {
		close(s.drained)
	}
}

// wait blocks until all operations admitted before stop have completed,
// returning ctx.Err() if ctx is done first. stop must have been called.
func (s *spanProcessorState) wait(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-s.drained:
	case <-ctx.Done():
		// Prefer a completed drain when it races with cancellation.
		select {
		case <-s.drained:
		default:
			return ctx.Err()
		}
	}
	return nil
}

// onStart, onEnd, and forceFlush invoke the processor method only if
// admission is open. The admission count is released even if the processor
// panics so that shutdown cannot be blocked forever.
func (s *spanProcessorState) onStart(ctx context.Context, rw ReadWriteSpan) {
	if !s.begin() {
		return
	}
	defer s.end()
	s.sp.OnStart(ctx, rw)
}

func (s *spanProcessorState) onEnd(snap ReadOnlySpan) {
	if !s.begin() {
		return
	}
	defer s.end()
	s.sp.OnEnd(snap)
}

func (s *spanProcessorState) forceFlush(ctx context.Context) error {
	if !s.begin() {
		return nil
	}
	defer s.end()
	return s.sp.ForceFlush(ctx)
}
