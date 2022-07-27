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

//go:build go1.18
// +build go1.18

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
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

	t.Run("Delta", func(t *testing.T) {
		incr := monoIncr
		eFunc := deltaExpecter[N](incr, true)
		t.Run("Monotonic", tester.Run(NewMonotonicDeltaSum[N](), incr, eFunc))

		incr = nonMonoIncr
		eFunc = deltaExpecter[N](incr, false)
		t.Run("NonMonotonic", tester.Run(NewNonMonotonicDeltaSum[N](), incr, eFunc))
	})

	t.Run("Cumulative", func(t *testing.T) {
		incr := monoIncr
		eFunc := cumuExpecter[N](incr, true)
		t.Run("Monotonic", tester.Run(NewMonotonicCumulativeSum[N](), incr, eFunc))

		incr = nonMonoIncr
		eFunc = cumuExpecter[N](incr, false)
		t.Run("NonMonotonic", tester.Run(NewNonMonotonicCumulativeSum[N](), incr, eFunc))
	})
}

func deltaExpecter[N int64 | float64](incr setMap, mono bool) expectFunc {
	sum := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality, IsMonotonic: mono}
	return func(m int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point[N](a, N(v*m)))
		}
		return sum
	}
}

func cumuExpecter[N int64 | float64](incr setMap, mono bool) expectFunc {
	var cycle int
	sum := metricdata.Sum[N]{Temporality: metricdata.CumulativeTemporality, IsMonotonic: mono}
	return func(m int) metricdata.Aggregation {
		cycle++
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point[N](a, N(v*cycle*m)))
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

	f := func(expect metricdata.Sum[N], a Aggregator[N]) func(*testing.T) {
		return func(t *testing.T) {
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			a.Aggregate(1, alice)
			expect.DataPoints = []metricdata.DataPoint[N]{point[N](alice, 1)}
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			// The attr set should be forgotten once Aggregations is called.
			expect.DataPoints = nil
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			// Aggregating another set should not affect the original (alice).
			a.Aggregate(1, bob)
			expect.DataPoints = []metricdata.DataPoint[N]{point[N](bob, 1)}
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())
		}
	}

	sum := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality}
	t.Run("NonMonotonic", f(sum, NewNonMonotonicDeltaSum[N]()))

	sum.IsMonotonic = true
	t.Run("Monotonic", f(sum, NewMonotonicDeltaSum[N]()))
}

func TestDeltaSumReset(t *testing.T) {
	t.Run("Int64", testDeltaSumReset[int64])
	t.Run("Float64", testDeltaSumReset[float64])
}

func testMonotonicError[N int64 | float64](t *testing.T) {
	f := func(a Aggregator[N]) func(t *testing.T) {
		var err error
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(e error) { err = e }))
		a.Aggregate(-1, alice) // Should error.
		return func(t *testing.T) { assert.ErrorIs(t, err, errNegVal) }
	}
	t.Run("Delta", f(NewMonotonicDeltaSum[N]()))
	t.Run("Cumulative", f(NewMonotonicCumulativeSum[N]()))
}

func TestMonotonicError(t *testing.T) {
	t.Run("Int64", testMonotonicError[int64])
	t.Run("Float64", testMonotonicError[float64])
}

func BenchmarkSum(b *testing.B) {
	b.Run("Int64", benchmarkSum[int64])
	b.Run("Float64", benchmarkSum[float64])
}

func benchmarkSum[N int64 | float64](b *testing.B) {
	b.Run("Delta", func(b *testing.B) {
		b.Run("Monotonic", benchmarkAggregator(NewMonotonicDeltaSum[int64]))
		b.Run("NonMonotonic", benchmarkAggregator(NewNonMonotonicDeltaSum[int64]))
	})
	b.Run("Cumulative", func(b *testing.B) {
		b.Run("Monotonic", benchmarkAggregator(NewMonotonicCumulativeSum[int64]))
		b.Run("NonMonotonic", benchmarkAggregator(NewNonMonotonicCumulativeSum[int64]))
	})
}
