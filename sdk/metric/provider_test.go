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
		_, _ = fmt.Fprint(&buf, prefix, args)
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

	type metricSpec struct {
		kind InstrumentKind
		name string
	}

	type metricExpectation struct {
		name       string
		wantPoints int
	}

	type testCase struct {
		name        string
		selector    CardinalityLimitSelector
		globalLimit int
		metrics     []metricSpec
		expect      []metricExpectation
	}

	testCases := []testCase{
		{
			name: "counter uses counter-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindCounter {
					return 3, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindCounter, name: "counter-metric"},
			},
			expect: []metricExpectation{
				{name: "counter-metric", wantPoints: 3},
			},
		},
		{
			name: "histogram uses histogram-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindHistogram {
					return 4, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindHistogram, name: "histogram-metric"},
			},
			expect: []metricExpectation{
				{name: "histogram-metric", wantPoints: 4},
			},
		},
		{
			name: "gauge uses gauge-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindGauge {
					return 5, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindGauge, name: "gauge-metric"},
			},
			expect: []metricExpectation{
				{name: "gauge-metric", wantPoints: 5},
			},
		},
		{
			name: "up down counter uses updowncounter-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindUpDownCounter {
					return 2, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindUpDownCounter, name: "updowncounter-metric"},
			},
			expect: []metricExpectation{
				{name: "updowncounter-metric", wantPoints: 2},
			},
		},
		{
			name: "instrument without specific limit falls back to global limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindCounter {
					return 3, false
				}
				return 0, true // fall back to global limit
			},
			globalLimit: 6,
			metrics: []metricSpec{
				{kind: InstrumentKindCounter, name: "counter-metric"},
				{kind: InstrumentKindHistogram, name: "histogram-metric"},
			},
			expect: []metricExpectation{
				{name: "counter-metric", wantPoints: 3},
				{name: "histogram-metric", wantPoints: 6},
			},
		},
		{
			name: "observable counter uses observable-counter-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindObservableCounter {
					return 4, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindObservableCounter, name: "observable-counter-metric"},
			},
			expect: []metricExpectation{
				{name: "observable-counter-metric", wantPoints: 4},
			},
		},
		{
			name: "observable gauge uses observable-gauge-specific limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindObservableGauge {
					return 5, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindObservableGauge, name: "observable-gauge-metric"},
			},
			expect: []metricExpectation{
				{name: "observable-gauge-metric", wantPoints: 5},
			},
		},
		{
			name: "observable up down counter uses limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindObservableUpDownCounter {
					return 3, false
				}
				return 0, true
			},
			globalLimit: 8,
			metrics: []metricSpec{
				{kind: InstrumentKindObservableUpDownCounter, name: "observable-updowncounter-metric"},
			},
			expect: []metricExpectation{
				{name: "observable-updowncounter-metric", wantPoints: 3},
			},
		},
		{
			name: "selector can set specific kind to unlimited while others use global limit",
			selector: func(kind InstrumentKind) (int, bool) {
				if kind == InstrumentKindCounter {
					return 0, false // unlimited for counter only
				}
				return 0, true // fallback to global limit
			},
			globalLimit: 3,
			metrics: []metricSpec{
				{kind: InstrumentKindCounter, name: "counter-metric"},
				{kind: InstrumentKindHistogram, name: "histogram-metric"},
			},
			expect: []metricExpectation{
				{name: "counter-metric", wantPoints: uniqueAttributesCount},
				{name: "histogram-metric", wantPoints: 3},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := NewManualReader(
				WithCardinalityLimitSelector(tc.selector),
			)
			mp := NewMeterProvider(
				WithReader(reader),
				WithCardinalityLimit(tc.globalLimit),
			)

			meter := mp.Meter("test-meter")

			// Create instruments and record/observe data based on metric specs.
			for _, ms := range tc.metrics {
				switch ms.kind {
				case InstrumentKindCounter:
					counter, err := meter.Int64Counter(ms.name)
					require.NoError(t, err)
					for i := range uniqueAttributesCount {
						counter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
					}
				case InstrumentKindHistogram:
					histogram, err := meter.Int64Histogram(ms.name)
					require.NoError(t, err)
					for i := range uniqueAttributesCount {
						histogram.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
					}
				case InstrumentKindGauge:
					gauge, err := meter.Int64Gauge(ms.name)
					require.NoError(t, err)
					for i := range uniqueAttributesCount {
						gauge.Record(t.Context(), int64(i), api.WithAttributes(attribute.Int("key", i)))
					}
				case InstrumentKindUpDownCounter:
					upDownCounter, err := meter.Int64UpDownCounter(ms.name)
					require.NoError(t, err)
					for i := range uniqueAttributesCount {
						upDownCounter.Add(t.Context(), 1, api.WithAttributes(attribute.Int("key", i)))
					}
				case InstrumentKindObservableCounter:
					obs, err := meter.Int64ObservableCounter(
						ms.name,
						api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
							for i := range uniqueAttributesCount {
								o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
							}
							return nil
						}),
					)
					require.NoError(t, err)
					require.NotNil(t, obs)
				case InstrumentKindObservableGauge:
					obs, err := meter.Int64ObservableGauge(
						ms.name,
						api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
							for i := range uniqueAttributesCount {
								o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
							}
							return nil
						}),
					)
					require.NoError(t, err)
					require.NotNil(t, obs)
				case InstrumentKindObservableUpDownCounter:
					obs, err := meter.Int64ObservableUpDownCounter(
						ms.name,
						api.WithInt64Callback(func(_ context.Context, o api.Int64Observer) error {
							for i := range uniqueAttributesCount {
								o.Observe(int64(i), api.WithAttributes(attribute.Int("key", i)))
							}
							return nil
						}),
					)
					require.NoError(t, err)
					require.NotNil(t, obs)
				default:
					t.Fatalf("unexpected instrument kind %v", ms.kind)
				}
			}

			var rm metricdata.ResourceMetrics
			err := reader.Collect(t.Context(), &rm)
			require.NoError(t, err)

			require.NotEmpty(t, rm.ScopeMetrics)
			metrics := rm.ScopeMetrics[0].Metrics

			for _, exp := range tc.expect {
				var found bool
				for _, m := range metrics {
					if m.Name != exp.name {
						continue
					}
					found = true
					switch data := m.Data.(type) {
					case metricdata.Sum[int64]:
						assert.Len(t, data.DataPoints, exp.wantPoints, tc.name)
					case metricdata.Histogram[int64]:
						assert.Len(t, data.DataPoints, exp.wantPoints, tc.name)
					case metricdata.Gauge[int64]:
						assert.Len(t, data.DataPoints, exp.wantPoints, tc.name)
					default:
						t.Fatalf("unexpected data type %T", data)
					}
				}
				require.Truef(t, found, "metric %q not found in results", exp.name)
			}
		})
	}
}
