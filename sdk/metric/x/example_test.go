// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	x "go.opentelemetry.io/otel/sdk/metric/x"
)

// ExampleSetMeterConfigurator demonstrates how to update the MeterConfigurator
// on a MeterProvider at runtime. The type assertion against [x.MeterConfiguratorUpdater]
// is intentional: it checks whether the provider supports runtime configurator updates
// without coupling the caller to the SDK implementation.
func ExampleSetMeterConfigurator() {
	mp := sdkmetric.NewMeterProvider()
	defer mp.Shutdown(context.Background()) //nolint:errcheck

	// Register as the global provider. In real applications the provider is
	// often obtained via otel.GetMeterProvider, hiding the concrete type.
	otel.SetMeterProvider(mp)

	// Obtain the provider through the global, concrete type is opaque here.
	p := otel.GetMeterProvider()

	// Type-assert to check if the provider supports runtime configurator updates.
	// This succeeds for *sdkmetric.MeterProvider but not for noop or other implementations.
	u, ok := p.(x.MeterConfiguratorUpdater)
	if !ok {
		// Provider does not support MeterConfiguratorUpdater; skip configuration.
		return
	}

	x.SetMeterConfigurator(u, func(s instrumentation.Scope) x.MeterConfig {
		// Disable meters from a known chatty third-party library.
		if s.Name == "com.example.chatty-library" {
			return x.NewMeterConfig(x.WithMeterEnabled(false))
		}
		return x.MeterConfig{}
	})
}
