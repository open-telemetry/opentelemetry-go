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

	"go.opentelemetry.io/otel/metric"
)

type testMeterProvider struct {
	count int
}

func (p *testMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.count++

	return &testMeter{}
}

type testMeter struct {
	afCounter       int
	afUpDownCounter int
	afGauge         int

	aiCounter       int
	aiUpDownCounter int
	aiGauge         int

	sfCounter       int
	sfUpDownCounter int
	sfHistogram     int

	siCounter       int
	siUpDownCounter int
	siHistogram     int

	callbacks []metric.Callback
}

func (m *testMeter) Float64Counter(name string, opts ...metric.InstrumentOption) (metric.Float64Counter, error) {
	m.sfCounter++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Float64UpDownCounter, error) {
	m.sfUpDownCounter++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64Histogram(name string, opts ...metric.InstrumentOption) (metric.Float64Histogram, error) {
	m.sfHistogram++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableCounter, error) {
	m.afCounter++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableUpDownCounter, error) {
	m.afUpDownCounter++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Float64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Float64ObservableGauge, error) {
	m.afGauge++
	return &testCountingFloatInstrument{}, nil
}

func (m *testMeter) Int64Counter(name string, opts ...metric.InstrumentOption) (metric.Int64Counter, error) {
	m.siCounter++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Int64UpDownCounter, error) {
	m.siUpDownCounter++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64Histogram(name string, opts ...metric.InstrumentOption) (metric.Int64Histogram, error) {
	m.siHistogram++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableCounter, error) {
	m.aiCounter++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableUpDownCounter, error) {
	m.aiUpDownCounter++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) Int64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Int64ObservableGauge, error) {
	m.aiGauge++
	return &testCountingIntInstrument{}, nil
}

func (m *testMeter) RegisterCallback(f metric.Callback, instrument metric.Observable, additional ...metric.Observable) (metric.Unregisterer, error) {
	m.callbacks = append(m.callbacks, f)
	return unregisterer{}, nil
}

// This enables async collection.
func (m *testMeter) collect() {
	ctx := context.Background()
	for _, f := range m.callbacks {
		_ = f(ctx)
	}
}
