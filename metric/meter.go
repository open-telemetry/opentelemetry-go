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
	// AsyncInt64 is the namespace for the Asynchronous Integer instruments.
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncInt64() asyncint64.InstrumentProvider

	// AsyncFloat64 is the namespace for the Asynchronous Float instruments
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncFloat64() asyncfloat64.InstrumentProvider

	// SyncInt64 is the namespace for the Synchronous Integer instruments
	SyncInt64() syncint64.InstrumentProvider
	// SyncFloat64 is the namespace for the Synchronous Float instruments
	SyncFloat64() syncfloat64.InstrumentProvider

	// RegisterCallback registers f to be called during the collection of a
	// measurement cycle.
	//
	// If Unregister of the returned Registration is called, f needs to be
	// unregistered and not called during collection.
	//
	// The instruments f is registered with are the only instruments that f may
	// observe values for.
	//
	// If no instruments are passed, f should not be registered nor called
	// during collection.
	RegisterCallback(instruments []instrument.Asynchronous, f func(context.Context)) (Registration, error)
}

// Registration is an token representing the unique registration of a callback
// for a set of instruments with a Meter.
type Registration interface {
	// Unregister removes the callback registration from a Meter.
	//
	// This method needs to be idempotent and concurrent safe.
	Unregister() error
}
