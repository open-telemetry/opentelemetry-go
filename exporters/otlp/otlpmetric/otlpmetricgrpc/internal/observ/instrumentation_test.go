// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/semconv/v1.40.0/otelconv"

	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/transform"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc/internal/x"
)

func TestNewInstrumentation_Disabled(t *testing.T) {
	// Ensure feature is disabled
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "false")

	em, err := NewInstrumentation(0, "dns:///localhost:4317")
	require.NoError(t, err)
	assert.Nil(t, em, "metrics should be nil when feature flag is false")

	// Tracking should be no-op when disabled
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	testRm := createTestResourceMetrics()
	otlpRm, _ := transform.ResourceMetrics(testRm)
	em.TrackExport(t.Context(), otlpRm).End(nil)
	em.TrackExport(t.Context(), otlpRm).End(errors.New("test error"))

	// Verify no metrics were recorded when disabled
	rm := &metricdata.ResourceMetrics{}
	err = reader.Collect(t.Context(), rm)
	require.NoError(t, err, "failed to collect metrics")

	totalMetrics := 0
	for _, sm := range rm.ScopeMetrics {
		totalMetrics += len(sm.Metrics)
	}
	assert.Zero(t, totalMetrics, "expected no metrics when disabled")
}

func TestNewInstrumentation_Enabled(t *testing.T) {
	// Enable feature
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Set up a test meter provider
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(provider)

	em, err := NewInstrumentation(0, "dns:///example.com:4317")
	require.NoError(t, err)
	require.NotNil(t, em, "metrics should not be nil when feature flag is true")

	// Verify attributes are set correctly
	assert.Len(t, em.attrs, 4, "attributes length mismatch")

	// Find and verify each attribute
	var componentName, componentType, serverAddress string
	var serverPort int
	for _, attr := range em.attrs {
		switch attr.Key {
		case semconv.OTelComponentNameKey:
			componentName = attr.Value.AsString()
		case semconv.OTelComponentTypeKey:
			componentType = attr.Value.AsString()
		case semconv.ServerAddressKey:
			serverAddress = attr.Value.AsString()
		case semconv.ServerPortKey:
			serverPort = int(attr.Value.AsInt64())
		}
	}

	assert.True(
		t,
		strings.HasPrefix(componentName, "otlp_grpc_metric_exporter/"),
		"component name should start with otlp_grpc_metric_exporter/",
	)
	assert.Equal(t, "otlp_grpc_metric_exporter", componentType, "component type mismatch")
	assert.Equal(t, "example.com", serverAddress, "server address mismatch")
	assert.Equal(t, 4317, serverPort, "server port mismatch")
}

func TestTrackExport(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantMetrics func(actualComponentName string) []metricdata.Metrics
	}{
		{
			name: "success",
			err:  nil,
			wantMetrics: func(actualComponentName string) []metricdata.Metrics {
				return []metricdata.Metrics{
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Value: 10,
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
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
										semconv.OTelComponentName(actualComponentName),
										semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
										semconv.RPCResponseStatusCode(codes.OK.String()),
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Count: 1,
								},
							},
						},
					},
				}
			},
		},
		{
			name: "error",
			err:  errors.New("export failed"),
			wantMetrics: func(actualComponentName string) []metricdata.Metrics {
				return []metricdata.Metrics{
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Value: 0,
								},
								{
									Attributes: attribute.NewSet(
										semconv.ErrorType(errors.New("export failed")),
										semconv.OTelComponentName(actualComponentName),
										semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Value: 10,
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
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
										semconv.ErrorType(errors.New("export failed")),
										semconv.OTelComponentName(actualComponentName),
										semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
										semconv.RPCResponseStatusCode("Unknown"),
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Count: 1,
								},
							},
						},
					},
				}
			},
		},
		{
			name: "partial_success",
			err: internal.PartialSuccess{
				ErrorMessage:  "some points rejected",
				RejectedItems: 3,
				RejectedKind:  "metric data points",
			},
			wantMetrics: func(actualComponentName string) []metricdata.Metrics {
				return []metricdata.Metrics{
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Value: 7, // 10 total - 3 rejected
								},
								{
									Attributes: attribute.NewSet(
										semconv.ErrorType(internal.PartialSuccess{
											ErrorMessage:  "some points rejected",
											RejectedItems: 3,
											RejectedKind:  "metric data points",
										}),
										semconv.OTelComponentName(actualComponentName),
										semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Value: 3, // 3 rejected
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
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
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
										semconv.ErrorType(internal.PartialSuccess{
											ErrorMessage:  "some points rejected",
											RejectedItems: 3,
											RejectedKind:  "metric data points",
										}),
										semconv.OTelComponentName(actualComponentName),
										semconv.OTelComponentTypeKey.String("otlp_grpc_metric_exporter"),
										semconv.RPCResponseStatusCode("Unknown"),
										semconv.ServerAddress("localhost"),
										semconv.ServerPort(4317),
									),
									Count: 1,
								},
							},
						},
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

			orig := otel.GetMeterProvider()
			t.Cleanup(func() { otel.SetMeterProvider(orig) })

			dropReaderMetrics := metric.NewView(
				metric.Instrument{
					Scope: instrumentation.Scope{Name: "go.opentelemetry.io/otel/sdk/metric/internal/observ"},
				},
				metric.Stream{Aggregation: metric.AggregationDrop{}},
			)

			reader := metric.NewManualReader()
			provider := metric.NewMeterProvider(
				metric.WithReader(reader),
				metric.WithView(dropReaderMetrics),
			)
			otel.SetMeterProvider(provider)

			em, err := NewInstrumentation(0, "dns:///localhost:4317")
			require.NoError(t, err)
			require.NotNil(t, em)
			rm := createTestResourceMetrics()

			otlpRm, err := transform.ResourceMetrics(rm)
			require.NoError(t, err)
			em.TrackExport(t.Context(), otlpRm).End(tt.err)

			var got metricdata.ResourceMetrics
			err = reader.Collect(t.Context(), &got)
			require.NoError(t, err)
			require.Len(t, got.ScopeMetrics, 1)

			actualComponentName := extractComponentName(got.ScopeMetrics[0])

			want := metricdata.ScopeMetrics{
				Scope: instrumentation.Scope{
					Name:      "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
					Version:   sdk.Version(),
					SchemaURL: semconv.SchemaURL,
				},
				Metrics: tt.wantMetrics(actualComponentName),
			}

			assert.Equal(t, want.Scope, got.ScopeMetrics[0].Scope)
			require.Len(t, got.ScopeMetrics[0].Metrics, len(want.Metrics))
			for i := range want.Metrics {
				opts := []metricdatatest.Option{metricdatatest.IgnoreTimestamp()}
				if strings.Contains(want.Metrics[i].Name, "duration") {
					opts = append(opts, metricdatatest.IgnoreValue())
				}
				metricdatatest.AssertEqual(t, want.Metrics[i], got.ScopeMetrics[0].Metrics[i], opts...)
			}
		})
	}
}

