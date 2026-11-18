// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// boundEntry holds a series value stored in [boundLimitedSyncMap] and tracks
// whether it is bound (survives Clear) and whether it has recorded at least
// one measurement in the current collection cycle.
type boundEntry[V any] struct {
	val      V
	measured atomic.Bool
	bound    bool
}

// boundLimitedSyncMap extends [limitedSyncMap] so bound series remain resident
// in the underlying [sync.Map] across [boundLimitedSyncMap.Clear] while unbound
// series are deleted.
type boundLimitedSyncMap[V any] struct {
	limitedSyncMap[*boundEntry[V]]
	boundOverflow    bool
	hasOverflowEntry bool
}

func (m *boundLimitedSyncMap[V]) atLimitLocked() bool {
	if m.aggLimit <= 0 {
		return false
	}
	concrete := m.len
	if m.hasOverflowEntry {
		concrete--
	}
	return concrete >= m.aggLimit-1
}

func (m *boundLimitedSyncMap[V]) LoadOrStoreAttr(
	lazy lazyFilteredAttributes,
	initValue func(*V, attribute.Set),
) *boundEntry[V] {
	distinct := lazy.Distinct()
	if actual, loaded := m.Load(distinct); loaded {
		return actual.(*boundEntry[V])
	}
	if m.aggLimit > 0 && m.overflow.Load() {
		if actual, loaded := m.Load(overflowSet.Equivalent()); loaded {
			return actual.(*boundEntry[V])
		}
	}

	m.lenMux.Lock()
	defer m.lenMux.Unlock()

	if actual, loaded := m.Load(distinct); loaded {
		return actual.(*boundEntry[V])
	}

	var fltrAttr attribute.Set
	overflow := m.atLimitLocked()
	if overflow {
		fltrAttr = overflowSet
		distinct = overflowSet.Equivalent()
		if actual, loaded := m.Load(distinct); loaded {
			m.overflow.Store(true)
			return actual.(*boundEntry[V])
		}
	} else {
		fltrAttr = lazy.Set()
	}

	e := &boundEntry[V]{}
	initValue(&e.val, fltrAttr)
	actual, loaded := m.LoadOrStore(distinct, e)
	if !loaded {
		m.len++
	}
	if overflow {
		m.hasOverflowEntry = true
		m.overflow.Store(true)
	}
	return actual.(*boundEntry[V])
}

type pinPlan[V any] struct {
	existing        *boundEntry[V]
	d               attribute.Distinct
	attrs           attribute.Set
	overflow        bool
	setOverflowFlag bool
}

func (m *boundLimitedSyncMap[V]) preparePinLocked(lazy lazyFilteredAttributes) pinPlan[V] {
	d := lazy.Distinct()
	if actual, loaded := m.Load(d); loaded {
		e := actual.(*boundEntry[V])
		overflow := m.hasOverflowEntry && d == overflowSet.Equivalent()
		attrs := overflowSet
		if !overflow {
			attrs = lazy.Set()
		}
		return pinPlan[V]{
			existing: e,
			d:        d,
			attrs:    attrs,
			overflow: overflow,
		}
	}
	if m.atLimitLocked() {
		d = overflowSet.Equivalent()
		if actual, loaded := m.Load(d); loaded {
			e := actual.(*boundEntry[V])
			return pinPlan[V]{
				existing:        e,
				d:               d,
				attrs:           overflowSet,
				overflow:        m.boundOverflow || !e.bound,
				setOverflowFlag: true,
			}
		}
		return pinPlan[V]{
			d:        d,
			attrs:    overflowSet,
			overflow: true,
		}
	}
	return pinPlan[V]{
		d:     d,
		attrs: lazy.Set(),
	}
}

