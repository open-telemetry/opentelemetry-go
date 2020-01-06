package internal

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
)

// This file contains the forwarding implementation of metric.Provider
// used as the default global instance.  Metric events using instruments
// provided by this implementation are no-ops until the first Meter
// implementation is set as the global provider.
//
// The implementation here uses Mutexes to maintain a list of active
// Meters in the Provider and Instruments in each Meter, under the
// assumption that these interfaces are not performance-critical.
//
// We have the invariant that setDelegate() will be called before a
// new metric.Provider implementation is registered as the global
// provider.  Mutexes in the Provider and Meters ensure that each
// instrument has a delegate before the global provider is set.
//
// LabelSets are implemented by delegating to the Meter instance using
// the metric.LabelSetDelegator interface.
//
// Bound instrument operations are implemented by delegating to the
// instrument after it is registered, with a sync.Once initializer to
// protect against races with Release().

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

	delegate unsafe.Pointer // (*metric.Meter)
}

type instImpl struct {
	name  string
	mkind metricKind
	nkind core.NumberKind
	opts  interface{}

	delegate unsafe.Pointer // (*metric.InstrumentImpl)
}

type labelSet struct {
	meter *meter
	value []core.KeyValue

	initialize sync.Once
	delegate   unsafe.Pointer // (* metric.LabelSet)
}

type instHandle struct {
	inst   *instImpl
	labels metric.LabelSet

	initialize sync.Once
	delegate   unsafe.Pointer // (*metric.HandleImpl)
}

var _ metric.Provider = &meterProvider{}
var _ metric.Meter = &meter{}
var _ metric.LabelSet = &labelSet{}
var _ metric.LabelSetDelegate = &labelSet{}
var _ metric.InstrumentImpl = &instImpl{}
var _ metric.BoundInstrumentImpl = &instHandle{}

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
		}
		return m.NewFloat64Counter(name, opts.([]metric.CounterOptionApplier)...).Impl()
	case gaugeKind:
		if nkind == core.Int64NumberKind {
			return m.NewInt64Gauge(name, opts.([]metric.GaugeOptionApplier)...).Impl()
		}
		return m.NewFloat64Gauge(name, opts.([]metric.GaugeOptionApplier)...).Impl()
	case measureKind:
		if nkind == core.Int64NumberKind {
			return m.NewInt64Measure(name, opts.([]metric.MeasureOptionApplier)...).Impl()
		}
		return m.NewFloat64Measure(name, opts.([]metric.MeasureOptionApplier)...).Impl()
	}
	return nil
}

// Instrument delegation

func (inst *instImpl) setDelegate(d metric.Meter) {
	implPtr := new(metric.InstrumentImpl)

	*implPtr = newInstDelegate(d, inst.name, inst.mkind, inst.nkind, inst.opts)

	atomic.StorePointer(&inst.delegate, unsafe.Pointer(implPtr))
}

func (inst *instImpl) Bind(labels metric.LabelSet) metric.BoundInstrumentImpl {
	if implPtr := (*metric.InstrumentImpl)(atomic.LoadPointer(&inst.delegate)); implPtr != nil {
		return (*implPtr).Bind(labels)
	}
	return &instHandle{
		inst:   inst,
		labels: labels,
	}
}

func (bound *instHandle) Unbind() {
	bound.initialize.Do(func() {})

	implPtr := (*metric.BoundInstrumentImpl)(atomic.LoadPointer(&bound.delegate))

	if implPtr == nil {
		return
	}

	(*implPtr).Unbind()
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

// Bound instrument initialization

func (bound *instHandle) RecordOne(ctx context.Context, number core.Number) {
	instPtr := (*metric.InstrumentImpl)(atomic.LoadPointer(&bound.inst.delegate))
	if instPtr == nil {
		return
	}
	var implPtr *metric.BoundInstrumentImpl
	bound.initialize.Do(func() {
		implPtr = new(metric.BoundInstrumentImpl)
		*implPtr = (*instPtr).Bind(bound.labels)
		atomic.StorePointer(&bound.delegate, unsafe.Pointer(implPtr))
	})
	if implPtr == nil {
		implPtr = (*metric.BoundInstrumentImpl)(atomic.LoadPointer(&bound.delegate))
	}
	(*implPtr).RecordOne(ctx, number)
}

// LabelSet initialization

func (m *meter) Labels(labels ...core.KeyValue) metric.LabelSet {
	return &labelSet{
		meter: m,
		value: labels,
	}
}

func (labels *labelSet) Delegate() metric.LabelSet {
	meterPtr := (*metric.Meter)(atomic.LoadPointer(&labels.meter.delegate))
	if meterPtr == nil {
		// This is technically impossible, provided the global
		// Meter is updated after the meters and instruments
		// have been delegated.
		return labels
	}
	var implPtr *metric.LabelSet
	labels.initialize.Do(func() {
		implPtr = new(metric.LabelSet)
		*implPtr = (*meterPtr).Labels(labels.value...)
		atomic.StorePointer(&labels.delegate, unsafe.Pointer(implPtr))
	})
	if implPtr == nil {
		implPtr = (*metric.LabelSet)(atomic.LoadPointer(&labels.delegate))
	}
	return (*implPtr)
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
