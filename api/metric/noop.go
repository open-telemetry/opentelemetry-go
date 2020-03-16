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
type NoopSynchronous struct{ noopInstrument }
type NoopAsynchronous struct{ noopInstrument }

var _ Provider = NoopProvider{}
var _ Meter = NoopMeter{}
var _ SynchronousImpl = NoopSynchronous{}
var _ BoundSynchronousImpl = noopBoundInstrument{}
var _ LabelSet = noopLabelSet{}
var _ AsynchronousImpl = NoopAsynchronous{}

func (NoopProvider) Meter(name string) Meter {
	return NoopMeter{}
}

func (noopInstrument) Interface() interface{} {
	return nil
}

func (noopInstrument) Descriptor() Descriptor {
	return Descriptor{}
}

func (noopBoundInstrument) RecordOne(context.Context, core.Number) {
}

func (noopBoundInstrument) Unbind() {
}

func (NoopSynchronous) Bind(LabelSet) BoundSynchronousImpl {
	return noopBoundInstrument{}
}

func (NoopSynchronous) RecordOne(context.Context, core.Number, LabelSet) {
}

func (NoopMeter) Labels(...core.KeyValue) LabelSet {
	return noopLabelSet{}
}

func (NoopMeter) RecordBatch(context.Context, LabelSet, ...Measurement) {
}

func (NoopMeter) NewInt64Counter(string, ...Option) (Int64Counter, error) {
	return Int64Counter{synchronousInstrument{NoopSynchronous{}}}, nil
}

func (NoopMeter) NewFloat64Counter(string, ...Option) (Float64Counter, error) {
	return Float64Counter{synchronousInstrument{NoopSynchronous{}}}, nil
}

func (NoopMeter) NewInt64Measure(string, ...Option) (Int64Measure, error) {
	return Int64Measure{synchronousInstrument{NoopSynchronous{}}}, nil
}

func (NoopMeter) NewFloat64Measure(string, ...Option) (Float64Measure, error) {
	return Float64Measure{synchronousInstrument{NoopSynchronous{}}}, nil
}

func (NoopMeter) RegisterInt64Observer(string, Int64ObserverCallback, ...Option) (Int64Observer, error) {
	return Int64Observer{asynchronousInstrument{NoopAsynchronous{}}}, nil
}

func (NoopMeter) RegisterFloat64Observer(string, Float64ObserverCallback, ...Option) (Float64Observer, error) {
	return Float64Observer{asynchronousInstrument{NoopAsynchronous{}}}, nil
}
