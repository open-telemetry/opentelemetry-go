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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// This is an aggregator that has a stable output, used for testing. It does not
// follow any spec prescribed aggregation.
type testStableAggregator[N int64 | float64] struct {
	sync.Mutex
	values []metricdata.DataPoint[N]
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (a *testStableAggregator[N]) Aggregate(measurement N, attr attribute.Set) {
	a.Lock()
	defer a.Unlock()

	a.values = append(a.values, metricdata.DataPoint[N]{
		Attributes: attr,
		Value:      measurement,
	})
}

// Aggregation returns an Aggregation, for all the aggregated
// measurements made and ends an aggregation cycle.
func (a *testStableAggregator[N]) Aggregation() metricdata.Aggregation {
	return metricdata.Gauge[N]{
		DataPoints: a.values,
	}
}

func testNewFilterNoFilter[N int64 | float64](t *testing.T, agg Aggregator[N]) {
	filter := NewFilter(agg, nil)
	assert.Equal(t, agg, filter)
}

func testNewFilter[N int64 | float64](t *testing.T, agg Aggregator[N]) {
	f := NewFilter(agg, testAttributeFilter)
	require.IsType(t, &filter[N]{}, f)
	filt := f.(*filter[N])
	assert.Equal(t, agg, filt.aggregator)
}

func testAttributeFilter(input attribute.Set) attribute.Set {
	out, _ := input.Filter(func(kv attribute.KeyValue) bool {
		return kv.Key == "power-level"
	})
	return out
}

func TestNewFilter(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		agg := &testStableAggregator[int64]{}
		testNewFilterNoFilter[int64](t, agg)
		testNewFilter[int64](t, agg)
	})
	t.Run("float64", func(t *testing.T) {
		agg := &testStableAggregator[float64]{}
		testNewFilterNoFilter[float64](t, agg)
		testNewFilter[float64](t, agg)
	})
}

func testDataPoint[N int64 | float64](attr attribute.Set) metricdata.DataPoint[N] {
	return metricdata.DataPoint[N]{
		Attributes: attr,
		Value:      1,
	}
}

func testFilterAggregate[N int64 | float64](t *testing.T) {
	testCases := []struct {
		name      string
		inputAttr []attribute.Set
		output    []metricdata.DataPoint[N]
	}{
		{
			name: "Will filter all out",
			inputAttr: []attribute.Set{
				attribute.NewSet(
					attribute.String("foo", "bar"),
					attribute.Float64("lifeUniverseEverything", 42.0),
				),
			},
			output: []metricdata.DataPoint[N]{
				testDataPoint[N](*attribute.EmptySet()),
			},
		},
		{
			name: "Will keep appropriate attributes",
			inputAttr: []attribute.Set{
				attribute.NewSet(
					attribute.String("foo", "bar"),
					attribute.Int("power-level", 9001),
					attribute.Float64("lifeUniverseEverything", 42.0),
				),
				attribute.NewSet(
					attribute.String("foo", "bar"),
					attribute.Int("power-level", 9001),
				),
			},
			output: []metricdata.DataPoint[N]{
				// A real Aggregator will combine these, the testAggregator doesn't for list stability.
				testDataPoint[N](attribute.NewSet(attribute.Int("power-level", 9001))),
				testDataPoint[N](attribute.NewSet(attribute.Int("power-level", 9001))),
			},
		},
		{
			name: "Will combine Aggregations",
			inputAttr: []attribute.Set{
				attribute.NewSet(
					attribute.String("foo", "bar"),
				),
				attribute.NewSet(
					attribute.Float64("lifeUniverseEverything", 42.0),
				),
			},
			output: []metricdata.DataPoint[N]{
				// A real Aggregator will combine these, the testAggregator doesn't for list stability.
				testDataPoint[N](*attribute.EmptySet()),
				testDataPoint[N](*attribute.EmptySet()),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFilter[N](&testStableAggregator[N]{}, testAttributeFilter)
			for _, set := range tt.inputAttr {
				f.Aggregate(1, set)
			}
			out := f.Aggregation().(metricdata.Gauge[N])
			assert.Equal(t, tt.output, out.DataPoints)
		})
	}
}

func TestFilterAggregate(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		testFilterAggregate[int64](t)
	})
	t.Run("float64", func(t *testing.T) {
		testFilterAggregate[float64](t)
	})
}

func testFilterConcurrent[N int64 | float64](t *testing.T) {
	f := NewFilter[N](&testStableAggregator[N]{}, testAttributeFilter)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		f.Aggregate(1, attribute.NewSet(
			attribute.String("foo", "bar"),
		))
		wg.Done()
	}()

	go func() {
		f.Aggregate(1, attribute.NewSet(
			attribute.Int("power-level", 9001),
		))
		wg.Done()
	}()

	wg.Wait()
}

func TestFilterConcurrent(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		testFilterConcurrent[int64](t)
	})
	t.Run("float64", func(t *testing.T) {
		testFilterConcurrent[float64](t)
	})
}
