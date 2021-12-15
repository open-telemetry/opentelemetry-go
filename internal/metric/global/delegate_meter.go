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

package global // import "go.opentelemetry.io/otel/internal/metric/global"

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

var noopMeter = metric.WrapMeterImpl(nil)

type delegatedInstrument interface {
	setDelegate(meter metric.Meter) error
}

type meterDelegate struct {
	instrumentationName string
	options             []metric.MeterOption

	lock        sync.Mutex
	delegate    unsafe.Pointer // (metric.Meter)
	instruments []delegatedInstrument
}

func newMeterDelegate(instrumentationName string, options ...metric.MeterOption) metric.Meter {
	return &meterDelegate{
		instrumentationName: instrumentationName,
		options:             options,
	}
}

func (m *meterDelegate) loadDelegatePtr() *metric.Meter {
	return (*metric.Meter)(atomic.LoadPointer(&m.delegate))
}

func (m *meterDelegate) RecordBatch(ctx context.Context, ls []attribute.KeyValue, ms ...metric.Measurement) {
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		return
	}
	(*dPtr).RecordBatch(ctx, ls, ms...)
}

func (m *meterDelegate) NewBatchObserver(callback metric.BatchObserverFunc) metric.BatchObserver {
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		// TODO: Correctly implement this.
		return nil
	}
	return (*dPtr).NewBatchObserver(callback)
}

func (m *meterDelegate) NewInt64Counter(name string, options ...metric.InstrumentOption) (metric.Int64Counter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64CounterDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64Counter(name, options...)
}

func (m *meterDelegate) NewFloat64Counter(name string, options ...metric.InstrumentOption) (metric.Float64Counter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64CounterDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64Counter(name, options...)
}

func (m *meterDelegate) NewInt64UpDownCounter(name string, options ...metric.InstrumentOption) (metric.Int64UpDownCounter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64UpDownCounterDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64UpDownCounter(name, options...)
}

func (m *meterDelegate) NewFloat64UpDownCounter(name string, options ...metric.InstrumentOption) (metric.Float64UpDownCounter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64UpDownCounterDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64UpDownCounter(name, options...)
}

func (m *meterDelegate) NewInt64Histogram(name string, options ...metric.InstrumentOption) (metric.Int64Histogram, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64HistogramDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64Histogram(name, options...)
}

func (m *meterDelegate) NewFloat64Histogram(name string, options ...metric.InstrumentOption) (metric.Float64Histogram, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64HistogramDelegate(name, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64Histogram(name, options...)
}

func (m *meterDelegate) NewInt64GaugeObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64GaugeObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64GaugeObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64GaugeObserver(name, callback, options...)
}

func (m *meterDelegate) NewFloat64GaugeObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64GaugeObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64GaugeObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64GaugeObserver(name, callback, options...)
}

func (m *meterDelegate) NewInt64CounterObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64CounterObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64CounterObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64CounterObserver(name, callback, options...)
}

func (m *meterDelegate) NewFloat64CounterObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64CounterObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64CounterObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64CounterObserver(name, callback, options...)
}

func (m *meterDelegate) NewInt64UpDownCounterObserver(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) (metric.Int64UpDownCounterObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newInt64UpDownCounterObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewInt64UpDownCounterObserver(name, callback, options...)
}

func (m *meterDelegate) NewFloat64UpDownCounterObserver(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) (metric.Float64UpDownCounterObserver, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		ret := newFloat64UpDownCounterObserverDelegate(name, callback, options...)
		m.instruments = append(m.instruments, ret.(delegatedInstrument))
		return ret, nil
	}
	return (*dPtr).NewFloat64UpDownCounterObserver(name, callback, options...)
}

func (m *meterDelegate) MeterImpl() sdkapi.MeterImpl {
	dPtr := m.loadDelegatePtr()
	if dPtr == nil {
		noopMeter.MeterImpl()
	}
	return (*dPtr).MeterImpl()
}

func (m *meterDelegate) setDelegate(delegate metric.MeterProvider) {
	m.lock.Lock()
	defer m.lock.Unlock()

	impl := delegate.Meter(m.instrumentationName, m.options...)
	atomic.StorePointer(&m.delegate, unsafe.Pointer(&impl))

	for _, inst := range m.instruments {
		if err := inst.setDelegate(impl); err != nil {
			// TODO: There is no standard way to deliver this error to the user.
			// See https://github.com/open-telemetry/opentelemetry-go/issues/514
			// Note that the default SDK will not generate any errors yet, this is
			// only for added safety.
			panic(err)
		}
	}
}
