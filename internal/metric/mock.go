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
		LabelSet   *LabelSet
	}

	LabelSet struct {
		Impl   *MeterImpl
		Labels map[core.Key]core.Value
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		LabelSet     *LabelSet
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
		Number     core.Number
		Instrument apimetric.InstrumentImpl
	}

	Instrument struct {
		meter      *MeterImpl
		descriptor apimetric.Descriptor
	}

	Async struct {
		Instrument

		callback func(func(core.Number, apimetric.LabelSet))
	}

	Sync struct {
		Instrument
	}
)

var (
	_ apimetric.SyncImpl      = &Sync{}
	_ apimetric.BoundSyncImpl = &Handle{}
	_ apimetric.LabelSet      = &LabelSet{}
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

func (s *Sync) Bind(labels apimetric.LabelSet) apimetric.BoundSyncImpl {
	if ld, ok := labels.(apimetric.LabelSetDelegate); ok {
		labels = ld.Delegate()
	}
	return &Handle{
		Instrument: s,
		LabelSet:   labels.(*LabelSet),
	}
}

func (s *Sync) RecordOne(ctx context.Context, number core.Number, labels apimetric.LabelSet) {
	if ld, ok := labels.(apimetric.LabelSetDelegate); ok {
		labels = ld.Delegate()
	}
	s.meter.doRecordSingle(ctx, labels.(*LabelSet), s, number)
}

func (h *Handle) RecordOne(ctx context.Context, number core.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.LabelSet, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *MeterImpl) doRecordSingle(ctx context.Context, labelSet *LabelSet, instrument apimetric.InstrumentImpl, number core.Number) {
	m.recordMockBatch(ctx, labelSet, Measurement{
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

func (m *MeterImpl) Labels(labels ...core.KeyValue) apimetric.LabelSet {
	ul := make(map[core.Key]core.Value)
	for _, kv := range labels {
		ul[kv.Key] = kv.Value
	}
	return &LabelSet{
		Impl:   m,
		Labels: ul,
	}
}

func (m *MeterImpl) NewSyncInstrument(descriptor metric.Descriptor) (apimetric.SyncImpl, error) {
	return &Sync{
		Instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *MeterImpl) NewAsyncInstrument(descriptor metric.Descriptor, callback func(func(core.Number, apimetric.LabelSet))) (apimetric.AsyncImpl, error) {
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

func (m *MeterImpl) RecordBatch(ctx context.Context, labels apimetric.LabelSet, measurements ...apimetric.Measurement) {
	ourLabelSet := labels.(*LabelSet)
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.SyncImpl().(*Sync),
			Number:     m.Number(),
		}
	}
	m.recordMockBatch(ctx, ourLabelSet, mm...)
}

func (m *MeterImpl) recordMockBatch(ctx context.Context, labelSet *LabelSet, measurements ...Measurement) {
	m.MeasurementBatches = append(m.MeasurementBatches, Batch{
		Ctx:          ctx,
		LabelSet:     labelSet,
		Measurements: measurements,
	})
}

func (m *MeterImpl) RunAsyncInstruments() {
	for _, observer := range m.AsyncInstruments {
		observer.callback(func(n core.Number, labels apimetric.LabelSet) {

			if ld, ok := labels.(apimetric.LabelSetDelegate); ok {
				labels = ld.Delegate()
			}

			m.doRecordSingle(context.Background(), labels.(*LabelSet), observer, n)
		})
	}
}
