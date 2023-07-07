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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestSum(t *testing.T) {
	t.Cleanup(mockTime(now))
	t.Run("Int64", testSum[int64])
	t.Run("Float64", testSum[float64])
}

func testSum[N int64 | float64](t *testing.T) {
	tester := &aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}
	totalMeasurements := defaultGoroutines * defaultMeasurements

	t.Run("Delta", func(t *testing.T) {
		incr, mono := monoIncr[N](), true
		eFunc := deltaExpecter[N](incr, mono)
		t.Run("Monotonic", tester.Run(newDeltaSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr[N](), false
		eFunc = deltaExpecter[N](incr, mono)
		t.Run("NonMonotonic", tester.Run(newDeltaSum[N](mono), incr, eFunc))
	})

	t.Run("Cumulative", func(t *testing.T) {
		incr, mono := monoIncr[N](), true
		eFunc := cumuExpecter[N](incr, mono)
		t.Run("Monotonic", tester.Run(newCumulativeSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr[N](), false
		eFunc = cumuExpecter[N](incr, mono)
		t.Run("NonMonotonic", tester.Run(newCumulativeSum[N](mono), incr, eFunc))
	})

	t.Run("PreComputedDelta", func(t *testing.T) {
		incr, mono := monoIncr[N](), true
		eFunc := preDeltaExpecter[N](incr, mono, N(totalMeasurements))
		t.Run("Monotonic", tester.Run(newPrecomputedDeltaSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr[N](), false
		eFunc = preDeltaExpecter[N](incr, mono, N(totalMeasurements))
		t.Run("NonMonotonic", tester.Run(newPrecomputedDeltaSum[N](mono), incr, eFunc))
	})

	t.Run("PreComputedCumulative", func(t *testing.T) {
		incr, mono := monoIncr[N](), true
		eFunc := preCumuExpecter[N](incr, mono, N(totalMeasurements))
		t.Run("Monotonic", tester.Run(newPrecomputedCumulativeSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr[N](), false
		eFunc = preCumuExpecter[N](incr, mono, N(totalMeasurements))
		t.Run("NonMonotonic", tester.Run(newPrecomputedCumulativeSum[N](mono), incr, eFunc))
	})
}

func deltaExpecter[N int64 | float64](incr setMap[N], mono bool) expectFunc {
	sum := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality, IsMonotonic: mono}
	return func(m int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point(a, v*N(m)))
		}
		return sum
	}
}

func cumuExpecter[N int64 | float64](incr setMap[N], mono bool) expectFunc {
	var cycle N
	sum := metricdata.Sum[N]{Temporality: metricdata.CumulativeTemporality, IsMonotonic: mono}
	return func(m int) metricdata.Aggregation {
		cycle++
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point(a, v*cycle*N(m)))
		}
		return sum
	}
}

func preDeltaExpecter[N int64 | float64](incr setMap[N], mono bool, totalMeasurements N) expectFunc {
	sum := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality, IsMonotonic: mono}
	last := make(map[attribute.Set]N)
	return func(int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			l := last[a]
			sum.DataPoints = append(sum.DataPoints, point(a, totalMeasurements*(N(v)-l)))
			last[a] = N(v)
		}
		return sum
	}
}

func preCumuExpecter[N int64 | float64](incr setMap[N], mono bool, totalMeasurements N) expectFunc {
	sum := metricdata.Sum[N]{Temporality: metricdata.CumulativeTemporality, IsMonotonic: mono}
	return func(int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point(a, totalMeasurements*N(v)))
		}
		return sum
	}
}

// point returns a DataPoint that started and ended now.
func point[N int64 | float64](a attribute.Set, v N) metricdata.DataPoint[N] {
	return metricdata.DataPoint[N]{
		Attributes: a,
		StartTime:  now(),
		Time:       now(),
		Value:      N(v),
	}
}

func testDeltaSumReset[N int64 | float64](t *testing.T) {
	t.Cleanup(mockTime(now))

	a := newDeltaSum[N](false)
	assert.Nil(t, a.Aggregation())

	a.Aggregate(1, alice)
	expect := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality}
	expect.DataPoints = []metricdata.DataPoint[N]{point[N](alice, 1)}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

	// The attr set should be forgotten once Aggregations is called.
	expect.DataPoints = nil
	assert.Nil(t, a.Aggregation())

	// Aggregating another set should not affect the original (alice).
	a.Aggregate(1, bob)
	expect.DataPoints = []metricdata.DataPoint[N]{point[N](bob, 1)}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())
}

func TestDeltaSumReset(t *testing.T) {
	t.Run("Int64", testDeltaSumReset[int64])
	t.Run("Float64", testDeltaSumReset[float64])
}

