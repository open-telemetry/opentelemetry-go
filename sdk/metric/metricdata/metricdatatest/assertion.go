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
	metricdata.DataPoint[float64] |
		metricdata.DataPoint[int64] |
		metricdata.Gauge[float64] |
		metricdata.Gauge[int64] |
		metricdata.Histogram |
		metricdata.HistogramDataPoint |
		metricdata.Metrics |
		metricdata.ResourceMetrics |
		metricdata.ScopeMetrics |
		metricdata.Sum[float64] |
		metricdata.Sum[int64]

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
	case metricdata.DataPoint[int64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[int64]))
	case metricdata.DataPoint[float64]:
		r = equalDataPoints(e, aIface.(metricdata.DataPoint[float64]))
	case metricdata.Gauge[int64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[int64]))
	case metricdata.Gauge[float64]:
		r = equalGauges(e, aIface.(metricdata.Gauge[float64]))
	case metricdata.Histogram:
		r = equalHistograms(e, aIface.(metricdata.Histogram))
	case metricdata.HistogramDataPoint:
		r = equalHistogramDataPoints(e, aIface.(metricdata.HistogramDataPoint))
	case metricdata.Metrics:
		r = equalMetrics(e, aIface.(metricdata.Metrics))
	case metricdata.ResourceMetrics:
		r = equalResourceMetrics(e, aIface.(metricdata.ResourceMetrics))
	case metricdata.ScopeMetrics:
		r = equalScopeMetrics(e, aIface.(metricdata.ScopeMetrics))
	case metricdata.Sum[int64]:
		r = equalSums(e, aIface.(metricdata.Sum[int64]))
	case metricdata.Sum[float64]:
		r = equalSums(e, aIface.(metricdata.Sum[float64]))
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
