// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package processconv provides types and functionality for OpenTelemetry semantic
// conventions in the "process" namespace.
package processconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
)

// CPUModeAttr is an attribute conforming to the cpu.mode semantic conventions.
// It represents the CPU mode for this data point.
type CPUModeAttr string

var (
	// CPUModeUser is the user.
	CPUModeUser CPUModeAttr = "user"
	// CPUModeSystem is the system.
	CPUModeSystem CPUModeAttr = "system"
	// CPUModeNice is the nice.
	CPUModeNice CPUModeAttr = "nice"
	// CPUModeIdle is the idle.
	CPUModeIdle CPUModeAttr = "idle"
	// CPUModeIOWait is the IO Wait.
	CPUModeIOWait CPUModeAttr = "iowait"
	// CPUModeInterrupt is the interrupt.
	CPUModeInterrupt CPUModeAttr = "interrupt"
	// CPUModeSteal is the steal.
	CPUModeSteal CPUModeAttr = "steal"
)

// DiskIODirectionAttr is an attribute conforming to the disk.io.direction
// semantic conventions. It represents the disk IO operation direction.
type DiskIODirectionAttr string

var (
	// DiskIODirectionRead is the standardized value "read" of DiskIODirectionAttr.
	DiskIODirectionRead DiskIODirectionAttr = "read"
	// DiskIODirectionWrite is the standardized value "write" of
	// DiskIODirectionAttr.
	DiskIODirectionWrite DiskIODirectionAttr = "write"
)

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the standardized value "transmit" of
	// NetworkIODirectionAttr.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the standardized value "receive" of
	// NetworkIODirectionAttr.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
)

// ContextSwitchTypeAttr is an attribute conforming to the
// process.context_switch.type semantic conventions. It represents the specifies
// whether the context switches for this data point were voluntary or
// involuntary.
type ContextSwitchTypeAttr string

var (
	// ContextSwitchTypeVoluntary is the standardized value "voluntary" of
	// ContextSwitchTypeAttr.
	ContextSwitchTypeVoluntary ContextSwitchTypeAttr = "voluntary"
	// ContextSwitchTypeInvoluntary is the standardized value "involuntary" of
	// ContextSwitchTypeAttr.
	ContextSwitchTypeInvoluntary ContextSwitchTypeAttr = "involuntary"
)

// SystemPagingFaultTypeAttr is an attribute conforming to the
// system.paging.fault.type semantic conventions. It represents the type of
// paging fault. Value MUST be either `major` or `minor`. If the metric is
// reported without this attribute, it should be the sum of major and minor page
// faults.
type SystemPagingFaultTypeAttr string

var (
	// SystemPagingFaultTypeMajor is the standardized value "major" of
	// SystemPagingFaultTypeAttr.
	SystemPagingFaultTypeMajor SystemPagingFaultTypeAttr = "major"
	// SystemPagingFaultTypeMinor is the standardized value "minor" of
	// SystemPagingFaultTypeAttr.
	SystemPagingFaultTypeMinor SystemPagingFaultTypeAttr = "minor"
)

// ContextSwitches is an instrument used to record metric values conforming to
// the "process.context_switches" semantic conventions. It represents the number
// of times the process has been context switched.
type ContextSwitches struct {
	metric.Int64Counter
}

var newContextSwitchesOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of times the process has been context switched."),
	metric.WithUnit("{context_switch}"),
}

// NewContextSwitches returns a new ContextSwitches instrument.
func NewContextSwitches(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ContextSwitches, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContextSwitches{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContextSwitchesOpts
	} else {
		opt = append(opt, newContextSwitchesOpts...)
	}

	i, err := m.Int64Counter(
		"process.context_switches",
		opt...,
	)
	if err != nil {
		return ContextSwitches{noop.Int64Counter{}}, err
	}
	return ContextSwitches{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContextSwitches) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ContextSwitches) Name() string {
	return "process.context_switches"
}

// Unit returns the semantic convention unit of the instrument
func (ContextSwitches) Unit() string {
	return "{context_switch}"
}

// Description returns the semantic convention description of the instrument
func (ContextSwitches) Description() string {
	return "Number of times the process has been context switched."
}

