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
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type processFunc func(context.Context, export.Record) error

type benchFixture struct {
	meter       metric.MeterMust
	accumulator *sdk.Accumulator
	B           *testing.B
	pcb         processFunc
}

func newFixture(b *testing.B) *benchFixture {
	b.ReportAllocs()
	bf := &benchFixture{
		B: b,
	}

	bf.accumulator = sdk.NewAccumulator(bf)
	bf.meter = metric.Must(metric.WrapMeterImpl(bf.accumulator, "benchmarks"))
	return bf
}

func (f *benchFixture) setProcessCallback(cb processFunc) {
	f.pcb = cb
}

func (*benchFixture) AggregatorFor(descriptor *metric.Descriptor) export.Aggregator {
	name := descriptor.Name()
	switch {
	case strings.HasSuffix(name, "counter"):
		return sum.New()
	case strings.HasSuffix(name, "lastvalue"):
		return lastvalue.New()
	default:
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

func (f *benchFixture) Process(ctx context.Context, rec export.Record) error {
	if f.pcb == nil {
		return nil
	}
	return f.pcb(ctx, rec)
}

func (*benchFixture) CheckpointSet() export.CheckpointSet {
	return nil
}

func (*benchFixture) FinishedCollection() {
}

func makeManyLabels(n int) [][]kv.KeyValue {
	r := make([][]kv.KeyValue, n)

	for i := 0; i < n; i++ {
		r[i] = makeLabels(1)
	}

	return r
}

func makeLabels(n int) []kv.KeyValue {
	used := map[string]bool{}
	l := make([]kv.KeyValue, n)
	for i := 0; i < n; i++ {
		var k string
		for {
			k = fmt.Sprint("k", rand.Intn(1000000000))
			if !used[k] {
				used[k] = true
				break
			}
		}
		l[i] = kv.Key(k).String(fmt.Sprint("v", rand.Intn(1000000000)))
	}
	return l
}

func benchmarkLabels(b *testing.B, n int) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(n)
	cnt := fix.meter.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkInt64CounterAddWithLabels_1(b *testing.B) {
	benchmarkLabels(b, 1)
}

func BenchmarkInt64CounterAddWithLabels_2(b *testing.B) {
	benchmarkLabels(b, 2)
}

func BenchmarkInt64CounterAddWithLabels_4(b *testing.B) {
	benchmarkLabels(b, 4)
}

func BenchmarkInt64CounterAddWithLabels_8(b *testing.B) {
	benchmarkLabels(b, 8)
}

func BenchmarkInt64CounterAddWithLabels_16(b *testing.B) {
	benchmarkLabels(b, 16)
}

// Note: performance does not depend on label set size for the
// benchmarks below--all are benchmarked for a single label.

func BenchmarkAcquireNewHandle(b *testing.B) {
	fix := newFixture(b)
	labelSets := makeManyLabels(b.N)
	cnt := fix.meter.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(labelSets[i]...)
	}
}

func BenchmarkAcquireExistingHandle(b *testing.B) {
	fix := newFixture(b)
	labelSets := makeManyLabels(b.N)
	cnt := fix.meter.NewInt64Counter("int64.counter")

	for i := 0; i < b.N; i++ {
		cnt.Bind(labelSets[i]...).Unbind()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(labelSets[i]...)
	}
}

func BenchmarkAcquireReleaseExistingHandle(b *testing.B) {
	fix := newFixture(b)
	labelSets := makeManyLabels(b.N)
	cnt := fix.meter.NewInt64Counter("int64.counter")

	for i := 0; i < b.N; i++ {
		cnt.Bind(labelSets[i]...).Unbind()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Bind(labelSets[i]...).Unbind()
	}
}

// Iterators

var benchmarkIteratorVar kv.KeyValue

func benchmarkIterator(b *testing.B, n int) {
	fix := newFixture(b)
	fix.setProcessCallback(func(ctx context.Context, rec export.Record) error {
		var kv kv.KeyValue
		li := rec.Labels().Iter()
		fix.B.StartTimer()
		for i := 0; i < fix.B.N; i++ {
			iter := li
			// test getting only the first element
			if iter.Next() {
				kv = iter.Label()
			}
		}
		fix.B.StopTimer()
		benchmarkIteratorVar = kv
		return nil
	})
	cnt := fix.meter.NewInt64Counter("int64.counter")
	ctx := context.Background()
	cnt.Add(ctx, 1, makeLabels(n)...)

	b.ResetTimer()
	fix.accumulator.Collect(ctx)
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

func BenchmarkInt64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	cnt := fix.meter.NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkInt64CounterHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	cnt := fix.meter.NewInt64Counter("int64.counter")
	handle := cnt.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1)
	}
}

func BenchmarkFloat64CounterAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	cnt := fix.meter.NewFloat64Counter("float64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1.1, labs...)
	}
}

func BenchmarkFloat64CounterHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	cnt := fix.meter.NewFloat64Counter("float64.counter")
	handle := cnt.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Add(ctx, 1.1)
	}
}

// LastValue

func BenchmarkInt64LastValueAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewInt64ValueRecorder("int64.lastvalue")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i), labs...)
	}
}

func BenchmarkInt64LastValueHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewInt64ValueRecorder("int64.lastvalue")
	handle := mea.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, int64(i))
	}
}

func BenchmarkFloat64LastValueAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewFloat64ValueRecorder("float64.lastvalue")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i), labs...)
	}
}

func BenchmarkFloat64LastValueHandleAdd(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewFloat64ValueRecorder("float64.lastvalue")
	handle := mea.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, float64(i))
	}
}

// ValueRecorders

func benchmarkInt64ValueRecorderAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewInt64ValueRecorder(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, int64(i), labs...)
	}
}

