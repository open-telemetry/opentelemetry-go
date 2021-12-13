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

// MeterMust is a wrapper for Meter interfaces that panics when any
// instrument constructor encounters an error.
type MeterMust struct {
	meter Meter
}

// BatchObserverMust is a wrapper for BatchObserver that panics when
// any instrument constructor encounters an error.
type BatchObserverMust struct {
	batch BatchObserver
}

// Must constructs a MeterMust implementation from a Meter, allowing
// the application to panic when any instrument constructor yields an
// error.
func Must(meter Meter) MeterMust {
	return MeterMust{meter: meter}
}

// NewInt64Counter calls `Meter.NewInt64Counter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64Counter(name string, cos ...InstrumentOption) Int64Counter {
	if inst, err := mm.meter.NewInt64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64Counter calls `Meter.NewFloat64Counter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64Counter(name string, cos ...InstrumentOption) Float64Counter {
	if inst, err := mm.meter.NewFloat64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64UpDownCounter calls `Meter.NewInt64UpDownCounter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64UpDownCounter(name string, cos ...InstrumentOption) Int64UpDownCounter {
	if inst, err := mm.meter.NewInt64UpDownCounter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64UpDownCounter calls `Meter.NewFloat64UpDownCounter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64UpDownCounter(name string, cos ...InstrumentOption) Float64UpDownCounter {
	if inst, err := mm.meter.NewFloat64UpDownCounter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64Histogram calls `Meter.NewInt64Histogram` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64Histogram(name string, mos ...InstrumentOption) Int64Histogram {
	if inst, err := mm.meter.NewInt64Histogram(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64Histogram calls `Meter.NewFloat64Histogram` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64Histogram(name string, mos ...InstrumentOption) Float64Histogram {
	if inst, err := mm.meter.NewFloat64Histogram(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64GaugeObserver calls `Meter.NewInt64GaugeObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64GaugeObserver(name string, callback Int64ObserverFunc, oos ...InstrumentOption) Int64GaugeObserver {
	if inst, err := mm.meter.NewInt64GaugeObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64GaugeObserver calls `Meter.NewFloat64GaugeObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64GaugeObserver(name string, callback Float64ObserverFunc, oos ...InstrumentOption) Float64GaugeObserver {
	if inst, err := mm.meter.NewFloat64GaugeObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64CounterObserver calls `Meter.NewInt64CounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64CounterObserver(name string, callback Int64ObserverFunc, oos ...InstrumentOption) Int64CounterObserver {
	if inst, err := mm.meter.NewInt64CounterObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64CounterObserver calls `Meter.NewFloat64CounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64CounterObserver(name string, callback Float64ObserverFunc, oos ...InstrumentOption) Float64CounterObserver {
	if inst, err := mm.meter.NewFloat64CounterObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64UpDownCounterObserver calls `Meter.NewInt64UpDownCounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64UpDownCounterObserver(name string, callback Int64ObserverFunc, oos ...InstrumentOption) Int64UpDownCounterObserver {
	if inst, err := mm.meter.NewInt64UpDownCounterObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64UpDownCounterObserver calls `Meter.NewFloat64UpDownCounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64UpDownCounterObserver(name string, callback Float64ObserverFunc, oos ...InstrumentOption) Float64UpDownCounterObserver {
	if inst, err := mm.meter.NewFloat64UpDownCounterObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewBatchObserver returns a wrapper around BatchObserver that panics
// when any instrument constructor returns an error.
func (mm MeterMust) NewBatchObserver(callback BatchObserverFunc) BatchObserverMust {
	return BatchObserverMust{
		batch: mm.meter.NewBatchObserver(callback),
	}
}

// NewInt64GaugeObserver calls `BatchObserver.NewInt64GaugeObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewInt64GaugeObserver(name string, oos ...InstrumentOption) Int64GaugeObserver {
	if inst, err := bm.batch.NewInt64GaugeObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64GaugeObserver calls `BatchObserver.NewFloat64GaugeObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewFloat64GaugeObserver(name string, oos ...InstrumentOption) Float64GaugeObserver {
	if inst, err := bm.batch.NewFloat64GaugeObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64CounterObserver calls `BatchObserver.NewInt64CounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewInt64CounterObserver(name string, oos ...InstrumentOption) Int64CounterObserver {
	if inst, err := bm.batch.NewInt64CounterObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64CounterObserver calls `BatchObserver.NewFloat64CounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewFloat64CounterObserver(name string, oos ...InstrumentOption) Float64CounterObserver {
	if inst, err := bm.batch.NewFloat64CounterObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64UpDownCounterObserver calls `BatchObserver.NewInt64UpDownCounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewInt64UpDownCounterObserver(name string, oos ...InstrumentOption) Int64UpDownCounterObserver {
	if inst, err := bm.batch.NewInt64UpDownCounterObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64UpDownCounterObserver calls `BatchObserver.NewFloat64UpDownCounterObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) NewFloat64UpDownCounterObserver(name string, oos ...InstrumentOption) Float64UpDownCounterObserver {
	if inst, err := bm.batch.NewFloat64UpDownCounterObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}
