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

package oteltest

import (
	"context"
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	apimetric "go.opentelemetry.io/otel"
	internalmetric "go.opentelemetry.io/otel/internal/metric"
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

		asyncInstruments *internalmetric.AsyncInstrumentState
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

		runner apimetric.AsyncRunner
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

func (s *Sync) Bind(labels []label.KeyValue) apimetric.BoundSyncImpl {
	return &Handle{
		Instrument: s,
		Labels:     labels,
	}
}

func (s *Sync) RecordOne(ctx context.Context, number apimetric.Number, labels []label.KeyValue) {
	s.meter.doRecordSingle(ctx, labels, s, number)
}

func (h *Handle) RecordOne(ctx context.Context, number apimetric.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.Labels, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *MeterImpl) doRecordSingle(ctx context.Context, labels []label.KeyValue, instrument apimetric.InstrumentImpl, number apimetric.Number) {
	m.collect(ctx, labels, []Measurement{{
		Instrument: instrument,
		Number:     number,
	}})
}

func NewMeterProvider() (*MeterImpl, apimetric.MeterProvider) {
	impl := &MeterImpl{
		asyncInstruments: internalmetric.NewAsyncInstrumentState(),
	}
	return impl, registry.NewMeterProvider(impl)
}

func NewMeter() (*MeterImpl, apimetric.Meter) {
	impl, p := NewMeterProvider()
	return impl, p.Meter("mock")
}

func (m *MeterImpl) NewSyncInstrument(descriptor otel.Descriptor) (apimetric.SyncImpl, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	return &Sync{
		Instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

func (m *MeterImpl) NewAsyncInstrument(descriptor otel.Descriptor, runner otel.AsyncRunner) (apimetric.AsyncImpl, error) {
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

func (m *MeterImpl) RecordBatch(ctx context.Context, labels []label.KeyValue, measurements ...apimetric.Measurement) {
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

// Measured is the helper struct which provides flat representation of recorded measurements
// to simplify testing
type Measured struct {
	Name                   string
	InstrumentationName    string
	InstrumentationVersion string
	Labels                 map[label.Key]label.Value
	Number                 otel.Number
}

// LabelsToMap converts label set to keyValue map, to be easily used in tests
func LabelsToMap(kvs ...label.KeyValue) map[label.Key]label.Value {
	m := map[label.Key]label.Value{}
	for _, label := range kvs {
		m[label.Key] = label.Value
	}
	return m
}

// AsStructs converts recorded batches to array of flat, readable Measured helper structures
func AsStructs(batches []Batch) []Measured {
	var r []Measured
	for _, batch := range batches {
		for _, m := range batch.Measurements {
			r = append(r, Measured{
				Name:                   m.Instrument.Descriptor().Name(),
				InstrumentationName:    m.Instrument.Descriptor().InstrumentationName(),
				InstrumentationVersion: m.Instrument.Descriptor().InstrumentationVersion(),
				Labels:                 LabelsToMap(batch.Labels...),
				Number:                 m.Number,
			})
		}
	}
	return r
}

// ResolveNumberByKind takes defined metric descriptor creates a concrete typed metric number
func ResolveNumberByKind(t *testing.T, kind otel.NumberKind, value float64) otel.Number {
	t.Helper()
	switch kind {
	case otel.Int64NumberKind:
		return otel.NewInt64Number(int64(value))
	case otel.Float64NumberKind:
		return otel.NewFloat64Number(value)
	}
	panic("invalid number kind")
}