func benchmarkInt64ValueRecorderHandleAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewInt64ValueRecorder(name)
	handle := mea.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, int64(i))
	}
}

func benchmarkFloat64ValueRecorderAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewFloat64ValueRecorder(name)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mea.Record(ctx, float64(i), labs...)
	}
}

func benchmarkFloat64ValueRecorderHandleAdd(b *testing.B, name string) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	mea := fix.meter.NewFloat64ValueRecorder(name)
	handle := mea.Bind(labs...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		handle.Record(ctx, float64(i))
	}
}

// Observers

func BenchmarkObserverRegistration(b *testing.B) {
	fix := newFixture(b)
	names := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		names = append(names, fmt.Sprintf("test.valueobserver.%d", i))
	}
	cb := func(_ context.Context, result metric.Int64ObserverResult) {}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fix.meter.NewInt64ValueObserver(names[i], cb)
	}
}

func BenchmarkValueObserverObservationInt64(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	_ = fix.meter.NewInt64ValueObserver("test.valueobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		for i := 0; i < b.N; i++ {
			result.Observe((int64)(i), labs...)
		}
	})

	b.ResetTimer()

	fix.accumulator.Collect(ctx)
}

func BenchmarkValueObserverObservationFloat64(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(1)
	_ = fix.meter.NewFloat64ValueObserver("test.valueobserver", func(_ context.Context, result metric.Float64ObserverResult) {
		for i := 0; i < b.N; i++ {
			result.Observe((float64)(i), labs...)
		}
	})

	b.ResetTimer()

	fix.accumulator.Collect(ctx)
}

// MaxSumCount

func BenchmarkInt64MaxSumCountAdd(b *testing.B) {
	benchmarkInt64ValueRecorderAdd(b, "int64.minmaxsumcount")
}

func BenchmarkInt64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkInt64ValueRecorderHandleAdd(b, "int64.minmaxsumcount")
}

func BenchmarkFloat64MaxSumCountAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderAdd(b, "float64.minmaxsumcount")
}

func BenchmarkFloat64MaxSumCountHandleAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderHandleAdd(b, "float64.minmaxsumcount")
}

// DDSketch

func BenchmarkInt64DDSketchAdd(b *testing.B) {
	benchmarkInt64ValueRecorderAdd(b, "int64.ddsketch")
}

func BenchmarkInt64DDSketchHandleAdd(b *testing.B) {
	benchmarkInt64ValueRecorderHandleAdd(b, "int64.ddsketch")
}

func BenchmarkFloat64DDSketchAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderAdd(b, "float64.ddsketch")
}

func BenchmarkFloat64DDSketchHandleAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderHandleAdd(b, "float64.ddsketch")
}

// Array

func BenchmarkInt64ArrayAdd(b *testing.B) {
	benchmarkInt64ValueRecorderAdd(b, "int64.array")
}

func BenchmarkInt64ArrayHandleAdd(b *testing.B) {
	benchmarkInt64ValueRecorderHandleAdd(b, "int64.array")
}

func BenchmarkFloat64ArrayAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderAdd(b, "float64.array")
}

func BenchmarkFloat64ArrayHandleAdd(b *testing.B) {
	benchmarkFloat64ValueRecorderHandleAdd(b, "float64.array")
}

// BatchRecord

func benchmarkBatchRecord8Labels(b *testing.B, numInst int) {
	const numLabels = 8
	ctx := context.Background()
	fix := newFixture(b)
	labs := makeLabels(numLabels)
	var meas []metric.Measurement

	for i := 0; i < numInst; i++ {
		inst := fix.meter.NewInt64Counter(fmt.Sprint("int64.counter.", i))
		meas = append(meas, inst.Measurement(1))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fix.accumulator.RecordBatch(ctx, labs, meas...)
	}
}

func BenchmarkBatchRecord8Labels_1Instrument(b *testing.B) {
	benchmarkBatchRecord8Labels(b, 1)
}

func BenchmarkBatchRecord_8Labels_2Instruments(b *testing.B) {
	benchmarkBatchRecord8Labels(b, 2)
}

func BenchmarkBatchRecord_8Labels_4Instruments(b *testing.B) {
	benchmarkBatchRecord8Labels(b, 4)
}

func BenchmarkBatchRecord_8Labels_8Instruments(b *testing.B) {
	benchmarkBatchRecord8Labels(b, 8)
}

// Record creation

func BenchmarkRepeatedDirectCalls(b *testing.B) {
	ctx := context.Background()
	fix := newFixture(b)
	encoder := label.DefaultEncoder()
	fix.pcb = func(_ context.Context, rec export.Record) error {
		_ = rec.Labels().Encoded(encoder)
		return nil
	}

	c := fix.meter.NewInt64Counter("int64.counter")
	k := kv.String("bench", "true")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.Add(ctx, 1, k)
		fix.accumulator.Collect(ctx)
	}
}

// LabelIterator

func BenchmarkLabelIterator(b *testing.B) {
	const labelCount = 1024
	ctx := context.Background()
	fix := newFixture(b)

	var rec export.Record
	fix.pcb = func(_ context.Context, processRec export.Record) error {
		rec = processRec
		return nil
	}

	keyValues := makeLabels(labelCount)
	counter := fix.meter.NewInt64Counter("test.counter")
	counter.Add(ctx, 1, keyValues...)

	fix.accumulator.Collect(ctx)

	b.ResetTimer()

	labels := rec.Labels()
	iter := labels.Iter()
	var val kv.KeyValue
	for i := 0; i < b.N; i++ {
		if !iter.Next() {
			iter = labels.Iter()
			iter.Next()
		}
		val = iter.Label()
	}
	if false {
		fmt.Println(val)
	}
}
