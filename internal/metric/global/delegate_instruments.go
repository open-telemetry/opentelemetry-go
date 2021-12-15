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

package global // import "go.opentelemetry.io/otel/internal/metric/global"

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

var noopFloat64Counter, _ = noopMeter.NewFloat64Counter("")
var noopInt64Counter, _ = noopMeter.NewInt64Counter("")
var noopFloat64UpDownCounter, _ = noopMeter.NewFloat64UpDownCounter("")
var noopInt64UpDownCounter, _ = noopMeter.NewInt64UpDownCounter("")
var noopFloat64Histogram, _ = noopMeter.NewFloat64Histogram("")
var noopInt64Histogram, _ = noopMeter.NewInt64Histogram("")

var noopFloat64GaugeObserver, _ = noopMeter.NewFloat64GaugeObserver("", nil)
var noopInt64GaugeObserver, _ = noopMeter.NewInt64GaugeObserver("", nil)
var noopFloat64CounterObserver, _ = noopMeter.NewFloat64CounterObserver("", nil)
var noopInt64CounterObserver, _ = noopMeter.NewInt64CounterObserver("", nil)
var noopFloat64UpDownCounterObserver, _ = noopMeter.NewFloat64UpDownCounterObserver("", nil)
var noopInt64UpDownCounterObserver, _ = noopMeter.NewInt64UpDownCounterObserver("", nil)

type float64CounterImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64Counter
}

func newFloat64CounterDelegate(name string, options ...metric.InstrumentOption) metric.Float64Counter {
	ret := &float64CounterImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopFloat64Counter)
	return ret
}

func (c *float64CounterImpl) Measurement(value float64) metric.Measurement {
	return c.delegate.Load().(metric.Float64Counter).Measurement(value)
}

func (c *float64CounterImpl) Add(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Float64Counter).Add(ctx, value, labels...)
}

func (c *float64CounterImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Float64Counter).SyncImpl()
}

func (c *float64CounterImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64Counter(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64CounterImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64Counter
}

func newInt64CounterDelegate(name string, options ...metric.InstrumentOption) metric.Int64Counter {
	ret := &int64CounterImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopInt64Counter)
	return ret
}

func (c *int64CounterImpl) Measurement(value int64) metric.Measurement {
	return c.delegate.Load().(metric.Int64Counter).Measurement(value)
}

func (c *int64CounterImpl) Add(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Int64Counter).Add(ctx, value, labels...)
}

func (c *int64CounterImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Int64Counter).SyncImpl()
}

func (c *int64CounterImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64Counter(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type float64UpDownCounterImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64UpDownCounter
}

func newFloat64UpDownCounterDelegate(name string, options ...metric.InstrumentOption) metric.Float64UpDownCounter {
	ret := &float64UpDownCounterImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopFloat64UpDownCounter)
	return ret
}

func (c *float64UpDownCounterImpl) Measurement(value float64) metric.Measurement {
	return c.delegate.Load().(metric.Float64UpDownCounter).Measurement(value)
}

func (c *float64UpDownCounterImpl) Add(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Float64UpDownCounter).Add(ctx, value, labels...)
}

func (c *float64UpDownCounterImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Float64UpDownCounter).SyncImpl()
}

func (c *float64UpDownCounterImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64UpDownCounter(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64UpDownCounterImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64UpDownCounter
}

func newInt64UpDownCounterDelegate(name string, options ...metric.InstrumentOption) metric.Int64UpDownCounter {
	ret := &int64UpDownCounterImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopInt64UpDownCounter)
	return ret
}

func (c *int64UpDownCounterImpl) Measurement(value int64) metric.Measurement {
	return c.delegate.Load().(metric.Int64UpDownCounter).Measurement(value)
}

func (c *int64UpDownCounterImpl) Add(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Int64UpDownCounter).Add(ctx, value, labels...)
}

func (c *int64UpDownCounterImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Int64UpDownCounter).SyncImpl()
}

func (c *int64UpDownCounterImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64UpDownCounter(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type float64HistogramImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64Histogram
}

func newFloat64HistogramDelegate(name string, options ...metric.InstrumentOption) metric.Float64Histogram {
	ret := &float64HistogramImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopFloat64Histogram)
	return ret
}

func (c *float64HistogramImpl) Measurement(value float64) metric.Measurement {
	return c.delegate.Load().(metric.Float64Histogram).Measurement(value)
}

func (c *float64HistogramImpl) Record(ctx context.Context, value float64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Float64Histogram).Record(ctx, value, labels...)
}

func (c *float64HistogramImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Float64Histogram).SyncImpl()
}

