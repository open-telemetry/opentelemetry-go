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

	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// MeterProvider provides access to named Meter instances, for instrumenting
// an application or library.
//
// Warning: methods may be added to this interface in minor releases.
type MeterProvider interface {
	// Meter creates an instance of a `Meter` interface. The instrumentationName
	// must be the name of the library providing instrumentation. This name may
	// be the same as the instrumented code only if that code provides built-in
	// instrumentation. If the instrumentationName is empty, then a
	// implementation defined default name will be used instead.
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter provides access to instrument instances for recording metrics.
//
// Warning: methods may be added to this interface in minor releases.
type Meter interface {
	// Int64Counter returns a new instrument identified by name and configured
	// with options. The instrument is used to synchronously record increasing
	// int64 measurements during a computational operation.
	Int64Counter(name string, options ...instrument.Option) (syncint64.Counter, error)
	// Int64UpDownCounter returns a new instrument identified by name and
	// configured with options. The instrument is used to synchronously record
	// int64 measurements during a computational operation.
	Int64UpDownCounter(name string, options ...instrument.Option) (syncint64.UpDownCounter, error)
	// Int64Histogram returns a new instrument identified by name and
	// configured with options. The instrument is used to synchronously record
	// the distribution of int64 measurements during a computational operation.
	Int64Histogram(name string, options ...instrument.Option) (syncint64.Histogram, error)
	// Int64ObservableCounter returns a new instrument identified by name and
	// configured with options. The instrument is used to asynchronously record
	// increasing int64 measurements once per a measurement collection cycle.
	Int64ObservableCounter(name string, options ...instrument.Option) (asyncint64.Counter, error)
	// Int64ObservableUpDownCounter returns a new instrument identified by name
	// and configured with options. The instrument is used to asynchronously
	// record int64 measurements once per a measurement collection cycle.
	Int64ObservableUpDownCounter(name string, options ...instrument.Option) (asyncint64.UpDownCounter, error)
	// Int64ObservableGauge returns a new instrument identified by name and
	// configured with options. The instrument is used to asynchronously record
	// instantaneous int64 measurements once per a measurement collection
	// cycle.
	Int64ObservableGauge(name string, options ...instrument.Option) (asyncint64.Gauge, error)

	// Float64Counter returns a new instrument identified by name and
	// configured with options. The instrument is used to synchronously record
	// increasing float64 measurements during a computational operation.
	Float64Counter(name string, options ...instrument.Option) (syncfloat64.Counter, error)
	// Float64UpDownCounter returns a new instrument identified by name and
	// configured with options. The instrument is used to synchronously record
	// float64 measurements during a computational operation.
	Float64UpDownCounter(name string, options ...instrument.Option) (syncfloat64.UpDownCounter, error)
	// Float64Histogram returns a new instrument identified by name and
	// configured with options. The instrument is used to synchronously record
	// the distribution of float64 measurements during a computational
	// operation.
	Float64Histogram(name string, options ...instrument.Option) (syncfloat64.Histogram, error)
	// Float64ObservableCounter returns a new instrument identified by name and
	// configured with options. The instrument is used to asynchronously record
	// increasing float64 measurements once per a measurement collection cycle.
	Float64ObservableCounter(name string, options ...instrument.Option) (asyncfloat64.Counter, error)
	// Float64ObservableUpDownCounter returns a new instrument identified by
	// name and configured with options. The instrument is used to
	// asynchronously record float64 measurements once per a measurement
	// collection cycle.
	Float64ObservableUpDownCounter(name string, options ...instrument.Option) (asyncfloat64.UpDownCounter, error)
	// Float64ObservableGauge returns a new instrument identified by name and
	// configured with options. The instrument is used to asynchronously record
	// instantaneous float64 measurements once per a measurement collection
	// cycle.
	Float64ObservableGauge(name string, options ...instrument.Option) (asyncfloat64.Gauge, error)

	// RegisterCallback captures the function that will be called during Collect.
	//
	// It is only valid to call Observe within the scope of the passed function,
	// and only on the instruments that were registered with this call.
	RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error
}
