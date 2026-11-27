// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func ExampleWithBinding() {
	ctx := context.Background()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetricx.WithBinding(),
	)
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Print(err)
		}
	}()

	// Create instruments directly from the opted-in provider. Global proxy
	// instruments created before provider registration cannot acquire Binder.
	counter, err := provider.Meter("example").Int64Counter("http.server.requests")
	if err != nil {
		log.Print(err)
		return
	}
	if binder, ok := counter.(metricx.Int64CounterBinder); ok {
		bound := binder.Bind(attribute.String("http.route", "/orders"))
		bound.Add(ctx, 1)
		bound.Add(ctx, 2)
	}

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		log.Print(err)
		return
	}
	sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	fmt.Println("Count:", sum.DataPoints[0].Value)
	// Output:
	// Count: 3
}
