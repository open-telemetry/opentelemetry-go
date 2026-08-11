// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// testMeterConfig satisfies meterConfigReader without importing sdk/metric/x.
type testMeterConfig struct{ enabled bool }

func (c testMeterConfig) Enabled() bool { return c.enabled }

// testConfiguratorOpt is a test implementation of meterConfiguratorOption.
type testConfiguratorOpt struct {
	Option
	fn       func(instrumentation.Scope) any
	onUpdate func(func())
}

func (testConfiguratorOpt) Experimental() {}

func (o testConfiguratorOpt) MeterConfigurator() func(instrumentation.Scope) any { return o.fn }

func (o testConfiguratorOpt) RegisterOnUpdate(cb func()) {
	if o.onUpdate != nil {
		o.onUpdate(cb)
	}
}

// disablingConfiguratorFn disables meters whose scope name is "disabled".
func disablingConfiguratorFn(s instrumentation.Scope) any {
	return testMeterConfig{enabled: s.Name != "disabled"}
}

func TestConfiguratorNewMeter(t *testing.T) {
	for _, tc := range []struct {
		name            string
		configuratorOpt Option
		scopeName       string
		wantEnabled     bool
	}{
		{
			name:            "configurator/none",
			configuratorOpt: nil,
			scopeName:       "any",
			wantEnabled:     true,
		},
		{
			name:            "configurator/scope/disabled",
			configuratorOpt: testConfiguratorOpt{fn: disablingConfiguratorFn},
			scopeName:       "disabled",
			wantEnabled:     false,
		},
		{
			name:            "configurator/scope/enabled",
			configuratorOpt: testConfiguratorOpt{fn: disablingConfiguratorFn},
			scopeName:       "other",
			wantEnabled:     true,
		},
		{
			name:            "configurator/set-before-provider",
			configuratorOpt: testConfiguratorOpt{fn: disablingConfiguratorFn},
			scopeName:       "disabled",
			wantEnabled:     false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var configuratorOpts []Option
			if tc.configuratorOpt != nil {
				configuratorOpts = append(configuratorOpts, tc.configuratorOpt)
			}

			mp := NewMeterProvider(configuratorOpts...)
			defer mp.Shutdown(t.Context()) //nolint:errcheck

			_ = mp.Meter(tc.scopeName)
			m := mp.meters.Lookup(instrumentation.Scope{Name: tc.scopeName}, func() *meter {
				return newMeter(instrumentation.Scope{Name: tc.scopeName}, mp.pipes)
			})
			require.NotNil(t, m)
			assert.Equal(t, tc.wantEnabled, m.enabled.Load())
		})
	}
}

func TestConfiguratorCacheWalkUpdatesCachedMeter(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	mp := NewMeterProvider(configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	// Create and cache a meter before updating the configurator.
	_ = mp.Meter("test")
	cachedMeter := mp.meters.Lookup(instrumentation.Scope{Name: "test"}, func() *meter {
		return newMeter(instrumentation.Scope{Name: "test"}, mp.pipes)
	})
	assert.True(t, cachedMeter.enabled.Load(), "meter should be enabled before configurator update")

	// Swap configurator to disable all scopes and simulate handle.Set().
	mp.configurator = func(_ instrumentation.Scope) any {
		return testMeterConfig{enabled: false}
	}
	if storedCallback != nil {
		storedCallback()
	}

	// Cached meter is updated by the cache walk triggered via onUpdate.
	assert.False(t, cachedMeter.enabled.Load(), "cached meter should be updated by cache walk")
}

func TestInstrumentEnabledReflectsConfigurator(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	m := mp.Meter("disabled")
	ctr, err := m.Int64Counter("ctr")
	require.NoError(t, err)
	assert.False(t, ctr.Enabled(t.Context()), "instrument in disabled scope should report Enabled=false")

	// Re-enable via configurator update; Enabled() should reflect the live
	// meter state rather than a value captured at instrument-creation time.
	mp.configurator = func(instrumentation.Scope) any {
		return testMeterConfig{enabled: true}
	}
	require.NotNil(t, storedCallback)
	storedCallback()
	assert.True(t, ctr.Enabled(t.Context()), "instrument should reflect re-enabled meter")
}

func TestInstrumentAddGatedByConfigurator(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	ctr, err := mp.Meter("disabled").Int64Counter("ctr")
	require.NoError(t, err)

	ctr.Add(t.Context(), 5)

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.Empty(t, rm.ScopeMetrics, "Add on a disabled meter should not reach the aggregator")

	// Re-enable and confirm Add() now reaches the aggregator.
	mp.configurator = func(instrumentation.Scope) any {
		return testMeterConfig{enabled: true}
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	ctr.Add(t.Context(), 7)

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
}

func TestObservableCallbackGatedByConfigurator(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	_, err := mp.Meter("disabled").Int64ObservableCounter("ctr", metric.WithInt64Callback(
		func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(5)
			return nil
		},
	))
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.Empty(t, rm.ScopeMetrics, "callback Observe on a disabled meter should not reach the aggregator")

	// Re-enable and confirm the callback's Observe() now reaches the aggregator.
	mp.configurator = func(instrumentation.Scope) any {
		return testMeterConfig{enabled: true}
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
}

func TestObserverObserveGatedByConfigurator(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	m := mp.Meter("disabled")
	ctr, err := m.Int64ObservableCounter("ctr")
	require.NoError(t, err)

	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveInt64(ctr, 5)
		return nil
	}, ctr)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.Empty(t, rm.ScopeMetrics, "ObserveInt64 on a disabled meter should not reach the aggregator")

	// Re-enable and confirm ObserveInt64 now reaches the aggregator.
	mp.configurator = func(instrumentation.Scope) any {
		return testMeterConfig{enabled: true}
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
}
