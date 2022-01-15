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
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
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

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex
	}

	instrument struct {
		descriptor sdkapi.Descriptor
		current    sync.Map // map[attribute.Fingerprint]*group
		cfactory   viewstate.CollectorFactory
	}

	group struct {
		refMapped   refcountMapped
		fingerprint uint64
		instrument  *instrument
		first       record
	}

	record struct {
		// updateCount is incremented on every Update.
		updateCount int64

		// collectedCount is set to updateCount on collection,
		// supports checking for no updates during a round.
		collectedCount int64

		group      *group
		attributes []attribute.KeyValue
		collector  viewstate.Collector
		next       unsafe.Pointer
	}

	common struct {
		accumulator *Accumulator
		registry    *registry.State
		views       *viewstate.State
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
	_ apiInstrument.Synchronous = &instrument{}
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

func (a *Accumulator) Int64Instruments(reg *registry.State, views *viewstate.State) syncint64.Instruments {
	return Int64Instruments{
		common: common{
			accumulator: a,
			registry:    reg,
			views:       views,
		},
	}
}

func (a *Accumulator) Float64Instruments(reg *registry.State, views *viewstate.State) syncfloat64.Instruments {
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

func (inst *instrument) Synchronous() {}

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
		func(desc sdkapi.Descriptor) (*instrument, error) {
			cfactory, err := c.views.NewFactory(desc)
			if err != nil {
				return nil, err
			}
			inst := &instrument{
				descriptor: desc,
				cfactory:   cfactory,
			}

			c.accumulator.instrumentsLock.Lock()
			defer c.accumulator.instrumentsLock.Unlock()

			c.accumulator.instruments = append(c.accumulator.instruments, inst)
			return inst, nil
		})
}

func (a *Accumulator) Collect() {
	a.collectLock.Lock()
	defer a.collectLock.Unlock()

	a.collectInstruments()
}

func (a *Accumulator) collectInstruments() {
	a.instrumentsLock.Lock()
	instruments := a.instruments
	a.instrumentsLock.Unlock()

	for _, inst := range instruments {
		inst.current.Range(func(_ interface{}, value interface{}) bool {
			grp := value.(*group)
			any := a.checkpointGroup(grp, false)

			if any != 0 {
				return true
			}
			// Having no updates since last collection, try to unmap:
			if unmapped := grp.refMapped.tryUnmap(); !unmapped {
				// The record is referenced by a binding, continue.
				return true
			}

			// If any other goroutines are now trying to re-insert this
			// entry in the map, they are busy calling Gosched() awaiting
			// this deletion:
			inst.current.Delete(grp.fingerprint)

			// Last we'll see of this.
			_ = a.checkpointGroup(grp, true)
			return true
		})
	}
}

func (a *Accumulator) checkpointGroup(grp *group, final bool) int {
	var checkpointed int
	for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {

		mods := atomic.LoadInt64(&rec.updateCount)
		coll := rec.collectedCount

		if mods != coll {
			// Updates happened in this interval,
			// checkpoint and continue.
			checkpointed += a.checkpointRecord(rec, final)
			rec.collectedCount = mods
		}
	}
	return checkpointed
}

func (a *Accumulator) checkpointRecord(r *record, final bool) int {
	// Note: We could use the `final` bit here to signal to the
	// receiver of this aggregation that it is the last in a
	// sequence and it should feel encouraged to forget its state
	// because a new collector factory will be built to continue
	// this stream (w/ a new *record).
	_ = final

	if r.collector == nil {
		return 0
	}
	if err := r.collector.Send(r.group.instrument.cfactory); err != nil {
		otel.Handle(err)
		return 0
	}

	return 1
}

func capture[N number.Any, Traits traits.Any[N]](_ context.Context, inst *instrument, num N, attrs []attribute.KeyValue) {
	// TODO: Here, this is the place to use context, extract baggage.

	rec, updater := acquireRecord[N](inst, attribute.Fingerprint(attrs...))
	defer rec.group.refMapped.unref()

	if err := aggregator.RangeTest[N, Traits](num, &rec.group.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	updater.Update(num)

	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&rec.updateCount, 1)
}

// acquireRecord gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func acquireRecord[N number.Any](inst *instrument, attrs attribute.Attributes) (*record, viewstate.Updater[N]) {
	var mk interface{} = attrs.Fingerprint
	if lookup, ok := inst.current.Load(mk); ok {
		// Existing record case.
		grp := lookup.(*group)

		if grp.refMapped.ref() {
			// At this moment it is guaranteed that the
			// group is in the map and will not be removed.
			return findOrCreate[N](inst, grp, attrs)
		}
		// This group is no longer mapped, try
		// to add a new group below.
	}

	newGrp := &group{
		refMapped:   refcountMapped{value: 2},
		instrument:  inst,
		fingerprint: attrs.Fingerprint,
	}

	for {
		if found, loaded := inst.current.LoadOrStore(mk, newGrp); loaded {
			oldGrp := found.(*group)
			if oldGrp.refMapped.ref() {
				return findOrCreate[N](inst, oldGrp, attrs)
			}
			runtime.Gosched()
			continue
		}
		break
	}

	rec := &newGrp.first
	return rec, initRecord[N](inst, newGrp, rec, attrs)
}

func initRecord[N number.Any](inst *instrument, grp *group, rec *record, attrs attribute.Attributes) viewstate.Updater[N] {
	rec.group = grp
	rec.attributes = attrs.KeyValues
	rec.collector = inst.cfactory.New(attrs.KeyValues, &inst.descriptor)

	// This conversion must be safe or else there is a bug.
	return rec.collector.(viewstate.Updater[N])
}

func findOrCreate[N number.Any](inst *instrument, grp *group, attrs attribute.Attributes) (*record, viewstate.Updater[N]) {
	var newRec *record

	for {
		var last *record

		for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {
			// TODO: Here an even-faster path option:
			// disregard the equality test and return the
			// first match.
			if attrs.Equals(attribute.Attributes{
				Fingerprint: grp.fingerprint,
				KeyValues:   rec.attributes,
			}) {
				return rec, rec.collector.(viewstate.Updater[N])
			}
			last = rec
		}

		if newRec == nil {
			newRec = &record{}
			_ = initRecord[N](inst, grp, newRec, attrs)
		}

		if !atomic.CompareAndSwapPointer(&last.next, nil, unsafe.Pointer(newRec)) {
			continue
		}

		return newRec, newRec.collector.(viewstate.Updater[N])
	}
}
