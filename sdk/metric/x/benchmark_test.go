// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func BenchmarkInt64Counter(b *testing.B) {
	ctx := b.Context()
	attrs := attribute.NewSet(attribute.String("route", "/orders"))
	b.Run("stable/unbound", func(b *testing.B) {
		reader := sdkmetric.NewManualReader()
		provider := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(reader),
			sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
		)
		counter, _ := provider.Meter("bench").Int64Counter("requests")
		option := metric.WithAttributeSet(attrs)
		b.ReportAllocs()
		for b.Loop() {
			counter.Add(ctx, 1, option)
		}
	})
	b.Run("experimental/unbound", func(b *testing.B) {
		reader := sdkmetricx.NewManualReader()
		provider := sdkmetricx.NewMeterProvider(
			sdkmetricx.WithReader(reader),
			sdkmetricx.WithExemplarFilter(exemplar.AlwaysOffFilter),
		)
		counter, _ := provider.Meter("bench").Int64Counter("requests")
		option := metric.WithAttributeSet(attrs)
		b.ReportAllocs()
		for b.Loop() {
			counter.Add(ctx, 1, option)
		}
	})
	b.Run("experimental/bound", func(b *testing.B) {
		reader := sdkmetricx.NewManualReader()
		provider := sdkmetricx.NewMeterProvider(
			sdkmetricx.WithReader(reader),
			sdkmetricx.WithExemplarFilter(exemplar.AlwaysOffFilter),
		)
		counter, _ := provider.Meter("bench").Int64Counter("requests")
		bound := counter.(metricx.Int64CounterBinder).Bind(attribute.String("route", "/orders"))
		b.ReportAllocs()
		for b.Loop() {
			bound.Add(ctx, 1)
		}
	})
}

func BenchmarkBoundInt64CounterTemporalities(b *testing.B) {
	for _, benchmark := range []struct {
		name     string
		selector sdkmetricx.TemporalitySelector
	}{
		{name: "cumulative", selector: sdkmetricx.CumulativeTemporalitySelector},
		{name: "delta", selector: sdkmetricx.DeltaTemporalitySelector},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			reader := sdkmetricx.NewManualReader(
				sdkmetricx.WithTemporalitySelector(benchmark.selector),
			)
			provider := sdkmetricx.NewMeterProvider(
				sdkmetricx.WithReader(reader),
				sdkmetricx.WithExemplarFilter(exemplar.AlwaysOffFilter),
			)
			counter, _ := provider.Meter("bench").Int64Counter("requests")
			bound := counter.(metricx.Int64CounterBinder).Bind(
				attribute.String("route", "/orders"),
			)
			b.ReportAllocs()
			for b.Loop() {
				bound.Add(b.Context(), 1)
			}
		})
	}
}
