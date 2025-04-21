// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/process"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// CPUModeAttr is an attribute conforming to the cpu.mode semantic conventions.
// It represents a process SHOULD be characterized *either* by data points with
// no `mode` labels, *or only* data points with `mode` labels.
type CPUModeAttr string

var (
	// CPUModeUser is the none.
	CPUModeUser CPUModeAttr = "user"
	// CPUModeSystem is the none.
	CPUModeSystem CPUModeAttr = "system"
	// CPUModeNice is the none.
	CPUModeNice CPUModeAttr = "nice"
	// CPUModeIdle is the none.
	CPUModeIdle CPUModeAttr = "idle"
	// CPUModeIowait is the none.
	CPUModeIowait CPUModeAttr = "iowait"
	// CPUModeInterrupt is the none.
	CPUModeInterrupt CPUModeAttr = "interrupt"
	// CPUModeSteal is the none.
	CPUModeSteal CPUModeAttr = "steal"
	// CPUModeKernel is the none.
	CPUModeKernel CPUModeAttr = "kernel"
)

// DiskIoDirectionAttr is an attribute conforming to the disk.io.direction
// semantic conventions. It represents the disk IO operation direction.
type DiskIoDirectionAttr string

var (
	// DiskIoDirectionRead is the none.
	DiskIoDirectionRead DiskIoDirectionAttr = "read"
	// DiskIoDirectionWrite is the none.
	DiskIoDirectionWrite DiskIoDirectionAttr = "write"
)

// NetworkIoDirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIoDirectionAttr string

var (
	// NetworkIoDirectionTransmit is the none.
	NetworkIoDirectionTransmit NetworkIoDirectionAttr = "transmit"
	// NetworkIoDirectionReceive is the none.
	NetworkIoDirectionReceive NetworkIoDirectionAttr = "receive"
)

// ContextSwitchTypeAttr is an attribute conforming to the
// process.context_switch_type semantic conventions. It represents the specifies
// whether the context switches for this data point were voluntary or
// involuntary.
type ContextSwitchTypeAttr string

var (
	// ContextSwitchTypeVoluntary is the none.
	ContextSwitchTypeVoluntary ContextSwitchTypeAttr = "voluntary"
	// ContextSwitchTypeInvoluntary is the none.
	ContextSwitchTypeInvoluntary ContextSwitchTypeAttr = "involuntary"
)

// PagingFaultTypeAttr is an attribute conforming to the
// process.paging.fault_type semantic conventions. It represents the type of page
// fault for this data point. Type `major` is for major/hard page faults, and
// `minor` is for minor/soft page faults.
type PagingFaultTypeAttr string

var (
	// PagingFaultTypeMajor is the none.
	PagingFaultTypeMajor PagingFaultTypeAttr = "major"
	// PagingFaultTypeMinor is the none.
	PagingFaultTypeMinor PagingFaultTypeAttr = "minor"
)

// ProcessContextSwitches is an instrument used to record metric values
// conforming to the "process.context_switches" semantic conventions. It
// represents the number of times the process has been context switched.
type ContextSwitches struct {
	inst metric.Int64Counter
}

// NewContextSwitches returns a new ContextSwitches instrument.
func NewContextSwitches(m metric.Meter) (ContextSwitches, error) {
	i, err := m.Int64Counter(
	    "process.context_switches",
	    metric.WithDescription("Number of times the process has been context switched."),
	    metric.WithUnit("{context_switch}"),
	)
	if err != nil {
	    return ContextSwitches{}, err
	}
	return ContextSwitches{i}, nil
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

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ContextSwitches) Add(
    ctx context.Context,
    incr int64,
	attrs ...ContextSwitchesAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ContextSwitches) conv(in []ContextSwitchesAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.contextSwitchesAttr()
	}
	return out
}

// ContextSwitchesAttr is an optional attribute for the ContextSwitches
// instrument.
type ContextSwitchesAttr interface {
    contextSwitchesAttr() attribute.KeyValue
}

type contextSwitchesAttr struct {
	kv attribute.KeyValue
}

func (a contextSwitchesAttr) contextSwitchesAttr() attribute.KeyValue {
    return a.kv
}

// ContextSwitchType returns an optional attribute for the
// "process.context_switch_type" semantic convention. It represents the specifies
// whether the context switches for this data point were voluntary or
// involuntary.
func (ContextSwitches) ContextSwitchType(val ContextSwitchTypeAttr) ContextSwitchesAttr {
	return contextSwitchesAttr{kv: attribute.String("process.context_switch_type", string(val))}
}

// ProcessCPUTime is an instrument used to record metric values conforming to the
// "process.cpu.time" semantic conventions. It represents the total CPU seconds
// broken down by different states.
type CPUTime struct {
	inst metric.Float64Counter
}

// NewCPUTime returns a new CPUTime instrument.
func NewCPUTime(m metric.Meter) (CPUTime, error) {
	i, err := m.Float64Counter(
	    "process.cpu.time",
	    metric.WithDescription("Total CPU seconds broken down by different states."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return CPUTime{}, err
	}
	return CPUTime{i}, nil
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

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m CPUTime) Add(
    ctx context.Context,
    incr float64,
	attrs ...CPUTimeAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CPUTime) conv(in []CPUTimeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cpuTimeAttr()
	}
	return out
}

