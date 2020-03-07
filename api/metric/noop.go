package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type NoopProvider struct{}

type NoopMeter struct {
	NoopMeasureConstructors
	NoopMeasureConstructorsMust
	NoopObserverConstructors
	NoopObserverConstructorsMust
}

type NoopMeasureConstructors struct{}
type NoopMeasureConstructorsMust struct{}
type NoopObserverConstructors struct{}
type NoopObserverConstructorsMust struct{}

type noopBoundInstrument struct{}
type noopLabelSet struct{}
type noopInstrument struct{}
type noopInt64Observer struct{}
type noopFloat64Observer struct{}

var _ Provider = NoopProvider{}
var _ Meter = NoopMeter{}

var _ MeasureConstructors = NoopMeasureConstructors{}
var _ MeasureConstructorsMust = NoopMeasureConstructorsMust{}
var _ ObserverConstructors = NoopObserverConstructors{}
var _ ObserverConstructorsMust = NoopObserverConstructorsMust{}

var _ InstrumentImpl = noopInstrument{}
var _ BoundInstrumentImpl = noopBoundInstrument{}
var _ LabelSet = noopLabelSet{}
var _ Int64Observer = noopInt64Observer{}
var _ Float64Observer = noopFloat64Observer{}

func (NoopProvider) Meter(name string) Meter {
	return NoopMeter{}
}

func (noopBoundInstrument) RecordOne(context.Context, core.Number) {
}

func (noopBoundInstrument) Unbind() {
}

func (noopInstrument) Bind(LabelSet) BoundInstrumentImpl {
	return noopBoundInstrument{}
}

func (noopInstrument) RecordOne(context.Context, core.Number, LabelSet) {
}

func (noopInstrument) Meter() Meter {
	return NoopMeter{}
}

func (noopInt64Observer) Unregister() {
}

func (noopFloat64Observer) Unregister() {
}

func (NoopMeter) Labels(...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (NoopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

// MeasureConstructors

func (NoopMeasureConstructors) NewInt64Counter(name string, cos ...CounterOptionApplier) (Int64Counter, error) {
	return WrapInt64CounterInstrument(noopInstrument{}), nil
}

func (NoopMeasureConstructors) NewFloat64Counter(name string, cos ...CounterOptionApplier) (Float64Counter, error) {
	return WrapFloat64CounterInstrument(noopInstrument{}), nil
}

func (NoopMeasureConstructors) NewInt64Gauge(name string, gos ...GaugeOptionApplier) (Int64Gauge, error) {
	return WrapInt64GaugeInstrument(noopInstrument{}), nil
}

func (NoopMeasureConstructors) NewFloat64Gauge(name string, gos ...GaugeOptionApplier) (Float64Gauge, error) {
	return WrapFloat64GaugeInstrument(noopInstrument{}), nil
}

func (NoopMeasureConstructors) NewInt64Measure(name string, mos ...MeasureOptionApplier) (Int64Measure, error) {
	return WrapInt64MeasureInstrument(noopInstrument{}), nil
}

func (NoopMeasureConstructors) NewFloat64Measure(name string, mos ...MeasureOptionApplier) (Float64Measure, error) {
	return WrapFloat64MeasureInstrument(noopInstrument{}), nil
}

// MeasureConstructorsMust

func (NoopMeasureConstructorsMust) MustNewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return WrapInt64CounterInstrument(noopInstrument{})
}

func (NoopMeasureConstructorsMust) MustNewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return WrapFloat64CounterInstrument(noopInstrument{})
}

func (NoopMeasureConstructorsMust) MustNewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	return WrapInt64GaugeInstrument(noopInstrument{})
}

func (NoopMeasureConstructorsMust) MustNewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	return WrapFloat64GaugeInstrument(noopInstrument{})
}

func (NoopMeasureConstructorsMust) MustNewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return WrapInt64MeasureInstrument(noopInstrument{})
}

func (NoopMeasureConstructorsMust) MustNewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return WrapFloat64MeasureInstrument(noopInstrument{})
}

// ObserverConstructors

func (NoopObserverConstructors) RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) (Int64Observer, error) {
	return noopInt64Observer{}, nil
}

func (NoopObserverConstructors) RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) (Float64Observer, error) {
	return noopFloat64Observer{}, nil
}

// ObserverConstructorsMust

func (NoopObserverConstructorsMust) MustRegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer {
	return noopInt64Observer{}
}

func (NoopObserverConstructorsMust) MustRegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer {
	return noopFloat64Observer{}
}
