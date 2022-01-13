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
	_ "runtime"
	"sync"
	_ "sync/atomic"
	"unsafe"

	_ "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	_ "go.opentelemetry.io/otel/metric"
	apiInstrument "go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"

	_ "go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
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

	syncInstruments[
		N number.Any,
		// Insts anyInstruments[C, U, H],
		C anyCounter[N],
		U, H any] struct {
		*Accumulator
	}

	// int64Instruments struct {}
	// float64Instruments struct {}

	// anyInstruments[C, U, H any] interface {
	// 	Counter(name string, opts ...apiInstrument.Option) (C, error)
	// 	UpDownCounter(name string, opts ...apiInstrument.Option) (U, error)
	// 	Histogram(name string, opts ...apiInstrument.Option) (H, error)
	// }

	// anyInstruments[C, U, H any] interface {
	// 	Counter(name string, opts ...apiInstrument.Option) (C, error)
	// 	UpDownCounter(name string, opts ...apiInstrument.Option) (U, error)
	// 	Histogram(name string, opts ...apiInstrument.Option) (H, error)
	// }

	// int64Instruments struct {
	// 	Counter(name string, opts ...apiInstrument.Option) (syncint64.Counter, error)
	// 	UpDownCounter(name string, opts ...apiInstrument.Option) (syncint64.UpDownCounter, error)
	// 	Histogram(name string, opts ...apiInstrument.Option) (syncint64.Histogram, error)
	// }

	anyCounter[N number.Any] interface {
		Add(ctx context.Context, incr N, attrs ...attribute.KeyValue)

		Synchronous()
	}

	// anyHistogram[N number.Any] interface {
	// 	Record(ctx context.Context, incr N, attrs ...attribute.KeyValue)

	// 	Synchronous()
	// }

	counter[N number.Any] struct {
		// inst *instrument
	}
)

// func (int64Instruments) Counter(c *counter[int64]) syncint64.Counter {
// 	return c
// }
// func (int64Instruments) UpDownCounter(c *counter[int64]) syncint64.UpDownCounter {
// 	return c
// }
// func (int64Instruments) Histogram(*counter[int64]) syncint64.Histogram {
// 	return nil
// }

// func (float64Instruments) Counter(c *counter[float64]) syncfloat64.Counter {
// 	return c
// }
// func (float64Instruments) UpDownCounter(c *counter[float64]) syncfloat64.UpDownCounter {
// 	return c
// }
// func (float64Instruments) Histogram(*counter[float64]) syncfloat64.Histogram {
// 	return nil
// }

var (
	_ apiInstrument.Synchronous = &instrument{}
	_ anyCounter[int64] = &counter[int64]{}
	_ anyCounter[float64] = &counter[float64]{}
	_ syncint64.Counter = &counter[int64]{}
	_ syncint64.UpDownCounter = &counter[int64]{}
)

func (inst *instrument) Synchronous() {}

func New() *Accumulator {
	return &Accumulator{}
}

func (a *Accumulator) SyncInt64() syncint64.Instruments {
	return syncInstruments[
		int64,
		syncint64.Counter,
		syncint64.UpDownCounter,
		syncint64.Histogram,
	]{
		Accumulator: a,
	}
}

func (a *Accumulator) SyncFloat64() syncfloat64.Instruments {
	return syncInstruments[
		float64,
		syncfloat64.Counter,
		syncfloat64.UpDownCounter,
		syncfloat64.Histogram,
	]{
		Accumulator: a,
	}
}

func (c counter[N]) Add(ctx context.Context, incr N, attrs ...attribute.KeyValue) {
	// @@@
}

func (c counter[N]) Synchronous() {}

var _ syncint64.Counter = &counter[int64]{}

func (a syncInstruments[N, C, U, H]) Counter(name string, opts ...apiInstrument.Option) (C, error) {
	// c := &counter[N]{}
	var c C
	c = counter[N]{}
	return c, nil
}

func (a syncInstruments[N, C, U, H]) UpDownCounter(name string, opts ...apiInstrument.Option) (U, error) {
	//return counter[N]{a}
	var u U
	return u, nil
}

func (a syncInstruments[N, C, U, H]) Histogram(name string, opts ...apiInstrument.Option) (H, error) {
	//return counter[N]{a}
	var h H
	return h, nil
}

// func (a *Accumulator) NewInstrument(descriptor sdkapi.Descriptor, cfactory viewstate.CollectorFactory) (sdkapi.Instrument, error) {
// 	inst := &instrument{
// 		descriptor: descriptor,
// 		cfactory:   cfactory,
// 	}

// 	a.instrumentsLock.Lock()
// 	defer a.instrumentsLock.Unlock()
// 	a.instruments = append(a.instruments, inst)
// 	return inst, nil
// }

// func (a *Accumulator) Collect() {
// 	a.collectLock.Lock()
// 	defer a.collectLock.Unlock()

// 	a.collectInstruments()
// }

