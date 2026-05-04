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
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
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
							Time:       y2kPlus(2),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      -11,
						},
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
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
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
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
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
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
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
				{ctx, 10, alice},
				{ctx, 3, bob},
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
				{ctx, 1, carol},
				{ctx, 1, dave},
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
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
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
							Time:       y2kPlus(2),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
							Value:      7,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
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
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(4),
							Time:       y2kPlus(5),
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
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
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
							Time:       y2kPlus(2),
							Value:      4,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(2),
							Value:      -11,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
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
							Time:       y2kPlus(3),
							Value:      11,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(3),
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
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
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
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(5),
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(5),
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

func TestDeltaSumLazyCleanupGoals(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64", testDeltaSumLazyCleanupGoals[int64]())
	c.Reset()
	t.Run("Float64", testDeltaSumLazyCleanupGoals[float64]())
}

func testDeltaSumLazyCleanupGoals[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(1),
							Value:      1,
						},
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(1),
							Value:      1,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      1,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
							Value:      1,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: attribute.NewSet(userCarol),
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      1,
						},
						{
							Attributes: attribute.NewSet(userDave),
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      1,
						},
					},
				},
			},
		},
	})
}

func TestDeltaSumLazyCleanupEarlyOverflow(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64", testDeltaSumLazyCleanupEarlyOverflow[int64]())
	c.Reset()
	t.Run("Float64", testDeltaSumLazyCleanupEarlyOverflow[float64]())
}

func testDeltaSumLazyCleanupEarlyOverflow[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(1),
							Value:      1,
						},
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(1),
							Value:      1,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, alice},
			},
			expect: output{
				n: 1,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      1,
						},
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 2,
				agg: metricdata.Sum[N]{
					IsMonotonic: mono,
					Temporality: metricdata.DeltaTemporality,
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: attribute.NewSet(userCarol),
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
							Value:      1,
						},
						{
							Attributes: attribute.NewSet(userDave),
							StartTime:  y2kPlus(2),
							Time:       y2kPlus(3),
							Value:      1,
						},
					},
				},
			},
		},
	})
}

func TestDeltaSumLazyCleanupExistingOverflow(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64", testDeltaSumLazyCleanupExistingOverflow[int64]())
	c.Reset()
	t.Run("Float64", testDeltaSumLazyCleanupExistingOverflow[float64]())
}

func testDeltaSumLazyCleanupExistingOverflow[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		mono := false
		in, out := Builder[N]{
			Temporality:      metricdata.DeltaTemporality,
			Filter:           attrFltr,
			AggregationLimit: 3,
		}.Sum(mono)
		ctx := t.Context()

		// Step 1: Measure A.
		in(ctx, 1, alice)
		got := new(metricdata.Aggregation)
		out(got)

		// Step 2: Measure B, C.
		in(ctx, 1, bob)
		in(ctx, 1, carol)
		out(got)

		// Step 3: Measure carol, dave, then alice.
		in(ctx, 1, carol)
		in(ctx, 1, dave)
		in(ctx, 1, alice)
		out(got)

		s, ok := (*got).(metricdata.Sum[N])
		require.True(t, ok)

		// Normalize exemplars to nil for comparison.
		for i := range s.DataPoints {
			s.DataPoints[i].Exemplars = nil
		}

		expected := []metricdata.DataPoint[N]{
			{
				Attributes: attribute.NewSet(userCarol),
				StartTime:  y2kPlus(2),
				Time:       y2kPlus(3),
				Value:      1,
			},
			{
				Attributes: attribute.NewSet(userDave),
				StartTime:  y2kPlus(2),
				Time:       y2kPlus(3),
				Value:      1,
			},
			{
				Attributes: overflowSet,
				StartTime:  y2kPlus(2),
				Time:       y2kPlus(3),
				Value:      1,
			},
		}

		require.Len(t, s.DataPoints, 3, "incorrect data size")
		assert.ElementsMatch(t, expected, s.DataPoints)
	}
}

type noopRes[N int64 | float64] struct{}

func (noopRes[N]) Offer(context.Context, N, []attribute.KeyValue) {}
func (noopRes[N]) Collect(*[]exemplar.Exemplar)                   {}

type sumBenchmarker interface {
	measure(ctx context.Context, value int64, fltrAttr attribute.Set)
	pseudoCollect()
}

type lazyBenchmarker struct {
	*deltaSum[int64]
}

func (b *lazyBenchmarker) pseudoCollect() {
	readIdx := b.hcwg.swapHotAndWait()
	b.hotColdValMap[readIdx].values.Range(func(_, _ any) bool { return true })
	b.hotColdValMap[readIdx].values.Clear()
}

func (b *lazyBenchmarker) measure(ctx context.Context, value int64, fltrAttr attribute.Set) {
	b.deltaSum.measure(ctx, value, fltrAttr, nil)
}

type limitedBenchmarker struct {
	hcwg          hotColdWaitGroup
	hotColdValMap [2]sumValueMap[int64]
}

func (b *limitedBenchmarker) measure(ctx context.Context, value int64, fltrAttr attribute.Set) {
	hotIdx := b.hcwg.start()
	defer b.hcwg.done(hotIdx)
	b.hotColdValMap[hotIdx].measure(ctx, value, fltrAttr, nil)
}

func (b *limitedBenchmarker) pseudoCollect() {
	readIdx := b.hcwg.swapHotAndWait()
	b.hotColdValMap[readIdx].values.Range(func(_, _ any) bool { return true })
	b.hotColdValMap[readIdx].values.Clear()
}

func BenchmarkDeltaSum(b *testing.B) {
	tests := []struct {
		name string
		make func() sumBenchmarker
	}{
		{"limited", func() sumBenchmarker {
			return &limitedBenchmarker{
				hotColdValMap: [2]sumValueMap[int64]{
					{
						values: limitedSyncMap{aggLimit: 2000},
						newRes: func(attribute.Set) FilteredExemplarReservoir[int64] { return noopRes[int64]{} },
					},
					{
						values: limitedSyncMap{aggLimit: 2000},
						newRes: func(attribute.Set) FilteredExemplarReservoir[int64] { return noopRes[int64]{} },
					},
				},
			}
		}},
		{"lazy", func() sumBenchmarker {
			return &lazyBenchmarker{
				deltaSum: newDeltaSum(
					true,
					2000,
					func(attribute.Set) FilteredExemplarReservoir[int64] { return noopRes[int64]{} },
				),
			}
		}},
	}

	ctx := b.Context()
	attr := attribute.NewSet(attribute.String("key", "value"))

	for _, tt := range tests {
		b.Run(tt.name+"/MeasureNoCollect", func(b *testing.B) {
			m := tt.make()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.measure(ctx, 1, attr)
			}
		})

		b.Run(tt.name+"/MeasureWithCollect", func(b *testing.B) {
			m := tt.make()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.pseudoCollect()
				m.measure(ctx, 1, attr)
			}
		})

		b.Run(tt.name+"/OnlyCollect", func(b *testing.B) {
			m := tt.make()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.pseudoCollect()
			}
		})
	}
}
