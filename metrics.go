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

package otel

import (
	"context"
)

// Provider supports named Meter instances.
type Provider interface {
	// GetMeter gets a named Meter interface.  If the name is an
	// empty string, the provider uses a default name.
	GetMeter(name string) Meter
}

// LabelSet is an implementation-level interface that represents a
// []KeyValue for use as pre-defined labels in the metrics API.
type LabelSet interface{}

// MetricOptions contains some options for metrics of any kind.
type MetricOptions struct {
	// Description is an optional field describing the metric
	// instrument.
	Description string
	// Unit is an optional field describing the metric instrument.
	Unit Unit
	// Keys are recommended keys determined in the handles
	// obtained for the metric.
	Keys []Key
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
	ApplyCounterOption(*MetricOptions)
}

// GaugeOptionApplier is an interface for applying metric options that
// are valid only for gauge metrics.
type GaugeOptionApplier interface {
	// ApplyGaugeOption is used to make some general or
	// gauge-specific changes in the Options.
	ApplyGaugeOption(*MetricOptions)
}

// MeasureOptionApplier is an interface for applying metric options
// that are valid only for measure metrics.
type MeasureOptionApplier interface {
	// ApplyMeasureOption is used to make some general or
	// measure-specific changes in the Options.
	ApplyMeasureOption(*MetricOptions)
}

// Measurement is used for reporting a batch of metric
// values. Instances of this type should be created by instruments
// (e.g., Int64Counter.Measurement()).
type Measurement struct {
	instrument Instrument
	number     Number
}

// Instrument returns the instrument that created this measurement.
// This returns an implementation-level object for use by the SDK,
// users should not refer to this.
func (m Measurement) Instrument() Instrument {
	return m.instrument
}

// Number returns a number recorded in this measurement.
func (m Measurement) Number() Number {
	return m.number
}

// Meter is an interface to the metrics portion of the OpenTelemetry SDK.
type Meter interface {
	// Labels returns a reference to a set of labels that cannot
	// be read by the application.
	Labels(...KeyValue) LabelSet

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

	// RecordBatch atomically records a batch of measurements.
	RecordBatch(context.Context, LabelSet, ...Measurement)
}

// MetricOption supports specifying the various metric options.
type MetricOption func(*MetricOptions)

// OptionApplier is an interface for applying metric options that are
// valid for all the kinds of metrics.
type OptionApplier interface {
	CounterOptionApplier
	GaugeOptionApplier
	MeasureOptionApplier
	// ApplyOption is used to make some general changes in the
	// Options.
	ApplyOption(*MetricOptions)
}

// CounterGaugeOptionApplier is an interface for applying metric
// options that are valid for counter or gauge metrics.
type CounterGaugeOptionApplier interface {
	CounterOptionApplier
	GaugeOptionApplier
}

type optionWrapper struct {
	F MetricOption
}

type counterOptionWrapper struct {
	F MetricOption
}

type gaugeOptionWrapper struct {
	F MetricOption
}

type measureOptionWrapper struct {
	F MetricOption
}

type counterGaugeOptionWrapper struct {
	FC MetricOption
	FG MetricOption
}

var (
	_ OptionApplier        = optionWrapper{}
	_ CounterOptionApplier = counterOptionWrapper{}
	_ GaugeOptionApplier   = gaugeOptionWrapper{}
	_ MeasureOptionApplier = measureOptionWrapper{}
)

func (o optionWrapper) ApplyCounterOption(opts *MetricOptions) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyGaugeOption(opts *MetricOptions) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyMeasureOption(opts *MetricOptions) {
	o.ApplyOption(opts)
}

func (o optionWrapper) ApplyOption(opts *MetricOptions) {
	o.F(opts)
}

func (o counterOptionWrapper) ApplyCounterOption(opts *MetricOptions) {
	o.F(opts)
}

func (o gaugeOptionWrapper) ApplyGaugeOption(opts *MetricOptions) {
	o.F(opts)
}

func (o measureOptionWrapper) ApplyMeasureOption(opts *MetricOptions) {
	o.F(opts)
}

func (o counterGaugeOptionWrapper) ApplyCounterOption(opts *MetricOptions) {
	o.FC(opts)
}

func (o counterGaugeOptionWrapper) ApplyGaugeOption(opts *MetricOptions) {
	o.FG(opts)
}

// WithDescription applies provided description.
func WithDescription(desc string) OptionApplier {
	return optionWrapper{
		F: func(opts *MetricOptions) {
			opts.Description = desc
		},
	}
}

// WithUnit applies provided unit.
func WithUnit(unit Unit) OptionApplier {
	return optionWrapper{
		F: func(opts *MetricOptions) {
			opts.Unit = unit
		},
	}
}

// WithKeys applies recommended label keys. Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...Key) OptionApplier {
	return optionWrapper{
		F: func(opts *MetricOptions) {
			opts.Keys = append(opts.Keys, keys...)
		},
	}
}

// WithMonotonic sets whether a counter or a gauge is not permitted to
// go down.
func WithMonotonic(monotonic bool) CounterGaugeOptionApplier {
	return counterGaugeOptionWrapper{
		FC: func(opts *MetricOptions) {
			opts.Alternate = !monotonic
		},
		FG: func(opts *MetricOptions) {
			opts.Alternate = monotonic
		},
	}
}

// WithAbsolute sets whether a measure is not permitted to be
// negative.
func WithAbsolute(absolute bool) MeasureOptionApplier {
	return measureOptionWrapper{
		F: func(opts *MetricOptions) {
			opts.Alternate = !absolute
		},
	}
}
