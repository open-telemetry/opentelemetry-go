// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetricgrpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestSelfObservability_Disabled verifies that when OTEL_GO_X_OBSERVABILITY is not set,
// no observability metrics are collected.
func TestSelfObservability_Disabled(t *testing.T) {
	require.NoError(t, os.Unsetenv("OTEL_GO_X_OBSERVABILITY"))

	ctx := t.Context()
	exporter, err := New(ctx, WithEndpointURL("http://localhost:4317"), WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, exporter)
	defer func() {
		require.NoError(t, exporter.Shutdown(ctx))
	}()

	// Verify instrumentation is nil when disabled
	assert.Nil(t, exporter.instrumentation)
}

// TestSelfObservability_Enabled verifies that when OTEL_GO_X_OBSERVABILITY=true,
// observability metrics are initialized and collected.
func TestSelfObservability_Enabled(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	ctx := t.Context()

	// Setup a metric reader to collect the self-observability metrics
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))

	// Must set the global meter provider before creating the exporter
	// so that the instrumentation can use it
	oldProvider := otel.GetMeterProvider()
	otel.SetMeterProvider(provider)
	defer otel.SetMeterProvider(oldProvider)

	exporter, err := New(ctx, WithEndpointURL("http://localhost:4317"), WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, exporter)
	defer func() {
		require.NoError(t, exporter.Shutdown(ctx))
	}()

	// Verify instrumentation is initialized when enabled
	assert.NotNil(t, exporter.instrumentation)

	// Create some test metric data
	rm := &metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Metrics: []metricdata.Metrics{
					{
						Name: "test.counter",
						Data: metricdata.Sum[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{Value: 1},
								{Value: 2},
								{Value: 3},
							},
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
						},
					},
				},
			},
		},
	}

	// Export the metrics (this will fail to connect but that's ok for this test)
	// We're just verifying the observability metrics are recorded
	_ = exporter.Export(ctx, rm)

	// Give a small delay for metrics to be recorded
	time.Sleep(10 * time.Millisecond)

	// Collect the self-observability metrics
	var collected metricdata.ResourceMetrics
	err = reader.Collect(ctx, &collected)
	require.NoError(t, err)

	// Verify we have some metrics collected
	foundInflight := false
	foundExported := false
	foundDuration := false

	for _, sm := range collected.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch m.Name {
			case "otel.sdk.exporter.metric_data_point.inflight":
				foundInflight = true
			case "otel.sdk.exporter.metric_data_point.exported":
				foundExported = true
			case "otel.sdk.exporter.operation.duration":
				foundDuration = true
			}
		}
	}

	// All three metrics should be present
	assert.True(t, foundInflight, "inflight metric should be recorded")
	assert.True(t, foundExported, "exported metric should be recorded")
	assert.True(t, foundDuration, "duration metric should be recorded")
}

// TestSelfObservability_ExportError verifies that observability metrics correctly
// track failures when exports encounter errors.
func TestSelfObservability_ExportError(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	ctx := t.Context()

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))

	oldProvider := otel.GetMeterProvider()
	otel.SetMeterProvider(provider)
	defer otel.SetMeterProvider(oldProvider)

	// Use an invalid endpoint to force an error
	exporter, err := New(ctx, WithEndpointURL("http://invalid-endpoint:1234"), WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, exporter)
	defer func() {
		require.NoError(t, exporter.Shutdown(ctx))
	}()

	rm := &metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Metrics: []metricdata.Metrics{
					{
						Name: "test.counter",
						Data: metricdata.Sum[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{Value: 1},
							},
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
						},
					},
				},
			},
		},
	}

	// This export should fail but metrics should still be recorded
	exportCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	exportErr := exporter.Export(exportCtx, rm)
	assert.Error(t, exportErr, "export should fail with invalid endpoint")

	// Collect the self-observability metrics
	var collected metricdata.ResourceMetrics
	err = reader.Collect(ctx, &collected)
	require.NoError(t, err)

	// Verify error.type attribute is present in exported metric
	foundErrorMetric := false
	for _, sm := range collected.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "otel.sdk.exporter.metric_data_point.exported" {
				if data, ok := m.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range data.DataPoints {
						for _, attr := range dp.Attributes.ToSlice() {
							if string(attr.Key) == "error.type" {
								foundErrorMetric = true
								break
							}
						}
					}
				}
			}
		}
	}

	assert.True(t, foundErrorMetric, "error.type attribute should be present in exported metric on failure")
}
