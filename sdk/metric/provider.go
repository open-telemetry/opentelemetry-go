package metric

import (
	"context"
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
		readers []*reader.Reader
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
		reader      *reader.Reader
		lastCollect time.Time
	}

	instrumentIface interface {
		Descriptor() sdkinstrument.Descriptor
		Collect(r *reader.Reader, seq reader.Sequence, output *[]reader.Instrument)
	}

	meter struct {
		library  instrumentation.Library
		provider *Provider
		names    map[string][]instrumentIface
		views    *viewstate.Compiler

		lock        sync.Mutex
		instruments []instrumentIface
		callbacks   []*asyncstate.Callback
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
		reader.Exporter().Register(p.producerFor(reader))
	}
	return p
}

func (p *Provider) producerFor(r *reader.Reader) reader.Producer {
	return &providerProducer{
		provider:    p,
		reader:      r,
		lastCollect: p.startTime,
	}
}

func resetMetrics(m *reader.Metrics) {
	for i := range m.Scopes {
		resetScope(&m.Scopes[i])
	}
	m.Scopes = m.Scopes[0:0:cap(m.Scopes)]
}

func resetScope(s *reader.Scope) {
	for i := range s.Instruments {
		resetInstrument(&s.Instruments[i])
	}
	s.Instruments = s.Instruments[0:0:cap(s.Instruments)]
}

func resetInstrument(inst *reader.Instrument) {
	inst.Series = inst.Series[0:0:cap(inst.Series)]
}

func appendScope(scopes *[]reader.Scope) *reader.Scope {
	// Note: there's a generic form of this logic in internal/viewstate,
	// should this use it?
	if len(*scopes) < cap(*scopes) {
		(*scopes) = (*scopes)[0 : len(*scopes)+1 : cap(*scopes)]
	} else {
		(*scopes) = append(*scopes, reader.Scope{})
	}
	return &(*scopes)[len(*scopes)-1]
}

func (pp *providerProducer) Produce(inout *reader.Metrics) reader.Metrics {
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
		resetMetrics(inout)
		output = *inout
	}

	output.Resource = pp.provider.cfg.res

	sequence := reader.Sequence{
		Start: pp.provider.startTime,
		Last:  lastTime,
		Now:   nowTime,
	}

	// TODO: Add a timeout to the context.
	ctx := context.Background()

	for _, meter := range ordered {
		meter.lock.Lock()
		callbacks := meter.callbacks
		instruments := meter.instruments
		meter.lock.Unlock()

		for _, cb := range callbacks {
			cb.Run(ctx, pp.reader)
		}

		scope := appendScope(&output.Scopes)

		scope.Library = meter.library

		for _, inst := range instruments {
			inst.Collect(pp.reader, sequence, &scope.Instruments)
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
		names:    map[string][]instrumentIface{},
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

func nameLookup[T instrumentIface](
	m *meter,
	name string,
	opts []instrument.Option,
	nk number.Kind,
	ik sdkinstrument.Kind,
	f func(desc sdkinstrument.Descriptor) T,
) (T, error) {
	cfg := instrument.NewConfig(opts...)
	desc := sdkinstrument.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())

	m.lock.Lock()
	defer m.lock.Unlock()
	lookup := m.names[name]

	for _, found := range lookup {
		match, ok := found.(T)
		if !ok {
			continue
		}

		exist := found.Descriptor()

		if exist.NumberKind != nk || exist.Kind != ik || exist.Unit != cfg.Unit() {
			continue
		}

		// Exact match (ignores description)
		return match, nil
	}
	value := f(desc)
	m.names[name] = append(m.names[name], value)
	return value, nil
}
