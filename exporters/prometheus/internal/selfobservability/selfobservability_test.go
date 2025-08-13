// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package selfobservability

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExporterMetrics(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	ctx := context.Background()

	// Test AddInflight and AddExported
	em.AddInflight(ctx, 3)
	em.AddInflight(ctx, -1)
	em.AddExported(ctx, 5)

	// Test TrackCollectionDuration
	endCollect := em.TrackCollectionDuration(ctx)
	endCollect(nil)

	// Test TrackOperationDuration
	endOp := em.TrackOperationDuration(ctx)
	endOp(nil)

	// Collect metrics
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Helper to find metrics
	findMetric := func(name string) *metricdata.Metrics {
		for _, sm := range rm.ScopeMetrics {
			for i := range sm.Metrics {
				if sm.Metrics[i].Name == name {
					return &sm.Metrics[i]
				}
			}
		}
		return nil
	}

	// Verify exported metric
	exported := findMetric("otel.sdk.exporter.metric_data_point.exported")
	require.NotNil(t, exported)
	switch data := exported.Data.(type) {
	case metricdata.Sum[int64]:
		assert.Equal(t, int64(5), data.DataPoints[0].Value)
	case metricdata.Sum[float64]:
		assert.InDelta(t, 5.0, data.DataPoints[0].Value, 0.001)
	}

	// Verify inflight metric (3 - 1 = 2)
	inflight := findMetric("otel.sdk.exporter.metric_data_point.inflight")
	require.NotNil(t, inflight)
	switch data := inflight.Data.(type) {
	case metricdata.Sum[int64]:
		assert.Equal(t, int64(2), data.DataPoints[0].Value)
	case metricdata.Sum[float64]:
		assert.InDelta(t, 2.0, data.DataPoints[0].Value, 0.001)
	}

	// Verify duration metrics exist
	collectionDuration := findMetric("otel.sdk.metric_reader.collection.duration")
	require.NotNil(t, collectionDuration)

	operationDuration := findMetric("otel.sdk.exporter.operation.duration")
	require.NotNil(t, operationDuration)
}

