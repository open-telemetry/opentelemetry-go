// Copyright The OpenTelemetry Authors
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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

// ErrSDKReturnedNilImpl is returned when a new `MeterImpl` returns nil.
var ErrSDKReturnedNilImpl = errors.New("SDK returned a nil implementation")

// Int64ObserverFunc is a type of callback that integral
// observers run.
type Int64ObserverFunc func(context.Context, Int64ObserverResult)

// Float64ObserverFunc is a type of callback that floating point
// observers run.
type Float64ObserverFunc func(context.Context, Float64ObserverResult)

// BatchObserverFunc is a callback argument for use with any
// Observer instrument that will be reported as a batch of
// observations.
type BatchObserverFunc func(context.Context, BatchObserverResult)

// Int64ObserverResult is passed to an observer callback to capture
// observations for one asynchronous integer metric instrument.
type Int64ObserverResult struct {
	instrument sdkapi.AsyncImpl
	function   func([]attribute.KeyValue, ...Observation)
}

// Float64ObserverResult is passed to an observer callback to capture
// observations for one asynchronous floating point metric instrument.
type Float64ObserverResult struct {
	instrument sdkapi.AsyncImpl
	function   func([]attribute.KeyValue, ...Observation)
}

// BatchObserverResult is passed to a batch observer callback to
// capture observations for multiple asynchronous instruments.
type BatchObserverResult struct {
	function func([]attribute.KeyValue, ...Observation)
}

// Observe captures a single integer value from the associated
// instrument callback, with the given labels.
func (ir Int64ObserverResult) Observe(value int64, labels ...attribute.KeyValue) {
	ir.function(labels, sdkapi.NewObservation(ir.instrument, number.NewInt64Number(value)))
}

// Observe captures a single floating point value from the associated
// instrument callback, with the given labels.
func (fr Float64ObserverResult) Observe(value float64, labels ...attribute.KeyValue) {
	fr.function(labels, sdkapi.NewObservation(fr.instrument, number.NewFloat64Number(value)))
}

// Observe captures a multiple observations from the associated batch
// instrument callback, with the given labels.
func (br BatchObserverResult) Observe(labels []attribute.KeyValue, obs ...Observation) {
	br.function(labels, obs...)
}

var _ sdkapi.AsyncSingleRunner = (*Int64ObserverFunc)(nil)
var _ sdkapi.AsyncSingleRunner = (*Float64ObserverFunc)(nil)
var _ sdkapi.AsyncBatchRunner = (*BatchObserverFunc)(nil)

// newInt64AsyncRunner returns a single-observer callback for integer Observer instruments.
func newInt64AsyncRunner(c Int64ObserverFunc) sdkapi.AsyncSingleRunner {
	return &c
}

// newFloat64AsyncRunner returns a single-observer callback for floating point Observer instruments.
func newFloat64AsyncRunner(c Float64ObserverFunc) sdkapi.AsyncSingleRunner {
	return &c
}

// newBatchAsyncRunner returns a batch-observer callback use with multiple Observer instruments.
func newBatchAsyncRunner(c BatchObserverFunc) sdkapi.AsyncBatchRunner {
	return &c
}

// AnyRunner implements AsyncRunner.
func (*Int64ObserverFunc) AnyRunner() {}

// AnyRunner implements AsyncRunner.
func (*Float64ObserverFunc) AnyRunner() {}

// AnyRunner implements AsyncRunner.
func (*BatchObserverFunc) AnyRunner() {}

