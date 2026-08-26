// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type meterConfiguratorOptionExtractor interface {
	MeterConfiguratorSnapshot() func() (func(instrumentation.Scope) any, uint64)
}

type meterConfiguratorOnUpdateRegistrar interface {
	RegisterOnUpdate(func())
}

func TestWithMeterConfiguratorImplementsExperimentalOption(t *testing.T) {
	h := NewMeterConfiguratorHandle()
	opt := WithMeterConfigurator(h)

	type experimental interface{ Experimental() }
	_, ok := opt.(experimental)
	require.True(t, ok, "must implement Experimental()")

	_, ok = opt.(meterConfiguratorOptionExtractor)
	require.True(t, ok, "must implement MeterConfiguratorSnapshot() func() (func(scope) any, uint64)")

	_, ok = opt.(meterConfiguratorOnUpdateRegistrar)
	require.True(t, ok, "must implement RegisterOnUpdate(func())")
}

func TestWithMeterConfiguratorReflectsSet(t *testing.T) {
	h := NewMeterConfiguratorHandle()
	h.Set(func(s instrumentation.Scope) MeterConfig {
		return NewMeterConfig(WithMeterEnabled(s.Name != "disabled"))
	})

	opt := WithMeterConfigurator(h)
	ex := opt.(meterConfiguratorOptionExtractor)
	snapshot := ex.MeterConfiguratorSnapshot()
	fn, version := snapshot()
	assert.Equal(t, uint64(1), version, "first Set call must produce version 1")

	for _, tc := range []struct {
		name    string
		scope   instrumentation.Scope
		enabled bool
	}{
		{"scope/enabled", instrumentation.Scope{Name: "test"}, true},
		{"scope/disabled", instrumentation.Scope{Name: "disabled"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := fn(tc.scope)
			cfg, ok := result.(interface{ Enabled() bool })
			require.True(t, ok, "result must implement Enabled() bool")
			assert.Equal(t, tc.enabled, cfg.Enabled())
		})
	}
}

func TestMeterConfiguratorHandleSetTriggersOnUpdate(t *testing.T) {
	h := NewMeterConfiguratorHandle()

	walked := false
	opt := WithMeterConfigurator(h)
	opt.(meterConfiguratorOnUpdateRegistrar).RegisterOnUpdate(func() { walked = true })

	h.Set(func(_ instrumentation.Scope) MeterConfig { return MeterConfig{} })
	assert.True(t, walked, "Set must trigger the registered onUpdate callback")
}

func TestMeterConfiguratorHandleSetSerializesConcurrentCalls(t *testing.T) {
	h := NewMeterConfiguratorHandle()
	opt := WithMeterConfigurator(h)

	var (
		mu      sync.Mutex
		active  int
		overlap bool
	)
	opt.(meterConfiguratorOnUpdateRegistrar).RegisterOnUpdate(func() {
		mu.Lock()
		active++
		if active > 1 {
			overlap = true
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		active--
		mu.Unlock()
	})

	const goroutines = 10
	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			h.Set(func(instrumentation.Scope) MeterConfig { return MeterConfig{} })
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		require.Fail(t, "timeout waiting for concurrent Set calls")
	}

	assert.False(t, overlap, "concurrent Set calls' onUpdate callbacks overlapped")
}

func TestMeterConfiguratorHandleSetNoConfigurator(t *testing.T) {
	h := NewMeterConfiguratorHandle()
	opt := WithMeterConfigurator(h)
	ex := opt.(meterConfiguratorOptionExtractor)

	// no Set called; closure must return zero MeterConfig, not panic
	snapshot := ex.MeterConfiguratorSnapshot()
	fn, version := snapshot()
	assert.Equal(t, uint64(0), version, "never-configured handle must report version 0")
	result := fn(instrumentation.Scope{Name: "test"})
	cfg, ok := result.(interface{ Enabled() bool })
	require.True(t, ok)
	assert.True(t, cfg.Enabled(), "zero MeterConfig must be enabled")
}

func TestWithMeterConfiguratorNilHandle(t *testing.T) {
	opt := WithMeterConfigurator(nil)
	ex := opt.(meterConfiguratorOptionExtractor)

	// nil handle must not panic; must behave as if the option were omitted.
	snapshot := ex.MeterConfiguratorSnapshot()
	fn, version := snapshot()
	assert.Equal(t, uint64(0), version, "nil handle must report version 0")
	result := fn(instrumentation.Scope{Name: "test"})
	cfg, ok := result.(interface{ Enabled() bool })
	require.True(t, ok)
	assert.True(t, cfg.Enabled(), "nil handle must fall back to zero MeterConfig")

	// RegisterOnUpdate on a nil handle must not panic.
	assert.NotPanics(t, func() {
		opt.(meterConfiguratorOnUpdateRegistrar).RegisterOnUpdate(func() {})
	})
}

func TestMeterConfiguratorHandleSetNilClears(t *testing.T) {
	h := NewMeterConfiguratorHandle()
	opt := WithMeterConfigurator(h)
	ex := opt.(meterConfiguratorOptionExtractor)

	h.Set(func(s instrumentation.Scope) MeterConfig {
		return NewMeterConfig(WithMeterEnabled(s.Name != "disabled"))
	})
	h.Set(nil)

	// Set(nil) must clear the configurator, not store a nil func; the
	// closure must fall back to the zero MeterConfig instead of panicking
	// on a nil func call.
	snapshot := ex.MeterConfiguratorSnapshot()
	fn, version := snapshot()
	// TODO: this is 0, not 2, because Set(nil) stores a bare nil pointer
	// rather than a versionedConfigurator, so its version is lost. See
	// discussion: this undercounts and can let a stale write win a CAS
	// against a meter that already observed a later real version.
	assert.Equal(t, uint64(0), version)
	result := fn(instrumentation.Scope{Name: "disabled"})
	cfg, ok := result.(interface{ Enabled() bool })
	require.True(t, ok)
	assert.True(t, cfg.Enabled(), "cleared configurator must fall back to zero MeterConfig")
}
