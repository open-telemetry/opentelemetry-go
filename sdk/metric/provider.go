package metric

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/pipeline"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/view"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Config struct {
		res     *resource.Resource
		readers []Reader
		views   []*view.Views
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
		pipe        int
		lastCollect time.Time
	}

	meter struct {
		library   instrumentation.Library
		provider  *Provider
		compilers pipeline.Register[*viewstate.Compiler]

		lock       sync.Mutex
		byDesc     map[sdkinstrument.Descriptor]interface{}
		syncInsts  []*syncstate.Instrument
		asyncInsts []*asyncstate.Instrument
		callbacks  []*asyncstate.Callback
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

func WithReader(r Reader, opts ...view.Option) Option {
	return func(cfg *Config) {
		cfg.readers = append(cfg.readers, r)
		cfg.views = append(cfg.views, view.New(r.String(), opts...))
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
	for pipe := 0; pipe < len(cfg.readers); pipe++ {
		cfg.readers[pipe].Register(p.producerFor(pipe))
	}
	return p
}

func (p *Provider) producerFor(pipe int) Producer {
	return &providerProducer{
		provider:    p,
		pipe:        pipe,
		lastCollect: p.startTime,
	}
}

func (pp *providerProducer) Produce(inout *data.Metrics) data.Metrics {
	ordered := pp.provider.getOrdered()

	// Note: the Last time is only used in delta-temporality
	// scenarios.  This lock protects the only stateful change in
	// `pp` but does not prevent concurrent collection.  If a
	// delta-temporality reader were to call Produce
	// concurrently, the results would be be recorded with
	// non-overlapping timestamps but would have been collected in
	// an overlapping way.
	pp.lock.Lock()
	lastTime := pp.lastCollect
	nowTime := time.Now()
	pp.lastCollect = nowTime
	pp.lock.Unlock()

	var output data.Metrics
	if inout != nil {
		inout.Reset()
		output = *inout
	}

	output.Resource = pp.provider.cfg.res

	sequence := data.Sequence{
		Start: pp.provider.startTime,
		Last:  lastTime,
		Now:   nowTime,
	}

	// TODO: Add a timeout to the context.
	ctx := context.Background()

	for _, meter := range ordered {
		meter.collectFor(
			ctx,
			pp.pipe,
			sequence,
			&output,
		)
	}

	return output
}

func (m *meter) collectFor(ctx context.Context, pipe int, seq data.Sequence, output *data.Metrics) {
	// Use m.lock to briefly access the current lists: syncInsts, asyncInsts, callbacks
	m.lock.Lock()
	syncInsts := m.syncInsts
	asyncInsts := m.asyncInsts
	callbacks := m.callbacks
	m.lock.Unlock()

	asyncState := asyncstate.NewState(pipe)

	for _, cb := range callbacks {
		cb.Run(ctx, asyncState)
	}

	for _, inst := range syncInsts {
		inst.SnapshotAndProcess()
	}

	for _, inst := range asyncInsts {
		inst.SnapshotAndProcess(asyncState)
	}

	scope := data.ReallocateFrom(&output.Scopes)
	scope.Library = m.library

	for _, coll := range m.compilers[pipe].Collectors() {
		coll.Collect(seq, &scope.Instruments)
	}
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
		provider:  p,
		library:   lib,
		byDesc:    map[sdkinstrument.Descriptor]interface{}{},
		compilers: pipeline.NewRegister[*viewstate.Compiler](len(p.cfg.readers)),
	}
	for pipe := range m.compilers {
		m.compilers[pipe] = viewstate.New(lib, p.cfg.views[pipe])
	}
	p.ordered = append(p.ordered, m)
	p.meters[lib] = m
	return m
}

func (m *meter) RegisterCallback(insts []instrument.Asynchronous, function func(context.Context)) error {
	cb, err := asyncstate.NewCallback(insts, m, function)

	if err == nil {
		m.lock.Lock()
		defer m.lock.Unlock()
		m.callbacks = append(m.callbacks, cb)
	}
	return err
}
