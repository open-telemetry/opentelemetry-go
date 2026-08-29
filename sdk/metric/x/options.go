// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"errors"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// MeterConfiguratorHandle holds a [MeterConfigurator] that can be updated at
// runtime. Pass it to [WithMeterConfigurator] at [sdkmetric.MeterProvider]
// construction. Calls to [MeterConfiguratorHandle.Set] are reflected
// immediately across all existing meters via a synchronous cache walk.
//
// A Handle belongs to exactly one MeterProvider for its lifetime. Passing an
// already-claimed Handle to another [WithMeterConfigurator] call is not
// supported: the first MeterProvider keeps the claim, the later one does not
// receive Set updates, and the attempt is logged as an error. A Handle whose
// MeterProvider has been shut down is released and may be claimed again.
type MeterConfiguratorHandle struct {
	mu           sync.Mutex // serializes Set calls; see Set's doc comment
	configurator atomic.Pointer[versionedConfigurator]
	onUpdate     atomic.Pointer[func()] // To avoid race between handle.Set() and RegisterOnUpdate
	version      atomic.Uint64          // bumped once per Set call; see Set
	registered   atomic.Bool            // claimed by a MeterProvider; see RegisterOnUpdate/Unregister
}

// errHandleAlreadyRegistered is logged when a MeterConfiguratorHandle already
// claimed by one MeterProvider is passed to WithMeterConfigurator for another.
var errHandleAlreadyRegistered = errors.New("MeterConfiguratorHandle already registered with a MeterProvider")

// versionedConfigurator pairs a MeterConfigurator with the version it was set
// under, so the two are always stored and read together and can never be
// observed mismatched.
type versionedConfigurator struct {
	fn      MeterConfigurator
	version uint64
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
// The walk only covers Meters that already exist when Set is called. A Meter
// concurrently under construction may finish applying this call's result a
// moment after Set has already returned to its caller; it still always
// converges to this call's result or a newer one, never a stale value.
//
// Concurrent calls to Set are serialized: a Set call blocks until any
// already-in-progress Set, including its cache walk, has completed. This
// keeps one Set's cache walk from partially overwriting another's result
// across different meters. The callback registered via RegisterOnUpdate must
// not call Set on the same handle; doing so deadlocks, since this lock is
// not reentrant.
//
// Passing a nil fn clears the configurator, reverting to the same default
// behavior as a handle that has never had Set called on it (the Meter
// enabled; see [MeterConfig]'s zero value).
func (h *MeterConfiguratorHandle) Set(fn MeterConfigurator) {
	h.mu.Lock()
	defer h.mu.Unlock()
	v := h.version.Add(1)
	h.configurator.Store(&versionedConfigurator{fn: fn, version: v})
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
//
// Each Handle must be passed to at most one MeterProvider; see
// [MeterConfiguratorHandle]'s doc comment.
//
// A nil h is a no-op, equivalent to omitting this option entirely.
func WithMeterConfigurator(h *MeterConfiguratorHandle) sdkmetric.Option {
	return meterConfiguratorProviderOption{handle: h}
}

// MeterConfiguratorSnapshot returns a closure over the handle so sdk/metric
// can call it via duck-type without importing this package. Calling the
// returned function takes one atomic snapshot of the currently installed
// configurator and the version it was set under. Callers evaluating multiple
// scopes against the same snapshot should call it once and reuse the
// returned per-scope function and version, rather than calling it again per
// scope.
func (o meterConfiguratorProviderOption) MeterConfiguratorSnapshot() func() (func(instrumentation.Scope) any, uint64) {
	defaultFn := func(instrumentation.Scope) any { return MeterConfig{} }
	return func() (func(instrumentation.Scope) any, uint64) {
		if o.handle == nil {
			return defaultFn, 0
		}
		vc := o.handle.configurator.Load()
		if vc == nil {
			// Set has never been called on this handle.
			return defaultFn, 0
		}
		if vc.fn == nil {
			// Set(nil) was called; the version it bumped still applies.
			return defaultFn, vc.version
		}
		return func(s instrumentation.Scope) any { return vc.fn(s) }, vc.version
	}
}

// RegisterOnUpdate is called once at construction by sdk/metric during [sdkmetric.NewMeterProvider]
// to register the cache walk callback.
// Subsequent [MeterConfiguratorHandle.Set] calls trigger it. Reports whether this call
// claimed the handle (fn is not stored in both cases):
//   - false if the handle is nil, or
//   - already claimed by another MeterProvider (see [MeterConfiguratorHandle]'s doc comment).
func (o meterConfiguratorProviderOption) RegisterOnUpdate(fn func()) bool {
	if o.handle == nil {
		return false
	}
	if !o.handle.registered.CompareAndSwap(false, true) {
		global.Error(errHandleAlreadyRegistered, "did not register MeterConfiguratorHandle")
		return false
	}
	o.handle.onUpdate.Store(&fn)
	return true
}

// Unregister releases this option's claim on the handle, clearing onUpdate so
// a Set call afterward no longer walks this (now presumably shut down) MeterProvider,
// and allowing a future MeterProvider to claim the handle via RegisterOnUpdate.
//
// Called by sdk/metric during [sdkmetric.MeterProvider.Shutdown],
// only for a provider whose RegisterOnUpdate call actually claimed it.
func (o meterConfiguratorProviderOption) Unregister() {
	if o.handle == nil {
		return
	}
	o.handle.onUpdate.Store(nil)
	o.handle.registered.Store(false)
}
