// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/system"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

// LinuxMemorySlabStateAttr is an attribute conforming to the
// linux.memory.slab.state semantic conventions. It represents the Linux Slab
// memory state.
type LinuxMemorySlabStateAttr string

var (
	// LinuxMemorySlabStateReclaimable is the none.
	LinuxMemorySlabStateReclaimable LinuxMemorySlabStateAttr = "reclaimable"
	// LinuxMemorySlabStateUnreclaimable is the none.
	LinuxMemorySlabStateUnreclaimable LinuxMemorySlabStateAttr = "unreclaimable"
)

// NetworkConnectionStateAttr is an attribute conforming to the
// network.connection.state semantic conventions. It represents the state of
// network connection.
type NetworkConnectionStateAttr string

var (
	// NetworkConnectionStateClosed is the none.
	NetworkConnectionStateClosed NetworkConnectionStateAttr = "closed"
	// NetworkConnectionStateCloseWait is the none.
	NetworkConnectionStateCloseWait NetworkConnectionStateAttr = "close_wait"
	// NetworkConnectionStateClosing is the none.
	NetworkConnectionStateClosing NetworkConnectionStateAttr = "closing"
	// NetworkConnectionStateEstablished is the none.
	NetworkConnectionStateEstablished NetworkConnectionStateAttr = "established"
	// NetworkConnectionStateFinWait1 is the none.
	NetworkConnectionStateFinWait1 NetworkConnectionStateAttr = "fin_wait_1"
	// NetworkConnectionStateFinWait2 is the none.
	NetworkConnectionStateFinWait2 NetworkConnectionStateAttr = "fin_wait_2"
	// NetworkConnectionStateLastAck is the none.
	NetworkConnectionStateLastAck NetworkConnectionStateAttr = "last_ack"
	// NetworkConnectionStateListen is the none.
	NetworkConnectionStateListen NetworkConnectionStateAttr = "listen"
	// NetworkConnectionStateSynReceived is the none.
	NetworkConnectionStateSynReceived NetworkConnectionStateAttr = "syn_received"
	// NetworkConnectionStateSynSent is the none.
	NetworkConnectionStateSynSent NetworkConnectionStateAttr = "syn_sent"
	// NetworkConnectionStateTimeWait is the none.
	NetworkConnectionStateTimeWait NetworkConnectionStateAttr = "time_wait"
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

// NetworkTransportAttr is an attribute conforming to the network.transport
// semantic conventions. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
type NetworkTransportAttr string

var (
	// NetworkTransportTCP is the TCP.
	NetworkTransportTCP NetworkTransportAttr = "tcp"
	// NetworkTransportUDP is the UDP.
	NetworkTransportUDP NetworkTransportAttr = "udp"
	// NetworkTransportPipe is the named or anonymous pipe.
	NetworkTransportPipe NetworkTransportAttr = "pipe"
	// NetworkTransportUnix is the unix domain socket.
	NetworkTransportUnix NetworkTransportAttr = "unix"
	// NetworkTransportQUIC is the QUIC.
	NetworkTransportQUIC NetworkTransportAttr = "quic"
)

// FilesystemStateAttr is an attribute conforming to the system.filesystem.state
// semantic conventions. It represents the filesystem state.
type FilesystemStateAttr string

var (
	// FilesystemStateUsed is the none.
	FilesystemStateUsed FilesystemStateAttr = "used"
	// FilesystemStateFree is the none.
	FilesystemStateFree FilesystemStateAttr = "free"
	// FilesystemStateReserved is the none.
	FilesystemStateReserved FilesystemStateAttr = "reserved"
)

// FilesystemTypeAttr is an attribute conforming to the system.filesystem.type
// semantic conventions. It represents the filesystem type.
type FilesystemTypeAttr string

var (
	// FilesystemTypeFat32 is the none.
	FilesystemTypeFat32 FilesystemTypeAttr = "fat32"
	// FilesystemTypeExfat is the none.
	FilesystemTypeExfat FilesystemTypeAttr = "exfat"
	// FilesystemTypeNtfs is the none.
	FilesystemTypeNtfs FilesystemTypeAttr = "ntfs"
	// FilesystemTypeRefs is the none.
	FilesystemTypeRefs FilesystemTypeAttr = "refs"
	// FilesystemTypeHfsplus is the none.
	FilesystemTypeHfsplus FilesystemTypeAttr = "hfsplus"
	// FilesystemTypeExt4 is the none.
	FilesystemTypeExt4 FilesystemTypeAttr = "ext4"
)

// MemoryStateAttr is an attribute conforming to the system.memory.state semantic
// conventions. It represents the memory state.
type MemoryStateAttr string

var (
	// MemoryStateUsed is the none.
	MemoryStateUsed MemoryStateAttr = "used"
	// MemoryStateFree is the none.
	MemoryStateFree MemoryStateAttr = "free"
	// MemoryStateBuffers is the none.
	MemoryStateBuffers MemoryStateAttr = "buffers"
	// MemoryStateCached is the none.
	MemoryStateCached MemoryStateAttr = "cached"
)

// PagingDirectionAttr is an attribute conforming to the system.paging.direction
// semantic conventions. It represents the paging access direction.
type PagingDirectionAttr string

var (
	// PagingDirectionIn is the none.
	PagingDirectionIn PagingDirectionAttr = "in"
	// PagingDirectionOut is the none.
	PagingDirectionOut PagingDirectionAttr = "out"
)

// PagingStateAttr is an attribute conforming to the system.paging.state semantic
// conventions. It represents the memory paging state.
type PagingStateAttr string

var (
	// PagingStateUsed is the none.
	PagingStateUsed PagingStateAttr = "used"
	// PagingStateFree is the none.
	PagingStateFree PagingStateAttr = "free"
)

// PagingTypeAttr is an attribute conforming to the system.paging.type semantic
// conventions. It represents the memory paging type.
type PagingTypeAttr string

var (
	// PagingTypeMajor is the none.
	PagingTypeMajor PagingTypeAttr = "major"
	// PagingTypeMinor is the none.
	PagingTypeMinor PagingTypeAttr = "minor"
)

// ProcessStatusAttr is an attribute conforming to the system.process.status
// semantic conventions. It represents the process state, e.g.,
// [Linux Process State Codes].
//
// [Linux Process State Codes]: https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES
type ProcessStatusAttr string

var (
	// ProcessStatusRunning is the none.
	ProcessStatusRunning ProcessStatusAttr = "running"
	// ProcessStatusSleeping is the none.
	ProcessStatusSleeping ProcessStatusAttr = "sleeping"
	// ProcessStatusStopped is the none.
	ProcessStatusStopped ProcessStatusAttr = "stopped"
	// ProcessStatusDefunct is the none.
	ProcessStatusDefunct ProcessStatusAttr = "defunct"
)

// CPUFrequency is an instrument used to record metric values conforming to the
// "system.cpu.frequency" semantic conventions. It represents the deprecated. Use
// `cpu.frequency` instead.
type CPUFrequency struct {
	inst metric.Int64Gauge
}

// NewCPUFrequency returns a new CPUFrequency instrument.
func NewCPUFrequency(m metric.Meter) (CPUFrequency, error) {
	i, err := m.Int64Gauge(
	    "system.cpu.frequency",
	    metric.WithDescription("Deprecated. Use `cpu.frequency` instead."),
	    metric.WithUnit("{Hz}"),
	)
	if err != nil {
	    return CPUFrequency{}, err
	}
	return CPUFrequency{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUFrequency) Name() string {
	return "system.cpu.frequency"
}

// Unit returns the semantic convention unit of the instrument
func (CPUFrequency) Unit() string {
	return "{Hz}"
}

// Description returns the semantic convention description of the instrument
func (CPUFrequency) Description() string {
	return "Deprecated. Use `cpu.frequency` instead."
}

func (m CPUFrequency) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// CPULogicalCount is an instrument used to record metric values conforming to
// the "system.cpu.logical.count" semantic conventions. It represents the reports
// the number of logical (virtual) processor cores created by the operating
// system to manage multitasking.
type CPULogicalCount struct {
	inst metric.Int64UpDownCounter
}

// NewCPULogicalCount returns a new CPULogicalCount instrument.
func NewCPULogicalCount(m metric.Meter) (CPULogicalCount, error) {
	i, err := m.Int64UpDownCounter(
	    "system.cpu.logical.count",
	    metric.WithDescription("Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"),
	    metric.WithUnit("{cpu}"),
	)
	if err != nil {
	    return CPULogicalCount{}, err
	}
	return CPULogicalCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPULogicalCount) Name() string {
	return "system.cpu.logical.count"
}

// Unit returns the semantic convention unit of the instrument
func (CPULogicalCount) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (CPULogicalCount) Description() string {
	return "Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"
}

func (m CPULogicalCount) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// CPUPhysicalCount is an instrument used to record metric values conforming to
// the "system.cpu.physical.count" semantic conventions. It represents the
// reports the number of actual physical processor cores on the hardware.
type CPUPhysicalCount struct {
	inst metric.Int64UpDownCounter
}

// NewCPUPhysicalCount returns a new CPUPhysicalCount instrument.
func NewCPUPhysicalCount(m metric.Meter) (CPUPhysicalCount, error) {
	i, err := m.Int64UpDownCounter(
	    "system.cpu.physical.count",
	    metric.WithDescription("Reports the number of actual physical processor cores on the hardware"),
	    metric.WithUnit("{cpu}"),
	)
	if err != nil {
	    return CPUPhysicalCount{}, err
	}
	return CPUPhysicalCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUPhysicalCount) Name() string {
	return "system.cpu.physical.count"
}

// Unit returns the semantic convention unit of the instrument
func (CPUPhysicalCount) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (CPUPhysicalCount) Description() string {
	return "Reports the number of actual physical processor cores on the hardware"
}

func (m CPUPhysicalCount) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// CPUTime is an instrument used to record metric values conforming to the
// "system.cpu.time" semantic conventions. It represents the deprecated. Use
// `cpu.time` instead.
type CPUTime struct {
	inst metric.Float64Counter
}

// NewCPUTime returns a new CPUTime instrument.
func NewCPUTime(m metric.Meter) (CPUTime, error) {
	i, err := m.Float64Counter(
	    "system.cpu.time",
	    metric.WithDescription("Deprecated. Use `cpu.time` instead."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return CPUTime{}, err
	}
	return CPUTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUTime) Name() string {
	return "system.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (CPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (CPUTime) Description() string {
	return "Deprecated. Use `cpu.time` instead."
}

func (m CPUTime) Add(ctx context.Context, incr float64) {
    m.inst.Add(ctx, incr)
}

// CPUUtilization is an instrument used to record metric values conforming to the
// "system.cpu.utilization" semantic conventions. It represents the deprecated.
// Use `cpu.utilization` instead.
type CPUUtilization struct {
	inst metric.Int64Gauge
}

// NewCPUUtilization returns a new CPUUtilization instrument.
func NewCPUUtilization(m metric.Meter) (CPUUtilization, error) {
	i, err := m.Int64Gauge(
	    "system.cpu.utilization",
	    metric.WithDescription("Deprecated. Use `cpu.utilization` instead."),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return CPUUtilization{}, err
	}
	return CPUUtilization{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CPUUtilization) Name() string {
	return "system.cpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (CPUUtilization) Description() string {
	return "Deprecated. Use `cpu.utilization` instead."
}

func (m CPUUtilization) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// DiskIo is an instrument used to record metric values conforming to the
// "system.disk.io" semantic conventions.
type DiskIo struct {
	inst metric.Int64Counter
}

// NewDiskIo returns a new DiskIo instrument.
func NewDiskIo(m metric.Meter) (DiskIo, error) {
	i, err := m.Int64Counter(
	    "system.disk.io",
	    metric.WithDescription(""),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return DiskIo{}, err
	}
	return DiskIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskIo) Name() string {
	return "system.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIo) Unit() string {
	return "By"
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

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskIo) Device(val string) DiskIoAttr {
	return diskIoAttr{kv: attribute.String("system.device", val)}
}

// DiskIoTime is an instrument used to record metric values conforming to the
// "system.disk.io_time" semantic conventions. It represents the time disk spent
// activated.
type DiskIoTime struct {
	inst metric.Float64Counter
}

// NewDiskIoTime returns a new DiskIoTime instrument.
func NewDiskIoTime(m metric.Meter) (DiskIoTime, error) {
	i, err := m.Float64Counter(
	    "system.disk.io_time",
	    metric.WithDescription("Time disk spent activated"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return DiskIoTime{}, err
	}
	return DiskIoTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskIoTime) Name() string {
	return "system.disk.io_time"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIoTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (DiskIoTime) Description() string {
	return "Time disk spent activated"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskIoTime) Add(
    ctx context.Context,
    incr float64,
	attrs ...DiskIoTimeAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskIoTime) conv(in []DiskIoTimeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskIoTimeAttr()
	}
	return out
}

// DiskIoTimeAttr is an optional attribute for the DiskIoTime instrument.
type DiskIoTimeAttr interface {
    diskIoTimeAttr() attribute.KeyValue
}

type diskIoTimeAttr struct {
	kv attribute.KeyValue
}

func (a diskIoTimeAttr) diskIoTimeAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskIoTime) Device(val string) DiskIoTimeAttr {
	return diskIoTimeAttr{kv: attribute.String("system.device", val)}
}

// DiskLimit is an instrument used to record metric values conforming to the
// "system.disk.limit" semantic conventions. It represents the total storage
// capacity of the disk.
type DiskLimit struct {
	inst metric.Int64UpDownCounter
}

// NewDiskLimit returns a new DiskLimit instrument.
func NewDiskLimit(m metric.Meter) (DiskLimit, error) {
	i, err := m.Int64UpDownCounter(
	    "system.disk.limit",
	    metric.WithDescription("The total storage capacity of the disk"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return DiskLimit{}, err
	}
	return DiskLimit{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskLimit) Name() string {
	return "system.disk.limit"
}

// Unit returns the semantic convention unit of the instrument
func (DiskLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskLimit) Description() string {
	return "The total storage capacity of the disk"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskLimit) Add(
    ctx context.Context,
    incr int64,
	attrs ...DiskLimitAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskLimit) conv(in []DiskLimitAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskLimitAttr()
	}
	return out
}

// DiskLimitAttr is an optional attribute for the DiskLimit instrument.
type DiskLimitAttr interface {
    diskLimitAttr() attribute.KeyValue
}

type diskLimitAttr struct {
	kv attribute.KeyValue
}

func (a diskLimitAttr) diskLimitAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskLimit) Device(val string) DiskLimitAttr {
	return diskLimitAttr{kv: attribute.String("system.device", val)}
}

// DiskMerged is an instrument used to record metric values conforming to the
// "system.disk.merged" semantic conventions.
type DiskMerged struct {
	inst metric.Int64Counter
}

// NewDiskMerged returns a new DiskMerged instrument.
func NewDiskMerged(m metric.Meter) (DiskMerged, error) {
	i, err := m.Int64Counter(
	    "system.disk.merged",
	    metric.WithDescription(""),
	    metric.WithUnit("{operation}"),
	)
	if err != nil {
	    return DiskMerged{}, err
	}
	return DiskMerged{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskMerged) Name() string {
	return "system.disk.merged"
}

// Unit returns the semantic convention unit of the instrument
func (DiskMerged) Unit() string {
	return "{operation}"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskMerged) Add(
    ctx context.Context,
    incr int64,
	attrs ...DiskMergedAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskMerged) conv(in []DiskMergedAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskMergedAttr()
	}
	return out
}

// DiskMergedAttr is an optional attribute for the DiskMerged instrument.
type DiskMergedAttr interface {
    diskMergedAttr() attribute.KeyValue
}

type diskMergedAttr struct {
	kv attribute.KeyValue
}

func (a diskMergedAttr) diskMergedAttr() attribute.KeyValue {
    return a.kv
}

// DiskIoDirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskMerged) DiskIoDirection(val DiskIoDirectionAttr) DiskMergedAttr {
	return diskMergedAttr{kv: attribute.String("disk.io.direction", string(val))}
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskMerged) Device(val string) DiskMergedAttr {
	return diskMergedAttr{kv: attribute.String("system.device", val)}
}

// DiskOperationTime is an instrument used to record metric values conforming to
// the "system.disk.operation_time" semantic conventions. It represents the sum
// of the time each operation took to complete.
type DiskOperationTime struct {
	inst metric.Float64Counter
}

// NewDiskOperationTime returns a new DiskOperationTime instrument.
func NewDiskOperationTime(m metric.Meter) (DiskOperationTime, error) {
	i, err := m.Float64Counter(
	    "system.disk.operation_time",
	    metric.WithDescription("Sum of the time each operation took to complete"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return DiskOperationTime{}, err
	}
	return DiskOperationTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskOperationTime) Name() string {
	return "system.disk.operation_time"
}

// Unit returns the semantic convention unit of the instrument
func (DiskOperationTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (DiskOperationTime) Description() string {
	return "Sum of the time each operation took to complete"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskOperationTime) Add(
    ctx context.Context,
    incr float64,
	attrs ...DiskOperationTimeAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskOperationTime) conv(in []DiskOperationTimeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskOperationTimeAttr()
	}
	return out
}

// DiskOperationTimeAttr is an optional attribute for the DiskOperationTime
// instrument.
type DiskOperationTimeAttr interface {
    diskOperationTimeAttr() attribute.KeyValue
}

type diskOperationTimeAttr struct {
	kv attribute.KeyValue
}

func (a diskOperationTimeAttr) diskOperationTimeAttr() attribute.KeyValue {
    return a.kv
}

// DiskIoDirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskOperationTime) DiskIoDirection(val DiskIoDirectionAttr) DiskOperationTimeAttr {
	return diskOperationTimeAttr{kv: attribute.String("disk.io.direction", string(val))}
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskOperationTime) Device(val string) DiskOperationTimeAttr {
	return diskOperationTimeAttr{kv: attribute.String("system.device", val)}
}

// DiskOperations is an instrument used to record metric values conforming to the
// "system.disk.operations" semantic conventions.
type DiskOperations struct {
	inst metric.Int64Counter
}

// NewDiskOperations returns a new DiskOperations instrument.
func NewDiskOperations(m metric.Meter) (DiskOperations, error) {
	i, err := m.Int64Counter(
	    "system.disk.operations",
	    metric.WithDescription(""),
	    metric.WithUnit("{operation}"),
	)
	if err != nil {
	    return DiskOperations{}, err
	}
	return DiskOperations{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DiskOperations) Name() string {
	return "system.disk.operations"
}

// Unit returns the semantic convention unit of the instrument
func (DiskOperations) Unit() string {
	return "{operation}"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskOperations) Add(
    ctx context.Context,
    incr int64,
	attrs ...DiskOperationsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m DiskOperations) conv(in []DiskOperationsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.diskOperationsAttr()
	}
	return out
}

// DiskOperationsAttr is an optional attribute for the DiskOperations instrument.
type DiskOperationsAttr interface {
    diskOperationsAttr() attribute.KeyValue
}

type diskOperationsAttr struct {
	kv attribute.KeyValue
}

func (a diskOperationsAttr) diskOperationsAttr() attribute.KeyValue {
    return a.kv
}

// DiskIoDirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskOperations) DiskIoDirection(val DiskIoDirectionAttr) DiskOperationsAttr {
	return diskOperationsAttr{kv: attribute.String("disk.io.direction", string(val))}
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskOperations) Device(val string) DiskOperationsAttr {
	return diskOperationsAttr{kv: attribute.String("system.device", val)}
}

// FilesystemLimit is an instrument used to record metric values conforming to
// the "system.filesystem.limit" semantic conventions. It represents the total
// storage capacity of the filesystem.
type FilesystemLimit struct {
	inst metric.Int64UpDownCounter
}

// NewFilesystemLimit returns a new FilesystemLimit instrument.
func NewFilesystemLimit(m metric.Meter) (FilesystemLimit, error) {
	i, err := m.Int64UpDownCounter(
	    "system.filesystem.limit",
	    metric.WithDescription("The total storage capacity of the filesystem"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return FilesystemLimit{}, err
	}
	return FilesystemLimit{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (FilesystemLimit) Name() string {
	return "system.filesystem.limit"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemLimit) Description() string {
	return "The total storage capacity of the filesystem"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m FilesystemLimit) Add(
    ctx context.Context,
    incr int64,
	attrs ...FilesystemLimitAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m FilesystemLimit) conv(in []FilesystemLimitAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.filesystemLimitAttr()
	}
	return out
}

// FilesystemLimitAttr is an optional attribute for the FilesystemLimit
// instrument.
type FilesystemLimitAttr interface {
    filesystemLimitAttr() attribute.KeyValue
}

type filesystemLimitAttr struct {
	kv attribute.KeyValue
}

func (a filesystemLimitAttr) filesystemLimitAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemLimit) Device(val string) FilesystemLimitAttr {
	return filesystemLimitAttr{kv: attribute.String("system.device", val)}
}

// FilesystemMode returns an optional attribute for the "system.filesystem.mode"
// semantic convention. It represents the filesystem mode.
func (FilesystemLimit) FilesystemMode(val string) FilesystemLimitAttr {
	return filesystemLimitAttr{kv: attribute.String("system.filesystem.mode", val)}
}

// FilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemLimit) FilesystemMountpoint(val string) FilesystemLimitAttr {
	return filesystemLimitAttr{kv: attribute.String("system.filesystem.mountpoint", val)}
}

// FilesystemType returns an optional attribute for the "system.filesystem.type"
// semantic convention. It represents the filesystem type.
func (FilesystemLimit) FilesystemType(val FilesystemTypeAttr) FilesystemLimitAttr {
	return filesystemLimitAttr{kv: attribute.String("system.filesystem.type", string(val))}
}

// FilesystemUsage is an instrument used to record metric values conforming to
// the "system.filesystem.usage" semantic conventions. It represents the reports
// a filesystem's space usage across different states.
type FilesystemUsage struct {
	inst metric.Int64UpDownCounter
}

// NewFilesystemUsage returns a new FilesystemUsage instrument.
func NewFilesystemUsage(m metric.Meter) (FilesystemUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "system.filesystem.usage",
	    metric.WithDescription("Reports a filesystem's space usage across different states."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return FilesystemUsage{}, err
	}
	return FilesystemUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (FilesystemUsage) Name() string {
	return "system.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemUsage) Description() string {
	return "Reports a filesystem's space usage across different states."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m FilesystemUsage) Add(
    ctx context.Context,
    incr int64,
	attrs ...FilesystemUsageAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m FilesystemUsage) conv(in []FilesystemUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.filesystemUsageAttr()
	}
	return out
}

// FilesystemUsageAttr is an optional attribute for the FilesystemUsage
// instrument.
type FilesystemUsageAttr interface {
    filesystemUsageAttr() attribute.KeyValue
}

type filesystemUsageAttr struct {
	kv attribute.KeyValue
}

func (a filesystemUsageAttr) filesystemUsageAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemUsage) Device(val string) FilesystemUsageAttr {
	return filesystemUsageAttr{kv: attribute.String("system.device", val)}
}

// FilesystemMode returns an optional attribute for the "system.filesystem.mode"
// semantic convention. It represents the filesystem mode.
func (FilesystemUsage) FilesystemMode(val string) FilesystemUsageAttr {
	return filesystemUsageAttr{kv: attribute.String("system.filesystem.mode", val)}
}

// FilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemUsage) FilesystemMountpoint(val string) FilesystemUsageAttr {
	return filesystemUsageAttr{kv: attribute.String("system.filesystem.mountpoint", val)}
}

