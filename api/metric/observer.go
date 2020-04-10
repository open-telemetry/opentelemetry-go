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

import "go.opentelemetry.io/otel/api/core"

type Int64ObserverResult interface {
	Observe(value int64, labels ...core.KeyValue)
}
type Float64ObserverResult interface {
	Observe(value float64, labels ...core.KeyValue)
}
type BatchObserverResult interface {
	Observe(labels []core.KeyValue, observations ...Observation)
}

type int64ObserverResult struct {
	instrument AsyncImpl
	observe    func(AsyncImpl, core.Number, []core.KeyValue)
}

type float64ObserverResult struct {
	instrument AsyncImpl
	observe    func(AsyncImpl, core.Number, []core.KeyValue)
}

type batchObserverResult struct {
	// observe func(AsyncImpl, core.Number, []core.KeyValue)
}

var _ Int64ObserverResult = int64ObserverResult{}
var _ Float64ObserverResult = float64ObserverResult{}
var _ BatchObserverResult = batchObserverResult{}

func (ir int64ObserverResult) Observe(value int64, labels ...core.KeyValue) {
	ir.observe(ir.instrument, core.NewInt64Number(value), labels)
}

func (fr float64ObserverResult) Observe(value float64, labels ...core.KeyValue) {
	fr.observe(fr.instrument, core.NewFloat64Number(value), labels)
}

func (br batchObserverResult) Observe(labels []core.KeyValue, obs ...Observation) {
	// fr.observe(fr.instrument, core.NewFloat64Number(value), labels)
}

type AsyncRunner interface {
	anyRunner()
}

type AsyncSingleRunner interface {
	Run(func(AsyncImpl, core.Number, []core.KeyValue))
	anyRunner()
}

type AsyncBatchRunner interface {
	Run(func([]core.KeyValue, []Observation))
	anyRunner()
}

type Int64ObserverCallback interface {
	intObserver()
}
type Float64ObserverCallback interface {
	floatObserver()
}
type BatchObserverCallback interface {
	Run(func([]core.KeyValue, []Observation))
	intObserver()
	floatObserver()
}

type int64ObserverCallback func(result Int64ObserverResult)
type float64ObserverCallback func(result Float64ObserverResult)
type batchObserverCallback func(result BatchObserverResult)

var _ Int64ObserverCallback = (*int64ObserverCallback)(nil)
var _ Int64ObserverCallback = (*batchObserverCallback)(nil)
var _ Float64ObserverCallback = (*float64ObserverCallback)(nil)
var _ Float64ObserverCallback = (*batchObserverCallback)(nil)

func (*int64ObserverCallback) intObserver()     {}
func (*float64ObserverCallback) floatObserver() {}
func (*batchObserverCallback) intObserver()     {}
func (*batchObserverCallback) floatObserver()   {}

func (i *int64ObserverCallback) Run(function func(core.Number, []core.KeyValue)) {
	//(*i)(function)
}

func (f *float64ObserverCallback) Run(function func(core.Number, []core.KeyValue)) {
	// (*f)(function)
}

func (f *batchObserverCallback) Run(function func([]core.KeyValue, []Observation)) {
	// (*f)(function)
}

func NewInt64ObserverCallback(f func(Int64ObserverResult)) Int64ObserverCallback {
	var c int64ObserverCallback = f
	return &c
}

func NewFloat64ObserverCallback(f func(Float64ObserverResult)) Float64ObserverCallback {
	var c float64ObserverCallback = f
	return &c
}

func NewBatchObserverCallback(f func(BatchObserverResult)) BatchObserverCallback {
	var c batchObserverCallback = f
	return &c
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
