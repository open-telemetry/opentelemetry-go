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
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// meterRegistry keeps a record of initialized meters for instrumentation
// scopes. A meter is unique to an instrumentation scope and if multiple
// requests for that meter are made a meterRegistry ensure the same instance
// is used.
//
// The zero meterRegistry is empty and ready for use.
//
// A meterRegistry must not be copied after first use.
//
// All methods of a meterRegistry are safe to call concurrently.
type meterRegistry struct {
	sync.Mutex

	meters map[instrumentation.Scope]*meter

	registry *pipelineRegistry
}

// Get returns a registered meter matching the instrumentation scope if it
// exists in the meterRegistry. Otherwise, a new meter configured for the
// instrumentation scope is registered and then returned.
//
// Get is safe to call concurrently.
func (r *meterRegistry) Get(s instrumentation.Scope) *meter {
	r.Lock()
	defer r.Unlock()

	if r.meters == nil {
		m := &meter{
			Scope:    s,
			registry: r.registry,
		}
		r.meters = map[instrumentation.Scope]*meter{s: m}
		return m
	}

	m, ok := r.meters[s]
	if ok {
		return m
	}

	m = &meter{
		Scope:    s,
		registry: r.registry,
	}
	r.meters[s] = m
	return m
}

// Range calls f sequentially for each meter present in the meterRegistry. If
// f returns false, the iteration is stopped.
//
// Range is safe to call concurrently.
func (r *meterRegistry) Range(f func(*meter) bool) {
	r.Lock()
	defer r.Unlock()

	for _, m := range r.meters {
		if !f(m) {
			return
		}
	}
}

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	instrumentation.Scope

	registry *pipelineRegistry
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// AsyncInt64 returns the asynchronous integer instrument provider.
func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	return asyncInt64Provider{scope: m.Scope, registry: m.registry}
}

// AsyncFloat64 returns the asynchronous floating-point instrument provider.
func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	return asyncFloat64Provider{scope: m.Scope, registry: m.registry}
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(insts []instrument.Asynchronous, f func(context.Context)) error {
	m.registry.registerCallback(f)
	return nil
}

// SyncInt64 returns the synchronous integer instrument provider.
func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	return syncInt64Provider{scope: m.Scope, registry: m.registry}
}

// SyncFloat64 returns the synchronous floating-point instrument provider.
func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	return syncFloat64Provider{scope: m.Scope, registry: m.registry}
}
