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
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/internal/observ"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	ID            = int64(42)
	ComponentType = "test-reader"
)

var Scope = instrumentation.Scope{
	Name:      observ.ScopeName,
	Version:   sdk.Version(),
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

func (m *errMeter) Float64Histogram(string, ...mapi.Float64HistogramOption) (mapi.Float64Histogram, error) {
	return nil, m.err
}

func TestNewInstrumentationObservabilityErrors(t *testing.T) {
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	_, err := observ.NewInstrumentation(ComponentType, ID)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors should be joined")

	assert.ErrorContains(t, err, "collection duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := observ.NewInstrumentation(ComponentType, ID)
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

	inst, err := observ.NewInstrumentation(ComponentType, ID)
	require.NoError(t, err)
	require.NotNil(t, inst)

	return inst, func() metricdata.ScopeMetrics {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

func baseAttrs(err error) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(observ.ComponentName(ComponentType, ID)),
		semconv.OTelComponentTypeKey.String(ComponentType),
	}
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attrs
}

func set(err error) attribute.Set {
	return attribute.NewSet(baseAttrs(err)...)
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

func assertCollectionMetrics(t *testing.T, got metricdata.ScopeMetrics, err error) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 1, "expected 1 metric (collection duration)")

	want := collectionDuration(err)
	metricdatatest.AssertEqual(t, want, m[0], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
}

func TestInstrumentationCollectMetricsSuccess(t *testing.T) {
	inst, collect := setup(t)

	inst.CollectMetrics(t.Context()).End(nil)

	assertCollectionMetrics(t, collect(), nil)
}

func TestInstrumentationCollectMetricsError(t *testing.T) {
	inst, collect := setup(t)

	wantErr := assert.AnError
	inst.CollectMetrics(t.Context()).End(wantErr)

	assertCollectionMetrics(t, collect(), wantErr)
}

func TestComponentName(t *testing.T) {
	tests := []struct {
		componentType string
		id            int64
		want          string
	}{
		{componentType: "periodic_metric_reader", id: 0, want: "periodic_metric_reader/0"},
		{componentType: "periodic_metric_reader", id: 1, want: "periodic_metric_reader/1"},
		{componentType: "periodic_metric_reader", id: 42, want: "periodic_metric_reader/42"},
		{componentType: "periodic_metric_reader", id: -1, want: "periodic_metric_reader/-1"},
		{componentType: "manual_metric_reader", id: 0, want: "manual_metric_reader/0"},
	}

	for _, tt := range tests {
		got := observ.ComponentName(tt.componentType, tt.id)
		assert.Equal(t, tt.want, got)
	}
}

func setupBench(b *testing.B) *observ.Instrumentation {
	b.Helper()
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Set up a proper MeterProvider for benchmarks
	original := otel.GetMeterProvider()
	b.Cleanup(func() { otel.SetMeterProvider(original) })

	r := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(r))
	otel.SetMeterProvider(mp)

	inst, err := observ.NewInstrumentation(ComponentType, ID)
	if err != nil {
		b.Fatalf("failed to create instrumentation: %v", err)
	}
	if inst == nil {
		b.Fatal("expected instrumentation, got nil")
	}
	return inst
}

func BenchmarkInstrumentationCollectMetrics(b *testing.B) {
	run := func(err error) func(*testing.B) {
		inst := setupBench(b)
		return func(b *testing.B) {
			ctx := b.Context()
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.CollectMetrics(ctx).End(err)
			}
		}
	}

	err := errors.New("benchmark error")
	b.Run("NoError", run(nil))
	b.Run("Error", run(err))
}
