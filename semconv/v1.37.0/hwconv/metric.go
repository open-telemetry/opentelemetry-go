// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "hw" namespace.
package hwconv

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

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the type of error encountered by the component.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// BatteryStateAttr is an attribute conforming to the hw.battery.state semantic
// conventions. It represents the current state of the battery.
type BatteryStateAttr string

var (
	// BatteryStateCharging is the charging.
	BatteryStateCharging BatteryStateAttr = "charging"
	// BatteryStateDischarging is the discharging.
	BatteryStateDischarging BatteryStateAttr = "discharging"
)

// GpuTaskAttr is an attribute conforming to the hw.gpu.task semantic
// conventions. It represents the type of task the GPU is performing.
type GpuTaskAttr string

var (
	// GpuTaskDecoder is the decoder.
	GpuTaskDecoder GpuTaskAttr = "decoder"
	// GpuTaskEncoder is the encoder.
	GpuTaskEncoder GpuTaskAttr = "encoder"
	// GpuTaskGeneral is the general.
	GpuTaskGeneral GpuTaskAttr = "general"
)

// LimitTypeAttr is an attribute conforming to the hw.limit_type semantic
// conventions. It represents the represents battery charge level thresholds
// relevant to device operation and health. Each `limit_type` denotes a specific
// charge limit such as the minimum or maximum optimal charge, the shutdown
// threshold, or energy-saving thresholds. These values are typically provided by
// the hardware or firmware to guide safe and efficient battery usage.
type LimitTypeAttr string

var (
	// LimitTypeCritical is the critical.
	LimitTypeCritical LimitTypeAttr = "critical"
	// LimitTypeDegraded is the degraded.
	LimitTypeDegraded LimitTypeAttr = "degraded"
	// LimitTypeHighCritical is the high Critical.
	LimitTypeHighCritical LimitTypeAttr = "high.critical"
	// LimitTypeHighDegraded is the high Degraded.
	LimitTypeHighDegraded LimitTypeAttr = "high.degraded"
	// LimitTypeLowCritical is the low Critical.
	LimitTypeLowCritical LimitTypeAttr = "low.critical"
	// LimitTypeLowDegraded is the low Degraded.
	LimitTypeLowDegraded LimitTypeAttr = "low.degraded"
	// LimitTypeMax is the maximum.
	LimitTypeMax LimitTypeAttr = "max"
	// LimitTypeThrottled is the throttled.
	LimitTypeThrottled LimitTypeAttr = "throttled"
	// LimitTypeTurbo is the turbo.
	LimitTypeTurbo LimitTypeAttr = "turbo"
)

// LogicalDiskStateAttr is an attribute conforming to the hw.logical_disk.state
// semantic conventions. It represents the state of the logical disk space usage.
type LogicalDiskStateAttr string

var (
	// LogicalDiskStateUsed is the used.
	LogicalDiskStateUsed LogicalDiskStateAttr = "used"
	// LogicalDiskStateFree is the free.
	LogicalDiskStateFree LogicalDiskStateAttr = "free"
)

// PhysicalDiskStateAttr is an attribute conforming to the hw.physical_disk.state
// semantic conventions. It represents the state of the physical disk endurance
// utilization.
type PhysicalDiskStateAttr string

var (
	// PhysicalDiskStateRemaining is the remaining.
	PhysicalDiskStateRemaining PhysicalDiskStateAttr = "remaining"
)

// StateAttr is an attribute conforming to the hw.state semantic conventions. It
// represents the current state of the component.
type StateAttr string

var (
	// StateDegraded is the degraded.
	StateDegraded StateAttr = "degraded"
	// StateFailed is the failed.
	StateFailed StateAttr = "failed"
	// StateNeedsCleaning is the needs Cleaning.
	StateNeedsCleaning StateAttr = "needs_cleaning"
	// StateOk is the OK.
	StateOk StateAttr = "ok"
	// StatePredictedFailure is the predicted Failure.
	StatePredictedFailure StateAttr = "predicted_failure"
)

// TapeDriveOperationTypeAttr is an attribute conforming to the
// hw.tape_drive.operation_type semantic conventions. It represents the type of
// tape drive operation.
type TapeDriveOperationTypeAttr string

var (
	// TapeDriveOperationTypeMount is the mount.
	TapeDriveOperationTypeMount TapeDriveOperationTypeAttr = "mount"
	// TapeDriveOperationTypeUnmount is the unmount.
	TapeDriveOperationTypeUnmount TapeDriveOperationTypeAttr = "unmount"
	// TapeDriveOperationTypeClean is the clean.
	TapeDriveOperationTypeClean TapeDriveOperationTypeAttr = "clean"
)

// TypeAttr is an attribute conforming to the hw.type semantic conventions. It
// represents the type of the component.
type TypeAttr string

var (
	// TypeBattery is the battery.
	TypeBattery TypeAttr = "battery"
	// TypeCPU is the CPU.
	TypeCPU TypeAttr = "cpu"
	// TypeDiskController is the disk controller.
	TypeDiskController TypeAttr = "disk_controller"
	// TypeEnclosure is the enclosure.
	TypeEnclosure TypeAttr = "enclosure"
	// TypeFan is the fan.
	TypeFan TypeAttr = "fan"
	// TypeGpu is the GPU.
	TypeGpu TypeAttr = "gpu"
	// TypeLogicalDisk is the logical disk.
	TypeLogicalDisk TypeAttr = "logical_disk"
	// TypeMemory is the memory.
	TypeMemory TypeAttr = "memory"
	// TypeNetwork is the network.
	TypeNetwork TypeAttr = "network"
	// TypePhysicalDisk is the physical disk.
	TypePhysicalDisk TypeAttr = "physical_disk"
	// TypePowerSupply is the power supply.
	TypePowerSupply TypeAttr = "power_supply"
	// TypeTapeDrive is the tape drive.
	TypeTapeDrive TypeAttr = "tape_drive"
	// TypeTemperature is the temperature.
	TypeTemperature TypeAttr = "temperature"
	// TypeVoltage is the voltage.
	TypeVoltage TypeAttr = "voltage"
)

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the direction of network traffic for
// network errors.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the standardized value "transmit" of
	// NetworkIODirectionAttr.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the standardized value "receive" of
	// NetworkIODirectionAttr.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
)

