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

// This test is too large for the race detector.  This SDK uses no locks
// that the race detector would help with, anyway.
// +build !race

package metric_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	api "go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
)

type (
	Batcher struct {
		// stop has to be aligned for 64-bit atomic operations.
		stop     int64
		expected sync.Map
		received sync.Map // Note: doesn't require synchronization
		wg       sync.WaitGroup
		impl     testImpl
		T        *testing.T

		lock  sync.Mutex
		lused map[string]bool

		dupCheck  map[testKey]int
		totalDups int64
	}

	htestImpl struct {
		newInstrument  func(meter api.Meter, name string) withImpl
		getUpdateValue func() core.Number
		operate        func(interface{}, context.Context, core.Number, api.LabelSet)
		newStore       func() interface{}

		// storeCollect and storeExpect are the same for
		// counters, different for gauges, to ensure we are
		// testing the timestamps correctly.
		storeCollect func(store interface{}, value core.Number, ts time.Time)
		storeExpect  func(store interface{}, value core.Number)
		readStore    func(store interface{}) core.Number
		equalValues  func(a, b core.Number) bool
	}

	// gaugeState supports merging gauge values, for the case
	// where a race condition causes duplicate records.  We always
	// take the later timestamp.
	histogramState struct {
		// raw has to be aligned for 64-bit atomic operations.
		count core.Number
	}
)

func (f *Batcher) someLabels() []core.KeyValue {
	n := 1 + rand.Intn(3)
	l := make([]core.KeyValue, n)

	for {
		oused := map[string]bool{}
		for i := 0; i < n; i++ {
			var k string
			for {
				k = fmt.Sprint("k", rand.Intn(1000000000))
				if !oused[k] {
					oused[k] = true
					break
				}
			}
			l[i] = key.New(k).String(fmt.Sprint("v", rand.Intn(1000000000)))
		}
		lc := canonicalizeLabels(l)
		f.lock.Lock()
		avail := !f.lused[lc]
		if avail {
			f.lused[lc] = true
			f.lock.Unlock()
			return l
		}
		f.lock.Unlock()
	}
}

func (f *Batcher) startWorker(sdk *sdk.SDK, wg *sync.WaitGroup, i int) {
	ctx := context.Background()
	name := fmt.Sprint("test_", i)
	instrument := f.impl.newInstrument(sdk, name)
	descriptor := sdk.GetDescriptor(instrument.Impl())
	kvs := f.someLabels()
	clabs := canonicalizeLabels(kvs)
	labs := sdk.Labels(kvs...)
	dur := getPeriod()
	key := testKey{
		labels:     clabs,
		descriptor: descriptor,
	}
	for {
		sleep := time.Duration(rand.ExpFloat64() * float64(dur))
		time.Sleep(sleep)
		value := f.impl.getUpdateValue()
		f.impl.operate(instrument, ctx, value, labs)

		actual, _ := f.expected.LoadOrStore(key, f.impl.newStore())

		f.impl.storeExpect(actual, value)

		if atomic.LoadInt64(&f.stop) != 0 {
			wg.Done()
			return
		}
	}
}

func (f *Batcher) assertTest(numCollect int) {
	csize := 0
	f.received.Range(func(key, gstore interface{}) bool {
		csize++
		gvalue := f.impl.readStore(gstore)

		estore, loaded := f.expected.Load(key)
		if !loaded {
			f.T.Error("Could not locate expected key: ", key)
		}
		evalue := f.impl.readStore(estore)

		if !f.impl.equalValues(evalue, gvalue) {
			f.T.Error("Expected value mismatch: ",
				evalue, "!=", gvalue, " for ", key)
		}
		return true
	})
	rsize := 0
	f.expected.Range(func(key, value interface{}) bool {
		rsize++
		if _, loaded := f.received.Load(key); !loaded {
			f.T.Error("Did not receive expected key: ", key)
		}
		return true
	})
	if rsize != csize {
		f.T.Error("Did not receive the correct set of metrics: Received != Expected", rsize, csize)
	}

	// Note: It's useful to know the test triggers this condition,
	// but we can't assert it.  Infrequently no duplicates are
	// found, and we can't really force a race to happen.
	//
	// fmt.Printf("Test duplicate records seen: %.1f%%\n",
	// 	float64(100*f.totalDups/int64(numCollect*concurrency())))
}

