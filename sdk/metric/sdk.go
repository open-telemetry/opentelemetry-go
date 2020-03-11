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

		// observers is a set of `*observer` instances
		observers sync.Map

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
	}

	instrument struct {
		descriptor *export.Descriptor
		meter      *SDK
	}

	// orderedLabels is a variable-size array of core.KeyValue
	// suitable for use as a map key.
	orderedLabels interface{}

	// labels implements the OpenTelemetry LabelSet API,
	// represents an internalized set of labels that may be used
	// repeatedly.
	labels struct {
		meter *SDK

		// slice is a slice of `ordered`.
		slice sortedLabels

		// ordered is the output of sorting and deduplicating
		// the labels, copied into an array of the correct
		// size for use as a map key.
		ordered orderedLabels
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its LabelSet.
	mapkey struct {
		descriptor *export.Descriptor
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

		// descriptor describes the metric instrument.
		descriptor *export.Descriptor

		// recorder implements the actual RecordOne() API,
		// depending on the type of aggregation.  If nil, the
		// metric was disabled by the exporter.
		recorder export.Aggregator
	}

	observerResult struct {
		observer *observer
	}

	int64ObserverResult struct {
		result observerResult
	}

	float64ObserverResult struct {
		result observerResult
	}

	observerCallback func(result observerResult)

	observer struct {
		meter      *SDK
		descriptor *export.Descriptor
		// recorders maps ordered labels to the pair of
		// labelset and recorder
		recorders map[orderedLabels]labeledRecorder
		callback  observerCallback
	}

	labeledRecorder struct {
		recorder      export.Aggregator
		labels        *labels
		modifiedEpoch int64
	}

	int64Observer struct {
		observer *observer
	}

	float64Observer struct {
		observer *observer
	}

	ErrorHandler func(error)
)

var (
	_ api.Meter                 = &SDK{}
	_ api.LabelSet              = &labels{}
	_ api.InstrumentImpl        = &instrument{}
	_ api.BoundInstrumentImpl   = &record{}
	_ api.Int64Observer         = int64Observer{}
	_ api.Float64Observer       = float64Observer{}
	_ api.Int64ObserverResult   = int64ObserverResult{}
	_ api.Float64ObserverResult = float64ObserverResult{}

	kvType = reflect.TypeOf(core.KeyValue{})
)

func (r observerResult) observe(number core.Number, ls api.LabelSet) {
	r.observer.recordOne(number, ls)
}

func (o *observer) recordOne(number core.Number, ls api.LabelSet) {
	if err := aggregator.RangeTest(number, o.descriptor); err != nil {
		o.meter.errorHandler(err)
		return
	}
	recorder := o.getRecorder(ls)
	if recorder == nil {
		// The instrument is disabled according to the
		// AggregationSelector.
		return
	}
	if err := recorder.Update(context.Background(), number, o.descriptor); err != nil {
		o.meter.errorHandler(err)
		return
	}
}

func (o *observer) getRecorder(ls api.LabelSet) export.Aggregator {
	labels := o.meter.labsFor(ls)
	lrec, ok := o.recorders[labels.ordered]
	if ok {
		lrec.modifiedEpoch = o.meter.currentEpoch
		o.recorders[labels.ordered] = lrec
		return lrec.recorder
	}
	rec := o.meter.batcher.AggregatorFor(o.descriptor)
	if o.recorders == nil {
		o.recorders = make(map[orderedLabels]labeledRecorder)
	}
	// This may store nil recorder in the map, thus disabling the
	// observer for the labelset for good. This is intentional,
	// but will be revisited later.
	o.recorders[labels.ordered] = labeledRecorder{
		recorder:      rec,
		labels:        labels,
		modifiedEpoch: o.meter.currentEpoch,
	}
	return rec
}

func (o *observer) unregister() {
	o.meter.observers.Delete(o)
}

func (r int64ObserverResult) Observe(value int64, labels api.LabelSet) {
	r.result.observe(core.NewInt64Number(value), labels)
}

