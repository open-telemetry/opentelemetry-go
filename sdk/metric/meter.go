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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	pipes pipelines

	instProviderInt64   *instProvider[int64]
	instProviderFloat64 *instProvider[float64]
}

func newMeter(s instrumentation.Scope, p pipelines) *meter {
	// viewCache ensures instrument conflicts, including number conflicts, this
	// meter is asked to create are logged to the user.
	var viewCache cache[string, instrumentID]

	// Passing nil as the ac parameter to newInstrumentCache will have each
	// create its own aggregator cache.
	ic := newInstrumentCache[int64](nil, &viewCache)
	fc := newInstrumentCache[float64](nil, &viewCache)

	return &meter{
		pipes:               p,
		instProviderInt64:   newInstProvider(s, p, ic),
		instProviderFloat64: newInstProvider(s, p, fc),
	}
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// AsyncInt64 returns the asynchronous integer instrument provider.
func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	return asyncInt64Provider{m.instProviderInt64}
}

// AsyncFloat64 returns the asynchronous floating-point instrument provider.
func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	return asyncFloat64Provider{m.instProviderFloat64}
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(insts []instrument.Asynchronous, f func(context.Context)) error {
	for _, inst := range insts {
		// Only register if at least one instrument has a non-drop aggregation.
		// Otherwise, calling f during collection will be wasted computation.
		switch t := inst.(type) {
		case *instrumentImpl[int64]:
			if len(t.aggregators) > 0 {
				return m.registerCallback(f)
			}
		case *instrumentImpl[float64]:
			if len(t.aggregators) > 0 {
				return m.registerCallback(f)
			}
		default:
			// Instrument external to the SDK. For example, an instrument from
			// the "go.opentelemetry.io/otel/metric/internal/global" package.
			//
			// Fail gracefully here, assume a valid instrument.
			return m.registerCallback(f)
		}
	}
	// All insts use drop aggregation.
	return nil
}

func (m *meter) registerCallback(f func(context.Context)) error {
	m.pipes.registerCallback(f)
	return nil
}

// SyncInt64 returns the synchronous integer instrument provider.
func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	return syncInt64Provider{m.instProviderInt64}
}

// SyncFloat64 returns the synchronous floating-point instrument provider.
func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	return syncFloat64Provider{m.instProviderFloat64}
}
