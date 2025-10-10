// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"slices"
	"sort"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type histogramPoint[N int64 | float64] struct {
	attrs attribute.Set
	res   FilteredExemplarReservoir[N]
	histogramPointCounters[N]
}

// histogramPointCounters contains only the atomic counter data, and is used by
// both histogramPoint and hotColdHistogramPoint.
type histogramPointCounters[N int64 | float64] struct {
	counts []atomic.Uint64
	total  atomicCounter[N]
	minMax atomicMinMax[N]
}

func (b *histogramPointCounters[N]) sum(value N) { b.total.add(value) }

func (b *histogramPointCounters[N]) bin(bounds []float64, value N) {
	// This search will return an index in the range [0, len(s.bounds)], where
	// it will return len(s.bounds) if value is greater than the last element
	// of s.bounds. This aligns with the histogramPoint in that the length of histogramPoint
	// is len(s.bounds)+1, with the last bucket representing:
	// (s.bounds[len(s.bounds)-1], +∞).
	idx := sort.SearchFloat64s(bounds, float64(value))
	b.counts[idx].Add(1)
}

func (b *histogramPointCounters[N]) loadCounts() ([]uint64, uint64) {
	// TODO (#3047): Making copies for bounds and counts incurs a large
	// memory allocation footprint. Alternatives should be explored.
	counts := make([]uint64, len(b.counts))
	count := uint64(0)
	for i := range counts {
		c := b.counts[i].Load()
		counts[i] = c
		count += c
	}
	return counts, count
}

// mergeIntoAndReset merges this set of histogram counter data into another,
// and resets the state of this set of counters. This is used by
// hotColdHistogramPoint to ensure that the cumulative counters continue to
// accumulate after being read.
func (b *histogramPointCounters[N]) mergeIntoAndReset( // nolint:revive // Intentional internal control flag
	into *histogramPointCounters[N],
	noMinMax, noSum bool,
) {
	for i := range b.counts {
		into.counts[i].Add(b.counts[i].Load())
		b.counts[i].Store(0)
	}

	if !noMinMax {
		// Do not reset min or max because cumulative min and max only ever grow
		// smaller or larger respectively.

		if b.minMax.set.Load() {
			into.minMax.Update(b.minMax.minimum.Load())
			into.minMax.Update(b.minMax.maximum.Load())
		}
	}
	if !noSum {
		into.total.add(b.total.load())
		b.total.reset()
	}
}

// deltaHistogram is a histogram whose internal storage is reset when it is
// collected.
type deltaHistogram[N int64 | float64] struct {
	hcwg          hotColdWaitGroup
	hotColdValMap [2]limitedSyncMap

	start    time.Time
	noMinMax bool
	noSum    bool
	bounds   []float64
	newRes   func(attribute.Set) FilteredExemplarReservoir[N]
}

func (s *deltaHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	hotIdx := s.hcwg.start()
	defer s.hcwg.done(hotIdx)
	h := s.hotColdValMap[hotIdx].LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		hPt := &histogramPoint[N]{
			res:   s.newRes(attr),
			attrs: attr,
			// N+1 buckets. For example:
			//
			//   bounds = [0, 5, 10]
			//
			// Then,
			//
			//   buckets = (-∞, 0], (0, 5.0], (5.0, 10.0], (10.0, +∞)
			histogramPointCounters: histogramPointCounters[N]{counts: make([]atomic.Uint64, len(s.bounds)+1)},
		}
		return hPt
	}).(*histogramPoint[N])

	h.bin(s.bounds, value)
	if !s.noMinMax {
		h.minMax.Update(value)
	}
	if !s.noSum {
		h.sum(value)
	}
	h.res.Offer(ctx, value, droppedAttr)
}

