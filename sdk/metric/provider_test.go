// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMeterConcurrentSafe(*testing.T) {
	const name = "TestMeterConcurrentSafe meter"
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Meter(name)
	}()

	_ = mp.Meter(name)
	<-done
}

func TestForceFlushConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.ForceFlush(t.Context())
	}()

	_ = mp.ForceFlush(t.Context())
	<-done
}

func TestShutdownConcurrentSafe(t *testing.T) {
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Shutdown(t.Context())
	}()

	_ = mp.Shutdown(t.Context())
	<-done
}

func TestMeterAndShutdownConcurrentSafe(t *testing.T) {
	const name = "TestMeterAndShutdownConcurrentSafe meter"
	mp := NewMeterProvider()

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Shutdown(t.Context())
	}()

	_ = mp.Meter(name)
	<-done
}

func TestMeterDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Meter("") })
}

func TestForceFlushDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.ForceFlush(t.Context()) })
}

func TestShutdownDoesNotPanicForEmptyMeterProvider(t *testing.T) {
	mp := MeterProvider{}
	assert.NotPanics(t, func() { _ = mp.Shutdown(t.Context()) })
}

func TestMeterProviderReturnsSameMeter(t *testing.T) {
	mp := MeterProvider{}
	mtr := mp.Meter("")

	assert.Same(t, mtr, mp.Meter(""))
	assert.NotSame(t, mtr, mp.Meter("diff"))
	assert.NotSame(t, mtr, mp.Meter("", api.WithInstrumentationAttributes(attribute.String("k", "v"))))
}

func TestEmptyMeterName(t *testing.T) {
	var buf strings.Builder
	warnLevel := 1
	l := funcr.New(func(prefix, args string) {
		_, _ = buf.WriteString(fmt.Sprint(prefix, args))
	}, funcr.Options{Verbosity: warnLevel})
	otel.SetLogger(l)
	mp := NewMeterProvider()

	mp.Meter("")

	assert.Contains(t, buf.String(), `"level"=1 "msg"="Invalid Meter name." "name"=""`)
}

func TestMeterProviderReturnsNoopMeterAfterShutdown(t *testing.T) {
	mp := NewMeterProvider()

	m := mp.Meter("")
	_, ok := m.(noop.Meter)
	assert.False(t, ok, "Meter from running MeterProvider is NoOp")

	require.NoError(t, mp.Shutdown(t.Context()))

	m = mp.Meter("")
	_, ok = m.(noop.Meter)
	assert.Truef(t, ok, "Meter from shutdown MeterProvider is not NoOp: %T", m)
}

func TestMeterProviderMixingOnRegisterErrors(t *testing.T) {
	otel.SetLogger(testr.New(t))

	rdr0 := NewManualReader()
	mp0 := NewMeterProvider(WithReader(rdr0))

	rdr1 := NewManualReader()
	mp1 := NewMeterProvider(WithReader(rdr1))

	// Meters with the same scope but different MeterProviders.
	m0 := mp0.Meter("TestMeterProviderMixingOnRegisterErrors")
	m1 := mp1.Meter("TestMeterProviderMixingOnRegisterErrors")

	m0Gauge, err := m0.Float64ObservableGauge("float64Gauge")
	require.NoError(t, err)

	m1Gauge, err := m1.Int64ObservableGauge("int64Gauge")
	require.NoError(t, err)

	_, err = m0.RegisterCallback(
		func(_ context.Context, o api.Observer) error {
			o.ObserveFloat64(m0Gauge, 2)
			// Observe an instrument from a different MeterProvider.
			o.ObserveInt64(m1Gauge, 1)

			return nil
		},
		m0Gauge, m1Gauge,
	)
	assert.Error(
		t,
		err,
		"Instrument registered with Meter from different MeterProvider",
	)

	var data metricdata.ResourceMetrics
	_ = rdr0.Collect(t.Context(), &data)
	// Only the metrics from mp0 should be produced.
	assert.Len(t, data.ScopeMetrics, 1)

	err = rdr1.Collect(t.Context(), &data)
	assert.NoError(t, err, "Errored when collect should be a noop")
	assert.Empty(
		t, data.ScopeMetrics,
		"Metrics produced for instrument collected by different MeterProvider",
	)
}

func TestMeterProviderCardinalityLimit(t *testing.T) {
	const uniqueAttributesCount = 10

	tests := []struct {
		name           string
		options        []Option
		wantDataPoints int
	}{
		{
			name:           "no limit (default)",
			options:        nil,
			wantDataPoints: uniqueAttributesCount,
		},
		{
			name:           "no limit (limit=0)",
			options:        []Option{WithCardinalityLimit(0)},
			wantDataPoints: uniqueAttributesCount,
		},
		{
			name:           "no limit (negative)",
			options:        []Option{WithCardinalityLimit(-5)},
			wantDataPoints: uniqueAttributesCount,
		},
		{
			name:           "limit=5",
			options:        []Option{WithCardinalityLimit(5)},
			wantDataPoints: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewManualReader()

			opts := append(tt.options, WithReader(reader))
			mp := NewMeterProvider(opts...)

			meter := mp.Meter("test-meter")
			counter, err := meter.Int64Counter("metric")
			require.NoError(t, err, "failed to create counter")

			for i := range uniqueAttributesCount {
				counter.Add(
					t.Context(),
					1,
					api.WithAttributes(attribute.Int("key", i)),
				)
			}

			var rm metricdata.ResourceMetrics
			err = reader.Collect(t.Context(), &rm)
			require.NoError(t, err, "failed to collect metrics")

			require.Len(t, rm.ScopeMetrics, 1, "expected 1 ScopeMetrics")
			require.Len(t, rm.ScopeMetrics[0].Metrics, 1, "expected 1 Metric")

			data := rm.ScopeMetrics[0].Metrics[0].Data
			require.IsType(t, metricdata.Sum[int64]{}, data, "expected metricdata.Sum[int64]")

			sumData := data.(metricdata.Sum[int64])
			assert.Len(
				t,
				sumData.DataPoints,
				tt.wantDataPoints,
				"unexpected number of data points",
			)
		})
	}
}

