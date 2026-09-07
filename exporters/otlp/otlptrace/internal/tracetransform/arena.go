// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

// Arena provides preallocated storage for protobuf messages to reduce heap
// allocations during OTLP transformation. It is intended to be reused across
// exports when the caller guarantees synchronous consumption (see SyncClient).
type Arena struct {
	kvs chunkedStorage[commonpb.KeyValue]
	avs chunkedStorage[commonpb.AnyValue]

	avStrValues   []commonpb.AnyValue_StringValue
	avBoolValues  []commonpb.AnyValue_BoolValue
	avIntValues   []commonpb.AnyValue_IntValue
	avFloatValues []commonpb.AnyValue_DoubleValue
}

const defaultAttributesPerSpan = 8

// NewArena creates a new Arena sized for the given number of spans.
func NewArena(size int) *Arena {
	size = max(1, size)
	chunkSize := size * defaultAttributesPerSpan
	return &Arena{
		kvs: chunkedStorage[commonpb.KeyValue]{
			chunkSize: chunkSize,
			resetFn:   func(m *commonpb.KeyValue) { m.Reset() },
		},
		avs: chunkedStorage[commonpb.AnyValue]{
			chunkSize: chunkSize,
			resetFn:   func(m *commonpb.AnyValue) { m.Reset() },
		},
		avStrValues:   make([]commonpb.AnyValue_StringValue, 0, chunkSize),
		avBoolValues:  make([]commonpb.AnyValue_BoolValue, 0, size),
		avIntValues:   make([]commonpb.AnyValue_IntValue, 0, size),
		avFloatValues: make([]commonpb.AnyValue_DoubleValue, 0, size),
	}
}

// Reset clears the arena for reuse. It must only be called after the
// previously allocated data is no longer referenced.
func (a *Arena) Reset() {
	a.kvs.reset()
	a.avs.reset()
	clear(a.avStrValues)
	a.avStrValues = a.avStrValues[:0]
	a.avBoolValues = a.avBoolValues[:0]
	a.avIntValues = a.avIntValues[:0]
	a.avFloatValues = a.avFloatValues[:0]
}

// Cap returns the total capacity of the arena in number of key-value slots.
func (a *Arena) Cap() int {
	return len(a.kvs.chunks) * a.kvs.chunkSize
}

// Exceeds reports whether the arena has grown beyond the given span limit.
// It is used to implement bounded retention: large arenas are discarded
// instead of being returned to a pool.
func (a *Arena) Exceeds(maxSpans int) bool {
	if maxSpans <= 0 {
		return false
	}
	// Approximate capacity in spans.
	maxCap := maxSpans * defaultAttributesPerSpan
	// Cap rounds up to whole chunks, so a batch using exactly maxCap slots
	// plus resource/scope overhead needs one more chunk. Allow that slack
	// so a default-size batch is retained instead of discarded.
	if a.Cap() > maxCap+a.kvs.chunkSize {
		return true
	}
	// Slices grow with Go append overshoot (up to ~25%) plus the same
	// chunk rounding. Each span can hold up to 8 scalars of one type,
	// so compare against maxCap rather than maxSpans.
	slack := a.kvs.chunkSize + maxCap/4
	if cap(a.avStrValues) > maxCap+slack {
		return true
	}
	if cap(a.avBoolValues) > maxCap+slack {
		return true
	}
	if cap(a.avIntValues) > maxCap+slack {
		return true
	}
	if cap(a.avFloatValues) > maxCap+slack {
		return true
	}
	return false
}

type chunkedStorage[T any] struct {
	chunkSize int
	chunks    [][]T
	idx       int
	resetFn   func(*T)
}

func (s *chunkedStorage[T]) alloc() *T {
	chunk := s.idx / s.chunkSize
	pos := s.idx % s.chunkSize
	if chunk >= len(s.chunks) {
		s.chunks = append(s.chunks, make([]T, s.chunkSize))
	}
	s.idx++
	return &s.chunks[chunk][pos]
}

func (s *chunkedStorage[T]) reset() {
	for i := range s.idx {
		chunk := i / s.chunkSize
		pos := i % s.chunkSize
		s.resetFn(&s.chunks[chunk][pos])
	}
	s.idx = 0
}
