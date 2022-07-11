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

import "testing"

func TestLastValue(t *testing.T) {
	t.Run("Int64", testAggregator(NewLastValue[int64](), lastValueExpecter[int64]))
	t.Run("Float64", testAggregator(NewLastValue[float64](), lastValueExpecter[float64]))
}

func lastValueExpecter[N int64 | float64](incr setMap[N]) func(int) setMap[N] {
	expect := make(setMap[N], len(incr))
	for actor, incr := range incr {
		expect[actor] = incr
	}
	return func(int) setMap[N] {
		return expect
	}
}

func testLastValueReset[N int64 | float64](a Aggregator[N]) func(*testing.T) {
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

func TestLastValueReset(t *testing.T) {
	t.Run("Int64", testLastValueReset(NewLastValue[int64]()))
	t.Run("Float64", testLastValueReset(NewLastValue[float64]()))
}
