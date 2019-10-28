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
		name      string
		kind      mockKind
		valueKind ValueKind
		opts      Options
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
		value      MeasurementValue
	}
)

var (
	_ Instrument = &mockInstrument{}
	_ Handle     = &mockHandle{}
	_ LabelSet   = &mockLabelSet{}
	_ Meter      = &mockMeter{}
)

const (
	mockKindCounter mockKind = iota
	mockKindGauge
	mockKindMeasure
)

func (i *mockInstrument) AcquireHandle(labels LabelSet) Handle {
	return &mockHandle{
		instrument: i,
		labelSet:   labels.(*mockLabelSet),
	}
}

func (i *mockInstrument) RecordOne(ctx context.Context, value MeasurementValue, labels LabelSet) {
	doRecordBatch(labels.(*mockLabelSet), ctx, i, value)
}

func (h *mockHandle) RecordOne(ctx context.Context, value MeasurementValue) {
	doRecordBatch(h.labelSet, ctx, h.instrument, value)
}

func (h *mockHandle) Release() {
}

func doRecordBatch(labelSet *mockLabelSet, ctx context.Context, instrument *mockInstrument, value MeasurementValue) {
	labelSet.meter.recordMockBatch(ctx, labelSet, mockMeasurement{
		instrument: instrument,
		value:      value,
	})
}

func (s *mockLabelSet) Meter() Meter {
	return s.meter
}

func newMockMeter() *mockMeter {
	return &mockMeter{}
}

func (m *mockMeter) Labels(ctx context.Context, labels ...core.KeyValue) LabelSet {
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
	instrument := m.newCounterInstrument(name, Int64ValueKind, cos...)
	return WrapInt64CounterInstrument(instrument)
}

func (m *mockMeter) NewFloat64Counter(name string, cos ...CounterOptionApplier) Float64Counter {
	instrument := m.newCounterInstrument(name, Float64ValueKind, cos...)
	return WrapFloat64CounterInstrument(instrument)
}

func (m *mockMeter) newCounterInstrument(name string, valueKind ValueKind, cos ...CounterOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyCounterOptions(&opts, cos...)
	return &mockInstrument{
		name:      name,
		kind:      mockKindCounter,
		valueKind: valueKind,
		opts:      opts,
	}
}

func (m *mockMeter) NewInt64Gauge(name string, gos ...GaugeOptionApplier) Int64Gauge {
	instrument := m.newGaugeInstrument(name, Int64ValueKind, gos...)
	return WrapInt64GaugeInstrument(instrument)
}

func (m *mockMeter) NewFloat64Gauge(name string, gos ...GaugeOptionApplier) Float64Gauge {
	instrument := m.newGaugeInstrument(name, Float64ValueKind, gos...)
	return WrapFloat64GaugeInstrument(instrument)
}

func (m *mockMeter) newGaugeInstrument(name string, valueKind ValueKind, gos ...GaugeOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyGaugeOptions(&opts, gos...)
	return &mockInstrument{
		name:      name,
		kind:      mockKindGauge,
		valueKind: valueKind,
		opts:      opts,
	}
}

func (m *mockMeter) NewInt64Measure(name string, mos ...MeasureOptionApplier) Int64Measure {
	instrument := m.newMeasureInstrument(name, Int64ValueKind, mos...)
	return WrapInt64MeasureInstrument(instrument)
}

func (m *mockMeter) NewFloat64Measure(name string, mos ...MeasureOptionApplier) Float64Measure {
	instrument := m.newMeasureInstrument(name, Float64ValueKind, mos...)
	return WrapFloat64MeasureInstrument(instrument)
}

func (m *mockMeter) newMeasureInstrument(name string, valueKind ValueKind, mos ...MeasureOptionApplier) *mockInstrument {
	opts := Options{}
	ApplyMeasureOptions(&opts, mos...)
	return &mockInstrument{
		name:      name,
		kind:      mockKindMeasure,
		valueKind: valueKind,
		opts:      opts,
	}
}

func (m *mockMeter) RecordBatch(ctx context.Context, labels LabelSet, measurements ...Measurement) {
	ourLabelSet := labels.(*mockLabelSet)
	mm := make([]mockMeasurement, len(measurements))
	for i := 0; i < len(measurements); i++ {
		m := measurements[i]
		mm[i] = mockMeasurement{
			instrument: m.Instrument().(*mockInstrument),
			value:      m.Value(),
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
