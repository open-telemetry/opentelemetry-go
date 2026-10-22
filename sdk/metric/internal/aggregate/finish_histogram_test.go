// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"fmt"
	"math"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var finishHistogramBounds = []float64{0, 5, 10}

func finishHistogramPoints[N int64 | float64](
	t *testing.T,
	compute ComputeAggregation,
) []metricdata.HistogramDataPoint[N] {
	t.Helper()
	var data metricdata.Aggregation
	count := compute(&data)
	hist, ok := data.(metricdata.Histogram[N])
	require.True(t, ok, "aggregation type %T", data)
	assert.Len(t, hist.DataPoints, count)
	return hist.DataPoints
}

func histogramPointWithAttrs[N int64 | float64](
	t *testing.T,
	points []metricdata.HistogramDataPoint[N],
	attrs attribute.Set,
) metricdata.HistogramDataPoint[N] {
	t.Helper()
	index := slices.IndexFunc(
		points,
		func(point metricdata.HistogramDataPoint[N]) bool { return point.Attributes.Equals(&attrs) },
	)
	require.NotEqual(t, -1, index)
	return points[index]
}

func TestFinishHistogramLifecycle(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testFinishHistogramLifecycle(t, int64(7), int64(-3)) })
	t.Run("Float64", func(t *testing.T) { testFinishHistogramLifecycle(t, 7.25, -3.5) })
}

func testFinishHistogramLifecycle[N int64 | float64](t *testing.T, first, last N) {
	t.Helper()
	t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
	var c clock
	t.Cleanup(c.Register())
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			c.Reset()
			agg := Builder[N]{
				Temporality: temporality,
			}.FinishExplicitBucketHistogram(
				finishHistogramBounds,
				false,
				false,
			)
			agg.Finish(bob.Equivalent(), y2k)
			assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))
			agg.Measure(t.Context(), first, alice)
			agg.Measure(t.Context(), last, alice)
			agg.Finish(alice.Equivalent(), y2kPlus(10))
			agg.Finish(alice.Equivalent(), y2kPlus(20))
			points := finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			point := points[0]
			assert.Equal(t, alice, point.Attributes)
			assert.Equal(t, y2kPlus(10), point.Time)
			assert.Equal(t, uint64(2), point.Count)
			assert.Equal(t, []uint64{1, 0, 1, 0}, point.BucketCounts)
			assert.Equal(t, finishHistogramBounds, point.Bounds)
			assert.Equal(t, first+last, point.Sum)
			assert.Equal(t, metricdata.NewExtrema(last), point.Min)
			assert.Equal(t, metricdata.NewExtrema(first), point.Max)
			start := point.StartTime
			assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))

			agg.Measure(t.Context(), first, alice)
			points = finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, uint64(1), points[0].Count)
			assert.Equal(t, first, points[0].Sum)
			assert.True(t, points[0].StartTime.After(start), "recreated series needs a new start")
			agg.Finish(alice.Equivalent(), y2kPlus(30))
			points = finishHistogramPoints[N](t, agg.ComputeAggregation)
			if temporality == metricdata.DeltaTemporality {
				assert.Empty(t, points, "Finish must not synthesize an empty delta point")
			} else {
				require.Len(t, points, 1)
				assert.Equal(t, first, points[0].Sum)
				assert.Equal(t, y2kPlus(30), points[0].Time)
			}
		})
		t.Run(temporality.String()+"/Reactivate", func(t *testing.T) {
			c.Reset()
			agg := Builder[N]{
				Temporality: temporality,
			}.FinishExplicitBucketHistogram(
				finishHistogramBounds,
				false,
				false,
			)
			agg.Measure(t.Context(), first, alice)
			agg.Finish(alice.Equivalent(), y2kPlus(10))
			agg.Measure(t.Context(), last, alice)
			points := finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, first+last, points[0].Sum)
			assert.Equal(t, uint64(2), points[0].Count)
			assert.NotEqual(t, y2kPlus(10), points[0].Time)
			points = finishHistogramPoints[N](t, agg.ComputeAggregation)
			if temporality == metricdata.CumulativeTemporality {
				require.Len(t, points, 1)
				assert.Equal(t, first+last, points[0].Sum)
			} else {
				assert.Empty(t, points)
			}
		})
		t.Run(temporality.String()+"/CollectBeforeFinish", func(t *testing.T) {
			c.Reset()
			agg := Builder[N]{
				Temporality: temporality,
			}.FinishExplicitBucketHistogram(
				finishHistogramBounds,
				false,
				false,
			)
			agg.Measure(t.Context(), first, alice)
			before := finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, before, 1)
			agg.Measure(t.Context(), last, alice)
			agg.Finish(alice.Equivalent(), y2kPlus(10))
			points := finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			if temporality == metricdata.CumulativeTemporality {
				assert.Equal(t, before[0].StartTime, points[0].StartTime)
				assert.Equal(t, uint64(2), points[0].Count)
				assert.Equal(t, first+last, points[0].Sum)
			} else {
				assert.Equal(t, before[0].Time, points[0].StartTime)
				assert.Equal(t, uint64(1), points[0].Count)
				assert.Equal(t, last, points[0].Sum)
				assert.Equal(t, metricdata.NewExtrema(last), points[0].Max)
			}
			assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))
		})
	}
}

