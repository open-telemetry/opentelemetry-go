// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	benchMethods = []string{"GET", "POST", "PUT", "DELETE"}
	benchRoutes  = []string{"/api/v1/users", "/api/v1/items"}
)

func benchmarkAttrSets() []otelmetric.MeasurementOption {
	opts := make([]otelmetric.MeasurementOption, 0, 5)
	for s := range 5 {
		opts = append(opts, otelmetric.WithAttributes(
			attribute.String("http.request.method", benchMethods[s%len(benchMethods)]),
			attribute.String("http.route", benchRoutes[s%len(benchRoutes)]),
			attribute.Int("http.response.status_code", 200+(s%3)*100),
			attribute.String("network.protocol.name", "http"),
			attribute.String("server.address", "localhost"),
		))
	}
	return opts
}

func run(n int) func(b *testing.B) {
	return func(b *testing.B) {
		ctx := b.Context()
		registry := prometheus.NewRegistry()
		exporter, err := New(WithRegisterer(registry))
		require.NoError(b, err)
		provider := metric.NewMeterProvider(
			metric.WithReader(exporter),
			metric.WithView(
				metric.NewView(
					metric.Instrument{Name: "http.client.request.duration.*"},
					metric.Stream{Aggregation: metric.AggregationBase2ExponentialHistogram{
						MaxSize:  160,
						MaxScale: 20,
					}},
				),
			),
		)
		meter := provider.Meter(
			"http.server",
			otelmetric.WithInstrumentationVersion("v1.0.0"),
			otelmetric.WithInstrumentationAttributes(attribute.String("scope.env", "prod")),
		)

		attrOpts := benchmarkAttrSets()
		for i := range n {
			switch i % 4 {
			case 0:
				counter, err := meter.Float64Counter(
					fmt.Sprintf("http.server.request.count.%d", i),
					otelmetric.WithDescription("Total HTTP requests"),
					otelmetric.WithUnit("{request}"),
				)
				require.NoError(b, err)
				for _, opt := range attrOpts {
					counter.Add(ctx, float64(i+1), opt)
				}
			case 1:
				gauge, err := meter.Float64Gauge(
					fmt.Sprintf("http.server.active_requests.%d", i),
					otelmetric.WithDescription("Active HTTP requests"),
					otelmetric.WithUnit("{request}"),
				)
				require.NoError(b, err)
				for _, opt := range attrOpts {
					gauge.Record(ctx, float64(i+1), opt)
				}
			case 2:
				hist, err := meter.Float64Histogram(
					fmt.Sprintf("http.server.request.duration.%d", i),
					otelmetric.WithDescription("HTTP server request duration"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(b, err)
				for s, opt := range attrOpts {
					hist.Record(ctx, float64((s+1)*10), opt)
				}
			case 3:
				expHist, err := meter.Float64Histogram(
					fmt.Sprintf("http.client.request.duration.%d", i),
					otelmetric.WithDescription("HTTP client request duration"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(b, err)
				for s, opt := range attrOpts {
					expHist.Record(ctx, float64((s+1)*10), opt)
				}
			}
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := registry.Gather()
			require.NoError(b, err)
		}
	}
}

func benchmarkCollect(b *testing.B, n int) {
	b.Run("ObservabilityDisabled", run(n))
	b.Run("ObservabilityEnabled", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		bmark := run(n)
		bmark(b)
	})
}

func BenchmarkCollect1(b *testing.B)    { benchmarkCollect(b, 1) }
func BenchmarkCollect10(b *testing.B)   { benchmarkCollect(b, 10) }
func BenchmarkCollect100(b *testing.B)  { benchmarkCollect(b, 100) }
func BenchmarkCollect1000(b *testing.B) { benchmarkCollect(b, 1000) }
