// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type finishLastValuePoint[N int64 | float64] struct {
	lifecycle     finish.Lifecycle
	value         atomicN[N]
	attrs         attribute.Set
	start         time.Time
	dropExemplars bool
	reservoir     FilteredExemplarReservoir[N]
}

func newFinishLastValuePoint[N int64 | float64](
	attrs attribute.Set,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *finishLastValuePoint[N] {
	r := reservoir(attrs)
	_, drop := r.(*dropRes[N])
	point := &finishLastValuePoint[N]{
		attrs:         attrs,
		start:         now(),
		dropExemplars: drop,
		reservoir:     r,
	}
	return point
}

func (v *finishLastValuePoint[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) bool {
	measurement, ok := v.lifecycle.AcquireMeasurement()
	if !ok {
		return false
	}
	v.measureAcquired(ctx, value, lazy, measurement)
	return true
}

func (v *finishLastValuePoint[N]) measureOverflow(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) bool {
	measurement, ok := v.lifecycle.AcquireSharedMeasurement()
	if !ok {
		return false
	}
	v.measureAcquired(ctx, value, lazy, measurement)
	return true
}

func (v *finishLastValuePoint[N]) measureAcquired(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
	measurement finish.Measurement,
) {
	v.value.Store(value)
	if v.dropExemplars {
		measurement.Release()
		return
	}
	defer measurement.Release()
	v.reservoir.Offer(ctx, value, lazy)
}

func (v *finishLastValuePoint[N]) finish(t time.Time) {
	v.lifecycle.Finish(t)
}

func (v *finishLastValuePoint[N]) collectCumulative(
	t time.Time,
) (metricdata.DataPoint[N], bool, bool) {
	collection, ok := v.lifecycle.BeginCumulativeCollection(t)
	if !ok {
		return metricdata.DataPoint[N]{}, false, false
	}
	dp, retire := v.collect(v.start, collection)
	return dp, true, retire
}

func (v *finishLastValuePoint[N]) collectDelta(
	start time.Time,
	t time.Time,
) (metricdata.DataPoint[N], bool, bool) {
	collection, ok := v.lifecycle.BeginDeltaCollection(t)
	if !ok {
		return metricdata.DataPoint[N]{}, false, false
	}
	// A deferred series can predate the previous collection's boundary.
	if v.start.Before(start) {
		start = v.start
	}
	dp, retire := v.collect(start, collection)
	return dp, true, retire
}

func (v *finishLastValuePoint[N]) collect(
	start time.Time,
	collection finish.Collection,
) (metricdata.DataPoint[N], bool) {
	defer collection.Complete()

	dp := metricdata.DataPoint[N]{
		Attributes: v.attrs,
		StartTime:  start,
		Time:       collection.Time(),
		Value:      v.value.Load(),
	}
	collectExemplars(&dp.Exemplars, v.reservoir.Collect)
	return dp, collection.ShouldRetire()
}

func (v *finishLastValuePoint[N]) shutdown(ctx context.Context) error {
	return v.lifecycle.Retire(ctx)
}

// FinishLastValue contains the operations of a finish-aware LastValue aggregation.
// Finish accepts a post-filter series identity and the timestamp for its final point.
//
// Stop disables subsequent Measure, Finish, and ComputeAggregation calls.
// Wait calls Stop, retires series observed during cleanup, and clears the map.
// A recording overlapping Stop may complete after Wait returns. If it publishes
// a series missed by cleanup, it retires that series when recording completes.
// Wait returns a context error if cleanup is interrupted by cancellation.
type FinishLastValue[N int64 | float64] struct {
	Measure            Measure[N]
	ComputeAggregation ComputeAggregation
	Finish             func(attribute.Distinct, time.Time)
	Stop               func()
	Wait               func(context.Context) error
}

type finishLastValue[N int64 | float64] struct {
	collectSem chan struct{}
	stopped    atomic.Bool
	points     []*finishLastValuePoint[N]

	values      limitedSyncMap[*finishLastValuePoint[N]]
	start       time.Time
	temporality metricdata.Temporality
	reservoir   func(attribute.Set) FilteredExemplarReservoir[N]
}

func newFinishLastValue[N int64 | float64](
	temporality metricdata.Temporality,
	limit int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *finishLastValue[N] {
	var start time.Time
	if temporality != metricdata.DeltaTemporality {
		temporality = metricdata.CumulativeTemporality
	} else {
		start = now()
	}
	return &finishLastValue[N]{
		collectSem: make(chan struct{}, 1),
		values: limitedSyncMap[*finishLastValuePoint[N]]{
			aggLimit: limit,
		},
		start:       start,
		temporality: temporality,
		reservoir:   reservoir,
	}
}

func (s *finishLastValue[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) {
	for {
		if s.stopped.Load() {
			return
		}
		var initial finish.Measurement
		point, loaded, overflowed := s.values.LoadOrStoreAttrReclaiming(lazy, func(
			attrs attribute.Set,
			overflow bool,
		) *finishLastValuePoint[N] {
			point := newFinishLastValuePoint(attrs, s.reservoir)
			// The point has not been published, so its new lifecycle cannot be
			// retired and this initial measurement admission cannot fail.
			if overflow {
				initial, _ = point.lifecycle.AcquireSharedMeasurement()
			} else {
				initial, _ = point.lifecycle.AcquireMeasurement()
			}
			return point
		})
		if !loaded {
			point.measureAcquired(ctx, value, lazy, initial)
			if s.stopped.Load() {
				s.retireAndDelete(point)
			}
			return
		}
		if overflowed {
			if point.measureOverflow(ctx, value, lazy) {
				if s.stopped.Load() {
					s.retireAndDelete(point)
				}
				return
			}
			continue
		}
		if point.measure(ctx, value, lazy) {
			if s.stopped.Load() {
				s.retireAndDelete(point)
			}
			return
		}
		s.values.CompareAndDelete(lazy.Distinct(), point)
	}
}

func (s *finishLastValue[N]) retireAndDelete(point *finishLastValuePoint[N]) {
	_ = point.shutdown(context.Background())
	s.values.CompareAndDelete(point.attrs.Equivalent(), point)
}

func (s *finishLastValue[N]) finish(
	distinct attribute.Distinct,
	t time.Time,
) {
	if s.stopped.Load() {
		return
	}
	value, ok := s.values.Load(distinct)
	if !ok {
		return
	}
	point := value.(*finishLastValuePoint[N])
	point.finish(t)
}

func (s *finishLastValue[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // Required by ComputeAggregation.
) int {
	s.collectSem <- struct{}{}
	defer func() { <-s.collectSem }()
	if s.stopped.Load() {
		return 0
	}

	// Every included series must start at or before the collection timestamp.
	snapshot := s.snapshot()
	defer clear(snapshot)
	t := now()
	start := s.start
	if s.temporality == metricdata.DeltaTemporality {
		s.start = t
	}
	gData, _ := (*dest).(metricdata.Gauge[N])
	points := reset(gData.DataPoints, 0, len(snapshot))

	for _, point := range snapshot {
		var dp metricdata.DataPoint[N]
		var emit, retire bool
		if s.temporality == metricdata.DeltaTemporality {
			dp, emit, retire = point.collectDelta(start, t)
		} else {
			dp, emit, retire = point.collectCumulative(t)
		}
		if !emit {
			continue
		}
		points = append(points, dp)

		if retire {
			s.values.CompareAndDelete(point.attrs.Equivalent(), point)
		}
	}

	gData.DataPoints = points
	*dest = gData
	return len(points)
}

// Snapshot before retiring series so newly reclaimed slots cannot expand this
// collection beyond the cardinality limit. collectSem protects the scratch slice.
func (s *finishLastValue[N]) snapshot() []*finishLastValuePoint[N] {
	s.values.lenMux.Lock()
	defer s.values.lenMux.Unlock()
	s.points = reset(s.points, 0, s.values.len)
	s.values.Range(func(_, raw any) bool {
		s.points = append(s.points, raw.(*finishLastValuePoint[N]))
		return true
	})
	return s.points
}

func (s *finishLastValue[N]) stop() {
	s.stopped.Store(true)
}

func (s *finishLastValue[N]) wait(ctx context.Context) error {
	s.stop()
	if err := s.lockCollection(ctx); err != nil {
		return err
	}
	defer func() { <-s.collectSem }()
	s.points = nil

	var err error
	s.values.Range(func(_, raw any) bool {
		err = raw.(*finishLastValuePoint[N]).shutdown(ctx)
		return err == nil
	})
	if err == nil {
		err = s.values.ClearContext(ctx)
	}
	return err
}

func (s *finishLastValue[N]) lockCollection(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case s.collectSem <- struct{}{}:
		if err := ctx.Err(); err != nil {
			<-s.collectSem
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FinishLastValue returns a LastValue aggregation with exact-attribute lifecycle support.
func (b Builder[N]) FinishLastValue() FinishLastValue[N] {
	store := newFinishLastValue(
		b.Temporality,
		b.AggregationLimit,
		b.resFunc(),
	)
	return FinishLastValue[N]{
		Measure:            b.filter(store.measure),
		ComputeAggregation: store.collect,
		Finish:             store.finish,
		Stop:               store.stop,
		Wait:               store.wait,
	}
}
