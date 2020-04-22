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
	"reflect"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	api "go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/resource"
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

		// resource represents the entity producing telemetry.
		resource resource.Resource

		// asyncSortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation.  It is cleared after use.
		asyncSortSlice sortedLabels
	}

	syncInstrument struct {
		instrument
	}

	// orderedLabels is a variable-size array of core.KeyValue
	// suitable for use as a map key.
	orderedLabels interface{}

	// labels represents an internalized set of labels that have been
	// sorted and deduplicated.
	labels struct {
		// cachedEncoderID needs to be aligned for atomic access
		cachedEncoderID int64
		// cachedEncoded is an encoded version of ordered
		// labels
		cachedEncoded string

		// ordered is the output of sorting and deduplicating
		// the labels, copied into an array of the correct
		// size for use as a map key.
		ordered orderedLabels
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its labels.
	mapkey struct {
		descriptor *metric.Descriptor
		ordered    orderedLabels
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

		// labels is the processed label set for this record.
		//
		// labels has to be aligned for 64-bit atomic operations.
		labels labels

		// sortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation.
		sortSlice sortedLabels

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
		recorders map[orderedLabels]labeledRecorder

		callback func(func(core.Number, []core.KeyValue))
	}

	labeledRecorder struct {
		observedEpoch int64
		labels        labels
		recorder      export.Aggregator
	}

	ErrorHandler func(error)
)