func (m *boundLimitedSyncMap[V]) commitPinLocked(plan pinPlan[V], e *boundEntry[V]) {
	if plan.existing != nil {
		if !e.bound {
			e.bound = true
			if plan.overflow {
				m.boundOverflow = true
			}
		}
		if plan.setOverflowFlag {
			m.overflow.Store(true)
		}
		return
	}
	m.Store(plan.d, e)
	m.len++
	if plan.overflow {
		m.boundOverflow = true
		m.hasOverflowEntry = true
		m.overflow.Store(m.atLimitLocked())
	}
}

func (m *boundLimitedSyncMap[V]) pin(
	lazy lazyFilteredAttributes,
	initValue func(*V, attribute.Set),
) *boundEntry[V] {
	m.lenMux.Lock()
	defer m.lenMux.Unlock()
	plan := m.preparePinLocked(lazy)
	e := plan.existing
	if e == nil {
		e = &boundEntry[V]{bound: true}
		initValue(&e.val, plan.attrs)
	}
	m.commitPinLocked(plan, e)
	return e
}

func (m *boundLimitedSyncMap[V]) Clear() {
	m.lenMux.Lock()
	defer m.lenMux.Unlock()
	m.Range(func(key, value any) bool {
		e := value.(*boundEntry[V])
		if e.bound {
			if e.measured.Load() {
				e.measured.Store(false)
			}
		} else {
			m.Delete(key)
			m.len--
		}
		return true
	})
	m.hasOverflowEntry = m.boundOverflow
	m.overflow.Store(m.boundOverflow && m.atLimitLocked())
}

// boundHotColdMap manages two [boundLimitedSyncMap] instances for delta
// aggregations with bound instrument handles.
type boundHotColdMap[V any] struct {
	hotColdWaitGroup
	pinMux        sync.Mutex
	hotColdValMap [2]boundLimitedSyncMap[V]
}

func (m *boundHotColdMap[V]) init(limit int) {
	m.hotColdValMap[0].aggLimit = limit
	m.hotColdValMap[1].aggLimit = limit
}

func (m *boundHotColdMap[V]) hot(i hotIdx) *boundLimitedSyncMap[V] {
	return &m.hotColdValMap[i]
}

func (m *boundHotColdMap[V]) Clear(readIdx coldIdx) {
	m.hotColdValMap[readIdx].Clear()
}

func (m *boundHotColdMap[V]) pin(
	lazy lazyFilteredAttributes,
	initValue func(*V, attribute.Set),
) [2]*boundEntry[V] {
	m.pinMux.Lock()
	defer m.pinMux.Unlock()

	h := m.start()
	defer m.done(h)

	// Lock both buffers in canonical index order (0 then 1) so no partial
	// state is published if initValue (ReservoirFunc) panics on either buffer.
	// Deadlock-free because pinMux serializes pin(), and all other callers
	// hold at most one buffer's lenMux.
	m.hotColdValMap[0].lenMux.Lock()
	defer m.hotColdValMap[0].lenMux.Unlock()
	m.hotColdValMap[1].lenMux.Lock()
	defer m.hotColdValMap[1].lenMux.Unlock()

	hot := &m.hotColdValMap[h]
	cold := &m.hotColdValMap[1-h]

	planHot := hot.preparePinLocked(lazy)
	var eColdExisting *boundEntry[V]
	coldActual, coldLoaded := cold.Load(planHot.d)
	if coldLoaded {
		if e := coldActual.(*boundEntry[V]); e.bound {
			eColdExisting = e
		}
	}

	// Phase 1: run initValue (external ReservoirFunc) before mutating either
	// buffer, preserving hot-then-cold call order. If either call panics,
	// both maps remain untouched.
	eHot := planHot.existing
	if eHot == nil {
		eHot = &boundEntry[V]{bound: true}
		initValue(&eHot.val, planHot.attrs)
	}
	eCold := eColdExisting
	if eCold == nil {
		eCold = &boundEntry[V]{bound: true}
		initValue(&eCold.val, planHot.attrs)
	}

	// Phase 2: commit both buffers atomically (no panicking operations).
	hot.commitPinLocked(planHot, eHot)
	if eColdExisting == nil {
		cold.Store(planHot.d, eCold)
		if !coldLoaded {
			cold.len++
		}
	}
	if planHot.overflow {
		cold.boundOverflow = true
		cold.hasOverflowEntry = true
		cold.overflow.Store(cold.atLimitLocked())
	}

	var entries [2]*boundEntry[V]
	entries[h] = eHot
	entries[1-h] = eCold
	return entries
}