// BatteryCharge is an instrument used to record metric values conforming to the
// "hw.battery.charge" semantic conventions. It represents the remaining fraction
// of battery charge.
type BatteryCharge struct {
	metric.Int64Gauge
}

// NewBatteryCharge returns a new BatteryCharge instrument.
func NewBatteryCharge(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (BatteryCharge, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryCharge{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.battery.charge",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Remaining fraction of battery charge."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return BatteryCharge{noop.Int64Gauge{}}, err
	}
	return BatteryCharge{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryCharge) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryCharge) Name() string {
	return "hw.battery.charge"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryCharge) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (BatteryCharge) Description() string {
	return "Remaining fraction of battery charge."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m BatteryCharge) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryCharge) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Amper-hours.
func (BatteryCharge) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryCharge) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryCharge) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryCharge) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryCharge) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryCharge) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// BatteryChargeLimit is an instrument used to record metric values conforming to
// the "hw.battery.charge.limit" semantic conventions. It represents the lower
// limit of battery charge fraction to ensure proper operation.
type BatteryChargeLimit struct {
	metric.Int64Gauge
}

// NewBatteryChargeLimit returns a new BatteryChargeLimit instrument.
func NewBatteryChargeLimit(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (BatteryChargeLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryChargeLimit{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.battery.charge.limit",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Lower limit of battery charge fraction to ensure proper operation."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return BatteryChargeLimit{noop.Int64Gauge{}}, err
	}
	return BatteryChargeLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryChargeLimit) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryChargeLimit) Name() string {
	return "hw.battery.charge.limit"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryChargeLimit) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (BatteryChargeLimit) Description() string {
	return "Lower limit of battery charge fraction to ensure proper operation."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m BatteryChargeLimit) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryChargeLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Amper-hours.
func (BatteryChargeLimit) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryChargeLimit) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the represents battery charge level thresholds
// relevant to device operation and health. Each `limit_type` denotes a specific
// charge limit such as the minimum or maximum optimal charge, the shutdown
// threshold, or energy-saving thresholds. These values are typically provided by
// the hardware or firmware to guide safe and efficient battery usage.
func (BatteryChargeLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryChargeLimit) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryChargeLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryChargeLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryChargeLimit) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// BatteryTimeLeft is an instrument used to record metric values conforming to
// the "hw.battery.time_left" semantic conventions. It represents the time left
// before battery is completely charged or discharged.
type BatteryTimeLeft struct {
	metric.Float64Gauge
}

// NewBatteryTimeLeft returns a new BatteryTimeLeft instrument.
func NewBatteryTimeLeft(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (BatteryTimeLeft, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryTimeLeft{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"hw.battery.time_left",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("Time left before battery is completely charged or discharged."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return BatteryTimeLeft{noop.Float64Gauge{}}, err
	}
	return BatteryTimeLeft{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryTimeLeft) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryTimeLeft) Name() string {
	return "hw.battery.time_left"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryTimeLeft) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (BatteryTimeLeft) Description() string {
	return "Time left before battery is completely charged or discharged."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The state is the the current state of the component
//
// All additional attrs passed are included in the recorded value.
func (m BatteryTimeLeft) Record(
	ctx context.Context,
	val float64,
	id string,
	state StateAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.state", string(state)),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryTimeLeft) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// AttrBatteryState returns an optional attribute for the "hw.battery.state"
// semantic convention. It represents the current state of the battery.
func (BatteryTimeLeft) AttrBatteryState(val BatteryStateAttr) attribute.KeyValue {
	return attribute.String("hw.battery.state", string(val))
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Amper-hours.
func (BatteryTimeLeft) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryTimeLeft) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryTimeLeft) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryTimeLeft) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryTimeLeft) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryTimeLeft) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// CPUSpeed is an instrument used to record metric values conforming to the
// "hw.cpu.speed" semantic conventions. It represents the CPU current frequency.
type CPUSpeed struct {
	metric.Int64Gauge
}

// NewCPUSpeed returns a new CPUSpeed instrument.
func NewCPUSpeed(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (CPUSpeed, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUSpeed{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.cpu.speed",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("CPU current frequency."),
			metric.WithUnit("Hz"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUSpeed{noop.Int64Gauge{}}, err
	}
	return CPUSpeed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUSpeed) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (CPUSpeed) Name() string {
	return "hw.cpu.speed"
}

// Unit returns the semantic convention unit of the instrument
func (CPUSpeed) Unit() string {
	return "Hz"
}

// Description returns the semantic convention description of the instrument
func (CPUSpeed) Description() string {
	return "CPU current frequency."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m CPUSpeed) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m CPUSpeed) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (CPUSpeed) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (CPUSpeed) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (CPUSpeed) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (CPUSpeed) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// CPUSpeedLimit is an instrument used to record metric values conforming to the
// "hw.cpu.speed.limit" semantic conventions. It represents the CPU maximum
// frequency.
type CPUSpeedLimit struct {
	metric.Int64Gauge
}

// NewCPUSpeedLimit returns a new CPUSpeedLimit instrument.
func NewCPUSpeedLimit(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (CPUSpeedLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUSpeedLimit{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.cpu.speed.limit",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("CPU maximum frequency."),
			metric.WithUnit("Hz"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUSpeedLimit{noop.Int64Gauge{}}, err
	}
	return CPUSpeedLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUSpeedLimit) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (CPUSpeedLimit) Name() string {
	return "hw.cpu.speed.limit"
}

// Unit returns the semantic convention unit of the instrument
func (CPUSpeedLimit) Unit() string {
	return "Hz"
}

// Description returns the semantic convention description of the instrument
func (CPUSpeedLimit) Description() string {
	return "CPU maximum frequency."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m CPUSpeedLimit) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m CPUSpeedLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (CPUSpeedLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (CPUSpeedLimit) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (CPUSpeedLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (CPUSpeedLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (CPUSpeedLimit) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Energy is an instrument used to record metric values conforming to the
// "hw.energy" semantic conventions. It represents the energy consumed by the
// component.
type Energy struct {
	metric.Int64Counter
}

// NewEnergy returns a new Energy instrument.
func NewEnergy(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Energy, error) {
	// Check if the meter is nil.
	if m == nil {
		return Energy{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"hw.energy",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Energy consumed by the component."),
			metric.WithUnit("J"),
		}, opt...)...,
	)
	if err != nil {
	    return Energy{noop.Int64Counter{}}, err
	}
	return Energy{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Energy) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Energy) Name() string {
	return "hw.energy"
}

// Unit returns the semantic convention unit of the instrument
func (Energy) Unit() string {
	return "J"
}

// Description returns the semantic convention description of the instrument
func (Energy) Description() string {
	return "Energy consumed by the component."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
func (m Energy) Add(
	ctx context.Context,
	incr int64,
	id string,
	hwType TypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Energy) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Energy) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Energy) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// Errors is an instrument used to record metric values conforming to the
// "hw.errors" semantic conventions. It represents the number of errors
// encountered by the component.
type Errors struct {
	metric.Int64Counter
}

// NewErrors returns a new Errors instrument.
func NewErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Errors, error) {
	// Check if the meter is nil.
	if m == nil {
		return Errors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"hw.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of errors encountered by the component."),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return Errors{noop.Int64Counter{}}, err
	}
	return Errors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Errors) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Errors) Name() string {
	return "hw.errors"
}

// Unit returns the semantic convention unit of the instrument
func (Errors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (Errors) Description() string {
	return "Number of errors encountered by the component."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
func (m Errors) Add(
	ctx context.Context,
	incr int64,
	id string,
	hwType TypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Errors) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the type of error encountered by the component.
func (Errors) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Errors) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Errors) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// network traffic for network errors.
func (Errors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// FanSpeed is an instrument used to record metric values conforming to the
// "hw.fan.speed" semantic conventions. It represents the fan speed in
// revolutions per minute.
type FanSpeed struct {
	metric.Int64Gauge
}

// NewFanSpeed returns a new FanSpeed instrument.
func NewFanSpeed(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (FanSpeed, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeed{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Fan speed in revolutions per minute."),
			metric.WithUnit("rpm"),
		}, opt...)...,
	)
	if err != nil {
	    return FanSpeed{noop.Int64Gauge{}}, err
	}
	return FanSpeed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeed) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeed) Name() string {
	return "hw.fan.speed"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeed) Unit() string {
	return "rpm"
}

