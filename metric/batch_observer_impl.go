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
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type batchObserverImpl struct {
	meter  meterImpl
	runner sdkapi.AsyncBatchRunner
}

// NewInt64GaugeObserver creates a new integer GaugeObserver instrument
// with the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewInt64GaugeObserver(name string, opts ...InstrumentOption) (Int64GaugeObserver, error) {
	if b.runner == nil {
		return wrapInt64GaugeObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64GaugeObserverInstrument(
		b.meter.newAsync(name, sdkapi.GaugeObserverInstrumentKind, number.Int64Kind, opts, b.runner))
}

// NewFloat64GaugeObserver creates a new floating point GaugeObserver with
// the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewFloat64GaugeObserver(name string, opts ...InstrumentOption) (Float64GaugeObserver, error) {
	if b.runner == nil {
		return wrapFloat64GaugeObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64GaugeObserverInstrument(
		b.meter.newAsync(name, sdkapi.GaugeObserverInstrumentKind, number.Float64Kind, opts, b.runner))
}

// NewInt64CounterObserver creates a new integer CounterObserver instrument
// with the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewInt64CounterObserver(name string, opts ...InstrumentOption) (Int64CounterObserver, error) {
	if b.runner == nil {
		return wrapInt64CounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64CounterObserverInstrument(
		b.meter.newAsync(name, sdkapi.CounterObserverInstrumentKind, number.Int64Kind, opts, b.runner))
}

// NewFloat64CounterObserver creates a new floating point CounterObserver with
// the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewFloat64CounterObserver(name string, opts ...InstrumentOption) (Float64CounterObserver, error) {
	if b.runner == nil {
		return wrapFloat64CounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64CounterObserverInstrument(
		b.meter.newAsync(name, sdkapi.CounterObserverInstrumentKind, number.Float64Kind, opts, b.runner))
}

// NewInt64UpDownCounterObserver creates a new integer UpDownCounterObserver instrument
// with the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewInt64UpDownCounterObserver(name string, opts ...InstrumentOption) (Int64UpDownCounterObserver, error) {
	if b.runner == nil {
		return wrapInt64UpDownCounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapInt64UpDownCounterObserverInstrument(
		b.meter.newAsync(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Int64Kind, opts, b.runner))
}

// NewFloat64UpDownCounterObserver creates a new floating point UpDownCounterObserver with
// the given name, running in a batch callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (b batchObserverImpl) NewFloat64UpDownCounterObserver(name string, opts ...InstrumentOption) (Float64UpDownCounterObserver, error) {
	if b.runner == nil {
		return wrapFloat64UpDownCounterObserverInstrument(sdkapi.NewNoopAsyncInstrument(), nil)
	}
	return wrapFloat64UpDownCounterObserverInstrument(
		b.meter.newAsync(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Float64Kind, opts, b.runner))
}
