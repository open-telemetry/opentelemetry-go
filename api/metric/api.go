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

// Options contains some options for metrics of any kind.
type Options struct {
	// Description is an optional field describing the metric
	// instrument.
	Description string
	// Unit is an optional field describing the metric instrument.
	Unit unit.Unit
	// Keys are recommended keys determined in the handles
	// obtained for the metric.
	Keys []core.Key
	// Alternate defines the property of metric value dependent on
	// a metric type.
	//
	// - for Counter, true implies that the metric is an up-down
	//   Counter
	//
	// - for Gauge, true implies that the metric is a
	//   non-descending Gauge
	//
	// - for Measure, true implies that the metric supports
	//   positive and negative values
	Alternate bool
}

// CounterOptionApplier is an interface for applying metric options
// that are valid only for counter metrics.
type CounterOptionApplier interface {
	// ApplyCounterOption is used to make some general or
	// counter-specific changes in the Options.
	ApplyCounterOption(*Options)
}

// GaugeOptionApplier is an interface for applying metric options that
// are valid only for gauge metrics.
type GaugeOptionApplier interface {
	// ApplyGaugeOption is used to make some general or
	// gauge-specific changes in the Options.
	ApplyGaugeOption(*Options)
}

// MeasureOptionApplier is an interface for applying metric options
// that are valid only for measure metrics.
type MeasureOptionApplier interface {
	// ApplyMeasureOption is used to make some general or
	// measure-specific changes in the Options.
	ApplyMeasureOption(*Options)
}

// ObserverOptionApplier is an interface for applying metric options
// that are valid only for observer metrics.
type ObserverOptionApplier interface {
	// ApplyObserverOption is used to make some general or
	// observer-specific changes in the Options.
	ApplyObserverOption(*Options)
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

	// NewInt64Counter creates a new integral counter with a given
	// name and customized with passed options.
	NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter
	// NewFloat64Counter creates a new floating point counter with
	// a given name and customized with passed options.
	NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter
	// NewInt64Gauge creates a new integral gauge with a given
	// name and customized with passed options.
	NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge
	// NewFloat64Gauge creates a new floating point gauge with a
	// given name and customized with passed options.
	NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge
	// NewInt64Measure creates a new integral measure with a given
	// name and customized with passed options.
	NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure
	// NewFloat64Measure creates a new floating point measure with
	// a given name and customized with passed options.
	NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure

	// RegisterInt64Observer creates a new integral observer with a
	// given name, running a given callback, and customized with passed
	// options. Callback can be nil.
	RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer
	// RegisterFloat64Observer creates a new floating point observer
	// with a given name, running a given callback, and customized with
	// passed options. Callback can be nil.
	RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer

	// RecordBatch atomically records a batch of measurements.
	RecordBatch(context.Context, LabelSet, ...Measurement)
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

// Option supports specifying the various metric options.
type Option func(*Options)

// OptionApplier is an interface for applying metric options that are
// valid for all the kinds of metrics.
type OptionApplier interface {
	CounterOptionApplier
	GaugeOptionApplier
	MeasureOptionApplier
	ObserverOptionApplier
	// ApplyOption is used to make some general changes in the
	// Options.
	ApplyOption(*Options)
}

// CounterGaugeObserverOptionApplier is an interface for applying
// metric options that are valid for counter, gauge or observer
// metrics.
type CounterGaugeObserverOptionApplier interface {
	CounterOptionApplier
	GaugeOptionApplier
	ObserverOptionApplier
}

type optionWrapper struct {
	F Option
}

type counterOptionWrapper struct {
	F Option
}

type gaugeOptionWrapper struct {
	F Option
}

type measureOptionWrapper struct {
	F Option
}

type observerOptionWrapper struct {
	F Option
}

type counterGaugeObserverOptionWrapper struct {
	FC Option
	FG Option
	FO Option
}

var (
	_ OptionApplier         = optionWrapper{}
	_ CounterOptionApplier  = counterOptionWrapper{}
	_ GaugeOptionApplier    = gaugeOptionWrapper{}
	_ MeasureOptionApplier  = measureOptionWrapper{}
	_ ObserverOptionApplier = observerOptionWrapper{}
)

func (o optionWrapper) ApplyCounterOption(opts *Options) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyGaugeOption(opts *Options) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyMeasureOption(opts *Options) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyObserverOption(opts *Options) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyOption(opts *Options) {
	o.F(opts)
}

func (o counterOptionWrapper) ApplyCounterOption(opts *Options) {
	o.F(opts)
}

func (o gaugeOptionWrapper) ApplyGaugeOption(opts *Options) {
	o.F(opts)
}

func (o measureOptionWrapper) ApplyMeasureOption(opts *Options) {
	o.F(opts)
}

func (o counterGaugeObserverOptionWrapper) ApplyCounterOption(opts *Options) {
	o.FC(opts)
}

func (o counterGaugeObserverOptionWrapper) ApplyGaugeOption(opts *Options) {
	o.FG(opts)
}

func (o counterGaugeObserverOptionWrapper) ApplyObserverOption(opts *Options) {
	o.FO(opts)
}

func (o observerOptionWrapper) ApplyObserverOption(opts *Options) {
	o.F(opts)
}

// WithDescription applies provided description.
func WithDescription(desc string) OptionApplier {
	return optionWrapper{
		F: func(opts *Options) {
			opts.Description = desc
		},
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) OptionApplier {
	return optionWrapper{
		F: func(opts *Options) {
			opts.Unit = unit
		},
	}
}

// WithKeys applies recommended label keys. Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...core.Key) OptionApplier {
	return optionWrapper{
		F: func(opts *Options) {
			opts.Keys = append(opts.Keys, keys...)
		},
	}
}

// WithMonotonic sets whether a counter, a gauge or an observer is not
// permitted to go down.
func WithMonotonic(monotonic bool) CounterGaugeObserverOptionApplier {
	return counterGaugeObserverOptionWrapper{
		FC: func(opts *Options) {
			opts.Alternate = !monotonic
		},
		FG: func(opts *Options) {
			opts.Alternate = monotonic
		},
		FO: func(opts *Options) {
			opts.Alternate = monotonic
		},
	}
}

// WithAbsolute sets whether a measure is not permitted to be
// negative.
func WithAbsolute(absolute bool) MeasureOptionApplier {
	return measureOptionWrapper{
		F: func(opts *Options) {
			opts.Alternate = !absolute
		},
	}
}
