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

//go:generate stringer -type=Kind,ValueKind

package metric

import (
	"context"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/unit"
)

// Kind categorizes different kinds of metric.
type Kind int

const (
	Invalid     Kind = iota
	CounterKind      // Supports Add()
	GaugeKind        // Supports Set()
	MeasureKind      // Supports Record()
	ObserverKind
)

// Handle is the implementation-level interface to Set/Add/Record
// individual metrics.
type Handle interface {
	// RecordOne allows the SDK to observe a single metric event
	RecordOne(ctx context.Context, value MeasurementValue)
}

// LabelSet represents a []core.KeyValue for use as pre-defined labels
// in the metrics API.
//
// TODO this belongs outside the metrics API, in some sense, but that
// might create a dependency.  Putting this here means we can't re-use
// a LabelSet between metrics and tracing, even when they are the same
// SDK.
type LabelSet interface {
	Meter() Meter
}

// ObservationCallback defines a type of the callback the observer
// will use to report the measurement
type ObservationCallback func(LabelSet, MeasurementValue)

// ObserverCallback defines a type of the callback SDK will call for
// the registered observers.
type ObserverCallback func(Meter, Observer, ObservationCallback)

// Meter is an interface to the metrics portion of the OpenTelemetry SDK.
type Meter interface {
	// DefineLabels returns a reference to a set of labels that
	// cannot be read by the application.
	DefineLabels(context.Context, ...core.KeyValue) LabelSet

	NewHandle(*Descriptor, LabelSet) Handle
	DeleteHandle(Handle)

	// RecordBatch atomically records a batch of measurements..
	RecordBatch(context.Context, LabelSet, ...Measurement)

	RegisterObserver(Observer, ObserverCallback)
	UnregisterObserver(Observer)
}

type DescriptorID uint64
type ValueKind int8

const (
	AnyValueKind ValueKind = iota
	Int64ValueKind
	Float64ValueKind
)

// Descriptor represents a named metric with recommended
// local-aggregation keys.
type Descriptor struct {
	name        string
	kind        Kind
	keys        []core.Key
	id          DescriptorID
	description string
	unit        unit.Unit
	valueKind   ValueKind
	alternate   bool
}

// Name is a required field describing this metric descriptor, should
// have length > 0.
func (d *Descriptor) Name() string {
	return d.name
}

// Kind is the metric kind of this descriptor.
func (d *Descriptor) Kind() Kind {
	return d.kind
}

// Keys are recommended keys determined in the handles obtained for
// this metric.
func (d *Descriptor) Keys() []core.Key {
	return d.keys
}

// ID is uniquely assigned to support per-SDK registration.
func (d *Descriptor) ID() DescriptorID {
	return d.id
}

// Description is an optional field describing this metric descriptor.
func (d *Descriptor) Description() string {
	return d.description
}

// Unit is an optional field describing this metric descriptor.
func (d *Descriptor) Unit() unit.Unit {
	return d.unit
}

// ValueKind describes the type of values the metric produces.
func (d *Descriptor) ValueKind() ValueKind {
	return d.valueKind
}

// Meaning depends on the metric type:
//
// - for Counter it implies that this is an up-down Counter
//
// - for Gauge/Observer it implies that this is non-descending
//   Gauge/Observer
//
// - for Measure it implies that it supports positive and negative
//   values
func (d *Descriptor) Alternate() bool {
	return d.alternate
}

// Measurement is used for reporting a batch of metric values.
type Measurement struct {
	Descriptor *Descriptor
	Value      MeasurementValue
}

// Option supports specifying the various metric options.
type Option func(*Descriptor)

type OptionApplier interface {
	CounterOptionApplier
	GaugeOptionApplier
	MeasureOptionApplier
	ApplyOption(*Descriptor)
}

type optionWrapper struct {
	F Option
}

var _ OptionApplier = optionWrapper{}

func (o optionWrapper) ApplyCounterOption(d *Descriptor) {
	o.ApplyOption(d)
}

func (o optionWrapper) ApplyGaugeOption(d *Descriptor) {
	o.ApplyOption(d)
}

func (o optionWrapper) ApplyMeasureOption(d *Descriptor) {
	o.ApplyOption(d)
}

func (o optionWrapper) ApplyOption(d *Descriptor) {
	o.F(d)
}

// WithDescription applies provided description.
func WithDescription(desc string) OptionApplier {
	return optionWrapper{
		F: func(d *Descriptor) {
			d.description = desc
		},
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) OptionApplier {
	return optionWrapper{
		F: func(d *Descriptor) {
			d.unit = unit
		},
	}
}

// WithKeys applies required label keys.  Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...core.Key) OptionApplier {
	return optionWrapper{
		F: func(d *Descriptor) {
			d.keys = append(d.keys, keys...)
		},
	}
}

// WithNonMonotonic sets whether a counter is permitted to go up AND down.
func WithNonMonotonic(nm bool) CounterOptionApplier {
	return counterOptionWrapper{
		F: func(d *Descriptor) {
			d.alternate = nm
		},
	}
}

// WithMonotonic sets whether a gauge is not permitted to go down.
func WithMonotonic(m bool) GaugeOptionApplier {
	return gaugeOptionWrapper{
		F: func(d *Descriptor) {
			d.alternate = m
		},
	}
}

// WithSigned sets whether a measure is permitted to be negative.
func WithSigned(s bool) MeasureOptionApplier {
	return measureOptionWrapper{
		F: func(d *Descriptor) {
			d.alternate = s
		},
	}
}

// Defined returns true when the descriptor has been registered.
func (d Descriptor) Defined() bool {
	return len(d.name) != 0
}

// RecordBatch reports to the global Meter.
func RecordBatch(ctx context.Context, labels LabelSet, batch ...Measurement) {
	GlobalMeter().RecordBatch(ctx, labels, batch...)
}

// Int64ObservationCallback defines a type of the callback the
// observer will use to report the int64 measurement.
type Int64ObservationCallback func(LabelSet, int64)

// Int64ObserverCallback defines a type of the callback SDK will call
// for the registered int64 observers.
type Int64ObserverCallback func(Meter, Int64Observer, Int64ObservationCallback)

func RegisterInt64Observer(meter Meter, observer Int64Observer, callback Int64ObserverCallback) {
	cb := func(m Meter, o Observer, ocb ObservationCallback) {
		iocb := func(l LabelSet, i int64) {
			ocb(l, NewInt64MeasurementValue(i))
		}
		callback(m, Int64Observer{o}, iocb)
	}
	meter.RegisterObserver(observer.Observer, cb)
}

func UnregisterInt64Observer(meter Meter, observer Int64Observer) {
	meter.UnregisterObserver(observer.Observer)
}

// Float64ObservationCallback defines a type of the callback the
// observer will use to report the float64 measurement.
type Float64ObservationCallback func(LabelSet, float64)

// Float64ObserverCallback defines a type of the callback SDK will
// call for the registered float64 observers.
type Float64ObserverCallback func(Meter, Float64Observer, Float64ObservationCallback)

func RegisterFloat64Observer(meter Meter, observer Float64Observer, callback Float64ObserverCallback) {
	cb := func(m Meter, o Observer, ocb ObservationCallback) {
		focb := func(l LabelSet, f float64) {
			ocb(l, NewFloat64MeasurementValue(f))
		}
		callback(m, Float64Observer{o}, focb)
	}
	meter.RegisterObserver(observer.Observer, cb)
}

func UnregisterFloat64Observer(meter Meter, observer Float64Observer) {
	meter.UnregisterObserver(observer.Observer)
}
