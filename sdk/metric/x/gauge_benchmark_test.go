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

func BenchmarkInt64GaugeRecord(b *testing.B) {
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
			gauge, err := provider.Meter("bench").Int64Gauge("requests")
			if err != nil {
				b.Fatal(err)
			}
			recordOptions := []metric.RecordOption{metric.WithAttributeSet(attrs)}
			b.ReportAllocs()
			for b.Loop() {
				gauge.Record(b.Context(), 1, recordOptions...)
			}
		})
	}
}

func BenchmarkInt64GaugeFinish(b *testing.B) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	gauge, err := provider.Meter("bench").Int64Gauge("requests")
	if err != nil {
		b.Fatal(err)
	}
	finisher := gauge.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	}
	gauge.Record(b.Context(), 1, metric.WithAttributes(attrs...))

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

func BenchmarkFloat64GaugeRecord(b *testing.B) {
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
			gauge, err := provider.Meter("bench").Float64Gauge("requests")
			if err != nil {
				b.Fatal(err)
			}
			recordOptions := []metric.RecordOption{metric.WithAttributeSet(attrs)}
			b.ReportAllocs()
			for b.Loop() {
				gauge.Record(b.Context(), 1, recordOptions...)
			}
		})
	}
}

func BenchmarkFloat64GaugeFinish(b *testing.B) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetricx.WithFinish(),
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	gauge, err := provider.Meter("bench").Float64Gauge("requests")
	if err != nil {
		b.Fatal(err)
	}
	finisher := gauge.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	}
	gauge.Record(b.Context(), 1, metric.WithAttributes(attrs...))

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
