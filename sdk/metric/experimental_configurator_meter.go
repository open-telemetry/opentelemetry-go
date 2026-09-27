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
