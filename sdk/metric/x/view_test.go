// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/metric/x"
)

func TestNewView(t *testing.T) {
	inst := sdkmetric.Instrument{
		Name:        "requests",
		Description: "total requests",
		Kind:        sdkmetric.InstrumentKindCounter,
		Unit:        "{request}",
		Scope: instrumentation.Scope{
			Name:      "scope",
			Version:   "v1.0.0",
			SchemaURL: "https://opentelemetry.io/schemas/1.37.0",
		},
	}

	t.Run("MatchAndReplace", func(t *testing.T) {
		enabled := true
		v := x.NewView(
			sdkmetric.Instrument{Name: "req*"},
			x.Stream{
				Description: "overridden",
				Unit:        "1",
				Aggregation: sdkmetric.AggregationSum{},
				Enabled:     &enabled,
			},
		)
		// Mutating the caller's bool after NewView must not affect the view.
		enabled = false

		got, match := v(inst)
		require.True(t, match)
		require.NotNil(t, got.Enabled)
		assert.True(t, *got.Enabled)
		assert.Equal(t, "requests", got.Name)
		assert.Equal(t, "overridden", got.Description)
		assert.Equal(t, "1", got.Unit)
		assert.Equal(t, sdkmetric.AggregationSum{}, got.Aggregation)

		// Mutating the returned Stream.Enabled must not affect subsequent calls.
		*got.Enabled = false
		got2, match2 := v(inst)
		require.True(t, match2)
		require.NotNil(t, got2.Enabled)
		assert.True(t, *got2.Enabled)
	})

	t.Run("NoMatch", func(t *testing.T) {
		v := x.NewView(
			sdkmetric.Instrument{Name: "other"},
			x.Stream{Enabled: new(true)},
		)
		got, match := v(inst)
		assert.False(t, match)
		assert.Empty(t, got)
	})

	t.Run("EmptyCriteriaDropped", func(t *testing.T) {
		v := x.NewView(sdkmetric.Instrument{}, x.Stream{Enabled: new(true)})
		_, match := v(inst)
		assert.False(t, match)
	})
}

