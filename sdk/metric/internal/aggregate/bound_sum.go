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
	once     sync.Once
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
	newValue func(attribute.Set) V,
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

	e := &boundEntry[V]{val: newValue(fltrAttr)}
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

func (m *boundLimitedSyncMap[V]) pinColdLocked( // nolint:revive // Intentional internal control flag
	d attribute.Distinct,
	attrs attribute.Set,
	overflow bool,
	newValue func(attribute.Set) V,
) *boundEntry[V] {
	if actual, loaded := m.Load(d); loaded {
		if e := actual.(*boundEntry[V]); e.bound {
			return e
		}
	} else {
		m.len++
	}
	e := &boundEntry[V]{val: newValue(attrs), bound: true}
	m.Store(d, e)
	if overflow {
		m.boundOverflow = true
		m.hasOverflowEntry = true
		m.overflow.Store(m.atLimitLocked())
	}
	return e
}

func (m *boundLimitedSyncMap[V]) pinLocked(
	lazy lazyFilteredAttributes,
	newValue func(attribute.Set) V,
) (*boundEntry[V], attribute.Distinct, attribute.Set, bool) {
	d := lazy.Distinct()
	if actual, loaded := m.Load(d); loaded {
		e := actual.(*boundEntry[V])
		overflow := m.hasOverflowEntry && d == overflowSet.Equivalent()
		attrs := overflowSet
		if !overflow {
			attrs = lazy.Set()
		}
		if !e.bound {
			e.bound = true
			if overflow {
				m.boundOverflow = true
			}
		}
		return e, d, attrs, overflow
	}
	if m.atLimitLocked() {
		d = overflowSet.Equivalent()
		if actual, loaded := m.Load(d); loaded {
			e := actual.(*boundEntry[V])
			overflow := m.boundOverflow
			if !e.bound {
				e.bound = true
				m.boundOverflow = true
				overflow = true
			}
			m.overflow.Store(true)
			return e, d, overflowSet, overflow
		}
		return m.pinColdLocked(d, overflowSet, true, newValue), d, overflowSet, true
	}
	attrs := lazy.Set()
	return m.pinColdLocked(d, attrs, false, newValue), d, attrs, false
}

func (m *boundLimitedSyncMap[V]) pin(
	lazy lazyFilteredAttributes,
	newValue func(attribute.Set) V,
) *boundEntry[V] {
	m.lenMux.Lock()
	defer m.lenMux.Unlock()
	e, _, _, _ := m.pinLocked(lazy, newValue)
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
	newValue func(attribute.Set) V,
) [2]*boundEntry[V] {
	m.pinMux.Lock()
	defer m.pinMux.Unlock()

	h := m.start()
	defer m.done(h)

	m.hotColdValMap[h].lenMux.Lock()
	eHot, d, attrs, overflow := m.hotColdValMap[h].pinLocked(lazy, newValue)
	m.hotColdValMap[h].lenMux.Unlock()

	m.hotColdValMap[1-h].lenMux.Lock()
	eCold := m.hotColdValMap[1-h].pinColdLocked(d, attrs, overflow, newValue)
	m.hotColdValMap[1-h].lenMux.Unlock()

	var entries [2]*boundEntry[V]
	entries[h] = eHot
	entries[1-h] = eCold
	return entries
}

func newBoundSumValue[N int64 | float64](
	attr attribute.Set,
	newRes func(attribute.Set) FilteredExemplarReservoir[N],
) sumValue[N] {
	r := newRes(attr)
	_, isDrop := r.(*dropRes[N])
	return sumValue[N]{
		res:           r,
		attrs:         attr,
		dropExemplars: isDrop,
	}
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

func (s *boundDeltaSum[N]) newValue(attr attribute.Set) sumValue[N] {
	return newBoundSumValue(attr, s.newRes)
}

func (s *boundDeltaSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	hotIdx := s.vals.start()
	defer s.vals.done(hotIdx)

	e := s.vals.hot(hotIdx).LoadOrStoreAttr(lazy, s.newValue)
	e.val.n.add(value)
	if !e.measured.Load() {
		e.measured.Store(true)
	}
	if !e.val.dropExemplars {
		e.val.res.Offer(ctx, value, lazy)
	}
}

func (s *boundDeltaSum[N]) bind(lazy lazyFilteredAttributes) BoundMeasure[N] {
	entries := s.vals.pin(lazy, s.newValue)
	if entries[0].val.dropExemplars && entries[1].val.dropExemplars {
		return func(_ context.Context, val N) {
			hotIdx := s.vals.start()
			e := entries[hotIdx]
			e.val.n.add(val)
			if !e.measured.Load() {
				e.measured.Store(true)
			}
			s.vals.done(hotIdx)
		}
	}
	return func(ctx context.Context, val N) {
		hotIdx := s.vals.start()
		e := entries[hotIdx]
		e.val.n.add(val)
		if !e.measured.Load() {
			e.measured.Store(true)
		}
		if !e.val.dropExemplars {
			e.val.res.Offer(ctx, val, lazy)
		}
		s.vals.done(hotIdx)
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

type boundCumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	values boundLimitedSyncMap[sumValue[N]]
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
		values: boundLimitedSyncMap[sumValue[N]]{
			limitedSyncMap: limitedSyncMap[*boundEntry[sumValue[N]]]{aggLimit: limit},
		},
		newRes: r,
	}
}

func (s *boundCumulativeSum[N]) newValue(attr attribute.Set) sumValue[N] {
	return newBoundSumValue(attr, s.newRes)
}

func markCumulativeMeasured[N int64 | float64](e *boundEntry[sumValue[N]]) {
	e.once.Do(func() {
		e.val.startTime = now()
		e.measured.Store(true)
	})
}

func (s *boundCumulativeSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	e := s.values.LoadOrStoreAttr(lazy, s.newValue)
	e.val.n.add(value)
	if !e.measured.Load() {
		markCumulativeMeasured(e)
	}
	if !e.val.dropExemplars {
		e.val.res.Offer(ctx, value, lazy)
	}
}

func (s *boundCumulativeSum[N]) bind(lazy lazyFilteredAttributes) BoundMeasure[N] {
	e := s.values.pin(lazy, s.newValue)
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
		e := value.(*boundEntry[sumValue[N]])
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
