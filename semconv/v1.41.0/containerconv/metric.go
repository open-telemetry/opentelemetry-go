// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package containerconv provides types and functionality for OpenTelemetry semantic
// conventions in the "container" namespace.
package containerconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
)

// CPUModeAttr is an attribute conforming to the cpu.mode semantic conventions.
// It represents the CPU mode for this data point. A container's CPU metric
// SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
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

// CPUTime is an instrument used to record metric values conforming to the
// "container.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type CPUTime struct {
	metric.Float64Counter
}

var newCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newCPUTimeOpts
	} else {
		opt = append(opt, newCPUTimeOpts...)
	}

	i, err := m.Float64Counter(
		"container.cpu.time",
		opt...,
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
	return "Total CPU time consumed."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Total CPU time consumed by the specific container on all available CPU cores
func (m CPUTime) Add(
	ctx context.Context,
	incr float64,
	attrs ...attribute.KeyValue,
) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
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

	m.Float64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Total CPU time consumed by the specific container on all available CPU cores
func (m CPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUTime) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUTimeObservable is an instrument used to record metric values conforming to
// the "container.cpu.time" semantic conventions. It represents the total CPU
// time consumed.
type CPUTimeObservable struct {
	metric.Float64ObservableCounter
}

var newCPUTimeObservableOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
}

// NewCPUTimeObservable returns a new CPUTimeObservable instrument.
func NewCPUTimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (CPUTimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUTimeObservable{noop.Float64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUTimeObservableOpts
	} else {
		opt = append(opt, newCPUTimeObservableOpts...)
	}

	i, err := m.Float64ObservableCounter(
		"container.cpu.time",
		opt...,
	)
	if err != nil {
		return CPUTimeObservable{noop.Float64ObservableCounter{}}, err
	}
	return CPUTimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUTimeObservable) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (CPUTimeObservable) Name() string {
	return "container.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (CPUTimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (CPUTimeObservable) Description() string {
	return "Total CPU time consumed."
}

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUTimeObservable) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUUsage is an instrument used to record metric values conforming to the
// "container.cpu.usage" semantic conventions. It represents the container's CPU
// usage, measured in cpus. Range from 0 to the number of allocatable CPUs.
type CPUUsage struct {
	metric.Int64Gauge
}

var newCPUUsageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
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

	if len(opt) == 0 {
		opt = newCPUUsageOpts
	} else {
		opt = append(opt, newCPUUsageOpts...)
	}

	i, err := m.Int64Gauge(
		"container.cpu.usage",
		opt...,
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
	return "Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."
}

// Record records val to the current distribution for attrs.
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific container on all available CPU cores, averaged over
// the sample window
func (m CPUUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUUsage) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// CPUUsageObservable is an instrument used to record metric values conforming to
// the "container.cpu.usage" semantic conventions. It represents the container's
// CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs.
type CPUUsageObservable struct {
	metric.Int64ObservableGauge
}

var newCPUUsageObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewCPUUsageObservable returns a new CPUUsageObservable instrument.
func NewCPUUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (CPUUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUsageObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUUsageObservableOpts
	} else {
		opt = append(opt, newCPUUsageObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"container.cpu.usage",
		opt...,
	)
	if err != nil {
		return CPUUsageObservable{noop.Int64ObservableGauge{}}, err
	}
	return CPUUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUsageObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (CPUUsageObservable) Name() string {
	return "container.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUsageObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (CPUUsageObservable) Description() string {
	return "Container's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."
}

// AttrCPUMode returns an optional attribute for the "cpu.mode" semantic
// convention. It represents the CPU mode for this data point. A container's CPU
// metric SHOULD be characterized *either* by data points with no `mode` labels,
// *or only* data points with `mode` labels.
func (CPUUsageObservable) AttrCPUMode(val CPUModeAttr) attribute.KeyValue {
	return attribute.String("cpu.mode", string(val))
}

// DiskIO is an instrument used to record metric values conforming to the
// "container.disk.io" semantic conventions. It represents the disk bytes for the
// container.
type DiskIO struct {
	metric.Int64Counter
}

var newDiskIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Disk bytes for the container."),
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
		"container.disk.io",
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

// Add adds incr to the existing count for attrs.
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
//
// The total number of bytes read/written successfully (aggregated from all
// disks).
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

// DiskIOObservable is an instrument used to record metric values conforming to
// the "container.disk.io" semantic conventions. It represents the disk bytes for
// the container.
type DiskIOObservable struct {
	metric.Int64ObservableCounter
}

var newDiskIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Disk bytes for the container."),
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
		"container.disk.io",
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
	return "container.disk.io"
}

