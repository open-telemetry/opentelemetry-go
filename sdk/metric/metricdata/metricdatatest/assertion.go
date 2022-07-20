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

	var r []string
	switch e := interface{}(expected).(type) {
	case metricdata.DataPoint:
		a := interface{}(actual).(metricdata.DataPoint)
		r = equalDataPoints(e, a)
	case metricdata.Float64:
		a := interface{}(actual).(metricdata.Float64)
		r = equalFloat64(e, a)
	case metricdata.Gauge:
		a := interface{}(actual).(metricdata.Gauge)
		r = equalGauges(e, a)
	case metricdata.Histogram:
		a := interface{}(actual).(metricdata.Histogram)
		r = equalHistograms(e, a)
	case metricdata.HistogramDataPoint:
		a := interface{}(actual).(metricdata.HistogramDataPoint)
		r = equalHistogramDataPoints(e, a)
	case metricdata.Int64:
		a := interface{}(actual).(metricdata.Int64)
		r = equalInt64(e, a)
	case metricdata.Metrics:
		a := interface{}(actual).(metricdata.Metrics)
		r = equalMetrics(e, a)
	case metricdata.ResourceMetrics:
		a := interface{}(actual).(metricdata.ResourceMetrics)
		r = equalResourceMetrics(e, a)
	case metricdata.ScopeMetrics:
		a := interface{}(actual).(metricdata.ScopeMetrics)
		r = equalScopeMetrics(e, a)
	case metricdata.Sum:
		a := interface{}(actual).(metricdata.Sum)
		r = equalSums(e, a)
	default:
		// We control all types passed to this, panic to signal developers
		// early they changed things in an incompatible way.
		panic(fmt.Sprintf("unknown types: %T", expected))
	}

	if len(r) > 0 {
		t.Error(strings.Join(r, "\n"))
		return false
	}
	return true
}

// AssertAggregationsEqual asserts that two Aggregations are equal.
func AssertAggregationsEqual(t *testing.T, expected, actual metricdata.Aggregation) bool {
	t.Helper()
	if r := equalAggregations(expected, actual); len(r) > 0 {
		t.Error(strings.Join(r, "\n"))
		return false
	}
	return true
}

// AssertValuesEqual asserts that two Values are equal.
func AssertValuesEqual(t *testing.T, expected, actual metricdata.Value) bool {
	t.Helper()
	if r := equalValues(expected, actual); len(r) > 0 {
		t.Error(strings.Join(r, "\n"))
		return false
	}
	return true
}
