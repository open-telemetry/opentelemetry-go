// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"sync/atomic"
	"testing"
	"time"

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

// cachedConfiguratorMeter returns the gated meter mp caches for scope name.
func cachedConfiguratorMeter(t *testing.T, mp *MeterProvider, name string) *configuratorMeter {
	t.Helper()
	gm, ok := mp.Meter(name).(*configuratorMeter)
	require.True(t, ok, "a provider with a configurator must return a *configuratorMeter")
	return gm
}

func TestConfiguratorMeterType(t *testing.T) {
	t.Run("WithConfigurator", func(t *testing.T) {
		mp := NewMeterProvider(testConfiguratorOpt{fn: disablingConfiguratorFn})
		m := mp.Meter("scope")
		assert.IsType(t, &configuratorMeter{}, m, "a provider with a configurator must return a *configuratorMeter")
		assert.Same(t, m, mp.Meter("scope"), "the same scope must return the same *configuratorMeter")
	})
	t.Run("WithoutConfigurator", func(t *testing.T) {
		mp := NewMeterProvider()
		m := mp.Meter("scope")
		assert.IsType(t, &meter{}, m, "a provider without a configurator must return the plain *meter")
		assert.Same(t, m, mp.Meter("scope"), "the same scope must return the same *meter")
	})
}

func TestConfiguratorSyncInstrumentTypes(t *testing.T) {
	tests := []struct {
		name  string
		int64 bool
		build func(metric.Meter) (any, error)
	}{
		{"Int64Counter", true, func(m metric.Meter) (any, error) { return m.Int64Counter("inst") }},
		{"Int64UpDownCounter", true, func(m metric.Meter) (any, error) { return m.Int64UpDownCounter("inst") }},
		{"Int64Histogram", true, func(m metric.Meter) (any, error) { return m.Int64Histogram("inst") }},
		{"Int64Gauge", true, func(m metric.Meter) (any, error) { return m.Int64Gauge("inst") }},
		{"Float64Counter", false, func(m metric.Meter) (any, error) { return m.Float64Counter("inst") }},
		{"Float64UpDownCounter", false, func(m metric.Meter) (any, error) { return m.Float64UpDownCounter("inst") }},
		{"Float64Histogram", false, func(m metric.Meter) (any, error) { return m.Float64Histogram("inst") }},
		{"Float64Gauge", false, func(m metric.Meter) (any, error) { return m.Float64Gauge("inst") }},
	}
	for _, tt := range tests {
		t.Run(tt.name+"/WithConfigurator", func(t *testing.T) {
			m := NewMeterProvider(testConfiguratorOpt{fn: disablingConfiguratorFn}).Meter("scope")
			i1, err := tt.build(m)
			require.NoError(t, err)
			if tt.int64 {
				assert.IsType(t, &configuratorInt64Inst{}, i1)
			} else {
				assert.IsType(t, &configuratorFloat64Inst{}, i1)
			}
			i2, err := tt.build(m)
			require.NoError(t, err)
			assert.Same(t, i1, i2, "the same instrument must return the same wrapper")
		})
		t.Run(tt.name+"/WithoutConfigurator", func(t *testing.T) {
			m := NewMeterProvider().Meter("scope")
			i, err := tt.build(m)
			require.NoError(t, err)
			if tt.int64 {
				assert.IsType(t, &int64Inst{}, i)
			} else {
				assert.IsType(t, &float64Inst{}, i)
			}
		})
	}
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

			var enabled bool
			switch m := mp.Meter(tc.scopeName).(type) {
			case *configuratorMeter:
				enabled = m.config.Load()
			case *meter:
				// A meter from a provider without a configurator is never
				// gated, so it always records.
				enabled = true
			default:
				t.Fatalf("unexpected meter type %T", m)
			}
			assert.Equal(t, tc.wantEnabled, enabled)
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
	assert.True(t, unregistered, "Shutdown must retire the configurator's handle")
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

		assert.False(t, cachedConfiguratorMeter(t, mp, "race").config.Load(),
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
		assert.False(t, cachedConfiguratorMeter(t, mp, "race").config.Load(),
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

	assert.False(t, cachedConfiguratorMeter(t, mp, "race").config.Load(),
		"final state must match the newer Set() walk, not the stale value the delayed apply step read")
}

// TestConfiguratorConcurrentMeterWaitsForInitialConfig asserts that a meter
// returned to a concurrent caller for the same scope has its initial
// configuration applied before it can record.
func TestConfiguratorConcurrentMeterWaitsForInitialConfig(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var firstCall atomic.Bool
	fn := func(instrumentation.Scope) any {
		if firstCall.CompareAndSwap(false, true) {
			close(entered)
			<-release
		}
		return testMeterConfig{enabled: false}
	}

	rdr := NewManualReader()
	mp := NewMeterProvider(WithReader(rdr), testConfiguratorOpt{fn: fn})
	defer mp.Shutdown(t.Context()) //nolint:errcheck

	creatorDone := make(chan struct{})
	go func() {
		defer close(creatorDone)
		_ = mp.Meter("race")
	}()
	<-entered // the meter is created and its configurator call is stalled

	recorded := make(chan struct{})
	go func() {
		defer close(recorded)
		c, err := mp.Meter("race").Int64Counter("c")
		assert.NoError(t, err)
		c.Add(t.Context(), 7)
	}()

	// Give the concurrent caller the chance to record while the initial
	// configuration is still pending. A correct provider blocks it instead.
	select {
	case <-recorded:
	case <-time.After(100 * time.Millisecond):
	}
	close(release)
	<-creatorDone
	<-recorded

	var rm metricdata.ResourceMetrics
	require.NoError(t, rdr.Collect(t.Context(), &rm))
	assert.Nil(t, findMetricByName(&rm, "c"),
		"a measurement made before the initial configuration was applied must not be exported")
}

// Asserts the callback must not run, and its error must not propagate,
// while the meter is disabled.
// Asserts the callback must not run, and its error must not propagate,
// while the meter is disabled.
