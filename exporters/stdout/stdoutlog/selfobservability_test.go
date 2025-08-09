// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestSelfObservability(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics)
	}{
		{
			name: "inflight",
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log 1"),
				}
				record1 := rf.NewRecord()

				rf.Body = log.StringValue("test log 2")
				record2 := rf.NewRecord()

				records := []sdklog.Record{record1, record2}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				require.NotEmpty(t, got.Metrics)

				var inflightMetric metricdata.Metrics
				for _, m := range got.Metrics {
					if m.Name == "otel.sdk.exporter.log.inflight" {
						inflightMetric = m
						break
					}
				}
				require.NotEmpty(t, inflightMetric.Name, "inflight metric should be present")

				sum, ok := inflightMetric.Data.(metricdata.Sum[int64])
				require.True(t, ok, "inflight metric should be a sum")
				require.Len(t, sum.DataPoints, 1, "should have one data point")
				require.Equal(t, int64(2), sum.DataPoints[0].Value, "should record 2 inflight records")
			},
		},
		{
			name: "exported",
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log 1"),
				}
				record1 := rf.NewRecord()

				rf.Body = log.StringValue("test log 2")
				record2 := rf.NewRecord()

				rf.Body = log.StringValue("test log 3")
				record3 := rf.NewRecord()

				records := []sdklog.Record{record1, record2, record3}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				require.NotEmpty(t, got.Metrics)

				var exportedMetric metricdata.Metrics
				for _, m := range got.Metrics {
					if m.Name == "otel.sdk.exporter.log.exported" {
						exportedMetric = m
						break
					}
				}
				require.NotEmpty(t, exportedMetric.Name, "exported metric should be present")

				sum, ok := exportedMetric.Data.(metricdata.Sum[int64])
				require.True(t, ok, "exported metric should be a sum")
				require.Len(t, sum.DataPoints, 1, "should have one data point")
				require.Equal(t, int64(3), sum.DataPoints[0].Value, "should record 3 exported records")
			},
		},
		{
			name: "duration",
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				rf := logtest.RecordFactory{
					Timestamp: time.Now(),
					Body:      log.StringValue("test log"),
				}
				record := rf.NewRecord()
				records := []sdklog.Record{record}

				err = exporter.Export(context.Background(), records)
				require.NoError(t, err)

				got := scopeMetrics()
				require.NotEmpty(t, got.Metrics)

				var durationMetric metricdata.Metrics
				for _, m := range got.Metrics {
					if m.Name == "otel.sdk.exporter.operation.duration" {
						durationMetric = m
						break
					}
				}
				require.NotEmpty(t, durationMetric.Name, "duration metric should be present")

				histogram, ok := durationMetric.Data.(metricdata.Histogram[float64])
				require.True(t, ok, "duration metric should be a histogram")
				require.Len(t, histogram.DataPoints, 1, "should have one data point")
				require.Greater(t, histogram.DataPoints[0].Sum, 0.0, "duration should be greater than 0")
			},
		},
		{
			name: "multiple_exports",
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				for i := 0; i < 3; i++ {
					rf := logtest.RecordFactory{
						Timestamp: time.Now(),
						Body:      log.StringValue("test log"),
					}
					record := rf.NewRecord()
					records := []sdklog.Record{record}
					err = exporter.Export(context.Background(), records)
					require.NoError(t, err)
				}

				got := scopeMetrics()
				require.NotEmpty(t, got.Metrics)

				var exportedMetric metricdata.Metrics
				for _, m := range got.Metrics {
					if m.Name == "otel.sdk.exporter.log.exported" {
						exportedMetric = m
						break
					}
				}
				require.NotEmpty(t, exportedMetric.Name, "exported metric should be present")

				sum, ok := exportedMetric.Data.(metricdata.Sum[int64])
				require.True(t, ok, "exported metric should be a sum")
				require.Len(t, sum.DataPoints, 1, "should have one data point")
				require.Equal(t, int64(3), sum.DataPoints[0].Value, "should record 3 total exported records")
			},
		},
		{
			name: "empty_records",
			test: func(t *testing.T, scopeMetrics func() metricdata.ScopeMetrics) {
				exporter, err := New()
				require.NoError(t, err)

				err = exporter.Export(context.Background(), []sdklog.Record{})
				require.NoError(t, err)

				got := scopeMetrics()
				require.Empty(t, got.Metrics, "no metrics should be recorded for empty records")
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
			prev := otel.GetMeterProvider()
			defer otel.SetMeterProvider(prev)
			r := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(r))
			otel.SetMeterProvider(mp)

			scopeMetrics := func() metricdata.ScopeMetrics {
				var got metricdata.ResourceMetrics
				err := r.Collect(context.Background(), &got)
				require.NoError(t, err)
				if len(got.ScopeMetrics) == 0 {
					return metricdata.ScopeMetrics{}
				}
				return got.ScopeMetrics[0]
			}
			tc.test(t, scopeMetrics)
		})
	}
}