// FilesystemState returns an optional attribute for the
// "system.filesystem.state" semantic convention. It represents the filesystem
// state.
func (FilesystemUsage) FilesystemState(val FilesystemStateAttr) FilesystemUsageAttr {
	return filesystemUsageAttr{kv: attribute.String("system.filesystem.state", string(val))}
}

// FilesystemType returns an optional attribute for the "system.filesystem.type"
// semantic convention. It represents the filesystem type.
func (FilesystemUsage) FilesystemType(val FilesystemTypeAttr) FilesystemUsageAttr {
	return filesystemUsageAttr{kv: attribute.String("system.filesystem.type", string(val))}
}

// FilesystemUtilization is an instrument used to record metric values conforming
// to the "system.filesystem.utilization" semantic conventions.
type FilesystemUtilization struct {
	inst metric.Int64Gauge
}

// NewFilesystemUtilization returns a new FilesystemUtilization instrument.
func NewFilesystemUtilization(m metric.Meter) (FilesystemUtilization, error) {
	i, err := m.Int64Gauge(
	    "system.filesystem.utilization",
	    metric.WithDescription(""),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return FilesystemUtilization{}, err
	}
	return FilesystemUtilization{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (FilesystemUtilization) Name() string {
	return "system.filesystem.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemUtilization) Unit() string {
	return "1"
}

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m FilesystemUtilization) Record(
    ctx context.Context,
    val int64,
	attrs ...FilesystemUtilizationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m FilesystemUtilization) conv(in []FilesystemUtilizationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.filesystemUtilizationAttr()
	}
	return out
}

