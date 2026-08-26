// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"sync/atomic"

	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/attrnorm"
)

type meterConfigurator func(instrumentation.Scope) any

type meterConfigReader interface{ Enabled() bool }

type meterConfiguratorOption interface {
	Experimental()
	MeterConfigurator() func(instrumentation.Scope) any
	RegisterOnUpdate(func())
}

// MeterProvider handles the creation and coordination of Meters. All Meters
// created by a MeterProvider will be associated with the same Resource, have
// the same Views applied to them, and have their produced metric telemetry
// passed to the configured Readers.
type MeterProvider struct {
	embedded.MeterProvider

	pipes  pipelines
	meters cache[instrumentation.Scope, *meter]
	// configurator is written once, in NewMeterProvider before mp is returned to
	// any caller, and never reassigned after. No synchronization is needed for
	// that single write or the reads that follow it.
	configurator meterConfigurator

	forceFlush, shutdown func(context.Context) error
	stopped              atomic.Bool
}

// Compile-time check MeterProvider implements metric.MeterProvider.
var _ metric.MeterProvider = (*MeterProvider)(nil)

// NewMeterProvider returns a new and configured MeterProvider.
//
// By default, the returned MeterProvider is configured with the default
// Resource and no Readers. Readers cannot be added after a MeterProvider is
// created. This means the returned MeterProvider, one created with no
// Readers, will perform no operations.
func NewMeterProvider(options ...Option) *MeterProvider {
	conf := newConfig(options)
	flush, sdown := conf.readerSignals()

	mp := &MeterProvider{
		pipes:      newPipelines(conf.res, conf.readers, conf.views, conf.exemplarFilter, conf.cardinalityLimit),
		forceFlush: flush,
		shutdown:   sdown,
	}

	var mco meterConfiguratorOption
	for _, o := range options {
		if m, ok := o.(meterConfiguratorOption); ok {
			mco = m
		}
	}
	if mco != nil {
		mp.configurator = mco.MeterConfigurator()
		mco.RegisterOnUpdate(func() {
			mp.meters.Range(func(s instrumentation.Scope, m *meter) {
				if cr, ok := mp.configurator(s).(meterConfigReader); ok {
					m.setEnabled(cr.Enabled())
				}
			})
		})
	}

	// Log after creation so all readers show correctly they are registered.
	global.Info(
		"MeterProvider created",
		"Resource", conf.res,
		"Readers", conf.readers,
		"Views", len(conf.views),
	)
	return mp
}

// Meter returns a Meter with the given name and configured with options.
//
// The name should be the name of the instrumentation scope creating
// telemetry. This name may be the same as the instrumented code only if that
// code provides built-in instrumentation.
//
// Calls to the Meter method after Shutdown has been called will return Meters
// that perform no operations.
//
// This method is safe to call concurrently.
func (mp *MeterProvider) Meter(name string, options ...metric.MeterOption) metric.Meter {
	if name == "" {
		global.Warn("Invalid Meter name.", "name", name)
	}

	if mp.stopped.Load() {
		return noop.Meter{}
	}

	c := metric.NewMeterConfig(options...)
	attrs, _ := attrnorm.Set(c.InstrumentationAttributes())
	s := instrumentation.Scope{
		Name:       name,
		Version:    c.InstrumentationVersion(),
		SchemaURL:  c.SchemaURL(),
		Attributes: attrs,
	}

	global.Info(
		"Meter created",
		"Name", s.Name,
		"Version", s.Version,
		"SchemaURL", s.SchemaURL,
		"Attributes", s.Attributes,
	)

	var meterCreated bool
	m := mp.meters.Lookup(s, func() *meter {
		meterCreated = true
		m := newMeter(s, mp.pipes)
		m.setEnabled(true)
		return m
	})
	// Apply the configurator outside the cache lock: it runs arbitrary user
	// code, and a configurator that calls Meter for the same MeterProvider
	// would deadlock on the cache's non-reentrant lock otherwise.
	// An already-cached meter either already has it applied, or will get it from a concurrent Set walk.
	if meterCreated && mp.configurator != nil {
		if cr, ok := mp.configurator(s).(meterConfigReader); ok {
			m.setEnabled(cr.Enabled())
		}
	}
	return m
}

// ForceFlush flushes all pending telemetry.
//
// This method honors the deadline or cancellation of ctx. An appropriate
// error will be returned in these situations. There is no guaranteed that all
// telemetry be flushed or all resources have been released in these
// situations.
//
// ForceFlush calls ForceFlush(context.Context) error
// on all Readers that implements this method.
//
// This method is safe to call concurrently.
func (mp *MeterProvider) ForceFlush(ctx context.Context) error {
	if mp.forceFlush != nil {
		return mp.forceFlush(ctx)
	}
	return nil
}

// Shutdown shuts down the MeterProvider flushing all pending telemetry and
// releasing any held computational resources.
//
// This call is idempotent. The first call will perform all flush and
// releasing operations. Subsequent calls will perform no action and will
// return an error stating this.
//
// Measurements made by instruments from meters this MeterProvider created
// will not be exported after Shutdown is called.
//
// This method honors the deadline or cancellation of ctx. An appropriate
// error will be returned in these situations. There is no guaranteed that all
// telemetry be flushed or all resources have been released in these
// situations.
//
// This method is safe to call concurrently.
func (mp *MeterProvider) Shutdown(ctx context.Context) error {
	// Even though it may seem like there is a synchronization issue between the
	// call to `Store` and checking `shutdown`, the Go concurrency model ensures
	// that is not the case, as all the atomic operations executed in a program
	// behave as though executed in some sequentially consistent order. This
	// definition provides the same semantics as C++'s sequentially consistent
	// atomics and Java's volatile variables.
	// See https://go.dev/ref/mem#atomic and https://pkg.go.dev/sync/atomic.

	mp.stopped.Store(true)
	if mp.shutdown != nil {
		return mp.shutdown(ctx)
	}
	return nil
}
