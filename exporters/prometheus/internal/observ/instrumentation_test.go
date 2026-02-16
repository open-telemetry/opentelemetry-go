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
	"go.opentelemetry.io/otel/exporters/prometheus/internal/observ"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

const ID = 0

var Scope = instrumentation.Scope{
	Name:      observ.ScopeName,
	Version:   observ.Version,
	SchemaURL: observ.SchemaURL,
}

type errMeterProvider struct {
	mapi.MeterProvider
	err error
}

func (m *errMeterProvider) Meter(string, ...mapi.MeterOption) mapi.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	mapi.Meter
	err error
}

func (m *errMeter) Int64UpDownCounter(string, ...mapi.Int64UpDownCounterOption) (mapi.Int64UpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) Int64Counter(string, ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Float64Histogram(string, ...mapi.Float64HistogramOption) (mapi.Float64Histogram, error) {
	return nil, m.err
}

func TestNewInstrumentationObservabilityErrors(t *testing.T) {
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	_, err := observ.NewInstrumentation(ID)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors should be joined")

	assert.ErrorContains(t, err, "inflight metric")
	assert.ErrorContains(t, err, "exported metric")
	assert.ErrorContains(t, err, "operation duration metric")
	assert.ErrorContains(t, err, "collection duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := observ.NewInstrumentation(ID)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

// setup installs a ManualReader MeterProvider and returns an instantiated
// Instrumentation plus a collector that returns the single ScopeMetrics group.
func setup(t *testing.T) (*observ.Instrumentation, func() metricdata.ScopeMetrics) {
	t.Helper()

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	original := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(original) })

	r := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(r))
	otel.SetMeterProvider(mp)

	inst, err := observ.NewInstrumentation(ID)
	require.NoError(t, err)
	require.NotNil(t, inst)

	return inst, func() metricdata.ScopeMetrics {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

func set(err error) attribute.Set {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(observ.ComponentName(ID)),
		semconv.OTelComponentTypeKey.String(observ.ComponentType),
	}
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attribute.NewSet(attrs...)
}

func scrapeInflight() metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
		Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
		Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: set(nil), Value: 0},
			},
		},
	}
}

func scrapeExported(success, total int64, err error) metricdata.Metrics {
	dps := []metricdata.DataPoint[int64]{
		{Attributes: set(nil), Value: success},
	}
	if err != nil {
		dps = append(dps, metricdata.DataPoint[int64]{
			Attributes: set(err),
			Value:      total - success,
		})
	}
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
		Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
		Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dps,
		},
	}
}

func operationDuration(err error) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterOperationDuration{}.Name(),
		Description: otelconv.SDKExporterOperationDuration{}.Description(),
		Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: set(err)},
			},
		},
	}
}

func collectionDuration(err error) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKMetricReaderCollectionDuration{}.Name(),
		Description: otelconv.SDKMetricReaderCollectionDuration{}.Description(),
		Unit:        otelconv.SDKMetricReaderCollectionDuration{}.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: set(err)},
			},
		},
	}
}

func assertExportMetricsMetrics(t *testing.T, got metricdata.ScopeMetrics, total, success int64, err error) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics (inflight, exported, operationDuration)")

	o := metricdatatest.IgnoreTimestamp()

	want := scrapeInflight()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = scrapeExported(success, total, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = operationDuration(err)
	metricdatatest.AssertEqual(t, want, m[2], o, metricdatatest.IgnoreValue())
}

func assertCollectionOnly(t *testing.T, got metricdata.ScopeMetrics, err error) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 1, "expected only collectionDuration metric")

	o := metricdatatest.IgnoreTimestamp()
	want := collectionDuration(err)
	metricdatatest.AssertEqual(t, want, m[0], o, metricdatatest.IgnoreValue())
}

func TestInstrumentationExportMetricsSuccess(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	timer := inst.RecordOperationDuration(t.Context())
	inst.ExportMetrics(t.Context(), n).End(n, nil)
	timer.Stop(nil)

	assertExportMetricsMetrics(t, collect(), n, n, nil)
}

func TestInstrumentationExportMetricsAllErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	err := assert.AnError

	timer := inst.RecordOperationDuration(t.Context())
	op := inst.ExportMetrics(t.Context(), n)

	const success = 0
	op.End(success, err)
	timer.Stop(err)

	assertExportMetricsMetrics(t, collect(), n, success, err)
}

func TestInstrumentationExportMetricsPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	err := assert.AnError

	timer := inst.RecordOperationDuration(t.Context())
	op := inst.ExportMetrics(t.Context(), n)

	const success = 5
	op.End(success, err)
	timer.Stop(err)

	assertExportMetricsMetrics(t, collect(), n, success, err)
}

func TestRecordCollectionDurationSuccess(t *testing.T) {
	inst, collect := setup(t)

	inst.RecordCollectionDuration(t.Context()).Stop(nil)

	assertCollectionOnly(t, collect(), nil)
}

func TestRecordCollectionDurationError(t *testing.T) {
	inst, collect := setup(t)

	wantErr := assert.AnError
	inst.RecordCollectionDuration(t.Context()).Stop(wantErr)

	assertCollectionOnly(t, collect(), wantErr)
}

func setupBench(b *testing.B) *observ.Instrumentation {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	inst, err := observ.NewInstrumentation(ID)
	if err != nil {
		b.Fatalf("failed to create instrumentation: %v", err)
	}
	if inst == nil {
		b.Fatal("expected instrumentation, got nil")
	}
	return inst
}

func BenchmarkInstrumentationExportMetrics(b *testing.B) {
	const nSpans = 10
	run := func(success int64, err error) func(*testing.B) {
		inst := setupBench(b)
		return func(b *testing.B) {
			ctx := b.Context()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.ExportMetrics(ctx, nSpans).End(success, err)
			}
		}
	}

	err := errors.New("benchmark error")
	b.Run("NoError", run(nSpans, nil))
	b.Run("AllError", run(0, err))
	b.Run("PartialError", run(4, err))
}

func BenchmarkInstrumentationRecordOperationDuration(b *testing.B) {
	run := func(err error) func(*testing.B) {
		inst := setupBench(b)
		return func(b *testing.B) {
			ctx := b.Context()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.RecordOperationDuration(ctx).Stop(err)
			}
		}
	}

	err := errors.New("benchmark error")
	b.Run("NoError", run(nil))
	b.Run("Error", run(err))
}

func BenchmarkInstrumentationRecordCollectionDuration(b *testing.B) {
	run := func(err error) func(*testing.B) {
		inst := setupBench(b)
		return func(b *testing.B) {
			ctx := b.Context()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.RecordCollectionDuration(ctx).Stop(err)
			}
		}
	}

	err := errors.New("benchmark error")
	b.Run("NoError", run(nil))
	b.Run("Error", run(err))
}
