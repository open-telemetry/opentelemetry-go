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
	"bytes"
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/api/core"
	api "go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/sdk/export"
)

var _ api.Meter = &SDK{}
var _ api.LabelSet = &labels{}
var hazardRecord = &record{}

type (
	// SDK implements the OpenTelemetry Meter API.
	SDK struct {
		// current maps `mapkey` to *record.
		current sync.Map

		// pool is a pool of labelset builders.
		pool sync.Pool // *bytes.Buffer

		// empty is the (singleton) result of DefineLabels()
		// w/ zero arguments.
		empty labels

		// records is the head of both the primary and the
		// reclaim records lists.
		records doublePtr

		// currentEpoch is the current epoch number. It is
		// incremented in `Collect()`.
		currentEpoch int64

		// exporter is the configured exporter+configuration.
		exporter export.MetricBatcher

		// collectLock prevents simultaneous calls to Collect().
		collectLock sync.Mutex
	}

	// sortedLabels are used to de-duplicate and canonicalize labels.
	sortedLabels []core.KeyValue

	// labels implements the OpenTelemetry LabelSet API,
	// represents an internalized set of labels that may be used
	// repeatedly.
	labels struct {
		meter   *SDK
		sorted  []core.KeyValue
		encoded string
	}

	// mapkey uniquely describes a metric instrument in terms of
	// its InstrumentID and the encoded form of its LabelSet.
	mapkey struct {
		id      api.DescriptorID
		encoded string
	}

	// record maintains the state of one metric instrument.  Due
	// the use of lock-free algorithms, there may be more than one
	// `record` in existence at a time, although at most one can
	// be referenced from the `SDK.current` map.
	record struct {
		// labels is the LabelSet passed by the user.
		labels *labels

		// descriptor describes the metric instrument.
		descriptor *api.Descriptor

		// refcount counts the number of active handles on
		// referring to this record.  active handles prevent
		// removing the record from the current map.
		refcount int64

		// collectedEpoch is the epoch number for which this
		// record has been exported.  This is modified by the
		// `Collect()` method.
		collectedEpoch int64

		// modifiedEpoch is the latest epoch number for which
		// this record was updated.  Generally, if
		// modifiedEpoch is less than collectedEpoch, this
		// record is due for reclaimation.
		modifiedEpoch int64

		// reclaim is an atomic to control the start of reclaiming.
		reclaim int64

		// recorder implements the actual RecordOne() API,
		// depending on the type of aggregation.  If nil, the
		// metric was disabled by the exporter.
		recorder export.MetricAggregator

		// next contains the next pointer for both the primary
		// and the reclaim lists.
		next doublePtr
	}

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

// New constructs a new *SDK that implements intermediate storage for
// metric instrument events, suporting configurable aggregation.
//
// The SDK does not start any background process to periodically
// collect metrics.  It is the caller's responsibility to arrange for
// periodic collection by calling the `Collect()` method.
func New(exporter export.MetricBatcher) *SDK {
	m := &SDK{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
		exporter: exporter,
	}
	m.empty.meter = m
	return m
}

// DefineLabels returns a LabelSet corresponding to the arguments.
// Labels are de-duplicated, with last-value-wins semantics.
func (m *SDK) DefineLabels(_ context.Context, kvs ...core.KeyValue) api.LabelSet {
	// Note: This computes a canonical encoding of the labels to
	// use as a map key.  It happens to use the encoding used by
	// statsd for labels, allowing an optimization for statsd
	// exporters.  This could be made configurable in the
	// constructor, to support the same optimization for different
	// exporters.

	// Check for empty set.
	if len(kvs) == 0 {
		return &m.empty
	}

	// Sort and de-duplicate.
	sorted := sortedLabels(kvs)
	sort.Stable(&sorted)
	oi := 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i-1].Key == sorted[i].Key {
			sorted[oi-1] = sorted[i]
			continue
		}
		sorted[oi] = sorted[i]
		oi++
	}
	sorted = sorted[0:oi]

	// Serialize.
	buf := m.pool.Get().(*bytes.Buffer)
	defer m.pool.Put(buf)
	buf.Reset()
	_, _ = buf.WriteRune('|')
	delimiter := '#'
	for _, kv := range sorted {
		_, _ = buf.WriteRune(delimiter)
		_, _ = buf.WriteString(string(kv.Key))
		_, _ = buf.WriteRune(':')
		_, _ = buf.WriteString(kv.Value.Emit())
		delimiter = ','
	}

	return &labels{
		meter:   m,
		sorted:  sorted,
		encoded: buf.String(),
	}
}