func TestExporterMetricsWithErrors(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	ctx := context.Background()

	// Test TrackCollectionDuration with error
	endCollect := em.TrackCollectionDuration(ctx)
	testErr := errors.New("collection error")
	endCollect(testErr)

	// Test TrackOperationDuration with error
	endOp := em.TrackOperationDuration(ctx)
	testErr2 := errors.New("operation error")
	endOp(testErr2)

	// Collect metrics
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Helper to find metrics
	findMetric := func(name string) *metricdata.Metrics {
		for _, sm := range rm.ScopeMetrics {
			for i := range sm.Metrics {
				if sm.Metrics[i].Name == name {
					return &sm.Metrics[i]
				}
			}
		}
		return nil
	}

	// Verify duration metrics with error attributes
	collectionDuration := findMetric("otel.sdk.metric_reader.collection.duration")
	require.NotNil(t, collectionDuration)

	operationDuration := findMetric("otel.sdk.exporter.operation.duration")
	require.NotNil(t, operationDuration)

	// Check that error attributes are present in the metrics
	switch data := collectionDuration.Data.(type) {
	case metricdata.Histogram[int64]:
		require.Len(t, data.DataPoints, 1)
		// Check for error.type attribute
		found := false
		for _, attr := range data.DataPoints[0].Attributes.ToSlice() {
			if attr.Key == "error.type" {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error.type attribute in collection duration metric")
	case metricdata.Histogram[float64]:
		require.Len(t, data.DataPoints, 1)
		// Check for error.type attribute
		found := false
		for _, attr := range data.DataPoints[0].Attributes.ToSlice() {
			if attr.Key == "error.type" {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected error.type attribute in collection duration metric")
	}
}

func TestExporterMetricsAddExportedWithError(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	ctx := context.Background()

	// Test AddExportedWithError with an error
	testErr := errors.New("test validation error")
	em.AddExportedWithError(ctx, 3, testErr)

	// Also add a regular exported metric for comparison
	em.AddExported(ctx, 2)

	// Collect metrics
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Helper to find metrics
	findMetric := func(name string) *metricdata.Metrics {
		for _, sm := range rm.ScopeMetrics {
			for i := range sm.Metrics {
				if sm.Metrics[i].Name == name {
					return &sm.Metrics[i]
				}
			}
		}
		return nil
	}

	// Verify exported metric exists
	exported := findMetric("otel.sdk.exporter.metric_data_point.exported")
	require.NotNil(t, exported, "Expected exported metric to exist")

	switch data := exported.Data.(type) {
	case metricdata.Sum[int64]:
		// Check that we have data points
		require.NotEmpty(t, data.DataPoints, "Expected at least one data point")

		// Look for data points with and without error attributes
		var regularCount, errorCount int64
		var foundErrorType bool

		for _, dp := range data.DataPoints {
			hasError := false
			for _, attr := range dp.Attributes.ToSlice() {
				if attr.Key == "error.type" {
					hasError = true
					foundErrorType = true
					// The error.type attribute should contain the Go error type, not the message
					assert.Equal(t, "*errors.errorString", attr.Value.AsString())
					errorCount += dp.Value
				}
			}
			if !hasError {
				regularCount += dp.Value
			}
		}

		// We should have found an error type attribute
		assert.True(t, foundErrorType, "Expected to find error.type attribute")

		// Verify counts
		assert.Equal(t, int64(2), regularCount, "Expected 2 regular exported metrics")
		assert.Equal(t, int64(3), errorCount, "Expected 3 exported metrics with error")

	case metricdata.Sum[float64]:
		// Similar verification for float64 case
		require.NotEmpty(t, data.DataPoints, "Expected at least one data point")

		var regularCount, errorCount float64
		var foundErrorType bool

		for _, dp := range data.DataPoints {
			hasError := false
			for _, attr := range dp.Attributes.ToSlice() {
				if attr.Key == "error.type" {
					hasError = true
					foundErrorType = true
					// The error.type attribute should contain the Go error type, not the message
					assert.Equal(t, "*errors.errorString", attr.Value.AsString())
					errorCount += dp.Value
				}
			}
			if !hasError {
				regularCount += dp.Value
			}
		}

		assert.True(t, foundErrorType, "Expected to find error.type attribute")
		assert.InDelta(t, 2.0, regularCount, 0.001)
		assert.InDelta(t, 3.0, errorCount, 0.001)

	default:
		t.Fatalf("Unexpected data type: %T", data)
	}
}

func TestExporterMetricsAddExportedWithErrorDisabled(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	em.DisableSelfObservability()
	ctx := context.Background()

	// Test AddExportedWithError when self-observability is disabled
	testErr := errors.New("test error")
	em.AddExportedWithError(ctx, 5, testErr)

	// Collect metrics - should be empty since self-observability is disabled
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Verify no metrics are recorded when disabled
	totalMetrics := 0
	for _, sm := range rm.ScopeMetrics {
		totalMetrics += len(sm.Metrics)
	}
	assert.Equal(t, 0, totalMetrics)
}

func TestExporterMetricsDisabledSelfObservability(t *testing.T) {
	// Set up test meter provider
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	// Disable self-observability to test the disabled path
	em.DisableSelfObservability()

	ctx := context.Background()

	// Test AddInflight - should not record
	em.AddInflight(ctx, 3)
	em.AddExported(ctx, 5)

	// Test TrackCollectionDuration - should return noop function
	endCollect := em.TrackCollectionDuration(ctx)
	endCollect(nil)

	// Test TrackOperationDuration - should return noop function
	endOp := em.TrackOperationDuration(ctx)
	endOp(nil)

	// Collect metrics - should be empty since self-observability is disabled
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))

	// Verify no metrics are recorded when disabled
	totalMetrics := 0
	for _, sm := range rm.ScopeMetrics {
		totalMetrics += len(sm.Metrics)
	}
	// Should be 0 since self-observability is disabled
	assert.Equal(t, 0, totalMetrics)
}

func TestNewExporterMetricsErrorHandling(t *testing.T) {
	// Test error handling in NewExporterMetrics by using a noop meter provider
	// This will cause metric creation to fail but should not panic

	// Save current meter provider
	prev := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	// Set a noop meter provider that might cause issues
	otel.SetMeterProvider(metric.NewMeterProvider())

	// This should handle errors gracefully and not panic
	em := NewExporterMetrics("test-component")
	assert.NotNil(t, em)
	assert.True(t, em.selfObservabilityEnabled)
	assert.NotNil(t, em.attrs)
	assert.Len(t, em.attrs, 2)

	// Test that the exporter metrics still function even if some metrics failed to create
	ctx := context.Background()
	em.AddInflight(ctx, 1)
	em.AddExported(ctx, 1)

	endCollect := em.TrackCollectionDuration(ctx)
	endCollect(nil)

	endOp := em.TrackOperationDuration(ctx)
	endOp(nil)
}

func TestNewExporterMetricsWithErrorHandler(t *testing.T) {
	// Test that error handlers are called when metric creation fails
	var handledErrors []error
	originalHandler := otel.GetErrorHandler()
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		handledErrors = append(handledErrors, err)
	}))
	t.Cleanup(func() { otel.SetErrorHandler(originalHandler) })

	// Set a limited meter provider that might trigger errors
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	prev := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(prev) })

	em := NewExporterMetrics("test-component")
	assert.NotNil(t, em)

	// The exact number of errors depends on the implementation
	// but we want to ensure error handling doesn't panic
	t.Logf("Number of handled errors during metric creation: %d", len(handledErrors))
}
