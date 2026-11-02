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
}

// newDeltaSum returns an aggregator that summarizes a set of measurements as
// their arithmetic sum. Each sum is scoped by attributes and the aggregation
// cycle the measurements were made in.
func newDeltaSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *deltaSum[N] {
	s := &deltaSum[N]{
		monotonic: monotonic,
		start:     now(),
		newRes:    r,
	}
	s.vals.init(limit)
	return s
}

// deltaSum is the storage for sums which resets every collection interval.
type deltaSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	vals   hotColdMap[*sumValue[N]]
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *deltaSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	hotIdx := s.vals.start()
	defer s.vals.done(hotIdx)

	sv := s.vals.hot(hotIdx).LoadOrStoreAttr(lazy, func(attr attribute.Set) *sumValue[N] {
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

	// delta always clears values on collection
	readIdx := s.vals.swapHotAndWait()
	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := s.vals.Len(readIdx)
	dPts := reset(sData.DataPoints, n, n)

	var i int
	s.vals.Range(readIdx, func(_, value any) bool {
		val := value.(*sumValue[N])
		collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		dPts[i].Attributes = val.attrs
		dPts[i].StartTime = s.start
		dPts[i].Time = t
		dPts[i].Value = val.n.load()
		i++
		return true
	})
	// Unused attribute sets do not report.
	s.vals.Clear(readIdx)

	// The delta collection cycle resets.
	s.start = t

	sData.DataPoints = dPts
	*dest = sData

	return i
}

// newCumulativeSum returns an aggregator that summarizes a set of measurements
// as their arithmetic sum. Each sum is scoped by attributes and the
// aggregation cycle the measurements were made in.
func newCumulativeSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeSum[N] {
	s := &cumulativeSum[N]{
		monotonic: monotonic,
		start:     now(),
	}
	var zero N
	if _, isFloat64 := any(zero).(float64); isFloat64 {
		// Only float64 needs separate buffers; int64 snapshots are atomic.
		s.float64Values = &float64CumulativeSum[N]{
			newRes: r,
			values: limitedSyncMap[*cumulativeSumValue[N]]{aggLimit: limit},
		}
	} else {
		s.values = limitedSyncMap[*sumValue[N]]{aggLimit: limit}
		s.newRes = r
	}
	return s
}

type cumulativeSumValue[N int64 | float64] struct {
	hcwg          hotColdWaitGroup
	counters      [2]atomicCounter[N]
	res           FilteredExemplarReservoir[N]
	attrs         attribute.Set
	startTime     time.Time
	dropExemplars bool
}

// float64CumulativeSum uses a hot and a cold counter for each series so
// collection can read a stable snapshot while measurements continue.
type float64CumulativeSum[N int64 | float64] struct {
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	values limitedSyncMap[*cumulativeSumValue[N]]
}

func (s *float64CumulativeSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	sum := s.values.LoadOrStoreAttr(lazy, func(attr attribute.Set) *cumulativeSumValue[N] {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &cumulativeSumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
		}
	})

	hotIdx := sum.hcwg.start()
	sum.counters[hotIdx].add(value)
	sum.hcwg.done(hotIdx)
	if !sum.dropExemplars {
		sum.res.Offer(ctx, value, lazy)
	}
}

// cumulativeSum is the storage for sums which never reset.
type cumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	values        limitedSyncMap[*sumValue[N]]
	newRes        func(attribute.Set) FilteredExemplarReservoir[N]
	float64Values *float64CumulativeSum[N]
}

func (s *cumulativeSum[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	if s.float64Values != nil {
		s.float64Values.measure(ctx, value, lazy)
		return
	}
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

	values := &s.values.Map
	n := s.values.Len()
	if s.float64Values != nil {
		values = &s.float64Values.values.Map
		n = s.float64Values.values.Len()
	}
	// Values may be added while we iterate, so only use the current length for
	// capacity.
	dPts := reset(sData.DataPoints, 0, n)

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	values.Range(func(_, value any) bool {
		var (
			attrs       attribute.Set
			seriesStart time.Time
			res         FilteredExemplarReservoir[N]
			total       N
		)
		if s.float64Values != nil {
			val := value.(*cumulativeSumValue[N])
			readIdx := val.hcwg.swapHotAndWait()
			total = val.counters[readIdx].load()
			hotIdx := (readIdx + 1) % 2
			val.counters[hotIdx].add(total)
			val.counters[readIdx].reset()
			attrs = val.attrs
			seriesStart = val.startTime
			res = val.res
		} else {
			val := value.(*sumValue[N])
			total = val.n.load()
			attrs = val.attrs
			seriesStart = val.startTime
			res = val.res
		}

		startTime := s.start
		if perSeriesStartTimeEnabled {
			startTime = seriesStart
		}
		newPt := metricdata.DataPoint[N]{
			Attributes: attrs,
			StartTime:  startTime,
			Time:       t,
			Value:      total,
		}
		collectExemplars(&newPt.Exemplars, res.Collect)
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

	// delta always clears values on collection
	readIdx := s.vals.swapHotAndWait()
	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := s.vals.Len(readIdx)
	dPts := reset(sData.DataPoints, n, n)

	var i int
	s.vals.Range(readIdx, func(key, value any) bool {
		val := value.(*sumValue[N])
		n := val.n.load()

		delta := n - s.reported[key]
		collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		dPts[i].Attributes = val.attrs
		dPts[i].StartTime = s.start
		dPts[i].Time = t
		dPts[i].Value = delta
		newReported[key] = n
		i++
		return true
	})
	// Unused attribute sets do not report.
	s.vals.Clear(readIdx)
	s.reported = newReported
	// The delta collection cycle resets.
	s.start = t

	sData.DataPoints = dPts
	*dest = sData

	return i
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

	// cumulative precomputed always clears values on collection
	readIdx := s.vals.swapHotAndWait()
	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := s.vals.Len(readIdx)
	dPts := reset(sData.DataPoints, n, n)

	var i int
	s.vals.Range(readIdx, func(_, value any) bool {
		val := value.(*sumValue[N])
		collectExemplars(&dPts[i].Exemplars, val.res.Collect)
		dPts[i].Attributes = val.attrs
		dPts[i].StartTime = s.start
		dPts[i].Time = t
		dPts[i].Value = val.n.load()
		i++
		return true
	})
	// Unused attribute sets do not report.
	s.vals.Clear(readIdx)

	sData.DataPoints = dPts
	*dest = sData

	return i
}
