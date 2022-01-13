package metric

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
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
		provider   *provider
		registry   *registry.Registry
		views      *viewstate.State
		syncAccum  *syncstate.Accumulator
		asyncAccum *asyncstate.Accumulator
	}
)

var (
	_ metric.Meter            = &meter{}
	_ syncint64.Instruments   = syncint64Instruments{}
	_ syncfloat64.Instruments = syncfloat64Instruments{}
)

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
	if m != nil {
		return m
	}
	m = &meter{
		provider:   p,
		library:    lib,
		registry:   registry.New(),
		views:      viewstate.New(lib, p.cfg.views, p.cfg.hasDefaultView),
		syncAccum:  syncstate.New(),
		asyncAccum: asyncstate.New(),
	}
	p.meters[lib] = m
	return m
}

func (m *meter) SyncInt64() syncint64.Instruments {
	return m.syncstate.Int64Instruments(m)
}

func (m *meter) SyncFloat64() syncfloat64.Instruments {
	return m.syncstate.Float64Instruments(m)
}

func (m *meter) AsyncInt64() asyncint64.Instruments {
	// return asyncint64Instruments{meter: m}
	return nil
}

func (m *meter) AsyncFloat64() asyncfloat64.Instruments {
	// return asyncfloat64Instruments{meter: m}
	return nil
}

func (m *meter) newInstrument(name string, opts []apiInstrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) {

// cfg := apiInstrument.NewConfig(opts...)
// descriptor := sdkapi.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())
// name string, opts []apiInstrument.Option, nk number.Kind, ik sdkapi.InstrumentKind)
// inst := a.instruments[name]
// if inst != nil {
// 	if inst.descriptor.NumberKind() == nk && inst.descriptor.InstrumentKind() == ik && inst.descriptor.Unit() == cfg.Unit() {
// 		return inst, nil
// 	}
// 	return nil, ErrIncompatibleInstruments
// }

// ErrIncompatibleInstruments = fmt.Errorf("incompatible instrument registration")

// func (m *meter) NewInstrument(descriptor sdkapi.Descriptor) (sdkapi.Instrument, error) {
// 	cfactory, err := m.views.NewFactory(descriptor)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if descriptor.InstrumentKind().Synchronous() {
// 		return m.syncAccum.NewInstrument(descriptor, cfactory)
// 	}
// 	return m.asyncAccum.NewInstrument(descriptor, cfactory)
// }

func (m *meter) NewCallback(insts []sdkapi.Instrument, callback func(context.Context) error) (sdkapi.Callback, error) {
	return m.asyncAccum.NewCallback(insts, callback)
}
