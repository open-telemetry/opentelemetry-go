// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
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

	// boundFloat64 caches the bound instrument to avoid allocations on the fast path.
	// It is only populated when N is float64.
	boundFloat64 metric.Float64Counter

	// boundInt64 caches the bound instrument to avoid allocations on the fast path.
	// It is only populated when N is int64.
	boundInt64 metric.Int64Counter
}

// boundFloat64SumValue implements metric.Float64Counter for a specific sumValue.
type boundFloat64SumValue struct {
	sv *sumValue[float64]
	embedded.Float64Counter
}

func (b *boundFloat64SumValue) Add(ctx context.Context, val float64, _ ...metric.AddOption) {
	b.sv.n.add(val)
	if !b.sv.dropExemplars {
		b.sv.res.Offer(ctx, val, nil)
	}
}

func (*boundFloat64SumValue) Enabled(_ context.Context) bool {
	return true
}

// boundInt64SumValue implements metric.Int64Counter for a specific sumValue.
type boundInt64SumValue struct {
	sv *sumValue[int64]
	embedded.Int64Counter
}

func (b *boundInt64SumValue) Add(ctx context.Context, val int64, _ ...metric.AddOption) {
	b.sv.n.add(val)
	if !b.sv.dropExemplars {
		b.sv.res.Offer(ctx, val, nil)
	}
}

func (*boundInt64SumValue) Enabled(_ context.Context) bool {
	return true
}

type sumValueMap[N int64 | float64] struct {
	values limitedSyncMap
	newRes func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *sumValueMap[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	sv := s.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
		}
	}).(*sumValue[N])
	sv.n.add(value)
	// It is possible for collection to race with measurement and observe the
	// exemplar in the batch of metrics after the add() for cumulative sums.
	// This is an accepted tradeoff to avoid locking during measurement.
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
	filter attribute.Filter,
) *deltaSum[N] {
	return &deltaSum[N]{
		monotonic: monotonic,
		start:     now(),
		vals:      NewHotColdSyncMap(limit),
		reported:  make(map[any]N),
		filter:    filter,
		newRes:    r,
	}
}

// deltaSum is the storage for sums which resets every collection interval.
type deltaSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	vals *HotColdSyncMap

	reported map[any]N // last reported values for pinned instruments
	filter   attribute.Filter
	newRes   func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *deltaSum[N]) measure(ctx context.Context, value N, fltrAttr attribute.Set, droppedAttr []attribute.KeyValue) {
	sv := s.vals.LoadOrStoreUnbound(fltrAttr, func(attr attribute.Set) any {
		r := s.newRes(attr)
		_, isDrop := r.(*dropRes[N])
		return &sumValue[N]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
		}
	}).(*sumValue[N])
	
	sv.n.add(value)
	if !sv.dropExemplars {
		sv.res.Offer(ctx, value, droppedAttr)
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
	readIdx := s.vals.SwapHotAndWait()
	hotIdx := 1 - readIdx
	
	// We don't know the total count ahead of time easily because we only collect
	// bound entries from hot map, and all from cold map.
	dPts := reset(sData.DataPoints, 0, s.vals.Len(readIdx))

	// 1. Collect from cold map
	s.vals.Range(readIdx, func(key, value any) bool {
		val := value.(*sumValue[N])
		
		if val.isBound {
			// Pinned/Bound entry: calculate delta
			n := val.n.load()
			last := s.reported[key]
			delta := n - last
			
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      delta,
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
			
			s.reported[key] = n // Update reported value
			// Do NOT delete from map!
		} else {
			// Non-bound entry: collect and delete
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      val.n.load(),
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
			
			s.vals.Delete(readIdx, key)
		}
		return true
	})

	// 2. Collect ONLY bound entries from hot map
	s.vals.Range(hotIdx, func(key, value any) bool {
		val := value.(*sumValue[N])
		
		if val.isBound {
			n := val.n.load()
			last := s.reported[key]
			delta := n - last
			
			newPt := metricdata.DataPoint[N]{
				Attributes: val.attrs,
				StartTime:  s.start,
				Time:       t,
				Value:      delta,
			}
			collectExemplars(&newPt.Exemplars, val.res.Collect)
			dPts = append(dPts, newPt)
			
			s.reported[key] = n // Update reported value
			// Do NOT delete from map!
		}
		return true
	})

	// The delta collection cycle resets.
	s.start = t

	sData.DataPoints = dPts
	*dest = sData

	return len(dPts)
}

// LookupBoundMeasure returns a Float64Counter that can be used to record measurements
// for the given attributes without map lookups.
func (s *deltaSum[N]) LookupBoundMeasure(attrs []attribute.KeyValue) metric.Float64Counter {
	sFloat, ok := any(s).(*deltaSum[float64])
	if !ok {
		return nil
	}

	_, compacted := attribute.NewDistinctWithFilter(attrs, nil)
	compactedSet := attribute.NewSet(compacted...)
	fltrAttr, _ := compactedSet.Filter(sFloat.filter)

	sv := sFloat.vals.LoadOrStoreBound(fltrAttr, func(attr attribute.Set) any {
		r := sFloat.newRes(attr)
		_, isDrop := r.(*dropRes[float64])
		newSV := &sumValue[float64]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
			isBound:       true,
		}
		newSV.boundFloat64 = &boundFloat64SumValue{sv: newSV}
		return newSV
	}).(*sumValue[float64])

	sv.isBound = true // Mark as bound in case it was loaded as unbound!
	return sv.boundFloat64
}

