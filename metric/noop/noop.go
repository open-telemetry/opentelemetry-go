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

// Package noop provides an implementation of the OpenTelemetry metric API that
// produces no telemetry and minimizes used computation resources.
//
// Using this package to implement the OpenTelemetry metric API will
// effectively disable OpenTelemetry.
//
// This implementation can be embedded in other implementations of the
// OpenTelemetry metric API. Doing so will mean the implementation defaults to
// no operation for methods it does not implement.
package noop // import "go.opentelemetry.io/otel/metric/noop"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
)

var (
	// Compile-time check this implements the OpenTelemetry API.

	_ metric.MeterProvider                        = MeterProvider{}
	_ metric.Meter                                = Meter{}
	_ metric.Observer                             = Observer{}
	_ metric.Registration                         = Registration{}
	_ instrument.Counter[int64]                   = Counter[int64]{}
	_ instrument.UpDownCounter[int64]             = UpDownCounter[int64]{}
	_ instrument.Histogram[int64]                 = Histogram[int64]{}
	_ instrument.ObservableCounter[int64]         = ObservableCounter[int64]{}
	_ instrument.ObservableGauge[int64]           = ObservableGauge[int64]{}
	_ instrument.ObservableUpDownCounter[int64]   = ObservableUpDownCounter[int64]{}
	_ instrument.ObserverT[int64]                 = ObserverT[int64]{}
	_ instrument.Counter[float64]                 = Counter[float64]{}
	_ instrument.UpDownCounter[float64]           = UpDownCounter[float64]{}
	_ instrument.Histogram[float64]               = Histogram[float64]{}
	_ instrument.ObservableCounter[float64]       = ObservableCounter[float64]{}
	_ instrument.ObservableGauge[float64]         = ObservableGauge[float64]{}
	_ instrument.ObservableUpDownCounter[float64] = ObservableUpDownCounter[float64]{}
	_ instrument.ObserverT[float64]               = ObserverT[float64]{}
)

// MeterProvider is an OpenTelemetry No-Op MeterProvider.
type MeterProvider struct{ embedded.MeterProvider }

// NewMeterProvider returns a MeterProvider that does not record any telemetry.
func NewMeterProvider() MeterProvider {
	return MeterProvider{}
}

// Meter returns an OpenTelemetry Meter that does not record any telemetry.
func (MeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return Meter{}
}

// Meter is an OpenTelemetry No-Op Meter.
type Meter struct{ embedded.Meter }

// Int64Counter returns a Counter used to record int64 measurements that
// produces no telemetry.
func (Meter) Int64Counter(string, ...instrument.CounterOption[int64]) (instrument.Counter[int64], error) {
	return Counter[int64]{}, nil
}

// Int64UpDownCounter returns an UpDownCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Int64UpDownCounter(string, ...instrument.UpDownCounterOption[int64]) (instrument.UpDownCounter[int64], error) {
	return UpDownCounter[int64]{}, nil
}

// Int64Histogram returns a Histogram used to record int64 measurements that
// produces no telemetry.
func (Meter) Int64Histogram(string, ...instrument.HistogramOption[int64]) (instrument.Histogram[int64], error) {
	return Histogram[int64]{}, nil
}

// Int64ObservableCounter returns an ObservableCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Int64ObservableCounter(string, ...instrument.ObservableCounterOption[int64]) (instrument.ObservableCounter[int64], error) {
	return ObservableCounter[int64]{}, nil
}

// Int64ObservableUpDownCounter returns an ObservableUpDownCounter used to
// record int64 measurements that produces no telemetry.
func (Meter) Int64ObservableUpDownCounter(string, ...instrument.ObservableUpDownCounterOption[int64]) (instrument.ObservableUpDownCounter[int64], error) {
	return ObservableUpDownCounter[int64]{}, nil
}

// Int64ObservableGauge returns an ObservableGauge used to record int64
// measurements that produces no telemetry.
func (Meter) Int64ObservableGauge(string, ...instrument.ObservableGaugeOption[int64]) (instrument.ObservableGauge[int64], error) {
	return ObservableGauge[int64]{}, nil
}

// Float64Counter returns a Counter used to record int64 measurements that
// produces no telemetry.
func (Meter) Float64Counter(string, ...instrument.CounterOption[float64]) (instrument.Counter[float64], error) {
	return Counter[float64]{}, nil
}

