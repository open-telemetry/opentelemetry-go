// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
	t.Run("HistogramInt64", testFailDatatype(histogramInt64A, histogramInt64B))
	t.Run("HistogramFloat64", testFailDatatype(histogramFloat64A, histogramFloat64B))
	t.Run("SumInt64", testFailDatatype(sumInt64A, sumInt64B))
	t.Run("SumFloat64", testFailDatatype(sumFloat64A, sumFloat64B))
	t.Run("GaugeInt64", testFailDatatype(gaugeInt64A, gaugeInt64B))
	t.Run("GaugeFloat64", testFailDatatype(gaugeFloat64A, gaugeFloat64B))
	t.Run("HistogramDataPointInt64", testFailDatatype(histogramDataPointInt64A, histogramDataPointInt64B))
	t.Run("HistogramDataPointFloat64", testFailDatatype(histogramDataPointFloat64A, histogramDataPointFloat64B))
	t.Run("DataPointInt64", testFailDatatype(dataPointInt64A, dataPointInt64B))
	t.Run("DataPointFloat64", testFailDatatype(dataPointFloat64A, dataPointFloat64B))
	t.Run("ExemplarInt64", testFailDatatype(exemplarInt64A, exemplarInt64B))
	t.Run("ExemplarFloat64", testFailDatatype(exemplarFloat64A, exemplarFloat64B))
	t.Run("Extrema", testFailDatatype(minA, minB))
}

func TestFailAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, sumInt64A, nil)
	AssertAggregationsEqual(t, sumFloat64A, gaugeFloat64A)
	AssertAggregationsEqual(t, unknownAggregation{}, unknownAggregation{})
	AssertAggregationsEqual(t, sumInt64A, sumInt64B)
	AssertAggregationsEqual(t, sumFloat64A, sumFloat64B)
	AssertAggregationsEqual(t, gaugeInt64A, gaugeInt64B)
	AssertAggregationsEqual(t, gaugeFloat64A, gaugeFloat64B)
	AssertAggregationsEqual(t, histogramInt64A, histogramInt64B)
	AssertAggregationsEqual(t, histogramFloat64A, histogramFloat64B)
}

func TestFailAssertAttribute(t *testing.T) {
	AssertHasAttributes(t, exemplarInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, exemplarFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, dataPointInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, dataPointFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, gaugeInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, gaugeFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, sumInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, sumFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, histogramDataPointInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, histogramDataPointFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, histogramInt64A, attribute.Bool("A", false))
	AssertHasAttributes(t, histogramFloat64A, attribute.Bool("B", true))
	AssertHasAttributes(t, metricsA, attribute.Bool("A", false))
	AssertHasAttributes(t, metricsA, attribute.Bool("B", true))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("A", false))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("B", true))
}
