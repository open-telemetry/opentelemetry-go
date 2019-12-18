package internal

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
)

type meterProvider struct {
	delegate metric.Provider
}

type meter struct {
	provider *meterProvider
	name     string

	delegate metric.Meter
}

type labelSet struct {
	meter  *meter
	labels []core.KeyValue

	delegate metric.LabelSet
}

type instImpl struct {
	meter *meter
	name  string
	opts  interface{}

	delegate metric.InstrumentImpl
}

type instHandle struct {
	inst   *instImpl
	labels metric.LabelSet

	delegate metric.HandleImpl
}

var _ metric.Provider = &meterProvider{}
var _ metric.Meter = &meter{}
var _ metric.LabelSet = &labelSet{}
var _ metric.InstrumentImpl = &instImpl{}
var _ metric.HandleImpl = &instHandle{}

func (p *meterProvider) delegateTo(provider metric.Provider) {
	// HERE YOU ARE @@@
}

func (p *meterProvider) Meter(name string) metric.Meter {
	return &meter{
		provider: p,
		name:     name,
	}
}

func (m *meter) Labels(labels ...core.KeyValue) metric.LabelSet {
	return &labelSet{
		meter:  m,
		labels: labels,
	}
}

func (m *meter) NewInt64Counter(name string, opts ...metric.CounterOptionApplier) metric.Int64Counter {
	return metric.WrapInt64CounterInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (m *meter) NewFloat64Counter(name string, opts ...metric.CounterOptionApplier) metric.Float64Counter {
	return metric.WrapFloat64CounterInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (m *meter) NewInt64Gauge(name string, opts ...metric.GaugeOptionApplier) metric.Int64Gauge {
	return metric.WrapInt64GaugeInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (m *meter) NewFloat64Gauge(name string, opts ...metric.GaugeOptionApplier) metric.Float64Gauge {
	return metric.WrapFloat64GaugeInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (m *meter) NewInt64Measure(name string, opts ...metric.MeasureOptionApplier) metric.Int64Measure {
	return metric.WrapInt64MeasureInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (m *meter) NewFloat64Measure(name string, opts ...metric.MeasureOptionApplier) metric.Float64Measure {
	return metric.WrapFloat64MeasureInstrument(&instImpl{
		meter: m,
		name:  name,
		opts:  opts,
	})
}

func (inst *instImpl) AcquireHandle(labels metric.LabelSet) metric.HandleImpl {
	return &instHandle{
		inst:   inst,
		labels: labels,
	}
}

func (b *instHandle) Release() {

}

func (*meter) RecordBatch(context.Context, metric.LabelSet, ...metric.Measurement) {
	// This is a no-op
}

func (*instImpl) RecordOne(ctx context.Context, number core.Number, labels metric.LabelSet) {
	// This is a no-op
}

func (*instHandle) RecordOne(ctx context.Context, number core.Number) {
	// This is a no-op
}
