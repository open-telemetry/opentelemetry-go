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

		// empty is the (singleton) result of Labels()
		// w/ zero arguments.
		empty labels

		// currentEpoch is the current epoch number. It is
		// incremented in `Collect()`.
		currentEpoch int64

		// batcher is the configured batcher+configuration.
		batcher export.Batcher

		// lencoder determines how labels are uniquely encoded.
		labelEncoder export.LabelEncoder

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex

		// errorHandler supports delivering errors to the user.
		errorHandler ErrorHandler

		// resource represents the entity producing telemetry.
		resource resource.Resource
	}

	syncInstrument struct {
		instrument
	}

	// orderedLabels is a variable-size array of core.KeyValue
	// suitable for use as a map key.
	orderedLabels interface{}

	// labels implements the OpenTelemetry LabelSet API,
	// represents an internalized set of labels that may be used
	// repeatedly.
	labels struct {
		meter *SDK
		// ordered is the output of sorting and deduplicating
		// the labels, copied into an array of the correct
		// size for use as a map key.
		ordered orderedLabels

		// sortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation
		sortSlice sortedLabels
		// cachedValue contains a `reflect.Value` of the `ordered`
		// member
		cachedValue reflect.Value
		// cachedEncoded contains an encoded version of the
		// `ordered` member
		cachedEncoded string
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its LabelSet.
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

		// modified is an atomic boolean that tracks if the current record
		// was modified since the last Collect().
		//
		// modified has to be aligned for 64-bit atomic operations.
		modified int64

		// labels is the LabelSet passed by the user.
		labels *labels

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

		callback func(func(core.Number, api.LabelSet))
	}

	labeledRecorder struct {
		recorder      export.Aggregator
		labels        *labels
		modifiedEpoch int64
	}

	ErrorHandler func(error)
)

