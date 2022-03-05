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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// Instruments provides a simple instrument to register
type Instruments struct {
	Meter
}

// SyncInt64Counter provides a simple SyncInt64 Counter
func (i Instruments) SyncInt64Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	return i.SyncInt64().Counter(name, opts...)
}

// SyncInt64UpDownCounter provides a simple SyncInt64 UpDownCounter
func (i Instruments) SyncInt64UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	return i.SyncInt64().UpDownCounter(name, opts...)
}

// SyncInt64Histogram provides a simple SyncInt64 Histogram
func (i Instruments) SyncInt64Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	return i.SyncInt64().Histogram(name, opts...)
}

// SyncFloat64Counter provides a simple SyncFloat64 Counter
func (i Instruments) SyncFloat64Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	return i.SyncFloat64().Counter(name, opts...)
}

// SyncFloat64UpDownCounter provides a simple SyncFloat64 UpDownCounter
func (i Instruments) SyncFloat64UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	return i.SyncFloat64().UpDownCounter(name, opts...)
}

// SyncFloat64Histogram provides a simple SyncFloat64 Histogram
func (i Instruments) SyncFloat64Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	return i.SyncFloat64().Histogram(name, opts...)
}

// AsyncInt64Counter provides a simple AsyncInt64C Counter
func (i Instruments) AsyncInt64Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	return i.AsyncInt64().Counter(name, opts...)
}

// AsyncInt64UpDownCounter provides a simple AsyncInt64 UpDownCounter
func (i Instruments) AsyncInt64UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	return i.AsyncInt64().UpDownCounter(name, opts...)
}

// AsyncInt64Gauge provides a simple AsyncInt64 Gauge
func (i Instruments) AsyncInt64Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	return i.AsyncInt64().Gauge(name, opts...)
}

// AsyncFloat64Counter provides a simple AsyncFloat64 Counter
func (i Instruments) AsyncFloat64Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	return i.AsyncFloat64().Counter(name, opts...)
}

// AsyncFloat64UpDownCounter provides a simple AsyncFloat64 UpDownCounter
func (i Instruments) AsyncFloat64UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	return i.AsyncFloat64().UpDownCounter(name, opts...)
}

// AsyncFloat64Gauge provides a simple AsyncFloat64 Gauge
func (i Instruments) AsyncFloat64Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	return i.AsyncFloat64().Gauge(name, opts...)
}