// Description returns the semantic convention description of the instrument
func (FanSpeed) Description() string {
	return "Fan speed in revolutions per minute."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m FanSpeed) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeed) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeed) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeed) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeed) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// FanSpeedLimit is an instrument used to record metric values conforming to the
// "hw.fan.speed.limit" semantic conventions. It represents the speed limit in
// rpm.
type FanSpeedLimit struct {
	metric.Int64Gauge
}

// NewFanSpeedLimit returns a new FanSpeedLimit instrument.
func NewFanSpeedLimit(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (FanSpeedLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeedLimit{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed.limit",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Speed limit in rpm."),
			metric.WithUnit("rpm"),
		}, opt...)...,
	)
	if err != nil {
	    return FanSpeedLimit{noop.Int64Gauge{}}, err
	}
	return FanSpeedLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeedLimit) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeedLimit) Name() string {
	return "hw.fan.speed.limit"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeedLimit) Unit() string {
	return "rpm"
}

// Description returns the semantic convention description of the instrument
func (FanSpeedLimit) Description() string {
	return "Speed limit in rpm."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m FanSpeedLimit) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeedLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (FanSpeedLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeedLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeedLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeedLimit) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// FanSpeedRatio is an instrument used to record metric values conforming to the
// "hw.fan.speed_ratio" semantic conventions. It represents the fan speed
// expressed as a fraction of its maximum speed.
type FanSpeedRatio struct {
	metric.Int64Gauge
}

// NewFanSpeedRatio returns a new FanSpeedRatio instrument.
func NewFanSpeedRatio(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (FanSpeedRatio, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeedRatio{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed_ratio",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Fan speed expressed as a fraction of its maximum speed."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return FanSpeedRatio{noop.Int64Gauge{}}, err
	}
	return FanSpeedRatio{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeedRatio) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeedRatio) Name() string {
	return "hw.fan.speed_ratio"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeedRatio) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (FanSpeedRatio) Description() string {
	return "Fan speed expressed as a fraction of its maximum speed."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m FanSpeedRatio) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeedRatio) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeedRatio) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeedRatio) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeedRatio) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// GpuIO is an instrument used to record metric values conforming to the
// "hw.gpu.io" semantic conventions. It represents the received and transmitted
// bytes by the GPU.
type GpuIO struct {
	metric.Int64Counter
}

// NewGpuIO returns a new GpuIO instrument.
func NewGpuIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (GpuIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuIO{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"hw.gpu.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Received and transmitted bytes by the GPU."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return GpuIO{noop.Int64Counter{}}, err
	}
	return GpuIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (GpuIO) Name() string {
	return "hw.gpu.io"
}

// Unit returns the semantic convention unit of the instrument
func (GpuIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuIO) Description() string {
	return "Received and transmitted bytes by the GPU."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The networkIoDirection is the the network IO operation direction.
//
// All additional attrs passed are included in the recorded value.
func (m GpuIO) Add(
	ctx context.Context,
	incr int64,
	id string,
	networkIoDirection NetworkIODirectionAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuIO) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuIO) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuIO) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuIO) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuIO) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuIO) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuIO) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryLimit is an instrument used to record metric values conforming to the
// "hw.gpu.memory.limit" semantic conventions. It represents the size of the GPU
// memory.
type GpuMemoryLimit struct {
	metric.Int64UpDownCounter
}

// NewGpuMemoryLimit returns a new GpuMemoryLimit instrument.
func NewGpuMemoryLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (GpuMemoryLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.gpu.memory.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Size of the GPU memory."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return GpuMemoryLimit{noop.Int64UpDownCounter{}}, err
	}
	return GpuMemoryLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryLimit) Name() string {
	return "hw.gpu.memory.limit"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryLimit) Description() string {
	return "Size of the GPU memory."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m GpuMemoryLimit) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuMemoryLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryLimit) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryLimit) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryLimit) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryLimit) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryLimit) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryUsage is an instrument used to record metric values conforming to the
// "hw.gpu.memory.usage" semantic conventions. It represents the GPU memory used.
type GpuMemoryUsage struct {
	metric.Int64UpDownCounter
}

