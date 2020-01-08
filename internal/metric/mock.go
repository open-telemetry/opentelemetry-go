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

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/core"
	apimetric "go.opentelemetry.io/otel/api/metric"
)

type (
	Handle struct {
		Instrument *Instrument
		Labels     label.Set
	}

	Instrument struct {
		Meter      *Meter
		Name       string
		Kind       Kind
		NumberKind core.NumberKind
		Opts       apimetric.Options
	}

	Batch struct {
		// Measurement needs to be aligned for 64-bit atomic operations.
		Measurements []Measurement
		Labels       label.Set
	}

	Meter struct {
		MeasurementBatches []Batch
	}

	Kind int8

	Measurement struct {
		// Number needs to be aligned for 64-bit atomic operations.
		Number     core.Number
		Instrument *Instrument
	}
)

var (
	_ apimetric.InstrumentImpl      = &Instrument{}
	_ apimetric.BoundInstrumentImpl = &Handle{}
	_ apimetric.Meter               = &Meter{}
)

const (
	KindCounter Kind = iota
	KindGauge
	KindMeasure
)

func (i *Instrument) Bind(ctx context.Context, labels []core.KeyValue) apimetric.BoundInstrumentImpl {
	return &Handle{
		Instrument: i,
		Labels:     scope.Current(ctx).AddResources(labels...).Resources(),
	}
}

func (i *Instrument) RecordOne(ctx context.Context, number core.Number, labels []core.KeyValue) {
	doRecordBatch(ctx, scope.Current(ctx).AddResources(labels...).Resources(), i, number)
}

func (h *Handle) RecordOne(ctx context.Context, number core.Number) {
	doRecordBatch(ctx, h.Labels, h.Instrument, number)
}

func (h *Handle) Unbind() {
}

func doRecordBatch(ctx context.Context, labelSet label.Set, instrument *Instrument, number core.Number) {
	instrument.Meter.recordMockBatch(ctx, labelSet, Measurement{
		Instrument: instrument,
		Number:     number,
	})
}

func NewMeter() *Meter {
	return &Meter{}
}

func (m *Meter) NewInt64Counter(name string, cos ...apimetric.CounterOptionApplier) apimetric.Int64Counter {
	instrument := m.newCounterInstrument(name, core.Int64NumberKind, cos...)
	return apimetric.WrapInt64CounterInstrument(instrument)
}

func (m *Meter) NewFloat64Counter(name string, cos ...apimetric.CounterOptionApplier) apimetric.Float64Counter {
	instrument := m.newCounterInstrument(name, core.Float64NumberKind, cos...)
	return apimetric.WrapFloat64CounterInstrument(instrument)
}

func (m *Meter) newCounterInstrument(name string, numberKind core.NumberKind, cos ...apimetric.CounterOptionApplier) *Instrument {
	opts := apimetric.Options{}
	apimetric.ApplyCounterOptions(&opts, cos...)
	return &Instrument{
		Meter:      m,
		Name:       name,
		Kind:       KindCounter,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) NewInt64Gauge(name string, gos ...apimetric.GaugeOptionApplier) apimetric.Int64Gauge {
	instrument := m.newGaugeInstrument(name, core.Int64NumberKind, gos...)
	return apimetric.WrapInt64GaugeInstrument(instrument)
}

func (m *Meter) NewFloat64Gauge(name string, gos ...apimetric.GaugeOptionApplier) apimetric.Float64Gauge {
	instrument := m.newGaugeInstrument(name, core.Float64NumberKind, gos...)
	return apimetric.WrapFloat64GaugeInstrument(instrument)
}

func (m *Meter) newGaugeInstrument(name string, numberKind core.NumberKind, gos ...apimetric.GaugeOptionApplier) *Instrument {
	opts := apimetric.Options{}
	apimetric.ApplyGaugeOptions(&opts, gos...)
	return &Instrument{
		Meter:      m,
		Name:       name,
		Kind:       KindGauge,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) NewInt64Measure(name string, mos ...apimetric.MeasureOptionApplier) apimetric.Int64Measure {
	instrument := m.newMeasureInstrument(name, core.Int64NumberKind, mos...)
	return apimetric.WrapInt64MeasureInstrument(instrument)
}

func (m *Meter) NewFloat64Measure(name string, mos ...apimetric.MeasureOptionApplier) apimetric.Float64Measure {
	instrument := m.newMeasureInstrument(name, core.Float64NumberKind, mos...)
	return apimetric.WrapFloat64MeasureInstrument(instrument)
}

func (m *Meter) newMeasureInstrument(name string, numberKind core.NumberKind, mos ...apimetric.MeasureOptionApplier) *Instrument {
	opts := apimetric.Options{}
	apimetric.ApplyMeasureOptions(&opts, mos...)
	return &Instrument{
		Meter:      m,
		Name:       name,
		Kind:       KindMeasure,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) RecordBatch(ctx context.Context, labels []core.KeyValue, measurements ...apimetric.Measurement) {
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.InstrumentImpl().(*Instrument),
			Number:     m.Number(),
		}
	}
	m.recordMockBatch(ctx, scope.Labels(ctx, labels...), mm...)
}

func (m *Meter) recordMockBatch(ctx context.Context, labels label.Set, measurements ...Measurement) {
	m.MeasurementBatches = append(m.MeasurementBatches, Batch{
		Labels:       labels,
		Measurements: measurements,
	})
}
