package metric

import (
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/asyncfloat64"
	"go.opentelemetry.io/otel/metric/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/syncfloat64"
	"go.opentelemetry.io/otel/metric/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res     *resource.Resource
		readers []reader.Reader
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
		registry   *registry.State
		syncAccum  *syncstate.Accumulator
		asyncAccum *asyncstate.Accumulator
		views      *viewstate.State
	}
)

var (
	_ metric.Meter = &meter{}
)

func WithResource(res *resource.Resource) Option {
	return func(cfg *Config) {
		cfg.res = res
	}
}

func WithReader(r reader.Reader) Option {
	return func(cfg *Config) {
		cfg.readers = append(cfg.readers, r)
	}
}

func New(opts ...Option) metric.MeterProvider {
	cfg := Config{
		res: resource.Default(),
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
		syncAccum:  syncstate.New(),
		asyncAccum: asyncstate.New(),
		views:      viewstate.New(lib, p.cfg.readers),
	}
	p.meters[lib] = m
	return m
}

func (m *meter) SyncInt64() syncint64.Instruments {
	return m.syncAccum.Int64Instruments(m.registry, m.views)
}

func (m *meter) SyncFloat64() syncfloat64.Instruments {
	return m.syncAccum.Float64Instruments(m.registry, m.views)
}

func (m *meter) AsyncInt64() asyncint64.Instruments {
	return m.asyncAccum.Int64Instruments(m.registry, m.views)
}

func (m *meter) AsyncFloat64() asyncfloat64.Instruments {
	return m.asyncAccum.Float64Instruments(m.registry, m.views)
}

func (m *meter) NewCallback(insts []instrument.Asynchronous, function metric.CallbackFunc) (metric.Callback, error) {
	return m.asyncAccum.NewCallback(insts, function)
}
