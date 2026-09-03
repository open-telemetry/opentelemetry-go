// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import "go.opentelemetry.io/otel/sdk/instrumentation"

// MeterConfig contains SDK runtime configuration for a Meter.
// It is returned by [MeterConfigurator] and controls whether a Meter
// records measurements.
//
// The zero value is a valid configuration with the Meter enabled.
type MeterConfig struct {
	disabled bool
}

// Enabled reports whether the Meter is enabled.
func (c MeterConfig) Enabled() bool {
	return !c.disabled
}

// NewMeterConfig returns a MeterConfig with opts applied.
func NewMeterConfig(opts ...MeterConfigOption) MeterConfig {
	var c MeterConfig
	for _, opt := range opts {
		opt.applyMeterConfig(&c)
	}
	return c
}

// MeterConfigOption is an option for [MeterConfig].
type MeterConfigOption interface {
	applyMeterConfig(*MeterConfig)
}

// WithMeterEnabled sets whether the Meter is enabled.
// A disabled Meter does not record any measurements.
func WithMeterEnabled(enabled bool) MeterConfigOption {
	return meterEnabledOption(enabled)
}

type meterEnabledOption bool

func (o meterEnabledOption) applyMeterConfig(c *MeterConfig) {
	c.disabled = !bool(o)
}

// MeterConfigurator is called by a MeterProvider when a Meter is created, and
// again for every existing Meter each time [MeterConfiguratorHandle.Set] is
// called. It receives the instrumentation scope and returns the runtime
// configuration for that Meter.
//
// Implementations must return quickly and must not block: a slow or blocked
// call can stall Meter creation and configuration updates for the whole
// MeterProvider, not just the caller that triggered it. Implementations must
// also not panic; a panic is not recovered and propagates to the caller of
// Meter or [MeterConfiguratorHandle.Set].
//
// A MeterConfigurator may be called more than once for the same
// instrumentation scope over the lifetime of a MeterProvider, and should be
// efficient and idempotent to call repeatedly for the same scope.
type MeterConfigurator func(instrumentation.Scope) MeterConfig
