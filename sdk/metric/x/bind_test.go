// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdkmetricx "go.opentelemetry.io/otel/sdk/metric/x"
	"go.opentelemetry.io/otel/trace"
)

func TestWithBinding(t *testing.T) {
	opt := sdkmetricx.WithBinding()
	_, experimental := opt.(interface{ Experimental() })
	assert.True(t, experimental)

	stableProvider := sdkmetric.NewMeterProvider()
	t.Cleanup(func() { assert.NoError(t, stableProvider.Shutdown(context.WithoutCancel(t.Context()))) })
	stableMeter := stableProvider.Meter("test")

	ic, err := stableMeter.Int64Counter("ic")
	require.NoError(t, err)
	_, ok := ic.(metricx.Int64CounterBinder)
	assert.False(t, ok)

	fc, err := stableMeter.Float64Counter("fc")
	require.NoError(t, err)
	_, ok = fc.(metricx.Float64CounterBinder)
	assert.False(t, ok)

	iudc, err := stableMeter.Int64UpDownCounter("iudc")
	require.NoError(t, err)
	_, ok = iudc.(metricx.Int64UpDownCounterBinder)
	assert.False(t, ok)

	fudc, err := stableMeter.Float64UpDownCounter("fudc")
	require.NoError(t, err)
	_, ok = fudc.(metricx.Float64UpDownCounterBinder)
	assert.False(t, ok)

	bindProvider := sdkmetric.NewMeterProvider(sdkmetricx.WithBinding())
	t.Cleanup(func() { assert.NoError(t, bindProvider.Shutdown(context.WithoutCancel(t.Context()))) })
	bindMeter := bindProvider.Meter("test")

	ic, err = bindMeter.Int64Counter("ic")
	require.NoError(t, err)
	_, ok = ic.(metricx.Int64CounterBinder)
	assert.True(t, ok)

	fc, err = bindMeter.Float64Counter("fc")
	require.NoError(t, err)
	_, ok = fc.(metricx.Float64CounterBinder)
	assert.True(t, ok)

	iudc, err = bindMeter.Int64UpDownCounter("iudc")
	require.NoError(t, err)
	_, ok = iudc.(metricx.Int64UpDownCounterBinder)
	assert.True(t, ok)

	fudc, err = bindMeter.Float64UpDownCounter("fudc")
	require.NoError(t, err)
	_, ok = fudc.(metricx.Float64UpDownCounterBinder)
	assert.True(t, ok)

	// Instruments without bound aggregation support remain unaffected.
	ig, err := bindMeter.Int64Gauge("ig")
	require.NoError(t, err)
	_, ok = ig.(metricx.Int64GaugeBinder)
	assert.False(t, ok)

	fg, err := bindMeter.Float64Gauge("fg")
	require.NoError(t, err)
	_, ok = fg.(metricx.Float64GaugeBinder)
	assert.False(t, ok)

	ih, err := bindMeter.Int64Histogram("ih")
	require.NoError(t, err)
	_, ok = ih.(metricx.Int64HistogramBinder)
	assert.False(t, ok)

	fh, err := bindMeter.Float64Histogram("fh")
	require.NoError(t, err)
	_, ok = fh.(metricx.Float64HistogramBinder)
	assert.False(t, ok)
}

