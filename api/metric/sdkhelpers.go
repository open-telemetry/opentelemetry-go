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

type MeterImpl interface {
	// Labels returns a reference to a set of labels that cannot
	// be read by the application.
	Labels(...core.KeyValue) LabelSet

	// RecordBatch atomically records a batch of measurements.
	RecordBatch(context.Context, LabelSet, ...Measurement)

	NewSynchronousInstrument(descriptor Descriptor) (SynchronousImpl, error)

	NewAsynchronousInstrument(descriptor Descriptor, callback func(func(core.Number, LabelSet))) (AsynchronousImpl, error)
}

// LabelSetDelegate is a general-purpose delegating implementation of
// the LabelSet interface.  This is implemented by the default
// Provider returned by api/global.SetMeterProvider(), and should be
// tested for by implementations before converting a LabelSet to their
// private concrete type.
type LabelSetDelegate interface {
	Delegate() LabelSet
}

type InstrumentImpl interface {
	Interface() interface{}
	Descriptor() Descriptor
}

// SynchronousImpl is the implementation-level interface Set/Add/Record
// individual metrics without precomputed labels.
type SynchronousImpl interface {
	InstrumentImpl

	// Bind creates a Bound Instrument to record metrics with
	// precomputed labels.
	Bind(labels LabelSet) BoundSynchronousImpl

	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number, labels LabelSet)
}

// BoundSynchronousImpl is the implementation-level interface to Set/Add/Record
// individual metrics with precomputed labels.
type BoundSynchronousImpl interface {

	// RecordOne allows the SDK to observe a single metric event.
	RecordOne(ctx context.Context, number core.Number)

	// Unbind frees the resources associated with this bound instrument. It
	// does not affect the metric this bound instrument was created through.
	Unbind()
}

type AsynchronousImpl interface {
	InstrumentImpl

	Unregister()
}

type wrappedMeterImpl struct {
	impl MeterImpl
}

type int64ObserverResult struct {
	observe func(core.Number, LabelSet)
}

type float64ObserverResult struct {
	observe func(core.Number, LabelSet)
}

var _ Int64ObserverResult = int64ObserverResult{}
var _ Float64ObserverResult = float64ObserverResult{}
var _ Meter = (*wrappedMeterImpl)(nil)

// Configure is a helper that applies all the options to a Config.
func Configure(opts []Option) Config {
	var config Config
	for _, o := range opts {
		o.Apply(&config)
	}
	return config
}

func WrapMeterImpl(impl MeterImpl) Meter {
	return &wrappedMeterImpl{
		impl: impl,
	}
}

func UnwrapImpl(meter Meter) (MeterImpl, bool) {
	if wrap, ok := meter.(*wrappedMeterImpl); ok {
		return wrap.impl, true
	}
	return nil, false
}

func (m *wrappedMeterImpl) Labels(labels ...core.KeyValue) LabelSet {
	return m.impl.Labels(labels...)
}

func (m *wrappedMeterImpl) RecordBatch(ctx context.Context, ls LabelSet, ms ...Measurement) {
	m.impl.RecordBatch(ctx, ls, ms...)
}

func makeDescriptor(name string, metricKind Kind, numberKind core.NumberKind, opts []Option) Descriptor {
	return Descriptor{
		Name:       name,
		Kind:       metricKind,
		NumberKind: numberKind,
		Config:     Configure(opts),
	}
}

func (m *wrappedMeterImpl) NewInt64Counter(name string, opts ...Option) (Int64Counter, error) {
	common, err := m.makeSynchronous(makeDescriptor(name, CounterKind, core.Int64NumberKind, opts))
	return Int64Counter{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewFloat64Counter(name string, opts ...Option) (Float64Counter, error) {
	common, err := m.makeSynchronous(makeDescriptor(name, CounterKind, core.Float64NumberKind, opts))
	return Float64Counter{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewInt64Measure(name string, opts ...Option) (Int64Measure, error) {
	common, err := m.makeSynchronous(makeDescriptor(name, MeasureKind, core.Int64NumberKind, opts))
	return Int64Measure{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewFloat64Measure(name string, opts ...Option) (Float64Measure, error) {
	common, err := m.makeSynchronous(makeDescriptor(name, MeasureKind, core.Float64NumberKind, opts))
	return Float64Measure{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) RegisterInt64Observer(name string, callback Int64ObserverCallback, opts ...Option) (Int64Observer, error) {
	if callback == nil {
		return NoopMeter{}.RegisterInt64Observer(name, callback, opts...)
	}
	common, err := m.makeAsynchronous(
		makeDescriptor(name, ObserverKind, core.Int64NumberKind, opts),
		func(observe func(core.Number, LabelSet)) {
			callback(int64ObserverResult{observe})
		},
	)
	return Int64Observer{asynchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) RegisterFloat64Observer(name string, callback Float64ObserverCallback, opts ...Option) (Float64Observer, error) {
	if callback == nil {
		return NoopMeter{}.RegisterFloat64Observer(name, callback, opts...)
	}
	common, err := m.makeAsynchronous(
		makeDescriptor(name, ObserverKind, core.Float64NumberKind, opts),
		func(observe func(core.Number, LabelSet)) {
			callback(float64ObserverResult{observe})
		},
	)
	return Float64Observer{asynchronousInstrument: common}, err
}

func (io int64ObserverResult) Observe(value int64, labels LabelSet) {
	io.observe(core.NewInt64Number(value), labels)
}

func (fo float64ObserverResult) Observe(value float64, labels LabelSet) {
	fo.observe(core.NewFloat64Number(value), labels)
}