// labsFor sanitizes the input LabelSet.  The input will be rejected
// if it was created by another Meter instance, for example.
func (m *SDK) labsFor(ls api.LabelSet) *labels {
	if l, _ := ls.(*labels); l != nil && l.meter == m {
		return l
	}
	return &m.empty
}

// NewHandle returns an interface used to implement handles in the
// OpenTelemetry API.  This takes a reference and will persist until
// the recorder is deleted.
func (m *SDK) NewHandle(inst api.ExplicitReportingMetric, ls api.LabelSet) api.Handle {
	desc := inst.Descriptor()
	// Create lookup key for sync.Map
	labs := m.labsFor(ls)
	mk := mapkey{
		id:      desc.ID(),
		encoded: labs.encoded,
	}

	// There's a memory allocation here.
	rec := &record{
		labels:         labs,
		descriptor:     desc,
		refcount:       1,
		collectedEpoch: -1,
		modifiedEpoch:  0,
	}

	// Load/Store: there's a memory allocation to place `mk` into an interface here.
	if actual, loaded := m.current.LoadOrStore(mk, rec); loaded {
		// Existing record case.
		rec = actual.(*record)
		atomic.AddInt64(&rec.refcount, 1)
		return rec
	}
	rec.recorder = m.exporter.AggregatorFor(rec)

	m.addPrimary(rec)
	return rec
}

// DeleteHandle removes a reference on the underlying record.  This
// method checks for a race with `Collect()`, which may have (tried to)
// reclaim an idle recorder, and restores the situation.
func (m *SDK) DeleteHandle(r api.Handle) {
	rec, _ := r.(*record)
	if rec == nil {
		return
	}

	for {
		collected := atomic.LoadInt64(&rec.collectedEpoch)
		modified := atomic.LoadInt64(&rec.modifiedEpoch)

		updated := collected + 1

		if modified == updated {
			// No change
			break
		}
		if !atomic.CompareAndSwapInt64(&rec.modifiedEpoch, modified, updated) {
			continue
		}

		if modified < collected {
			// This record could have been reclaimed.
			m.saveFromReclaim(rec)
		}

		break
	}

	_ = atomic.AddInt64(&rec.refcount, -1)
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
// each instrument.
func (m *SDK) Collect(ctx context.Context) {
	// This logic relies on being single-threaded, enforce this.
	m.collectLock.Lock()
	defer m.collectLock.Unlock()

	var next *record
	for inuse := m.records.primary.swapNil(); inuse != nil; inuse = next {
		next = inuse.next.primary.load()

		refcount := atomic.LoadInt64(&inuse.refcount)

		if refcount > 0 {
			m.collect(ctx, inuse)
			m.addPrimary(inuse)
			continue
		}

		modified := atomic.LoadInt64(&inuse.modifiedEpoch)
		collected := atomic.LoadInt64(&inuse.collectedEpoch)
		m.collect(ctx, inuse)

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
			m.collect(ctx, chances)
			m.addPrimary(chances)
		}
	}

	m.currentEpoch++
}

func (m *SDK) collect(ctx context.Context, r *record) {
	if r.recorder != nil {
		r.recorder.Collect(ctx, r, m.exporter)
	}
}

// RecordSingle enters a single metric event.
func (m *SDK) RecordOne(ctx context.Context, ls api.LabelSet, measurement api.Measurement) {
	r := m.NewHandle(measurement.Instrument, ls)
	defer m.DeleteHandle(r)
	r.RecordOne(ctx, measurement.Value)
}

// RecordBatch enters a batch of metric events.
func (m *SDK) RecordBatch(ctx context.Context, ls api.LabelSet, measurements ...api.Measurement) {
	for _, meas := range measurements {
		m.RecordOne(ctx, ls, meas)
	}
}

func (l *labels) Meter() api.Meter {
	return l.meter
}

func (r *record) RecordOne(ctx context.Context, value api.MeasurementValue) {
	if r.recorder != nil {
		r.recorder.Update(ctx, value, r)
	}
}

func (r *record) mapkey() mapkey {
	return mapkey{
		id:      r.descriptor.ID(),
		encoded: r.labels.encoded,
	}
}

func (r *record) Descriptor() *api.Descriptor {
	return r.descriptor
}

func (r *record) Labels() []core.KeyValue {
	return r.labels.sorted
}
