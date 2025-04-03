// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/internal/global"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"
)

// Below, an alternate meter provider is constructed specifically to
// test the asynchronous instrument path.  The alternative SDK uses
// no-op implementations for its synchronous instruments, and the six
// asynchronous instrument types are created here to test that
// instruments and callbacks are unwrapped inside this library.

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

type testAiCounter struct {
	meter *altMeter
	embedded.Int64ObservableCounter
	metric.Int64Observable
}

var _ metric.Int64ObservableCounter = &testAiCounter{}

type testAfCounter struct {
	meter *altMeter
	embedded.Float64ObservableCounter
	metric.Float64Observable
}

var _ metric.Float64ObservableCounter = &testAfCounter{}

type testAiUpDownCounter struct {
	meter *altMeter
	embedded.Int64ObservableUpDownCounter
	metric.Int64Observable
}

var _ metric.Int64ObservableUpDownCounter = &testAiUpDownCounter{}

type testAfUpDownCounter struct {
	meter *altMeter
	embedded.Float64ObservableUpDownCounter
	metric.Float64Observable
}

var _ metric.Float64ObservableUpDownCounter = &testAfUpDownCounter{}

type testAiGauge struct {
	meter *altMeter
	embedded.Int64ObservableGauge
	metric.Int64Observable
}

var _ metric.Int64ObservableGauge = &testAiGauge{}

type testAfGauge struct {
	meter *altMeter
	embedded.Float64ObservableGauge
	metric.Float64Observable
}

var _ metric.Float64ObservableGauge = &testAfGauge{}

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

func (am *altMeter) Int64Counter(name string, _ ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return noop.NewMeterProvider().Meter("noop").Int64Counter(name)
}

func (am *altMeter) Int64UpDownCounter(
	name string,
	_ ...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	return noop.NewMeterProvider().Meter("noop").Int64UpDownCounter(name)
}

func (am *altMeter) Int64Histogram(name string, _ ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return noop.NewMeterProvider().Meter("noop").Int64Histogram(name)
}

func (am *altMeter) Int64Gauge(name string, _ ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return noop.NewMeterProvider().Meter("noop").Int64Gauge(name)
}

func (am *altMeter) Int64ObservableCounter(
	name string,
	options ...metric.Int64ObservableCounterOption,
) (metric.Int64ObservableCounter, error) {
	return &testAiCounter{
		meter: am,
	}, nil
}

func (am *altMeter) Int64ObservableUpDownCounter(
	name string,
	options ...metric.Int64ObservableUpDownCounterOption,
) (metric.Int64ObservableUpDownCounter, error) {
	return &testAiUpDownCounter{
		meter: am,
	}, nil
}

func (am *altMeter) Int64ObservableGauge(
	name string,
	options ...metric.Int64ObservableGaugeOption,
) (metric.Int64ObservableGauge, error) {
	return &testAiGauge{
		meter: am,
	}, nil
}

func (am *altMeter) Float64Counter(name string, _ ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return noop.NewMeterProvider().Meter("noop").Float64Counter(name)
}

func (am *altMeter) Float64UpDownCounter(
	name string,
	_ ...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	return noop.NewMeterProvider().Meter("noop").Float64UpDownCounter(name)
}

func (am *altMeter) Float64Histogram(
	name string,
	options ...metric.Float64HistogramOption,
) (metric.Float64Histogram, error) {
	return noop.NewMeterProvider().Meter("noop").Float64Histogram(name)
}

func (am *altMeter) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return noop.NewMeterProvider().Meter("noop").Float64Gauge(name)
}

func (am *altMeter) Float64ObservableCounter(
	name string,
	options ...metric.Float64ObservableCounterOption,
) (metric.Float64ObservableCounter, error) {
	return &testAfCounter{
		meter: am,
	}, nil
}

func (am *altMeter) Float64ObservableUpDownCounter(
	name string,
	options ...metric.Float64ObservableUpDownCounterOption,
) (metric.Float64ObservableUpDownCounter, error) {
	return &testAfUpDownCounter{
		meter: am,
	}, nil
}

func (am *altMeter) Float64ObservableGauge(
	name string,
	options ...metric.Float64ObservableGaugeOption,
) (metric.Float64ObservableGauge, error) {
	return &testAfGauge{
		meter: am,
	}, nil
}

func (am *altMeter) RegisterCallback(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
	for _, inst := range instruments {
		switch inst.(type) {
		case *testAiCounter, *testAfCounter,
			*testAiUpDownCounter, *testAfUpDownCounter,
			*testAiGauge, *testAfGauge:
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
	case *testAiCounter, *testAfCounter,
		*testAiUpDownCounter, *testAfUpDownCounter,
		*testAiGauge, *testAfGauge:
		// OK!
	default:
		ao.t.Errorf("unexpected type %T", inst)
	}
}

func TestMeterDelegation(t *testing.T) {
	ResetForTest(t)

	amp := &altMeterProvider{t: t}

	gm := MeterProvider().Meter("test")
	aic, err := gm.Int64ObservableCounter("test_counter_i")
	require.NoError(t, err)
	afc, err := gm.Float64ObservableCounter("test_counter_f")
	require.NoError(t, err)
	aiu, err := gm.Int64ObservableUpDownCounter("test_updowncounter_i")
	require.NoError(t, err)
	afu, err := gm.Float64ObservableUpDownCounter("test_updowncounter_f")
	require.NoError(t, err)
	aig, err := gm.Int64ObservableGauge("test_gauge_i")
	require.NoError(t, err)
	afg, err := gm.Float64ObservableGauge("test_gauge_f")
	require.NoError(t, err)

	_, err = gm.RegisterCallback(func(_ context.Context, obs metric.Observer) error {
		obs.ObserveInt64(aic, 10)
		obs.ObserveFloat64(afc, 10)
		obs.ObserveInt64(aiu, 10)
		obs.ObserveFloat64(afu, 10)
		obs.ObserveInt64(aig, 10)
		obs.ObserveFloat64(afg, 10)
		return nil
	}, aic, afc, aiu, afu, aig, afg)
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
