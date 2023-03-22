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

// Package embedded provides interfaces embedded within the OpenTelemetry
// metric API.
package embedded // import "go.opentelemetry.io/otel/metric/embedded"

// MeterProvider is embedded in the OpenTelemetry metric API MeterProvider.
type MeterProvider interface{ meterProvider() }

// Meter is embedded in the OpenTelemetry metric API Meter.
type Meter interface{ meter() }

// Float64Observer is embedded in the OpenTelemetry metric API Float64Observer.
type Float64Observer interface{ float64Observer() }

// Int64Observer is embedded in the OpenTelemetry metric API Int64Observer.
type Int64Observer interface{ int64Observer() }

// Observer is embedded in the OpenTelemetry metric API Observer.
type Observer interface{ observer() }

// Registration is embedded in the OpenTelemetry metric API Registration.
type Registration interface{ registration() }

// Float64Counter is embedded in the OpenTelemetry metric API Float64Counter.
type Float64Counter interface{ float64Counter() }

// Float64Histogram is embedded in the OpenTelemetry metric API
// Float64Histogram.
type Float64Histogram interface{ float64Histogram() }

// Float64ObservableCounter is embedded in the OpenTelemetry metric API
// Float64ObservableCounter.
type Float64ObservableCounter interface{ float64ObservableCounter() }

// Float64ObservableGauge is embedded in the OpenTelemetry metric API
// Float64ObservableGauge.
type Float64ObservableGauge interface{ float64ObservableGauge() }

// Float64ObservableUpDownCounter is embedded in the OpenTelemetry metric API
// Float64ObservableUpDownCounter.
type Float64ObservableUpDownCounter interface{ float64ObservableUpDownCounter() }

// Float64UpDownCounter is embedded in the OpenTelemetry metric API
// Float64UpDownCounter.
type Float64UpDownCounter interface{ float64UpDownCounter() }

// Int64Counter is embedded in the OpenTelemetry metric API Int64Counter.
type Int64Counter interface{ int64Counter() }

// Int64Histogram is embedded in the OpenTelemetry metric API Int64Histogram.
type Int64Histogram interface{ int64Histogram() }

// Int64ObservableCounter is embedded in the OpenTelemetry metric API
// Int64ObservableCounter.
type Int64ObservableCounter interface{ int64ObservableCounter() }

// Int64ObservableGauge is embedded in the OpenTelemetry metric API
// Int64ObservableGauge.
type Int64ObservableGauge interface{ int64ObservableGauge() }

// Int64ObservableUpDownCounter is embedded in the OpenTelemetry metric API
// Int64ObservableUpDownCounter.
type Int64ObservableUpDownCounter interface{ int64ObservableUpDownCounter() }

// Int64UpDownCounter is embedded in the OpenTelemetry metric API
// Int64UpDownCounter.
type Int64UpDownCounter interface{ int64UpDownCounter() }
