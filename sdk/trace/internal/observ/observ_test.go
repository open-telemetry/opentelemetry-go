// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/trace/internal/observ"
)

func setup(t *testing.T) func() metricdata.ScopeMetrics {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	return func() metricdata.ScopeMetrics {
		var got metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(t.Context(), &got))
		if len(got.ScopeMetrics) != 1 {
			return metricdata.ScopeMetrics{}
		}
		return got.ScopeMetrics[0]
	}
}

func scopeMetrics(metrics ...metricdata.Metrics) metricdata.ScopeMetrics {
	return metricdata.ScopeMetrics{
		Scope: instrumentation.Scope{
			Name:      observ.ScopeName,
			Version:   sdk.Version(),
			SchemaURL: observ.SchemaURL,
		},
		Metrics: metrics,
	}
}

func check(t *testing.T, got metricdata.ScopeMetrics, want ...metricdata.Metrics) {
	o := []metricdatatest.Option{
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreExemplars(),
	}
	metricdatatest.AssertEqual(t, scopeMetrics(want...), got, o...)
}

func dPt(set attribute.Set, value int64) metricdata.DataPoint[int64] {
	return metricdata.DataPoint[int64]{Attributes: set, Value: value}
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

func (m *errMeter) Int64ObservableUpDownCounter(
	string,
	...mapi.Int64ObservableUpDownCounterOption,
) (mapi.Int64ObservableUpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) RegisterCallback(mapi.Callback, ...mapi.Observable) (mapi.Registration, error) {
	return nil, m.err
}
