// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"slices"
	"sort"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/finish"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type finishHistogramPoint[N int64 | float64] struct {
	lifecycle finish.Lifecycle
	histogramPointCounters[N]
	attrs         attribute.Set
	start         time.Time
	dropExemplars bool
	reservoir     FilteredExemplarReservoir[N]
}

func newFinishHistogramPoint[N int64 | float64](
	attrs attribute.Set,
	buckets int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *finishHistogramPoint[N] {
	r := reservoir(attrs)
	_, drop := r.(*dropRes[N])
	return &finishHistogramPoint[N]{
		histogramPointCounters: histogramPointCounters[N]{counts: make([]atomic.Uint64, buckets)},
		attrs:                  attrs,
		start:                  now(),
		dropExemplars:          drop,
		reservoir:              r,
	}
}

func (p *finishHistogramPoint[N]) shutdown(ctx context.Context) error {
	return p.lifecycle.Retire(ctx)
}

// FinishExplicitBucketHistogram contains the operations of a Finish-aware
// explicit bucket Histogram aggregation. Finish accepts a post-filter series
// identity and the timestamp for its final point.
//
// Stop disables subsequent Measure, Finish, and ComputeAggregation calls.
// Wait calls Stop, retires series observed during cleanup, and clears the map.
// A recording overlapping Stop may complete after Wait returns. If it publishes
// a series missed by cleanup, it retires that series when recording completes.
// Wait returns a context error if cleanup is interrupted by cancellation.
type FinishExplicitBucketHistogram[N int64 | float64] struct {
	Measure            Measure[N]
	ComputeAggregation ComputeAggregation
	Finish             func(attribute.Distinct, time.Time)
	Stop               func()
	Wait               func(context.Context) error
}

type finishHistogram[N int64 | float64] struct {
	collectSem  chan struct{}
	stopped     atomic.Bool
	points      []*finishHistogramPoint[N]
	values      limitedSyncMap[*finishHistogramPoint[N]]
	start       time.Time
	temporality metricdata.Temporality
	bounds      []float64
	noMinMax    bool
	noSum       bool
	reservoir   func(attribute.Set) FilteredExemplarReservoir[N]
}

func newFinishHistogram[N int64 | float64](
	boundaries []float64,
	noMinMax, noSum bool,
	temporality metricdata.Temporality,
	limit int,
	reservoir func(attribute.Set) FilteredExemplarReservoir[N],
) *finishHistogram[N] {
	var start time.Time
	if temporality != metricdata.DeltaTemporality {
		temporality = metricdata.CumulativeTemporality
	} else {
		start = now()
	}
	bounds := slices.Clone(boundaries)
	slices.Sort(bounds)
	return &finishHistogram[N]{
		collectSem:  make(chan struct{}, 1),
		values:      limitedSyncMap[*finishHistogramPoint[N]]{aggLimit: limit},
		start:       start,
		temporality: temporality,
		bounds:      bounds,
		noMinMax:    noMinMax,
		noSum:       noSum,
		reservoir:   reservoir,
	}
}

func (s *finishHistogram[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	for {
		if s.stopped.Load() {
			return
		}
		var measurement finish.Measurement
		point, loaded, overflowed := s.values.LoadOrStoreAttrReclaiming(
			lazy,
			func(attrs attribute.Set, overflow bool) *finishHistogramPoint[N] {
				point := newFinishHistogramPoint(attrs, len(s.bounds)+1, s.reservoir)
				// Admit the initial measurement before publication so collection cannot
				// observe an empty point. The unpublished lifecycle cannot be retired.
				if overflow {
					measurement, _ = point.lifecycle.AcquireSharedMeasurement()
				} else {
					measurement, _ = point.lifecycle.AcquireMeasurement()
				}
				return point
			},
		)
		if loaded {
			var ok bool
			if overflowed {
				measurement, ok = point.lifecycle.AcquireSharedMeasurement()
			} else {
				measurement, ok = point.lifecycle.AcquireMeasurement()
			}
			if !ok {
				s.values.CompareAndDelete(lazy.Distinct(), point)
				continue
			}
		}
		s.measureAcquired(ctx, value, lazy, point, measurement)
		if s.stopped.Load() {
			s.retireAndDelete(point)
		}
		return
	}
}

func (s *finishHistogram[N]) measureAcquired(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
	point *finishHistogramPoint[N],
	measurement finish.Measurement,
) {
	index := sort.SearchFloat64s(s.bounds, float64(value))
	point.counts[index].Add(1)
	if !s.noMinMax {
		point.minMax.Update(value)
	}
	if !s.noSum {
		point.total.add(value)
	}
	if point.dropExemplars {
		measurement.Release()
		return
	}
	defer measurement.Release()
	point.reservoir.Offer(ctx, value, lazy)
}

func (s *finishHistogram[N]) retireAndDelete(point *finishHistogramPoint[N]) {
	_ = point.shutdown(context.Background())
	s.values.CompareAndDelete(point.attrs.Equivalent(), point)
}

func (s *finishHistogram[N]) finish(distinct attribute.Distinct, at time.Time) {
	if s.stopped.Load() {
		return
	}
	raw, ok := s.values.Load(distinct)
	if ok {
		point := raw.(*finishHistogramPoint[N])
		// A delayed Finish must not retire a lifetime created after it began.
		if at.Before(point.start) {
			return
		}
		point.lifecycle.Finish(at)
	}
}

func (s *finishHistogram[N]) collect(
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
	hist, _ := (*dest).(metricdata.Histogram[N])
	hist.Temporality = s.temporality
	points := reset(hist.DataPoints, len(snapshot), len(snapshot))
	// Exported bounds must not alias the aggregation's configuration.
	var bounds []float64
	if len(snapshot) > 0 {
		bounds = slices.Clone(s.bounds)
	}
	n := 0
	for _, point := range snapshot {
		emit, retire := s.collectPoint(point, &points[n], start, t, bounds)
		if !emit {
			continue
		}
		n++
		if retire {
			s.values.CompareAndDelete(point.attrs.Equivalent(), point)
		}
	}
	hist.DataPoints = points[:n]
	*dest = hist
	return n
}

func (s *finishHistogram[N]) collectPoint(
	point *finishHistogramPoint[N],
	dest *metricdata.HistogramDataPoint[N],
	start, at time.Time,
	bounds []float64,
) (bool, bool) {
	var collection finish.Collection
	var ok bool
	if s.temporality == metricdata.DeltaTemporality {
		collection, ok = point.lifecycle.BeginDeltaCollection(at)
		// A deferred series can predate the previous collection's boundary.
		if point.start.Before(start) {
			start = point.start
		}
	} else {
		collection, ok = point.lifecycle.BeginCumulativeCollection(at)
		start = point.start
	}
	if !ok {
		return false, false
	}
	defer collection.Complete()
	// Lifecycle admission keeps the bucket counts, sum, extrema, and exemplars
	// in one measurement snapshot without an additional hot/cold buffer.
	dest.Attributes = point.attrs
	dest.StartTime = start
	dest.Time = collection.Time()
	dest.Bounds = bounds
	dest.Count = point.loadCountsInto(&dest.BucketCounts)
	dest.Sum = 0
	if !s.noSum {
		dest.Sum = point.total.load()
	}
	dest.Min = metricdata.Extrema[N]{}
	dest.Max = metricdata.Extrema[N]{}
	if !s.noMinMax && point.minMax.set.Load() {
		dest.Min = metricdata.NewExtrema(point.minMax.minimum.Load())
		dest.Max = metricdata.NewExtrema(point.minMax.maximum.Load())
	}
	collectExemplars(&dest.Exemplars, point.reservoir.Collect)
	return true, collection.ShouldRetire()
}

// Snapshot before retiring series so newly reclaimed slots cannot expand this
// collection beyond the cardinality limit. collectSem protects the scratch slice.
func (s *finishHistogram[N]) snapshot() []*finishHistogramPoint[N] {
	s.values.lenMux.Lock()
	defer s.values.lenMux.Unlock()
	s.points = reset(s.points, 0, s.values.len)
	s.values.Range(func(_, raw any) bool {
		s.points = append(s.points, raw.(*finishHistogramPoint[N]))
		return true
	})
	return s.points
}

func (s *finishHistogram[N]) stop() {
	s.stopped.Store(true)
}

func (s *finishHistogram[N]) wait(ctx context.Context) error {
	s.stop()
	if err := s.lockCollection(ctx); err != nil {
		return err
	}
	defer func() { <-s.collectSem }()
	s.points = nil

	var err error
	s.values.Range(func(_, raw any) bool {
		err = raw.(*finishHistogramPoint[N]).shutdown(ctx)
		return err == nil
	})
	if err == nil {
		err = s.values.ClearContext(ctx)
	}
	return err
}

func (s *finishHistogram[N]) lockCollection(ctx context.Context) error {
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

// FinishExplicitBucketHistogram returns an explicit bucket Histogram aggregation
// with exact-attribute lifecycle support. Boundaries are copied and sorted.
func (b Builder[N]) FinishExplicitBucketHistogram(
	boundaries []float64,
	noMinMax, noSum bool,
) FinishExplicitBucketHistogram[N] {
	store := newFinishHistogram(boundaries, noMinMax, noSum, b.Temporality, b.AggregationLimit, b.resFunc())
	return FinishExplicitBucketHistogram[N]{
		Measure:            b.filter(store.measure),
		ComputeAggregation: store.collect,
		Finish:             store.finish,
		Stop:               store.stop,
		Wait:               store.wait,
	}
}