// NewGpuMemoryUsage returns a new GpuMemoryUsage instrument.
func NewGpuMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (GpuMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.gpu.memory.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("GPU memory used."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return GpuMemoryUsage{noop.Int64UpDownCounter{}}, err
	}
	return GpuMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryUsage) Name() string {
	return "hw.gpu.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryUsage) Description() string {
	return "GPU memory used."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m GpuMemoryUsage) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuMemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryUsage) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryUsage) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryUsage) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryUsage) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryUsage) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryUsage) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryUsage) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryUtilization is an instrument used to record metric values conforming
// to the "hw.gpu.memory.utilization" semantic conventions. It represents the
// fraction of GPU memory used.
type GpuMemoryUtilization struct {
	metric.Int64Gauge
}

// NewGpuMemoryUtilization returns a new GpuMemoryUtilization instrument.
func NewGpuMemoryUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (GpuMemoryUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.gpu.memory.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Fraction of GPU memory used."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return GpuMemoryUtilization{noop.Int64Gauge{}}, err
	}
	return GpuMemoryUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryUtilization) Name() string {
	return "hw.gpu.memory.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryUtilization) Description() string {
	return "Fraction of GPU memory used."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m GpuMemoryUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m GpuMemoryUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryUtilization) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryUtilization) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryUtilization) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryUtilization) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryUtilization) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuUtilization is an instrument used to record metric values conforming to the
// "hw.gpu.utilization" semantic conventions. It represents the fraction of time
// spent in a specific task.
type GpuUtilization struct {
	metric.Int64Gauge
}

// NewGpuUtilization returns a new GpuUtilization instrument.
func NewGpuUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (GpuUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.gpu.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Fraction of time spent in a specific task."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return GpuUtilization{noop.Int64Gauge{}}, err
	}
	return GpuUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (GpuUtilization) Name() string {
	return "hw.gpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (GpuUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (GpuUtilization) Description() string {
	return "Fraction of time spent in a specific task."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m GpuUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m GpuUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuUtilization) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuUtilization) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrGpuTask returns an optional attribute for the "hw.gpu.task" semantic
// convention. It represents the type of task the GPU is performing.
func (GpuUtilization) AttrGpuTask(val GpuTaskAttr) attribute.KeyValue {
	return attribute.String("hw.gpu.task", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuUtilization) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuUtilization) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuUtilization) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// HostAmbientTemperature is an instrument used to record metric values
// conforming to the "hw.host.ambient_temperature" semantic conventions. It
// represents the ambient (external) temperature of the physical host.
type HostAmbientTemperature struct {
	metric.Int64Gauge
}

// NewHostAmbientTemperature returns a new HostAmbientTemperature instrument.
func NewHostAmbientTemperature(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HostAmbientTemperature, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostAmbientTemperature{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.host.ambient_temperature",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Ambient (external) temperature of the physical host."),
			metric.WithUnit("Cel"),
		}, opt...)...,
	)
	if err != nil {
	    return HostAmbientTemperature{noop.Int64Gauge{}}, err
	}
	return HostAmbientTemperature{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostAmbientTemperature) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HostAmbientTemperature) Name() string {
	return "hw.host.ambient_temperature"
}

// Unit returns the semantic convention unit of the instrument
func (HostAmbientTemperature) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (HostAmbientTemperature) Description() string {
	return "Ambient (external) temperature of the physical host."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostAmbientTemperature) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m HostAmbientTemperature) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostAmbientTemperature) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostAmbientTemperature) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostEnergy is an instrument used to record metric values conforming to the
// "hw.host.energy" semantic conventions. It represents the total energy consumed
// by the entire physical host, in joules.
type HostEnergy struct {
	metric.Int64Counter
}

// NewHostEnergy returns a new HostEnergy instrument.
func NewHostEnergy(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (HostEnergy, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostEnergy{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"hw.host.energy",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Total energy consumed by the entire physical host, in joules."),
			metric.WithUnit("J"),
		}, opt...)...,
	)
	if err != nil {
	    return HostEnergy{noop.Int64Counter{}}, err
	}
	return HostEnergy{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostEnergy) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (HostEnergy) Name() string {
	return "hw.host.energy"
}

// Unit returns the semantic convention unit of the instrument
func (HostEnergy) Unit() string {
	return "J"
}

// Description returns the semantic convention description of the instrument
func (HostEnergy) Description() string {
	return "Total energy consumed by the entire physical host, in joules."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
//
// The overall energy usage of a host MUST be reported using the specific
// `hw.host.energy` and `hw.host.power` metrics **only**, instead of the generic
// `hw.energy` and `hw.power` described in the previous section, to prevent
// summing up overlapping values.
func (m HostEnergy) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// The overall energy usage of a host MUST be reported using the specific
// `hw.host.energy` and `hw.host.power` metrics **only**, instead of the generic
// `hw.energy` and `hw.power` described in the previous section, to prevent
// summing up overlapping values.
func (m HostEnergy) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostEnergy) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostEnergy) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostHeatingMargin is an instrument used to record metric values conforming to
// the "hw.host.heating_margin" semantic conventions. It represents the by how
// many degrees Celsius the temperature of the physical host can be increased,
// before reaching a warning threshold on one of the internal sensors.
type HostHeatingMargin struct {
	metric.Int64Gauge
}

// NewHostHeatingMargin returns a new HostHeatingMargin instrument.
func NewHostHeatingMargin(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HostHeatingMargin, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostHeatingMargin{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.host.heating_margin",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors."),
			metric.WithUnit("Cel"),
		}, opt...)...,
	)
	if err != nil {
	    return HostHeatingMargin{noop.Int64Gauge{}}, err
	}
	return HostHeatingMargin{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostHeatingMargin) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HostHeatingMargin) Name() string {
	return "hw.host.heating_margin"
}

// Unit returns the semantic convention unit of the instrument
func (HostHeatingMargin) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (HostHeatingMargin) Description() string {
	return "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostHeatingMargin) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m HostHeatingMargin) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostHeatingMargin) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostHeatingMargin) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostPower is an instrument used to record metric values conforming to the
// "hw.host.power" semantic conventions. It represents the instantaneous power
// consumed by the entire physical host in Watts (`hw.host.energy` is preferred).
type HostPower struct {
	metric.Int64Gauge
}

// NewHostPower returns a new HostPower instrument.
func NewHostPower(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HostPower, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostPower{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.host.power",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)."),
			metric.WithUnit("W"),
		}, opt...)...,
	)
	if err != nil {
	    return HostPower{noop.Int64Gauge{}}, err
	}
	return HostPower{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostPower) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HostPower) Name() string {
	return "hw.host.power"
}