func initBoundSumValue[N int64 | float64](
	v *sumValue[N],
	attr attribute.Set,
	newRes func(attribute.Set) FilteredExemplarReservoir[N],
) {
	r := newRes(attr)
	_, isDrop := r.(*dropRes[N])
	v.res = r
	v.attrs = attr
	v.dropExemplars = isDrop
}

// BoundSum contains the operations of a bound-aware Sum aggregation.
//
// Measurements recorded via Measure (unbound) and Bind (bound) for the same
// attribute set are aggregated into a single data point on collection and count
// as a single series toward the configured AggregationLimit.
type BoundSum[N int64 | float64] struct {
	Measure            Measure[N]
	ComputeAggregation ComputeAggregation
	Bind               Binder[N]
}

type boundDeltaSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	vals   boundHotColdMap[sumValue[N]]
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func newBoundDeltaSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *boundDeltaSum[N] {
	s := &boundDeltaSum[N]{
		monotonic: monotonic,
		start:     now(),
		newRes:    r,
	}
	s.vals.init(limit)
	return s
}

func (s *boundDeltaSum[N]) initValue(v *sumValue[N], attr attribute.Set) {
	initBoundSumValue(v, attr, s.newRes)
}

func (s *boundDeltaSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	hotIdx := s.vals.start()
	defer s.vals.done(hotIdx)

	e := s.vals.hot(hotIdx).LoadOrStoreAttr(lazy, s.initValue)
	e.val.n.add(value)
	if !e.measured.Load() {
		e.measured.Store(true)
	}
	if !e.val.dropExemplars {
		e.val.res.Offer(ctx, value, lazy)
	}
}

func (s *boundDeltaSum[N]) bind(lazy lazyFilteredAttributes) BoundMeasure[N] {
	entries := s.vals.pin(lazy, s.initValue)
	if entries[0].val.dropExemplars && entries[1].val.dropExemplars {
		return func(_ context.Context, val N) {
			hotIdx := s.vals.start()
			e := entries[hotIdx]
			e.val.n.add(val)
			if !e.measured.Load() {
				e.measured.Store(true)
			}
			// Direct done call: nothing above invokes external code or can panic,
			// avoiding deferwrap overhead on this hot path.
			s.vals.done(hotIdx)
		}
	}
	return func(ctx context.Context, val N) {
		hotIdx := s.vals.start()
		defer s.vals.done(hotIdx)
		e := entries[hotIdx]
		e.val.n.add(val)
		if !e.measured.Load() {
			e.measured.Store(true)
		}
		if !e.val.dropExemplars {
			e.val.res.Offer(ctx, val, lazy)
		}
	}
}

func (s *boundDeltaSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.DeltaTemporality
	sData.IsMonotonic = s.monotonic

	s.vals.pinMux.Lock()
	defer s.vals.pinMux.Unlock()

	readIdx := s.vals.swapHotAndWait()
	cold := &s.vals.hotColdValMap[readIdx]
	cold.lenMux.Lock()
	defer cold.lenMux.Unlock()

	n := cold.len
	dPts := reset(sData.DataPoints, n, n)

	var i int
	cold.Range(func(key, value any) bool {
		e := value.(*boundEntry[sumValue[N]])
		if e.measured.Load() {
			collectExemplars(&dPts[i].Exemplars, e.val.res.Collect)
			dPts[i].Attributes = e.val.attrs
			dPts[i].StartTime = s.start
			dPts[i].Time = t
			dPts[i].Value = e.val.n.load()
			i++
			if e.bound {
				e.val.n.reset()
				e.measured.Store(false)
			}
		}
		if !e.bound {
			cold.Delete(key)
			cold.len--
		}
		return true
	})
	cold.hasOverflowEntry = cold.boundOverflow
	cold.overflow.Store(cold.boundOverflow && cold.atLimitLocked())

	s.start = t

	sData.DataPoints = dPts[:i]
	*dest = sData

	return i
}

