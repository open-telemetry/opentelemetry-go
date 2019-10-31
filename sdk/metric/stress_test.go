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
	"math"
	"math/rand"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	api "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

const (
	concurrencyPerCPU = 100
	reclaimPeriod     = time.Millisecond * 100
	testRun           = time.Second
	epsilon           = 1e-10
)

type (
	testFixture struct {
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

	testKey struct {
		labels     string
		descriptor *export.Descriptor
	}

	testImpl struct {
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

	withImpl interface {
		Impl() metric.InstrumentImpl
	}

	// gaugeState supports merging gauge values, for the case
	// where a race condition causes duplicate records.  We always
	// take the later timestamp.
	gaugeState struct {
		raw core.Number
		ts  time.Time
	}
)

func concurrency() int {
	return concurrencyPerCPU * runtime.NumCPU()
}

func canonicalizeLabels(ls []core.KeyValue) string {
	copy := append(ls[0:0:0], ls...)
	sort.SliceStable(copy, func(i, j int) bool {
		return copy[i].Key < copy[j].Key
	})
	var b strings.Builder
	for _, kv := range copy {
		b.WriteString(string(kv.Key))
		b.WriteString("=")
		b.WriteString(kv.Value.Emit())
		b.WriteString("$")
	}
	return b.String()
}

func getPeriod() time.Duration {
	dur := math.Max(
		float64(reclaimPeriod)/10,
		float64(reclaimPeriod)*(1+0.1*rand.NormFloat64()),
	)
	return time.Duration(dur)
}

func (f *testFixture) someLabels() []core.KeyValue {
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

func (f *testFixture) startWorker(sdk *sdk.SDK, wg *sync.WaitGroup, i int) {
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

func (f *testFixture) assertTest(numCollect int) {
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

func (f *testFixture) preCollect() {
	// Collect calls Export in a single-threaded context. No need
	// to lock this struct.
	f.dupCheck = map[testKey]int{}
}

func (f *testFixture) AggregatorFor(record export.MetricRecord) export.MetricAggregator {
	switch record.Descriptor().MetricKind() {
	case export.CounterMetricKind:
		return counter.New()
	case export.GaugeMetricKind:
		return gauge.New()
	default:
		panic("Not implemented for this test")
	}
}

func (f *testFixture) Export(ctx context.Context, record export.MetricRecord, agg export.MetricAggregator) {
	desc := record.Descriptor()
	key := testKey{
		labels:     canonicalizeLabels(record.Labels()),
		descriptor: desc,
	}
	if f.dupCheck[key] == 0 {
		f.dupCheck[key]++
	} else {
		f.totalDups++
	}

	actual, _ := f.received.LoadOrStore(key, f.impl.newStore())

	switch desc.MetricKind() {
	case export.CounterMetricKind:
		f.impl.storeCollect(actual, agg.(*counter.Aggregator).AsNumber(), time.Time{})
	case export.GaugeMetricKind:
		gauge := agg.(*gauge.Aggregator)
		f.impl.storeCollect(actual, gauge.AsNumber(), gauge.Timestamp())
	default:
		panic("Not used in this test")
	}
}

func stressTest(t *testing.T, impl testImpl) {
	ctx := context.Background()
	t.Parallel()
	fixture := &testFixture{
		T:     t,
		impl:  impl,
		lused: map[string]bool{},
	}
	cc := concurrency()
	sdk := sdk.New(fixture)
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

func int64sEqual(a, b core.Number) bool {
	return a.AsInt64() == b.AsInt64()
}

func float64sEqual(a, b core.Number) bool {
	diff := math.Abs(a.AsFloat64() - b.AsFloat64())
	return diff < math.Abs(a.AsFloat64())*epsilon
}

// Counters

func intCounterTestImpl(nonMonotonic bool) testImpl {
	return testImpl{
		newInstrument: func(meter api.Meter, name string) withImpl {
			return meter.NewInt64Counter(name, api.WithMonotonic(!nonMonotonic))
		},
		getUpdateValue: func() core.Number {
			var offset int64
			if nonMonotonic {
				offset = -50
			}
			for {
				x := offset + int64(rand.Intn(100))
				if x != 0 {
					return core.NewInt64Number(x)
				}
			}
		},
		operate: func(inst interface{}, ctx context.Context, value core.Number, labels api.LabelSet) {
			counter := inst.(api.Int64Counter)
			counter.Add(ctx, value.AsInt64(), labels)
		},
		newStore: func() interface{} {
			n := core.NewInt64Number(0)
			return &n
		},
		storeCollect: func(store interface{}, value core.Number, _ time.Time) {
			store.(*core.Number).AddInt64Atomic(value.AsInt64())
		},
		storeExpect: func(store interface{}, value core.Number) {
			store.(*core.Number).AddInt64Atomic(value.AsInt64())
		},
		readStore: func(store interface{}) core.Number {
			return store.(*core.Number).AsNumberAtomic()
		},
		equalValues: int64sEqual,
	}
}

func TestStressInt64CounterNormal(t *testing.T) {
	stressTest(t, intCounterTestImpl(false))
}

func TestStressInt64CounterNonMonotonic(t *testing.T) {
	stressTest(t, intCounterTestImpl(true))
}

func floatCounterTestImpl(nonMonotonic bool) testImpl {
	return testImpl{
		newInstrument: func(meter api.Meter, name string) withImpl {
			return meter.NewFloat64Counter(name, api.WithMonotonic(!nonMonotonic))
		},
		getUpdateValue: func() core.Number {
			var offset float64
			if nonMonotonic {
				offset = -0.5
			}
			for {
				x := offset + rand.Float64()
				if x != 0 {
					return core.NewFloat64Number(x)
				}
			}
		},
		operate: func(inst interface{}, ctx context.Context, value core.Number, labels api.LabelSet) {
			counter := inst.(api.Float64Counter)
			counter.Add(ctx, value.AsFloat64(), labels)
		},
		newStore: func() interface{} {
			n := core.NewFloat64Number(0.0)
			return &n
		},
		storeCollect: func(store interface{}, value core.Number, _ time.Time) {
			store.(*core.Number).AddFloat64Atomic(value.AsFloat64())
		},
		storeExpect: func(store interface{}, value core.Number) {
			store.(*core.Number).AddFloat64Atomic(value.AsFloat64())
		},
		readStore: func(store interface{}) core.Number {
			return store.(*core.Number).AsNumberAtomic()
		},
		equalValues: float64sEqual,
	}
}

func TestStressFloat64CounterNormal(t *testing.T) {
	stressTest(t, floatCounterTestImpl(false))
}

func TestStressFloat64CounterNonMonotonic(t *testing.T) {
	stressTest(t, floatCounterTestImpl(true))
}

// Gauges

func intGaugeTestImpl(monotonic bool) testImpl {
	// (Now()-startTime) is used as a free monotonic source
	startTime := time.Now()

	return testImpl{
		newInstrument: func(meter api.Meter, name string) withImpl {
			return meter.NewInt64Gauge(name, api.WithMonotonic(monotonic))
		},
		getUpdateValue: func() core.Number {
			if !monotonic {
				r1 := rand.Int63()
				return core.NewInt64Number(rand.Int63() - r1)
			}
			return core.NewInt64Number(int64(time.Since(startTime)))
		},
		operate: func(inst interface{}, ctx context.Context, value core.Number, labels api.LabelSet) {
			gauge := inst.(api.Int64Gauge)
			gauge.Set(ctx, value.AsInt64(), labels)
		},
		newStore: func() interface{} {
			return &gaugeState{
				raw: core.NewInt64Number(0),
			}
		},
		storeCollect: func(store interface{}, value core.Number, ts time.Time) {
			gs := store.(*gaugeState)

			if !ts.Before(gs.ts) {
				gs.ts = ts
				gs.raw.SetInt64Atomic(value.AsInt64())
			}
		},
		storeExpect: func(store interface{}, value core.Number) {
			gs := store.(*gaugeState)
			gs.raw.SetInt64Atomic(value.AsInt64())
		},
		readStore: func(store interface{}) core.Number {
			gs := store.(*gaugeState)
			return gs.raw.AsNumberAtomic()
		},
		equalValues: int64sEqual,
	}
}

func TestStressInt64GaugeNormal(t *testing.T) {
	stressTest(t, intGaugeTestImpl(false))
}

func TestStressInt64GaugeMonotonic(t *testing.T) {
	stressTest(t, intGaugeTestImpl(true))
}

func floatGaugeTestImpl(monotonic bool) testImpl {
	// (Now()-startTime) is used as a free monotonic source
	startTime := time.Now()

	return testImpl{
		newInstrument: func(meter api.Meter, name string) withImpl {
			return meter.NewFloat64Gauge(name, api.WithMonotonic(monotonic))
		},
		getUpdateValue: func() core.Number {
			if !monotonic {
				return core.NewFloat64Number((-0.5 + rand.Float64()) * 100000)
			}
			return core.NewFloat64Number(float64(time.Since(startTime)))
		},
		operate: func(inst interface{}, ctx context.Context, value core.Number, labels api.LabelSet) {
			gauge := inst.(api.Float64Gauge)
			gauge.Set(ctx, value.AsFloat64(), labels)
		},
		newStore: func() interface{} {
			return &gaugeState{
				raw: core.NewFloat64Number(0),
			}
		},
		storeCollect: func(store interface{}, value core.Number, ts time.Time) {
			gs := store.(*gaugeState)

			if !ts.Before(gs.ts) {
				gs.ts = ts
				gs.raw.SetFloat64Atomic(value.AsFloat64())
			}
		},
		storeExpect: func(store interface{}, value core.Number) {
			gs := store.(*gaugeState)
			gs.raw.SetFloat64Atomic(value.AsFloat64())
		},
		readStore: func(store interface{}) core.Number {
			gs := store.(*gaugeState)
			return gs.raw.AsNumberAtomic()
		},
		equalValues: float64sEqual,
	}
}

func TestStressFloat64GaugeNormal(t *testing.T) {
	stressTest(t, floatGaugeTestImpl(false))
}

func TestStressFloat64GaugeMonotonic(t *testing.T) {
	stressTest(t, floatGaugeTestImpl(true))
}
