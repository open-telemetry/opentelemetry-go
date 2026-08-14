// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	metricapi "go.opentelemetry.io/otel/metric"
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

	t.Run("WithOptions", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Int64Binder)
		require.True(t, ok)

		bound := binder.Bind(attrs...)
		extraAttr := attribute.String("extra", "attr")
		bound.Add(t.Context(), 5, metricapi.WithAttributes(extraAttr))

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		expectedSet := attribute.NewSet(attribute.String("K", "V"), extraAttr)

		var found bool
		for _, dp := range sum.DataPoints {
			if dp.Attributes.Equals(&expectedSet) {
				assert.Equal(t, int64(5), dp.Value)
				found = true
			}
		}
		assert.True(t, found, "expected data point with merged attributes not found")
	})

	t.Run("Enabled", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(
			WithReader(r),
			WithView(NewView(
				Instrument{Name: "drop.counter"},
				Stream{Aggregation: AggregationDrop{}},
			)),
		)
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)
		bound := counter.(x.Int64Binder).Bind(attrs...)
		assert.True(t, bound.Enabled(t.Context()))

		droppedCounter, err := meter.Int64Counter("drop.counter")
		require.NoError(t, err)
		boundDropped := droppedCounter.(x.Int64Binder).Bind(attrs...)
		assert.False(t, boundDropped.Enabled(t.Context()))
	})
}

func TestBoundInstrumentFloat64(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.String("K", "V")}
	set := attribute.NewSet(attrs...)

	t.Run("Cumulative", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Float64Binder)
		require.True(t, ok, "counter does not implement x.Float64Binder")

		bound := binder.Bind(attrs...)
		bound.Add(t.Context(), 1.5)

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
		assert.Equal(t, float64(1.5), dp.Value)
		assert.Equal(t, set, dp.Attributes)

		// Record again
		bound.Add(t.Context(), 2.5)

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		assert.Equal(t, float64(4.0), sum.DataPoints[0].Value) // Cumulative
	})

	t.Run("Delta", func(t *testing.T) {
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
		bound.Add(t.Context(), 1.5)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		assert.Equal(t, float64(1.5), sum.DataPoints[0].Value)

		// Record again
		bound.Add(t.Context(), 2.5)

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		assert.Equal(t, float64(2.5), sum.DataPoints[0].Value) // Delta
	})

	t.Run("WithOptions", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Float64Binder)
		require.True(t, ok)

		bound := binder.Bind(attrs...)
		extraAttr := attribute.String("extra", "attr")
		bound.Add(t.Context(), 5.5, metricapi.WithAttributes(extraAttr))

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		expectedSet := attribute.NewSet(attribute.String("K", "V"), extraAttr)

		var found bool
		for _, dp := range sum.DataPoints {
			if dp.Attributes.Equals(&expectedSet) {
				assert.Equal(t, float64(5.5), dp.Value)
				found = true
			}
		}
		assert.True(t, found, "expected data point with merged attributes not found")
	})

	t.Run("Enabled", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(
			WithReader(r),
			WithView(NewView(
				Instrument{Name: "drop.counter"},
				Stream{Aggregation: AggregationDrop{}},
			)),
		)
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)
		bound := counter.(x.Float64Binder).Bind(attrs...)
		assert.True(t, bound.Enabled(t.Context()))

		droppedCounter, err := meter.Float64Counter("drop.counter")
		require.NoError(t, err)
		boundDropped := droppedCounter.(x.Float64Binder).Bind(attrs...)
		assert.False(t, boundDropped.Enabled(t.Context()))
	})
}
