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

// This test is too large for the race detector.  This SDK uses no locks
// that the race detector would help with, anyway.
// +build !race

package metric

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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
)

const (
	concurrencyPerCPU = 100
	reclaimPeriod     = time.Millisecond * 100
	testRun           = 5 * time.Second
	epsilon           = 1e-10
)

var Must = otel.Must

type (
	testFixture struct {
		// stop has to be aligned for 64-bit atomic operations.
		stop     int64
		expected sync.Map
		received sync.Map // Note: doesn't require synchronization
		wg       sync.WaitGroup
		impl     testImpl
		T        *testing.T

		export.AggregatorSelector

		lock  sync.Mutex
		lused map[string]bool

		dupCheck  map[testKey]int
		totalDups int64
	}

	testKey struct {
		labels     string
		descriptor *otel.Descriptor
	}

	testImpl struct {
		newInstrument  func(meter otel.Meter, name string) SyncImpler
		getUpdateValue func() otel.Number
		operate        func(interface{}, context.Context, otel.Number, []label.KeyValue)
		newStore       func() interface{}

		// storeCollect and storeExpect are the same for
		// counters, different for lastValues, to ensure we are
		// testing the timestamps correctly.
		storeCollect func(store interface{}, value otel.Number, ts time.Time)
		storeExpect  func(store interface{}, value otel.Number)
		readStore    func(store interface{}) otel.Number
		equalValues  func(a, b otel.Number) bool
	}

	SyncImpler interface {
		SyncImpl() otel.SyncImpl
	}

	// lastValueState supports merging lastValue values, for the case
	// where a race condition causes duplicate records.  We always
	// take the later timestamp.
	lastValueState struct {
		// raw has to be aligned for 64-bit atomic operations.
		raw otel.Number
		ts  time.Time
	}
)

func concurrency() int {
	return concurrencyPerCPU * runtime.NumCPU()
}