func deltaTemporalitySelector(sdkmetric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func cumulativeTemporalitySelector(sdkmetric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func collectSumPoints[N int64 | float64](t *testing.T, reader *sdkmetric.ManualReader) []metricdata.DataPoint[N] {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	if len(rm.ScopeMetrics) == 0 || len(rm.ScopeMetrics[0].Metrics) == 0 {
		return nil
	}
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[N])
	require.True(t, ok)
	return sum.DataPoints
}

func sumPointsByAttr[N int64 | float64](pts []metricdata.DataPoint[N]) map[attribute.Set]N {
	out := make(map[attribute.Set]N, len(pts))
	for _, pt := range pts {
		out[pt.Attributes] = pt.Value
	}
	return out
}

func TestBoundCountersAndUpDownCounters(t *testing.T) {
	temporalities := []struct {
		name string
		temp metricdata.Temporality
	}{
		{name: "Cumulative", temp: metricdata.CumulativeTemporality},
		{name: "Delta", temp: metricdata.DeltaTemporality},
	}

	for _, tc := range temporalities {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("Int64Counter", func(t *testing.T) {
				testBoundSumBehavior(
					t,
					tc.temp,
					true,
					func(m metric.Meter, name string) (metric.Int64Counter, error) {
						return m.Int64Counter(name)
					},
					func(inst metric.Int64Counter, attrs ...attribute.KeyValue) metric.Int64Counter {
						return inst.(metricx.Int64CounterBinder).Bind(attrs...)
					},
					func(ctx context.Context, inst metric.Int64Counter, v int64, opts ...metric.AddOption) {
						inst.Add(ctx, v, opts...)
					},
					int64(2),
					int64(3),
				)
			})
			t.Run("Float64Counter", func(t *testing.T) {
				testBoundSumBehavior(
					t,
					tc.temp,
					true,
					func(m metric.Meter, name string) (metric.Float64Counter, error) {
						return m.Float64Counter(name)
					},
					func(inst metric.Float64Counter, attrs ...attribute.KeyValue) metric.Float64Counter {
						return inst.(metricx.Float64CounterBinder).Bind(attrs...)
					},
					func(ctx context.Context, inst metric.Float64Counter, v float64, opts ...metric.AddOption) {
						inst.Add(ctx, v, opts...)
					},
					float64(1.5),
					float64(2.25),
				)
			})
			t.Run("Int64UpDownCounter", func(t *testing.T) {
				testBoundSumBehavior(
					t,
					tc.temp,
					false,
					func(m metric.Meter, name string) (metric.Int64UpDownCounter, error) {
						return m.Int64UpDownCounter(name)
					},
					func(inst metric.Int64UpDownCounter, attrs ...attribute.KeyValue) metric.Int64UpDownCounter {
						return inst.(metricx.Int64UpDownCounterBinder).Bind(attrs...)
					},
					func(ctx context.Context, inst metric.Int64UpDownCounter, v int64, opts ...metric.AddOption) {
						inst.Add(ctx, v, opts...)
					},
					int64(5),
					int64(-2),
				)
			})
			t.Run("Float64UpDownCounter", func(t *testing.T) {
				testBoundSumBehavior(
					t,
					tc.temp,
					false,
					func(m metric.Meter, name string) (metric.Float64UpDownCounter, error) {
						return m.Float64UpDownCounter(name)
					},
					func(inst metric.Float64UpDownCounter, attrs ...attribute.KeyValue) metric.Float64UpDownCounter {
						return inst.(metricx.Float64UpDownCounterBinder).Bind(attrs...)
					},
					func(ctx context.Context, inst metric.Float64UpDownCounter, v float64, opts ...metric.AddOption) {
						inst.Add(ctx, v, opts...)
					},
					float64(3.5),
					float64(-1.25),
				)
			})
		})
	}
}

