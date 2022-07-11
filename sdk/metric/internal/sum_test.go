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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestSum(t *testing.T) {
	t.Run("Delta", func(t *testing.T) {
		t.Run("Int64", testAggregator(NewDeltaSum[int64](), deltaSumExpecter[int64]))
		t.Run("Float64", testAggregator(NewDeltaSum[float64](), deltaSumExpecter[float64]))
	})

	t.Run("Cumulative", func(t *testing.T) {
		t.Run("Int64", testAggregator(NewCumulativeSum[int64](), cumuSumExpecter[int64]))
		t.Run("Float64", testAggregator(NewCumulativeSum[float64](), cumuSumExpecter[float64]))
	})
}

func deltaSumExpecter[N int64 | float64](incr setMap[N]) func(m int) setMap[N] {
	expect := make(setMap[N], len(incr))
	return func(m int) setMap[N] {
		for actor, incr := range incr {
			expect[actor] = incr * N(m)
		}
		return expect
	}
}

func cumuSumExpecter[N int64 | float64](incr setMap[N]) func(m int) setMap[N] {
	var cycle int
	expect := make(setMap[N], len(incr))
	return func(m int) setMap[N] {
		cycle++
		for actor := range incr {
			expect[actor] = incr[actor] * N(cycle) * N(m)
		}
		return expect
	}
}

func testDeltaSumReset[N int64 | float64](a Aggregator[N]) func(*testing.T) {
	return func(t *testing.T) {
		expect := make(setMap[N])
		assertSetMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		a.Aggregate(1, alice)
		expect[alice] = 1
		assertSetMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		// The attr set should be forgotten once Aggregations is called.
		delete(expect, alice)
		assertSetMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		// Aggregating another set should not affect the original (alice).
		a.Aggregate(1, bob)
		expect[bob] = 1
		assertSetMap(t, expect, aggregationsToMap[N](a.Aggregations()))
	}
}

func TestDeltaSumReset(t *testing.T) {
	t.Run("Int64", testDeltaSumReset(NewDeltaSum[int64]()))
	t.Run("Float64", testDeltaSumReset(NewDeltaSum[float64]()))
}

var result []Aggregation

func benchmarkAggregatorN[N int64 | float64](b *testing.B, factory func() Aggregator[N], count int) {
	attrs := make([]attribute.Set, count)
	for i := range attrs {
		attrs[i] = attribute.NewSet(attribute.Int("value", i))
	}

	b.Run("Aggregate", func(b *testing.B) {
		agg := factory()
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			for _, attr := range attrs {
				agg.Aggregate(1, attr)
			}
		}
		assert.Len(b, agg.Aggregations(), count)
	})

	b.Run("Aggregations", func(b *testing.B) {
		aggs := make([]Aggregator[N], b.N)
		for n := range aggs {
			a := factory()
			for _, attr := range attrs {
				a.Aggregate(1, attr)
			}
			aggs[n] = a
		}

		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			result = aggs[n].Aggregations()
		}
	})
}

func benchmarkAggregator[N int64 | float64](factory func() Aggregator[N]) func(*testing.B) {
	counts := []int{1, 10, 100}
	return func(b *testing.B) {
		for _, n := range counts {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				benchmarkAggregatorN(b, factory, n)
			})
		}
	}
}

func BenchmarkSum(b *testing.B) {
	b.Run("Delta", func(b *testing.B) {
		b.Run("Int64", benchmarkAggregator(NewDeltaSum[int64]))
		b.Run("Float64", benchmarkAggregator(NewDeltaSum[float64]))
	})
	b.Run("Cumulative", func(b *testing.B) {
		b.Run("Int64", benchmarkAggregator(NewCumulativeSum[int64]))
		b.Run("Float64", benchmarkAggregator(NewCumulativeSum[float64]))
	})
}
