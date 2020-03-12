package internal

import (
	"context"
	"errors"
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
	measureKind
)

type meterProvider struct {
	delegate metric.Provider

	lock   sync.Mutex
	meters []*meter
}

type meter struct {
	delegate unsafe.Pointer // (*metric.Meter)

	provider *meterProvider
	name     string

	lock          sync.Mutex
	instruments   []*instImpl
	liveObservers map[*obsImpl]struct{}
	// orderedObservers slice contains observers in their order of
	// registration. It may also contain unregistered
	// observers. The liveObservers map should be consulted to
	// check if the observer is registered or not.
	orderedObservers []*obsImpl
}

type instImpl struct {
	delegate unsafe.Pointer // (*metric.InstrumentImpl)

	name  string
	mkind metricKind
	nkind core.NumberKind
	opts  []metric.Option
}

type obsImpl struct {
	delegate unsafe.Pointer // (*metric.Int64Observer or *metric.Float64Observer)

	name     string
	nkind    core.NumberKind
	opts     []metric.Option
	meter    *meter
	callback interface{}
}

type hasImpl interface {
	Impl() metric.InstrumentImpl
}

type int64ObsImpl struct {
	observer *obsImpl
}

type float64ObsImpl struct {
	observer *obsImpl
}

// this is a common subset of the metric observers interfaces
type observerUnregister interface {
	Unregister()
}

type labelSet struct {
	delegate unsafe.Pointer // (* metric.LabelSet)

	meter *meter
	value []core.KeyValue

	initialize sync.Once
}

type instHandle struct {
	delegate unsafe.Pointer // (*metric.HandleImpl)

	inst   *instImpl
	labels metric.LabelSet

	initialize sync.Once
}

var _ metric.Provider = &meterProvider{}
var _ metric.Meter = &meter{}
var _ metric.LabelSet = &labelSet{}
var _ metric.LabelSetDelegate = &labelSet{}
var _ metric.InstrumentImpl = &instImpl{}
var _ metric.BoundInstrumentImpl = &instHandle{}
var _ metric.Int64Observer = int64ObsImpl{}
var _ metric.Float64Observer = float64ObsImpl{}
var _ observerUnregister = (metric.Int64Observer)(nil)
var _ observerUnregister = (metric.Float64Observer)(nil)

var errInvalidMetricKind = errors.New("Invalid Metric kind")

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
	for _, obs := range m.orderedObservers {
		if _, ok := m.liveObservers[obs]; ok {
			obs.setDelegate(*d)
		}
	}
	m.liveObservers = nil
	m.orderedObservers = nil
}

func (m *meter) newInst(name string, mkind metricKind, nkind core.NumberKind, opts []metric.Option) (metric.InstrumentImpl, error) {
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
	return inst, nil
}

func delegateCheck(has hasImpl, err error) (metric.InstrumentImpl, error) {
	if has != nil {
		return has.Impl(), err
	}
	if err == nil {
		err = metric.ErrSDKReturnedNilImpl
	}
	return nil, err
}

func newInstDelegate(m metric.Meter, name string, mkind metricKind, nkind core.NumberKind, opts []metric.Option) (metric.InstrumentImpl, error) {
	switch mkind {
	case counterKind:
		if nkind == core.Int64NumberKind {
			return delegateCheck(m.NewInt64Counter(name, opts...))
		}
		return delegateCheck(m.NewFloat64Counter(name, opts...))
	case measureKind:
		if nkind == core.Int64NumberKind {
			return delegateCheck(m.NewInt64Measure(name, opts...))
		}
		return delegateCheck(m.NewFloat64Measure(name, opts...))
	}
	return nil, errInvalidMetricKind
}

// Instrument delegation