func TestPreComputedDeltaSum(t *testing.T) {
	var mono bool
	agg := newPrecomputedDeltaSum[int64](mono)
	require.Implements(t, (*aggregator[int64])(nil), agg)

	attrs := attribute.NewSet(attribute.String("key", "val"))
	agg.Aggregate(1, attrs)
	want := metricdata.Sum[int64]{
		IsMonotonic: mono,
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.DataPoint[int64]{point[int64](attrs, 1)},
	}
	opt := metricdatatest.IgnoreTimestamp()
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// No observation results in an empty aggregation, and causes previous
	// observations to be forgotten.
	metricdatatest.AssertAggregationsEqual(t, nil, agg.Aggregation(), opt)

	agg.Aggregate(1, attrs)
	// measured(+): 1, previous(-): 0
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 1)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Duplicate observations add
	agg.Aggregate(2, attrs)
	agg.Aggregate(5, attrs)
	agg.Aggregate(3, attrs)
	agg.Aggregate(10, attrs)
	// measured(+): 20, previous(-): 1
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 19)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Filtered values should not persist.
	agg.Aggregate(5, attrs)
	// measured(+): 5, previous(-): 20
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, -15)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Order should not affect measure.
	// Filtered should add.
	agg.Aggregate(3, attrs)
	agg.Aggregate(7, attrs)
	agg.Aggregate(10, attrs)
	// measured(+): 20, previous(-): 5
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 15)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)
	agg.Aggregate(7, attrs)
	// measured(+): 7, previous(-): 20
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, -13)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)
}

func TestPreComputedCumulativeSum(t *testing.T) {
	var mono bool
	agg := newPrecomputedCumulativeSum[int64](mono)
	require.Implements(t, (*aggregator[int64])(nil), agg)

	attrs := attribute.NewSet(attribute.String("key", "val"))
	agg.Aggregate(1, attrs)
	want := metricdata.Sum[int64]{
		IsMonotonic: mono,
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.DataPoint[int64]{point[int64](attrs, 1)},
	}
	opt := metricdatatest.IgnoreTimestamp()
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Cumulative values should not persist.
	metricdatatest.AssertAggregationsEqual(t, nil, agg.Aggregation(), opt)

	agg.Aggregate(1, attrs)
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 1)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Duplicate measurements add
	agg.Aggregate(5, attrs)
	agg.Aggregate(3, attrs)
	agg.Aggregate(10, attrs)
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 18)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)

	// Order should not affect measure.
	// Duplicates should add.
	agg.Aggregate(3, attrs)
	agg.Aggregate(7, attrs)
	agg.Aggregate(10, attrs)
	want.DataPoints = []metricdata.DataPoint[int64]{point[int64](attrs, 20)}
	metricdatatest.AssertAggregationsEqual(t, want, agg.Aggregation(), opt)
}

func TestEmptySumNilAggregation(t *testing.T) {
	assert.Nil(t, newCumulativeSum[int64](true).Aggregation())
	assert.Nil(t, newCumulativeSum[int64](false).Aggregation())
	assert.Nil(t, newCumulativeSum[float64](true).Aggregation())
	assert.Nil(t, newCumulativeSum[float64](false).Aggregation())
	assert.Nil(t, newDeltaSum[int64](true).Aggregation())
	assert.Nil(t, newDeltaSum[int64](false).Aggregation())
	assert.Nil(t, newDeltaSum[float64](true).Aggregation())
	assert.Nil(t, newDeltaSum[float64](false).Aggregation())
	assert.Nil(t, newPrecomputedCumulativeSum[int64](true).Aggregation())
	assert.Nil(t, newPrecomputedCumulativeSum[int64](false).Aggregation())
	assert.Nil(t, newPrecomputedCumulativeSum[float64](true).Aggregation())
	assert.Nil(t, newPrecomputedCumulativeSum[float64](false).Aggregation())
	assert.Nil(t, newPrecomputedDeltaSum[int64](true).Aggregation())
	assert.Nil(t, newPrecomputedDeltaSum[int64](false).Aggregation())
	assert.Nil(t, newPrecomputedDeltaSum[float64](true).Aggregation())
	assert.Nil(t, newPrecomputedDeltaSum[float64](false).Aggregation())
}

func BenchmarkSum(b *testing.B) {
	b.Run("Int64", benchmarkSum[int64])
	b.Run("Float64", benchmarkSum[float64])
}

func benchmarkSum[N int64 | float64](b *testing.B) {
	// The monotonic argument is only used to annotate the Sum returned from
	// the Aggregation method. It should not have an effect on operational
	// performance, therefore, only monotonic=false is benchmarked here.
	factory := func() aggregator[N] { return newDeltaSum[N](false) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() aggregator[N] { return newCumulativeSum[N](false) }
	b.Run("Cumulative", benchmarkAggregator(factory))
}
