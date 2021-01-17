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
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	internal "go.opentelemetry.io/otel/internal/metric"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Accumulator implements the OpenTelemetry Meter API.  The
	// Accumulator is bound to a single export.Processor in
	// `NewAccumulator()`.
	//
	// The Accumulator supports a Collect() API to gather and export
	// current data.  Collect() should be arranged according to
	// the processor model.  Push-based processors will setup a
	// timer to call Collect() periodically.  Pull-based processors
	// will call Collect() when a pull request arrives.
	Accumulator struct {
		// current maps `mapkey` to *record.
		current sync.Map

		// asyncInstruments is a set of
		// `*asyncInstrument` instances
		asyncLock        sync.Mutex
		asyncInstruments *internal.AsyncInstrumentState

		// currentEpoch is the current epoch number. It is
		// incremented in `Collect()`.
		currentEpoch int64

		// processor is the configured processor+configuration.
		processor export.Processor

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex

		// asyncSortSlice has a single purpose - as a temporary
		// place for sorting during labels creation to avoid
		// allocation.  It is cleared after use.
		asyncSortSlice label.Sortable

		// resource is applied to all records in this Accumulator.
		resource *resource.Resource
	}

	syncInstrument struct {
		instrument

		// views is a logical set of registered views for this instrument.
		views sync.Map
	}

	// syncInstrumentViewKey is used only in syncInstrument as a key representation of a view.
	syncInstrumentViewKey struct {
		selector           metric.AggregatorSelector
		labelConfig        metric.ViewLabelConfig
		labelKeysSignature string
	}

	// syncInstrumentViewValue is used only in syncInstrument as a value representation of a view.
	syncInstrumentViewValue struct {
		labelKeys map[label.Key]struct{}
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its labels.
	mapkey struct {
		descriptor *metric.Descriptor
		ordered    label.Distinct
	}

	// boundInstrument is a snapshot of the view data for a given instrument.
	boundInstrument struct {
		records []*record
	}

	// record maintains the state of one metric instrument.  Due
	// the use of lock-free algorithms, there may be more than one
	// `record` in existence at a time, although at most one can
	// be referenced from the `Accumulator.current` map.
	record struct {
		// refMapped keeps track of refcounts and the mapping state to the
		// Accumulator.current map.
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

		// viewSelector is a reference to the corresponding view's aggregator selector.
		viewSelector metric.AggregatorSelector

		// current implements the actual RecordOne() API,
		// depending on the type of aggregation.  If nil, the
		// metric was disabled by the exporter.
		current    metric.Aggregator
		checkpoint metric.Aggregator
	}

	instrument struct {
		meter      *Accumulator
		descriptor metric.Descriptor
	}

	asyncInstrument struct {
		instrument
		// recorders maps ordered labels to the pair of
		// labelset and recorder
		recorders map[label.Distinct]*labeledRecorder
	}

	labeledRecorder struct {
		observedEpoch int64
		labels        *label.Set
		observed      metric.Aggregator
	}
)

var (
	_ metric.MeterImpl     = &Accumulator{}
	_ metric.AsyncImpl     = &asyncInstrument{}
	_ metric.SyncImpl      = &syncInstrument{}
	_ metric.BoundSyncImpl = &record{}
	_ metric.BoundSyncImpl = &boundInstrument{}

	ErrUninitializedInstrument = fmt.Errorf("use of an uninitialized instrument")
)

func (inst *instrument) Descriptor() metric.Descriptor {
	return inst.descriptor
}

func (a *asyncInstrument) Implementation() interface{} {
	return a
}

func (s *syncInstrument) registerView(labelConfig metric.ViewLabelConfig, labelKeys []label.Key, selector metric.AggregatorSelector) {
	key, val := newSyncInstrumentView(labelConfig, labelKeys, selector)

	// First-write-win: ignore duplicated view silently.
	// maybe return/record an error?
	s.views.LoadOrStore(key, val)
}

