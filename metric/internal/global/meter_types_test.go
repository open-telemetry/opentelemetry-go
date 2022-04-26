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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

type testMeterProvider struct {
	count int
}

func (p *testMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.count++

	return &testMeter{}
}

type testMeter struct {
	afCount int
	aiCount int
	sfCount int
	siCount int

	callbacks []func(context.Context)
}

// AsyncInt64 is the namespace for the Asynchronous Integer instruments.
//
// To Observe data with instruments it must be registered in a callback.
func (m *testMeter) AsyncInt64() asyncint64.InstrumentProvider {
	m.aiCount++
	return &testAIInstrumentProvider{}
}

// AsyncFloat64 is the namespace for the Asynchronous Float instruments
//
// To Observe data with instruments it must be registered in a callback.
func (m *testMeter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	m.afCount++
	return &testAFInstrumentProvider{}
}

// RegisterCallback captures the function that will be called during Collect.
//
// It is only valid to call Observe within the scope of the passed function,
// and only on the instruments that were registered with this call.
func (m *testMeter) RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error {
	m.callbacks = append(m.callbacks, function)
	return nil
}

// SyncInt64 is the namespace for the Synchronous Integer instruments.
func (m *testMeter) SyncInt64() syncint64.InstrumentProvider {
	m.siCount++
	return &testSIInstrumentProvider{}
}

// SyncFloat64 is the namespace for the Synchronous Float instruments.
func (m *testMeter) SyncFloat64() syncfloat64.InstrumentProvider {
	m.sfCount++
	return &testSFInstrumentProvider{}
}

// This enables async collection.
func (m *testMeter) collect() {
	ctx := context.Background()
	for _, f := range m.callbacks {
		f(ctx)
	}
}

type testAFInstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip testAFInstrumentProvider) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	return &testCountingFloatInstrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip testAFInstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	return &testCountingFloatInstrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip testAFInstrumentProvider) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	return &testCountingFloatInstrument{}, nil
}

type testAIInstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip testAIInstrumentProvider) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	return &testCountingIntInstrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip testAIInstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	return &testCountingIntInstrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip testAIInstrumentProvider) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	return &testCountingIntInstrument{}, nil
}

type testSFInstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip testSFInstrumentProvider) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	return &testCountingFloatInstrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip testSFInstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	return &testCountingFloatInstrument{}, nil
}

// Histogram creates an instrument for recording a distribution of values.
func (ip testSFInstrumentProvider) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	return &testCountingFloatInstrument{}, nil
}

type testSIInstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip testSIInstrumentProvider) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	return &testCountingIntInstrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip testSIInstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	return &testCountingIntInstrument{}, nil
}

// Histogram creates an instrument for recording a distribution of values.
func (ip testSIInstrumentProvider) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	return &testCountingIntInstrument{}, nil
}
