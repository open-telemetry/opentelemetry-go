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

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func TestLastValue(t *testing.T) {
	t.Cleanup(mockTime(now))

	t.Run("Int64", testLastValue[int64]())
	t.Run("Float64", testLastValue[float64]())
}

func testLastValue[N int64 | float64]() func(*testing.T) {
	tester := &aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}

	eFunc := func(increments setMap) expectFunc {
		data := make([]metricdata.DataPoint[N], 0, len(increments))
		for a, v := range increments {
			point := metricdata.DataPoint[N]{Attributes: a, Time: now(), Value: N(v)}
			data = append(data, point)
		}
		gauge := metricdata.Gauge[N]{DataPoints: data}
		return func(int) metricdata.Aggregation { return gauge }
	}
	incr := monoIncr
	return tester.Run(NewLastValue[N](), incr, eFunc(incr))
}

func testLastValueReset[N int64 | float64](t *testing.T) {
	t.Cleanup(mockTime(now))

	a := NewLastValue[N]()
	expect := metricdata.Gauge[N]{}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())

	a.Aggregate(1, alice)
	expect.DataPoints = []metricdata.DataPoint[N]{{
		Attributes: alice,
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
		Time:       now(),
		Value:      1,
	}}
	metricdatatest.AssertAggregationsEqual(t, expect, a.Aggregation())
}

func TestLastValueReset(t *testing.T) {
	t.Run("Int64", testLastValueReset[int64])
	t.Run("Float64", testLastValueReset[float64])
}

func BenchmarkLastValue(b *testing.B) {
	b.Run("Int64", benchmarkAggregator(NewLastValue[int64]))
	b.Run("Float64", benchmarkAggregator(NewLastValue[float64]))
}
