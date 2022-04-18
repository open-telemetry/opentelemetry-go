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

package metric_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type benchFixture struct {
	meter       metric.Meter
	accumulator *sdk.Accumulator
	B           *testing.B
	export.AggregatorSelector
}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B:                  b,
		AggregatorSelector: processortest.AggregatorSelector(),
	}

	bf.accumulator = sdk.NewAccumulator(bf)
	bf.meter = sdkapi.WrapMeterImpl(bf.accumulator)
	return bf
}

func (f *benchFixture) Process(export.Accumulation) error {
	return nil
}

func (f *benchFixture) Meter(_ string, _ ...metric.MeterOption) metric.Meter {
	return f.meter
}

func (f *benchFixture) iCounter(name string) syncint64.Counter {
	ctr, err := f.meter.SyncInt64().Counter(name)
	if err != nil {
		f.B.Error(err)
	}
	return ctr
}
func (f *benchFixture) fCounter(name string) syncfloat64.Counter {
	ctr, err := f.meter.SyncFloat64().Counter(name)
	if err != nil {
		f.B.Error(err)
	}
	return ctr
}
func (f *benchFixture) iHistogram(name string) syncint64.Histogram {
	ctr, err := f.meter.SyncInt64().Histogram(name)
	if err != nil {
		f.B.Error(err)
	}
	return ctr
}
func (f *benchFixture) fHistogram(name string) syncfloat64.Histogram {
	ctr, err := f.meter.SyncFloat64().Histogram(name)
	if err != nil {
		f.B.Error(err)
	}
	return ctr
}

func makeAttrs(n int) []attribute.KeyValue {
	used := map[string]bool{}
	l := make([]attribute.KeyValue, n)
	for i := 0; i < n; i++ {
		var k string
		for {
			k = fmt.Sprint("k", rand.Intn(1000000000))
			if !used[k] {
				used[k] = true
				break
			}
		}
		l[i] = attribute.String(k, fmt.Sprint("v", rand.Intn(1000000000)))
	}
	return l
}

func benchmarkAttrs(b *testing.B, n int) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(n)
	cnt := fix.iCounter("int64.sum")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkInt64CounterAddWithAttrs_1(b *testing.B) {
	benchmarkAttrs(b, 1)
}

func BenchmarkInt64CounterAddWithAttrs_2(b *testing.B) {
	benchmarkAttrs(b, 2)
}

func BenchmarkInt64CounterAddWithAttrs_4(b *testing.B) {
	benchmarkAttrs(b, 4)
}

func BenchmarkInt64CounterAddWithAttrs_8(b *testing.B) {
	benchmarkAttrs(b, 8)
}

func BenchmarkInt64CounterAddWithAttrs_16(b *testing.B) {
	benchmarkAttrs(b, 16)
}

// Note: performance does not depend on attribute set size for the benchmarks
// below--all are benchmarked for a single attribute.

// Iterators

var benchmarkIteratorVar attribute.KeyValue

func benchmarkIterator(b *testing.B, n int) {
	attrs := attribute.NewSet(makeAttrs(n)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := attrs.Iter()
		for iter.Next() {
			benchmarkIteratorVar = iter.Attribute()
		}
	}
}

func BenchmarkIterator_0(b *testing.B) {
	benchmarkIterator(b, 0)
}

func BenchmarkIterator_1(b *testing.B) {
	benchmarkIterator(b, 1)
}

func BenchmarkIterator_2(b *testing.B) {
	benchmarkIterator(b, 2)
}

func BenchmarkIterator_4(b *testing.B) {
	benchmarkIterator(b, 4)
}

func BenchmarkIterator_8(b *testing.B) {
	benchmarkIterator(b, 8)
}

func BenchmarkIterator_16(b *testing.B) {
	benchmarkIterator(b, 16)
}

// Counters

func BenchmarkGlobalInt64CounterAddWithSDK(b *testing.B) {
	// Compare with BenchmarkInt64CounterAdd() to see overhead of global
	// package. This is in the SDK to avoid the API from depending on the
	// SDK.
	ctx := context.Background()
	fix := newFixture(b)

	global.SetMeterProvider(fix)

	labs := []attribute.KeyValue{attribute.String("A", "B")}

	cnt := fix.iCounter("int64.sum")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkInt64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	cnt := fix.iCounter("int64.sum")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkFloat64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	cnt := fix.fCounter("float64.sum")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1.1, labs...)
	}
}

// LastValue

func BenchmarkInt64LastValueAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	mea := fix.iHistogram("int64.lastvalue")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i), labs...)
	}
}

func BenchmarkFloat64LastValueAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	mea := fix.fHistogram("float64.lastvalue")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i), labs...)
	}
}

// Histograms

func BenchmarkInt64HistogramAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	mea := fix.iHistogram("int64.histogram")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i), labs...)
	}
}

func BenchmarkFloat64HistogramAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	mea := fix.fHistogram("float64.histogram")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i), labs...)
	}
}

// Observers

func BenchmarkObserverRegistration(b *testing.B) {
	fix := newFixture(b)
	names := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		names = append(names, fmt.Sprintf("test.%d.lastvalue", i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctr, _ := fix.meter.AsyncInt64().Counter(names[i])
		_ = fix.meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(context.Context) {})
	}
}

func BenchmarkGaugeObserverObservationInt64(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	ctr, _ := fix.meter.AsyncInt64().Counter("test.lastvalue")
	err := fix.meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
		for i := 0; i < b.N; i++ {
			ctr.Observe(ctx, (int64)(i), labs...)
		}
	})
	if err != nil {
		b.Errorf("could not register callback: %v", err)
		b.FailNow()
	}

	b.ResetTimer()

	fix.accumulator.Collect(ctx)
}

func BenchmarkGaugeObserverObservationFloat64(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(1)
	ctr, _ := fix.meter.AsyncFloat64().Counter("test.lastvalue")
	err := fix.meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
		for i := 0; i < b.N; i++ {
			ctr.Observe(ctx, (float64)(i), labs...)
		}
	})
	if err != nil {
		b.Errorf("could not register callback: %v", err)
		b.FailNow()
	}

	b.ResetTimer()

	fix.accumulator.Collect(ctx)
}

// BatchRecord

func benchmarkBatchRecord8Attrs(b *testing.B, numInst int) {
	const numAttrs = 8
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeAttrs(numAttrs)
	var meas []syncint64.Counter

	for i := 0; i < numInst; i++ {
		meas = append(meas, fix.iCounter(fmt.Sprintf("int64.%d.sum", i)))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, ctr := range meas {
			ctr.Add(ctx, 1, labs...)
		}
	}
}

func BenchmarkBatchRecord8Attrs_1Instrument(b *testing.B) {
	benchmarkBatchRecord8Attrs(b, 1)
}

func BenchmarkBatchRecord_8Attrs_2Instruments(b *testing.B) {
	benchmarkBatchRecord8Attrs(b, 2)
}

func BenchmarkBatchRecord_8Attrs_4Instruments(b *testing.B) {
	benchmarkBatchRecord8Attrs(b, 4)
}

func BenchmarkBatchRecord_8Attrs_8Instruments(b *testing.B) {
	benchmarkBatchRecord8Attrs(b, 8)
}

// Record creation

func BenchmarkRepeatedDirectCalls(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)

	c := fix.iCounter("int64.sum")
	k := attribute.String("bench", "true")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Add(ctx, 1, k)
		fix.accumulator.Collect(ctx)
	}
}