// LookupBoundMeasureInt64 returns an Int64Counter that can be used to record measurements
// for the given attributes without map lookups.
func (s *deltaSum[N]) LookupBoundMeasureInt64(attrs []attribute.KeyValue) metric.Int64Counter {
	sInt, ok := any(s).(*deltaSum[int64])
	if !ok {
		return nil
	}

	_, compacted := attribute.NewDistinctWithFilter(attrs, nil)
	compactedSet := attribute.NewSet(compacted...)
	fltrAttr, _ := compactedSet.Filter(sInt.filter)

	sv := sInt.vals.LoadOrStoreBound(fltrAttr, func(attr attribute.Set) any {
		r := sInt.newRes(attr)
		_, isDrop := r.(*dropRes[int64])
		newSV := &sumValue[int64]{
			res:           r,
			attrs:         attr,
			startTime:     now(),
			dropExemplars: isDrop,
			isBound:       true,
		}
		newSV.boundInt64 = &boundInt64SumValue{sv: newSV}
		return newSV
	}).(*sumValue[int64])

	sv.isBound = true // Mark as bound in case it was loaded as unbound!
	return sv.boundInt64
}

// newCumulativeSum returns an aggregator that summarizes a set of measurements
// as their arithmetic sum. Each sum is scoped by attributes and the
// aggregation cycle the measurements were made in.
func newCumulativeSum[N int64 | float64](
	monotonic bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
	filter attribute.Filter,
) *cumulativeSum[N] {
	state := &cardinalityState{limit: limit}
	return &cumulativeSum[N]{
		monotonic: monotonic,
		start:     now(),
		sumValueMap: sumValueMap[N]{
			values: limitedSyncMap{state: state},
			newRes: r,
		},
		cardinality: state,
		filter:      filter,
	}
}

// cumulativeSum is the storage for sums which never reset.
type cumulativeSum[N int64 | float64] struct {
	monotonic bool
	start     time.Time

	sumValueMap[N]
	cardinality *cardinalityState
	filter      attribute.Filter
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
	dPts := reset(sData.DataPoints, 0, s.values.Len())

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	s.values.Range(func(_, value any) bool {
		val := value.(*sumValue[N])

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

// LookupBoundMeasure returns a Float64Counter that can be used to record measurements
// for the given attributes without map lookups.
func (s *cumulativeSum[N]) LookupBoundMeasure(attrs []attribute.KeyValue) metric.Float64Counter {
	sFloat, ok := any(s).(*cumulativeSum[float64])
	if !ok {
		return nil
	}

	_, compacted := attribute.NewDistinctWithFilter(attrs, nil)
	compactedSet := attribute.NewSet(compacted...)
	fltrAttr, _ := compactedSet.Filter(sFloat.filter)
	d := fltrAttr.Equivalent()

	var sv *sumValue[float64]
	if actual, loaded := sFloat.values.LoadByDistinct(d); loaded {
		sv = actual.(*sumValue[float64])
		sv.isBound = true
	} else {
		sv = sFloat.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
			r := sFloat.newRes(attr)
			_, isDrop := r.(*dropRes[float64])
			newSV := &sumValue[float64]{
				res:           r,
				attrs:         attr,
				startTime:     now(),
				dropExemplars: isDrop,
				isBound:       true,
			}
			newSV.boundFloat64 = &boundFloat64SumValue{sv: newSV}
			return newSV
		}).(*sumValue[float64])
	}

	return sv.boundFloat64
}

// LookupBoundMeasureInt64 returns an Int64Counter that can be used to record measurements
// for the given attributes without map lookups.
func (s *cumulativeSum[N]) LookupBoundMeasureInt64(attrs []attribute.KeyValue) metric.Int64Counter {
	sInt, ok := any(s).(*cumulativeSum[int64])
	if !ok {
		return nil
	}

	_, compacted := attribute.NewDistinctWithFilter(attrs, nil)
	compactedSet := attribute.NewSet(compacted...)
	fltrAttr, _ := compactedSet.Filter(sInt.filter)
	d := fltrAttr.Equivalent()

	var sv *sumValue[int64]
	if actual, loaded := sInt.values.LoadByDistinct(d); loaded {
		sv = actual.(*sumValue[int64])
		sv.isBound = true
	} else {
		sv = sInt.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
			r := sInt.newRes(attr)
			_, isDrop := r.(*dropRes[int64])
			newSV := &sumValue[int64]{
				res:           r,
				attrs:         attr,
				startTime:     now(),
				dropExemplars: isDrop,
				isBound:       true,
			}
			newSV.boundInt64 = &boundInt64SumValue{sv: newSV}
			return newSV
		}).(*sumValue[int64])
	}

	return sv.boundInt64
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
		deltaSum: newDeltaSum[N](monotonic, limit, r, nil),
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
	readIdx := s.vals.SwapHotAndWait()
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
	readIdx := s.vals.SwapHotAndWait()
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
	s.vals.Clear(readIdx)

	sData.DataPoints = dPts
	*dest = sData

	return i
}