func TestOptInAllInstrumentKinds(t *testing.T) {
	type instrumentCase struct {
		name     string
		instName string
		isSync   bool
		record   func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (enabled bool, callbackCalled *atomic.Bool)
		wantData metricdata.Aggregation
	}

	cases := []instrumentCase{
		{
			name:     "Int64Counter",
			instName: "int64.counter",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Int64CounterOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Int64Counter("int64.counter", cOpts...)
				require.NoError(t, err)
				inst.Add(t.Context(), 5)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: 5}},
			},
		},
		{
			name:     "Float64Counter",
			instName: "float64.counter",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Float64CounterOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Float64Counter("float64.counter", cOpts...)
				require.NoError(t, err)
				inst.Add(t.Context(), 5.5)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Sum[float64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints:  []metricdata.DataPoint[float64]{{Value: 5.5}},
			},
		},
		{
			name:     "Int64UpDownCounter",
			instName: "int64.updowncounter",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Int64UpDownCounterOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Int64UpDownCounter("int64.updowncounter", cOpts...)
				require.NoError(t, err)
				inst.Add(t.Context(), -3)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: -3}},
			},
		},
		{
			name:     "Float64UpDownCounter",
			instName: "float64.updowncounter",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Float64UpDownCounterOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Float64UpDownCounter("float64.updowncounter", cOpts...)
				require.NoError(t, err)
				inst.Add(t.Context(), -3.5)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Sum[float64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[float64]{{Value: -3.5}},
			},
		},
		{
			name:     "Int64Gauge",
			instName: "int64.gauge",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Int64GaugeOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Int64Gauge("int64.gauge", cOpts...)
				require.NoError(t, err)
				inst.Record(t.Context(), 42)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Gauge[int64]{
				DataPoints: []metricdata.DataPoint[int64]{{Value: 42}},
			},
		},
		{
			name:     "Float64Gauge",
			instName: "float64.gauge",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Float64GaugeOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Float64Gauge("float64.gauge", cOpts...)
				require.NoError(t, err)
				inst.Record(t.Context(), 42.5)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Gauge[float64]{
				DataPoints: []metricdata.DataPoint[float64]{{Value: 42.5}},
			},
		},
		{
			name:     "Int64Histogram",
			instName: "int64.histogram",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Int64HistogramOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Int64Histogram("int64.histogram", cOpts...)
				require.NoError(t, err)
				inst.Record(t.Context(), 7)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Histogram[int64]{
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.HistogramDataPoint[int64]{
					{
						Count: 1,
						Bounds: []float64{
							0,
							5,
							10,
							25,
							50,
							75,
							100,
							250,
							500,
							750,
							1000,
							2500,
							5000,
							7500,
							10000,
						},
						BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						Min:          metricdata.NewExtrema[int64](7),
						Max:          metricdata.NewExtrema[int64](7),
						Sum:          7,
					},
				},
			},
		},
		{
			name:     "Float64Histogram",
			instName: "float64.histogram",
			isSync:   true,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				var cOpts []metric.Float64HistogramOption
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				inst, err := m.Float64Histogram("float64.histogram", cOpts...)
				require.NoError(t, err)
				inst.Record(t.Context(), 7.0)
				return inst.Enabled(t.Context()), nil
			},
			wantData: metricdata.Histogram[float64]{
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.HistogramDataPoint[float64]{
					{
						Count: 1,
						Bounds: []float64{
							0,
							5,
							10,
							25,
							50,
							75,
							100,
							250,
							500,
							750,
							1000,
							2500,
							5000,
							7500,
							10000,
						},
						BucketCounts: []uint64{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						Min:          metricdata.NewExtrema[float64](7.0),
						Max:          metricdata.NewExtrema[float64](7.0),
						Sum:          7.0,
					},
				},
			},
		},
		{
			name:     "Int64ObservableCounter",
			instName: "int64.observable.counter",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Int64ObservableCounterOption{
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						called.Store(true)
						o.Observe(9)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Int64ObservableCounter("int64.observable.counter", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: 9}},
			},
		},
		{
			name:     "Float64ObservableCounter",
			instName: "float64.observable.counter",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Float64ObservableCounterOption{
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						called.Store(true)
						o.Observe(9.5)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Float64ObservableCounter("float64.observable.counter", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Sum[float64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: true,
				DataPoints:  []metricdata.DataPoint[float64]{{Value: 9.5}},
			},
		},
		{
			name:     "Int64ObservableUpDownCounter",
			instName: "int64.observable.updowncounter",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Int64ObservableUpDownCounterOption{
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						called.Store(true)
						o.Observe(-4)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Int64ObservableUpDownCounter("int64.observable.updowncounter", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Sum[int64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[int64]{{Value: -4}},
			},
		},
		{
			name:     "Float64ObservableUpDownCounter",
			instName: "float64.observable.updowncounter",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Float64ObservableUpDownCounterOption{
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						called.Store(true)
						o.Observe(-4.5)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Float64ObservableUpDownCounter("float64.observable.updowncounter", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Sum[float64]{
				Temporality: metricdata.CumulativeTemporality,
				IsMonotonic: false,
				DataPoints:  []metricdata.DataPoint[float64]{{Value: -4.5}},
			},
		},
		{
			name:     "Int64ObservableGauge",
			instName: "int64.observable.gauge",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Int64ObservableGaugeOption{
					metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
						called.Store(true)
						o.Observe(12)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Int64ObservableGauge("int64.observable.gauge", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Gauge[int64]{
				DataPoints: []metricdata.DataPoint[int64]{{Value: 12}},
			},
		},
		{
			name:     "Float64ObservableGauge",
			instName: "float64.observable.gauge",
			isSync:   false,
			record: func(t *testing.T, m metric.Meter, opts ...metric.InstrumentOption) (bool, *atomic.Bool) {
				called := &atomic.Bool{}
				cOpts := []metric.Float64ObservableGaugeOption{
					metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
						called.Store(true)
						o.Observe(12.5)
						return nil
					}),
				}
				for _, o := range opts {
					cOpts = append(cOpts, o)
				}
				_, err := m.Float64ObservableGauge("float64.observable.gauge", cOpts...)
				require.NoError(t, err)
				return false, called
			},
			wantData: metricdata.Gauge[float64]{
				DataPoints: []metricdata.DataPoint[float64]{{Value: 12.5}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("OptIn_DisabledByDefault", func(t *testing.T) {
				rdr := sdkmetric.NewManualReader()
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr))
				m := mp.Meter("test")

				enabled, cbCalled := tc.record(t, m, metricx.WithOptIn())
				if tc.isSync {
					assert.False(t, enabled, "OptIn instrument should be disabled by default")
				}

				var rm metricdata.ResourceMetrics
				require.NoError(t, rdr.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics, "OptIn instrument should not export metrics by default")
				if cbCalled != nil {
					assert.False(t, cbCalled.Load(), "OptIn observable callback should not be invoked when disabled")
				}
			})

			t.Run("OptIn_MatchingViewWithNilEnabledStillDrops", func(t *testing.T) {
				rdr := sdkmetric.NewManualReader()
				// Even when a View matches and sets Description or Aggregation,
				// if Enabled is nil (unset), an OptIn instrument remains disabled.
				view := x.NewView(
					sdkmetric.Instrument{Name: tc.instName},
					x.Stream{
						Description: "custom description",
						Aggregation: sdkmetric.AggregationDefault{},
					},
				)
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), x.WithView(view))
				m := mp.Meter("test")

				enabled, cbCalled := tc.record(t, m, metricx.WithOptIn())
				if tc.isSync {
					assert.False(t, enabled, "OptIn instrument with Enabled=nil view should remain disabled")
				}

				var rm metricdata.ResourceMetrics
				require.NoError(t, rdr.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics)
				if cbCalled != nil {
					assert.False(t, cbCalled.Load())
				}
			})

			t.Run("OptIn_EnabledByView", func(t *testing.T) {
				rdr := sdkmetric.NewManualReader()
				view := x.NewView(
					sdkmetric.Instrument{Name: tc.instName},
					x.Stream{Enabled: new(true)},
				)
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), x.WithView(view))
				m := mp.Meter("test")

				enabled, cbCalled := tc.record(t, m, metricx.WithOptIn())
				if tc.isSync {
					assert.True(t, enabled, "OptIn instrument with Enabled=true view should be enabled")
				}

				var rm metricdata.ResourceMetrics
				require.NoError(t, rdr.Collect(t.Context(), &rm))
				require.Len(t, rm.ScopeMetrics, 1)
				require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
				metricdatatest.AssertEqual(
					t,
					metricdata.Metrics{
						Name: tc.instName,
						Data: tc.wantData,
					},
					rm.ScopeMetrics[0].Metrics[0],
					metricdatatest.IgnoreTimestamp(),
					metricdatatest.IgnoreExemplars(),
				)
				if cbCalled != nil {
					assert.True(t, cbCalled.Load())
				}
			})

			t.Run("DefaultInstrument_DisabledByView", func(t *testing.T) {
				rdr := sdkmetric.NewManualReader()
				// Even if Aggregation is explicitly set on the View, Enabled=false drops the stream.
				view := x.NewView(
					sdkmetric.Instrument{Name: tc.instName},
					x.Stream{
						Aggregation: sdkmetric.AggregationDefault{},
						Enabled:     new(false),
					},
				)
				mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), x.WithView(view))
				m := mp.Meter("test")

				enabled, cbCalled := tc.record(t, m)
				if tc.isSync {
					assert.False(t, enabled, "instrument with Enabled=false view should be disabled")
				}

				var rm metricdata.ResourceMetrics
				require.NoError(t, rdr.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics)
				if cbCalled != nil {
					assert.False(t, cbCalled.Load())
				}
			})
		})
	}
}

