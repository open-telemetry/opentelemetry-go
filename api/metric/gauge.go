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

type Float64Gauge struct {
	Gauge
}

type Int64Gauge struct {
	Gauge
}

type GaugeHandle struct {
	Handle
}

type Float64GaugeHandle struct {
	GaugeHandle
}

type Int64GaugeHandle struct {
	GaugeHandle
}

type GaugeOptionApplier interface {
	ApplyGaugeOption(*Descriptor)
}

type gaugeOptionWrapper struct {
	F Option
}

var _ GaugeOptionApplier = gaugeOptionWrapper{}

func (o gaugeOptionWrapper) ApplyGaugeOption(d *Descriptor) {
	o.F(d)
}

func NewGauge(name string, valueKind ValueKind, mos ...GaugeOptionApplier) (g Gauge) {
	registerDescriptor(name, GaugeKind, valueKind, &g.Descriptor)
	for _, opt := range mos {
		opt.ApplyGaugeOption(&g.Descriptor)
	}
	return
}

func NewFloat64Gauge(name string, mos ...GaugeOptionApplier) (g Float64Gauge) {
	g.Gauge = NewGauge(name, Float64ValueKind, mos...)
	return
}

func NewInt64Gauge(name string, mos ...GaugeOptionApplier) (g Int64Gauge) {
	g.Gauge = NewGauge(name, Int64ValueKind, mos...)
	return
}

func (g *Gauge) GetHandle(ctx context.Context, labels LabelSet) (h GaugeHandle) {
	h.Handle = labels.Meter().NewHandle(ctx, g.Descriptor, labels)
	return
}

func (g *Float64Gauge) GetHandle(ctx context.Context, labels LabelSet) (h Float64GaugeHandle) {
	h.GaugeHandle = g.Gauge.GetHandle(ctx, labels)
	return
}

func (g *Int64Gauge) GetHandle(ctx context.Context, labels LabelSet) (h Int64GaugeHandle) {
	h.GaugeHandle = g.Gauge.GetHandle(ctx, labels)
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

func (g *Float64Gauge) Measurement(value float64) Measurement {
	return g.Gauge.Float64Measurement(value)
}

func (g *Int64Gauge) Measurement(value int64) Measurement {
	return g.Gauge.Int64Measurement(value)
}

func (g *Gauge) Set(ctx context.Context, value MeasurementValue, labels LabelSet) {
	labels.Meter().RecordBatch(ctx, labels, Measurement{
		Descriptor: g.Descriptor,
		Value:      value,
	})
}

func (g *Float64Gauge) Set(ctx context.Context, value float64, labels LabelSet) {
	g.Gauge.Set(ctx, NewFloat64MeasurementValue(value), labels)
}

func (g *Int64Gauge) Set(ctx context.Context, value int64, labels LabelSet) {
	g.Gauge.Set(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *GaugeHandle) Set(ctx context.Context, value MeasurementValue) {
	h.RecordOne(ctx, value)
}

func (h *Float64GaugeHandle) Set(ctx context.Context, value float64) {
	h.GaugeHandle.Set(ctx, NewFloat64MeasurementValue(value))
}

func (h *Int64GaugeHandle) Set(ctx context.Context, value int64) {
	h.GaugeHandle.Set(ctx, NewInt64MeasurementValue(value))
}
