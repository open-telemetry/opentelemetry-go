// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/internal/global"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

type altMeterProvider struct {
	t      *testing.T
	meters []*altMeter
	embedded.MeterProvider
}

var _ metric.MeterProvider = &altMeterProvider{}

func (amp *altMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	am := &altMeter{
		provider: amp,
	}
	amp.meters = append(amp.meters, am)
	return am
}

type altMeter struct {
	provider *altMeterProvider
	cbs      []metric.Callback
	embedded.Meter
}

var _ metric.Meter = &altMeter{}

type testAsyncCounter struct {
	meter *altMeter
	embedded.Int64ObservableCounter
	metric.Int64Observable
}

var _ metric.Int64ObservableCounter = &testAsyncCounter{}

type altRegistration struct {
	cb metric.Callback
	embedded.Registration
}

type altObserver struct {
	t *testing.T
	embedded.Observer
}

func (*altRegistration) Unregister() error {
	return nil
}

func (am *altMeter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return nil, nil
}

func (am *altMeter) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return nil, nil
}

func (am *altMeter) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return nil, nil
}

func (am *altMeter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return nil, nil
}

func (am *altMeter) Int64ObservableCounter(name string, options ...metric.Int64ObservableCounterOption) (metric.Int64ObservableCounter, error) {
	return &testAsyncCounter{
		meter: am,
	}, nil
}

func (am *altMeter) Int64ObservableUpDownCounter(name string, options ...metric.Int64ObservableUpDownCounterOption) (metric.Int64ObservableUpDownCounter, error) {
	return nil, nil
}

func (am *altMeter) Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
	return nil, nil
}

func (am *altMeter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return nil, nil
}

func (am *altMeter) Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return nil, nil
}

func (am *altMeter) Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return nil, nil
}

func (am *altMeter) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return nil, nil
}

func (am *altMeter) Float64ObservableCounter(name string, options ...metric.Float64ObservableCounterOption) (metric.Float64ObservableCounter, error) {
	return nil, nil
}

func (am *altMeter) Float64ObservableUpDownCounter(name string, options ...metric.Float64ObservableUpDownCounterOption) (metric.Float64ObservableUpDownCounter, error) {
	return nil, nil
}

func (am *altMeter) Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
	// Note: The global delegation also breaks when we return nil in one of these!
	return nil, nil
}

func (am *altMeter) RegisterCallback(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
	for _, inst := range instruments {
		switch inst.(type) {
		case *testAsyncCounter:
			// OK!
		default:
			am.provider.t.Errorf("unexpected type %T", inst)
		}
	}
	am.cbs = append(am.cbs, f)
	return &altRegistration{cb: f}, nil
}

func (ao *altObserver) ObserveFloat64(inst metric.Float64Observable, _ float64, _ ...metric.ObserveOption) {
	ao.observe(inst)
}

func (ao *altObserver) ObserveInt64(inst metric.Int64Observable, _ int64, _ ...metric.ObserveOption) {
	ao.observe(inst)
}

func (ao *altObserver) observe(inst any) {
	switch inst.(type) {
	case *testAsyncCounter:
		// OK!
	default:
		ao.t.Errorf("unexpected type %T", inst)
	}
}

func TestMeterDelegation(t *testing.T) {
	ResetForTest(t)

	amp := &altMeterProvider{t: t}

	gm := MeterProvider().Meter("test")
	ai, err := gm.Int64ObservableCounter("test_counter")
	require.NoError(t, err)

	_, err = gm.RegisterCallback(func(_ context.Context, obs metric.Observer) error {
		obs.ObserveInt64(ai, 10)
		return nil
	}, ai)
	require.NoError(t, err)

	SetMeterProvider(amp)

	ctx := context.Background()
	ao := &altObserver{t: t}
	for _, meter := range amp.meters {
		for _, cb := range meter.cbs {
			require.NoError(t, cb(ctx, ao))
		}
	}
}
