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

package syncstate

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

// Performance note: there is still 1 obligatory allocation in the
// fast path of this code due to the sync.Map key.  Assuming Go will
// give us a generic form of sync.Map some time soon, the allocation
// cost of instrument.Current will be reduced to zero allocs in the
// fast path.  See also https://github.com/a8m/syncmap.

type (
	Accumulator struct {
		instrumentsLock sync.Mutex
		instruments     []*instrument
	}

	instrument struct {
		apiInstrument.Synchronous
		
		descriptor sdkapi.Descriptor
		current    sync.Map // map[attribute.Set]*record
		compiled   viewstate.Instrument
	}

	record struct {
		refMapped   refcountMapped
		instrument  *instrument

		// updateCount is incremented on every Update.
		updateCount int64

		// collectedCount is set to updateCount on collection,
		// supports checking for no updates during a round.
		collectedCount int64

		distinct   attribute.Set
		attributes []attribute.KeyValue
		collector  viewstate.Collector
		next       unsafe.Pointer
	}

	common struct {
		accumulator *Accumulator
		registry    *registry.State
		views       *viewstate.Compiler
	}


	Int64Instruments struct { common }
	Float64Instruments struct { common }

	counter[N number.Any, Traits traits.Any[N]] struct {
		*instrument
	}

	histogram[N number.Any, Traits traits.Any[N]] struct {
		*instrument
	}
)

var (
	_ syncint64.Counter         = counter[int64, traits.Int64]{}
	_ syncint64.UpDownCounter   = counter[int64, traits.Int64]{}
	_ syncint64.Histogram       = histogram[int64, traits.Int64]{}
	_ syncfloat64.Counter       = counter[float64, traits.Float64]{}
	_ syncfloat64.UpDownCounter = counter[float64, traits.Float64]{}
	_ syncfloat64.Histogram     = histogram[float64, traits.Float64]{}
)

func New() *Accumulator {
	return &Accumulator{}
}

func (a *Accumulator) Int64Instruments(reg *registry.State, views *viewstate.Compiler) syncint64.InstrumentProvider {
	return Int64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (a *Accumulator) Float64Instruments(reg *registry.State, views *viewstate.Compiler) syncfloat64.InstrumentProvider {
	return Float64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (i Int64Instruments) Counter(name string, opts ...apiInstrument.Option) (syncint64.Counter, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.CounterInstrumentKind)
	return counter[int64, traits.Int64]{instrument: inst}, err
}

func (i Int64Instruments) UpDownCounter(name string, opts ...apiInstrument.Option) (syncint64.UpDownCounter, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.UpDownCounterInstrumentKind)
	return counter[int64, traits.Int64]{instrument: inst}, err
}

func (i Int64Instruments) Histogram(name string, opts ...apiInstrument.Option) (syncint64.Histogram, error) {
	inst, err := i.newInstrument(name, opts, number.Int64Kind, sdkapi.HistogramInstrumentKind)
	return histogram[int64, traits.Int64]{instrument: inst}, err
}

func (f Float64Instruments) Counter(name string, opts ...apiInstrument.Option) (syncfloat64.Counter, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.CounterInstrumentKind)
	return counter[float64, traits.Float64]{instrument: inst}, err
}

func (f Float64Instruments) UpDownCounter(name string, opts ...apiInstrument.Option) (syncfloat64.UpDownCounter, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.UpDownCounterInstrumentKind)
	return counter[float64, traits.Float64]{instrument: inst}, err
}

func (f Float64Instruments) Histogram(name string, opts ...apiInstrument.Option) (syncfloat64.Histogram, error) {
	inst, err := f.newInstrument(name, opts, number.Float64Kind, sdkapi.HistogramInstrumentKind)
	return histogram[float64, traits.Float64]{instrument: inst}, err
}

// implements registry.hasDescriptor
func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (c counter[N, Traits]) Add(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	if c.instrument != nil {
		capture[N, Traits](ctx, c.instrument, incr, attrs)
	}
}

