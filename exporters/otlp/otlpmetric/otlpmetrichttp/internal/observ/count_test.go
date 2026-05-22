// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

func TestCountDataPoints(t *testing.T) {
	tests := []struct {
		name string
		rm   *metricpb.ResourceMetrics
		want int64
	}{
		{
			name: "nil",
			rm:   nil,
			want: 0,
		},
		{
			name: "empty",
			rm:   &metricpb.ResourceMetrics{},
			want: 0,
		},
		{
			name: "gauge",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Gauge{
									Gauge: &metricpb.Gauge{
										DataPoints: []*metricpb.NumberDataPoint{{}, {}},
									},
								},
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "sum",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										DataPoints: []*metricpb.NumberDataPoint{{}, {}, {}},
									},
								},
							},
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "histogram",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Histogram{
									Histogram: &metricpb.Histogram{
										DataPoints: []*metricpb.HistogramDataPoint{{}},
									},
								},
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "exponential_histogram",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_ExponentialHistogram{
									ExponentialHistogram: &metricpb.ExponentialHistogram{
										DataPoints: []*metricpb.ExponentialHistogramDataPoint{{}, {}},
									},
								},
							},
						},
					},
				},
			},
			want: 2,
		},
		{
			name: "summary",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Summary{
									Summary: &metricpb.Summary{
										DataPoints: []*metricpb.SummaryDataPoint{{}, {}, {}, {}},
									},
								},
							},
						},
					},
				},
			},
			want: 4,
		},
		{
			name: "multiple",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Gauge{
									Gauge: &metricpb.Gauge{
										DataPoints: []*metricpb.NumberDataPoint{{}},
									},
								},
							},
							{
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										DataPoints: []*metricpb.NumberDataPoint{{}, {}},
									},
								},
							},
						},
					},
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Histogram{
									Histogram: &metricpb.Histogram{
										DataPoints: []*metricpb.HistogramDataPoint{{}, {}, {}},
									},
								},
							},
						},
					},
				},
			},
			want: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countDataPoints(tt.rm)
			assert.Equal(t, tt.want, got)
		})
	}
}