func TestMeterProviderPerInstrumentCardinalityLimits(t *testing.T) {
	const uniqueAttributesCount = 10

	t.Run("counter uses counter-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8), // global limit
			WithCounterCardinalityLimit(3),
		)

		meter := mp.Meter("test-meter")
		counter, err := meter.Int64Counter("counter-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			counter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		sumData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Len(t, sumData.DataPoints, 3, "counter should use counter-specific limit of 3")
	})

	t.Run("histogram uses histogram-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithHistogramCardinalityLimit(4),
		)

		meter := mp.Meter("test-meter")
		histogram, err := meter.Int64Histogram("histogram-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			histogram.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		histData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[int64])
		assert.Len(t, histData.DataPoints, 4, "histogram should use histogram-specific limit of 4")
	})

	t.Run("gauge uses gauge-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithGaugeCardinalityLimit(5),
		)

		meter := mp.Meter("test-meter")
		gauge, err := meter.Int64Gauge("gauge-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			gauge.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		gaugeData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[int64])
		assert.Len(t, gaugeData.DataPoints, 5, "gauge should use gauge-specific limit of 5")
	})

	t.Run("up down counter uses updowncounter-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithUpDownCounterCardinalityLimit(2),
		)

		meter := mp.Meter("test-meter")
		upDownCounter, err := meter.Int64UpDownCounter("updowncounter-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			upDownCounter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		sumData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Len(t, sumData.DataPoints, 2, "up down counter should use updowncounter-specific limit of 2")
	})

	t.Run("instrument without specific limit falls back to global limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(6),
			WithCounterCardinalityLimit(3), // only counter has specific limit
		)

		meter := mp.Meter("test-meter")

		// Counter should use its specific limit
		counter, err := meter.Int64Counter("counter-metric")
		require.NoError(t, err)

		// Histogram should fall back to global limit
		histogram, err := meter.Int64Histogram("histogram-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			counter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
			histogram.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

		for _, m := range rm.ScopeMetrics[0].Metrics {
			switch m.Name {
			case "counter-metric":
				sumData := m.Data.(metricdata.Sum[int64])
				assert.Len(t, sumData.DataPoints, 3, "counter should use counter-specific limit of 3")
			case "histogram-metric":
				histData := m.Data.(metricdata.Histogram[int64])
				assert.Len(t, histData.DataPoints, 6, "histogram should fall back to global limit of 6")
			}
		}
	})

	t.Run("zero per-instrument limit disables limit for that instrument", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(3),        // global limit
			WithCounterCardinalityLimit(0), // explicitly disable limit for counters
		)

		meter := mp.Meter("test-meter")

		// Counter should have no limit (0 means disabled)
		counter, err := meter.Int64Counter("counter-metric")
		require.NoError(t, err)

		// Histogram should use global limit
		histogram, err := meter.Int64Histogram("histogram-metric")
		require.NoError(t, err)

		for i := range uniqueAttributesCount {
			counter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
			histogram.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
		}

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

		for _, m := range rm.ScopeMetrics[0].Metrics {
			switch m.Name {
			case "counter-metric":
				sumData := m.Data.(metricdata.Sum[int64])
				assert.Len(t, sumData.DataPoints, uniqueAttributesCount, "counter should have no limit (0 disables)")
			case "histogram-metric":
				histData := m.Data.(metricdata.Histogram[int64])
				assert.Len(t, histData.DataPoints, 3, "histogram should use global limit of 3")
			}
		}
	})

	t.Run("observable counter uses observable-counter-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithObservableCounterCardinalityLimit(4),
		)

		meter := mp.Meter("test-meter")
		observableCounter, err := meter.Int64ObservableCounter(
			"observable-counter-metric",
			api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
				for i := range uniqueAttributesCount {
					o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
				}
				return nil
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, observableCounter)

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		sumData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Len(t, sumData.DataPoints, 4, "observable counter should use observable-counter-specific limit of 4")
	})

	t.Run("observable gauge uses observable-gauge-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithObservableGaugeCardinalityLimit(5),
		)

		meter := mp.Meter("test-meter")
		observableGauge, err := meter.Int64ObservableGauge(
			"observable-gauge-metric",
			api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
				for i := range uniqueAttributesCount {
					o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
				}
				return nil
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, observableGauge)

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		gaugeData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Gauge[int64])
		assert.Len(t, gaugeData.DataPoints, 5, "observable gauge should use observable-gauge-specific limit of 5")
	})

	t.Run("observable up down counter uses observable-updowncounter-specific limit", func(t *testing.T) {
		reader := NewManualReader()
		mp := NewMeterProvider(
			WithReader(reader),
			WithCardinalityLimit(8),
			WithObservableUpDownCounterCardinalityLimit(3),
		)

		meter := mp.Meter("test-meter")
		observableUpDownCounter, err := meter.Int64ObservableUpDownCounter(
			"observable-updowncounter-metric",
			api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
				for i := range uniqueAttributesCount {
					o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
				}
				return nil
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, observableUpDownCounter)

		var rm metricdata.ResourceMetrics
		err = reader.Collect(t.Context(), &rm)
		require.NoError(t, err)

		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

		sumData := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		assert.Len(
			t,
			sumData.DataPoints,
			3,
			"observable up down counter should use observable-updowncounter-specific limit of 3",
		)
	})
}
