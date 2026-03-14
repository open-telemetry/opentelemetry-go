// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestLastValue(t *testing.T) {
	c := new(clock)
	t.Cleanup(c.Register())

	t.Run("Int64/DeltaLastValue", testDeltaLastValue[int64]())
	c.Reset()
	t.Run("Float64/DeltaLastValue", testDeltaLastValue[float64]())
	c.Reset()

	t.Run("Int64/CumulativeLastValue", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		testCumulativeLastValue[int64](false)(t)
	})
	c.Reset()

	t.Run("Int64/CumulativeLastValue/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		testCumulativeLastValue[int64](true)(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeLastValue", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "false")
		testCumulativeLastValue[float64](false)(t)
	})
	c.Reset()

	t.Run("Float64/CumulativeLastValue/PerSeriesStartTimeEnabled", func(t *testing.T) {
		t.Setenv("OTEL_GO_X_PER_SERIES_START_TIMESTAMPS", "true")
		testCumulativeLastValue[float64](true)(t)
	})
	c.Reset()

	t.Run("Int64/DeltaPrecomputedLastValue", testDeltaPrecomputedLastValue[int64]())
	c.Reset()
	t.Run("Float64/DeltaPrecomputedLastValue", testDeltaPrecomputedLastValue[float64]())
	c.Reset()

	t.Run("Int64/CumulativePrecomputedLastValue", testCumulativePrecomputedLastValue[int64]())
	c.Reset()
	t.Run("Float64/CumulativePrecomputedLastValue", testCumulativePrecomputedLastValue[float64]())
}

func testDeltaLastValue[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.LastValue()
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			// Empty output if nothing is measured.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      2,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      -10,
						},
					},
				},
			},
		}, {
			// Everything resets, do not report old measurements.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      3,
						},
					},
				},
			},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Gauge[N]{
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
							Value:      1,
						},
					},
				},
			},
		},
	})
}

//nolint:revive // perSeriesStartTimeEnabled used to test a feature gate.
func testCumulativeLastValue[N int64 | float64](perSeriesStartTimeEnabled bool) func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.LastValue()

	aliceStartTime := y2kPlus(0)
	bobStartTime := y2kPlus(0)
	overflowStartTime := y2kPlus(0)

	timeStep1 := y2kPlus(2)
	timeStep2 := y2kPlus(3)
	timeStep3 := y2kPlus(4)
	timeStep4 := y2kPlus(5)

	if perSeriesStartTimeEnabled {
		aliceStartTime = y2kPlus(2)
		bobStartTime = y2kPlus(3)
		overflowStartTime = y2kPlus(7)

		timeStep1 = y2kPlus(4)
		timeStep2 = y2kPlus(5)
		timeStep3 = y2kPlus(6)
		timeStep4 = y2kPlus(8)
	}

	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			// Empty output if nothing is measured.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       timeStep1,
							Value:      2,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       timeStep1,
							Value:      -10,
						},
					},
				},
			},
		}, {
			// Cumulative temporality means no resets.
			input: []arg[N]{},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       timeStep2,
							Value:      2,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       timeStep2,
							Value:      -10,
						},
					},
				},
			},
		}, {
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       timeStep3,
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       timeStep3,
							Value:      3,
						},
					},
				},
			},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  aliceStartTime,
							Time:       timeStep4,
							Value:      1,
						},
						{
							Attributes: fltrBob,
							StartTime:  bobStartTime,
							Time:       timeStep4,
							Value:      1,
						},
						{
							Attributes: overflowSet,
							StartTime:  overflowStartTime,
							Time:       timeStep4,
							Value:      1,
						},
					},
				},
			},
		},
	})
}

func testDeltaPrecomputedLastValue[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedLastValue()
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			// Empty output if nothing is measured.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      2,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(1),
							Time:       y2kPlus(2),
							Value:      -10,
						},
					},
				},
			},
		}, {
			// Everything resets, do not report old measurements.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(3),
							Time:       y2kPlus(4),
							Value:      3,
						},
					},
				},
			},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Gauge[N]{
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
							Value:      1,
						},
					},
				},
			},
		},
	})
}

