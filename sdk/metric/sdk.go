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

package metric

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	api "go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// SDK implements the OpenTelemetry Meter API.  The SDK is
	// bound to a single export.Batcher in `New()`.
	//
	// The SDK supports a Collect() API to gather and export
	// current data.  Collect() should be arranged according to
	// the batcher model.  Push-based batchers will setup a
	// timer to call Collect() periodically.  Pull-based batchers
	// will call Collect() when a pull request arrives.
	SDK struct {
		// current maps `mapkey` to *record.
		current sync.Map

		// asyncInstruments is a set of
		// `*asyncInstrument` instances
		asyncInstruments sync.Map

		// currentEpoch is the current epoch number. It is
		// incremented in `Collect()`.
		currentEpoch int64

		// batcher is the configured batcher+configuration.
		batcher export.Batcher

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex

		// errorHandler supports delivering errors to the user.
		errorHandler ErrorHandler

		// asyncSortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation.  It is cleared after use.
		asyncSortSlice label.Sortable
	}

	syncInstrument struct {
		instrument
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its labels.
	mapkey struct {
		descriptor *metric.Descriptor
		ordered    label.Distinct
	}

	// record maintains the state of one metric instrument.  Due
	// the use of lock-free algorithms, there may be more than one
	// `record` in existence at a time, although at most one can
	// be referenced from the `SDK.current` map.
	record struct {
		// refMapped keeps track of refcounts and the mapping state to the
		// SDK.current map.
		refMapped refcountMapped

		// updateCount is incremented on every Update.
		updateCount int64

		// collectedCount is set to updateCount on collection,
		// supports checking for no updates during a round.
		collectedCount int64

		// storage is the stored label set for this record,
		// except in cases where a label set is shared due to
		// batch recording.
		storage label.Set

		// labels is the processed label set for this record.
		// this may refer to the `storage` field in another
		// record if this label set is shared resulting from
		// `RecordBatch`.
		labels *label.Set

		// sortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation.
		sortSlice label.Sortable

		// inst is a pointer to the corresponding instrument.
		inst *syncInstrument

		// recorder implements the actual RecordOne() API,
		// depending on the type of aggregation.  If nil, the
		// metric was disabled by the exporter.
		recorder export.Aggregator
	}

	instrument struct {
		meter      *SDK
		descriptor metric.Descriptor
	}

	asyncInstrument struct {
		instrument
		// recorders maps ordered labels to the pair of
		// labelset and recorder
		recorders map[label.Distinct]*labeledRecorder

		callback func(func(core.Number, []core.KeyValue))
	}

	labeledRecorder struct {
		observedEpoch int64
		labels        *label.Set
		recorder      export.Aggregator
	}

	ErrorHandler func(error)
)

var (
	_ api.MeterImpl     = &SDK{}
	_ api.AsyncImpl     = &asyncInstrument{}
	_ api.SyncImpl      = &syncInstrument{}
	_ api.BoundSyncImpl = &record{}
)

func (inst *instrument) Descriptor() api.Descriptor {
	return inst.descriptor
}

func (a *asyncInstrument) Implementation() interface{} {
	return a
}

func (s *syncInstrument) Implementation() interface{} {
	return s
}

func (a *asyncInstrument) observe(number core.Number, labels []core.KeyValue) {
	if err := aggregator.RangeTest(number, &a.descriptor); err != nil {
		a.meter.errorHandler(err)
		return
	}
	recorder := a.getRecorder(labels)
	if recorder == nil {
		// The instrument is disabled according to the
		// AggregationSelector.
		return
	}
	if err := recorder.Update(context.Background(), number, &a.descriptor); err != nil {
		a.meter.errorHandler(err)
		return
	}
}

func (a *asyncInstrument) getRecorder(kvs []core.KeyValue) export.Aggregator {
	// We are in a single-threaded context.  Note: this assumption
	// could be violated if the user added concurrency within
	// their callback.
	labels := label.NewSetWithSortable(kvs, &a.meter.asyncSortSlice)

	lrec, ok := a.recorders[labels.Equivalent()]
	if ok {
		if lrec.observedEpoch == a.meter.currentEpoch {
			// last value wins for Observers, so if we see the same labels
			// in the current epoch, we replace the old recorder
			lrec.recorder = a.meter.batcher.AggregatorFor(&a.descriptor)
		} else {
			lrec.observedEpoch = a.meter.currentEpoch
		}
		a.recorders[labels.Equivalent()] = lrec
		return lrec.recorder
	}
	rec := a.meter.batcher.AggregatorFor(&a.descriptor)
	if a.recorders == nil {
		a.recorders = make(map[label.Distinct]*labeledRecorder)
	}
	// This may store nil recorder in the map, thus disabling the
	// asyncInstrument for the labelset for good. This is intentional,
	// but will be revisited later.
	a.recorders[labels.Equivalent()] = &labeledRecorder{
		recorder:      rec,
		labels:        &labels,
		observedEpoch: a.meter.currentEpoch,
	}
	return rec
}

