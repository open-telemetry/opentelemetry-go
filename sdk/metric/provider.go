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

// MeterProvider handles the creation and coordination of Meters. All Meters
// created by a MeterProvider will be associated with the same Resource, have
// the same Views applied to them, and have their produced metric telemetry
// passed to the configured Readers.
type MeterProvider struct {
	cfg       Config
	startTime time.Time
	lock      sync.Mutex
	ordered   []*meter
	meters    map[instrumentation.Library]*meter
}

// Compile-time check MeterProvider implements metric.MeterProvider.
var _ metric.MeterProvider = (*MeterProvider)(nil)

// NewMeterProvider returns a new and configured MeterProvider.
//
// By default, the returned MeterProvider is configured with the default
// Resource and no Readers. Readers cannot be added after a MeterProvider is
// created. This means the returned MeterProvider, one created with no
// Readers, will be perform no operations.
func NewMeterProvider(options ...Option) *MeterProvider {
	cfg := Config{
		res: resource.Default(),
	}
	for _, opt := range options {
		opt(&cfg)
	}
	p := &MeterProvider{
		cfg:       cfg,
		startTime: time.Now(),
		meters:    map[instrumentation.Library]*meter{},
	}
	for pipe := 0; pipe < len(cfg.readers); pipe++ {
		cfg.readers[pipe].Register(p.producerFor(pipe))
	}
	return p
}

// Meter returns a Meter with the given name and configured with options.
//
// The name should be the name of the instrumentation scope creating
// telemetry. This name may be the same as the instrumented code only if that
// code provides built-in instrumentation.
//
// If name is empty, the default (go.opentelemetry.io/otel/sdk/meter) will be
// used.
//
// Calls to the Meter method after Shutdown has been called will return Meters
// that perform no operations.
//
// This method is safe to call concurrently.
func (mp *MeterProvider) Meter(name string, options ...metric.MeterOption) metric.Meter {
	cfg := metric.NewMeterConfig(opts...)
	lib := instrumentation.Library{
		Name:      name,
		Version:   cfg.InstrumentationVersion(),
		SchemaURL: cfg.SchemaURL(),
	}

	mp.lock.Lock()
	defer mp.lock.Unlock()

	m := mp.meters[lib]
	if m != nil {
		return m
	}
	m = &meter{
		provider:  p,
		library:   lib,
		byDesc:    map[sdkinstrument.Descriptor]interface{}{},
		compilers: pipeline.NewRegister[*viewstate.Compiler](len(mp.cfg.readers)),
	}
	for pipe := range m.compilers {
		m.compilers[pipe] = viewstate.New(lib, mp.cfg.views[pipe])
	}
	mp.ordered = append(mp.ordered, m)
	mp.meters[lib] = m
	return m
}

type (
	Config struct {
		res     *resource.Resource
		readers []Reader
		views   []*view.Views
	}

	Option func(cfg *Config)

	providerProducer struct {
		lock        sync.Mutex
		provider    *MeterProvider
		pipe        int
		lastCollect time.Time
	}

	meter struct {
		library   instrumentation.Library
		provider  *MeterProvider
		compilers pipeline.Register[*viewstate.Compiler]

		lock       sync.Mutex
		byDesc     map[sdkinstrument.Descriptor]interface{}
		syncInsts  []*syncstate.Instrument
		asyncInsts []*asyncstate.Instrument
		callbacks  []*asyncstate.Callback
	}
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

func (mp *MeterProvider) producerFor(pipe int) Producer {
	return &providerProducer{
		provider:    p,
		pipe:        pipe,
		lastCollect: mp.startTime,
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

func (p *MeterProvider) getOrdered() []*meter {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.ordered
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
