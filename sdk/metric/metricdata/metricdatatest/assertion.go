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

	// Generic types cannot be type asserted. Use an interface instead.
	aIface := interface{}(actual)

	var r []string
	switch e := interface{}(expected).(type) {
	case metricdata.DataPoint:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint))
	case metricdata.Float64:
		r = equalFloat64(e, aIface.(metricdata.Float64))
	case metricdata.Gauge:
		r = equalGauges(e, aIface.(metricdata.Gauge))
	case metricdata.Histogram:
		r = equalHistograms(e, aIface.(metricdata.Histogram))
	case metricdata.HistogramDataPoint:
		r = equalHistogramDataPoints(e, aIface.(metricdata.HistogramDataPoint))
	case metricdata.Int64:
		r = equalInt64(e, aIface.(metricdata.Int64))
	case metricdata.Metrics:
		r = equalMetrics(e, aIface.(metricdata.Metrics))
	case metricdata.ResourceMetrics:
		r = equalResourceMetrics(e, aIface.(metricdata.ResourceMetrics))
	case metricdata.ScopeMetrics:
		r = equalScopeMetrics(e, aIface.(metricdata.ScopeMetrics))
	case metricdata.Sum:
		r = equalSums(e, aIface.(metricdata.Sum))
	default:
		// We control all types passed to this, panic to signal developers
		// early they changed things in an incompatible way.
		panic(fmt.Sprintf("unknown types: %T", expected))
	}

	if len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation) bool {
	t.Helper()
	if r := equalAggregations(expected, actual); len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}

// AssertValuesEqual asserts that two Values are equal.
func AssertValuesEqual(t *testing.T, expected, actual metricdata.Value) bool {
	t.Helper()
	if r := equalValues(expected, actual); len(r) > 0 {
		t.Error(r)
		return false
	}
	return true
}
