package metric

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/registry"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/views"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res     *resource.Resource
		readers []*reader.Reader
		views   []views.View
	}

	Option func(cfg *Config)

	provider struct {
		cfg         Config
		collectLock sync.Mutex
		lock        sync.Mutex
		ordered     []*meter
		meters      map[instrumentation.Library]*meter
	}

	providerProducer struct {
		lock     sync.Mutex
		provider *provider
		reader   *reader.Reader
	}

	meter struct {
		library    instrumentation.Library
		provider   *provider
		registry   *registry.State
		syncAccum  *syncstate.Accumulator
		asyncAccum *asyncstate.Accumulator
		views      *viewstate.Compiler
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

func WithReader(r *reader.Reader) Option {
	return func(cfg *Config) {
		cfg.readers = append(cfg.readers, r)
	}
}

func WithViews(vs ...views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, vs...)
	}
}

func New(opts ...Option) metric.MeterProvider {
	cfg := Config{
		res: resource.Default(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	p := &provider{
		cfg:    cfg,
		meters: map[instrumentation.Library]*meter{},
	}
	for _, reader := range cfg.readers {
		reader.Exporter().Register(p.producerFor(reader))
	}
	return p
}

func (p *provider) producerFor(reader *reader.Reader) reader.Producer {
	return &providerProducer{
		provider: p,
		reader:   reader,
	}
}

func (pp *providerProducer) Produce() reader.Metrics {
	pp.lock.Lock()
	defer pp.lock.Unlock()

	ordered := pp.provider.getOrdered()

	output := reader.Metrics{
		Resource: pp.provider.cfg.res,
		Scopes:   make([]reader.Scope, len(ordered)),
	}

	for idx, meter := range ordered {
		output.Scopes[idx].Library = meter.library

		meter.asyncAccum.Collect(pp.reader, &output.Scopes[idx].Instruments)
		meter.syncAccum.Collect(pp.reader, &output.Scopes[idx].Instruments)
	}

	return output
}

func (p *provider) getOrdered() []*meter {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.ordered
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
		views:      viewstate.New(lib, p.cfg.views, p.cfg.readers),
	}
	p.ordered = append(p.ordered, m)
	p.meters[lib] = m
	return m
}

func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	return m.syncAccum.Int64Instruments(m.registry, m.views)
}

func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	return m.syncAccum.Float64Instruments(m.registry, m.views)
}

func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	return m.asyncAccum.Int64Instruments(m.registry, m.views)
}

func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	return m.asyncAccum.Float64Instruments(m.registry, m.views)
}

func (m *meter) RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error {
	return m.asyncAccum.RegisterCallback(insts, function)
}
