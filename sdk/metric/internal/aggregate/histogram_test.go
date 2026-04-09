// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var (
	bounds   = []float64{1, 5}
	noMinMax = false
)

func TestHistogram(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64/Delta/Sum", testDeltaHist[int64](conf[int64]{hPt: hPointSummed[int64]}))
	c.Reset()
	t.Run("Int64/Delta/NoSum", testDeltaHist[int64](conf[int64]{noSum: true, hPt: hPoint[int64]}))
	c.Reset()
	t.Run("Float64/Delta/Sum", testDeltaHist[float64](conf[float64]{hPt: hPointSummed[float64]}))
	c.Reset()
	t.Run("Float64/Delta/NoSum", testDeltaHist[float64](conf[float64]{noSum: true, hPt: hPoint[float64]}))
	c.Reset()

	t.Run("Int64/Cumulative/Sum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[int64](conf[int64]{hPt: hPointSummed[int64]})(t)
	})
	c.Reset()

	t.Run("Int64/Cumulative/Sum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[int64](conf[int64]{hPt: hPointSummed[int64]})(t)
	})
	c.Reset()

	t.Run("Int64/Cumulative/NoSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[int64](conf[int64]{noSum: true, hPt: hPoint[int64]})(t)
	})
	c.Reset()

	t.Run("Int64/Cumulative/NoSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[int64](conf[int64]{noSum: true, hPt: hPoint[int64]})(t)
	})
	c.Reset()

	t.Run("Float64/Cumulative/Sum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[float64](conf[float64]{hPt: hPointSummed[float64]})(t)
	})
	c.Reset()

	t.Run("Float64/Cumulative/Sum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[float64](conf[float64]{hPt: hPointSummed[float64]})(t)
	})
	c.Reset()

	t.Run("Float64/Cumulative/NoSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[float64](
			conf[float64]{noSum: true, hPt: hPoint[float64]},
		)(
			t,
		)
	})
	c.Reset()

	t.Run("Float64/Cumulative/NoSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeHist[float64](
			conf[float64]{noSum: true, hPt: hPoint[float64]},
		)(
			t,
		)
	})
	c.Reset()
}

type conf[N int64 | float64] struct {
	noSum bool
	hPt   func(attribute.Set, N, uint64, time.Time, time.Time) metricdata.HistogramDataPoint[N]
}

func testDeltaHist[N int64 | float64](c conf[N]) func(t *testing.T) {
	in, _, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExplicitBucketHistogram(bounds, noMinMax, c.noSum)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice, false},
				{ctx, 10, bob, false},
				{ctx, 2, alice, false},
				{ctx, 2, alice, false},
				{ctx, 10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 3, y2kPlus(1), y2kPlus(2)),
						c.hPt(fltrBob, 10, 2, y2kPlus(1), y2kPlus(2)),
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice, false},
				{ctx, 3, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 10, 1, y2kPlus(2), y2kPlus(3)),
						c.hPt(fltrBob, 3, 1, y2kPlus(2), y2kPlus(3)),
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Delta histograms are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice, false},
				{ctx, 1, bob, false},
				// These will exceed cardinality limit.
				{ctx, 1, carol, false},
				{ctx, 1, dave, false},
			},
			expect: output{
				n: 3,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 1, 1, y2kPlus(4), y2kPlus(5)),
						c.hPt(fltrBob, 1, 1, y2kPlus(4), y2kPlus(5)),
						c.hPt(overflowSet, 1, 2, y2kPlus(4), y2kPlus(5)),
					},
				},
			},
		},
	})
}

