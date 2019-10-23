package metric

import (
	"context"

	"go.opentelemetry.io/api/core"
)

type noopMeter struct{}
type noopHandle struct{}
type noopLabelSet struct{}
type noopInstrument struct{}

var _ Meter = noopMeter{}
var _ Instrument = noopInstrument{}
var _ Handle = noopHandle{}
var _ LabelSet = noopLabelSet{}

func (noopHandle) RecordOne(context.Context, MeasurementValue) {
}

func (noopHandle) Release() {
}

func (noopInstrument) AcquireHandle(LabelSet) Handle {
	return noopHandle{}
}

func (noopInstrument) RecordOne(context.Context, MeasurementValue, LabelSet) {
}

func (noopLabelSet) Meter() Meter {
	return noopMeter{}
}

func (noopMeter) Labels(context.Context, ...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (noopMeter) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	return WrapInt64CounterInstrument(noopInstrument{})
}

func (noopMeter) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	return WrapFloat64CounterInstrument(noopInstrument{})
}

func (noopMeter) NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	return WrapInt64GaugeInstrument(noopInstrument{})
}

func (noopMeter) NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	return WrapFloat64GaugeInstrument(noopInstrument{})
}

func (noopMeter) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	return WrapInt64MeasureInstrument(noopInstrument{})
}

func (noopMeter) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	return WrapFloat64MeasureInstrument(noopInstrument{})
}

func (noopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}
