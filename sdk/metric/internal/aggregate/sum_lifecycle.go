// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	sdkx "go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type lifecycleSumPoint[N int64 | float64] struct {
	mu        sync.RWMutex
	value     atomicCounter[N]
	attrs     attribute.Set
	start     time.Time
	active    bool
	bound     atomic.Bool
	reservoir FilteredExemplarReservoir[N]
}

func newLifecycleSumPoint[N int64 | float64](
	attrs attribute.Set,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *lifecycleSumPoint[N] {
	return &lifecycleSumPoint[N]{
		attrs:     attrs,
		start:     now(),
		active:    true,
		reservoir: reservoir(attrs),
	}
}

func (p *lifecycleSumPoint[N]) add(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if !p.active {
		return false
	}
	p.value.add(value)
	p.reservoir.Offer(ctx, value, lazy)
	return true
}

func (p *lifecycleSumPoint[N]) deactivate() {
	p.mu.Lock()
	p.active = false
	p.mu.Unlock()
}

type lifecycleSum[N int64 | float64] struct {
	mu          sync.Mutex
	active      map[attribute.Distinct]*lifecycleSumPoint[N]
	final       map[attribute.Distinct][]*lifecycleSumPoint[N]
	temporality metricdata.Temporality
	monotonic   bool
	limit       int
	reservoir   func(attribute.Set) FilteredExemplarReservoir[N]
	start       time.Time
}

func newLifecycleSum[N int64 | float64](
	monotonic bool,
	temporality metricdata.Temporality,
	limit int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *lifecycleSum[N] {
	if temporality != metricdata.DeltaTemporality {
		temporality = metricdata.CumulativeTemporality
	}
	return &lifecycleSum[N]{
		active:      make(map[attribute.Distinct]*lifecycleSumPoint[N]),
		final:       make(map[attribute.Distinct][]*lifecycleSumPoint[N]),
		temporality: temporality,
		monotonic:   monotonic,
		limit:       limit,
		reservoir:   reservoir,
		start:       now(),
	}
}

func (s *lifecycleSum[N]) point(lazy lazyFilteredAttributes) *lifecycleSumPoint[N] {
	key := lazy.Distinct()
	if point, ok := s.active[key]; ok {
		return point
	}

	concrete := len(s.active)
	if _, ok := s.active[overflowSet.Equivalent()]; ok {
		concrete--
	}
	var attrs attribute.Set
	if s.limit > 0 && concrete >= s.limit-1 {
		attrs = overflowSet
		key = attrs.Equivalent()
		if point, ok := s.active[key]; ok {
			return point
		}
	} else {
		attrs = lazy.Set()
	}
	point := newLifecycleSumPoint(attrs, s.reservoir)
	s.active[key] = point
	return point
}

func (s *lifecycleSum[N]) acquire(lazy lazyFilteredAttributes) *lifecycleSumPoint[N] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.point(lazy)
}

func (s *lifecycleSum[N]) bindPoint(lazy lazyFilteredAttributes) *lifecycleSumPoint[N] {
	s.mu.Lock()
	defer s.mu.Unlock()
	point := s.point(lazy)
	point.bound.Store(true)
	return point
}

func (s *lifecycleSum[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) {
	for {
		if s.acquire(lazy).add(ctx, value, lazy) {
			return
		}
	}
}

type lifecycleSumMeasure[N int64 | float64] struct {
	store    *lifecycleSum[N]
	lazy     lazyFilteredAttributes
	point    atomic.Pointer[lifecycleSumPoint[N]]
	slowPath sync.Mutex
}

func (s *lifecycleSum[N]) bind(lazy lazyFilteredAttributes) func(context.Context, N) {
	handle := &lifecycleSumMeasure[N]{
		store: s,
		lazy:  lazy,
	}
	handle.point.Store(s.bindPoint(lazy))
	return handle.measure
}

func (h *lifecycleSumMeasure[N]) measure(ctx context.Context, value N) {
	point := h.point.Load()
	if point != nil && point.add(ctx, value, h.lazy) {
		return
	}
	h.slowPath.Lock()
	defer h.slowPath.Unlock()
	point = h.point.Load()
	if point != nil && point.add(ctx, value, h.lazy) {
		return
	}
	point = h.store.bindPoint(h.lazy)
	h.point.Store(point)
	_ = point.add(ctx, value, h.lazy)
}

func (s *lifecycleSum[N]) finish(lazy lazyFilteredAttributes) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := lazy.Distinct()
	point, ok := s.active[key]
	if !ok || point.attrs.Equals(&overflowSet) {
		return
	}
	point.deactivate()
	delete(s.active, key)
	s.final[key] = append(s.final[key], point)
}

func (s *lifecycleSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // Required by ComputeAggregation.
) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := now()
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = s.temporality
	sData.IsMonotonic = s.monotonic
	points := reset(sData.DataPoints, 0, len(s.final)+len(s.active))
	for _, values := range s.final {
		points = append(points, s.finalDataPoint(values, t))
	}
	for key, point := range s.active {
		if _, pendingFinal := s.final[key]; pendingFinal {
			continue
		}
		points = append(points, s.activeDataPoint(point, t))
		if s.temporality == metricdata.DeltaTemporality && !point.bound.Load() {
			point.deactivate()
			delete(s.active, key)
		}
	}
	clear(s.final)
	s.start = t
	sData.DataPoints = points
	*dest = sData
	return len(points)
}

