// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
)

func TestFloat64CounterWithFinish(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		options := []sdkmetric.Option{}
		if enabled {
			options = append(options, sdkmetricx.WithFinish())
		}
		provider := sdkmetric.NewMeterProvider(options...)
		t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
		meter := provider.Meter("test")
		counter, err := meter.Float64Counter("requests")
		require.NoError(t, err)
		_, ok := counter.(metricx.Finisher)
		assert.Equal(t, enabled, ok)
		again, err := meter.Float64Counter("requests")
		require.NoError(t, err)
		assert.Same(t, counter, again)
	}
}

func TestFloat64CounterFinishLifecycle(t *testing.T) {
	for _, temporality := range []metricdata.Temporality{
		metricdata.CumulativeTemporality, metricdata.DeltaTemporality,
	} {
		t.Run(temporality.String(), func(t *testing.T) {
			reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(
				func(sdkmetric.InstrumentKind) metricdata.Temporality { return temporality },
			))
			provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetricx.WithFinish())
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
			counter, err := provider.Meter("test").Float64Counter("requests")
			require.NoError(t, err)
			finisher, ok := counter.(metricx.Finisher)
			require.True(t, ok)
			attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
			option := metric.WithAttributes(attrs...)

			finisher.Finish(t.Context(), attrs...) // Unknown series.
			counter.Add(t.Context(), 1.25, option)
			finisher.Finish(t.Context(), attrs...)
			counter.Add(t.Context(), 0.5, option) // Reactivate before collection.
			points := collectFloat64Sum(t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, 1.75, points[0].Value)
			start := points[0].StartTime

			counter.Add(t.Context(), 0.25, option)
			before := time.Now()
			finisher.Finish(t.Context(), attrs...)
			after := time.Now()
			finisher.Finish(t.Context(), attrs...) // Keep the first Finish timestamp.
			points = collectFloat64Sum(t, reader)
			require.Len(t, points, 1)
			want := 2.0
			if temporality == metricdata.DeltaTemporality {
				want = 0.25
			}
			assert.Equal(t, want, points[0].Value)
			assert.False(t, points[0].Time.Before(before))
			assert.False(t, points[0].Time.After(after))
			if temporality == metricdata.CumulativeTemporality {
				assert.Equal(t, start, points[0].StartTime)
			}
			assert.Empty(t, collectFloat64Sum(t, reader))

			before = time.Now()
			counter.Add(t.Context(), 0.75, option) // Start a new lifetime.
			points = collectFloat64Sum(t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, 0.75, points[0].Value)
			if temporality == metricdata.CumulativeTemporality {
				assert.False(t, points[0].StartTime.Before(before))
			}
			finisher.Finish(t.Context(), attrs...)
			points = collectFloat64Sum(t, reader)
			if temporality == metricdata.DeltaTemporality {
				assert.Empty(t, points, "Finish must not synthesize an empty delta point")
			} else {
				require.Len(t, points, 1)
				assert.Equal(t, 0.75, points[0].Value)
			}
			assert.Empty(t, collectFloat64Sum(t, reader))
		})
	}
}

func TestFloat64CounterFinishFanout(t *testing.T) {
	readers := []*sdkmetric.ManualReader{sdkmetric.NewManualReader(), sdkmetric.NewManualReader()}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(readers[0]), sdkmetric.WithReader(readers[1]),
		sdkmetricx.WithFinish(),
		sdkmetric.WithView(
			sdkmetric.NewView(sdkmetric.Instrument{Name: "requests"}, sdkmetric.Stream{
				Name: "by_route", AttributeFilter: attribute.NewAllowKeysFilter("route"),
			}),
			sdkmetric.NewView(sdkmetric.Instrument{Name: "requests"}, sdkmetric.Stream{
				Name: "by_zone", AttributeFilter: attribute.NewAllowKeysFilter("zone"),
			}),
		),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	counter, err := provider.Meter("test").Float64Counter("requests")
	require.NoError(t, err)
	finisher, ok := counter.(metricx.Finisher)
	require.True(t, ok)
	counter.Add(t.Context(), 1.25, metric.WithAttributes(
		attribute.String("route", "/orders"), attribute.String("zone", "east"),
	))
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"), attribute.String("route", "/old"),
		attribute.String("route", "/orders"),
	}
	wantInput := slices.Clone(attrs)
	finisher.Finish(t.Context(), attrs...)
	assert.Equal(t, wantInput, attrs)
	var finishTime time.Time
	for _, reader := range readers {
		var rm metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)
		for _, m := range rm.ScopeMetrics[0].Metrics {
			sum, ok := m.Data.(metricdata.Sum[float64])
			require.True(t, ok)
			assert.True(t, sum.IsMonotonic)
			require.Len(t, sum.DataPoints, 1)
			point := sum.DataPoints[0]
			assert.Equal(t, 1.25, point.Value)
			want := attribute.NewSet(attribute.String("route", "/orders"))
			if m.Name == "by_zone" {
				want = attribute.NewSet(attribute.String("zone", "east"))
			} else {
				assert.Equal(t, "by_route", m.Name)
			}
			assert.Equal(t, want, point.Attributes)
			if finishTime.IsZero() {
				finishTime = point.Time
			}
			assert.Equal(t, finishTime, point.Time)
		}
		assert.Empty(t, collectFloat64Sum(t, reader))
	}
}

func TestFloat64CounterFinishUnsupportedAggregation(t *testing.T) {
	for _, agg := range []sdkmetric.Aggregation{
		sdkmetric.AggregationDrop{}, sdkmetric.AggregationExplicitBucketHistogram{},
	} {
		t.Run(fmt.Sprintf("%T", agg), func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			provider := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(reader), sdkmetricx.WithFinish(),
				sdkmetric.WithView(sdkmetric.NewView(sdkmetric.Instrument{Name: "requests"},
					sdkmetric.Stream{Aggregation: agg})),
			)
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
			counter, err := provider.Meter("test").Float64Counter("requests")
			require.NoError(t, err)
			finisher, ok := counter.(metricx.Finisher)
			require.True(t, ok)
			counter.Add(t.Context(), 1.25)
			finisher.Finish(t.Context())
			var before, after metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &before))
			require.NoError(t, reader.Collect(t.Context(), &after))
			if _, drop := agg.(sdkmetric.AggregationDrop); drop {
				assert.Empty(t, before.ScopeMetrics)
				assert.Empty(t, after.ScopeMetrics)
			} else {
				require.Len(t, after.ScopeMetrics, 1)
				require.Len(t, after.ScopeMetrics[0].Metrics, 1)
				hist, ok := after.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
				require.True(t, ok)
				require.Len(t, hist.DataPoints, 1)
				assert.Equal(t, uint64(1), hist.DataPoints[0].Count)
				assert.Equal(t, 1.25, hist.DataPoints[0].Sum)
			}
		})
	}
}

func collectFloat64Sum(t *testing.T, reader *sdkmetric.ManualReader) []metricdata.DataPoint[float64] {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	if len(rm.ScopeMetrics) == 0 {
		return nil
	}
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[float64])
	require.True(t, ok)
	assert.True(t, sum.IsMonotonic)
	return sum.DataPoints
}
