// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus

import (
	"context"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric"
)

func benchmarkCollect(b *testing.B, n int) {
	ctx := context.Background()
	registry := prometheus.NewRegistry()
	exporter, err := New(WithRegisterer(registry))
	require.NoError(b, err)
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("testmeter")

	for i := 0; i < n; i++ {
		counter, err := meter.Float64Counter(fmt.Sprintf("foo_%d", i))
		require.NoError(b, err)
		counter.Add(ctx, float64(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Gather()
		require.NoError(b, err)
	}
}

func BenchmarkCollect1(b *testing.B)     { benchmarkCollect(b, 1) }
func BenchmarkCollect10(b *testing.B)    { benchmarkCollect(b, 10) }
func BenchmarkCollect100(b *testing.B)   { benchmarkCollect(b, 100) }
func BenchmarkCollect1000(b *testing.B)  { benchmarkCollect(b, 1000) }
func BenchmarkCollect10000(b *testing.B) { benchmarkCollect(b, 10000) }