func testBoundSumBehavior[I any, N int64 | float64](
	t *testing.T,
	temporality metricdata.Temporality,
	wantMonotonic bool,
	create func(metric.Meter, string) (I, error),
	bind func(I, ...attribute.KeyValue) I,
	add func(context.Context, I, N, ...metric.AddOption),
	v1, v2 N,
) {
	t.Helper()

	selector := func(sdkmetric.InstrumentKind) metricdata.Temporality { return temporality }
	reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(selector))
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	inst, err := create(provider.Meter("test"), "requests")
	require.NoError(t, err)

	// Input slice must not be modified or retained by Bind.
	rawAttrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/old"),
		attribute.String("route", "/orders"),
	}
	wantRaw := slices.Clone(rawAttrs)
	bound := bind(inst, rawAttrs...)
	assert.Equal(t, wantRaw, rawAttrs, "Bind modified caller slice")
	rawAttrs[0] = attribute.String("zone", "mutated")

	// Idle bound instrument emits nothing until measured.
	assert.Empty(t, collectSumPoints[N](t, reader))

	// Bound + unbound measurements for the same series merge into one point.
	add(t.Context(), bound, v1)
	add(t.Context(), inst, v2, metric.WithAttributes(
		attribute.String("zone", "east"),
		attribute.String("route", "/orders"),
	))

	// Empty AddOption still uses the bound series.
	add(t.Context(), bound, v1, metric.WithAttributes())

	// Record-time attributes merge over bound attributes (record-time wins).
	add(t.Context(), bound, v2, metric.WithAttributes(
		attribute.String("route", "/checkout"),
		attribute.String("method", "POST"),
	))

	// Empty Bind() with no attributes.
	emptyBound := bind(inst)
	add(t.Context(), emptyBound, v1)
	add(t.Context(), emptyBound, v2, metric.WithAttributes(attribute.String("only", "extra")))

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[N])
	require.True(t, ok)
	assert.Equal(t, wantMonotonic, sum.IsMonotonic)

	boundSet := attribute.NewSet(
		attribute.String("route", "/orders"),
		attribute.String("zone", "east"),
	)
	mergedSet := attribute.NewSet(
		attribute.String("method", "POST"),
		attribute.String("route", "/checkout"),
		attribute.String("zone", "east"),
	)
	extraOnlySet := attribute.NewSet(attribute.String("only", "extra"))

	got := sumPointsByAttr(sum.DataPoints)
	assert.Equal(t, map[attribute.Set]N{
		boundSet:           v1 + v2 + v1,
		mergedSet:          v2,
		attribute.NewSet(): v1,
		extraOnlySet:       v2,
	}, got)

	// Second collection: delta emits nothing while idle, then resumes cleanly.
	if temporality == metricdata.DeltaTemporality {
		assert.Empty(t, collectSumPoints[N](t, reader))
		add(t.Context(), bound, v2)
		got2 := sumPointsByAttr(collectSumPoints[N](t, reader))
		assert.Equal(t, map[attribute.Set]N{boundSet: v2}, got2)
	} else {
		add(t.Context(), bound, v2)
		got2 := sumPointsByAttr(collectSumPoints[N](t, reader))
		assert.Equal(t, v1+v2+v1+v2, got2[boundSet])
	}
}

func TestBoundNestedMapNormalization(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)

	bound := counter.(metricx.Int64CounterBinder).Bind(
		attribute.Map("meta", attribute.String("k", "old"), attribute.String("k", "new")),
	)
	bound.Add(t.Context(), 2)
	// Same normalized map via unbound Add must aggregate into the same series.
	counter.Add(t.Context(), 3, metric.WithAttributes(
		attribute.Map("meta", attribute.String("k", "new")),
	))
	// Record-time map with duplicate keys merged onto bound attributes.
	bound.Add(t.Context(), 4, metric.WithAttributes(
		attribute.Map("extra", attribute.String("x", "1"), attribute.String("x", "2")),
	))

	got := sumPointsByAttr(collectSumPoints[int64](t, reader))
	assert.Equal(t, map[attribute.Set]int64{
		attribute.NewSet(attribute.Map("meta", attribute.String("k", "new"))): 5,
		attribute.NewSet(
			attribute.Map("meta", attribute.String("k", "new")),
			attribute.Map("extra", attribute.String("x", "2")),
		): 4,
	}, got)
}

func TestBoundCardinalityLimitEvaluatedAtBindTime(t *testing.T) {
	// Cardinality limit 2 => 1 concrete series + 1 overflow series.
	reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(deltaTemporalitySelector))
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithCardinalityLimit(2),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	binder := counter.(metricx.Int64CounterBinder)

	// Fill the single concrete slot with an unbound series.
	counter.Add(t.Context(), 1, metric.WithAttributes(attribute.String("id", "unbound")))

	// Bind while at limit: must resolve to overflow and remain fixed across cycles.
	overflowBound := binder.Bind(attribute.String("id", "late"))
	overflowBound.Add(t.Context(), 2)

	overflowAttr := attribute.NewSet(attribute.Bool("otel.metric.overflow", true))
	unboundAttr := attribute.NewSet(attribute.String("id", "unbound"))
	lateAttr := attribute.NewSet(attribute.String("id", "late"))

	got1 := sumPointsByAttr(collectSumPoints[int64](t, reader))
	assert.Equal(t, map[attribute.Set]int64{
		unboundAttr:  1,
		overflowAttr: 2,
	}, got1)

	// After delta collection clears the unbound series, a NEW Bind call for the
	// same attributes is evaluated independently and gets the concrete slot,
	// while the earlier bound handle stays on overflow (including when invoked
	// with empty attribute options).
	concreteBound := binder.Bind(attribute.String("id", "late"))
	overflowBound.Add(t.Context(), 5)
	overflowBound.Add(t.Context(), 3, metric.WithAttributes())
	overflowBound.Add(t.Context(), 2, metric.WithAttributeSet(attribute.NewSet()))
	concreteBound.Add(t.Context(), 7)

	got2 := sumPointsByAttr(collectSumPoints[int64](t, reader))
	assert.Equal(t, map[attribute.Set]int64{
		overflowAttr: 10,
		lateAttr:     7,
	}, got2)
}

