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

type finishGauge[N int64 | float64] interface {
	Record(context.Context, N, ...metric.RecordOption)
	Enabled(context.Context) bool
}

func newFinishGauge[N int64 | float64](t testing.TB, meter metric.Meter, name string) finishGauge[N] {
	t.Helper()
	var gauge any
	var err error
	switch any(N(0)).(type) {
	case int64:
		gauge, err = meter.Int64Gauge(name)
	case float64:
		gauge, err = meter.Float64Gauge(name)
	}
	require.NoError(t, err)
	return gauge.(finishGauge[N])
}

func collectGaugePoints[N int64 | float64](t *testing.T, reader *sdkmetric.ManualReader) []metricdata.DataPoint[N] {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	if len(rm.ScopeMetrics) == 0 {
		return nil
	}
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	data, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[N])
	require.True(t, ok)
	return data.DataPoints
}

func TestGaugeFinish(t *testing.T) {
	t.Run("Int64", func(t *testing.T) { testGaugeFinish[int64](t, 1) })
	t.Run("Float64", func(t *testing.T) { testGaugeFinish[float64](t, 0.25) })
}

func testGaugeFinish[N int64 | float64](t *testing.T, step N) {
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
				gauge := newFinishGauge[N](t, meter, "active")
				_, ok := gauge.(metricx.Finisher)
				assert.Equal(t, enabled, ok)
				assert.Same(t, gauge, newFinishGauge[N](t, meter, "active"))
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
			gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
			finisher, ok := gauge.(metricx.Finisher)
			require.True(t, ok)
			attrs := []attribute.KeyValue{attribute.String("entity", "one")}
			option := metric.WithAttributes(attrs...)

			finisher.Finish(t.Context(), attrs...) // Unknown series.
			gauge.Record(t.Context(), 5*step, option)
			finisher.Finish(t.Context(), attrs...)
			gauge.Record(t.Context(), -2*step, option) // Cancel the pending Finish.
			points := collectGaugePoints[N](t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, N(-2)*step, points[0].Value)
			start := points[0].StartTime

			gauge.Record(t.Context(), -4*step, option)
			before := time.Now()
			finisher.Finish(t.Context(), attrs...)
			after := time.Now()
			finisher.Finish(t.Context(), attrs...)
			points = collectGaugePoints[N](t, reader)
			require.Len(t, points, 1)
			want := N(-4) * step
			assert.Equal(t, want, points[0].Value)
			assert.False(t, points[0].Time.Before(before))
			assert.False(t, points[0].Time.After(after))
			if temporality == metricdata.CumulativeTemporality {
				assert.Equal(t, start, points[0].StartTime)
			}
			assert.Empty(t, collectGaugePoints[N](t, reader))

			before = time.Now()
			gauge.Record(t.Context(), -3*step, option)
			points = collectGaugePoints[N](t, reader)
			require.Len(t, points, 1)
			assert.Equal(t, N(-3)*step, points[0].Value)
			if temporality == metricdata.CumulativeTemporality {
				assert.False(t, points[0].StartTime.Before(before))
			}
			finisher.Finish(t.Context(), attrs...)
			points = collectGaugePoints[N](t, reader)
			if temporality == metricdata.DeltaTemporality {
				assert.Empty(t, points, "Finish must not synthesize an empty delta point")
			} else {
				require.Len(t, points, 1)
				assert.Equal(t, N(-3)*step, points[0].Value)
			}
			assert.Empty(t, collectGaugePoints[N](t, reader))
		})
	}
	t.Run("StaleFinish", testGaugeStaleFinish[N])
	t.Run("Fanout", testGaugeFinishFanout[N])
	t.Run("Overflow", testGaugeFinishOverflow[N])
	t.Run("Fallback", testGaugeFinishFallback[N])
	t.Run("Concurrent", testGaugeFinishConcurrent[N])
}

