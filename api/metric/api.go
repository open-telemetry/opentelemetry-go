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
)

// Recorder is the implementation-level interface to Set/Add/Record individual metrics.
type Recorder interface {
	// Record allows the SDK to observe a single metric event
	Record(ctx context.Context, value float64)
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

// Meter is an interface to the metrics portion of the OpenTelemetry SDK.
type Meter interface {
	// DefineLabels returns a reference to a set of labels that
	// cannot be read by the application.
	DefineLabels(context.Context, ...core.KeyValue) LabelSet

	// RecorderFor returns a handle for observing single measurements.
	RecorderFor(context.Context, LabelSet, Instrument) Recorder

	// RecordSingle records a single measurement without computing a handle.
	RecordSingle(context.Context, LabelSet, Measurement)

	// RecordBatch atomically records a batch of measurements.  An
	// implementation may elect to call `RecordSingle` on each
	// measurement, or it could choose a more-optimized approach.
	RecordBatch(context.Context, LabelSet, ...Measurement)
}

type InstrumentID uint64

// Instrument represents a named metric with recommended local-aggregation keys.
type Instrument struct {
	// Name is a required field describing this metric instrument,
	// should have length > 0.
	Name string

	// ID is uniquely assigned to support per-SDK registration.
	ID InstrumentID

	// Description is an optional field describing this metric instrument.
	Description string

	// Unit is an optional field describing this metric instrument.
	Unit unit.Unit

	// Kind is the metric kind of this instrument.
	Kind Kind

	// NonMonotonic implies this is an up-down Counter.
	NonMonotonic bool

	// Monotonic implies this is a non-descending Gauge.
	Monotonic bool

	// NonNegative implies this is a non-negative Measure.
	NonNegative bool

	// Disabled implies this instrument is disabled by default.
	Disabled bool

	// Keys are required keys determined in the handles
	// obtained for this metric.
	Keys []core.Key
}

// Handle contains a Recorder to support the implementation-defined
// behavior of reporting a single metric with pre-determined label
// values.
type Handle struct {
	Recorder
}

// Measurement is used for reporting a batch of metric values.
type Measurement struct {
	Instrument Instrument
	Value      float64
}

// Option supports specifying the various metric options.
type Option func(*Instrument)

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(inst *Instrument) {
		inst.Description = desc
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(inst *Instrument) {
		inst.Unit = unit
	}
}

// WithNonMonotonic sets whether a counter is permitted to go up AND down.
func WithNonMonotonic(nm bool) Option {
	return func(inst *Instrument) {
		inst.NonMonotonic = nm
	}
}

// WithMonotonic sets whether a gauge is not permitted to go down.
func WithMonotonic(m bool) Option {
	return func(inst *Instrument) {
		inst.Monotonic = m
	}
}

// WithNonNegative sets whether a measure is not permitted to be negative.
func WithNonNegative(non bool) Option {
	return func(inst *Instrument) {
		inst.NonNegative = non
	}
}

// WithDisabled sets whether a measure is disabled by default
func WithDisabled(dis bool) Option {
	return func(inst *Instrument) {
		inst.Disabled = dis
	}
}

// WithKeys applies required label keys.  Multiple `WithKeys`
// options accumulate.
func WithKeys(keys ...core.Key) Option {
	return func(m *Instrument) {
		m.Keys = append(m.Keys, keys...)
	}
}

// Defined returns true when the instrument has been registered.
func (inst Instrument) Defined() bool {
	return len(inst.Name) != 0
}

// RecordSingle reports to the global Meter.
func RecordSingle(ctx context.Context, labels LabelSet, measurement Measurement) {
	GlobalMeter().RecordSingle(ctx, labels, measurement)
}

// RecordBatch reports to the global Meter.
func RecordBatch(ctx context.Context, labels LabelSet, batch ...Measurement) {
	GlobalMeter().RecordBatch(ctx, labels, batch...)
}
