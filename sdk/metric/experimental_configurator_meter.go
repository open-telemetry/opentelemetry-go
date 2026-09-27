// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

// configuratorMeter is the Meter returned by a MeterProvider built with a
// MeterConfigurator. Only such providers create it, so a provider without a
// configurator keeps returning the plain meter.
type configuratorMeter struct {
	*meter

	// gate holds the configurator's enabled decision for the meter's scope.
	gate *versionedEnabled

	// int64Wrappers and float64Wrappers return the same wrapper for the same
	// base instrument, so repeated instrument creation returns an identical
	// instrument, as it does without a configurator.
	int64Wrappers   cache[*int64Inst, *configuratorInt64Inst]
	float64Wrappers cache[*float64Inst, *configuratorFloat64Inst]
}

// Compile-time check configuratorMeter implements metric.Meter.
var _ metric.Meter = (*configuratorMeter)(nil)

func newConfiguratorMeter(m *meter) *configuratorMeter {
	// The gate is still the base meter's own enabled state, which the stable
	// instruments read; later steps move it onto configuratorMeter.
	return &configuratorMeter{meter: m, gate: &m.enabled}
}

// setEnabledIfNewer applies enabled to the gate if version is at least as new
// as the version already stored; see versionedEnabled.StoreIfNewer.
func (m *configuratorMeter) setEnabledIfNewer( // nolint:revive  // enabled is not a control flag.
	version uint64,
	enabled bool,
) {
	m.gate.StoreIfNewer(version, enabled)
}

// wrapInt64 returns the gated wrapper for inst. It reports false, and callers
// return inst unchanged, if inst isn't an *int64Inst.
func (m *configuratorMeter) wrapInt64(inst any) (*configuratorInt64Inst, bool) {
	base, ok := inst.(*int64Inst)
	if !ok {
		return nil, false
	}
	return m.int64Wrappers.Lookup(base, func() *configuratorInt64Inst {
		return &configuratorInt64Inst{int64Inst: base, gate: m.gate}
	}), true
}

// wrapFloat64 returns the gated wrapper for inst. It reports false, and
// callers return inst unchanged, if inst isn't a *float64Inst.
func (m *configuratorMeter) wrapFloat64(inst any) (*configuratorFloat64Inst, bool) {
	base, ok := inst.(*float64Inst)
	if !ok {
		return nil, false
	}
	return m.float64Wrappers.Lookup(base, func() *configuratorFloat64Inst {
		return &configuratorFloat64Inst{float64Inst: base, gate: m.gate}
	}), true
}

