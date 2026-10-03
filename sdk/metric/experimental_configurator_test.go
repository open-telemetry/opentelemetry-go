// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// newToggledConfiguratorProvider returns a MeterProvider whose configurator
// enables or disables every scope as set by setEnabled, which applies the
// change similar to MeterConfiguratorHandle.Set.
func newToggledConfiguratorProvider(t *testing.T) (*MeterProvider, *ManualReader, func(bool)) {
	t.Helper()
	enabled := new(atomic.Bool)
	enabled.Store(true)
	version := new(atomic.Uint64)
	var walk func()
	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), testConfiguratorOpt{
		fn:       func(instrumentation.Scope) any { return testMeterConfig{enabled: enabled.Load()} },
		version:  version,
		onUpdate: func(cb func()) { walk = cb },
	})
	return mp, rdr, func(e bool) {
		enabled.Store(e)
		version.Add(1)
		walk()
	}
}

// configuratorCollected reports whether collecting rdr exports the metric
// named name. A metric with no measurements is not exported.
func configuratorCollected(t *testing.T, rdr *ManualReader, name string) bool {
	t.Helper()
	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	return findMetricByName(&rm, name) != nil
}

func TestConfiguratorSyncInstrumentsGateRecording(t *testing.T) {
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
		{"Int64Histogram", func(m metric.Meter) (syncInst, error) {
			i, err := m.Int64Histogram(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64Counter", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64Counter(name)
			return syncInst{func(ctx context.Context) { i.Add(ctx, 1) }, i.Enabled}, err
		}},
		{"Float64Histogram", func(m metric.Meter) (syncInst, error) {
			i, err := m.Float64Histogram(name)
			return syncInst{func(ctx context.Context) { i.Record(ctx, 1) }, i.Enabled}, err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			mp, rdr, setEnabled := newToggledConfiguratorProvider(t)
			inst, err := tt.build(mp.Meter("scope"))
			require.NoError(t, err)

			setEnabled(false)
			assert.False(t, inst.enabled(ctx))
			inst.record(ctx)
			assert.False(t, configuratorCollected(t, rdr, name), "a measurement while disabled must be dropped")

			setEnabled(true)
			assert.True(t, inst.enabled(ctx))
			inst.record(ctx)
			assert.True(t, configuratorCollected(t, rdr, name), "a measurement while enabled must be recorded")
		})
	}
}

func TestConfiguratorObservableCallbacksGated(t *testing.T) {
	const name = "obs"
	int64Cb := func(calls *int) metric.Int64Callback {
		return func(_ context.Context, o metric.Int64Observer) error {
			*calls++
			o.Observe(1)
			return nil
		}
	}
	float64Cb := func(calls *int) metric.Float64Callback {
		return func(_ context.Context, o metric.Float64Observer) error {
			*calls++
			o.Observe(1)
			return nil
		}
	}
	tests := []struct {
		name  string
		build func(m metric.Meter, calls *int) error
	}{
		{"Int64ObservableCounter", func(m metric.Meter, calls *int) error {
			_, err := m.Int64ObservableCounter(name, metric.WithInt64Callback(int64Cb(calls)))
			return err
		}},
		{"Int64ObservableUpDownCounter", func(m metric.Meter, calls *int) error {
			_, err := m.Int64ObservableUpDownCounter(name, metric.WithInt64Callback(int64Cb(calls)))
			return err
		}},
		{"Int64ObservableGauge", func(m metric.Meter, calls *int) error {
			_, err := m.Int64ObservableGauge(name, metric.WithInt64Callback(int64Cb(calls)))
			return err
		}},
		{"Float64ObservableCounter", func(m metric.Meter, calls *int) error {
			_, err := m.Float64ObservableCounter(name, metric.WithFloat64Callback(float64Cb(calls)))
			return err
		}},
		{"Float64ObservableUpDownCounter", func(m metric.Meter, calls *int) error {
			_, err := m.Float64ObservableUpDownCounter(name, metric.WithFloat64Callback(float64Cb(calls)))
			return err
		}},
		{"Float64ObservableGauge", func(m metric.Meter, calls *int) error {
			_, err := m.Float64ObservableGauge(name, metric.WithFloat64Callback(float64Cb(calls)))
			return err
		}},
		{"RegisterCallback", func(m metric.Meter, calls *int) error {
			inst, err := m.Int64ObservableCounter(name)
			if err != nil {
				return err
			}
			_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				*calls++
				o.ObserveInt64(inst, 1)
				return nil
			}, inst)
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp, rdr, setEnabled := newToggledConfiguratorProvider(t)
			calls := 0
			require.NoError(t, tt.build(mp.Meter("scope"), &calls))

			setEnabled(false)
			assert.False(t, configuratorCollected(t, rdr, name), "a disabled meter must produce no observations")
			assert.Zero(t, calls, "a disabled meter's callback must not be invoked")

			setEnabled(true)
			assert.True(t, configuratorCollected(t, rdr, name), "an enabled meter must produce observations")
			assert.Equal(t, 1, calls, "an enabled meter's callback must be invoked once per collection")
		})
	}
}

func TestConfiguratorMeterWrapUnexpectedType(t *testing.T) {
	m := newConfiguratorMeter(newMeter(instrumentation.Scope{Name: "scope"}, nil))

	i64, ok := m.wrapInt64(noopInstrument{})
	assert.False(t, ok, "a non-SDK int64 instrument must not be wrapped")
	assert.Nil(t, i64)

	f64, ok := m.wrapFloat64(noopInstrument{})
	assert.False(t, ok, "a non-SDK float64 instrument must not be wrapped")
	assert.Nil(t, f64)
}

// noopInstrument is a value that is not an SDK instrument.
type noopInstrument struct{}