func (h histogram[N, Traits]) Record(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	if h.instrument != nil {
		capture[N, Traits](ctx, h.instrument, incr, attrs)
	}
}

func (c common) newInstrument(name string, opts []apiInstrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) (*instrument, error) {
	return registry.Lookup(
		c.registry,
		name, opts, nk, ik,
		func(desc sdkapi.Descriptor) *instrument{
			compiled := c.views.Compile(desc)
			inst := &instrument{
				descriptor: desc,
				compiled:   compiled,
			}

			c.accumulator.instrumentsLock.Lock()
			defer c.accumulator.instrumentsLock.Unlock()

			c.accumulator.instruments = append(c.accumulator.instruments, inst)
			return inst
		})
}

func (a *Accumulator) Collect(reader *reader.Reader, instruments *[]reader.Instrument) {
	a.instrumentsLock.Lock()
	defer a.instrumentsLock.Unlock()

	for _, inst := range a.instruments {
		inst.current.Range(func(key interface{}, value interface{}) bool {
			rec := value.(*record)
			any := a.collectRecord(rec, false)

			if any != 0 {
				return true
			}
			// Having no updates since last collection, try to unmap:
			if unmapped := rec.refMapped.tryUnmap(); !unmapped {
				// The record is referenced by a binding, continue.
				return true
			}

			// If any other goroutines are now trying to re-insert this
			// entry in the map, they are busy calling Gosched() awaiting
			// this deletion:
			inst.current.Delete(key)

			// Last we'll see of this.
			_ = a.collectRecord(rec, true)
			return true
		})
	}
}

func (a *Accumulator) collectRecord(rec *record, final bool) int {
	mods := atomic.LoadInt64(&rec.updateCount)
	coll := rec.collectedCount

	if mods == coll {
		return 0
	}
	// Updates happened in this interval,
	// collect and continue.
	rec.collectedCount = mods

	// Note: We could use the `final` bit here to signal to the
	// receiver of this aggregation that it is the last in a
	// sequence and it should feel encouraged to forget its state
	// because a new collector factory will be built to continue
	// this stream (w/ a new *record).
	_ = final

	if rec.collector == nil {
		return 0
	}
	rec.collector.Collect()
	return 1
}

func capture[N number.Any, Traits traits.Any[N]](_ context.Context, inst *instrument, num N, attrs []attribute.KeyValue) {
	// TODO: Here, this is the place to use context, extract baggage.

	rec, updater := acquireRecord[N](inst, attrs)
	defer rec.refMapped.unref()

	if err := aggregator.RangeTest[N, Traits](num, &rec.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	updater.(viewstate.CollectorUpdater[N]).Update(num)

	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&rec.updateCount, 1)
}

// acquireRecord gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func acquireRecord[N number.Any](inst *instrument, attrs []attribute.KeyValue) (*record, viewstate.Updater[N]) {
	aset := attribute.NewSet(attrs...)
	if lookup, ok := inst.current.Load(aset); ok {
		// Existing record case.
		rec := lookup.(*record)

		if rec.refMapped.ref() {
			// At this moment it is guaranteed that the
			// record is in the map and will not be removed.
			return rec, rec.collector.(viewstate.Updater[N])
		}
		// This group is no longer mapped, try
		// to add a new group below.
	}

	newRec := &record{
		refMapped:   refcountMapped{value: 2},
		instrument:  inst,
		distinct: aset,
		attributes: attrs,
	}

	for {
		if found, loaded := inst.current.LoadOrStore(aset, newRec); loaded {
			oldRec := found.(*record)
			if oldRec.refMapped.ref() {
				return oldRec, oldRec.collector.(viewstate.Updater[N])
			}
			runtime.Gosched()
			continue
		}
		break
	}


	return newRec, initRecord[N](inst, newRec, attrs)
}

func initRecord[N number.Any](inst *instrument, rec *record, attrs []attribute.KeyValue) viewstate.Updater[N] {
	rec.collector = inst.compiled.NewCollector(attrs, nil)
	return rec.collector.(viewstate.Updater[N])
}
