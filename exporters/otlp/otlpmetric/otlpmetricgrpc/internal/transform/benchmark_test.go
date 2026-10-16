// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package transform

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// benchResourceMetrics builds a steady-state workload with attribute sets
// reused across exports.
func benchResourceMetrics(series int) *metricdata.ResourceMetrics {
	now := time.Now()
	dps := make([]metricdata.DataPoint[float64], 0, series)
	for i := range series {
		attrs := attribute.NewSet(
			attribute.String("csfversion", "1.183.0"),
			attribute.String("env", "intl-shared-prod"),
			attribute.String("projectname", "intl/central-resolution-hub"),
			attribute.String("host", "i-0bdf11c292977661a"),
			attribute.String("configname", "production"),
			attribute.String("datadog.metric.type", "gauge"),
			attribute.Int("series", i),
		)
		dps = append(dps, metricdata.DataPoint[float64]{
			Attributes: attrs,
			StartTime:  now,
			Time:       now,
			Value:      float64(i),
		})
	}

	return &metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{{
			Metrics: []metricdata.Metrics{{
				Name: "central_resolution_hub.query.result_count.max",
				Data: metricdata.Gauge[float64]{DataPoints: dps},
			}},
		}},
	}
}

func BenchmarkResourceMetricsSteadyState(b *testing.B) {
	rm := benchResourceMetrics(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ResourceMetrics(rm)
	}
}
