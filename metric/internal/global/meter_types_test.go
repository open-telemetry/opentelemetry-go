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

package global // import "go.opentelemetry.io/otel/metric/internal/global"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type testMeterProvider struct {
	count int
}

func (p *testMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.count++

	return &testMeter{}
}

type testMeter struct {
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

func (m *testMeter) Int64Counter(name string, options ...instrument.Int64Option) (instrument.Int64Counter, error) {
	m.siCount++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64UpDownCounter(name string, options ...instrument.Int64Option) (instrument.Int64UpDownCounter, error) {
	m.siUDCount++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64Histogram(name string, options ...instrument.Int64Option) (instrument.Int64Histogram, error) {
	m.siHist++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableCounter, error) {
	m.aiCount++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableUpDownCounter, error) {
	m.aiUDCount++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableGauge, error) {
	m.aiGauge++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Float64Counter(name string, options ...instrument.Float64Option) (instrument.Float64Counter, error) {
	m.sfCount++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64UpDownCounter(name string, options ...instrument.Float64Option) (instrument.Float64UpDownCounter, error) {
	m.sfUDCount++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64Histogram(name string, options ...instrument.Float64Option) (instrument.Float64Histogram, error) {
	m.sfHist++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableCounter, error) {
	m.afCount++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableUpDownCounter, error) {
	m.afUDCount++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableGauge, error) {
	m.afGauge++
	return &testCountingFloatInstrument{}, nil
}

// RegisterCallback captures the function that will be called during Collect.
func (m *testMeter) RegisterCallback(f metric.Callback, i ...instrument.Asynchronous) (metric.Registration, error) {
	m.callbacks = append(m.callbacks, f)
	return testReg{
		f: func(idx int) func() {
			return func() { m.callbacks[idx] = nil }
		}(len(m.callbacks) - 1),
	}, nil
}

type testReg struct {
	f func()
}

func (r testReg) Unregister() error {
	r.f()
	return nil
}

// This enables async collection.
func (m *testMeter) collect() {
	ctx := context.Background()
	o := observationRecorder{ctx}
	for _, f := range m.callbacks {
		if f == nil {
			// Unregister.
			continue
		}
		_ = f(ctx, o)
	}
}

type observationRecorder struct {
	ctx context.Context
}

func (o observationRecorder) ObserveFloat64(i instrument.Float64Observable, value float64, attr ...attribute.KeyValue) {
	iImpl, ok := i.(*testCountingFloatInstrument)
	if ok {
		iImpl.observe()
	}
}

func (o observationRecorder) ObserveInt64(i instrument.Int64Observable, value int64, attr ...attribute.KeyValue) {
	iImpl, ok := i.(*testCountingIntInstrument)
	if ok {
		iImpl.observe()
	}
}