func (m *SDK) SetErrorHandler(f ErrorHandler) {
	m.errorHandler = f
}

// acquireHandle gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func (s *syncInstrument) acquireHandle(kvs []core.KeyValue, labelPtr *label.Set) *record {
	var rec *record
	var equiv label.Distinct

	if labelPtr == nil {
		// This memory allocation may not be used, but it's
		// needed for the `sortSlice` field, to avoid an
		// allocation while sorting.
		rec = &record{}
		rec.storage = label.NewSetWithSortable(kvs, &rec.sortSlice)
		rec.labels = &rec.storage
		equiv = rec.storage.Equivalent()
	} else {
		equiv = labelPtr.Equivalent()
	}

	// Create lookup key for sync.Map (one allocation, as this
	// passes through an interface{})
	mk := mapkey{
		descriptor: &s.descriptor,
		ordered:    equiv,
	}

	if actual, ok := s.meter.current.Load(mk); ok {
		// Existing record case.
		existingRec := actual.(*record)
		if existingRec.refMapped.ref() {
			// At this moment it is guaranteed that the entry is in
			// the map and will not be removed.
			return existingRec
		}
		// This entry is no longer mapped, try to add a new entry.
	}

	if rec == nil {
		rec = &record{}
		rec.labels = labelPtr
	}
	rec.refMapped = refcountMapped{value: 2}
	rec.inst = s
	rec.recorder = s.meter.batcher.AggregatorFor(&s.descriptor)

	for {
		// Load/Store: there's a memory allocation to place `mk` into
		// an interface here.
		if actual, loaded := s.meter.current.LoadOrStore(mk, rec); loaded {
			// Existing record case. Cannot change rec here because if fail
			// will try to add rec again to avoid new allocations.
			oldRec := actual.(*record)
			if oldRec.refMapped.ref() {
				// At this moment it is guaranteed that the entry is in
				// the map and will not be removed.
				return oldRec
			}
			// This loaded entry is marked as unmapped (so Collect will remove
			// it from the map immediately), try again - this is a busy waiting
			// strategy to wait until Collect() removes this entry from the map.
			//
			// This can be improved by having a list of "Unmapped" entries for
			// one time only usages, OR we can make this a blocking path and use
			// a Mutex that protects the delete operation (delete only if the old
			// record is associated with the key).

			// Let collector get work done to remove the entry from the map.
			runtime.Gosched()
			continue
		}
		// The new entry was added to the map, good to go.
		return rec
	}
}

func (s *syncInstrument) Bind(kvs []core.KeyValue) api.BoundSyncImpl {
	return s.acquireHandle(kvs, nil)
}

func (s *syncInstrument) RecordOne(ctx context.Context, number core.Number, kvs []core.KeyValue) {
	h := s.acquireHandle(kvs, nil)
	defer h.Unbind()
	h.RecordOne(ctx, number)
}

// New constructs a new SDK for the given batcher.  This SDK supports
// only a single batcher.
//
// The SDK does not start any background process to collect itself
// periodically, this responsbility lies with the batcher, typically,
// depending on the type of export.  For example, a pull-based
// batcher will call Collect() when it receives a request to scrape
// current metric values.  A push-based batcher should configure its
// own periodic collection.
func New(batcher export.Batcher, opts ...Option) *SDK {
	c := &Config{ErrorHandler: DefaultErrorHandler}
	for _, opt := range opts {
		opt.Apply(c)
	}

	return &SDK{
		batcher:      batcher,
		errorHandler: c.ErrorHandler,
	}
}

func DefaultErrorHandler(err error) {
	fmt.Fprintln(os.Stderr, "Metrics SDK error:", err)
}

