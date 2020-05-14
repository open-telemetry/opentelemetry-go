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

	"go.opentelemetry.io/otel/api/kv"

	apitrace "go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func BenchmarkStartEndSpan(b *testing.B) {
	traceBenchmark(b, "Benchmark StartEndSpan", func(b *testing.B, t apitrace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_4(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With 4 Attributes", func(b *testing.B, t apitrace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				kv.Key("key1").Bool(false),
				kv.Key("key2").String("hello"),
				kv.Key("key3").Uint64(123),
				kv.Key("key4").Float64(123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_8(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With 8 Attributes", func(b *testing.B, t apitrace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				kv.Key("key1").Bool(false),
				kv.Key("key2").String("hello"),
				kv.Key("key3").Uint64(123),
				kv.Key("key4").Float64(123.456),
				kv.Key("key21").Bool(false),
				kv.Key("key22").String("hello"),
				kv.Key("key23").Uint64(123),
				kv.Key("key24").Float64(123.456),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With all Attribute types", func(b *testing.B, t apitrace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				kv.Key("key1").Bool(false),
				kv.Key("key2").String("hello"),
				kv.Key("key3").Int64(123),
				kv.Key("key4").Uint64(123),
				kv.Key("key5").Int32(123),
				kv.Key("key6").Uint32(123),
				kv.Key("key7").Float64(123.456),
				kv.Key("key8").Float32(123.456),
				kv.Key("key9").Int(123),
				kv.Key("key10").Uint(123),
			)
			span.End()
		}
	})
}

func BenchmarkSpanWithAttributes_all_2x(b *testing.B) {
	traceBenchmark(b, "Benchmark Start With all Attributes types twice", func(b *testing.B, t apitrace.Tracer) {
		ctx := context.Background()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.SetAttributes(
				kv.Key("key1").Bool(false),
				kv.Key("key2").String("hello"),
				kv.Key("key3").Int64(123),
				kv.Key("key4").Uint64(123),
				kv.Key("key5").Int32(123),
				kv.Key("key6").Uint32(123),
				kv.Key("key7").Float64(123.456),
				kv.Key("key8").Float32(123.456),
				kv.Key("key10").Int(123),
				kv.Key("key11").Uint(123),
				kv.Key("key21").Bool(false),
				kv.Key("key22").String("hello"),
				kv.Key("key23").Int64(123),
				kv.Key("key24").Uint64(123),
				kv.Key("key25").Int32(123),
				kv.Key("key26").Uint32(123),
				kv.Key("key27").Float64(123.456),
				kv.Key("key28").Float32(123.456),
				kv.Key("key210").Int(123),
				kv.Key("key211").Uint(123),
			)
			span.End()
		}
	})
}

func BenchmarkTraceID_DotString(b *testing.B) {
	t, _ := apitrace.IDFromHex("0000000000000001000000000000002a")
	sc := apitrace.SpanContext{TraceID: t}

	want := "0000000000000001000000000000002a"
	for i := 0; i < b.N; i++ {
		if got := sc.TraceID.String(); got != want {
			b.Fatalf("got = %q want = %q", got, want)
		}
	}
}

func BenchmarkSpanID_DotString(b *testing.B) {
	sc := apitrace.SpanContext{SpanID: apitrace.SpanID{1}}
	want := "0100000000000000"
	for i := 0; i < b.N; i++ {
		if got := sc.SpanID.String(); got != want {
			b.Fatalf("got = %q want = %q", got, want)
		}
	}
}

func traceBenchmark(b *testing.B, name string, fn func(*testing.B, apitrace.Tracer)) {
	b.Run("AlwaysSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b, tracer(b, name, sdktrace.AlwaysSample()))
	})
	b.Run("NeverSample", func(b *testing.B) {
		b.ReportAllocs()
		fn(b, tracer(b, name, sdktrace.NeverSample()))
	})
}

func tracer(b *testing.B, name string, sampler sdktrace.Sampler) apitrace.Tracer {
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}))
	if err != nil {
		b.Fatalf("Failed to create trace provider for test %s\n", name)
	}
	return tp.Tracer(name)
}
