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
	MeterConfigurator() func(instrumentation.Scope) any
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
	require.True(t, ok, "must implement MeterConfigurator() func(scope) any")

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

	for _, tc := range []struct {
		name    string
		scope   instrumentation.Scope
		enabled bool
	}{
		{"scope/enabled", instrumentation.Scope{Name: "test"}, true},
		{"scope/disabled", instrumentation.Scope{Name: "disabled"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := ex.MeterConfigurator()(tc.scope)
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
	result := ex.MeterConfigurator()(instrumentation.Scope{Name: "test"})
	cfg, ok := result.(interface{ Enabled() bool })
	require.True(t, ok)
	assert.True(t, cfg.Enabled(), "zero MeterConfig must be enabled")
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
	result := ex.MeterConfigurator()(instrumentation.Scope{Name: "disabled"})
	cfg, ok := result.(interface{ Enabled() bool })
	require.True(t, ok)
	assert.True(t, cfg.Enabled(), "cleared configurator must fall back to zero MeterConfig")
}
