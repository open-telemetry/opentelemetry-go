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

// MeterConfigurator is called by a MeterProvider when a Meter is created.
// It receives the instrumentation scope and returns the runtime configuration
// for that Meter.
type MeterConfigurator func(instrumentation.Scope) MeterConfig