// FilesystemUtilizationAttr is an optional attribute for the
// FilesystemUtilization instrument.
type FilesystemUtilizationAttr interface {
    filesystemUtilizationAttr() attribute.KeyValue
}

type filesystemUtilizationAttr struct {
	kv attribute.KeyValue
}

func (a filesystemUtilizationAttr) filesystemUtilizationAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemUtilization) Device(val string) FilesystemUtilizationAttr {
	return filesystemUtilizationAttr{kv: attribute.String("system.device", val)}
}

// FilesystemMode returns an optional attribute for the "system.filesystem.mode"
// semantic convention. It represents the filesystem mode.
func (FilesystemUtilization) FilesystemMode(val string) FilesystemUtilizationAttr {
	return filesystemUtilizationAttr{kv: attribute.String("system.filesystem.mode", val)}
}

// FilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemUtilization) FilesystemMountpoint(val string) FilesystemUtilizationAttr {
	return filesystemUtilizationAttr{kv: attribute.String("system.filesystem.mountpoint", val)}
}

// FilesystemState returns an optional attribute for the
// "system.filesystem.state" semantic convention. It represents the filesystem
// state.
func (FilesystemUtilization) FilesystemState(val FilesystemStateAttr) FilesystemUtilizationAttr {
	return filesystemUtilizationAttr{kv: attribute.String("system.filesystem.state", string(val))}
}