var (
	_ api.MeterImpl       = &SDK{}
	_ api.LabelSet        = &labels{}
	_ api.AsyncImpl       = &asyncInstrument{}
	_ api.SyncImpl        = &syncInstrument{}
	_ api.BoundSyncImpl   = &record{}
	_ api.Resourcer       = &SDK{}
	_ export.LabelStorage = &labels{}

	kvType = reflect.TypeOf(core.KeyValue{})
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

func (a *asyncInstrument) observe(number core.Number, ls api.LabelSet) {
	if err := aggregator.RangeTest(number, &a.descriptor); err != nil {
		a.meter.errorHandler(err)
		return
	}
	recorder := a.getRecorder(ls)
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

func (a *asyncInstrument) getRecorder(ls api.LabelSet) export.Aggregator {
	labels := a.meter.labsFor(ls)
	lrec, ok := a.recorders[labels.ordered]
	if ok {
		lrec.modifiedEpoch = a.meter.currentEpoch
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
		modifiedEpoch: a.meter.currentEpoch,
	}
	return rec
}

func (m *SDK) SetErrorHandler(f ErrorHandler) {
	m.errorHandler = f
}

func (s *syncInstrument) acquireHandle(ls *labels) *record {
	// Create lookup key for sync.Map (one allocation)
	mk := mapkey{
		descriptor: &s.descriptor,
		ordered:    ls.ordered,
	}

	if actual, ok := s.meter.current.Load(mk); ok {
		// Existing record case, only one allocation so far.
		rec := actual.(*record)
		if rec.refMapped.ref() {
			// At this moment it is guaranteed that the entry is in
			// the map and will not be removed.
			return rec
		}
		// This entry is no longer mapped, try to add a new entry.
	}

	// There's a memory allocation here.
	rec := &record{
		labels:    ls,
		inst:      s,
		refMapped: refcountMapped{value: 2},
		modified:  0,
		recorder:  s.meter.batcher.AggregatorFor(&s.descriptor),
	}

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

func (s *syncInstrument) Bind(ls api.LabelSet) api.BoundSyncImpl {
	labs := s.meter.labsFor(ls)
	return s.acquireHandle(labs)
}

func (s *syncInstrument) RecordOne(ctx context.Context, number core.Number, ls api.LabelSet) {
	ourLs := s.meter.labsFor(ls)
	h := s.acquireHandle(ourLs)
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
func New(batcher export.Batcher, labelEncoder export.LabelEncoder, opts ...Option) *SDK {
	c := &Config{ErrorHandler: DefaultErrorHandler}
	for _, opt := range opts {
		opt.Apply(c)
	}

	m := &SDK{
		empty: labels{
			ordered: [0]core.KeyValue{},
		},
		batcher:      batcher,
		labelEncoder: labelEncoder,
		errorHandler: c.ErrorHandler,
		resource:     c.Resource,
	}
	m.empty.meter = m
	m.empty.cachedValue = reflect.ValueOf(m.empty.ordered)
	return m
}

func DefaultErrorHandler(err error) {
	fmt.Fprintln(os.Stderr, "Metrics SDK error:", err)
}

// Labels returns a LabelSet corresponding to the arguments.  Passed labels
// are sorted and de-duplicated, with last-value-wins semantics.  Note that
// sorting and deduplicating happens in-place to avoid allocation, so the
// passed slice will be modified.
func (m *SDK) Labels(kvs ...core.KeyValue) api.LabelSet {
	// Check for empty set.
	if len(kvs) == 0 {
		return &m.empty
	}

	ls := &labels{ // allocation
		meter:     m,
		sortSlice: kvs,
	}

	// Sort and de-duplicate.  Note: this use of `ls.sortSlice`
	// avoids an allocation by using the address-able field rather
	// than `kvs`.
	sort.Stable(&ls.sortSlice)
	ls.sortSlice = nil

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
	ls.computeOrdered(kvs)
	return ls
}

func (ls *labels) NumLabels() int {
	return ls.cachedValue.Len()
}

func (ls *labels) GetLabel(idx int) core.KeyValue {
	return ls.cachedValue.Index(idx).Interface().(core.KeyValue)
}

func (ls *labels) computeOrdered(kvs []core.KeyValue) {
	ls.ordered = computeOrderedFixed(kvs)
	if ls.ordered == nil {
		ls.ordered = computeOrderedReflect(kvs)
	}
	ls.cachedValue = reflect.ValueOf(ls.ordered)
}

func (ls *labels) ensureEncoded(encoder export.LabelEncoder) {
	if ls.cachedEncoded != "" {
		return
	}
	iter := export.NewLabelIterator(ls)
	ls.cachedEncoded = encoder.Encode(iter)
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

// labsFor sanitizes the input LabelSet.  The input will be rejected
// if it was created by another Meter instance, for example.
func (m *SDK) labsFor(ls api.LabelSet) *labels {
	if del, ok := ls.(api.LabelSetDelegate); ok {
		ls = del.Delegate()
	}
	if l, _ := ls.(*labels); l != nil && l.meter == m {
		return l
	}
	return &m.empty
}

func (m *SDK) NewSyncInstrument(descriptor api.Descriptor) (api.SyncImpl, error) {
	return &syncInstrument{
		instrument: instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *SDK) NewAsyncInstrument(descriptor api.Descriptor, callback func(func(core.Number, api.LabelSet))) (api.AsyncImpl, error) {
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
		inuse := value.(*record)
		unmapped := inuse.refMapped.tryUnmap()
		// If able to unmap then remove the record from the current Map.
		if unmapped {
			m.current.Delete(inuse.mapkey())
		}

		// Always report the values if a reference to the Record is active,
		// this is to keep the previous behavior.
		// TODO: Reconsider this logic.
		if inuse.refMapped.inUse() || atomic.LoadInt64(&inuse.modified) != 0 {
			atomic.StoreInt64(&inuse.modified, 0)
			checkpointed += m.checkpointRecord(ctx, inuse)
		}

		// Always continue to iterate over the entire map.
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
		epochDiff := m.currentEpoch - lrec.modifiedEpoch
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

func (m *SDK) checkpoint(ctx context.Context, descriptor *metric.Descriptor, recorder export.Aggregator, labels *labels) int {
	if recorder == nil {
		return 0
	}
	recorder.Checkpoint(ctx, descriptor)

	labels.ensureEncoded(m.labelEncoder)
	exportLabels := export.NewLabels(labels, labels.cachedEncoded, m.labelEncoder)
	exportRecord := export.NewRecord(descriptor, exportLabels, recorder)
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
func (m *SDK) RecordBatch(ctx context.Context, ls api.LabelSet, measurements ...api.Measurement) {
	for _, meas := range measurements {
		meas.SyncImpl().RecordOne(ctx, meas.Number(), ls)
	}
}

// GetDescriptor returns a pointer to the descriptor of an instrument,
// which is not part of the public metric API.  This is for testing.  Use
// SyncImpl().Descriptor() to get a copy of the descriptor.
func (m *SDK) GetDescriptor(inst api.SyncImpl) *metric.Descriptor {
	if ii, ok := inst.(*syncInstrument); ok {
		return &ii.descriptor
	}
	return nil
}

func (r *record) RecordOne(ctx context.Context, number core.Number) {
	if r.recorder == nil {
		// The instrument is disabled according to the AggregationSelector.
		return
	}
	if err := aggregator.RangeTest(number, &r.inst.descriptor); err != nil {
		r.labels.meter.errorHandler(err)
		return
	}
	if err := r.recorder.Update(ctx, number, &r.inst.descriptor); err != nil {
		r.labels.meter.errorHandler(err)
		return
	}
}

func (r *record) Unbind() {
	// Record was modified, inform the Collect() that things need to be collected.
	// TODO: Reconsider if we should marked as modified when an Update happens and
	// collect only when updates happened even for Bounds.
	atomic.StoreInt64(&r.modified, 1)
	r.refMapped.unref()
}

func (r *record) mapkey() mapkey {
	return mapkey{
		descriptor: &r.inst.descriptor,
		ordered:    r.labels.ordered,
	}
}