// CPUTimeAttr is an optional attribute for the CPUTime instrument.
type CPUTimeAttr interface {
    cpuTimeAttr() attribute.KeyValue
}

type cpuTimeAttr struct {
	kv attribute.KeyValue
}

func (a cpuTimeAttr) cpuTimeAttr() attribute.KeyValue {
    return a.kv
}

// CPUMode returns an optional attribute for the "cpu.mode" semantic convention.
// It represents a process SHOULD be characterized *either* by data points with
// no `mode` labels, *or only* data points with `mode` labels.
func (CPUTime) CPUMode(val CPUModeAttr) CPUTimeAttr {
	return cpuTimeAttr{kv: attribute.String("cpu.mode", string(val))}
}

// ProcessCPUUtilization is an instrument used to record metric values conforming
// to the "process.cpu.utilization" semantic conventions. It represents the
// difference in process.cpu.time since the last measurement, divided by the
// elapsed time and number of CPUs available to the process.
type CPUUtilization struct {
	inst metric.Int64Gauge
}

// NewCPUUtilization returns a new CPUUtilization instrument.
func NewCPUUtilization(m metric.Meter) (CPUUtilization, error) {
	i, err := m.Int64Gauge(
	    "process.cpu.utilization",
	    metric.WithDescription("Difference in process.cpu.time since the last measurement, divided by the elapsed time and number of CPUs available to the process."),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return CPUUtilization{}, err
	}
	return CPUUtilization{i}, nil
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

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m CPUUtilization) Record(
    ctx context.Context,
    val int64,
	attrs ...CPUUtilizationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CPUUtilization) conv(in []CPUUtilizationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cpuUtilizationAttr()
	}
	return out
}

// CPUUtilizationAttr is an optional attribute for the CPUUtilization instrument.
type CPUUtilizationAttr interface {
    cpuUtilizationAttr() attribute.KeyValue
}

type cpuUtilizationAttr struct {
	kv attribute.KeyValue
}

func (a cpuUtilizationAttr) cpuUtilizationAttr() attribute.KeyValue {
    return a.kv
}

// CPUMode returns an optional attribute for the "cpu.mode" semantic convention.
// It represents a process SHOULD be characterized *either* by data points with
// no `mode` labels, *or only* data points with `mode` labels.
func (CPUUtilization) CPUMode(val CPUModeAttr) CPUUtilizationAttr {
	return cpuUtilizationAttr{kv: attribute.String("cpu.mode", string(val))}
}

// ProcessDiskIo is an instrument used to record metric values conforming to the
// "process.disk.io" semantic conventions. It represents the disk bytes
// transferred.
type DiskIo struct {
	inst metric.Int64Counter
}

// NewDiskIo returns a new DiskIo instrument.
func NewDiskIo(m metric.Meter) (DiskIo, error) {
	i, err := m.Int64Counter(
	    "process.disk.io",
	    metric.WithDescription("Disk bytes transferred."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return DiskIo{}, err
	}
	return DiskIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskIo) Name() string {
	return "process.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIo) Description() string {
	return "Disk bytes transferred."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskIo) Add(
    ctx context.Context,
    incr int64,
	attrs ...DiskIoAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskIo) conv(in []DiskIoAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskIoAttr()
	}
	return out
}

// DiskIoAttr is an optional attribute for the DiskIo instrument.
type DiskIoAttr interface {
    diskIoAttr() attribute.KeyValue
}

type diskIoAttr struct {
	kv attribute.KeyValue
}

func (a diskIoAttr) diskIoAttr() attribute.KeyValue {
    return a.kv
}

// DiskIoDirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskIo) DiskIoDirection(val DiskIoDirectionAttr) DiskIoAttr {
	return diskIoAttr{kv: attribute.String("disk.io.direction", string(val))}
}

