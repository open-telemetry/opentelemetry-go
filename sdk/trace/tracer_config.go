// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import "go.opentelemetry.io/otel/sdk/instrumentation"

// TracerConfig defines configurable aspects of a Tracer's behavior.
type TracerConfig struct {
	enabled bool
}

// Enabled reports whether the Tracer is enabled.
//
// If c is nil, true is returned. This is the default TracerConfig behavior.
func (c *TracerConfig) Enabled() bool {
	if c == nil {
		return true
	}
	return c.enabled
}

// TracerConfigOption applies an option to a TracerConfig.
type TracerConfigOption interface {
	apply(TracerConfig) TracerConfig
}

type tracerConfigOptionFunc func(TracerConfig) TracerConfig

func (fn tracerConfigOptionFunc) apply(cfg TracerConfig) TracerConfig {
	return fn(cfg)
}

// NewTracerConfig returns a new TracerConfig with opts applied.
//
// By default, Tracers are enabled.
func NewTracerConfig(opts ...TracerConfigOption) *TracerConfig {
	cfg := TracerConfig{enabled: true}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return &cfg
}

// WithTracerConfigEnabled returns a TracerConfigOption that sets the enabled
// state of a TracerConfig.
func WithTracerConfigEnabled(enabled bool) TracerConfigOption {
	return tracerConfigOptionFunc(func(cfg TracerConfig) TracerConfig {
		cfg.enabled = enabled
		return cfg
	})
}

// TracerConfigurator computes the TracerConfig for a Tracer.
//
// It is called when a Tracer is first created, and for each outstanding Tracer
// when a TracerProvider's TracerConfigurator is updated.
//
// Returning nil indicates the default TracerConfig should be used.
type TracerConfigurator func(scope instrumentation.Scope) *TracerConfig

type tracerConfiguratorOption interface {
	experimentalOption
	TracerConfigurator() TracerConfigurator
}
