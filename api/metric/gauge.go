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

// Float64Gauge is a metric that stores the last float64 value.
type Float64Gauge struct {
	CommonMetric
}

// Int64Gauge is a metric that stores the last int64 value.
type Int64Gauge struct {
	CommonMetric
}

// Float64GaugeHandle is a handle for Float64Gauge.
type Float64GaugeHandle struct {
	Handle
}

// Int64GaugeHandle is a handle for Int64Gauge.
type Int64GaugeHandle struct {
	Handle
}

// GaugeOptionApplier is an interface for applying metric options that
// are valid only for gauge metrics.
type GaugeOptionApplier interface {
	// ApplyGaugeOption is used to make some gauge-specific
	// changes in the Descriptor.
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

// NewFloat64Gauge creates a new gauge for float64.
func NewFloat64Gauge(name string, mos ...GaugeOptionApplier) (g Float64Gauge) {
	g.CommonMetric = newGauge(name, Float64ValueKind, mos...)
	return
}

// NewInt64Gauge creates a new gauge for int64.
func NewInt64Gauge(name string, mos ...GaugeOptionApplier) (g Int64Gauge) {
	g.CommonMetric = newGauge(name, Int64ValueKind, mos...)
	return
}

// GetHandle creates a handle for this gauge. The labels should
// contain the keys and values specified in the gauge with the
// WithKeys option.
func (g *Float64Gauge) GetHandle(labels LabelSet) (h Float64GaugeHandle) {
	h.Handle = g.getHandle(labels)
	return
}

// GetHandle creates a handle for this gauge. The labels should
// contain the keys and values specified in the gauge with the
// WithKeys option.
func (g *Int64Gauge) GetHandle(labels LabelSet) (h Int64GaugeHandle) {
	h.Handle = g.getHandle(labels)
	return
}

// Measurement creates a Measurement object to use with batch
// recording.
func (g *Float64Gauge) Measurement(value float64) Measurement {
	return g.float64Measurement(value)
}

// Measurement creates a Measurement object to use with batch
// recording.
func (g *Int64Gauge) Measurement(value int64) Measurement {
	return g.int64Measurement(value)
}

// Set sets the value of the gauge to the passed value. The labels
// should contain the keys and values specified in the gauge with the
// WithKeys option.
func (g *Float64Gauge) Set(ctx context.Context, value float64, labels LabelSet) {
	g.recordOne(ctx, NewFloat64MeasurementValue(value), labels)
}

// Set sets the value of the gauge to the passed value. The labels
// should contain the keys and values specified in the gauge with the
// WithKeys option.
func (g *Int64Gauge) Set(ctx context.Context, value int64, labels LabelSet) {
	g.recordOne(ctx, NewInt64MeasurementValue(value), labels)
}

// Set sets the value of the gauge to the passed value.
func (h *Float64GaugeHandle) Set(ctx context.Context, value float64) {
	h.RecordOne(ctx, NewFloat64MeasurementValue(value))
}

// Set sets the value of the gauge to the passed value.
func (h *Int64GaugeHandle) Set(ctx context.Context, value int64) {
	h.RecordOne(ctx, NewInt64MeasurementValue(value))
}
