// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"

import (
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"

	"go.opentelemetry.io/otel/trace"
)

const (
	attributesPerSpan = 8
	// idsSizePerSpan pre-allocated space for span ids: 16(TraceID)+8(SpanID)+8(ParentSpanID)*2(extra space for links).
	idsSizePerSpan = (16 + 8 + 8) * 2
)

// Arena has slices for allocating small pb msg instead of one-by-one allocation.
type Arena struct {
	kvs chunkedStorage[commonpb.KeyValue]
	avs chunkedStorage[commonpb.AnyValue]

	avStrValues   []commonpb.AnyValue_StringValue
	avBoolValues  []commonpb.AnyValue_BoolValue
	avIntValues   []commonpb.AnyValue_IntValue
	avFloatValues []commonpb.AnyValue_DoubleValue

	ids []byte
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
		ids:           make([]byte, 0, size*idsSizePerSpan),
	}
}

func (a *Arena) allocTraceID(tid trace.TraceID) []byte {
	start := len(a.ids)
	a.ids = append(a.ids, tid[:]...)
	return a.ids[start:]
}

func (a *Arena) allocSpanID(sid trace.SpanID) []byte {
	start := len(a.ids)
	a.ids = append(a.ids, sid[:]...)
	return a.ids[start:]
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
