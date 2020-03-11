// Copyright 2020, OpenTelemetry Authors
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
)

// MeterMust is a wrapper for Meter interfaces that panics when any
// instrument constructor encounters an error.  It implements the
// `MeterMethods` interface but the signature of the constructors is
// different.
type MeterMust struct {
	meter Meter
}

var _ MeterMethods = MeterMust{}

// Must constructs a MeterMust implementation from a Meter, allowing
// the application to panic when any instrument constructor yields an
// error.
func Must(meter Meter) MeterMust {
	return MeterMust{meter: meter}
}

// MeasureConstructorsMust are part of the `MeterMust` API, a Meter
// where `New` methods panic instead of return an error.
type MeasureConstructorsMust interface {
	// NewInt64Counter creates a new integral counter with a given
	// name and customized with passed options.
	NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter
	// NewFloat64Counter creates a new floating point counter with
	// a given name and customized with passed options.
	NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter
	// NewInt64Measure creates a new integral measure with a given
	// name and customized with passed options.
	NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure
	// NewFloat64Measure creates a new floating point measure with
	// a given name and customized with passed options.
	NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure
}

// ObserverConstructorsMust are part of the `MeterMust` API, a `Meter`
// where `Register` methods panic instead or return an error.
type ObserverConstructorsMust interface {
	// RegisterInt64Observer creates a new integral observer with a
	// given name, running a given callback, and customized with passed
	// options. Callback can be nil.
	RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer
	// RegisterFloat64Observer creates a new floating point observer
	// with a given name, running a given callback, and customized with
	// passed options. Callback can be nil.
	RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer
}

// Labels implements MeterMust.
func (mm MeterMust) Labels(kvs ...core.KeyValue) LabelSet {
	return mm.meter.Labels(kvs...)
}

// RecordBatch implements MeterMust.
func (mm MeterMust) RecordBatch(ctx context.Context, ls LabelSet, ms ...Measurement) {
	mm.meter.RecordBatch(ctx, ls, ms...)
}

// NewInt64Counter implements MeterMust.
func (mm MeterMust) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	if inst, err := mm.meter.NewInt64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64Counter implements MeterMust.
func (mm MeterMust) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	if inst, err := mm.meter.NewFloat64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewInt64Measure implements MeterMust.
func (mm MeterMust) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	if inst, err := mm.meter.NewInt64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// NewFloat64Measure implements MeterMust.
func (mm MeterMust) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	if inst, err := mm.meter.NewFloat64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// RegisterInt64Observer implements MeterMust.
func (mm MeterMust) RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer {
	if inst, err := mm.meter.RegisterInt64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// RegisterFloat64Observer implements MeterMust.
func (mm MeterMust) RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer {
	if inst, err := mm.meter.RegisterFloat64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}
