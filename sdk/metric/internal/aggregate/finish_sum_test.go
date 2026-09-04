// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestFinishSumLifecycle(t *testing.T) {
	var c clock
	unregister := c.Register()
	defer unregister()

	t.Run("Cumulative", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{
			Temporality:   metricdata.CumulativeTemporality,
			ReservoirFunc: dropExemplars[int64],
		}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		agg.Finish(alice.Equivalent(), y2kPlus(20))

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(2), points[0].Value)
		assert.Equal(t, y2kPlus(0), points[0].StartTime)
		assert.Equal(t, y2kPlus(10), points[0].Time)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		agg.Measure(t.Context(), 3, alice)
		points = finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(3), points[0].Value)
		assert.Equal(t, y2kPlus(3), points[0].StartTime)
	})

	t.Run("CancelPending", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		agg.Measure(t.Context(), 3, alice)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(5), points[0].Value)
		assert.Equal(t, y2kPlus(0), points[0].StartTime)
		assert.Equal(t, y2kPlus(1), points[0].Time)
		require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 1)
	})

	t.Run("Delta", func(t *testing.T) {
		c.Reset()
		agg := Builder[int64]{
			Temporality:   metricdata.DeltaTemporality,
			ReservoirFunc: dropExemplars[int64],
		}.FinishSum(true)
		agg.Measure(t.Context(), 2, alice)
		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(2), points[0].Value)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		// A series retired by a prior delta collection cannot synthesize a
		// zero-valued final point.
		agg.Finish(alice.Equivalent(), y2kPlus(10))
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))

		agg.Measure(t.Context(), 3, alice)
		agg.Finish(alice.Equivalent(), y2kPlus(20))
		points = finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(3), points[0].Value)
		assert.Equal(t, y2kPlus(20), points[0].Time)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})
}

func TestFinishSumAttributes(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		empty := attribute.NewSet()
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 1, empty)
		agg.Finish(empty.Equivalent(), y2k)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, empty, points[0].Attributes)
	})

	t.Run("Unknown", func(t *testing.T) {
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Finish(alice.Equivalent(), y2k)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})

	t.Run("FilteredExemplar", func(t *testing.T) {
		mockRes := &notConcurrentSafeReservoir{}
		agg := Builder[int64]{
			Filter: attrFltr,
			ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
				return NewFilteredExemplarReservoir[int64](exemplar.AlwaysOnFilter, mockRes)
			},
		}.FinishSum(true)
		agg.Measure(t.Context(), 7, alice)
		agg.Finish(fltrAlice.Equivalent(), y2k)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, fltrAlice, points[0].Attributes)
		require.Len(t, points[0].Exemplars, 1)
		assert.Equal(t, int64(7), points[0].Exemplars[0].Value)
		assert.Equal(t, []attribute.KeyValue{adminTrue}, points[0].Exemplars[0].FilteredAttributes)
	})
}

func TestFinishSumCardinality(t *testing.T) {
	agg := Builder[int64]{
		AggregationLimit: 3,
		ReservoirFunc:    dropExemplars[int64],
	}.FinishSum(true)
	carol := attribute.NewSet(userCarol)
	dave := attribute.NewSet(userDave)

	agg.Measure(t.Context(), 1, alice)
	agg.Measure(t.Context(), 2, bob)
	agg.Measure(t.Context(), 3, carol) // Overflow.
	agg.Finish(alice.Equivalent(), y2k)
	points := finishSumPoints(t, agg.ComputeAggregation)
	require.Len(t, points, 3)

	// The finalized Alice series released a normal slot even though the
	// shared overflow series remains active.
	agg.Measure(t.Context(), 4, dave)
	points = finishSumPoints(t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	assert.Equal(t, int64(4), pointWithAttrs(t, points, dave).Value)
	assert.Equal(t, int64(3), pointWithAttrs(t, points, overflowSet).Value)

	// The shared overflow series cannot be finalized, even when addressed by
	// its reserved attributes.
	agg.Finish(overflowSet.Equivalent(), y2k)
	require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 3)
}

func TestFinishSumOverflowProvenance(t *testing.T) {
	t.Run("UnlimitedMarkerSeries", func(t *testing.T) {
		agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Finish(overflowSet.Equivalent(), y2k)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, int64(1), points[0].Value)
		assert.Empty(t, finishSumPoints(t, agg.ComputeAggregation))
	})

	t.Run("MarkerSeriesBecomesOverflow", func(t *testing.T) {
		agg := Builder[int64]{
			AggregationLimit: 3,
			ReservoirFunc:    dropExemplars[int64],
		}.FinishSum(true)
		carol := attribute.NewSet(userCarol)

		agg.Measure(t.Context(), 1, overflowSet)
		agg.Measure(t.Context(), 2, alice)
		agg.Measure(t.Context(), 3, bob)
		agg.Measure(t.Context(), 4, carol) // Routed to the marker series.
		agg.Finish(overflowSet.Equivalent(), y2k)

		points := finishSumPoints(t, agg.ComputeAggregation)
		require.Len(t, points, 3)
		assert.Equal(t, int64(5), pointWithAttrs(t, points, overflowSet).Value)
		require.Len(t, finishSumPoints(t, agg.ComputeAggregation), 3)
	})
}

