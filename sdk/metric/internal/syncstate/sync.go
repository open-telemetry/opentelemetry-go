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
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/sdkapi/number"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
)

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
)

var (
	_ sdkapi.Instrument = &instrument{}
)

func (inst *instrument) Implementation() interface{} {
	return inst
}

func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (inst *instrument) initRecord(grp *group, rec *record, attrs attribute.Attributes) {
	rec.group = grp
	rec.attributes = attrs.KeyValues
	rec.collector = inst.cfactory.New(attrs.KeyValues)
}

func (inst *instrument) findOrCreate(grp *group, attrs attribute.Attributes) *record {
	var newRec *record

	for {
		var last *record

		for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {
			// TODO: Fast path: disregard the following and return the
			// first match, HERE.

			if attrs.Equals(attribute.Attributes{
				Fingerprint: grp.fingerprint,
				KeyValues:   rec.attributes,
			}) {
				return rec
			}
			last = rec
		}

		if newRec == nil {
			newRec = &record{}
			inst.initRecord(grp, newRec, attrs)
		}

		if !atomic.CompareAndSwapPointer(&last.next, nil, unsafe.Pointer(newRec)) {
			continue
		}

		return newRec
	}
}

// acquireRecord gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func (inst *instrument) acquireRecord(attrs attribute.Attributes) *record {
	var mk interface{} = attrs.Fingerprint
	if lookup, ok := inst.current.Load(mk); ok {
		// Existing record case.
		grp := lookup.(*group)

		if grp.refMapped.ref() {
			// At this moment it is guaranteed that the
			// group is in the map and will not be removed.
			return inst.findOrCreate(grp, attrs)
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
				return inst.findOrCreate(oldGrp, attrs)
			}
			runtime.Gosched()
			continue
		}
		break
	}

	rec := &newGrp.first
	inst.initRecord(newGrp, rec, attrs)
	return rec
}

func (inst *instrument) Capture(_ context.Context, num number.Number, attrs []attribute.KeyValue) {
	// TODO This is the place to use context, extract baggage.

	r := inst.acquireRecord(attribute.Fingerprint(attrs...))
	defer r.group.refMapped.unref()

	if r.collector == nil {
		// The instrument is disabled.
		return
	}
	if err := aggregator.RangeTest(num, &r.group.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	r.collector.Update(num, &r.group.instrument.descriptor)
	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&r.updateCount, 1)
}

func New() *Accumulator {
	return &Accumulator{}
}

func (a *Accumulator) NewInstrument(descriptor sdkapi.Descriptor, cfactory viewstate.CollectorFactory) (sdkapi.Instrument, error) {
	inst := &instrument{
		descriptor: descriptor,
		cfactory:   cfactory,
	}

	a.instrumentsLock.Lock()
	defer a.instrumentsLock.Unlock()
	a.instruments = append(a.instruments, inst)
	return inst, nil
}

func (a *Accumulator) Collect() {
	a.collectLock.Lock()
	defer a.collectLock.Unlock()

	a.collectInstruments()
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

func (a *Accumulator) checkpointRecord(r *record, final bool) int {
	if r.collector == nil {
		return 0
	}
	if err := r.collector.Send(final); err != nil {
		otel.Handle(err)
		return 0
	}

	return 1
}
