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
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/view"
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

func (m *meter) Float64Counter(name string, opts ...metric.InstrumentOption) (metric.Float64Counter, error) {
	return m.instProviderFloat64.Lookup(view.SyncCounter, name, opts)
}

func (m *meter) Float64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Float64UpDownCounter, error) {
	return m.instProviderFloat64.Lookup(view.SyncUpDownCounter, name, opts)
}

func (m *meter) Float64Histogram(name string, opts ...metric.InstrumentOption) (metric.Float64Histogram, error) {
	return m.instProviderFloat64.Lookup(view.SyncHistogram, name, opts)
}

func (m *meter) Float64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableCounter, error) {
	return m.instProviderFloat64.LookupObservable(view.AsyncCounter, name, opts)
}

func (m *meter) Float64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Float64ObservableUpDownCounter, error) {
	return m.instProviderFloat64.LookupObservable(view.AsyncUpDownCounter, name, opts)
}

func (m *meter) Float64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Float64ObservableGauge, error) {
	return m.instProviderFloat64.LookupObservable(view.AsyncGauge, name, opts)
}

func (m *meter) Int64Counter(name string, opts ...metric.InstrumentOption) (metric.Int64Counter, error) {
	return m.instProviderInt64.Lookup(view.SyncCounter, name, opts)
}

func (m *meter) Int64UpDownCounter(name string, opts ...metric.InstrumentOption) (metric.Int64UpDownCounter, error) {
	return m.instProviderInt64.Lookup(view.SyncUpDownCounter, name, opts)
}

func (m *meter) Int64Histogram(name string, opts ...metric.InstrumentOption) (metric.Int64Histogram, error) {
	return m.instProviderInt64.Lookup(view.SyncHistogram, name, opts)
}

func (m *meter) Int64ObservableCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableCounter, error) {
	return m.instProviderInt64.LookupObservable(view.AsyncCounter, name, opts)
}

func (m *meter) Int64ObservableUpDownCounter(name string, opts ...metric.ObservableOption) (metric.Int64ObservableUpDownCounter, error) {
	return m.instProviderInt64.LookupObservable(view.AsyncUpDownCounter, name, opts)
}

func (m *meter) Int64ObservableGauge(name string, opts ...metric.ObservableOption) (metric.Int64ObservableGauge, error) {
	return m.instProviderInt64.LookupObservable(view.AsyncGauge, name, opts)
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(f metric.Callback, instrument metric.Observable, additional ...metric.Observable) (metric.Unregisterer, error) {
	return m.pipes.registerCallback(f)
}
