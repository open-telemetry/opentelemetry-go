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

// MeterMust is a variation on Meter that panics when any instrument
// constructor encounters an error.
type MeterMust interface {
	MeterMethods

	MeasureConstructorsMust
	ObserverConstructorsMust
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

type meterMust struct {
	meter Meter
}

func Must(meter Meter) MeterMust {
	return &meterMust{meter: meter}
}

func (mm *meterMust) Labels(kvs ...core.KeyValue) LabelSet {
	return mm.meter.Labels(kvs...)
}

func (mm *meterMust) RecordBatch(ctx context.Context, ls LabelSet, ms ...Measurement) {
	mm.meter.RecordBatch(ctx, ls, ms...)
}

func (mm *meterMust) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	if inst, err := mm.meter.NewInt64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm *meterMust) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	if inst, err := mm.meter.NewFloat64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm *meterMust) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	if inst, err := mm.meter.NewInt64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm *meterMust) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	if inst, err := mm.meter.NewFloat64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm *meterMust) RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer {
	if inst, err := mm.meter.RegisterInt64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm *meterMust) RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer {
	if inst, err := mm.meter.RegisterFloat64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}
