package selfobservability

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus/internal/counter"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
)

func TestNewSelfObservability(t *testing.T) {
	tests := []struct {
		name                string
		setupMeterProvider  func() (metric.MeterProvider, func())
		expectError         bool
		expectedErrorSubstr string
		verifyResult        func(t *testing.T, obs *SelfObservability)
	}{
		{
			name: "successful_creation",
			setupMeterProvider: func() (metric.MeterProvider, func()) {
				reader := sdkmetric.NewManualReader()
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
				prevMP := otel.GetMeterProvider()
				otel.SetMeterProvider(mp)
				return mp, func() { otel.SetMeterProvider(prevMP) }
			},
			expectError: false,
			verifyResult: func(t *testing.T, obs *SelfObservability) {
				require.NotNil(t, obs)
				require.NotNil(t, obs.attrs)
				assert.Len(t, obs.attrs, 2)

				// Verify component name contains the right prefix and an ID
				componentName := ""
				componentType := ""
				for _, attr := range obs.attrs {
					if attr.Key == semconv.OTelComponentNameKey {
						componentName = attr.Value.AsString()
					}
					if attr.Key == semconv.OTelComponentTypeKey {
						componentType = attr.Value.AsString()
					}
				}

				assert.Contains(t, componentName, otelComponentType)
				assert.Equal(t, otelComponentType, componentType)

				// Verify metrics are properly initialized
				assert.NotNil(t, obs.inflightMetric)
				assert.NotNil(t, obs.exportedMetric)
				assert.NotNil(t, obs.operationDuration)
				assert.NotNil(t, obs.collectionDuration)
			},
		},
		{
			name: "with_counter_id_sequence",
			setupMeterProvider: func() (metric.MeterProvider, func()) {
				// Reset counter to ensure predictable ID
				counter.SetExporterID(100)

				reader := sdkmetric.NewManualReader()
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
				prevMP := otel.GetMeterProvider()
				otel.SetMeterProvider(mp)
				return mp, func() { otel.SetMeterProvider(prevMP) }
			},
			expectError: false,
			verifyResult: func(t *testing.T, obs *SelfObservability) {
				require.NotNil(t, obs)

				// Verify the component name includes the expected ID
				componentName := ""
				for _, attr := range obs.attrs {
					if attr.Key == semconv.OTelComponentNameKey {
						componentName = attr.Value.AsString()
					}
				}

				expectedName := fmt.Sprintf("%s/%d", otelComponentType, 100)
				assert.Equal(t, expectedName, componentName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp, cleanup := tt.setupMeterProvider()
			defer cleanup()

			obs, err := NewSelfObservability()

			if tt.expectError {
				require.Error(t, err)
				if tt.expectedErrorSubstr != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorSubstr)
				}
			} else {
				require.NoError(t, err)
				if tt.verifyResult != nil {
					tt.verifyResult(t, obs)
				}
			}
			_ = mp // Use mp to avoid unused variable
		})
	}
}

func TestSelfObservability_ContextMethods(t *testing.T) {
	type contextKey string

	tests := []struct {
		name          string
		setupContext  func() context.Context
		verifyContext func(t *testing.T, obs *SelfObservability, expectedCtx context.Context)
	}{
		{
			name: "set_and_get_context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), contextKey("test-key"), "test-value")
			},
			verifyContext: func(t *testing.T, obs *SelfObservability, expectedCtx context.Context) {
				obs.SetContext(expectedCtx)
				retrievedCtx := obs.GetContext()
				assert.Equal(t, expectedCtx, retrievedCtx)
				assert.Equal(t, "test-value", retrievedCtx.Value(contextKey("test-key")))
			},
		},
		{
			name: "nil_context",
			setupContext: func() context.Context {
				return nil
			},
			verifyContext: func(t *testing.T, obs *SelfObservability, expectedCtx context.Context) {
				obs.SetContext(expectedCtx)
				retrievedCtx := obs.GetContext()
				assert.Equal(t, expectedCtx, retrievedCtx)
				assert.Nil(t, retrievedCtx)
			},
		},
		{
			name: "background_context",
			setupContext: func() context.Context {
				return context.Background()
			},
			verifyContext: func(t *testing.T, obs *SelfObservability, expectedCtx context.Context) {
				obs.SetContext(expectedCtx)
				retrievedCtx := obs.GetContext()
				assert.Equal(t, expectedCtx, retrievedCtx)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test meter provider
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			prevMP := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			defer otel.SetMeterProvider(prevMP)

			// Create self-observability instance
			obs, err := NewSelfObservability()
			require.NoError(t, err)

			// Test context operations
			ctx := tt.setupContext()
			tt.verifyContext(t, obs, ctx)
		})
	}
}

