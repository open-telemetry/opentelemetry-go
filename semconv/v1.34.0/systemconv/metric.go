// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "system" namespace.
package systemconv

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

// DiskIODirectionAttr is an attribute conforming to the disk.io.direction
// semantic conventions. It represents the disk IO operation direction.
type DiskIODirectionAttr string

var (
	// DiskIODirectionRead is the none.
	DiskIODirectionRead DiskIODirectionAttr = "read"
	// DiskIODirectionWrite is the none.
	DiskIODirectionWrite DiskIODirectionAttr = "write"
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

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the none.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the none.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
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

// CPULogicalCount is an instrument used to record metric values conforming to
// the "system.cpu.logical.count" semantic conventions. It represents the reports
// the number of logical (virtual) processor cores created by the operating
// system to manage multitasking.
type CPULogicalCount struct {
	metric.Int64UpDownCounter
}

// NewCPULogicalCount returns a new CPULogicalCount instrument.
func NewCPULogicalCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (CPULogicalCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPULogicalCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.cpu.logical.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Reports the number of logical (virtual) processor cores created by the operating system to manage multitasking"),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return CPULogicalCount{noop.Int64UpDownCounter{}}, err
	}
	return CPULogicalCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPULogicalCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// Calculated by multiplying the number of sockets by the number of cores per
// socket, and then by the number of threads per core
func (m CPULogicalCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// CPUPhysicalCount is an instrument used to record metric values conforming to
// the "system.cpu.physical.count" semantic conventions. It represents the
// reports the number of actual physical processor cores on the hardware.
type CPUPhysicalCount struct {
	metric.Int64UpDownCounter
}

// NewCPUPhysicalCount returns a new CPUPhysicalCount instrument.
func NewCPUPhysicalCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (CPUPhysicalCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUPhysicalCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.cpu.physical.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Reports the number of actual physical processor cores on the hardware"),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUPhysicalCount{noop.Int64UpDownCounter{}}, err
	}
	return CPUPhysicalCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUPhysicalCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// Calculated by multiplying the number of sockets by the number of cores per
// socket
func (m CPUPhysicalCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DiskIO is an instrument used to record metric values conforming to the
// "system.disk.io" semantic conventions.
type DiskIO struct {
	metric.Int64Counter
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

	i, err := m.Int64Counter(
		"system.disk.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("By"),
		}, opt...)...,
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
	return "system.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIO) Unit() string {
	return "By"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m DiskIO) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskIO) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskIO) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// DiskIOTime is an instrument used to record metric values conforming to the
// "system.disk.io_time" semantic conventions. It represents the time disk spent
// activated.
type DiskIOTime struct {
	metric.Float64Counter
}

// NewDiskIOTime returns a new DiskIOTime instrument.
func NewDiskIOTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (DiskIOTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskIOTime{noop.Float64Counter{}}, nil
	}

	i, err := m.Float64Counter(
		"system.disk.io_time",
		append([]metric.Float64CounterOption{
			metric.WithDescription("Time disk spent activated"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return DiskIOTime{noop.Float64Counter{}}, err
	}
	return DiskIOTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskIOTime) Inst() metric.Float64Counter {
	return m.Float64Counter
}

// Name returns the semantic convention name of the instrument.
func (DiskIOTime) Name() string {
	return "system.disk.io_time"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIOTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (DiskIOTime) Description() string {
	return "Time disk spent activated"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// The real elapsed time ("wall clock") used in the I/O path (time from
// operations running in parallel are not counted). Measured as:
//
//   - Linux: Field 13 from [procfs-diskstats]
//   - Windows: The complement of
//     ["Disk% Idle Time"]
//     performance counter: `uptime * (100 - "Disk\% Idle Time") / 100`
//
//
// [procfs-diskstats]: https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// ["Disk% Idle Time"]: https://learn.microsoft.com/archive/blogs/askcore/windows-performance-monitor-disk-counters-explained#windows-performance-monitor-disk-counters-explained
func (m DiskIOTime) Add(
	ctx context.Context,
	incr float64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Float64Counter.Add(ctx, incr, *o...)
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskIOTime) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// DiskLimit is an instrument used to record metric values conforming to the
// "system.disk.limit" semantic conventions. It represents the total storage
// capacity of the disk.
type DiskLimit struct {
	metric.Int64UpDownCounter
}

// NewDiskLimit returns a new DiskLimit instrument.
func NewDiskLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DiskLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.disk.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The total storage capacity of the disk"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return DiskLimit{noop.Int64UpDownCounter{}}, err
	}
	return DiskLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskLimit) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// DiskMerged is an instrument used to record metric values conforming to the
// "system.disk.merged" semantic conventions.
type DiskMerged struct {
	metric.Int64Counter
}

// NewDiskMerged returns a new DiskMerged instrument.
func NewDiskMerged(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (DiskMerged, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskMerged{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.disk.merged",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{operation}"),
		}, opt...)...,
	)
	if err != nil {
	    return DiskMerged{noop.Int64Counter{}}, err
	}
	return DiskMerged{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskMerged) Inst() metric.Int64Counter {
	return m.Int64Counter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskMerged) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskMerged) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// DiskOperationTime is an instrument used to record metric values conforming to
