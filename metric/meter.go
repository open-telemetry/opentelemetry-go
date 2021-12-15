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

package metric

import (
	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"
)

// MeterProvider supports creating named Meter instances, for instrumenting
// an application containing multiple libraries of code.
type MeterProvider interface {
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter is an instance of an OpenTelemetry metrics interface for an
// individual named library of code.  This is the top-level entry
// point for creating instruments.
type Meter struct {
	sdkapi.MeterImpl
}

type AsyncFloat64Instruments struct {
	sdkapi.MeterImpl
}

type AsyncInt64Instruments struct {
	sdkapi.MeterImpl
}

type SyncFloat64Instruments struct {
	sdkapi.MeterImpl
}

type SyncInt64Instruments struct {
	sdkapi.MeterImpl
}

func (m Meter) AsyncInt64() AsyncInt64Instruments {
	return AsyncInt64Instruments{m.MeterImpl}
}

func (m Meter) AsyncFloat64() AsyncFloat64Instruments {
	return AsyncFloat64Instruments{m.MeterImpl}
}

func (m Meter) SyncInt64() SyncInt64Instruments {
	return SyncInt64Instruments{m.MeterImpl}
}

func (m Meter) SyncFloat64() SyncFloat64Instruments {
	return SyncFloat64Instruments{m.MeterImpl}
}

func (m AsyncFloat64Instruments) Counter(name string, opts ...InstrumentOption) (asyncfloat64.Counter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.CounterObserverInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncfloat64.Counter{Instrument: inst}, err
}

func (m AsyncFloat64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (asyncfloat64.UpDownCounter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.UpDownCounterObserverInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncfloat64.UpDownCounter{Instrument: inst}, err
}

func (m AsyncFloat64Instruments) Gauge(name string, opts ...InstrumentOption) (asyncfloat64.Gauge, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.GaugeObserverInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncfloat64.Gauge{Instrument: inst}, err
}

func (m AsyncInt64Instruments) Counter(name string, opts ...InstrumentOption) (asyncint64.Counter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.CounterObserverInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncint64.Counter{Instrument: inst}, err
}

func (m AsyncInt64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (asyncint64.UpDownCounter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.UpDownCounterObserverInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncint64.UpDownCounter{Instrument: inst}, err
}

func (m AsyncInt64Instruments) Gauge(name string, opts ...InstrumentOption) (asyncint64.Gauge, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.GaugeObserverInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return asyncint64.Gauge{Instrument: inst}, err
}

func (m SyncFloat64Instruments) Counter(name string, opts ...InstrumentOption) (syncfloat64.Counter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.CounterInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncfloat64.Counter{Instrument: inst}, err
}

func (m SyncFloat64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (syncfloat64.UpDownCounter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.UpDownCounterInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncfloat64.UpDownCounter{Instrument: inst}, err
}

func (m SyncFloat64Instruments) Histogram(name string, opts ...InstrumentOption) (syncfloat64.Histogram, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.HistogramInstrumentKind,
		number.Float64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncfloat64.Histogram{Instrument: inst}, err
}

func (m SyncInt64Instruments) Counter(name string, opts ...InstrumentOption) (syncint64.Counter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.CounterInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncint64.Counter{Instrument: inst}, err
}

func (m SyncInt64Instruments) UpDownCounter(name string, opts ...InstrumentOption) (syncint64.UpDownCounter, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.UpDownCounterInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncint64.UpDownCounter{Instrument: inst}, err
}

func (m SyncInt64Instruments) Histogram(name string, opts ...InstrumentOption) (syncint64.Histogram, error) {
	icfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(
		name,
		sdkapi.HistogramInstrumentKind,
		number.Int64Kind,
		icfg.Description(),
		icfg.Unit(),
	)
	inst, err := m.NewInstrument(desc)
	return syncint64.Histogram{Instrument: inst}, err
}
