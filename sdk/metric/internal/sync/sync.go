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

package sync // import "go.opentelemetry.io/otel/sdk/metric/internal/sync"

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/internal/number"
	"go.opentelemetry.io/otel/sdk/metric/internal/sdkapi"
)

type (
	Accumulator struct {
		instrumentsLock sync.Mutex
		instruments     []*instrument

		// collectLock prevents simultaneous calls to Collect().
		collectLock       sync.Mutex
		collectorSelector export.CollectorSelector
	}

	instrument struct {
		descriptor sdkapi.Descriptor
		accum      *Accumulator
		current    sync.Map // map[attribute.Fingerprint]*group
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
		collector  export.Collector
		next       unsafe.Pointer
	}
)

var (
	_ sdkapi.Instrument = &instrument{}
)

func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (accum *Accumulator) initRecord(inst *instrument, grp *group, rec *record, attrs attribute.Attributes) {
	rec.group = grp
	rec.attributes = attrs.KeyValues
	rec.collector = accum.collectorSelector.CollectorFor(&inst.descriptor)
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
			inst.accum.initRecord(inst, grp, newRec, attrs)
		}

		if !atomic.CompareAndSwapPointer(&last.next, nil, unsafe.Pointer(newRec)) {
			continue
		}

		return newRec
	}
}

// acquireHandle gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func (inst *instrument) acquireHandle(attrs attribute.Attributes) *record {
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
	inst.accum.initRecord(inst, newGrp, rec, attrs)
	return rec
}

//
func (s *instrument) Capture(ctx context.Context, num number.Number, attrs []attribute.KeyValue) {
	h := s.acquireHandle(attribute.Fingerprint(attrs...))
	defer h.unbind()
	h.record(ctx, num)
}

func newSyncAccumulator(cs export.CollectorSelector) *Accumulator {
	return &Accumulator{
		collectorSelector: cs,
	}
}

// NewInstrument implements sdkapi.MetricImpl.
func (m *Accumulator) NewInstrument(descriptor sdkapi.Descriptor) (sdkapi.Instrument, error) {
	inst := &instrument{
		descriptor: descriptor,
		accum:      m,
	}

	m.instrumentsLock.Lock()
	defer m.instrumentsLock.Unlock()
	m.instruments = append(m.instruments, inst)
	return inst, nil
}

func (m *Accumulator) Collect() {
	m.collectLock.Lock()
	defer m.collectLock.Unlock()

	m.collectInstruments()
}

func (m *Accumulator) checkpointGroup(grp *group) int {
	var checkpointed int
	for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {

		mods := atomic.LoadInt64(&rec.updateCount)
		coll := rec.collectedCount

		if mods != coll {
			// Updates happened in this interval,
			// checkpoint and continue.
			checkpointed += m.checkpointRecord(rec)
			rec.collectedCount = mods
		}
	}
	return checkpointed
}

func (m *Accumulator) collectInstruments() {
	m.instrumentsLock.Lock()
	instruments := m.instruments
	m.instrumentsLock.Unlock()

	for _, inst := range instruments {
		inst.current.Range(func(_ interface{}, value interface{}) bool {
			grp := value.(*group)
			any := m.checkpointGroup(grp)

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

			return true
		})
	}
}

func (m *Accumulator) checkpointRecord(r *record) int {
	if r.collector == nil {
		return 0
	}
	if err := r.collector.Send(); err != nil {
		otel.Handle(err)
		return 0
	}

	return 1
}

func (r *record) record(ctx context.Context, num number.Number) {
	if r.collector == nil {
		// The instrument is disabled.
		return
	}
	if err := aggregator.RangeTest(num, &r.group.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	if err := r.collector.Update(ctx, num, &r.group.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&r.updateCount, 1)
}

func (r *record) unbind() {
	r.group.refMapped.unref()
}
