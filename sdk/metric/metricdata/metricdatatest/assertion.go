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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Datatypes are the concrete data-types the metricdata package provides.
type Datatypes interface {
	metricdata.DataPoint |
		metricdata.Float64 |
		metricdata.Gauge |
		metricdata.Histogram |
		metricdata.HistogramDataPoint |
		metricdata.Int64 |
		metricdata.Metrics |
		metricdata.ResourceMetrics |
		metricdata.ScopeMetrics |
		metricdata.Sum

	// Interface types are not allowed in union types, therefore the
	// Aggregation and Value type from metricdata are not included here.
}

// AssertEqual asserts that the two concrete data-types from the metricdata
// package are equal.
func AssertEqual[T Datatypes](t *testing.T, expected, actual T) bool {
	t.Helper()
	// Generic types cannot be type switched on. Convert them to interfaces by
	// passing to assertEqual, which performs the correct functionality based
	// on the type.
	//
	// This function exists, instead of just exporting assertEqual, to ensure
	// the expected and actual types are not any and match.
	return assertEqual(t, expected, actual)
}

func assertEqual(t *testing.T, expected, actual interface{}) bool {
	t.Helper()
	switch e := expected.(type) {
	case metricdata.DataPoint:
		a := actual.(metricdata.DataPoint)
		return assertCompare(equalDataPoints(e, a))(t)
	case metricdata.Float64:
		a := actual.(metricdata.Float64)
		return assertCompare(equalFloat64(e, a))(t)
	case metricdata.Gauge:
		a := actual.(metricdata.Gauge)
		return assertCompare(equalGauges(e, a))(t)
	case metricdata.Histogram:
		a := actual.(metricdata.Histogram)
		return assertCompare(equalHistograms(e, a))(t)
	case metricdata.HistogramDataPoint:
		a := actual.(metricdata.HistogramDataPoint)
		return assertCompare(equalHistogramDataPoints(e, a))(t)
	case metricdata.Int64:
		a := actual.(metricdata.Int64)
		return assertCompare(equalInt64(e, a))(t)
	case metricdata.Metrics:
		a := actual.(metricdata.Metrics)
		return assertCompare(equalMetrics(e, a))(t)
	case metricdata.ResourceMetrics:
		a := actual.(metricdata.ResourceMetrics)
		return assertCompare(equalResourceMetrics(e, a))(t)
	case metricdata.ScopeMetrics:
		a := actual.(metricdata.ScopeMetrics)
		return assertCompare(equalScopeMetrics(e, a))(t)
	case metricdata.Sum:
		a := actual.(metricdata.Sum)
		return assertCompare(equalSums(e, a))(t)
	default:
		// assertEqual is unexported and we control all types passed to this
		// with AssertEqual, panic early to signal to developers when we
		// change things in an incompatible way early.
		panic(fmt.Sprintf("unknown types: %T", expected))
	}
}

// assertCompare evaluates the return value of an equality check function. The
// return function will produce an appropriate testing error if equal is
// false.
func assertCompare(equal bool, reasons []string) func(*testing.T) bool { // nolint: revive  // equal is not a control flag.
	return func(t *testing.T) bool {
		t.Helper()
		if !equal {
			if len(reasons) > 0 {
				t.Error(strings.Join(reasons, "\n"))
			} else {
				t.Fail()
			}
		}
		return equal
	}
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation) bool {
	t.Helper()
	return assertCompare(equalAggregations(expected, actual))(t)
}

// AssertValuesEqual asserts that two Values are equal.
func AssertValuesEqual(t *testing.T, expected, actual metricdata.Value) bool {
	t.Helper()
	return assertCompare(equalValues(expected, actual))(t)
}
