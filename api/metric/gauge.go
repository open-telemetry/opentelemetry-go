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
	CommonMetric
}

type Int64Gauge struct {
	CommonMetric
}

type Float64GaugeHandle struct {
	Handle
}

type Int64GaugeHandle struct {
	Handle
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

func newGauge(name string, valueKind ValueKind, mos ...GaugeOptionApplier) CommonMetric {
	m := registerCommonMetric(name, GaugeKind, valueKind)
	for _, opt := range mos {
		opt.ApplyGaugeOption(m.Descriptor)
	}
	return m
}

func NewFloat64Gauge(name string, mos ...GaugeOptionApplier) (g Float64Gauge) {
	g.CommonMetric = newGauge(name, Float64ValueKind, mos...)
	return
}

func NewInt64Gauge(name string, mos ...GaugeOptionApplier) (g Int64Gauge) {
	g.CommonMetric = newGauge(name, Int64ValueKind, mos...)
	return
}

func (g *Float64Gauge) GetHandle(labels LabelSet) (h Float64GaugeHandle) {
	h.Handle = g.getHandle(labels)
	return
}

func (g *Int64Gauge) GetHandle(labels LabelSet) (h Int64GaugeHandle) {
	h.Handle = g.getHandle(labels)
	return
}

func (g *Float64Gauge) Measurement(value float64) Measurement {
	return g.float64Measurement(value)
}

func (g *Int64Gauge) Measurement(value int64) Measurement {
	return g.int64Measurement(value)
}

func (g *Float64Gauge) Set(ctx context.Context, value float64, labels LabelSet) {
	g.recordOne(ctx, NewFloat64MeasurementValue(value), labels)
}

func (g *Int64Gauge) Set(ctx context.Context, value int64, labels LabelSet) {
	g.recordOne(ctx, NewInt64MeasurementValue(value), labels)
}

func (h *Float64GaugeHandle) Set(ctx context.Context, value float64) {
	h.RecordOne(ctx, NewFloat64MeasurementValue(value))
}

func (h *Int64GaugeHandle) Set(ctx context.Context, value int64) {
	h.RecordOne(ctx, NewInt64MeasurementValue(value))
}
