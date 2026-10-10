// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type testAggregationMetricFilter struct {
	ReaderOption
	testMetric func(instrumentation.Scope, string, metricdata.Aggregation, string) int
	testAttrs  func(instrumentation.Scope, string, metricdata.Aggregation, string, attribute.Set) int
}

func (testAggregationMetricFilter) Experimental() {}

func (f testAggregationMetricFilter) TestMetric(
	scope instrumentation.Scope,
	name string,
	data metricdata.Aggregation,
	unit string,
) int {
	return f.testMetric(scope, name, data, unit)
}

func (f testAggregationMetricFilter) TestAttributes(
	scope instrumentation.Scope,
	name string,
	data metricdata.Aggregation,
	unit string,
	attrs attribute.Set,
) int {
	return f.testAttrs(scope, name, data, unit, attrs)
}

func TestFilterScopeMetrics(t *testing.T) {
	keep := attribute.NewSet(attribute.String("keep", "true"))
	drop := attribute.NewSet(attribute.String("drop", "true"))

	metrics := func(name string) []metricdata.Metrics {
		return []metricdata.Metrics{{
			Name: name,
			Data: metricdata.Sum[int64]{
				DataPoints: []metricdata.DataPoint[int64]{
					{Attributes: keep, Value: 1},
					{Attributes: drop, Value: 2},
				},
			},
		}}
	}

	tests := []struct {
		name         string
		scopeMetrics []metricdata.ScopeMetrics
		filter       testAggregationMetricFilter
		assert       func(*testing.T, []metricdata.ScopeMetrics)
	}{
		{
			name: "Accept",
			scopeMetrics: []metricdata.ScopeMetrics{{
				Metrics: metrics("metric"),
			}},
			filter: testAggregationMetricFilter{
				testMetric: func(instrumentation.Scope, string, metricdata.Aggregation, string) int {
					return metricFilterAccept
				},
			},
			assert: func(t *testing.T, got []metricdata.ScopeMetrics) {
				require.Len(t, got, 1)
				require.Len(t, got[0].Metrics, 1)
				assert.Equal(t, 2, dataPointCount(got[0].Metrics[0].Data))
			},
		},
		{
			name: "Drop",
			scopeMetrics: []metricdata.ScopeMetrics{{
				Metrics: metrics("metric"),
			}},
			filter: testAggregationMetricFilter{
				testMetric: func(instrumentation.Scope, string, metricdata.Aggregation, string) int {
					return metricFilterDrop
				},
			},
			assert: func(t *testing.T, got []metricdata.ScopeMetrics) {
				assert.Empty(t, got)
			},
		},
		{
			name: "AcceptPartial",
			scopeMetrics: []metricdata.ScopeMetrics{{
				Metrics: metrics("metric"),
			}},
			filter: testAggregationMetricFilter{
				testMetric: func(instrumentation.Scope, string, metricdata.Aggregation, string) int {
					return metricFilterAcceptPartial
				},
				testAttrs: func(_ instrumentation.Scope, _ string, _ metricdata.Aggregation, _ string, attrs attribute.Set) int {
					if attrs.Equals(&drop) {
						return metricFilterAttrDrop
					}
					return metricFilterAttrAccept
				},
			},
			assert: func(t *testing.T, got []metricdata.ScopeMetrics) {
				require.Len(t, got, 1)
				require.Len(t, got[0].Metrics, 1)
				assert.Equal(t, 1, dataPointCount(got[0].Metrics[0].Data))
			},
		},
		{
			name: "AcceptPartial drops all data points",
			scopeMetrics: []metricdata.ScopeMetrics{{
				Metrics: metrics("metric"),
			}},
			filter: testAggregationMetricFilter{
				testMetric: func(instrumentation.Scope, string, metricdata.Aggregation, string) int {
					return metricFilterAcceptPartial
				},
				testAttrs: func(instrumentation.Scope, string, metricdata.Aggregation, string, attribute.Set) int {
					return metricFilterAttrDrop
				},
			},
			assert: func(t *testing.T, got []metricdata.ScopeMetrics) {
				assert.Empty(t, got)
			},
		},
		{
			name: "Drop empty scope",
			scopeMetrics: []metricdata.ScopeMetrics{
				{Scope: instrumentation.Scope{Name: "drop"}, Metrics: metrics("drop")},
				{Scope: instrumentation.Scope{Name: "keep"}, Metrics: metrics("keep")},
			},
			filter: testAggregationMetricFilter{
				testMetric: func(_ instrumentation.Scope, name string, _ metricdata.Aggregation, _ string) int {
					if name == "drop" {
						return metricFilterDrop
					}
					return metricFilterAccept
				},
			},
			assert: func(t *testing.T, got []metricdata.ScopeMetrics) {
				require.Len(t, got, 1)
				require.Len(t, got[0].Metrics, 1)
				assert.Equal(t, "keep", got[0].Scope.Name)
				assert.Equal(t, 2, dataPointCount(got[0].Metrics[0].Data))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterScopeMetrics(tt.scopeMetrics, tt.filter)
			tt.assert(t, got)
		})
	}
}

