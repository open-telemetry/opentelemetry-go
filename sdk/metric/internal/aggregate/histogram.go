// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type histogramPoint[N int64 | float64] struct {
	sync.Mutex
	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	counts   []uint64
	count    uint64
	total    N
	min, max N
}

func (b *histogramPoint[N]) sum(value N) { b.total += value }

func (b *histogramPoint[N]) bin(idx int) {
	b.counts[idx]++
	b.count++
}

func (b *histogramPoint[N]) minMax(value N) {
	if value < b.min {
		b.min = value
	} else if value > b.max {
		b.max = value
	}
}

// histogramValueMap summarizes a set of measurements as an histogramValueMap with
// explicitly defined buckets.
type histogramValueMap[N int64 | float64] struct {
	noMinMax bool
	noSum    bool
	bounds   []float64

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	values limitedSyncMap
}

func newHistogramValueMap[N int64 | float64](
	bounds []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) histogramValueMap[N] {
	// The responsibility of keeping all histogramPoint correctly associated with the
	// passed boundaries is ultimately this type's responsibility. Make a copy
	// here so we can always guarantee this. Or, in the case of failure, have
	// complete control over the fix.
	b := slices.Clone(bounds)
	slices.Sort(b)
	return histogramValueMap[N]{
		noMinMax: noMinMax,
		noSum:    noSum,
		bounds:   b,
		newRes:   r,
		values:   limitedSyncMap{aggLimit: limit},
	}
}

func (s *histogramValueMap[N]) measure(ctx context.Context, value N, fltrAttr attribute.Set, droppedAttr []attribute.KeyValue) {
	// This search will return an index in the range [0, len(s.bounds)], where
	// it will return len(s.bounds) if value is greater than the last element
	// of s.bounds. This aligns with the histogramPoint in that the length of histogramPoint
	// is len(s.bounds)+1, with the last bucket representing:
	// (s.bounds[len(s.bounds)-1], +∞).
	idx := sort.SearchFloat64s(s.bounds, float64(value))
	h := s.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		return &histogramPoint[N]{
			res:    s.newRes(attr),
			attrs:  attr,
			min:    value,
			max:    value,
			counts: make([]uint64, len(s.bounds)+1),
		}
	}).(*histogramPoint[N])
	h.Lock()
	defer h.Unlock()

	h.bin(idx)
	if !s.noMinMax {
		h.minMax(value)
	}
	if !s.noSum {
		h.sum(value)
	}
	h.res.Offer(ctx, value, droppedAttr)
}

// deltaHistogram TODO
type deltaHistogram[N int64 | float64] struct {
	hcwg          hotColdWaitGroup
	hotColdValMap [2]histogramValueMap[N]

	start time.Time
}

// measure TODO
func (s *deltaHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	hotIdx := s.hcwg.start()
	defer s.hcwg.done(hotIdx)
	s.hotColdValMap[hotIdx].measure(ctx, value, fltrAttr, droppedAttr)
}

// newDeltaHistogram TODO
func newDeltaHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *deltaHistogram[N] {
	return &deltaHistogram[N]{
		start: now(),
		hotColdValMap: [2]histogramValueMap[N]{
			newHistogramValueMap(boundaries, noMinMax, noSum, limit, r),
			newHistogramValueMap(boundaries, noMinMax, noSum, limit, r),
		},
	}
}

func (s *deltaHistogram[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Histogram, memory reuse is missed. In that
	// case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.Histogram[N])
	h.Temporality = metricdata.DeltaTemporality

	// delta always clears values on collection
	readIdx := s.hcwg.swapHotAndWait()

	// Do not allow modification of our copy of bounds.
	bounds := slices.Clone(s.hotColdValMap[readIdx].bounds)

	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := s.hotColdValMap[readIdx].values.Len()
	hDPts := reset(h.DataPoints, n, n)

	var i int
	s.hotColdValMap[readIdx].values.Range(func(key, value any) bool {
		val := value.(*histogramPoint[N])
		val.Lock()
		defer val.Unlock()
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = s.start
		hDPts[i].Time = t
		hDPts[i].Count = val.count
		hDPts[i].Bounds = bounds
		hDPts[i].BucketCounts = val.counts

		if !s.hotColdValMap[readIdx].noSum {
			hDPts[i].Sum = val.total
		}

		if !s.hotColdValMap[readIdx].noMinMax {
			hDPts[i].Min = metricdata.NewExtrema(val.min)
			hDPts[i].Max = metricdata.NewExtrema(val.max)
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
		return true
	})
	// Unused attribute sets do not report.
	s.hotColdValMap[readIdx].values.Clear()
	// The delta collection cycle resets.
	s.start = t

	h.DataPoints = hDPts
	*dest = h

	return n
}

// cumulativeHistogram summarizes a set of measurements as an histogram with explicitly
// defined histogramPoint.
type cumulativeHistogram[N int64 | float64] struct {
	histogramValueMap[N]

	start time.Time
}

// newDeltaHistogram TODO
func newCumulativeHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeHistogram[N] {
	return &cumulativeHistogram[N]{
		start:             now(),
		histogramValueMap: newHistogramValueMap(boundaries, noMinMax, noSum, limit, r),
	}
}

func (s *cumulativeHistogram[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Histogram, memory reuse is missed. In that
	// case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.Histogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	// Do not allow modification of our copy of bounds.
	bounds := slices.Clone(s.bounds)

	// Values are being concurrently written while we iterate, so only use the
	// current length for capacity.
	hDPts := reset(h.DataPoints, 0, s.values.Len())

	var i int
	s.values.Range(func(key, value any) bool {
		val := value.(*histogramPoint[N])
		val.Lock()
		defer val.Unlock()
		newPt := metricdata.HistogramDataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Count:      val.count,
			Bounds:     bounds,
			// The HistogramDataPoint field values returned need to be copies of
			// the histogramPoint value as we will keep updating them.
			//
			// TODO (#3047): Making copies for bounds and counts incurs a large
			// memory allocation footprint. Alternatives should be explored.
			BucketCounts: slices.Clone(val.counts),
		}

		if !s.noSum {
			newPt.Sum = val.total
		}

		if !s.noMinMax {
			newPt.Min = metricdata.NewExtrema(val.min)
			newPt.Max = metricdata.NewExtrema(val.max)
		}

		collectExemplars(&newPt.Exemplars, val.res.Collect)
		hDPts = append(hDPts, newPt)

		i++
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
		return true
	})

	h.DataPoints = hDPts
	*dest = h

	return i
}