func TestFinishHistogramFinishTimestamp(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testFinishHistogramFinishTimestamp(t, int64(7)) })
	t.Run("Float64", func(t *testing.T) { testFinishHistogramFinishTimestamp(t, 1.25) })
}

func testFinishHistogramFinishTimestamp[N int64 | float64](t *testing.T, value N) {
	t.Helper()
	var c clock
	t.Cleanup(c.Register())
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		for _, recreate := range []bool{false, true} {
			for _, offset := range []int64{-1, 0, 1} {
				t.Run(fmt.Sprintf("%s/recreate=%t/offset=%d", temporality, recreate, offset), func(t *testing.T) {
					c.Reset()
					store := newFinishHistogram(finishHistogramBounds, false, false, temporality, 0, dropExemplars[N])
					lazy := newLazyFilteredAttributes(alice, nil)
					if recreate {
						store.measure(t.Context(), value, lazy)
						store.finish(alice.Equivalent(), now())
						require.Len(t, finishHistogramPoints[N](t, store.collect), 1)
						require.Empty(t, finishHistogramPoints[N](t, store.collect))
					}
					store.measure(t.Context(), value, lazy)
					raw, ok := store.values.Load(alice.Equivalent())
					require.True(t, ok)
					start := raw.(*finishHistogramPoint[N]).start
					finishTime := start.Add(time.Duration(offset) * time.Nanosecond)
					store.finish(alice.Equivalent(), finishTime)
					points := finishHistogramPoints[N](t, store.collect)
					require.Len(t, points, 1)
					assert.Equal(t, uint64(1), points[0].Count)
					assert.Equal(t, value, points[0].Sum)
					if offset < 0 {
						assert.False(t, points[0].Time.Before(start))
						if temporality == metricdata.CumulativeTemporality {
							require.Len(t, finishHistogramPoints[N](t, store.collect), 1,
								"stale Finish must not retire the new lifetime")
							store.finish(alice.Equivalent(), now())
							require.Len(t, finishHistogramPoints[N](t, store.collect), 1)
						}
					} else {
						assert.Equal(t, finishTime, points[0].Time)
					}
					assert.Empty(t, finishHistogramPoints[N](t, store.collect))
				})
			}
		}
	}
}

func TestFinishHistogramPayload(t *testing.T) {
	t.Run("Int64", testFinishHistogramPayload[int64])
	t.Run("Float64", testFinishHistogramPayload[float64])
}

