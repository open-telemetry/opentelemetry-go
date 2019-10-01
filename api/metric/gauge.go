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

type Gauge struct {
	Descriptor
}

type GaugeHandle struct {
	Handle
}

func NewGauge(name string, mos ...Option) (g Gauge) {
	registerDescriptor(name, GaugeKind, mos, &g.Descriptor)
	return
}

func (g *Gauge) GetHandle(ctx context.Context, labels LabelSet) (h GaugeHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, g.Descriptor, labels)
	return
}

func (g *Gauge) Float64Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: g.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (g *Gauge) Int64Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: g.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (g *Gauge) Set(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: g.Descriptor,
		Value:      value,
	})
}

func (g *Gauge) SetFloat64(ctx context.Context, value float64, labels LabelSet) {
	g.Set(ctx, NewFloat64MeasurementValue(value), labels)
}

func (g *Gauge) SetInt64(ctx context.Context, value int64, labels LabelSet) {
	g.Set(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *GaugeHandle) Set(ctx context.Context, value MeasurementValue) {
	h.Record(ctx, value)
}

func (h *GaugeHandle) SetFloat64(ctx context.Context, value float64) {
	h.Set(ctx, NewFloat64MeasurementValue(value))
}

func (h *GaugeHandle) SetInt64(ctx context.Context, value int64) {
	h.Set(ctx, NewInt64MeasurementValue(value))
}
