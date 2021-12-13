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
	"go.opentelemetry.io/otel/metric/sdkapi"
)

// MeterProvider supports named Meter instances.
type MeterProvider interface {
	// Meter creates an implementation of the Meter interface.
	// The instrumentationName must be the name of the library providing
	// instrumentation. This name may be the same as the instrumented code
	// only if that code provides built-in instrumentation. If the
	// instrumentationName is empty, then a implementation defined default
	// name will be used instead.
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter is the creator of metric instruments.
//
// Warning: methods may be added to this interface in minor releases.
type Meter interface {
	// RecordBatch atomically records a batch of measurements.
	RecordBatch(ctx context.Context, ls []attribute.KeyValue, ms ...Measurement)

	// NewBatchObserver creates a new BatchObserver that supports
	// making batches of observations for multiple instruments.
	NewBatchObserver(callback BatchObserverFunc) BatchObserver

	// NewInt64Counter creates a new integer Counter instrument with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewInt64Counter(name string, options ...InstrumentOption) (Int64Counter, error)

	// NewFloat64Counter creates a new floating point Counter with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewFloat64Counter(name string, options ...InstrumentOption) (Float64Counter, error)

	// NewInt64UpDownCounter creates a new integer UpDownCounter instrument with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewInt64UpDownCounter(name string, options ...InstrumentOption) (Int64UpDownCounter, error)

	// NewFloat64UpDownCounter creates a new floating point UpDownCounter with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewFloat64UpDownCounter(name string, options ...InstrumentOption) (Float64UpDownCounter, error)

	// NewInt64Histogram creates a new integer Histogram instrument with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewInt64Histogram(name string, opts ...InstrumentOption) (Int64Histogram, error)

	// NewFloat64Histogram creates a new floating point Histogram with the
	// given name, customized with options.  May return an error if the
	// name is invalid (e.g., empty) or improperly registered (e.g.,
	// duplicate registration).
	NewFloat64Histogram(name string, opts ...InstrumentOption) (Float64Histogram, error)

	// NewInt64GaugeObserver creates a new integer GaugeObserver instrument
	// with the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64GaugeObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64GaugeObserver, error)

	// NewFloat64GaugeObserver creates a new floating point GaugeObserver with
	// the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64GaugeObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64GaugeObserver, error)

	// NewInt64CounterObserver creates a new integer CounterObserver instrument
	// with the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64CounterObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64CounterObserver, error)

	// NewFloat64CounterObserver creates a new floating point CounterObserver with
	// the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64CounterObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64CounterObserver, error)

	// NewInt64UpDownCounterObserver creates a new integer UpDownCounterObserver instrument
	// with the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64UpDownCounterObserver(name string, callback Int64ObserverFunc, opts ...InstrumentOption) (Int64UpDownCounterObserver, error)

	// NewFloat64UpDownCounterObserver creates a new floating point UpDownCounterObserver with
	// the given name, running a given callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64UpDownCounterObserver(name string, callback Float64ObserverFunc, opts ...InstrumentOption) (Float64UpDownCounterObserver, error)

	// MeterImpl returns the underlying MeterImpl of this Meter.
	MeterImpl() sdkapi.MeterImpl
}

// BatchObserver represents an Observer callback that can report
// observations for multiple instruments.
//
// Warning: methods may be added to this interface in minor releases.
type BatchObserver interface {
	// NewInt64GaugeObserver creates a new integer GaugeObserver instrument
	// with the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64GaugeObserver(name string, opts ...InstrumentOption) (Int64GaugeObserver, error)

	// NewFloat64GaugeObserver creates a new floating point GaugeObserver with
	// the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64GaugeObserver(name string, opts ...InstrumentOption) (Float64GaugeObserver, error)

	// NewInt64CounterObserver creates a new integer CounterObserver instrument
	// with the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64CounterObserver(name string, opts ...InstrumentOption) (Int64CounterObserver, error)

	// NewFloat64CounterObserver creates a new floating point CounterObserver with
	// the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64CounterObserver(name string, opts ...InstrumentOption) (Float64CounterObserver, error)

	// NewInt64UpDownCounterObserver creates a new integer UpDownCounterObserver instrument
	// with the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewInt64UpDownCounterObserver(name string, opts ...InstrumentOption) (Int64UpDownCounterObserver, error)

	// NewFloat64UpDownCounterObserver creates a new floating point UpDownCounterObserver with
	// the given name, running in a batch callback, and customized with
	// options.  May return an error if the name is invalid (e.g., empty)
	// or improperly registered (e.g., duplicate registration).
	NewFloat64UpDownCounterObserver(name string, opts ...InstrumentOption) (Float64UpDownCounterObserver, error)
}

// Measurement is used for reporting a synchronous batch of metric
// values. Instances of this type should be created by synchronous
// instruments (e.g., Int64Counter.Measurement()).
//
// Note: This is an alias because it is a first-class member of the
// API but is also part of the lower-level sdkapi interface.
type Measurement = sdkapi.Measurement

// Observation is used for reporting an asynchronous  batch of metric
// values. Instances of this type should be created by asynchronous
// instruments (e.g., Int64GaugeObserver.Observation()).
//
// Note: This is an alias because it is a first-class member of the
// API but is also part of the lower-level sdkapi interface.
type Observation = sdkapi.Observation
