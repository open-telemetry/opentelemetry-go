// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptrace_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// benchNoopSyncClient isolates transform and arena cost: it consumes uploads
// synchronously without marshal or I/O.
type benchNoopSyncClient struct{}

var (
	_ otlptrace.SyncClient = (*benchNoopSyncClient)(nil)
	_ otlptrace.Client     = (*benchNoopSyncClient)(nil)
)

func (*benchNoopSyncClient) Start(context.Context) error { return nil }
func (*benchNoopSyncClient) Stop(context.Context) error  { return nil }

func (*benchNoopSyncClient) UploadTraces(context.Context, []*tracepb.ResourceSpans) error {
	return nil
}

func (*benchNoopSyncClient) UploadTracesSync(context.Context, []*tracepb.ResourceSpans) error {
	return nil
}

// benchHeapClient is the heap-path baseline: existing Exporter with a noop
// UploadTraces.
type benchHeapClient struct{}

var _ otlptrace.Client = (*benchHeapClient)(nil)

func (*benchHeapClient) Start(context.Context) error { return nil }
func (*benchHeapClient) Stop(context.Context) error  { return nil }

func (*benchHeapClient) UploadTraces(context.Context, []*tracepb.ResourceSpans) error {
	return nil
}

var benchTime = time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

func benchLinkContext(id byte) trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{id, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		SpanID:  trace.SpanID{id, 1, 2, 3, 4, 5, 6, 7},
	})
}

func benchSparseBatch(n int) []sdktrace.ReadOnlySpan {
	res := resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.43.0",
		attribute.String("service.name", "bench-sparse"),
		attribute.String("deployment.environment", "test"),
	)
	scope := instrumentation.Scope{
		Name:       "bench/sparse",
		Version:    "v0.0.1",
		SchemaURL:  "https://opentelemetry.io/schemas/1.43.0",
		Attributes: attribute.NewSet(attribute.String("scope.attr", "s1")),
	}
	stubs := make([]tracetest.SpanStub, n)
	for i := range stubs {
		stubs[i] = tracetest.SpanStub{
			Name: "sparse-span",
			Attributes: []attribute.KeyValue{
				attribute.String("http.method", "GET"),
				attribute.Int64("http.status_code", 200),
			},
			Events: []sdktrace.Event{
				{
					Name: "cache-miss",
					Time: benchTime,
					Attributes: []attribute.KeyValue{
						attribute.String("cache.key", "k"),
					},
				},
			},
			Links: []sdktrace.Link{
				{
					SpanContext: benchLinkContext(byte(i)),
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

func benchDenseBatch(n int) []sdktrace.ReadOnlySpan {
	res := resource.NewWithAttributes(
		"https://opentelemetry.io/schemas/1.43.0",
		attribute.String("service.name", "bench-dense"),
		attribute.String("deployment.environment", "test"),
		attribute.Int64("service.instance.id", 7),
		attribute.Bool("service.tls", true),
		attribute.StringSlice("service.tags", []string{"a", "b", "c"}),
	)
	scope := instrumentation.Scope{
		Name:      "bench/dense",
		Version:   "v0.1.0",
		SchemaURL: "https://opentelemetry.io/schemas/1.43.0",
		Attributes: attribute.NewSet(
			attribute.String("scope.attr", "d1"),
			attribute.Int64("scope.revision", 3),
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
			Events: []sdktrace.Event{
				{
					Name: "db.query",
					Time: benchTime,
					Attributes: []attribute.KeyValue{
						attribute.String("db.system", "postgres"),
						attribute.String("db.statement", "SELECT 1"),
						attribute.Int64("db.rows", 42),
						attribute.Bool("db.cached", true),
					},
				},
				{
					Name: "cache.store",
					Time: benchTime,
					Attributes: []attribute.KeyValue{
						attribute.String("cache.key", "k2"),
						attribute.Int64("cache.ttl_ms", 60000),
					},
				},
			},
			Links: []sdktrace.Link{
				{
					SpanContext: benchLinkContext(byte(i)),
					Attributes: []attribute.KeyValue{
						attribute.String("link.type", "parent"),
						attribute.Int64("link.weight", 3),
						attribute.Bool("link.sampled", true),
					},
				},
				{
					SpanContext: benchLinkContext(byte(i + 1)),
					Attributes: []attribute.KeyValue{
						attribute.String("link.type", "child"),
					},
				},
			},
			Resource:             res,
			InstrumentationScope: scope,
		}
	}
	return tracetest.SpanStubs(stubs).Snapshots()
}

// runSyncBenchmark exports the same batch in a sequential loop so the pool
// reuse (retain) or realloc (discard) behavior reaches steady state. A single
// warmup export runs before the timer so retained numbers reflect reuse, not
// the initial arena allocation. Sequential (not parallel) keeps pool ownership
// deterministic.
func runSyncBenchmark(b *testing.B, exp *otlptrace.SyncExporter, spans []sdktrace.ReadOnlySpan) {
	b.Helper()
	ctx := b.Context()
	require.NoError(b, exp.ExportSpans(ctx, spans))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := exp.ExportSpans(ctx, spans); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSyncExporterSparseRetained(b *testing.B) {
	spans := benchSparseBatch(100)
	exp, err := otlptrace.NewSync(
		b.Context(),
		&benchNoopSyncClient{},
		otlptrace.WithInitialBatchSize(128),
		otlptrace.WithMaxRetainedBatchSize(512),
	)
	require.NoError(b, err)
	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
	b.Cleanup(func() { _ = exp.Shutdown(context.Background()) })
	runSyncBenchmark(b, exp, spans)
}

func BenchmarkSyncExporterDenseRetained(b *testing.B) {
	spans := benchDenseBatch(100)
	exp, err := otlptrace.NewSync(
		b.Context(),
		&benchNoopSyncClient{},
		otlptrace.WithInitialBatchSize(128),
		otlptrace.WithMaxRetainedBatchSize(512),
	)
	require.NoError(b, err)
	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
	b.Cleanup(func() { _ = exp.Shutdown(context.Background()) })
	runSyncBenchmark(b, exp, spans)
}

func BenchmarkSyncExporterDenseDiscarded(b *testing.B) {
	spans := benchDenseBatch(100)
	exp, err := otlptrace.NewSync(
		b.Context(),
		&benchNoopSyncClient{},
		otlptrace.WithInitialBatchSize(128),
		otlptrace.WithMaxRetainedBatchSize(1),
	)
	require.NoError(b, err)
	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
	b.Cleanup(func() { _ = exp.Shutdown(context.Background()) })
	runSyncBenchmark(b, exp, spans)
}

func BenchmarkExporterDenseHeap(b *testing.B) {
	spans := benchDenseBatch(100)
	exp, err := otlptrace.New(b.Context(), &benchHeapClient{})
	require.NoError(b, err)
	//nolint:usetesting // required to avoid getting a canceled context at cleanup.
	b.Cleanup(func() { _ = exp.Shutdown(context.Background()) })
	ctx := b.Context()
	require.NoError(b, exp.ExportSpans(ctx, spans))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := exp.ExportSpans(ctx, spans); err != nil {
			b.Fatal(err)
		}
	}
}
