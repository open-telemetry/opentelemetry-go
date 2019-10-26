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

	"go.opentelemetry.io/api/core"
)

type (
	mockHandle struct {
		instrument *mockInstrument
		labelSet   *mockLabelSet
	}

	mockInstrument struct {
		name       string
		kind       mockKind
		numberKind core.NumberKind
		opts       Options
	}

	mockLabelSet struct {
		meter  *mockMeter
		labels map[core.Key]core.Value
	}

	batch struct {
		ctx          context.Context
		labelSet     *mockLabelSet
		measurements []mockMeasurement
	}

	mockMeter struct {
		measurementBatches []batch
	}

	mockKind int8

	mockMeasurement struct {
		instrument *mockInstrument
		number     core.Number
	}
)

var (
	_ InstrumentImpl = &mockInstrument{}
	_ HandleImpl     = &mockHandle{}
	_ LabelSet       = &mockLabelSet{}
	_ Meter          = &mockMeter{}
)

const (
	mockKindCounter mockKind = iota
	mockKindGauge
	mockKindMeasure
)

func (i *mockInstrument) AcquireHandle(labels LabelSet) HandleImpl {
	return &mockHandle{
		instrument: i,
		labelSet:   labels.(*mockLabelSet),
	}
}

func (i *mockInstrument) RecordOne(ctx context.Context, number core.Number, labels LabelSet) {
	doRecordBatch(labels.(*mockLabelSet), ctx, i, number)
}

func (h *mockHandle) RecordOne(ctx context.Context, number core.Number) {
	doRecordBatch(h.labelSet, ctx, h.instrument, number)
}

func (h *mockHandle) Release() {
}

func doRecordBatch(labelSet *mockLabelSet, ctx context.Context, instrument *mockInstrument, number core.Number) {
	labelSet.meter.recordMockBatch(ctx, labelSet, mockMeasurement{
		instrument: instrument,
		number:     number,
	})
}

func (s *mockLabelSet) Meter() Meter {
	return s.meter
}

func newMockMeter() *mockMeter {
	return &mockMeter{}
}

func (m *mockMeter) Labels(labels ...core.KeyValue) LabelSet {
	ul := make(map[core.Key]core.Value)
	for _, kv := range labels {
		ul[kv.Key] = kv.Value
	}
	return &mockLabelSet{
		meter:  m,
		labels: ul,
	}
}

func (m *mockMeter) NewInt64Counter(name string, cos ...CounterOptionApplier) Int64Counter {
	instrument := m.newCounterInstrument(name, core.Int64NumberKind, cos...)
	return WrapInt64CounterInstrument(instrument)
}

func (m *mockMeter) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	instrument := m.newCounterInstrument(name, core.Float64NumberKind, cos...)
	return WrapFloat64CounterInstrument(instrument)
}

func (m *mockMeter) newCounterInstrument(name string, numberKind core.NumberKind, cos ...CounterOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyCounterOptions(&opts, cos...)
	return &mockInstrument{
		name:       name,
		kind:       mockKindCounter,
		numberKind: numberKind,
		opts:       opts,
	}
}

func (m *mockMeter) NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	instrument := m.newGaugeInstrument(name, core.Int64NumberKind, gos...)
	return WrapInt64GaugeInstrument(instrument)
}

func (m *mockMeter) NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	instrument := m.newGaugeInstrument(name, core.Float64NumberKind, gos...)
	return WrapFloat64GaugeInstrument(instrument)
}

func (m *mockMeter) newGaugeInstrument(name string, numberKind core.NumberKind, gos ...GaugeOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyGaugeOptions(&opts, gos...)
	return &mockInstrument{
		name:       name,
		kind:       mockKindGauge,
		numberKind: numberKind,
		opts:       opts,
	}
}

func (m *mockMeter) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	instrument := m.newMeasureInstrument(name, core.Int64NumberKind, mos...)
	return WrapInt64MeasureInstrument(instrument)
}

func (m *mockMeter) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	instrument := m.newMeasureInstrument(name, core.Float64NumberKind, mos...)
	return WrapFloat64MeasureInstrument(instrument)
}

func (m *mockMeter) newMeasureInstrument(name string, numberKind core.NumberKind, mos ...MeasureOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyMeasureOptions(&opts, mos...)
	return &mockInstrument{
		name:       name,
		kind:       mockKindMeasure,
		numberKind: numberKind,
		opts:       opts,
	}
}

func (m *mockMeter) RecordBatch(ctx context.Context, labels LabelSet, measurements ...Measurement) {
	ourLabelSet := labels.(*mockLabelSet)
	mm := make([]mockMeasurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = mockMeasurement{
			instrument: m.InstrumentImpl().(*mockInstrument),
			number:     m.Number(),
		}
	}
	m.recordMockBatch(ctx, ourLabelSet, mm...)
}

func (m *mockMeter) recordMockBatch(ctx context.Context, labelSet *mockLabelSet, measurements ...mockMeasurement) {
	m.measurementBatches = append(m.measurementBatches, batch{
		ctx:          ctx,
		labelSet:     labelSet,
		measurements: measurements,
	})
}