// FilesystemType returns an optional attribute for the "system.filesystem.type"
// semantic convention. It represents the filesystem type.
func (FilesystemUtilization) FilesystemType(val FilesystemTypeAttr) FilesystemUtilizationAttr {
	return filesystemUtilizationAttr{kv: attribute.String("system.filesystem.type", string(val))}
}

// LinuxMemoryAvailable is an instrument used to record metric values conforming
// to the "system.linux.memory.available" semantic conventions. It represents an
// estimate of how much memory is available for starting new applications,
// without causing swapping.
type LinuxMemoryAvailable struct {
	inst metric.Int64UpDownCounter
}

// NewLinuxMemoryAvailable returns a new LinuxMemoryAvailable instrument.
func NewLinuxMemoryAvailable(m metric.Meter) (LinuxMemoryAvailable, error) {
	i, err := m.Int64UpDownCounter(
	    "system.linux.memory.available",
	    metric.WithDescription("An estimate of how much memory is available for starting new applications, without causing swapping"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return LinuxMemoryAvailable{}, err
	}
	return LinuxMemoryAvailable{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (LinuxMemoryAvailable) Name() string {
	return "system.linux.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (LinuxMemoryAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LinuxMemoryAvailable) Description() string {
	return "An estimate of how much memory is available for starting new applications, without causing swapping"
}

func (m LinuxMemoryAvailable) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// LinuxMemorySlabUsage is an instrument used to record metric values conforming
// to the "system.linux.memory.slab.usage" semantic conventions. It represents
// the reports the memory used by the Linux kernel for managing caches of
// frequently used objects.
type LinuxMemorySlabUsage struct {
	inst metric.Int64UpDownCounter
}

// NewLinuxMemorySlabUsage returns a new LinuxMemorySlabUsage instrument.
func NewLinuxMemorySlabUsage(m metric.Meter) (LinuxMemorySlabUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "system.linux.memory.slab.usage",
	    metric.WithDescription("Reports the memory used by the Linux kernel for managing caches of frequently used objects."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return LinuxMemorySlabUsage{}, err
	}
	return LinuxMemorySlabUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (LinuxMemorySlabUsage) Name() string {
	return "system.linux.memory.slab.usage"
}

// Unit returns the semantic convention unit of the instrument
func (LinuxMemorySlabUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LinuxMemorySlabUsage) Description() string {
	return "Reports the memory used by the Linux kernel for managing caches of frequently used objects."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m LinuxMemorySlabUsage) Add(
    ctx context.Context,
    incr int64,
	attrs ...LinuxMemorySlabUsageAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m LinuxMemorySlabUsage) conv(in []LinuxMemorySlabUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.linuxMemorySlabUsageAttr()
	}
	return out
}

// LinuxMemorySlabUsageAttr is an optional attribute for the LinuxMemorySlabUsage
// instrument.
type LinuxMemorySlabUsageAttr interface {
    linuxMemorySlabUsageAttr() attribute.KeyValue
}

type linuxMemorySlabUsageAttr struct {
	kv attribute.KeyValue
}

func (a linuxMemorySlabUsageAttr) linuxMemorySlabUsageAttr() attribute.KeyValue {
    return a.kv
}

// LinuxMemorySlabState returns an optional attribute for the
// "linux.memory.slab.state" semantic convention. It represents the Linux Slab
// memory state.
func (LinuxMemorySlabUsage) LinuxMemorySlabState(val LinuxMemorySlabStateAttr) LinuxMemorySlabUsageAttr {
	return linuxMemorySlabUsageAttr{kv: attribute.String("linux.memory.slab.state", string(val))}
}

// MemoryLimit is an instrument used to record metric values conforming to the
// "system.memory.limit" semantic conventions. It represents the total memory
// available in the system.
type MemoryLimit struct {
	inst metric.Int64UpDownCounter
}

// NewMemoryLimit returns a new MemoryLimit instrument.
func NewMemoryLimit(m metric.Meter) (MemoryLimit, error) {
	i, err := m.Int64UpDownCounter(
	    "system.memory.limit",
	    metric.WithDescription("Total memory available in the system."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryLimit{}, err
	}
	return MemoryLimit{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (MemoryLimit) Name() string {
	return "system.memory.limit"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryLimit) Description() string {
	return "Total memory available in the system."
}

func (m MemoryLimit) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// MemoryShared is an instrument used to record metric values conforming to the
// "system.memory.shared" semantic conventions. It represents the shared memory
// used (mostly by tmpfs).
type MemoryShared struct {
	inst metric.Int64UpDownCounter
}

// NewMemoryShared returns a new MemoryShared instrument.
func NewMemoryShared(m metric.Meter) (MemoryShared, error) {
	i, err := m.Int64UpDownCounter(
	    "system.memory.shared",
	    metric.WithDescription("Shared memory used (mostly by tmpfs)."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryShared{}, err
	}
	return MemoryShared{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (MemoryShared) Name() string {
	return "system.memory.shared"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryShared) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryShared) Description() string {
	return "Shared memory used (mostly by tmpfs)."
}

func (m MemoryShared) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "system.memory.usage" semantic conventions. It represents the reports memory
// in use by state.
type MemoryUsage struct {
	inst metric.Int64UpDownCounter
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(m metric.Meter) (MemoryUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "system.memory.usage",
	    metric.WithDescription("Reports memory in use by state."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemoryUsage{}, err
	}
	return MemoryUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (MemoryUsage) Name() string {
	return "system.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryUsage) Description() string {
	return "Reports memory in use by state."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m MemoryUsage) Add(
    ctx context.Context,
    incr int64,
	attrs ...MemoryUsageAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m MemoryUsage) conv(in []MemoryUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.memoryUsageAttr()
	}
	return out
}

// MemoryUsageAttr is an optional attribute for the MemoryUsage instrument.
type MemoryUsageAttr interface {
    memoryUsageAttr() attribute.KeyValue
}

type memoryUsageAttr struct {
	kv attribute.KeyValue
}

func (a memoryUsageAttr) memoryUsageAttr() attribute.KeyValue {
    return a.kv
}

// MemoryState returns an optional attribute for the "system.memory.state"
// semantic convention. It represents the memory state.
func (MemoryUsage) MemoryState(val MemoryStateAttr) MemoryUsageAttr {
	return memoryUsageAttr{kv: attribute.String("system.memory.state", string(val))}
}

// MemoryUtilization is an instrument used to record metric values conforming to
// the "system.memory.utilization" semantic conventions.
type MemoryUtilization struct {
	inst metric.Int64Gauge
}

// NewMemoryUtilization returns a new MemoryUtilization instrument.
func NewMemoryUtilization(m metric.Meter) (MemoryUtilization, error) {
	i, err := m.Int64Gauge(
	    "system.memory.utilization",
	    metric.WithDescription(""),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return MemoryUtilization{}, err
	}
	return MemoryUtilization{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (MemoryUtilization) Name() string {
	return "system.memory.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUtilization) Unit() string {
	return "1"
}

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m MemoryUtilization) Record(
    ctx context.Context,
    val int64,
	attrs ...MemoryUtilizationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m MemoryUtilization) conv(in []MemoryUtilizationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.memoryUtilizationAttr()
	}
	return out
}

// MemoryUtilizationAttr is an optional attribute for the MemoryUtilization
// instrument.
type MemoryUtilizationAttr interface {
    memoryUtilizationAttr() attribute.KeyValue
}

type memoryUtilizationAttr struct {
	kv attribute.KeyValue
}

func (a memoryUtilizationAttr) memoryUtilizationAttr() attribute.KeyValue {
    return a.kv
}

// MemoryState returns an optional attribute for the "system.memory.state"
// semantic convention. It represents the memory state.
func (MemoryUtilization) MemoryState(val MemoryStateAttr) MemoryUtilizationAttr {
	return memoryUtilizationAttr{kv: attribute.String("system.memory.state", string(val))}
}

// NetworkConnections is an instrument used to record metric values conforming to
// the "system.network.connections" semantic conventions.
type NetworkConnections struct {
	inst metric.Int64UpDownCounter
}

// NewNetworkConnections returns a new NetworkConnections instrument.
func NewNetworkConnections(m metric.Meter) (NetworkConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "system.network.connections",
	    metric.WithDescription(""),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return NetworkConnections{}, err
	}
	return NetworkConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkConnections) Name() string {
	return "system.network.connections"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkConnections) Unit() string {
	return "{connection}"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkConnections) Add(
    ctx context.Context,
    incr int64,
	attrs ...NetworkConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetworkConnections) conv(in []NetworkConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.networkConnectionsAttr()
	}
	return out
}

// NetworkConnectionsAttr is an optional attribute for the NetworkConnections
// instrument.
type NetworkConnectionsAttr interface {
    networkConnectionsAttr() attribute.KeyValue
}

type networkConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a networkConnectionsAttr) networkConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkConnectionState returns an optional attribute for the
// "network.connection.state" semantic convention. It represents the state of
// network connection.
func (NetworkConnections) NetworkConnectionState(val NetworkConnectionStateAttr) NetworkConnectionsAttr {
	return networkConnectionsAttr{kv: attribute.String("network.connection.state", string(val))}
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkConnections) NetworkInterfaceName(val string) NetworkConnectionsAttr {
	return networkConnectionsAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (NetworkConnections) NetworkTransport(val NetworkTransportAttr) NetworkConnectionsAttr {
	return networkConnectionsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkDropped is an instrument used to record metric values conforming to the
// "system.network.dropped" semantic conventions. It represents the count of
// packets that are dropped or discarded even though there was no error.
type NetworkDropped struct {
	inst metric.Int64Counter
}

// NewNetworkDropped returns a new NetworkDropped instrument.
func NewNetworkDropped(m metric.Meter) (NetworkDropped, error) {
	i, err := m.Int64Counter(
	    "system.network.dropped",
	    metric.WithDescription("Count of packets that are dropped or discarded even though there was no error"),
	    metric.WithUnit("{packet}"),
	)
	if err != nil {
	    return NetworkDropped{}, err
	}
	return NetworkDropped{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkDropped) Name() string {
	return "system.network.dropped"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkDropped) Unit() string {
	return "{packet}"
}

// Description returns the semantic convention description of the instrument
func (NetworkDropped) Description() string {
	return "Count of packets that are dropped or discarded even though there was no error"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkDropped) Add(
    ctx context.Context,
    incr int64,
	attrs ...NetworkDroppedAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetworkDropped) conv(in []NetworkDroppedAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.networkDroppedAttr()
	}
	return out
}

// NetworkDroppedAttr is an optional attribute for the NetworkDropped instrument.
type NetworkDroppedAttr interface {
    networkDroppedAttr() attribute.KeyValue
}

type networkDroppedAttr struct {
	kv attribute.KeyValue
}

func (a networkDroppedAttr) networkDroppedAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkDropped) NetworkInterfaceName(val string) NetworkDroppedAttr {
	return networkDroppedAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkDropped) NetworkIoDirection(val NetworkIoDirectionAttr) NetworkDroppedAttr {
	return networkDroppedAttr{kv: attribute.String("network.io.direction", string(val))}
}

// NetworkErrors is an instrument used to record metric values conforming to the
// "system.network.errors" semantic conventions. It represents the count of
// network errors detected.
type NetworkErrors struct {
	inst metric.Int64Counter
}

// NewNetworkErrors returns a new NetworkErrors instrument.
func NewNetworkErrors(m metric.Meter) (NetworkErrors, error) {
	i, err := m.Int64Counter(
	    "system.network.errors",
	    metric.WithDescription("Count of network errors detected"),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return NetworkErrors{}, err
	}
	return NetworkErrors{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkErrors) Name() string {
	return "system.network.errors"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkErrors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (NetworkErrors) Description() string {
	return "Count of network errors detected"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkErrors) Add(
    ctx context.Context,
    incr int64,
	attrs ...NetworkErrorsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetworkErrors) conv(in []NetworkErrorsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.networkErrorsAttr()
	}
	return out
}

// NetworkErrorsAttr is an optional attribute for the NetworkErrors instrument.
type NetworkErrorsAttr interface {
    networkErrorsAttr() attribute.KeyValue
}

type networkErrorsAttr struct {
	kv attribute.KeyValue
}

func (a networkErrorsAttr) networkErrorsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkErrors) NetworkInterfaceName(val string) NetworkErrorsAttr {
	return networkErrorsAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkErrors) NetworkIoDirection(val NetworkIoDirectionAttr) NetworkErrorsAttr {
	return networkErrorsAttr{kv: attribute.String("network.io.direction", string(val))}
}

// NetworkIo is an instrument used to record metric values conforming to the
// "system.network.io" semantic conventions.
type NetworkIo struct {
	inst metric.Int64Counter
}

// NewNetworkIo returns a new NetworkIo instrument.
func NewNetworkIo(m metric.Meter) (NetworkIo, error) {
	i, err := m.Int64Counter(
	    "system.network.io",
	    metric.WithDescription(""),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NetworkIo{}, err
	}
	return NetworkIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkIo) Name() string {
	return "system.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIo) Unit() string {
	return "By"
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
func (NetworkIo) NetworkInterfaceName(val string) NetworkIoAttr {
	return networkIoAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIo) NetworkIoDirection(val NetworkIoDirectionAttr) NetworkIoAttr {
	return networkIoAttr{kv: attribute.String("network.io.direction", string(val))}
}

// NetworkPackets is an instrument used to record metric values conforming to the
// "system.network.packets" semantic conventions.
type NetworkPackets struct {
	inst metric.Int64Counter
}

// NewNetworkPackets returns a new NetworkPackets instrument.
func NewNetworkPackets(m metric.Meter) (NetworkPackets, error) {
	i, err := m.Int64Counter(
	    "system.network.packets",
	    metric.WithDescription(""),
	    metric.WithUnit("{packet}"),
	)
	if err != nil {
	    return NetworkPackets{}, err
	}
	return NetworkPackets{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetworkPackets) Name() string {
	return "system.network.packets"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkPackets) Unit() string {
	return "{packet}"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkPackets) Add(
    ctx context.Context,
    incr int64,
	attrs ...NetworkPacketsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetworkPackets) conv(in []NetworkPacketsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.networkPacketsAttr()
	}
	return out
}

// NetworkPacketsAttr is an optional attribute for the NetworkPackets instrument.
type NetworkPacketsAttr interface {
    networkPacketsAttr() attribute.KeyValue
}

type networkPacketsAttr struct {
	kv attribute.KeyValue
}

func (a networkPacketsAttr) networkPacketsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkPackets) NetworkIoDirection(val NetworkIoDirectionAttr) NetworkPacketsAttr {
	return networkPacketsAttr{kv: attribute.String("network.io.direction", string(val))}
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (NetworkPackets) Device(val string) NetworkPacketsAttr {
	return networkPacketsAttr{kv: attribute.String("system.device", val)}
}

// PagingFaults is an instrument used to record metric values conforming to the
// "system.paging.faults" semantic conventions.
type PagingFaults struct {
	inst metric.Int64Counter
}

// NewPagingFaults returns a new PagingFaults instrument.
func NewPagingFaults(m metric.Meter) (PagingFaults, error) {
	i, err := m.Int64Counter(
	    "system.paging.faults",
	    metric.WithDescription(""),
	    metric.WithUnit("{fault}"),
	)
	if err != nil {
	    return PagingFaults{}, err
	}
	return PagingFaults{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PagingFaults) Name() string {
	return "system.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PagingFaults) Unit() string {
	return "{fault}"
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

// PagingType returns an optional attribute for the "system.paging.type" semantic
// convention. It represents the memory paging type.
func (PagingFaults) PagingType(val PagingTypeAttr) PagingFaultsAttr {
	return pagingFaultsAttr{kv: attribute.String("system.paging.type", string(val))}
}

// PagingOperations is an instrument used to record metric values conforming to
// the "system.paging.operations" semantic conventions.
type PagingOperations struct {
	inst metric.Int64Counter
}

// NewPagingOperations returns a new PagingOperations instrument.
func NewPagingOperations(m metric.Meter) (PagingOperations, error) {
	i, err := m.Int64Counter(
	    "system.paging.operations",
	    metric.WithDescription(""),
	    metric.WithUnit("{operation}"),
	)
	if err != nil {
	    return PagingOperations{}, err
	}
	return PagingOperations{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PagingOperations) Name() string {
	return "system.paging.operations"
}

// Unit returns the semantic convention unit of the instrument
func (PagingOperations) Unit() string {
	return "{operation}"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PagingOperations) Add(
    ctx context.Context,
    incr int64,
	attrs ...PagingOperationsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PagingOperations) conv(in []PagingOperationsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.pagingOperationsAttr()
	}
	return out
}

// PagingOperationsAttr is an optional attribute for the PagingOperations
// instrument.
type PagingOperationsAttr interface {
    pagingOperationsAttr() attribute.KeyValue
}

type pagingOperationsAttr struct {
	kv attribute.KeyValue
}

func (a pagingOperationsAttr) pagingOperationsAttr() attribute.KeyValue {
    return a.kv
}

// PagingDirection returns an optional attribute for the
// "system.paging.direction" semantic convention. It represents the paging access
// direction.
func (PagingOperations) PagingDirection(val PagingDirectionAttr) PagingOperationsAttr {
	return pagingOperationsAttr{kv: attribute.String("system.paging.direction", string(val))}
}

// PagingType returns an optional attribute for the "system.paging.type" semantic
// convention. It represents the memory paging type.
func (PagingOperations) PagingType(val PagingTypeAttr) PagingOperationsAttr {
	return pagingOperationsAttr{kv: attribute.String("system.paging.type", string(val))}
}

// PagingUsage is an instrument used to record metric values conforming to the
// "system.paging.usage" semantic conventions. It represents the unix swap or
// windows pagefile usage.
type PagingUsage struct {
	inst metric.Int64UpDownCounter
}

// NewPagingUsage returns a new PagingUsage instrument.
func NewPagingUsage(m metric.Meter) (PagingUsage, error) {
	i, err := m.Int64UpDownCounter(
	    "system.paging.usage",
	    metric.WithDescription("Unix swap or windows pagefile usage"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return PagingUsage{}, err
	}
	return PagingUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PagingUsage) Name() string {
	return "system.paging.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PagingUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PagingUsage) Description() string {
	return "Unix swap or windows pagefile usage"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PagingUsage) Add(
    ctx context.Context,
    incr int64,
	attrs ...PagingUsageAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PagingUsage) conv(in []PagingUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.pagingUsageAttr()
	}
	return out
}

// PagingUsageAttr is an optional attribute for the PagingUsage instrument.
type PagingUsageAttr interface {
    pagingUsageAttr() attribute.KeyValue
}

type pagingUsageAttr struct {
	kv attribute.KeyValue
}

func (a pagingUsageAttr) pagingUsageAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the unique identifier for the device responsible for
// managing paging operations.
func (PagingUsage) Device(val string) PagingUsageAttr {
	return pagingUsageAttr{kv: attribute.String("system.device", val)}
}

// PagingState returns an optional attribute for the "system.paging.state"
// semantic convention. It represents the memory paging state.
func (PagingUsage) PagingState(val PagingStateAttr) PagingUsageAttr {
	return pagingUsageAttr{kv: attribute.String("system.paging.state", string(val))}
}

// PagingUtilization is an instrument used to record metric values conforming to
// the "system.paging.utilization" semantic conventions.
type PagingUtilization struct {
	inst metric.Int64Gauge
}

// NewPagingUtilization returns a new PagingUtilization instrument.
func NewPagingUtilization(m metric.Meter) (PagingUtilization, error) {
	i, err := m.Int64Gauge(
	    "system.paging.utilization",
	    metric.WithDescription(""),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return PagingUtilization{}, err
	}
	return PagingUtilization{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PagingUtilization) Name() string {
	return "system.paging.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PagingUtilization) Unit() string {
	return "1"
}

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PagingUtilization) Record(
    ctx context.Context,
    val int64,
	attrs ...PagingUtilizationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PagingUtilization) conv(in []PagingUtilizationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.pagingUtilizationAttr()
	}
	return out
}

// PagingUtilizationAttr is an optional attribute for the PagingUtilization
// instrument.
type PagingUtilizationAttr interface {
    pagingUtilizationAttr() attribute.KeyValue
}

type pagingUtilizationAttr struct {
	kv attribute.KeyValue
}

func (a pagingUtilizationAttr) pagingUtilizationAttr() attribute.KeyValue {
    return a.kv
}

// Device returns an optional attribute for the "system.device" semantic
// convention. It represents the unique identifier for the device responsible for
// managing paging operations.
func (PagingUtilization) Device(val string) PagingUtilizationAttr {
	return pagingUtilizationAttr{kv: attribute.String("system.device", val)}
}

// PagingState returns an optional attribute for the "system.paging.state"
// semantic convention. It represents the memory paging state.
func (PagingUtilization) PagingState(val PagingStateAttr) PagingUtilizationAttr {
	return pagingUtilizationAttr{kv: attribute.String("system.paging.state", string(val))}
}

// ProcessCount is an instrument used to record metric values conforming to the
// "system.process.count" semantic conventions. It represents the total number of
// processes in each state.
type ProcessCount struct {
	inst metric.Int64UpDownCounter
}

// NewProcessCount returns a new ProcessCount instrument.
func NewProcessCount(m metric.Meter) (ProcessCount, error) {
	i, err := m.Int64UpDownCounter(
	    "system.process.count",
	    metric.WithDescription("Total number of processes in each state"),
	    metric.WithUnit("{process}"),
	)
	if err != nil {
	    return ProcessCount{}, err
	}
	return ProcessCount{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ProcessCount) Name() string {
	return "system.process.count"
}

// Unit returns the semantic convention unit of the instrument
func (ProcessCount) Unit() string {
	return "{process}"
}

// Description returns the semantic convention description of the instrument
func (ProcessCount) Description() string {
	return "Total number of processes in each state"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ProcessCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...ProcessCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ProcessCount) conv(in []ProcessCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.processCountAttr()
	}
	return out
}

// ProcessCountAttr is an optional attribute for the ProcessCount instrument.
type ProcessCountAttr interface {
    processCountAttr() attribute.KeyValue
}

type processCountAttr struct {
	kv attribute.KeyValue
}

func (a processCountAttr) processCountAttr() attribute.KeyValue {
    return a.kv
}

// ProcessStatus returns an optional attribute for the "system.process.status"
// semantic convention. It represents the process state, e.g.,
// [Linux Process State Codes].
//
// [Linux Process State Codes]: https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES
func (ProcessCount) ProcessStatus(val ProcessStatusAttr) ProcessCountAttr {
	return processCountAttr{kv: attribute.String("system.process.status", string(val))}
}

// ProcessCreated is an instrument used to record metric values conforming to the
// "system.process.created" semantic conventions. It represents the total number
// of processes created over uptime of the host.
type ProcessCreated struct {
	inst metric.Int64Counter
}

// NewProcessCreated returns a new ProcessCreated instrument.
func NewProcessCreated(m metric.Meter) (ProcessCreated, error) {
	i, err := m.Int64Counter(
	    "system.process.created",
	    metric.WithDescription("Total number of processes created over uptime of the host"),
	    metric.WithUnit("{process}"),
	)
	if err != nil {
	    return ProcessCreated{}, err
	}
	return ProcessCreated{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ProcessCreated) Name() string {
	return "system.process.created"
}

// Unit returns the semantic convention unit of the instrument
func (ProcessCreated) Unit() string {
	return "{process}"
}

// Description returns the semantic convention description of the instrument
func (ProcessCreated) Description() string {
	return "Total number of processes created over uptime of the host"
}

func (m ProcessCreated) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// Uptime is an instrument used to record metric values conforming to the
// "system.uptime" semantic conventions. It represents the time the system has
// been running.
type Uptime struct {
	inst metric.Float64Gauge
}

// NewUptime returns a new Uptime instrument.
func NewUptime(m metric.Meter) (Uptime, error) {
	i, err := m.Float64Gauge(
	    "system.uptime",
	    metric.WithDescription("The time the system has been running"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return Uptime{}, err
	}
	return Uptime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (Uptime) Name() string {
	return "system.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (Uptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (Uptime) Description() string {
	return "The time the system has been running"
}

func (m Uptime) Record(ctx context.Context, val float64) {
    m.inst.Record(ctx, val)
}