func testFinishHistogramPayload[N int64 | float64](t *testing.T) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			for _, omitMinMax := range []bool{false, true} {
				for _, omitSum := range []bool{false, true} {
					t.Run(fmt.Sprintf("noMinMax=%t/noSum=%t", omitMinMax, omitSum), func(t *testing.T) {
						input := []float64{10, 0, 5}
						agg := Builder[N]{
							Temporality: temporality,
						}.FinishExplicitBucketHistogram(
							input,
							omitMinMax,
							omitSum,
						)
						assert.Equal(t, []float64{10, 0, 5}, input, "construction modified caller-owned boundaries")
						input[0] = 1000
						for _, value := range []N{-1, 0, 1, 5, 10, 11} {
							agg.Measure(t.Context(), value, alice)
						}
						agg.Finish(alice.Equivalent(), now())
						data := metricdata.Aggregation(
							metricdata.Histogram[N]{
								DataPoints: []metricdata.HistogramDataPoint[N]{
									{
										Sum:          N(999),
										Min:          metricdata.NewExtrema(N(999)),
										Max:          metricdata.NewExtrema(N(999)),
										BucketCounts: []uint64{999, 999, 999, 999},
										Exemplars:    []metricdata.Exemplar[N]{{Value: N(999)}},
									},
								},
							},
						)
						require.Equal(t, 1, agg.ComputeAggregation(&data))
						hist := data.(metricdata.Histogram[N])
						assert.Equal(t, temporality, hist.Temporality)
						point := hist.DataPoints[0]
						assert.Equal(t, uint64(6), point.Count)
						assert.Equal(t, []uint64{2, 2, 1, 1}, point.BucketCounts)
						assert.Equal(t, finishHistogramBounds, point.Bounds)
						assert.Empty(t, point.Exemplars)
						if omitSum {
							assert.Zero(t, point.Sum)
						} else {
							assert.Equal(t, N(26), point.Sum)
						}
						if omitMinMax {
							assert.Equal(t, metricdata.Extrema[N]{}, point.Min)
							assert.Equal(t, metricdata.Extrema[N]{}, point.Max)
						} else {
							assert.Equal(t, metricdata.NewExtrema(N(-1)), point.Min)
							assert.Equal(t, metricdata.NewExtrema(N(11)), point.Max)
						}
						point.Bounds[0] = -1000
						point.BucketCounts[0] = 999
						agg.Measure(t.Context(), 1, alice)
						next := finishHistogramPoints[N](t, agg.ComputeAggregation)
						require.Len(t, next, 1)
						assert.Equal(t, finishHistogramBounds, next[0].Bounds)
						assert.Equal(t, []uint64{0, 1, 0, 0}, next[0].BucketCounts)
					})
				}
			}
		})
	}
	t.Run("EmptyAttributesAndBoundaries", func(t *testing.T) {
		agg := Builder[N]{}.FinishExplicitBucketHistogram(nil, false, false)
		attrs := attribute.NewSet()
		agg.Measure(t.Context(), 0, attrs)
		agg.Finish(attrs.Equivalent(), now())
		points := finishHistogramPoints[N](t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, attrs, points[0].Attributes)
		assert.Zero(t, points[0].Sum)
		assert.Equal(t, []uint64{1}, points[0].BucketCounts)
		assert.Empty(t, points[0].Bounds)
		assert.Equal(t, metricdata.NewExtrema(N(0)), points[0].Min)
		assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))
	})
	t.Run("FilteredExemplar", func(t *testing.T) {
		mockRes := &notConcurrentSafeReservoir{}
		agg := Builder[N]{Filter: attrFltr, ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[N] {
			return NewFilteredExemplarReservoir[N](exemplar.AlwaysOnFilter, mockRes)
		}}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
		agg.Measure(t.Context(), 7, alice)
		agg.Finish(fltrAlice.Equivalent(), now())
		points := finishHistogramPoints[N](t, agg.ComputeAggregation)
		require.Len(t, points, 1)
		assert.Equal(t, fltrAlice, points[0].Attributes)
		require.Len(t, points[0].Exemplars, 1)
		assert.Equal(t, N(7), points[0].Exemplars[0].Value)
		assert.Equal(t, []attribute.KeyValue{adminTrue}, points[0].Exemplars[0].FilteredAttributes)
		assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))
	})
}

func TestFinishHistogramFloatSpecialValues(t *testing.T) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			builder := Builder[float64]{Temporality: temporality}
			agg := builder.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
			stable, compute := builder.ExplicitBucketHistogram(finishHistogramBounds, false, false)
			for _, value := range []float64{math.NaN(), math.Inf(-1), math.Inf(1), -1, 0, 0.5, 5, 10} {
				agg.Measure(t.Context(), value, alice)
				stable(t.Context(), value, alice)
			}
			got := finishHistogramPoints[float64](t, agg.ComputeAggregation)
			want := finishHistogramPoints[float64](t, compute)
			require.Len(t, got, 1)
			require.Len(t, want, 1)
			assert.Equal(t, want[0].BucketCounts, got[0].BucketCounts)
			assert.Equal(t, want[0].Count, got[0].Count)
			assert.True(t, math.IsNaN(got[0].Sum))
			for _, pair := range [][2]metricdata.Extrema[float64]{{want[0].Min, got[0].Min}, {want[0].Max, got[0].Max}} {
				wantValue, wantDefined := pair[0].Value()
				gotValue, gotDefined := pair[1].Value()
				assert.Equal(t, wantDefined, gotDefined)
				if math.IsNaN(wantValue) {
					assert.True(t, math.IsNaN(gotValue))
				} else {
					assert.Equal(t, wantValue, gotValue)
				}
			}
		})
	}
}

