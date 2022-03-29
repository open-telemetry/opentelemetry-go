package metric

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/internal/asyncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/syncstate"
	"go.opentelemetry.io/otel/sdk/metric/internal/viewstate"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
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
		cfg       Config
		startTime time.Time
		lock      sync.Mutex
		ordered   []*meter
		meters    map[instrumentation.Library]*meter
	}

	providerProducer struct {
		lock        sync.Mutex
		provider    *provider
		reader      *reader.Reader
		lastCollect time.Time
	}

	instrumentIface interface {
		Descriptor() sdkapi.Descriptor
		Collect(r *reader.Reader, seq reader.Sequence, 
	}

	meter struct {
		library  instrumentation.Library
		provider *provider
		names    map[string][]instrumentIface
		views    *viewstate.Compiler

		lock        sync.Mutex
		instruments []instrumentIface
		callbacks   []*asyncstate.Callback
	}

	asyncint64Instruments   struct{ *meter }
	asyncfloat64Instruments struct{ *meter }
	syncint64Instruments    struct{ *meter }
	syncfloat64Instruments  struct{ *meter }
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
		cfg:       cfg,
		startTime: time.Now(),
		meters:    map[instrumentation.Library]*meter{},
	}
	for _, reader := range cfg.readers {
		reader.Exporter().Register(p.producerFor(reader))
	}
	return p
}

func (p *provider) producerFor(r *reader.Reader) reader.Producer {
	return &providerProducer{
		provider: p,
		reader:   r,
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

	sequence := viewstate.Sequence{
		Start: pp.provider.startTime,
		Last:  pp.lastCollect,
		Now:   time.Now(),
	}

	// TODO: Add a timeout to the context.
	ctx := context.Background()

	for meterIdx, meter := range ordered {
		// Lock
		meter.lock.Lock()
		callbacks := meter.callbacks
		instruments := meter.instruments
		meter.lock.Unlock()

		for _, cb := range callbacks {
			cb.Run(ctx, pp.reader)
		}

		output.Scopes[meterIdx].Library = meter.library

		// Note: the number of output instruments is
		// determined by the views, not the number of actual
		// instruments.

		for _, inst := range instruments {
			inst.Collect(pp.reader, sequence, &output.Scopes[meterIdx].Instruments)
		}
	}

	pp.lastCollect = sequence.Now

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
		provider: p,
		library:  lib,
		names:    map[string][]instrumentIface{},
		views:    viewstate.New(lib, p.cfg.views, p.cfg.readers),
	}
	p.ordered = append(p.ordered, m)
	p.meters[lib] = m
	return m
}

func (m *meter) SyncInt64() syncint64.InstrumentProvider {
	return syncint64Instruments{m}
}

func (m *meter) SyncFloat64() syncfloat64.InstrumentProvider {
	return syncfloat64Instruments{m}
}

func (m *meter) AsyncInt64() asyncint64.InstrumentProvider {
	return asyncint64Instruments{m}
}

func (m *meter) AsyncFloat64() asyncfloat64.InstrumentProvider {
	return asyncfloat64Instruments{m}
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

func (m *meter) newAsyncInst(name string, opts []instrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) (*asyncstate.Instrument, error) {
	return nameLookup(
		m, name, opts, nk, ik,
		func(desc sdkapi.Descriptor) *asyncstate.Instrument {
			compiled := m.views.Compile(desc)
			inst := asyncstate.NewInstrument(desc, compiled)
			m.instruments = append(m.instruments, inst)
			return inst
		})
}

func (i asyncint64Instruments) Counter(name string, opts ...instrument.Option) (asyncint64.Counter, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.CounterObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (i asyncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncint64.UpDownCounter, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (i asyncint64Instruments) Gauge(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
	inst, err := i.newAsyncInst(name, opts, number.Int64Kind, sdkapi.GaugeObserverInstrumentKind)
	return asyncstate.NewObserver[int64, traits.Int64](inst), err
}

func (f asyncfloat64Instruments) Counter(name string, opts ...instrument.Option) (asyncfloat64.Counter, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.CounterObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}

func (f asyncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (asyncfloat64.UpDownCounter, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.UpDownCounterObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}

func (f asyncfloat64Instruments) Gauge(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
	inst, err := f.newAsyncInst(name, opts, number.Float64Kind, sdkapi.GaugeObserverInstrumentKind)
	return asyncstate.NewObserver[float64, traits.Float64](inst), err
}

func nameLookup[T instrumentIface](
	m *meter,
	name string,
	opts []instrument.Option,
	nk number.Kind,
	ik sdkapi.InstrumentKind,
	f func(desc sdkapi.Descriptor) T,
) (T, error) {
	cfg := instrument.NewConfig(opts...)
	desc := sdkapi.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())

	m.lock.Lock()
	defer m.lock.Unlock()
	lookup := m.names[name]

	for _, found := range lookup {
		match, ok := found.(T)
		if !ok {
			continue
		}

		exist := found.Descriptor()

		if exist.NumberKind() != nk || exist.InstrumentKind() != ik || exist.Unit() != cfg.Unit() {
			continue
		}

		// Exact match (ignores description)
		return match, nil
	}
	value := f(desc)
	m.names[name] = append(m.names[name], value)
	return value, nil
}

func (m *meter) newSyncInst(name string, opts []instrument.Option, nk number.Kind, ik sdkapi.InstrumentKind) (*syncstate.Instrument, error) {
	return nameLookup(
		m, name, opts, nk, ik,
		func(desc sdkapi.Descriptor) *syncstate.Instrument {
			compiled := m.views.Compile(desc)
			inst := syncstate.NewInstrument(desc, compiled)

			m.instruments = append(m.instruments, inst)
			return inst
		})
}

func (i syncint64Instruments) Counter(name string, opts ...instrument.Option) (syncint64.Counter, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkapi.CounterInstrumentKind)
	return syncstate.NewCounter[int64, traits.Int64](inst), err
}

func (i syncint64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncint64.UpDownCounter, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkapi.UpDownCounterInstrumentKind)
	return syncstate.NewCounter[int64, traits.Int64](inst), err
}

func (i syncint64Instruments) Histogram(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
	inst, err := i.newSyncInst(name, opts, number.Int64Kind, sdkapi.HistogramInstrumentKind)
	return syncstate.NewHistogram[int64, traits.Int64](inst), err
}

func (f syncfloat64Instruments) Counter(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkapi.CounterInstrumentKind)
	return syncstate.NewCounter[float64, traits.Float64](inst), err
}

func (f syncfloat64Instruments) UpDownCounter(name string, opts ...instrument.Option) (syncfloat64.UpDownCounter, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkapi.UpDownCounterInstrumentKind)
	return syncstate.NewCounter[float64, traits.Float64](inst), err
}

func (f syncfloat64Instruments) Histogram(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
	inst, err := f.newSyncInst(name, opts, number.Float64Kind, sdkapi.HistogramInstrumentKind)
	return syncstate.NewHistogram[float64, traits.Float64](inst), err
}
