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

type Float64Gauge struct {
	Descriptor
}

type Int64Gauge struct {
	Descriptor
}

type Float64GaugeHandle struct {
	Handle
}

type Int64GaugeHandle struct {
	Handle
}

func NewFloat64Gauge(name string, mos ...Option) (g Float64Gauge) {
	registerDescriptor(name, GaugeKind, mos, &g.Descriptor)
	return
}

func NewInt64Gauge(name string, mos ...Option) (g Int64Gauge) {
	registerDescriptor(name, GaugeKind, mos, &g.Descriptor)
	return
}

func (g *Float64Gauge) GetHandle(ctx context.Context, labels LabelSet) (h Float64GaugeHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, g.Descriptor, labels)
	return
}

func (g *Int64Gauge) GetHandle(ctx context.Context, labels LabelSet) (h Int64GaugeHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, g.Descriptor, labels)
	return
}

func (g *Float64Gauge) Measurement(value float64) Measurement {
	return Measurement{
		Descriptor: g.Descriptor,
		Value:      NewFloat64MeasurementValue(value),
	}
}

func (g *Int64Gauge) Measurement(value int64) Measurement {
	return Measurement{
		Descriptor: g.Descriptor,
		Value:      NewInt64MeasurementValue(value),
	}
}

func (g *Float64Gauge) Set(ctx context.Context, value float64, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, g.Measurement(value))
}

func (g *Int64Gauge) Set(ctx context.Context, value int64, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, g.Measurement(value))
}

func (h *Float64GaugeHandle) Set(ctx context.Context, value float64) {
	h.RecordFloat(ctx, value)
}

func (h *Int64GaugeHandle) Set(ctx context.Context, value int64) {
	h.RecordInt(ctx, value)
}
