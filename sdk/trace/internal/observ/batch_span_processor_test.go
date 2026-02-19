// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace/internal/observ"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

const id = 0

func TestBSPComponentName(t *testing.T) {
	got := observ.BSPComponentName(42)
	want := semconv.OTelComponentName("batching_span_processor/42")
	assert.Equal(t, want, got)
}

func TestNewBSPDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY
	bsp, err := observ.NewBSP(id, nil, 0)
	assert.NoError(t, err)
	assert.Nil(t, bsp)
}

func TestNewBSPErrors(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	_, err := observ.NewBSP(id, nil, 0)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "create BSP queue capacity metric")
	assert.ErrorContains(t, err, "create BSP queue size metric")
	assert.ErrorContains(t, err, "register BSP queue size/capacity callback")
	assert.ErrorContains(t, err, "create BSP processed spans metric")
}

func bspSet(attrs ...attribute.KeyValue) attribute.Set {
	return attribute.NewSet(append([]attribute.KeyValue{
		semconv.OTelComponentTypeBatchingSpanProcessor,
		observ.BSPComponentName(id),
	}, attrs...)...)
}

func qCap(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorSpanQueueCapacity{}.Name(),
		Description: otelconv.SDKProcessorSpanQueueCapacity{}.Description(),
		Unit:        otelconv.SDKProcessorSpanQueueCapacity{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: bspSet(), Value: v},
			},
		},
	}
}

func qSize(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorSpanQueueSize{}.Name(),
		Description: otelconv.SDKProcessorSpanQueueSize{}.Description(),
		Unit:        otelconv.SDKProcessorSpanQueueSize{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: bspSet(), Value: v},
			},
		},
	}
}

func TestBSPCallback(t *testing.T) {
	collect := setup(t)

	var n int64 = 3
	bsp, err := observ.NewBSP(id, func() int64 { return n }, 5)
	require.NoError(t, err)
	require.NotNil(t, bsp)

	check(t, collect(), qSize(n), qCap(5))

	n = 4
	check(t, collect(), qSize(n), qCap(5))

	require.NoError(t, bsp.Shutdown())
	got := collect()
	assert.Empty(t, got.Metrics, "no metrics after shutdown")
}

func processed(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorSpanProcessed{}.Name(),
		Description: otelconv.SDKProcessorSpanProcessed{}.Description(),
		Unit:        otelconv.SDKProcessorSpanProcessed{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dPts,
		},
	}
}

func TestBSPProcessed(t *testing.T) {
	collect := setup(t)

	bsp, err := observ.NewBSP(id, nil, 0)
	require.NoError(t, err)
	require.NotNil(t, bsp)
	require.NoError(t, bsp.Shutdown()) // Unregister callback.

	ctx := t.Context()
	const p0 int64 = 10
	bsp.Processed(ctx, p0)
	const e0 int64 = 1
	bsp.ProcessedQueueFull(ctx, e0)
	check(t, collect(), processed(
		dPt(bspSet(), p0),
		dPt(bspSet(observ.ErrQueueFull), e0),
	))

	const p1 int64 = 20
	bsp.Processed(ctx, p1)
	const e1 int64 = 2
	bsp.ProcessedQueueFull(ctx, e1)
	check(t, collect(), processed(
		dPt(bspSet(), p0+p1),
		dPt(bspSet(observ.ErrQueueFull), e0+e1),
	))
}

func BenchmarkBSP(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	newBSP := func(b *testing.B) *observ.BSP {
		b.Helper()
		bsp, err := observ.NewBSP(id, func() int64 { return 3 }, 5)
		require.NoError(b, err)
		require.NotNil(b, bsp)
		b.Cleanup(func() {
			if err := bsp.Shutdown(); err != nil {
				b.Errorf("Shutdown: %v", err)
			}
		})
		return bsp
	}
	ctx := b.Context()

	b.Run("Processed", func(b *testing.B) {
		orig := otel.GetMeterProvider()
		b.Cleanup(func() { otel.SetMeterProvider(orig) })

		// Ensure deterministic benchmark by using noop meter.
		otel.SetMeterProvider(noop.NewMeterProvider())

		bsp := newBSP(b)

		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bsp.Processed(ctx, 10)
			}
		})
	})
	b.Run("ProcessedQueueFull", func(b *testing.B) {
		orig := otel.GetMeterProvider()
		b.Cleanup(func() { otel.SetMeterProvider(orig) })

		// Ensure deterministic benchmark by using noop meter.
		otel.SetMeterProvider(noop.NewMeterProvider())

		bsp := newBSP(b)

		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bsp.ProcessedQueueFull(ctx, 1)
			}
		})
	})
	b.Run("Callback", func(b *testing.B) {
		orig := otel.GetMeterProvider()
		b.Cleanup(func() { otel.SetMeterProvider(orig) })

		reader := metric.NewManualReader()
		mp := metric.NewMeterProvider(metric.WithReader(reader))
		otel.SetMeterProvider(mp)

		bsp := newBSP(b)
		var got metricdata.ResourceMetrics

		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			_ = reader.Collect(ctx, &got)
		}

		_ = got
		_ = bsp
	})
}
