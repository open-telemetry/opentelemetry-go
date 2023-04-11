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
// This implementation can be metricembed in other implementations of the
// OpenTelemetry metric API. Doing so will mean the implementation defaults to
// no operation for methods it does not implement.
package noop // import "go.opentelemetry.io/otel/metric/noop"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/metricembed"
)

var (
	// Compile-time check this implements the OpenTelemetry API.

	_ metric.MeterProvider                      = MeterProvider{}
	_ metric.Meter                              = Meter{}
	_ metric.Observer                           = Observer{}
	_ metric.Registration                       = Registration{}
	_ instrument.Int64Counter                   = Int64Counter{}
	_ instrument.Float64Counter                 = Float64Counter{}
	_ instrument.Int64UpDownCounter             = Int64UpDownCounter{}
	_ instrument.Float64UpDownCounter           = Float64UpDownCounter{}
	_ instrument.Int64Histogram                 = Int64Histogram{}
	_ instrument.Float64Histogram               = Float64Histogram{}
	_ instrument.Int64ObservableCounter         = Int64ObservableCounter{}
	_ instrument.Float64ObservableCounter       = Float64ObservableCounter{}
	_ instrument.Int64ObservableGauge           = Int64ObservableGauge{}
	_ instrument.Float64ObservableGauge         = Float64ObservableGauge{}
	_ instrument.Int64ObservableUpDownCounter   = Int64ObservableUpDownCounter{}
	_ instrument.Float64ObservableUpDownCounter = Float64ObservableUpDownCounter{}
	_ instrument.Int64Observer                  = Int64Observer{}
	_ instrument.Float64Observer                = Float64Observer{}
)

// MeterProvider is an OpenTelemetry No-Op MeterProvider.
type MeterProvider struct{ metricembed.MeterProvider }

// NewMeterProvider returns a MeterProvider that does not record any telemetry.
func NewMeterProvider() MeterProvider {
	return MeterProvider{}
}

// Meter returns an OpenTelemetry Meter that does not record any telemetry.
func (MeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return Meter{}
}

// Meter is an OpenTelemetry No-Op Meter.
type Meter struct{ metricembed.Meter }

// Int64Counter returns a Counter used to record int64 measurements that
// produces no telemetry.
func (Meter) Int64Counter(string, ...instrument.Int64CounterOption) (instrument.Int64Counter, error) {
	return Int64Counter{}, nil
}

// Int64UpDownCounter returns an UpDownCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Int64UpDownCounter(string, ...instrument.Int64UpDownCounterOption) (instrument.Int64UpDownCounter, error) {
	return Int64UpDownCounter{}, nil
}

// Int64Histogram returns a Histogram used to record int64 measurements that
// produces no telemetry.
func (Meter) Int64Histogram(string, ...instrument.Int64HistogramOption) (instrument.Int64Histogram, error) {
	return Int64Histogram{}, nil
}

// Int64ObservableCounter returns an ObservableCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Int64ObservableCounter(string, ...instrument.Int64ObservableCounterOption) (instrument.Int64ObservableCounter, error) {
	return Int64ObservableCounter{}, nil
}

// Int64ObservableUpDownCounter returns an ObservableUpDownCounter used to
// record int64 measurements that produces no telemetry.
func (Meter) Int64ObservableUpDownCounter(string, ...instrument.Int64ObservableUpDownCounterOption) (instrument.Int64ObservableUpDownCounter, error) {
	return Int64ObservableUpDownCounter{}, nil
}

// Int64ObservableGauge returns an ObservableGauge used to record int64
// measurements that produces no telemetry.
func (Meter) Int64ObservableGauge(string, ...instrument.Int64ObservableGaugeOption) (instrument.Int64ObservableGauge, error) {
	return Int64ObservableGauge{}, nil
}

// Float64Counter returns a Counter used to record int64 measurements that
// produces no telemetry.
func (Meter) Float64Counter(string, ...instrument.Float64CounterOption) (instrument.Float64Counter, error) {
	return Float64Counter{}, nil
}

// Float64UpDownCounter returns an UpDownCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Float64UpDownCounter(string, ...instrument.Float64UpDownCounterOption) (instrument.Float64UpDownCounter, error) {
	return Float64UpDownCounter{}, nil
}

// Float64Histogram returns a Histogram used to record int64 measurements that
// produces no telemetry.
func (Meter) Float64Histogram(string, ...instrument.Float64HistogramOption) (instrument.Float64Histogram, error) {
	return Float64Histogram{}, nil
}

// Float64ObservableCounter returns an ObservableCounter used to record int64
// measurements that produces no telemetry.
func (Meter) Float64ObservableCounter(string, ...instrument.Float64ObservableCounterOption) (instrument.Float64ObservableCounter, error) {
	return Float64ObservableCounter{}, nil
}

// Float64ObservableUpDownCounter returns an ObservableUpDownCounter used to
// record int64 measurements that produces no telemetry.
func (Meter) Float64ObservableUpDownCounter(string, ...instrument.Float64ObservableUpDownCounterOption) (instrument.Float64ObservableUpDownCounter, error) {
	return Float64ObservableUpDownCounter{}, nil
}

