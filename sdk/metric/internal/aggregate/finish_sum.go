// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type finishSumValue[N int64 | float64] struct {
	lifecycle     finish.Lifecycle
	value         atomicCounter[N]
	attrs         attribute.Set
	start         time.Time
	overflow      atomic.Bool
	dropExemplars bool
	reservoir     FilteredExemplarReservoir[N]
}

func newFinishSumValue[N int64 | float64](
	attrs attribute.Set,
	overflow bool,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) (*finishSumValue[N], finish.Measurement) {
	r := reservoir(attrs)
	_, drop := r.(*dropRes[N])
	point := &finishSumValue[N]{
		attrs:         attrs,
		start:         now(),
		dropExemplars: drop,
		reservoir:     r,
	}
	point.overflow.Store(overflow)
	measurement, _ := point.lifecycle.AcquireMeasurement()
	return point, measurement
}

func (v *finishSumValue[N]) measure(
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

func (v *finishSumValue[N]) measureAcquired(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
	measurement finish.Measurement,
) {
	v.value.add(value)
	if v.dropExemplars {
		measurement.Release()
		return
	}
	defer measurement.Release()
	v.reservoir.Offer(ctx, value, lazy)
}

func (v *finishSumValue[N]) finish(t time.Time) {
	if !v.overflow.Load() {
		v.lifecycle.Finish(t)
	}
}

func (v *finishSumValue[N]) collectCumulative(
	t time.Time,
) (metricdata.DataPoint[N], bool, bool) {
	return v.collect(v.lifecycle.BeginCumulativeCollection(t))
}

func (v *finishSumValue[N]) collectDelta(
	t time.Time,
) (metricdata.DataPoint[N], bool, bool) {
	return v.collect(v.lifecycle.BeginDeltaCollection(t))
}

func (v *finishSumValue[N]) collect(
	collection finish.Collection,
) (metricdata.DataPoint[N], bool, bool) {
	if !collection.ShouldEmit() {
		return metricdata.DataPoint[N]{}, false, false
	}
	defer collection.Complete()

	dp := metricdata.DataPoint[N]{
		Attributes: v.attrs,
		StartTime:  v.start,
		Time:       collection.Time(),
		Value:      v.value.load(),
	}
	collectExemplars(&dp.Exemplars, v.reservoir.Collect)
	return dp, true, collection.ShouldRetire()
}

func (v *finishSumValue[N]) shutdown() {
	v.lifecycle.Retire()
}

// FinishSum contains the operations of a finish-aware Sum aggregation.
type FinishSum[N int64 | float64] struct {
	Measure            Measure[N]
	ComputeAggregation ComputeAggregation
	Finish             func(attribute.Distinct, time.Time)
	Shutdown           func()
}

type finishSum[N int64 | float64] struct {
	collectMu    sync.Mutex
	shutdownOnce sync.Once
	stopped      atomic.Bool

	values      limitedSyncMap[*finishSumValue[N]]
	temporality metricdata.Temporality
	monotonic   bool
	reservoir   func(attribute.Set) FilteredExemplarReservoir[N]
}

func newFinishSum[N int64 | float64](
	monotonic bool,
	temporality metricdata.Temporality,
	limit int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *finishSum[N] {
	if temporality != metricdata.DeltaTemporality {
		temporality = metricdata.CumulativeTemporality
	}
	return &finishSum[N]{
		values: limitedSyncMap[*finishSumValue[N]]{
			aggLimit: limit,
		},
		temporality: temporality,
		monotonic:   monotonic,
		reservoir:   reservoir,
	}
}

func (s *finishSum[N]) measure(
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
		) *finishSumValue[N] {
			var point *finishSumValue[N]
			point, initial = newFinishSumValue(attrs, overflow, s.reservoir)
			return point
		})
		if overflowed {
			point.overflow.Store(true)
		}
		if !loaded {
			point.measureAcquired(ctx, value, lazy, initial)
			if s.stopped.Load() {
				s.retireAndDelete(point)
			}
			return
		}
		if point.measure(ctx, value, lazy) {
			if s.stopped.Load() {
				s.retireAndDelete(point)
			}
			return
		}
	}
}

func (s *finishSum[N]) retireAndDelete(point *finishSumValue[N]) {
	point.shutdown()
	s.values.CompareAndDelete(point.attrs.Equivalent(), point)
}

func (s *finishSum[N]) finish(
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
	point := value.(*finishSumValue[N])
	point.finish(t)
}

func (s *finishSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // Required by ComputeAggregation.
) int {
	s.collectMu.Lock()
	defer s.collectMu.Unlock()
	if s.stopped.Load() {
		return 0
	}

	t := now()
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = s.temporality
	sData.IsMonotonic = s.monotonic
	points := reset(sData.DataPoints, 0, s.values.Len())

	s.values.Range(func(key, raw any) bool {
		point := raw.(*finishSumValue[N])
		var dp metricdata.DataPoint[N]
		var emit, retire bool
		if s.temporality == metricdata.DeltaTemporality {
			dp, emit, retire = point.collectDelta(t)
		} else {
			dp, emit, retire = point.collectCumulative(t)
		}
		if !emit {
			return true
		}
		points = append(points, dp)

		if retire {
			s.values.CompareAndDelete(key.(attribute.Distinct), point)
		}
		return true
	})

	sData.DataPoints = points
	*dest = sData
	return len(points)
}

func (s *finishSum[N]) shutdown() {
	s.shutdownOnce.Do(func() {
		s.stopped.Store(true)
		s.collectMu.Lock()
		defer s.collectMu.Unlock()
		s.values.Range(func(_, raw any) bool {
			raw.(*finishSumValue[N]).shutdown()
			return true
		})
		s.values.Clear()
	})
}

// FinishSum returns a Sum aggregation with exact-attribute lifecycle support.
func (b Builder[N]) FinishSum(monotonic bool) FinishSum[N] {
	store := newFinishSum(
		monotonic,
		b.Temporality,
		b.AggregationLimit,
		b.resFunc(),
	)
	return FinishSum[N]{
		Measure:            b.filter(store.measure),
		ComputeAggregation: store.collect,
		Finish:             store.finish,
		Shutdown:           store.shutdown,
	}
}
