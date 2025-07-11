// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "cpu" namespace.
package cpuconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
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
	// ModeIOWait is the none.
	ModeIOWait ModeAttr = "iowait"
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
	metric.Int64Gauge
}

// NewFrequency returns a new Frequency instrument.
func NewFrequency(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (Frequency, error) {
	// Check if the meter is nil.
	if m == nil {
		return Frequency{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"cpu.frequency",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Operating frequency of the logical CPU in Hertz."),
			metric.WithUnit("Hz"),
		}, opt...)...,
	)
	if err != nil {
	    return Frequency{noop.Int64Gauge{}}, err
	}
	return Frequency{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Frequency) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
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
	metric.Float64ObservableCounter
}

// NewTime returns a new Time instrument.
func NewTime(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (Time, error) {
	// Check if the meter is nil.
	if m == nil {
		return Time{noop.Float64ObservableCounter{}}, nil
	}

	i, err := m.Float64ObservableCounter(
		"cpu.time",
		append([]metric.Float64ObservableCounterOption{
			metric.WithDescription("Seconds each logical CPU spent on each mode"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return Time{noop.Float64ObservableCounter{}}, err
	}
	return Time{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Time) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
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
	metric.Int64Gauge
}

// NewUtilization returns a new Utilization instrument.
func NewUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (Utilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return Utilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"cpu.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("For each logical CPU, the utilization is calculated as the change in cumulative CPU time (cpu.time) over a measurement interval, divided by the elapsed time."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return Utilization{noop.Int64Gauge{}}, err
	}
	return Utilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Utilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
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