// Unit returns the semantic convention unit of the instrument
func (HostPower) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (HostPower) Description() string {
	return "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
//
// The overall energy usage of a host MUST be reported using the specific
// `hw.host.energy` and `hw.host.power` metrics **only**, instead of the generic
// `hw.energy` and `hw.power` described in the previous section, to prevent
// summing up overlapping values.
func (m HostPower) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The overall energy usage of a host MUST be reported using the specific
// `hw.host.energy` and `hw.host.power` metrics **only**, instead of the generic
// `hw.energy` and `hw.power` described in the previous section, to prevent
// summing up overlapping values.
func (m HostPower) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostPower) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostPower) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskLimit is an instrument used to record metric values conforming to
// the "hw.logical_disk.limit" semantic conventions. It represents the size of
// the logical disk.
type LogicalDiskLimit struct {
	metric.Int64UpDownCounter
}

// NewLogicalDiskLimit returns a new LogicalDiskLimit instrument.
func NewLogicalDiskLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (LogicalDiskLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.logical_disk.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Size of the logical disk."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return LogicalDiskLimit{noop.Int64UpDownCounter{}}, err
	}
	return LogicalDiskLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskLimit) Name() string {
	return "hw.logical_disk.limit"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskLimit) Description() string {
	return "Size of the logical disk."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m LogicalDiskLimit) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m LogicalDiskLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskLimit) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskUsage is an instrument used to record metric values conforming to
// the "hw.logical_disk.usage" semantic conventions. It represents the logical
// disk space usage.
type LogicalDiskUsage struct {
	metric.Int64UpDownCounter
}

// NewLogicalDiskUsage returns a new LogicalDiskUsage instrument.
func NewLogicalDiskUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (LogicalDiskUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.logical_disk.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Logical disk space usage."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return LogicalDiskUsage{noop.Int64UpDownCounter{}}, err
	}
	return LogicalDiskUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskUsage) Name() string {
	return "hw.logical_disk.usage"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskUsage) Description() string {
	return "Logical disk space usage."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The logicalDiskState is the state of the logical disk space usage
//
// All additional attrs passed are included in the recorded value.
func (m LogicalDiskUsage) Add(
	ctx context.Context,
	incr int64,
	id string,
	logicalDiskState LogicalDiskStateAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.logical_disk.state", string(logicalDiskState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m LogicalDiskUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskUsage) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskUsage) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskUsage) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskUtilization is an instrument used to record metric values
// conforming to the "hw.logical_disk.utilization" semantic conventions. It
// represents the logical disk space utilization as a fraction.
type LogicalDiskUtilization struct {
	metric.Int64Gauge
}

// NewLogicalDiskUtilization returns a new LogicalDiskUtilization instrument.
func NewLogicalDiskUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (LogicalDiskUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.logical_disk.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Logical disk space utilization as a fraction."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return LogicalDiskUtilization{noop.Int64Gauge{}}, err
	}
	return LogicalDiskUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskUtilization) Name() string {
	return "hw.logical_disk.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskUtilization) Description() string {
	return "Logical disk space utilization as a fraction."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The logicalDiskState is the state of the logical disk space usage
//
// All additional attrs passed are included in the recorded value.
func (m LogicalDiskUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	logicalDiskState LogicalDiskStateAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.logical_disk.state", string(logicalDiskState)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m LogicalDiskUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskUtilization) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// MemorySize is an instrument used to record metric values conforming to the
// "hw.memory.size" semantic conventions. It represents the size of the memory
// module.
type MemorySize struct {
	metric.Int64UpDownCounter
}

// NewMemorySize returns a new MemorySize instrument.
func NewMemorySize(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (MemorySize, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemorySize{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.memory.size",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Size of the memory module."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return MemorySize{noop.Int64UpDownCounter{}}, err
	}
	return MemorySize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemorySize) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemorySize) Name() string {
	return "hw.memory.size"
}

// Unit returns the semantic convention unit of the instrument
func (MemorySize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemorySize) Description() string {
	return "Size of the memory module."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m MemorySize) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m MemorySize) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrMemoryType returns an optional attribute for the "hw.memory.type" semantic
// convention. It represents the type of the memory module.
func (MemorySize) AttrMemoryType(val string) attribute.KeyValue {
	return attribute.String("hw.memory.type", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (MemorySize) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (MemorySize) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (MemorySize) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (MemorySize) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (MemorySize) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkBandwidthLimit is an instrument used to record metric values conforming
// to the "hw.network.bandwidth.limit" semantic conventions. It represents the
// link speed.
type NetworkBandwidthLimit struct {
	metric.Int64UpDownCounter
}

// NewNetworkBandwidthLimit returns a new NetworkBandwidthLimit instrument.
func NewNetworkBandwidthLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NetworkBandwidthLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkBandwidthLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.network.bandwidth.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Link speed."),
			metric.WithUnit("By/s"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkBandwidthLimit{noop.Int64UpDownCounter{}}, err
	}
	return NetworkBandwidthLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkBandwidthLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkBandwidthLimit) Name() string {
	return "hw.network.bandwidth.limit"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkBandwidthLimit) Unit() string {
	return "By/s"
}

// Description returns the semantic convention description of the instrument
func (NetworkBandwidthLimit) Description() string {
	return "Link speed."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m NetworkBandwidthLimit) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkBandwidthLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkBandwidthLimit) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkBandwidthLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkBandwidthLimit) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkBandwidthLimit) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkBandwidthLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkBandwidthLimit) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkBandwidthLimit) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkBandwidthUtilization is an instrument used to record metric values
// conforming to the "hw.network.bandwidth.utilization" semantic conventions. It
// represents the utilization of the network bandwidth as a fraction.
type NetworkBandwidthUtilization struct {
	metric.Int64Gauge
}