func (f *Batcher) preCollect() {
	// Collect calls Process in a single-threaded context. No need
	// to lock this struct.
	f.dupCheck = map[testKey]int{}
}

func (*Batcher) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	switch descriptor.MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	case export.MeasureKind:
		return histogram.New(descriptor, []core.Number{core.NewInt64Number(25), core.NewInt64Number(50), core.NewInt64Number(75)})
	default:
		panic("Not implemented for this test")
	}
}

func (*Batcher) CheckpointSet() export.CheckpointSet {
	return nil
}

func (*Batcher) FinishedCollection() {
}

func (f *Batcher) Process(_ context.Context, record export.Record) error {
	key := testKey{
		labels:     canonicalizeLabels(record.Labels().Ordered()),
		descriptor: record.Descriptor(),
	}
	if f.dupCheck[key] == 0 {
		f.dupCheck[key]++
	} else {
		f.totalDups++
	}

	actual, _ := f.received.LoadOrStore(key, f.impl.newStore())

	agg := record.Aggregator()
	switch record.Descriptor().MetricKind() {
	case export.MeasureKind:
		hist := agg.(aggregator.Histogram)
		b, err := hist.Histogram()
		if err != nil {
			f.T.Fatal("Sum error: ", err)
		}

		var count core.Number
		for _, n := range b.Counts {
			count.AddUint64(n.AsUint64())
		}

		f.impl.storeCollect(actual, count, time.Time{})
	default:
		panic("Not used in this test")
	}
	return nil
}

func histogramStressTest(t *testing.T, impl testImpl) {
	ctx := context.Background()
	t.Parallel()
	fixture := &Batcher{
		T:     t,
		impl:  impl,
		lused: map[string]bool{},
	}
	cc := concurrency()
	sdk := sdk.New(fixture, sdk.NewDefaultLabelEncoder())
	fixture.wg.Add(cc + 1)

	for i := 0; i < cc; i++ {
		go fixture.startWorker(sdk, &fixture.wg, i)
	}

	numCollect := 0

	go func() {
		for {
			time.Sleep(reclaimPeriod)
			fixture.preCollect()
			sdk.Collect(ctx)
			numCollect++
			if atomic.LoadInt64(&fixture.stop) != 0 {
				fixture.wg.Done()
				return
			}
		}
	}()

	time.Sleep(testRun)
	atomic.StoreInt64(&fixture.stop, 1)
	fixture.wg.Wait()
	fixture.preCollect()
	sdk.Collect(ctx)
	numCollect++

	fixture.assertTest(numCollect)
}

func TestStressInt64Histogram(t *testing.T) {
	timpl := testImpl{
		newInstrument: func(meter api.Meter, name string) withImpl {
			return meter.NewInt64Measure(name)
		},
		getUpdateValue: func() core.Number {
			for {
				x := int64(rand.Intn(100))
				if x != 0 {
					return core.NewInt64Number(x)
				}
			}
		},
		operate: func(inst interface{}, ctx context.Context, value core.Number, labels api.LabelSet) {
			counter := inst.(api.Int64Measure)
			counter.Record(ctx, value.AsInt64(), labels)
		},
		newStore: func() interface{} {
			n := core.NewInt64Number(0)
			return &n
		},
		storeCollect: func(store interface{}, value core.Number, _ time.Time) {
			store.(*core.Number).AddInt64Atomic(value.AsInt64())
		},
		storeExpect: func(store interface{}, value core.Number) {
			store.(*core.Number).AddInt64Atomic(1)
		},
		readStore: func(store interface{}) core.Number {
			return store.(*core.Number).AsNumberAtomic()
		},
		equalValues: int64sEqual,
	}
	histogramStressTest(t, timpl)
}
