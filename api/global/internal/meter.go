package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
)

type metricKind int8

const (
	counterKind metricKind = iota
	gaugeKind
	measureKind
)

type meterProvider struct {
	lock     sync.Mutex
	meters   []*meter
	delegate metric.Provider
}

type meter struct {
	provider *meterProvider
	name     string

	lock        sync.Mutex
	instruments []*instImpl
	delegate    unsafe.Pointer // (*metric.Meter)
}

type labelSet struct {
	meter  *meter
	labels []core.KeyValue

	delegate metric.LabelSet
}

type instImpl struct {
	meter *meter
	name  string
	mkind metricKind
	nkind core.NumberKind
	opts  interface{}

	delegate unsafe.Pointer // (*metric.InstrumentImpl)
}

type instHandle struct {
	inst   *instImpl
	labels metric.LabelSet

	delegate unsafe.Pointer // (*metric.HandleImpl)
}

var _ metric.Provider = &meterProvider{}
var _ metric.Meter = &meter{}
var _ metric.LabelSet = &labelSet{}
var _ metric.InstrumentImpl = &instImpl{}
var _ metric.HandleImpl = &instHandle{}

// Provider interface and delegation

func (p *meterProvider) setDelegate(provider metric.Provider) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.delegate = provider
	for _, m := range p.meters {
		m.setDelegate(provider)
	}
	p.meters = nil
}

func (p *meterProvider) Meter(name string) metric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.delegate != nil {
		return p.delegate.Meter(name)
	}

	m := &meter{
		provider: p,
		name:     name,
	}
	p.meters = append(p.meters, m)
	return m
}

// Meter interface and delegation

func (m *meter) setDelegate(provider metric.Provider) {
	m.lock.Lock()
	defer m.lock.Unlock()

	d := new(metric.Meter)
	*d = provider.Meter(m.name)
	m.delegate = unsafe.Pointer(d)

	for _, inst := range m.instruments {
		inst.setDelegate(*d)
	}
	m.instruments = nil
}

func (m *meter) newInst(name string, mkind metricKind, nkind core.NumberKind, opts interface{}) metric.InstrumentImpl {
	m.lock.Lock()
	defer m.lock.Unlock()

	if meterPtr := (*metric.Meter)(atomic.LoadPointer(&m.delegate)); meterPtr != nil {
		return newInstDelegate(*meterPtr, name, mkind, nkind, opts)
	}

	inst := &instImpl{
		meter: m,
		name:  name,
		mkind: mkind,
		nkind: nkind,
		opts:  opts,
	}
	m.instruments = append(m.instruments, inst)
	return inst
}

func newInstDelegate(m metric.Meter, name string, mkind metricKind, nkind core.NumberKind, opts interface{}) metric.InstrumentImpl {
	switch mkind {
	case counterKind:
		if nkind == core.Int64NumberKind {
			return m.NewInt64Counter(name, opts.([]metric.CounterOptionApplier)...).Impl()
		} else {
			return m.NewFloat64Counter(name, opts.([]metric.CounterOptionApplier)...).Impl()
		}
	case gaugeKind:
		if nkind == core.Int64NumberKind {
			return m.NewInt64Gauge(name, opts.([]metric.GaugeOptionApplier)...).Impl()
		} else {
			return m.NewFloat64Gauge(name, opts.([]metric.GaugeOptionApplier)...).Impl()
		}
	case measureKind:
		if nkind == core.Int64NumberKind {
			return m.NewInt64Measure(name, opts.([]metric.MeasureOptionApplier)...).Impl()
		} else {
			return m.NewFloat64Measure(name, opts.([]metric.MeasureOptionApplier)...).Impl()
		}
	}
	return nil
}

// Instrument delegation

func (inst *instImpl) setDelegate(d metric.Meter) {
	impl := new(metric.InstrumentImpl)

	*impl = newInstDelegate(d, inst.name, inst.mkind, inst.nkind, inst.opts)

	atomic.StorePointer(&inst.delegate, unsafe.Pointer(impl))
}

func (inst *instImpl) AcquireHandle(labels metric.LabelSet) metric.HandleImpl {
	if implPtr := (*metric.InstrumentImpl)(atomic.LoadPointer(&inst.delegate)); implPtr != nil {
		return (*implPtr).AcquireHandle(labels)
	}

	return &instHandle{
		inst:   inst,
		labels: labels,
	}
}

func (bound *instHandle) Release() {
	// TODO
}

// Metric updates

func (m *meter) RecordBatch(ctx context.Context, labels metric.LabelSet, measurements ...metric.Measurement) {
	if delegatePtr := (*metric.Meter)(atomic.LoadPointer(&m.delegate)); delegatePtr != nil {
		(*delegatePtr).RecordBatch(ctx, labels, measurements...)
	}
}

func (inst *instImpl) RecordOne(ctx context.Context, number core.Number, labels metric.LabelSet) {
	if instPtr := (*metric.InstrumentImpl)(atomic.LoadPointer(&inst.delegate)); instPtr != nil {
		(*instPtr).RecordOne(ctx, number, labels)
	}
}

func (*instHandle) RecordOne(ctx context.Context, number core.Number) {
	// TODO
}

// Constructors

func (m *meter) NewInt64Counter(name string, opts ...metric.CounterOptionApplier) metric.Int64Counter {
	return metric.WrapInt64CounterInstrument(m.newInst(name, counterKind, core.Int64NumberKind, opts))
}

func (m *meter) NewFloat64Counter(name string, opts ...metric.CounterOptionApplier) metric.Float64Counter {
	return metric.WrapFloat64CounterInstrument(m.newInst(name, counterKind, core.Float64NumberKind, opts))
}

func (m *meter) NewInt64Gauge(name string, opts ...metric.GaugeOptionApplier) metric.Int64Gauge {
	return metric.WrapInt64GaugeInstrument(m.newInst(name, gaugeKind, core.Int64NumberKind, opts))
}

func (m *meter) NewFloat64Gauge(name string, opts ...metric.GaugeOptionApplier) metric.Float64Gauge {
	return metric.WrapFloat64GaugeInstrument(m.newInst(name, gaugeKind, core.Float64NumberKind, opts))
}

func (m *meter) NewInt64Measure(name string, opts ...metric.MeasureOptionApplier) metric.Int64Measure {
	return metric.WrapInt64MeasureInstrument(m.newInst(name, measureKind, core.Int64NumberKind, opts))
}

func (m *meter) NewFloat64Measure(name string, opts ...metric.MeasureOptionApplier) metric.Float64Measure {
	return metric.WrapFloat64MeasureInstrument(m.newInst(name, measureKind, core.Float64NumberKind, opts))
}

// TODO

func (m *meter) Labels(labels ...core.KeyValue) metric.LabelSet {
	return &labelSet{
		meter:  m,
		labels: labels,
	}
}
