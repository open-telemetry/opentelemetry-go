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

type Float64Counter struct {
	CommonMetric
}

type Int64Counter struct {
	CommonMetric
}

type Float64CounterHandle struct {
	Handle
}

type Int64CounterHandle struct {
	Handle
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

func newCounter(name string, valueKind ValueKind, mos ...CounterOptionApplier) CommonMetric {
	m := registerCommonMetric(name, CounterKind, valueKind)
	for _, opt := range mos {
		opt.ApplyCounterOption(m.Descriptor)
	}
	return m
}

func NewFloat64Counter(name string, mos ...CounterOptionApplier) (c Float64Counter) {
	c.CommonMetric = newCounter(name, Float64ValueKind, mos...)
	return
}

func NewInt64Counter(name string, mos ...CounterOptionApplier) (c Int64Counter) {
	c.CommonMetric = newCounter(name, Int64ValueKind, mos...)
	return
}

func (c *Float64Counter) GetHandle(labels LabelSet) (h Float64CounterHandle) {
	h.Handle = c.getHandle(labels)
	return
}

func (c *Int64Counter) GetHandle(labels LabelSet) (h Int64CounterHandle) {
	h.Handle = c.getHandle(labels)
	return
}

func (c *Float64Counter) Measurement(value float64) Measurement {
	return c.float64Measurement(value)
}

func (c *Int64Counter) Measurement(value int64) Measurement {
	return c.int64Measurement(value)
}

func (c *Float64Counter) Add(ctx context.Context, value float64, labels LabelSet) {
	c.recordOne(ctx, NewFloat64MeasurementValue(value), labels)
}

func (c *Int64Counter) Add(ctx context.Context, value int64, labels LabelSet) {
	c.recordOne(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *Float64CounterHandle) Add(ctx context.Context, value float64) {
	h.RecordOne(ctx, NewFloat64MeasurementValue(value))
}

func (h *Int64CounterHandle) Add(ctx context.Context, value int64) {
	h.RecordOne(ctx, NewInt64MeasurementValue(value))
}