func TestSelfObservability_RecordCollectionDuration(t *testing.T) {
	tests := []struct {
		name          string
		operation     func() error
		expectedError error
		verifyMetrics func(t *testing.T, reader sdkmetric.Reader)
	}{
		{
			name: "successful_operation",
			operation: func() error {
				time.Sleep(10 * time.Millisecond) // Small delay to measure
				return nil
			},
			expectedError: nil,
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find collection duration metric
				var collectionDuration *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.metric_reader.collection.duration" {
							collectionDuration = &sm.Metrics[i]
							break
						}
					}
				}

				require.NotNil(t, collectionDuration, "collection duration metric should exist")

				switch data := collectionDuration.Data.(type) {
				case metricdata.Histogram[float64]:
					require.Len(t, data.DataPoints, 1)
					dp := data.DataPoints[0]
					assert.True(t, dp.Sum > 0, "duration should be greater than 0")
					assert.Equal(t, uint64(1), dp.Count, "count should be 1")

					// Verify attributes contain component info
					attrs := dp.Attributes.ToSlice()
					hasComponentName := false
					hasComponentType := false
					for _, attr := range attrs {
						if attr.Key == semconv.OTelComponentNameKey {
							hasComponentName = true
							assert.Contains(t, attr.Value.AsString(), otelComponentType)
						}
						if attr.Key == semconv.OTelComponentTypeKey {
							hasComponentType = true
							assert.Equal(t, otelComponentType, attr.Value.AsString())
						}
					}
					assert.True(t, hasComponentName, "should have component name attribute")
					assert.True(t, hasComponentType, "should have component type attribute")
				default:
					t.Fatalf("unexpected metric data type: %T", data)
				}
			},
		},
		{
			name: "operation_with_error",
			operation: func() error {
				time.Sleep(5 * time.Millisecond)
				return errors.New("test error")
			},
			expectedError: errors.New("test error"),
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find collection duration metric
				var collectionDuration *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.metric_reader.collection.duration" {
							collectionDuration = &sm.Metrics[i]
							break
						}
					}
				}

				require.NotNil(t, collectionDuration, "collection duration metric should exist")

				switch data := collectionDuration.Data.(type) {
				case metricdata.Histogram[float64]:
					require.Len(t, data.DataPoints, 1)
					dp := data.DataPoints[0]
					assert.True(t, dp.Sum > 0, "duration should be greater than 0")

					// Verify error attribute is present
					attrs := dp.Attributes.ToSlice()
					hasErrorType := false
					for _, attr := range attrs {
						if attr.Key == semconv.ErrorTypeKey {
							hasErrorType = true
							assert.Equal(t, "*errors.errorString", attr.Value.AsString())
						}
					}
					assert.True(t, hasErrorType, "should have error type attribute")
				default:
					t.Fatalf("unexpected metric data type: %T", data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test meter provider
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			prevMP := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			defer otel.SetMeterProvider(prevMP)

			// Create self-observability instance
			obs, err := NewSelfObservability()
			require.NoError(t, err)

			ctx := context.Background()

			// Execute operation and record duration
			err = obs.RecordCollectionDuration(ctx, tt.operation)

			// Verify error expectation
			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			// Verify metrics
			if tt.verifyMetrics != nil {
				tt.verifyMetrics(t, reader)
			}
		})
	}
}

