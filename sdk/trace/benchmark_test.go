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

package trace_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func BenchmarkStartEndSpan(b *testing.B) {
	traceBenchmark(b, "Benchmark StartEndSpan", func(b *testing.B, t trace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_4(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With 4 Attributes", func(b *testing.B, t trace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				label.Bool("key1", false),
				label.String("key2", "hello"),
				label.Uint64("key3", 123),
				label.Float64("key4", 123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_8(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With 8 Attributes", func(b *testing.B, t trace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				label.Bool("key1", false),
				label.String("key2", "hello"),
				label.Uint64("key3", 123),
				label.Float64("key4", 123.456),
				label.Bool("key21", false),
				label.String("key22", "hello"),
				label.Uint64("key23", 123),
				label.Float64("key24", 123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With all Attribute types", func(b *testing.B, t trace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				label.Bool("key1", false),
				label.String("key2", "hello"),
				label.Int64("key3", 123),
				label.Uint64("key4", 123),
				label.Int32("key5", 123),
				label.Uint32("key6", 123),
				label.Float64("key7", 123.456),
				label.Float32("key8", 123.456),
				label.Int("key9", 123),
				label.Uint("key10", 123),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all_2x(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With all Attributes types twice", func(b *testing.B, t trace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				label.Bool("key1", false),
				label.String("key2", "hello"),
				label.Int64("key3", 123),
				label.Uint64("key4", 123),
				label.Int32("key5", 123),
				label.Uint32("key6", 123),
				label.Float64("key7", 123.456),
				label.Float32("key8", 123.456),
				label.Int("key10", 123),
				label.Uint("key11", 123),
				label.Bool("key21", false),
				label.String("key22", "hello"),
				label.Int64("key23", 123),
				label.Uint64("key24", 123),
				label.Int32("key25", 123),
				label.Uint32("key26", 123),
				label.Float64("key27", 123.456),
				label.Float32("key28", 123.456),
				label.Int("key210", 123),
				label.Uint("key211", 123),
			)
			span.End()
		}
	})
}

func BenchmarkTraceID_DotString(b *testing.B) {
	t, _ := trace.TraceIDFromHex("0000000000000001000000000000002a")
	sc := trace.SpanContext{TraceID: t}

	want := "0000000000000001000000000000002a"
	for i := 0; i < b.N; i++ {
		if got := sc.TraceID.String(); got != want {
			b.Fatalf("got = %q want = %q", got, want)
		}
	}
}

func BenchmarkSpanID_DotString(b *testing.B) {
	sc := trace.SpanContext{SpanID: trace.SpanID{1}}
	want := "0100000000000000"
	for i := 0; i < b.N; i++ {
		if got := sc.SpanID.String(); got != want {
			b.Fatalf("got = %q want = %q", got, want)
		}
	}
}

func traceBenchmark(b *testing.B, name string, fn func(*testing.B, trace.Tracer)) {
	b.Run("AlwaysSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b, tracer(b, name, sdktrace.AlwaysSample()))
	})
	b.Run("NeverSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b, tracer(b, name, sdktrace.NeverSample()))
	})
}

func tracer(b *testing.B, name string, sampler sdktrace.Sampler) trace.Tracer {
	tp := sdktrace.NewTracerProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}))
	return tp.Tracer(name)
}