type boundCumulativeSumValue[N int64 | float64] struct {
	sumValue[N]
	once sync.Once
}

type boundCumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	values boundLimitedSyncMap[boundCumulativeSumValue[N]]
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func newBoundCumulativeSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *boundCumulativeSum[N] {
	return &boundCumulativeSum[N]{
		monotonic: monotonic,
		start:     now(),
		values: boundLimitedSyncMap[boundCumulativeSumValue[N]]{
			limitedSyncMap: limitedSyncMap[*boundEntry[boundCumulativeSumValue[N]]]{aggLimit: limit},
		},
		newRes: r,
	}
}

func (s *boundCumulativeSum[N]) initValue(v *boundCumulativeSumValue[N], attr attribute.Set) {
	initBoundSumValue(&v.sumValue, attr, s.newRes)
}

func markCumulativeMeasured[N int64 | float64](e *boundEntry[boundCumulativeSumValue[N]]) {
	e.val.once.Do(func() {
		e.val.startTime = now()
		e.measured.Store(true)
	})
}

func (s *boundCumulativeSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	e := s.values.LoadOrStoreAttr(lazy, s.initValue)
	e.val.n.add(value)
	if !e.measured.Load() {
		markCumulativeMeasured(e)
	}
	if !e.val.dropExemplars {
		e.val.res.Offer(ctx, value, lazy)
	}
}

func (s *boundCumulativeSum[N]) bind(lazy lazyFilteredAttributes) BoundMeasure[N] {
	e := s.values.pin(lazy, s.initValue)
	if e.val.dropExemplars {
		return func(_ context.Context, val N) {
			e.val.n.add(val)
			if !e.measured.Load() {
				markCumulativeMeasured(e)
			}
		}
	}
	return func(ctx context.Context, val N) {
		e.val.n.add(val)
		if !e.measured.Load() {
			markCumulativeMeasured(e)
		}
		e.val.res.Offer(ctx, val, lazy)
	}
}

func (s *boundCumulativeSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.CumulativeTemporality
	sData.IsMonotonic = s.monotonic

	dPts := reset(sData.DataPoints, 0, s.values.Len())

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	s.values.Range(func(_, value any) bool {
		e := value.(*boundEntry[boundCumulativeSumValue[N]])
		if !e.measured.Load() {
			return true
		}

		startTime := s.start
		if perSeriesStartTimeEnabled {
			startTime = e.val.startTime
		}
		newPt := metricdata.DataPoint[N]{
			Attributes: e.val.attrs,
			StartTime:  startTime,
			Time:       t,
			Value:      e.val.n.load(),
		}
		collectExemplars(&newPt.Exemplars, e.val.res.Collect)
		dPts = append(dPts, newPt)
		i++
		return true
	})

	sData.DataPoints = dPts
	*dest = sData

	return i
}

// BoundSum returns a Sum aggregation with bound-instrument support.
func (b Builder[N]) BoundSum(monotonic bool) BoundSum[N] {
	switch b.Temporality {
	case metricdata.DeltaTemporality:
		s := newBoundDeltaSum[N](monotonic, b.AggregationLimit, b.resFunc())
		return BoundSum[N]{
			Measure:            b.filter(s.measure),
			ComputeAggregation: s.collect,
			Bind:               b.bind(s.bind),
		}
	default:
		s := newBoundCumulativeSum[N](monotonic, b.AggregationLimit, b.resFunc())
		return BoundSum[N]{
			Measure:            b.filter(s.measure),
			ComputeAggregation: s.collect,
			Bind:               b.bind(s.bind),
		}
	}
}
