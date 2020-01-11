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

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
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
	bf.sdk = sdk.New(bf, label.NewDefaultEncoder())
	return bf
}

func (*benchFixture) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	case export.MeasureKind:
		if strings.HasSuffix(descriptor.Name().String(), "minmaxsumcount") {
			return minmaxsumcount.New(descriptor)
		} else if strings.HasSuffix(descriptor.Name().String(), "ddsketch") {
			return ddsketch.New(ddsketch.NewDefaultConfig(), descriptor)
		} else if strings.HasSuffix(descriptor.Name().String(), "array") {
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

func makeLabelSets(n int) [][]core.KeyValue {
	r := make([][]core.KeyValue, n)

	for i := 0; i < n; i++ {
		r[i] = makeLabels(1)
	}

	return r
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
	labs := makeLabels(n)
	enc := label.NewDefaultEncoder()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		label.NewSet(labs...).Encoded(enc)
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

func BenchmarkAcquireNewHandle(b *testing.B) {
	labelSets := makeLabelSets(b.N)

	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk)

	cnt := scx.Meter().NewInt64Counter("int64.counter")
	ctxs := make([]context.Context, b.N)

	for i := 0; i < b.N; i++ {
		ctxs[i] = scx.AddResources(labelSets[i]...).InContext(context.Background())
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(ctxs[i])
	}
}

func BenchmarkAcquireExistingHandle(b *testing.B) {
	labelSets := makeLabelSets(b.N)

	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk)

	cnt := scx.Meter().NewInt64Counter("int64.counter")
	ctxs := make([]context.Context, b.N)

	for i := 0; i < b.N; i++ {
		ctxs[i] = scx.AddResources(labelSets[i]...).InContext(context.Background())
		cnt.Bind(ctxs[i]).Unbind()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(ctxs[i])
	}
}

func BenchmarkAcquireReleaseExistingHandle(b *testing.B) {
	labelSets := makeLabelSets(b.N)

	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk)

	cnt := scx.Meter().NewInt64Counter("int64.counter")
	ctxs := make([]context.Context, b.N)

	for i := 0; i < b.N; i++ {
		ctxs[i] = scx.AddResources(labelSets[i]...).InContext(context.Background())
		cnt.Bind(ctxs[i]).Unbind()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(ctxs[i]).Unbind()
	}
}

// Counters

func BenchmarkInt64CounterAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	cnt := scx.Meter().NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1)
	}
}

func BenchmarkInt64CounterHandleAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	cnt := scx.Meter().NewInt64Counter("int64.counter")

	handle := cnt.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1)
	}
}

func BenchmarkFloat64CounterAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	cnt := scx.Meter().NewFloat64Counter("float64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1.1)
	}
}

func BenchmarkFloat64CounterHandleAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	cnt := scx.Meter().NewFloat64Counter("float64.counter")
	handle := cnt.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1.1)
	}
}

// Gauges

func BenchmarkInt64GaugeAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	gau := scx.Meter().NewInt64Gauge("int64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gau.Set(ctx, int64(i))
	}
}

func BenchmarkInt64GaugeHandleAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	gau := scx.Meter().NewInt64Gauge("int64.gauge")
	handle := gau.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Set(ctx, int64(i))
	}
}

func BenchmarkFloat64GaugeAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	gau := scx.Meter().NewFloat64Gauge("float64.gauge")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gau.Set(ctx, float64(i))
	}
}

func BenchmarkFloat64GaugeHandleAdd(b *testing.B) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	gau := scx.Meter().NewFloat64Gauge("float64.gauge")
	handle := gau.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Set(ctx, float64(i))
	}
}

// Measures

func benchmarkInt64MeasureAdd(b *testing.B, name string) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	mea := scx.Meter().NewInt64Measure("int64.measure")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i))
	}
}

func benchmarkInt64MeasureHandleAdd(b *testing.B, name string) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	mea := scx.Meter().NewInt64Measure("int64.measure")
	handle := mea.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, int64(i))
	}
}

func benchmarkFloat64MeasureAdd(b *testing.B, name string) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	mea := scx.Meter().NewFloat64Measure("float64.measure")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i))
	}
}

func benchmarkFloat64MeasureHandleAdd(b *testing.B, name string) {
	fix := newFixture(b)
	scx := scope.WithMeterSDK(fix.sdk).AddResources(makeLabels(1)...)

	ctx := scx.InContext(context.Background())
	mea := scx.Meter().NewFloat64Measure("float64.measure")
	handle := mea.Bind(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, float64(i))
	}
}

// MaxSumCount

func BenchmarkInt64MaxSumCountAdd(b *testing.B) {
	benchmarkInt64MeasureAdd(b, "int64.minmaxsumcount")
}

func BenchmarkInt64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkInt64MeasureHandleAdd(b, "int64.minmaxsumcount")
}

func BenchmarkFloat64MaxSumCountAdd(b *testing.B) {
	benchmarkFloat64MeasureAdd(b, "float64.minmaxsumcount")
}

func BenchmarkFloat64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkFloat64MeasureHandleAdd(b, "float64.minmaxsumcount")
}

// DDSketch

func BenchmarkInt64DDSketchAdd(b *testing.B) {
	benchmarkInt64MeasureAdd(b, "int64.ddsketch")
}

func BenchmarkInt64DDSketchHandleAdd(b *testing.B) {
	benchmarkInt64MeasureHandleAdd(b, "int64.ddsketch")
}

func BenchmarkFloat64DDSketchAdd(b *testing.B) {
	benchmarkFloat64MeasureAdd(b, "float64.ddsketch")
}

func BenchmarkFloat64DDSketchHandleAdd(b *testing.B) {
	benchmarkFloat64MeasureHandleAdd(b, "float64.ddsketch")
}

// Array

func BenchmarkInt64ArrayAdd(b *testing.B) {
	benchmarkInt64MeasureAdd(b, "int64.array")
}

func BenchmarkInt64ArrayHandleAdd(b *testing.B) {
	benchmarkInt64MeasureHandleAdd(b, "int64.array")
}

func BenchmarkFloat64ArrayAdd(b *testing.B) {
	benchmarkFloat64MeasureAdd(b, "float64.array")
}

func BenchmarkFloat64ArrayHandleAdd(b *testing.B) {
	benchmarkFloat64MeasureHandleAdd(b, "float64.array")
}
