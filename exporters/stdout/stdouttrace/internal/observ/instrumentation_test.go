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
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/observ"
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
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "inflight metric")
	assert.ErrorContains(t, err, "span exported metric")
	assert.ErrorContains(t, err, "operation duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := observ.NewInstrumentation(ID)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

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

func spanInflight() metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterSpanInflight{}.Name(),
		Description: otelconv.SDKExporterSpanInflight{}.Description(),
		Unit:        otelconv.SDKExporterSpanInflight{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: set(nil), Value: 0},
			},
		},
	}
}

func spanExported(success, total int64, err error) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{Attributes: set(nil), Value: success},
	}
	if err != nil {
		dp = append(dp, metricdata.DataPoint[int64]{
			Attributes: set(err),
			Value:      total - success,
		})
	}
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterSpanExported{}.Name(),
		Description: otelconv.SDKExporterSpanExported{}.Description(),
		Unit:        otelconv.SDKExporterSpanExported{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
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

func assertMetrics(t *testing.T, got metricdata.ScopeMetrics, spans, success int64, err error) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := spanInflight()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = spanExported(success, spans, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = operationDuration(err)
	metricdatatest.AssertEqual(t, want, m[2], o, metricdatatest.IgnoreValue())
}

func TestInstrumentationExportSpans(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	inst.ExportSpans(t.Context(), n).End(n, nil)

	assertMetrics(t, collect(), n, n, nil)
}

func TestInstrumentationExportSpansAllErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	const success = 0
	inst.ExportSpans(t.Context(), n).End(success, assert.AnError)

	assertMetrics(t, collect(), n, success, assert.AnError)
}

func TestInstrumentationExportSpansPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	const success = 5
	inst.ExportSpans(t.Context(), n).End(success, assert.AnError)

	assertMetrics(t, collect(), n, success, assert.AnError)
}

func BenchmarkInstrumentationExportSpans(b *testing.B) {
	setup := func(b *testing.B) *observ.Instrumentation {
		b.Helper()
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := observ.NewInstrumentation(ID)
		if err != nil {
			b.Fatalf("failed to create instrumentation: %v", err)
		}
		return inst
	}

	const nSpans = 10
	err := errors.New("benchmark error")
	run := func(n int64, err error) func(*testing.B) {
		return func(b *testing.B) {
			inst := setup(b)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.ExportSpans(b.Context(), nSpans).End(n, err)
			}
		}
	}

	b.Run("NoError", run(nSpans, nil))
	b.Run("PartialError", run(4, err))
	b.Run("FullError", run(0, err))
}