// newDeltaHistogram returns a histogram that is reset each time it is
// collected.
func newDeltaHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *deltaHistogram[N] {
	// The responsibility of keeping all histogramPoint correctly associated with the
	// passed boundaries is ultimately this type's responsibility. Make a copy
	// here so we can always guarantee this. Or, in the case of failure, have
	// complete control over the fix.
	b := slices.Clone(boundaries)
	slices.Sort(b)
	return &deltaHistogram[N]{
		start:    now(),
		noMinMax: noMinMax,
		noSum:    noSum,
		bounds:   b,
		newRes:   r,
		hotColdValMap: [2]limitedSyncMap{
			{aggLimit: limit},
			{aggLimit: limit},
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
	bounds := slices.Clone(s.bounds)

	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := s.hotColdValMap[readIdx].Len()
	hDPts := reset(h.DataPoints, n, n)

	var i int
	s.hotColdValMap[readIdx].Range(func(_, value any) bool {
		val := value.(*histogramPoint[N])
		bucketCounts, count := val.loadCounts()
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = s.start
		hDPts[i].Time = t
		hDPts[i].Count = count
		hDPts[i].Bounds = bounds
		hDPts[i].BucketCounts = bucketCounts

		if !s.noSum {
			hDPts[i].Sum = val.total.load()
		}

		if !s.noMinMax {
			if val.minMax.set.Load() {
				hDPts[i].Min = metricdata.NewExtrema(val.minMax.minimum.Load())
				hDPts[i].Max = metricdata.NewExtrema(val.minMax.maximum.Load())
			}
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
		return true
	})
	// Unused attribute sets do not report.
	s.hotColdValMap[readIdx].Clear()
	// The delta collection cycle resets.
	s.start = t

	h.DataPoints = hDPts
	*dest = h

	return n
}

// cumulativeHistogram summarizes a set of measurements as an histogram with explicitly
// defined histogramPoint.
type cumulativeHistogram[N int64 | float64] struct {
	values limitedSyncMap

	start    time.Time
	noMinMax bool
	noSum    bool
	bounds   []float64
	newRes   func(attribute.Set) FilteredExemplarReservoir[N]
}

// newCumulativeHistogram returns a histogram that accumulates measurements
// into a histogram data structure. It is never reset.
func newCumulativeHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeHistogram[N] {
	// The responsibility of keeping all histogramPoint correctly associated with the
	// passed boundaries is ultimately this type's responsibility. Make a copy
	// here so we can always guarantee this. Or, in the case of failure, have
	// complete control over the fix.
	b := slices.Clone(boundaries)
	slices.Sort(b)
	return &cumulativeHistogram[N]{
		start:    now(),
		noMinMax: noMinMax,
		noSum:    noSum,
		bounds:   b,
		newRes:   r,
		values:   limitedSyncMap{aggLimit: limit},
	}
}

type hotColdHistogramPoint[N int64 | float64] struct {
	hcwg         hotColdWaitGroup
	hotColdPoint [2]histogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]
}

func (s *cumulativeHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	h := s.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		hPt := &hotColdHistogramPoint[N]{
			res:   s.newRes(attr),
			attrs: attr,
			// N+1 buckets. For example:
			//
			//   bounds = [0, 5, 10]
			//
			// Then,
			//
			//   buckets = (-∞, 0], (0, 5.0], (5.0, 10.0], (10.0, +∞)
			hotColdPoint: [2]histogramPointCounters[N]{
				{
					counts: make([]atomic.Uint64, len(s.bounds)+1),
				},
				{
					counts: make([]atomic.Uint64, len(s.bounds)+1),
				},
			},
		}
		return hPt
	}).(*hotColdHistogramPoint[N])

	hotIdx := h.hcwg.start()
	defer h.hcwg.done(hotIdx)

	h.hotColdPoint[hotIdx].bin(s.bounds, value)
	if !s.noMinMax {
		h.hotColdPoint[hotIdx].minMax.Update(value)
	}
	if !s.noSum {
		h.hotColdPoint[hotIdx].sum(value)
	}
	h.res.Offer(ctx, value, droppedAttr)
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
	s.values.Range(func(_, value any) bool {
		val := value.(*hotColdHistogramPoint[N])
		// swap, observe, and clear the point
		readIdx := val.hcwg.swapHotAndWait()
		bucketCounts, count := val.hotColdPoint[readIdx].loadCounts()
		newPt := metricdata.HistogramDataPoint[N]{
			Attributes: val.attrs,
			StartTime:  s.start,
			Time:       t,
			Count:      count,
			Bounds:     bounds,
			// The HistogramDataPoint field values returned need to be copies of
			// the histogramPoint value as we will keep updating them.
			BucketCounts: bucketCounts,
		}

		if !s.noSum {
			newPt.Sum = val.hotColdPoint[readIdx].total.load()
		}
		if !s.noMinMax {
			if val.hotColdPoint[readIdx].minMax.set.Load() {
				newPt.Min = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.minimum.Load())
				newPt.Max = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.maximum.Load())
			}
		}
		// Once we've read the point, merge it back into the hot histogram
		// point since it is cumulative.
		hotIdx := (readIdx + 1) % 2
		val.hotColdPoint[readIdx].mergeIntoAndReset(&val.hotColdPoint[hotIdx], s.noMinMax, s.noSum)

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