func TestFilterScopeMetricsAcceptPartialAggregations(t *testing.T) {
	keep := attribute.NewSet(attribute.String("keep", "true"))
	drop := attribute.NewSet(attribute.String("drop", "true"))
	filter := testAggregationMetricFilter{
		testMetric: func(instrumentation.Scope, string, metricdata.Aggregation, string) int {
			return metricFilterAcceptPartial
		},
		testAttrs: func(_ instrumentation.Scope, _ string, _ metricdata.Aggregation, _ string, attrs attribute.Set) int {
			if attrs.Equals(&drop) {
				return metricFilterAttrDrop
			}
			return metricFilterAttrAccept
		},
	}

	tests := []struct {
		name string
		data metricdata.Aggregation
	}{
		{
			"Sum int64",
			metricdata.Sum[int64]{DataPoints: []metricdata.DataPoint[int64]{{Attributes: keep}, {Attributes: drop}}},
		},
		{
			"Sum float64",
			metricdata.Sum[float64]{
				DataPoints: []metricdata.DataPoint[float64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"Gauge int64",
			metricdata.Gauge[int64]{DataPoints: []metricdata.DataPoint[int64]{{Attributes: keep}, {Attributes: drop}}},
		},
		{
			"Gauge float64",
			metricdata.Gauge[float64]{
				DataPoints: []metricdata.DataPoint[float64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"Histogram int64",
			metricdata.Histogram[int64]{
				DataPoints: []metricdata.HistogramDataPoint[int64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"Histogram float64",
			metricdata.Histogram[float64]{
				DataPoints: []metricdata.HistogramDataPoint[float64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"ExponentialHistogram int64",
			metricdata.ExponentialHistogram[int64]{
				DataPoints: []metricdata.ExponentialHistogramDataPoint[int64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"ExponentialHistogram float64",
			metricdata.ExponentialHistogram[float64]{
				DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{{Attributes: keep}, {Attributes: drop}},
			},
		},
		{
			"Summary",
			metricdata.Summary{DataPoints: []metricdata.SummaryDataPoint{{Attributes: keep}, {Attributes: drop}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterScopeMetrics([]metricdata.ScopeMetrics{{
				Metrics: []metricdata.Metrics{{Name: "metric", Data: tt.data}},
			}}, filter)
			require.Len(t, got, 1)
			require.Len(t, got[0].Metrics, 1)
			assert.Equal(t, 1, dataPointCount(got[0].Metrics[0].Data))
		})
	}
}

// testMetricFilterOption is a fake ReaderOption that also satisfies the
// internal metricFilter interface. It lets stable-package tests exercise the
// MetricFilter plumbing without importing the experimental x package.
type testMetricFilterOption struct {
	ReaderOption
	testMetric     func(instrumentation.Scope, string, metricdata.Aggregation, string) int
	testAttributes func(instrumentation.Scope, string, metricdata.Aggregation, string, attribute.Set) int
}

func (testMetricFilterOption) Experimental() {}

func (o testMetricFilterOption) TestMetric(
	scope instrumentation.Scope,
	name string,
	kind metricdata.Aggregation,
	unit string,
) int {
	return o.testMetric(scope, name, kind, unit)
}

func (o testMetricFilterOption) TestAttributes(
	scope instrumentation.Scope,
	name string,
	kind metricdata.Aggregation,
	unit string,
	attrs attribute.Set,
) int {
	return o.testAttributes(scope, name, kind, unit, attrs)
}

// testExporter is a minimal Exporter for tests that need a PeriodicReader.
type testExporter struct{}

func (testExporter) Temporality(InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (testExporter) Aggregation(InstrumentKind) Aggregation { return AggregationDefault{} }

func (testExporter) Export(context.Context, *metricdata.ResourceMetrics) error {
	return nil
}
func (testExporter) ForceFlush(context.Context) error { return nil }
func (testExporter) Shutdown(context.Context) error   { return nil }

// experimentalReaderOption embeds ReaderOption and only implements
// Experimental(). It is used to verify that experimental reader options do not
// panic when the embedded option interface is never applied.
type experimentalReaderOption struct{ ReaderOption }

func (experimentalReaderOption) Experimental() {}

func sumDataPointCount(rm *metricdata.ResourceMetrics, name string) int {
	m := findMetricByName(rm, name)
	if m == nil {
		return 0
	}
	sum, ok := m.Data.(metricdata.Sum[int64])
	if !ok {
		return -1
	}
	return len(sum.DataPoints)
}
