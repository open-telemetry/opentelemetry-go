// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type sumValue[N int64 | float64] struct {
	n             atomicCounter[N]
	res           FilteredExemplarReservoir[N]
	attrs         attribute.Set
	startTime     time.Time
	dropExemplars bool
	isBound       bool // true if this entry was created by or used by a bound instrument
	lastReported  N    // last reported value for pinned instruments (delta only)
}

// cumulativeSum is the storage for sums which never reset.
type cumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	newRes      func(attribute.Set) FilteredExemplarReservoir[N]
	values      limitedSyncMap[*sumValue[N]]
	cardinality *cardinalityState
}

func newCumulativeSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeSum[N] {
	state := &cardinalityState{limit: limit}
	return &cumulativeSum[N]{
		monotonic:   monotonic,
		start:       now(),
		newRes:      r,
		values:      limitedSyncMap[*sumValue[N]]{state: state},
		cardinality: state,
	}
}

func (s *cumulativeSum[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	sv := s.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) *sumValue[N] {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
		}
	})
	sv.n.add(value)
	if !sv.dropExemplars {
		sv.res.Offer(ctx, value, droppedAttr)
	}
}

// newDeltaSum returns an aggregator that summarizes a set of measurements as
// their arithmetic sum. Each sum is scoped by attributes and the aggregation
// cycle the measurements were made in.
func newDeltaSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *deltaSum[N] {
	return &deltaSum[N]{
		monotonic: monotonic,
		start:     now(),
		vals:      newHotColdMap[*sumValue[N]](limit),
		newRes:    r,
	}
}

// deltaSum is the storage for sums which resets every collection interval.
type deltaSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	vals *hotColdMap[*sumValue[N]]

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *deltaSum[N]) measure(ctx context.Context, value N, fltrAttr attribute.Set, droppedAttr []attribute.KeyValue) {
	s.vals.WriteUnbound(fltrAttr, func(attr attribute.Set) *sumValue[N] {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
		}
	}, func(sv *sumValue[N]) {
		sv.n.add(value)
		if !sv.dropExemplars {
			sv.res.Offer(ctx, value, droppedAttr)
		}
	})
}

func (s *deltaSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Sum, memory reuse is missed. In that case,
	// use the zero-value sData and hope for better alignment next cycle.
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.DeltaTemporality
	sData.IsMonotonic = s.monotonic

	readIdx := s.vals.SwapHotAndWait()

	// We don't know the total count ahead of time easily because we only collect
	dPts := reset(sData.DataPoints, 0, s.vals.Len(readIdx))

	// 1. Collect from cold map (unbound only)
	s.vals.Collect(readIdx, func(val *sumValue[N]) bool { return val.isBound }, func(_ any, val *sumValue[N]) bool {
		newPt := metricdata.DataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Value:      val.n.load(),
		}
		collectExemplars(&newPt.Exemplars, val.res.Collect)
		dPts = append(dPts, newPt)
		return true
	})

	// 2. Collect from pinned registry (calculating delta using lastReported)
	s.vals.RangePinned(func(_ any, val *sumValue[N]) bool {
		n := val.n.load()
		delta := n - val.lastReported

		newPt := metricdata.DataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Value:      delta,
		}
		collectExemplars(&newPt.Exemplars, val.res.Collect)
		dPts = append(dPts, newPt)

		val.lastReported = n // Update reported value inside entry
		return true
	})

	// The delta collection cycle resets.
	s.start = t

	sData.DataPoints = dPts
	*dest = sData

	return len(dPts)
}

func (s *deltaSum[N]) Bind(attrs attribute.Set) BoundMeasure[N] {
	sv := s.vals.Bind(attrs, func(attr attribute.Set) *sumValue[N] {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
			isBound:       true,
		}
	})
	return sv.boundMeasure()
}

// newCumulativeSum returns an aggregator that summarizes a set of measurements
// as their arithmetic sum. Each sum is scoped by attributes and the
// aggregation cycle the measurements were made in.

