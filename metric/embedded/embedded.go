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

// Package embedded provides interfaces embedded within the [OpenTelemetry
// metric API].
//
// Implementers of the [OpenTelemetry metric API] can embed the relevant type
// from this package into their implementation directly. Doing so will result
// in a compilation error for users when the [OpenTelemetry metric API] is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [OpenTelemetry metric API]: go.opentelemetry.io/otel/metric
package embedded // import "go.opentelemetry.io/otel/metric/embedded"

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

// ObserverT is embedded in the OpenTelemetry metric API [ObserverT].
//
// Embed this interface in your implementation of the [Observer] if you want
// users to experience a compilation error, signaling they need to update to
// your latest implementation, when the [Observer] interface is extended (which
// is something that can happen without a major version bump of the API
// package).
//
// [Observer]: go.opentelemetry.io/otel/metric.Observer
type ObserverT[N int64 | float64] interface{ observerT() }

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

// Counter is embedded in the OpenTelemetry metric API [Counter].
//
// Embed this interface in your implementation of the [Counter] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Counter] interface is extended
// (which is something that can happen without a major version bump of the API
// package).
//
// [Counter]: go.opentelemetry.io/otel/metric.Counter
type Counter[N int64 | float64] interface{ counter() }

// Histogram is embedded in the OpenTelemetry metric API [Histogram].
//
// Embed this interface in your implementation of the [Histogram] if you
// want users to experience a compilation error, signaling they need to update
// to your latest implementation, when the [Histogram] interface is
// extended (which is something that can happen without a major version bump of
// the API package).
//
// [Histogram]: go.opentelemetry.io/otel/metric.Histogram
type Histogram[N int64 | float64] interface{ histogram() }

// ObservableCounter is embedded in the OpenTelemetry metric API
// [ObservableCounter].
//
// Embed this interface in your implementation of the [ObservableCounter]
// if you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [ObservableCounter]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [ObservableCounter]: go.opentelemetry.io/otel/metric.ObservableCounter
type ObservableCounter[N int64 | float64] interface{ observableCounter() }

// ObservableGauge is embedded in the OpenTelemetry metric API
// [ObservableGauge].
//
// Embed this interface in your implementation of the [ObservableGauge] if
// you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [ObservableGauge]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [ObservableGauge]: go.opentelemetry.io/otel/metric.ObservableGauge
type ObservableGauge[N int64 | float64] interface{ observableGauge() }

// ObservableUpDownCounter is embedded in the OpenTelemetry metric API
// [ObservableUpDownCounter].
//
// Embed this interface in your implementation of the
// [ObservableUpDownCounter] if you want users to experience a compilation
// error, signaling they need to update to your latest implementation, when the
// [ObservableUpDownCounter] interface is extended (which is something
// that can happen without a major version bump of the API package).
//
// [ObservableUpDownCounter]: go.opentelemetry.io/otel/metric.ObservableUpDownCounter
type ObservableUpDownCounter[N int64 | float64] interface{ observableUpDownCounter() }

// UpDownCounter is embedded in the OpenTelemetry metric API
// [UpDownCounter].
//
// Embed this interface in your implementation of the [UpDownCounter] if
// you want users to experience a compilation error, signaling they need to
// update to your latest implementation, when the [UpDownCounter]
// interface is extended (which is something that can happen without a major
// version bump of the API package).
//
// [UpDownCounter]: go.opentelemetry.io/otel/metric.UpDownCounter
type UpDownCounter[N int64 | float64] interface{ upDownCounter() }
