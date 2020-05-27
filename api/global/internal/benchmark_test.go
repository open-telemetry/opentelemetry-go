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

package internal_test

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/global/internal"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var Must = metric.Must

// benchFixture is copied from sdk/metric/benchmark_test.go.
// TODO refactor to share this code.
type benchFixture struct {
	accumulator *sdk.Accumulator
	meter       metric.Meter
	B           *testing.B
}

var _ metric.Provider = &benchFixture{}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B: b,
	}

	bf.accumulator = sdk.NewAccumulator(bf)
	bf.meter = metric.WrapMeterImpl(bf.accumulator, "test")
	return bf
}

func (*benchFixture) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case metric.CounterKind:
		return sum.New()
	case metric.ValueRecorderKind:
		if strings.HasSuffix(descriptor.Name(), "minmaxsumcount") {
			return minmaxsumcount.New(descriptor)
		} else if strings.HasSuffix(descriptor.Name(), "ddsketch") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		} else if strings.HasSuffix(descriptor.Name(), "array") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		}
	}
	return nil
}

func (*benchFixture) Process(context.Context, export.Record) error {
	return nil
}

func (*benchFixture) CheckpointSet() export.CheckpointSet {
	return nil
}

func (*benchFixture) FinishedCollection() {
}

func (fix *benchFixture) Meter(name string) metric.Meter {
	return fix.meter
}

func BenchmarkGlobalInt64CounterAddNoSDK(b *testing.B) {
	internal.ResetForTest()
	ctx := context.Background()
	sdk := global.Meter("test")
	labs := []kv.KeyValue{kv.String("A", "B")}
	cnt := Must(sdk).NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkGlobalInt64CounterAddWithSDK(b *testing.B) {
	// Comapare with BenchmarkInt64CounterAdd() in ../../sdk/meter/benchmark_test.go
	ctx := context.Background()
	fix := newFixture(b)

	sdk := global.Meter("test")

	global.SetMeterProvider(fix)

	labs := []kv.KeyValue{kv.String("A", "B")}
	cnt := Must(sdk).NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkStartEndSpan(b *testing.B) {
	// Comapare with BenchmarkStartEndSpan() in ../../sdk/trace/benchmark_test.go
	traceBenchmark(b, func(b *testing.B) {
		t := global.Tracer("Benchmark StartEndSpan")
		ctx := context.Background()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, span := t.Start(ctx, "/foo")
			span.End()
		}
	})
}

func traceBenchmark(b *testing.B, fn func(*testing.B)) {
	internal.ResetForTest()
	b.Run("No SDK", func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
	b.Run("Default SDK (AlwaysSample)", func(b *testing.B) {
		b.ReportAllocs()
		global.SetTraceProvider(traceProvider(b, sdktrace.AlwaysSample()))
		fn(b)
	})
	b.Run("Default SDK (NeverSample)", func(b *testing.B) {
		b.ReportAllocs()
		global.SetTraceProvider(traceProvider(b, sdktrace.NeverSample()))
		fn(b)
	})
}

func traceProvider(b *testing.B, sampler sdktrace.Sampler) trace.Provider {
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}))
	if err != nil {
		b.Fatalf("Failed to create trace provider with sampler: %v", err)
	}
	return tp
}
