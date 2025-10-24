// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

type testSetup struct {
	reader *sdkmetric.ManualReader
	mp     *sdkmetric.MeterProvider
	ctx    context.Context
	em     *Instrumentation
}

func setupTestMeterProvider(t *testing.T) *testSetup {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	originalMP := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(originalMP) })

	em, err := NewInstrumentation(0)
	assert.NoError(t, err)

	return &testSetup{
		reader: reader,
		mp:     mp,
		ctx:    t.Context(),
		em:     em,
	}
}

func collectMetrics(t *testing.T, setup *testSetup) metricdata.ResourceMetrics {
	var rm metricdata.ResourceMetrics
	err := setup.reader.Collect(setup.ctx, &rm)
	assert.NoError(t, err)
	return rm
}

func findMetric(rm metricdata.ResourceMetrics, name string) (metricdata.Metrics, bool) {
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name == name {
				return m, true
			}
		}
	}
	return metricdata.Metrics{}, false
}

func TestInstrumentationExportMetrics(t *testing.T) {
	setup := setupTestMeterProvider(t)

	op1 := setup.em.ExportMetrics(setup.ctx, 2)
	op2 := setup.em.ExportMetrics(setup.ctx, 3)
	op3 := setup.em.ExportMetrics(setup.ctx, 1)
	time.Sleep(5 * time.Millisecond)
	op2.End(nil)
	op1.End(errors.New("failed"))
	op3.End(nil)

	rm := collectMetrics(t, setup)
	assert.NotEmpty(t, rm.ScopeMetrics)

	inflight, found := findMetric(rm, otelconv.SDKExporterMetricDataPointInflight{}.Name())
	assert.True(t, found)
	var totalInflightValue int64
	if sum, ok := inflight.Data.(metricdata.Sum[int64]); ok {
		for _, dp := range sum.DataPoints {
			totalInflightValue += dp.Value
		}
	}

	exported, found := findMetric(rm, otelconv.SDKExporterMetricDataPointExported{}.Name())
	assert.True(t, found)
	var totalExported int64
	if sum, ok := exported.Data.(metricdata.Sum[int64]); ok {
		for _, dp := range sum.DataPoints {
			totalExported += dp.Value
		}
	}

	duration, found := findMetric(rm, otelconv.SDKExporterOperationDuration{}.Name())
	assert.True(t, found)
	var operationCount uint64
	if hist, ok := duration.Data.(metricdata.Histogram[float64]); ok {
		for _, dp := range hist.DataPoints {
			operationCount += dp.Count
			assert.Positive(t, dp.Sum)
		}
	}

	assert.Equal(t, int64(6), totalExported)
	assert.Equal(t, uint64(3), operationCount)
	assert.Equal(t, int64(0), totalInflightValue)
}

func TestInstrumentationExportMetricsWithError(t *testing.T) {
	setup := setupTestMeterProvider(t)
	count := int64(3)
	testErr := errors.New("export failed")

	op := setup.em.ExportMetrics(setup.ctx, count)
	op.End(testErr)

	rm := collectMetrics(t, setup)
	assert.NotEmpty(t, rm.ScopeMetrics)

	exported, found := findMetric(rm, otelconv.SDKExporterMetricDataPointExported{}.Name())
	assert.True(t, found)
	if sum, ok := exported.Data.(metricdata.Sum[int64]); ok {
		attr, hasErrorAttr := sum.DataPoints[0].Attributes.Value(semconv.ErrorTypeKey)
		assert.True(t, hasErrorAttr)
		assert.Equal(t, "*errors.errorString", attr.AsString())
	}
}

func TestInstrumentationExportMetricsInflightTracking(t *testing.T) {
	setup := setupTestMeterProvider(t)
	count := int64(10)

	op := setup.em.ExportMetrics(setup.ctx, count)
	rm := collectMetrics(t, setup)
	inflight, found := findMetric(rm, otelconv.SDKExporterMetricDataPointInflight{}.Name())
	assert.True(t, found)

	var inflightValue int64
	if sum, ok := inflight.Data.(metricdata.Sum[int64]); ok {
		for _, dp := range sum.DataPoints {
			inflightValue = dp.Value
		}
	}
	assert.Equal(t, count, inflightValue)

	op.End(nil)
	rm = collectMetrics(t, setup)
	inflight, found = findMetric(rm, otelconv.SDKExporterMetricDataPointInflight{}.Name())
	assert.True(t, found)
	if sum, ok := inflight.Data.(metricdata.Sum[int64]); ok {
		for _, dp := range sum.DataPoints {
			assert.Equal(t, int64(0), dp.Value)
		}
	}
}

func TestInstrumentationExportMetricsAttributesNotPermanentlyModified(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	em, err := NewInstrumentation(42)
	assert.NoError(t, err)

	// Should have component.name and component.type attributes
	assert.Len(t, em.attrs, 2)
	expectedComponentName := semconv.OTelComponentName(
		"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric.exporter/42",
	)
	expectedComponentType := semconv.OTelComponentTypeKey.String(componentType)
	assert.Contains(t, em.attrs, expectedComponentName)
	assert.Contains(t, em.attrs, expectedComponentType)

	op := em.ExportMetrics(t.Context(), 1)
	op.End(errors.New("test error"))
	op = em.ExportMetrics(t.Context(), 1)
	op.End(nil)

	// Attributes should not be modified after tracking exports
	assert.Len(t, em.attrs, 2)
	assert.Contains(t, em.attrs, expectedComponentName)
	assert.Contains(t, em.attrs, expectedComponentType)
}

func BenchmarkExportMetrics(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	orig := otel.GetMeterProvider()
	b.Cleanup(func() {
		otel.SetMeterProvider(orig)
	})

	// Ensure deterministic benchmark by using noop meter.
	otel.SetMeterProvider(noop.NewMeterProvider())

	newExp := func(b *testing.B) *Instrumentation {
		b.Helper()
		em, err := NewInstrumentation(0)
		require.NoError(b, err)
		require.NotNil(b, em)
		return em
	}

	b.Run("Success", func(b *testing.B) {
		em := newExp(b)
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				op := em.ExportMetrics(b.Context(), 10)
				op.End(nil)
			}
		})
	})

	b.Run("WithError", func(b *testing.B) {
		em := newExp(b)
		testErr := errors.New("export failed")
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				op := em.ExportMetrics(b.Context(), 10)
				op.End(testErr)
			}
		})
	})
}