func (r float64ObserverResult) Observe(value float64, labels api.LabelSet) {
	r.result.observe(core.NewFloat64Number(value), labels)
}

func (o int64Observer) Unregister() {
	o.observer.unregister()
}

func (o float64Observer) Unregister() {
	o.observer.unregister()
}

func (i *instrument) Meter() api.Meter {
	return i.meter
}

func (m *SDK) SetErrorHandler(f ErrorHandler) {
	m.errorHandler = f
}

func (i *instrument) acquireHandle(ls *labels) *record {
	// Create lookup key for sync.Map (one allocation)
	mk := mapkey{
		descriptor: i.descriptor,
		ordered:    ls.ordered,
	}

	if actual, ok := i.meter.current.Load(mk); ok {
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
		labels:     ls,
		descriptor: i.descriptor,
		refMapped:  refcountMapped{value: 2},
		modified:   0,
		recorder:   i.meter.batcher.AggregatorFor(i.descriptor),
	}

	for {
		// Load/Store: there's a memory allocation to place `mk` into
		// an interface here.
		if actual, loaded := i.meter.current.LoadOrStore(mk, rec); loaded {
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

func (i *instrument) Bind(ls api.LabelSet) api.BoundInstrumentImpl {
	labs := i.meter.labsFor(ls)
	return i.acquireHandle(labs)
}

func (i *instrument) RecordOne(ctx context.Context, number core.Number, ls api.LabelSet) {
	ourLs := i.meter.labsFor(ls)
	h := i.acquireHandle(ourLs)
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
func New(batcher export.Batcher, labelEncoder export.LabelEncoder) *SDK {
	m := &SDK{
		batcher:      batcher,
		labelEncoder: labelEncoder,
		errorHandler: DefaultErrorHandler,
	}
	m.empty.meter = m
	return m
}

func DefaultErrorHandler(err error) {
	fmt.Fprintln(os.Stderr, "Metrics SDK error:", err)
}

// Labels returns a LabelSet corresponding to the arguments.  Passed
// labels are de-duplicated, with last-value-wins semantics.
func (m *SDK) Labels(kvs ...core.KeyValue) api.LabelSet {
	// Check for empty set.
	if len(kvs) == 0 {
		return &m.empty
	}

	ls := &labels{ // allocation
		meter: m,
		slice: kvs,
	}

	// Sort and de-duplicate.  Note: this use of `ls.slice` avoids
	// an allocation by using the address-able field rather than
	// `kvs`.  Labels retains a copy of this slice, i.e., the
	// initial allocation at the varargs call site.
	//
	// Note that `ls.slice` continues to refer to this memory,
	// even though a new array is allocated for `ls.ordered`.  It
	// is possible for the `slice` to refer to the same memory,
	// although in the reflection code path of `computeOrdered` it
	// costs an allocation to yield a slice through
	// `(reflect.Value).Interface()`.
	//
	// TODO: There is a possibility that the caller passes values
	// without an allocation (e.g., `meter.Labels(kvs...)`), and
	// that the user could later modify the slice, leading to
	// incorrect results.  This is indeed a risk, one that should
	// be quickly addressed via the following TODO.
	//
	// TODO: It would be better overall if the export.Labels interface
	// did not expose a slice via `Ordered()`, if instead it exposed
	// getter methods like `Len()` and `Order(i int)`.  Then we would
	// just implement the interface using the `orderedLabels` array.
	sort.Stable(&ls.slice)

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
	ls.slice = kvs
	ls.computeOrdered(kvs)
	return ls
}

func (ls *labels) computeOrdered(kvs []core.KeyValue) {
	switch len(kvs) {
	case 1:
		ptr := new([1]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 2:
		ptr := new([2]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 3:
		ptr := new([3]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 4:
		ptr := new([4]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 5:
		ptr := new([5]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 6:
		ptr := new([6]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 7:
		ptr := new([7]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 8:
		ptr := new([8]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 9:
		ptr := new([9]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	case 10:
		ptr := new([10]core.KeyValue)
		copy((*ptr)[:], kvs)
		ls.ordered = *ptr
	default:
		at := reflect.New(reflect.ArrayOf(len(kvs), kvType)).Elem()

		for i := 0; i < len(kvs); i++ {
			*(at.Index(i).Addr().Interface().(*core.KeyValue)) = kvs[i]
		}

		ls.ordered = at.Interface()
	}
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

func newDescriptor(name string, metricKind export.Kind, numberKind core.NumberKind, opts *api.Options) *export.Descriptor {
	return export.NewDescriptor(
		name,
		metricKind,
		opts.Keys,
		opts.Description,
		opts.Unit,
		numberKind,
		opts.Alternate)
}

func (m *SDK) newInstrument(name string, metricKind export.Kind, numberKind core.NumberKind, opts *api.Options) *instrument {
	descriptor := newDescriptor(name, metricKind, numberKind, opts)
	return &instrument{
		descriptor: descriptor,
		meter:      m,
	}
}

func (m *SDK) newCounterInstrument(name string, numberKind core.NumberKind, cos ...api.CounterOptionApplier) *instrument {
	opts := api.Options{}
	api.ApplyCounterOptions(&opts, cos...)
	return m.newInstrument(name, export.CounterKind, numberKind, &opts)
}

func (m *SDK) newMeasureInstrument(name string, numberKind core.NumberKind, mos ...api.MeasureOptionApplier) *instrument {
	opts := api.Options{}
	api.ApplyMeasureOptions(&opts, mos...)
	return m.newInstrument(name, export.MeasureKind, numberKind, &opts)
}

func (m *SDK) NewInt64Counter(name string, cos ...api.CounterOptionApplier) (api.Int64Counter, error) {
	return api.WrapInt64CounterInstrument(m.newCounterInstrument(name, core.Int64NumberKind, cos...), nil)
}

func (m *SDK) NewFloat64Counter(name string, cos ...api.CounterOptionApplier) (api.Float64Counter, error) {
	return api.WrapFloat64CounterInstrument(m.newCounterInstrument(name, core.Float64NumberKind, cos...), nil)
}

func (m *SDK) NewInt64Measure(name string, mos ...api.MeasureOptionApplier) (api.Int64Measure, error) {
	return api.WrapInt64MeasureInstrument(m.newMeasureInstrument(name, core.Int64NumberKind, mos...), nil)
}

func (m *SDK) NewFloat64Measure(name string, mos ...api.MeasureOptionApplier) (api.Float64Measure, error) {
	return api.WrapFloat64MeasureInstrument(m.newMeasureInstrument(name, core.Float64NumberKind, mos...), nil)
}

func (m *SDK) RegisterInt64Observer(name string, callback api.Int64ObserverCallback, oos ...api.ObserverOptionApplier) (api.Int64Observer, error) {
	if callback == nil {
		return api.NoopMeter{}.RegisterInt64Observer("", nil)
	}
	opts := api.Options{}
	api.ApplyObserverOptions(&opts, oos...)
	descriptor := newDescriptor(name, export.ObserverKind, core.Int64NumberKind, &opts)
	cb := wrapInt64ObserverCallback(callback)
	obs := m.newObserver(descriptor, cb)
	return int64Observer{
		observer: obs,
	}, nil
}

func wrapInt64ObserverCallback(callback api.Int64ObserverCallback) observerCallback {
	return func(result observerResult) {
		typeSafeResult := int64ObserverResult{
			result: result,
		}
		callback(typeSafeResult)
	}
}

func (m *SDK) RegisterFloat64Observer(name string, callback api.Float64ObserverCallback, oos ...api.ObserverOptionApplier) (api.Float64Observer, error) {
	if callback == nil {
		return api.NoopMeter{}.RegisterFloat64Observer("", nil)
	}
	opts := api.Options{}
	api.ApplyObserverOptions(&opts, oos...)
	descriptor := newDescriptor(name, export.ObserverKind, core.Float64NumberKind, &opts)
	cb := wrapFloat64ObserverCallback(callback)
	obs := m.newObserver(descriptor, cb)
	return float64Observer{
		observer: obs,
	}, nil
}

func wrapFloat64ObserverCallback(callback api.Float64ObserverCallback) observerCallback {
	return func(result observerResult) {
		typeSafeResult := float64ObserverResult{
			result: result,
		}
		callback(typeSafeResult)
	}
}

func (m *SDK) newObserver(descriptor *export.Descriptor, callback observerCallback) *observer {
	obs := &observer{
		meter:      m,
		descriptor: descriptor,
		recorders:  nil,
		callback:   callback,
	}
	m.observers.Store(obs, nil)
	return obs
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
	checkpointed += m.collectObservers(ctx)
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

func (m *SDK) collectObservers(ctx context.Context) int {
	checkpointed := 0

	m.observers.Range(func(key, value interface{}) bool {
		obs := key.(*observer)
		result := observerResult{
			observer: obs,
		}
		obs.callback(result)
		checkpointed += m.checkpointObserver(ctx, obs)
		return true
	})

	return checkpointed
}

func (m *SDK) checkpointRecord(ctx context.Context, r *record) int {
	return m.checkpoint(ctx, r.descriptor, r.recorder, r.labels)
}

func (m *SDK) checkpointObserver(ctx context.Context, obs *observer) int {
	if len(obs.recorders) == 0 {
		return 0
	}
	checkpointed := 0
	for encodedLabels, lrec := range obs.recorders {
		epochDiff := m.currentEpoch - lrec.modifiedEpoch
		if epochDiff == 0 {
			checkpointed += m.checkpoint(ctx, obs.descriptor, lrec.recorder, lrec.labels)
		} else if epochDiff > 1 {
			// This is second collection cycle with no
			// observations for this labelset. Remove the
			// recorder.
			delete(obs.recorders, encodedLabels)
		}
	}
	if len(obs.recorders) == 0 {
		obs.recorders = nil
	}
	return checkpointed
}

func (m *SDK) checkpoint(ctx context.Context, descriptor *export.Descriptor, recorder export.Aggregator, labels *labels) int {
	if recorder == nil {
		return 0
	}
	recorder.Checkpoint(ctx, descriptor)

	// TODO Labels are encoded once per collection interval,
	// instead of once per bound instrument lifetime.  This can be
	// addressed similarly to OTEP 78, see
	// https://github.com/jmacd/opentelemetry-go/blob/8bed2e14df7f9f4688fbab141924bb786dc9a3a1/api/context/internal/set.go#L89
	exportLabels := export.NewLabels(labels.slice, m.labelEncoder.Encode(labels.slice), m.labelEncoder)
	exportRecord := export.NewRecord(descriptor, exportLabels, recorder)
	err := m.batcher.Process(ctx, exportRecord)
	if err != nil {
		m.errorHandler(err)
	}
	return 1
}

// RecordBatch enters a batch of metric events.
func (m *SDK) RecordBatch(ctx context.Context, ls api.LabelSet, measurements ...api.Measurement) {
	for _, meas := range measurements {
		meas.InstrumentImpl().RecordOne(ctx, meas.Number(), ls)
	}
}

// GetDescriptor returns the descriptor of an instrument, which is not
// part of the public metric API.
func (m *SDK) GetDescriptor(inst metric.InstrumentImpl) *export.Descriptor {
	if ii, ok := inst.(*instrument); ok {
		return ii.descriptor
	}
	return nil
}

func (r *record) RecordOne(ctx context.Context, number core.Number) {
	if r.recorder == nil {
		// The instrument is disabled according to the AggregationSelector.
		return
	}
	if err := aggregator.RangeTest(number, r.descriptor); err != nil {
		r.labels.meter.errorHandler(err)
		return
	}
	if err := r.recorder.Update(ctx, number, r.descriptor); err != nil {
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
		descriptor: r.descriptor,
		ordered:    r.labels.ordered,
	}
}
