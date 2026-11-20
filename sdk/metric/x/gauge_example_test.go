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

func ExampleWithFinish_gauge() {
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
	gauge, err := provider.Meter("example").Int64Gauge("queue.depth")
	if err != nil {
		log.Print(err)
		return
	}
	attrs := []attribute.KeyValue{attribute.String("server", "one")}
	option := metric.WithAttributes(attrs...)
	gauge.Record(ctx, 3, option)
	gauge.Record(ctx, 1, option)
	if finisher, ok := gauge.(metricx.Finisher); ok {
		finisher.Finish(ctx, attrs...)
	}

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &rm); err != nil {
		log.Print(err)
		return
	}
	data := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[int64])
	fmt.Println("Final value:", data.DataPoints[0].Value)
	if err := reader.Collect(ctx, &rm); err != nil {
		log.Print(err)
		return
	}
	fmt.Println("Remaining scopes:", len(rm.ScopeMetrics))
	// Output:
	// Final value: 1
	// Remaining scopes: 0
}
