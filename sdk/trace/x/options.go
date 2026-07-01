// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x contains experimental trace SDK options.
package x // import "go.opentelemetry.io/otel/sdk/trace/x"

import "go.opentelemetry.io/otel/sdk/trace"

type tracerConfiguratorOption struct {
	trace.TracerProviderOption
	configurator trace.TracerConfigurator
}

// Experimental prevents the API from panicking when the option is used.
func (tracerConfiguratorOption) Experimental() {}

func (o tracerConfiguratorOption) TracerConfigurator() trace.TracerConfigurator {
	return o.configurator
}

// WithTracerConfigurator returns a TracerProviderOption that configures the
// TracerConfigurator used when creating Tracers.
//
// This is an experimental feature.
func WithTracerConfigurator(cfg trace.TracerConfigurator) trace.TracerProviderOption {
	return tracerConfiguratorOption{
		configurator: cfg,
	}
}