var (
	_ api.MeterImpl       = &SDK{}
	_ api.AsyncImpl       = &asyncInstrument{}
	_ api.SyncImpl        = &syncInstrument{}
	_ api.BoundSyncImpl   = &record{}
	_ api.Resourcer       = &SDK{}
	_ export.LabelStorage = &labels{}
	_ export.Labels       = &labels{}

	kvType = reflect.TypeOf(core.KeyValue{})

	emptyLabels = labels{
		ordered: [0]core.KeyValue{},
	}
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
	labels := a.meter.makeLabels(kvs, &a.meter.asyncSortSlice)

	lrec, ok := a.recorders[labels.ordered]
	if ok {
		if lrec.observedEpoch == a.meter.currentEpoch {
			// last value wins for Observers, so if we see the same labels
			// in the current epoch, we replace the old recorder
			lrec.recorder = a.meter.batcher.AggregatorFor(&a.descriptor)
		} else {
			lrec.observedEpoch = a.meter.currentEpoch
		}
		a.recorders[labels.ordered] = lrec
		return lrec.recorder
	}
	rec := a.meter.batcher.AggregatorFor(&a.descriptor)
	if a.recorders == nil {
		a.recorders = make(map[orderedLabels]labeledRecorder)
	}
	// This may store nil recorder in the map, thus disabling the
	// asyncInstrument for the labelset for good. This is intentional,
	// but will be revisited later.
	a.recorders[labels.ordered] = labeledRecorder{
		recorder:      rec,
		labels:        labels,
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
func (s *syncInstrument) acquireHandle(kvs []core.KeyValue, lptr *labels) *record {
	var rec *record
	var labels labels

	if lptr == nil || lptr.ordered == nil {
		// This memory allocation may not be used, but it's
		// needed for the `sortSlice` field, to avoid an
		// allocation while sorting.
		rec = &record{}
		labels = s.meter.makeLabels(kvs, &rec.sortSlice)
	} else {
		labels = *lptr
	}

	// Create lookup key for sync.Map (one allocation, as this
	// passes through an interface{})
	mk := mapkey{
		descriptor: &s.descriptor,
		ordered:    labels.ordered,
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
	}
	rec.refMapped = refcountMapped{value: 2}
	rec.labels = labels
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
		resource:     c.Resource,
	}
}

func DefaultErrorHandler(err error) {
	fmt.Fprintln(os.Stderr, "Metrics SDK error:", err)
}

// makeLabels returns a `labels` corresponding to the arguments.  Labels
// are sorted and de-duplicated, with last-value-wins semantics.  Note that
// sorting and deduplicating happens in-place to avoid allocation, so the
// passed slice will be modified.  The `sortSlice` argument refers to a memory
// location used temporarily while sorting the slice, to avoid a memory
// allocation.
func (m *SDK) makeLabels(kvs []core.KeyValue, sortSlice *sortedLabels) labels {
	// Check for empty set.
	if len(kvs) == 0 {
		return emptyLabels
	}

	*sortSlice = kvs

	// Sort and de-duplicate.  Note: this use of `sortSlice`
	// avoids an allocation because it is a pointer.
	sort.Stable(sortSlice)

	*sortSlice = nil

	oi := 1
	for i := 1; i < len(kvs); i++ {
		if kvs[i-1].Key == kvs[i].Key {
			// Overwrite the value for "last-value wins".
			kvs[oi-1].Value = kvs[i].Value
			continue
		}
		kvs[oi] = kvs[i]
		oi++
	}
	kvs = kvs[0:oi]
	return computeOrderedLabels(kvs)
}

// NumLabels is a part of an implementation of the export.LabelStorage
// interface.
func (ls *labels) NumLabels() int {
	return reflect.ValueOf(ls.ordered).Len()
}

// GetLabel is a part of an implementation of the export.LabelStorage
// interface.
func (ls *labels) GetLabel(idx int) core.KeyValue {
	// Note: The Go compiler successfully avoids an allocation for
	// the interface{} conversion here:
	return reflect.ValueOf(ls.ordered).Index(idx).Interface().(core.KeyValue)
}

// Iter is a part of an implementation of the export.Labels interface.
func (ls *labels) Iter() export.LabelIterator {
	return export.NewLabelIterator(ls)
}

// Encoded is a part of an implementation of the export.Labels
// interface.
func (ls *labels) Encoded(encoder export.LabelEncoder) string {
	id := encoder.ID()
	if id <= 0 {
		// Punish misbehaving encoders by not even trying to
		// cache them
		return encoder.Encode(ls.Iter())
	}
	cachedID := atomic.LoadInt64(&ls.cachedEncoderID)
	// If cached ID is less than zero, it means that other
	// goroutine is currently caching the encoded labels and the
	// ID of the encoder. Wait until it's done - it's a
	// nonblocking op.
	for cachedID < 0 {
		// Let other goroutine finish its work.
		runtime.Gosched()
		cachedID = atomic.LoadInt64(&ls.cachedEncoderID)
	}
	// At this point, cachedID is either 0 (nothing cached) or
	// some other number.
	//
	// If cached ID is the same as ID of the passed encoder, we've
	// got the fast path.
	if cachedID == id {
		return ls.cachedEncoded
	}
	// If we are here, either some other encoder cached its
	// encoded labels or the cache is still for the taking. Either
	// way, we need to compute the encoded labels anyway.
	encoded := encoder.Encode(ls.Iter())
	// If some other encoder took the cache, then we just return
	// our encoded labels. That's a slow path.
	if cachedID > 0 {
		return encoded
	}
	// Try to take the cache for ourselves. This is the place
	// where other encoders may be "blocked".
	if atomic.CompareAndSwapInt64(&ls.cachedEncoderID, 0, -1) {
		// The cache is ours.
		ls.cachedEncoded = encoded
		atomic.StoreInt64(&ls.cachedEncoderID, id)
	}
	return encoded
}

func computeOrderedLabels(kvs []core.KeyValue) labels {
	var ls labels
	ls.ordered = computeOrderedFixed(kvs)
	if ls.ordered == nil {
		ls.ordered = computeOrderedReflect(kvs)
	}
	return ls
}

func computeOrderedFixed(kvs []core.KeyValue) orderedLabels {
	switch len(kvs) {
	case 1:
		ptr := new([1]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 2:
		ptr := new([2]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 3:
		ptr := new([3]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 4:
		ptr := new([4]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 5:
		ptr := new([5]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 6:
		ptr := new([6]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 7:
		ptr := new([7]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 8:
		ptr := new([8]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 9:
		ptr := new([9]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 10:
		ptr := new([10]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	default:
		return nil
	}
}

func computeOrderedReflect(kvs []core.KeyValue) interface{} {
	at := reflect.New(reflect.ArrayOf(len(kvs), kvType)).Elem()
	for i, kv := range kvs {
		*(at.Index(i).Addr().Interface().(*core.KeyValue)) = kv
	}
	return at.Interface()
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
	return m.checkpoint(ctx, &r.inst.descriptor, r.recorder, &r.labels)
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
			checkpointed += m.checkpoint(ctx, &a.descriptor, lrec.recorder, &lrec.labels)
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

func (m *SDK) checkpoint(ctx context.Context, descriptor *metric.Descriptor, recorder export.Aggregator, labels *labels) int {
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

// Resource returns the Resource this SDK was created with describing the
// entity for which it creates instruments for.
//
// Resource means that the SDK implements the Resourcer interface and
// therefore all metric instruments it creates will inherit its
// Resource by default unless explicitly overwritten.
func (m *SDK) Resource() resource.Resource {
	return m.resource
}

// RecordBatch enters a batch of metric events.
func (m *SDK) RecordBatch(ctx context.Context, kvs []core.KeyValue, measurements ...api.Measurement) {
	// Labels will be computed the first time acquireHandle is
	// called.  Subsequent calls to acquireHandle will re-use the
	// previously computed value instead of recomputing the
	// ordered labels.
	var labels labels
	for i, meas := range measurements {
		s := meas.SyncImpl().(*syncInstrument)

		h := s.acquireHandle(kvs, &labels)

		// Re-use labels for the next measurement.
		if i == 0 {
			labels = h.labels
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
		ordered:    r.labels.ordered,
	}
}
