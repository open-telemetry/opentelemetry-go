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

type CounterHandle struct {
	Handle
}

func NewCounter(name string, mos ...Option) (c Counter) {
	registerDescriptor(name, CounterKind, mos, &c.Descriptor)
	return
}

func (c *Counter) GetHandle(ctx context.Context, labels LabelSet) (h CounterHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, c.Descriptor, labels)
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

func (c *Counter) Add(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: c.Descriptor,
		Value:      value,
	})
}

func (c *Counter) AddFloat64(ctx context.Context, value float64, labels LabelSet) {
	c.Add(ctx, NewFloat64MeasurementValue(value), labels)
}

func (c *Counter) AddInt64(ctx context.Context, value int64, labels LabelSet) {
	c.Add(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *CounterHandle) Add(ctx context.Context, value MeasurementValue) {
	h.Record(ctx, value)
}

func (h *CounterHandle) AddFloat64(ctx context.Context, value float64) {
	h.Add(ctx, NewFloat64MeasurementValue(value))
}

func (h *CounterHandle) AddInt64(ctx context.Context, value int64) {
	h.Add(ctx, NewInt64MeasurementValue(value))
}