// NewNetworkBandwidthUtilization returns a new NetworkBandwidthUtilization
// instrument.
func NewNetworkBandwidthUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (NetworkBandwidthUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkBandwidthUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.network.bandwidth.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Utilization of the network bandwidth as a fraction."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkBandwidthUtilization{noop.Int64Gauge{}}, err
	}
	return NetworkBandwidthUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkBandwidthUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (NetworkBandwidthUtilization) Name() string {
	return "hw.network.bandwidth.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkBandwidthUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (NetworkBandwidthUtilization) Description() string {
	return "Utilization of the network bandwidth as a fraction."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m NetworkBandwidthUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m NetworkBandwidthUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkBandwidthUtilization) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkBandwidthUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkBandwidthUtilization) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkBandwidthUtilization) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkBandwidthUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkBandwidthUtilization) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkBandwidthUtilization) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkIO is an instrument used to record metric values conforming to the
// "hw.network.io" semantic conventions. It represents the received and
// transmitted network traffic in bytes.
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
		"hw.network.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Received and transmitted network traffic in bytes."),
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
	return "hw.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIO) Description() string {
	return "Received and transmitted network traffic in bytes."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The networkIoDirection is the the network IO operation direction.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkIO) Add(
	ctx context.Context,
	incr int64,
	id string,
	networkIoDirection NetworkIODirectionAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkIO) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkIO) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkIO) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkIO) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkIO) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkIO) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkIO) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkPackets is an instrument used to record metric values conforming to the
// "hw.network.packets" semantic conventions. It represents the received and
// transmitted network traffic in packets (or frames).
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
		"hw.network.packets",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Received and transmitted network traffic in packets (or frames)."),
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
	return "hw.network.packets"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkPackets) Unit() string {
	return "{packet}"
}

// Description returns the semantic convention description of the instrument
func (NetworkPackets) Description() string {
	return "Received and transmitted network traffic in packets (or frames)."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The networkIoDirection is the the network IO operation direction.
//
// All additional attrs passed are included in the recorded value.
func (m NetworkPackets) Add(
	ctx context.Context,
	incr int64,
	id string,
	networkIoDirection NetworkIODirectionAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkPackets) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkPackets) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkPackets) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkPackets) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkPackets) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkPackets) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkPackets) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkPackets) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkUp is an instrument used to record metric values conforming to the
// "hw.network.up" semantic conventions. It represents the link status: `1` (up)
// or `0` (down).
type NetworkUp struct {
	metric.Int64UpDownCounter
}

// NewNetworkUp returns a new NetworkUp instrument.
func NewNetworkUp(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NetworkUp, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkUp{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.network.up",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Link status: `1` (up) or `0` (down)."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return NetworkUp{noop.Int64UpDownCounter{}}, err
	}
	return NetworkUp{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkUp) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkUp) Name() string {
	return "hw.network.up"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkUp) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (NetworkUp) Description() string {
	return "Link status: `1` (up) or `0` (down)."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m NetworkUp) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkUp) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkUp) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkUp) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkUp) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkUp) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkUp) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkUp) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkUp) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PhysicalDiskEnduranceUtilization is an instrument used to record metric values
// conforming to the "hw.physical_disk.endurance_utilization" semantic
// conventions. It represents the endurance remaining for this SSD disk.
type PhysicalDiskEnduranceUtilization struct {
	metric.Int64Gauge
}

// NewPhysicalDiskEnduranceUtilization returns a new
// PhysicalDiskEnduranceUtilization instrument.
func NewPhysicalDiskEnduranceUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PhysicalDiskEnduranceUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskEnduranceUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.physical_disk.endurance_utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Endurance remaining for this SSD disk."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return PhysicalDiskEnduranceUtilization{noop.Int64Gauge{}}, err
	}
	return PhysicalDiskEnduranceUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskEnduranceUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskEnduranceUtilization) Name() string {
	return "hw.physical_disk.endurance_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskEnduranceUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskEnduranceUtilization) Description() string {
	return "Endurance remaining for this SSD disk."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The physicalDiskState is the state of the physical disk endurance utilization
//
// All additional attrs passed are included in the recorded value.
func (m PhysicalDiskEnduranceUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	physicalDiskState PhysicalDiskStateAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.physical_disk.state", string(physicalDiskState)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PhysicalDiskEnduranceUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskEnduranceUtilization) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskEnduranceUtilization) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskEnduranceUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskEnduranceUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskEnduranceUtilization) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskEnduranceUtilization) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskEnduranceUtilization) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PhysicalDiskSize is an instrument used to record metric values conforming to
// the "hw.physical_disk.size" semantic conventions. It represents the size of
// the disk.
type PhysicalDiskSize struct {
	metric.Int64UpDownCounter
}

// NewPhysicalDiskSize returns a new PhysicalDiskSize instrument.
func NewPhysicalDiskSize(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PhysicalDiskSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskSize{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.physical_disk.size",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Size of the disk."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return PhysicalDiskSize{noop.Int64UpDownCounter{}}, err
	}
	return PhysicalDiskSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskSize) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskSize) Name() string {
	return "hw.physical_disk.size"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskSize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskSize) Description() string {
	return "Size of the disk."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m PhysicalDiskSize) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PhysicalDiskSize) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskSize) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskSize) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskSize) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskSize) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskSize) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskSize) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskSize) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PhysicalDiskSmart is an instrument used to record metric values conforming to
// the "hw.physical_disk.smart" semantic conventions. It represents the value of
// the corresponding [S.M.A.R.T.] (Self-Monitoring, Analysis, and Reporting
// Technology) attribute.
//
// [S.M.A.R.T.]: https://wikipedia.org/wiki/S.M.A.R.T.
type PhysicalDiskSmart struct {
	metric.Int64Gauge
}

