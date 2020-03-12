// Copyright 2019, OpenTelemetry Authors
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
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
)

// Provider supports named Meter instances.
type Provider interface {
	// Meter gets a named Meter interface.  If the name is an
	// empty string, the provider uses a default name.
	Meter(name string) Meter
}

// LabelSet is an implementation-level interface that represents a
// []core.KeyValue for use as pre-defined labels in the metrics API.
type LabelSet interface {
}

// Config contains some options for metrics of any kind.
type Config struct {
	// Description is an optional field describing the metric
	// instrument.
	Description string
	// Unit is an optional field describing the metric instrument.
	Unit unit.Unit
	// Keys are recommended keys determined in the handles
	// obtained for the metric.
	Keys []core.Key
}

// Option is an interface for applying metric options.
type Option interface {
	// Apply is used to set the Option value of a Config.
	Apply(*Config)
}

// Measurement is used for reporting a batch of metric
// values. Instances of this type should be created by instruments
// (e.g., Int64Counter.Measurement()).
type Measurement struct {
	// number needs to be aligned for 64-bit atomic operations.
	number     core.Number
	instrument InstrumentImpl
}

// Instrument returns the instrument that created this measurement.
// This returns an implementation-level object for use by the SDK,
// users should not refer to this.
func (m Measurement) InstrumentImpl() InstrumentImpl {
	return m.instrument
}

// Number returns a number recorded in this measurement.
func (m Measurement) Number() core.Number {
	return m.number
}

// Meter is an interface to the metrics portion of the OpenTelemetry SDK.
type Meter interface {
	// Labels returns a reference to a set of labels that cannot
	// be read by the application.
	Labels(...core.KeyValue) LabelSet

	// RecordBatch atomically records a batch of measurements.
	RecordBatch(context.Context, LabelSet, ...Measurement)

	// All instrument constructors may return an error for
	// conditions such as:
	//   `name` is an empty string
	//   `name` was previously registered as a different kind of instrument
	//          for a given named `Meter`.

	// NewInt64Counter creates a new integral counter with a given
	// name and customized with passed options.
	NewInt64Counter(name string, opts ...Option) (Int64Counter, error)
	// NewFloat64Counter creates a new floating point counter with
	// a given name and customized with passed options.
	NewFloat64Counter(name string, opts ...Option) (Float64Counter, error)
	// NewInt64Measure creates a new integral measure with a given
	// name and customized with passed options.
	NewInt64Measure(name string, opts ...Option) (Int64Measure, error)
	// NewFloat64Measure creates a new floating point measure with
	// a given name and customized with passed options.
	NewFloat64Measure(name string, opts ...Option) (Float64Measure, error)

	// RegisterInt64Observer creates a new integral observer with a
	// given name, running a given callback, and customized with passed
	// options. Callback can be nil.
	RegisterInt64Observer(name string, callback Int64ObserverCallback, opts ...Option) (Int64Observer, error)
	// RegisterFloat64Observer creates a new floating point observer
	// with a given name, running a given callback, and customized with
	// passed options. Callback can be nil.
	RegisterFloat64Observer(name string, callback Float64ObserverCallback, opts ...Option) (Float64Observer, error)
}

// Int64ObserverResult is an interface for reporting integral
// observations.
type Int64ObserverResult interface {
	Observe(value int64, labels LabelSet)
}

// Float64ObserverResult is an interface for reporting floating point
// observations.
type Float64ObserverResult interface {
	Observe(value float64, labels LabelSet)
}

// Int64ObserverCallback is a type of callback that integral
// observers run.
type Int64ObserverCallback func(result Int64ObserverResult)

// Float64ObserverCallback is a type of callback that floating point
// observers run.
type Float64ObserverCallback func(result Float64ObserverResult)

// Int64Observer is a metric that captures a set of int64 values at a
// point in time.
type Int64Observer interface {
	Unregister()
}

// Float64Observer is a metric that captures a set of float64 values
// at a point in time.
type Float64Observer interface {
	Unregister()
}

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return descriptionOption(desc)
}

type descriptionOption string

func (d descriptionOption) Apply(config *Config) {
	config.Description = string(d)
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return unitOption(unit)
}

type unitOption unit.Unit

func (u unitOption) Apply(config *Config) {
	config.Unit = unit.Unit(u)
}

// WithKeys applies recommended label keys. Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...core.Key) Option {
	return keysOption(keys)
}

type keysOption []core.Key

func (k keysOption) Apply(config *Config) {
	config.Keys = append(config.Keys, k...)
}
