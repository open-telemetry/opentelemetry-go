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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

const (
	goroutines   = 5
	measurements = 30
	cycles       = 3
)

var (
	alice = attribute.NewSet(attribute.String("user", "alice"), attribute.Bool("admin", true))
	bob   = attribute.NewSet(attribute.String("user", "bob"), attribute.Bool("admin", false))
	carol = attribute.NewSet(attribute.String("user", "carol"), attribute.Bool("admin", false))
)

func TestSum(t *testing.T) {
	t.Run("Delta", func(t *testing.T) {
		t.Run("Int64", testSum(NewDeltaSum[int64](), deltaExpecter[int64]))
		t.Run("Float64", testSum(NewDeltaSum[float64](), deltaExpecter[float64]))
	})

	t.Run("Cumulative", func(t *testing.T) {
		t.Run("Int64", testSum(NewCumulativeSum[int64](), cumulativeExpecter[int64]))
		t.Run("Float64", testSum(NewCumulativeSum[float64](), cumulativeExpecter[float64]))
	})
}

// expectFunc returns a function that will return a map of expected values of
// a cycle. Each call advances the cycle.
type expectFunc[N int64 | float64] func(increments map[attribute.Set]N) func() map[attribute.Set]N

func testSum[N int64 | float64](a Aggregator[N], expecter expectFunc[N]) func(*testing.T) {
	increments := map[attribute.Set]N{alice: 1, bob: -1, carol: 2}
	f := expecter(increments)
	return func(t *testing.T) {
		for i := 0; i < cycles; i++ {
			var wg sync.WaitGroup
			wg.Add(goroutines)
			for i := 0; i < goroutines; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < measurements; j++ {
						for attrs, n := range increments {
							a.Aggregate(n, attrs)
						}
					}
				}()
			}
			wg.Wait()

			assertMap(t, f(), aggregationsToMap[N](a.Aggregations()))
		}
	}
}

func aggregationsToMap[N int64 | float64](a []Aggregation) map[attribute.Set]N {
	m := make(map[attribute.Set]N)
	for _, a := range a {
		m[a.Attributes] = a.Value.(SingleValue[N]).Value
	}
	return m
}

// assertMap asserts expected equals actual. The testify assert.Equal function
// does not give clear error messages for maps, this attempts to do so.
func assertMap[N int64 | float64](t *testing.T, expected, actual map[attribute.Set]N) {
	extra := make(map[attribute.Set]struct{})
	for attr := range actual {
		extra[attr] = struct{}{}
	}

	for attr, v := range expected {
		name := attr.Encoded(attribute.DefaultEncoder())
		t.Run(name, func(t *testing.T) {
			require.Contains(t, actual, attr)
			delete(extra, attr)
			assert.Equal(t, v, actual[attr])
		})
	}

	assert.Lenf(t, extra, 0, "unknown values added: %v", extra)
}

func deltaExpecter[N int64 | float64](incr map[attribute.Set]N) func() map[attribute.Set]N {
	expect := make(map[attribute.Set]N, len(incr))
	for actor, incr := range incr {
		expect[actor] = incr * measurements * goroutines
	}
	return func() map[attribute.Set]N { return expect }
}

func cumulativeExpecter[N int64 | float64](incr map[attribute.Set]N) func() map[attribute.Set]N {
	var cycle int
	base := make(map[attribute.Set]N, len(incr))
	for actor, incr := range incr {
		base[actor] = incr * measurements * goroutines
	}

	expect := make(map[attribute.Set]N, len(incr))
	return func() map[attribute.Set]N {
		cycle++
		for actor := range base {
			expect[actor] = base[actor] * N(cycle)
		}
		return expect
	}
}

func testDeltaSumReset[N int64 | float64](a Aggregator[N]) func(*testing.T) {
	return func(t *testing.T) {
		expect := make(map[attribute.Set]N)
		assertMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		a.Aggregate(1, alice)
		expect[alice] = 1
		assertMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		// The sum should be reset to zero once Aggregations is called.
		expect[alice] = 0
		assertMap(t, expect, aggregationsToMap[N](a.Aggregations()))

		// Aggregating another set should not affect the original (alice).
		a.Aggregate(1, bob)
		expect[bob] = 1
		assertMap(t, expect, aggregationsToMap[N](a.Aggregations()))
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
