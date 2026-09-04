// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func TestSpansWithArenaMatchesSpans(t *testing.T) {
	res := resource.NewSchemaless(attribute.String("service.name", "test"))
	span := tracetest.SpanStub{
		Name:        "test-span",
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1, 2, 3}, SpanID: trace.SpanID{4, 5, 6}}),
		Parent:      trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1, 2, 3}, SpanID: trace.SpanID{7, 8, 9}}),
		SpanKind:    trace.SpanKindServer,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Second),
		Attributes: []attribute.KeyValue{
			attribute.String("k", "v"),
			attribute.Int("i", 42),
			attribute.Bool("b", true),
			attribute.Float64("f", 3.14),
			attribute.StringSlice("ss", []string{"a", "b"}),
			attribute.Int64Slice("is", []int64{1, 2}),
		},
		Events: []tracesdk.Event{
			{Name: "event1", Attributes: []attribute.KeyValue{attribute.String("ek", "ev")}},
		},
		Links: []tracesdk.Link{
			{SpanContext: trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{9, 9, 9}, SpanID: trace.SpanID{8, 8, 8}}), Attributes: []attribute.KeyValue{attribute.String("lk", "lv")}},
		},
		Resource:             res,
		InstrumentationScope: instrumentation.Scope{Name: "test", Version: "1.0"},
	}.Snapshot()

	heapSpans := Spans([]tracesdk.ReadOnlySpan{span})
	arena := NewArena(4)
	arenaSpans := SpansWithArena([]tracesdk.ReadOnlySpan{span}, arena)

	require.Len(t, heapSpans, 1)
	require.Len(t, arenaSpans, 1)
	require.Len(t, heapSpans[0].ScopeSpans, 1)
	require.Len(t, arenaSpans[0].ScopeSpans, 1)

	hs := heapSpans[0].ScopeSpans[0].Spans[0]
	as := arenaSpans[0].ScopeSpans[0].Spans[0]

	assert.Equal(t, hs.Name, as.Name)
	assert.Equal(t, hs.TraceId, as.TraceId)
	assert.Equal(t, hs.SpanId, as.SpanId)
	assert.Equal(t, hs.Kind, as.Kind)
	assert.Equal(t, hs.Status, as.Status)
	assert.ElementsMatch(t, hs.Attributes, as.Attributes)
	assert.Equal(t, hs.Events, as.Events)
	assert.Equal(t, hs.Links, as.Links)
}

func TestArenaReuse(t *testing.T) {
	res := resource.NewSchemaless(attribute.String("k", "v"))
	arena := NewArena(2)

	s1 := tracetest.SpanStub{
		Name:                 "s1",
		SpanContext:          trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{1}}),
		Resource:             res,
		InstrumentationScope: instrumentation.Scope{Name: "a"},
	}.Snapshot()
	r1 := SpansWithArena([]tracesdk.ReadOnlySpan{s1}, arena)
	require.Len(t, r1, 1)
	name1 := r1[0].ScopeSpans[0].Spans[0].Name

	arena.Reset()

	s2 := tracetest.SpanStub{
		Name:                 "s2",
		SpanContext:          trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{2}, SpanID: trace.SpanID{2}}),
		Resource:             res,
		InstrumentationScope: instrumentation.Scope{Name: "a"},
	}.Snapshot()
	r2 := SpansWithArena([]tracesdk.ReadOnlySpan{s2}, arena)
	require.Len(t, r2, 1)
	assert.Equal(t, "s2", r2[0].ScopeSpans[0].Spans[0].Name)
	// Ensure first result would be corrupted if retained (demonstrates why sync is needed)
	// After Reset, r1's underlying arena memory is reused, so name1 should not be trusted after Reset.
	// We just verify r2 is correct.
	_ = name1
}

func TestArenaExceeds(t *testing.T) {
	arena := NewArena(2)
	assert.False(t, arena.Exceeds(10))
	// Grow arena by using many spans
	res := resource.NewSchemaless()
	spans := make([]tracesdk.ReadOnlySpan, 100)
	for i := range spans {
		spans[i] = tracetest.SpanStub{
			Name:       "s",
			Attributes: []attribute.KeyValue{attribute.String("k", "v"), attribute.Int("i", i)},
			Resource:   res,
		}.Snapshot()
	}
	_ = SpansWithArena(spans, arena)
	// After large batch, cap should exceed small limit
	assert.True(t, arena.Exceeds(5))
	assert.False(t, arena.Exceeds(1000))
}