// NewPhysicalDiskSmart returns a new PhysicalDiskSmart instrument.
func NewPhysicalDiskSmart(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PhysicalDiskSmart, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskSmart{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.physical_disk.smart",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return PhysicalDiskSmart{noop.Int64Gauge{}}, err
	}
	return PhysicalDiskSmart{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskSmart) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskSmart) Name() string {
	return "hw.physical_disk.smart"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskSmart) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskSmart) Description() string {
	return "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m PhysicalDiskSmart) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PhysicalDiskSmart) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskSmart) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskSmart) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskSmart) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskSmart) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskSmartAttribute returns an optional attribute for the
// "hw.physical_disk.smart_attribute" semantic convention. It represents the
// [S.M.A.R.T.] (Self-Monitoring, Analysis, and Reporting Technology) attribute
// of the physical disk.
//
// [S.M.A.R.T.]: https://wikipedia.org/wiki/S.M.A.R.T.
func (PhysicalDiskSmart) AttrPhysicalDiskSmartAttribute(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.smart_attribute", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskSmart) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskSmart) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskSmart) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Power is an instrument used to record metric values conforming to the
// "hw.power" semantic conventions. It represents the instantaneous power
// consumed by the component.
type Power struct {
	metric.Int64Gauge
}

// NewPower returns a new Power instrument.
func NewPower(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (Power, error) {
	// Check if the meter is nil.
	if m == nil {
		return Power{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.power",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Instantaneous power consumed by the component."),
			metric.WithUnit("W"),
		}, opt...)...,
	)
	if err != nil {
	    return Power{noop.Int64Gauge{}}, err
	}
	return Power{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Power) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (Power) Name() string {
	return "hw.power"
}

// Unit returns the semantic convention unit of the instrument
func (Power) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (Power) Description() string {
	return "Instantaneous power consumed by the component."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
//
// It is recommended to report `hw.energy` instead of `hw.power` when possible.
func (m Power) Record(
	ctx context.Context,
	val int64,
	id string,
	hwType TypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// It is recommended to report `hw.energy` instead of `hw.power` when possible.
func (m Power) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Power) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Power) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// PowerSupplyLimit is an instrument used to record metric values conforming to
// the "hw.power_supply.limit" semantic conventions. It represents the maximum
// power output of the power supply.
type PowerSupplyLimit struct {
	metric.Int64UpDownCounter
}

// NewPowerSupplyLimit returns a new PowerSupplyLimit instrument.
func NewPowerSupplyLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PowerSupplyLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.power_supply.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Maximum power output of the power supply."),
			metric.WithUnit("W"),
		}, opt...)...,
	)
	if err != nil {
	    return PowerSupplyLimit{noop.Int64UpDownCounter{}}, err
	}
	return PowerSupplyLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyLimit) Name() string {
	return "hw.power_supply.limit"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyLimit) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyLimit) Description() string {
	return "Maximum power output of the power supply."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m PowerSupplyLimit) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PowerSupplyLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (PowerSupplyLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyLimit) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyLimit) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyLimit) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PowerSupplyUsage is an instrument used to record metric values conforming to
// the "hw.power_supply.usage" semantic conventions. It represents the current
// power output of the power supply.
type PowerSupplyUsage struct {
	metric.Int64UpDownCounter
}

// NewPowerSupplyUsage returns a new PowerSupplyUsage instrument.
func NewPowerSupplyUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PowerSupplyUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.power_supply.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Current power output of the power supply."),
			metric.WithUnit("W"),
		}, opt...)...,
	)
	if err != nil {
	    return PowerSupplyUsage{noop.Int64UpDownCounter{}}, err
	}
	return PowerSupplyUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyUsage) Name() string {
	return "hw.power_supply.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyUsage) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyUsage) Description() string {
	return "Current power output of the power supply."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m PowerSupplyUsage) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PowerSupplyUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyUsage) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyUsage) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyUsage) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyUsage) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyUsage) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PowerSupplyUtilization is an instrument used to record metric values
// conforming to the "hw.power_supply.utilization" semantic conventions. It
// represents the utilization of the power supply as a fraction of its maximum
// output.
type PowerSupplyUtilization struct {
	metric.Int64Gauge
}

// NewPowerSupplyUtilization returns a new PowerSupplyUtilization instrument.
func NewPowerSupplyUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PowerSupplyUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyUtilization{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.power_supply.utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Utilization of the power supply as a fraction of its maximum output."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return PowerSupplyUtilization{noop.Int64Gauge{}}, err
	}
	return PowerSupplyUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyUtilization) Name() string {
	return "hw.power_supply.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyUtilization) Description() string {
	return "Utilization of the power supply as a fraction of its maximum output."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m PowerSupplyUtilization) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PowerSupplyUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyUtilization) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyUtilization) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyUtilization) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyUtilization) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyUtilization) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Status is an instrument used to record metric values conforming to the
// "hw.status" semantic conventions. It represents the operational status: `1`
// (true) or `0` (false) for each of the possible states.
type Status struct {
	metric.Int64UpDownCounter
}

// NewStatus returns a new Status instrument.
func NewStatus(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (Status, error) {
	// Check if the meter is nil.
	if m == nil {
		return Status{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"hw.status",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Operational status: `1` (true) or `0` (false) for each of the possible states."),
			metric.WithUnit("1"),
		}, opt...)...,
	)
	if err != nil {
	    return Status{noop.Int64UpDownCounter{}}, err
	}
	return Status{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Status) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (Status) Name() string {
	return "hw.status"
}

// Unit returns the semantic convention unit of the instrument
func (Status) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (Status) Description() string {
	return "Operational status: `1` (true) or `0` (false) for each of the possible states."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// The state is the the current state of the component
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
//
// `hw.status` is currently specified as an *UpDownCounter* but would ideally be
// represented using a [*StateSet* as defined in OpenMetrics]. This semantic
// convention will be updated once *StateSet* is specified in OpenTelemetry. This
// planned change is not expected to have any consequence on the way users query
// their timeseries backend to retrieve the values of `hw.status` over time.
//
// [ [*StateSet* as defined in OpenMetrics]: https://github.com/prometheus/OpenMetrics/blob/v1.0.0/specification/OpenMetrics.md#stateset
func (m Status) Add(
	ctx context.Context,
	incr int64,
	id string,
	state StateAttr,
	hwType TypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
				attribute.String("hw.state", string(state)),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// `hw.status` is currently specified as an *UpDownCounter* but would ideally be
// represented using a [*StateSet* as defined in OpenMetrics]. This semantic
// convention will be updated once *StateSet* is specified in OpenTelemetry. This
// planned change is not expected to have any consequence on the way users query
// their timeseries backend to retrieve the values of `hw.status` over time.
//
// [ [*StateSet* as defined in OpenMetrics]: https://github.com/prometheus/OpenMetrics/blob/v1.0.0/specification/OpenMetrics.md#stateset
func (m Status) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Status) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Status) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// TapeDriveOperations is an instrument used to record metric values conforming
// to the "hw.tape_drive.operations" semantic conventions. It represents the
// operations performed by the tape drive.
type TapeDriveOperations struct {
	metric.Int64Counter
}

// NewTapeDriveOperations returns a new TapeDriveOperations instrument.
func NewTapeDriveOperations(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (TapeDriveOperations, error) {
	// Check if the meter is nil.
	if m == nil {
		return TapeDriveOperations{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"hw.tape_drive.operations",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Operations performed by the tape drive."),
			metric.WithUnit("{operation}"),
		}, opt...)...,
	)
	if err != nil {
	    return TapeDriveOperations{noop.Int64Counter{}}, err
	}
	return TapeDriveOperations{i}, nil
}

// Inst returns the underlying metric instrument.
func (m TapeDriveOperations) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (TapeDriveOperations) Name() string {
	return "hw.tape_drive.operations"
}

// Unit returns the semantic convention unit of the instrument
func (TapeDriveOperations) Unit() string {
	return "{operation}"
}

// Description returns the semantic convention description of the instrument
func (TapeDriveOperations) Description() string {
	return "Operations performed by the tape drive."
}

// Add adds incr to the existing count for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m TapeDriveOperations) Add(
	ctx context.Context,
	incr int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m TapeDriveOperations) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (TapeDriveOperations) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (TapeDriveOperations) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (TapeDriveOperations) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (TapeDriveOperations) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrTapeDriveOperationType returns an optional attribute for the
