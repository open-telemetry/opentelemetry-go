// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

// Arena has slices for allocating small pb msg instead of one-by-one allocation.
type Arena struct {
	kvs chunkedStorage[commonpb.KeyValue]
	avs chunkedStorage[commonpb.AnyValue]

	avStrValues   []commonpb.AnyValue_StringValue
	avBoolValues  []commonpb.AnyValue_BoolValue
	avIntValues   []commonpb.AnyValue_IntValue
	avFloatValues []commonpb.AnyValue_DoubleValue
}

// NewArena creates new Arena for spans transformation.
func NewArena(size int) *Arena {
	return &Arena{
		kvs: chunkedStorage[commonpb.KeyValue]{chunkSize: size, resetFn: func(m *commonpb.KeyValue) {
			m.Reset()
		}},
		avs: chunkedStorage[commonpb.AnyValue]{chunkSize: size, resetFn: func(m *commonpb.AnyValue) {
			m.Reset()
		}},
		avStrValues:   make([]commonpb.AnyValue_StringValue, 0, size),
		avBoolValues:  make([]commonpb.AnyValue_BoolValue, 0, size),
		avIntValues:   make([]commonpb.AnyValue_IntValue, 0, size),
		avFloatValues: make([]commonpb.AnyValue_DoubleValue, 0, size),
	}
}

// Reset resets Arena should be called after UploadTraces is done and structs allocated with arena won't be used anymore.
func (a *Arena) Reset() {
	a.kvs.reset()
	a.avs.reset()
	// strings allocated on heap so we clear avStrValues for GC to collect it properly
	clear(a.avStrValues)
	a.avStrValues = a.avStrValues[:0]
	a.avBoolValues = a.avBoolValues[:0]
	a.avIntValues = a.avIntValues[:0]
	a.avFloatValues = a.avFloatValues[:0]
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
