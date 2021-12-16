// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package global // import "go.opentelemetry.io/otel/internal/metric/global"

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/metric/registry"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

// This file contains the forwarding implementation of MeterProvider used as
// the default global instance.  Metric events using instruments provided by
// this implementation are no-ops until the first Meter implementation is set
// as the global provider.
//
// The implementation here uses Mutexes to maintain a list of active Meters in
// the MeterProvider and Instruments in each Meter, under the assumption that
// these interfaces are not performance-critical.
//
// We have the invariant that setDelegate() will be called before a new
// MeterProvider implementation is registered as the global provider.  Mutexes
// in the MeterProvider and Meters ensure that each instrument has a delegate
// before the global provider is set.
//
// Metric uniqueness checking is implemented by calling the exported
// methods of the api/metric/registry package.

type meterKey struct {
	InstrumentationName    string
	InstrumentationVersion string
	SchemaURL              string
}

type meterProvider struct {
	delegate metric.MeterProvider

	// lock protects `delegate` and `meters`.
	lock sync.Mutex

	// meters maintains a unique entry for every named Meter
	// that has been registered through the global instance.
	meters map[meterKey]*meterEntry
}

type meterImpl struct {
	delegate unsafe.Pointer // (*metric.MeterImpl)

	lock      sync.Mutex
	insts     []*instrument
	callbacks []*callback
}

type meterEntry struct {
	unique sdkapi.MeterImpl
	impl   meterImpl
}

type instrument struct {
	delegate unsafe.Pointer // (*sdkapi.Instrument)

	descriptor sdkapi.Descriptor
}

type callback struct {
	delegate unsafe.Pointer // (*sdkapi.Callback)

	insts    []sdkapi.Instrument
	function func(context.Context) error
}

var _ metric.MeterProvider = &meterProvider{}
var _ sdkapi.MeterImpl = &meterImpl{}
var _ sdkapi.Instrument = &instrument{}
var _ sdkapi.Callback = &callback{}

func (inst *instrument) Descriptor() sdkapi.Descriptor {
	return inst.descriptor
}

// MeterProvider interface and delegation

func newMeterProvider() *meterProvider {
	return &meterProvider{
		meters: map[meterKey]*meterEntry{},
	}
}

func (p *meterProvider) setDelegate(provider metric.MeterProvider) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.delegate = provider
	for key, entry := range p.meters {
		entry.impl.setDelegate(key, provider)
	}
	p.meters = nil
}

func (p *meterProvider) Meter(instrumentationName string, opts ...metric.MeterOption) metric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.delegate != nil {
		return p.delegate.Meter(instrumentationName, opts...)
	}

	cfg := metric.NewMeterConfig(opts...)
	key := meterKey{
		InstrumentationName:    instrumentationName,
		InstrumentationVersion: cfg.InstrumentationVersion(),
		SchemaURL:              cfg.SchemaURL(),
	}
	entry, ok := p.meters[key]
	if !ok {
		entry = &meterEntry{}
		// Note: This code implements its own MeterProvider
		// name-uniqueness logic because there is
		// synchronization required at the moment of
		// delegation.  We use the same instrument-uniqueness
		// checking the real SDK uses here:
		entry.unique = registry.NewUniqueInstrumentMeterImpl(&entry.impl)
		p.meters[key] = entry
	}
	return metric.Meter{entry.unique}
}

// Meter interface and delegation

func (m *meterImpl) setDelegate(key meterKey, provider metric.MeterProvider) {
	m.lock.Lock()
	defer m.lock.Unlock()

	d := new(sdkapi.MeterImpl)
	*d = provider.Meter(
		key.InstrumentationName,
		metric.WithInstrumentationVersion(key.InstrumentationVersion),
		metric.WithSchemaURL(key.SchemaURL),
	).MeterImpl
	m.delegate = unsafe.Pointer(d)

	for _, inst := range m.insts {
		inst.setDelegate(*d)
	}
	m.insts = nil
	for _, cb := range m.callbacks {
		cb.setDelegate(*d)
	}
	m.callbacks = nil
}

func (m *meterImpl) NewInstrument(desc sdkapi.Descriptor) (sdkapi.Instrument, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if meterPtr := (*sdkapi.MeterImpl)(atomic.LoadPointer(&m.delegate)); meterPtr != nil {
		return (*meterPtr).NewInstrument(desc)
	}

	inst := &instrument{
		descriptor: desc,
	}
	m.insts = append(m.insts, inst)
	return inst, nil
}

// Synchronous delegation

func (inst *instrument) setDelegate(d sdkapi.MeterImpl) {
	implPtr := new(sdkapi.Instrument)

	var err error
	*implPtr, err = d.NewInstrument(inst.descriptor)

	if err != nil {
		// TODO: There is no standard way to deliver this error to the user.
		// See https://github.com/open-telemetry/opentelemetry-go/issues/514
		// Note that the default SDK will not generate any errors yet, this is
		// only for added safety.
		panic(err)
	}

	atomic.StorePointer(&inst.delegate, unsafe.Pointer(implPtr))
}

func (inst *instrument) Implementation() interface{} {
	if implPtr := (*sdkapi.Instrument)(atomic.LoadPointer(&inst.delegate)); implPtr != nil {
		return (*implPtr).Implementation()
	}
	return inst
}

// Async delegation

func (m *meterImpl) NewCallback(insts []sdkapi.Instrument, f func(context.Context) error) (sdkapi.Callback, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if meterPtr := (*sdkapi.MeterImpl)(atomic.LoadPointer(&m.delegate)); meterPtr != nil {
		return (*meterPtr).NewCallback(insts, f)
	}

	cb := &callback{
		insts:    insts,
		function: f,
	}
	m.callbacks = append(m.callbacks, cb)
	return cb, nil
}

func (cb *callback) Instruments() []sdkapi.Instrument {
	if implPtr := (*sdkapi.Callback)(atomic.LoadPointer(&cb.delegate)); implPtr != nil {
		return (*implPtr).Instruments()
	}
	return cb.insts
}

func (cb *callback) setDelegate(d sdkapi.MeterImpl) {
	implPtr := new(sdkapi.Callback)

	var err error
	*implPtr, err = d.NewCallback(cb.insts, cb.function)

	if err != nil {
		panic(err)
	}

	atomic.StorePointer(&cb.delegate, unsafe.Pointer(implPtr))
}

// Metric updates

func (inst *instrument) RecordOne(ctx context.Context, number number.Number, attrs attribute.Attributes) {
	if instPtr := (*sdkapi.Instrument)(atomic.LoadPointer(&inst.delegate)); instPtr != nil {
		(*instPtr).RecordOne(ctx, number, attrs)
	}
}

func AtomicFieldOffsets() map[string]uintptr {
	return map[string]uintptr{
		"meterProvider.delegate": unsafe.Offsetof(meterProvider{}.delegate),
		"meterImpl.delegate":     unsafe.Offsetof(meterImpl{}.delegate),
		"instrument.delegate":    unsafe.Offsetof(instrument{}.delegate),
		"callback.delegate":      unsafe.Offsetof(instrument{}.delegate),
	}
}
