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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

type (
	Accumulator struct {
		instrumentsLock sync.Mutex
		instruments     []*instrument

		callbacksLock sync.Mutex
		callbacks     []*callback

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex
		aggSelector export.AggregatorSelector
	}

	instrument struct {
		descriptor sdkapi.Descriptor
		accum      *Accumulator
		current    sync.Map // map[attribute.Fingerprint]*group
	}

	group struct {
		refMapped   refcountMapped
		fingerprint interface{} // attribute.Fingerprint, current key
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
		live       export.Aggregator
		snapshot   export.Aggregator
		next       unsafe.Pointer
	}

	callback struct {
		function func(context.Context) error
	}
)

var (
	_ sdkapi.MeterImpl  = &Accumulator{}
	_ sdkapi.Instrument = &instrument{}

	// ErrUninitializedInstrument is returned when an instrument is used when uninitialized.
	ErrUninitializedInstrument = fmt.Errorf("use of an uninitialized instrument")
)

func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

func (s *instrument) Implementation() interface{} {
	return s
}

func attributesAreEqual(a, b []attribute.KeyValue) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Key != b[i].Key {
			return false
		}
		if a[i].Value.Type() != b[i].Value.Type() {
			return false
		}
		switch a[i].Value.Type() {
		case attribute.INVALID, attribute.BOOL, attribute.INT64,
			attribute.FLOAT64, attribute.STRING:
			if a[i].Value != b[i].Value {
				return false
			}
		case attribute.BOOLSLICE:
			as := a[i].Value.AsBoolSlice()
			bs := b[i].Value.AsBoolSlice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case attribute.INT64SLICE:
			as := a[i].Value.AsInt64Slice()
			bs := b[i].Value.AsInt64Slice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case attribute.FLOAT64SLICE:
			as := a[i].Value.AsFloat64Slice()
			bs := b[i].Value.AsFloat64Slice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		case attribute.STRINGSLICE:
			as := a[i].Value.AsStringSlice()
			bs := b[i].Value.AsStringSlice()
			if len(as) != len(bs) {
				return false
			}
			for j := 0; j < len(as); j++ {
				if as[j] != bs[j] {
					return false
				}
			}
		}
	}
	return true
}

func (inst *instrument) findOrCreate(grp *group, kvs []attribute.KeyValue) *record {
	var newRec *record

	for {
		var last *record

		for rec := &grp.first; rec != nil; rec = (*record)(atomic.LoadPointer(&rec.next)) {
			// TODO: Fast path: disregard the following and return the
			// first match, HERE.

			if attributesAreEqual(kvs, rec.attributes) {
				return rec
			}
			last = rec
		}

		if newRec == nil {
			newRec = &record{}
			newRec.group = grp
			newRec.attributes = kvs

			inst.accum.aggSelector.AggregatorFor(
				&inst.descriptor, &newRec.live, &newRec.snapshot,
			)
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
func (inst *instrument) acquireHandle(mk interface{}, kvs []attribute.KeyValue) *record {
	if lookup, ok := inst.current.Load(mk); ok {
		// Existing record case.
		grp := lookup.(*group)

		if grp.refMapped.ref() {
			// At this moment it is guaranteed that the
			// group is in the map and will not be removed.
			return inst.findOrCreate(grp, kvs)
		}
		// This group is no longer mapped, try
		// to add a new group below.
	}

	newGrp := &group{
		refMapped: refcountMapped{value: 2},
	}

	for {
		if found, loaded := inst.current.LoadOrStore(mk, newGrp); loaded {
			oldGrp := found.(*group)
			if oldGrp.refMapped.ref() {
				return inst.findOrCreate(oldGrp, kvs)
			}
			runtime.Gosched()
			continue
		}
		break
	}

	rec := &newGrp.first
	rec.group = newGrp
	rec.attributes = kvs

	inst.accum.aggSelector.AggregatorFor(
		&inst.descriptor, &rec.live, &rec.snapshot,
	)
	return rec
}

//
func (s *instrument) RecordOne(ctx context.Context, num number.Number, kvs []attribute.KeyValue) {
	h := s.acquireHandle(attribute.Hash(kvs...), kvs)
	defer h.unbind()
	h.RecordOne(ctx, num)
}

// NewAccumulator constructs a new Accumulator for the given
// processor.  This Accumulator supports only a single processor.
//
// The Accumulator does not start any background process to collect itself
// periodically, this responsibility lies with the processor, typically,
// depending on the type of export.  For example, a pull-based
// processor will call Collect() when it receives a request to scrape
// current metric values.  A push-based processor should configure its
// own periodic collection.
func NewAccumulator(processor export.Processor) *Accumulator {
	return &Accumulator{
		aggSelector: processor,
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

func (m *Accumulator) NewCallback(insts []sdkapi.Instrument, function func(context.Context) error) (sdkapi.Callback, error) {
	cb := &callback{}

	m.callbacksLock.Lock()
	defer m.callbacksLock.Unlock()
	m.callbacks = append(m.callbacsk, cb)
	return cb, nil
}

func (cb *callback) Instruments() []sdkapi.Instrument {
}

// Collect traverses the list of active records and observers and
// exports data for each active instrument.  Collect() may not be
// called concurrently.
//
// During the collection pass, the export.Processor will receive
// one Export() call per current aggregation.
//
// Returns the number of records that were checkpointed.
func (m *Accumulator) Collect(ctx context.Context) int {
	m.collectLock.Lock()
	defer m.collectLock.Unlock()

	m.observeAsyncInstruments(ctx)

	return m.collectSyncInstruments()
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

func (m *Accumulator) collectSyncInstruments() int {
	checkpointed := 0

	m.instrumentsLock.Lock()
	instruments := m.instruments
	m.instrumentsLock.Unlock()

	for _, inst := range instruments {
		inst.current.Range(func(key interface{}, value interface{}) bool {
			grp := value.(*group)
			any := m.checkpointGroup(grp)
			checkpointed += any

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

			checkpointed += m.checkpointGroup(grp)
			return true
		})
	}

	return checkpointed
}

func (m *Accumulator) observeAsyncInstruments(ctx context.Context) {
	m.callbacksLock.Lock()
	callbacks := m.callbacks
	m.callbacksLock.Unlock()

	for _, cb := range callbacks {
		cb.function(ctx)
	}
}

func (m *Accumulator) checkpointRecord(r *record) int {
	if r.live == nil {
		return 0
	}
	err := r.live.SynchronizedMove(r.snapshot, &r.group.instrument.descriptor)
	if err != nil {
		otel.Handle(err)
		return 0
	}

	// @@@
	// a := export.NewAccumulation(&r.group.instrument.descriptor, r.attributes, r.snapshot)
	// err = m.processor.Process(a)
	// if err != nil {
	// 	otel.Handle(err)
	// }
	return 1
}

// RecordOne implements sdkapi.SyncImpl.
func (r *record) RecordOne(ctx context.Context, num number.Number) {
	if r.live == nil {
		// The instrument is disabled according to the AggregatorSelector.
		return
	}
	if err := aggregator.RangeTest(num, &r.group.instrument.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	if err := r.live.Update(ctx, num, &r.group.instrument.descriptor); err != nil {
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

func (m *Accumulator) fromSDK(inst sdkapi.Instrument) *instrument {
	if inst != nil {
		if ii, ok := inst.Implementation().(*instrument); ok {
			return ii
		}
	}
	otel.Handle(ErrUninitializedInstrument)
	return nil
}