func (s *cumulativeSum[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Sum, memory reuse is missed. In that case,
	// use the zero-value sData and hope for better alignment next cycle.
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.CumulativeTemporality
	sData.IsMonotonic = s.monotonic

	// Values are being concurrently written while we iterate, so only use the
	// current length for capacity.
	dPts := reset(sData.DataPoints, 0, s.values.Len())

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	s.values.Range(func(_ any, val *sumValue[N]) bool {
		startTime := s.start
		if perSeriesStartTimeEnabled {
			startTime = val.startTime
		}
		newPt := metricdata.DataPoint[N]{
			Attributes: val.attrs,
			StartTime:  startTime,
			Time:       t,
			Value:      val.n.load(),
		}
		collectExemplars(&newPt.Exemplars, val.res.Collect)
		dPts = append(dPts, newPt)
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
		i++
		return true
	})

	sData.DataPoints = dPts
	*dest = sData

	return i
}

func (s *cumulativeSum[N]) Bind(attrs attribute.Set) BoundMeasure[N] {
	sv := s.values.LoadOrStoreAttr(attrs, func(attr attribute.Set) *sumValue[N] {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
			isBound:       true,
		}
	})
	return sv.boundMeasure()
}

// boundMeasure returns the BoundMeasure recording to v. The exemplar
// decision is made once here rather than per-measurement so the returned
// hot path does not read v's fields, which share a cache line with the
// concurrently written counter.
func (v *sumValue[N]) boundMeasure() BoundMeasure[N] {
	if v.dropExemplars {
		n := &v.n
		return func(_ context.Context, val N) {
			n.add(val)
		}
	}
	return func(ctx context.Context, val N) {
		v.n.add(val)
		v.res.Offer(ctx, val, nil)
	}
}

// newPrecomputedSum returns an aggregator that summarizes a set of
// observations as their arithmetic sum. Each sum is scoped by attributes and
// the aggregation cycle the measurements were made in.
func newPrecomputedSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *precomputedSum[N] {
	return &precomputedSum[N]{
		deltaSum: newDeltaSum[N](monotonic, limit, r),
	}
}

// precomputedSum summarizes a set of observations as their arithmetic sum.
type precomputedSum[N int64 | float64] struct {
	*deltaSum[N]

	reported map[any]N
}

func (s *precomputedSum[N]) delta(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()
	newReported := make(map[any]N)

	// If *dest is not a metricdata.Sum, memory reuse is missed. In that case,
	// use the zero-value sData and hope for better alignment next cycle.
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.DeltaTemporality
	sData.IsMonotonic = s.monotonic

	readIdx := s.vals.SwapHotAndWait()
	n := s.vals.Len(readIdx)
	dPts := reset(sData.DataPoints, 0, n)

	s.vals.Collect(readIdx, func(*sumValue[N]) bool { return false }, func(key any, val *sumValue[N]) bool {
		n := val.n.load()
		delta := n - s.reported[key]

		newPt := metricdata.DataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Value:      delta,
		}
		collectExemplars(&newPt.Exemplars, val.res.Collect)
		dPts = append(dPts, newPt)
		newReported[key] = n
		return true
	})
	s.reported = newReported
	// The delta collection cycle resets.
	s.start = t

	sData.DataPoints = dPts
	*dest = sData

	return len(dPts)
}

func (s *precomputedSum[N]) cumulative(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Sum, memory reuse is missed. In that case,
	// use the zero-value sData and hope for better alignment next cycle.
	sData, _ := (*dest).(metricdata.Sum[N])
	sData.Temporality = metricdata.CumulativeTemporality
	sData.IsMonotonic = s.monotonic

	readIdx := s.vals.SwapHotAndWait()
	n := s.vals.Len(readIdx)
	dPts := reset(sData.DataPoints, 0, n)

	s.vals.Collect(readIdx, func(*sumValue[N]) bool { return false }, func(_ any, val *sumValue[N]) bool {
		newPt := metricdata.DataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Value:      val.n.load(),
		}
		collectExemplars(&newPt.Exemplars, val.res.Collect)
		dPts = append(dPts, newPt)
		return true
	})

	sData.DataPoints = dPts
	*dest = sData

	return len(dPts)
}
