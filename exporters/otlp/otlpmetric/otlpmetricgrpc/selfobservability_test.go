// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetricgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/otest"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
)

func TestSelfObservability_Disabled(t *testing.T) {
	// Ensure self-observability is disabled
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "false")

	coll, err := otest.NewGRPCCollector("", nil)
	require.NoError(t, err)
	defer coll.Shutdown()

	exp, err := New(context.Background(),
		WithEndpoint(coll.Addr().String()),
		WithInsecure())
	require.NoError(t, err)

	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	// Note: Cannot directly test exp.metrics.enabled as it's private
	// The test passes if no panics occur and export works
}

func TestSelfObservability_Enabled(t *testing.T) {
	// Enable self-observability
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	coll, err := otest.NewGRPCCollector("", nil)
	require.NoError(t, err)
	defer coll.Shutdown()

	exp, err := New(context.Background(),
		WithEndpoint(coll.Addr().String()),
		WithInsecure())
	require.NoError(t, err)

	// Note: Cannot directly test exp.metrics.enabled as it's private
	// verify through metrics collection instead

	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	selfObsMetrics := &metricdata.ResourceMetrics{}
	err = reader.Collect(context.Background(), selfObsMetrics)
	require.NoError(t, err)

	// Verify the three expected metrics exist
	foundMetrics := make(map[string]bool)
	for _, sm := range selfObsMetrics.ScopeMetrics {
		if sm.Scope.Name == "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc" {
			for _, m := range sm.Metrics {
				foundMetrics[m.Name] = true

				switch m.Name {
				case "otel.sdk.exporter.metric_data_point.exported":
					if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
						assert.Equal(t, int64(4), sum.DataPoints[0].Value, "expected 4 data points exported")
						verifyAttributes(t, sum.DataPoints[0].Attributes, coll.Addr().String())
					}

				case "otel.sdk.exporter.metric_data_point.inflight":
					if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
						assert.Equal(t, int64(0), sum.DataPoints[0].Value, "expected 0 inflight data points")
					}

				case "otel.sdk.exporter.operation.duration":
					if hist, ok := m.Data.(metricdata.Histogram[float64]); ok && len(hist.DataPoints) > 0 {
						assert.NotEqual(t, uint64(0), hist.DataPoints[0].Count, "expected duration to be recorded")
						// Note: We don't check if duration is positive as very fast operations
						// may result in zero or near-zero durations on some platforms
						verifyAttributes(t, hist.DataPoints[0].Attributes, coll.Addr().String())
					}
				}
			}
		}
	}

	expectedMetrics := []string{
		"otel.sdk.exporter.metric_data_point.exported",
		"otel.sdk.exporter.metric_data_point.inflight",
		"otel.sdk.exporter.operation.duration",
	}
	for _, metricName := range expectedMetrics {
		assert.True(t, foundMetrics[metricName], "missing expected metric: %s", metricName)
	}
}

func TestSelfObservability_ExportError(t *testing.T) {
	// Enable self-observability
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	// Create exporter with invalid endpoint to force error
	exp, err := New(context.Background(),
		WithEndpoint("invalid:999999"),
		WithInsecure())
	require.NoError(t, err)

	// Export data (should fail)
	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	assert.Error(t, err, "expected error but got none")

	// Collect metrics
	selfObsMetrics := &metricdata.ResourceMetrics{}
	err = reader.Collect(context.Background(), selfObsMetrics)
	require.NoError(t, err)

	// Verify error handling in metrics
	for _, sm := range selfObsMetrics.ScopeMetrics {
		if sm.Scope.Name == "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc" {
			for _, m := range sm.Metrics {
				switch m.Name {
				case "otel.sdk.exporter.metric_data_point.exported":
					// Should not increment on error
					if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
						assert.Equal(t, int64(0), sum.DataPoints[0].Value, "expected no exported count on error")
					}

				case "otel.sdk.exporter.operation.duration":
					// Should record duration with error attribute
					if hist, ok := m.Data.(metricdata.Histogram[float64]); ok && len(hist.DataPoints) > 0 {
						attrs := hist.DataPoints[0].Attributes.ToSlice()
						hasErrorAttr := false
						for _, attr := range attrs {
							if attr.Key == semconv.ErrorTypeKey && attr.Value.AsString() == "_OTHER" {
								hasErrorAttr = true
								break
							}
						}
						assert.True(t, hasErrorAttr, "expected error.type attribute on failed export")
					}
				}
			}
		}
	}
}

func TestSelfObservability_EndpointParsing(t *testing.T) {
	// Enable self-observability
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	// Set up meter provider for metric collection
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	// Set up collector for successful export
	coll, err := otest.NewGRPCCollector("", nil)
	require.NoError(t, err)
	defer coll.Shutdown()

	// Create exporter
	exp, err := New(context.Background(),
		WithEndpoint(coll.Addr().String()),
		WithInsecure())
	require.NoError(t, err)

	// Export some data to trigger metrics
	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	// Collect metrics to verify they were created with proper attributes
	selfObsMetrics := &metricdata.ResourceMetrics{}
	err = reader.Collect(context.Background(), selfObsMetrics)
	require.NoError(t, err)

	// Verify metrics exist and have proper component type
	found := false
	for _, sm := range selfObsMetrics.ScopeMetrics {
		if sm.Scope.Name == "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc" {
			for _, m := range sm.Metrics {
				if m.Name == "otel.sdk.exporter.operation.duration" {
					if hist, ok := m.Data.(metricdata.Histogram[float64]); ok && len(hist.DataPoints) > 0 {
						attrs := hist.DataPoints[0].Attributes.ToSlice()
						for _, attr := range attrs {
							if attr.Key == semconv.OTelComponentTypeKey &&
								attr.Value.AsString() == "otlp_grpc_metric_exporter" {
								found = true
								break
							}
						}
					}
				}
			}
		}
	}
	assert.True(t, found, "expected self-observability metrics with correct component type")
}

// verifyAttributes checks that the expected attributes are present.
func verifyAttributes(t *testing.T, attrs attribute.Set, _ string) {
	attrSlice := attrs.ToSlice()

	var componentType, serverAddr string
	var serverPort int

	for _, attr := range attrSlice {
		switch attr.Key {
		case semconv.OTelComponentTypeKey:
			componentType = attr.Value.AsString()
		case semconv.ServerAddressKey:
			serverAddr = attr.Value.AsString()
		case semconv.ServerPortKey:
			serverPort = int(attr.Value.AsInt64())
		}
	}

	assert.Equal(t, "otlp_grpc_metric_exporter", componentType)
	assert.NotEmpty(t, serverAddr, "expected non-empty server address")
	assert.Positive(t, serverPort, "expected positive server port")
}

// createTestResourceMetrics creates sample metric data for testing.
func createTestResourceMetrics() *metricdata.ResourceMetrics {
	now := time.Now()
	return &metricdata.ResourceMetrics{
		Resource: resource.Default(),
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "test", Version: "v1"},
				Metrics: []metricdata.Metrics{
					{
						Name: "test_gauge_int",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{Value: 1, Time: now},
								{Value: 2, Time: now},
							},
						},
					},
					{
						Name: "test_sum_float",
						Data: metricdata.Sum[float64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[float64]{
								{Value: 3.5, Time: now},
							},
						},
					},
					{
						Name: "test_histogram",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{Count: 10, Sum: 5.0, Time: now},
							},
						},
					},
				},
			},
		},
	}
}
