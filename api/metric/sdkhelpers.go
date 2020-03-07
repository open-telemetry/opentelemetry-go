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
)

// LabelSetDelegate is a general-purpose delegating implementation of
// the LabelSet interface.  This is implemented by the default
// Provider returned by api/global.SetMeterProvider(), and should be
// tested for by implementations before converting a LabelSet to their
// private concrete type.
type LabelSetDelegate interface {
	Delegate() LabelSet
}

// InstrumentImpl is the implementation-level interface Set/Add/Record
// individual metrics without precomputed labels.
type InstrumentImpl interface {
	// Bind creates a Bound Instrument to record metrics with
	// precomputed labels.
	Bind(labels LabelSet) BoundInstrumentImpl

	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number, labels LabelSet)
}

// BoundInstrumentImpl is the implementation-level interface to Set/Add/Record
// individual metrics with precomputed labels.
type BoundInstrumentImpl interface {
	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number)

	// Unbind frees the resources associated with this bound instrument. It
	// does not affect the metric this bound instrument was created through.
	Unbind()
}

// MeasureConstructorsMustImpl adds support for the
// MeasureConstructorsMust interface to a Meter implementation.
type MeasureConstructorsMustImpl struct {
	ctors MeasureConstructors
}

// ObserverConstructorsMustImpl adds support for the
// ObserverConstructorsMust interface to a Meter implementation.
type ObserverConstructorsMustImpl struct {
	ctors ObserverConstructors
}

var _ MeasureConstructorsMust = MeasureConstructorsMustImpl{}
var _ ObserverConstructorsMust = ObserverConstructorsMustImpl{}

// WrapInt64CounterInstrument wraps the instrument in the type-safe
// wrapper as an integral counter.
//
// It is mostly intended for SDKs.
func WrapInt64CounterInstrument(instrument InstrumentImpl, err error) (Int64Counter, error) {
	common, err := newCommonMetric(instrument, err)
	return Int64Counter{commonMetric: common}, err
}

// WrapFloat64CounterInstrument wraps the instrument in the type-safe
// wrapper as an floating point counter.
//
// It is mostly intended for SDKs.
func WrapFloat64CounterInstrument(instrument InstrumentImpl, err error) (Float64Counter, error) {
	common, err := newCommonMetric(instrument, err)
	return Float64Counter{commonMetric: common}, err
}

// WrapInt64GaugeInstrument wraps the instrument in the type-safe
// wrapper as an integral gauge.
//
// It is mostly intended for SDKs.
func WrapInt64GaugeInstrument(instrument InstrumentImpl, err error) (Int64Gauge, error) {
	common, err := newCommonMetric(instrument, err)
	return Int64Gauge{commonMetric: common}, err
}

// WrapFloat64GaugeInstrument wraps the instrument in the type-safe
// wrapper as an floating point gauge.
//
// It is mostly intended for SDKs.
func WrapFloat64GaugeInstrument(instrument InstrumentImpl, err error) (Float64Gauge, error) {
	common, err := newCommonMetric(instrument, err)
	return Float64Gauge{commonMetric: common}, err
}

// WrapInt64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an integral measure.
//
// It is mostly intended for SDKs.
func WrapInt64MeasureInstrument(instrument InstrumentImpl, err error) (Int64Measure, error) {
	common, err := newCommonMetric(instrument, err)
	return Int64Measure{commonMetric: common}, err
}

// WrapFloat64MeasureInstrument wraps the instrument in the type-safe
// wrapper as an floating point measure.
//
// It is mostly intended for SDKs.
func WrapFloat64MeasureInstrument(instrument InstrumentImpl, err error) (Float64Measure, error) {
	common, err := newCommonMetric(instrument, err)
	return Float64Measure{commonMetric: common}, err
}

// ApplyCounterOptions is a helper that applies all the counter
// options to passed opts.
func ApplyCounterOptions(opts *Options, cos ...CounterOptionApplier) {
	for _, o := range cos {
		o.ApplyCounterOption(opts)
	}
}

// ApplyGaugeOptions is a helper that applies all the gauge options to
// passed opts.
func ApplyGaugeOptions(opts *Options, gos ...GaugeOptionApplier) {
	for _, o := range gos {
		o.ApplyGaugeOption(opts)
	}
}

// ApplyMeasureOptions is a helper that applies all the measure
// options to passed opts.
func ApplyMeasureOptions(opts *Options, mos ...MeasureOptionApplier) {
	for _, o := range mos {
		o.ApplyMeasureOption(opts)
	}
}

// ApplyObserverOptions is a helper that applies all the observer
// options to passed opts.
func ApplyObserverOptions(opts *Options, mos ...ObserverOptionApplier) {
	for _, o := range mos {
		o.ApplyObserverOption(opts)
	}
}

// MakeMeasureConstructorsMust adds support for
// MeasureConstructorsMust to a Meter.
func MakeMeasureConstructorsMust(ctors MeasureConstructors) MeasureConstructorsMustImpl {
	return MeasureConstructorsMustImpl{
		ctors: ctors,
	}
}

// MakeObserverConstructorsMust adds support for
// ObserverConstructorsMust to a Meter.
func MakeObserverConstructorsMust(ctors ObserverConstructors) ObserverConstructorsMustImpl {
	return ObserverConstructorsMustImpl{
		ctors: ctors,
	}
}

// MeasureConstructorsMustImpl

func (mm MeasureConstructorsMustImpl) MustNewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	if inst, err := mm.ctors.NewInt64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm MeasureConstructorsMustImpl) MustNewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	if inst, err := mm.ctors.NewFloat64Counter(name, cos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm MeasureConstructorsMustImpl) MustNewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	if inst, err := mm.ctors.NewInt64Gauge(name, gos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm MeasureConstructorsMustImpl) MustNewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	if inst, err := mm.ctors.NewFloat64Gauge(name, gos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm MeasureConstructorsMustImpl) MustNewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	if inst, err := mm.ctors.NewInt64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (mm MeasureConstructorsMustImpl) MustNewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	if inst, err := mm.ctors.NewFloat64Measure(name, mos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

// ObserverConstructorsMustImpl

func (om ObserverConstructorsMustImpl) MustRegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer {
	if inst, err := om.ctors.RegisterInt64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}

func (om ObserverConstructorsMustImpl) MustRegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer {
	if inst, err := om.ctors.RegisterFloat64Observer(name, callback, oos...); err != nil {
		panic(err)
	} else {
		return inst
	}
}
