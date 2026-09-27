// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// These tests pin the configurator's gating behavior through the public API
// only, so they hold regardless of how the SDK implements the gate.

func gatingMeterProvider(t *testing.T) (*sdkmetric.MeterProvider, *sdkmetric.ManualReader, *MeterConfiguratorHandle) {
	t.Helper()
	h := NewMeterConfiguratorHandle()
	rdr := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(rdr), WithMeterConfigurator(h))
	t.Cleanup(func() { _ = mp.Shutdown(context.Background()) })
	return mp, rdr, h
}

// gatingSetAll enables or disables every scope.
func gatingSetAll(h *MeterConfiguratorHandle, enabled bool) {
	h.Set(func(instrumentation.Scope) MeterConfig {
		return NewMeterConfig(WithMeterEnabled(enabled))
	})
}

// gatingPoints collects rdr and returns the number of data points per metric
// name, across all scopes.
func gatingPoints(t *testing.T, rdr *sdkmetric.ManualReader) map[string]int {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	points := map[string]int{}
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			switch data := m.Data.(type) {
			case metricdata.Sum[int64]:
				points[m.Name] += len(data.DataPoints)
			case metricdata.Sum[float64]:
				points[m.Name] += len(data.DataPoints)
			case metricdata.Gauge[int64]:
				points[m.Name] += len(data.DataPoints)
			case metricdata.Gauge[float64]:
				points[m.Name] += len(data.DataPoints)
			case metricdata.Histogram[int64]:
				points[m.Name] += len(data.DataPoints)
			case metricdata.Histogram[float64]:
				points[m.Name] += len(data.DataPoints)
			default:
				t.Fatalf("unexpected data type %T for %q", m.Data, m.Name)
			}
		}
	}
	return points
}

// gatingOrder is when every scope is disabled, relative to creating the meter
// and its instrument.
type gatingOrder int

const (
	// gatingDisableExisting disables after the meter and instrument exist, so
	// the Set walk applies it.
	gatingDisableExisting gatingOrder = iota
	// gatingDisableBeforeMeter disables before the meter is created, so it is
	// applied when the meter is created.
	gatingDisableBeforeMeter
	// gatingDisableBeforeInstrument disables after the meter is created but
	// before the instrument is, so the instrument must pick up the meter's
	// current state.
	gatingDisableBeforeInstrument
)

var gatingOrders = []struct {
	name  string
	order gatingOrder
}{
	{"DisableExisting", gatingDisableExisting},
	{"DisableBeforeMeter", gatingDisableBeforeMeter},
	{"DisableBeforeInstrument", gatingDisableBeforeInstrument},
}

// gatingCreate creates the meter for scope "scope" and passes it to create,
// disabling every scope at the point order names. Every scope is disabled
// when it returns.
func gatingCreate(mp *sdkmetric.MeterProvider, h *MeterConfiguratorHandle, order gatingOrder, create func(metric.Meter)) {
	if order == gatingDisableBeforeMeter {
		gatingSetAll(h, false)
	}
	m := mp.Meter("scope")
	if order == gatingDisableBeforeInstrument {
		gatingSetAll(h, false)
	}
	create(m)
	if order == gatingDisableExisting {
		gatingSetAll(h, false)
	}
}