func (m *SDK) NewSyncInstrument(descriptor api.Descriptor) (api.SyncImpl, error) {
	return &syncInstrument{
		instrument: instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *SDK) NewAsyncInstrument(descriptor api.Descriptor, callback func(func(core.Number, []core.KeyValue))) (api.AsyncImpl, error) {
	a := &asyncInstrument{
		instrument: instrument{
			descriptor: descriptor,
			meter:      m,
		},
		callback: callback,
	}
	m.asyncInstruments.Store(a, nil)
	return a, nil
}

// Collect traverses the list of active records and observers and
// exports data for each active instrument.  Collect() may not be
// called concurrently.
//
// During the collection pass, the export.Batcher will receive
// one Export() call per current aggregation.
//
// Returns the number of records that were checkpointed.
func (m *SDK) Collect(ctx context.Context) int {
	m.collectLock.Lock()
	defer m.collectLock.Unlock()

	checkpointed := m.collectRecords(ctx)
	checkpointed += m.collectAsync(ctx)
	m.currentEpoch++
	return checkpointed
}

func (m *SDK) collectRecords(ctx context.Context) int {
	checkpointed := 0

	m.current.Range(func(key interface{}, value interface{}) bool {
		// Note: always continue to iterate over the entire
		// map by returning `true` in this function.
		inuse := value.(*record)

		mods := atomic.LoadInt64(&inuse.updateCount)
		coll := inuse.collectedCount

		if mods != coll {
			// Updates happened in this interval,
			// checkpoint and continue.
			checkpointed += m.checkpointRecord(ctx, inuse)
			inuse.collectedCount = mods
			return true
		}

		// Having no updates since last collection, try to unmap:
		if unmapped := inuse.refMapped.tryUnmap(); !unmapped {
			// The record is referenced by a binding, continue.
			return true
		}

		// If any other goroutines are now trying to re-insert this
		// entry in the map, they are busy calling Gosched() awaiting
		// this deletion:
		m.current.Delete(inuse.mapkey())

		// There's a potential race between `LoadInt64` and
		// `tryUnmap` in this function.  Since this is the
		// last we'll see of this record, checkpoint
		mods = atomic.LoadInt64(&inuse.updateCount)
		if mods != coll {
			checkpointed += m.checkpointRecord(ctx, inuse)
		}
		return true
	})

	return checkpointed
}

func (m *SDK) collectAsync(ctx context.Context) int {
	checkpointed := 0

	m.asyncInstruments.Range(func(key, value interface{}) bool {
		a := key.(*asyncInstrument)
		a.callback(a.observe)
		checkpointed += m.checkpointAsync(ctx, a)
		return true
	})

	return checkpointed
}

func (m *SDK) checkpointRecord(ctx context.Context, r *record) int {
	return m.checkpoint(ctx, &r.inst.descriptor, r.recorder, r.labels)
}

func (m *SDK) checkpointAsync(ctx context.Context, a *asyncInstrument) int {
	if len(a.recorders) == 0 {
		return 0
	}
	checkpointed := 0
	for encodedLabels, lrec := range a.recorders {
		lrec := lrec
		epochDiff := m.currentEpoch - lrec.observedEpoch
		if epochDiff == 0 {
			checkpointed += m.checkpoint(ctx, &a.descriptor, lrec.recorder, lrec.labels)
		} else if epochDiff > 1 {
			// This is second collection cycle with no
			// observations for this labelset. Remove the
			// recorder.
			delete(a.recorders, encodedLabels)
		}
	}
	if len(a.recorders) == 0 {
		a.recorders = nil
	}
	return checkpointed
}

func (m *SDK) checkpoint(ctx context.Context, descriptor *metric.Descriptor, recorder export.Aggregator, labels *label.Set) int {
	if recorder == nil {
		return 0
	}
	recorder.Checkpoint(ctx, descriptor)

	exportRecord := export.NewRecord(descriptor, labels, recorder)
	err := m.batcher.Process(ctx, exportRecord)
	if err != nil {
		m.errorHandler(err)
	}
	return 1
}

// RecordBatch enters a batch of metric events.
func (m *SDK) RecordBatch(ctx context.Context, kvs []core.KeyValue, measurements ...api.Measurement) {
	// Labels will be computed the first time acquireHandle is
	// called.  Subsequent calls to acquireHandle will re-use the
	// previously computed value instead of recomputing the
	// ordered labels.
	var labelsPtr *label.Set
	for i, meas := range measurements {
		s := meas.SyncImpl().(*syncInstrument)

		h := s.acquireHandle(kvs, labelsPtr)

		// Re-use labels for the next measurement.
		if i == 0 {
			labelsPtr = h.labels
		}

		defer h.Unbind()
		h.RecordOne(ctx, meas.Number())
	}
}

func (r *record) RecordOne(ctx context.Context, number core.Number) {
	if r.recorder == nil {
		// The instrument is disabled according to the AggregationSelector.
		return
	}
	if err := aggregator.RangeTest(number, &r.inst.descriptor); err != nil {
		r.inst.meter.errorHandler(err)
		return
	}
	if err := r.recorder.Update(ctx, number, &r.inst.descriptor); err != nil {
		r.inst.meter.errorHandler(err)
		return
	}
	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&r.updateCount, 1)
}

func (r *record) Unbind() {
	r.refMapped.unref()
}

func (r *record) mapkey() mapkey {
	return mapkey{
		descriptor: &r.inst.descriptor,
		ordered:    r.labels.Equivalent(),
	}
}