func (c *float64HistogramImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64Histogram(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64HistogramImpl struct {
	name     string
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64Histogram
}

func newInt64HistogramDelegate(name string, options ...metric.InstrumentOption) metric.Int64Histogram {
	ret := &int64HistogramImpl{
		name:    name,
		options: options,
	}
	ret.delegate.Store(noopInt64Histogram)
	return ret
}

func (c *int64HistogramImpl) Measurement(value int64) metric.Measurement {
	return c.delegate.Load().(metric.Int64Histogram).Measurement(value)
}

func (c *int64HistogramImpl) Record(ctx context.Context, value int64, labels ...attribute.KeyValue) {
	c.delegate.Load().(metric.Int64Histogram).Record(ctx, value, labels...)
}

func (c *int64HistogramImpl) SyncImpl() sdkapi.SyncImpl {
	return c.delegate.Load().(metric.Int64Histogram).SyncImpl()
}

func (c *int64HistogramImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64Histogram(c.name, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type float64GaugeObserverImpl struct {
	name     string
	callback metric.Float64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64GaugeObserver
}

func newFloat64GaugeObserverDelegate(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) metric.Float64GaugeObserver {
	ret := &float64GaugeObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopFloat64GaugeObserver)
	return ret
}

func (c *float64GaugeObserverImpl) Observation(v float64) metric.Observation {
	return c.delegate.Load().(metric.Float64GaugeObserver).Observation(v)
}

func (c *float64GaugeObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Float64GaugeObserver).AsyncImpl()
}

func (c *float64GaugeObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64GaugeObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64GaugeObserverImpl struct {
	name     string
	callback metric.Int64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64GaugeObserver
}

func newInt64GaugeObserverDelegate(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) metric.Int64GaugeObserver {
	ret := &int64GaugeObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopInt64GaugeObserver)
	return ret
}

func (c *int64GaugeObserverImpl) Observation(v int64) metric.Observation {
	return c.delegate.Load().(metric.Int64GaugeObserver).Observation(v)
}

func (c *int64GaugeObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Int64GaugeObserver).AsyncImpl()
}

func (c *int64GaugeObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64GaugeObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type float64CounterObserverImpl struct {
	name     string
	callback metric.Float64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64CounterObserver
}

func newFloat64CounterObserverDelegate(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) metric.Float64CounterObserver {
	ret := &float64CounterObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopFloat64CounterObserver)
	return ret
}

func (c *float64CounterObserverImpl) Observation(v float64) metric.Observation {
	return c.delegate.Load().(metric.Float64CounterObserver).Observation(v)
}

func (c *float64CounterObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Float64CounterObserver).AsyncImpl()
}

func (c *float64CounterObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64CounterObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64CounterObserverImpl struct {
	name     string
	callback metric.Int64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64CounterObserver
}

func newInt64CounterObserverDelegate(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) metric.Int64CounterObserver {
	ret := &int64CounterObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopInt64CounterObserver)
	return ret
}

func (c *int64CounterObserverImpl) Observation(v int64) metric.Observation {
	return c.delegate.Load().(metric.Int64CounterObserver).Observation(v)
}

func (c *int64CounterObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Int64CounterObserver).AsyncImpl()
}

func (c *int64CounterObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64CounterObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type float64UpDownCounterObserverImpl struct {
	name     string
	callback metric.Float64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Float64UpDownCounterObserver
}

func newFloat64UpDownCounterObserverDelegate(name string, callback metric.Float64ObserverFunc, options ...metric.InstrumentOption) metric.Float64UpDownCounterObserver {
	ret := &float64UpDownCounterObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopFloat64UpDownCounterObserver)
	return ret
}

func (c *float64UpDownCounterObserverImpl) Observation(v float64) metric.Observation {
	return c.delegate.Load().(metric.Float64UpDownCounterObserver).Observation(v)
}

func (c *float64UpDownCounterObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Float64UpDownCounterObserver).AsyncImpl()
}

func (c *float64UpDownCounterObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewFloat64UpDownCounterObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}

type int64UpDownCounterObserverImpl struct {
	name     string
	callback metric.Int64ObserverFunc
	options  []metric.InstrumentOption
	delegate atomic.Value // metric.Int64UpDownCounterObserver
}

func newInt64UpDownCounterObserverDelegate(name string, callback metric.Int64ObserverFunc, options ...metric.InstrumentOption) metric.Int64UpDownCounterObserver {
	ret := &int64UpDownCounterObserverImpl{
		name:     name,
		callback: callback,
		options:  options,
	}
	ret.delegate.Store(noopInt64UpDownCounterObserver)
	return ret
}

func (c *int64UpDownCounterObserverImpl) Observation(v int64) metric.Observation {
	return c.delegate.Load().(metric.Int64UpDownCounterObserver).Observation(v)
}

func (c *int64UpDownCounterObserverImpl) AsyncImpl() sdkapi.AsyncImpl {
	return c.delegate.Load().(metric.Int64UpDownCounterObserver).AsyncImpl()
}

func (c *int64UpDownCounterObserverImpl) setDelegate(meter metric.Meter) error {
	impl, err := meter.NewInt64UpDownCounterObserver(c.name, c.callback, c.options...)
	if err != nil {
		return err
	}
	c.delegate.Store(impl)
	return nil
}
