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

	"go.opentelemetry.io/otel"
)

type (
	Handle struct {
		Instrument *Instrument
		LabelSet   *LabelSet
	}

	Instrument struct {
		Name       string
		Kind       Kind
		NumberKind otel.NumberKind
		Opts       otel.MetricOptions
	}

	LabelSet struct {
		TheMeter *Meter
		Labels   map[otel.Key]otel.Value
	}

	Batch struct {
		Ctx          context.Context
		LabelSet     *LabelSet
		Measurements []Measurement
	}

	Meter struct {
		MeasurementBatches []Batch
	}

	Kind int8

	Measurement struct {
		Instrument *Instrument
		Number     otel.Number
	}
)

var (
	_ otel.Instrument = &Instrument{}
	_ otel.Handle     = &Handle{}
	_ otel.LabelSet   = &LabelSet{}
	_ otel.Meter      = &Meter{}
)

const (
	KindCounter Kind = iota
	KindGauge
	KindMeasure
)

func (i *Instrument) AcquireHandle(labels otel.LabelSet) otel.Handle {
	return &Handle{
		Instrument: i,
		LabelSet:   labels.(*LabelSet),
	}
}

func (i *Instrument) RecordOne(ctx context.Context, number otel.Number, labels otel.LabelSet) {
	doRecordBatch(ctx, labels.(*LabelSet), i, number)
}

func (h *Handle) RecordOne(ctx context.Context, number otel.Number) {
	doRecordBatch(ctx, h.LabelSet, h.Instrument, number)
}

func (h *Handle) Release() {
}

func doRecordBatch(ctx context.Context, labelSet *LabelSet, instrument *Instrument, number otel.Number) {
	labelSet.TheMeter.recordMockBatch(ctx, labelSet, Measurement{
		Instrument: instrument,
		Number:     number,
	})
}

func (s *LabelSet) Meter() otel.Meter {
	return s.TheMeter
}

func NewMeter() *Meter {
	return &Meter{}
}

func (m *Meter) Labels(labels ...otel.KeyValue) otel.LabelSet {
	ul := make(map[otel.Key]otel.Value)
	for _, kv := range labels {
		ul[kv.Key] = kv.Value
	}
	return &LabelSet{
		TheMeter: m,
		Labels:   ul,
	}
}

func (m *Meter) NewInt64Counter(name string, cos ...otel.CounterOptionApplier) otel.Int64Counter {
	instrument := m.newCounterInstrument(name, otel.Int64NumberKind, cos...)
	return otel.WrapInt64CounterInstrument(instrument)
}

func (m *Meter) NewFloat64Counter(name string, cos ...otel.CounterOptionApplier) otel.Float64Counter {
	instrument := m.newCounterInstrument(name, otel.Float64NumberKind, cos...)
	return otel.WrapFloat64CounterInstrument(instrument)
}

func (m *Meter) newCounterInstrument(name string, numberKind otel.NumberKind, cos ...otel.CounterOptionApplier) *Instrument {
	opts := otel.MetricOptions{}
	otel.ApplyCounterOptions(&opts, cos...)
	return &Instrument{
		Name:       name,
		Kind:       KindCounter,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) NewInt64Gauge(name string, gos ...otel.GaugeOptionApplier) otel.Int64Gauge {
	instrument := m.newGaugeInstrument(name, otel.Int64NumberKind, gos...)
	return otel.WrapInt64GaugeInstrument(instrument)
}

func (m *Meter) NewFloat64Gauge(name string, gos ...otel.GaugeOptionApplier) otel.Float64Gauge {
	instrument := m.newGaugeInstrument(name, otel.Float64NumberKind, gos...)
	return otel.WrapFloat64GaugeInstrument(instrument)
}

func (m *Meter) newGaugeInstrument(name string, numberKind otel.NumberKind, gos ...otel.GaugeOptionApplier) *Instrument {
	opts := otel.MetricOptions{}
	otel.ApplyGaugeOptions(&opts, gos...)
	return &Instrument{
		Name:       name,
		Kind:       KindGauge,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) NewInt64Measure(name string, mos ...otel.MeasureOptionApplier) otel.Int64Measure {
	instrument := m.newMeasureInstrument(name, otel.Int64NumberKind, mos...)
	return otel.WrapInt64MeasureInstrument(instrument)
}

func (m *Meter) NewFloat64Measure(name string, mos ...otel.MeasureOptionApplier) otel.Float64Measure {
	instrument := m.newMeasureInstrument(name, otel.Float64NumberKind, mos...)
	return otel.WrapFloat64MeasureInstrument(instrument)
}

func (m *Meter) newMeasureInstrument(name string, numberKind otel.NumberKind, mos ...otel.MeasureOptionApplier) *Instrument {
	opts := otel.MetricOptions{}
	otel.ApplyMeasureOptions(&opts, mos...)
	return &Instrument{
		Name:       name,
		Kind:       KindMeasure,
		NumberKind: numberKind,
		Opts:       opts,
	}
}

func (m *Meter) RecordBatch(ctx context.Context, labels otel.LabelSet, measurements ...otel.Measurement) {
	ourLabelSet := labels.(*LabelSet)
	mm := make([]Measurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = Measurement{
			Instrument: m.Instrument().(*Instrument),
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
