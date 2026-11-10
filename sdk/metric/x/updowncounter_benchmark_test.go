// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"slices"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func BenchmarkInt64UpDownCounterRecord(b *testing.B) {
	attrs := attribute.NewSet(attribute.String("route", "/orders"))
	for _, bench := range []struct {
		name    string
		options []sdkmetric.Option
	}{
		{name: "stable"},
		{name: "finish", options: []sdkmetric.Option{sdkmetricx.WithFinish()}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			reader := sdkmetric.NewManualReader()
			options := append(
				slices.Clone(bench.options),
				sdkmetric.WithReader(reader),
				sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
			)
			provider := sdkmetric.NewMeterProvider(options...)
			b.Cleanup(func() { _ = provider.Shutdown(context.WithoutCancel(b.Context())) })
			counter, err := provider.Meter("bench").Int64UpDownCounter("requests")
			if err != nil {
				b.Fatal(err)
			}
			addOptions := []metric.AddOption{metric.WithAttributeSet(attrs)}
			b.ReportAllocs()
			for b.Loop() {
				counter.Add(b.Context(), 1, addOptions...)
			}
		})
	}
}

func BenchmarkInt64UpDownCounterFinish(b *testing.B) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("bench").Int64UpDownCounter("requests")
	if err != nil {
		b.Fatal(err)
	}
	finisher := counter.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	}
	counter.Add(b.Context(), 1, metric.WithAttributes(attrs...))

	b.Run("reused-slice", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			finisher.Finish(b.Context(), attrs...)
		}
	})
	b.Run("inline", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			finisher.Finish(
				b.Context(),
				attribute.String("zone", "east"),
				attribute.String("route", "/orders"),
			)
		}
	})
}

func BenchmarkFloat64UpDownCounterRecord(b *testing.B) {
	attrs := attribute.NewSet(attribute.String("route", "/orders"))
	for _, bench := range []struct {
		name    string
		options []sdkmetric.Option
	}{
		{name: "stable"},
		{name: "finish", options: []sdkmetric.Option{sdkmetricx.WithFinish()}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			reader := sdkmetric.NewManualReader()
			options := append(
				slices.Clone(bench.options),
				sdkmetric.WithReader(reader),
				sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
			)
			provider := sdkmetric.NewMeterProvider(options...)
			b.Cleanup(func() { _ = provider.Shutdown(context.WithoutCancel(b.Context())) })
			counter, err := provider.Meter("bench").Float64UpDownCounter("requests")
			if err != nil {
				b.Fatal(err)
			}
			addOptions := []metric.AddOption{metric.WithAttributeSet(attrs)}
			b.ReportAllocs()
			for b.Loop() {
				counter.Add(b.Context(), 1, addOptions...)
			}
		})
	}
}

func BenchmarkFloat64UpDownCounterFinish(b *testing.B) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("bench").Float64UpDownCounter("requests")
	if err != nil {
		b.Fatal(err)
	}
	finisher := counter.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	}
	counter.Add(b.Context(), 1, metric.WithAttributes(attrs...))

	b.Run("reused-slice", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			finisher.Finish(b.Context(), attrs...)
		}
	})
	b.Run("inline", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			finisher.Finish(
				b.Context(),
				attribute.String("zone", "east"),
				attribute.String("route", "/orders"),
			)
		}
	})
}