func (s *syncInstrument) rangeViews(kvs []label.KeyValue, f func([]label.KeyValue, metric.AggregatorSelector)) {
	empty := true
	s.views.Range(func(key, value interface{}) bool {
		empty = false
		return false
	})
	if empty {
		s.registerView(metric.Ungroup, nil, nil)
	}

	// compute on need
	var labelKvs map[label.Key]label.KeyValue

	s.views.Range(func(key, value interface{}) bool {
		viewKey := key.(*syncInstrumentViewKey)
		viewValue := value.(*syncInstrumentViewValue)

		var viewLabelKeyValues []label.KeyValue
		switch viewKey.labelConfig {
		case metric.Ungroup:
			viewLabelKeyValues = kvs

		case metric.LabelKeys:
			if labelKvs == nil {
				labelKvs = map[label.Key]label.KeyValue{}
				for _, kv := range kvs {
					labelKvs[kv.Key] = kv
				}
			}
			for k, kv := range labelKvs {
				if _, p := viewValue.labelKeys[k]; p {
					viewLabelKeyValues = append(viewLabelKeyValues, kv)
				}
			}
			for k := range viewValue.labelKeys {
				if _, p := labelKvs[k]; !p {
					viewLabelKeyValues = append(viewLabelKeyValues, k.Empty())
				}
			}

		case metric.DropAll:
			// keep nil
		default:
			panic(fmt.Sprintf("impossible label config: %s", viewKey.labelConfig))
		}

		f(viewLabelKeyValues, viewKey.selector)
		return true
	})
}

// newSyncInstrumentView returns a view key/value representation for the given parameters.
func newSyncInstrumentView(labelConfig metric.ViewLabelConfig, labelKeys []label.Key, selector metric.AggregatorSelector) (*syncInstrumentViewKey, *syncInstrumentViewValue) {
	key := &syncInstrumentViewKey{
		selector:           selector,
		labelConfig:        labelConfig,
		labelKeysSignature: "",
	}
	value := &syncInstrumentViewValue{
		labelKeys: nil,
	}

	if labelConfig == metric.LabelKeys {
		if len(labelKeys) > 0 {
			value.labelKeys = map[label.Key]struct{}{}
			for _, k := range labelKeys {
				value.labelKeys[k] = struct{}{}
			}

			keys := make([]string, 0, len(value.labelKeys))
			for k := range value.labelKeys {
				keys = append(keys, string(k))
			}
			sort.Strings(keys)
			key.labelKeysSignature = strings.Join(keys, "|")
		} else {
			// make sure there is only one representation for aggregations without any labels
			key.labelConfig = metric.DropAll
		}
	}

	return key, value
}

func (s *syncInstrument) Implementation() interface{} {
	return s
}