func testCumulativeHist[N int64 | float64](c conf[N]) func(t *testing.T) {
	in, _, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExplicitBucketHistogram(bounds, noMinMax, c.noSum)

	aliceStartTime := y2kPlus(0)
	bobStartTime := y2kPlus(0)
	overflowStartTime := y2kPlus(0)

	if x.PerSeriesStartTimestamps.Enabled() {
		aliceStartTime = y2kPlus(2)
		bobStartTime = y2kPlus(3)
		overflowStartTime = y2kPlus(7)
	}

	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice, false},
				{ctx, 10, bob, false},
				{ctx, 2, alice, false},
				{ctx, 2, alice, false},
				{ctx, 10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 3, aliceStartTime, y2kPlus(4)),
						c.hPt(fltrBob, 10, 2, bobStartTime, y2kPlus(4)),
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice, false},
				{ctx, 10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 4, aliceStartTime, y2kPlus(5)),
						c.hPt(fltrBob, 10, 3, bobStartTime, y2kPlus(5)),
					},
				},
			},
		},
		{
			input: []arg[N]{},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 4, aliceStartTime, y2kPlus(6)),
						c.hPt(fltrBob, 10, 3, bobStartTime, y2kPlus(6)),
					},
				},
			},
		},
		{
			input: []arg[N]{
				// These will exceed cardinality limit.
				{ctx, 1, carol, false},
				{ctx, 1, dave, false},
			},
			expect: output{
				n: 3,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 4, aliceStartTime, y2kPlus(8)),
						c.hPt(fltrBob, 10, 3, bobStartTime, y2kPlus(8)),
						c.hPt(overflowSet, 1, 2, overflowStartTime, y2kPlus(8)),
					},
				},
			},
		},
	})
}

func TestHistogramConcurrentSafe(t *testing.T) {
	t.Run("Int64/Delta", testDeltaHistConcurrentSafe[int64]())
	t.Run("Float64/Delta", testDeltaHistConcurrentSafe[float64]())
	t.Run("Int64/Cumulative", testCumulativeHistConcurrentSafe[int64]())
	t.Run("Float64/Cumulative", testCumulativeHistConcurrentSafe[float64]())
}

func validateHistogram[N int64 | float64](t *testing.T, aggs []metricdata.Aggregation) {
	sums := make(map[attribute.Set]N)
	counts := make(map[attribute.Set]uint64)
	bucketCounts := make(map[attribute.Set][]uint64)

	for i, agg := range aggs {
		s, ok := agg.(metricdata.Histogram[N])
		require.True(t, ok)
		require.LessOrEqual(t, len(s.DataPoints), 3, "AggregationLimit of 3 exceeded in a single cycle")
		for _, dp := range s.DataPoints {
			if s.Temporality == metricdata.DeltaTemporality {
				sums[dp.Attributes] += dp.Sum
				counts[dp.Attributes] += dp.Count
				if bucketCounts[dp.Attributes] == nil {
					bucketCounts[dp.Attributes] = make([]uint64, len(dp.BucketCounts))
				}
				for idx, c := range dp.BucketCounts {
					bucketCounts[dp.Attributes][idx] += c
				}
			} else if i == len(aggs)-1 {
				sums[dp.Attributes] = dp.Sum
				counts[dp.Attributes] = dp.Count
				bucketCounts[dp.Attributes] = make([]uint64, len(dp.BucketCounts))
				copy(bucketCounts[dp.Attributes], dp.BucketCounts)
			}
		}
	}

	var totalSum N
	var totalCount uint64
	totalBuckets := make([]uint64, 4)

	for _, val := range sums {
		totalSum += val
	}
	for _, val := range counts {
		totalCount += val
	}
	for _, bc := range bucketCounts {
		for idx, c := range bc {
			if idx < len(totalBuckets) {
				totalBuckets[idx] += c
			}
		}
	}

	assertSumEqual[N](t, expectedConcurrentSum[N](), totalSum)
	assert.Equal(t, expectedConcurrentCount, totalCount)

	var expectedBuckets []uint64
	switch any(*new(N)).(type) {
	case float64:
		// Float sequence: 2.5, 6.1, 4.4, 10.0, 22.0, -3.5, -6.5, 3.0, -6.0
		// Bounds {0, 2, 4}:
		// (-inf, 0]: -3.5, -6.5, -6.0 (3x)
		// (0, 2]: none (0x)
		// (2, 4]: 2.5, 3.0 (2x)
		// (4, +inf): 6.1, 4.4, 10.0, 22.0 (4x)
		// 10 full loops per goroutine * 10 goroutines = 100x
		expectedBuckets = []uint64{300, 0, 200, 400}
	default:
		// Int sequence: 2, 6, 4, 10, 22, -3, -6, 3, -6
		// Bounds {0, 2, 4}:
		// (-inf, 0]: -3, -6, -6 (3x)
		// (0, 2]: 2 (1x)
		// (2, 4]: 4, 3 (2x)
		// (4, +inf): 6, 10, 22 (3x)
		// 10 full loops per goroutine * 10 goroutines = 100x
		expectedBuckets = []uint64{300, 100, 200, 300}
	}
	assert.Equal(t, expectedBuckets, totalBuckets)
}

func testCumulativeHistConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, _, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExplicitBucketHistogram([]float64{0, 2, 4}, false, false)
	return testAggregationConcurrentSafe[N](in, out, validateHistogram[N])
}

func testDeltaHistConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, _, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.ExplicitBucketHistogram([]float64{0, 2, 4}, false, false)
	return testAggregationConcurrentSafe[N](in, out, validateHistogram[N])
}

// hPointSummed returns an HistogramDataPoint that started and ended now with
// multi number of measurements values v. It includes a min and max (set to v).
func hPointSummed[N int64 | float64](
	a attribute.Set,
	v N,
	multi uint64,
	start, t time.Time,
) metricdata.HistogramDataPoint[N] {
	idx := sort.SearchFloat64s(bounds, float64(v))
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint[N]{
		Attributes:   a,
		StartTime:    start,
		Time:         t,
		Count:        multi,
		Bounds:       bounds,
		BucketCounts: counts,
		Min:          metricdata.NewExtrema(v),
		Max:          metricdata.NewExtrema(v),
		Sum:          v * N(multi),
	}
}

// hPoint returns an HistogramDataPoint that started and ended now with multi
// number of measurements values v. It includes a min and max (set to v).
func hPoint[N int64 | float64](
	a attribute.Set,
	v N,
	multi uint64,
	start, t time.Time,
) metricdata.HistogramDataPoint[N] {
	idx := sort.SearchFloat64s(bounds, float64(v))
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint[N]{
		Attributes:   a,
		StartTime:    start,
		Time:         t,
		Count:        multi,
		Bounds:       bounds,
		BucketCounts: counts,
		Min:          metricdata.NewExtrema(v),
		Max:          metricdata.NewExtrema(v),
	}
}

func TestHistogramImmutableBounds(t *testing.T) {
	b := []float64{0, 1, 2}
	cpB := make([]float64, len(b))
	copy(cpB, b)

	h := newCumulativeHistogram[int64](b, false, false, 0, dropExemplars[int64])
	require.Equal(t, cpB, h.bounds)

	b[0] = 10
	assert.Equal(t, cpB, h.bounds, "modifying the bounds argument should not change the bounds")

	h.measure(t.Context(), 5, alice, nil, false)

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	h.collect(&data)
	hdp := data.(metricdata.Histogram[int64]).DataPoints[0]
	hdp.Bounds[1] = 10
	assert.Equal(t, cpB, h.bounds, "modifying the Aggregation bounds should not change the bounds")
}

func TestCumulativeHistogramImmutableCounts(t *testing.T) {
	h := newCumulativeHistogram[int64](bounds, noMinMax, false, 0, dropExemplars[int64])
	h.measure(t.Context(), 5, alice, nil, false)

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	h.collect(&data)
	hdp := data.(metricdata.Histogram[int64]).DataPoints[0]

	hPt, ok := h.values.Load(alice.Equivalent())
	require.True(t, ok)
	hcHistPt := hPt.(*hotColdHistogramPoint[int64])
	readIdx := hcHistPt.hcwg.swapHotAndWait()
	var bucketCounts []uint64
	hcHistPt.hotColdPoint[readIdx].loadCountsInto(&bucketCounts)
	require.Equal(t, hdp.BucketCounts, bucketCounts)
	hotIdx := (readIdx + 1) % 2
	hcHistPt.hotColdPoint[readIdx].mergeIntoAndReset(&hcHistPt.hotColdPoint[hotIdx], noMinMax, false)

	cpCounts := make([]uint64, len(hdp.BucketCounts))
	copy(cpCounts, hdp.BucketCounts)
	hdp.BucketCounts[0] = 10
	hPt, ok = h.values.Load(alice.Equivalent())
	require.True(t, ok)
	hcHistPt = hPt.(*hotColdHistogramPoint[int64])
	readIdx = hcHistPt.hcwg.swapHotAndWait()
	hcHistPt.hotColdPoint[readIdx].loadCountsInto(&bucketCounts)
	assert.Equal(
		t,
		cpCounts,
		bucketCounts,
		"modifying the Aggregator bucket counts should not change the Aggregator",
	)
}

