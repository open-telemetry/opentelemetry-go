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
)

type Measure struct {
	Descriptor
}

type MeasureHandle struct {
	Handle
}

func NewMeasure(name string, mos ...Option) (m Measure) {
	registerDescriptor(name, MeasureKind, mos, &m.Descriptor)
	return
}

func (m *Measure) GetHandle(ctx context.Context, labels LabelSet) (h MeasureHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, m.Descriptor, labels)
	return
}

func (m *Measure) Float64Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (m *Measure) Int64Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (m *Measure) Record(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: m.Descriptor,
		Value:      value,
	})
}

func (m *Measure) RecordFloat64(ctx context.Context, value float64, labels LabelSet) {
	m.Record(ctx, NewFloat64MeasurementValue(value), labels)
}

func (m *Measure) RecordInt64(ctx context.Context, value int64, labels LabelSet) {
	m.Record(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *MeasureHandle) Record(ctx context.Context, value MeasurementValue) {
	h.Handle.Record(ctx, value)
}

func (h *MeasureHandle) RecordFloat64(ctx context.Context, value float64) {
	h.Record(ctx, NewFloat64MeasurementValue(value))
}

func (h *MeasureHandle) RecordInt64(ctx context.Context, value int64) {
	h.Record(ctx, NewInt64MeasurementValue(value))
}
