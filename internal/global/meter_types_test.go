// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package global // import "go.opentelemetry.io/otel/internal/global"

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

type testMeterProvider struct {
	embedded.MeterProvider

	count int
}

func (p *testMeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	p.count++
	return &testMeter{}
}

type testMeter struct {
	embedded.Meter

	afCount   int
	afUDCount int
	afGauge   int

	aiCount   int
	aiUDCount int
	aiGauge   int

	sfCount   int
	sfUDCount int
	sfHist    int
	sfGauge   int

	siCount   int
	siUDCount int
	siHist    int
	siGauge   int

	callbacks []metric.Callback
}

func (m *testMeter) Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	m.siCount++
	return &testInt64Counter{}, nil
}

func (m *testMeter) Int64UpDownCounter(
	string,
	...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	m.siUDCount++
	return &testInt64UpDownCounter{}, nil
}

func (m *testMeter) Int64Histogram(string, ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	m.siHist++
	return &testInt64Histogram{}, nil
}

func (m *testMeter) Int64Gauge(string, ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	m.siGauge++
	return &testInt64Gauge{}, nil
}

func (m *testMeter) Int64ObservableCounter(
	string,
	...metric.Int64ObservableCounterOption,
) (metric.Int64ObservableCounter, error) {
	m.aiCount++
	return &testInt64Observable{}, nil
}

func (m *testMeter) Int64ObservableUpDownCounter(
	string,
	...metric.Int64ObservableUpDownCounterOption,
) (metric.Int64ObservableUpDownCounter, error) {
	m.aiUDCount++
	return &testInt64Observable{}, nil
}

func (m *testMeter) Int64ObservableGauge(
	string,
	...metric.Int64ObservableGaugeOption,
) (metric.Int64ObservableGauge, error) {
	m.aiGauge++
	return &testInt64Observable{}, nil
}

func (m *testMeter) Float64Counter(string, ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	m.sfCount++
	return &testFloat64Counter{}, nil
}

func (m *testMeter) Float64UpDownCounter(
	string,
	...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	m.sfUDCount++
	return &testFloat64UpDownCounter{}, nil
}

func (m *testMeter) Float64Histogram(
	string,
	...metric.Float64HistogramOption,
) (metric.Float64Histogram, error) {
	m.sfHist++
	return &testFloat64Histogram{}, nil
}

func (m *testMeter) Float64Gauge(string, ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	m.sfGauge++
	return &testFloat64Gauge{}, nil
}

func (m *testMeter) Float64ObservableCounter(
	string,
	...metric.Float64ObservableCounterOption,
) (metric.Float64ObservableCounter, error) {
	m.afCount++
	return &testFloat64Observable{}, nil
}

func (m *testMeter) Float64ObservableUpDownCounter(
	string,
	...metric.Float64ObservableUpDownCounterOption,
) (metric.Float64ObservableUpDownCounter, error) {
	m.afUDCount++
	return &testFloat64Observable{}, nil
}

func (m *testMeter) Float64ObservableGauge(
	string,
	...metric.Float64ObservableGaugeOption,
) (metric.Float64ObservableGauge, error) {
	m.afGauge++
	return &testFloat64Observable{}, nil
}

// RegisterCallback captures the function that will be called during Collect.
func (m *testMeter) RegisterCallback(f metric.Callback, _ ...metric.Observable) (metric.Registration, error) {
	m.callbacks = append(m.callbacks, f)
	return testReg{
		f: func(idx int) func() {
			return func() { m.callbacks[idx] = nil }
		}(len(m.callbacks) - 1),
	}, nil
}

type testReg struct {
	embedded.Registration

	f func()
}

func (r testReg) Unregister() error {
	r.f()
	return nil
}

// This enables async collection.
func (m *testMeter) collect() {
	ctx := context.Background()
	o := observationRecorder{ctx: ctx}
	for _, f := range m.callbacks {
		if f == nil {
			// Unregister.
			continue
		}
		_ = f(ctx, o)
	}
}

type observationRecorder struct {
	embedded.Observer

	ctx context.Context
}

func (observationRecorder) ObserveFloat64(i metric.Float64Observable, _ float64, _ ...metric.ObserveOption) {
	iImpl, ok := i.(*testFloat64Observable)
	if ok {
		iImpl.observe()
	}
}

func (observationRecorder) ObserveInt64(i metric.Int64Observable, _ int64, _ ...metric.ObserveOption) {
	iImpl, ok := i.(*testInt64Observable)
	if ok {
		iImpl.observe()
	}
}
