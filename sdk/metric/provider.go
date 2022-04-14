package metric

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/views"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res     *resource.Resource
		readers []*reader.ReaderConfig
		views   []views.View
	}

	Option func(cfg *Config)

	Provider struct {
		cfg       Config
		startTime time.Time
		lock      sync.Mutex
		ordered   []*meter
		meters    map[instrumentation.Library]*meter
	}

	providerProducer struct {
		lock        sync.Mutex
		provider    *Provider
		reader      *reader.ReaderConfig
		lastCollect time.Time
	}

	meter struct {
		library  instrumentation.Library
		provider *Provider
		views    *viewstate.Compiler

		lock        sync.Mutex
		byDesc      map[sdkinstrument.Descriptor]instrumentIface
		callbacks   []*asyncstate.Callback
		instruments []instrumentIface
	}

	instrumentIface interface {
		AccumulateFor(*reader.ReaderConfig)
	}
)

var (
	_ metric.Meter         = &meter{}
	_ metric.MeterProvider = &Provider{}
)

func WithResource(res *resource.Resource) Option {
	return func(cfg *Config) {
		cfg.res = res
	}
}

func WithReader(r reader.Reader, opts ...reader.Option) Option {
	return func(cfg *Config) {
		rConfig := reader.NewConfig(r, opts...)
		cfg.readers = append(cfg.readers, rConfig)
	}
}

func WithViews(vs ...views.View) Option {
	return func(cfg *Config) {
		cfg.views = append(cfg.views, vs...)
	}
}

func New(opts ...Option) *Provider {
	cfg := Config{
		res: resource.Default(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	p := &Provider{
		cfg:       cfg,
		startTime: time.Now(),
		meters:    map[instrumentation.Library]*meter{},
	}
	for _, reader := range cfg.readers {
		reader.Reader().Register(p.producerFor(reader))
	}
	return p
}

func (p *Provider) producerFor(r *reader.ReaderConfig) reader.Producer {
	return &providerProducer{
		provider:    p,
		reader:      r,
		lastCollect: p.startTime,
	}
}

func (pp *providerProducer) Produce(ctx context.Context, inout *reader.Metrics) reader.Metrics {
	ordered := pp.provider.getOrdered()

	// Note: the Last time is only used in delta-temporality
	// scenarios.  This lock protects the only stateful change in
	// `pp` but does not prevent concurrent collection.  If a
	// delta-temporality exporter were to call Produce
	// concurrently, the results would be be recorded with
	// non-overlapping timestamps but would have been collected in
	// an overlapping way.
	pp.lock.Lock()
	lastTime := pp.lastCollect
	nowTime := time.Now()
	pp.lastCollect = nowTime
	pp.lock.Unlock()

	var output reader.Metrics
	if inout != nil {
		inout.Reset()
		output = *inout
	}

	output.Resource = pp.provider.cfg.res

	sequence := reader.Sequence{
		Start: pp.provider.startTime,
		Last:  lastTime,
		Now:   nowTime,
	}

	for _, meter := range ordered {
		meter.lock.Lock()
		callbacks := meter.callbacks
		instruments := meter.instruments
		meter.lock.Unlock()

		for _, cb := range callbacks {
			cb.Run(ctx, pp.reader)
		}

		for _, inst := range instruments {
			inst.AccumulateFor(pp.reader)
		}

		scope := reader.Reallocate(&output.Scopes)

		scope.Library = meter.library

		for _, coll := range meter.views.Collectors(pp.reader) {
			coll.Collect(pp.reader, sequence, &scope.Instruments)
		}
	}

	return output
}

func (p *Provider) getOrdered() []*meter {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.ordered
}

func (p *Provider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
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
		provider: p,
		library:  lib,
		byDesc:   map[sdkinstrument.Descriptor]instrumentIface{},
		views:    viewstate.New(lib, p.cfg.views, p.cfg.readers),
	}
	p.ordered = append(p.ordered, m)
	p.meters[lib] = m
	return m
}

func (m *meter) RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error {
	cb, err := asyncstate.NewCallback(insts, function)

	if err == nil {
		m.lock.Lock()
		defer m.lock.Unlock()
		m.callbacks = append(m.callbacks, cb)
	}
	return err
}

var errLookupConflicts = fmt.Errorf("caller should lookup conflict information")

func configureInstrument[T instrumentIface](
	m *meter,
	name string,
	opts []instrument.Option,
	nk number.Kind,
	ik sdkinstrument.Kind,
	create func(desc sdkinstrument.Descriptor) (T, error),
) (T, error) {
	cfg := instrument.NewConfig(opts...)
	desc := sdkinstrument.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())
	m.lock.Lock()
	defer m.lock.Unlock()
	if lookup, has := m.byDesc[desc]; has {
		// Note: Recomputing conflicts
		_, err := create(desc)
		return lookup.(T), err
	}
	inst, err := create(desc)
	m.byDesc[desc] = inst
	m.instruments = append(m.instruments, inst)
	return inst, err
}
