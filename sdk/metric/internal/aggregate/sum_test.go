// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestSum(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64/DeltaSum", testDeltaSum[int64]())
	c.Reset()

	t.Run("Float64/DeltaSum", testDeltaSum[float64]())
	c.Reset()

	t.Run("Int64/CumulativeSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeSum[int64]()(t)
	})
	c.Reset()

	t.Run("Int64/CumulativeSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeSum[int64]()(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeSum", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		assert.False(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeSum[float64]()(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeSum/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		assert.True(t, x.PerSeriesStartTimestamps.Enabled())
		testCumulativeSum[float64]()(t)
	})
	c.Reset()

	t.Run("Int64/DeltaPrecomputedSum", testDeltaPrecomputedSum[int64]())
	c.Reset()

	t.Run("Float64/DeltaPrecomputedSum", testDeltaPrecomputedSum[float64]())
	c.Reset()

	t.Run("Int64/CumulativePrecomputedSum", testCumulativePrecomputedSum[int64]())
	c.Reset()

	t.Run("Float64/CumulativePrecomputedSum", testCumulativePrecomputedSum[float64]())
}

func testDeltaSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice, false},
				{ctx, -1, bob, false},
				{ctx, 1, alice, false},
				{ctx, 2, alice, false},
				{ctx, -10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(4),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(4),
							Value:      -11,
						},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
							Value:      3,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Delta sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      2,
						},
					},
				},
			},
		},
	})
}

func testCumulativeSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(mono)

	aliceStartTime := y2kPlus(0)
	bobStartTime := y2kPlus(0)
	overflowStartTime := y2kPlus(0)

	if x.PerSeriesStartTimestamps.Enabled() {
		aliceStartTime = y2kPlus(2)
		bobStartTime = y2kPlus(3)
		overflowStartTime = y2kPlus(6)
	}

	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice, false},
				{ctx, -1, bob, false},
				{ctx, 1, alice, false},
				{ctx, 2, alice, false},
				{ctx, -10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       y2kPlus(4),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       y2kPlus(4),
							Value:      -11,
						},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       y2kPlus(5),
							Value:      14,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       y2kPlus(5),
							Value:      -8,
						},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       y2kPlus(7),
							Value:      14,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       y2kPlus(7),
							Value:      -8,
						},
						{
							Attributes: overflowSet,
							StartTime:  overflowStartTime,
							Time:       y2kPlus(7),
							Value:      2,
						},
					},
				},
			},
		},
	})
}

func testDeltaPrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice, false},
				{ctx, -1, bob, false},
				{ctx, 1, fltrAlice, false},
				{ctx, 2, alice, false},
				{ctx, -10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(4),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(4),
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice, false},
				{ctx, 10, alice, false},
				{ctx, 3, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
							Value:      7,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(7),
							Value:      14,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Precomputed sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(8),
							Time:       y2kPlus(12),
							Value:      2,
						},
					},
				},
			},
		},
	})
}

func testCumulativePrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{},
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice, false},
				{ctx, -1, bob, false},
				{ctx, 1, fltrAlice, false},
				{ctx, 2, alice, false},
				{ctx, -10, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(4),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(4),
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice, false},
				{ctx, 10, alice, false},
				{ctx, 3, bob, false},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(7),
							Value:      11,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(7),
							Value:      3,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{},
			// Precomputed sums are expected to reset.
			expect: output{
				n: 0,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints:  []metricdata.DataPoint[N]{},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(12),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(12),
							Value:      2,
						},
					},
				},
			},
		},
	})
}

func TestSumConcurrentSafe(t *testing.T) {
	t.Run("Int64/DeltaSum", testDeltaSumConcurrentSafe[int64]())
	t.Run("Float64/DeltaSum", testDeltaSumConcurrentSafe[float64]())
	t.Run("Int64/CumulativeSum", testCumulativeSumConcurrentSafe[int64]())
	t.Run("Float64/CumulativeSum", testCumulativeSumConcurrentSafe[float64]())
	t.Run("Int64/DeltaPrecomputedSum", testDeltaPrecomputedSumConcurrentSafe[int64]())
	t.Run("Float64/DeltaPrecomputedSum", testDeltaPrecomputedSumConcurrentSafe[float64]())
	t.Run("Int64/CumulativePrecomputedSum", testCumulativePrecomputedSumConcurrentSafe[int64]())
	t.Run("Float64/CumulativePrecomputedSum", testCumulativePrecomputedSumConcurrentSafe[float64]())
}