// Run implements AsyncSingleRunner.
func (i *Int64ObserverFunc) Run(ctx context.Context, impl sdkapi.AsyncImpl, function func([]attribute.KeyValue, ...Observation)) {
	(*i)(ctx, Int64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

// Run implements AsyncSingleRunner.
func (f *Float64ObserverFunc) Run(ctx context.Context, impl sdkapi.AsyncImpl, function func([]attribute.KeyValue, ...Observation)) {
	(*f)(ctx, Float64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

// Run implements AsyncBatchRunner.
func (b *BatchObserverFunc) Run(ctx context.Context, function func([]attribute.KeyValue, ...Observation)) {
	(*b)(ctx, BatchObserverResult{
		function: function,
	})
}

// syncInstrument contains a SyncImpl.
type syncInstrument interface {
	// SyncImpl returns the implementation object for synchronous instruments.
	SyncImpl() sdkapi.SyncImpl
}

// asyncInstrument contains a AsyncImpl.
type asyncInstrument interface {
	// AsyncImpl returns the implementation object for asynchronous instruments.
	AsyncImpl() sdkapi.AsyncImpl
}

// Int64GaugeObserver is a metric that captures a set of int64 values at a
// point in time.
//
// Warning: methods may be added to this interface in minor releases.
type Int64GaugeObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v int64) Observation
}

// Float64GaugeObserver is a metric that captures a set of float64 values
// at a point in time.
//
// Warning: methods may be added to this interface in minor releases.
type Float64GaugeObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v float64) Observation
}

// Int64CounterObserver is a metric that captures a precomputed sum of
// int64 values at a point in time.
//
// Warning: methods may be added to this interface in minor releases.
type Int64CounterObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v int64) Observation
}

// Float64CounterObserver is a metric that captures a precomputed sum of
// float64 values at a point in time.
//
// Warning: methods may be added to this interface in minor releases.
type Float64CounterObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v float64) Observation
}

// Int64UpDownCounterObserver is a metric that captures a precomputed sum of
// int64 values at a point in time.
type Int64UpDownCounterObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v int64) Observation
}

// Float64UpDownCounterObserver is a metric that captures a precomputed sum of
// float64 values at a point in time.
//
// Warning: methods may be added to this interface in minor releases.
type Float64UpDownCounterObserver interface {
	asyncInstrument

	// Observation returns an Observation, a BatchObserverFunc
	// argument, for an asynchronous integer instrument.
	// This returns an implementation-level object for use by the SDK,
	// users should not refer to this.
	Observation(v float64) Observation
}

// Float64Counter is a metric that accumulates float64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64Counter interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch
	// recording.
	Measurement(value float64) Measurement

	// Add adds the value to the counter's sum. The labels should contain
	// the keys and values to be associated with this value.
	Add(ctx context.Context, value float64, labels ...attribute.KeyValue)
}

// Int64Counter is a metric that accumulates int64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64Counter interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch
	// recording.
	Measurement(value int64) Measurement

	// Add adds the value to the counter's sum. The labels should contain
	// the keys and values to be associated with this value.
	Add(ctx context.Context, value int64, labels ...attribute.KeyValue)
}

// Float64UpDownCounter is a metric instrument that sums floating
// point values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64UpDownCounter interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch recording.
	Measurement(value float64) Measurement

	// Add adds the value to the counter's sum. The labels should contain
	// the keys and values to be associated with this value.
	Add(ctx context.Context, value float64, labels ...attribute.KeyValue)
}

// Int64UpDownCounter is a metric instrument that sums integer values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64UpDownCounter interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch
	// recording.
	Measurement(value int64) Measurement

	// Add adds the value to the counter's sum. The labels should contain
	// the keys and values to be associated with this value.
	Add(ctx context.Context, value int64, labels ...attribute.KeyValue)
}

// Float64Histogram is a metric that records float64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Float64Histogram interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch
	// recording.
	Measurement(value float64) Measurement

	// Record adds a new value to the list of Histogram's records. The
	// labels should contain the keys and values to be associated with
	// this value.
	Record(ctx context.Context, value float64, labels ...attribute.KeyValue)
}

// Int64Histogram is a metric that records int64 values.
//
// Warning: methods may be added to this interface in minor releases.
type Int64Histogram interface {
	syncInstrument

	// Measurement creates a Measurement object to use with batch
	// recording.
	Measurement(value int64) Measurement

	// Record adds a new value to the Histogram's distribution. The
	// labels should contain the keys and values to be associated with
	// this value.
	Record(ctx context.Context, value int64, labels ...attribute.KeyValue)
}
