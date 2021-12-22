package metric

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/views"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res            *resource.Resource
		views          []views.View
		hasDefaultView bool
	}

	Option func(cfg *Config)

	provider struct {
		cfg    Config
		lock   sync.Mutex
		meters map[instrumentation.Library]*meter
	}

	meter struct {
		library    instrumentation.Library
		uniqueImpl *registry.UniqueInstrumentMeterImpl
		provider   *provider
		views      *viewstate.State
		syncAccum  *syncstate.Accumulator
		asyncAccum *asyncstate.Accumulator
	}
)

var _ sdkapi.MeterImpl = &meter{}

func WithResource(res *resource.Resource) Option {
	return func(cfg *Config) {
		cfg.res = res
	}
}

func WithView(view views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, view)
	}
}

func WithDefaultView(hasDefaultView bool) Option {
	return func(cfg *Config) {
		cfg.hasDefaultView = hasDefaultView
	}
}

func New(opts ...Option) metric.MeterProvider {
	cfg := Config{
		res:            resource.Default(),
		hasDefaultView: true,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &provider{
		cfg:    cfg,
		meters: map[instrumentation.Library]*meter{},
	}
}

func (p *provider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	cfg := metric.NewMeterConfig(opts...)
	lib := instrumentation.Library{
		Name:      name,
		Version:   cfg.InstrumentationVersion(),
		SchemaURL: cfg.SchemaURL(),
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	m := p.meters[lib]
	if m == nil {
		m = &meter{
			provider:   p,
			library:    lib,
			views:      viewstate.New(lib, p.cfg.views, p.cfg.hasDefaultView),
			syncAccum:  syncstate.New(),
			asyncAccum: asyncstate.New(),
		}
		m.uniqueImpl = registry.NewUniqueInstrumentMeterImpl(m)
		p.meters[lib] = m
	}
	return metric.Meter{
		MeterImpl: m.uniqueImpl,
	}
}

func (m *meter) NewInstrument(descriptor sdkapi.Descriptor) (sdkapi.Instrument, error) {
	cfactory, err := m.views.NewFactory(descriptor)
	if err != nil {
		return nil, err
	}

	if descriptor.InstrumentKind().Synchronous() {
		return m.syncAccum.NewInstrument(descriptor, cfactory)
	}
	return m.asyncAccum.NewInstrument(descriptor, cfactory)
}

func (m *meter) NewCallback(insts []sdkapi.Instrument, callback func(context.Context) error) (sdkapi.Callback, error) {
	return m.asyncAccum.NewCallback(insts, callback)
}