func (s *lifecycleSum[N]) finalDataPoint(
	points []*lifecycleSumPoint[N],
	t time.Time,
) metricdata.DataPoint[N] {
	dp := metricdata.DataPoint[N]{Time: t, StartTime: s.start}
	perSeriesStart := s.temporality == metricdata.DeltaTemporality ||
		sdkx.PerSeriesStartTimestamps.Enabled()
	for index, point := range points {
		point.mu.Lock()
		if index == 0 {
			dp.Attributes = point.attrs
			if perSeriesStart {
				dp.StartTime = point.start
			}
		} else if perSeriesStart && point.start.Before(dp.StartTime) {
			dp.StartTime = point.start
		}
		dp.Value += point.value.load()
		appendSumExemplars(&dp.Exemplars, point.reservoir)
		point.mu.Unlock()
	}
	return dp
}

func (s *lifecycleSum[N]) activeDataPoint(
	point *lifecycleSumPoint[N],
	t time.Time,
) metricdata.DataPoint[N] {
	point.mu.Lock()
	defer point.mu.Unlock()
	start := s.start
	if s.temporality == metricdata.DeltaTemporality || sdkx.PerSeriesStartTimestamps.Enabled() {
		start = point.start
	}
	dp := metricdata.DataPoint[N]{
		Attributes: point.attrs,
		StartTime:  start,
		Time:       t,
		Value:      point.value.load(),
	}
	appendSumExemplars(&dp.Exemplars, point.reservoir)
	if s.temporality == metricdata.DeltaTemporality {
		point.value.reset()
		point.start = t
	}
	return dp
}

func appendSumExemplars[N int64 | float64](
	dest *[]metricdata.Exemplar[N],
	reservoir FilteredExemplarReservoir[N],
) {
	var exemplars []metricdata.Exemplar[N]
	collectExemplars(&exemplars, reservoir.Collect)
	*dest = append(*dest, exemplars...)
}

type boundSumAggregator[N int64 | float64] struct {
	measure Measure[N]
	store   *lifecycleSum[N]
	filter  attribute.Filter
}

func (a *boundSumAggregator[N]) Measure() Measure[N] {
	return a.measure
}

func (a *boundSumAggregator[N]) ComputeAggregation() ComputeAggregation {
	return a.store.collect
}

func (a *boundSumAggregator[N]) Bind(attrs attribute.Set) func(context.Context, N) {
	return a.store.bind(newLazyFilteredAttributes(attrs, a.filter))
}

func (a *boundSumAggregator[N]) Finish(attrs attribute.Set) {
	a.store.finish(newLazyFilteredAttributes(attrs, a.filter))
}