func (a *asyncInstrument) observe(num number.Number, labels *label.Set) {
	if err := aggregator.RangeTest(num, &a.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	recorder := a.getRecorder(labels)
	if recorder == nil {
		// The instrument is disabled according to the
		// AggregatorSelector.
		return
	}
	if err := recorder.Update(context.Background(), num, &a.descriptor); err != nil {
		otel.Handle(err)
		return
	}
}

func (a *asyncInstrument) getRecorder(labels *label.Set) metric.Aggregator {
	lrec, ok := a.recorders[labels.Equivalent()]
	if ok {
		// Note: SynchronizedMove(nil) can't return an error
		_ = lrec.observed.SynchronizedMove(nil, &a.descriptor)
		lrec.observedEpoch = a.meter.currentEpoch
		a.recorders[labels.Equivalent()] = lrec
		return lrec.observed
	}
	var rec metric.Aggregator
	a.meter.processor.AggregatorFor(&a.descriptor, &rec)
	if a.recorders == nil {
		a.recorders = make(map[label.Distinct]*labeledRecorder)
	}
	// This may store nil recorder in the map, thus disabling the
	// asyncInstrument for the labelset for good. This is intentional,
	// but will be revisited later.
	a.recorders[labels.Equivalent()] = &labeledRecorder{
		observed:      rec,
		labels:        labels,
		observedEpoch: a.meter.currentEpoch,
	}
	return rec
}

// acquireHandle gets or creates a `*record` corresponding to `kvs`,
// the input labels.  The second argument `labels` is passed in to
// support re-use of the orderedLabels computed by a previous
// measurement in the same batch.   This performs two allocations
// in the common case.
func (s *syncInstrument) acquireHandle(kvs []label.KeyValue, selector metric.AggregatorSelector) *record {
	var rec *record
	var equiv label.Distinct

	// This memory allocation may not be used, but it's
	// needed for the `sortSlice` field, to avoid an
	// allocation while sorting.
	rec = &record{}
	rec.storage = label.NewSetWithSortable(kvs, &rec.sortSlice)
	rec.labels = &rec.storage
	rec.viewSelector = selector
	equiv = rec.storage.Equivalent()

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

	rec.refMapped = refcountMapped{value: 2}
	rec.inst = s

	if selector != nil {
		selector.AggregatorFor(&s.descriptor, &rec.current, &rec.checkpoint)
	} else {
		s.meter.processor.AggregatorFor(&s.descriptor, &rec.current, &rec.checkpoint)
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

func (s *syncInstrument) Bind(kvs []label.KeyValue) metric.BoundSyncImpl {
	result := &boundInstrument{}
	s.rangeViews(kvs, func(viewKvs []label.KeyValue, selector metric.AggregatorSelector) {
		result.records = append(result.records, s.acquireHandle(viewKvs, selector))
	})
	return result
}

func (s *syncInstrument) RecordOne(ctx context.Context, num number.Number, kvs []label.KeyValue) {
	s.rangeViews(kvs, func(viewKvs []label.KeyValue, selector metric.AggregatorSelector) {
		h := s.acquireHandle(viewKvs, selector)
		h.RecordOne(ctx, num)
		h.Unbind()
	})
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
func NewAccumulator(processor export.Processor, resource *resource.Resource) *Accumulator {
	return &Accumulator{
		processor:        processor,
		asyncInstruments: internal.NewAsyncInstrumentState(),
		resource:         resource,
	}
}

func (m *Accumulator) RegisterView(v metric.View) {
	s := m.fromSync(v.SyncImpl())
	// maybe return an error if s is `nil`?
	if s != nil {
		s.registerView(v.LabelConfig(), v.LabelKeys(), v.AggregatorFactory())
	}
}

// NewSyncInstrument implements metric.MetricImpl.
func (m *Accumulator) NewSyncInstrument(descriptor metric.Descriptor) (metric.SyncImpl, error) {
	return &syncInstrument{
		instrument: instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

// NewAsyncInstrument implements metric.MetricImpl.
func (m *Accumulator) NewAsyncInstrument(descriptor metric.Descriptor, runner metric.AsyncRunner) (metric.AsyncImpl, error) {
	a := &asyncInstrument{
		instrument: instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}
	m.asyncLock.Lock()
	defer m.asyncLock.Unlock()
	m.asyncInstruments.Register(a, runner)
	return a, nil
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

	checkpointed := m.observeAsyncInstruments(ctx)
	checkpointed += m.collectSyncInstruments()
	m.currentEpoch++

	return checkpointed
}

func (m *Accumulator) collectSyncInstruments() int {
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
			checkpointed += m.checkpointRecord(inuse)
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
			checkpointed += m.checkpointRecord(inuse)
		}
		return true
	})

	return checkpointed
}

// CollectAsync implements internal.AsyncCollector.
func (m *Accumulator) CollectAsync(kv []label.KeyValue, obs ...metric.Observation) {
	labels := label.NewSetWithSortable(kv, &m.asyncSortSlice)

	for _, ob := range obs {
		if a := m.fromAsync(ob.AsyncImpl()); a != nil {
			a.observe(ob.Number(), &labels)
		}
	}
}

func (m *Accumulator) observeAsyncInstruments(ctx context.Context) int {
	m.asyncLock.Lock()
	defer m.asyncLock.Unlock()

	asyncCollected := 0

	m.asyncInstruments.Run(ctx, m)

	for _, inst := range m.asyncInstruments.Instruments() {
		if a := m.fromAsync(inst); a != nil {
			asyncCollected += m.checkpointAsync(a)
		}
	}

	return asyncCollected
}

func (m *Accumulator) checkpointRecord(r *record) int {
	if r.current == nil {
		return 0
	}
	err := r.current.SynchronizedMove(r.checkpoint, &r.inst.descriptor)
	if err != nil {
		otel.Handle(err)
		return 0
	}

	a := export.NewAccumulation(&r.inst.descriptor, r.labels, m.resource, r.checkpoint, r.viewSelector)
	err = m.processor.Process(a)
	if err != nil {
		otel.Handle(err)
	}
	return 1
}

func (m *Accumulator) checkpointAsync(a *asyncInstrument) int {
	if len(a.recorders) == 0 {
		return 0
	}
	checkpointed := 0
	for encodedLabels, lrec := range a.recorders {
		lrec := lrec
		epochDiff := m.currentEpoch - lrec.observedEpoch
		if epochDiff == 0 {
			if lrec.observed != nil {
				a := export.NewAccumulation(&a.descriptor, lrec.labels, m.resource, lrec.observed, nil)
				err := m.processor.Process(a)
				if err != nil {
					otel.Handle(err)
				}
				checkpointed++
			}
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

// RecordBatch enters a batch of metric events.
func (m *Accumulator) RecordBatch(ctx context.Context, kvs []label.KeyValue, measurements ...metric.Measurement) {
	for _, meas := range measurements {
		s := m.fromSync(meas.SyncImpl())
		if s == nil {
			continue
		}
		s.RecordOne(ctx, meas.Number(), kvs)
	}
}

// RecordOne implements metric.SyncImpl.
func (r *record) RecordOne(ctx context.Context, num number.Number) {
	if r.current == nil {
		// The instrument is disabled according to the AggregatorSelector.
		return
	}
	if err := aggregator.RangeTest(num, &r.inst.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	if err := r.current.Update(ctx, num, &r.inst.descriptor); err != nil {
		otel.Handle(err)
		return
	}
	// Record was modified, inform the Collect() that things need
	// to be collected while the record is still mapped.
	atomic.AddInt64(&r.updateCount, 1)
}

// Unbind implements metric.SyncImpl.
func (r *record) Unbind() {
	r.refMapped.unref()
}

func (r *record) mapkey() mapkey {
	return mapkey{
		descriptor: &r.inst.descriptor,
		ordered:    r.labels.Equivalent(),
	}
}

func (b *boundInstrument) RecordOne(ctx context.Context, number number.Number) {
	for _, r := range b.records {
		r.RecordOne(ctx, number)
	}
}

func (b *boundInstrument) Unbind() {
	for _, r := range b.records {
		r.Unbind()
	}
}

// fromSync gets a sync implementation object, checking for
// uninitialized instruments and instruments created by another SDK.
func (m *Accumulator) fromSync(sync metric.SyncImpl) *syncInstrument {
	if sync != nil {
		if inst, ok := sync.Implementation().(*syncInstrument); ok {
			return inst
		}
	}
	otel.Handle(ErrUninitializedInstrument)
	return nil
}

// fromSync gets an async implementation object, checking for
// uninitialized instruments and instruments created by another SDK.
func (m *Accumulator) fromAsync(async metric.AsyncImpl) *asyncInstrument {
	if async != nil {
		if inst, ok := async.Implementation().(*asyncInstrument); ok {
			return inst
		}
	}
	otel.Handle(ErrUninitializedInstrument)
	return nil
}
