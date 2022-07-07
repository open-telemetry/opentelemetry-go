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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

const (
	goroutines   = 5
	measurements = 30
)

var (
	alice = attribute.NewSet(attribute.String("user", "alice"), attribute.Bool("admin", true))
	bob   = attribute.NewSet(attribute.String("user", "bob"), attribute.Bool("admin", false))
	carol = attribute.NewSet(attribute.String("user", "carol"), attribute.Bool("admin", false))
)

// apply aggregates all the incr values with agg.
func apply[N int64 | float64](incr map[attribute.Set]N, agg Aggregator[N]) {
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < measurements; j++ {
				for attrs, n := range incr {
					agg.Aggregate(n, attrs)
				}
			}
		}()
	}
	wg.Wait()
}

func check[N int64 | float64](t *testing.T, expected map[attribute.Set]N, actual []Aggregation) {
	extra := make(map[attribute.Set]struct{})
	// Convert []Aggregation to map[attribute.Set]N
	aMap := make(map[attribute.Set]N)
	for _, a := range actual {
		aMap[a.Attributes] = a.Value.(SingleValue[N]).Value
		extra[a.Attributes] = struct{}{}
	}

	for attr, v := range expected {
		name := attr.Encoded(attribute.DefaultEncoder())
		t.Run(name, func(t *testing.T) {
			require.Contains(t, aMap, attr)
			delete(extra, attr)
			assert.Equal(t, v, aMap[attr])
		})
	}

	assert.Lenf(t, extra, 0, "unknown values added: %v", extra)
}

func testDeltaSum[N int64 | float64](t *testing.T, agg Aggregator[N]) {
	increments := map[attribute.Set]N{alice: 1, bob: -1, carol: 2}
	apply(increments, agg)

	want := make(map[attribute.Set]N, len(increments))
	for actor, incr := range increments {
		want[actor] = incr * measurements * goroutines
	}
	check(t, want, agg.Aggregations())

	require.IsType(t, &deltaSum[N]{}, agg)
	ds := agg.(*deltaSum[N])
	assert.Len(t, ds.values, 0)

	apply(increments, agg)
	// Delta sums are expected to reset after each call to Aggregations.
	check(t, want, agg.Aggregations())
}

func testCumulativeSum[N int64 | float64](t *testing.T, agg Aggregator[N]) {
	increments := map[attribute.Set]N{alice: 1, bob: -1, carol: 2}
	apply(increments, agg)

	want := make(map[attribute.Set]N, len(increments))
	for actor, incr := range increments {
		want[actor] = incr * measurements * goroutines
	}
	check(t, want, agg.Aggregations())

	require.IsType(t, &cumulativeSum[N]{}, agg)
	ds := agg.(*cumulativeSum[N])
	assert.Len(t, ds.values, len(increments))

	apply(increments, agg)
	// Cumulative sums maintain state, this should double the value.
	for actor := range want {
		want[actor] += want[actor]
	}
	check(t, want, agg.Aggregations())
}

func TestInt64DeltaSum(t *testing.T) {
	testDeltaSum(t, NewDeltaSum[int64]())
}

func TestFloat64DeltaSum(t *testing.T) {
	testDeltaSum(t, NewDeltaSum[float64]())
}

func TestInt64CumulativeSum(t *testing.T) {
	testCumulativeSum(t, NewCumulativeSum[int64]())
}

func TestFloat64CumulativeSum(t *testing.T) {
	testCumulativeSum(t, NewCumulativeSum[float64]())
}