func (m *configuratorMeter) Int64Counter(
	name string,
	options ...metric.Int64CounterOption,
) (metric.Int64Counter, error) {
	i, err := m.meter.Int64Counter(name, options...)
	if w, ok := m.wrapInt64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Int64UpDownCounter(
	name string,
	options ...metric.Int64UpDownCounterOption,
) (metric.Int64UpDownCounter, error) {
	i, err := m.meter.Int64UpDownCounter(name, options...)
	if w, ok := m.wrapInt64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Int64Histogram(
	name string,
	options ...metric.Int64HistogramOption,
) (metric.Int64Histogram, error) {
	i, err := m.meter.Int64Histogram(name, options...)
	if w, ok := m.wrapInt64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	i, err := m.meter.Int64Gauge(name, options...)
	if w, ok := m.wrapInt64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Float64Counter(
	name string,
	options ...metric.Float64CounterOption,
) (metric.Float64Counter, error) {
	i, err := m.meter.Float64Counter(name, options...)
	if w, ok := m.wrapFloat64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Float64UpDownCounter(
	name string,
	options ...metric.Float64UpDownCounterOption,
) (metric.Float64UpDownCounter, error) {
	i, err := m.meter.Float64UpDownCounter(name, options...)
	if w, ok := m.wrapFloat64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Float64Histogram(
	name string,
	options ...metric.Float64HistogramOption,
) (metric.Float64Histogram, error) {
	i, err := m.meter.Float64Histogram(name, options...)
	if w, ok := m.wrapFloat64(i); ok {
		return w, err
	}
	return i, err
}

func (m *configuratorMeter) Float64Gauge(
	name string,
	options ...metric.Float64GaugeOption,
) (metric.Float64Gauge, error) {
	i, err := m.meter.Float64Gauge(name, options...)
	if w, ok := m.wrapFloat64(i); ok {
		return w, err
	}
	return i, err
}

// gateInt64Callbacks returns cbs with each callback skipped while the meter
// is disabled. The count is preserved, so the base meter's handling of
// repeated creation with callbacks is unchanged.
func (m *configuratorMeter) gateInt64Callbacks(cbs []metric.Int64Callback) []metric.Int64Callback {
	gated := make([]metric.Int64Callback, len(cbs))
	for i, cb := range cbs {
		gated[i] = func(ctx context.Context, o metric.Int64Observer) error {
			if !m.gate.Load() {
				return nil
			}
			return cb(ctx, o)
		}
	}
	return gated
}

// gateFloat64Callbacks returns cbs with each callback skipped while the meter
// is disabled. The count is preserved, so the base meter's handling of
// repeated creation with callbacks is unchanged.
func (m *configuratorMeter) gateFloat64Callbacks(cbs []metric.Float64Callback) []metric.Float64Callback {
	gated := make([]metric.Float64Callback, len(cbs))
	for i, cb := range cbs {
		gated[i] = func(ctx context.Context, o metric.Float64Observer) error {
			if !m.gate.Load() {
				return nil
			}
			return cb(ctx, o)
		}
	}
	return gated
}

// The observable constructors mirror meter's, but pass gated callbacks to the
// base meter. The instruments themselves come from the base meter's cache
// unwrapped, since only their callbacks need gating.

func (m *configuratorMeter) Int64ObservableCounter(
	name string,
	options ...metric.Int64ObservableCounterOption,
) (metric.Int64ObservableCounter, error) {
	cfg := metric.NewInt64ObservableCounterConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableCounter,
		Scope:       m.scope,
	}
	return m.int64ObservableInstrument(id, defaultAttributes(options), m.gateInt64Callbacks(cfg.Callbacks()))
}

func (m *configuratorMeter) Int64ObservableUpDownCounter(
	name string,
	options ...metric.Int64ObservableUpDownCounterOption,
) (metric.Int64ObservableUpDownCounter, error) {
	cfg := metric.NewInt64ObservableUpDownCounterConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableUpDownCounter,
		Scope:       m.scope,
	}
	return m.int64ObservableInstrument(id, defaultAttributes(options), m.gateInt64Callbacks(cfg.Callbacks()))
}

func (m *configuratorMeter) Int64ObservableGauge(
	name string,
	options ...metric.Int64ObservableGaugeOption,
) (metric.Int64ObservableGauge, error) {
	cfg := metric.NewInt64ObservableGaugeConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableGauge,
		Scope:       m.scope,
	}
	return m.int64ObservableInstrument(id, defaultAttributes(options), m.gateInt64Callbacks(cfg.Callbacks()))
}

func (m *configuratorMeter) Float64ObservableCounter(
	name string,
	options ...metric.Float64ObservableCounterOption,
) (metric.Float64ObservableCounter, error) {
	cfg := metric.NewFloat64ObservableCounterConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableCounter,
		Scope:       m.scope,
	}
	return m.float64ObservableInstrument(id, defaultAttributes(options), m.gateFloat64Callbacks(cfg.Callbacks()))
}

func (m *configuratorMeter) Float64ObservableUpDownCounter(
	name string,
	options ...metric.Float64ObservableUpDownCounterOption,
) (metric.Float64ObservableUpDownCounter, error) {
	cfg := metric.NewFloat64ObservableUpDownCounterConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableUpDownCounter,
		Scope:       m.scope,
	}
	return m.float64ObservableInstrument(id, defaultAttributes(options), m.gateFloat64Callbacks(cfg.Callbacks()))
}

func (m *configuratorMeter) Float64ObservableGauge(
	name string,
	options ...metric.Float64ObservableGaugeOption,
) (metric.Float64ObservableGauge, error) {
	cfg := metric.NewFloat64ObservableGaugeConfig(options...)
	id := Instrument{
		Name:        name,
		Description: cfg.Description(),
		Unit:        cfg.Unit(),
		Kind:        InstrumentKindObservableGauge,
		Scope:       m.scope,
	}
	return m.float64ObservableInstrument(id, defaultAttributes(options), m.gateFloat64Callbacks(cfg.Callbacks()))
}

// RegisterCallback registers f with the base meter, skipped while the meter is
// disabled.
func (m *configuratorMeter) RegisterCallback(
	f metric.Callback,
	insts ...metric.Observable,
) (metric.Registration, error) {
	return m.meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		if !m.gate.Load() {
			return nil
		}
		return f(ctx, o)
	}, insts...)
}