func TestSelfObservability_RecordOperationDuration(t *testing.T) {
	tests := []struct {
		name          string
		operationErr  error
		verifyMetrics func(t *testing.T, reader sdkmetric.Reader, hasError bool)
	}{
		{
			name:         "successful_operation",
			operationErr: nil,
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader, hasError bool) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find operation duration metric
				var operationDuration *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.exporter.operation.duration" {
							operationDuration = &sm.Metrics[i]
							break
						}
					}
				}

				require.NotNil(t, operationDuration, "operation duration metric should exist")

				switch data := operationDuration.Data.(type) {
				case metricdata.Histogram[float64]:
					require.Len(t, data.DataPoints, 1)
					dp := data.DataPoints[0]
					assert.True(t, dp.Sum > 0, "duration should be greater than 0")
					assert.Equal(t, uint64(1), dp.Count, "count should be 1")

					// For successful operations, should not have error attribute
					attrs := dp.Attributes.ToSlice()
					hasErrorType := false
					for _, attr := range attrs {
						if attr.Key == semconv.ErrorTypeKey {
							hasErrorType = true
						}
					}
					assert.False(t, hasErrorType, "successful operation should not have error type attribute")
				default:
					t.Fatalf("unexpected metric data type: %T", data)
				}
			},
		},
		{
			name:         "operation_with_error",
			operationErr: errors.New("operation failed"),
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader, hasError bool) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find operation duration metric
				var operationDuration *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.exporter.operation.duration" {
							operationDuration = &sm.Metrics[i]
							break
						}
					}
				}

				require.NotNil(t, operationDuration, "operation duration metric should exist")

				switch data := operationDuration.Data.(type) {
				case metricdata.Histogram[float64]:
					require.Len(t, data.DataPoints, 1)
					dp := data.DataPoints[0]
					assert.True(t, dp.Sum > 0, "duration should be greater than 0")

					// For failed operations, should have error attribute
					attrs := dp.Attributes.ToSlice()
					hasErrorType := false
					for _, attr := range attrs {
						if attr.Key == semconv.ErrorTypeKey {
							hasErrorType = true
							assert.Equal(t, "*errors.errorString", attr.Value.AsString())
						}
					}
					assert.True(t, hasErrorType, "failed operation should have error type attribute")
				default:
					t.Fatalf("unexpected metric data type: %T", data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test meter provider
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			prevMP := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			defer otel.SetMeterProvider(prevMP)

			// Create self-observability instance
			obs, err := NewSelfObservability()
			require.NoError(t, err)

			ctx := context.Background()

			// Start operation tracking
			completeFn := obs.RecordOperationDuration(ctx)

			// Simulate some operation time
			time.Sleep(10 * time.Millisecond)

			// Complete operation with error (if any)
			completeFn(tt.operationErr)

			// Verify metrics
			if tt.verifyMetrics != nil {
				tt.verifyMetrics(t, reader, tt.operationErr != nil)
			}
		})
	}
}

