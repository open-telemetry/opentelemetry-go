// Copyright (c) 2024 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package benchmark_test

import (
	"context"
	"testing"

	promsdk "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func BenchmarkPrometheusCounter(b *testing.B) {
	opts := promsdk.CounterOpts{
		Name: "test_counter",
		Help: "help",
	}
	cv := promsdk.NewCounterVec(opts, []string{"tag1"})
	counter := cv.WithLabelValues("value1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Add(1)
	}
}

func BenchmarkOTELCounter(b *testing.B) {
	counter := otelCounter(b)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Add(ctx, 1)
	}
}

func BenchmarkOTELCounterWithLabel(b *testing.B) {
	counter := otelCounter(b)
	ctx := context.Background()
	attrSet := attribute.NewSet(attribute.String("tag1", "value1"))
	attrOpt := metric.WithAttributeSet(attrSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Add(ctx, 1, attrOpt)
	}
}

func otelCounter(b *testing.B) metric.Int64Counter {
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewManualReader()),
	)

	meter := meterProvider.Meter("test")
	counter, err := meter.Int64Counter("test_counter")
	require.NoError(b, err)
	return counter
}

func BenchmarkNoOpOTELCounter(b *testing.B) {
	meter := noop.NewMeterProvider().Meter("test")
	counter, err := meter.Int64Counter("test_counter")
	require.NoError(b, err)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Add(ctx, 1)
	}
}

func BenchmarkNoOpOTELCounterWithLabel(b *testing.B) {
	meter := noop.NewMeterProvider().Meter("test")
	counter, err := meter.Int64Counter("test_counter")
	require.NoError(b, err)
	ctx := context.Background()
	attrSet := attribute.NewSet(attribute.String("tag1", "value1"))
	attrOpt := metric.WithAttributeSet(attrSet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Add(ctx, 1, attrOpt)
	}
}
