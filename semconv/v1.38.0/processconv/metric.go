// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package processconv provides types and functionality for OpenTelemetry semantic
// conventions in the "process" namespace.
package processconv

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

// CPUModeAttr is an attribute conforming to the cpu.mode semantic conventions.
// It represents a process SHOULD be characterized *either* by data points with
// no `mode` labels, *or only* data points with `mode` labels.
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
	// CPUModeKernel is the kernel.
	CPUModeKernel CPUModeAttr = "kernel"
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
// system.paging.fault.type semantic conventions. It represents the paging fault
// type.
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
	metric.Int64ObservableCounter
}

var newContextSwitchesOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Number of times the process has been context switched."),
	metric.WithUnit("{context_switch}"),
}

// NewContextSwitches returns a new ContextSwitches instrument.
func NewContextSwitches(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (ContextSwitches, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContextSwitches{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContextSwitchesOpts
	} else {
		opt = append(opt, newContextSwitchesOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.context_switches",
		opt...,
	)
	if err != nil {
		return ContextSwitches{noop.Int64ObservableCounter{}}, err
	}
	return ContextSwitches{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContextSwitches) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
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

// AttrContextSwitchType returns an optional attribute for the
// "process.context_switch.type" semantic convention. It represents the specifies
// whether the context switches for this data point were voluntary or
// involuntary.
func (ContextSwitches) AttrContextSwitchType(val ContextSwitchTypeAttr) attribute.KeyValue {
	return attribute.String("process.context_switch.type", string(val))
}

// CPUTime is an instrument used to record metric values conforming to the
// "process.cpu.time" semantic conventions. It represents the total CPU seconds
// broken down by different states.
type CPUTime struct {
	metric.Float64ObservableCounter
}

var newCPUTimeOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Total CPU seconds broken down by different states."),
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
	return "Total CPU seconds broken down by different states."
}

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents a process SHOULD be characterized *either* by data
// points with no `mode` labels, *or only* data points with `mode` labels.
func (CPUTime) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUUtilization is an instrument used to record metric values conforming to the
// "process.cpu.utilization" semantic conventions. It represents the difference
// in process.cpu.time since the last measurement, divided by the elapsed time
// and number of CPUs available to the process.
type CPUUtilization struct {
	metric.Int64ObservableGauge
}

var newCPUUtilizationOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."),
	metric.WithUnit("1"),
}

// NewCPUUtilization returns a new CPUUtilization instrument.
func NewCPUUtilization(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (CPUUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUtilization{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUUtilizationOpts
	} else {
		opt = append(opt, newCPUUtilizationOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"process.cpu.utilization",
		opt...,
	)
	if err != nil {
		return CPUUtilization{noop.Int64ObservableGauge{}}, err
	}
	return CPUUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUtilization) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
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

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents a process SHOULD be characterized *either* by data
// points with no `mode` labels, *or only* data points with `mode` labels.
func (CPUUtilization) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// DiskIO is an instrument used to record metric values conforming to the
// "process.disk.io" semantic conventions. It represents the disk bytes
// transferred.
type DiskIO struct {
	metric.Int64ObservableCounter
}

var newDiskIOOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Disk bytes transferred."),
	metric.WithUnit("By"),
}

// NewDiskIO returns a new DiskIO instrument.
func NewDiskIO(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (DiskIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskIO{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDiskIOOpts
	} else {
		opt = append(opt, newDiskIOOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.disk.io",
		opt...,
	)
	if err != nil {
		return DiskIO{noop.Int64ObservableCounter{}}, err
	}
	return DiskIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskIO) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
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

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskIO) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "process.memory.usage" semantic conventions. It represents the amount of
// physical memory in use.
type MemoryUsage struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryUsageOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The amount of physical memory in use."),
	metric.WithUnit("By"),
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUsage{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryUsageOpts
	} else {
		opt = append(opt, newMemoryUsageOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.memory.usage",
		opt...,
	)
	if err != nil {
		return MemoryUsage{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUsage) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
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

// MemoryVirtual is an instrument used to record metric values conforming to the
// "process.memory.virtual" semantic conventions. It represents the amount of
// committed virtual memory.
type MemoryVirtual struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryVirtualOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The amount of committed virtual memory."),
	metric.WithUnit("By"),
}

// NewMemoryVirtual returns a new MemoryVirtual instrument.
func NewMemoryVirtual(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryVirtual, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryVirtual{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryVirtualOpts
	} else {
		opt = append(opt, newMemoryVirtualOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.memory.virtual",
		opt...,
	)
	if err != nil {
		return MemoryVirtual{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryVirtual{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryVirtual) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
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

// NetworkIO is an instrument used to record metric values conforming to the
// "process.network.io" semantic conventions. It represents the network bytes
// transferred.
type NetworkIO struct {
	metric.Int64ObservableCounter
}

var newNetworkIOOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Network bytes transferred."),
	metric.WithUnit("By"),
}

// NewNetworkIO returns a new NetworkIO instrument.
func NewNetworkIO(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NetworkIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkIO{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkIOOpts
	} else {
		opt = append(opt, newNetworkIOOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.network.io",
		opt...,
	)
	if err != nil {
		return NetworkIO{noop.Int64ObservableCounter{}}, err
	}
	return NetworkIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkIO) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
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

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// OpenFileDescriptorCount is an instrument used to record metric values
// conforming to the "process.open_file_descriptor.count" semantic conventions.
// It represents the number of file descriptors in use by the process.
type OpenFileDescriptorCount struct {
	metric.Int64ObservableUpDownCounter
}

var newOpenFileDescriptorCountOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of file descriptors in use by the process."),
	metric.WithUnit("{file_descriptor}"),
}

// NewOpenFileDescriptorCount returns a new OpenFileDescriptorCount instrument.
func NewOpenFileDescriptorCount(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (OpenFileDescriptorCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return OpenFileDescriptorCount{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newOpenFileDescriptorCountOpts
	} else {
		opt = append(opt, newOpenFileDescriptorCountOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.open_file_descriptor.count",
		opt...,
	)
	if err != nil {
		return OpenFileDescriptorCount{noop.Int64ObservableUpDownCounter{}}, err
	}
	return OpenFileDescriptorCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m OpenFileDescriptorCount) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (OpenFileDescriptorCount) Name() string {
	return "process.open_file_descriptor.count"
}

// Unit returns the semantic convention unit of the instrument
func (OpenFileDescriptorCount) Unit() string {
	return "{file_descriptor}"
}

// Description returns the semantic convention description of the instrument
func (OpenFileDescriptorCount) Description() string {
	return "Number of file descriptors in use by the process."
}

// PagingFaults is an instrument used to record metric values conforming to the
// "process.paging.faults" semantic conventions. It represents the number of page
// faults the process has made.
type PagingFaults struct {
	metric.Int64ObservableCounter
}

var newPagingFaultsOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Number of page faults the process has made."),
	metric.WithUnit("{fault}"),
}

// NewPagingFaults returns a new PagingFaults instrument.
func NewPagingFaults(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (PagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingFaults{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPagingFaultsOpts
	} else {
		opt = append(opt, newPagingFaultsOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"process.paging.faults",
		opt...,
	)
	if err != nil {
		return PagingFaults{noop.Int64ObservableCounter{}}, err
	}
	return PagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingFaults) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
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

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (PagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// ThreadCount is an instrument used to record metric values conforming to the
// "process.thread.count" semantic conventions. It represents the process threads
// count.
type ThreadCount struct {
	metric.Int64ObservableUpDownCounter
}

var newThreadCountOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Process threads count."),
	metric.WithUnit("{thread}"),
}

// NewThreadCount returns a new ThreadCount instrument.
func NewThreadCount(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ThreadCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ThreadCount{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newThreadCountOpts
	} else {
		opt = append(opt, newThreadCountOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"process.thread.count",
		opt...,
	)
	if err != nil {
		return ThreadCount{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ThreadCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ThreadCount) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
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

// Uptime is an instrument used to record metric values conforming to the
// "process.uptime" semantic conventions. It represents the time the process has
// been running.
type Uptime struct {
	metric.Float64ObservableGauge
}

var newUptimeOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The time the process has been running."),
	metric.WithUnit("s"),
}

// NewUptime returns a new Uptime instrument.
func NewUptime(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (Uptime, error) {
	// Check if the meter is nil.
	if m == nil {
		return Uptime{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newUptimeOpts
	} else {
		opt = append(opt, newUptimeOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"process.uptime",
		opt...,
	)
	if err != nil {
		return Uptime{noop.Float64ObservableGauge{}}, err
	}
	return Uptime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Uptime) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
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