// Add adds incr to the existing count for attrs.
//
// The contextSwitchType is the specifies whether the context switches for this
// data point were voluntary or involuntary.
func (m ContextSwitches) Add(
	ctx context.Context,
	incr int64,
	contextSwitchType ContextSwitchTypeAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("process.context_switch.type", string(contextSwitchType)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("process.context_switch.type", string(contextSwitchType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ContextSwitches) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// ContextSwitchesObservable is an instrument used to record metric values
// conforming to the "process.context_switches" semantic conventions. It
// represents the number of times the process has been context switched.
type ContextSwitchesObservable struct {
	metric.Int64ObservableCounter
}

var newContextSwitchesObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Number of times the process has been context switched."),
	metric.WithUnit("{context_switch}"),
}

// NewContextSwitchesObservable returns a new ContextSwitchesObservable
// instrument.
func NewContextSwitchesObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (ContextSwitchesObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContextSwitchesObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContextSwitchesObservableOpts
	} else {
		opt = append(opt, newContextSwitchesObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.context_switches",
		opt...,
	)
	if err != nil {
		return ContextSwitchesObservable{noop.Int64ObservableCounter{}}, err
	}
	return ContextSwitchesObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContextSwitchesObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (ContextSwitchesObservable) Name() string {
	return "process.context_switches"
}

// Unit returns the semantic convention unit of the instrument
func (ContextSwitchesObservable) Unit() string {
	return "{context_switch}"
}

// Description returns the semantic convention description of the instrument
func (ContextSwitchesObservable) Description() string {
	return "Number of times the process has been context switched."
}

// AttrContextSwitchType returns a required attribute for the
// "process.context_switch.type" semantic convention. It represents the specifies
// whether the context switches for this data point were voluntary or
// involuntary.
func (ContextSwitchesObservable) AttrContextSwitchType(val ContextSwitchTypeAttr) attribute.KeyValue {
	return attribute.String("process.context_switch.type", string(val))
}

// CPUTime is an instrument used to record metric values conforming to the
// "process.cpu.time" semantic conventions. It represents the total CPU seconds
// broken down by different CPU modes.
type CPUTime struct {
	metric.Float64ObservableCounter
}

var newCPUTimeOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Total CPU seconds broken down by different CPU modes."),
	metric.WithUnit("s"),
}

// NewCPUTime returns a new CPUTime instrument.
func NewCPUTime(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (CPUTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUTime{noop.Float64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUTimeOpts
	} else {
		opt = append(opt, newCPUTimeOpts...)
	}

	i, err := m.Float64ObservableCounter(
		"process.cpu.time",
		opt...,
	)
	if err != nil {
		return CPUTime{noop.Float64ObservableCounter{}}, err
	}
	return CPUTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUTime) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (CPUTime) Name() string {
	return "process.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (CPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (CPUTime) Description() string {
	return "Total CPU seconds broken down by different CPU modes."
}

// AttrCPUMode returns a required attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point.
func (CPUTime) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUUtilization is an instrument used to record metric values conforming to the
// "process.cpu.utilization" semantic conventions. It represents the difference
// in process.cpu.time since the last measurement, divided by the elapsed time
// and number of CPUs available to the process.
type CPUUtilization struct {
	metric.Int64Gauge
}

var newCPUUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."),
	metric.WithUnit("1"),
}

// NewCPUUtilization returns a new CPUUtilization instrument.
func NewCPUUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (CPUUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUtilization{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUUtilizationOpts
	} else {
		opt = append(opt, newCPUUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"process.cpu.utilization",
		opt...,
	)
	if err != nil {
		return CPUUtilization{noop.Int64Gauge{}}, err
	}
	return CPUUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (CPUUtilization) Name() string {
	return "process.cpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (CPUUtilization) Description() string {
	return "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."
}

// Record records val to the current distribution for attrs.
//
// The cpuMode is the the CPU mode for this data point.
func (m CPUUtilization) Record(
	ctx context.Context,
	val int64,
	cpuMode CPUModeAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("cpu.mode", string(cpuMode)),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("cpu.mode", string(cpuMode)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m CPUUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// CPUUtilizationObservable is an instrument used to record metric values
// conforming to the "process.cpu.utilization" semantic conventions. It
// represents the difference in process.cpu.time since the last measurement,
// divided by the elapsed time and number of CPUs available to the process.
type CPUUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newCPUUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."),
	metric.WithUnit("1"),
}

// NewCPUUtilizationObservable returns a new CPUUtilizationObservable instrument.
func NewCPUUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (CPUUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUUtilizationObservableOpts
	} else {
		opt = append(opt, newCPUUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"process.cpu.utilization",
		opt...,
	)
	if err != nil {
		return CPUUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return CPUUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (CPUUtilizationObservable) Name() string {
	return "process.cpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (CPUUtilizationObservable) Description() string {
	return "Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."
}

// AttrCPUMode returns a required attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point.
func (CPUUtilizationObservable) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// DiskIO is an instrument used to record metric values conforming to the
// "process.disk.io" semantic conventions. It represents the disk bytes
// transferred.
type DiskIO struct {
	metric.Int64Counter
}

var newDiskIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Disk bytes transferred."),
	metric.WithUnit("By"),
}

// NewDiskIO returns a new DiskIO instrument.
func NewDiskIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (DiskIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskIO{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDiskIOOpts
	} else {
		opt = append(opt, newDiskIOOpts...)
	}

	i, err := m.Int64Counter(
		"process.disk.io",
		opt...,
	)
	if err != nil {
		return DiskIO{noop.Int64Counter{}}, err
	}
	return DiskIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (DiskIO) Name() string {
	return "process.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIO) Description() string {
	return "Disk bytes transferred."
}

// Add adds incr to the existing count for attrs.
//
// The diskIoDirection is the the disk IO operation direction.
func (m DiskIO) Add(
	ctx context.Context,
	incr int64,
	diskIoDirection DiskIODirectionAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("disk.io.direction", string(diskIoDirection)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("disk.io.direction", string(diskIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m DiskIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// DiskIOObservable is an instrument used to record metric values conforming to
// the "process.disk.io" semantic conventions. It represents the disk bytes
// transferred.
type DiskIOObservable struct {
	metric.Int64ObservableCounter
}

var newDiskIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Disk bytes transferred."),
	metric.WithUnit("By"),
}

// NewDiskIOObservable returns a new DiskIOObservable instrument.
func NewDiskIOObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (DiskIOObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskIOObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDiskIOObservableOpts
	} else {
		opt = append(opt, newDiskIOObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.disk.io",
		opt...,
	)
	if err != nil {
		return DiskIOObservable{noop.Int64ObservableCounter{}}, err
	}
	return DiskIOObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskIOObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (DiskIOObservable) Name() string {
	return "process.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIOObservable) Description() string {
	return "Disk bytes transferred."
}

// AttrDiskIODirection returns a required attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskIOObservable) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "process.memory.usage" semantic conventions. It represents the amount of
// physical memory in use.
type MemoryUsage struct {
	metric.Int64UpDownCounter
}

var newMemoryUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The amount of physical memory in use."),
	metric.WithUnit("By"),
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryUsageOpts
	} else {
		opt = append(opt, newMemoryUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"process.memory.usage",
		opt...,
	)
	if err != nil {
		return MemoryUsage{noop.Int64UpDownCounter{}}, err
	}
	return MemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryUsage) Name() string {
	return "process.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryUsage) Description() string {
	return "The amount of physical memory in use."
}

// Add adds incr to the existing count for attrs.
func (m MemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m MemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// MemoryUsageObservable is an instrument used to record metric values conforming
// to the "process.memory.usage" semantic conventions. It represents the amount
// of physical memory in use.
type MemoryUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The amount of physical memory in use."),
	metric.WithUnit("By"),
}

// NewMemoryUsageObservable returns a new MemoryUsageObservable instrument.
func NewMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryUsageObservableOpts
	} else {
		opt = append(opt, newMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.memory.usage",
		opt...,
	)
	if err != nil {
		return MemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryUsageObservable) Name() string {
	return "process.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryUsageObservable) Description() string {
	return "The amount of physical memory in use."
}

// MemoryVirtual is an instrument used to record metric values conforming to the
// "process.memory.virtual" semantic conventions. It represents the amount of
// committed virtual memory.
type MemoryVirtual struct {
	metric.Int64UpDownCounter
}

var newMemoryVirtualOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The amount of committed virtual memory."),
	metric.WithUnit("By"),
}

// NewMemoryVirtual returns a new MemoryVirtual instrument.
func NewMemoryVirtual(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryVirtual, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryVirtual{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryVirtualOpts
	} else {
		opt = append(opt, newMemoryVirtualOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"process.memory.virtual",
		opt...,
	)
	if err != nil {
		return MemoryVirtual{noop.Int64UpDownCounter{}}, err
	}
	return MemoryVirtual{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryVirtual) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryVirtual) Name() string {
	return "process.memory.virtual"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryVirtual) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryVirtual) Description() string {
	return "The amount of committed virtual memory."
}

// Add adds incr to the existing count for attrs.
func (m MemoryVirtual) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m MemoryVirtual) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// MemoryVirtualObservable is an instrument used to record metric values
// conforming to the "process.memory.virtual" semantic conventions. It represents
// the amount of committed virtual memory.
type MemoryVirtualObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryVirtualObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The amount of committed virtual memory."),
	metric.WithUnit("By"),
}

// NewMemoryVirtualObservable returns a new MemoryVirtualObservable instrument.
func NewMemoryVirtualObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryVirtualObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryVirtualObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryVirtualObservableOpts
	} else {
		opt = append(opt, newMemoryVirtualObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.memory.virtual",
		opt...,
	)
	if err != nil {
		return MemoryVirtualObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryVirtualObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryVirtualObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryVirtualObservable) Name() string {
	return "process.memory.virtual"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryVirtualObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryVirtualObservable) Description() string {
	return "The amount of committed virtual memory."
}

// NetworkIO is an instrument used to record metric values conforming to the
// "process.network.io" semantic conventions. It represents the network bytes
// transferred.
type NetworkIO struct {
	metric.Int64Counter
}

var newNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes transferred."),
	metric.WithUnit("By"),
}

// NewNetworkIO returns a new NetworkIO instrument.
func NewNetworkIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NetworkIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkIO{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkIOOpts
	} else {
		opt = append(opt, newNetworkIOOpts...)
	}

	i, err := m.Int64Counter(
		"process.network.io",
		opt...,
	)
	if err != nil {
		return NetworkIO{noop.Int64Counter{}}, err
	}
	return NetworkIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (NetworkIO) Name() string {
	return "process.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIO) Description() string {
	return "Network bytes transferred."
}

// Add adds incr to the existing count for attrs.
//
// The networkIoDirection is the the network IO operation direction.
func (m NetworkIO) Add(
	ctx context.Context,
	incr int64,
	networkIoDirection NetworkIODirectionAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("network.io.direction", string(networkIoDirection)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// NetworkIOObservable is an instrument used to record metric values conforming
// to the "process.network.io" semantic conventions. It represents the network
// bytes transferred.
type NetworkIOObservable struct {
	metric.Int64ObservableCounter
}

var newNetworkIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Network bytes transferred."),
	metric.WithUnit("By"),
}

// NewNetworkIOObservable returns a new NetworkIOObservable instrument.
func NewNetworkIOObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NetworkIOObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkIOObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkIOObservableOpts
	} else {
		opt = append(opt, newNetworkIOObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.network.io",
		opt...,
	)
	if err != nil {
		return NetworkIOObservable{noop.Int64ObservableCounter{}}, err
	}
	return NetworkIOObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkIOObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkIOObservable) Name() string {
	return "process.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIOObservable) Description() string {
	return "Network bytes transferred."
}

// AttrNetworkIODirection returns a required attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PagingFaults is an instrument used to record metric values conforming to the
// "process.paging.faults" semantic conventions. It represents the number of page
// faults the process has made.
type PagingFaults struct {
	metric.Int64Counter
}

var newPagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of page faults the process has made."),
	metric.WithUnit("{fault}"),
}

// NewPagingFaults returns a new PagingFaults instrument.
func NewPagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPagingFaultsOpts
	} else {
		opt = append(opt, newPagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"process.paging.faults",
		opt...,
	)
	if err != nil {
		return PagingFaults{noop.Int64Counter{}}, err
	}
	return PagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (PagingFaults) Name() string {
	return "process.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (PagingFaults) Description() string {
	return "Number of page faults the process has made."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m PagingFaults) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the type of
// paging fault. Value MUST be either `major` or `minor`. If the metric is
// reported without this attribute, it should be the sum of major and minor page
// faults.
func (PagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// PagingFaultsObservable is an instrument used to record metric values
// conforming to the "process.paging.faults" semantic conventions. It represents
// the number of page faults the process has made.
type PagingFaultsObservable struct {
	metric.Int64ObservableCounter
}

var newPagingFaultsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Number of page faults the process has made."),
	metric.WithUnit("{fault}"),
}

// NewPagingFaultsObservable returns a new PagingFaultsObservable instrument.
func NewPagingFaultsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (PagingFaultsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingFaultsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPagingFaultsObservableOpts
	} else {
		opt = append(opt, newPagingFaultsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.paging.faults",
		opt...,
	)
	if err != nil {
		return PagingFaultsObservable{noop.Int64ObservableCounter{}}, err
	}
	return PagingFaultsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingFaultsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (PagingFaultsObservable) Name() string {
	return "process.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PagingFaultsObservable) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (PagingFaultsObservable) Description() string {
	return "Number of page faults the process has made."
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the type of
// paging fault. Value MUST be either `major` or `minor`. If the metric is
// reported without this attribute, it should be the sum of major and minor page
// faults.
func (PagingFaultsObservable) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// ThreadCount is an instrument used to record metric values conforming to the
// "process.thread.count" semantic conventions. It represents the process threads
// count.
type ThreadCount struct {
	metric.Int64UpDownCounter
}

var newThreadCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Process threads count."),
	metric.WithUnit("{thread}"),
}

// NewThreadCount returns a new ThreadCount instrument.
func NewThreadCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ThreadCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ThreadCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newThreadCountOpts
	} else {
		opt = append(opt, newThreadCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"process.thread.count",
		opt...,
	)
	if err != nil {
		return ThreadCount{noop.Int64UpDownCounter{}}, err
	}
	return ThreadCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ThreadCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ThreadCount) Name() string {
	return "process.thread.count"
}

// Unit returns the semantic convention unit of the instrument
func (ThreadCount) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (ThreadCount) Description() string {
	return "Process threads count."
}

// Add adds incr to the existing count for attrs.
func (m ThreadCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ThreadCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ThreadCountObservable is an instrument used to record metric values conforming
// to the "process.thread.count" semantic conventions. It represents the process
// threads count.
type ThreadCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newThreadCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Process threads count."),
	metric.WithUnit("{thread}"),
}

// NewThreadCountObservable returns a new ThreadCountObservable instrument.
func NewThreadCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ThreadCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ThreadCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newThreadCountObservableOpts
	} else {
		opt = append(opt, newThreadCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.thread.count",
		opt...,
	)
	if err != nil {
		return ThreadCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ThreadCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ThreadCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ThreadCountObservable) Name() string {
	return "process.thread.count"
}

// Unit returns the semantic convention unit of the instrument
func (ThreadCountObservable) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (ThreadCountObservable) Description() string {
	return "Process threads count."
}

// UnixFileDescriptorCount is an instrument used to record metric values
// conforming to the "process.unix.file_descriptor.count" semantic conventions.
// It represents the number of unix file descriptors in use by the process.
type UnixFileDescriptorCount struct {
	metric.Int64UpDownCounter
}

var newUnixFileDescriptorCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of unix file descriptors in use by the process."),
	metric.WithUnit("{file_descriptor}"),
}

// NewUnixFileDescriptorCount returns a new UnixFileDescriptorCount instrument.
func NewUnixFileDescriptorCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (UnixFileDescriptorCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return UnixFileDescriptorCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newUnixFileDescriptorCountOpts
	} else {
		opt = append(opt, newUnixFileDescriptorCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"process.unix.file_descriptor.count",
		opt...,
	)
	if err != nil {
		return UnixFileDescriptorCount{noop.Int64UpDownCounter{}}, err
	}
	return UnixFileDescriptorCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m UnixFileDescriptorCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (UnixFileDescriptorCount) Name() string {
	return "process.unix.file_descriptor.count"
}

// Unit returns the semantic convention unit of the instrument
func (UnixFileDescriptorCount) Unit() string {
	return "{file_descriptor}"
}

// Description returns the semantic convention description of the instrument
func (UnixFileDescriptorCount) Description() string {
	return "Number of unix file descriptors in use by the process."
}

// Add adds incr to the existing count for attrs.
func (m UnixFileDescriptorCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m UnixFileDescriptorCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// UnixFileDescriptorCountObservable is an instrument used to record metric
// values conforming to the "process.unix.file_descriptor.count" semantic
// conventions. It represents the number of unix file descriptors in use by the
// process.
type UnixFileDescriptorCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newUnixFileDescriptorCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of unix file descriptors in use by the process."),
	metric.WithUnit("{file_descriptor}"),
}

// NewUnixFileDescriptorCountObservable returns a new
// UnixFileDescriptorCountObservable instrument.
func NewUnixFileDescriptorCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (UnixFileDescriptorCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return UnixFileDescriptorCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newUnixFileDescriptorCountObservableOpts
	} else {
		opt = append(opt, newUnixFileDescriptorCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.unix.file_descriptor.count",
		opt...,
	)
	if err != nil {
		return UnixFileDescriptorCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return UnixFileDescriptorCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m UnixFileDescriptorCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (UnixFileDescriptorCountObservable) Name() string {
	return "process.unix.file_descriptor.count"
}

// Unit returns the semantic convention unit of the instrument
func (UnixFileDescriptorCountObservable) Unit() string {
	return "{file_descriptor}"
}

// Description returns the semantic convention description of the instrument
func (UnixFileDescriptorCountObservable) Description() string {
	return "Number of unix file descriptors in use by the process."
}

// Uptime is an instrument used to record metric values conforming to the
// "process.uptime" semantic conventions. It represents the time the process has
// been running.
type Uptime struct {
	metric.Float64Gauge
}

var newUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the process has been running."),
	metric.WithUnit("s"),
}

// NewUptime returns a new Uptime instrument.
func NewUptime(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (Uptime, error) {
	// Check if the meter is nil.
	if m == nil {
		return Uptime{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newUptimeOpts
	} else {
		opt = append(opt, newUptimeOpts...)
	}

	i, err := m.Float64Gauge(
		"process.uptime",
		opt...,
	)
	if err != nil {
		return Uptime{noop.Float64Gauge{}}, err
	}
	return Uptime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Uptime) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (Uptime) Name() string {
	return "process.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (Uptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (Uptime) Description() string {
	return "The time the process has been running."
}

// Record records val to the current distribution for attrs.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m Uptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m Uptime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// UptimeObservable is an instrument used to record metric values conforming to
// the "process.uptime" semantic conventions. It represents the time the process
// has been running.
type UptimeObservable struct {
	metric.Float64ObservableGauge
}

var newUptimeObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The time the process has been running."),
	metric.WithUnit("s"),
}

// NewUptimeObservable returns a new UptimeObservable instrument.
func NewUptimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (UptimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return UptimeObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newUptimeObservableOpts
	} else {
		opt = append(opt, newUptimeObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"process.uptime",
		opt...,
	)
	if err != nil {
		return UptimeObservable{noop.Float64ObservableGauge{}}, err
	}
	return UptimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m UptimeObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (UptimeObservable) Name() string {
	return "process.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (UptimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (UptimeObservable) Description() string {
	return "The time the process has been running."
}

// WindowsHandleCount is an instrument used to record metric values conforming to
// the "process.windows.handle.count" semantic conventions. It represents the
// number of handles held by the process.
type WindowsHandleCount struct {
	metric.Int64UpDownCounter
}

var newWindowsHandleCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of handles held by the process."),
	metric.WithUnit("{handle}"),
}

// NewWindowsHandleCount returns a new WindowsHandleCount instrument.
func NewWindowsHandleCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (WindowsHandleCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return WindowsHandleCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newWindowsHandleCountOpts
	} else {
		opt = append(opt, newWindowsHandleCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"process.windows.handle.count",
		opt...,
	)
	if err != nil {
		return WindowsHandleCount{noop.Int64UpDownCounter{}}, err
	}
	return WindowsHandleCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m WindowsHandleCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (WindowsHandleCount) Name() string {
	return "process.windows.handle.count"
}

// Unit returns the semantic convention unit of the instrument
func (WindowsHandleCount) Unit() string {
	return "{handle}"
}

// Description returns the semantic convention description of the instrument
func (WindowsHandleCount) Description() string {
	return "Number of handles held by the process."
}

// Add adds incr to the existing count for attrs.
func (m WindowsHandleCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m WindowsHandleCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// WindowsHandleCountObservable is an instrument used to record metric values
// conforming to the "process.windows.handle.count" semantic conventions. It
// represents the number of handles held by the process.
type WindowsHandleCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newWindowsHandleCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of handles held by the process."),
	metric.WithUnit("{handle}"),
}

// NewWindowsHandleCountObservable returns a new WindowsHandleCountObservable
// instrument.
func NewWindowsHandleCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (WindowsHandleCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return WindowsHandleCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newWindowsHandleCountObservableOpts
	} else {
		opt = append(opt, newWindowsHandleCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.windows.handle.count",
		opt...,
	)
	if err != nil {
		return WindowsHandleCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return WindowsHandleCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m WindowsHandleCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (WindowsHandleCountObservable) Name() string {
	return "process.windows.handle.count"
}

// Unit returns the semantic convention unit of the instrument
func (WindowsHandleCountObservable) Unit() string {
	return "{handle}"
}

// Description returns the semantic convention description of the instrument
func (WindowsHandleCountObservable) Description() string {
	return "Number of handles held by the process."
}
