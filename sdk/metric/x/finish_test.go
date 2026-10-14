// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
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

func TestWithFinish(t *testing.T) {
	opt := sdkmetricx.WithFinish()
	_, experimental := opt.(interface{ Experimental() })
	assert.True(t, experimental)

	provider := sdkmetric.NewMeterProvider()
	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	_, ok := counter.(metricx.Finisher)
	assert.False(t, ok)

	provider = sdkmetric.NewMeterProvider(sdkmetricx.WithFinish())
	counter, err = provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	_, ok = counter.(metricx.Finisher)
	assert.True(t, ok)
}

func TestInt64CounterFinish(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	view := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "requests"},
		sdkmetric.Stream{
			AttributeFilter: attribute.NewAllowKeysFilter("route"),
		},
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(view),
		sdkmetricx.WithFinish(),
	)
	t.Cleanup(func() { _ = provider.Shutdown(t.Context()) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	finisher := counter.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/old"),
		attribute.String("route", "/orders"),
	}
	wantInput := slices.Clone(attrs)
	counter.Add(t.Context(), 3, metric.WithAttributes(attrs...))

	before := time.Now()
	finisher.Finish(t.Context(), attrs...)
	after := time.Now()
	assert.Equal(t, wantInput, attrs, "Finish modified its input")

	points := collectSum(t, reader)
	require.Len(t, points, 1)
	assert.Equal(t, int64(3), points[0].Value)
	assert.Equal(t, attribute.NewSet(attribute.String("route", "/orders")), points[0].Attributes)
	assert.False(t, points[0].Time.Before(before))
	assert.False(t, points[0].Time.After(after))
	assert.Empty(t, collectSum(t, reader))
}

func TestInt64CounterFinishReadersIndependent(t *testing.T) {
	reader1 := sdkmetric.NewManualReader()
	reader2 := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader1),
		sdkmetric.WithReader(reader2),
		sdkmetricx.WithFinish(),
	)
	t.Cleanup(func() { _ = provider.Shutdown(t.Context()) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	attrs := []attribute.KeyValue{attribute.String("route", "/orders")}
	counter.Add(t.Context(), 3, metric.WithAttributes(attrs...))
	counter.(metricx.Finisher).Finish(t.Context(), attrs...)

	require.Len(t, collectSum(t, reader1), 1)
	require.Len(t, collectSum(t, reader2), 1)
	assert.Empty(t, collectSum(t, reader1))
	assert.Empty(t, collectSum(t, reader2))
}

func TestInt64CounterFinishStreamsIndependently(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "requests"},
				sdkmetric.Stream{
					Name:            "requests.by_route",
					AttributeFilter: attribute.NewAllowKeysFilter("route"),
				},
			),
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "requests"},
				sdkmetric.Stream{
					Name:            "requests.by_zone",
					AttributeFilter: attribute.NewAllowKeysFilter("zone"),
				},
			),
		),
		sdkmetricx.WithFinish(),
	)
	t.Cleanup(func() { _ = provider.Shutdown(t.Context()) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	}
	counter.Add(t.Context(), 3, metric.WithAttributes(attrs...))
	counter.(metricx.Finisher).Finish(t.Context(), attrs...)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 2)
	got := make(map[string]attribute.Set, 2)
	var finishTime time.Time
	for _, m := range rm.ScopeMetrics[0].Metrics {
		sum, ok := m.Data.(metricdata.Sum[int64])
		require.True(t, ok)
		require.Len(t, sum.DataPoints, 1)
		got[m.Name] = sum.DataPoints[0].Attributes
		if finishTime.IsZero() {
			finishTime = sum.DataPoints[0].Time
		} else {
			assert.Equal(t, finishTime, sum.DataPoints[0].Time)
		}
	}
	assert.Equal(t, map[string]attribute.Set{
		"requests.by_route": attribute.NewSet(attribute.String("route", "/orders")),
		"requests.by_zone":  attribute.NewSet(attribute.String("zone", "east")),
	}, got)
	assert.Empty(t, collectSum(t, reader))
}

func collectSum(t *testing.T, reader *sdkmetric.ManualReader) []metricdata.DataPoint[int64] {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	if len(rm.ScopeMetrics) == 0 || len(rm.ScopeMetrics[0].Metrics) == 0 {
		return nil
	}
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	require.True(t, ok)
	return sum.DataPoints
}
