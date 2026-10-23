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
	"testing/synctest"
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

type finishUpDownCounter[N int64 | float64] interface {
	Add(context.Context, N, ...metric.AddOption)
	Enabled(context.Context) bool
}

func newFinishUpDownCounter[N int64 | float64](t testing.TB, meter metric.Meter, name string) finishUpDownCounter[N] {
	t.Helper()
	var counter any
	var err error
	switch any(N(0)).(type) {
	case int64:
		counter, err = meter.Int64UpDownCounter(name)
	case float64:
		counter, err = meter.Float64UpDownCounter(name)
	}
	require.NoError(t, err)
	return counter.(finishUpDownCounter[N])
}

func collectUpDownSum[N int64 | float64](t *testing.T, reader *sdkmetric.ManualReader) []metricdata.DataPoint[N] {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	if len(rm.ScopeMetrics) == 0 {
		return nil
	}
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[N])
	require.True(t, ok)
	assert.False(t, sum.IsMonotonic)
	return sum.DataPoints
}

func TestUpDownCounterFinish(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testUpDownCounterFinish[int64](t, 1) })
	t.Run("Float64", func(t *testing.T) { testUpDownCounterFinish[float64](t, 0.25) })
}

func TestUpDownCounterFinishConcurrentSafeLifetime(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testUpDownCounterFinishLifetime(t, int64(-3)) })
	t.Run("Float64", func(t *testing.T) { testUpDownCounterFinishLifetime(t, -0.75) })
}

func testUpDownCounterFinishLifetime[N int64 | float64](t *testing.T, value N) {
	t.Helper()
	for _, recreate := range []bool{false, true} {
		t.Run(fmt.Sprintf("recreate=%t", recreate), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var pause atomic.Bool
				entered := make(chan struct{})
				resume := make(chan struct{})
				filter := func(attribute.KeyValue) bool {
					if pause.CompareAndSwap(true, false) {
						close(entered)
						select {
						case <-resume:
						case <-t.Context().Done():
						}
					}
					return true
				}
				reader := sdkmetric.NewManualReader()
				provider := sdkmetric.NewMeterProvider(
					sdkmetric.WithReader(reader), sdkmetricx.WithFinish(),
					sdkmetric.WithView(sdkmetric.NewView(
						sdkmetric.Instrument{Name: "active"}, sdkmetric.Stream{AttributeFilter: filter},
					)),
				)
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
				finisher := counter.(metricx.Finisher)
				attrs := []attribute.KeyValue{attribute.String("entity", "one")}
				option := metric.WithAttributes(attrs...)
				if recreate {
					counter.Add(t.Context(), value, option)
					finisher.Finish(t.Context(), attrs...)
					require.Len(t, collectUpDownSum[N](t, reader), 1)
					require.Empty(t, collectUpDownSum[N](t, reader))
				}
				pause.Store(true)
				done := make(chan struct{})
				go func() {
					finisher.Finish(t.Context(), attrs...)
					close(done)
				}()
				t.Cleanup(func() { <-done })
				<-entered // Finish captures its timestamp before the View filter.
				// Advance the virtual clock so the new lifetime starts strictly
				// after Finish, even on platforms with coarse clock resolution.
				time.Sleep(time.Nanosecond)
				counter.Add(t.Context(), value, option)
				close(resume)
				<-done

				points := collectUpDownSum[N](t, reader)
				require.Len(t, points, 1)
				assert.Equal(t, value, points[0].Value)
				assert.False(t, points[0].Time.Before(points[0].StartTime))
				start := points[0].StartTime
				points = collectUpDownSum[N](t, reader)
				require.Len(t, points, 1, "stale Finish must not retire the new lifetime")
				assert.Equal(t, start, points[0].StartTime)
				assert.Equal(t, value, points[0].Value)
			})
		})
	}
}