// the "system.disk.operation_time" semantic conventions. It represents the sum
// of the time each operation took to complete.
type DiskOperationTime struct {
	metric.Float64Counter
}

// NewDiskOperationTime returns a new DiskOperationTime instrument.
func NewDiskOperationTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (DiskOperationTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskOperationTime{noop.Float64Counter{}}, nil
	}

	i, err := m.Float64Counter(
		"system.disk.operation_time",
		append([]metric.Float64CounterOption{
			metric.WithDescription("Sum of the time each operation took to complete"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return DiskOperationTime{noop.Float64Counter{}}, err
	}
	return DiskOperationTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskOperationTime) Inst() metric.Float64Counter {
	return m.Float64Counter
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
//
// Because it is the sum of time each request took, parallel-issued requests each
// contribute to make the count grow. Measured as:
//
//   - Linux: Fields 7 & 11 from [procfs-diskstats]
//   - Windows: "Avg. Disk sec/Read" perf counter multiplied by "Disk Reads/sec"
//     perf counter (similar for Writes)
//
//
// [procfs-diskstats]: https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
func (m DiskOperationTime) Add(
	ctx context.Context,
	incr float64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Float64Counter.Add(ctx, incr, *o...)
}

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskOperationTime) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskOperationTime) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// DiskOperations is an instrument used to record metric values conforming to the
// "system.disk.operations" semantic conventions.
type DiskOperations struct {
	metric.Int64Counter
}

// NewDiskOperations returns a new DiskOperations instrument.
func NewDiskOperations(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (DiskOperations, error) {
	// Check if the meter is nil.
	if m == nil {
		return DiskOperations{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.disk.operations",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{operation}"),
		}, opt...)...,
	)
	if err != nil {
	    return DiskOperations{noop.Int64Counter{}}, err
	}
	return DiskOperations{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DiskOperations) Inst() metric.Int64Counter {
	return m.Int64Counter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskOperations) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (DiskOperations) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// FilesystemLimit is an instrument used to record metric values conforming to
// the "system.filesystem.limit" semantic conventions. It represents the total
// storage capacity of the filesystem.
type FilesystemLimit struct {
	metric.Int64UpDownCounter
}

// NewFilesystemLimit returns a new FilesystemLimit instrument.
func NewFilesystemLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (FilesystemLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.filesystem.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The total storage capacity of the filesystem"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return FilesystemLimit{noop.Int64UpDownCounter{}}, err
	}
	return FilesystemLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemLimit) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// AttrFilesystemMode returns an optional attribute for the
// "system.filesystem.mode" semantic convention. It represents the filesystem
// mode.
func (FilesystemLimit) AttrFilesystemMode(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mode", val)
}

// AttrFilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemLimit) AttrFilesystemMountpoint(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mountpoint", val)
}

// AttrFilesystemType returns an optional attribute for the
// "system.filesystem.type" semantic convention. It represents the filesystem
// type.
func (FilesystemLimit) AttrFilesystemType(val FilesystemTypeAttr) attribute.KeyValue {
	return attribute.String("system.filesystem.type", string(val))
}

// FilesystemUsage is an instrument used to record metric values conforming to
// the "system.filesystem.usage" semantic conventions. It represents the reports
// a filesystem's space usage across different states.
type FilesystemUsage struct {
	metric.Int64UpDownCounter
}

// NewFilesystemUsage returns a new FilesystemUsage instrument.
func NewFilesystemUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (FilesystemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.filesystem.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Reports a filesystem's space usage across different states."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return FilesystemUsage{noop.Int64UpDownCounter{}}, err
	}
	return FilesystemUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
//
// The sum of all `system.filesystem.usage` values over the different
// `system.filesystem.state` attributes
// SHOULD equal the total storage capacity of the filesystem, that is
// `system.filesystem.limit`.
func (m FilesystemUsage) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemUsage) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// AttrFilesystemMode returns an optional attribute for the
// "system.filesystem.mode" semantic convention. It represents the filesystem
// mode.
func (FilesystemUsage) AttrFilesystemMode(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mode", val)
}

// AttrFilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemUsage) AttrFilesystemMountpoint(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mountpoint", val)
}

// AttrFilesystemState returns an optional attribute for the
// "system.filesystem.state" semantic convention. It represents the filesystem
// state.
func (FilesystemUsage) AttrFilesystemState(val FilesystemStateAttr) attribute.KeyValue {
	return attribute.String("system.filesystem.state", string(val))
}

// AttrFilesystemType returns an optional attribute for the
// "system.filesystem.type" semantic convention. It represents the filesystem
// type.
func (FilesystemUsage) AttrFilesystemType(val FilesystemTypeAttr) attribute.KeyValue {
	return attribute.String("system.filesystem.type", string(val))
}

// FilesystemUtilization is an instrument used to record metric values conforming
// to the "system.filesystem.utilization" semantic conventions.
type FilesystemUtilization struct {
	metric.Int64Gauge
}