func TestFinishHistogramCardinality(t *testing.T) {
	t.Run("Int64", testFinishHistogramCardinality[int64])
	t.Run("Float64", testFinishHistogramCardinality[float64])
}

func testFinishHistogramCardinality[N int64 | float64](t *testing.T) {
	agg := Builder[N]{AggregationLimit: 3}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
	carol := attribute.NewSet(userCarol)
	dave := attribute.NewSet(userDave)
	agg.Measure(t.Context(), 1, alice)
	agg.Measure(t.Context(), 2, bob)
	agg.Measure(t.Context(), 3, carol)
	agg.Measure(t.Context(), 4, carol)
	agg.Finish(carol.Equivalent(), now())
	agg.Finish(alice.Equivalent(), now())
	points := finishHistogramPoints[N](t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	overflow := histogramPointWithAttrs(t, points, overflowSet)
	assert.Equal(t, uint64(2), overflow.Count)
	assert.Equal(t, N(7), overflow.Sum)
	agg.Measure(t.Context(), 5, dave)
	agg.Finish(overflowSet.Equivalent(), now())
	points = finishHistogramPoints[N](t, agg.ComputeAggregation)
	require.Len(t, points, 3)
	assert.Equal(t, N(5), histogramPointWithAttrs(t, points, dave).Sum)
	assert.Equal(t, N(7), histogramPointWithAttrs(t, points, overflowSet).Sum)
	t.Run("UnlimitedMarkerSeries", func(t *testing.T) {
		agg := Builder[N]{}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Finish(overflowSet.Equivalent(), now())
		require.Len(t, finishHistogramPoints[N](t, agg.ComputeAggregation), 1)
		assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))
	})
	t.Run("MarkerSeriesBecomesOverflow", func(t *testing.T) {
		agg := Builder[N]{AggregationLimit: 3}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
		agg.Measure(t.Context(), 1, overflowSet)
		agg.Measure(t.Context(), 2, alice)
		agg.Measure(t.Context(), 3, bob)
		agg.Finish(overflowSet.Equivalent(), now())
		agg.Measure(t.Context(), 4, carol)
		points := finishHistogramPoints[N](t, agg.ComputeAggregation)
		require.Len(t, points, 3)
		assert.Equal(t, N(5), histogramPointWithAttrs(t, points, overflowSet).Sum)
		require.Len(t, finishHistogramPoints[N](t, agg.ComputeAggregation), 3)
	})
}

func TestFinishHistogramConcurrentSafe(t *testing.T) {
	t.Run("Int64", testFinishHistogramConcurrentSafe[int64])
	t.Run("Float64", testFinishHistogramConcurrentSafe[float64])
}