func TestBoundMultipleReadersAndViews(t *testing.T) {
	cumReader := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(cumulativeTemporalitySelector),
	)
	deltaReader := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(deltaTemporalitySelector),
	)
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(cumReader),
		sdkmetric.WithReader(deltaReader),
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
					Name:        "requests.hist",
					Aggregation: sdkmetric.AggregationExplicitBucketHistogram{Boundaries: []float64{5, 10}},
				},
			),
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "requests"},
				sdkmetric.Stream{
					Name:        "requests.dropped",
					Aggregation: sdkmetric.AggregationDrop{},
				},
			),
		),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	bound := counter.(metricx.Int64CounterBinder).Bind(
		attribute.String("route", "/orders"),
		attribute.String("zone", "east"),
	)
	bound.Add(t.Context(), 4)

	checkReader := func(r *sdkmetric.ManualReader, wantTemp metricdata.Temporality) {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

		byName := make(map[string]metricdata.Metrics, 2)
		for _, m := range rm.ScopeMetrics[0].Metrics {
			byName[m.Name] = m
		}

		sum, ok := byName["requests.by_route"].Data.(metricdata.Sum[int64])
		require.True(t, ok)
		assert.Equal(t, wantTemp, sum.Temporality)
		require.Len(t, sum.DataPoints, 1)
		assert.Equal(t, int64(4), sum.DataPoints[0].Value)
		assert.Equal(t, attribute.NewSet(attribute.String("route", "/orders")), sum.DataPoints[0].Attributes)

		hist, ok := byName["requests.hist"].Data.(metricdata.Histogram[int64])
		require.True(t, ok)
		assert.Equal(t, wantTemp, hist.Temporality)
		require.Len(t, hist.DataPoints, 1)
		assert.Equal(t, uint64(1), hist.DataPoints[0].Count)
		assert.Equal(t, int64(4), hist.DataPoints[0].Sum)
		assert.Equal(t, attribute.NewSet(
			attribute.String("route", "/orders"),
			attribute.String("zone", "east"),
		), hist.DataPoints[0].Attributes)
	}

	checkReader(cumReader, metricdata.CumulativeTemporality)
	checkReader(deltaReader, metricdata.DeltaTemporality)

	// Delta reader resets after collection; cumulative reader retains data.
	var rm metricdata.ResourceMetrics
	require.NoError(t, deltaReader.Collect(t.Context(), &rm))
	assert.Empty(t, rm.ScopeMetrics)

	checkReader(cumReader, metricdata.CumulativeTemporality)

	// Record-time attributes merge and fan out across both Sum and Histogram streams.
	bound.Add(t.Context(), 6, metric.WithAttributes(
		attribute.String("route", "/checkout"),
		attribute.String("zone", "west"),
	))
	require.NoError(t, deltaReader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 2)
	mergedByName := make(map[string]metricdata.Metrics, 2)
	for _, m := range rm.ScopeMetrics[0].Metrics {
		mergedByName[m.Name] = m
	}
	mergedSum := mergedByName["requests.by_route"].Data.(metricdata.Sum[int64])
	require.Len(t, mergedSum.DataPoints, 1)
	assert.Equal(t, int64(6), mergedSum.DataPoints[0].Value)
	assert.Equal(t, attribute.NewSet(attribute.String("route", "/checkout")), mergedSum.DataPoints[0].Attributes)

	mergedHist := mergedByName["requests.hist"].Data.(metricdata.Histogram[int64])
	require.Len(t, mergedHist.DataPoints, 1)
	assert.Equal(t, uint64(1), mergedHist.DataPoints[0].Count)
	assert.Equal(t, int64(6), mergedHist.DataPoints[0].Sum)
	assert.Equal(t, attribute.NewSet(
		attribute.String("route", "/checkout"),
		attribute.String("zone", "west"),
	), mergedHist.DataPoints[0].Attributes)
}

