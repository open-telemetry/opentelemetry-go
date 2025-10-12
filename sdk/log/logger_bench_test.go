// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func BenchmarkLoggerEmit(b *testing.B) {
	logger := newTestLogger(b)

	r := log.Record{}
	r.SetTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetObservedTimestamp(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC))
	r.SetBody(log.StringValue("testing body value"))
	r.SetSeverity(log.SeverityInfo)
	r.SetSeverityText("testing text")

	r.AddAttributes(
		log.String("k1", "str"),
		log.Float64("k2", 1.0),
		log.Int("k3", 2),
		log.Bool("k4", true),
		log.Bytes("k5", []byte{1}),
	)

	r10 := r
	r10.AddAttributes(
		log.String("k6", "str"),
		log.Float64("k7", 1.0),
		log.Int("k8", 2),
		log.Bool("k9", true),
		log.Bytes("k10", []byte{1}),
	)

	require.Equal(b, 5, r.AttributesLen())
	require.Equal(b, 10, r10.AttributesLen())

	b.Run("5 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Emit(b.Context(), r)
			}
		})
	})

	b.Run("10 attributes", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Emit(b.Context(), r10)
			}
		})
	})
}

func BenchmarkLoggerEmitObservability(b *testing.B) {
	r := log.Record{}

	orig := otel.GetMeterProvider()
	b.Cleanup(func() { otel.SetMeterProvider(orig) })
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	run := func(logger *logger) func(b *testing.B) {
		return func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					logger.Emit(b.Context(), r)
				}
			})
		}
	}

	lp := NewLoggerProvider()
	scope := instrumentation.Scope{}

	b.Run("Disabled", run(newLogger(lp, scope)))

	b.Run("Enabled", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

		run(newLogger(lp, scope))(b)
	})

	var rm metricdata.ResourceMetrics
	err := reader.Collect(b.Context(), &rm)
	require.NoError(b, err)
	require.Len(b, rm.ScopeMetrics, 1)
}

func BenchmarkLoggerEnabled(b *testing.B) {
	logger := newTestLogger(b)
	ctx := b.Context()
	param := log.EnabledParameters{Severity: log.SeverityDebug}
	var enabled bool

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		enabled = logger.Enabled(ctx, param)
	}

	_ = enabled
}

func newTestLogger(t testing.TB) log.Logger {
	provider := NewLoggerProvider(
		WithProcessor(newFltrProcessor("0", false)),
		WithProcessor(newFltrProcessor("1", true)),
	)
	return provider.Logger(t.Name())
}
