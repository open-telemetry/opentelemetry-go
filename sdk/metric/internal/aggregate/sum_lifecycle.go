// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"slices"
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
	dropped []attribute.KeyValue,
) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if !p.active {
		return false
	}
	p.value.add(value)
	p.reservoir.Offer(ctx, value, dropped)
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

func (s *lifecycleSum[N]) point(attrs attribute.Set) *lifecycleSumPoint[N] {
	key := attrs.Equivalent()
	if point, ok := s.active[key]; ok {
		return point
	}

	concrete := len(s.active)
	if _, ok := s.active[overflowSet.Equivalent()]; ok {
		concrete--
	}
	if s.limit > 0 && concrete >= s.limit-1 {
		attrs = overflowSet
		key = attrs.Equivalent()
		if point, ok := s.active[key]; ok {
			return point
		}
	}
	point := newLifecycleSumPoint(attrs, s.reservoir)
	s.active[key] = point
	return point
}

func (s *lifecycleSum[N]) acquire(attrs attribute.Set) *lifecycleSumPoint[N] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.point(attrs)
}

func (s *lifecycleSum[N]) bindPoint(attrs attribute.Set) *lifecycleSumPoint[N] {
	s.mu.Lock()
	defer s.mu.Unlock()
	point := s.point(attrs)
	point.bound.Store(true)
	return point
}

func (s *lifecycleSum[N]) measure(
	ctx context.Context,
	value N,
	attrs attribute.Set,
	dropped []attribute.KeyValue,
) {
	for {
		if s.acquire(attrs).add(ctx, value, dropped) {
			return
		}
	}
}

type lifecycleSumMeasure[N int64 | float64] struct {
	store    *lifecycleSum[N]
	attrs    attribute.Set
	dropped  []attribute.KeyValue
	point    atomic.Pointer[lifecycleSumPoint[N]]
	slowPath sync.Mutex
}

func (s *lifecycleSum[N]) bind(
	attrs attribute.Set,
	dropped []attribute.KeyValue,
) func(context.Context, N) {
	handle := &lifecycleSumMeasure[N]{
		store:   s,
		attrs:   attrs,
		dropped: slices.Clone(dropped),
	}
	handle.point.Store(s.bindPoint(attrs))
	return handle.measure
}

func (h *lifecycleSumMeasure[N]) measure(ctx context.Context, value N) {
	point := h.point.Load()
	if point != nil && point.add(ctx, value, h.dropped) {
		return
	}
	h.slowPath.Lock()
	defer h.slowPath.Unlock()
	point = h.point.Load()
	if point != nil && point.add(ctx, value, h.dropped) {
		return
	}
	point = h.store.bindPoint(h.attrs)
	h.point.Store(point)
	_ = point.add(ctx, value, h.dropped)
}

func (s *lifecycleSum[N]) finish(attrs attribute.Set) {
	if attrs.Equals(&overflowSet) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := attrs.Equivalent()
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
	var dropped []attribute.KeyValue
	if a.filter != nil {
		attrs, dropped = attrs.Filter(a.filter)
	}
	return a.store.bind(attrs, dropped)
}

func (a *boundSumAggregator[N]) Finish(attrs attribute.Set) {
	if a.filter != nil {
		attrs, _ = attrs.Filter(a.filter)
	}
	a.store.finish(attrs)
}
