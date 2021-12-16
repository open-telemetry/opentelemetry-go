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
	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"
)

// MeterProvider supports creating named Meter instances, for instrumenting
// an application containing multiple libraries of code.
type MeterProvider interface {
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter is an instance of an OpenTelemetry metrics interface for an
// individual named library of code.  This is the top-level entry
// point for creating instruments.
type Meter interface {
	AsyncInt64() AsyncInt64Instruments
	AsyncFloat64() AsyncFloat64Instruments
	SyncInt64() SyncInt64Instruments
	SyncFloat64() SyncFloat64Instruments
}

type AsyncFloat64Instruments interface {
	Counter(name string, opts ...InstrumentOption) (asyncfloat64.Counter, error)
	UpDownCounter(name string, opts ...InstrumentOption) (asyncfloat64.UpDownCounter, error)
	Gauge(name string, opts ...InstrumentOption) (asyncfloat64.Gauge, error)
}

type AsyncInt64Instruments interface {
	Counter(name string, opts ...InstrumentOption) (asyncint64.Counter, error)
	UpDownCounter(name string, opts ...InstrumentOption) (asyncint64.UpDownCounter, error)
	Gauge(name string, opts ...InstrumentOption) (asyncint64.Gauge, error)
}

type SyncFloat64Instruments interface {
	Counter(name string, opts ...InstrumentOption) (syncfloat64.Counter, error)
	UpDownCounter(name string, opts ...InstrumentOption) (syncfloat64.UpDownCounter, error)
	Histogram(name string, opts ...InstrumentOption) (syncfloat64.Histogram, error)
}

type SyncInt64Instruments interface {
	Counter(name string, opts ...InstrumentOption) (syncint64.Counter, error)
	UpDownCounter(name string, opts ...InstrumentOption) (syncint64.UpDownCounter, error)
	Histogram(name string, opts ...InstrumentOption) (syncint64.Histogram, error)
}
