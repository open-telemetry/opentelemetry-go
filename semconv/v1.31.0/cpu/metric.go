// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/cpu"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ModeAttr is an attribute conforming to the cpu.mode semantic conventions. It
// represents the mode of the CPU.
type ModeAttr string

var (
	// ModeUser is the none.
	ModeUser ModeAttr = "user"
	// ModeSystem is the none.
	ModeSystem ModeAttr = "system"
	// ModeNice is the none.
	ModeNice ModeAttr = "nice"
	// ModeIdle is the none.
	ModeIdle ModeAttr = "idle"
	// ModeIowait is the none.
	ModeIowait ModeAttr = "iowait"
	// ModeInterrupt is the none.
	ModeInterrupt ModeAttr = "interrupt"
	// ModeSteal is the none.
	ModeSteal ModeAttr = "steal"
	// ModeKernel is the none.
	ModeKernel ModeAttr = "kernel"
)

// Frequency is an instrument used to record metric values conforming to the
// "cpu.frequency" semantic conventions. It represents the operating frequency of
// the logical CPU in Hertz.
type Frequency struct {
	inst metric.Int64Gauge
}

// NewFrequency returns a new Frequency instrument.
func NewFrequency(m metric.Meter) (Frequency, error) {
	i, err := m.Int64Gauge(
	    "cpu.frequency",
	    metric.WithDescription("Operating frequency of the logical CPU in Hertz."),
	    metric.WithUnit("Hz"),
	)
	if err != nil {
	    return Frequency{}, err
	}
	return Frequency{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (Frequency) Name() string {
	return "cpu.frequency"
}

// Unit returns the semantic convention unit of the instrument
func (Frequency) Unit() string {
	return "Hz"
}

// Description returns the semantic convention description of the instrument
func (Frequency) Description() string {
	return "Operating frequency of the logical CPU in Hertz."
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m Frequency) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrLogicalNumber returns an optional attribute for the "cpu.logical_number"
// semantic convention. It represents the logical CPU number [0..n-1].
func (Frequency) AttrLogicalNumber(val int) attribute.KeyValue {
	return attribute.Int("cpu.logical_number", val)
}

// Time is an instrument used to record metric values conforming to the
// "cpu.time" semantic conventions. It represents the seconds each logical CPU
// spent on each mode.
type Time struct {
	inst metric.Float64Counter
}

// NewTime returns a new Time instrument.
func NewTime(m metric.Meter) (Time, error) {
	i, err := m.Float64Counter(
	    "cpu.time",
	    metric.WithDescription("Seconds each logical CPU spent on each mode"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return Time{}, err
	}
	return Time{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (Time) Name() string {
	return "cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (Time) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (Time) Description() string {
	return "Seconds each logical CPU spent on each mode"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m Time) Add(
	ctx context.Context,
	incr float64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrLogicalNumber returns an optional attribute for the "cpu.logical_number"
// semantic convention. It represents the logical CPU number [0..n-1].
func (Time) AttrLogicalNumber(val int) attribute.KeyValue {
	return attribute.Int("cpu.logical_number", val)
}

// AttrMode returns an optional attribute for the "cpu.mode" semantic convention.
// It represents the mode of the CPU.
func (Time) AttrMode(val ModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// Utilization is an instrument used to record metric values conforming to the
// "cpu.utilization" semantic conventions. It represents the for each logical
// CPU, the utilization is calculated as the change in cumulative CPU time
// (cpu.time) over a measurement interval, divided by the elapsed time.
type Utilization struct {
	inst metric.Int64Gauge
}

// NewUtilization returns a new Utilization instrument.
func NewUtilization(m metric.Meter) (Utilization, error) {
	i, err := m.Int64Gauge(
	    "cpu.utilization",
	    metric.WithDescription("For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time."),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return Utilization{}, err
	}
	return Utilization{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (Utilization) Name() string {
	return "cpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (Utilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (Utilization) Description() string {
	return "For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time."
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m Utilization) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrLogicalNumber returns an optional attribute for the "cpu.logical_number"
// semantic convention. It represents the logical CPU number [0..n-1].
func (Utilization) AttrLogicalNumber(val int) attribute.KeyValue {
	return attribute.Int("cpu.logical_number", val)
}

// AttrMode returns an optional attribute for the "cpu.mode" semantic convention.
// It represents the mode of the CPU.
func (Utilization) AttrMode(val ModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}