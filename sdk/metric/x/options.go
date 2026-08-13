// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// MeterConfiguratorHandle holds a [MeterConfigurator] that can be updated at
// runtime. Pass it to [WithMeterConfigurator] at [sdkmetric.MeterProvider]
// construction. Calls to [MeterConfiguratorHandle.Set] are reflected
// immediately across all existing meters via a synchronous cache walk.
type MeterConfiguratorHandle struct {
	mu           sync.Mutex // serializes Set calls; see Set's doc comment
	configurator atomic.Pointer[MeterConfigurator]
	onUpdate     atomic.Pointer[func()] // To avoid race between handle.Set() and RegisterOnUpdate
}

// NewMeterConfiguratorHandle returns a new [MeterConfiguratorHandle] with no
// configurator set.
func NewMeterConfiguratorHandle() *MeterConfiguratorHandle {
	return &MeterConfiguratorHandle{}
}

// Set updates the [MeterConfigurator] and triggers a synchronous cache walk on
// the [sdkmetric.MeterProvider] registered via [WithMeterConfigurator]. Set
// does not return until that walk completes. Because the walk invokes fn once
// per existing Meter, Set's duration scales with both the number of Meters on
// the MeterProvider and fn's own latency. See [MeterConfigurator] for the
// requirements this places on fn.
//
// Concurrent calls to Set are serialized: a Set call blocks until any
// already-in-progress Set, including its cache walk, has completed. This
// keeps one Set's cache walk from partially overwriting another's result
// across different meters. The callback registered via RegisterOnUpdate must
// not call Set on the same handle; doing so deadlocks, since this lock is
// not reentrant.
//
// Passing a nil fn clears the configurator, reverting to the same default
// behavior as a handle that has never had Set called on it.
func (h *MeterConfiguratorHandle) Set(fn MeterConfigurator) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if fn == nil {
		h.configurator.Store(nil)
	} else {
		h.configurator.Store(&fn)
	}
	if cb := h.onUpdate.Load(); cb != nil {
		(*cb)()
	}
}

type meterConfiguratorProviderOption struct {
	// nil embed; skip guard in newConfig prevents apply from being called.
	sdkmetric.Option
	handle *MeterConfiguratorHandle
}

// Experimental marks this as an experimental option so the skip guard in
// newConfig skips calling the nil embedded apply.
func (meterConfiguratorProviderOption) Experimental() {}

// WithMeterConfigurator returns an [sdkmetric.Option] that wires a
// [MeterConfiguratorHandle] into a [sdkmetric.MeterProvider]. The handle must
// be passed at construction; runtime configurator updates via
// [MeterConfiguratorHandle.Set] are only supported when a handle is registered
// here. Providers created without this option cannot have a configurator added
// later.
func WithMeterConfigurator(h *MeterConfiguratorHandle) sdkmetric.Option {
	return meterConfiguratorProviderOption{handle: h}
}

// MeterConfigurator returns a closure over the handle so sdk/metric can call
// it via duck-type without importing this package.
func (o meterConfiguratorProviderOption) MeterConfigurator() func(instrumentation.Scope) any {
	return func(s instrumentation.Scope) any {
		if p := o.handle.configurator.Load(); p != nil {
			return (*p)(s)
		}
		return MeterConfig{}
	}
}

// RegisterOnUpdate is called by sdk/metric during [sdkmetric.NewMeterProvider]
// to register the cache walk callback. Called once at construction; subsequent
// [MeterConfiguratorHandle.Set] calls trigger it.
func (o meterConfiguratorProviderOption) RegisterOnUpdate(fn func()) {
	o.handle.onUpdate.Store(&fn)
}
