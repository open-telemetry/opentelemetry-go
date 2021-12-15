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

package registry // import "go.opentelemetry.io/otel/internal/metric/registry"

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type instrumentWithType struct {
	impl           interface{}
	instrumentKind sdkapi.InstrumentKind
	numberKind     number.Kind
}

// UniqueInstrumentMeter implements the metric.MeterImpl interface, adding
// uniqueness checking for instrument descriptors.
type UniqueInstrumentMeter struct {
	metric.Meter
	lock  sync.Mutex
	state map[string]instrumentWithType
}

var _ metric.Meter = (*UniqueInstrumentMeter)(nil)

// NewMetricKindMismatchError formats an error that describes a
// mismatched metric instrument definition.
func newMetricKindMismatchError(name string, instrumentKind sdkapi.InstrumentKind, numberKind number.Kind) error {
	return fmt.Errorf("metric %s registered as %s %s: %w",
		name,
		numberKind,
		instrumentKind,
		ErrMetricKindMismatch)
}

// NewUniqueInstrumentMeter returns a wrapped metric.MeterImpl
// with the addition of instrument name uniqueness checking.
func NewUniqueInstrumentMeter(delegate metric.Meter) *UniqueInstrumentMeter {
	return &UniqueInstrumentMeter{
		Meter: delegate,
		state: make(map[string]instrumentWithType),
	}
}

func (m *UniqueInstrumentMeter) NewInt64Counter(name string, options ...metric.InstrumentOption) (metric.Int64Counter, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64Counter(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.CounterInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.CounterInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64Counter), nil
}

func (m *UniqueInstrumentMeter) NewFloat64Counter(name string, options ...metric.InstrumentOption) (metric.Float64Counter, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64Counter(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.CounterInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.CounterInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64Counter), nil
}

func (m *UniqueInstrumentMeter) NewInt64UpDownCounter(name string, options ...metric.InstrumentOption) (metric.Int64UpDownCounter, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64UpDownCounter(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.UpDownCounterInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.UpDownCounterInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64UpDownCounter), nil
}

func (m *UniqueInstrumentMeter) NewFloat64UpDownCounter(name string, options ...metric.InstrumentOption) (metric.Float64UpDownCounter, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64UpDownCounter(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.UpDownCounterInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.UpDownCounterInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64UpDownCounter), nil
}

func (m *UniqueInstrumentMeter) NewInt64Histogram(name string, options ...metric.InstrumentOption) (metric.Int64Histogram, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64Histogram(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.HistogramInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.HistogramInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64Histogram), nil

}

func (m *UniqueInstrumentMeter) NewFloat64Histogram(name string, options ...metric.InstrumentOption) (metric.Float64Histogram, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64Histogram(name, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.HistogramInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.HistogramInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64Histogram), nil
}

func (m *UniqueInstrumentMeter) NewInt64GaugeObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64GaugeObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64GaugeObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.GaugeObserverInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.GaugeObserverInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64GaugeObserver), nil
}

func (m *UniqueInstrumentMeter) NewFloat64GaugeObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64GaugeObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64GaugeObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.GaugeObserverInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.GaugeObserverInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64GaugeObserver), nil
}

func (m *UniqueInstrumentMeter) NewInt64CounterObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64CounterObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64CounterObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.CounterObserverInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.CounterObserverInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64CounterObserver), nil
}

func (m *UniqueInstrumentMeter) NewFloat64CounterObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64CounterObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64CounterObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.CounterObserverInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.CounterObserverInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64CounterObserver), nil
}

func (m *UniqueInstrumentMeter) NewInt64UpDownCounterObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64UpDownCounterObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewInt64UpDownCounterObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.UpDownCounterObserverInstrumentKind, numberKind: number.Int64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Int64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Int64UpDownCounterObserver), nil
}

func (m *UniqueInstrumentMeter) NewFloat64UpDownCounterObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64UpDownCounterObserver, error) {
	inst, ok := m.state[name]
	if !ok {
		syncInst, err := m.Meter.NewFloat64UpDownCounterObserver(name, callback, options...)
		if err != nil {
			return nil, err
		}
		m.state[name] = instrumentWithType{impl: syncInst, instrumentKind: sdkapi.UpDownCounterObserverInstrumentKind, numberKind: number.Float64Kind}
		return syncInst, nil
	}

	if err := compatibleInstrument(name, sdkapi.UpDownCounterObserverInstrumentKind, number.Float64Kind, inst); err != nil {
		return nil, err
	}

	return inst.impl.(metric.Float64UpDownCounterObserver), nil
}

// compatible determines whether two sdkapi.Descriptors are considered
// the same for the purpose of uniqueness checking.
func compatibleInstrument(name string, instrumentKind sdkapi.InstrumentKind, numberKind number.Kind, existing instrumentWithType) error {
	if instrumentKind != existing.instrumentKind || numberKind != existing.numberKind {
		// Return an ErrMetricKindMismatch error if there is a conflict between
		// a descriptor that was already registered and the `descriptor` argument
		return newMetricKindMismatchError(name, sdkapi.CounterInstrumentKind, number.Float64Kind)
	}
	return nil
}
