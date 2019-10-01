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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/unit"
)

// Kind categorizes different kinds of metric.
type Kind int

//go:generate stringer -type=Kind
const (
	Invalid     Kind = iota
	CounterKind      // Supports Add()
	GaugeKind        // Supports Set()
	MeasureKind      // Supports Record()
	ObserverKind
)

// Recorder is the implementation-level interface to Set/Add/Record individual metrics.
type Handle interface {
	// Record allows the SDK to observe a single metric event
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

// ObserverCallback defines a type of the callback SDK will call for
// the registered observers.
type ObserverCallback func(Meter, Observer) (LabelSet, MeasurementValue)

// Meter is an interface to the metrics portion of the OpenTelemetry SDK.
type Meter interface {
	// DefineLabels returns a reference to a set of labels that
	// cannot be read by the application.
	DefineLabels(context.Context, ...core.KeyValue) LabelSet

	NewHandle(context.Context, Descriptor, LabelSet) Handle
	DeleteHandle(Handle)

	// RecordBatch atomically records a batch of measurements..
	RecordBatch(context.Context, LabelSet, ...Measurement)

	RegisterObserver(Observer, ObserverCallback)
	UnregisterObserver(Observer)
}

type DescriptorID uint64

// Descriptor represents a named metric with recommended local-aggregation keys.
type Descriptor struct {
	// Name is a required field describing this metric descriptor,
	// should have length > 0.
	Name string

	// Kind is the metric kind of this descriptor.
	Kind Kind

	// Keys are required keys determined in the handles
	// obtained for this metric.
	Keys []core.Key

	// ID is uniquely assigned to support per-SDK registration.
	ID DescriptorID

	// Description is an optional field describing this metric descriptor.
	Description string

	// Unit is an optional field describing this metric descriptor.
	Unit unit.Unit

	// Disabled implies this descriptor is disabled by default.
	Disabled bool

	// NonMonotonic implies this is an up-down Counter.
	NonMonotonic bool

	// Monotonic implies this is a non-descending Gauge/Observer.
	Monotonic bool

	// Signed implies this is a Measure that supports positive and
	// negative values.
	Signed bool
}

// Measurement is used for reporting a batch of metric values.
type Measurement struct {
	Descriptor Descriptor
	Value      MeasurementValue
}

// Option supports specifying the various metric options.
type Option func(*Descriptor)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(d *Descriptor) {
		d.Description = desc
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(d *Descriptor) {
		d.Unit = unit
	}
}

// WithDisabled sets whether a metric is disabled by default
func WithDisabled(dis bool) Option {
	return func(d *Descriptor) {
		d.Disabled = dis
	}
}

// WithKeys applies required label keys.  Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...core.Key) Option {
	return func(d *Descriptor) {
		d.Keys = append(d.Keys, keys...)
	}
}

// WithNonMonotonic sets whether a counter is permitted to go up AND down.
func WithNonMonotonic(nm bool) Option {
	return func(d *Descriptor) {
		d.NonMonotonic = nm
	}
}

// WithMonotonic sets whether a gauge is not permitted to go down.
func WithMonotonic(m bool) Option {
	return func(d *Descriptor) {
		d.Monotonic = m
	}
}

// WithSigned sets whether a measure is permitted to be negative.
func WithSigned(s bool) Option {
	return func(d *Descriptor) {
		d.Signed = s
	}
}

// Defined returns true when the descriptor has been registered.
func (d Descriptor) Defined() bool {
	return len(d.Name) != 0
}

// RecordBatch reports to the global Meter.
func RecordBatch(ctx context.Context, labels LabelSet, batch ...Measurement) {
	GlobalMeter().RecordBatch(ctx, labels, batch...)
}