func TestBoundFallbackOnlyAggregation(t *testing.T) {
	reader := sdkmetric.NewManualReader(sdkmetric.WithAggregationSelector(
		func(sdkmetric.InstrumentKind) sdkmetric.Aggregation {
			return sdkmetric.AggregationBase2ExponentialHistogram{MaxSize: 160, MaxScale: 20}
		},
	))
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	bound := counter.(metricx.Int64CounterBinder).Bind(attribute.String("route", "/orders"))
	assert.True(t, bound.Enabled(t.Context()))

	bound.Add(t.Context(), 4)
	bound.Add(t.Context(), 8, metric.WithAttributes(attribute.String("zone", "east")))

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	expHist, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.ExponentialHistogram[int64])
	require.True(t, ok)

	got := make(map[attribute.Set]int64, len(expHist.DataPoints))
	for _, dp := range expHist.DataPoints {
		got[dp.Attributes] = dp.Sum
	}
	assert.Equal(t, map[attribute.Set]int64{
		attribute.NewSet(attribute.String("route", "/orders")): 4,
		attribute.NewSet(
			attribute.String("route", "/orders"),
			attribute.String("zone", "east"),
		): 8,
	}, got)
}

func TestBoundComposableViewsAndOptIn(t *testing.T) {
	for _, bindFirst := range []bool{false, true} {
		t.Run(fmt.Sprintf("bindFirst=%t", bindFirst), func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			options := []sdkmetric.Option{
				sdkmetricx.WithBinding(),
				sdkmetricx.WithViewMatchingMode(sdkmetricx.ViewMatchingModeComposable),
			}
			if !bindFirst {
				slices.Reverse(options)
			}
			options = append(options,
				sdkmetric.WithReader(reader),
				sdkmetricx.WithView(
					sdkmetricx.NewView(
						sdkmetric.Instrument{Name: "requests"},
						sdkmetricx.Stream{
							Name:    "requests.enabled",
							Enabled: new(true),
						},
					),
					sdkmetricx.NewView(
						sdkmetric.Instrument{Name: "requests"},
						sdkmetricx.Stream{
							AttributeFilter: attribute.NewAllowKeysFilter("route"),
						},
					),
				),
			)
			provider := sdkmetric.NewMeterProvider(options...)
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

			meter := provider.Meter("test")
			enabledCounter, err := meter.Int64Counter("requests", metricx.WithOptIn())
			require.NoError(t, err)
			assert.True(t, enabledCounter.Enabled(t.Context()))

			disabledCounter, err := meter.Int64Counter("disabled", metricx.WithOptIn())
			require.NoError(t, err)
			assert.False(t, disabledCounter.Enabled(t.Context()))
			disabledBound := disabledCounter.(metricx.Int64CounterBinder).Bind(attribute.String("route", "/x"))
			assert.False(t, disabledBound.Enabled(t.Context()))
			disabledBound.Add(t.Context(), 1)

			bound := enabledCounter.(metricx.Int64CounterBinder).Bind(
				attribute.String("route", "/orders"),
				attribute.String("zone", "east"),
			)
			bound.Add(t.Context(), 3)

			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &rm))
			require.Len(t, rm.ScopeMetrics, 1)
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
			got := rm.ScopeMetrics[0].Metrics[0]
			assert.Equal(t, "requests.enabled", got.Name)
			sum, ok := got.Data.(metricdata.Sum[int64])
			require.True(t, ok)
			require.Len(t, sum.DataPoints, 1)
			assert.Equal(t, int64(3), sum.DataPoints[0].Value)
			assert.Equal(t, attribute.NewSet(attribute.String("route", "/orders")), sum.DataPoints[0].Attributes)
		})
	}
}

