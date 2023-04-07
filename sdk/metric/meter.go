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
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	embedded.Meter

	scope instrumentation.Scope
	pipes pipelines

	int64IP   *instProvider[int64]
	float64IP *instProvider[float64]
}

func newMeter(s instrumentation.Scope, p pipelines) *meter {
	// viewCache ensures instrument conflicts, including number conflicts, this
	// meter is asked to create are logged to the user.
	var viewCache cache[string, streamID]

	return &meter{
		scope:     s,
		pipes:     p,
		int64IP:   newInstProvider[int64](s, p, &viewCache),
		float64IP: newInstProvider[float64](s, p, &viewCache),
	}
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// Int64Counter returns a new instrument identified by name and configured with
// options. The instrument is used to synchronously record increasing int64
// measurements during a computational operation.
func (m *meter) Int64Counter(name string, options ...instrument.CounterOption[int64]) (instrument.Counter[int64], error) {
	cfg := instrument.NewCounterConfig(options...)
	const kind = InstrumentKindCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// int64 measurements during a computational operation.
func (m *meter) Int64UpDownCounter(name string, options ...instrument.UpDownCounterOption[int64]) (instrument.UpDownCounter[int64], error) {
	cfg := instrument.NewUpDownCounterConfig(options...)
	const kind = InstrumentKindUpDownCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of int64 measurements during a computational operation.
func (m *meter) Int64Histogram(name string, options ...instrument.HistogramOption[int64]) (instrument.Histogram[int64], error) {
	cfg := instrument.NewHistogramConfig(options...)
	const kind = InstrumentKindHistogram
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableCounter(name string, options ...instrument.ObservableCounterOption[int64]) (instrument.ObservableCounter[int64], error) {
	cfg := instrument.NewObservableCounterConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := observProvider[int64]{m.int64IP}
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
func (m *meter) Int64ObservableUpDownCounter(name string, options ...instrument.ObservableUpDownCounterOption[int64]) (instrument.ObservableUpDownCounter[int64], error) {
	cfg := instrument.NewObservableUpDownCounterConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := observProvider[int64]{m.int64IP}
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
func (m *meter) Int64ObservableGauge(name string, options ...instrument.ObservableGaugeOption[int64]) (instrument.ObservableGauge[int64], error) {
	cfg := instrument.NewObservableGaugeConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := observProvider[int64]{m.int64IP}
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
func (m *meter) Float64Counter(name string, options ...instrument.CounterOption[float64]) (instrument.Counter[float64], error) {
	cfg := instrument.NewCounterConfig(options...)
	const kind = InstrumentKindCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// float64 measurements during a computational operation.
func (m *meter) Float64UpDownCounter(name string, options ...instrument.UpDownCounterOption[float64]) (instrument.UpDownCounter[float64], error) {
	cfg := instrument.NewUpDownCounterConfig(options...)
	const kind = InstrumentKindUpDownCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of float64 measurements during a computational operation.
func (m *meter) Float64Histogram(name string, options ...instrument.HistogramOption[float64]) (instrument.Histogram[float64], error) {
	cfg := instrument.NewHistogramConfig(options...)
	const kind = InstrumentKindHistogram
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableCounter(name string, options ...instrument.ObservableCounterOption[float64]) (instrument.ObservableCounter[float64], error) {
	cfg := instrument.NewObservableCounterConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := observProvider[float64]{m.float64IP}
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
func (m *meter) Float64ObservableUpDownCounter(name string, options ...instrument.ObservableUpDownCounterOption[float64]) (instrument.ObservableUpDownCounter[float64], error) {
	cfg := instrument.NewObservableUpDownCounterConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := observProvider[float64]{m.float64IP}
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
func (m *meter) Float64ObservableGauge(name string, options ...instrument.ObservableGaugeOption[float64]) (instrument.ObservableGauge[float64], error) {
	cfg := instrument.NewObservableGaugeConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := observProvider[float64]{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, nil
}

// RegisterCallback registers f to be called each collection cycle so it will
// make observations for insts during those cycles.
//
// The only instruments f can make observations for are insts. All other
// observations will be dropped and an error will be logged.
//
// Only instruments from this meter can be registered with f, an error is
// returned if other instrument are provided.
//
// The returned Registration can be used to unregister f.
func (m *meter) RegisterCallback(f metric.Callback, insts ...instrument.Observable) (metric.Registration, error) {
	if len(insts) == 0 {
		// Don't allocate a observer if not needed.
		return noopRegister{}, nil
	}

	reg := newObserver()
	var errs multierror
	for _, inst := range insts {
		// Unwrap any global.
		if u, ok := inst.(interface {
			Unwrap() instrument.Observable
		}); ok {
			inst = u.Unwrap()
		}

		switch o := inst.(type) {
		case *observable[int64]:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerInt64(o.observablID)
		case *observable[float64]:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerFloat64(o.observablID)
		default:
			// Instrument external to the SDK.
			return nil, fmt.Errorf("invalid observable: from different implementation: %T", inst)
		}
	}

	if err := errs.errorOrNil(); err != nil {
		return nil, err
	}

	if reg.len() == 0 {
		// All insts use drop aggregation.
		return noopRegister{}, nil
	}

	cback := func(ctx context.Context) error {
		return f(ctx, reg)
	}
	return m.pipes.registerMultiCallback(cback), nil
}

type observer struct {
	embedded.Observer

	float64 map[observablID[float64]]struct{}
	int64   map[observablID[int64]]struct{}
}

func newObserver() observer {
	return observer{
		float64: make(map[observablID[float64]]struct{}),
		int64:   make(map[observablID[int64]]struct{}),
	}
}

func (r observer) len() int {
	return len(r.float64) + len(r.int64)
}

func (r observer) registerFloat64(id observablID[float64]) {
	r.float64[id] = struct{}{}
}

func (r observer) registerInt64(id observablID[int64]) {
	r.int64[id] = struct{}{}
}

var (
	errUnknownObserver = errors.New("unknown observable instrument")
	errUnregObserver   = errors.New("observable instrument not registered for callback")
)

func (r observer) ObserveFloat64(o instrument.ObservableT[float64], v float64, a ...attribute.KeyValue) {
	var oImpl *observable[float64]
	switch conv := o.(type) {
	case *observable[float64]:
		oImpl = conv
	case interface {
		Unwrap() instrument.Observable
	}:
		// Unwrap any global.
		async := conv.Unwrap()
		var ok bool
		if oImpl, ok = async.(*observable[float64]); !ok {
			global.Error(errUnknownObserver, "failed to record asynchronous")
			return
		}
	default:
		global.Error(errUnknownObserver, "failed to record")
		return
	}

	if _, registered := r.float64[oImpl.observablID]; !registered {
		global.Error(errUnregObserver, "failed to record",
			"name", oImpl.name,
			"description", oImpl.description,
			"unit", oImpl.unit,
			"number", fmt.Sprintf("%T", float64(0)),
		)
		return
	}
	oImpl.observe(v, a)
}

func (r observer) ObserveInt64(o instrument.ObservableT[int64], v int64, a ...attribute.KeyValue) {
	var oImpl *observable[int64]
	switch conv := o.(type) {
	case *observable[int64]:
		oImpl = conv
	case interface {
		Unwrap() instrument.Observable
	}:
		// Unwrap any global.
		async := conv.Unwrap()
		var ok bool
		if oImpl, ok = async.(*observable[int64]); !ok {
			global.Error(errUnknownObserver, "failed to record asynchronous")
			return
		}
	default:
		global.Error(errUnknownObserver, "failed to record")
		return
	}

	if _, registered := r.int64[oImpl.observablID]; !registered {
		global.Error(errUnregObserver, "failed to record",
			"name", oImpl.name,
			"description", oImpl.description,
			"unit", oImpl.unit,
			"number", fmt.Sprintf("%T", int64(0)),
		)
		return
	}
	oImpl.observe(v, a)
}

type noopRegister struct{ embedded.Registration }

func (noopRegister) Unregister() error {
	return nil
}

// instProvider provides all OpenTelemetry instruments.
type instProvider[N int64 | float64] struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[N]
}

func newInstProvider[N int64 | float64](s instrumentation.Scope, p pipelines, c *cache[string, streamID]) *instProvider[N] {
	return &instProvider[N]{scope: s, pipes: p, resolve: newResolver[N](p, c)}
}

func (p *instProvider[N]) aggs(kind InstrumentKind, name, desc, u string) ([]internal.Aggregator[N], error) {
	inst := Instrument{
		Name:        name,
		Description: desc,
		Unit:        u,
		Kind:        kind,
		Scope:       p.scope,
	}
	return p.resolve.Aggregators(inst)
}

// lookup returns the resolved instrumentImpl.
func (p *instProvider[N]) lookup(kind InstrumentKind, name, desc, u string) (*instrumentImpl[N], error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return &instrumentImpl[N]{aggregators: aggs}, err
}

type observProvider[N int64 | float64] struct {
	*instProvider[N]
}

func (p observProvider[N]) lookup(kind InstrumentKind, name, desc, u string) (*observable[N], error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return newObservable(p.scope, kind, name, desc, u, aggs), err
}

func (p observProvider[N]) registerCallbacks(inst *observable[N], cBacks []instrument.Callback[N]) {
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p observProvider[N]) callback(i *observable[N], f instrument.Callback[N]) func(context.Context) error {
	inst := observerT[N]{observable: i}
	return func(ctx context.Context) error { return f(ctx, inst) }
}

type observerT[N int64 | float64] struct {
	embedded.ObserverT[N]
	*observable[N]
}

func (o observerT[N]) Observe(val N, attrs ...attribute.KeyValue) {
	o.observe(val, attrs)
}