func TestOptInPreservesOtherAdvisories(t *testing.T) {
	t.Run("ExplicitBucketBoundaries", func(t *testing.T) {
		rdr := sdkmetric.NewManualReader()
		view := x.NewView(
			sdkmetric.Instrument{Name: "optin.hist"},
			x.Stream{Enabled: new(true)},
		)
		mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), x.WithView(view))
		m := mp.Meter("test")

		hist, err := m.Float64Histogram(
			"optin.hist",
			metricx.WithOptIn(),
			metric.WithExplicitBucketBoundaries(1, 5, 10),
		)
		require.NoError(t, err)
		assert.True(t, hist.Enabled(t.Context()))
		hist.Record(t.Context(), 3)

		var rm metricdata.ResourceMetrics
		require.NoError(t, rdr.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		gotHist, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64])
		require.True(t, ok)
		require.Len(t, gotHist.DataPoints, 1)
		assert.Equal(t, []float64{1, 5, 10}, gotHist.DataPoints[0].Bounds)
	})

	t.Run("DefaultAttributes", func(t *testing.T) {
		k1 := attribute.Key("k1")
		k2 := attribute.Key("k2")

		rdr := sdkmetric.NewManualReader()
		view := x.NewView(
			sdkmetric.Instrument{Name: "optin.counter"},
			x.Stream{Enabled: new(true)},
		)
		mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), x.WithView(view))
		m := mp.Meter("test")

		ctr, err := m.Int64Counter(
			"optin.counter",
			metricx.WithOptIn(),
			metricx.WithDefaultAttributes(k1),
		)
		require.NoError(t, err)
		assert.True(t, ctr.Enabled(t.Context()))
		ctr.Add(t.Context(), 1, metric.WithAttributes(k1.String("v1"), k2.String("v2")))

		var rm metricdata.ResourceMetrics
		require.NoError(t, rdr.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		gotSum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
		require.True(t, ok)
		require.Len(t, gotSum.DataPoints, 1)
		assert.Equal(t, attribute.NewSet(k1.String("v1")), gotSum.DataPoints[0].Attributes)
	})
}