// NewFilesystemUtilization returns a new FilesystemUtilization instrument.
func NewFilesystemUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (FilesystemUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"system.filesystem.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription(""),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return FilesystemUtilization{noop.Int64Gauge{}}, err
	}
	return FilesystemUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (FilesystemUtilization) Name() string {
	return "system.filesystem.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemUtilization) Unit() string {
	return "1"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m FilesystemUtilization) Record(
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

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the identifier for the device where the filesystem
// resides.
func (FilesystemUtilization) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// AttrFilesystemMode returns an optional attribute for the
// "system.filesystem.mode" semantic convention. It represents the filesystem
// mode.
func (FilesystemUtilization) AttrFilesystemMode(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mode", val)
}

// AttrFilesystemMountpoint returns an optional attribute for the
// "system.filesystem.mountpoint" semantic convention. It represents the
// filesystem mount path.
func (FilesystemUtilization) AttrFilesystemMountpoint(val string) attribute.KeyValue {
	return attribute.String("system.filesystem.mountpoint", val)
}

// AttrFilesystemState returns an optional attribute for the
// "system.filesystem.state" semantic convention. It represents the filesystem
// state.
func (FilesystemUtilization) AttrFilesystemState(val FilesystemStateAttr) attribute.KeyValue {
	return attribute.String("system.filesystem.state", string(val))
}

// AttrFilesystemType returns an optional attribute for the
// "system.filesystem.type" semantic convention. It represents the filesystem
// type.
func (FilesystemUtilization) AttrFilesystemType(val FilesystemTypeAttr) attribute.KeyValue {
	return attribute.String("system.filesystem.type", string(val))
}

// LinuxMemoryAvailable is an instrument used to record metric values conforming
// to the "system.linux.memory.available" semantic conventions. It represents an
// estimate of how much memory is available for starting new applications,
// without causing swapping.
type LinuxMemoryAvailable struct {
	metric.Int64UpDownCounter
}

// NewLinuxMemoryAvailable returns a new LinuxMemoryAvailable instrument.
func NewLinuxMemoryAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (LinuxMemoryAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return LinuxMemoryAvailable{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.linux.memory.available",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("An estimate of how much memory is available for starting new applications, without causing swapping"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return LinuxMemoryAvailable{noop.Int64UpDownCounter{}}, err
	}
	return LinuxMemoryAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LinuxMemoryAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This is an alternative to `system.memory.usage` metric with `state=free`.
// Linux starting from 3.14 exports "available" memory. It takes "free" memory as
// a baseline, and then factors in kernel-specific values.
// This is supposed to be more accurate than just "free" memory.
// For reference, see the calculations [here].
// See also `MemAvailable` in [/proc/meminfo].
//
// [here]: https://superuser.com/a/980821
// [/proc/meminfo]: https://man7.org/linux/man-pages/man5/proc.5.html
func (m LinuxMemoryAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// LinuxMemorySlabUsage is an instrument used to record metric values conforming
// to the "system.linux.memory.slab.usage" semantic conventions. It represents
// the reports the memory used by the Linux kernel for managing caches of
// frequently used objects.
type LinuxMemorySlabUsage struct {
	metric.Int64UpDownCounter
}

// NewLinuxMemorySlabUsage returns a new LinuxMemorySlabUsage instrument.
func NewLinuxMemorySlabUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (LinuxMemorySlabUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return LinuxMemorySlabUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.linux.memory.slab.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Reports the memory used by the Linux kernel for managing caches of frequently used objects."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return LinuxMemorySlabUsage{noop.Int64UpDownCounter{}}, err
	}
	return LinuxMemorySlabUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LinuxMemorySlabUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
//
// The sum over the `reclaimable` and `unreclaimable` state values in
// `linux.memory.slab.usage` SHOULD be equal to the total slab memory available
// on the system.
// Note that the total slab memory is not constant and may vary over time.
// See also the [Slab allocator] and `Slab` in [/proc/meminfo].
//
// [Slab allocator]: https://blogs.oracle.com/linux/post/understanding-linux-kernel-memory-statistics
// [/proc/meminfo]: https://man7.org/linux/man-pages/man5/proc.5.html
func (m LinuxMemorySlabUsage) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrLinuxMemorySlabState returns an optional attribute for the
// "linux.memory.slab.state" semantic convention. It represents the Linux Slab
// memory state.
func (LinuxMemorySlabUsage) AttrLinuxMemorySlabState(val LinuxMemorySlabStateAttr) attribute.KeyValue {
	return attribute.String("linux.memory.slab.state", string(val))
}

// MemoryLimit is an instrument used to record metric values conforming to the
// "system.memory.limit" semantic conventions. It represents the total memory
// available in the system.
type MemoryLimit struct {
	metric.Int64UpDownCounter
}

// NewMemoryLimit returns a new MemoryLimit instrument.
func NewMemoryLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.memory.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Total memory available in the system."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return MemoryLimit{noop.Int64UpDownCounter{}}, err
	}
	return MemoryLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// Its value SHOULD equal the sum of `system.memory.state` over all states.
func (m MemoryLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// MemoryShared is an instrument used to record metric values conforming to the
// "system.memory.shared" semantic conventions. It represents the shared memory
// used (mostly by tmpfs).
type MemoryShared struct {
	metric.Int64UpDownCounter
}

// NewMemoryShared returns a new MemoryShared instrument.
func NewMemoryShared(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryShared, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryShared{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.memory.shared",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Shared memory used (mostly by tmpfs)."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return MemoryShared{noop.Int64UpDownCounter{}}, err
	}
	return MemoryShared{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryShared) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// Equivalent of `shared` from [`free` command] or
// `Shmem` from [`/proc/meminfo`]"
//
// [`free` command]: https://man7.org/linux/man-pages/man1/free.1.html
// [`/proc/meminfo`]: https://man7.org/linux/man-pages/man5/proc.5.html
func (m MemoryShared) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "system.memory.usage" semantic conventions. It represents the reports memory
// in use by state.
type MemoryUsage struct {
	metric.Int64ObservableUpDownCounter
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

	i, err := m.Int64ObservableUpDownCounter(
		"system.memory.usage",
		append([]metric.Int64ObservableUpDownCounterOption{
			metric.WithDescription("Reports memory in use by state."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// AttrMemoryState returns an optional attribute for the "system.memory.state"
// semantic convention. It represents the memory state.
func (MemoryUsage) AttrMemoryState(val MemoryStateAttr) attribute.KeyValue {
	return attribute.String("system.memory.state", string(val))
}

// MemoryUtilization is an instrument used to record metric values conforming to
// the "system.memory.utilization" semantic conventions.
type MemoryUtilization struct {
	metric.Float64ObservableGauge
}

// NewMemoryUtilization returns a new MemoryUtilization instrument.
func NewMemoryUtilization(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (MemoryUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUtilization{noop.Float64ObservableGauge{}}, nil
	}

	i, err := m.Float64ObservableGauge(
		"system.memory.utilization",
		append([]metric.Float64ObservableGaugeOption{
			metric.WithDescription(""),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return MemoryUtilization{noop.Float64ObservableGauge{}}, err
	}
	return MemoryUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUtilization) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (MemoryUtilization) Name() string {
	return "system.memory.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUtilization) Unit() string {
	return "1"
}

// AttrMemoryState returns an optional attribute for the "system.memory.state"
// semantic convention. It represents the memory state.
func (MemoryUtilization) AttrMemoryState(val MemoryStateAttr) attribute.KeyValue {
	return attribute.String("system.memory.state", string(val))
}

// NetworkConnections is an instrument used to record metric values conforming to
// the "system.network.connections" semantic conventions.
type NetworkConnections struct {
	metric.Int64UpDownCounter
}

// NewNetworkConnections returns a new NetworkConnections instrument.
func NewNetworkConnections(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NetworkConnections, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkConnections{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.network.connections",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{connection}"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkConnections{noop.Int64UpDownCounter{}}, err
	}
	return NetworkConnections{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkConnections) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrNetworkConnectionState returns an optional attribute for the
// "network.connection.state" semantic convention. It represents the state of
// network connection.
func (NetworkConnections) AttrNetworkConnectionState(val NetworkConnectionStateAttr) attribute.KeyValue {
	return attribute.String("network.connection.state", string(val))
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkConnections) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (NetworkConnections) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// NetworkDropped is an instrument used to record metric values conforming to the
// "system.network.dropped" semantic conventions. It represents the count of
// packets that are dropped or discarded even though there was no error.
type NetworkDropped struct {
	metric.Int64Counter
}

// NewNetworkDropped returns a new NetworkDropped instrument.
func NewNetworkDropped(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NetworkDropped, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkDropped{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.network.dropped",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Count of packets that are dropped or discarded even though there was no error"),
			metric.WithUnit("{packet}"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkDropped{noop.Int64Counter{}}, err
	}
	return NetworkDropped{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkDropped) Inst() metric.Int64Counter {
	return m.Int64Counter
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
//
// Measured as:
//
//   - Linux: the `drop` column in `/proc/dev/net` ([source])
//   - Windows: [`InDiscards`/`OutDiscards`]
//     from [`GetIfEntry2`]
//
//
// [source]: https://web.archive.org/web/20180321091318/http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html
// [`InDiscards`/`OutDiscards`]: https://docs.microsoft.com/windows/win32/api/netioapi/ns-netioapi-mib_if_row2
// [`GetIfEntry2`]: https://docs.microsoft.com/windows/win32/api/netioapi/nf-netioapi-getifentry2
func (m NetworkDropped) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkDropped) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkDropped) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NetworkErrors is an instrument used to record metric values conforming to the
// "system.network.errors" semantic conventions. It represents the count of
// network errors detected.
type NetworkErrors struct {
	metric.Int64Counter
}

// NewNetworkErrors returns a new NetworkErrors instrument.
func NewNetworkErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NetworkErrors, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkErrors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.network.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Count of network errors detected"),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkErrors{noop.Int64Counter{}}, err
	}
	return NetworkErrors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkErrors) Inst() metric.Int64Counter {
	return m.Int64Counter
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
//
// Measured as:
//
//   - Linux: the `errs` column in `/proc/dev/net` ([source]).
//   - Windows: [`InErrors`/`OutErrors`]
//     from [`GetIfEntry2`].
//
//
// [source]: https://web.archive.org/web/20180321091318/http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html
// [`InErrors`/`OutErrors`]: https://docs.microsoft.com/windows/win32/api/netioapi/ns-netioapi-mib_if_row2
// [`GetIfEntry2`]: https://docs.microsoft.com/windows/win32/api/netioapi/nf-netioapi-getifentry2
func (m NetworkErrors) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkErrors) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkErrors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NetworkIO is an instrument used to record metric values conforming to the
// "system.network.io" semantic conventions.
type NetworkIO struct {
	metric.Int64ObservableCounter
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

	i, err := m.Int64ObservableCounter(
		"system.network.io",
		append([]metric.Int64ObservableCounterOption{
			metric.WithDescription(""),
			metric.WithUnit("By"),
		}, opt...)...,
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
	return "system.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIO) Unit() string {
	return "By"
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkIO) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NetworkPackets is an instrument used to record metric values conforming to the
// "system.network.packets" semantic conventions.
type NetworkPackets struct {
	metric.Int64Counter
}

// NewNetworkPackets returns a new NetworkPackets instrument.
func NewNetworkPackets(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NetworkPackets, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkPackets{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.network.packets",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{packet}"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkPackets{noop.Int64Counter{}}, err
	}
	return NetworkPackets{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkPackets) Inst() metric.Int64Counter {
	return m.Int64Counter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkPackets) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the device identifier.
func (NetworkPackets) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// PagingFaults is an instrument used to record metric values conforming to the
// "system.paging.faults" semantic conventions.
type PagingFaults struct {
	metric.Int64Counter
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

	i, err := m.Int64Counter(
		"system.paging.faults",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{fault}"),
		}, opt...)...,
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrPagingType returns an optional attribute for the "system.paging.type"
// semantic convention. It represents the memory paging type.
func (PagingFaults) AttrPagingType(val PagingTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.type", string(val))
}

// PagingOperations is an instrument used to record metric values conforming to
// the "system.paging.operations" semantic conventions.
type PagingOperations struct {
	metric.Int64Counter
}

// NewPagingOperations returns a new PagingOperations instrument.
func NewPagingOperations(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PagingOperations, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingOperations{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.paging.operations",
		append([]metric.Int64CounterOption{
			metric.WithDescription(""),
			metric.WithUnit("{operation}"),
		}, opt...)...,
	)
	if err != nil {
	    return PagingOperations{noop.Int64Counter{}}, err
	}
	return PagingOperations{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingOperations) Inst() metric.Int64Counter {
	return m.Int64Counter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrPagingDirection returns an optional attribute for the
// "system.paging.direction" semantic convention. It represents the paging access
// direction.
func (PagingOperations) AttrPagingDirection(val PagingDirectionAttr) attribute.KeyValue {
	return attribute.String("system.paging.direction", string(val))
}

// AttrPagingType returns an optional attribute for the "system.paging.type"
// semantic convention. It represents the memory paging type.
func (PagingOperations) AttrPagingType(val PagingTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.type", string(val))
}

// PagingUsage is an instrument used to record metric values conforming to the
// "system.paging.usage" semantic conventions. It represents the unix swap or
// windows pagefile usage.
type PagingUsage struct {
	metric.Int64UpDownCounter
}

// NewPagingUsage returns a new PagingUsage instrument.
func NewPagingUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PagingUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.paging.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Unix swap or windows pagefile usage"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return PagingUsage{noop.Int64UpDownCounter{}}, err
	}
	return PagingUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the unique identifier for the device responsible for
// managing paging operations.
func (PagingUsage) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// AttrPagingState returns an optional attribute for the "system.paging.state"
// semantic convention. It represents the memory paging state.
func (PagingUsage) AttrPagingState(val PagingStateAttr) attribute.KeyValue {
	return attribute.String("system.paging.state", string(val))
}

// PagingUtilization is an instrument used to record metric values conforming to
// the "system.paging.utilization" semantic conventions.
type PagingUtilization struct {
	metric.Int64Gauge
}

// NewPagingUtilization returns a new PagingUtilization instrument.
func NewPagingUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PagingUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PagingUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"system.paging.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription(""),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return PagingUtilization{noop.Int64Gauge{}}, err
	}
	return PagingUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PagingUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PagingUtilization) Name() string {
	return "system.paging.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PagingUtilization) Unit() string {
	return "1"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m PagingUtilization) Record(
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

// AttrDevice returns an optional attribute for the "system.device" semantic
// convention. It represents the unique identifier for the device responsible for
// managing paging operations.
func (PagingUtilization) AttrDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// AttrPagingState returns an optional attribute for the "system.paging.state"
// semantic convention. It represents the memory paging state.
func (PagingUtilization) AttrPagingState(val PagingStateAttr) attribute.KeyValue {
	return attribute.String("system.paging.state", string(val))
}

// ProcessCount is an instrument used to record metric values conforming to the
// "system.process.count" semantic conventions. It represents the total number of
// processes in each state.
type ProcessCount struct {
	metric.Int64UpDownCounter
}

// NewProcessCount returns a new ProcessCount instrument.
func NewProcessCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ProcessCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ProcessCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"system.process.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Total number of processes in each state"),
			metric.WithUnit("{process}"),
		}, opt...)...,
	)
	if err != nil {
	    return ProcessCount{noop.Int64UpDownCounter{}}, err
	}
	return ProcessCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ProcessCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrProcessStatus returns an optional attribute for the
// "system.process.status" semantic convention. It represents the process state,
// e.g., [Linux Process State Codes].
//
// [Linux Process State Codes]: https://man7.org/linux/man-pages/man1/ps.1.html#PROCESS_STATE_CODES
func (ProcessCount) AttrProcessStatus(val ProcessStatusAttr) attribute.KeyValue {
	return attribute.String("system.process.status", string(val))
}

// ProcessCreated is an instrument used to record metric values conforming to the
// "system.process.created" semantic conventions. It represents the total number
// of processes created over uptime of the host.
type ProcessCreated struct {
	metric.Int64Counter
}

// NewProcessCreated returns a new ProcessCreated instrument.
func NewProcessCreated(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ProcessCreated, error) {
	// Check if the meter is nil.
	if m == nil {
		return ProcessCreated{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"system.process.created",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Total number of processes created over uptime of the host"),
			metric.WithUnit("{process}"),
		}, opt...)...,
	)
	if err != nil {
	    return ProcessCreated{noop.Int64Counter{}}, err
	}
	return ProcessCreated{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ProcessCreated) Inst() metric.Int64Counter {
	return m.Int64Counter
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

// Add adds incr to the existing count.
func (m ProcessCreated) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// Uptime is an instrument used to record metric values conforming to the
// "system.uptime" semantic conventions. It represents the time the system has
// been running.
type Uptime struct {
	metric.Float64Gauge
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

	i, err := m.Float64Gauge(
		"system.uptime",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The time the system has been running"),
			metric.WithUnit("s"),
		}, opt...)...,
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

// Record records val to the current distribution.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m Uptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}