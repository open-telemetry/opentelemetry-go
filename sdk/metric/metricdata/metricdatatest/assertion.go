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

// Package metricdatatest provides testing functionality for use with the
// metricdata package.
package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"fmt"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Datatypes are the concrete data-types the metricdata package provides.
type Datatypes interface {
	metricdata.DataPoint | metricdata.Float64 | metricdata.Gauge | metricdata.Histogram | metricdata.HistogramDataPoint | metricdata.Int64 | metricdata.Metrics | metricdata.ResourceMetrics | metricdata.ScopeMetrics | metricdata.Sum

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
		return assertCompare(equalDataPoints(e, actual.(metricdata.DataPoint)))(t)
	case metricdata.Float64:
		return assertCompare(equalFloat64(e, actual.(metricdata.Float64)))(t)
	case metricdata.Gauge:
		return assertCompare(equalGauges(e, actual.(metricdata.Gauge)))(t)
	case metricdata.Histogram:
		return assertCompare(equalHistograms(e, actual.(metricdata.Histogram)))(t)
	case metricdata.HistogramDataPoint:
		return assertCompare(equalHistogramDataPoints(e, actual.(metricdata.HistogramDataPoint)))(t)
	case metricdata.Int64:
		return assertCompare(equalInt64(e, actual.(metricdata.Int64)))(t)
	case metricdata.Metrics:
		return assertCompare(equalMetrics(e, actual.(metricdata.Metrics)))(t)
	case metricdata.ResourceMetrics:
		return assertCompare(equalResourceMetrics(e, actual.(metricdata.ResourceMetrics)))(t)
	case metricdata.ScopeMetrics:
		return assertCompare(equalScopeMetrics(e, actual.(metricdata.ScopeMetrics)))(t)
	case metricdata.Sum:
		return assertCompare(equalSums(e, actual.(metricdata.Sum)))(t)
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
func assertCompare(reasons []string) func(*testing.T) bool {
	return func(t *testing.T) bool {
		t.Helper()
		if len(reasons) > 0 {
			t.Error(strings.Join(reasons, "\n"))
			return false
		}
		return true
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
