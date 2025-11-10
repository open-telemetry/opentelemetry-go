// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/observ"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

type testSetup struct {
	reader *sdkmetric.ManualReader
	em     *observ.Instrumentation
}

func setupTestMeterProvider(t *testing.T) *testSetup {
	t.Helper()
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// drop metric reader metrics as we are only testing for stdoutmetric exporter
	dropReaderMetrics := sdkmetric.NewView(
		sdkmetric.Instrument{
			Scope: instrumentation.Scope{Name: "go.opentelemetry.io/otel/sdk/metric/internal/observ"},
		},
		sdkmetric.Stream{Aggregation: sdkmetric.AggregationDrop{}},
	)

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(dropReaderMetrics),
	)

	originalMP := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(originalMP) })

	em, err := observ.NewInstrumentation(0)
	assert.NoError(t, err)

	return &testSetup{
		reader: reader,
		em:     em,
	}
}

func collectMetrics(t *testing.T, setup *testSetup) metricdata.ResourceMetrics {
	var rm metricdata.ResourceMetrics
	err := setup.reader.Collect(t.Context(), &rm)
	assert.NoError(t, err)
	return rm
}

const exporterComponentID = 0

var errExport = errors.New("export failed")

func exporterSet(attrs ...attribute.KeyValue) attribute.Set {
	return attribute.NewSet(append([]attribute.KeyValue{
		observ.ExporterComponentName(exporterComponentID),
		semconv.OTelComponentTypeKey.String(observ.ComponentType),
	}, attrs...)...)
}

func dPt(set attribute.Set, value int64) metricdata.DataPoint[int64] {
	return metricdata.DataPoint[int64]{Attributes: set, Value: value}
}

func exported(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
		Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
		Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dPts,
		},
	}
}

func inflight(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
		Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
		Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints:  dPts,
		},
	}
}

func duration(dPts ...metricdata.HistogramDataPoint[float64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterOperationDuration{}.Name(),
		Description: otelconv.SDKExporterOperationDuration{}.Description(),
		Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  dPts,
		},
	}
}

func histDPt(set attribute.Set) metricdata.HistogramDataPoint[float64] {
	return metricdata.HistogramDataPoint[float64]{
		Attributes: set,
	}
}

func checkMetrics(
	t *testing.T,
	rm metricdata.ResourceMetrics,
	wantInflight, wantExported, wantDuration metricdata.Metrics,
) {
	t.Helper()
	require.Len(t, rm.ScopeMetrics, 1)

	m := rm.ScopeMetrics[0].Metrics
	require.Len(t, m, 3)

	opts := metricdatatest.IgnoreTimestamp()

	metricdatatest.AssertEqual(t, wantInflight, m[0], opts)
	metricdatatest.AssertEqual(t, wantExported, m[1], opts)
	// ignoring values for histogram since duration is not deterministic
	metricdatatest.AssertEqual(t, wantDuration, m[2], opts, metricdatatest.IgnoreValue())
}

func checkInflight(t *testing.T, rm metricdata.ResourceMetrics, wantInflight metricdata.Metrics) {
	t.Helper()
	require.Len(t, rm.ScopeMetrics, 1)

	m := rm.ScopeMetrics[0].Metrics
	require.NotEmpty(t, m)

	inflightName := otelconv.SDKExporterMetricDataPointInflight{}.Name()
	var inflightMetric metricdata.Metrics
	found := false
	for _, metric := range m {
		if metric.Name == inflightName {
			inflightMetric = metric
			found = true
			break
		}
	}
	require.True(t, found)
	metricdatatest.AssertEqual(t, wantInflight, inflightMetric, metricdatatest.IgnoreTimestamp())
}

func TestInstrumentationExportMetrics(t *testing.T) {
	setup := setupTestMeterProvider(t)

	ctx := t.Context()
	op1 := setup.em.ExportMetrics(ctx, 2)
	op2 := setup.em.ExportMetrics(ctx, 3)
	op3 := setup.em.ExportMetrics(ctx, 1)
	totalMetrics := int64(6)
	checkInflight(t, collectMetrics(t, setup), inflight(dPt(exporterSet(), totalMetrics)))

	op2.End(nil)
	op1.End(errExport)
	op3.End(nil)

	successExported := int64(4)
	erroredExported := int64(2)
	checkMetrics(
		t,
		collectMetrics(t, setup),
		inflight(dPt(exporterSet(), 0)),
		exported(
			dPt(exporterSet(), successExported),
			dPt(exporterSet(semconv.ErrorType(errExport)), erroredExported),
		),
		duration(
			histDPt(exporterSet()),
			histDPt(exporterSet(semconv.ErrorType(errExport))),
		),
	)
}

func BenchmarkExportMetrics(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	orig := otel.GetMeterProvider()
	b.Cleanup(func() {
		otel.SetMeterProvider(orig)
	})

	// Ensure deterministic benchmark by using noop meter.
	otel.SetMeterProvider(noop.NewMeterProvider())

	newExp := func(b *testing.B) *observ.Instrumentation {
		b.Helper()
		em, err := observ.NewInstrumentation(0)
		require.NoError(b, err)
		require.NotNil(b, em)
		return em
	}

	b.Run("NoError", func(b *testing.B) {
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
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				op := em.ExportMetrics(b.Context(), 10)
				op.End(errExport)
			}
		})
	})
}
