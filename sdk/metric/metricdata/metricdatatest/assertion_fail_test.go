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

//go:build tests_fail
// +build tests_fail

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// These tests are used to develop the failure messages of this package's
// assertions. They can be run with the following.
//
//   go test -tags tests_fail ./...

func testFailDatatype[T Datatypes](a, b T) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, b)
	}
}

func TestFailAssertEqual(t *testing.T) {
	t.Run("ResourceMetrics", testFailDatatype(resourceMetricsA, resourceMetricsB))
	t.Run("ScopeMetrics", testFailDatatype(scopeMetricsA, scopeMetricsB))
	t.Run("Metrics", testFailDatatype(metricsA, metricsB))
	t.Run("Histogram", testFailDatatype(histogramA, histogramB))
	t.Run("SumInt64", testFailDatatype(sumInt64A, sumInt64B))
	t.Run("SumFloat64", testFailDatatype(sumFloat64A, sumFloat64B))
	t.Run("GaugeInt64", testFailDatatype(gaugeInt64A, gaugeInt64B))
	t.Run("GaugeFloat64", testFailDatatype(gaugeFloat64A, gaugeFloat64B))
	t.Run("HistogramDataPoint", testFailDatatype(histogramDataPointA, histogramDataPointB))
	t.Run("DataPointInt64", testFailDatatype(dataPointInt64A, dataPointInt64B))
	t.Run("DataPointFloat64", testFailDatatype(dataPointFloat64A, dataPointFloat64B))

}

func TestFailAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, sumInt64A, nil)
	AssertAggregationsEqual(t, sumFloat64A, gaugeFloat64A)
	AssertAggregationsEqual(t, unknownAggregation{}, unknownAggregation{})
	AssertAggregationsEqual(t, sumInt64A, sumInt64B)
	AssertAggregationsEqual(t, sumFloat64A, sumFloat64B)
	AssertAggregationsEqual(t, gaugeInt64A, gaugeInt64B)
	AssertAggregationsEqual(t, gaugeFloat64A, gaugeFloat64B)
	AssertAggregationsEqual(t, histogramA, histogramB)
}

func TestFailAssertAttribute(t *testing.T) {
	AssertHasAttributes(t, dataPointInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, dataPointFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, gaugeInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, gaugeFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, sumInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, sumFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, histogramDataPointA, attribute.Bool("A", false))
	AssertHasAttributes(t, histogramDataPointA, attribute.Bool("B", true))
	AssertHasAttributes(t, histogramA, attribute.Bool("A", false))
	AssertHasAttributes(t, histogramA, attribute.Bool("B", true))
	AssertHasAttributes(t, metricsA, attribute.Bool("A", false))
	AssertHasAttributes(t, metricsA, attribute.Bool("B", true))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("A", false))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("B", true))
}