func testCumulativePrecomputedLastValue[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedLastValue()
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			// Empty output if nothing is measured.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(2),
							Value:      2,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(2),
							Value:      -10,
						},
					},
				},
			},
		}, {
			// Everything resets, do not report old measurements.
			input:  []arg[N]{},
			expect: output{n: 0, agg: metricdata.Gauge[N]{}},
		}, {
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			expect: output{
				n: 2,
				agg: metricdata.Gauge[N]{
					DataPoints: []metricdata.DataPoint[N]{
						{
							Attributes: fltrAlice,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(4),
							Value:      10,
						},
						{
							Attributes: fltrBob,
							StartTime:  y2kPlus(0),
							Time:       y2kPlus(4),
							Value:      3,
						},
					},
				},
			},
		}, {
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, 1, bob},
				// These will exceed cardinality limit.
				{ctx, 1, carol},
				{ctx, 1, dave},
			},
			expect: output{
				n: 3,
				agg: metricdata.Gauge[N]{
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
							Value:      1,
						},
					},
				},
			},
		},
	})
}

func TestLastValueConcurrentSafe(t *testing.T) {
	t.Run("Int64/DeltaLastValue", testDeltaLastValueConcurrentSafe[int64]())
	t.Run("Float64/DeltaLastValue", testDeltaLastValueConcurrentSafe[float64]())
	t.Run("Int64/CumulativeLastValue", testCumulativeLastValueConcurrentSafe[int64]())
	t.Run("Float64/CumulativeLastValue", testCumulativeLastValueConcurrentSafe[float64]())
	t.Run("Int64/DeltaPrecomputedLastValue", testDeltaPrecomputedLastValueConcurrentSafe[int64]())
	t.Run("Float64/DeltaPrecomputedLastValue", testDeltaPrecomputedLastValueConcurrentSafe[float64]())
	t.Run("Int64/CumulativePrecomputedLastValue", testCumulativePrecomputedLastValueConcurrentSafe[int64]())
	t.Run("Float64/CumulativePrecomputedLastValue", testCumulativePrecomputedLastValueConcurrentSafe[float64]())
}

func validateGauge[N int64 | float64](t *testing.T, aggs []metricdata.Aggregation) {
	// A gauge takes the *last* recorded value.
	// During high concurrency, reading the Gauge can snap any value in the
	// iteration cycle of the corresponding Goroutines.
	valid := make(map[N]bool)
	for _, v := range getConcurrentVals[N]() {
		valid[v] = true
	}
	// TODO(dashpole): Fix a concurrency bug where a gauge can be collected with
	// the value zero even when no zero-value measurements have been recorded.
	valid[0] = true

	for _, agg := range aggs {
		s, ok := agg.(metricdata.Gauge[N])
		require.True(t, ok)
		require.LessOrEqual(t, len(s.DataPoints), 3, "AggregationLimit of 3 exceeded")
		for _, dp := range s.DataPoints {
			assert.True(t, valid[dp.Value], "Unexpected gauge value: %v", dp.Value)
		}
	}
}

func testCumulativeLastValueConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.LastValue()
	return testAggregationConcurrentSafe[N](in, out, validateGauge[N])
}

func testDeltaLastValueConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.LastValue()
	return testAggregationConcurrentSafe[N](in, out, validateGauge[N])
}

func testDeltaPrecomputedLastValueConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.DeltaTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedLastValue()
	return testAggregationConcurrentSafe[N](in, out, validateGauge[N])
}

func testCumulativePrecomputedLastValueConcurrentSafe[N int64 | float64]() func(*testing.T) {
	in, out := Builder[N]{
		Temporality:      metricdata.CumulativeTemporality,
		Filter:           attrFltr,
		AggregationLimit: 3,
	}.PrecomputedLastValue()
	return testAggregationConcurrentSafe[N](in, out, validateGauge[N])
}

func BenchmarkLastValue(b *testing.B) {
	b.Run("Int64", benchmarkAggregate(Builder[int64]{}.PrecomputedLastValue))
	b.Run("Float64", benchmarkAggregate(Builder[float64]{}.PrecomputedLastValue))
}