func testFinishHistogramConcurrentSafe[N int64 | float64](t *testing.T) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			agg := Builder[N]{
				Temporality:      temporality,
				Filter:           attrFltr,
				AggregationLimit: 3,
			}.FinishExplicitBucketHistogram(
				[]float64{0, 2, 4},
				false,
				false,
			)
			testAggregationConcurrentSafe(agg.Measure, agg.ComputeAggregation, validateHistogram[N])(t)
		})
	}
	t.Run("FinishAndCollection", func(t *testing.T) {
		agg := Builder[N]{
			Temporality: metricdata.DeltaTemporality,
		}.FinishExplicitBucketHistogram(
			finishHistogramBounds,
			false,
			false,
		)
		const measurements = 2000
		var stop atomic.Bool
		var wg sync.WaitGroup
		wg.Go(func() {
			for !stop.Load() {
				agg.Finish(alice.Equivalent(), time.Now())
			}
		})
		var count uint64
		var sum N
		wg.Go(func() {
			for !stop.Load() {
				var data metricdata.Aggregation
				agg.ComputeAggregation(&data)
				hist, ok := data.(metricdata.Histogram[N])
				if !assert.True(t, ok) {
					return
				}
				for _, point := range hist.DataPoints {
					count += point.Count
					sum += point.Sum
					assert.Equal(t, []uint64{0, 0, point.Count, 0}, point.BucketCounts)
					assert.Equal(t, N(7)*N(point.Count), point.Sum)
					assert.Equal(t, metricdata.NewExtrema(N(7)), point.Min)
					assert.Equal(t, metricdata.NewExtrema(N(7)), point.Max)
				}
			}
		})
		for range measurements {
			agg.Measure(t.Context(), N(7), alice)
		}
		stop.Store(true)
		wg.Wait()
		for _, point := range finishHistogramPoints[N](t, agg.ComputeAggregation) {
			count += point.Count
			sum += point.Sum
		}
		assert.Equal(t, uint64(measurements), count)
		assert.Equal(t, N(7*measurements), sum)
	})
}

func TestFinishHistogramInitialMeasurementAdmission(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		reservoir := &blockingFinishSumReservoir{offered: make(chan struct{}), release: make(chan struct{})}
		agg := Builder[int64]{
			ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] { return reservoir },
		}.FinishExplicitBucketHistogram(
			finishHistogramBounds,
			false,
			false,
		)
		measured := make(chan struct{})
		go func() { defer close(measured); agg.Measure(t.Context(), 7, alice) }()
		<-reservoir.offered
		result := make(chan metricdata.Aggregation, 1)
		go func() { var data metricdata.Aggregation; agg.ComputeAggregation(&data); result <- data }()
		synctest.Wait()
		assert.Empty(t, result, "collection completed before the admitted measurement")
		close(reservoir.release)
		<-measured
		hist := (<-result).(metricdata.Histogram[int64])
		require.Len(t, hist.DataPoints, 1)
		assert.Equal(t, uint64(1), hist.DataPoints[0].Count)
		assert.Equal(t, int64(7), hist.DataPoints[0].Sum)
	})
}

func TestFinishHistogramCollectionTimestamp(t *testing.T) {
	t.Run("Int64", testFinishHistogramCollectionTimestamp[int64])
	t.Run("Float64", testFinishHistogramCollectionTimestamp[float64])
}

func testFinishHistogramCollectionTimestamp[N int64 | float64](t *testing.T) {
	t.Helper()
	var c clock
	t.Cleanup(c.Register())
	agg := Builder[N]{}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
	agg.Measure(t.Context(), N(1), alice)

	now = func() time.Time {
		timestamp := c.Now()
		now = c.Now
		// Create a series after the collection timestamp has been sampled.
		agg.Measure(t.Context(), N(2), bob)
		return timestamp
	}
	points := finishHistogramPoints[N](t, agg.ComputeAggregation)
	for _, point := range points {
		assert.False(t, point.StartTime.After(point.Time), "collection emitted an inverted interval")
	}
	require.Len(t, points, 1, "series created after timestamp capture must wait for the next collection")
	assert.Equal(t, alice, points[0].Attributes)
	assert.Equal(t, y2k, points[0].StartTime)
	assert.Equal(t, y2kPlus(1), points[0].Time)

	points = finishHistogramPoints[N](t, agg.ComputeAggregation)
	require.Len(t, points, 2)
	i := slices.IndexFunc(points, func(point metricdata.HistogramDataPoint[N]) bool {
		return point.Attributes.Equals(&bob)
	})
	require.NotEqual(t, -1, i)
	assert.Equal(t, y2kPlus(2), points[i].StartTime)
	assert.Equal(t, y2kPlus(3), points[i].Time)
	assert.Equal(t, N(2), points[i].Sum)
	assert.Empty(t, points[i].Exemplars)
}

func TestFinishHistogramDeltaDeferredSeries(t *testing.T) {
	t.Run("Int64", testFinishHistogramDeltaDeferredSeries[int64])
	t.Run("Float64", testFinishHistogramDeltaDeferredSeries[float64])
}