//nolint:revive // isPrecomputed is used for configuring validation
func validateSum[N int64 | float64](isPrecomputed bool) func(t *testing.T, aggs []metricdata.Aggregation) {
	return func(t *testing.T, aggs []metricdata.Aggregation) {
		sums := make(map[attribute.Set]N)
		for i, agg := range aggs {
			s, ok := agg.(metricdata.Sum[N])
			require.True(t, ok)
			require.LessOrEqual(t, len(s.DataPoints), 3, "AggregationLimit of 3 exceeded in a single cycle")
			for _, dp := range s.DataPoints {
				if s.Temporality == metricdata.DeltaTemporality {
					sums[dp.Attributes] += dp.Value
				} else if i == len(aggs)-1 {
					sums[dp.Attributes] = dp.Value
				}
			}
		}

		if isPrecomputed {
			// Precomputed Sums clear the state when collected concurrently. Due to hot/cold overlap
			// during flush, the sum drops intermediate updates, so the final calculation won't cleanly
			// add up to the total number of operations performed by the workers. Therefore, skip exact
			// invariant check, verifying only that limits and map updates occurred safely.
			return
		}

		var total N
		for _, val := range sums {
			total += val
		}

		assertSumEqual[N](t, expectedConcurrentSum[N](), total)
	}
}

func testDeltaSumConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(false)
	return testAggregationConcurrentSafe[N](in, out, validateSum[N](false))
}

func testCumulativeSumConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(false)
	return testAggregationConcurrentSafe[N](in, out, validateSum[N](false))
}

func testDeltaPrecomputedSumConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedSum(false)
	return testAggregationConcurrentSafe[N](in, out, validateSum[N](true))
}

func testCumulativePrecomputedSumConcurrentSafe[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedSum(false)
	return testAggregationConcurrentSafe[N](in, out, validateSum[N](true))
}

func BenchmarkSum(b *testing.B) {
	// The monotonic argument is only used to annotate the Sum returned from
	// the Aggregation method. It should not have an effect on operational
	// performance, therefore, only monotonic=false is benchmarked here.
	b.Run("Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.Sum(false)
	}))
	b.Run("Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.Sum(false)
	}))
	b.Run("Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.Sum(false)
	}))
	b.Run("Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.Sum(false)
	}))

	b.Run("Precomputed/Int64/Cumulative", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.CumulativeTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Int64/Delta", benchmarkAggregate(func() (Measure[int64], ComputeAggregation) {
		return Builder[int64]{
			Temporality: metricdata.DeltaTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Float64/Cumulative", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.CumulativeTemporality,
		}.PrecomputedSum(false)
	}))
	b.Run("Precomputed/Float64/Delta", benchmarkAggregate(func() (Measure[float64], ComputeAggregation) {
		return Builder[float64]{
			Temporality: metricdata.DeltaTemporality,
		}.PrecomputedSum(false)
	}))
}

func TestCumulativeSumFinishResetsStartTime(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())
	t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
	assert.True(t, x.PerSeriesStartTimestamps.Enabled())

	in, out := Builder[int64]{
		Temporality: metricdata.CumulativeTemporality,
		Filter:      attrFltr,
	}.Sum(false)

	ctx := t.Context()
	in(ctx, 1, alice, false)

	var got metricdata.Aggregation = metricdata.Sum[int64]{}
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2kPlus(1),
				Time:       y2kPlus(2),
				Value:      1,
			},
		},
	}, got)

	in(ctx, 0, alice, true)
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2kPlus(1),
				Time:       y2kPlus(3),
				Value:      1,
			},
		},
	}, got)

	assert.Equal(t, 0, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.DataPoint[int64]{},
	}, got)

	in(ctx, 3, alice, false)
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2kPlus(5),
				Time:       y2kPlus(6),
				Value:      3,
			},
		},
	}, got)
}

func TestDeltaSumFinishExportsFinalPoint(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	in, out := Builder[int64]{
		Temporality: metricdata.DeltaTemporality,
		Filter:      attrFltr,
	}.Sum(false)

	ctx := t.Context()
	in(ctx, 1, alice, false)
	in(ctx, 0, alice, true)

	var got metricdata.Aggregation = metricdata.Sum[int64]{}
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2k,
				Time:       y2kPlus(2),
				Value:      1,
			},
		},
	}, got)

	assert.Equal(t, 0, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.DataPoint[int64]{},
	}, got)
}

func TestDeltaSumFinishRevivePreservesData(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	in, out := Builder[int64]{
		Temporality: metricdata.DeltaTemporality,
		Filter:      attrFltr,
	}.Sum(false)

	ctx := t.Context()
	in(ctx, 1, alice, false)
	in(ctx, 0, alice, true)
	in(ctx, 2, alice, false)

	var got metricdata.Aggregation = metricdata.Sum[int64]{}
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2k,
				Time:       y2kPlus(2),
				Value:      3,
			},
		},
	}, got)
}

func TestCumulativeSumFinishRevivePreservesData(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())
	t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
	assert.True(t, x.PerSeriesStartTimestamps.Enabled())

	in, out := Builder[int64]{
		Temporality: metricdata.CumulativeTemporality,
		Filter:      attrFltr,
	}.Sum(false)

	ctx := t.Context()
	in(ctx, 1, alice, false)
	in(ctx, 0, alice, true)
	in(ctx, 2, alice, false)

	var got metricdata.Aggregation = metricdata.Sum[int64]{}
	assert.Equal(t, 1, out(&got))
	metricdatatest.AssertAggregationsEqual(t, metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints: []metricdata.DataPoint[int64]{
			{
				Attributes: fltrAlice,
				StartTime:  y2kPlus(1),
				Time:       y2kPlus(2),
				Value:      3,
			},
		},
	}, got)
}