func testUpDownCounterFinish[N int64 | float64](t *testing.T, step N) {
	t.Run("Activation", func(t *testing.T) {
		for _, enabled := range []bool{false, true} {
			t.Run(fmt.Sprint(enabled), func(t *testing.T) {
				options := []sdkmetric.Option{}
				if enabled {
					options = append(options, sdkmetricx.WithFinish())
				}
				provider := sdkmetric.NewMeterProvider(options...)
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				meter := provider.Meter("test")
				counter := newFinishUpDownCounter[N](t, meter, "active")
				_, ok := counter.(metricx.Finisher)
				assert.Equal(t, enabled, ok)
				assert.Same(t, counter, newFinishUpDownCounter[N](t, meter, "active"))
			})
		}
	})
	for _, temporality := range []metricdata.Temporality{metricdata.CumulativeTemporality, metricdata.DeltaTemporality} {
		t.Run(temporality.String(), func(t *testing.T) {
			reader := sdkmetric.NewManualReader(
				sdkmetric.WithTemporalitySelector(
					func(sdkmetric.InstrumentKind) metricdata.Temporality { return temporality },
				),
			)
			provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetricx.WithFinish())
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
			counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
			finisher, ok := counter.(metricx.Finisher)
			require.True(t, ok)
			attrs := []attribute.KeyValue{attribute.String("entity", "one")}
			option := metric.WithAttributes(attrs...)

			finisher.Finish(t.Context(), attrs...) // Unknown series.
			counter.Add(t.Context(), 5*step, option)
			finisher.Finish(t.Context(), attrs...)
			counter.Add(t.Context(), -2*step, option) // Cancel the pending Finish.
			points := collectUpDownSum[N](t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, N(3)*step, points[0].Value)
			start := points[0].StartTime

			counter.Add(t.Context(), -4*step, option)
			before := time.Now()
			finisher.Finish(t.Context(), attrs...)
			after := time.Now()
			finisher.Finish(t.Context(), attrs...)
			points = collectUpDownSum[N](t, reader)
			require.Len(t, points, 1)
			want := N(-1) * step
			if temporality == metricdata.DeltaTemporality {
				want = -4 * step
			}
			assert.Equal(t, want, points[0].Value)
			assert.False(t, points[0].Time.Before(before))
			assert.False(t, points[0].Time.After(after))
			if temporality == metricdata.CumulativeTemporality {
				assert.Equal(t, start, points[0].StartTime)
			}
			assert.Empty(t, collectUpDownSum[N](t, reader))

			before = time.Now()
			counter.Add(t.Context(), -3*step, option)
			points = collectUpDownSum[N](t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, N(-3)*step, points[0].Value)
			if temporality == metricdata.CumulativeTemporality {
				assert.False(t, points[0].StartTime.Before(before))
			}
			finisher.Finish(t.Context(), attrs...)
			points = collectUpDownSum[N](t, reader)
			if temporality == metricdata.DeltaTemporality {
				assert.Empty(t, points, "Finish must not synthesize an empty delta point")
			} else {
				require.Len(t, points, 1)
				assert.Equal(t, N(-3)*step, points[0].Value)
			}
			assert.Empty(t, collectUpDownSum[N](t, reader))
		})
	}
	t.Run("Fanout", testUpDownCounterFinishFanout[N])
	t.Run("Overflow", testUpDownCounterFinishOverflow[N])
	t.Run("Fallback", testUpDownCounterFinishFallback[N])
	t.Run("Concurrent", testUpDownCounterFinishConcurrent[N])
}

func testUpDownCounterFinishFanout[N int64 | float64](t *testing.T) {
	readers := []*sdkmetric.ManualReader{sdkmetric.NewManualReader(), sdkmetric.NewManualReader()}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(readers[0]), sdkmetric.WithReader(readers[1]), sdkmetricx.WithFinish(),
		sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "active"},
				sdkmetric.Stream{Name: "by_entity", AttributeFilter: attribute.NewAllowKeysFilter("entity")},
			),
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "active"},
				sdkmetric.Stream{Name: "by_zone", AttributeFilter: attribute.NewAllowKeysFilter("zone")},
			),
		),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
	finisher, ok := counter.(metricx.Finisher)
	require.True(t, ok)
	counter.Add(
		t.Context(),
		-3,
		metric.WithAttributes(attribute.String("entity", "one"), attribute.String("zone", "east")),
	)
	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("entity", "old"),
		attribute.String("entity", "one"),
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
			sum, ok := m.Data.(metricdata.Sum[N])
			require.True(t, ok)
			assert.False(t, sum.IsMonotonic)
			require.Len(t, sum.DataPoints, 1)
			point := sum.DataPoints[0]
			assert.Equal(t, N(-3), point.Value)
			want := attribute.NewSet(attribute.String("entity", "one"))
			if m.Name == "by_zone" {
				want = attribute.NewSet(attribute.String("zone", "east"))
			} else {
				assert.Equal(t, "by_entity", m.Name)
			}
			assert.Equal(t, want, point.Attributes)
			if finishTime.IsZero() {
				finishTime = point.Time
			}
			assert.Equal(t, finishTime, point.Time)
		}
		require.NoError(t, reader.Collect(t.Context(), &rm))
		assert.Empty(t, rm.ScopeMetrics)
	}
}