func TestMultipleViewsWithOptInAndEnabled(t *testing.T) {
	rdr := sdkmetric.NewManualReader()
	// View 1 is a wildcard view without Enabled (drops OptIn instruments).
	// View 2 specifically enables "optin.counter" with the same description.
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(rdr),
		sdkmetric.WithView(
			sdkmetric.NewView(
				sdkmetric.Instrument{Name: "*"},
				sdkmetric.Stream{Description: "shared"},
			),
		),
		x.WithView(
			nil,
			x.NewView(
				sdkmetric.Instrument{Name: "optin.counter"},
				x.Stream{
					Description: "shared",
					Enabled:     new(true),
				},
			),
		),
	)
	m := mp.Meter("test")

	ctr, err := m.Int64Counter("optin.counter", metricx.WithOptIn())
	require.NoError(t, err)
	assert.True(t, ctr.Enabled(t.Context()))
	ctr.Add(t.Context(), 10)

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	assert.Equal(t, "optin.counter", rm.ScopeMetrics[0].Metrics[0].Name)
	assert.Equal(t, "shared", rm.ScopeMetrics[0].Metrics[0].Description)
}

func TestOptInWithRegisterCallback(t *testing.T) {
	rdr := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr))
	m := mp.Meter("test")

	obs, err := m.Int64ObservableGauge("optin.obs", metricx.WithOptIn())
	require.NoError(t, err)

	var called atomic.Bool
	reg, err := m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		called.Store(true)
		o.ObserveInt64(obs, 1)
		return nil
	}, obs)
	require.NoError(t, err)
	defer func() { _ = reg.Unregister() }()

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.Empty(t, rm.ScopeMetrics)
	assert.False(t, called.Load(), "callback should not be invoked when all registered observables are dropped")
}

func TestOptInWithFinishAndMeterConfigurator(t *testing.T) {
	t.Run("WithFinish_DisabledByDefault", func(t *testing.T) {
		rdr := sdkmetric.NewManualReader()
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(rdr),
			x.WithFinish(),
		)
		m := mp.Meter("test")

		ic, err := m.Int64Counter("optin.int64", metricx.WithOptIn())
		require.NoError(t, err)
		assert.False(t, ic.Enabled(t.Context()))
		ic.Add(t.Context(), 1)
		if f, ok := ic.(metricx.Finisher); ok {
			f.Finish(t.Context())
		}

		fc, err := m.Float64Counter("optin.float64", metricx.WithOptIn())
		require.NoError(t, err)
		assert.False(t, fc.Enabled(t.Context()))
		fc.Add(t.Context(), 1.5)
		if f, ok := fc.(metricx.Finisher); ok {
			f.Finish(t.Context())
		}

		var rm metricdata.ResourceMetrics
		require.NoError(t, rdr.Collect(t.Context(), &rm))
		assert.Empty(t, rm.ScopeMetrics)
	})

	t.Run("WithMeterConfigurator", func(t *testing.T) {
		handle := x.NewMeterConfiguratorHandle()
		rdr := sdkmetric.NewManualReader()
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(rdr),
			x.WithMeterConfigurator(handle),
			x.WithView(
				x.NewView(
					sdkmetric.Instrument{Name: "enabled.optin"},
					x.Stream{Enabled: new(true)},
				),
			),
		)
		m := mp.Meter("scope-a")

		disabledOptIn, err := m.Int64Counter("disabled.optin", metricx.WithOptIn())
		require.NoError(t, err)
		assert.False(t, disabledOptIn.Enabled(t.Context()))

		enabledOptIn, err := m.Int64Counter("enabled.optin", metricx.WithOptIn())
		require.NoError(t, err)
		assert.True(t, enabledOptIn.Enabled(t.Context()))

		var obsCalled atomic.Bool
		_, err = m.Int64ObservableCounter(
			"disabled.obs.optin",
			metricx.WithOptIn(),
			metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
				obsCalled.Store(true)
				o.Observe(1)
				return nil
			}),
		)
		require.NoError(t, err)

		var rm metricdata.ResourceMetrics
		require.NoError(t, rdr.Collect(t.Context(), &rm))
		assert.False(t, obsCalled.Load())

		handle.Set(func(instrumentation.Scope) x.MeterConfig {
			return x.NewMeterConfig(x.WithMeterEnabled(false))
		})
		assert.False(t, enabledOptIn.Enabled(t.Context()))
	})
}
