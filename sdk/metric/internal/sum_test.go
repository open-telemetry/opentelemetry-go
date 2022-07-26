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
	"fmt"
	"testing"

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
	expecter := sumExpecterFactory[N]{}

	t.Run("Delta", func(t *testing.T) {
		temp := metricdata.DeltaTemporality
		t.Run("Monotonic", func(t *testing.T) {
			incr := monoIncr
			eFunc := expecter.ExpecterFunc(incr, temp, true)
			tester.Run(NewMonotonicDeltaSum[N](), incr, eFunc)
		})
		t.Run("NonMonotonic", func(t *testing.T) {
			incr := nonMonoIncr
			eFunc := expecter.ExpecterFunc(incr, temp, false)
			tester.Run(NewNonMonotonicDeltaSum[N](), incr, eFunc)
		})
	})

	t.Run("Cumulative", func(t *testing.T) {
		temp := metricdata.CumulativeTemporality
		t.Run("Monotonic", func(t *testing.T) {
			incr := monoIncr
			eFunc := expecter.ExpecterFunc(incr, temp, true)
			tester.Run(NewMonotonicCumulativeSum[N](), incr, eFunc)
		})
		t.Run("NonMonotonic", func(t *testing.T) {
			incr := nonMonoIncr
			eFunc := expecter.ExpecterFunc(incr, temp, false)
			tester.Run(NewNonMonotonicCumulativeSum[N](), incr, eFunc)
		})
	})
}

type sumExpecterFactory[N int64 | float64] struct{}

func (s *sumExpecterFactory[N]) ExpecterFunc(increments setMap, t metricdata.Temporality, monotonic bool) expectFunc {
	sum := metricdata.Sum[N]{
		Temporality: t,
		IsMonotonic: monotonic,
	}

	switch t {
	case metricdata.DeltaTemporality:
		return func(m int) metricdata.Aggregation {
			sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(increments))
			for actor, incr := range increments {
				sum.DataPoints = append(sum.DataPoints, metricdata.DataPoint[N]{
					Attributes: actor,
					StartTime:  now(),
					Time:       now(),
					Value:      N(incr * m),
				})
			}
			return sum
		}
	case metricdata.CumulativeTemporality:
		var cycle int
		return func(m int) metricdata.Aggregation {
			cycle++
			sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(increments))
			for actor, incr := range increments {
				sum.DataPoints = append(sum.DataPoints, metricdata.DataPoint[N]{
					Attributes: actor,
					StartTime:  now(),
					Time:       now(),
					Value:      N(incr * cycle * m),
				})
			}
			return sum
		}
	default:
		panic(fmt.Sprintf("unsupported temporality: %v", t))
	}
}

func testDeltaSumReset[N int64 | float64](t *testing.T) {
	t.Cleanup(mockTime(now))

	f := func(expect metricdata.Sum[N], a Aggregator[N]) func(*testing.T) {
		return func(t *testing.T) {
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			a.Aggregate(1, alice)
			expect.DataPoints = []metricdata.DataPoint[N]{{
				Attributes: alice,
				StartTime:  now(),
				Time:       now(),
				Value:      1,
			}}
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			// The attr set should be forgotten once Aggregations is called.
			expect.DataPoints = nil
			metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

			// Aggregating another set should not affect the original (alice).
			a.Aggregate(1, bob)
			expect.DataPoints = []metricdata.DataPoint[N]{{
				Attributes: bob,
				StartTime:  now(),
				Time:       now(),
				Value:      1,
			}}
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
