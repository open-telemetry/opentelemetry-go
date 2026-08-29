// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"errors"
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
	version  *atomic.Uint64 // nil => every snapshot reports version 0
	onUpdate func(func())
	// rejectRegistration simulates RegisterOnUpdate finding the handle
	// already claimed by another MeterProvider; zero value (false) claims
	// normally, matching every pre-existing test.
	rejectRegistration bool
	unregister         func()
}

func (testConfiguratorOpt) Experimental() {}

func (o testConfiguratorOpt) MeterConfiguratorSnapshot() func() (func(instrumentation.Scope) any, uint64) {
	return func() (func(instrumentation.Scope) any, uint64) {
		var v uint64
		if o.version != nil {
			v = o.version.Load()
		}
		return o.fn, v
	}
}

func (o testConfiguratorOpt) RegisterOnUpdate(cb func()) bool {
	if o.rejectRegistration {
		return false
	}
	if o.onUpdate != nil {
		o.onUpdate(cb)
	}
	return true
}

func (o testConfiguratorOpt) Unregister() {
	if o.unregister != nil {
		o.unregister()
	}
}

// disablingConfiguratorFn disables meters whose scope name is "disabled".
func disablingConfiguratorFn(s instrumentation.Scope) any {
	return testMeterConfig{enabled: s.Name != "disabled"}
}

// errCallbackShouldNotRun is returned by callbacks in the NotInvokedWhileDisabled
// tests below; seeing it propagate means the callback ran when it shouldn't have.
var errCallbackShouldNotRun = errors.New("callback should not run while meter is disabled")

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

func TestConfiguratorShutdownReleasesHandle(t *testing.T) {
	var unregistered bool
	configuratorOpt := testConfiguratorOpt{
		fn:         disablingConfiguratorFn,
		unregister: func() { unregistered = true },
	}

	mp := NewMeterProvider(configuratorOpt)
	assert.False(t, unregistered, "Unregister must not run before Shutdown")

	require.NoError(t, mp.Shutdown(t.Context()))
	assert.True(t, unregistered, "Shutdown must release the configurator's claim on its handle")
}

// This test guards against a provider whose registration was rejected (handle already claimed
// elsewhere) releasing someone else's active claim on its own Shutdown.
func TestConfiguratorShutdownSkipsUnregisterWhenNotClaimed(t *testing.T) {
	var unregistered bool
	configuratorOpt := testConfiguratorOpt{
		fn:                 disablingConfiguratorFn,
		rejectRegistration: true,
		unregister:         func() { unregistered = true },
	}

	mp := NewMeterProvider(configuratorOpt)
	require.NoError(t, mp.Shutdown(t.Context()))
	assert.False(t, unregistered, "Shutdown must not unregister a claim this provider never held")
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
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(_ instrumentation.Scope) any {
			return testMeterConfig{enabled: false}
		}, 0
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

// TestConfiguratorStaleApplyLosesRaceToNewerSet exercises the case the
// version-stamped CAS exists for: a new meter's creation-time apply step
// reads an old configurator snapshot, then stalls before writing it, and only
// resumes after a concurrent Set() walk has already landed a newer decision
// for the same meter. The final state must match the newer Set() walk, never
// the stale value the delayed apply step read.
func TestConfiguratorStaleApplyLosesRaceToNewerSet(t *testing.T) {
	version := new(atomic.Uint64)
	version.Store(1) // the "old" configurator's version

	started := make(chan struct{})
	release := make(chan struct{})
	var firstCall atomic.Bool

	fn := func(s instrumentation.Scope) any {
		if s.Name != "race" {
			return testMeterConfig{enabled: true}
		}
		if firstCall.CompareAndSwap(false, true) {
			close(started)
			<-release                             // hold until the concurrent Set() walk has fully landed
			return testMeterConfig{enabled: true} // the old, now-stale decision
		}
		return testMeterConfig{enabled: false} // the newer Set()'s decision
	}

	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       fn,
		version:  version,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	mp := NewMeterProvider(configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = mp.Meter("race")
	}()
	<-started // creation's apply step captured version 1 and is now stalled

	// test's stand-in for what a real handle.Set() does internally (h.version.Add(1))
	version.Store(2)
	require.NotNil(t, storedCallback)
	// simulates the concurrent Set() walk: "race" is already
	// cached, so it lands version 2 directly.
	storedCallback()

	close(release) // let the stale apply step resume and try to write version 1
	<-done

	cachedMeter := mp.meters.Lookup(instrumentation.Scope{Name: "race"}, func() *meter {
		return newMeter(instrumentation.Scope{Name: "race"}, mp.pipes)
	})
	assert.False(t, cachedMeter.enabled.Load(),
		"final state must match the newer Set() walk, not the stale value the delayed apply step read")
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
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
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
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
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
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
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
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	rm = metricdata.ResourceMetrics{}
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
}

// Asserts the callback must not run, and its error must not propagate,
// while the meter is disabled.
func TestObservableCallbackNotInvokedWhileDisabled(t *testing.T) {
	var storedCallback func()
	configuratorOpt := testConfiguratorOpt{
		fn:       disablingConfiguratorFn,
		onUpdate: func(cb func()) { storedCallback = cb },
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), configuratorOpt)
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	var invoked bool
	_, err := mp.Meter("disabled").Int64ObservableCounter("ctr", metric.WithInt64Callback(
		func(_ context.Context, o metric.Int64Observer) error {
			invoked = true
			o.Observe(5)
			return errCallbackShouldNotRun
		},
	))
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	err = rdr.Collect(t.Context(), &rm)
	assert.NoError(t, err, "callback error must not propagate while meter is disabled")
	assert.False(t, invoked, "callback must not run at all while meter is disabled")
	assert.Empty(t, rm.ScopeMetrics)

	// Re-enable and confirm the callback resumes on the very next collection,
	// with no re-registration needed.
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	rm = metricdata.ResourceMetrics{}
	err = rdr.Collect(t.Context(), &rm)
	assert.ErrorIs(t, err, errCallbackShouldNotRun, "callback error must propagate once re-enabled")
	assert.True(t, invoked, "callback must run once re-enabled")
}

// Asserts the callback must not run, and its error must not propagate,
// while the meter is disabled.
func TestRegisterCallbackNotInvokedWhileDisabled(t *testing.T) {
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

	var invoked bool
	_, err = m.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		invoked = true
		o.ObserveInt64(ctr, 5)
		return errCallbackShouldNotRun
	}, ctr)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	err = rdr.Collect(t.Context(), &rm)
	assert.NoError(t, err, "callback error must not propagate while meter is disabled")
	assert.False(t, invoked, "callback must not run at all while meter is disabled")
	assert.Empty(t, rm.ScopeMetrics)

	// Re-enable and confirm the callback resumes on the very next collection,
	// with no re-registration needed.
	mp.configurator = func() (func(instrumentation.Scope) any, uint64) {
		return func(instrumentation.Scope) any {
			return testMeterConfig{enabled: true}
		}, 0
	}
	require.NotNil(t, storedCallback)
	storedCallback()

	rm = metricdata.ResourceMetrics{}
	err = rdr.Collect(t.Context(), &rm)
	assert.ErrorIs(t, err, errCallbackShouldNotRun, "callback error must propagate once re-enabled")
	assert.True(t, invoked, "callback must run once re-enabled")
}
