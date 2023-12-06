// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var (
	bounds   = []float64{1, 5}
	noMinMax = false
)

func TestHistogram(t *testing.T) {
	t.Cleanup(mockTime(now))

	t.Run("Int64/Delta/Sum", testDeltaHist[int64](conf[int64]{hPt: hPointSummed[int64]}))
	t.Run("Int64/Delta/NoSum", testDeltaHist[int64](conf[int64]{noSum: true, hPt: hPoint[int64]}))
	t.Run("Float64/Delta/Sum", testDeltaHist[float64](conf[float64]{hPt: hPointSummed[float64]}))
	t.Run("Float64/Delta/NoSum", testDeltaHist[float64](conf[float64]{noSum: true, hPt: hPoint[float64]}))

	t.Run("Int64/Cumulative/Sum", testCumulativeHist[int64](conf[int64]{hPt: hPointSummed[int64]}))
	t.Run("Int64/Cumulative/NoSum", testCumulativeHist[int64](conf[int64]{noSum: true, hPt: hPoint[int64]}))
	t.Run("Float64/Cumulative/Sum", testCumulativeHist[float64](conf[float64]{hPt: hPointSummed[float64]}))
	t.Run("Float64/Cumulative/NoSum", testCumulativeHist[float64](conf[float64]{noSum: true, hPt: hPoint[float64]}))
}

type conf[N int64 | float64] struct {
	noSum bool
	hPt   func(attribute.Set, N, uint64) metricdata.HistogramDataPoint[N]
}

func testDeltaHist[N int64 | float64](c conf[N]) func(t *testing.T) {
	in, out := Builder[N]{
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
				{ctx, 2, alice},
				{ctx, 10, bob},
				{ctx, 2, alice},
				{ctx, 2, alice},
				{ctx, 10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 3),
						c.hPt(fltrBob, 10, 2),
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 10, 1),
						c.hPt(fltrBob, 3, 1),
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
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 1, 1),
						c.hPt(fltrBob, 1, 1),
						c.hPt(overflowSet, 1, 2),
					},
				},
			},
		},
	})
}

func testCumulativeHist[N int64 | float64](c conf[N]) func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
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
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.HistogramDataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice},
				{ctx, 10, bob},
				{ctx, 2, alice},
				{ctx, 2, alice},
				{ctx, 10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 3),
						c.hPt(fltrBob, 10, 2),
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 2, alice},
				{ctx, 10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 4),
						c.hPt(fltrBob, 10, 3),
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
						c.hPt(fltrAlice, 2, 4),
						c.hPt(fltrBob, 10, 3),
					},
				},
			},
		},
		{
			input: []arg[N]{
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Histogram[N]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[N]{
						c.hPt(fltrAlice, 2, 4),
						c.hPt(fltrBob, 10, 3),
						c.hPt(overflowSet, 1, 2),
					},
				},
			},
		},
	})
}

