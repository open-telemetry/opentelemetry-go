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

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"testing"

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
		incr, mono := monoIncr, true
		eFunc := deltaExpecter[N](incr, mono)
		t.Run("Monotonic", tester.Run(NewDeltaSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr, false
		eFunc = deltaExpecter[N](incr, mono)
		t.Run("NonMonotonic", tester.Run(NewDeltaSum[N](mono), incr, eFunc))
	})

	t.Run("Cumulative", func(t *testing.T) {
		incr, mono := monoIncr, true
		eFunc := cumuExpecter[N](incr, mono)
		t.Run("Monotonic", tester.Run(NewCumulativeSum[N](mono), incr, eFunc))

		incr, mono = nonMonoIncr, false
		eFunc = cumuExpecter[N](incr, mono)
		t.Run("NonMonotonic", tester.Run(NewCumulativeSum[N](mono), incr, eFunc))
	})

	t.Run("PreComputed", func(t *testing.T) {
		incr, mono, temp := monoIncr, true, metricdata.DeltaTemporality
		eFunc := preExpecter[N](incr, mono, temp)
		t.Run("Monotonic/Delta", tester.Run(NewPrecomputedSum[N](mono, temp), incr, eFunc))

		temp = metricdata.CumulativeTemporality
		eFunc = preExpecter[N](incr, mono, temp)
		t.Run("Monotonic/Cumulative", tester.Run(NewPrecomputedSum[N](mono, temp), incr, eFunc))

		incr, mono, temp = nonMonoIncr, false, metricdata.DeltaTemporality
		eFunc = preExpecter[N](incr, mono, temp)
		t.Run("NonMonotonic/Delta", tester.Run(NewPrecomputedSum[N](mono, temp), incr, eFunc))

		temp = metricdata.CumulativeTemporality
		eFunc = preExpecter[N](incr, mono, temp)
		t.Run("NonMonotonic/Cumulative", tester.Run(NewPrecomputedSum[N](mono, temp), incr, eFunc))
	})
}

func deltaExpecter[N int64 | float64](incr setMap, mono bool) expectFunc {
	sum := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality, IsMonotonic: mono}
	return func(m int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point(a, N(v*m)))
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
			sum.DataPoints = append(sum.DataPoints, point(a, N(v*cycle*m)))
		}
		return sum
	}
}

func preExpecter[N int64 | float64](incr setMap, mono bool, temp metricdata.Temporality) expectFunc {
	sum := metricdata.Sum[N]{Temporality: temp, IsMonotonic: mono}
	return func(int) metricdata.Aggregation {
		sum.DataPoints = make([]metricdata.DataPoint[N], 0, len(incr))
		for a, v := range incr {
			sum.DataPoints = append(sum.DataPoints, point(a, N(v)))
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

	expect := metricdata.Sum[N]{Temporality: metricdata.DeltaTemporality}
	a := NewDeltaSum[N](false)
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

func TestDeltaSumReset(t *testing.T) {
	t.Run("Int64", testDeltaSumReset[int64])
	t.Run("Float64", testDeltaSumReset[float64])
}

func BenchmarkSum(b *testing.B) {
	b.Run("Int64", benchmarkSum[int64])
	b.Run("Float64", benchmarkSum[float64])
}

func benchmarkSum[N int64 | float64](b *testing.B) {
	// The monotonic argument is only used to annotate the Sum returned from
	// the Aggregation method. It should not have an effect on operational
	// performance, therefore, only monotonic=false is benchmarked here.
	factory := func() Aggregator[N] { return NewDeltaSum[N](false) }
	b.Run("Delta", benchmarkAggregator(factory))
	factory = func() Aggregator[N] { return NewCumulativeSum[N](false) }
	b.Run("Cumulative", benchmarkAggregator(factory))
}
