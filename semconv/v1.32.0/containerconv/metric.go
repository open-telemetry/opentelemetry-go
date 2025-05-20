// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "container" namespace.
package containerconv

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
	// CPUModeIOWait is the none.
	CPUModeIOWait CPUModeAttr = "iowait"
	// CPUModeInterrupt is the none.
	CPUModeInterrupt CPUModeAttr = "interrupt"
	// CPUModeSteal is the none.
	CPUModeSteal CPUModeAttr = "steal"
	// CPUModeKernel is the none.
	CPUModeKernel CPUModeAttr = "kernel"
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

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the none.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the none.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
)

// CPUTime is an instrument used to record metric values conforming to the
// "container.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type CPUTime struct {
	metric.Float64Counter
}

// NewCPUTime returns a new CPUTime instrument.
func NewCPUTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (CPUTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUTime{noop.Float64Counter{}}, nil
	}

	i, err := m.Float64Counter(
		"container.cpu.time",
		append([]metric.Float64CounterOption{
			metric.WithDescription("Total CPU time consumed"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUTime{noop.Float64Counter{}}, err
	}
	return CPUTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUTime) Inst() metric.Float64Counter {
	return m.Float64Counter
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
//
// Total CPU time consumed by the specific container on all available CPU cores
func (m CPUTime) Add(
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

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUTime) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUUsage is an instrument used to record metric values conforming to the
// "container.cpu.usage" semantic conventions. It represents the container's CPU
// usage, measured in cpus. Range from 0 to the number of allocatable CPUs.
type CPUUsage struct {
	metric.Int64Gauge
}

// NewCPUUsage returns a new CPUUsage instrument.
func NewCPUUsage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (CPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUsage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"container.cpu.usage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUUsage{noop.Int64Gauge{}}, err
	}
	return CPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUsage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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
//
// CPU usage of the specific container on all available CPU cores, averaged over
// the sample window
func (m CPUUsage) Record(
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

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUUsage) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// DiskIO is an instrument used to record metric values conforming to the
// "container.disk.io" semantic conventions. It represents the disk bytes for the
// container.
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
		"container.disk.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Disk bytes for the container."),
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
	return "container.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIO) Description() string {
	return "Disk bytes for the container."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// The total number of bytes read/written successfully (aggregated from all
// disks).
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

// AttrSystemDevice returns an optional attribute for the "system.device"
// semantic convention. It represents the device identifier.
func (DiskIO) AttrSystemDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "container.memory.usage" semantic conventions. It represents the memory usage
// of the container.
type MemoryUsage struct {
	metric.Int64Counter
}

// NewMemoryUsage returns a new MemoryUsage instrument.
func NewMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (MemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUsage{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"container.memory.usage",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Memory usage of the container."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return MemoryUsage{noop.Int64Counter{}}, err
	}
	return MemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUsage) Inst() metric.Int64Counter {
	return m.Int64Counter
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

// Add adds incr to the existing count.
//
// Memory usage of the container.
func (m MemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// NetworkIO is an instrument used to record metric values conforming to the
// "container.network.io" semantic conventions. It represents the network bytes
// for the container.
type NetworkIO struct {
	metric.Int64Counter
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

	i, err := m.Int64Counter(
		"container.network.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Network bytes for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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
	return "container.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIO) Description() string {
	return "Network bytes for the container."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// The number of bytes sent/received on all network interfaces by the container.
func (m NetworkIO) Add(
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
func (NetworkIO) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// Uptime is an instrument used to record metric values conforming to the
// "container.uptime" semantic conventions. It represents the time the container
// has been running.
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
		"container.uptime",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The time the container has been running"),
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