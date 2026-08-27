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

type boundSumPoint[N int64 | float64] struct {
	mu        sync.RWMutex
	value     atomicCounter[N]
	attrs     attribute.Set
	start     time.Time
	active    bool
	bound     atomic.Bool
	reservoir FilteredExemplarReservoir[N]
}

func newBoundSumPoint[N int64 | float64](
	attrs attribute.Set,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *boundSumPoint[N] {
	return &boundSumPoint[N]{
		attrs:     attrs,
		start:     now(),
		active:    true,
		reservoir: reservoir(attrs),
	}
}

func (p *boundSumPoint[N]) add(
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

func (p *boundSumPoint[N]) deactivate() {
	p.mu.Lock()
	p.active = false
	p.mu.Unlock()
}

type boundSumSnapshotMode uint8

const (
	activeSnapshot boundSumSnapshotMode = iota
	finalSnapshot
)

func (p *boundSumPoint[N]) snapshot(
	t time.Time,
	temporality metricdata.Temporality,
	providerStart time.Time,
	mode boundSumSnapshotMode,
) metricdata.DataPoint[N] {
	p.mu.Lock()
	defer p.mu.Unlock()
	value := p.value.load()
	if temporality == metricdata.DeltaTemporality && mode != finalSnapshot {
		p.value.reset()
	}
	start := providerStart
	if temporality == metricdata.DeltaTemporality || sdkx.PerSeriesStartTimestamps.Enabled() {
		start = p.start
	}
	dp := metricdata.DataPoint[N]{
		Attributes: p.attrs,
		StartTime:  start,
		Time:       t,
		Value:      value,
	}
	collectExemplars(&dp.Exemplars, p.reservoir.Collect)
	if temporality == metricdata.DeltaTemporality && mode != finalSnapshot {
		p.start = t
	}
	return dp
}

type boundSum[N int64 | float64] struct {
	mu          sync.Mutex
	active      map[attribute.Distinct]*boundSumPoint[N]
	final       map[attribute.Distinct]*boundSumPoint[N]
	temporality metricdata.Temporality
	monotonic   bool
	limit       int
	reservoir   func(attribute.Set) FilteredExemplarReservoir[N]
	start       time.Time
}

func newBoundSum[N int64 | float64](
	monotonic bool,
	temporality metricdata.Temporality,
	limit int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *boundSum[N] {
	if temporality != metricdata.DeltaTemporality {
		temporality = metricdata.CumulativeTemporality
	}
	return &boundSum[N]{
		active:      make(map[attribute.Distinct]*boundSumPoint[N]),
		final:       make(map[attribute.Distinct]*boundSumPoint[N]),
		temporality: temporality,
		monotonic:   monotonic,
		limit:       limit,
		reservoir:   reservoir,
		start:       now(),
	}
}

type boundSumPointMode uint8

const (
	unboundPoint boundSumPointMode = iota
	boundPoint
)

func (s *boundSum[N]) point(attrs attribute.Set, mode boundSumPointMode) *boundSumPoint[N] {
	distinct := attrs.Equivalent()
	if p, ok := s.active[distinct]; ok {
		if mode == boundPoint {
			p.bound.Store(true)
		}
		return p
	}

	concrete := len(s.active)
	if _, ok := s.active[overflowSet.Equivalent()]; ok {
		concrete--
	}
	if s.limit > 0 && concrete >= s.limit-1 {
		attrs = overflowSet
		distinct = attrs.Equivalent()
		if p, ok := s.active[distinct]; ok {
			if mode == boundPoint {
				p.bound.Store(true)
			}
			return p
		}
	}
	p := newBoundSumPoint(attrs, s.reservoir)
	p.bound.Store(mode == boundPoint)
	s.active[distinct] = p
	return p
}

func (s *boundSum[N]) acquire(attrs attribute.Set, mode boundSumPointMode) *boundSumPoint[N] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.point(attrs, mode)
}

func (s *boundSum[N]) measure(
	ctx context.Context,
	value N,
	attrs attribute.Set,
	dropped []attribute.KeyValue,
) {
	for {
		p := s.acquire(attrs, unboundPoint)
		if p.add(ctx, value, dropped) {
			return
		}
	}
}

type boundSumMeasure[N int64 | float64] struct {
	store    *boundSum[N]
	attrs    attribute.Set
	dropped  []attribute.KeyValue
	point    atomic.Pointer[boundSumPoint[N]]
	slowPath sync.Mutex
}

func (s *boundSum[N]) bind(
	attrs attribute.Set,
	dropped []attribute.KeyValue,
) *boundSumMeasure[N] {
	h := &boundSumMeasure[N]{store: s, attrs: attrs, dropped: slices.Clone(dropped)}
	h.point.Store(s.acquire(attrs, boundPoint))
	return h
}

func (h *boundSumMeasure[N]) Measure(ctx context.Context, value N) {
	p := h.point.Load()
	if p != nil && p.add(ctx, value, h.dropped) {
		return
	}
	h.slowPath.Lock()
	defer h.slowPath.Unlock()
	p = h.point.Load()
	if p != nil && p.add(ctx, value, h.dropped) {
		return
	}
	p = h.store.acquire(h.attrs, boundPoint)
	h.point.Store(p)
	_ = p.add(ctx, value, h.dropped)
}

func (s *boundSum[N]) finish(attrs attribute.Set) {
	if attrs.Equals(&overflowSet) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := attrs.Equivalent()
	p, ok := s.active[key]
	if !ok || p.attrs.Equals(&overflowSet) {
		return
	}
	p.deactivate()
	delete(s.active, key)
	if old, ok := s.final[key]; ok {
		old.value.add(p.value.load())
		return
	}
	s.final[key] = p
}

func (s *boundSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // Required by ComputeAggregation.
) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := now()
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = s.temporality
	sData.IsMonotonic = s.monotonic
	points := reset(sData.DataPoints, 0, len(s.final)+len(s.active))
	for _, p := range s.final {
		points = append(points, p.snapshot(t, s.temporality, s.start, finalSnapshot))
	}
	for key, p := range s.active {
		if _, pendingFinal := s.final[key]; pendingFinal {
			continue
		}
		points = append(points, p.snapshot(t, s.temporality, s.start, activeSnapshot))
		if s.temporality == metricdata.DeltaTemporality && !p.bound.Load() {
			p.deactivate()
			delete(s.active, key)
		}
	}
	clear(s.final)
	sData.DataPoints = points
	*dest = sData
	return len(points)
}

type boundSumBinder[N int64 | float64] struct {
	store  *boundSum[N]
	filter attribute.Filter
}

func (b boundSumBinder[N]) Bind(attrs attribute.Set) BoundMeasure[N] {
	var dropped []attribute.KeyValue
	if b.filter != nil {
		attrs, dropped = attrs.Filter(b.filter)
	}
	return b.store.bind(attrs, dropped)
}

type boundSumFinisher[N int64 | float64] struct {
	store  *boundSum[N]
	filter attribute.Filter
}

func (f boundSumFinisher[N]) Finish(attrs attribute.Set) {
	if f.filter != nil {
		attrs, _ = attrs.Filter(f.filter)
	}
	f.store.finish(attrs)
}
