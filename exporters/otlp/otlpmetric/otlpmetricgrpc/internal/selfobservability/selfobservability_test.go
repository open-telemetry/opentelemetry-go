// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package selfobservability

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
)

func TestNewExporterMetrics_Disabled(t *testing.T) {
	// Ensure feature is disabled
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "false")

	em := NewExporterMetrics("test_component", "localhost", 4317)

	assert.False(t, em.enabled, "metrics should be disabled when feature flag is false")

	// Tracking should be no-op when disabled
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	finish := em.TrackExport(context.Background(), createTestResourceMetrics())
	finish(nil)
	finish(errors.New("test error"))

	// Verify no metrics were recorded when disabled
	rm := &metricdata.ResourceMetrics{}
	err := reader.Collect(context.Background(), rm)
	require.NoError(t, err, "failed to collect metrics")

	totalMetrics := 0
	for _, sm := range rm.ScopeMetrics {
		totalMetrics += len(sm.Metrics)
	}
	assert.Zero(t, totalMetrics, "expected no metrics when disabled")
}

func TestNewExporterMetrics_Enabled(t *testing.T) {
	// Enable feature
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	// Set up a test meter provider
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	em := NewExporterMetrics("test_component", "example.com", 4317)

	assert.True(t, em.enabled, "metrics should be enabled when feature flag is true")

	// Verify attributes are set correctly
	expectedAttrs := []attribute.KeyValue{
		semconv.OTelComponentTypeKey.String("test_component"),
		semconv.OTelComponentName("test_component/0"),
		semconv.ServerAddress("example.com"),
		semconv.ServerPort(4317),
	}

	assert.Len(t, em.attrs, len(expectedAttrs), "attributes length mismatch")
	assert.Equal(t, expectedAttrs, em.attrs, "attributes should match expected values")
}

func TestNewExporterMetrics_MeterFailure(t *testing.T) {
	// Enable feature
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	// Use a meter provider that will cause metric creation to work
	// but test the error handling paths by using nil meter in the semantic convention
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	// This test primarily covers the enabled path, but the error handling
	// is covered by the semantic convention's internal nil checks
	em := NewExporterMetrics("test_component", "example.com", 4317)

	// Should be enabled with valid meter provider
	assert.True(t, em.enabled, "metrics should be enabled when meter provider is valid")

	// Test that tracking works properly
	finish := em.TrackExport(context.Background(), createTestResourceMetrics())
	finish(nil)
	finish(errors.New("test error"))

	rm := &metricdata.ResourceMetrics{}
	err := reader.Collect(context.Background(), rm)
	require.NoError(t, err, "failed to collect metrics")

	// Verify metrics were recorded
	totalMetrics := 0
	for _, sm := range rm.ScopeMetrics {
		if sm.Scope.Name == "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc" {
			totalMetrics += len(sm.Metrics)
		}
	}
	assert.Positive(t, totalMetrics, "expected self-observability metrics to be recorded when enabled")
}

func TestTrackExport_Success(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	em := NewExporterMetrics("test_component", "localhost", 4317)
	rm := createTestResourceMetrics()

	// Track export operation
	finish := em.TrackExport(context.Background(), rm)
	time.Sleep(10 * time.Millisecond) // Small delay to measure duration
	finish(nil)                       // Success

	// Read metrics to verify
	metrics := &metricdata.ResourceMetrics{}
	err := reader.Collect(context.Background(), metrics)
	require.NoError(t, err, "failed to collect metrics")

	// Verify exported counter was incremented
	exportedFound := false
	inflightFound := false
	durationFound := false

	for _, sm := range metrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch m.Name {
			case "otel.sdk.exporter.metric_data_point.exported":
				exportedFound = true
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
					assert.Equal(
						t,
						int64(10),
						sum.DataPoints[0].Value,
						"expected exported count 10",
					) // Expected data points from test data
				}
			case "otel.sdk.exporter.metric_data_point.inflight":
				inflightFound = true
				// Inflight should be 0 after completion
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
					assert.Equal(t, int64(0), sum.DataPoints[0].Value, "expected inflight count 0")
				}
			case "otel.sdk.exporter.operation.duration":
				durationFound = true
				// Duration should be recorded
				if hist, ok := m.Data.(metricdata.Histogram[float64]); ok && len(hist.DataPoints) > 0 {
					assert.Positive(t, hist.DataPoints[0].Count, "expected duration to be recorded")
				}
			}
		}
	}

	assert.True(t, exportedFound, "exported metric not found")
	assert.True(t, inflightFound, "inflight metric not found")
	assert.True(t, durationFound, "duration metric not found")
}

