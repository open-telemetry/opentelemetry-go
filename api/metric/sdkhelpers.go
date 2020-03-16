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

	// Is this _really_ needed?
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

var (
	_ Meter                 = (*wrappedMeterImpl)(nil)
	_ Int64ObserverResult   = (*int64ObserverResult)(nil)
	_ Float64ObserverResult = (*float64ObserverResult)(nil)
)

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

func (m *wrappedMeterImpl) newSynchronous(name string, metricKind Kind, numberKind core.NumberKind, opts []Option) (SynchronousImpl, error) {
	return m.impl.NewSynchronousInstrument(makeDescriptor(name, metricKind, numberKind, opts))
}

func (m *wrappedMeterImpl) NewInt64Counter(name string, opts ...Option) (Int64Counter, error) {
	return WrapInt64CounterInstrument(
		m.newSynchronous(name, CounterKind, core.Int64NumberKind, opts))
}

func WrapInt64CounterInstrument(syncInst SynchronousImpl, err error) (Int64Counter, error) {
	common, err := checkSynchronous(syncInst, err)
	return Int64Counter{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewFloat64Counter(name string, opts ...Option) (Float64Counter, error) {
	return WrapFloat64CounterInstrument(
		m.newSynchronous(name, CounterKind, core.Float64NumberKind, opts))
}

func WrapFloat64CounterInstrument(syncInst SynchronousImpl, err error) (Float64Counter, error) {
	common, err := checkSynchronous(syncInst, err)
	return Float64Counter{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewInt64Measure(name string, opts ...Option) (Int64Measure, error) {
	return WrapInt64MeasureInstrument(
		m.newSynchronous(name, MeasureKind, core.Int64NumberKind, opts))
}

func WrapInt64MeasureInstrument(syncInst SynchronousImpl, err error) (Int64Measure, error) {
	common, err := checkSynchronous(syncInst, err)
	return Int64Measure{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) NewFloat64Measure(name string, opts ...Option) (Float64Measure, error) {
	return WrapFloat64MeasureInstrument(
		m.newSynchronous(name, MeasureKind, core.Float64NumberKind, opts))
}

func WrapFloat64MeasureInstrument(syncInst SynchronousImpl, err error) (Float64Measure, error) {
	common, err := checkSynchronous(syncInst, err)
	return Float64Measure{synchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) newAsynchronous(name string, mkind Kind, nkind core.NumberKind, opts []Option, callback func(func(core.Number, LabelSet))) (AsynchronousImpl, error) {
	return m.impl.NewAsynchronousInstrument(
		makeDescriptor(name, mkind, nkind, opts),
		callback)
}

func (m *wrappedMeterImpl) RegisterInt64Observer(name string, callback Int64ObserverCallback, opts ...Option) (Int64Observer, error) {
	if callback == nil {
		return NoopMeter{}.RegisterInt64Observer("", nil)
	}
	result := &int64ObserverResult{}
	return WrapInt64ObserverInstrument(
		m.newAsynchronous(name, ObserverKind, core.Int64NumberKind, opts,
			func(observe func(core.Number, LabelSet)) {
				result.observe = observe
				callback(result)
			}))
}

func WrapInt64ObserverInstrument(asyncInst AsynchronousImpl, err error) (Int64Observer, error) {
	common, err := checkAsynchronous(asyncInst, err)
	return Int64Observer{asynchronousInstrument: common}, err
}

func (m *wrappedMeterImpl) RegisterFloat64Observer(name string, callback Float64ObserverCallback, opts ...Option) (Float64Observer, error) {
	if callback == nil {
		return NoopMeter{}.RegisterFloat64Observer("", nil)
	}
	result := &float64ObserverResult{}
	return WrapFloat64ObserverInstrument(
		m.newAsynchronous(name, ObserverKind, core.Float64NumberKind, opts,
			func(observe func(core.Number, LabelSet)) {
				result.observe = observe
				callback(result)
			}))
}

func WrapFloat64ObserverInstrument(asyncInst AsynchronousImpl, err error) (Float64Observer, error) {
	common, err := checkAsynchronous(asyncInst, err)
	return Float64Observer{asynchronousInstrument: common}, err
}

func (io *int64ObserverResult) Observe(value int64, labels LabelSet) {
	io.observe(core.NewInt64Number(value), labels)
}

func (fo *float64ObserverResult) Observe(value float64, labels LabelSet) {
	fo.observe(core.NewFloat64Number(value), labels)
}
