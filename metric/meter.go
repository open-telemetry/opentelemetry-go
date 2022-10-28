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
)

// MeterProvider provides access to named Meter instances, for instrumenting
// an application or library.
type MeterProvider interface {
	// Meter creates an instance of a `Meter` interface. The instrumentationName
	// must be the name of the library providing instrumentation. This name may
	// be the same as the instrumented code only if that code provides built-in
	// instrumentation. If the instrumentationName is empty, then a
	// implementation defined default name will be used instead.
	Meter(instrumentationName string, opts ...MeterOption) Meter
}

// Meter provides access to instrument instances for recording metrics.
type Meter interface {
	// AsyncInt64 is the namespace for the Asynchronous Integer instruments.
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncInt64() instrument.AsyncInstrumentProvider[int64]

	// AsyncFloat64 is the namespace for the Asynchronous Float instruments
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncFloat64() instrument.AsyncInstrumentProvider[float64]

	// RegisterCallback captures the function that will be called during Collect.
	//
	// It is only valid to call Observe within the scope of the passed function,
	// and only on the instruments that were registered with this call.
	RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error

	// SyncInt64 is the namespace for the Synchronous Integer instruments
	SyncInt64() instrument.SyncInstrumentProvider[int64]
	// SyncFloat64 is the namespace for the Synchronous Float instruments
	SyncFloat64() instrument.SyncInstrumentProvider[float64]
}
