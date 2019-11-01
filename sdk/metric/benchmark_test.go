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

package metric_test

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/sdk/export"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/maxsumcount"
)

type benchFixture struct {
	sdk *sdk.SDK
	B   *testing.B
}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B: b,
	}
	bf.sdk = sdk.New(bf)
	return bf
}

func (bf *benchFixture) AggregatorFor(rec export.MetricRecord) export.MetricAggregator {
	switch rec.Descriptor().MetricKind() {
	case export.CounterMetricKind:
		return counter.New()
	case export.GaugeMetricKind:
		return gauge.New()
	case export.MeasureMetricKind:
		if strings.HasSuffix(rec.Descriptor().Name(), "maxsumcount") {
			return maxsumcount.New()
		} else if strings.HasSuffix(rec.Descriptor().Name(), "ddsketch") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), rec.Descriptor())
		}
	}
	return nil
}

func (bf *benchFixture) Export(ctx context.Context, rec export.MetricRecord, agg export.MetricAggregator) {
}

func makeLabels(n int) []core.KeyValue {
	used := map[string]bool{}
	l := make([]core.KeyValue, n)
	for i := 0; i < n; i++ {
		var k string
		for {
			k = fmt.Sprint("k", rand.Intn(1000000000))
			if !used[k] {
				used[k] = true
				break
			}
		}
		l[i] = key.New(k).String(fmt.Sprint("v", rand.Intn(1000000000)))
	}
	return l
}

func benchmarkLabels(b *testing.B, n int) {
	fix := newFixture(b)
	labs := makeLabels(n)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fix.sdk.Labels(labs...)
	}
}

func BenchmarkLabels_1(b *testing.B) {
	benchmarkLabels(b, 1)
}

func BenchmarkLabels_2(b *testing.B) {
	benchmarkLabels(b, 2)
}

func BenchmarkLabels_4(b *testing.B) {
	benchmarkLabels(b, 4)
}

func BenchmarkLabels_8(b *testing.B) {
	benchmarkLabels(b, 8)
}

func BenchmarkLabels_16(b *testing.B) {
	benchmarkLabels(b, 16)
}

// Note: performance does not depend on label set size for the
// benchmarks below.

func BenchmarkInt64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs)
	}
}

func BenchmarkInt64CounterAcquireHandle(b *testing.B) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := cnt.AcquireHandle(labs)
		handle.Release()
	}
}

func BenchmarkInt64CounterHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewInt64Counter("int64.counter")
	handle := cnt.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1)
	}
}

func BenchmarkFloat64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewFloat64Counter("float64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1.1, labs)
	}
}

func BenchmarkFloat64CounterAcquireHandle(b *testing.B) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewFloat64Counter("float64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := cnt.AcquireHandle(labs)
		handle.Release()
	}
}

func BenchmarkFloat64CounterHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	cnt := fix.sdk.NewFloat64Counter("float64.counter")
	handle := cnt.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1.1)
	}
}

// Gauges

func BenchmarkInt64GaugeAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewInt64Gauge("int64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gau.Set(ctx, int64(i), labs)
	}
}

func BenchmarkInt64GaugeAcquireHandle(b *testing.B) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewInt64Gauge("int64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := gau.AcquireHandle(labs)
		handle.Release()
	}
}

func BenchmarkInt64GaugeHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewInt64Gauge("int64.gauge")
	handle := gau.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Set(ctx, int64(i))
	}
}

func BenchmarkFloat64GaugeAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewFloat64Gauge("float64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gau.Set(ctx, float64(i), labs)
	}
}

func BenchmarkFloat64GaugeAcquireHandle(b *testing.B) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewFloat64Gauge("float64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := gau.AcquireHandle(labs)
		handle.Release()
	}
}

func BenchmarkFloat64GaugeHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	gau := fix.sdk.NewFloat64Gauge("float64.gauge")
	handle := gau.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Set(ctx, float64(i))
	}
}

// Measures

func benchmarkInt64MeasureAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewInt64Measure(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i), labs)
	}
}

func benchmarkInt64MeasureAcquireHandle(b *testing.B, name string) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewInt64Measure(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := mea.AcquireHandle(labs)
		handle.Release()
	}
}

func benchmarkInt64MeasureHandleAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewInt64Measure(name)
	handle := mea.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, int64(i))
	}
}

func benchmarkFloat64MeasureAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewFloat64Measure(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i), labs)
	}
}

func benchmarkFloat64MeasureAcquireHandle(b *testing.B, name string) {
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewFloat64Measure(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle := mea.AcquireHandle(labs)
		handle.Release()
	}
}

func benchmarkFloat64MeasureHandleAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := fix.sdk.Labels(makeLabels(1)...)
	mea := fix.sdk.NewFloat64Measure(name)
	handle := mea.AcquireHandle(labs)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, float64(i))
	}
}

// MaxSumCount

func BenchmarkInt64MaxSumCountAdd(b *testing.B) {
	benchmarkInt64MeasureAdd(b, "int64.maxsumcount")
}

func BenchmarkInt64MaxSumCountAcquireHandle(b *testing.B) {
	benchmarkInt64MeasureAcquireHandle(b, "int64.maxsumcount")
}

func BenchmarkInt64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkInt64MeasureHandleAdd(b, "int64.maxsumcount")
}

func BenchmarkFloat64MaxSumCountAdd(b *testing.B) {
	benchmarkFloat64MeasureAdd(b, "float64.maxsumcount")
}

func BenchmarkFloat64MaxSumCountAcquireHandle(b *testing.B) {
	benchmarkFloat64MeasureAcquireHandle(b, "float64.maxsumcount")
}

func BenchmarkFloat64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkFloat64MeasureHandleAdd(b, "float64.maxsumcount")
}

// DDSketch

func BenchmarkInt64DDSketchAdd(b *testing.B) {
	benchmarkInt64MeasureAdd(b, "int64.ddsketch")
}

func BenchmarkInt64DDSketchAcquireHandle(b *testing.B) {
	benchmarkInt64MeasureAcquireHandle(b, "int64.ddsketch")
}

func BenchmarkInt64DDSketchHandleAdd(b *testing.B) {
	benchmarkInt64MeasureHandleAdd(b, "int64.ddsketch")
}

func BenchmarkFloat64DDSketchAdd(b *testing.B) {
	benchmarkFloat64MeasureAdd(b, "float64.ddsketch")
}

func BenchmarkFloat64DDSketchAcquireHandle(b *testing.B) {
	benchmarkFloat64MeasureAcquireHandle(b, "float64.ddsketch")
}

func BenchmarkFloat64DDSketchHandleAdd(b *testing.B) {
	benchmarkFloat64MeasureHandleAdd(b, "float64.ddsketch")
}