func testFinishHistogramDeltaDeferredSeries[N int64 | float64](t *testing.T) {
	t.Helper()
	for _, finished := range []bool{false, true} {
		t.Run(fmt.Sprint("finished=", finished), func(t *testing.T) {
			var c clock
			t.Cleanup(c.Register())
			agg := Builder[N]{
				Temporality: metricdata.DeltaTemporality,
			}.FinishExplicitBucketHistogram(
				finishHistogramBounds,
				false,
				false,
			)
			agg.Measure(t.Context(), N(1), alice)

			now = func() time.Time {
				now = c.Now
				// Record after snapshotting but before sampling collection time.
				agg.Measure(t.Context(), N(2), bob)
				finishTime := c.Now()
				if finished {
					agg.Finish(bob.Equivalent(), finishTime)
				}
				return c.Now()
			}
			points := finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, alice, points[0].Attributes)
			assert.Equal(t, y2k, points[0].StartTime)
			assert.Equal(t, y2kPlus(4), points[0].Time)

			points = finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, bob, points[0].Attributes)
			assert.Equal(t, N(2), points[0].Sum)
			assert.Equal(t, y2kPlus(2), points[0].StartTime)
			end := y2kPlus(5)
			if finished {
				end = y2kPlus(3)
			}
			assert.Equal(t, end, points[0].Time)
			assert.False(t, points[0].StartTime.After(points[0].Time))
			assert.Empty(t, finishHistogramPoints[N](t, agg.ComputeAggregation))

			agg.Measure(t.Context(), N(3), bob)
			points = finishHistogramPoints[N](t, agg.ComputeAggregation)
			require.Len(t, points, 1)
			assert.Equal(t, N(3), points[0].Sum)
			assert.Equal(t, y2kPlus(6), points[0].StartTime)
			assert.Equal(t, y2kPlus(8), points[0].Time)
		})
	}
}

func TestFinishHistogramCollectDefersNewSeries(t *testing.T) {
	barrier := &lastValueCollectionBarrier{entered: make(chan struct{}), release: make(chan struct{})}
	agg := Builder[int64]{
		Temporality: metricdata.DeltaTemporality, AggregationLimit: 4,
		ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] { return barrier },
	}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
	for _, attrs := range []attribute.Set{alice, bob, attribute.NewSet(userCarol)} {
		agg.Measure(t.Context(), 1, attrs)
	}
	result := make(chan metricdata.Aggregation, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		var data metricdata.Aggregation
		agg.ComputeAggregation(&data)
		result <- data
	}()
	var once sync.Once
	release := func() { once.Do(func() { close(barrier.release) }) }
	defer func() { release(); <-done }()
	select {
	case <-barrier.entered:
	case <-time.After(time.Second):
		t.Fatal("collection did not reach the second series")
	}
	dave := attribute.NewSet(userDave)
	agg.Measure(t.Context(), 123, dave)
	release()
	hist := (<-result).(metricdata.Histogram[int64])
	require.Len(t, hist.DataPoints, 3)
	points := finishHistogramPoints[int64](t, agg.ComputeAggregation)
	require.Len(t, points, 1)
	assert.Equal(t, dave, points[0].Attributes)
	assert.Equal(t, int64(123), points[0].Sum)
}

func TestFinishHistogramWaitCancellationDuringCreation(t *testing.T) {
	for _, populated := range []bool{false, true} {
		t.Run(fmt.Sprint("populated=", populated), func(t *testing.T) {
			entered, release := make(chan struct{}), make(chan struct{})
			store := newFinishHistogram(finishHistogramBounds, false, false, metricdata.CumulativeTemporality, 0,
				func(attrs attribute.Set) FilteredExemplarReservoir[int64] {
					if attrs.Equals(&bob) {
						close(entered)
						<-release
					}
					return DropReservoir[int64](attrs)
				})
			if populated {
				store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
			}

			measured := make(chan struct{})
			done := make(chan struct{})
			var shutdownErr error
			started := false
			go func() {
				defer close(measured)
				store.measure(t.Context(), 1, newLazyFilteredAttributes(bob, nil))
			}()
			defer func() {
				close(release)
				<-measured
				if started {
					<-done
				}
				if !populated {
					assert.Zero(t, store.values.Len(), "late publication retained a stopped series")
				}
				require.NoError(t, store.wait(t.Context()))
				assert.Zero(t, store.values.Len())
			}()
			select {
			case <-entered:
			case <-time.After(time.Second):
				t.Fatal("measurement did not enter reservoir construction")
			}

			ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
			defer cancel()
			started = true
			go func() {
				defer close(done)
				shutdownErr = store.wait(ctx)
			}()
			select {
			case <-done:
				assert.ErrorIs(t, shutdownErr, context.DeadlineExceeded)
			case <-time.After(time.Second):
				t.Error("shutdown ignored cancellation during reservoir construction")
			}
			assert.True(t, store.stopped.Load())
		})
	}
}

