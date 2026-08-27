// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func Example_boundCounter() {
	reader := sdkmetricx.NewManualReader()
	provider := sdkmetricx.NewMeterProvider(sdkmetricx.WithReader(reader))
	counter, _ := provider.Meter("example").Int64Counter("requests")

	binder := counter.(metricx.Int64CounterBinder)
	bound := binder.Bind(attribute.String("route", "/orders"))
	bound.Add(context.Background(), 1)

	finisher := counter.(metricx.Finisher)
	finisher.Finish(context.Background(), attribute.String("route", "/orders"))
	// A later Add through bound lazily starts a new series lifetime.
	bound.Add(context.Background(), 1)
}
