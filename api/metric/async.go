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

package metric

import "go.opentelemetry.io/otel/api/kv"

// The file is organized as follows:
//
//  - Observation type
//  - Three kinds of Observer callback (int64, float64, batch)
//  - Three kinds of Observer result (int64, float64, batch)
//  - Three kinds of Observe() function (int64, float64, batch)
//  - Three kinds of AsyncRunner interface (abstract, single, batch)
//  - Two kinds of Observer constructor (int64, float64)
//  - Two kinds of Observation() function (int64, float64)
//  - Various internals

// Observation is used for reporting an asynchronous  batch of metric
// values. Instances of this type should be created by asynchronous
// instruments (e.g., Int64ValueObserver.Observation()).
type Observation struct {
	// number needs to be aligned for 64-bit atomic operations.
	number     Number
	instrument AsyncImpl
}

// Int64ObserverCallback is a type of callback that integral
// observers run.
type Int64ObserverCallback func(Int64ObserverResult)

// Float64ObserverCallback is a type of callback that floating point
// observers run.
type Float64ObserverCallback func(Float64ObserverResult)

// BatchObserverCallback is a callback argument for use with any
// Observer instrument that will be reported as a batch of
// observations.
type BatchObserverCallback func(BatchObserverResult)

// Int64ObserverResult is passed to an observer callback to capture
// observations for one asynchronous integer metric instrument.
type Int64ObserverResult struct {
	instrument AsyncImpl
	function   func([]kv.KeyValue, ...Observation)
}

// Float64ObserverResult is passed to an observer callback to capture
// observations for one asynchronous floating point metric instrument.
type Float64ObserverResult struct {
	instrument AsyncImpl
	function   func([]kv.KeyValue, ...Observation)
}

// BatchObserverResult is passed to a batch observer callback to
// capture observations for multiple asynchronous instruments.
type BatchObserverResult struct {
	function func([]kv.KeyValue, ...Observation)
}

// Observe captures a single integer value from the associated
// instrument callback, with the given labels.
func (ir Int64ObserverResult) Observe(value int64, labels ...kv.KeyValue) {
	ir.function(labels, Observation{
		instrument: ir.instrument,
		number:     NewInt64Number(value),
	})
}

// Observe captures a single floating point value from the associated
// instrument callback, with the given labels.
func (fr Float64ObserverResult) Observe(value float64, labels ...kv.KeyValue) {
	fr.function(labels, Observation{
		instrument: fr.instrument,
		number:     NewFloat64Number(value),
	})
}

// Observe captures a multiple observations from the associated batch
// instrument callback, with the given labels.
func (br BatchObserverResult) Observe(labels []kv.KeyValue, obs ...Observation) {
	br.function(labels, obs...)
}

// AsyncRunner is expected to convert into an AsyncSingleRunner or an
// AsyncBatchRunner.  SDKs will encounter an error if the AsyncRunner
// does not satisfy one of these interfaces.
type AsyncRunner interface {
	// AnyRunner() is a non-exported method with no functional use
	// other than to make this a non-empty interface.
	AnyRunner()
}

// AsyncSingleRunner is an interface implemented by single-observer
// callbacks.
type AsyncSingleRunner interface {
	// Run accepts a single instrument and function for capturing
	// observations of that instrument.  Each call to the function
	// receives one captured observation.  (The function accepts
	// multiple observations so the same implementation can be
	// used for batch runners.)
	Run(single AsyncImpl, capture func([]kv.KeyValue, ...Observation))

	AsyncRunner
}

// AsyncBatchRunner is an interface implemented by batch-observer
// callbacks.
type AsyncBatchRunner interface {
	// Run accepts a function for capturing observations of
	// multiple instruments.
	Run(capture func([]kv.KeyValue, ...Observation))

	AsyncRunner
}

var _ AsyncSingleRunner = (*Int64ObserverCallback)(nil)
var _ AsyncSingleRunner = (*Float64ObserverCallback)(nil)
var _ AsyncBatchRunner = (*BatchObserverCallback)(nil)

// newInt64AsyncRunner returns a single-observer callback for integer Observer instruments.
func newInt64AsyncRunner(c Int64ObserverCallback) AsyncSingleRunner {
	return &c
}

// newFloat64AsyncRunner returns a single-observer callback for floating point Observer instruments.
func newFloat64AsyncRunner(c Float64ObserverCallback) AsyncSingleRunner {
	return &c
}

// newBatchAsyncRunner returns a batch-observer callback use with multiple Observer instruments.
func newBatchAsyncRunner(c BatchObserverCallback) AsyncBatchRunner {
	return &c
}

// AnyRunner implements AsyncRunner.
func (*Int64ObserverCallback) AnyRunner() {}

// AnyRunner implements AsyncRunner.
func (*Float64ObserverCallback) AnyRunner() {}

// AnyRunner implements AsyncRunner.
func (*BatchObserverCallback) AnyRunner() {}

// Run implements AsyncSingleRunner.
func (i *Int64ObserverCallback) Run(impl AsyncImpl, function func([]kv.KeyValue, ...Observation)) {
	(*i)(Int64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

// Run implements AsyncSingleRunner.
func (f *Float64ObserverCallback) Run(impl AsyncImpl, function func([]kv.KeyValue, ...Observation)) {
	(*f)(Float64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

// Run implements AsyncBatchRunner.
func (b *BatchObserverCallback) Run(function func([]kv.KeyValue, ...Observation)) {
	(*b)(BatchObserverResult{
		function: function,
	})
}

// wrapInt64ValueObserverInstrument returns an `Int64ValueObserver` from a
// `AsyncImpl`.  An error will be generated if the
// `AsyncImpl` is nil (in which case a No-op is substituted),
// otherwise the error passes through.
func wrapInt64ValueObserverInstrument(asyncInst AsyncImpl, err error) (Int64ValueObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return Int64ValueObserver{asyncInstrument: common}, err
}

// wrapFloat64ValueObserverInstrument returns an `Float64ValueObserver` from a
// `AsyncImpl`.  An error will be generated if the
// `AsyncImpl` is nil (in which case a No-op is substituted),
// otherwise the error passes through.
func wrapFloat64ValueObserverInstrument(asyncInst AsyncImpl, err error) (Float64ValueObserver, error) {
	common, err := checkNewAsync(asyncInst, err)
	return Float64ValueObserver{asyncInstrument: common}, err
}
