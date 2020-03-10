package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type NoopProvider struct{}
type NoopMeter struct{}
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

func (NoopMeter) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return WrapInt64CounterInstrument(noopInstrument{})
}

func (NoopMeter) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return WrapFloat64CounterInstrument(noopInstrument{})
}

func (NoopMeter) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return WrapInt64MeasureInstrument(noopInstrument{})
}

func (NoopMeter) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return WrapFloat64MeasureInstrument(noopInstrument{})
}

func (NoopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

func (NoopMeter) RegisterInt64Observer(name string, callback Int64ObserverCallback, oos ...ObserverOptionApplier) Int64Observer {
	return noopInt64Observer{}
}

func (NoopMeter) RegisterFloat64Observer(name string, callback Float64ObserverCallback, oos ...ObserverOptionApplier) Float64Observer {
	return noopFloat64Observer{}
}
