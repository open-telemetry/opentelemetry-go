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

package metrictest

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/registry"
)

type (
	Handle struct {
		Instrument *Sync
		Labels     []label.KeyValue
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		Labels       []label.KeyValue
		LibraryName  string
	}

	// MeterImpl is an OpenTelemetry Meter implementation used for testing.
	MeterImpl struct {
		lock sync.Mutex

		MeasurementBatches []Batch

		asyncInstruments *AsyncInstrumentState
	}

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     otel.Number
		Instrument otel.InstrumentImpl
	}

	Instrument struct {
		meter      *MeterImpl
		descriptor otel.Descriptor
	}

	Async struct {
		Instrument

		runner otel.AsyncRunner
	}

	Sync struct {
		Instrument
	}
)

var (
	_ otel.SyncImpl      = &Sync{}
	_ otel.BoundSyncImpl = &Handle{}
	_ otel.MeterImpl     = &MeterImpl{}
	_ otel.AsyncImpl     = &Async{}
)

func (i Instrument) Descriptor() otel.Descriptor {
	return i.descriptor
}

func (a *Async) Implementation() interface{} {
	return a
}

func (s *Sync) Implementation() interface{} {
	return s
}

func (s *Sync) Bind(labels []label.KeyValue) otel.BoundSyncImpl {
	return &Handle{
		Instrument: s,
		Labels:     labels,
	}
}

func (s *Sync) RecordOne(ctx context.Context, number otel.Number, labels []label.KeyValue) {
	s.meter.doRecordSingle(ctx, labels, s, number)
}

func (h *Handle) RecordOne(ctx context.Context, number otel.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.Labels, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *MeterImpl) doRecordSingle(ctx context.Context, labels []label.KeyValue, instrument otel.InstrumentImpl, number otel.Number) {
	m.collect(ctx, labels, []Measurement{{
		Instrument: instrument,
		Number:     number,
	}})
}

func NewProvider() (*MeterImpl, otel.MeterProvider) {
	impl := &MeterImpl{
		asyncInstruments: NewAsyncInstrumentState(),
	}
	return impl, registry.NewProvider(impl)
}

func NewMeter() (*MeterImpl, otel.Meter) {
	impl, p := NewProvider()
	return impl, p.Meter("mock")
}

func (m *MeterImpl) NewSyncInstrument(descriptor otel.Descriptor) (otel.SyncImpl, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return &Sync{
		Instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *MeterImpl) NewAsyncInstrument(descriptor otel.Descriptor, runner otel.AsyncRunner) (otel.AsyncImpl, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	a := &Async{
		Instrument: Instrument{
			descriptor: descriptor,
			meter:      m,
		},
		runner: runner,
	}
	m.asyncInstruments.Register(a, runner)
	return a, nil
}

func (m *MeterImpl) RecordBatch(ctx context.Context, labels []label.KeyValue, measurements ...otel.Measurement) {
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.SyncImpl().Implementation().(*Sync),
			Number:     m.Number(),
		}
	}
	m.collect(ctx, labels, mm)
}

func (m *MeterImpl) CollectAsync(labels []label.KeyValue, obs ...otel.Observation) {
	mm := make([]Measurement, len(obs))
	for i := 0; i < len(obs); i++ {
		o := obs[i]
		mm[i] = Measurement{
			Instrument: o.AsyncImpl(),
			Number:     o.Number(),
		}
	}
	m.collect(context.Background(), labels, mm)
}

func (m *MeterImpl) collect(ctx context.Context, labels []label.KeyValue, measurements []Measurement) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.MeasurementBatches = append(m.MeasurementBatches, Batch{
		Ctx:          ctx,
		Labels:       labels,
		Measurements: measurements,
	})
}

func (m *MeterImpl) RunAsyncInstruments() {
	m.asyncInstruments.Run(context.Background(), m)
}
