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

//go:build go1.18 && tests_fail
// +build go1.18,tests_fail

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"testing"
)

// These tests are used to develop the failure messages of this packages
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
	t.Run("Sum", testFailDatatype(sumA, sumB))
	t.Run("Gauge", testFailDatatype(gaugeA, gaugeB))
	t.Run("HistogramDataPoint", testFailDatatype(histogramDataPointA, histogramDataPointB))
	t.Run("DataPoint", testFailDatatype(dataPointsA, dataPointsB))
	t.Run("Int64", testFailDatatype(int64A, int64B))
	t.Run("Float64", testFailDatatype(float64A, float64B))
}

func TestFailAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, sumA, nil)
	AssertAggregationsEqual(t, sumA, gaugeA)
	AssertAggregationsEqual(t, unknownAggregation{}, unknownAggregation{})
	AssertAggregationsEqual(t, sumA, sumB)
	AssertAggregationsEqual(t, gaugeA, gaugeB)
	AssertAggregationsEqual(t, histogramA, histogramB)
}

func TestFailAssertValuesEqual(t *testing.T) {
	AssertValuesEqual(t, int64A, nil)
	AssertValuesEqual(t, int64A, float64A)
	AssertValuesEqual(t, unknownValue{}, unknownValue{})
	AssertValuesEqual(t, int64A, int64B)
	AssertValuesEqual(t, float64A, float64B)
}