func TestBoundExemplars(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithExemplarFilter(exemplar.TraceBasedFilter),
		sdkmetric.WithView(sdkmetric.NewView(
			sdkmetric.Instrument{Name: "requests"},
			sdkmetric.Stream{AttributeFilter: attribute.NewAllowKeysFilter("route")},
		)),
		sdkmetricx.WithBinding(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	bound := counter.(metricx.Int64CounterBinder).Bind(
		attribute.String("route", "/orders"),
		attribute.String("zone", "east"),
	)

	// Unsampled context produces no exemplar.
	bound.Add(t.Context(), 1)
	pts := collectSumPoints[int64](t, reader)
	require.Len(t, pts, 1)
	assert.Empty(t, pts[0].Exemplars)

	// Sampled recording context produces an exemplar with trace context and
	// view-filtered attributes.
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01, 0x02, 0x03, 0x04},
		SpanID:     trace.SpanID{0x05, 0x06, 0x07, 0x08},
		TraceFlags: trace.FlagsSampled,
	})
	sampledCtx := trace.ContextWithSpanContext(t.Context(), sc)
	bound.Add(sampledCtx, 2)

	pts = collectSumPoints[int64](t, reader)
	require.Len(t, pts, 1)
	assert.Equal(t, int64(3), pts[0].Value)
	require.Len(t, pts[0].Exemplars, 1)
	ex := pts[0].Exemplars[0]
	assert.Equal(t, int64(2), ex.Value)
	assert.Equal(t, sc.TraceID(), trace.TraceID(ex.TraceID))
	assert.Equal(t, sc.SpanID(), trace.SpanID(ex.SpanID))
	assert.Equal(t, []attribute.KeyValue{attribute.String("zone", "east")}, ex.FilteredAttributes)
}

func TestBoundConcurrentSafe(t *testing.T) {
	const (
		workers    = 8
		iterations = 100
	)
	for _, tc := range []struct {
		name     string
		selector sdkmetric.TemporalitySelector
		isDelta  bool
	}{
		{name: "Cumulative", selector: cumulativeTemporalitySelector, isDelta: false},
		{name: "Delta", selector: deltaTemporalitySelector, isDelta: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reader := sdkmetric.NewManualReader(sdkmetric.WithTemporalitySelector(tc.selector))
			provider := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(reader),
				sdkmetricx.WithBinding(),
			)
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })

			counter, err := provider.Meter("test").Int64Counter("requests")
			require.NoError(t, err)
			binder := counter.(metricx.Int64CounterBinder)
			shared := binder.Bind(attribute.String("route", "/shared"))

			var (
				workersDone atomic.Bool
				workerWG    sync.WaitGroup
				collectWG   sync.WaitGroup
				totals      = make(map[attribute.Set]int64)
			)
			collectOnce := func() {
				var rm metricdata.ResourceMetrics
				require.NoError(t, reader.Collect(t.Context(), &rm))
				if len(rm.ScopeMetrics) == 0 || len(rm.ScopeMetrics[0].Metrics) == 0 {
					return
				}
				sum := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
				if tc.isDelta {
					for _, dp := range sum.DataPoints {
						totals[dp.Attributes] += dp.Value
					}
				} else {
					clear(totals)
					for _, dp := range sum.DataPoints {
						totals[dp.Attributes] = dp.Value
					}
				}
			}

			collectWG.Go(func() {
				for !workersDone.Load() {
					collectOnce()
				}
			})
			for i := range workers {
				workerWG.Go(func() {
					local := binder.Bind(attribute.Int("worker", i))
					for range iterations {
						shared.Add(t.Context(), 1)
						local.Add(t.Context(), 1)
						counter.Add(t.Context(), 1, metric.WithAttributes(attribute.String("route", "/shared")))
						local.Add(t.Context(), 1, metric.WithAttributes(attribute.String("extra", "yes")))
					}
				})
			}
			workerWG.Wait()
			workersDone.Store(true)
			collectWG.Wait()
			collectOnce()

			want := map[attribute.Set]int64{
				attribute.NewSet(attribute.String("route", "/shared")): workers * iterations * 2,
			}
			for i := range workers {
				want[attribute.NewSet(attribute.Int("worker", i))] = iterations
				want[attribute.NewSet(
					attribute.Int("worker", i),
					attribute.String("extra", "yes"),
				)] = iterations
			}
			assert.Equal(t, want, totals)
		})
	}
}
