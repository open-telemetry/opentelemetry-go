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

type Float64Measure struct {
	Descriptor
}

type Int64Measure struct {
	Descriptor
}

type Float64MeasureHandle struct {
	Handle
}

type Int64MeasureHandle struct {
	Handle
}

func NewFloat64Measure(name string, mos ...Option) (m Float64Measure) {
	registerDescriptor(name, MeasureKind, mos, &m.Descriptor)
	return
}

func NewInt64Measure(name string, mos ...Option) (m Int64Measure) {
	registerDescriptor(name, MeasureKind, mos, &m.Descriptor)
	return
}

func (m *Float64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Float64MeasureHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, m.Descriptor, labels)
	return
}

func (m *Int64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Int64MeasureHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, m.Descriptor, labels)
	return
}

func (m *Float64Measure) Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (m *Int64Measure) Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: m.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (m *Float64Measure) Record(ctx context.Context, value float64, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, m.Measurement(value))
}

func (m *Int64Measure) Record(ctx context.Context, value int64, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, m.Measurement(value))
}

func (h *Float64MeasureHandle) Record(ctx context.Context, value float64) {
	h.RecordFloat(ctx, value)
}

func (h *Int64MeasureHandle) Record(ctx context.Context, value int64) {
	h.RecordInt(ctx, value)
}
