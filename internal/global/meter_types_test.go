// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global // import "go.opentelemetry.io/otel/internal/global"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
)

type testMeterProvider struct {
	embedded.MeterProvider

	count int
}

func (p *testMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
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

	siCount   int
	siUDCount int
	siHist    int

	callbacks []metric.Callback
}

func (m *testMeter) Int64Counter(name string, options ...instrument.CounterOption[int64]) (instrument.Counter[int64], error) {
	m.siCount++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Int64UpDownCounter(name string, options ...instrument.UpDownCounterOption[int64]) (instrument.UpDownCounter[int64], error) {
	m.siUDCount++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Int64Histogram(name string, options ...instrument.HistogramOption[int64]) (instrument.Histogram[int64], error) {
	m.siHist++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Int64ObservableCounter(name string, options ...instrument.ObservableCounterOption[int64]) (instrument.ObservableCounter[int64], error) {
	m.aiCount++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Int64ObservableUpDownCounter(name string, options ...instrument.ObservableUpDownCounterOption[int64]) (instrument.ObservableUpDownCounter[int64], error) {
	m.aiUDCount++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Int64ObservableGauge(name string, options ...instrument.ObservableGaugeOption[int64]) (instrument.ObservableGauge[int64], error) {
	m.aiGauge++
	return &testCountingInstrument[int64]{}, nil
}

func (m *testMeter) Float64Counter(name string, options ...instrument.CounterOption[float64]) (instrument.Counter[float64], error) {
	m.sfCount++
	return &testCountingInstrument[float64]{}, nil
}

func (m *testMeter) Float64UpDownCounter(name string, options ...instrument.UpDownCounterOption[float64]) (instrument.UpDownCounter[float64], error) {
	m.sfUDCount++
	return &testCountingInstrument[float64]{}, nil
}

func (m *testMeter) Float64Histogram(name string, options ...instrument.HistogramOption[float64]) (instrument.Histogram[float64], error) {
	m.sfHist++
	return &testCountingInstrument[float64]{}, nil
}

func (m *testMeter) Float64ObservableCounter(name string, options ...instrument.ObservableCounterOption[float64]) (instrument.ObservableCounter[float64], error) {
	m.afCount++
	return &testCountingInstrument[float64]{}, nil
}

func (m *testMeter) Float64ObservableUpDownCounter(name string, options ...instrument.ObservableUpDownCounterOption[float64]) (instrument.ObservableUpDownCounter[float64], error) {
	m.afUDCount++
	return &testCountingInstrument[float64]{}, nil
}

func (m *testMeter) Float64ObservableGauge(name string, options ...instrument.ObservableGaugeOption[float64]) (instrument.ObservableGauge[float64], error) {
	m.afGauge++
	return &testCountingInstrument[float64]{}, nil
}

// RegisterCallback captures the function that will be called during Collect.
func (m *testMeter) RegisterCallback(f metric.Callback, i ...instrument.Observable) (metric.Registration, error) {
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

func (o observationRecorder) ObserveFloat64(i instrument.ObservableT[float64], value float64, attr ...attribute.KeyValue) {
	iImpl, ok := i.(*testCountingInstrument[float64])
	if ok {
		iImpl.observe()
	}
}

func (o observationRecorder) ObserveInt64(i instrument.ObservableT[int64], value int64, attr ...attribute.KeyValue) {
	iImpl, ok := i.(*testCountingInstrument[int64])
	if ok {
		iImpl.observe()
	}
}
