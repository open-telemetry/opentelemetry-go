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

type Counter struct {
	Descriptor
}

type Float64Counter struct {
	Counter
}

type Int64Counter struct {
	Counter
}

type CounterHandle struct {
	Handle
}

type Float64CounterHandle struct {
	CounterHandle
}

type Int64CounterHandle struct {
	CounterHandle
}

type CounterOptionApplier interface {
	ApplyCounterOption(*Descriptor)
}

type counterOptionWrapper struct {
	F Option
}

var _ CounterOptionApplier = counterOptionWrapper{}

func (o counterOptionWrapper) ApplyCounterOption(d *Descriptor) {
	o.F(d)
}

func NewCounter(name string, valueKind ValueKind, mos ...CounterOptionApplier) (c Counter) {
	registerDescriptor(name, CounterKind, valueKind, &c.Descriptor)
	for _, opt := range mos {
		opt.ApplyCounterOption(&c.Descriptor)
	}
	return
}

func NewFloat64Counter(name string, mos ...CounterOptionApplier) (c Float64Counter) {
	c.Counter = NewCounter(name, Float64ValueKind, mos...)
	return
}

func NewInt64Counter(name string, mos ...CounterOptionApplier) (c Int64Counter) {
	c.Counter = NewCounter(name, Int64ValueKind, mos...)
	return
}

func (c *Counter) GetHandle(ctx context.Context, labels LabelSet) (h CounterHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, c.Descriptor, labels)
	return
}

func (c *Float64Counter) GetHandle(ctx context.Context, labels LabelSet) (h Float64CounterHandle) {
	h.CounterHandle = c.Counter.GetHandle(ctx, labels)
	return
}

func (c *Int64Counter) GetHandle(ctx context.Context, labels LabelSet) (h Int64CounterHandle) {
	h.CounterHandle = c.Counter.GetHandle(ctx, labels)
	return
}

func (c *Counter) Float64Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: c.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (c *Counter) Int64Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: c.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (c *Float64Counter) Measurement(value float64) Measurement {
	return c.Counter.Float64Measurement(value)
}

func (c *Int64Counter) Measurement(value int64) Measurement {
	return c.Counter.Int64Measurement(value)
}

func (c *Counter) Add(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: c.Descriptor,
		Value:      value,
	})
}

func (c *Float64Counter) Add(ctx context.Context, value float64, labels LabelSet) {
	c.Counter.Add(ctx, NewFloat64MeasurementValue(value), labels)
}

func (c *Int64Counter) Add(ctx context.Context, value int64, labels LabelSet) {
	c.Counter.Add(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *CounterHandle) Add(ctx context.Context, value MeasurementValue) {
	h.RecordOne(ctx, value)
}

func (h *Float64CounterHandle) Add(ctx context.Context, value float64) {
	h.CounterHandle.Add(ctx, NewFloat64MeasurementValue(value))
}

func (h *Int64CounterHandle) Add(ctx context.Context, value int64) {
	h.CounterHandle.Add(ctx, NewInt64MeasurementValue(value))
}
