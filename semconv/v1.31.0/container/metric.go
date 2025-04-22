// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/container"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// CPUModeAttr is an attribute conforming to the cpu.mode semantic conventions.
// It represents the CPU mode for this data point. A container's CPU metric
// SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
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

// CPUTime is an instrument used to record metric values conforming to the
// "container.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type CPUTime struct {
	inst metric.Float64Counter
}

// NewCPUTime returns a new CPUTime instrument.
func NewCPUTime(m metric.Meter) (CPUTime, error) {
	i, err := m.Float64Counter(
	    "container.cpu.time",
	    metric.WithDescription("Total CPU time consumed"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return CPUTime{}, err
	}
	return CPUTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUTime) Name() string {
	return "container.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (CPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (CPUTime) Description() string {
	return "Total CPU time consumed"
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
// It represents the CPU mode for this data point. A container's CPU metric
// SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUTime) CPUModeAttr(val CPUModeAttr) CPUTimeAttr {
	return cpuTimeAttr{kv: attribute.String("cpu.mode", string(val))}
}

// CPUUsage is an instrument used to record metric values conforming to the
// "container.cpu.usage" semantic conventions. It represents the container's CPU
// usage, measured in cpus. Range from 0 to the number of allocatable CPUs.
type CPUUsage struct {
	inst metric.Int64Gauge
}

// NewCPUUsage returns a new CPUUsage instrument.
func NewCPUUsage(m metric.Meter) (CPUUsage, error) {
	i, err := m.Int64Gauge(
	    "container.cpu.usage",
	    metric.WithDescription("Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
	    metric.WithUnit("{cpu}"),
	)
	if err != nil {
	    return CPUUsage{}, err
	}
	return CPUUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUUsage) Name() string {
	return "container.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUsage) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (CPUUsage) Description() string {
	return "Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m CPUUsage) Record(
    ctx context.Context,
    val int64,
	attrs ...CPUUsageAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CPUUsage) conv(in []CPUUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cpuUsageAttr()
	}
	return out
}

// CPUUsageAttr is an optional attribute for the CPUUsage instrument.
type CPUUsageAttr interface {
    cpuUsageAttr() attribute.KeyValue
}

type cpuUsageAttr struct {
	kv attribute.KeyValue
}

func (a cpuUsageAttr) cpuUsageAttr() attribute.KeyValue {
    return a.kv
}

// CPUMode returns an optional attribute for the "cpu.mode" semantic convention.
// It represents the CPU mode for this data point. A container's CPU metric
// SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUUsage) CPUModeAttr(val CPUModeAttr) CPUUsageAttr {
	return cpuUsageAttr{kv: attribute.String("cpu.mode", string(val))}
}

// DiskIo is an instrument used to record metric values conforming to the
// "container.disk.io" semantic conventions. It represents the disk bytes for the
// container.
type DiskIo struct {
	inst metric.Int64Counter
}

// NewDiskIo returns a new DiskIo instrument.
func NewDiskIo(m metric.Meter) (DiskIo, error) {
	i, err := m.Int64Counter(
	    "container.disk.io",
	    metric.WithDescription("Disk bytes for the container."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return DiskIo{}, err
	}
	return DiskIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskIo) Name() string {
	return "container.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIo) Description() string {
	return "Disk bytes for the container."
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
func (DiskIo) DiskIoDirectionAttr(val DiskIoDirectionAttr) DiskIoAttr {
	return diskIoAttr{kv: attribute.String("disk.io.direction", string(val))}
}

// SystemDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskIo) SystemDeviceAttr(val string) DiskIoAttr {
	return diskIoAttr{kv: attribute.String("system.device", val)}
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "container.memory.usage" semantic conventions. It represents the memory usage
// of the container.
type MemoryUsage struct {
	inst metric.Int64Counter
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(m metric.Meter) (MemoryUsage, error) {
	i, err := m.Int64Counter(
	    "container.memory.usage",
	    metric.WithDescription("Memory usage of the container."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryUsage{}, err
	}
	return MemoryUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (MemoryUsage) Name() string {
	return "container.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryUsage) Description() string {
	return "Memory usage of the container."
}

func (m MemoryUsage) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// NetworkIo is an instrument used to record metric values conforming to the
// "container.network.io" semantic conventions. It represents the network bytes
// for the container.
type NetworkIo struct {
	inst metric.Int64Counter
}

// NewNetworkIo returns a new NetworkIo instrument.
func NewNetworkIo(m metric.Meter) (NetworkIo, error) {
	i, err := m.Int64Counter(
	    "container.network.io",
	    metric.WithDescription("Network bytes for the container."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NetworkIo{}, err
	}
	return NetworkIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkIo) Name() string {
	return "container.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIo) Description() string {
	return "Network bytes for the container."
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

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkIo) NetworkInterfaceNameAttr(val string) NetworkIoAttr {
	return networkIoAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIo) NetworkIoDirectionAttr(val NetworkIoDirectionAttr) NetworkIoAttr {
	return networkIoAttr{kv: attribute.String("network.io.direction", string(val))}
}

// Uptime is an instrument used to record metric values conforming to the
// "container.uptime" semantic conventions. It represents the time the container
// has been running.
type Uptime struct {
	inst metric.Float64Gauge
}

// NewUptime returns a new Uptime instrument.
func NewUptime(m metric.Meter) (Uptime, error) {
	i, err := m.Float64Gauge(
	    "container.uptime",
	    metric.WithDescription("The time the container has been running"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return Uptime{}, err
	}
	return Uptime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (Uptime) Name() string {
	return "container.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (Uptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (Uptime) Description() string {
	return "The time the container has been running"
}

func (m Uptime) Record(ctx context.Context, val float64) {
    m.inst.Record(ctx, val)
}