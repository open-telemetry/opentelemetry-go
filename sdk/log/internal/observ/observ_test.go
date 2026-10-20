// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/internal/observ"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
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
		for _, sm := range got.ScopeMetrics {
			if sm.Scope.Name == observ.ScopeName {
				return sm
			}
		}
		return metricdata.ScopeMetrics{}
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

type meterProvider struct {
	mapi.MeterProvider
	m mapi.Meter
}

func (p meterProvider) Meter(_ string, _ ...mapi.MeterOption) mapi.Meter { return p.m }

// errOnNthObsCounterMeter fails Int64ObservableUpDownCounter on the nth call.
type errOnNthObsCounterMeter struct {
	noop.Meter
	n, cnt int
	err    error
}

func (m *errOnNthObsCounterMeter) Int64ObservableUpDownCounter(
	name string,
	opts ...mapi.Int64ObservableUpDownCounterOption,
) (mapi.Int64ObservableUpDownCounter, error) {
	m.cnt++
	if m.cnt == m.n {
		return nil, m.err
	}
	return m.Meter.Int64ObservableUpDownCounter(name, opts...)
}

type errCallbackMeter struct {
	noop.Meter
	err error
}

func (m *errCallbackMeter) RegisterCallback(mapi.Callback, ...mapi.Observable) (mapi.Registration, error) {
	return nil, m.err
}

type errCounterMeter struct {
	noop.Meter
	err error
}

func (m *errCounterMeter) Int64Counter(_ string, _ ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}
