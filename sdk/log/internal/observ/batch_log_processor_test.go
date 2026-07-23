// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/log/internal/observ"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/otelconv"
)

const id = 0

func TestBLPComponentName(t *testing.T) {
	got := observ.BLPComponentName(42)
	want := semconv.OTelComponentName("batching_log_processor/42")
	assert.Equal(t, want, got)
}

func TestNewBLPDisabled(t *testing.T) {
	blp, err := observ.NewBLP(id, nil, 0)
	assert.NoError(t, err)
	assert.Nil(t, blp)
}

func TestNewBLPNilQLen(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	blp, err := observ.NewBLP(id, nil, 0)
	assert.Nil(t, blp)
	assert.ErrorContains(t, err, "qLen must not be nil")
}

func TestNewBLPErrors(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	check := func(t *testing.T, wantMsg string) {
		t.Helper()
		_, err := observ.NewBLP(id, func() int64 { return 0 }, 0)
		require.ErrorIs(t, err, assert.AnError)
		assert.ErrorContains(t, err, wantMsg)
	}

	t.Run("qCap", func(t *testing.T) {
		otel.SetMeterProvider(meterProvider{m: &errOnNthObsCounterMeter{n: 1, err: assert.AnError}})
		check(t, "create BLP queue capacity metric")
	})
	t.Run("qSize", func(t *testing.T) {
		otel.SetMeterProvider(meterProvider{m: &errOnNthObsCounterMeter{n: 2, err: assert.AnError}})
		check(t, "create BLP queue size metric")
	})
	t.Run("callback", func(t *testing.T) {
		otel.SetMeterProvider(meterProvider{m: &errCallbackMeter{err: assert.AnError}})
		check(t, "register BLP queue size/capacity callback")
	})
	t.Run("processed", func(t *testing.T) {
		otel.SetMeterProvider(meterProvider{m: &errCounterMeter{err: assert.AnError}})
		check(t, "create BLP processed logs metric")
	})
}

func blpSet(attrs ...attribute.KeyValue) attribute.Set {
	return attribute.NewSet(append([]attribute.KeyValue{
		semconv.OTelComponentTypeBatchingLogProcessor,
		observ.BLPComponentName(id),
	}, attrs...)...)
}

func qCap(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogQueueCapacity{}.Name(),
		Description: otelconv.SDKProcessorLogQueueCapacity{}.Description(),
		Unit:        otelconv.SDKProcessorLogQueueCapacity{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: blpSet(), Value: v},
			},
		},
	}
}

func qSize(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogQueueSize{}.Name(),
		Description: otelconv.SDKProcessorLogQueueSize{}.Description(),
		Unit:        otelconv.SDKProcessorLogQueueSize{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: blpSet(), Value: v},
			},
		},
	}
}

func TestBLPCallback(t *testing.T) {
	collect := setup(t)

	var n int64 = 3
	blp, err := observ.NewBLP(id, func() int64 { return n }, 5)
	require.NoError(t, err)
	require.NotNil(t, blp)

	check(t, collect(), qSize(n), qCap(5))

	n = 4
	check(t, collect(), qSize(n), qCap(5))

	require.NoError(t, blp.Shutdown())
	got := collect()
	assert.Empty(t, got.Metrics, "no metrics after shutdown")
}

func processed(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogProcessed{}.Name(),
		Description: otelconv.SDKProcessorLogProcessed{}.Description(),
		Unit:        otelconv.SDKProcessorLogProcessed{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dPts,
		},
	}
}

func TestBLPProcessed(t *testing.T) {
	collect := setup(t)

	blp, err := observ.NewBLP(id, func() int64 { return 0 }, 0)
	require.NoError(t, err)
	require.NotNil(t, blp)
	require.NoError(t, blp.Shutdown()) // Unregister callback.

	ctx := t.Context()
	const p0 int64 = 10
	blp.Processed(ctx, p0)
	const e0 int64 = 1
	blp.ProcessedQueueFull(ctx, e0)
	check(t, collect(), processed(
		dPt(blpSet(), p0),
		dPt(blpSet(observ.ErrQueueFull), e0),
	))

	const p1 int64 = 20
	blp.Processed(ctx, p1)
	const e1 int64 = 2
	blp.ProcessedQueueFull(ctx, e1)
	check(t, collect(), processed(
		dPt(blpSet(), p0+p1),
		dPt(blpSet(observ.ErrQueueFull), e0+e1),
	))
}

func BenchmarkBLP(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	b.Cleanup(func() { otel.SetMeterProvider(orig) })

	newBLP := func(b *testing.B) *observ.BLP {
		b.Helper()
		blp, err := observ.NewBLP(id, func() int64 { return 3 }, 5)
		require.NoError(b, err)
		require.NotNil(b, blp)
		b.Cleanup(func() {
			if err := blp.Shutdown(); err != nil {
				b.Errorf("Shutdown: %v", err)
			}
		})
		return blp
	}
	ctx := b.Context()

	for _, tt := range []struct {
		name string
		fn   func(*observ.BLP, context.Context)
	}{
		{"Processed", func(blp *observ.BLP, ctx context.Context) { blp.Processed(ctx, 10) }},
		{"ProcessedQueueFull", func(blp *observ.BLP, ctx context.Context) { blp.ProcessedQueueFull(ctx, 1) }},
	} {
		b.Run(tt.name, func(b *testing.B) {
			otel.SetMeterProvider(noop.NewMeterProvider())
			blp := newBLP(b)
			b.ResetTimer()
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					tt.fn(blp, ctx)
				}
			})
		})
	}

	b.Run("Callback", func(b *testing.B) {
		reader := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(reader))
		otel.SetMeterProvider(mp)

		blp := newBLP(b)
		var got metricdata.ResourceMetrics

		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			_ = reader.Collect(ctx, &got)
		}
		_ = blp
	})
}
