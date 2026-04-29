// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestBoundInstrumentFloat64(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("K", "V")}
	set := attribute.NewSet(attrs...)

	// Test bound instrument (Cumulative)
	t.Run("Bound/Cumulative", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Float64Binder)
		require.True(t, ok, "counter does not implement x.Float64Binder")

		bound := binder.Bind(attrs...)
		bound.Add(t.Context(), 1)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		m := rm.ScopeMetrics[0].Metrics[0]
		assert.Equal(t, "test.counter", m.Name)

		sum, ok := m.Data.(metricdata.Sum[float64])
		require.True(t, ok)
		require.Len(t, sum.DataPoints, 1)
		dp := sum.DataPoints[0]
		assert.Equal(t, float64(1), dp.Value)
		assert.Equal(t, set, dp.Attributes)
	})

	// Test bound instrument (Delta)
	t.Run("Bound/Delta", func(t *testing.T) {
		r := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}))
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Float64Binder)
		require.True(t, ok, "counter does not implement x.Float64Binder")

		bound := binder.Bind(attrs...)
		bound.Add(t.Context(), 1)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		assert.Equal(t, float64(1), sum.DataPoints[0].Value)

		// Record again on the bound instrument!
		bound.Add(t.Context(), 2)

		// Collect again. The value should be 2 (Delta!), not 3 (Cumulative!).
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		assert.Equal(t, float64(2), sum.DataPoints[0].Value)
	})
}

func TestBoundInstrumentInt64(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("K", "V")}
	set := attribute.NewSet(attrs...)

	// Test bound instrument (Cumulative)
	t.Run("Bound/Cumulative", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Int64Binder)
		require.True(t, ok, "counter does not implement x.Int64Binder")

		bound := binder.Bind(attrs...)
		bound.Add(t.Context(), 1)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		m := rm.ScopeMetrics[0].Metrics[0]
		assert.Equal(t, "test.counter", m.Name)

		sum, ok := m.Data.(metricdata.Sum[int64])
		require.True(t, ok)
		require.Len(t, sum.DataPoints, 1)
		dp := sum.DataPoints[0]
		assert.Equal(t, int64(1), dp.Value)
		assert.Equal(t, set, dp.Attributes)
	})

	// Test bound instrument (Delta)
	t.Run("Bound/Delta", func(t *testing.T) {
		r := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}))
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Int64Binder)
		require.True(t, ok, "counter does not implement x.Int64Binder")

		bound := binder.Bind(attrs...)
		bound.Add(t.Context(), 1)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Equal(t, int64(1), sum.DataPoints[0].Value)

		// Record again on the bound instrument!
		bound.Add(t.Context(), 2)

		// Collect again. The value should be 2 (Delta!), not 3 (Cumulative!).
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Equal(t, int64(2), sum.DataPoints[0].Value)
	})
}

func TestBoundInstrumentMultipleReaders(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("K", "V")}
	r1 := NewManualReader()
	r2 := NewManualReader()
	mp := NewMeterProvider(WithReader(r1), WithReader(r2))
	meter := mp.Meter("test")

	counter, err := meter.Int64Counter("test.counter")
	require.NoError(t, err)

	binder, ok := counter.(x.Int64Binder)
	require.True(t, ok, "counter does not implement x.Int64Binder")

	// This triggers the "slow path" in Bind since len(aggregators) > 1
	bound := binder.Bind(attrs...)
	bound.Add(t.Context(), 1)

	// Verify Reader 1
	var rm1 metricdata.ResourceMetrics
	require.NoError(t, r1.Collect(t.Context(), &rm1))
	require.Len(t, rm1.ScopeMetrics, 1)
	sum1 := rm1.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	assert.Equal(t, int64(1), sum1.DataPoints[0].Value)

	// Verify Reader 2
	var rm2 metricdata.ResourceMetrics
	require.NoError(t, r2.Collect(t.Context(), &rm2))
	require.Len(t, rm2.ScopeMetrics, 1)
	sum2 := rm2.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	assert.Equal(t, int64(1), sum2.DataPoints[0].Value)
}

func TestBoundInstrumentDeltaConcurrency(t *testing.T) {
	r := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
		return metricdata.DeltaTemporality
	}))
	mp := NewMeterProvider(WithReader(r))
	meter := mp.Meter("test")

	counter, err := meter.Int64Counter("test.counter")
	require.NoError(t, err)

	binder, ok := counter.(x.Int64Binder)
	require.True(t, ok, "counter does not implement x.Int64Binder")

	bound := binder.Bind(attribute.String("K", "V"))

	// Number of goroutines and operations per goroutine
	numWorkers := 10
	opsPerWorker := 1000

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	ctx := t.Context()
	for range numWorkers {
		go func() {
			defer wg.Done()
			for j := range opsPerWorker {
				bound.Add(ctx, 1)
				// Occasionally collect to trigger Delta map clears and cycle pointer resets
				if j%100 == 0 {
					var rm metricdata.ResourceMetrics
					_ = r.Collect(ctx, &rm)
				}
			}
		}()
	}

	wg.Wait()

	// Final collection to ensure no panic occurs during cleanup.
	var rm metricdata.ResourceMetrics
	require.NoError(t, r.Collect(ctx, &rm))
}
