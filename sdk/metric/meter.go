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
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	pipes pipelines

	int64IP   *instProvider[int64]
	float64IP *instProvider[float64]
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
		pipes:     p,
		int64IP:   newInstProvider(s, p, ic),
		float64IP: newInstProvider(s, p, fc),
	}
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// Int64Counter returns a new instrument identified by name and configured with
// options. The instrument is used to synchronously record increasing int64
// measurements during a computational operation.
func (m *meter) Int64Counter(name string, options ...instrument.Int64Option) (syncint64.Counter, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// int64 measurements during a computational operation.
func (m *meter) Int64UpDownCounter(name string, options ...instrument.Int64Option) (syncint64.UpDownCounter, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindUpDownCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of int64 measurements during a computational operation.
func (m *meter) Int64Histogram(name string, options ...instrument.Int64Option) (syncint64.Histogram, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindHistogram
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (asyncint64.Counter, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Int64ObservableUpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (asyncint64.UpDownCounter, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Int64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (asyncint64.Gauge, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Float64Counter returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record increasing
// float64 measurements during a computational operation.
func (m *meter) Float64Counter(name string, options ...instrument.Float64Option) (syncfloat64.Counter, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// float64 measurements during a computational operation.
func (m *meter) Float64UpDownCounter(name string, options ...instrument.Float64Option) (syncfloat64.UpDownCounter, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindUpDownCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of float64 measurements during a computational operation.
func (m *meter) Float64Histogram(name string, options ...instrument.Float64Option) (syncfloat64.Histogram, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindHistogram
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (asyncfloat64.Counter, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Float64ObservableUpDownCounter returns a new instrument identified by name
// and configured with options. The instrument is used to asynchronously record
// float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (asyncfloat64.UpDownCounter, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// Float64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (asyncfloat64.Gauge, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// RegisterCallback registers the function f to be called when any of the
// insts Collect method is called.
func (m *meter) RegisterCallback(insts []instrument.Asynchronous, f metric.Callback) (metric.Registration, error) {
	for _, inst := range insts {
		// Only register if at least one instrument has a non-drop aggregation.
		// Otherwise, calling f during collection will be wasted computation.
		switch t := inst.(type) {
		case *instrumentImpl[int64]:
			if len(t.aggregators) > 0 {
				return m.registerMultiCallback(f)
			}
		case *instrumentImpl[float64]:
			if len(t.aggregators) > 0 {
				return m.registerMultiCallback(f)
			}
		default:
			// Instrument external to the SDK. For example, an instrument from
			// the "go.opentelemetry.io/otel/metric/internal/global" package.
			//
			// Fail gracefully here, assume a valid instrument.
			return m.registerMultiCallback(f)
		}
	}
	// All insts use drop aggregation.
	return noopRegister{}, nil
}

type noopRegister struct{}

func (noopRegister) Unregister() error {
	return nil
}

func (m *meter) registerMultiCallback(c metric.Callback) (metric.Registration, error) {
	return m.pipes.registerMultiCallback(c), nil
}

// instProvider provides all OpenTelemetry instruments.
type instProvider[N int64 | float64] struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[N]
}

func newInstProvider[N int64 | float64](s instrumentation.Scope, p pipelines, c instrumentCache[N]) *instProvider[N] {
	return &instProvider[N]{scope: s, pipes: p, resolve: newResolver(p, c)}
}

// lookup returns the resolved instrumentImpl.
func (p *instProvider[N]) lookup(kind InstrumentKind, name, desc string, u unit.Unit) (*instrumentImpl[N], error) {
	inst := Instrument{
		Name:        name,
		Description: desc,
		Unit:        u,
		Kind:        kind,
		Scope:       p.scope,
	}
	aggs, err := p.resolve.Aggregators(inst)
	return &instrumentImpl[N]{aggregators: aggs}, err
}

type int64ObservProvider struct{ *instProvider[int64] }

func (p int64ObservProvider) registerCallbacks(inst *instrumentImpl[int64], cBacks []instrument.Int64Callback) {
	if inst == nil {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p int64ObservProvider) callback(i *instrumentImpl[int64], f instrument.Int64Callback) func(context.Context) error {
	return func(ctx context.Context) error { return f(ctx, i) }
}

type float64ObservProvider struct{ *instProvider[float64] }

func (p float64ObservProvider) registerCallbacks(inst *instrumentImpl[float64], cBacks []instrument.Float64Callback) {
	if inst == nil {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p float64ObservProvider) callback(i *instrumentImpl[float64], f instrument.Float64Callback) func(context.Context) error {
	return func(ctx context.Context) error { return f(ctx, i) }
}