func testUpDownCounterFinishOverflow[N int64 | float64](t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithCardinalityLimit(2),
		sdkmetricx.WithFinish(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
	finisher, ok := counter.(metricx.Finisher)
	require.True(t, ok)
	counter.Add(t.Context(), 1)
	attrs := []attribute.KeyValue{attribute.String("entity", "overflow")}
	counter.Add(t.Context(), -2, metric.WithAttributes(attrs...))
	finisher.Finish(t.Context(), attrs...)
	finisher.Finish(t.Context(), attribute.Bool("otel.metric.overflow", true))
	finisher.Finish(t.Context()) // Only the attribute-less series.
	points := collectUpDownSum[N](t, reader)
	require.Len(t, points, 2)
	points = collectUpDownSum[N](t, reader)
	require.Len(t, points, 1)
	assert.Equal(t, attribute.NewSet(attribute.Bool("otel.metric.overflow", true)), points[0].Attributes)
	assert.Equal(t, N(-2), points[0].Value)
}

func testUpDownCounterFinishFallback[N int64 | float64](t *testing.T) {
	for _, agg := range []sdkmetric.Aggregation{sdkmetric.AggregationDrop{}, sdkmetric.AggregationExplicitBucketHistogram{}, sdkmetric.AggregationBase2ExponentialHistogram{MaxSize: 160, MaxScale: 20}} {
		t.Run(fmt.Sprintf("%T", agg), func(t *testing.T) {
			reader := sdkmetric.NewManualReader()
			provider := sdkmetric.NewMeterProvider(
				sdkmetric.WithReader(reader),
				sdkmetricx.WithFinish(),
				sdkmetric.WithView(
					sdkmetric.NewView(sdkmetric.Instrument{Name: "active"}, sdkmetric.Stream{Aggregation: agg}),
				),
			)
			t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
			counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
			finisher, ok := counter.(metricx.Finisher)
			require.True(t, ok)
			counter.Add(t.Context(), -3)
			finisher.Finish(t.Context())
			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &rm))
			require.NoError(t, reader.Collect(t.Context(), &rm))
			if _, drop := agg.(sdkmetric.AggregationDrop); drop {
				assert.False(t, counter.Enabled(t.Context()))
				assert.Empty(t, rm.ScopeMetrics)
				return
			}
			require.Len(t, rm.ScopeMetrics, 1)
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
			switch hist := rm.ScopeMetrics[0].Metrics[0].Data.(type) {
			case metricdata.Histogram[N]:
				require.Len(t, hist.DataPoints, 1)
				assert.Equal(t, uint64(1), hist.DataPoints[0].Count)
				assert.Zero(t, hist.DataPoints[0].Sum, "negative-capable instruments must omit the sum")
			case metricdata.ExponentialHistogram[N]:
				require.Len(t, hist.DataPoints, 1)
				assert.Equal(t, uint64(1), hist.DataPoints[0].Count)
				assert.Zero(t, hist.DataPoints[0].Sum)
			default:
				t.Fatalf("unexpected aggregation %T", hist)
			}
		})
	}
}

func testUpDownCounterFinishConcurrent[N int64 | float64](t *testing.T) {
	reader := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(
			func(sdkmetric.InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality },
		),
	)
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetricx.WithFinish())
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	counter := newFinishUpDownCounter[N](t, provider.Meter("test"), "active")
	finisher, ok := counter.(metricx.Finisher)
	require.True(t, ok)
	var workers sync.WaitGroup
	done := make(chan struct{})
	workers.Go(func() {
		for range 1000 {
			counter.Add(t.Context(), -1)
		}
	})
	workers.Go(func() {
		for range 1000 {
			finisher.Finish(t.Context())
		}
	})
	go func() { workers.Wait(); close(done) }()
	var total N
	for {
		for _, point := range collectUpDownSum[N](t, reader) {
			total += point.Value
		}
		select {
		case <-done:
			for _, point := range collectUpDownSum[N](t, reader) {
				total += point.Value
			}
			assert.Equal(t, N(-1000), total, "each measurement must be collected exactly once")
			return
		default:
		}
	}
}