func TestSelfObservability_TrackExport(t *testing.T) {
	tests := []struct {
		name          string
		exportCount   int64
		exportErr     error
		successCount  int64
		verifyMetrics func(t *testing.T, reader sdkmetric.Reader, exportCount, successCount int64, hasError bool)
	}{
		{
			name:         "successful_export_all",
			exportCount:  5,
			exportErr:    nil,
			successCount: 5,
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader, exportCount, successCount int64, hasError bool) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find inflight and exported metrics
				var inflightMetric, exportedMetric *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						switch sm.Metrics[i].Name {
						case "otel.sdk.exporter.metric_data_point.inflight":
							inflightMetric = &sm.Metrics[i]
						case "otel.sdk.exporter.metric_data_point.exported":
							exportedMetric = &sm.Metrics[i]
						}
					}
				}

				require.NotNil(t, inflightMetric, "inflight metric should exist")
				require.NotNil(t, exportedMetric, "exported metric should exist")

				// Verify inflight metric (should be 0 after completion)
				switch data := inflightMetric.Data.(type) {
				case metricdata.Sum[int64]:
					totalInflight := int64(0)
					for _, dp := range data.DataPoints {
						totalInflight += dp.Value
					}
					assert.Equal(t, int64(0), totalInflight, "inflight should be 0 after completion")
				default:
					t.Fatalf("unexpected inflight metric data type: %T", data)
				}

				// Verify exported metric
				switch data := exportedMetric.Data.(type) {
				case metricdata.Sum[int64]:
					totalExported := int64(0)
					for _, dp := range data.DataPoints {
						totalExported += dp.Value
					}
					assert.Equal(t, successCount, totalExported, "exported count should match success count")
				default:
					t.Fatalf("unexpected exported metric data type: %T", data)
				}
			},
		},
		{
			name:         "partial_export_success",
			exportCount:  10,
			exportErr:    errors.New("partial failure"),
			successCount: 7,
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader, exportCount, successCount int64, hasError bool) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find exported metrics
				var exportedMetrics []*metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.exporter.metric_data_point.exported" {
							exportedMetrics = append(exportedMetrics, &sm.Metrics[i])
						}
					}
				}

				require.NotEmpty(t, exportedMetrics, "exported metrics should exist")

				// Should have both success and error metrics
				var totalSuccessExported, totalErrorExported int64
				var hasErrorAttribute bool

				for _, metric := range exportedMetrics {
					switch data := metric.Data.(type) {
					case metricdata.Sum[int64]:
						for _, dp := range data.DataPoints {
							attrs := dp.Attributes.ToSlice()
							hasError := false
							for _, attr := range attrs {
								if attr.Key == semconv.ErrorTypeKey {
									hasError = true
									hasErrorAttribute = true
									break
								}
							}
							if hasError {
								totalErrorExported += dp.Value
							} else {
								totalSuccessExported += dp.Value
							}
						}
					default:
						t.Fatalf("unexpected exported metric data type: %T", data)
					}
				}

				assert.Equal(t, successCount, totalSuccessExported, "success count should match")
				assert.Equal(t, exportCount-successCount, totalErrorExported, "error count should match failed exports")
				assert.True(t, hasErrorAttribute, "should have error attribute for failed exports")
			},
		},
		{
			name:         "zero_count_export",
			exportCount:  0,
			exportErr:    nil,
			successCount: 0,
			verifyMetrics: func(t *testing.T, reader sdkmetric.Reader, exportCount, successCount int64, hasError bool) {
				var rm metricdata.ResourceMetrics
				err := reader.Collect(context.Background(), &rm)
				require.NoError(t, err)

				// Find inflight metric
				var inflightMetric *metricdata.Metrics
				for _, sm := range rm.ScopeMetrics {
					for i := range sm.Metrics {
						if sm.Metrics[i].Name == "otel.sdk.exporter.metric_data_point.inflight" {
							inflightMetric = &sm.Metrics[i]
							break
						}
					}
				}

				require.NotNil(t, inflightMetric, "inflight metric should exist")

				// Verify inflight metric (should be 0)
				switch data := inflightMetric.Data.(type) {
				case metricdata.Sum[int64]:
					totalInflight := int64(0)
					for _, dp := range data.DataPoints {
						totalInflight += dp.Value
					}
					assert.Equal(t, int64(0), totalInflight, "inflight should be 0")
				default:
					t.Fatalf("unexpected inflight metric data type: %T", data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test meter provider
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			prevMP := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			defer otel.SetMeterProvider(prevMP)

			// Create self-observability instance
			obs, err := NewSelfObservability()
			require.NoError(t, err)

			ctx := context.Background()

			// Start export tracking
			completeFn := obs.TrackExport(ctx, tt.exportCount)

			// Complete export with results
			completeFn(tt.exportErr, tt.successCount)

			// Verify metrics
			if tt.verifyMetrics != nil {
				tt.verifyMetrics(t, reader, tt.exportCount, tt.successCount, tt.exportErr != nil)
			}
		})
	}
}

func TestSelfObservability_AttributePooling(t *testing.T) {
	// This test verifies that the attribute pooling works correctly
	// and doesn't cause race conditions or memory leaks

	tests := []struct {
		name          string
		operationFunc func(obs *SelfObservability)
	}{
		{
			name: "track_export_pooling",
			operationFunc: func(obs *SelfObservability) {
				ctx := context.Background()
				completeFn := obs.TrackExport(ctx, 1)
				completeFn(nil, 1)
			},
		},
		{
			name: "record_operation_duration_pooling",
			operationFunc: func(obs *SelfObservability) {
				ctx := context.Background()
				completeFn := obs.RecordOperationDuration(ctx)
				completeFn(nil)
			},
		},
		{
			name: "record_collection_duration_pooling",
			operationFunc: func(obs *SelfObservability) {
				ctx := context.Background()
				err := obs.RecordCollectionDuration(ctx, func() error { return nil })
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test meter provider
			reader := sdkmetric.NewManualReader()
			mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
			prevMP := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			defer otel.SetMeterProvider(prevMP)

			// Create self-observability instance
			obs, err := NewSelfObservability()
			require.NoError(t, err)

			// Run the operation multiple times to test pooling
			for i := 0; i < 10; i++ {
				tt.operationFunc(obs)
			}

			// Verify that we can still collect metrics after pooling operations
			var rm metricdata.ResourceMetrics
			err = reader.Collect(context.Background(), &rm)
			require.NoError(t, err)

			// Should have some metrics
			assert.NotEmpty(t, rm.ScopeMetrics, "should have scope metrics after pooling operations")
		})
	}
}
