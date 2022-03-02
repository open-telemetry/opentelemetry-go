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

type test_MeterProvider struct {
	count int
}

func (p *test_MeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.count++

	return &test_Meter{}
}

type test_Meter struct {
	afCount int
	aiCount int
	sfCount int
	siCount int

	callbacks []func(context.Context)
}

// AsyncInt64 is the namespace for the Asynchronous Integer instruments.
//
// To Observe data with instruments it must be registered in a callback.
func (m *test_Meter) AsyncInt64() asyncint64.InstrumentProvider {
	m.aiCount++
	return &test_ai_InstrumentProvider{}
}

// AsyncFloat64 is the namespace for the Asynchronous Float instruments
//
// To Observe data with instruments it must be registered in a callback.
func (m *test_Meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	m.afCount++
	return &test_af_InstrumentProvider{}
}

// RegisterCallback captures the function that will be called during Collect.
//
// It is only valid to call Observe within the scope of the passed function,
// and only on the instruments that were registered with this call.
func (m *test_Meter) RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error {
	m.callbacks = append(m.callbacks, function)
	return nil
}

// SyncInt64 is the namespace for the Synchronous Integer instruments
func (m *test_Meter) SyncInt64() syncint64.InstrumentProvider {
	m.siCount++
	return &test_si_InstrumentProvider{}
}

// SyncFloat64 is the namespace for the Synchronous Float instruments
func (m *test_Meter) SyncFloat64() syncfloat64.InstrumentProvider {
	m.sfCount++
	return &test_sf_InstrumentProvider{}
}

// This enables async collection
func (m *test_Meter) collect() {
	ctx := context.Background()
	for _, f := range m.callbacks {
		f(ctx)
	}
}

type test_af_InstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip test_af_InstrumentProvider) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	return &test_counting_float_instrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip test_af_InstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	return &test_counting_float_instrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip test_af_InstrumentProvider) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	return &test_counting_float_instrument{}, nil
}

type test_ai_InstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip test_ai_InstrumentProvider) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	return &test_counting_int_instrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip test_ai_InstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	return &test_counting_int_instrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip test_ai_InstrumentProvider) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	return &test_counting_int_instrument{}, nil
}

type test_sf_InstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip test_sf_InstrumentProvider) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	return &test_counting_float_instrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip test_sf_InstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	return &test_counting_float_instrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip test_sf_InstrumentProvider) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	return &test_counting_float_instrument{}, nil
}

type test_si_InstrumentProvider struct{}

// Counter creates an instrument for recording increasing values.
func (ip test_si_InstrumentProvider) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	return &test_counting_int_instrument{}, nil
}

// UpDownCounter creates an instrument for recording changes of a value.
func (ip test_si_InstrumentProvider) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	return &test_counting_int_instrument{}, nil
}

// Gauge creates an instrument for recording the current value.
func (ip test_si_InstrumentProvider) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	return &test_counting_int_instrument{}, nil
}
