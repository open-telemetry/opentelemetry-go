// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	apitrace "go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func BenchmarkStartEndSpan(b *testing.B) {
	t := getTracer(b, "Benchmark StartEndSpan")

	traceBenchmark(b, func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_4(b *testing.B) {
	t := getTracer(b, "Benchmark Start With 4 Attributes")

	traceBenchmark(b, func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				key.New("key1").Bool(false),
				key.New("key2").String("hello"),
				key.New("key3").Uint64(123),
				key.New("key4").Float64(123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_8(b *testing.B) {
	t := getTracer(b, "Benchmark Start With 8 Attributes")

	traceBenchmark(b, func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				key.New("key1").Bool(false),
				key.New("key2").String("hello"),
				key.New("key3").Uint64(123),
				key.New("key4").Float64(123.456),
				key.New("key21").Bool(false),
				key.New("key22").String("hello"),
				key.New("key23").Uint64(123),
				key.New("key24").Float64(123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all(b *testing.B) {
	t := getTracer(b, "Benchmark Start With all Attribute types")

	traceBenchmark(b, func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				key.New("key1").Bool(false),
				key.New("key2").String("hello"),
				key.New("key3").Int64(123),
				key.New("key4").Uint64(123),
				key.New("key5").Int32(123),
				key.New("key6").Uint32(123),
				key.New("key7").Float64(123.456),
				key.New("key8").Float32(123.456),
				key.New("key9").Int(123),
				key.New("key10").Uint(123),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all_2x(b *testing.B) {
	t := getTracer(b, "Benchmark Start With all Attributes types twice")
	traceBenchmark(b, func(b *testing.B) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				key.New("key1").Bool(false),
				key.New("key2").String("hello"),
				key.New("key3").Int64(123),
				key.New("key4").Uint64(123),
				key.New("key5").Int32(123),
				key.New("key6").Uint32(123),
				key.New("key7").Float64(123.456),
				key.New("key8").Float32(123.456),
				key.New("key10").Int(123),
				key.New("key11").Uint(123),
				key.New("key21").Bool(false),
				key.New("key22").String("hello"),
				key.New("key23").Int64(123),
				key.New("key24").Uint64(123),
				key.New("key25").Int32(123),
				key.New("key26").Uint32(123),
				key.New("key27").Float64(123.456),
				key.New("key28").Float32(123.456),
				key.New("key210").Int(123),
				key.New("key211").Uint(123),
			)
			span.End()
		}
	})
}

func BenchmarkTraceID_DotString(b *testing.B) {
	traceBenchmark(b, func(b *testing.B) {
		t, _ := core.TraceIDFromHex("0000000000000001000000000000002a")
		sc := core.SpanContext{TraceID: t}

		want := "0000000000000001000000000000002a"
		for i := 0; i < b.N; i++ {
			if got := sc.TraceIDString(); got != want {
				b.Fatalf("got = %q want = %q", got, want)
			}
		}
	})
}

func BenchmarkSpanID_DotString(b *testing.B) {
	traceBenchmark(b, func(b *testing.B) {
		sc := core.SpanContext{SpanID: core.SpanID{1}}
		want := "1000000000000000"
		for i := 0; i < b.N; i++ {
			if got := sc.SpanIDString(); got != want {
				b.Fatalf("got = %q want = %q", got, want)
			}
		}
	})
}

func traceBenchmark(b *testing.B, fn func(*testing.B)) {
	b.Run("AlwaysSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
	b.Run("NeverSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
}

func getTracer(b *testing.B, name string) apitrace.Tracer {
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(testConfig))
	if err != nil {
		b.Fatalf("Failed to create trace provider for test %s\n", name)
	}
	return tp.GetTracer(name)
}
