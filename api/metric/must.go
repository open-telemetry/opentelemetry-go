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
func (mm MeterMust) NewInt64Counter(name string, cos ...Option) Int64Counter {
	if inst, err := mm.meter.NewInt64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64Counter calls `Meter.NewFloat64Counter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64Counter(name string, cos ...Option) Float64Counter {
	if inst, err := mm.meter.NewFloat64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64UpDownCounter calls `Meter.NewInt64UpDownCounter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64UpDownCounter(name string, cos ...Option) Int64UpDownCounter {
	if inst, err := mm.meter.NewInt64UpDownCounter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64UpDownCounter calls `Meter.NewFloat64UpDownCounter` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64UpDownCounter(name string, cos ...Option) Float64UpDownCounter {
	if inst, err := mm.meter.NewFloat64UpDownCounter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64ValueRecorder calls `Meter.NewInt64ValueRecorder` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewInt64ValueRecorder(name string, mos ...Option) Int64ValueRecorder {
	if inst, err := mm.meter.NewInt64ValueRecorder(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64ValueRecorder calls `Meter.NewFloat64ValueRecorder` and returns the
// instrument, panicking if it encounters an error.
func (mm MeterMust) NewFloat64ValueRecorder(name string, mos ...Option) Float64ValueRecorder {
	if inst, err := mm.meter.NewFloat64ValueRecorder(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// RegisterInt64ValueObserver calls `Meter.RegisterInt64ValueObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) RegisterInt64ValueObserver(name string, callback Int64ObserverCallback, oos ...Option) Int64ValueObserver {
	if inst, err := mm.meter.RegisterInt64ValueObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// RegisterFloat64ValueObserver calls `Meter.RegisterFloat64ValueObserver` and
// returns the instrument, panicking if it encounters an error.
func (mm MeterMust) RegisterFloat64ValueObserver(name string, callback Float64ObserverCallback, oos ...Option) Float64ValueObserver {
	if inst, err := mm.meter.RegisterFloat64ValueObserver(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewBatchObserver returns a wrapper around BatchObserver that panics
// when any instrument constructor returns an error.
func (mm MeterMust) NewBatchObserver(callback BatchObserverCallback) BatchObserverMust {
	return BatchObserverMust{
		batch: mm.meter.NewBatchObserver(callback),
	}
}

// RegisterInt64ValueObserver calls `BatchObserver.RegisterInt64ValueObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) RegisterInt64ValueObserver(name string, oos ...Option) Int64ValueObserver {
	if inst, err := bm.batch.RegisterInt64ValueObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// RegisterFloat64ValueObserver calls `BatchObserver.RegisterFloat64ValueObserver` and
// returns the instrument, panicking if it encounters an error.
func (bm BatchObserverMust) RegisterFloat64ValueObserver(name string, oos ...Option) Float64ValueObserver {
	if inst, err := bm.batch.RegisterFloat64ValueObserver(name, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}
