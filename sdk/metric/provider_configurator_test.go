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

func TestConfiguratorMultipleOptionsLastWins(t *testing.T) {
	var registered1, registered2 bool
	opt1 := testConfiguratorOpt{
		fn:       func(instrumentation.Scope) any { return testMeterConfig{enabled: false} },
		onUpdate: func(func()) { registered1 = true },
	}
	opt2 := testConfiguratorOpt{
		fn:       func(instrumentation.Scope) any { return testMeterConfig{enabled: true} },
		onUpdate: func(func()) { registered2 = true },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), opt1, opt2)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	assert.False(t, registered1, "earlier configurator option must not be wired to this provider")
	assert.True(t, registered2, "only the last configurator option must be wired to this provider")

	ctr, err := mp.Meter("any").Int64Counter("ctr")
	require.NoError(t, err)
	assert.True(t, ctr.Enabled(t.Context()), "provider must use the last configurator option, not the first")
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

// TestConfiguratorNewMeterConvergesWithSetWalk proves both orderings of the
// cache-lock/walk consistency guarantee: whichever of a new meter's insertion
// or a concurrent Set() walk acquires the cache lock first, the meter never
// ends up stale.
func TestConfiguratorNewMeterConvergesWithSetWalk(t *testing.T) {
	newProvider := func() (mp *MeterProvider, enabled *atomic.Bool, walk func()) {
		enabled = new(atomic.Bool)
		enabled.Store(true)

		var storedCallback func()
		configuratorOpt := testConfiguratorOpt{
			fn: func(instrumentation.Scope) any {
				return testMeterConfig{enabled: enabled.Load()}
			},
			onUpdate: func(cb func()) { storedCallback = cb },
		}

		mp = NewMeterProvider(configuratorOpt)
		return mp, enabled, func() { storedCallback() }
	}

	cachedMeter := func(mp *MeterProvider, name string) *meter {
		return mp.meters.Lookup(instrumentation.Scope{Name: name}, func() *meter {
			return newMeter(instrumentation.Scope{Name: name}, mp.pipes)
		})
	}

	t.Run("insert_then_cfg_set", func(t *testing.T) {
		mp, enabled, walk := newProvider()
		defer mp.Shutdown(t.Context()) //nolint:errcheck

		inserted := make(chan struct{})
		go func() {
			defer close(inserted)
			_ = mp.Meter("race")
		}()
		<-inserted

		enabled.Store(false)
		walk()

		assert.False(t, cachedMeter(mp, "race").enabled.Load(),
			"walk started after new meter must observe it")
	})

	t.Run("cfg_set_then_new_meter", func(t *testing.T) {
		mp, enabled, walk := newProvider()
		defer mp.Shutdown(t.Context()) //nolint:errcheck

		walked := make(chan struct{})
		go func() {
			defer close(walked)
			enabled.Store(false)
			walk()
		}()
		<-walked

		_ = mp.Meter("race")
		assert.False(t, cachedMeter(mp, "race").enabled.Load(),
			"meter created after the walk must read the updated configurator directly")
	})
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
