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

package prometheus // import "go.opentelemetry.io/otel/bridge/prometheus"

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

const (
	traceIDStr = "4bf92f3577b34da6a3ce929d0e0e4736"
	spanIDStr  = "00f067aa0ba902b7"
)

func TestProduce(t *testing.T) {
	testCases := []struct {
		name     string
		testFn   func(*prometheus.Registry)
		expected []metricdata.ScopeMetrics
		wantErr  error
	}{
		{
			name:   "no metrics registered",
			testFn: func(*prometheus.Registry) {},
		},
		{
			name: "gauge",
			testFn: func(reg *prometheus.Registry) {
				metric := prometheus.NewGauge(prometheus.GaugeOpts{
					Name: "test_gauge_metric",
					Help: "A gauge metric for testing",
				})
				reg.MustRegister(metric)
				metric.Set(123.4)
			},
			expected: []metricdata.ScopeMetrics{{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "test_gauge_metric",
						Description: "A gauge metric for testing",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(),
									Value:      123.4,
								},
							},
						},
					},
				},
			}},
		},
		{
			name: "counter",
			testFn: func(reg *prometheus.Registry) {
				metric := prometheus.NewCounter(prometheus.CounterOpts{
					Name: "test_counter_metric",
					Help: "A counter metric for testing",
				})
				reg.MustRegister(metric)
				metric.(prometheus.ExemplarAdder).AddWithExemplar(
					245.3, prometheus.Labels{
						"trace_id":        traceIDStr,
						"span_id":         spanIDStr,
						"other_attribute": "abcd",
					},
				)
			},
			expected: []metricdata.ScopeMetrics{{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "test_counter_metric",
						Description: "A counter metric for testing",
						Data: metricdata.Sum[float64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(),
									Value:      245.3,
									Exemplars: []metricdata.Exemplar[float64]{
										{
											Value:              245.3,
											TraceID:            []byte(traceIDStr),
											SpanID:             []byte(spanIDStr),
											FilteredAttributes: []attribute.KeyValue{attribute.String("other_attribute", "abcd")},
										},
									},
								},
							},
						},
					},
				},
			}},
		},
		{
			name: "summary dropped",
			testFn: func(reg *prometheus.Registry) {
				metric := prometheus.NewSummary(prometheus.SummaryOpts{
					Name: "test_summary_metric",
					Help: "A summary metric for testing",
				})
				reg.MustRegister(metric)
				metric.Observe(15.0)
			},
			wantErr: errUnsupportedType,
		},
		{
			name: "histogram",
			testFn: func(reg *prometheus.Registry) {
				metric := prometheus.NewHistogram(prometheus.HistogramOpts{
					Name: "test_histogram_metric",
					Help: "A histogram metric for testing",
				})
				reg.MustRegister(metric)
				metric.(prometheus.ExemplarObserver).ObserveWithExemplar(
					578.3, prometheus.Labels{
						"trace_id":        traceIDStr,
						"span_id":         spanIDStr,
						"other_attribute": "efgh",
					},
				)
			},
			expected: []metricdata.ScopeMetrics{{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "test_histogram_metric",
						Description: "A histogram metric for testing",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{
									Count:        1,
									Sum:          578.3,
									Bounds:       []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
									BucketCounts: []uint64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
									Attributes:   attribute.NewSet(),
									Exemplars: []metricdata.Exemplar[float64]{
										{
											Value:   578.3,
											TraceID: []byte(traceIDStr),
											SpanID:  []byte(spanIDStr),
											FilteredAttributes: []attribute.KeyValue{
												attribute.String("other_attribute", "efgh"),
											},
										},
									},
								},
							},
						},
					},
				},
			}},
		},
		{
			name: "partial success",
			testFn: func(reg *prometheus.Registry) {
				metric := prometheus.NewGauge(prometheus.GaugeOpts{
					Name: "test_gauge_metric",
					Help: "A gauge metric for testing",
				})
				reg.MustRegister(metric)
				metric.Set(123.4)
				unsupportedMetric := prometheus.NewSummary(prometheus.SummaryOpts{
					Name: "test_summary_metric",
					Help: "A summary metric for testing",
				})
				reg.MustRegister(unsupportedMetric)
				unsupportedMetric.Observe(15.0)
			},
			expected: []metricdata.ScopeMetrics{{
				Scope: instrumentation.Scope{
					Name: scopeName,
				},
				Metrics: []metricdata.Metrics{
					{
						Name:        "test_gauge_metric",
						Description: "A gauge metric for testing",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(),
									Value:      123.4,
								},
							},
						},
					},
				},
			}},
			wantErr: errUnsupportedType,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			reg := prometheus.NewRegistry()
			tt.testFn(reg)
			p := NewMetricProducer(WithGatherer(reg))
			output, err := p.Produce(context.Background())
			if tt.wantErr == nil {
				assert.Nil(t, err)
			}
			require.Equal(t, len(output), len(tt.expected))
			for i := range output {
				metricdatatest.AssertEqual(t, tt.expected[i], output[i], metricdatatest.IgnoreTimestamp())
			}
		})
	}
}
