// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpmetricgrpc

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/otest"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

func TestSelfObservability_Disabled(t *testing.T) {
	// Ensure self-observability is disabled
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "false")

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

	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	// Verify that no self-observability metrics are reported
	selfObsMetrics := &metricdata.ResourceMetrics{}
	err = reader.Collect(context.Background(), selfObsMetrics)
	require.NoError(t, err)

	// Check that no self-observability metrics exist
	selfObsMetricCount := 0
	for _, sm := range selfObsMetrics.ScopeMetrics {
		if sm.Scope.Name == "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc" {
			selfObsMetricCount += len(sm.Metrics)
		}
	}
	assert.Equal(t, 0, selfObsMetricCount, "expected no self-observability metrics when disabled")
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

	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	var got metricdata.ResourceMetrics
	err = reader.Collect(context.Background(), &got)
	require.NoError(t, err)
	require.Len(t, got.ScopeMetrics, 1)

	serverAddr, serverPort := parseEndpoint(coll.Addr().String())

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
			Version:   sdk.Version(),
			SchemaURL: semconv.SchemaURL,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName("otlp_grpc_metric_exporter/0"),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
							Value: 4,
						},
					},
				},
			},
			{
				Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName("otlp_grpc_metric_exporter/0"),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
							Value: 0,
						},
					},
				},
			},
			{
				Name:        otelconv.SDKExporterOperationDuration{}.Name(),
				Description: otelconv.SDKExporterOperationDuration{}.Description(),
				Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
				Data: metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName("otlp_grpc_metric_exporter/0"),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
							Count: 1,
						},
					},
				},
			},
		},
	}

	metricdatatest.AssertEqual(t, want, got.ScopeMetrics[0],
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreValue())
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

	var got metricdata.ResourceMetrics
	err = reader.Collect(context.Background(), &got)
	require.NoError(t, err)
	require.Len(t, got.ScopeMetrics, 1)

	actualComponentName := extractComponentName(got.ScopeMetrics[0])

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
			Version:   sdk.Version(),
			SchemaURL: semconv.SchemaURL,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName(actualComponentName),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddress("invalid"),
								semconv.ServerPort(999999),
							),
						},
					},
				},
			},
			{
				Name:        otelconv.SDKExporterOperationDuration{}.Name(),
				Description: otelconv.SDKExporterOperationDuration{}.Description(),
				Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
				Data: metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes: attribute.NewSet(
								semconv.ErrorTypeOther,
								semconv.OTelComponentName(actualComponentName),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddress("invalid"),
								semconv.ServerPort(999999),
							),
						},
					},
				},
			},
		},
	}

	metricdatatest.AssertEqual(t, want, got.ScopeMetrics[0],
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreValue())
}

func TestSelfObservability_EndpointParsing(t *testing.T) {
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

	rm := createTestResourceMetrics()
	err = exp.Export(context.Background(), rm)
	require.NoError(t, err)

	var got metricdata.ResourceMetrics
	err = reader.Collect(context.Background(), &got)
	require.NoError(t, err)
	require.Len(t, got.ScopeMetrics, 1)

	serverAddr, serverPort := parseEndpoint(coll.Addr().String())

	var actualComponentName string
	if len(got.ScopeMetrics[0].Metrics) > 0 {
		if data, ok := got.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64]); ok && len(data.DataPoints) > 0 {
			attrs := data.DataPoints[0].Attributes.ToSlice()
			for _, attr := range attrs {
				if attr.Key == semconv.OTelComponentNameKey {
					actualComponentName = attr.Value.AsString()
					break
				}
			}
		}
	}

	want := metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
			Version:   sdk.Version(),
			SchemaURL: semconv.SchemaURL,
		},
		Metrics: []metricdata.Metrics{
			{
				Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: true,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName(actualComponentName),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
						},
					},
				},
			},
			{
				Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
				Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
				Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
				Data: metricdata.Sum[int64]{
					Temporality: metricdata.CumulativeTemporality,
					IsMonotonic: false,
					DataPoints: []metricdata.DataPoint[int64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName(actualComponentName),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
						},
					},
				},
			},
			{
				Name:        otelconv.SDKExporterOperationDuration{}.Name(),
				Description: otelconv.SDKExporterOperationDuration{}.Description(),
				Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
				Data: metricdata.Histogram[float64]{
					Temporality: metricdata.CumulativeTemporality,
					DataPoints: []metricdata.HistogramDataPoint[float64]{
						{
							Attributes: attribute.NewSet(
								semconv.OTelComponentName(actualComponentName),
								semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
								semconv.ServerAddressKey.String(serverAddr),
								semconv.ServerPortKey.Int(serverPort),
							),
						},
					},
				},
			},
		},
	}

	metricdatatest.AssertEqual(t, want, got.ScopeMetrics[0],
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreValue())
}

// parseEndpoint extracts server address and port from endpoint string.
func parseEndpoint(endpoint string) (string, int) {
	host, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "localhost", 4317
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 4317
	}

	return host, port
}

// extractComponentName extracts the component name from metrics data to handle dynamic counter.
func extractComponentName(scopeMetrics metricdata.ScopeMetrics) string {
	for _, m := range scopeMetrics.Metrics {
		switch data := m.Data.(type) {
		case metricdata.Sum[int64]:
			if len(data.DataPoints) > 0 {
				attrs := data.DataPoints[0].Attributes.ToSlice()
				for _, attr := range attrs {
					if attr.Key == semconv.OTelComponentNameKey {
						return attr.Value.AsString()
					}
				}
			}
		case metricdata.Histogram[float64]:
			if len(data.DataPoints) > 0 {
				attrs := data.DataPoints[0].Attributes.ToSlice()
				for _, attr := range attrs {
					if attr.Key == semconv.OTelComponentNameKey {
						return attr.Value.AsString()
					}
				}
			}
		}
	}
	return ""
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