func TestCountProtoDataPoints(t *testing.T) {
	tests := []struct {
		name     string
		rm       *metricpb.ResourceMetrics
		expected int64
	}{
		{
			name:     "nil resource metrics",
			rm:       nil,
			expected: 0,
		},
		{
			name:     "empty resource metrics",
			rm:       &metricpb.ResourceMetrics{},
			expected: 0,
		},
		{
			name: "test data",
			rm: &metricpb.ResourceMetrics{
				ScopeMetrics: []*metricpb.ScopeMetrics{
					{
						Metrics: []*metricpb.Metric{
							{
								Data: &metricpb.Metric_Gauge{
									Gauge: &metricpb.Gauge{
										DataPoints: []*metricpb.NumberDataPoint{
											{}, {},
										},
									},
								},
							},
							{
								Data: &metricpb.Metric_Sum{
									Sum: &metricpb.Sum{
										DataPoints: []*metricpb.NumberDataPoint{
											{},
										},
									},
								},
							},
							{
								Data: &metricpb.Metric_Histogram{
									Histogram: &metricpb.Histogram{
										DataPoints: []*metricpb.HistogramDataPoint{
											{}, {}, {},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: 6, // 2 gauge + 1 sum + 3 histogram
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := countProtoDataPoints(tt.rm)
			assert.Equal(t, tt.expected, count, "data points count mismatch")
		})
	}
}

func TestIsObservabilityEnabled(t *testing.T) {
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
			t.Setenv("OTEL_GO_X_OBSERVABILITY", tt.envValue)

			got := x.Observability.Enabled()
			assert.Equal(t, tt.want, got, "observability enabled state mismatch")
		})
	}
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

func BenchmarkInstrumentationTrackExport(b *testing.B) {
	run := func(enabled bool, err error) func(*testing.B) {
		return func(b *testing.B) {
			testRm := createTestResourceMetrics()
			otlpRm, _ := transform.ResourceMetrics(testRm)
			if enabled {
				b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

				orig := otel.GetMeterProvider()
				b.Cleanup(func() { otel.SetMeterProvider(orig) })

				reader := metric.NewManualReader()
				provider := metric.NewMeterProvider(metric.WithReader(reader))
				otel.SetMeterProvider(provider)
			} else {
				b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
			}
			inst, instErr := NewInstrumentation(0, "dns:///localhost:4317")
			if instErr != nil {
				b.Fatal(instErr)
			}
			if enabled {
				if inst == nil {
					b.Fatal("instrumentation should not be nil when enabled")
				}
			} else if inst != nil {
				b.Fatal("instrumentation should be nil when disabled")
			}
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.TrackExport(b.Context(), otlpRm).End(err)
			}
		}
	}

	b.Run("EnabledNoError", run(true, nil))
	b.Run("EnabledError", run(true, errors.New("export failed")))
	b.Run("Disabled", run(false, nil))
}