// ProcessMemoryUsage is an instrument used to record metric values conforming to
// the "process.memory.usage" semantic conventions. It represents the amount of
// physical memory in use.
type MemoryUsage struct {
	inst metric.Int64UpDownCounter
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(m metric.Meter) (MemoryUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "process.memory.usage",
	    metric.WithDescription("The amount of physical memory in use."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryUsage{}, err
	}
	return MemoryUsage{i}, nil
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

func (m MemoryUsage) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// ProcessMemoryVirtual is an instrument used to record metric values conforming
// to the "process.memory.virtual" semantic conventions. It represents the amount
// of committed virtual memory.
type MemoryVirtual struct {
	inst metric.Int64UpDownCounter
}

// NewMemoryVirtual returns a new MemoryVirtual instrument.
func NewMemoryVirtual(m metric.Meter) (MemoryVirtual, error) {
	i, err := m.Int64UpDownCounter(
	    "process.memory.virtual",
	    metric.WithDescription("The amount of committed virtual memory."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryVirtual{}, err
	}
	return MemoryVirtual{i}, nil
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

func (m MemoryVirtual) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// ProcessNetworkIo is an instrument used to record metric values conforming to
// the "process.network.io" semantic conventions. It represents the network bytes
// transferred.
type NetworkIo struct {
	inst metric.Int64Counter
}

// NewNetworkIo returns a new NetworkIo instrument.
func NewNetworkIo(m metric.Meter) (NetworkIo, error) {
	i, err := m.Int64Counter(
	    "process.network.io",
	    metric.WithDescription("Network bytes transferred."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NetworkIo{}, err
	}
	return NetworkIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkIo) Name() string {
	return "process.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIo) Description() string {
	return "Network bytes transferred."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkIo) Add(
    ctx context.Context,
    incr int64,
	attrs ...NetworkIoAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetworkIo) conv(in []NetworkIoAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.networkIoAttr()
	}
	return out
}

// NetworkIoAttr is an optional attribute for the NetworkIo instrument.
type NetworkIoAttr interface {
    networkIoAttr() attribute.KeyValue
}

type networkIoAttr struct {
	kv attribute.KeyValue
}

func (a networkIoAttr) networkIoAttr() attribute.KeyValue {
    return a.kv
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIo) NetworkIoDirection(val NetworkIoDirectionAttr) NetworkIoAttr {
	return networkIoAttr{kv: attribute.String("network.io.direction", string(val))}
}

// ProcessOpenFileDescriptorCount is an instrument used to record metric values
// conforming to the "process.open_file_descriptor.count" semantic conventions.
// It represents the number of file descriptors in use by the process.
type OpenFileDescriptorCount struct {
	inst metric.Int64UpDownCounter
}

// NewOpenFileDescriptorCount returns a new OpenFileDescriptorCount instrument.
func NewOpenFileDescriptorCount(m metric.Meter) (OpenFileDescriptorCount, error) {
	i, err := m.Int64UpDownCounter(
	    "process.open_file_descriptor.count",
	    metric.WithDescription("Number of file descriptors in use by the process."),
	    metric.WithUnit("{file_descriptor}"),
	)
	if err != nil {
	    return OpenFileDescriptorCount{}, err
	}
	return OpenFileDescriptorCount{i}, nil
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

func (m OpenFileDescriptorCount) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// ProcessPagingFaults is an instrument used to record metric values conforming
// to the "process.paging.faults" semantic conventions. It represents the number
// of page faults the process has made.
type PagingFaults struct {
	inst metric.Int64Counter
}

// NewPagingFaults returns a new PagingFaults instrument.
func NewPagingFaults(m metric.Meter) (PagingFaults, error) {
	i, err := m.Int64Counter(
	    "process.paging.faults",
	    metric.WithDescription("Number of page faults the process has made."),
	    metric.WithUnit("{fault}"),
	)
	if err != nil {
	    return PagingFaults{}, err
	}
	return PagingFaults{i}, nil
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

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PagingFaults) Add(
    ctx context.Context,
    incr int64,
	attrs ...PagingFaultsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PagingFaults) conv(in []PagingFaultsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.pagingFaultsAttr()
	}
	return out
}

// PagingFaultsAttr is an optional attribute for the PagingFaults instrument.
type PagingFaultsAttr interface {
    pagingFaultsAttr() attribute.KeyValue
}

type pagingFaultsAttr struct {
	kv attribute.KeyValue
}

func (a pagingFaultsAttr) pagingFaultsAttr() attribute.KeyValue {
    return a.kv
}

// PagingFaultType returns an optional attribute for the
// "process.paging.fault_type" semantic convention. It represents the type of
// page fault for this data point. Type `major` is for major/hard page faults,
// and `minor` is for minor/soft page faults.
func (PagingFaults) PagingFaultType(val PagingFaultTypeAttr) PagingFaultsAttr {
	return pagingFaultsAttr{kv: attribute.String("process.paging.fault_type", string(val))}
}

// ProcessThreadCount is an instrument used to record metric values conforming to
// the "process.thread.count" semantic conventions. It represents the process
// threads count.
type ThreadCount struct {
	inst metric.Int64UpDownCounter
}

// NewThreadCount returns a new ThreadCount instrument.
func NewThreadCount(m metric.Meter) (ThreadCount, error) {
	i, err := m.Int64UpDownCounter(
	    "process.thread.count",
	    metric.WithDescription("Process threads count."),
	    metric.WithUnit("{thread}"),
	)
	if err != nil {
	    return ThreadCount{}, err
	}
	return ThreadCount{i}, nil
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

func (m ThreadCount) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// ProcessUptime is an instrument used to record metric values conforming to the
// "process.uptime" semantic conventions. It represents the time the process has
// been running.
type Uptime struct {
	inst metric.Float64Gauge
}

// NewUptime returns a new Uptime instrument.
func NewUptime(m metric.Meter) (Uptime, error) {
	i, err := m.Float64Gauge(
	    "process.uptime",
	    metric.WithDescription("The time the process has been running."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return Uptime{}, err
	}
	return Uptime{i}, nil
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

func (m Uptime) Record(ctx context.Context, val float64) {
    m.inst.Record(ctx, val)
}