// Float64ObservableGauge returns an ObservableGauge used to record int64
// measurements that produces no telemetry.
func (Meter) Float64ObservableGauge(string, ...instrument.Float64ObservableGaugeOption) (instrument.Float64ObservableGauge, error) {
	return Float64ObservableGauge{}, nil
}

// RegisterCallback performs no operation.
func (Meter) RegisterCallback(metric.Callback, ...instrument.Observable) (metric.Registration, error) {
	return Registration{}, nil
}

// Observer acts as a recorder of measurements for multiple instruments in a
// Callback, it performing no operation.
type Observer struct{ metricembed.Observer }

// ObserveFloat64 performs no operation.
func (Observer) ObserveFloat64(instrument.Float64Observable, float64, ...attribute.KeyValue) {
}

// ObserveInt64 performs no operation.
func (Observer) ObserveInt64(instrument.Int64Observable, int64, ...attribute.KeyValue) {
}

// Registration is the registration of a Callback with a No-Op Meter.
type Registration struct{ metricembed.Registration }

// Unregister unregisters the Callback the Registration represents with the
// No-Op Meter. This will always return nil because the No-Op Meter performs no
// operation, including hold any record of registrations.
func (Registration) Unregister() error { return nil }

// Int64Counter is an OpenTelemetry Counter used to record int64 measurements.
// It produces no telemetry.
type Int64Counter struct{ metricembed.Int64Counter }

// Add performs no operation.
func (Int64Counter) Add(context.Context, int64, ...attribute.KeyValue) {}

// Float64Counter is an OpenTelemetry Counter used to record float64
// measurements. It produces no telemetry.
type Float64Counter struct{ metricembed.Float64Counter }

// Add performs no operation.
func (Float64Counter) Add(context.Context, float64, ...attribute.KeyValue) {}

// Int64UpDownCounter is an OpenTelemetry UpDownCounter used to record int64
// measurements. It produces no telemetry.
type Int64UpDownCounter struct{ metricembed.Int64UpDownCounter }

// Add performs no operation.
func (Int64UpDownCounter) Add(context.Context, int64, ...attribute.KeyValue) {}

// Float64UpDownCounter is an OpenTelemetry UpDownCounter used to record
// float64 measurements. It produces no telemetry.
type Float64UpDownCounter struct {
	metricembed.Float64UpDownCounter
}

// Add performs no operation.
func (Float64UpDownCounter) Add(context.Context, float64, ...attribute.KeyValue) {}

// Int64Histogram is an OpenTelemetry Histogram used to record int64
// measurements. It produces no telemetry.
type Int64Histogram struct{ metricembed.Int64Histogram }

// Record performs no operation.
func (Int64Histogram) Record(context.Context, int64, ...attribute.KeyValue) {}

// Float64Histogram is an OpenTelemetry Histogram used to record float64
// measurements. It produces no telemetry.
type Float64Histogram struct{ metricembed.Float64Histogram }

// Record performs no operation.
func (Float64Histogram) Record(context.Context, float64, ...attribute.KeyValue) {}

// Int64ObservableCounter is an OpenTelemetry ObservableCounter used to record
// int64 measurements. It produces no telemetry.
type Int64ObservableCounter struct {
	instrument.Int64Observable
	metricembed.Int64ObservableCounter
}

// Float64ObservableCounter is an OpenTelemetry ObservableCounter used to record
// float64 measurements. It produces no telemetry.
type Float64ObservableCounter struct {
	instrument.Float64Observable
	metricembed.Float64ObservableCounter
}

// Int64ObservableGauge is an OpenTelemetry ObservableGauge used to record
// int64 measurements. It produces no telemetry.
type Int64ObservableGauge struct {
	instrument.Int64Observable
	metricembed.Int64ObservableGauge
}

// Float64ObservableGauge is an OpenTelemetry ObservableGauge used to record
// float64 measurements. It produces no telemetry.
type Float64ObservableGauge struct {
	instrument.Float64Observable
	metricembed.Float64ObservableGauge
}

// Int64ObservableUpDownCounter is an OpenTelemetry ObservableUpDownCounter
// used to record int64 measurements. It produces no telemetry.
type Int64ObservableUpDownCounter struct {
	instrument.Int64Observable
	metricembed.Int64ObservableUpDownCounter
}

// Float64ObservableUpDownCounter is an OpenTelemetry ObservableUpDownCounter
// used to record float64 measurements. It produces no telemetry.
type Float64ObservableUpDownCounter struct {
	instrument.Float64Observable
	metricembed.Float64ObservableUpDownCounter
}

// Int64Observer is a recorder of int64 measurements that performs no operation.
type Int64Observer struct{ metricembed.Int64Observer }

// Observe performs no operation.
func (Int64Observer) Observe(int64, ...attribute.KeyValue) {}

// Float64Observer is a recorder of float64 measurements that performs no
// operation.
type Float64Observer struct{ metricembed.Float64Observer }

// Observe performs no operation.
func (Float64Observer) Observe(float64, ...attribute.KeyValue) {}
