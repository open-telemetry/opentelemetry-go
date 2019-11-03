package propagation

import (
	"context"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func BenchmarkInject(b *testing.B) {
	var t HTTPTraceContextPropagator

	injectSubBenchmarks(b, func(ctx context.Context, b *testing.B) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			t.Inject(ctx, req.Header)
		}
	})
}

func injectSubBenchmarks(b *testing.B, fn func(context.Context, *testing.B)) {
	b.Run("SampledSpanContext", func(b *testing.B) {
		var id uint64
		spanID, _ := core.SpanIDFromHex("00f067aa0ba902b7")
		traceID, _ := core.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")

		mockTracer := &mocktrace.MockTracer{
			Sampled:     false,
			StartSpanID: &id,
		}
		b.ReportAllocs()
		sc := core.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: core.TraceFlagsSampled,
		}
		ctx := context.Background()
		ctx, _ = mockTracer.Start(ctx, "inject", trace.ChildOf(sc))
		fn(ctx, b)
	})

	b.Run("WithoutSpanContext", func(b *testing.B) {
		b.ReportAllocs()
		ctx := context.Background()
		fn(ctx, b)
	})
}

func BenchmarkExtract(b *testing.B) {
	extractSubBenchmarks(b, func(b *testing.B, req *http.Request) {
		var propagator HTTPTraceContextPropagator
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			propagator.Extract(ctx, req.Header)
		}
	})
}

func extractSubBenchmarks(b *testing.B, fn func(*testing.B, *http.Request)) {
	b.Run("Sampled", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
		b.ReportAllocs()

		fn(b, req)
	})

	b.Run("BogusVersion", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("traceparent", "qw-00000000000000000000000000000000-0000000000000000-01")
		b.ReportAllocs()
		fn(b, req)
	})

	b.Run("FutureAdditionalData", func(b *testing.B) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		req.Header.Set("traceparent", "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09-XYZxsf09")
		b.ReportAllocs()
		fn(b, req)
	})
}
