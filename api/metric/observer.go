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

import (
	"go.opentelemetry.io/otel/api/core"
)

// Int64ObserverCallback is a callback argument for
// RegisterInt64Observer.  If reporting a single instrument, use
// NewInt64ObserverCallback().  If reporting a batch of instruments,
// use NewBatchObserverCallback().
type Int64ObserverCallback interface {
	intObserver()
	anyRunner()
}

// Float64ObserverCallback is a callback argument for
// RegisterFloat64Observer.  If reporting a single instrument, use
// NewFloat64ObserverCallback().  If reporting a batch of instruments,
// use NewBatchObserverCallback().
type Float64ObserverCallback interface {
	floatObserver()
	anyRunner()
}

// BatchObserverCallback is a callback argument for use with any
// Observer instrument that will be reported as a batch of
// observations.
type BatchObserverCallback interface {
	intObserver()
	floatObserver()
	anyRunner()
}

// NewInt64ObserverCallback returns a single-observer callback for integer Observer instruments.
func NewInt64ObserverCallback(f func(Int64ObserverResult)) Int64ObserverCallback {
	var c int64ObserverCallback = f
	return &c
}

// NewFloat64ObserverCallback returns a single-observer callback for floating point Observer instruments.
func NewFloat64ObserverCallback(f func(Float64ObserverResult)) Float64ObserverCallback {
	var c float64ObserverCallback = f
	return &c
}

// NewBatchObserverCallback returns a batch-observer callback use with multiple Observer instruments.
func NewBatchObserverCallback(f func(BatchObserverResult)) BatchObserverCallback {
	var c batchObserverCallback = f
	return &c
}

// Int64ObserverResult is passed to Int64ObserverCallback for
// reporting results from a single integer Observer instrument.
type Int64ObserverResult interface {
	Observe(value int64, labels ...core.KeyValue)
}

// Float64ObserverResult is passed to Float64ObserverCallback for
// reporting results from a single floating point Observer instrument.
type Float64ObserverResult interface {
	Observe(value float64, labels ...core.KeyValue)
}

// BatchObserverResult is passed to BatchObserverCallback for
// reporting results for a batch of Observers.
type BatchObserverResult interface {
	Observe(labels []core.KeyValue, observations ...Observation)
}

// AsyncRunner is an interface implemented by all Observer callbacks.
// SDKs should test whether the runner is an AsyncSingleRunner or
// AsyncBatchRunner and use the correct interface to run the observers.
type AsyncRunner interface {
	anyRunner()
}

// AsyncSingleRunner is an interface implemented by single-observer
// callbacks.
type AsyncSingleRunner interface {
	Run(AsyncImpl, func([]core.KeyValue, Observation))
}

// AsyncBatchRunner is an interface implemented by batch-observer
// callbacks.
type AsyncBatchRunner interface {
	Run(func([]core.KeyValue, []Observation))
}

type int64ObserverResult struct {
	instrument AsyncImpl
	function   func([]core.KeyValue, Observation)
}

type float64ObserverResult struct {
	instrument AsyncImpl
	function   func([]core.KeyValue, Observation)
}

type batchObserverResult struct {
	function func([]core.KeyValue, []Observation)
}

var _ Int64ObserverResult = int64ObserverResult{}
var _ Float64ObserverResult = float64ObserverResult{}
var _ BatchObserverResult = batchObserverResult{}

func (ir int64ObserverResult) Observe(value int64, labels ...core.KeyValue) {
	ir.function(labels, Observation{
		instrument: ir.instrument,
		number:     core.NewInt64Number(value),
	})
}

func (fr float64ObserverResult) Observe(value float64, labels ...core.KeyValue) {
	fr.function(labels, Observation{
		instrument: fr.instrument,
		number:     core.NewFloat64Number(value),
	})
}

func (br batchObserverResult) Observe(labels []core.KeyValue, obs ...Observation) {
	br.function(labels, obs)
}

type int64ObserverCallback func(result Int64ObserverResult)
type float64ObserverCallback func(result Float64ObserverResult)
type batchObserverCallback func(result BatchObserverResult)

var _ Int64ObserverCallback = (*int64ObserverCallback)(nil)
var _ Int64ObserverCallback = (*batchObserverCallback)(nil)
var _ Float64ObserverCallback = (*float64ObserverCallback)(nil)
var _ Float64ObserverCallback = (*batchObserverCallback)(nil)
var _ AsyncSingleRunner = (*int64ObserverCallback)(nil)
var _ AsyncSingleRunner = (*float64ObserverCallback)(nil)
var _ AsyncBatchRunner = (*batchObserverCallback)(nil)

func (*int64ObserverCallback) intObserver()     {}
func (*int64ObserverCallback) anyRunner()       {}
func (*float64ObserverCallback) floatObserver() {}
func (*float64ObserverCallback) anyRunner()     {}
func (*batchObserverCallback) intObserver()     {}
func (*batchObserverCallback) floatObserver()   {}
func (*batchObserverCallback) anyRunner()       {}

func (i *int64ObserverCallback) Run(impl AsyncImpl, function func([]core.KeyValue, Observation)) {
	(*i)(int64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

func (f *float64ObserverCallback) Run(impl AsyncImpl, function func([]core.KeyValue, Observation)) {
	(*f)(float64ObserverResult{
		instrument: impl,
		function:   function,
	})
}

func (b *batchObserverCallback) Run(function func([]core.KeyValue, []Observation)) {
	(*b)(batchObserverResult{
		function: function,
	})
}

// Int64Observer is a metric that captures a set of int64 values at a
// point in time.
type Int64Observer struct {
	asyncInstrument
}

func (i Int64Observer) Observation(v int64) Observation {
	return Observation{
		number:     core.NewInt64Number(v),
		instrument: i.instrument,
	}
}

// Float64Observer is a metric that captures a set of float64 values
// at a point in time.
type Float64Observer struct {
	asyncInstrument
}

func (f Float64Observer) Observation(v float64) Observation {
	return Observation{
		number:     core.NewFloat64Number(v),
		instrument: f.instrument,
	}
}
