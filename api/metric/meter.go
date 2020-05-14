package metric

import (
	"context"

	"go.opentelemetry.io/otel/api/kv"
)

// Provider supports named Meter instances.
type Provider interface {
	// Meter gets a named Meter interface.  If the name is an
	// empty string, the provider uses a default name.
	Meter(name string) Meter
}

// Meter is the OpenTelemetry metric API, based on a `MeterImpl`
// implementation and the `Meter` library name.
//
// An uninitialized Meter is a no-op implementation.
type Meter struct {
	impl        MeterImpl
	libraryName string
}

// WrapMeterImpl constructs a `Meter` implementation from a
// `MeterImpl` implementation.
func WrapMeterImpl(impl MeterImpl, libraryName string) Meter {
	return Meter{
		impl:        impl,
		libraryName: libraryName,
	}
}

// MeterImpl returns the underlying MeterImpl of this Meter.
func (m Meter) MeterImpl() MeterImpl {
	return m.impl
}

// RecordBatch atomically records a batch of measurements.
func (m Meter) RecordBatch(ctx context.Context, ls []kv.KeyValue, ms ...Measurement) {
	if m.impl == nil {
		return
	}
	m.impl.RecordBatch(ctx, ls, ms...)
}

// NewInt64Counter creates a new integer Counter instrument with the
// given name, customized with options.  May return an error if the
// name is invalid (e.g., empty) or improperly registered (e.g.,
// duplicate registration).
func (m Meter) NewInt64Counter(name string, options ...Option) (Int64Counter, error) {
	return wrapInt64CounterInstrument(
		m.newSync(name, CounterKind, Int64NumberKind, options))
}

// NewFloat64Counter creates a new floating point Counter with the
// given name, customized with options.  May return an error if the
// name is invalid (e.g., empty) or improperly registered (e.g.,
// duplicate registration).
func (m Meter) NewFloat64Counter(name string, options ...Option) (Float64Counter, error) {
	return wrapFloat64CounterInstrument(
		m.newSync(name, CounterKind, Float64NumberKind, options))
}

// NewInt64Measure creates a new integer Measure instrument with the
// given name, customized with options.  May return an error if the
// name is invalid (e.g., empty) or improperly registered (e.g.,
// duplicate registration).
func (m Meter) NewInt64Measure(name string, opts ...Option) (Int64Measure, error) {
	return wrapInt64MeasureInstrument(
		m.newSync(name, MeasureKind, Int64NumberKind, opts))
}

// NewFloat64Measure creates a new floating point Measure with the
// given name, customized with options.  May return an error if the
// name is invalid (e.g., empty) or improperly registered (e.g.,
// duplicate registration).
func (m Meter) NewFloat64Measure(name string, opts ...Option) (Float64Measure, error) {
	return wrapFloat64MeasureInstrument(
		m.newSync(name, MeasureKind, Float64NumberKind, opts))
}

// RegisterInt64Observer creates a new integer Observer instrument
// with the given name, running a given callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (m Meter) RegisterInt64Observer(name string, callback Int64ObserverCallback, opts ...Option) (Int64Observer, error) {
	if callback == nil {
		return wrapInt64ObserverInstrument(NoopAsync{}, nil)
	}
	return wrapInt64ObserverInstrument(
		m.newAsync(name, ObserverKind, Int64NumberKind, opts,
			newInt64AsyncRunner(callback)))
}

// RegisterFloat64Observer creates a new floating point Observer with
// the given name, running a given callback, and customized with
// options.  May return an error if the name is invalid (e.g., empty)
// or improperly registered (e.g., duplicate registration).
func (m Meter) RegisterFloat64Observer(name string, callback Float64ObserverCallback, opts ...Option) (Float64Observer, error) {
	if callback == nil {
		return wrapFloat64ObserverInstrument(NoopAsync{}, nil)
	}
	return wrapFloat64ObserverInstrument(
		m.newAsync(name, ObserverKind, Float64NumberKind, opts,
			newFloat64AsyncRunner(callback)))
}

// NewBatchObserver creates a new BatchObserver that supports
// making batches of observations for multiple instruments.
func (m Meter) NewBatchObserver(callback BatchObserverCallback) BatchObserver {
	return BatchObserver{
		meter:  m,
		runner: newBatchAsyncRunner(callback),
	}
}
