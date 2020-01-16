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
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"

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

		// empty is the (singleton) result of Labels()
		// w/ zero arguments.
		empty labels

		// records is the head of both the primary and the
		// reclaim records lists.
		records doublePtr

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

	// sortedLabels are used to de-duplicate and canonicalize labels.
	sortedLabels []core.KeyValue

	// labels implements the OpenTelemetry LabelSet API,
	// represents an internalized set of labels that may be used
	// repeatedly.
	labels struct {
		meter   *SDK
		sorted  sortedLabels
		encoded string
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its LabelSet.
	mapkey struct {
		descriptor *export.Descriptor
		encoded    string
	}

	// record maintains the state of one metric instrument.  Due
	// the use of lock-free algorithms, there may be more than one
	// `record` in existence at a time, although at most one can
	// be referenced from the `SDK.current` map.
	record struct {
		// refcount counts the number of active handles on
		// referring to this record.  active handles prevent
		// removing the record from the current map.
		//
		// refcount has to be aligned for 64-bit atomic operations.
		refcount int64

		// collectedEpoch is the epoch number for which this
		// record has been exported.  This is modified by the
		// `Collect()` method.
		//
		// collectedEpoch has to be aligned for 64-bit atomic operations.
		collectedEpoch int64

		// modifiedEpoch is the latest epoch number for which
		// this record was updated.  Generally, if
		// modifiedEpoch is less than collectedEpoch, this
		// record is due for reclaimation.
		//
		// modifiedEpoch has to be aligned for 64-bit atomic operations.
		modifiedEpoch int64

		// reclaim is an atomic to control the start of reclaiming.
		//
		// reclaim has to be aligned for 64-bit atomic operations.
		reclaim int64

		// labels is the LabelSet passed by the user.
		labels *labels

		// descriptor describes the metric instrument.
		descriptor *export.Descriptor

		// recorder implements the actual RecordOne() API,
		// depending on the type of aggregation.  If nil, the
		// metric was disabled by the exporter.
		recorder export.Aggregator

		// next contains the next pointer for both the primary
		// and the reclaim lists.
		next doublePtr
	}

	ErrorHandler func(error)

	// singlePointer wraps an unsafe.Pointer and supports basic
	// load(), store(), clear(), and swapNil() operations.
	singlePtr struct {
		ptr unsafe.Pointer
	}

	// doublePtr is used for the head and next links of two lists.
	doublePtr struct {
		primary singlePtr
		reclaim singlePtr
	}
)

var (
	_ api.Meter               = &SDK{}
	_ api.LabelSet            = &labels{}
	_ api.InstrumentImpl      = &instrument{}
	_ api.BoundInstrumentImpl = &record{}

	// hazardRecord is used as a pointer value that indicates the
	// value is not included in any list.  (`nil` would be
	// ambiguous, since the final element in a list has `nil` as
	// the next pointer).
	hazardRecord = &record{}
)

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
		encoded:    ls.encoded,
	}

	if actual, ok := i.meter.current.Load(mk); ok {
		// Existing record case, only one allocation so far.
		rec := actual.(*record)
		atomic.AddInt64(&rec.refcount, 1)
		return rec
	}

	// There's a memory allocation here.
	rec := &record{
		labels:         ls,
		descriptor:     i.descriptor,
		refcount:       1,
		collectedEpoch: -1,
		modifiedEpoch:  0,
		recorder:       i.meter.batcher.AggregatorFor(i.descriptor),
	}

	// Load/Store: there's a memory allocation to place `mk` into
	// an interface here.
	if actual, loaded := i.meter.current.LoadOrStore(mk, rec); loaded {
		// Existing record case.
		rec = actual.(*record)
		atomic.AddInt64(&rec.refcount, 1)
		return rec
	}

	i.meter.addPrimary(rec)
	return rec
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
	// Note: This computes a canonical encoding of the labels to
	// use as a map key.  It happens to use the encoding used by
	// statsd for labels, allowing an optimization for statsd
	// batchers.  This could be made configurable in the
	// constructor, to support the same optimization for different
	// batchers.

	// Check for empty set.
	if len(kvs) == 0 {
		return &m.empty
	}

	ls := &labels{
		meter:  m,
		sorted: kvs,
	}

	// Sort and de-duplicate.
	sort.Stable(&ls.sorted)
	oi := 1
	for i := 1; i < len(ls.sorted); i++ {
		if ls.sorted[i-1].Key == ls.sorted[i].Key {
			ls.sorted[oi-1] = ls.sorted[i]
			continue
		}
		ls.sorted[oi] = ls.sorted[i]
		oi++
	}
	ls.sorted = ls.sorted[0:oi]

	ls.encoded = m.labelEncoder.Encode(ls.sorted)

	return ls
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