func (inst *instImpl) setDelegate(d metric.Meter) {
	implPtr := new(metric.InstrumentImpl)

	var err error
	*implPtr, err = newInstDelegate(d, inst.name, inst.mkind, inst.nkind, inst.opts)

	if err != nil {
		// TODO: There is no standard way to deliver this error to the user.
		// See https://github.com/open-telemetry/opentelemetry-go/issues/514
		// Note that the default SDK will not generate any errors yet, this is
		// only for added safety.
		panic(err)
	}

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

// Any Observer delegation

func (obs *obsImpl) setDelegate(d metric.Meter) {
	if obs.nkind == core.Int64NumberKind {
		obs.setInt64Delegate(d)
	} else {
		obs.setFloat64Delegate(d)
	}
}

func (obs *obsImpl) unregister() {
	unreg := obs.getUnregister()
	if unreg != nil {
		unreg.Unregister()
		return
	}
	obs.meter.lock.Lock()
	defer obs.meter.lock.Unlock()
	delete(obs.meter.liveObservers, obs)
	if len(obs.meter.liveObservers) == 0 {
		obs.meter.liveObservers = nil
		obs.meter.orderedObservers = nil
	}
}

func (obs *obsImpl) getUnregister() observerUnregister {
	ptr := atomic.LoadPointer(&obs.delegate)
	if ptr == nil {
		return nil
	}
	if obs.nkind == core.Int64NumberKind {
		return *(*metric.Int64Observer)(ptr)
	}
	return *(*metric.Float64Observer)(ptr)
}

// Int64Observer delegation

func (obs *obsImpl) setInt64Delegate(d metric.Meter) {
	obsPtr := new(metric.Int64Observer)
	cb := obs.callback.(metric.Int64ObserverCallback)

	var err error
	*obsPtr, err = d.RegisterInt64Observer(obs.name, cb, obs.opts...)

	if err != nil {
		// TODO: There is no standard way to deliver this error to the user.
		// See https://github.com/open-telemetry/opentelemetry-go/issues/514
		// Note that the default SDK will not generate any errors yet, this is
		// only for added safety.
		panic(err)
	}

	atomic.StorePointer(&obs.delegate, unsafe.Pointer(obsPtr))
}

func (obs int64ObsImpl) Unregister() {
	obs.observer.unregister()
}

// Float64Observer delegation

func (obs *obsImpl) setFloat64Delegate(d metric.Meter) {
	obsPtr := new(metric.Float64Observer)
	cb := obs.callback.(metric.Float64ObserverCallback)

	var err error
	*obsPtr, err = d.RegisterFloat64Observer(obs.name, cb, obs.opts...)
	if err != nil {
		// TODO: There is no standard way to deliver this error to the user.
		// See https://github.com/open-telemetry/opentelemetry-go/issues/514
		// Note that the default SDK will not generate any errors yet, this is
		// only for added safety.
		panic(err)
	}

	atomic.StorePointer(&obs.delegate, unsafe.Pointer(obsPtr))
}

func (obs float64ObsImpl) Unregister() {
	obs.observer.unregister()
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
	// This may still be nil if instrument was created and bound
	// without a delegate, then the instrument was set to have a
	// delegate and unbound.
	if implPtr == nil {
		return
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

func (m *meter) NewInt64Counter(name string, opts ...metric.Option) (metric.Int64Counter, error) {
	return metric.WrapInt64CounterInstrument(m.newInst(name, counterKind, core.Int64NumberKind, opts))
}

func (m *meter) NewFloat64Counter(name string, opts ...metric.Option) (metric.Float64Counter, error) {
	return metric.WrapFloat64CounterInstrument(m.newInst(name, counterKind, core.Float64NumberKind, opts))
}

func (m *meter) NewInt64Measure(name string, opts ...metric.Option) (metric.Int64Measure, error) {
	return metric.WrapInt64MeasureInstrument(m.newInst(name, measureKind, core.Int64NumberKind, opts))
}

func (m *meter) NewFloat64Measure(name string, opts ...metric.Option) (metric.Float64Measure, error) {
	return metric.WrapFloat64MeasureInstrument(m.newInst(name, measureKind, core.Float64NumberKind, opts))
}

func (m *meter) RegisterInt64Observer(name string, callback metric.Int64ObserverCallback, opts ...metric.Option) (metric.Int64Observer, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if meterPtr := (*metric.Meter)(atomic.LoadPointer(&m.delegate)); meterPtr != nil {
		return (*meterPtr).RegisterInt64Observer(name, callback, opts...)
	}

	obs := &obsImpl{
		name:     name,
		nkind:    core.Int64NumberKind,
		opts:     opts,
		meter:    m,
		callback: callback,
	}
	m.addObserver(obs)
	return int64ObsImpl{
		observer: obs,
	}, nil
}

func (m *meter) RegisterFloat64Observer(name string, callback metric.Float64ObserverCallback, opts ...metric.Option) (metric.Float64Observer, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if meterPtr := (*metric.Meter)(atomic.LoadPointer(&m.delegate)); meterPtr != nil {
		return (*meterPtr).RegisterFloat64Observer(name, callback, opts...)
	}

	obs := &obsImpl{
		name:     name,
		nkind:    core.Float64NumberKind,
		opts:     opts,
		meter:    m,
		callback: callback,
	}
	m.addObserver(obs)
	return float64ObsImpl{
		observer: obs,
	}, nil
}

func (m *meter) addObserver(obs *obsImpl) {
	if m.liveObservers == nil {
		m.liveObservers = make(map[*obsImpl]struct{})
	}
	m.liveObservers[obs] = struct{}{}
	m.orderedObservers = append(m.orderedObservers, obs)
}

func AtomicFieldOffsets() map[string]uintptr {
	return map[string]uintptr{
		"meterProvider.delegate": unsafe.Offsetof(meterProvider{}.delegate),
		"meter.delegate":         unsafe.Offsetof(meter{}.delegate),
		"instImpl.delegate":      unsafe.Offsetof(instImpl{}.delegate),
		"obsImpl.delegate":       unsafe.Offsetof(obsImpl{}.delegate),
		"labelSet.delegate":      unsafe.Offsetof(labelSet{}.delegate),
		"instHandle.delegate":    unsafe.Offsetof(instHandle{}.delegate),
	}
}
