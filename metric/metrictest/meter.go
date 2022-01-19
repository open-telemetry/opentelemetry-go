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

package metrictest // import "go.opentelemetry.io/otel/metric/metrictest"

import (
	"context"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	internalmetric "go.opentelemetry.io/otel/internal/metric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

type (
	Handle struct {
		Instrument *Sync
		Labels     []attribute.KeyValue
	}

	// Library is the same as "sdk/instrumentation".Library but there is
	// a package cycle to use it.
	Library struct {
		InstrumentationName    string
		InstrumentationVersion string
		SchemaURL              string
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Ctx          context.Context
		Labels       []attribute.KeyValue
		Library      Library
	}

	// MeterImpl is an OpenTelemetry Meter implementation used for testing.
	MeterImpl struct {
		library          Library
		provider         *MeterProvider
		asyncInstruments *internalmetric.AsyncInstrumentState
	}

	// MeterProvider is a collection of named MeterImpls used for testing.
	MeterProvider struct {
		lock sync.Mutex

		MeasurementBatches []Batch
		impls              []*MeterImpl
	}

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     number.Number
		Instrument sdkapi.InstrumentImpl
	}

	Instrument struct {
		meter      *MeterImpl
		descriptor sdkapi.Descriptor
	}

	Async struct {
		Instrument

		runner sdkapi.AsyncRunner
	}

	Sync struct {
		Instrument
	}
)

var (
	_ sdkapi.SyncImpl  = &Sync{}
	_ sdkapi.MeterImpl = &MeterImpl{}
	_ sdkapi.AsyncImpl = &Async{}
)

// NewDescriptor is a test helper for constructing test metric
// descriptors using standard options.
func NewDescriptor(name string, ikind sdkapi.InstrumentKind, nkind number.Kind, opts ...metric.InstrumentOption) sdkapi.Descriptor {
	cfg := metric.NewInstrumentConfig(opts...)
	return sdkapi.NewDescriptor(name, ikind, nkind, cfg.Description(), cfg.Unit())
}

func (i Instrument) Descriptor() sdkapi.Descriptor {
	return i.descriptor
}

func (a *Async) Implementation() interface{} {
	return a
}

func (s *Sync) Implementation() interface{} {
	return s
}

func (s *Sync) RecordOne(ctx context.Context, number number.Number, labels []attribute.KeyValue) {
	s.meter.doRecordSingle(ctx, labels, s, number)
}

func (h *Handle) RecordOne(ctx context.Context, number number.Number) {
	h.Instrument.meter.doRecordSingle(ctx, h.Labels, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func (m *MeterImpl) doRecordSingle(ctx context.Context, labels []attribute.KeyValue, instrument sdkapi.InstrumentImpl, number number.Number) {
	m.collect(ctx, labels, []Measurement{{
		Instrument: instrument,
		Number:     number,
	}})
}

// NewMeterProvider returns a MeterProvider suitable for testing.
// When the test is complete, consult MeterProvider.MeasurementBatches.
func NewMeterProvider() *MeterProvider {
	return &MeterProvider{}
}

// Meter implements metric.MeterProvider.
func (p *MeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	p.lock.Lock()
	defer p.lock.Unlock()
	cfg := metric.NewMeterConfig(opts...)
	impl := &MeterImpl{
		library: Library{
			InstrumentationName:    name,
			InstrumentationVersion: cfg.InstrumentationVersion(),
			SchemaURL:              cfg.SchemaURL(),
		},
		provider:         p,
		asyncInstruments: internalmetric.NewAsyncInstrumentState(),
	}
	p.impls = append(p.impls, impl)
	return metric.WrapMeterImpl(impl)
}

// NewSyncInstrument implements sdkapi.MeterImpl.
func (m *MeterImpl) NewSyncInstrument(descriptor sdkapi.Descriptor) (sdkapi.SyncImpl, error) {
	return &Sync{
		Instrument{
			descriptor: descriptor,
			meter:      m,
		},
	}, nil
}

// NewAsyncInstrument implements sdkapi.MeterImpl.
func (m *MeterImpl) NewAsyncInstrument(descriptor sdkapi.Descriptor, runner sdkapi.AsyncRunner) (sdkapi.AsyncImpl, error) {
	a := &Async{
		Instrument: Instrument{
			descriptor: descriptor,
			meter:      m,
		},
		runner: runner,
	}
	m.provider.registerAsyncInstrument(a, m, runner)
	return a, nil
}

// RecordBatch implements sdkapi.MeterImpl.
func (m *MeterImpl) RecordBatch(ctx context.Context, labels []attribute.KeyValue, measurements ...sdkapi.Measurement) {
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

// CollectAsync is called from asyncInstruments.Run() with the lock held.
func (m *MeterImpl) CollectAsync(labels []attribute.KeyValue, obs ...sdkapi.Observation) {
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

// collect is called from CollectAsync() or RecordBatch() with the lock held.
func (m *MeterImpl) collect(ctx context.Context, labels []attribute.KeyValue, measurements []Measurement) {
	m.provider.addMeasurement(Batch{
		Ctx:          ctx,
		Labels:       labels,
		Measurements: measurements,
		Library:      m.library,
	})
}

// registerAsyncInstrument locks the provider and registers the new Async instrument.
func (p *MeterProvider) registerAsyncInstrument(a *Async, m *MeterImpl, runner sdkapi.AsyncRunner) {
	p.lock.Lock()
	defer p.lock.Unlock()

	m.asyncInstruments.Register(a, runner)
}

// addMeasurement locks the provider and adds the new measurement batch.
func (p *MeterProvider) addMeasurement(b Batch) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.MeasurementBatches = append(p.MeasurementBatches, b)
}

// copyImpls locks the provider and copies the current list of *MeterImpls.
func (p *MeterProvider) copyImpls() []*MeterImpl {
	p.lock.Lock()
	defer p.lock.Unlock()
	cpy := make([]*MeterImpl, len(p.impls))
	copy(cpy, p.impls)
	return cpy
}

// RunAsyncInstruments is used in tests to trigger collection from
// asynchronous instruments.
func (p *MeterProvider) RunAsyncInstruments() {
	for _, impl := range p.copyImpls() {
		impl.asyncInstruments.Run(context.Background(), impl)
	}
}

// Measured is the helper struct which provides flat representation of recorded measurements
// to simplify testing
type Measured struct {
	Name    string
	Labels  map[attribute.Key]attribute.Value
	Number  number.Number
	Library Library
}

// LabelsToMap converts label set to keyValue map, to be easily used in tests
func LabelsToMap(kvs ...attribute.KeyValue) map[attribute.Key]attribute.Value {
	m := map[attribute.Key]attribute.Value{}
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
				Name:    m.Instrument.Descriptor().Name(),
				Labels:  LabelsToMap(batch.Labels...),
				Number:  m.Number,
				Library: batch.Library,
			})
		}
	}
	return r
}

// ResolveNumberByKind takes defined metric descriptor creates a concrete typed metric number
func ResolveNumberByKind(t *testing.T, kind number.Kind, value float64) number.Number {
	t.Helper()
	switch kind {
	case number.Int64Kind:
		return number.NewInt64Number(int64(value))
	case number.Float64Kind:
		return number.NewFloat64Number(value)
	}
	panic("invalid number kind")
}