// hPointSummed returns an HistogramDataPoint that started and ended now with
// multi number of measurements values v. It includes a min and max (set to v).
func hPointSummed[N int64 | float64](a attribute.Set, v N, multi uint64) metricdata.HistogramDataPoint[N] {
	idx := sort.SearchFloat64s(bounds, float64(v))
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint[N]{
		Attributes:   a,
		StartTime:    now(),
		Time:         now(),
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
func hPoint[N int64 | float64](a attribute.Set, v N, multi uint64) metricdata.HistogramDataPoint[N] {
	idx := sort.SearchFloat64s(bounds, float64(v))
	counts := make([]uint64, len(bounds)+1)
	counts[idx] += multi
	return metricdata.HistogramDataPoint[N]{
		Attributes:   a,
		StartTime:    now(),
		Time:         now(),
		Count:        multi,
		Bounds:       bounds,
		BucketCounts: counts,
		Min:          metricdata.NewExtrema(v),
		Max:          metricdata.NewExtrema(v),
	}
}

func TestBucketsBin(t *testing.T) {
	t.Run("Int64", testBucketsBin[int64]())
	t.Run("Float64", testBucketsBin[float64]())
}

func testBucketsBin[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		b := newBuckets[N](3)
		assertB := func(counts []uint64, count uint64, min, max N) {
			t.Helper()
			assert.Equal(t, counts, b.counts)
			assert.Equal(t, count, b.count)
			assert.Equal(t, min, b.min)
			assert.Equal(t, max, b.max)
		}

		assertB([]uint64{0, 0, 0}, 0, 0, 0)
		b.bin(1, 2)
		assertB([]uint64{0, 1, 0}, 1, 0, 2)
		b.bin(0, -1)
		assertB([]uint64{1, 1, 0}, 2, -1, 2)
	}
}

func TestBucketsSum(t *testing.T) {
	t.Run("Int64", testBucketsSum[int64]())
	t.Run("Float64", testBucketsSum[float64]())
}

func testBucketsSum[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		b := newBuckets[N](3)

		var want N
		assert.Equal(t, want, b.total)

		b.sum(2)
		want = 2
		assert.Equal(t, want, b.total)

		b.sum(-1)
		want = 1
		assert.Equal(t, want, b.total)
	}
}

func TestHistogramImmutableBounds(t *testing.T) {
	b := []float64{0, 1, 2}
	cpB := make([]float64, len(b))
	copy(cpB, b)

	h := newHistogram[int64](b, false, false, 0)
	require.Equal(t, cpB, h.bounds)

	b[0] = 10
	assert.Equal(t, cpB, h.bounds, "modifying the bounds argument should not change the bounds")

	h.measure(context.Background(), 5, alice)

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	h.cumulative(&data)
	hdp := data.(metricdata.Histogram[int64]).DataPoints[0]
	hdp.Bounds[1] = 10
	assert.Equal(t, cpB, h.bounds, "modifying the Aggregation bounds should not change the bounds")
}

func TestCumulativeHistogramImutableCounts(t *testing.T) {
	h := newHistogram[int64](bounds, noMinMax, false, 0)
	h.measure(context.Background(), 5, alice)

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	h.cumulative(&data)
	hdp := data.(metricdata.Histogram[int64]).DataPoints[0]

	require.Equal(t, hdp.BucketCounts, h.values[alice].counts)

	cpCounts := make([]uint64, len(hdp.BucketCounts))
	copy(cpCounts, hdp.BucketCounts)
	hdp.BucketCounts[0] = 10
	assert.Equal(t, cpCounts, h.values[alice].counts, "modifying the Aggregator bucket counts should not change the Aggregator")
}

func TestDeltaHistogramReset(t *testing.T) {
	t.Cleanup(mockTime(now))

	h := newHistogram[int64](bounds, noMinMax, false, 0)

	var data metricdata.Aggregation = metricdata.Histogram[int64]{}
	require.Equal(t, 0, h.delta(&data))
	require.Len(t, data.(metricdata.Histogram[int64]).DataPoints, 0)

	h.measure(context.Background(), 1, alice)

	expect := metricdata.Histogram[int64]{Temporality: metricdata.DeltaTemporality}
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPointSummed[int64](alice, 1, 1)}
	h.delta(&data)
	metricdatatest.AssertAggregationsEqual(t, expect, data)

	// The attr set should be forgotten once Aggregations is called.
	expect.DataPoints = nil
	assert.Equal(t, 0, h.delta(&data))
	assert.Len(t, data.(metricdata.Histogram[int64]).DataPoints, 0)

	// Aggregating another set should not affect the original (alice).
	h.measure(context.Background(), 1, bob)
	expect.DataPoints = []metricdata.HistogramDataPoint[int64]{hPointSummed[int64](bob, 1, 1)}
	h.delta(&data)
	metricdatatest.AssertAggregationsEqual(t, expect, data)
}

func BenchmarkHistogram(b *testing.B) {
	b.Run("Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
	}))
	b.Run("Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
	}))
	b.Run("Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
	}))
	b.Run("Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.ExplicitBucketHistogram(bounds, noMinMax, false)
	}))
}
