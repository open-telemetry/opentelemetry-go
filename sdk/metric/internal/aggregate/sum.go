// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

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
	isBound       bool
	lastReported  N
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
	lazy := newLazyFilteredAttributes(v.attrs, nil)
	return func(ctx context.Context, val N) {
		v.n.add(val)
		v.res.Offer(ctx, val, lazy)
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

func (s *deltaSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	hotIdx := s.vals.start()
	defer s.vals.done(hotIdx)

	sv := s.vals.LoadOrStoreHot(hotIdx, lazy, func(attr attribute.Set) *sumValue[N] {
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
		sv.res.Offer(ctx, value, lazy)
	}
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

	n := s.vals.Len(readIdx) + s.vals.PinnedLen()
	dPts := reset(sData.DataPoints, n, n)

	var i int
	// 1. Collect from cold map (unbound only)
	s.vals.Collect(readIdx, func(val *sumValue[N]) bool { return val.isBound }, func(_ any, val *sumValue[N]) bool {
		if i < len(dPts) {
			dPts[i].Attributes = val.attrs
			dPts[i].StartTime = s.start
			dPts[i].Time = t
			dPts[i].Value = val.n.load()
			collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		} else {
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      val.n.load(),
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
		}
		i++
		return true
	})

	// 2. Collect from pinned registry (calculating delta using lastReported)
	s.vals.RangePinned(func(_ any, val *sumValue[N]) bool {
		n := val.n.load()
		delta := n - val.lastReported
		if delta == 0 {
			return true
		}

		if i < len(dPts) {
			dPts[i].Attributes = val.attrs
			dPts[i].StartTime = s.start
			dPts[i].Time = t
			dPts[i].Value = delta
			collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		} else {
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      delta,
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
		}
		i++

		val.lastReported = n // Update reported value inside entry
		return true
	})

	dPts = dPts[:i]

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
	}, func(val *sumValue[N]) {
		val.isBound = true
	})
	return sv.boundMeasure()
}

// newCumulativeSum returns an aggregator that summarizes a set of measurements
// as their arithmetic sum. Each sum is scoped by attributes and the
// aggregation cycle the measurements were made in.
func newCumulativeSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeSum[N] {
	return &cumulativeSum[N]{
		monotonic: monotonic,
		start:     now(),
		values:    limitedSyncMap[*sumValue[N]]{aggLimit: limit},
		newRes:    r,
	}
}

// cumulativeSum is the storage for sums which never reset.
type cumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	values limitedSyncMap[*sumValue[N]]
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *cumulativeSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	sv := s.values.LoadOrStoreAttr(lazy, func(attr attribute.Set) *sumValue[N] {
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
	// It is possible for collection to race with measurement and observe the
	// exemplar in the batch of metrics after the add() for cumulative sums.
	// This is an accepted tradeoff to avoid locking during measurement.
	if !sv.dropExemplars {
		sv.res.Offer(ctx, value, lazy)
	}
}

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
	n := s.values.Len()
	dPts := reset(sData.DataPoints, n, n)

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	s.values.Range(func(_, value any) bool {
		val := value.(*sumValue[N])

		startTime := s.start
		if perSeriesStartTimeEnabled {
			startTime = val.startTime
		}
		if i < len(dPts) {
			dPts[i].Attributes = val.attrs
			dPts[i].StartTime = startTime
			dPts[i].Time = t
			dPts[i].Value = val.n.load()
			collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		} else {
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  startTime,
				Time:       t,
				Value:      val.n.load(),
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
		}
		i++
		return true
	})
	dPts = dPts[:i]

	sData.DataPoints = dPts
	*dest = sData

	return i
}

func (s *cumulativeSum[N]) Bind(attrs attribute.Set) BoundMeasure[N] {
	sv := s.values.LoadOrStoreAttr(newLazyFilteredAttributes(attrs, nil), func(attr attribute.Set) *sumValue[N] {
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
	dPts := reset(sData.DataPoints, n, n)

	var i int
	s.vals.Collect(readIdx, func(*sumValue[N]) bool { return false }, func(key any, val *sumValue[N]) bool {
		n := val.n.load()
		delta := n - s.reported[key]

		if i < len(dPts) {
			dPts[i].Attributes = val.attrs
			dPts[i].StartTime = s.start
			dPts[i].Time = t
			dPts[i].Value = delta
			collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		} else {
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      delta,
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
		}
		i++
		newReported[key] = n
		return true
	})
	dPts = dPts[:i]
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
	dPts := reset(sData.DataPoints, n, n)

	var i int
	s.vals.Collect(readIdx, func(*sumValue[N]) bool { return false }, func(_ any, val *sumValue[N]) bool {
		if i < len(dPts) {
			dPts[i].Attributes = val.attrs
			dPts[i].StartTime = s.start
			dPts[i].Time = t
			dPts[i].Value = val.n.load()
			collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		} else {
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      val.n.load(),
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
		}
		i++
		return true
	})
	dPts = dPts[:i]

	sData.DataPoints = dPts
	*dest = sData

	return len(dPts)
}
