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
	"context"

	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument"
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
	AsyncInt64() asyncint64.Instruments
	AsyncFloat64() asyncfloat64.Instruments
	SyncInt64() syncint64.Instruments
	SyncFloat64() syncfloat64.Instruments

	NewCallback(insts []instrument.Asynchronous, function CallbackFunc) (Callback, error)
}

type CallbackFunc func(context.Context) error

type Callback interface {
	Instruments() []instrument.Asynchronous
}
