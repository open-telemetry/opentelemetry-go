// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func ExampleWithFinish_upDownCounter() {
	ctx := context.Background()
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetricx.WithFinish())
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Print(err)
			return
		}
	}()
	// Create instruments directly from the opted-in provider. Global proxy
	// instruments created before provider registration cannot acquire Finisher.
	counter, err := provider.Meter("example").Int64UpDownCounter("active.connections")
	if err != nil {
		log.Print(err)
		return
	}
	attrs := []attribute.KeyValue{attribute.String("server", "one")}
	option := metric.WithAttributes(attrs...)
	counter.Add(ctx, 3, option)
	counter.Add(ctx, -1, option)
	if finisher, ok := counter.(metricx.Finisher); ok {
		finisher.Finish(ctx, attrs...)
	}

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		log.Print(err)
		return
	}
	sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	fmt.Println("Final value:", sum.DataPoints[0].Value)
	if err := reader.Collect(ctx, &rm); err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Remaining scopes:", len(rm.ScopeMetrics))
	// Output:
	// Final value: 2
	// Remaining scopes: 0
}
