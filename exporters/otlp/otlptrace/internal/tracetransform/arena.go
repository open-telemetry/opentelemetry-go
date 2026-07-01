// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

const attributesPerSpan = 8

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
	size = max(1, size)
	return &Arena{
		kvs:           chunkedStorage[commonpb.KeyValue]{chunkSize: size * attributesPerSpan},
		avs:           chunkedStorage[commonpb.AnyValue]{chunkSize: size * attributesPerSpan},
		avStrValues:   make([]commonpb.AnyValue_StringValue, 0, size*attributesPerSpan),
		avBoolValues:  make([]commonpb.AnyValue_BoolValue, 0, size),
		avIntValues:   make([]commonpb.AnyValue_IntValue, 0, size),
		avFloatValues: make([]commonpb.AnyValue_DoubleValue, 0, size),
	}
}

type chunkedStorage[T any] struct {
	chunkSize int
	chunks    [][]T
	idx       int
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