func (m *SDK) newInstrument(name string, metricKind export.Kind, numberKind core.NumberKind, opts *api.Options) *instrument {
	descriptor := export.NewDescriptor(
		name,
		metricKind,
		opts.Keys,
		opts.Description,
		opts.Unit,
		numberKind,
		opts.Alternate)
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

func (m *SDK) newGaugeInstrument(name string, numberKind core.NumberKind, gos ...api.GaugeOptionApplier) *instrument {
	opts := api.Options{}
	api.ApplyGaugeOptions(&opts, gos...)
	return m.newInstrument(name, export.GaugeKind, numberKind, &opts)
}

func (m *SDK) newMeasureInstrument(name string, numberKind core.NumberKind, mos ...api.MeasureOptionApplier) *instrument {
	opts := api.Options{}
	api.ApplyMeasureOptions(&opts, mos...)
	return m.newInstrument(name, export.MeasureKind, numberKind, &opts)
}

func (m *SDK) NewInt64Counter(name string, cos ...api.CounterOptionApplier) api.Int64Counter {
	return api.WrapInt64CounterInstrument(m.newCounterInstrument(name, core.Int64NumberKind, cos...))
}

func (m *SDK) NewFloat64Counter(name string, cos ...api.CounterOptionApplier) api.Float64Counter {
	return api.WrapFloat64CounterInstrument(m.newCounterInstrument(name, core.Float64NumberKind, cos...))
}

func (m *SDK) NewInt64Gauge(name string, gos ...api.GaugeOptionApplier) api.Int64Gauge {
	return api.WrapInt64GaugeInstrument(m.newGaugeInstrument(name, core.Int64NumberKind, gos...))
}

func (m *SDK) NewFloat64Gauge(name string, gos ...api.GaugeOptionApplier) api.Float64Gauge {
	return api.WrapFloat64GaugeInstrument(m.newGaugeInstrument(name, core.Float64NumberKind, gos...))
}

func (m *SDK) NewInt64Measure(name string, mos ...api.MeasureOptionApplier) api.Int64Measure {
	return api.WrapInt64MeasureInstrument(m.newMeasureInstrument(name, core.Int64NumberKind, mos...))
}

func (m *SDK) NewFloat64Measure(name string, mos ...api.MeasureOptionApplier) api.Float64Measure {
	return api.WrapFloat64MeasureInstrument(m.newMeasureInstrument(name, core.Float64NumberKind, mos...))
}

// saveFromReclaim puts a record onto the "reclaim" list when it
// detects an attempt to delete the record while it is still in use.
func (m *SDK) saveFromReclaim(rec *record) {
	for {
		reclaimed := atomic.LoadInt64(&rec.reclaim)
		if reclaimed != 0 {
			return
		}
		if atomic.CompareAndSwapInt64(&rec.reclaim, 0, 1) {
			break
		}
	}

	m.addReclaim(rec)
}

// Collect traverses the list of active records and exports data for
// each active instrument.  Collect() may not be called concurrently.
//
// During the collection pass, the export.Batcher will receive
// one Export() call per current aggregation.
//
// Returns the number of records that were checkpointed.
func (m *SDK) Collect(ctx context.Context) int {
	m.collectLock.Lock()
	defer m.collectLock.Unlock()

	checkpointed := 0

	var next *record
	for inuse := m.records.primary.swapNil(); inuse != nil; inuse = next {
		next = inuse.next.primary.load()

		refcount := atomic.LoadInt64(&inuse.refcount)

		if refcount > 0 {
			checkpointed += m.checkpoint(ctx, inuse)
			m.addPrimary(inuse)
			continue
		}

		modified := atomic.LoadInt64(&inuse.modifiedEpoch)
		collected := atomic.LoadInt64(&inuse.collectedEpoch)
		checkpointed += m.checkpoint(ctx, inuse)

		if modified >= collected {
			atomic.StoreInt64(&inuse.collectedEpoch, m.currentEpoch)
			m.addPrimary(inuse)
			continue
		}

		// Remove this entry.
		m.current.Delete(inuse.mapkey())
		inuse.next.primary.store(hazardRecord)
	}

	for chances := m.records.reclaim.swapNil(); chances != nil; chances = next {
		atomic.StoreInt64(&chances.collectedEpoch, m.currentEpoch)

		next = chances.next.reclaim.load()
		chances.next.reclaim.clear()
		atomic.StoreInt64(&chances.reclaim, 0)

		if chances.next.primary.load() == hazardRecord {
			checkpointed += m.checkpoint(ctx, chances)
			m.addPrimary(chances)
		}
	}

	m.currentEpoch++
	return checkpointed
}

func (m *SDK) checkpoint(ctx context.Context, r *record) int {
	if r.recorder == nil {
		return 0
	}
	r.recorder.Checkpoint(ctx, r.descriptor)
	labels := export.NewLabels(r.labels.sorted, r.labels.encoded, m.labelEncoder)
	err := m.batcher.Process(ctx, export.NewRecord(r.descriptor, labels, r.recorder))

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
	for {
		collected := atomic.LoadInt64(&r.collectedEpoch)
		modified := atomic.LoadInt64(&r.modifiedEpoch)

		updated := collected + 1

		if modified == updated {
			// No change
			break
		}
		if !atomic.CompareAndSwapInt64(&r.modifiedEpoch, modified, updated) {
			continue
		}

		if modified < collected {
			// This record could have been reclaimed.
			r.labels.meter.saveFromReclaim(r)
		}

		break
	}

	_ = atomic.AddInt64(&r.refcount, -1)
}

func (r *record) mapkey() mapkey {
	return mapkey{
		descriptor: r.descriptor,
		encoded:    r.labels.encoded,
	}
}