func TestFinishSumInitialMeasurementAdmission(t *testing.T) {
	point := newFinishSumValue(alice, false, dropExemplars[int64])
	measurement, ok := point.lifecycle.AcquireMeasurement()
	require.True(t, ok)

	type collectionResult struct {
		point metricdata.DataPoint[int64]
		emit  bool
	}
	started := make(chan struct{})
	collected := make(chan collectionResult, 1)
	go func() {
		close(started)
		point, emit, _ := point.collectDelta(y2k)
		collected <- collectionResult{point: point, emit: emit}
	}()
	<-started
	select {
	case <-collected:
		t.Fatal("collection completed before the initial measurement")
	case <-time.After(10 * time.Millisecond):
	}

	point.value.add(1)
	measurement.Release()
	result := <-collected
	assert.True(t, result.emit)
	assert.Equal(t, int64(1), result.point.Value)
}

func TestFinishSumShutdown(t *testing.T) {
	agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
	agg.Measure(t.Context(), 1, alice)
	agg.Shutdown()
	agg.Shutdown()
	agg.Measure(t.Context(), 2, alice)
	agg.Finish(alice.Equivalent(), time.Now())
	var data metricdata.Aggregation
	assert.Zero(t, agg.ComputeAggregation(&data))
	assert.Nil(t, data)
}

func TestFinishSumRetireAndDelete(t *testing.T) {
	store := newFinishSum(
		true,
		metricdata.CumulativeTemporality,
		0,
		dropExemplars[int64],
	)
	lazy := newLazyFilteredAttributes(alice, nil)
	store.measure(t.Context(), 1, lazy)
	raw, ok := store.values.Load(alice.Equivalent())
	require.True(t, ok)
	point := raw.(*finishSumValue[int64])

	store.retireAndDelete(point)
	store.retireAndDelete(point)
	assert.Zero(t, store.values.Len())
	assert.False(t, point.measure(t.Context(), 1, lazy))

	_, emit, retire := point.collectCumulative(y2k)
	assert.False(t, emit)
	assert.False(t, retire)
}

func TestFinishSumConcurrentLifecycle(t *testing.T) {
	agg := Builder[int64]{
		Temporality:   metricdata.DeltaTemporality,
		ReservoirFunc: dropExemplars[int64],
	}.FinishSum(true)

	const measurements = 2_000
	var (
		collected atomic.Int64
		stop      atomic.Bool
		wg        sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for !stop.Load() {
			agg.Finish(alice.Equivalent(), time.Now())
		}
	}()
	go func() {
		defer wg.Done()
		for !stop.Load() {
			var data metricdata.Aggregation
			agg.ComputeAggregation(&data)
			for _, point := range data.(metricdata.Sum[int64]).DataPoints {
				collected.Add(point.Value)
			}
		}
	}()
	for range measurements {
		agg.Measure(t.Context(), 1, alice)
	}
	stop.Store(true)
	wg.Wait()
	for _, point := range finishSumPoints(t, agg.ComputeAggregation) {
		collected.Add(point.Value)
	}
	assert.Equal(t, int64(measurements), collected.Load())
}

func TestFinishSumConcurrentShutdown(t *testing.T) {
	agg := Builder[int64]{ReservoirFunc: dropExemplars[int64]}.FinishSum(true)
	var (
		stop atomic.Bool
		wg   sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for !stop.Load() {
			agg.Measure(t.Context(), 1, alice)
		}
	}()
	go func() {
		defer wg.Done()
		for !stop.Load() {
			agg.Finish(alice.Equivalent(), time.Now())
		}
	}()
	agg.Shutdown()
	stop.Store(true)
	wg.Wait()
	var data metricdata.Aggregation
	assert.Zero(t, agg.ComputeAggregation(&data))
}

func BenchmarkFinishSum(b *testing.B) {
	b.Run("Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		agg := Builder[int64]{Temporality: metricdata.CumulativeTemporality}.FinishSum(true)
		return agg.Measure, agg.ComputeAggregation
	}))
	b.Run("Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		agg := Builder[int64]{Temporality: metricdata.DeltaTemporality}.FinishSum(true)
		return agg.Measure, agg.ComputeAggregation
	}))
}

func BenchmarkFinishSumMeasure(b *testing.B) {
	benchmarks := []struct {
		name    string
		factory func() Measure[int64]
	}{
		{"false", func() Measure[int64] {
			measure, _ := Builder[int64]{Temporality: metricdata.CumulativeTemporality}.Sum(true)
			return measure
		}},
		{"true", func() Measure[int64] {
			return Builder[int64]{Temporality: metricdata.CumulativeTemporality}.FinishSum(true).Measure
		}},
	}
	for _, benchmark := range benchmarks {
		b.Run("finish="+benchmark.name, func(b *testing.B) {
			b.Run("mode=serial", func(b *testing.B) {
				measure := benchmark.factory()
				b.ReportAllocs()
				for b.Loop() {
					measure(b.Context(), 1, alice)
				}
			})
			b.Run("mode=parallel", func(b *testing.B) {
				measure := benchmark.factory()
				b.ReportAllocs()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						measure(b.Context(), 1, alice)
					}
				})
			})
		})
	}
}

func finishSumPoints(
	t testing.TB,
	compute ComputeAggregation,
) []metricdata.DataPoint[int64] {
	t.Helper()
	var data metricdata.Aggregation
	compute(&data)
	sum, ok := data.(metricdata.Sum[int64])
	require.True(t, ok, "aggregation type %T", data)
	return sum.DataPoints
}

func pointWithAttrs(
	t testing.TB,
	points []metricdata.DataPoint[int64],
	attrs attribute.Set,
) metricdata.DataPoint[int64] {
	t.Helper()
	for _, point := range points {
		if point.Attributes.Equals(&attrs) {
			return point
		}
	}
	t.Fatalf("no point with attributes %v in %#v", attrs, points)
	return metricdata.DataPoint[int64]{}
}
