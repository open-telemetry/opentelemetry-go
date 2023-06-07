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
	"regexp"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal"
)

var (
	instrumentNameRe         = regexp.MustCompile(`^([A-Za-z]){1}([A-Za-z0-9\_\-\.]){0,62}$`)
	ErrInvalidInstrumentName = errors.New("invalid instrument name. Instrument names must consist of 63 or fewer characters including alphanumeric, _, ., -, and start with a letter")
)

// meter handles the creation and coordination of all metric instruments. A
// meter represents a single instrumentation scope; all metric telemetry
// produced by an instrumentation scope will use metric instruments from a
// single meter.
type meter struct {
	embedded.Meter

	scope instrumentation.Scope
	pipes pipelines

	int64IP   *int64InstProvider
	float64IP *float64InstProvider
}

func newMeter(s instrumentation.Scope, p pipelines) *meter {
	// viewCache ensures instrument conflicts, including number conflicts, this
	// meter is asked to create are logged to the user.
	var viewCache cache[string, streamID]

	return &meter{
		scope:     s,
		pipes:     p,
		int64IP:   newInt64InstProvider(s, p, &viewCache),
		float64IP: newFloat64InstProvider(s, p, &viewCache),
	}
}

// Compile-time check meter implements metric.Meter.
var _ metric.Meter = (*meter)(nil)

