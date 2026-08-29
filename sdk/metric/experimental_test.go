// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type experimentalBindingOption struct{ Option }

func (experimentalBindingOption) WrapInt64CounterBinding(
	counter api.Int64Counter,
	bind func(...attribute.KeyValue) api.Int64Counter,
) api.Int64Counter {
	return &experimentalBindingCounter{Int64Counter: counter, bind: bind}
}

type experimentalFinishingOption struct{ Option }

func (experimentalFinishingOption) WrapInt64CounterFinishing(
	counter api.Int64Counter,
	finish func(context.Context, ...attribute.KeyValue),
) api.Int64Counter {
	if binder, ok := counter.(metricx.Int64CounterBinder); ok {
		return &experimentalBindingFinishingCounter{
			Int64Counter: counter,
			bind:         binder.Bind,
			finish:       finish,
		}
	}
	return &experimentalFinishingCounter{Int64Counter: counter, finish: finish}
}

type experimentalBindingCounter struct {
	api.Int64Counter
	bind func(...attribute.KeyValue) api.Int64Counter
}

func (c *experimentalBindingCounter) Bind(attrs ...attribute.KeyValue) api.Int64Counter {
	return c.bind(attrs...)
}

type experimentalFinishingCounter struct {
	api.Int64Counter
	finish func(context.Context, ...attribute.KeyValue)
}

func (c *experimentalFinishingCounter) Finish(ctx context.Context, attrs ...attribute.KeyValue) {
	c.finish(ctx, attrs...)
}

type experimentalBindingFinishingCounter struct {
	api.Int64Counter
	bind   func(...attribute.KeyValue) api.Int64Counter
	finish func(context.Context, ...attribute.KeyValue)
}

func (c *experimentalBindingFinishingCounter) Bind(attrs ...attribute.KeyValue) api.Int64Counter {
	return c.bind(attrs...)
}

func (c *experimentalBindingFinishingCounter) Finish(
	ctx context.Context,
	attrs ...attribute.KeyValue,
) {
	c.finish(ctx, attrs...)
}

func TestExperimentalCounterLifecycle(t *testing.T) {
	cumulative := NewManualReader()
	delta := NewManualReader(WithTemporalitySelector(DeltaTemporalitySelector))
	view := NewView(
		Instrument{Name: "requests"},
		Stream{AttributeFilter: attribute.NewAllowKeysFilter("route")},
	)
	provider := NewMeterProvider(
		experimentalBindingOption{Option: WithView()},
		experimentalFinishingOption{Option: WithView()},
		WithReader(cumulative),
		WithReader(delta),
		WithView(view),
		WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	binder := counter.(metricx.Int64CounterBinder)
	finisher := counter.(metricx.Finisher)
	attrs := []attribute.KeyValue{
		attribute.String("route", "/old"),
		attribute.String("secret", "drop"),
	}
	bound := binder.Bind(attrs...)
	attrs[0] = attribute.String("route", "/mutated")
	assert.True(t, bound.Enabled(t.Context()))

	bound.Add(t.Context(), 4)
	bound.Add(t.Context(), 1, api.WithAttributes())
	bound.Add(t.Context(), 2, api.WithAttributes(attribute.String("route", "/new")))
	finisher.Finish(t.Context(), attribute.String("route", "/old"))
	finisher.Finish(t.Context(), attribute.String("route", "/old"))
	bound.Add(t.Context(), 3)

	for _, reader := range []*ManualReader{cumulative, delta} {
		first := experimentalSumPoints(t, collectExperimental(t, reader))
		assert.Equal(t, int64(5), experimentalRouteValue(first, "/old"))
		assert.Equal(t, int64(2), experimentalRouteValue(first, "/new"))

		second := experimentalSumPoints(t, collectExperimental(t, reader))
		assert.Equal(t, int64(3), experimentalRouteValue(second, "/old"))
	}
}

func TestExperimentalCounterFallbackAggregation(t *testing.T) {
	reader := NewManualReader()
	view := NewView(
		Instrument{Name: "requests"},
		Stream{Aggregation: AggregationExplicitBucketHistogram{Boundaries: []float64{0, 10}}},
	)
	provider := NewMeterProvider(
		experimentalBindingOption{Option: WithView()},
		experimentalFinishingOption{Option: WithView()},
		WithReader(reader),
		WithView(view),
		WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	bound := counter.(metricx.Int64CounterBinder).Bind(attribute.String("route", "/orders"))
	bound.Add(t.Context(), 5)
	counter.(metricx.Finisher).Finish(t.Context(), attribute.String("route", "/orders"))

	rm := collectExperimental(t, reader)
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	histogram, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[int64])
	require.True(t, ok)
	require.Len(t, histogram.DataPoints, 1)
	assert.Equal(t, uint64(1), histogram.DataPoints[0].Count)
	assert.Equal(t, int64(5), histogram.DataPoints[0].Sum)
}

func TestExperimentalCounterFinishOnly(t *testing.T) {
	reader := NewManualReader()
	provider := NewMeterProvider(
		experimentalFinishingOption{Option: WithView()},
		WithReader(reader),
		WithExemplarFilter(exemplar.AlwaysOffFilter),
	)
	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	_, isBinder := counter.(metricx.Int64CounterBinder)
	assert.False(t, isBinder)

	counter.Add(t.Context(), 1, api.WithAttributes(attribute.String("route", "/orders")))
	counter.(metricx.Finisher).Finish(t.Context(), attribute.String("route", "/orders"))
	points := experimentalSumPoints(t, collectExperimental(t, reader))
	assert.Equal(t, int64(1), experimentalRouteValue(points, "/orders"))
}

func collectExperimental(t *testing.T, reader *ManualReader) metricdata.ResourceMetrics {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	return rm
}

func experimentalSumPoints(
	t *testing.T,
	rm metricdata.ResourceMetrics,
) []metricdata.DataPoint[int64] {
	t.Helper()
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	require.True(t, ok)
	return sum.DataPoints
}

func experimentalRouteValue(
	points []metricdata.DataPoint[int64],
	value string,
) int64 {
	for _, point := range points {
		if got, ok := point.Attributes.Value("route"); ok && got.AsString() == value {
			return point.Value
		}
	}
	return 0
}