// func (a *Accumulator) collectInstruments() {
// 	a.instrumentsLock.Lock()
// 	instruments := a.instruments
// 	a.instrumentsLock.Unlock()

// 	for _, inst := range instruments {
// 		inst.current.Range(func(_ interface{}, value interface{}) bool {
// 			grp := value.(*group)
// 			any := a.checkpointGroup(grp, false)

// 			if any != 0 {
// 				return true
// 			}
// 			// Having no updates since last collection, try to unmap:
// 			if unmapped := grp.refMapped.tryUnmap(); !unmapped {
// 				// The record is referenced by a binding, continue.
// 				return true
// 			}

// 			// If any other goroutines are now trying to re-insert this
// 			// entry in the map, they are busy calling Gosched() awaiting
// 			// this deletion:
// 			inst.current.Delete(grp.fingerprint)

// 			// Last we'll see of this.
// 			_ = a.checkpointGroup(grp, true)
// 			return true
// 		})
// 	}
// }

// func (a *Accumulator) checkpointGroup(grp *group, final bool) int {
// 	var checkpointed int
// 	for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {

// 		mods := atomic.LoadInt64(&rec.updateCount)
// 		coll := rec.collectedCount

// 		if mods != coll {
// 			// Updates happened in this interval,
// 			// checkpoint and continue.
// 			checkpointed += a.checkpointRecord(rec, final)
// 			rec.collectedCount = mods
// 		}
// 	}
// 	return checkpointed
// }

// func (a *Accumulator) checkpointRecord(r *record, final bool) int {
// 	if r.collector == nil {
// 		return 0
// 	}
// 	if err := r.collector.Send(final); err != nil {
// 		otel.Handle(err)
// 		return 0
// 	}

// 	return 1
// }

// func (inst *instrument) Capture(_ context.Context, num number.Number, attrs []attribute.KeyValue) {
// 	// TODO: Here, this is the place to use context, extract baggage.

// 	r := inst.acquireRecord(attribute.Fingerprint(attrs...))
// 	defer r.group.refMapped.unref()

// 	if r.collector == nil {
// 		// The instrument is disabled.
// 		return
// 	}
// 	if err := aggregator.RangeTest(num, &r.group.instrument.descriptor); err != nil {
// 		otel.Handle(err)
// 		return
// 	}
// 	r.collector.Update(num, &r.group.instrument.descriptor)
// 	// Record was modified, inform the Collect() that things need
// 	// to be collected while the record is still mapped.
// 	atomic.AddInt64(&r.updateCount, 1)
// }

// // acquireRecord gets or creates a `*record` corresponding to `kvs`,
// // the input labels.  The second argument `labels` is passed in to
// // support re-use of the orderedLabels computed by a previous
// // measurement in the same batch.   This performs two allocations
// // in the common case.
// func (inst *instrument) acquireRecord(attrs attribute.Attributes) *record {
// 	var mk interface{} = attrs.Fingerprint
// 	if lookup, ok := inst.current.Load(mk); ok {
// 		// Existing record case.
// 		grp := lookup.(*group)

// 		if grp.refMapped.ref() {
// 			// At this moment it is guaranteed that the
// 			// group is in the map and will not be removed.
// 			return inst.findOrCreate(grp, attrs)
// 		}
// 		// This group is no longer mapped, try
// 		// to add a new group below.
// 	}

// 	newGrp := &group{
// 		refMapped:   refcountMapped{value: 2},
// 		instrument:  inst,
// 		fingerprint: attrs.Fingerprint,
// 	}

// 	for {
// 		if found, loaded := inst.current.LoadOrStore(mk, newGrp); loaded {
// 			oldGrp := found.(*group)
// 			if oldGrp.refMapped.ref() {
// 				return inst.findOrCreate(oldGrp, attrs)
// 			}
// 			runtime.Gosched()
// 			continue
// 		}
// 		break
// 	}

// 	rec := &newGrp.first
// 	inst.initRecord(newGrp, rec, attrs)
// 	return rec
// }

// func (inst *instrument) initRecord(grp *group, rec *record, attrs attribute.Attributes) {
// 	rec.group = grp
// 	rec.attributes = attrs.KeyValues
// 	rec.collector = inst.cfactory.New(attrs.KeyValues)
// }

// func (inst *instrument) findOrCreate(grp *group, attrs attribute.Attributes) *record {
// 	var newRec *record

// 	for {
// 		var last *record

// 		for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {
// 			// TODO: Here an even-faster path option:
// 			// disregard the equality test and return the
// 			// first match.
// 			if attrs.Equals(attribute.Attributes{
// 				Fingerprint: grp.fingerprint,
// 				KeyValues:   rec.attributes,
// 			}) {
// 				return rec
// 			}
// 			last = rec
// 		}

// 		if newRec == nil {
// 			newRec = &record{}
// 			inst.initRecord(grp, newRec, attrs)
// 		}

// 		if !atomic.CompareAndSwapPointer(&last.next, nil, unsafe.Pointer(newRec)) {
// 			continue
// 		}

// 		return newRec
// 	}
// }