func testGaugeStaleFinish[N int64 | float64](t *testing.T) {
	for _, recreate := range []bool{false, true} {
		t.Run(fmt.Sprint("recreate=", recreate), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				entered, release := make(chan struct{}), make(chan struct{})
				var pause atomic.Bool
				reader := sdkmetric.NewManualReader()
				provider := sdkmetric.NewMeterProvider(
					sdkmetric.WithReader(reader), sdkmetricx.WithFinish(),
					sdkmetric.WithView(sdkmetric.NewView(
						sdkmetric.Instrument{Name: "active"},
						sdkmetric.Stream{AttributeFilter: func(attribute.KeyValue) bool {
							if pause.CompareAndSwap(true, false) {
								close(entered)
								<-release
							}
							return true
						}},
					)),
				)
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
				finisher := gauge.(metricx.Finisher)
				attrs := []attribute.KeyValue{attribute.String("entity", "one")}
				option := metric.WithAttributes(attrs...)
				if recreate {
					gauge.Record(t.Context(), 1, option)
					finisher.Finish(t.Context(), attrs...)
					require.Len(t, collectGaugePoints[N](t, reader), 1)
					require.Empty(t, collectGaugePoints[N](t, reader))
				}

				pause.Store(true)
				done := make(chan struct{})
				go func() {
					defer close(done)
					finisher.Finish(t.Context(), attrs...)
				}()
				var once sync.Once
				resume := func() { once.Do(func() { close(release) }) }
				defer func() { resume(); <-done }()
				select {
				case <-entered:
				case <-time.After(time.Second):
					t.Fatal("Finish did not reach the View filter")
				}
				// Advance the test clock so the new lifetime starts strictly after Finish.
				time.Sleep(time.Nanosecond)
				gauge.Record(t.Context(), -3, option)
				resume()
				<-done

				for range 2 {
					points := collectGaugePoints[N](t, reader)
					require.Len(t, points, 1, "stale Finish must not retire the new lifetime")
					assert.Equal(t, N(-3), points[0].Value)
					assert.False(t, points[0].Time.Before(points[0].StartTime))
				}
				finisher.Finish(t.Context(), attrs...)
				require.Len(t, collectGaugePoints[N](t, reader), 1)
				assert.Empty(t, collectGaugePoints[N](t, reader))
			})
		})
	}
}

func testGaugeFinishFanout[N int64 | float64](t *testing.T) {
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
	gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
	finisher, ok := gauge.(metricx.Finisher)
	require.True(t, ok)
	gauge.Record(
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
			data, ok := m.Data.(metricdata.Gauge[N])
			require.True(t, ok)
			require.Len(t, data.DataPoints, 1)
			point := data.DataPoints[0]
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

func testGaugeFinishOverflow[N int64 | float64](t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithCardinalityLimit(2),
		sdkmetricx.WithFinish(),
	)
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
	finisher, ok := gauge.(metricx.Finisher)
	require.True(t, ok)
	gauge.Record(t.Context(), 1)
	attrs := []attribute.KeyValue{attribute.String("entity", "overflow")}
	gauge.Record(t.Context(), -2, metric.WithAttributes(attrs...))
	finisher.Finish(t.Context(), attrs...)
	finisher.Finish(t.Context(), attribute.Bool("otel.metric.overflow", true))
	finisher.Finish(t.Context()) // Only the attribute-less series.
	points := collectGaugePoints[N](t, reader)
	require.Len(t, points, 2)
	points = collectGaugePoints[N](t, reader)
	require.Len(t, points, 1)
	assert.Equal(t, attribute.NewSet(attribute.Bool("otel.metric.overflow", true)), points[0].Attributes)
	assert.Equal(t, N(-2), points[0].Value)
}

func testGaugeFinishFallback[N int64 | float64](t *testing.T) {
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
			gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
			finisher, ok := gauge.(metricx.Finisher)
			require.True(t, ok)
			gauge.Record(t.Context(), -3)
			finisher.Finish(t.Context())
			var rm metricdata.ResourceMetrics
			require.NoError(t, reader.Collect(t.Context(), &rm))
			require.NoError(t, reader.Collect(t.Context(), &rm))
			if _, drop := agg.(sdkmetric.AggregationDrop); drop {
				assert.False(t, gauge.Enabled(t.Context()))
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

func testGaugeFinishConcurrent[N int64 | float64](t *testing.T) {
	reader := sdkmetric.NewManualReader(
		sdkmetric.WithTemporalitySelector(
			func(sdkmetric.InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality },
		),
	)
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader), sdkmetricx.WithFinish())
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
	gauge := newFinishGauge[N](t, provider.Meter("test"), "active")
	finisher, ok := gauge.(metricx.Finisher)
	require.True(t, ok)
	var workers sync.WaitGroup
	done := make(chan struct{})
	workers.Go(func() {
		for i := 1; i <= 1000; i++ {
			gauge.Record(t.Context(), -N(i))
		}
	})
	workers.Go(func() {
		for range 1000 {
			finisher.Finish(t.Context())
		}
	})
	go func() { workers.Wait(); close(done) }()
	var latest N
	collect := func() {
		for _, point := range collectGaugePoints[N](t, reader) {
			assert.False(t, point.Time.Before(point.StartTime))
			assert.Less(t, point.Value, N(0))
			assert.GreaterOrEqual(t, point.Value, N(-1000))
			if point.Value < latest {
				latest = point.Value
			}
		}
	}
	for {
		collect()
		select {
		case <-done:
			finisher.Finish(t.Context())
			collect()
			assert.Equal(t, N(-1000), latest, "the latest measurement must be collected")
			assert.Empty(t, collectGaugePoints[N](t, reader))
			return
		default:
		}
	}
}
