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

	t.Run("UnboundThenBound_Delta", func(t *testing.T) {
		r := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}))
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Int64Binder)
		require.True(t, ok)

		// 1. Unbound write
		counter.Add(t.Context(), 1, metricapi.WithAttributes(attrs...))

		// 2. Bind the same attributes
		bound := binder.Bind(attrs...)

		// 3. Bound write
		bound.Add(t.Context(), 2)

		// 4. Collect - should produce exactly 1 data point with value 3 (1+2), no duplicates
		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		require.Len(t, sum.DataPoints, 1, "must not produce duplicate data points")
		assert.Equal(t, int64(3), sum.DataPoints[0].Value)

		// 5. Subsequent cycle with both bound and unbound writes
		bound.Add(t.Context(), 4)
		counter.Add(t.Context(), 1, metricapi.WithAttributes(attrs...))

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		require.Len(t, sum.DataPoints, 1, "must not produce duplicate data points")
		assert.Equal(t, int64(5), sum.DataPoints[0].Value)
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

		// Test x.WithUnsafeAttributes
		unsafeAttr := attribute.String("unsafe", "val")
		bound.Add(t.Context(), 10, x.WithUnsafeAttributes(unsafeAttr))
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		expectedUnsafeSet := attribute.NewSet(attribute.String("K", "V"), unsafeAttr)
		found = false
		for _, dp := range sum.DataPoints {
			if dp.Attributes.Equals(&expectedUnsafeSet) {
				assert.Equal(t, int64(10), dp.Value)
				found = true
			}
		}
		assert.True(t, found, "expected data point with unsafe attributes not found")
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

	t.Run("MixedAggregations", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(
			WithReader(r),
			WithView(
				NewView(
					Instrument{Name: "test.counter"},
					Stream{
						Name:        "test.counter.hist",
						Aggregation: AggregationExplicitBucketHistogram{Boundaries: []float64{0, 10, 100}},
					},
				),
				NewView(
					Instrument{Name: "test.counter"},
					Stream{Aggregation: AggregationDefault{}},
				),
			),
		)
		meter := mp.Meter("test")

		counter, err := meter.Int64Counter("test.counter")
		require.NoError(t, err)
		bound := counter.(x.Int64Binder).Bind(attrs...)
		bound.Add(t.Context(), 5)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

		var sumFound, histFound bool
		for _, m := range rm.ScopeMetrics[0].Metrics {
			switch data := m.Data.(type) {
			case metricdata.Sum[int64]:
				assert.Equal(t, "test.counter", m.Name)
				require.Len(t, data.DataPoints, 1)
				assert.Equal(t, int64(5), data.DataPoints[0].Value)
				sumFound = true
			case metricdata.Histogram[int64]:
				assert.Equal(t, "test.counter.hist", m.Name)
				require.Len(t, data.DataPoints, 1)
				assert.Equal(t, uint64(1), data.DataPoints[0].Count)
				assert.Equal(t, int64(5), data.DataPoints[0].Sum)
				histFound = true
			}
		}
		assert.True(t, sumFound, "expected Sum metric stream")
		assert.True(t, histFound, "expected Histogram metric stream")
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

	t.Run("UnboundThenBound_Delta", func(t *testing.T) {
		r := NewManualReader(WithTemporalitySelector(func(InstrumentKind) metricdata.Temporality {
			return metricdata.DeltaTemporality
		}))
		mp := NewMeterProvider(WithReader(r))
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)

		binder, ok := counter.(x.Float64Binder)
		require.True(t, ok)

		// 1. Unbound write
		counter.Add(t.Context(), 1.5, metricapi.WithAttributes(attrs...))

		// 2. Bind the same attributes
		bound := binder.Bind(attrs...)

		// 3. Bound write
		bound.Add(t.Context(), 2.5)

		// 4. Collect - should produce exactly 1 data point with value 4.0 (1.5+2.5), no duplicates
		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		require.Len(t, sum.DataPoints, 1, "must not produce duplicate data points")
		assert.Equal(t, float64(4.0), sum.DataPoints[0].Value)

		// 5. Subsequent cycle with both bound and unbound writes
		bound.Add(t.Context(), 4.5)
		counter.Add(t.Context(), 1.5, metricapi.WithAttributes(attrs...))

		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		require.Len(t, sum.DataPoints, 1, "must not produce duplicate data points")
		assert.Equal(t, float64(6.0), sum.DataPoints[0].Value)
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

		// Test x.WithUnsafeAttributes
		unsafeAttr := attribute.String("unsafe", "val")
		bound.Add(t.Context(), 10.5, x.WithUnsafeAttributes(unsafeAttr))
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		sum = rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
		expectedUnsafeSet := attribute.NewSet(attribute.String("K", "V"), unsafeAttr)
		found = false
		for _, dp := range sum.DataPoints {
			if dp.Attributes.Equals(&expectedUnsafeSet) {
				assert.Equal(t, float64(10.5), dp.Value)
				found = true
			}
		}
		assert.True(t, found, "expected data point with unsafe attributes not found")
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

	t.Run("MixedAggregations", func(t *testing.T) {
		r := NewManualReader()
		mp := NewMeterProvider(
			WithReader(r),
			WithView(
				NewView(
					Instrument{Name: "test.counter"},
					Stream{
						Name:        "test.counter.hist",
						Aggregation: AggregationExplicitBucketHistogram{Boundaries: []float64{0, 10, 100}},
					},
				),
				NewView(
					Instrument{Name: "test.counter"},
					Stream{Aggregation: AggregationDefault{}},
				),
			),
		)
		meter := mp.Meter("test")

		counter, err := meter.Float64Counter("test.counter")
		require.NoError(t, err)
		bound := counter.(x.Float64Binder).Bind(attrs...)
		bound.Add(t.Context(), 5.5)

		var rm metricdata.ResourceMetrics
		err = r.Collect(t.Context(), &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

		var sumFound, histFound bool
		for _, m := range rm.ScopeMetrics[0].Metrics {
			switch data := m.Data.(type) {
			case metricdata.Sum[float64]:
				assert.Equal(t, "test.counter", m.Name)
				require.Len(t, data.DataPoints, 1)
				assert.Equal(t, float64(5.5), data.DataPoints[0].Value)
				sumFound = true
			case metricdata.Histogram[float64]:
				assert.Equal(t, "test.counter.hist", m.Name)
				require.Len(t, data.DataPoints, 1)
				assert.Equal(t, uint64(1), data.DataPoints[0].Count)
				assert.Equal(t, float64(5.5), data.DataPoints[0].Sum)
				histFound = true
			}
		}
		assert.True(t, sumFound, "expected Sum metric stream")
		assert.True(t, histFound, "expected Histogram metric stream")
	})
}

func TestBindDoesNotMutateAttributeSlice(t *testing.T) {
	r := NewManualReader()
	mp := NewMeterProvider(WithReader(r))
	meter := mp.Meter("test")

	counter, err := meter.Int64Counter("test.counter")
	require.NoError(t, err)

	attrs := []attribute.KeyValue{
		attribute.Int("b", 2),
		attribute.Int("a", 1),
		attribute.Int("a", 3),
	}
	original := []attribute.KeyValue{
		attribute.Int("b", 2),
		attribute.Int("a", 1),
		attribute.Int("a", 3),
	}

	_ = counter.(x.Int64Binder).Bind(attrs...)
	assert.Equal(t, original, attrs, "Bind mutated the caller's attribute slice")
}
