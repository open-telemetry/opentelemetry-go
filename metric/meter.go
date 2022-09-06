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

// MeterProvider provides Meters that are used by instrumentation code to
// create instruments that measure code operations.
//
// A MeterProvider is the collection destination of all measurements made from
// instruments the provided Meters created, it represents a unique telemetry
// collection pipeline. How that pipeline is defined, meaning how those
// measurements are collected, processed, and where they are exported, depends
// on its implementation. Instrumentation authors do not need to define this
// implementation, rather just use the provided Meters to instrument code.
//
// Commonly, instrumentation code will accept a MeterProvider implementation at
// runtime from its users or it can simply use the globally registered one (see
// https://pkg.go.dev/go.opentelemetry.io/otel/metric/global#MeterProvider).
//
// Warning: methods may be added to this interface in minor releases.
type MeterProvider interface {
	// Meter returns a unique Meter scoped to be used by instrumentation code
	// to measure code operations. The scope and identity of that
	// instrumentation code is uniquely defined by the name and options passed.
	//
	// The passed name needs to uniquely identify instrumentation code.
	// Therefore, it is recommended that name is the Go package name of the
	// library providing instrumentation (note: not the code being
	// instrumented). Instrumentation libraries can have multiple versions,
	// therefore, the WithInstrumentationVersion option should be used to
	// distinguish these different codebases. Additionally, instrumentation
	// libraries may sometimes use metric measurements to communicate different
	// domains of code operations data (i.e. using different Meters to
	// communicate user experience and back-end operations). If this is the
	// case, the WithScopeAttributes option should be used to uniquely identify
	// Meters that handle the different domains of code operations data.
	//
	// If the same name and options are passed multiple times, the same Meter
	// will be returned (it is up to the implementation if this will be the
	// same underlying instance of that Meter or not). It is not necessary to
	// call this multiple times with the same name and options to get an
	// up-to-date Meter. All implementations will ensure any MeterProvider
	// configuration changes are propagated to all provided Meters.
	//
	// If name is empty, then an implementation defined default name will be
	// used instead.
	//
	// This method is safe to call concurrently.
	Meter(name string, options ...MeterOption) Meter
}

// Meter provides access to instrument instances for recording metrics.
type Meter interface {
	// AsyncInt64 is the namespace for the Asynchronous Integer instruments.
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncInt64() asyncint64.InstrumentProvider

	// AsyncFloat64 is the namespace for the Asynchronous Float instruments
	//
	// To Observe data with instruments it must be registered in a callback.
	AsyncFloat64() asyncfloat64.InstrumentProvider

	// RegisterCallback captures the function that will be called during Collect.
	//
	// It is only valid to call Observe within the scope of the passed function,
	// and only on the instruments that were registered with this call.
	RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error

	// SyncInt64 is the namespace for the Synchronous Integer instruments
	SyncInt64() syncint64.InstrumentProvider
	// SyncFloat64 is the namespace for the Synchronous Float instruments
	SyncFloat64() syncfloat64.InstrumentProvider
}
