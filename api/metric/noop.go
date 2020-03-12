package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type NoopProvider struct{}

type NoopMeter struct {
}

type noopBoundInstrument struct{}
type noopLabelSet struct{}
type noopInstrument struct{}
type noopInt64Observer struct{}
type noopFloat64Observer struct{}

var _ Provider = NoopProvider{}
var _ Meter = NoopMeter{}
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

func (NoopMeter) NewInt64Counter(name string, cos ...Option) (Int64Counter, error) {
	return WrapInt64CounterInstrument(noopInstrument{}, nil)
}

func (NoopMeter) NewFloat64Counter(name string, cos ...Option) (Float64Counter, error) {
	return WrapFloat64CounterInstrument(noopInstrument{}, nil)
}

func (NoopMeter) NewInt64Measure(name string, mos ...Option) (Int64Measure, error) {
	return WrapInt64MeasureInstrument(noopInstrument{}, nil)
}

func (NoopMeter) NewFloat64Measure(name string, mos ...Option) (Float64Measure, error) {
	return WrapFloat64MeasureInstrument(noopInstrument{}, nil)
}

func (NoopMeter) RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...Option) (Int64Observer, error) {
	return noopInt64Observer{}, nil
}

func (NoopMeter) RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...Option) (Float64Observer, error) {
	return noopFloat64Observer{}, nil
}