// Float64UpDownCounter returns an UpDownCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Float64UpDownCounter(string, ...instrument.UpDownCounterOption[float64]) (instrument.UpDownCounter[float64], error) {
	return UpDownCounter[float64]{}, nil
}

// Float64Histogram returns a Histogram used to record int64 measurements that
// produces no telemetry.
func (Meter) Float64Histogram(string, ...instrument.HistogramOption[float64]) (instrument.Histogram[float64], error) {
	return Histogram[float64]{}, nil
}

// Float64ObservableCounter returns an ObservableCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Float64ObservableCounter(string, ...instrument.ObservableCounterOption[float64]) (instrument.ObservableCounter[float64], error) {
	return ObservableCounter[float64]{}, nil
}

// Float64ObservableUpDownCounter returns an ObservableUpDownCounter used to
// record int64 measurements that produces no telemetry.
func (Meter) Float64ObservableUpDownCounter(string, ...instrument.ObservableUpDownCounterOption[float64]) (instrument.ObservableUpDownCounter[float64], error) {
	return ObservableUpDownCounter[float64]{}, nil
}

// Float64ObservableGauge returns an ObservableGauge used to record int64
// measurements that produces no telemetry.
func (Meter) Float64ObservableGauge(string, ...instrument.ObservableGaugeOption[float64]) (instrument.ObservableGauge[float64], error) {
	return ObservableGauge[float64]{}, nil
}

// RegisterCallback performs no operation.
func (Meter) RegisterCallback(metric.Callback, ...instrument.Observable) (metric.Registration, error) {
	return Registration{}, nil
}

// Observer acts as a recorder of measurements for multiple instruments in a
// Callback, it performing no operation.
type Observer struct{ embedded.Observer }

// ObserveFloat64 performs no operation.
func (Observer) ObserveFloat64(instrument.ObservableT[float64], float64, ...attribute.KeyValue) {
}

// ObserveInt64 performs no operation.
func (Observer) ObserveInt64(instrument.ObservableT[int64], int64, ...attribute.KeyValue) {
}

// Registration is the registration of a Callback with a No-Op Meter.
type Registration struct{ embedded.Registration }

// Unregister unregisters the Callback the Registration represents with the
// No-Op Meter. This will always return nil because the No-Op Meter performs no
// operation, including hold any record of registrations.
func (Registration) Unregister() error { return nil }

// Counter is an OpenTelemetry Counter used to record measurements. It produces
// no telemetry.
type Counter[N int64 | float64] struct{ embedded.Counter[N] }

// Add performs no operation.
func (Counter[N]) Add(context.Context, N, ...attribute.KeyValue) {}

// UpDownCounter is an OpenTelemetry UpDownCounter used to record measurements.
// It produces no telemetry.
type UpDownCounter[N int64 | float64] struct{ embedded.UpDownCounter[N] }

// Add performs no operation.
func (UpDownCounter[N]) Add(context.Context, N, ...attribute.KeyValue) {}

// Histogram is an OpenTelemetry Histogram used to record measurements. It
// produces no telemetry.
type Histogram[N int64 | float64] struct{ embedded.Histogram[N] }

// Record performs no operation.
func (Histogram[N]) Record(context.Context, N, ...attribute.KeyValue) {}

// ObservableCounter is an OpenTelemetry ObservableCounter used to record
// measurements. It produces no telemetry.
type ObservableCounter[N int64 | float64] struct {
	instrument.ObservableT[N]
	embedded.ObservableCounter[N]
}

// ObservableGauge is an OpenTelemetry ObservableGauge used to record
// measurements. It produces no telemetry.
type ObservableGauge[N int64 | float64] struct {
	instrument.ObservableT[N]
	embedded.ObservableGauge[N]
}

// ObservableUpDownCounter is an OpenTelemetry ObservableUpDownCounter used to
// record measurements. It produces no telemetry.
type ObservableUpDownCounter[N int64 | float64] struct {
	instrument.ObservableT[N]
	embedded.ObservableUpDownCounter[N]
}

// ObserverT is a recorder of measurements that performs no operation.
type ObserverT[N int64 | float64] struct{ embedded.ObserverT[N] }

// Observe performs no operation.
func (ObserverT[N]) Observe(N, ...attribute.KeyValue) {}
