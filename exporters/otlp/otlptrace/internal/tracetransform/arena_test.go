// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestArena(t *testing.T) {
	t.Run("alloc traceID and spanID for size", func(t *testing.T) {
		const arenaSize = 16
		for n := range arenaSize * 3 {
			a := NewArena(arenaSize)

			spans := generateSpans(n)
			gotTIDs := make([][]byte, 0, len(spans))
			gotSIDs := make([][]byte, 0, len(spans))

			for _, s := range spans {
				sc := s.SpanContext()
				tid := a.allocTraceID(sc.TraceID())
				gotTIDs = append(gotTIDs, tid)
				sid := a.allocSpanID(sc.SpanID())
				gotSIDs = append(gotSIDs, sid)
			}

			for i, tid := range gotTIDs {
				sc := spans[i].SpanContext()
				sid := gotSIDs[i]
				var (
					traceID trace.TraceID
					spanID  trace.SpanID
				)

				assert.Len(t, tid, len(traceID))
				copy(traceID[:], tid)
				assert.Len(t, sid, len(spanID))
				copy(spanID[:], sid)

				assert.Equal(t, sc.TraceID(), traceID)
				assert.Equal(t, sc.SpanID(), spanID)
			}
		}
	})
}

func generateSpans(n int) []tracesdk.ReadOnlySpan {
	spans := make([]tracesdk.ReadOnlySpan, n)
	for i := range n {
		spans[i] = generateSpanWithRandomIDs()
	}
	return spans
}

func generateSpanWithRandomIDs() tracesdk.ReadOnlySpan {
	tid := trace.TraceID{}
	sid := trace.SpanID{}

	binary.NativeEndian.PutUint64(tid[:8], rand.Uint64())
	binary.NativeEndian.PutUint64(tid[8:], rand.Uint64())
	binary.NativeEndian.PutUint64(sid[:], rand.Uint64())

	return tracetest.SpanStub{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: tid,
			SpanID:  sid,
		}),
	}.Snapshot()
}
