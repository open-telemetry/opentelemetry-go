// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestBoundInstrumentInt64(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("K", "V")}
	set := attribute.NewSet(attrs...)

	t.Run("Cumulative", func(t *testing.T) {
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

		// Record again
		bound.Add(t.Context(), 2)

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Equal(t, int64(3), sum.DataPoints[0].Value) // Cumulative
	})

	t.Run("Delta", func(t *testing.T) {
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

		// Record again
		bound.Add(t.Context(), 2)

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Equal(t, int64(2), sum.DataPoints[0].Value) // Delta
	})
}
