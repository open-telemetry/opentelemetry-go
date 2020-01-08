package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type NoopMeter struct{}
type noopBoundInstrument struct{}
type noopInstrument struct{}

var _ Meter = NoopMeter{}
var _ InstrumentImpl = noopInstrument{}
var _ BoundInstrumentImpl = noopBoundInstrument{}

func (noopBoundInstrument) RecordOne(context.Context, core.Number) {
}

func (noopBoundInstrument) Unbind() {
}

func (noopInstrument) Bind(context.Context, []core.KeyValue) BoundInstrumentImpl {
	return noopBoundInstrument{}
}

func (noopInstrument) RecordOne(context.Context, core.Number, []core.KeyValue) {
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

func (NoopMeter) RecordBatch(context.Context, []core.KeyValue, ...Measurement) {
}
