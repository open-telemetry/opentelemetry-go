// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package propagation_test

import (
	"context"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func BenchmarkInject(b *testing.B) {
	var t propagation.TraceContext

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
		spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
		traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")

		mockTracer := oteltest.DefaultTracer()
		b.ReportAllocs()
		sc := trace.SpanContext{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
		}
		ctx := trace.ContextWithRemoteSpanContext(context.Background(), sc)
		ctx, _ = mockTracer.Start(ctx, "inject")
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
		var propagator propagation.TraceContext
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
