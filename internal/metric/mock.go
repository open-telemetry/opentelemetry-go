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

package metric

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	apimetric "go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/metric/registry"
)

type (
	Handle struct {
		Instrument *Sync
		Labels     []core.KeyValue
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		Labels       []core.KeyValue
		LibraryName  string
	}

	MeterProvider struct {
		lock       sync.Mutex
		impl       *MeterImpl
		unique     metric.MeterImpl
		registered map[string]apimetric.Meter
	}

	MeterImpl struct {
		MeasurementBatches []Batch
		AsyncInstruments   []*Async
	}

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     apimetric.Number
		Instrument apimetric.InstrumentImpl
	}

	Instrument struct {
		meter      *MeterImpl
		descriptor apimetric.Descriptor
	}

	Async struct {
		Instrument

		callback func(func(apimetric.Number, []core.KeyValue))
	}

	Sync struct {
		Instrument
	}
)

var (
	_ apimetric.SyncImpl      = &Sync{}
	_ apimetric.BoundSyncImpl = &Handle{}
	_ apimetric.MeterImpl     = &MeterImpl{}
	_ apimetric.AsyncImpl     = &Async{}
)

func (i Instrument) Descriptor() apimetric.Descriptor {
	return i.descriptor
}

func (a *Async) Implementation() interface{} {
	return a
}

func (s *Sync) Implementation() interface{} {
	return s
}

func (s *Sync) Bind(labels []core.KeyValue) apimetric.BoundSyncImpl {
	return &Handle{
		Instrument: s,
		Labels:     labels,
	}
}

func (s *Sync) RecordOne(ctx context.Context, number apimetric.Number, labels []core.KeyValue) {
	s.meter.doRecordSingle(ctx, labels, s, number)
}

func (h *Handle) RecordOne(ctx context.Context, number apimetric.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.Labels, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *MeterImpl) doRecordSingle(ctx context.Context, labels []core.KeyValue, instrument apimetric.InstrumentImpl, number apimetric.Number) {
	m.recordMockBatch(ctx, labels, Measurement{
		Instrument: instrument,
		Number:     number,
	})
}

func NewProvider() (*MeterImpl, apimetric.Provider) {
	impl := &MeterImpl{}
	p := &MeterProvider{
		impl:       impl,
		unique:     registry.NewUniqueInstrumentMeterImpl(impl),
		registered: map[string]apimetric.Meter{},
	}
	return impl, p
}

func (p *MeterProvider) Meter(name string) apimetric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if lookup, ok := p.registered[name]; ok {
		return lookup
	}
	m := apimetric.WrapMeterImpl(p.unique, name)
	p.registered[name] = m
	return m
}

func NewMeter() (*MeterImpl, apimetric.Meter) {
	impl, p := NewProvider()
	return impl, p.Meter("mock")
}

func (m *MeterImpl) NewSyncInstrument(descriptor metric.Descriptor) (apimetric.SyncImpl, error) {
	return &Sync{
		Instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *MeterImpl) NewAsyncInstrument(descriptor metric.Descriptor, callback func(func(apimetric.Number, []core.KeyValue))) (apimetric.AsyncImpl, error) {
	a := &Async{
		Instrument: Instrument{
			descriptor: descriptor,
			meter:      m,
		},
		callback: callback,
	}
	m.AsyncInstruments = append(m.AsyncInstruments, a)
	return a, nil
}

func (m *MeterImpl) RecordBatch(ctx context.Context, labels []core.KeyValue, measurements ...apimetric.Measurement) {
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.SyncImpl().Implementation().(*Sync),
			Number:     m.Number(),
		}
	}
	m.recordMockBatch(ctx, labels, mm...)
}

func (m *MeterImpl) recordMockBatch(ctx context.Context, labels []core.KeyValue, measurements ...Measurement) {
	m.MeasurementBatches = append(m.MeasurementBatches, Batch{
		Ctx:          ctx,
		Labels:       labels,
		Measurements: measurements,
	})
}

func (m *MeterImpl) RunAsyncInstruments() {
	for _, observer := range m.AsyncInstruments {
		observer.callback(func(n apimetric.Number, labels []core.KeyValue) {
			m.doRecordSingle(context.Background(), labels, observer, n)
		})
	}
}
