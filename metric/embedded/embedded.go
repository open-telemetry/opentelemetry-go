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

type MeterProvider interface{ meterProvider() }

type Meter interface{ meter() }

type Float64Observer interface{ float64Observer() }
type Int64Observer interface{ int64Observer() }
type Observer interface{ observer() }

type Registration interface{ registration() }

type Float64Counter interface{ float64Counter() }
type Float64Histogram interface{ float64Histogram() }
type Float64ObservableCounter interface{ float64ObservableCounter() }
type Float64ObservableGauge interface{ float64ObservableGauge() }
type Float64ObservableUpDownCounter interface{ float64ObservableUpDownCounter() }
type Float64UpDownCounter interface{ float64UpDownCounter() }
type Int64Counter interface{ int64Counter() }
type Int64Histogram interface{ int64Histogram() }
type Int64ObservableCounter interface{ int64ObservableCounter() }
type Int64ObservableGauge interface{ int64ObservableGauge() }
type Int64ObservableUpDownCounter interface{ int64ObservableUpDownCounter() }
type Int64UpDownCounter interface{ int64UpDownCounter() }
