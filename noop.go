package otel

import (
	"context"
)

type NoopProvider struct{}
type NoopMeter struct{}
type noopHandle struct{}
type noopLabelSet struct{}
type noopInstrument struct{}

var _ Provider = NoopProvider{}
var _ Meter = NoopMeter{}
var _ Instrument = noopInstrument{}
var _ Handle = noopHandle{}
var _ LabelSet = noopLabelSet{}

func (NoopProvider) GetMeter(name string) Meter {
	return NoopMeter{}
}

func (noopHandle) RecordOne(context.Context, Number) {
}

func (noopHandle) Release() {
}

func (noopInstrument) AcquireHandle(LabelSet) Handle {
	return noopHandle{}
}

func (noopInstrument) RecordOne(context.Context, Number, LabelSet) {
}

func (noopInstrument) Meter() Meter {
	return NoopMeter{}
}

func (NoopMeter) Labels(...KeyValue) LabelSet {
	return noopLabelSet{}
}

func (NoopMeter) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return WrapInt64CounterInstrument(noopInstrument{})
}

func (NoopMeter) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return WrapFloat64CounterInstrument(noopInstrument{})
}

func (NoopMeter) NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	return WrapInt64GaugeInstrument(noopInstrument{})
}

func (NoopMeter) NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	return WrapFloat64GaugeInstrument(noopInstrument{})
}

func (NoopMeter) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return WrapInt64MeasureInstrument(noopInstrument{})
}

func (NoopMeter) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return WrapFloat64MeasureInstrument(noopInstrument{})
}

func (NoopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}
