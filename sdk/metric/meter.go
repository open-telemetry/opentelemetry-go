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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	metric.Meter

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
func (m *meter) Int64Counter(name string, options ...instrument.Int64Option) (instrument.Int64Counter, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// int64 measurements during a computational operation.
func (m *meter) Int64UpDownCounter(name string, options ...instrument.Int64Option) (instrument.Int64UpDownCounter, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindUpDownCounter
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Int64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of int64 measurements during a computational operation.
func (m *meter) Int64Histogram(name string, options ...instrument.Int64Option) (instrument.Int64Histogram, error) {
	cfg := instrument.NewInt64Config(options...)
	const kind = InstrumentKindHistogram
	return m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

type int64ObservableCounter struct {
	instrument.Int64ObservableCounter
	*observable[int64]
}

// Int64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableCounter, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableCounter
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.int64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := int64ObservableCounter{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, int64Observer{observe: inst.observe})
		})
	}
	return inst, nil
}

type int64ObservableUpDownCounter struct {
	instrument.Int64ObservableUpDownCounter
	*observable[int64]
}

// Int64ObservableUpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableUpDownCounter(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableUpDownCounter, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.int64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := int64ObservableUpDownCounter{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, int64Observer{observe: inst.observe})
		})
	}
	return inst, nil
}

type int64ObservableGauge struct {
	instrument.Int64ObservableGauge
	*observable[int64]
}

// Int64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableGauge(name string, options ...instrument.Int64ObserverOption) (instrument.Int64ObservableGauge, error) {
	cfg := instrument.NewInt64ObserverConfig(options...)
	const kind = InstrumentKindObservableGauge
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.int64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := int64ObservableGauge{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, int64Observer{observe: inst.observe})
		})
	}
	return inst, nil
}

// Float64Counter returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record increasing
// float64 measurements during a computational operation.
func (m *meter) Float64Counter(name string, options ...instrument.Float64Option) (instrument.Float64Counter, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// float64 measurements during a computational operation.
func (m *meter) Float64UpDownCounter(name string, options ...instrument.Float64Option) (instrument.Float64UpDownCounter, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindUpDownCounter
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

// Float64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of float64 measurements during a computational operation.
func (m *meter) Float64Histogram(name string, options ...instrument.Float64Option) (instrument.Float64Histogram, error) {
	cfg := instrument.NewFloat64Config(options...)
	const kind = InstrumentKindHistogram
	return m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
}

type float64ObservableCounter struct {
	instrument.Float64ObservableCounter
	*observable[float64]
}

// Float64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableCounter, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableCounter
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.float64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := float64ObservableCounter{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, float64Observer{observe: inst.observe})
		})
	}
	return inst, nil
}

type float64ObservableUpDownCounter struct {
	instrument.Float64ObservableUpDownCounter
	*observable[float64]
}

// Float64ObservableUpDownCounter returns a new instrument identified by name
// and configured with options. The instrument is used to asynchronously record
// float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableUpDownCounter(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableUpDownCounter, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.float64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := float64ObservableUpDownCounter{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, float64Observer{observe: inst.observe})
		})
	}
	return inst, nil
}

type float64ObservableGauge struct {
	instrument.Float64ObservableGauge
	*observable[float64]
}

// Float64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableGauge(name string, options ...instrument.Float64ObserverOption) (instrument.Float64ObservableGauge, error) {
	cfg := instrument.NewFloat64ObserverConfig(options...)
	const kind = InstrumentKindObservableGauge
	desc := cfg.Description()
	u := cfg.Unit()

	aggs, err := m.float64IP.aggs(kind, name, desc, u)
	if err != nil {
		return nil, err
	}
	inst := float64ObservableGauge{
		observable: newObservable(m.scope, kind, name, desc, u, aggs),
	}
	if len(inst.aggregators) == 0 {
		// Drop aggregator.
		return inst, nil
	}
	for _, cBack := range cfg.Callbacks() {
		m.pipes.registerCallback(func(ctx context.Context) error {
			return cBack(ctx, float64Observer{observe: inst.observe})
		})
	}
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
		return noop.Registration{}, nil
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
		case int64ObservableCounter:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerInt64(o.observablID)
		case int64ObservableUpDownCounter:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerInt64(o.observablID)
		case int64ObservableGauge:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerInt64(o.observablID)
		case float64ObservableCounter:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerFloat64(o.observablID)
		case float64ObservableUpDownCounter:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerFloat64(o.observablID)
		case float64ObservableGauge:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerFloat64(o.observablID)
		default:
			// Instrument external to the SDK.
			return nil, fmt.Errorf("invalid observable: from different implementation")
		}
	}

	if err := errs.errorOrNil(); err != nil {
		return nil, err
	}

	if reg.len() == 0 {
		// All insts use drop aggregation.
		return noop.Registration{}, nil
	}

	cback := func(ctx context.Context) error {
		return f(ctx, reg)
	}
	return m.pipes.registerMultiCallback(cback), nil
}

type observer struct {
	metric.Observer

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

func (r observer) ObserveFloat64(o instrument.Float64Observable, v float64, a ...attribute.KeyValue) {
	var oImpl *observable[float64]
	switch conv := o.(type) {
	case float64ObservableCounter:
		oImpl = conv.observable
	case float64ObservableUpDownCounter:
		oImpl = conv.observable
	case float64ObservableGauge:
		oImpl = conv.observable
	case interface {
		Unwrap() instrument.Observable
	}:
		// Unwrap any global.
		switch unwrapped := conv.Unwrap().(type) {
		case float64ObservableCounter:
			oImpl = unwrapped.observable
		case float64ObservableUpDownCounter:
			oImpl = unwrapped.observable
		case float64ObservableGauge:
			oImpl = unwrapped.observable
		default:
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

func (r observer) ObserveInt64(o instrument.Int64Observable, v int64, a ...attribute.KeyValue) {
	var oImpl *observable[int64]
	switch conv := o.(type) {
	case int64ObservableCounter:
		oImpl = conv.observable
	case int64ObservableUpDownCounter:
		oImpl = conv.observable
	case int64ObservableGauge:
		oImpl = conv.observable
	case interface {
		Unwrap() instrument.Observable
	}:
		// Unwrap any global.
		switch unwrapped := conv.Unwrap().(type) {
		case int64ObservableCounter:
			oImpl = unwrapped.observable
		case int64ObservableUpDownCounter:
			oImpl = unwrapped.observable
		case int64ObservableGauge:
			oImpl = unwrapped.observable
		default:
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

type int64Observer struct {
	instrument.Int64Observer
	observe func(int64, []attribute.KeyValue)
}

func (o int64Observer) Observe(val int64, attrs ...attribute.KeyValue) {
	o.observe(val, attrs)
}

type float64Observer struct {
	instrument.Float64Observer
	observe func(float64, []attribute.KeyValue)
}

func (o float64Observer) Observe(val float64, attrs ...attribute.KeyValue) {
	o.observe(val, attrs)
}
