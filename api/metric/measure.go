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

type Float64Measure struct {
	Measure
}

type Int64Measure struct {
	Measure
}

type MeasureHandle struct {
	Handle
}

type Float64MeasureHandle struct {
	MeasureHandle
}

type Int64MeasureHandle struct {
	MeasureHandle
}

type MeasureOptionApplier interface {
	ApplyMeasureOption(*Descriptor)
}

type measureOptionWrapper struct {
	F Option
}

var _ MeasureOptionApplier = measureOptionWrapper{}

func (o measureOptionWrapper) ApplyMeasureOption(d *Descriptor) {
	o.F(d)
}

func NewMeasure(name string, valueKind ValueKind, mos ...MeasureOptionApplier) (m Measure) {
	registerDescriptor(name, MeasureKind, valueKind, &m.Descriptor)
	for _, opt := range mos {
		opt.ApplyMeasureOption(&m.Descriptor)
	}
	return
}

func NewFloat64Measure(name string, mos ...MeasureOptionApplier) (c Float64Measure) {
	c.Measure = NewMeasure(name, Float64ValueKind, mos...)
	return
}

func NewInt64Measure(name string, mos ...MeasureOptionApplier) (c Int64Measure) {
	c.Measure = NewMeasure(name, Int64ValueKind, mos...)
	return
}

func (m *Measure) GetHandle(ctx context.Context, labels LabelSet) (h MeasureHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, m.Descriptor, labels)
	return
}

func (c *Float64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Float64MeasureHandle) {
	h.MeasureHandle = c.Measure.GetHandle(ctx, labels)
	return
}

func (c *Int64Measure) GetHandle(ctx context.Context, labels LabelSet) (h Int64MeasureHandle) {
	h.MeasureHandle = c.Measure.GetHandle(ctx, labels)
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

func (c *Float64Measure) Measurement(value float64) Measurement {
	return c.Measure.Float64Measurement(value)
}

func (c *Int64Measure) Measurement(value int64) Measurement {
	return c.Measure.Int64Measurement(value)
}

func (m *Measure) Record(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: m.Descriptor,
		Value:      value,
	})
}

func (c *Float64Measure) Record(ctx context.Context, value float64, labels LabelSet) {
	c.Measure.Record(ctx, NewFloat64MeasurementValue(value), labels)
}

func (c *Int64Measure) Record(ctx context.Context, value int64, labels LabelSet) {
	c.Measure.Record(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *MeasureHandle) Record(ctx context.Context, value MeasurementValue) {
	h.RecordOne(ctx, value)
}

func (h *Float64MeasureHandle) Record(ctx context.Context, value float64) {
	h.MeasureHandle.Record(ctx, NewFloat64MeasurementValue(value))
}

func (h *Int64MeasureHandle) Record(ctx context.Context, value int64) {
	h.MeasureHandle.Record(ctx, NewInt64MeasurementValue(value))
}