func TestGatingSyncInstruments(t *testing.T) {
	const name = "inst"
	type syncInst struct {
		record  func(context.Context)
		enabled func(context.Context) bool
	}
	tests := []struct {
		name  string
		build func(metric.Meter) (syncInst, error)
	}{
		{"Int64Counter", func(m metric.Meter) (syncInst, error) {
			i, err := m.Int64Counter(name)
			return syncInst{func(ctx context.Context) { i.Add(ctx, 1) }, i.Enabled}, err
		}},
		{"Int64UpDownCounter", func(m metric.Meter) (syncInst, error) {
			i, err := m.Int64UpDownCounter(name)
			return syncInst{func(ctx context.Context) { i.Add(ctx, 1) }, i.Enabled}, err
		}},
		{"Int64Histogram", func(m metric.Meter) (syncInst, error) {
			i, err := m.Int64Histogram(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
		{"Int64Gauge", func(m metric.Meter) (syncInst, error) {
			i, err := m.Int64Gauge(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64Counter", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64Counter(name)
			return syncInst{func(ctx context.Context) { i.Add(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64UpDownCounter", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64UpDownCounter(name)
			return syncInst{func(ctx context.Context) { i.Add(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64Histogram", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64Histogram(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64Gauge", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64Gauge(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
	}
	for _, tt := range tests {
		for _, o := range gatingOrders {
			t.Run(tt.name+"/"+o.name, func(t *testing.T) {
				ctx := t.Context()
				mp, rdr, h := gatingMeterProvider(t)
				var inst syncInst
				gatingCreate(mp, h, o.order, func(m metric.Meter) {
					var err error
					inst, err = tt.build(m)
					require.NoError(t, err)
					if o.order == gatingDisableExisting {
						require.True(t, inst.enabled(ctx), "instrument must start enabled")
					}
				})

				assert.False(t, inst.enabled(ctx), "Enabled must report false while disabled")
				inst.record(ctx)
				assert.Zero(t, gatingPoints(t, rdr)[name], "a measurement while disabled must be dropped")

				gatingSetAll(h, true)
				assert.True(t, inst.enabled(ctx), "Enabled must follow a re-enabling Set")
				inst.record(ctx)
				assert.Equal(t, 1, gatingPoints(t, rdr)[name], "a measurement after re-enabling must be recorded")
			})
		}
	}
}

func TestGatingObservableCallbacks(t *testing.T) {
	const name = "obs"
	tests := []struct {
		name  string
		build func(m metric.Meter, called func()) error
	}{
		{"Int64ObservableCounter", func(m metric.Meter, called func()) error {
			_, err := m.Int64ObservableCounter(name, metric.WithInt64Callback(
				func(_ context.Context, o metric.Int64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
		{"Int64ObservableUpDownCounter", func(m metric.Meter, called func()) error {
			_, err := m.Int64ObservableUpDownCounter(name, metric.WithInt64Callback(
				func(_ context.Context, o metric.Int64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
		{"Int64ObservableGauge", func(m metric.Meter, called func()) error {
			_, err := m.Int64ObservableGauge(name, metric.WithInt64Callback(
				func(_ context.Context, o metric.Int64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
		{"Float64ObservableCounter", func(m metric.Meter, called func()) error {
			_, err := m.Float64ObservableCounter(name, metric.WithFloat64Callback(
				func(_ context.Context, o metric.Float64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
		{"Float64ObservableUpDownCounter", func(m metric.Meter, called func()) error {
			_, err := m.Float64ObservableUpDownCounter(name, metric.WithFloat64Callback(
				func(_ context.Context, o metric.Float64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
		{"Float64ObservableGauge", func(m metric.Meter, called func()) error {
			_, err := m.Float64ObservableGauge(name, metric.WithFloat64Callback(
				func(_ context.Context, o metric.Float64Observer) error { called(); o.Observe(1); return nil }))
			return err
		}},
	}
	for _, tt := range tests {
		for _, o := range gatingOrders {
			t.Run(tt.name+"/"+o.name, func(t *testing.T) {
				mp, rdr, h := gatingMeterProvider(t)
				calls := 0
				gatingCreate(mp, h, o.order, func(m metric.Meter) {
					require.NoError(t, tt.build(m, func() { calls++ }))
				})

				assert.Zero(t, gatingPoints(t, rdr)[name], "a disabled meter must produce no observations")
				assert.Zero(t, calls, "a disabled meter's callback must not be invoked")

				gatingSetAll(h, true)
				assert.Equal(t, 1, gatingPoints(t, rdr)[name], "a re-enabled meter must produce observations")
				assert.Equal(t, 1, calls, "a re-enabled meter's callback must be invoked once per collection")
			})
		}
	}
}

func TestGatingRegisterCallback(t *testing.T) {
	const name = "obs"
	mp, rdr, h := gatingMeterProvider(t)
	m := mp.Meter("scope")
	inst, err := m.Int64ObservableCounter(name)
	require.NoError(t, err)
	calls := 0
	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		calls++
		o.ObserveInt64(inst, 1)
		return nil
	}, inst)
	require.NoError(t, err)

	gatingSetAll(h, false)
	assert.Zero(t, gatingPoints(t, rdr)[name], "a disabled meter must produce no observations")
	assert.Zero(t, calls, "a disabled meter's registered callback must not be invoked")

	gatingSetAll(h, true)
	assert.Equal(t, 1, gatingPoints(t, rdr)[name], "a re-enabled meter must produce observations")
	assert.Equal(t, 1, calls, "a re-enabled meter's registered callback must be invoked once per collection")
}

func TestGatingScopeIsolation(t *testing.T) {
	disableOffScope := func(h *MeterConfiguratorHandle) {
		h.Set(func(s instrumentation.Scope) MeterConfig {
			return NewMeterConfig(WithMeterEnabled(s.Name != "off"))
		})
	}
	tests := []struct {
		name string
		// setBeforeMeters sets the configurator before either meter is
		// created, so each meter gets its scope's config at creation rather
		// than from the Set walk.
		setBeforeMeters bool
	}{
		{"SetAfterMeterCreation", false},
		{"SetBeforeMeterCreation", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mp, rdr, h := gatingMeterProvider(t)
			if tt.setBeforeMeters {
				disableOffScope(h)
			}
			on, err := mp.Meter("on").Int64Counter("on.counter")
			require.NoError(t, err)
			off, err := mp.Meter("off").Int64Counter("off.counter")
			require.NoError(t, err)
			if !tt.setBeforeMeters {
				disableOffScope(h)
			}

			on.Add(ctx, 1)
			off.Add(ctx, 1)

			points := gatingPoints(t, rdr)
			assert.Equal(t, 1, points["on.counter"], "an enabled scope must record")
			assert.Zero(t, points["off.counter"], "a disabled scope must not record")
			assert.True(t, on.Enabled(ctx), "disabling another scope must not affect this one")
			assert.False(t, off.Enabled(ctx))
		})
	}
}

func TestGatingInstrumentIdentity(t *testing.T) {
	mp, _, _ := gatingMeterProvider(t)

	m1, m2 := mp.Meter("scope"), mp.Meter("scope")
	assert.True(t, m1 == m2, "Meter must return the same meter for the same scope")

	c1, err := m1.Int64Counter("counter")
	require.NoError(t, err)
	c2, err := m1.Int64Counter("counter")
	require.NoError(t, err)
	assert.True(t, c1 == c2, "repeated sync instrument creation must return an equal instrument")

	h1, err := m1.Float64Histogram("histogram")
	require.NoError(t, err)
	h2, err := m1.Float64Histogram("histogram")
	require.NoError(t, err)
	assert.True(t, h1 == h2, "repeated sync instrument creation must return an equal instrument")

	o1, err := m1.Int64ObservableCounter("observable")
	require.NoError(t, err)
	o2, err := m1.Int64ObservableCounter("observable")
	require.NoError(t, err)
	assert.True(t, o1 == o2, "repeated async instrument creation must return an equal instrument")
}
