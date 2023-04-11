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

// Package metricembed provides interfaces embedded within the [OpenTelemetry
// metric API].
//
// Implementers of the [OpenTelemetry metric API] can embed the relevant type
// from this package into their implementation directly. Doing so will result
// in a compilation error for users when the [OpenTelemetry metric API] is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [OpenTelemetry metric API]: go.opentelemetry.io/otel/metric
package metricembed // import "go.opentelemetry.io/otel/metric/metricembed"

// MeterProvider is embedded in the OpenTelemetry metric API [MeterProvider].
//
// Embed this interface in your implementation of the [MeterProvider] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [MeterProvider] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [MeterProvider]: go.opentelemetry.io/otel/metric.MeterProvider
type MeterProvider interface{ meterProvider() }

// Meter is embedded in the OpenTelemetry metric API [Meter].
//
// Embed this interface in your implementation of the [Meter] if you want users
// to experience a compilation error, signaling they need to update to your
// latest implementation, when the [Meter] interface is extended (which is
// something that can happen without a major version bump of the API package).
//
// [Meter]: go.opentelemetry.io/otel/metric.Meter
type Meter interface{ meter() }

// Float64Observer is embedded in the OpenTelemetry metric API
// [Float64Observer].
//
// Embed this interface in your implementation of the [Float64Observer] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Float64Observer] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Float64Observer]: go.opentelemetry.io/otel/metric.Float64Observer
type Float64Observer interface{ float64Observer() }

// Int64Observer is embedded in the OpenTelemetry metric API [Int64Observer].
//
// Embed this interface in your implementation of the [Int64Observer] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Int64Observer] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Int64Observer]: go.opentelemetry.io/otel/metric.Int64Observer
type Int64Observer interface{ int64Observer() }

// Observer is embedded in the OpenTelemetry metric API [Observer].
//
// Embed this interface in your implementation of the [Observer] if you want
// users to experience a compilation error, signaling they need to update to
// your latest implementation, when the [Observer] interface is extended (which
// is something that can happen without a major version bump of the API
// package).
//
// [Observer]: go.opentelemetry.io/otel/metric.Observer
type Observer interface{ observer() }

// Registration is embedded in the OpenTelemetry metric API [Registration].
//
// Embed this interface in your implementation of the [Registration] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Registration] interface is extended
// (which is something that can happen without a major version bump of the API
// package).
//
// [Registration]: go.opentelemetry.io/otel/metric.Registration
type Registration interface{ registration() }

// Float64Counter is embedded in the OpenTelemetry metric API [Float64Counter].
//
// Embed this interface in your implementation of the [Float64Counter] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Float64Counter] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Float64Counter]: go.opentelemetry.io/otel/metric.Float64Counter
type Float64Counter interface{ float64Counter() }

// Float64Histogram is embedded in the OpenTelemetry metric API
// [Float64Histogram].
//
// Embed this interface in your implementation of the [Float64Histogram] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Float64Histogram] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Float64Histogram]: go.opentelemetry.io/otel/metric.Float64Histogram
type Float64Histogram interface{ float64Histogram() }

// Float64ObservableCounter is embedded in the OpenTelemetry metric API
// [Float64ObservableCounter].
//
// Embed this interface in your implementation of the
// [Float64ObservableCounter] if you want users to experience a compilation
// error, signaling they need to update to your latest implementation, when the
// [Float64ObservableCounter] interface is extended (which is something that
// can happen without a major version bump of the API package).
//
// [Float64ObservableCounter]: go.opentelemetry.io/otel/metric.Float64ObservableCounter
type Float64ObservableCounter interface{ float64ObservableCounter() }

// Float64ObservableGauge is embedded in the OpenTelemetry metric API
// [Float64ObservableGauge].
//
// Embed this interface in your implementation of the [Float64ObservableGauge]
// if you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [Float64ObservableGauge]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [Float64ObservableGauge]: go.opentelemetry.io/otel/metric.Float64ObservableGauge
type Float64ObservableGauge interface{ float64ObservableGauge() }

// Float64ObservableUpDownCounter is embedded in the OpenTelemetry metric API
// [Float64ObservableUpDownCounter].
//
// Embed this interface in your implementation of the
// [Float64ObservableUpDownCounter] if you want users to experience a
// compilation error, signaling they need to update to your latest
// implementation, when the [Float64ObservableUpDownCounter] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Float64ObservableUpDownCounter]: go.opentelemetry.io/otel/metric.Float64ObservableUpDownCounter
type Float64ObservableUpDownCounter interface{ float64ObservableUpDownCounter() }

// Float64UpDownCounter is embedded in the OpenTelemetry metric API
// [Float64UpDownCounter].
//
// Embed this interface in your implementation of the [Float64UpDownCounter] if
// you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [Float64UpDownCounter]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [Float64UpDownCounter]: go.opentelemetry.io/otel/metric.Float64UpDownCounter
type Float64UpDownCounter interface{ float64UpDownCounter() }

// Int64Counter is embedded in the OpenTelemetry metric API [Int64Counter].
//
// Embed this interface in your implementation of the [Int64Counter] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Int64Counter] interface is extended
// (which is something that can happen without a major version bump of the API
// package).
//
// [Int64Counter]: go.opentelemetry.io/otel/metric.Int64Counter
type Int64Counter interface{ int64Counter() }

// Int64Histogram is embedded in the OpenTelemetry metric API [Int64Histogram].
//
// Embed this interface in your implementation of the [Int64Histogram] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Int64Histogram] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Int64Histogram]: go.opentelemetry.io/otel/metric.Int64Histogram
type Int64Histogram interface{ int64Histogram() }

// Int64ObservableCounter is embedded in the OpenTelemetry metric API
// [Int64ObservableCounter].
//
// Embed this interface in your implementation of the [Int64ObservableCounter]
// if you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [Int64ObservableCounter]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [Int64ObservableCounter]: go.opentelemetry.io/otel/metric.Int64ObservableCounter
type Int64ObservableCounter interface{ int64ObservableCounter() }

// Int64ObservableGauge is embedded in the OpenTelemetry metric API
// [Int64ObservableGauge].
//
// Embed this interface in your implementation of the [Int64ObservableGauge] if
// you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [Int64ObservableGauge]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [Int64ObservableGauge]: go.opentelemetry.io/otel/metric.Int64ObservableGauge
type Int64ObservableGauge interface{ int64ObservableGauge() }

// Int64ObservableUpDownCounter is embedded in the OpenTelemetry metric API
// [Int64ObservableUpDownCounter].
//
// Embed this interface in your implementation of the
// [Int64ObservableUpDownCounter] if you want users to experience a compilation
// error, signaling they need to update to your latest implementation, when the
// [Int64ObservableUpDownCounter] interface is extended (which is something
// that can happen without a major version bump of the API package).
//
// [Int64ObservableUpDownCounter]: go.opentelemetry.io/otel/metric.Int64ObservableUpDownCounter
type Int64ObservableUpDownCounter interface{ int64ObservableUpDownCounter() }

// Int64UpDownCounter is embedded in the OpenTelemetry metric API
// [Int64UpDownCounter].
//
// Embed this interface in your implementation of the [Int64UpDownCounter] if
// you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [Int64UpDownCounter]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [Int64UpDownCounter]: go.opentelemetry.io/otel/metric.Int64UpDownCounter
type Int64UpDownCounter interface{ int64UpDownCounter() }