func TestDeltaHistogramReset(t *testing.T) {
	orig := now
	now = func() time.Time { return y2k }
	t.Cleanup(func() { now = orig })

	h := newDeltaHistogram[int64](bounds, noMinMax, false, 0, dropExemplars[int64])

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	require.Equal(t, 0, h.collect(&data))
	require.Empty(t, data.(metricdata.Histogram[int64]).DataPoints)

	h.measure(t.Context(), 1, alice, nil, false)

	expect := metricdata.Histogram[int64]{Temporality: metricdata.DeltaTemporality}
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPointSummed[int64](alice, 1, 1, now(), now())}
	h.collect(&data)
	metricdatatest.AssertAggregationsEqual(t, expect, data)

	// The attr set should be forgotten once Aggregations is called.
	expect.DataPoints = nil
	assert.Equal(t, 0, h.collect(&data))
	assert.Empty(t, data.(metricdata.Histogram[int64]).DataPoints)

	// Aggregating another set should not affect the original (alice).
	h.measure(t.Context(), 1, bob, nil, false)
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPointSummed[int64](bob, 1, 1, now(), now())}
	h.collect(&data)
	metricdatatest.AssertAggregationsEqual(t, expect, data)
}

func BenchmarkHistogram(b *testing.B) {
	b.Run("Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		in, _, out := Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
		return in, out
	}))
	b.Run("Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		in, _, out := Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
		return in, out
	}))
	b.Run("Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		in, _, out := Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
		return in, out
	}))
	b.Run("Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		in, _, out := Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
		return in, out
	}))
}

func TestCumulativeHistogramFinishResetsStartTime(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())
	t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
	assert.True(t, x.PerSeriesStartTimestamps.Enabled())

	h := newCumulativeHistogram[int64](bounds, noMinMax, false, 0, dropExemplars[int64])
	ctx := t.Context()

	h.measure(ctx, 5, alice, nil, false)

	var got metricdata.Aggregation = metricdata.Histogram[int64]{}
	assert.Equal(t, 1, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.HistogramDataPoint[int64]{
			hPointSummed[int64](alice, 5, 1, y2kPlus(1), y2kPlus(2)),
		},
	}, got)

	h.measure(ctx, 0, alice, nil, true)
	assert.Equal(t, 1, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.HistogramDataPoint[int64]{
			hPointSummed[int64](alice, 5, 1, y2kPlus(1), y2kPlus(3)),
		},
	}, got)

	assert.Equal(t, 0, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{},
	}, got)

	h.measure(ctx, 5, alice, nil, false)
	assert.Equal(t, 1, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.HistogramDataPoint[int64]{
			hPointSummed[int64](alice, 5, 1, y2kPlus(5), y2kPlus(6)),
		},
	}, got)
}

func TestDeltaHistogramFinishExportsFinalPoint(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	h := newDeltaHistogram[int64](bounds, noMinMax, false, 0, dropExemplars[int64])
	ctx := t.Context()

	h.measure(ctx, 5, alice, nil, false)
	h.measure(ctx, 0, alice, nil, true)

	var got metricdata.Aggregation = metricdata.Histogram[int64]{}
	assert.Equal(t, 1, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints: []metricdata.HistogramDataPoint[int64]{
			hPointSummed[int64](alice, 5, 1, y2k, y2kPlus(1)),
		},
	}, got)

	assert.Equal(t, 0, h.collect(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Histogram[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{},
	}, got)
}
