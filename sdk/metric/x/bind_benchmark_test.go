// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func BenchmarkInt64CounterBind(b *testing.B) {
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	attrSet := attribute.NewSet(attrs...)
	addOpts := []metric.AddOption{metric.WithAttributeSet(attrSet)}

	temporalities := []struct {
		name     string
		selector sdkmetric.TemporalitySelector
	}{
		{
			name: "Cumulative",
			selector: func(sdkmetric.InstrumentKind) metricdata.Temporality {
				return metricdata.CumulativeTemporality
			},
		},
		{
			name: "Delta",
			selector: func(sdkmetric.InstrumentKind) metricdata.Temporality {
				return metricdata.DeltaTemporality
			},
		},
	}

	for _, temp := range temporalities {
		b.Run(temp.name, func(b *testing.B) {
			b.Run("stable", func(b *testing.B) {
				reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(temp.selector))
				provider := sdkmetric.NewMeterProvider(
					sdkmetric.WithReader(reader),
					sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
				)
				b.Cleanup(func() { _ = provider.Shutdown(context.WithoutCancel(b.Context())) })
				counter, err := provider.Meter("bench").Int64Counter("requests")
				if err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					counter.Add(b.Context(), 1, addOpts...)
				}
			})

			b.Run("unbound", func(b *testing.B) {
				reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(temp.selector))
				provider := sdkmetric.NewMeterProvider(
					sdkmetricx.WithBinding(),
					sdkmetric.WithReader(reader),
					sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
				)
				b.Cleanup(func() { _ = provider.Shutdown(context.WithoutCancel(b.Context())) })
				counter, err := provider.Meter("bench").Int64Counter("requests")
				if err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					counter.Add(b.Context(), 1, addOpts...)
				}
			})

			b.Run("bound", func(b *testing.B) {
				reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(temp.selector))
				provider := sdkmetric.NewMeterProvider(
					sdkmetricx.WithBinding(),
					sdkmetric.WithReader(reader),
					sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
				)
				b.Cleanup(func() { _ = provider.Shutdown(context.WithoutCancel(b.Context())) })
				counter, err := provider.Meter("bench").Int64Counter("requests")
				if err != nil {
					b.Fatal(err)
				}
				bound := counter.(metricx.Int64CounterBinder).Bind(attrs...)
				b.ReportAllocs()
				for b.Loop() {
					bound.Add(b.Context(), 1)
				}
			})
		})
	}
}