// Unit returns the semantic convention unit of the instrument
func (DiskIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (DiskIOObservable) Description() string {
	return "Disk bytes for the container."
}

// AttrDiskIODirection returns an optional attribute for the "disk.io.direction"
// semantic convention. It represents the disk IO operation direction.
func (DiskIOObservable) AttrDiskIODirection(val DiskIODirectionAttr) attribute.KeyValue {
	return attribute.String("disk.io.direction", string(val))
}

// AttrSystemDevice returns an optional attribute for the "system.device"
// semantic convention. It represents the device identifier.
func (DiskIOObservable) AttrSystemDevice(val string) attribute.KeyValue {
	return attribute.String("system.device", val)
}

// FilesystemAvailable is an instrument used to record metric values conforming
// to the "container.filesystem.available" semantic conventions. It represents
// the container filesystem available bytes.
type FilesystemAvailable struct {
	metric.Int64UpDownCounter
}

var newFilesystemAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewFilesystemAvailable returns a new FilesystemAvailable instrument.
func NewFilesystemAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (FilesystemAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newFilesystemAvailableOpts
	} else {
		opt = append(opt, newFilesystemAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.filesystem.available",
		opt...,
	)
	if err != nil {
		return FilesystemAvailable{noop.Int64UpDownCounter{}}, err
	}
	return FilesystemAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (FilesystemAvailable) Name() string {
	return "container.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemAvailable) Description() string {
	return "Container filesystem available bytes."
}

// Add adds incr to the existing count for attrs.
//
// In K8s, this metric is derived from the
// [FsStats.AvailableBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// In K8s, this metric is derived from the
// [FsStats.AvailableBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// FilesystemAvailableObservable is an instrument used to record metric values
// conforming to the "container.filesystem.available" semantic conventions. It
// represents the container filesystem available bytes.
type FilesystemAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newFilesystemAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewFilesystemAvailableObservable returns a new FilesystemAvailableObservable
// instrument.
func NewFilesystemAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (FilesystemAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newFilesystemAvailableObservableOpts
	} else {
		opt = append(opt, newFilesystemAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.filesystem.available",
		opt...,
	)
	if err != nil {
		return FilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return FilesystemAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (FilesystemAvailableObservable) Name() string {
	return "container.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemAvailableObservable) Description() string {
	return "Container filesystem available bytes."
}

// FilesystemCapacity is an instrument used to record metric values conforming to
// the "container.filesystem.capacity" semantic conventions. It represents the
// container filesystem capacity.
type FilesystemCapacity struct {
	metric.Int64UpDownCounter
}

var newFilesystemCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container filesystem capacity."),
	metric.WithUnit("By"),
}

// NewFilesystemCapacity returns a new FilesystemCapacity instrument.
func NewFilesystemCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (FilesystemCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newFilesystemCapacityOpts
	} else {
		opt = append(opt, newFilesystemCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return FilesystemCapacity{noop.Int64UpDownCounter{}}, err
	}
	return FilesystemCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (FilesystemCapacity) Name() string {
	return "container.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemCapacity) Description() string {
	return "Container filesystem capacity."
}

// Add adds incr to the existing count for attrs.
//
// In K8s, this metric is derived from the
// [FsStats.CapacityBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// In K8s, this metric is derived from the
// [FsStats.CapacityBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// FilesystemCapacityObservable is an instrument used to record metric values
// conforming to the "container.filesystem.capacity" semantic conventions. It
// represents the container filesystem capacity.
type FilesystemCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newFilesystemCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container filesystem capacity."),
	metric.WithUnit("By"),
}

// NewFilesystemCapacityObservable returns a new FilesystemCapacityObservable
// instrument.
func NewFilesystemCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (FilesystemCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newFilesystemCapacityObservableOpts
	} else {
		opt = append(opt, newFilesystemCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return FilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return FilesystemCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (FilesystemCapacityObservable) Name() string {
	return "container.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemCapacityObservable) Description() string {
	return "Container filesystem capacity."
}

// FilesystemUsage is an instrument used to record metric values conforming to
// the "container.filesystem.usage" semantic conventions. It represents the
// container filesystem usage.
type FilesystemUsage struct {
	metric.Int64UpDownCounter
}

var newFilesystemUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container filesystem usage."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newFilesystemUsageOpts
	} else {
		opt = append(opt, newFilesystemUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.filesystem.usage",
		opt...,
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
	return "container.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemUsage) Description() string {
	return "Container filesystem usage."
}

// Add adds incr to the existing count for attrs.
//
// This may not equal capacity - available.
//
// In K8s, this metric is derived from the
// [FsStats.UsedBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// This may not equal capacity - available.
//
// In K8s, this metric is derived from the
// [FsStats.UsedBytes] field
// of the [ContainerStats.Rootfs]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [ContainerStats.Rootfs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#ContainerStats
func (m FilesystemUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// FilesystemUsageObservable is an instrument used to record metric values
// conforming to the "container.filesystem.usage" semantic conventions. It
// represents the container filesystem usage.
type FilesystemUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newFilesystemUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container filesystem usage."),
	metric.WithUnit("By"),
}

// NewFilesystemUsageObservable returns a new FilesystemUsageObservable
// instrument.
func NewFilesystemUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (FilesystemUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newFilesystemUsageObservableOpts
	} else {
		opt = append(opt, newFilesystemUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.filesystem.usage",
		opt...,
	)
	if err != nil {
		return FilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return FilesystemUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FilesystemUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (FilesystemUsageObservable) Name() string {
	return "container.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (FilesystemUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (FilesystemUsageObservable) Description() string {
	return "Container filesystem usage."
}

// MemoryAvailable is an instrument used to record metric values conforming to
// the "container.memory.available" semantic conventions. It represents the
// container memory available.
type MemoryAvailable struct {
	metric.Int64UpDownCounter
}

var newMemoryAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container memory available."),
	metric.WithUnit("By"),
}

// NewMemoryAvailable returns a new MemoryAvailable instrument.
func NewMemoryAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryAvailableOpts
	} else {
		opt = append(opt, newMemoryAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.memory.available",
		opt...,
	)
	if err != nil {
		return MemoryAvailable{noop.Int64UpDownCounter{}}, err
	}
	return MemoryAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryAvailable) Name() string {
	return "container.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryAvailable) Description() string {
	return "Container memory available."
}

// Add adds incr to the existing count for attrs.
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// In general, this metric can be derived from [cadvisor] and by subtracting the
// `container_memory_working_set_bytes` metric from the
// `container_spec_memory_limit_bytes` metric.
// In K8s, this metric is derived from the [MemoryStats.AvailableBytes] field of
// the [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// In general, this metric can be derived from [cadvisor] and by subtracting the
// `container_memory_working_set_bytes` metric from the
// `container_spec_memory_limit_bytes` metric.
// In K8s, this metric is derived from the [MemoryStats.AvailableBytes] field of
// the [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// MemoryAvailableObservable is an instrument used to record metric values
// conforming to the "container.memory.available" semantic conventions. It
// represents the container memory available.
type MemoryAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container memory available."),
	metric.WithUnit("By"),
}

// NewMemoryAvailableObservable returns a new MemoryAvailableObservable
// instrument.
func NewMemoryAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryAvailableObservableOpts
	} else {
		opt = append(opt, newMemoryAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.memory.available",
		opt...,
	)
	if err != nil {
		return MemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryAvailableObservable) Name() string {
	return "container.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryAvailableObservable) Description() string {
	return "Container memory available."
}

// MemoryPagingFaults is an instrument used to record metric values conforming to
// the "container.memory.paging.faults" semantic conventions. It represents the
// container memory paging faults.
type MemoryPagingFaults struct {
	metric.Int64Counter
}

var newMemoryPagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Container memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewMemoryPagingFaults returns a new MemoryPagingFaults instrument.
func NewMemoryPagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (MemoryPagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryPagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryPagingFaultsOpts
	} else {
		opt = append(opt, newMemoryPagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"container.memory.paging.faults",
		opt...,
	)
	if err != nil {
		return MemoryPagingFaults{noop.Int64Counter{}}, err
	}
	return MemoryPagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryPagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (MemoryPagingFaults) Name() string {
	return "container.memory.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryPagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (MemoryPagingFaults) Description() string {
	return "Container memory paging faults."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_failures_total{failure_type=pgfault, scope=container}` and
// `container_memory_failures_total{failure_type=pgmajfault, scope=container}`
// metric.
// In K8s, this metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the Kubelet's
// stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryPagingFaults) Add(
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
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_failures_total{failure_type=pgfault, scope=container}` and
// `container_memory_failures_total{failure_type=pgmajfault, scope=container}`
// metric.
// In K8s, this metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the Kubelet's
// stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryPagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (MemoryPagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// MemoryPagingFaultsObservable is an instrument used to record metric values
// conforming to the "container.memory.paging.faults" semantic conventions. It
// represents the container memory paging faults.
type MemoryPagingFaultsObservable struct {
	metric.Int64ObservableCounter
}

var newMemoryPagingFaultsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Container memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewMemoryPagingFaultsObservable returns a new MemoryPagingFaultsObservable
// instrument.
func NewMemoryPagingFaultsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (MemoryPagingFaultsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryPagingFaultsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryPagingFaultsObservableOpts
	} else {
		opt = append(opt, newMemoryPagingFaultsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"container.memory.paging.faults",
		opt...,
	)
	if err != nil {
		return MemoryPagingFaultsObservable{noop.Int64ObservableCounter{}}, err
	}
	return MemoryPagingFaultsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryPagingFaultsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryPagingFaultsObservable) Name() string {
	return "container.memory.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryPagingFaultsObservable) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (MemoryPagingFaultsObservable) Description() string {
	return "Container memory paging faults."
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (MemoryPagingFaultsObservable) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// MemoryRss is an instrument used to record metric values conforming to the
// "container.memory.rss" semantic conventions. It represents the container
// memory RSS.
type MemoryRss struct {
	metric.Int64UpDownCounter
}

var newMemoryRssOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container memory RSS."),
	metric.WithUnit("By"),
}

// NewMemoryRss returns a new MemoryRss instrument.
func NewMemoryRss(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryRss, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryRss{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryRssOpts
	} else {
		opt = append(opt, newMemoryRssOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.memory.rss",
		opt...,
	)
	if err != nil {
		return MemoryRss{noop.Int64UpDownCounter{}}, err
	}
	return MemoryRss{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryRss) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryRss) Name() string {
	return "container.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryRss) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryRss) Description() string {
	return "Container memory RSS."
}

// Add adds incr to the existing count for attrs.
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_rss` metric.
// In K8s, this metric is derived from the [MemoryStats.RSSBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryRss) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_rss` metric.
// In K8s, this metric is derived from the [MemoryStats.RSSBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryRss) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// MemoryRssObservable is an instrument used to record metric values conforming
// to the "container.memory.rss" semantic conventions. It represents the
// container memory RSS.
type MemoryRssObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryRssObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container memory RSS."),
	metric.WithUnit("By"),
}

// NewMemoryRssObservable returns a new MemoryRssObservable instrument.
func NewMemoryRssObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryRssObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryRssObservableOpts
	} else {
		opt = append(opt, newMemoryRssObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.memory.rss",
		opt...,
	)
	if err != nil {
		return MemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryRssObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryRssObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryRssObservable) Name() string {
	return "container.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryRssObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryRssObservable) Description() string {
	return "Container memory RSS."
}

// MemoryUsage is an instrument used to record metric values conforming to the
// "container.memory.usage" semantic conventions. It represents the memory usage
// of the container.
type MemoryUsage struct {
	metric.Int64Counter
}

var newMemoryUsageOpts = []metric.Int64CounterOption{
	metric.WithDescription("Memory usage of the container."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newMemoryUsageOpts
	} else {
		opt = append(opt, newMemoryUsageOpts...)
	}

	i, err := m.Int64Counter(
		"container.memory.usage",
		opt...,
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

// Add adds incr to the existing count for attrs.
//
// Memory usage of the container.
func (m MemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Memory usage of the container.
func (m MemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// MemoryUsageObservable is an instrument used to record metric values conforming
// to the "container.memory.usage" semantic conventions. It represents the memory
// usage of the container.
type MemoryUsageObservable struct {
	metric.Int64ObservableCounter
}

var newMemoryUsageObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Memory usage of the container."),
	metric.WithUnit("By"),
}

// NewMemoryUsageObservable returns a new MemoryUsageObservable instrument.
func NewMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (MemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryUsageObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryUsageObservableOpts
	} else {
		opt = append(opt, newMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"container.memory.usage",
		opt...,
	)
	if err != nil {
		return MemoryUsageObservable{noop.Int64ObservableCounter{}}, err
	}
	return MemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryUsageObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryUsageObservable) Name() string {
	return "container.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryUsageObservable) Description() string {
	return "Memory usage of the container."
}

// MemoryWorkingSet is an instrument used to record metric values conforming to
// the "container.memory.working_set" semantic conventions. It represents the
// container memory working set.
type MemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Container memory working set."),
	metric.WithUnit("By"),
}

// NewMemoryWorkingSet returns a new MemoryWorkingSet instrument.
func NewMemoryWorkingSet(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemoryWorkingSet, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryWorkingSet{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryWorkingSetOpts
	} else {
		opt = append(opt, newMemoryWorkingSetOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"container.memory.working_set",
		opt...,
	)
	if err != nil {
		return MemoryWorkingSet{noop.Int64UpDownCounter{}}, err
	}
	return MemoryWorkingSet{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryWorkingSet) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryWorkingSet) Name() string {
	return "container.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryWorkingSet) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryWorkingSet) Description() string {
	return "Container memory working set."
}

// Add adds incr to the existing count for attrs.
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_working_set_bytes` metric.
// In K8s, this metric is derived from the [MemoryStats.WorkingSetBytes] field of
// the [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
//
// In general, this metric can be derived from [cadvisor] and specifically the
// `container_memory_working_set_bytes` metric.
// In K8s, this metric is derived from the [MemoryStats.WorkingSetBytes] field of
// the [PodStats.Memory] of the Kubelet's stats API.
//
// [cadvisor]: https://github.com/google/cadvisor/blob/v0.53.0/docs/storage/prometheus.md#prometheus-container-metrics
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m MemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// MemoryWorkingSetObservable is an instrument used to record metric values
// conforming to the "container.memory.working_set" semantic conventions. It
// represents the container memory working set.
type MemoryWorkingSetObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemoryWorkingSetObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Container memory working set."),
	metric.WithUnit("By"),
}

// NewMemoryWorkingSetObservable returns a new MemoryWorkingSetObservable
// instrument.
func NewMemoryWorkingSetObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemoryWorkingSetObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemoryWorkingSetObservableOpts
	} else {
		opt = append(opt, newMemoryWorkingSetObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"container.memory.working_set",
		opt...,
	)
	if err != nil {
		return MemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemoryWorkingSetObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemoryWorkingSetObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemoryWorkingSetObservable) Name() string {
	return "container.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (MemoryWorkingSetObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemoryWorkingSetObservable) Description() string {
	return "Container memory working set."
}

// NetworkIO is an instrument used to record metric values conforming to the
// "container.network.io" semantic conventions. It represents the network bytes
// for the container.
type NetworkIO struct {
	metric.Int64Counter
}

var newNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes for the container."),
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
		"container.network.io",
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

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// The number of bytes sent/received on all network interfaces by the container.
func (m NetworkIO) Add(
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
//
// The number of bytes sent/received on all network interfaces by the container.
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

// NetworkIOObservable is an instrument used to record metric values conforming
// to the "container.network.io" semantic conventions. It represents the network
// bytes for the container.
type NetworkIOObservable struct {
	metric.Int64ObservableCounter
}

var newNetworkIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Network bytes for the container."),
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
		"container.network.io",
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
	return "container.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIOObservable) Description() string {
	return "Network bytes for the container."
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NetworkIOObservable) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// Uptime is an instrument used to record metric values conforming to the
// "container.uptime" semantic conventions. It represents the time the container
// has been running.
type Uptime struct {
	metric.Float64Gauge
}

var newUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the container has been running."),
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
		"container.uptime",
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
	return "container.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (Uptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (Uptime) Description() string {
	return "The time the container has been running."
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
// the "container.uptime" semantic conventions. It represents the time the
// container has been running.
type UptimeObservable struct {
	metric.Float64ObservableGauge
}

var newUptimeObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The time the container has been running."),
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
		"container.uptime",
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
	return "container.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (UptimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (UptimeObservable) Description() string {
	return "The time the container has been running."
}
