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
	defaultGoroutines   = 5
	defaultMeasurements = 30
	defaultCycles       = 3
)

var (
	alice = attribute.NewSet(attribute.String("user", "alice"), attribute.Bool("admin", true))
	bob   = attribute.NewSet(attribute.String("user", "bob"), attribute.Bool("admin", false))
	carol = attribute.NewSet(attribute.String("user", "carol"), attribute.Bool("admin", false))
)

// setMap maps atribute sets to a number.
type setMap[N int64 | float64] map[attribute.Set]N

// expectFunc returns a function that will return an setMap of expected
// values of a cycle that contains m measurements (total across all
// goroutines). Each call advances the cycle.
type expectFunc[N int64 | float64] func(increments setMap[N]) func(m int) setMap[N]

// testAggregator tests aggregator a produces the expecter defined values
// using an aggregatorTester.
func testAggregator[N int64 | float64](a Aggregator[N], expecter expectFunc[N]) func(*testing.T) {
	return (&aggregatorTester[N]{
		GoroutineN:   defaultGoroutines,
		MeasurementN: defaultMeasurements,
		CycleN:       defaultCycles,
	}).Run(a, expecter)
}

// aggregatorTester runs an acceptance test on an Aggregator. It will ask an
// Aggregator to aggregate a set of values as if they were real measurements
// made MeasurementN number of times. This will be done in GoroutineN number
// of different goroutines. After the Aggregator has been asked to aggregate
// all these measurements, it is validated using a passed expecterFunc. This
// set of operation is a signle cycle, and the the aggregatorTester will run
// CycleN number of cycles.
type aggregatorTester[N int64 | float64] struct {
	// GoroutineN is the number of goroutines aggregatorTester will use to run
	// the test with.
	GoroutineN int
	// MeasurementN is the number of measurements that are made each cycle a
	// goroutine runs the test.
	MeasurementN int
	// CycleN is the number of times a goroutine will make a set of
	// measurements.
	CycleN int
}

func (at *aggregatorTester[N]) Run(a Aggregator[N], expecter expectFunc[N]) func(*testing.T) {
	increments := map[attribute.Set]N{alice: 1, bob: -1, carol: 2}
	f := expecter(increments)
	m := at.MeasurementN * at.GoroutineN
	return func(t *testing.T) {
		for i := 0; i < at.CycleN; i++ {
			var wg sync.WaitGroup
			wg.Add(at.GoroutineN)
			for i := 0; i < at.GoroutineN; i++ {
				go func() {
					defer wg.Done()
					for j := 0; j < at.MeasurementN; j++ {
						for attrs, n := range increments {
							a.Aggregate(n, attrs)
						}
					}
				}()
			}
			wg.Wait()

			assertSetMap(t, f(m), aggregationsToMap[N](a.Aggregations()))
		}
	}
}

func aggregationsToMap[N int64 | float64](a []Aggregation) setMap[N] {
	m := make(setMap[N])
	for _, a := range a {
		m[a.Attributes] = a.Value.(SingleValue[N]).Value
	}
	return m
}

// assertSetMap asserts expected equals actual. The testify assert.Equal
// function does not give clear error messages for maps, this attempts to do
// so.
func assertSetMap[N int64 | float64](t *testing.T, expected, actual setMap[N]) {
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
