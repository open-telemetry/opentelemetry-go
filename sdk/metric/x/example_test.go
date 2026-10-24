// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	x "go.opentelemetry.io/otel/sdk/metric/x"
)

// ExampleMeterConfiguratorHandle demonstrates how to configure and update a
// MeterConfigurator at runtime. The handle is created before the MeterProvider
// and passed via WithMeterConfigurator; subsequent calls to Set are reflected
// immediately across all existing meters via a synchronous cache walk.
func ExampleMeterConfiguratorHandle_Set() {
	handle := x.NewMeterConfiguratorHandle()

	mp := sdkmetric.NewMeterProvider(
		x.WithMeterConfigurator(handle),
	)
	defer mp.Shutdown(context.Background()) //nolint:errcheck

	// Initial configuration, set before any Meter is created.
	handle.Set(func(s instrumentation.Scope) x.MeterConfig {
		if s.Name == "com.example.chatty-library" {
			return x.NewMeterConfig(x.WithMeterEnabled(false))
		}
		return x.MeterConfig{}
	})

	// Runtime update; triggers a synchronous cache walk on the MeterProvider.
	handle.Set(func(s instrumentation.Scope) x.MeterConfig {
		if s.Name == "com.example.other-library" {
			return x.NewMeterConfig(x.WithMeterEnabled(false))
		}
		return x.MeterConfig{}
	})
}
