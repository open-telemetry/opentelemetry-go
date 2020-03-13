// Copyright 2019, OpenTelemetry Authors
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
	apimetric "go.opentelemetry.io/otel/api/metric"
)

type (
	Handle struct {
		Instrument *Synchronous
		LabelSet   *LabelSet
	}

	LabelSet struct {
		TheMeter *Meter
		Labels   map[core.Key]core.Value
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		LabelSet     *LabelSet
	}

	MeterProvider struct {
		lock       sync.Mutex
		registered map[string]apimetric.Meter
	}

	Meter struct {
		// TODO add synchronization

		MeasurementBatches []Batch
		// Observers contains also unregistered
		// observers. Check the Dead field of the Observer to
		// figure out its status.
		Observers []*Asynchronous
	}

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     core.Number
		Instrument apimetric.InstrumentImpl
	}

	Instrument struct {
		meter      *Meter
		descriptor apimetric.Descriptor
	}

	Asynchronous struct {
		Instrument

		Dead bool

		callback func(func(core.Number, apimetric.LabelSet))
	}

	Synchronous struct {
		Instrument
	}
)

var (
	_ apimetric.SynchronousImpl      = &Synchronous{}
	_ apimetric.BoundSynchronousImpl = &Handle{}
	_ apimetric.LabelSet             = &LabelSet{}
	_ apimetric.MeterImpl            = &Meter{}
	_ apimetric.AsynchronousImpl     = &Asynchronous{}
)

func (i Instrument) Descriptor() apimetric.Descriptor {
	return i.descriptor
}

func (a *Asynchronous) Interface() interface{} {
	return a
}

func (s *Synchronous) Interface() interface{} {
	return s
}

func (a *Asynchronous) Unregister() {
	a.Dead = true
}

func (s *Synchronous) Bind(labels apimetric.LabelSet) apimetric.BoundSynchronousImpl {
	if ld, ok := labels.(apimetric.LabelSetDelegate); ok {
		labels = ld.Delegate()
	}
	return &Handle{
		Instrument: s,
		LabelSet:   labels.(*LabelSet),
	}
}

func (i *Synchronous) RecordOne(ctx context.Context, number core.Number, labels apimetric.LabelSet) {
	if ld, ok := labels.(apimetric.LabelSetDelegate); ok {
		labels = ld.Delegate()
	}
	i.meter.doRecordSingle(ctx, labels.(*LabelSet), i, number)
}

func (h *Handle) RecordOne(ctx context.Context, number core.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.LabelSet, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *Meter) doRecordSingle(ctx context.Context, labelSet *LabelSet, instrument apimetric.InstrumentImpl, number core.Number) {
	m.recordMockBatch(ctx, labelSet, Measurement{
		Instrument: instrument,
		Number:     number,
	})
}

func NewProvider() *MeterProvider {
	return &MeterProvider{
		registered: map[string]apimetric.Meter{},
	}
}

func (p *MeterProvider) Meter(name string) apimetric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if lookup, ok := p.registered[name]; ok {
		return lookup
	}
	_, m := NewMeter()
	p.registered[name] = m
	return m
}

func NewMeter() (*Meter, apimetric.Meter) {
	mock := &Meter{}
	return mock, apimetric.WrapMeterImpl(mock)
}

func (m *Meter) Labels(labels ...core.KeyValue) apimetric.LabelSet {
	ul := make(map[core.Key]core.Value)
	for _, kv := range labels {
		ul[kv.Key] = kv.Value
	}
	return &LabelSet{
		TheMeter: m,
		Labels:   ul,
	}
}

func (m *Meter) NewSynchronousInstrument(name string, metricKind apimetric.Kind, numberKind core.NumberKind, config apimetric.Config) (apimetric.SynchronousImpl, error) {
	return &Synchronous{
		Instrument{
			descriptor: apimetric.Descriptor{
				Name:       name,
				Kind:       metricKind,
				NumberKind: numberKind,
				Config:     config,
			},
			meter: m,
		},
	}, nil
}

func (m *Meter) NewAsynchronousInstrument(name string, metricKind apimetric.Kind, numberKind core.NumberKind, callback func(func(core.Number, apimetric.LabelSet)), config apimetric.Config) (apimetric.AsynchronousImpl, error) {
	a := &Asynchronous{
		Instrument: Instrument{
			descriptor: apimetric.Descriptor{
				Name:       name,
				Kind:       metricKind,
				NumberKind: numberKind,
				Config:     config,
			},
			meter: m,
		},
		callback: callback,
	}
	m.Observers = append(m.Observers, a)
	return a, nil
}

func (m *Meter) RecordBatch(ctx context.Context, labels apimetric.LabelSet, measurements ...apimetric.Measurement) {
	ourLabelSet := labels.(*LabelSet)
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.InstrumentImpl().Interface().(*Synchronous),
			Number:     m.Number(),
		}
	}
	m.recordMockBatch(ctx, ourLabelSet, mm...)
}

func (m *Meter) recordMockBatch(ctx context.Context, labelSet *LabelSet, measurements ...Measurement) {
	m.MeasurementBatches = append(m.MeasurementBatches, Batch{
		Ctx:          ctx,
		LabelSet:     labelSet,
		Measurements: measurements,
	})
}

func (m *Meter) RunObservers() {
	for _, observer := range m.Observers {
		if observer.Dead {
			continue
		}
		observer.callback(func(n core.Number, ls apimetric.LabelSet) {
			m.doRecordSingle(context.Background(), ls.(*LabelSet), observer, n)
		})
	}
}