// Int64Counter returns a new instrument identified by name and configured with
// options. The instrument is used to synchronously record increasing int64
// measurements during a computational operation.
func (m *meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {

	cfg := metric.NewInt64CounterConfig(options...)
	const kind = InstrumentKindCounter
	i, err := m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Int64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// int64 measurements during a computational operation.
func (m *meter) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	cfg := metric.NewInt64UpDownCounterConfig(options...)
	const kind = InstrumentKindUpDownCounter
	i, err := m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Int64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of int64 measurements during a computational operation.
func (m *meter) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	cfg := metric.NewInt64HistogramConfig(options...)
	const kind = InstrumentKindHistogram
	i, err := m.int64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Int64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableCounter(name string, options ...metric.Int64ObservableCounterOption) (metric.Int64ObservableCounter, error) {
	cfg := metric.NewInt64ObservableCounterConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

// Int64ObservableUpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableUpDownCounter(name string, options ...metric.Int64ObservableUpDownCounterOption) (metric.Int64ObservableUpDownCounter, error) {
	cfg := metric.NewInt64ObservableUpDownCounterConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

// Int64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous int64 measurements once per a measurement collection cycle.
func (m *meter) Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
	cfg := metric.NewInt64ObservableGaugeConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := int64ObservProvider{m.int64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

// Float64Counter returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record increasing
// float64 measurements during a computational operation.
func (m *meter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	cfg := metric.NewFloat64CounterConfig(options...)
	const kind = InstrumentKindCounter
	i, err := m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Float64UpDownCounter returns a new instrument identified by name and
// configured with options. The instrument is used to synchronously record
// float64 measurements during a computational operation.
func (m *meter) Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	cfg := metric.NewFloat64UpDownCounterConfig(options...)
	const kind = InstrumentKindUpDownCounter
	i, err := m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Float64Histogram returns a new instrument identified by name and configured
// with options. The instrument is used to synchronously record the
// distribution of float64 measurements during a computational operation.
func (m *meter) Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	cfg := metric.NewFloat64HistogramConfig(options...)
	const kind = InstrumentKindHistogram
	i, err := m.float64IP.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return i, err
	}

	return i, validateInstrumentName(name)
}

// Float64ObservableCounter returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// increasing float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableCounter(name string, options ...metric.Float64ObservableCounterOption) (metric.Float64ObservableCounter, error) {
	cfg := metric.NewFloat64ObservableCounterConfig(options...)
	const kind = InstrumentKindObservableCounter
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

// Float64ObservableUpDownCounter returns a new instrument identified by name
// and configured with options. The instrument is used to asynchronously record
// float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableUpDownCounter(name string, options ...metric.Float64ObservableUpDownCounterOption) (metric.Float64ObservableUpDownCounter, error) {
	cfg := metric.NewFloat64ObservableUpDownCounterConfig(options...)
	const kind = InstrumentKindObservableUpDownCounter
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

// Float64ObservableGauge returns a new instrument identified by name and
// configured with options. The instrument is used to asynchronously record
// instantaneous float64 measurements once per a measurement collection cycle.
func (m *meter) Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
	cfg := metric.NewFloat64ObservableGaugeConfig(options...)
	const kind = InstrumentKindObservableGauge
	p := float64ObservProvider{m.float64IP}
	inst, err := p.lookup(kind, name, cfg.Description(), cfg.Unit())
	if err != nil {
		return nil, err
	}
	p.registerCallbacks(inst, cfg.Callbacks())
	return inst, validateInstrumentName(name)
}

func validateInstrumentName(name string) error {
	if !instrumentNameRe.MatchString(name) {
		return ErrInvalidInstrumentName
	}
	return nil
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
func (m *meter) RegisterCallback(f metric.Callback, insts ...metric.Observable) (metric.Registration, error) {
	if len(insts) == 0 {
		// Don't allocate a observer if not needed.
		return noopRegister{}, nil
	}

	reg := newObserver()
	var errs multierror
	for _, inst := range insts {
		// Unwrap any global.
		if u, ok := inst.(interface {
			Unwrap() metric.Observable
		}); ok {
			inst = u.Unwrap()
		}

		switch o := inst.(type) {
		case int64Observable:
			if err := o.registerable(m.scope); err != nil {
				if !errors.Is(err, errEmptyAgg) {
					errs.append(err)
				}
				continue
			}
			reg.registerInt64(o.observablID)
		case float64Observable:
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

func (r observer) ObserveFloat64(o metric.Float64Observable, v float64, opts ...metric.ObserveOption) {
	var oImpl float64Observable
	switch conv := o.(type) {
	case float64Observable:
		oImpl = conv
	case interface {
		Unwrap() metric.Observable
	}:
		// Unwrap any global.
		async := conv.Unwrap()
		var ok bool
		if oImpl, ok = async.(float64Observable); !ok {
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
	c := metric.NewObserveConfig(opts)
	oImpl.observe(v, c.Attributes())
}

func (r observer) ObserveInt64(o metric.Int64Observable, v int64, opts ...metric.ObserveOption) {
	var oImpl int64Observable
	switch conv := o.(type) {
	case int64Observable:
		oImpl = conv
	case interface {
		Unwrap() metric.Observable
	}:
		// Unwrap any global.
		async := conv.Unwrap()
		var ok bool
		if oImpl, ok = async.(int64Observable); !ok {
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
	c := metric.NewObserveConfig(opts)
	oImpl.observe(v, c.Attributes())
}

type noopRegister struct{ embedded.Registration }

func (noopRegister) Unregister() error {
	return nil
}

// int64InstProvider provides int64 OpenTelemetry instruments.
type int64InstProvider struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[int64]
}

func newInt64InstProvider(s instrumentation.Scope, p pipelines, c *cache[string, streamID]) *int64InstProvider {
	return &int64InstProvider{scope: s, pipes: p, resolve: newResolver[int64](p, c)}
}

func (p *int64InstProvider) aggs(kind InstrumentKind, name, desc, u string) ([]internal.Aggregator[int64], error) {
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
func (p *int64InstProvider) lookup(kind InstrumentKind, name, desc, u string) (*int64Inst, error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return &int64Inst{aggregators: aggs}, err
}

// float64InstProvider provides float64 OpenTelemetry instruments.
type float64InstProvider struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[float64]
}

func newFloat64InstProvider(s instrumentation.Scope, p pipelines, c *cache[string, streamID]) *float64InstProvider {
	return &float64InstProvider{scope: s, pipes: p, resolve: newResolver[float64](p, c)}
}

func (p *float64InstProvider) aggs(kind InstrumentKind, name, desc, u string) ([]internal.Aggregator[float64], error) {
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
func (p *float64InstProvider) lookup(kind InstrumentKind, name, desc, u string) (*float64Inst, error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return &float64Inst{aggregators: aggs}, err
}

type int64ObservProvider struct{ *int64InstProvider }

func (p int64ObservProvider) lookup(kind InstrumentKind, name, desc, u string) (int64Observable, error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return newInt64Observable(p.scope, kind, name, desc, u, aggs), err
}

func (p int64ObservProvider) registerCallbacks(inst int64Observable, cBacks []metric.Int64Callback) {
	if inst.observable == nil || len(inst.aggregators) == 0 {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p int64ObservProvider) callback(i int64Observable, f metric.Int64Callback) func(context.Context) error {
	inst := int64Observer{int64Observable: i}
	return func(ctx context.Context) error { return f(ctx, inst) }
}

type int64Observer struct {
	embedded.Int64Observer
	int64Observable
}

func (o int64Observer) Observe(val int64, opts ...metric.ObserveOption) {
	c := metric.NewObserveConfig(opts)
	o.observe(val, c.Attributes())
}

type float64ObservProvider struct{ *float64InstProvider }

func (p float64ObservProvider) lookup(kind InstrumentKind, name, desc, u string) (float64Observable, error) {
	aggs, err := p.aggs(kind, name, desc, u)
	return newFloat64Observable(p.scope, kind, name, desc, u, aggs), err
}

func (p float64ObservProvider) registerCallbacks(inst float64Observable, cBacks []metric.Float64Callback) {
	if inst.observable == nil || len(inst.aggregators) == 0 {
		// Drop aggregator.
		return
	}

	for _, cBack := range cBacks {
		p.pipes.registerCallback(p.callback(inst, cBack))
	}
}

func (p float64ObservProvider) callback(i float64Observable, f metric.Float64Callback) func(context.Context) error {
	inst := float64Observer{float64Observable: i}
	return func(ctx context.Context) error { return f(ctx, inst) }
}

type float64Observer struct {
	embedded.Float64Observer
	float64Observable
}

func (o float64Observer) Observe(val float64, opts ...metric.ObserveOption) {
	c := metric.NewObserveConfig(opts)
	o.observe(val, c.Attributes())
}