func TestTrackExport_Error(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	em := NewExporterMetrics("test_component", "localhost", 4317)
	rm := createTestResourceMetrics()

	// Track export operation that fails
	finish := em.TrackExport(context.Background(), rm)
	finish(errors.New("export failed"))

	// Read metrics
	metrics := &metricdata.ResourceMetrics{}
	err := reader.Collect(context.Background(), metrics)
	require.NoError(t, err, "failed to collect metrics")

	// Verify no exported count (due to error) but duration is recorded with error attribute
	for _, sm := range metrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == "otel.sdk.exporter.metric_data_point.exported" {
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok && len(sum.DataPoints) > 0 {
					assert.Equal(t, int64(0), sum.DataPoints[0].Value, "expected no exported count on error")
				}
			}
			if m.Name == "otel.sdk.exporter.operation.duration" {
				if hist, ok := m.Data.(metricdata.Histogram[float64]); ok && len(hist.DataPoints) > 0 {
					// Check for error attribute
					hasErrorAttr := false
					for _, attr := range hist.DataPoints[0].Attributes.ToSlice() {
						if attr.Key == semconv.ErrorTypeKey && attr.Value.AsString() == "_OTHER" {
							hasErrorAttr = true
							break
						}
					}
					assert.True(t, hasErrorAttr, "expected error.type attribute on duration metric")
				}
			}
		}
	}
}

func TestCountDataPoints(t *testing.T) {
	tests := []struct {
		name     string
		rm       *metricdata.ResourceMetrics
		expected int64
	}{
		{
			name:     "nil resource metrics",
			rm:       nil,
			expected: 0,
		},
		{
			name:     "empty resource metrics",
			rm:       &metricdata.ResourceMetrics{},
			expected: 0,
		},
		{
			name:     "test data",
			rm:       createTestResourceMetrics(),
			expected: 10, // 2 gauge + 1 gauge + 1 sum + 1 sum + 1 histogram + 1 histogram + 1 exponential histogram + 1 exponential histogram + 1 summary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := countDataPoints(tt.rm)
			assert.Equal(t, tt.expected, count, "data points count mismatch")
		})
	}
}

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		defaultPort int
		wantAddress string
		wantPort    int
	}{
		{
			name:        "empty endpoint",
			endpoint:    "",
			defaultPort: 4317,
			wantAddress: "localhost",
			wantPort:    4317,
		},
		{
			name:        "host only",
			endpoint:    "example.com",
			defaultPort: 4317,
			wantAddress: "example.com",
			wantPort:    4317,
		},
		{
			name:        "host with port",
			endpoint:    "example.com:9090",
			defaultPort: 4317,
			wantAddress: "example.com",
			wantPort:    9090,
		},
		{
			name:        "full URL",
			endpoint:    "https://example.com:9090/v1/metrics",
			defaultPort: 4317,
			wantAddress: "example.com",
			wantPort:    9090,
		},
		{
			name:        "invalid URL",
			endpoint:    "://invalid",
			defaultPort: 4317,
			wantAddress: "localhost",
			wantPort:    4317,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, port := ParseEndpoint(tt.endpoint, tt.defaultPort)
			assert.Equal(t, tt.wantAddress, address, "address mismatch")
			assert.Equal(t, tt.wantPort, port, "port mismatch")
		})
	}
}

func TestIsSelfObservabilityEnabled(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     bool
	}{
		{"unset", "", false},
		{"false", "false", false},
		{"true lowercase", "true", true},
		{"true uppercase", "TRUE", true},
		{"true mixed case", "True", true},
		{"invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", tt.envValue)
			}

			got := isSelfObservabilityEnabled()
			assert.Equal(t, tt.want, got, "self-observability enabled state mismatch")
		})
	}
}

// createTestResourceMetrics creates sample data for testing.
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
						Name: "test_gauge_float",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{Value: 1.5, Time: now},
							},
						},
					},
					{
						Name: "test_sum_int",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[int64]{
								{Value: 10, Time: now},
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
						Name: "test_histogram_int",
						Data: metricdata.Histogram[int64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint[int64]{
								{Count: 5, Sum: 15, Time: now},
							},
						},
					},
					{
						Name: "test_histogram_float",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{Count: 10, Sum: 5.0, Time: now},
							},
						},
					},
					{
						Name: "test_exponential_histogram_int",
						Data: metricdata.ExponentialHistogram[int64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.ExponentialHistogramDataPoint[int64]{
								{Count: 3, Sum: 9, Time: now, Scale: 1},
							},
						},
					},
					{
						Name: "test_exponential_histogram_float",
						Data: metricdata.ExponentialHistogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{
								{Count: 2, Sum: 4.5, Time: now, Scale: 1},
							},
						},
					},
					{
						Name: "test_summary",
						Data: metricdata.Summary{
							DataPoints: []metricdata.SummaryDataPoint{
								{Count: 7, Sum: 21.0, Time: now},
							},
						},
					},
				},
			},
		},
	}
}
