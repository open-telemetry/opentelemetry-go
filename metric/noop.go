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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// NewNoopMeterProvider creates a MeterProvider that does not record any metrics.
func NewNoopMeterProvider() MeterProvider {
	return noopMeterProvider{}
}

type noopMeterProvider struct{}

func (noopMeterProvider) Meter(string, ...MeterOption) Meter {
	return noopMeter{}
}

// NewNoopMeter creates a Meter that does not record any metrics.
func NewNoopMeter() Meter {
	return noopMeter{}
}

type noopMeter struct{}

func (noopMeter) Int64Counter(string, ...instrument.Int64Option) (syncint64.Counter, error) {
	return nonrecordingSyncInt64Instrument{}, nil
}

func (noopMeter) Int64UpDownCounter(string, ...instrument.Int64Option) (syncint64.UpDownCounter, error) {
	return nonrecordingSyncInt64Instrument{}, nil
}

func (noopMeter) Int64Histogram(string, ...instrument.Int64Option) (syncint64.Histogram, error) {
	return nonrecordingSyncInt64Instrument{}, nil
}

func (noopMeter) Int64ObservableCounter(string, ...instrument.Int64ObserverOption) (asyncint64.Counter, error) {
	return nonrecordingAsyncInt64Instrument{}, nil
}

func (noopMeter) Int64ObservableUpDownCounter(string, ...instrument.Int64ObserverOption) (asyncint64.UpDownCounter, error) {
	return nonrecordingAsyncInt64Instrument{}, nil
}

func (noopMeter) Int64ObservableGauge(string, ...instrument.Int64ObserverOption) (asyncint64.Gauge, error) {
	return nonrecordingAsyncInt64Instrument{}, nil
}

func (noopMeter) Float64Counter(string, ...instrument.Float64Option) (syncfloat64.Counter, error) {
	return nonrecordingSyncFloat64Instrument{}, nil
}

func (noopMeter) Float64UpDownCounter(string, ...instrument.Float64Option) (syncfloat64.UpDownCounter, error) {
	return nonrecordingSyncFloat64Instrument{}, nil
}

func (noopMeter) Float64Histogram(string, ...instrument.Float64Option) (syncfloat64.Histogram, error) {
	return nonrecordingSyncFloat64Instrument{}, nil
}

func (noopMeter) Float64ObservableCounter(string, ...instrument.Float64ObserverOption) (asyncfloat64.Counter, error) {
	return nonrecordingAsyncFloat64Instrument{}, nil
}

func (noopMeter) Float64ObservableUpDownCounter(string, ...instrument.Float64ObserverOption) (asyncfloat64.UpDownCounter, error) {
	return nonrecordingAsyncFloat64Instrument{}, nil
}

func (noopMeter) Float64ObservableGauge(string, ...instrument.Float64ObserverOption) (asyncfloat64.Gauge, error) {
	return nonrecordingAsyncFloat64Instrument{}, nil
}

// RegisterCallback creates a register callback that does not record any metrics.
func (noopMeter) RegisterCallback([]instrument.Asynchronous, func(context.Context)) (Registration, error) {
	return noopReg{}, nil
}

type noopReg struct{}

func (noopReg) Unregister() error { return nil }

type nonrecordingAsyncFloat64Instrument struct {
	instrument.Asynchronous
}

var (
	_ asyncfloat64.Counter       = nonrecordingAsyncFloat64Instrument{}
	_ asyncfloat64.UpDownCounter = nonrecordingAsyncFloat64Instrument{}
	_ asyncfloat64.Gauge         = nonrecordingAsyncFloat64Instrument{}
)

func (n nonrecordingAsyncFloat64Instrument) Counter(string, ...instrument.Float64ObserverOption) (asyncfloat64.Counter, error) {
	return n, nil
}

func (n nonrecordingAsyncFloat64Instrument) UpDownCounter(string, ...instrument.Float64ObserverOption) (asyncfloat64.UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingAsyncFloat64Instrument) Gauge(string, ...instrument.Float64ObserverOption) (asyncfloat64.Gauge, error) {
	return n, nil
}

func (nonrecordingAsyncFloat64Instrument) Observe(context.Context, float64, ...attribute.KeyValue) {

}

type nonrecordingAsyncInt64Instrument struct {
	instrument.Asynchronous
}

var (
	_ asyncint64.Counter       = nonrecordingAsyncInt64Instrument{}
	_ asyncint64.UpDownCounter = nonrecordingAsyncInt64Instrument{}
	_ asyncint64.Gauge         = nonrecordingAsyncInt64Instrument{}
)

func (n nonrecordingAsyncInt64Instrument) Counter(string, ...instrument.Int64ObserverOption) (asyncint64.Counter, error) {
	return n, nil
}

func (n nonrecordingAsyncInt64Instrument) UpDownCounter(string, ...instrument.Int64ObserverOption) (asyncint64.UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingAsyncInt64Instrument) Gauge(string, ...instrument.Int64ObserverOption) (asyncint64.Gauge, error) {
	return n, nil
}

func (nonrecordingAsyncInt64Instrument) Observe(context.Context, int64, ...attribute.KeyValue) {
}

type nonrecordingSyncFloat64Instrument struct {
	instrument.Synchronous
}

var (
	_ syncfloat64.Counter       = nonrecordingSyncFloat64Instrument{}
	_ syncfloat64.UpDownCounter = nonrecordingSyncFloat64Instrument{}
	_ syncfloat64.Histogram     = nonrecordingSyncFloat64Instrument{}
)

func (n nonrecordingSyncFloat64Instrument) Counter(string, ...instrument.Float64Option) (syncfloat64.Counter, error) {
	return n, nil
}

func (n nonrecordingSyncFloat64Instrument) UpDownCounter(string, ...instrument.Float64Option) (syncfloat64.UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingSyncFloat64Instrument) Histogram(string, ...instrument.Float64Option) (syncfloat64.Histogram, error) {
	return n, nil
}

func (nonrecordingSyncFloat64Instrument) Add(context.Context, float64, ...attribute.KeyValue) {

}

func (nonrecordingSyncFloat64Instrument) Record(context.Context, float64, ...attribute.KeyValue) {

}

type nonrecordingSyncInt64Instrument struct {
	instrument.Synchronous
}

var (
	_ syncint64.Counter       = nonrecordingSyncInt64Instrument{}
	_ syncint64.UpDownCounter = nonrecordingSyncInt64Instrument{}
	_ syncint64.Histogram     = nonrecordingSyncInt64Instrument{}
)

func (n nonrecordingSyncInt64Instrument) Counter(string, ...instrument.Int64Option) (syncint64.Counter, error) {
	return n, nil
}

func (n nonrecordingSyncInt64Instrument) UpDownCounter(string, ...instrument.Int64Option) (syncint64.UpDownCounter, error) {
	return n, nil
}

func (n nonrecordingSyncInt64Instrument) Histogram(string, ...instrument.Int64Option) (syncint64.Histogram, error) {
	return n, nil
}

func (nonrecordingSyncInt64Instrument) Add(context.Context, int64, ...attribute.KeyValue) {
}
func (nonrecordingSyncInt64Instrument) Record(context.Context, int64, ...attribute.KeyValue) {
}