// "hw.tape_drive.operation_type" semantic convention. It represents the type of
// tape drive operation.
func (TapeDriveOperations) AttrTapeDriveOperationType(val TapeDriveOperationTypeAttr) attribute.KeyValue {
	return attribute.String("hw.tape_drive.operation_type", string(val))
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (TapeDriveOperations) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Temperature is an instrument used to record metric values conforming to the
// "hw.temperature" semantic conventions. It represents the temperature in
// degrees Celsius.
type Temperature struct {
	metric.Int64Gauge
}

// NewTemperature returns a new Temperature instrument.
func NewTemperature(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (Temperature, error) {
	// Check if the meter is nil.
	if m == nil {
		return Temperature{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.temperature",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Temperature in degrees Celsius."),
			metric.WithUnit("Cel"),
		}, opt...)...,
	)
	if err != nil {
	    return Temperature{noop.Int64Gauge{}}, err
	}
	return Temperature{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Temperature) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (Temperature) Name() string {
	return "hw.temperature"
}

// Unit returns the semantic convention unit of the instrument
func (Temperature) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (Temperature) Description() string {
	return "Temperature in degrees Celsius."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m Temperature) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m Temperature) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Temperature) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Temperature) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (Temperature) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// TemperatureLimit is an instrument used to record metric values conforming to
// the "hw.temperature.limit" semantic conventions. It represents the temperature
// limit in degrees Celsius.
type TemperatureLimit struct {
	metric.Int64Gauge
}

// NewTemperatureLimit returns a new TemperatureLimit instrument.
func NewTemperatureLimit(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (TemperatureLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return TemperatureLimit{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.temperature.limit",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Temperature limit in degrees Celsius."),
			metric.WithUnit("Cel"),
		}, opt...)...,
	)
	if err != nil {
	    return TemperatureLimit{noop.Int64Gauge{}}, err
	}
	return TemperatureLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m TemperatureLimit) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (TemperatureLimit) Name() string {
	return "hw.temperature.limit"
}

// Unit returns the semantic convention unit of the instrument
func (TemperatureLimit) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (TemperatureLimit) Description() string {
	return "Temperature limit in degrees Celsius."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m TemperatureLimit) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m TemperatureLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (TemperatureLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (TemperatureLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (TemperatureLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (TemperatureLimit) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// Voltage is an instrument used to record metric values conforming to the
// "hw.voltage" semantic conventions. It represents the voltage measured by the
// sensor.
type Voltage struct {
	metric.Int64Gauge
}

// NewVoltage returns a new Voltage instrument.
func NewVoltage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (Voltage, error) {
	// Check if the meter is nil.
	if m == nil {
		return Voltage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.voltage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Voltage measured by the sensor."),
			metric.WithUnit("V"),
		}, opt...)...,
	)
	if err != nil {
	    return Voltage{noop.Int64Gauge{}}, err
	}
	return Voltage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Voltage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (Voltage) Name() string {
	return "hw.voltage"
}

// Unit returns the semantic convention unit of the instrument
func (Voltage) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (Voltage) Description() string {
	return "Voltage measured by the sensor."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m Voltage) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m Voltage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (Voltage) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (Voltage) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (Voltage) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// VoltageLimit is an instrument used to record metric values conforming to the
// "hw.voltage.limit" semantic conventions. It represents the voltage limit in
// Volts.
type VoltageLimit struct {
	metric.Int64Gauge
}

// NewVoltageLimit returns a new VoltageLimit instrument.
func NewVoltageLimit(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (VoltageLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return VoltageLimit{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.voltage.limit",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Voltage limit in Volts."),
			metric.WithUnit("V"),
		}, opt...)...,
	)
	if err != nil {
	    return VoltageLimit{noop.Int64Gauge{}}, err
	}
	return VoltageLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m VoltageLimit) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (VoltageLimit) Name() string {
	return "hw.voltage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (VoltageLimit) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (VoltageLimit) Description() string {
	return "Voltage limit in Volts."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m VoltageLimit) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m VoltageLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (VoltageLimit) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (VoltageLimit) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (VoltageLimit) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (VoltageLimit) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// VoltageNominal is an instrument used to record metric values conforming to the
// "hw.voltage.nominal" semantic conventions. It represents the nominal
// (expected) voltage.
type VoltageNominal struct {
	metric.Int64Gauge
}

// NewVoltageNominal returns a new VoltageNominal instrument.
func NewVoltageNominal(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (VoltageNominal, error) {
	// Check if the meter is nil.
	if m == nil {
		return VoltageNominal{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"hw.voltage.nominal",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Nominal (expected) voltage."),
			metric.WithUnit("V"),
		}, opt...)...,
	)
	if err != nil {
	    return VoltageNominal{noop.Int64Gauge{}}, err
	}
	return VoltageNominal{i}, nil
}

// Inst returns the underlying metric instrument.
func (m VoltageNominal) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (VoltageNominal) Name() string {
	return "hw.voltage.nominal"
}

// Unit returns the semantic convention unit of the instrument
func (VoltageNominal) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (VoltageNominal) Description() string {
	return "Nominal (expected) voltage."
}

// Record records val to the current distribution for attrs.
//
// The id is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m VoltageNominal) Record(
	ctx context.Context,
	val int64,
	id string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs,
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m VoltageNominal) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (VoltageNominal) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (VoltageNominal) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (VoltageNominal) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}