// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"slices"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type buckets[N int64 | float64] struct {
	count    uint64
	counts   []uint64
	minMax   atomicMinMax[N]
	total    atomicSum[N]
	noSum    bool
	noMinMax bool

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]
}

func (b *buckets[N]) measure(
	ctx context.Context,
	value N,
	idx int,
	droppedAttr []attribute.KeyValue,
) {
	atomic.AddUint64(&b.counts[idx], 1)
	atomic.AddUint64(&b.count, 1)
	if !b.noMinMax {
		b.minMax.observe(value)
	}
	if !b.noSum {
		b.total.add(value)
	}
	b.res.Offer(ctx, value, droppedAttr)
}

// histValues summarizes a set of measurements as an histValues with
// explicitly defined buckets.
type histValues[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	bounds   []float64

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	limit  limiter[buckets[N]]
	values map[attribute.Distinct]*buckets[N]
	sync.RWMutex
}

func newHistValues[N int64 | float64](
	bounds []float64,
	noSum bool,
	noMinMax bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *histValues[N] {
	// The responsibility of keeping all buckets correctly associated with the
	// passed boundaries is ultimately this type's responsibility. Make a copy
	// here so we can always guarantee this. Or, in the case of failure, have
	// complete control over the fix.
	b := slices.Clone(bounds)
	slices.Sort(b)
	return &histValues[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		bounds:   b,
		newRes:   r,
		limit:    newLimiter[buckets[N]](limit),
		values:   make(map[attribute.Distinct]*buckets[N]),
	}
}

// Aggregate records the measurement value, scoped by attr, and aggregates it
// into a histogram.
func (s *histValues[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	// This search will return an index in the range [0, len(s.bounds)], where
	// it will return len(s.bounds) if value is greater than the last element
	// of s.bounds. This aligns with the buckets in that the length of buckets
	// is len(s.bounds)+1, with the last bucket representing:
	// (s.bounds[len(s.bounds)-1], +∞).
	idx := sort.SearchFloat64s(s.bounds, float64(value))

	// Hold the RLock even after we are done reading from the values map to
	// ensure we don't race with collection.
	s.RLock()
	attr := s.limit.Attributes(fltrAttr, s.values)
	b, ok := s.values[attr.Equivalent()]
	if ok {
		b.measure(ctx, value, idx, droppedAttr)
		s.RUnlock()
		return
	}
	s.RUnlock()
	// Switch to a full lock to add a new element to the map.
	s.Lock()
	defer s.Unlock()
	// Check that the element wasn't added since we last checked.
	b, ok = s.values[attr.Equivalent()]
	if ok {
		b.measure(ctx, value, idx, droppedAttr)
		return
	}
	b = &buckets[N]{
		attrs: attr,
		// N+1 buckets. For example:
		//
		//   bounds = [0, 5, 10]
		//
		// Then,
		//
		//   buckets = (-∞, 0], (0, 5.0], (5.0, 10.0], (10.0, +∞)
		counts:   make([]uint64, len(s.bounds)+1),
		res:      s.newRes(attr),
		noSum:    s.noSum,
		noMinMax: s.noMinMax,
	}
	b.measure(ctx, value, idx, droppedAttr)
	s.values[attr.Equivalent()] = b
}

// newHistogram returns an Aggregator that summarizes a set of measurements as
// an histogram.
func newHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *histogram[N] {
	return &histogram[N]{
		histValues: newHistValues[N](boundaries, noSum, noMinMax, limit, r),
		start:      now(),
	}
}

// histogram summarizes a set of measurements as an histogram with explicitly
// defined buckets.
type histogram[N int64 | float64] struct {
	*histValues[N]

	start time.Time
}

func (s *histogram[N]) delta(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Histogram, memory reuse is missed. In that
	// case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.Histogram[N])
	h.Temporality = metricdata.DeltaTemporality

	// Acquire a full lock to ensure there are no concurrent measure() calls.
	// If we only used a RLock, we could observe "partial" measurements, such
	// as a histogram count increment without a histogram total increment.
	s.Lock()
	defer s.Unlock()

	// Do not allow modification of our copy of bounds.
	bounds := slices.Clone(s.bounds)

	n := len(s.values)
	hDPts := reset(h.DataPoints, n, n)

	var i int
	for _, val := range s.values {
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = s.start
		hDPts[i].Time = t
		hDPts[i].Count = val.count
		hDPts[i].Bounds = bounds
		hDPts[i].BucketCounts = val.counts

		if !s.noSum {
			hDPts[i].Sum = val.total.load()
		}

		if !s.noMinMax {
			hDPts[i].Min = metricdata.NewExtrema(val.minMax.loadMin())
			hDPts[i].Max = metricdata.NewExtrema(val.minMax.loadMax())
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
	}
	// Unused attribute sets do not report.
	clear(s.values)
	// The delta collection cycle resets.
	s.start = t

	h.DataPoints = hDPts
	*dest = h

	return n
}

func (s *histogram[N]) cumulative(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.Histogram, memory reuse is missed. In that
	// case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.Histogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	// Acquire a full lock to ensure there are no concurrent measure() calls.
	// If we only used a RLock, we could observe "partial" measurements, such
	// as a histogram count increment without a histogram total increment.
	s.Lock()
	defer s.Unlock()

	// Do not allow modification of our copy of bounds.
	bounds := slices.Clone(s.bounds)

	n := len(s.values)
	hDPts := reset(h.DataPoints, n, n)

	var i int
	for _, val := range s.values {
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = s.start
		hDPts[i].Time = t
		hDPts[i].Count = val.count
		hDPts[i].Bounds = bounds

		// The HistogramDataPoint field values returned need to be copies of
		// the buckets value as we will keep updating them.
		//
		// TODO (#3047): Making copies for bounds and counts incurs a large
		// memory allocation footprint. Alternatives should be explored.
		hDPts[i].BucketCounts = slices.Clone(val.counts)

		if !s.noSum {
			hDPts[i].Sum = val.total.load()
		}

		if !s.noMinMax {
			hDPts[i].Min = metricdata.NewExtrema(val.minMax.loadMin())
			hDPts[i].Max = metricdata.NewExtrema(val.minMax.loadMax())
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}

	h.DataPoints = hDPts
	*dest = h

	return n
}