func canonicalizeLabels(ls []label.KeyValue) string {
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

func (f *testFixture) someLabels() []label.KeyValue {
	n := 1 + rand.Intn(3)
	l := make([]label.KeyValue, n)

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
			l[i] = label.String(k, fmt.Sprint("v", rand.Intn(1000000000)))
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

func (f *testFixture) startWorker(impl *Accumulator, meter otel.Meter, wg *sync.WaitGroup, i int) {
	ctx := context.Background()
	name := fmt.Sprint("test_", i)
	instrument := f.impl.newInstrument(meter, name)
	var descriptor *otel.Descriptor
	if ii, ok := instrument.SyncImpl().(*syncInstrument); ok {
		descriptor = &ii.descriptor
	}
	kvs := f.someLabels()
	clabs := canonicalizeLabels(kvs)
	dur := getPeriod()
	key := testKey{
		labels:     clabs,
		descriptor: descriptor,
	}
	for {
		sleep := time.Duration(rand.ExpFloat64() * float64(dur))
		time.Sleep(sleep)
		value := f.impl.getUpdateValue()
		f.impl.operate(instrument, ctx, value, kvs)

		actual, _ := f.expected.LoadOrStore(key, f.impl.newStore())

		f.impl.storeExpect(actual, value)

		if atomic.LoadInt64(&f.stop) != 0 {
			wg.Done()
			return
		}
	}
}

func (f *testFixture) assertTest(numCollect int) {
	var allErrs []func()
	csize := 0
	f.received.Range(func(key, gstore interface{}) bool {
		csize++
		gvalue := f.impl.readStore(gstore)

		estore, loaded := f.expected.Load(key)
		if !loaded {
			allErrs = append(allErrs, func() {
				f.T.Error("Could not locate expected key: ", key)
			})
			return true
		}
		evalue := f.impl.readStore(estore)

		if !f.impl.equalValues(evalue, gvalue) {
			allErrs = append(allErrs, func() {
				f.T.Error("Expected value mismatch: ",
					evalue, "!=", gvalue, " for ", key)
			})
		}
		return true
	})
	rsize := 0
	f.expected.Range(func(key, value interface{}) bool {
		rsize++
		if _, loaded := f.received.Load(key); !loaded {
			allErrs = append(allErrs, func() {
				f.T.Error("Did not receive expected key: ", key)
			})
		}
		return true
	})
	if rsize != csize {
		f.T.Error("Did not receive the correct set of metrics: Received != Expected", rsize, csize)
	}

	for _, anErr := range allErrs {
		anErr()
	}

	// Note: It's useful to know the test triggers this condition,
	// but we can't assert it.  Infrequently no duplicates are
	// found, and we can't really force a race to happen.
	//
	// fmt.Printf("Test duplicate records seen: %.1f%%\n",
	// 	float64(100*f.totalDups/int64(numCollect*concurrency())))
}

func (f *testFixture) preCollect() {
	// Collect calls Process in a single-threaded context. No need
	// to lock this struct.
	f.dupCheck = map[testKey]int{}
}

func (*testFixture) CheckpointSet() export.CheckpointSet {
	return nil
}

func (f *testFixture) Process(accumulation export.Accumulation) error {
	labels := accumulation.Labels().ToSlice()
	key := testKey{
		labels:     canonicalizeLabels(labels),
		descriptor: accumulation.Descriptor(),
	}
	if f.dupCheck[key] == 0 {
		f.dupCheck[key]++
	} else {
		f.totalDups++
	}

	actual, _ := f.received.LoadOrStore(key, f.impl.newStore())

	agg := accumulation.Aggregator()
	switch accumulation.Descriptor().InstrumentKind() {
	case otel.CounterInstrumentKind:
		sum, err := agg.(aggregation.Sum).Sum()
		if err != nil {
			f.T.Fatal("Sum error: ", err)
		}
		f.impl.storeCollect(actual, sum, time.Time{})
	case otel.ValueRecorderInstrumentKind:
		lv, ts, err := agg.(aggregation.LastValue).LastValue()
		if err != nil && err != aggregation.ErrNoData {
			f.T.Fatal("Last value error: ", err)
		}
		f.impl.storeCollect(actual, lv, ts)
	default:
		panic("Not used in this test")
	}
	return nil
}

func stressTest(t *testing.T, impl testImpl) {
	ctx := context.Background()
	t.Parallel()
	fixture := &testFixture{
		T:                  t,
		impl:               impl,
		lused:              map[string]bool{},
		AggregatorSelector: processortest.AggregatorSelector(),
	}
	cc := concurrency()
	sdk := NewAccumulator(fixture)
	meter := otel.WrapMeterImpl(sdk, "stress_test")
	fixture.wg.Add(cc + 1)

	for i := 0; i < cc; i++ {
		go fixture.startWorker(sdk, meter, &fixture.wg, i)
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

func int64sEqual(a, b otel.Number) bool {
	return a.AsInt64() == b.AsInt64()
}

func float64sEqual(a, b otel.Number) bool {
	diff := math.Abs(a.AsFloat64() - b.AsFloat64())
	return diff < math.Abs(a.AsFloat64())*epsilon
}

// Counters

func intCounterTestImpl() testImpl {
	return testImpl{
		newInstrument: func(meter otel.Meter, name string) SyncImpler {
			return Must(meter).NewInt64Counter(name + ".sum")
		},
		getUpdateValue: func() otel.Number {
			for {
				x := int64(rand.Intn(100))
				if x != 0 {
					return otel.NewInt64Number(x)
				}
			}
		},
		operate: func(inst interface{}, ctx context.Context, value otel.Number, labels []label.KeyValue) {
			counter := inst.(otel.Int64Counter)
			counter.Add(ctx, value.AsInt64(), labels...)
		},
		newStore: func() interface{} {
			n := otel.NewInt64Number(0)
			return &n
		},
		storeCollect: func(store interface{}, value otel.Number, _ time.Time) {
			store.(*otel.Number).AddInt64Atomic(value.AsInt64())
		},
		storeExpect: func(store interface{}, value otel.Number) {
			store.(*otel.Number).AddInt64Atomic(value.AsInt64())
		},
		readStore: func(store interface{}) otel.Number {
			return store.(*otel.Number).AsNumberAtomic()
		},
		equalValues: int64sEqual,
	}
}

func TestStressInt64Counter(t *testing.T) {
	stressTest(t, intCounterTestImpl())
}

func floatCounterTestImpl() testImpl {
	return testImpl{
		newInstrument: func(meter otel.Meter, name string) SyncImpler {
			return Must(meter).NewFloat64Counter(name + ".sum")
		},
		getUpdateValue: func() otel.Number {
			for {
				x := rand.Float64()
				if x != 0 {
					return otel.NewFloat64Number(x)
				}
			}
		},
		operate: func(inst interface{}, ctx context.Context, value otel.Number, labels []label.KeyValue) {
			counter := inst.(otel.Float64Counter)
			counter.Add(ctx, value.AsFloat64(), labels...)
		},
		newStore: func() interface{} {
			n := otel.NewFloat64Number(0.0)
			return &n
		},
		storeCollect: func(store interface{}, value otel.Number, _ time.Time) {
			store.(*otel.Number).AddFloat64Atomic(value.AsFloat64())
		},
		storeExpect: func(store interface{}, value otel.Number) {
			store.(*otel.Number).AddFloat64Atomic(value.AsFloat64())
		},
		readStore: func(store interface{}) otel.Number {
			return store.(*otel.Number).AsNumberAtomic()
		},
		equalValues: float64sEqual,
	}
}

func TestStressFloat64Counter(t *testing.T) {
	stressTest(t, floatCounterTestImpl())
}

// LastValue

func intLastValueTestImpl() testImpl {
	return testImpl{
		newInstrument: func(meter otel.Meter, name string) SyncImpler {
			return Must(meter).NewInt64ValueRecorder(name + ".lastvalue")
		},
		getUpdateValue: func() otel.Number {
			r1 := rand.Int63()
			return otel.NewInt64Number(rand.Int63() - r1)
		},
		operate: func(inst interface{}, ctx context.Context, value otel.Number, labels []label.KeyValue) {
			valuerecorder := inst.(otel.Int64ValueRecorder)
			valuerecorder.Record(ctx, value.AsInt64(), labels...)
		},
		newStore: func() interface{} {
			return &lastValueState{
				raw: otel.NewInt64Number(0),
			}
		},
		storeCollect: func(store interface{}, value otel.Number, ts time.Time) {
			gs := store.(*lastValueState)

			if !ts.Before(gs.ts) {
				gs.ts = ts
				gs.raw.SetInt64Atomic(value.AsInt64())
			}
		},
		storeExpect: func(store interface{}, value otel.Number) {
			gs := store.(*lastValueState)
			gs.raw.SetInt64Atomic(value.AsInt64())
		},
		readStore: func(store interface{}) otel.Number {
			gs := store.(*lastValueState)
			return gs.raw.AsNumberAtomic()
		},
		equalValues: int64sEqual,
	}
}

func TestStressInt64LastValue(t *testing.T) {
	stressTest(t, intLastValueTestImpl())
}

func floatLastValueTestImpl() testImpl {
	return testImpl{
		newInstrument: func(meter otel.Meter, name string) SyncImpler {
			return Must(meter).NewFloat64ValueRecorder(name + ".lastvalue")
		},
		getUpdateValue: func() otel.Number {
			return otel.NewFloat64Number((-0.5 + rand.Float64()) * 100000)
		},
		operate: func(inst interface{}, ctx context.Context, value otel.Number, labels []label.KeyValue) {
			valuerecorder := inst.(otel.Float64ValueRecorder)
			valuerecorder.Record(ctx, value.AsFloat64(), labels...)
		},
		newStore: func() interface{} {
			return &lastValueState{
				raw: otel.NewFloat64Number(0),
			}
		},
		storeCollect: func(store interface{}, value otel.Number, ts time.Time) {
			gs := store.(*lastValueState)

			if !ts.Before(gs.ts) {
				gs.ts = ts
				gs.raw.SetFloat64Atomic(value.AsFloat64())
			}
		},
		storeExpect: func(store interface{}, value otel.Number) {
			gs := store.(*lastValueState)
			gs.raw.SetFloat64Atomic(value.AsFloat64())
		},
		readStore: func(store interface{}) otel.Number {
			gs := store.(*lastValueState)
			return gs.raw.AsNumberAtomic()
		},
		equalValues: float64sEqual,
	}
}

func TestStressFloat64LastValue(t *testing.T) {
	stressTest(t, floatLastValueTestImpl())
}
