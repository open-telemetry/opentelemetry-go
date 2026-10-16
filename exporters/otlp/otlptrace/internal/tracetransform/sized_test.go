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

func TestSpansRetainsAfterNextExport(t *testing.T) {
	res := resource.NewSchemaless(attribute.String("service.name", "test"))
	scope := instrumentation.Scope{Name: "test"}

	s1 := tracetest.SpanStub{
		Name:       "s1",
		Attributes: []attribute.KeyValue{attribute.String("k", "v1")},
		Events: []tracesdk.Event{
			{Name: "e1", Attributes: []attribute.KeyValue{attribute.String("ek", "ev1")}},
		},
		Resource:             res,
		InstrumentationScope: scope,
	}.Snapshot()
	out1 := Spans([]tracesdk.ReadOnlySpan{s1})
	require.Len(t, out1, 1)
	require.Len(t, out1[0].ScopeSpans, 1)
	require.Len(t, out1[0].ScopeSpans[0].Spans, 1)
	firstAttrs := out1[0].ScopeSpans[0].Spans[0].Attributes
	require.Len(t, firstAttrs, 1)
	assert.Equal(t, "v1", firstAttrs[0].Value.GetStringValue())

	s2 := tracetest.SpanStub{
		Name:       "s2",
		Attributes: []attribute.KeyValue{attribute.String("k", "v2")},
		Events: []tracesdk.Event{
			{Name: "e1", Attributes: []attribute.KeyValue{attribute.String("ek", "ev2")}},
		},
		Resource:             res,
		InstrumentationScope: scope,
	}.Snapshot()
	out2 := Spans([]tracesdk.ReadOnlySpan{s2})
	require.Len(t, out2, 1)

	assert.Equal(t, "v1", out1[0].ScopeSpans[0].Spans[0].Attributes[0].Value.GetStringValue())
	assert.Equal(t, "ev1", out1[0].ScopeSpans[0].Spans[0].Events[0].Attributes[0].Value.GetStringValue())
	assert.Equal(t, "v2", out2[0].ScopeSpans[0].Spans[0].Attributes[0].Value.GetStringValue())
}

var sizedBenchTime = time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

func sizedSparseBatch(n int) []tracesdk.ReadOnlySpan {
	res := resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.43.0",
		attribute.String("service.name", "bench-sparse"),
	)
	scope := instrumentation.Scope{Name: "bench/sparse", Version: "v0.0.1"}
	stubs := make([]tracetest.SpanStub, n)
	for i := range stubs {
		stubs[i] = tracetest.SpanStub{
			Name: "sparse-span",
			Attributes: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.Int64("http.status_code", 200),
			},
			Events: []tracesdk.Event{
				{
					Name: "cache-miss",
					Time: sizedBenchTime,
					Attributes: []attribute.KeyValue{
						attribute.String("cache.key", "k"),
					},
				},
			},
			Links: []tracesdk.Link{
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID: trace.TraceID{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
						SpanID:  trace.SpanID{byte(i), 1, 2, 3, 4, 5, 6, 7},
					}),
				},
			},
			Resource:             res,
			InstrumentationScope: scope,
		}
	}
	return tracetest.SpanStubs(stubs).Snapshots()
}

func sizedDenseBatch(n int) []tracesdk.ReadOnlySpan {
	res := resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.43.0",
		attribute.String("service.name", "bench-dense"),
		attribute.Int64("service.instance.id", 7),
		attribute.StringSlice("service.tags", []string{"a", "b", "c"}),
	)
	scope := instrumentation.Scope{
		Name:    "bench/dense",
		Version: "v0.1.0",
		Attributes: attribute.NewSet(
			attribute.String("scope.attr", "d1"),
		),
	}
	stubs := make([]tracetest.SpanStub, n)
	for i := range stubs {
		stubs[i] = tracetest.SpanStub{
			Name: "dense-span",
			Attributes: []attribute.KeyValue{
				attribute.String("http.method", "POST"),
				attribute.String("http.route", "/api/v1/items"),
				attribute.Int64("http.status_code", 201),
				attribute.Bool("error", false),
				attribute.Float64("duration.ms", 12.5),
				attribute.StringSlice("tags", []string{"t1", "t2"}),
				attribute.Int("retry.count", 2),
				attribute.String("user.id", "u-123"),
			},
			Events: []tracesdk.Event{
				{
					Name: "db.query",
					Time: sizedBenchTime,
					Attributes: []attribute.KeyValue{
						attribute.String("db.system", "postgres"),
						attribute.String("db.statement", "SELECT 1"),
					},
				},
			},
			Links: []tracesdk.Link{
				{
					SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
						TraceID: trace.TraceID{byte(i), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
						SpanID:  trace.SpanID{byte(i), 1, 2, 3, 4, 5, 6, 7},
					}),
					Attributes: []attribute.KeyValue{
						attribute.String("link.type", "parent"),
					},
				},
			},
			Resource:             res,
			InstrumentationScope: scope,
		}
	}
	return tracetest.SpanStubs(stubs).Snapshots()
}

func BenchmarkSpansSparseSized(b *testing.B) {
	spans := sizedSparseBatch(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Spans(spans)
	}
}

func BenchmarkSpansDenseSized(b *testing.B) {
	spans := sizedDenseBatch(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Spans(spans)
	}
}
