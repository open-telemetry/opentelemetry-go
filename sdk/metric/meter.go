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
	"errors"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
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

	pipes pipelines
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
			Scope: s,
			pipes: r.pipes,
		}
		r.meters = map[instrumentation.Scope]*meter{s: m}
		return m
	}

	m, ok := r.meters[s]
	if ok {
		return m
	}

	m = &meter{
		Scope: s,
		pipes: r.pipes,
	}
	r.meters[s] = m
	return m
}

type cache map[instrumentID]any

type cacheResult[N int64 | float64] struct {
	aggregators []internal.Aggregator[N]
	err         error
}

type querier[N int64 | float64] struct {
	sync.Mutex

	c cache
}

func newQuerier[N int64 | float64](c cache) *querier[N] {
	return &querier[N]{c: c}
}

var (
	errCacheMiss = errors.New("cache miss")
	errExists    = errors.New("instrument already exists for different number type")
)

func (q *querier[N]) Get(key instrumentID) (r *cacheResult[N], err error) {
	q.Lock()
	defer q.Unlock()

	vIface, ok := q.c[key]
	if !ok {
		err = errCacheMiss
		return r, err
	}

	switch v := vIface.(type) {
	case *cacheResult[N]:
		r = v
	default:
		err = errExists
	}
	return r, err
}

func (q *querier[N]) Set(key instrumentID, val *cacheResult[N]) {
	q.Lock()
	defer q.Unlock()

	q.c[key] = val
}

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	instrumentation.Scope

	pipes pipelines
	cache *cache
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// AsyncInt64 returns the asynchronous integer instrument provider.
func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	q := newQuerier[int64](*m.cache)
	return asyncInt64Provider{scope: m.Scope, resolve: newResolver(m.pipes, q)}
}

// AsyncFloat64 returns the asynchronous floating-point instrument provider.
func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	q := newQuerier[float64](*m.cache)
	return asyncFloat64Provider{scope: m.Scope, resolve: newResolver(m.pipes, q)}
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(insts []instrument.Asynchronous, f func(context.Context)) error {
	m.pipes.registerCallback(f)
	return nil
}

// SyncInt64 returns the synchronous integer instrument provider.
func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	q := newQuerier[int64](*m.cache)
	return syncInt64Provider{scope: m.Scope, resolve: newResolver(m.pipes, q)}
}

// SyncFloat64 returns the synchronous floating-point instrument provider.
func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	q := newQuerier[float64](*m.cache)
	return syncFloat64Provider{scope: m.Scope, resolve: newResolver(m.pipes, q)}
}
