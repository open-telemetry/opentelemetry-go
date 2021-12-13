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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type meterImpl struct {
	impl sdkapi.MeterImpl
}

// WrapMeterImpl constructs a `Meter` implementation from a
// `MeterImpl` implementation.
func WrapMeterImpl(impl sdkapi.MeterImpl) Meter {
	return meterImpl{
		impl: impl,
	}
}

// RecordBatch atomically records a batch of measurements.
func (m meterImpl) RecordBatch(ctx context.Context, ls []attribute.KeyValue, ms ...Measurement) {
	if m.impl == nil {
		return
	}
	m.impl.RecordBatch(ctx, ls, ms...)
}

// NewBatchObserver creates a new BatchObserver that supports
// making batches of observations for multiple instruments.
func (m meterImpl) NewBatchObserver(callback BatchObserverFunc) BatchObserver {
	return batchObserverImpl{
		meter:  m,
		runner: newBatchAsyncRunner(callback),
	}
}

func (m meterImpl) NewInt64Counter(name string, options ...InstrumentOption) (Int64Counter, error) {
	return wrapInt64CounterInstrument(
		m.newSync(name, sdkapi.CounterInstrumentKind, number.Int64Kind, options))
}

func (m meterImpl) NewFloat64Counter(name string, options ...InstrumentOption) (Float64Counter, error) {
	return wrapFloat64CounterInstrument(
		m.newSync(name, sdkapi.CounterInstrumentKind, number.Float64Kind, options))
}

func (m meterImpl) NewInt64UpDownCounter(name string, options ...InstrumentOption) (Int64UpDownCounter, error) {
	return wrapInt64UpDownCounterInstrument(
		m.newSync(name, sdkapi.UpDownCounterInstrumentKind, number.Int64Kind, options))
}

func (m meterImpl) NewFloat64UpDownCounter(name string, options ...InstrumentOption) (Float64UpDownCounter, error) {
	return wrapFloat64UpDownCounterInstrument(
		m.newSync(name, sdkapi.UpDownCounterInstrumentKind, number.Float64Kind, options))
}

func (m meterImpl) NewInt64Histogram(name string, opts ...InstrumentOption) (Int64Histogram, error) {
	return wrapInt64HistogramInstrument(
		m.newSync(name, sdkapi.HistogramInstrumentKind, number.Int64Kind, opts))
}

func (m meterImpl) NewFloat64Histogram(name string, opts ...InstrumentOption) (Float64Histogram, error) {
	return wrapFloat64HistogramInstrument(
		m.newSync(name, sdkapi.HistogramInstrumentKind, number.Float64Kind, opts))
}

func (m meterImpl) NewInt64GaugeObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64GaugeObserver, error) {
	if callback == nil {
		return wrapInt64GaugeObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64GaugeObserverInstrument(
		m.newAsync(name, sdkapi.GaugeObserverInstrumentKind, number.Int64Kind, opts, newInt64AsyncRunner(callback)))
}

func (m meterImpl) NewFloat64GaugeObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64GaugeObserver, error) {
	if callback == nil {
		return wrapFloat64GaugeObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64GaugeObserverInstrument(
		m.newAsync(name, sdkapi.GaugeObserverInstrumentKind, number.Float64Kind, opts, newFloat64AsyncRunner(callback)))
}

func (m meterImpl) NewInt64CounterObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64CounterObserver, error) {
	if callback == nil {
		return wrapInt64CounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64CounterObserverInstrument(
		m.newAsync(name, sdkapi.CounterObserverInstrumentKind, number.Int64Kind, opts, newInt64AsyncRunner(callback)))
}

func (m meterImpl) NewFloat64CounterObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64CounterObserver, error) {
	if callback == nil {
		return wrapFloat64CounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64CounterObserverInstrument(
		m.newAsync(name, sdkapi.CounterObserverInstrumentKind, number.Float64Kind, opts, newFloat64AsyncRunner(callback)))
}

func (m meterImpl) NewInt64UpDownCounterObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64UpDownCounterObserver, error) {
	if callback == nil {
		return wrapInt64UpDownCounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64UpDownCounterObserverInstrument(
		m.newAsync(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Int64Kind, opts, newInt64AsyncRunner(callback)))
}

func (m meterImpl) NewFloat64UpDownCounterObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64UpDownCounterObserver, error) {
	if callback == nil {
		return wrapFloat64UpDownCounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64UpDownCounterObserverInstrument(
		m.newAsync(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Float64Kind, opts, newFloat64AsyncRunner(callback)))
}

func (m meterImpl) MeterImpl() sdkapi.MeterImpl {
	return m.impl
}

// newAsync constructs one new asynchronous instrument.
func (m meterImpl) newAsync(
	name string,
	mkind sdkapi.InstrumentKind,
	nkind number.Kind,
	opts []InstrumentOption,
	runner sdkapi.AsyncRunner,
) (
	sdkapi.AsyncImpl,
	error,
) {
	if m.impl == nil {
		return sdkapi.NewNoopAsyncInstrument(), nil
	}
	cfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(name, mkind, nkind, cfg.description, cfg.unit)
	return m.impl.NewAsyncInstrument(desc, runner)
}

// newSync constructs one new synchronous instrument.
func (m meterImpl) newSync(
	name string,
	metricKind sdkapi.InstrumentKind,
	numberKind number.Kind,
	opts []InstrumentOption,
) (
	sdkapi.SyncImpl,
	error,
) {
	if m.impl == nil {
		return sdkapi.NewNoopSyncInstrument(), nil
	}
	cfg := NewInstrumentConfig(opts...)
	desc := sdkapi.NewDescriptor(name, metricKind, numberKind, cfg.description, cfg.unit)
	return m.impl.NewSyncInstrument(desc)
}
