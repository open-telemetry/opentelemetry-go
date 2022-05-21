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

//go:build go1.17
// +build go1.17

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	// TODO (#2821, #2815, 2814): implement.
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// AsyncInt64 returns the asynchronous integer instrument provider.
func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	// TODO (#2815): implement.
	return nil
}

// AsyncFloat64 returns the asynchronous floating-point instrument provider.
func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	// TODO (#2815): implement.
	return nil
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(insts []instrument.Asynchronous, f func(context.Context)) error {
	// TODO (#2815): implement.
	return nil
}

// SyncInt64 returns the synchronous integer instrument provider.
func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	// TODO (#2814): implement.
	return nil
}

// SyncFloat64 returns the synchronous floating-point instrument provider.
func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	// TODO (#2814): implement.
	return nil
}
