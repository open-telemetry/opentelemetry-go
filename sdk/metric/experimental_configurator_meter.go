// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import "go.opentelemetry.io/otel/metric"

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
