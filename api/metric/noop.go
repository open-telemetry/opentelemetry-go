package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
)

type NoopProvider struct{}
type NoopMeter struct{}

type noopLabelSet struct{}
type noopInstrument struct{}
type noopBoundInstrument struct{}
type NoopSync struct{ noopInstrument }
type NoopAsync struct{ noopInstrument }

var _ Provider = NoopProvider{}
var _ Meter = NoopMeter{}
var _ SyncImpl = NoopSync{}
var _ BoundSyncImpl = noopBoundInstrument{}
var _ LabelSet = noopLabelSet{}
var _ AsyncImpl = NoopAsync{}

func (NoopProvider) Meter(name string) Meter {
	return NoopMeter{}
}

func (noopInstrument) Implementation() interface{} {
	return nil
}

func (noopInstrument) Descriptor() Descriptor {
	return Descriptor{}
}

func (noopBoundInstrument) RecordOne(context.Context, core.Number) {
}

func (noopBoundInstrument) Unbind() {
}

func (NoopSync) Bind(LabelSet) BoundSyncImpl {
	return noopBoundInstrument{}
}

func (NoopSync) RecordOne(context.Context, core.Number, LabelSet) {
}

func (NoopMeter) Labels(...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (NoopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

func (NoopMeter) NewInt64Counter(string, ...Option) (Int64Counter, error) {
	return Int64Counter{syncInstrument{NoopSync{}}}, nil
}

func (NoopMeter) NewFloat64Counter(string, ...Option) (Float64Counter, error) {
	return Float64Counter{syncInstrument{NoopSync{}}}, nil
}

func (NoopMeter) NewInt64Measure(string, ...Option) (Int64Measure, error) {
	return Int64Measure{syncInstrument{NoopSync{}}}, nil
}

func (NoopMeter) NewFloat64Measure(string, ...Option) (Float64Measure, error) {
	return Float64Measure{syncInstrument{NoopSync{}}}, nil
}

func (NoopMeter) RegisterInt64Observer(string, Int64ObserverCallback, ...Option) (Int64Observer, error) {
	return Int64Observer{asyncInstrument{NoopAsync{}}}, nil
}

func (NoopMeter) RegisterFloat64Observer(string, Float64ObserverCallback, ...Option) (Float64Observer, error) {
	return Float64Observer{asyncInstrument{NoopAsync{}}}, nil
}
