// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package selfobservability

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type testSetup struct {
	reader *sdkmetric.ManualReader
	mp     *sdkmetric.MeterProvider
	ctx    context.Context
	em     *ExporterMetrics
}

func setupTestMeterProvider(t *testing.T) *testSetup {
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	originalMP := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(originalMP) })

	componentName := semconv.OTelComponentName("test")
	componentType := semconv.OTelComponentTypeKey.String("exporter")
	em := NewExporterMetrics("go.opentelemetry.io/otel/exporters/stdout/stdoutmetric", componentName, componentType)

	return &testSetup{
		reader: reader,
		mp:     mp,
		ctx:    context.Background(),
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

func TestExporterMetrics_TrackExport(t *testing.T) {
	setup := setupTestMeterProvider(t)

	done1 := setup.em.TrackExport(setup.ctx, 2)
	done2 := setup.em.TrackExport(setup.ctx, 3)
	done3 := setup.em.TrackExport(setup.ctx, 1)
	time.Sleep(5 * time.Millisecond)
	done2(nil)
	done1(errors.New("failed"))
	done3(nil)

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

func TestExporterMetrics_TrackExport_WithError(t *testing.T) {
	setup := setupTestMeterProvider(t)
	count := int64(3)
	testErr := errors.New("export failed")

	done := setup.em.TrackExport(setup.ctx, count)
	done(testErr)

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

func TestExporterMetrics_TrackExport_InflightTracking(t *testing.T) {
	setup := setupTestMeterProvider(t)
	count := int64(10)

	done := setup.em.TrackExport(setup.ctx, count)
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

	done(nil)
	rm = collectMetrics(t, setup)
	inflight, found = findMetric(rm, otelconv.SDKExporterMetricDataPointInflight{}.Name())
	assert.True(t, found)
	if sum, ok := inflight.Data.(metricdata.Sum[int64]); ok {
		for _, dp := range sum.DataPoints {
			assert.Equal(t, int64(0), dp.Value)
		}
	}
}

func TestExporterMetrics_AttributesNotPermanentlyModified(t *testing.T) {
	componentName := semconv.OTelComponentName("test-component")
	componentType := semconv.OTelComponentTypeKey.String("test-exporter")
	em := NewExporterMetrics("test", componentName, componentType)

	assert.Len(t, em.attrs, 2)
	assert.Contains(t, em.attrs, componentName)
	assert.Contains(t, em.attrs, componentType)

	done := em.TrackExport(context.Background(), 1)
	done(errors.New("test error"))
	done = em.TrackExport(context.Background(), 1)
	done(nil)

	assert.Len(t, em.attrs, 2)
	assert.Contains(t, em.attrs, componentName)
	assert.Contains(t, em.attrs, componentType)
}