func TestFinishHistogramWaitCancellationDuringCollection(t *testing.T) {
	reservoir := &blockingFinishSumCollectionReservoir{
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}
	store := newFinishHistogram(finishHistogramBounds, false, false, metricdata.CumulativeTemporality, 0,
		func(attribute.Set) FilteredExemplarReservoir[int64] { return reservoir })
	store.measure(t.Context(), 1, newLazyFilteredAttributes(alice, nil))
	collected := make(chan struct{})
	done := make(chan struct{})
	var shutdownErr error
	started := false
	go func() {
		defer close(collected)
		var data metricdata.Aggregation
		store.collect(&data)
	}()
	defer func() {
		close(reservoir.release)
		<-collected
		if started {
			<-done
		}
		require.NoError(t, store.wait(t.Context()))
		assert.Zero(t, store.values.Len())
	}()
	select {
	case <-reservoir.entered:
	case <-time.After(time.Second):
		t.Fatal("collection did not enter reservoir")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 25*time.Millisecond)
	defer cancel()
	started = true
	go func() {
		defer close(done)
		shutdownErr = store.wait(ctx)
	}()
	select {
	case <-done:
		assert.ErrorIs(t, shutdownErr, context.DeadlineExceeded)
	case <-time.After(time.Second):
		t.Error("shutdown ignored cancellation during collection")
	}
	assert.True(t, store.stopped.Load())
}

func TestFinishHistogramShutdown(t *testing.T) {
	t.Run("Int64", testFinishHistogramShutdown[int64])
	t.Run("Float64", testFinishHistogramShutdown[float64])
}

func testFinishHistogramShutdown[N int64 | float64](t *testing.T) {
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			store := newFinishHistogram(finishHistogramBounds, false, false, temporality, 0, dropExemplars[N])
			store.measure(t.Context(), N(1), newLazyFilteredAttributes(alice, nil))
			var data metricdata.Aggregation
			require.Equal(t, 1, store.collect(&data))
			require.Len(t, store.points, 1)
			assert.Nil(t, store.points[0], "collection scratch must not retain series")
			store.measure(t.Context(), N(2), newLazyFilteredAttributes(alice, nil))
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			assert.ErrorIs(t, store.wait(ctx), context.Canceled)
			assert.True(t, store.stopped.Load())
			require.NoError(t, store.wait(t.Context()))
			require.NoError(t, store.wait(t.Context()))
			assert.Zero(t, store.values.Len())
			assert.Nil(t, store.points)
			store.measure(t.Context(), N(3), newLazyFilteredAttributes(bob, nil))
			store.finish(alice.Equivalent(), y2k)
			data = nil
			assert.Zero(t, store.collect(&data))
			assert.Nil(t, data)
		})
	}
}

func TestFinishHistogramIndependentAggregations(t *testing.T) {
	first := Builder[int64]{}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
	second := Builder[int64]{}.FinishExplicitBucketHistogram(finishHistogramBounds, false, false)
	first.Measure(t.Context(), 7, alice)
	second.Measure(t.Context(), 7, alice)
	finishTime := now()
	first.Finish(alice.Equivalent(), finishTime)
	second.Finish(alice.Equivalent(), finishTime)
	require.Len(t, finishHistogramPoints[int64](t, first.ComputeAggregation), 1)
	assert.Empty(t, finishHistogramPoints[int64](t, first.ComputeAggregation))
	points := finishHistogramPoints[int64](t, second.ComputeAggregation)
	require.Len(t, points, 1)
	assert.Equal(t, finishTime, points[0].Time)
	assert.Equal(t, int64(7), points[0].Sum)
	assert.Empty(t, finishHistogramPoints[int64](t, second.ComputeAggregation))
}
