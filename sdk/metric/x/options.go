// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type meterConfiguratorProviderOption struct {
	// nil embed, skip guard in newConfig prevents apply from being called
	sdkmetric.Option
	configurator MeterConfigurator
}

// Experimental marks this as an experimental option so the skip guard
// in newConfig() skips calling the nil embedded apply().
func (meterConfiguratorProviderOption) Experimental() {}

// WithMeterConfigurator returns an [sdkmetric.Option] that sets the
// MeterConfigurator on a MeterProvider.
func WithMeterConfigurator(fn MeterConfigurator) sdkmetric.Option {
	return meterConfiguratorProviderOption{configurator: fn}
}

// MeterConfigurator returns the configurator as func(scope) any so
// sdk/metric can extract it via duck-type without importing this package.
func (co meterConfiguratorProviderOption) MeterConfigurator() func(instrumentation.Scope) any {
	return func(s instrumentation.Scope) any {
		return co.configurator(s)
	}
}

// MeterConfiguratorUpdater is implemented by MeterProviders that support
// runtime configurator updates. Type-assert a MeterProvider to this interface
// to update its configurator after construction
type MeterConfiguratorUpdater interface {
	SetMeterConfigurator(func(instrumentation.Scope) any)
}

// SetMeterConfigurator updates the MeterConfigurator on MeterProvider which support that operation
func SetMeterConfigurator(meterProvider MeterConfiguratorUpdater, fn MeterConfigurator) {
	meterProvider.SetMeterConfigurator(func(s instrumentation.Scope) any {
		return fn(s)
	})
}
