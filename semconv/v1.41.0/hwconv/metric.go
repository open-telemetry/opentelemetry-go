// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package hwconv provides types and functionality for OpenTelemetry semantic
// conventions in the "hw" namespace.
package hwconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
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

var newBatteryChargeOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Remaining fraction of battery charge."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newBatteryChargeOpts
	} else {
		opt = append(opt, newBatteryChargeOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.battery.charge",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryCharge) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
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

// BatteryChargeObservable is an instrument used to record metric values
// conforming to the "hw.battery.charge" semantic conventions. It represents the
// remaining fraction of battery charge.
type BatteryChargeObservable struct {
	metric.Int64ObservableGauge
}

var newBatteryChargeObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Remaining fraction of battery charge."),
	metric.WithUnit("1"),
}

// NewBatteryChargeObservable returns a new BatteryChargeObservable instrument.
func NewBatteryChargeObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (BatteryChargeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryChargeObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newBatteryChargeObservableOpts
	} else {
		opt = append(opt, newBatteryChargeObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.battery.charge",
		opt...,
	)
	if err != nil {
		return BatteryChargeObservable{noop.Int64ObservableGauge{}}, err
	}
	return BatteryChargeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryChargeObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryChargeObservable) Name() string {
	return "hw.battery.charge"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryChargeObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (BatteryChargeObservable) Description() string {
	return "Remaining fraction of battery charge."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (BatteryChargeObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
func (BatteryChargeObservable) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryChargeObservable) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryChargeObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryChargeObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryChargeObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryChargeObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// BatteryChargeLimit is an instrument used to record metric values conforming to
// the "hw.battery.charge.limit" semantic conventions. It represents the lower
// limit of battery charge fraction to ensure proper operation.
type BatteryChargeLimit struct {
	metric.Int64Gauge
}

var newBatteryChargeLimitOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Lower limit of battery charge fraction to ensure proper operation."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newBatteryChargeLimitOpts
	} else {
		opt = append(opt, newBatteryChargeLimitOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.battery.charge.limit",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryChargeLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
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

// BatteryChargeLimitObservable is an instrument used to record metric values
// conforming to the "hw.battery.charge.limit" semantic conventions. It
// represents the lower limit of battery charge fraction to ensure proper
// operation.
type BatteryChargeLimitObservable struct {
	metric.Int64ObservableGauge
}

var newBatteryChargeLimitObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Lower limit of battery charge fraction to ensure proper operation."),
	metric.WithUnit("1"),
}

// NewBatteryChargeLimitObservable returns a new BatteryChargeLimitObservable
// instrument.
func NewBatteryChargeLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (BatteryChargeLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryChargeLimitObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newBatteryChargeLimitObservableOpts
	} else {
		opt = append(opt, newBatteryChargeLimitObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.battery.charge.limit",
		opt...,
	)
	if err != nil {
		return BatteryChargeLimitObservable{noop.Int64ObservableGauge{}}, err
	}
	return BatteryChargeLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryChargeLimitObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryChargeLimitObservable) Name() string {
	return "hw.battery.charge.limit"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryChargeLimitObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (BatteryChargeLimitObservable) Description() string {
	return "Lower limit of battery charge fraction to ensure proper operation."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (BatteryChargeLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
func (BatteryChargeLimitObservable) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryChargeLimitObservable) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the represents battery charge level thresholds
// relevant to device operation and health. Each `limit_type` denotes a specific
// charge limit such as the minimum or maximum optimal charge, the shutdown
// threshold, or energy-saving thresholds. These values are typically provided by
// the hardware or firmware to guide safe and efficient battery usage.
func (BatteryChargeLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryChargeLimitObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryChargeLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryChargeLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryChargeLimitObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// BatteryTimeLeft is an instrument used to record metric values conforming to
// the "hw.battery.time_left" semantic conventions. It represents the time left
// before battery is completely charged or discharged.
type BatteryTimeLeft struct {
	metric.Float64Gauge
}

var newBatteryTimeLeftOpts = []metric.Float64GaugeOption{
	metric.WithDescription("Time left before battery is completely charged or discharged."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newBatteryTimeLeftOpts
	} else {
		opt = append(opt, newBatteryTimeLeftOpts...)
	}

	i, err := m.Float64Gauge(
		"hw.battery.time_left",
		opt...,
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
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.state", string(state)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.state", string(state)),
			)...,
		),
	)

	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m BatteryTimeLeft) RecordSet(ctx context.Context, val float64, set attribute.Set) {
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

// AttrBatteryState returns an optional attribute for the "hw.battery.state"
// semantic convention. It represents the current state of the battery.
func (BatteryTimeLeft) AttrBatteryState(val BatteryStateAttr) attribute.KeyValue {
	return attribute.String("hw.battery.state", string(val))
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
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

// BatteryTimeLeftObservable is an instrument used to record metric values
// conforming to the "hw.battery.time_left" semantic conventions. It represents
// the time left before battery is completely charged or discharged.
type BatteryTimeLeftObservable struct {
	metric.Float64ObservableGauge
}

var newBatteryTimeLeftObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("Time left before battery is completely charged or discharged."),
	metric.WithUnit("s"),
}

// NewBatteryTimeLeftObservable returns a new BatteryTimeLeftObservable
// instrument.
func NewBatteryTimeLeftObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (BatteryTimeLeftObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return BatteryTimeLeftObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newBatteryTimeLeftObservableOpts
	} else {
		opt = append(opt, newBatteryTimeLeftObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"hw.battery.time_left",
		opt...,
	)
	if err != nil {
		return BatteryTimeLeftObservable{noop.Float64ObservableGauge{}}, err
	}
	return BatteryTimeLeftObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m BatteryTimeLeftObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (BatteryTimeLeftObservable) Name() string {
	return "hw.battery.time_left"
}

// Unit returns the semantic convention unit of the instrument
func (BatteryTimeLeftObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (BatteryTimeLeftObservable) Description() string {
	return "Time left before battery is completely charged or discharged."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (BatteryTimeLeftObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrState returns a required attribute for the "hw.state" semantic convention.
// It represents the current state of the component.
func (BatteryTimeLeftObservable) AttrState(val StateAttr) attribute.KeyValue {
	return attribute.String("hw.state", string(val))
}

// AttrBatteryState returns an optional attribute for the "hw.battery.state"
// semantic convention. It represents the current state of the battery.
func (BatteryTimeLeftObservable) AttrBatteryState(val BatteryStateAttr) attribute.KeyValue {
	return attribute.String("hw.battery.state", string(val))
}

// AttrBatteryCapacity returns an optional attribute for the
// "hw.battery.capacity" semantic convention. It represents the design capacity
// in Watts-hours or Ampere-hours.
func (BatteryTimeLeftObservable) AttrBatteryCapacity(val string) attribute.KeyValue {
	return attribute.String("hw.battery.capacity", val)
}

// AttrBatteryChemistry returns an optional attribute for the
// "hw.battery.chemistry" semantic convention. It represents the battery
// [chemistry], e.g. Lithium-Ion, Nickel-Cadmium, etc.
//
// [chemistry]: https://schemas.dmtf.org/wbem/cim-html/2.31.0/CIM_Battery.html
func (BatteryTimeLeftObservable) AttrBatteryChemistry(val string) attribute.KeyValue {
	return attribute.String("hw.battery.chemistry", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (BatteryTimeLeftObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (BatteryTimeLeftObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (BatteryTimeLeftObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (BatteryTimeLeftObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// CPUSpeed is an instrument used to record metric values conforming to the
// "hw.cpu.speed" semantic conventions. It represents the CPU current frequency.
type CPUSpeed struct {
	metric.Int64Gauge
}

var newCPUSpeedOpts = []metric.Int64GaugeOption{
	metric.WithDescription("CPU current frequency."),
	metric.WithUnit("Hz"),
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

	if len(opt) == 0 {
		opt = newCPUSpeedOpts
	} else {
		opt = append(opt, newCPUSpeedOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.cpu.speed",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m CPUSpeed) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// CPUSpeedObservable is an instrument used to record metric values conforming to
// the "hw.cpu.speed" semantic conventions. It represents the CPU current
// frequency.
type CPUSpeedObservable struct {
	metric.Int64ObservableGauge
}

var newCPUSpeedObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("CPU current frequency."),
	metric.WithUnit("Hz"),
}

// NewCPUSpeedObservable returns a new CPUSpeedObservable instrument.
func NewCPUSpeedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (CPUSpeedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUSpeedObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUSpeedObservableOpts
	} else {
		opt = append(opt, newCPUSpeedObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.cpu.speed",
		opt...,
	)
	if err != nil {
		return CPUSpeedObservable{noop.Int64ObservableGauge{}}, err
	}
	return CPUSpeedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUSpeedObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (CPUSpeedObservable) Name() string {
	return "hw.cpu.speed"
}

// Unit returns the semantic convention unit of the instrument
func (CPUSpeedObservable) Unit() string {
	return "Hz"
}

// Description returns the semantic convention description of the instrument
func (CPUSpeedObservable) Description() string {
	return "CPU current frequency."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (CPUSpeedObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (CPUSpeedObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (CPUSpeedObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (CPUSpeedObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (CPUSpeedObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// CPUSpeedLimit is an instrument used to record metric values conforming to the
// "hw.cpu.speed.limit" semantic conventions. It represents the CPU maximum
// frequency.
type CPUSpeedLimit struct {
	metric.Int64Gauge
}

var newCPUSpeedLimitOpts = []metric.Int64GaugeOption{
	metric.WithDescription("CPU maximum frequency."),
	metric.WithUnit("Hz"),
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

	if len(opt) == 0 {
		opt = newCPUSpeedLimitOpts
	} else {
		opt = append(opt, newCPUSpeedLimitOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.cpu.speed.limit",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m CPUSpeedLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// CPUSpeedLimitObservable is an instrument used to record metric values
// conforming to the "hw.cpu.speed.limit" semantic conventions. It represents the
// CPU maximum frequency.
type CPUSpeedLimitObservable struct {
	metric.Int64ObservableGauge
}

var newCPUSpeedLimitObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("CPU maximum frequency."),
	metric.WithUnit("Hz"),
}

// NewCPUSpeedLimitObservable returns a new CPUSpeedLimitObservable instrument.
func NewCPUSpeedLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (CPUSpeedLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUSpeedLimitObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newCPUSpeedLimitObservableOpts
	} else {
		opt = append(opt, newCPUSpeedLimitObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.cpu.speed.limit",
		opt...,
	)
	if err != nil {
		return CPUSpeedLimitObservable{noop.Int64ObservableGauge{}}, err
	}
	return CPUSpeedLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUSpeedLimitObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (CPUSpeedLimitObservable) Name() string {
	return "hw.cpu.speed.limit"
}

// Unit returns the semantic convention unit of the instrument
func (CPUSpeedLimitObservable) Unit() string {
	return "Hz"
}

// Description returns the semantic convention description of the instrument
func (CPUSpeedLimitObservable) Description() string {
	return "CPU maximum frequency."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (CPUSpeedLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (CPUSpeedLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (CPUSpeedLimitObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (CPUSpeedLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (CPUSpeedLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (CPUSpeedLimitObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Energy is an instrument used to record metric values conforming to the
// "hw.energy" semantic conventions. It represents the energy consumed by the
// component.
type Energy struct {
	metric.Int64Counter
}

var newEnergyOpts = []metric.Int64CounterOption{
	metric.WithDescription("Energy consumed by the component."),
	metric.WithUnit("J"),
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

	if len(opt) == 0 {
		opt = newEnergyOpts
	} else {
		opt = append(opt, newEnergyOpts...)
	}

	i, err := m.Int64Counter(
		"hw.energy",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.type", string(hwType)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Energy) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// EnergyObservable is an instrument used to record metric values conforming to
// the "hw.energy" semantic conventions. It represents the energy consumed by the
// component.
type EnergyObservable struct {
	metric.Int64ObservableCounter
}

var newEnergyObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Energy consumed by the component."),
	metric.WithUnit("J"),
}

// NewEnergyObservable returns a new EnergyObservable instrument.
func NewEnergyObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (EnergyObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return EnergyObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newEnergyObservableOpts
	} else {
		opt = append(opt, newEnergyObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.energy",
		opt...,
	)
	if err != nil {
		return EnergyObservable{noop.Int64ObservableCounter{}}, err
	}
	return EnergyObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m EnergyObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (EnergyObservable) Name() string {
	return "hw.energy"
}

// Unit returns the semantic convention unit of the instrument
func (EnergyObservable) Unit() string {
	return "J"
}

// Description returns the semantic convention description of the instrument
func (EnergyObservable) Description() string {
	return "Energy consumed by the component."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (EnergyObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrType returns a required attribute for the "hw.type" semantic convention.
// It represents the type of the component.
func (EnergyObservable) AttrType(val TypeAttr) attribute.KeyValue {
	return attribute.String("hw.type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (EnergyObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (EnergyObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// Errors is an instrument used to record metric values conforming to the
// "hw.errors" semantic conventions. It represents the number of errors
// encountered by the component.
type Errors struct {
	metric.Int64Counter
}

var newErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of errors encountered by the component."),
	metric.WithUnit("{error}"),
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

	if len(opt) == 0 {
		opt = newErrorsOpts
	} else {
		opt = append(opt, newErrorsOpts...)
	}

	i, err := m.Int64Counter(
		"hw.errors",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.type", string(hwType)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Errors) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ErrorsObservable is an instrument used to record metric values conforming to
// the "hw.errors" semantic conventions. It represents the number of errors
// encountered by the component.
type ErrorsObservable struct {
	metric.Int64ObservableCounter
}

var newErrorsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Number of errors encountered by the component."),
	metric.WithUnit("{error}"),
}

// NewErrorsObservable returns a new ErrorsObservable instrument.
func NewErrorsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (ErrorsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ErrorsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newErrorsObservableOpts
	} else {
		opt = append(opt, newErrorsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.errors",
		opt...,
	)
	if err != nil {
		return ErrorsObservable{noop.Int64ObservableCounter{}}, err
	}
	return ErrorsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ErrorsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (ErrorsObservable) Name() string {
	return "hw.errors"
}

// Unit returns the semantic convention unit of the instrument
func (ErrorsObservable) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (ErrorsObservable) Description() string {
	return "Number of errors encountered by the component."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (ErrorsObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrType returns a required attribute for the "hw.type" semantic convention.
// It represents the type of the component.
func (ErrorsObservable) AttrType(val TypeAttr) attribute.KeyValue {
	return attribute.String("hw.type", string(val))
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the type of error encountered by the component.
func (ErrorsObservable) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (ErrorsObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (ErrorsObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// network traffic for network errors.
func (ErrorsObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// FanSpeed is an instrument used to record metric values conforming to the
// "hw.fan.speed" semantic conventions. It represents the fan speed in
// revolutions per minute.
type FanSpeed struct {
	metric.Int64Gauge
}

var newFanSpeedOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Fan speed in revolutions per minute."),
	metric.WithUnit("rpm"),
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

	if len(opt) == 0 {
		opt = newFanSpeedOpts
	} else {
		opt = append(opt, newFanSpeedOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeed) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// FanSpeedObservable is an instrument used to record metric values conforming to
// the "hw.fan.speed" semantic conventions. It represents the fan speed in
// revolutions per minute.
type FanSpeedObservable struct {
	metric.Int64ObservableGauge
}

var newFanSpeedObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Fan speed in revolutions per minute."),
	metric.WithUnit("rpm"),
}

// NewFanSpeedObservable returns a new FanSpeedObservable instrument.
func NewFanSpeedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (FanSpeedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeedObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newFanSpeedObservableOpts
	} else {
		opt = append(opt, newFanSpeedObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.fan.speed",
		opt...,
	)
	if err != nil {
		return FanSpeedObservable{noop.Int64ObservableGauge{}}, err
	}
	return FanSpeedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeedObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeedObservable) Name() string {
	return "hw.fan.speed"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeedObservable) Unit() string {
	return "rpm"
}

// Description returns the semantic convention description of the instrument
func (FanSpeedObservable) Description() string {
	return "Fan speed in revolutions per minute."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (FanSpeedObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeedObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeedObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeedObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// FanSpeedLimit is an instrument used to record metric values conforming to the
// "hw.fan.speed.limit" semantic conventions. It represents the speed limit in
// rpm.
type FanSpeedLimit struct {
	metric.Int64Gauge
}

var newFanSpeedLimitOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Speed limit in rpm."),
	metric.WithUnit("rpm"),
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

	if len(opt) == 0 {
		opt = newFanSpeedLimitOpts
	} else {
		opt = append(opt, newFanSpeedLimitOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed.limit",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeedLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// FanSpeedLimitObservable is an instrument used to record metric values
// conforming to the "hw.fan.speed.limit" semantic conventions. It represents the
// speed limit in rpm.
type FanSpeedLimitObservable struct {
	metric.Int64ObservableGauge
}

var newFanSpeedLimitObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Speed limit in rpm."),
	metric.WithUnit("rpm"),
}

// NewFanSpeedLimitObservable returns a new FanSpeedLimitObservable instrument.
func NewFanSpeedLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (FanSpeedLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeedLimitObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newFanSpeedLimitObservableOpts
	} else {
		opt = append(opt, newFanSpeedLimitObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.fan.speed.limit",
		opt...,
	)
	if err != nil {
		return FanSpeedLimitObservable{noop.Int64ObservableGauge{}}, err
	}
	return FanSpeedLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeedLimitObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeedLimitObservable) Name() string {
	return "hw.fan.speed.limit"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeedLimitObservable) Unit() string {
	return "rpm"
}

// Description returns the semantic convention description of the instrument
func (FanSpeedLimitObservable) Description() string {
	return "Speed limit in rpm."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (FanSpeedLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (FanSpeedLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeedLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeedLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeedLimitObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// FanSpeedRatio is an instrument used to record metric values conforming to the
// "hw.fan.speed_ratio" semantic conventions. It represents the fan speed
// expressed as a fraction of its maximum speed.
type FanSpeedRatio struct {
	metric.Int64Gauge
}

var newFanSpeedRatioOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Fan speed expressed as a fraction of its maximum speed."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newFanSpeedRatioOpts
	} else {
		opt = append(opt, newFanSpeedRatioOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.fan.speed_ratio",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m FanSpeedRatio) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// FanSpeedRatioObservable is an instrument used to record metric values
// conforming to the "hw.fan.speed_ratio" semantic conventions. It represents the
// fan speed expressed as a fraction of its maximum speed.
type FanSpeedRatioObservable struct {
	metric.Int64ObservableGauge
}

var newFanSpeedRatioObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Fan speed expressed as a fraction of its maximum speed."),
	metric.WithUnit("1"),
}

// NewFanSpeedRatioObservable returns a new FanSpeedRatioObservable instrument.
func NewFanSpeedRatioObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (FanSpeedRatioObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return FanSpeedRatioObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newFanSpeedRatioObservableOpts
	} else {
		opt = append(opt, newFanSpeedRatioObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.fan.speed_ratio",
		opt...,
	)
	if err != nil {
		return FanSpeedRatioObservable{noop.Int64ObservableGauge{}}, err
	}
	return FanSpeedRatioObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m FanSpeedRatioObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (FanSpeedRatioObservable) Name() string {
	return "hw.fan.speed_ratio"
}

// Unit returns the semantic convention unit of the instrument
func (FanSpeedRatioObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (FanSpeedRatioObservable) Description() string {
	return "Fan speed expressed as a fraction of its maximum speed."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (FanSpeedRatioObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (FanSpeedRatioObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (FanSpeedRatioObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (FanSpeedRatioObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// GpuIO is an instrument used to record metric values conforming to the
// "hw.gpu.io" semantic conventions. It represents the received and transmitted
// bytes by the GPU.
type GpuIO struct {
	metric.Int64Counter
}

var newGpuIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Received and transmitted bytes by the GPU."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newGpuIOOpts
	} else {
		opt = append(opt, newGpuIOOpts...)
	}

	i, err := m.Int64Counter(
		"hw.gpu.io",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// GpuIOObservable is an instrument used to record metric values conforming to
// the "hw.gpu.io" semantic conventions. It represents the received and
// transmitted bytes by the GPU.
type GpuIOObservable struct {
	metric.Int64ObservableCounter
}

var newGpuIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Received and transmitted bytes by the GPU."),
	metric.WithUnit("By"),
}

// NewGpuIOObservable returns a new GpuIOObservable instrument.
func NewGpuIOObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (GpuIOObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuIOObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newGpuIOObservableOpts
	} else {
		opt = append(opt, newGpuIOObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.gpu.io",
		opt...,
	)
	if err != nil {
		return GpuIOObservable{noop.Int64ObservableCounter{}}, err
	}
	return GpuIOObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuIOObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (GpuIOObservable) Name() string {
	return "hw.gpu.io"
}

// Unit returns the semantic convention unit of the instrument
func (GpuIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuIOObservable) Description() string {
	return "Received and transmitted bytes by the GPU."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (GpuIOObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrNetworkIODirection returns a required attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (GpuIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuIOObservable) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuIOObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuIOObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuIOObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuIOObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuIOObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuIOObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryLimit is an instrument used to record metric values conforming to the
// "hw.gpu.memory.limit" semantic conventions. It represents the size of the GPU
// memory.
type GpuMemoryLimit struct {
	metric.Int64UpDownCounter
}

var newGpuMemoryLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Size of the GPU memory."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newGpuMemoryLimitOpts
	} else {
		opt = append(opt, newGpuMemoryLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.gpu.memory.limit",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuMemoryLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// GpuMemoryLimitObservable is an instrument used to record metric values
// conforming to the "hw.gpu.memory.limit" semantic conventions. It represents
// the size of the GPU memory.
type GpuMemoryLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newGpuMemoryLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Size of the GPU memory."),
	metric.WithUnit("By"),
}

// NewGpuMemoryLimitObservable returns a new GpuMemoryLimitObservable instrument.
func NewGpuMemoryLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (GpuMemoryLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newGpuMemoryLimitObservableOpts
	} else {
		opt = append(opt, newGpuMemoryLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.gpu.memory.limit",
		opt...,
	)
	if err != nil {
		return GpuMemoryLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return GpuMemoryLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryLimitObservable) Name() string {
	return "hw.gpu.memory.limit"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryLimitObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryLimitObservable) Description() string {
	return "Size of the GPU memory."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (GpuMemoryLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryLimitObservable) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryLimitObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryLimitObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryLimitObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryLimitObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryUsage is an instrument used to record metric values conforming to the
// "hw.gpu.memory.usage" semantic conventions. It represents the GPU memory used.
type GpuMemoryUsage struct {
	metric.Int64UpDownCounter
}

var newGpuMemoryUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("GPU memory used."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newGpuMemoryUsageOpts
	} else {
		opt = append(opt, newGpuMemoryUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.gpu.memory.usage",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m GpuMemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// GpuMemoryUsageObservable is an instrument used to record metric values
// conforming to the "hw.gpu.memory.usage" semantic conventions. It represents
// the GPU memory used.
type GpuMemoryUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newGpuMemoryUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("GPU memory used."),
	metric.WithUnit("By"),
}

// NewGpuMemoryUsageObservable returns a new GpuMemoryUsageObservable instrument.
func NewGpuMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (GpuMemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newGpuMemoryUsageObservableOpts
	} else {
		opt = append(opt, newGpuMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.gpu.memory.usage",
		opt...,
	)
	if err != nil {
		return GpuMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return GpuMemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryUsageObservable) Name() string {
	return "hw.gpu.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryUsageObservable) Description() string {
	return "GPU memory used."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (GpuMemoryUsageObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryUsageObservable) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryUsageObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryUsageObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryUsageObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryUsageObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryUsageObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryUsageObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuMemoryUtilization is an instrument used to record metric values conforming
// to the "hw.gpu.memory.utilization" semantic conventions. It represents the
// fraction of GPU memory used.
type GpuMemoryUtilization struct {
	metric.Int64Gauge
}

var newGpuMemoryUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Fraction of GPU memory used."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newGpuMemoryUtilizationOpts
	} else {
		opt = append(opt, newGpuMemoryUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.gpu.memory.utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m GpuMemoryUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// GpuMemoryUtilizationObservable is an instrument used to record metric values
// conforming to the "hw.gpu.memory.utilization" semantic conventions. It
// represents the fraction of GPU memory used.
type GpuMemoryUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newGpuMemoryUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Fraction of GPU memory used."),
	metric.WithUnit("1"),
}

// NewGpuMemoryUtilizationObservable returns a new GpuMemoryUtilizationObservable
// instrument.
func NewGpuMemoryUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (GpuMemoryUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuMemoryUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newGpuMemoryUtilizationObservableOpts
	} else {
		opt = append(opt, newGpuMemoryUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.gpu.memory.utilization",
		opt...,
	)
	if err != nil {
		return GpuMemoryUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return GpuMemoryUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuMemoryUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (GpuMemoryUtilizationObservable) Name() string {
	return "hw.gpu.memory.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (GpuMemoryUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (GpuMemoryUtilizationObservable) Description() string {
	return "Fraction of GPU memory used."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (GpuMemoryUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuMemoryUtilizationObservable) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuMemoryUtilizationObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuMemoryUtilizationObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuMemoryUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuMemoryUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuMemoryUtilizationObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuMemoryUtilizationObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// GpuUtilization is an instrument used to record metric values conforming to the
// "hw.gpu.utilization" semantic conventions. It represents the fraction of time
// spent in a specific task.
type GpuUtilization struct {
	metric.Int64Gauge
}

var newGpuUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Fraction of time spent in a specific task."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newGpuUtilizationOpts
	} else {
		opt = append(opt, newGpuUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.gpu.utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m GpuUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// GpuUtilizationObservable is an instrument used to record metric values
// conforming to the "hw.gpu.utilization" semantic conventions. It represents the
// fraction of time spent in a specific task.
type GpuUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newGpuUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Fraction of time spent in a specific task."),
	metric.WithUnit("1"),
}

// NewGpuUtilizationObservable returns a new GpuUtilizationObservable instrument.
func NewGpuUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (GpuUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return GpuUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newGpuUtilizationObservableOpts
	} else {
		opt = append(opt, newGpuUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.gpu.utilization",
		opt...,
	)
	if err != nil {
		return GpuUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return GpuUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m GpuUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (GpuUtilizationObservable) Name() string {
	return "hw.gpu.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (GpuUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (GpuUtilizationObservable) Description() string {
	return "Fraction of time spent in a specific task."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (GpuUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrDriverVersion returns an optional attribute for the "hw.driver_version"
// semantic convention. It represents the driver version for the hardware
// component.
func (GpuUtilizationObservable) AttrDriverVersion(val string) attribute.KeyValue {
	return attribute.String("hw.driver_version", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (GpuUtilizationObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrGpuTask returns an optional attribute for the "hw.gpu.task" semantic
// convention. It represents the type of task the GPU is performing.
func (GpuUtilizationObservable) AttrGpuTask(val GpuTaskAttr) attribute.KeyValue {
	return attribute.String("hw.gpu.task", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (GpuUtilizationObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (GpuUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (GpuUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (GpuUtilizationObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (GpuUtilizationObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// HostAmbientTemperature is an instrument used to record metric values
// conforming to the "hw.host.ambient_temperature" semantic conventions. It
// represents the ambient (external) temperature of the physical host.
type HostAmbientTemperature struct {
	metric.Int64Gauge
}

var newHostAmbientTemperatureOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Ambient (external) temperature of the physical host."),
	metric.WithUnit("Cel"),
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

	if len(opt) == 0 {
		opt = newHostAmbientTemperatureOpts
	} else {
		opt = append(opt, newHostAmbientTemperatureOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.host.ambient_temperature",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m HostAmbientTemperature) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// HostAmbientTemperatureObservable is an instrument used to record metric values
// conforming to the "hw.host.ambient_temperature" semantic conventions. It
// represents the ambient (external) temperature of the physical host.
type HostAmbientTemperatureObservable struct {
	metric.Int64ObservableGauge
}

var newHostAmbientTemperatureObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Ambient (external) temperature of the physical host."),
	metric.WithUnit("Cel"),
}

// NewHostAmbientTemperatureObservable returns a new
// HostAmbientTemperatureObservable instrument.
func NewHostAmbientTemperatureObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HostAmbientTemperatureObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostAmbientTemperatureObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHostAmbientTemperatureObservableOpts
	} else {
		opt = append(opt, newHostAmbientTemperatureObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.host.ambient_temperature",
		opt...,
	)
	if err != nil {
		return HostAmbientTemperatureObservable{noop.Int64ObservableGauge{}}, err
	}
	return HostAmbientTemperatureObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostAmbientTemperatureObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HostAmbientTemperatureObservable) Name() string {
	return "hw.host.ambient_temperature"
}

// Unit returns the semantic convention unit of the instrument
func (HostAmbientTemperatureObservable) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (HostAmbientTemperatureObservable) Description() string {
	return "Ambient (external) temperature of the physical host."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (HostAmbientTemperatureObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostAmbientTemperatureObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostAmbientTemperatureObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostEnergy is an instrument used to record metric values conforming to the
// "hw.host.energy" semantic conventions. It represents the total energy consumed
// by the entire physical host, in joules.
type HostEnergy struct {
	metric.Int64Counter
}

var newHostEnergyOpts = []metric.Int64CounterOption{
	metric.WithDescription("Total energy consumed by the entire physical host, in joules."),
	metric.WithUnit("J"),
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

	if len(opt) == 0 {
		opt = newHostEnergyOpts
	} else {
		opt = append(opt, newHostEnergyOpts...)
	}

	i, err := m.Int64Counter(
		"hw.host.energy",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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

// HostEnergyObservable is an instrument used to record metric values conforming
// to the "hw.host.energy" semantic conventions. It represents the total energy
// consumed by the entire physical host, in joules.
type HostEnergyObservable struct {
	metric.Int64ObservableCounter
}

var newHostEnergyObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Total energy consumed by the entire physical host, in joules."),
	metric.WithUnit("J"),
}

// NewHostEnergyObservable returns a new HostEnergyObservable instrument.
func NewHostEnergyObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (HostEnergyObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostEnergyObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHostEnergyObservableOpts
	} else {
		opt = append(opt, newHostEnergyObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.host.energy",
		opt...,
	)
	if err != nil {
		return HostEnergyObservable{noop.Int64ObservableCounter{}}, err
	}
	return HostEnergyObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostEnergyObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (HostEnergyObservable) Name() string {
	return "hw.host.energy"
}

// Unit returns the semantic convention unit of the instrument
func (HostEnergyObservable) Unit() string {
	return "J"
}

// Description returns the semantic convention description of the instrument
func (HostEnergyObservable) Description() string {
	return "Total energy consumed by the entire physical host, in joules."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (HostEnergyObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostEnergyObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostEnergyObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostHeatingMargin is an instrument used to record metric values conforming to
// the "hw.host.heating_margin" semantic conventions. It represents the by how
// many degrees Celsius the temperature of the physical host can be increased,
// before reaching a warning threshold on one of the internal sensors.
type HostHeatingMargin struct {
	metric.Int64Gauge
}

var newHostHeatingMarginOpts = []metric.Int64GaugeOption{
	metric.WithDescription("By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors."),
	metric.WithUnit("Cel"),
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

	if len(opt) == 0 {
		opt = newHostHeatingMarginOpts
	} else {
		opt = append(opt, newHostHeatingMarginOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.host.heating_margin",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m HostHeatingMargin) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// HostHeatingMarginObservable is an instrument used to record metric values
// conforming to the "hw.host.heating_margin" semantic conventions. It represents
// the by how many degrees Celsius the temperature of the physical host can be
// increased, before reaching a warning threshold on one of the internal sensors.
type HostHeatingMarginObservable struct {
	metric.Int64ObservableGauge
}

var newHostHeatingMarginObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors."),
	metric.WithUnit("Cel"),
}

// NewHostHeatingMarginObservable returns a new HostHeatingMarginObservable
// instrument.
func NewHostHeatingMarginObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HostHeatingMarginObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostHeatingMarginObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHostHeatingMarginObservableOpts
	} else {
		opt = append(opt, newHostHeatingMarginObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.host.heating_margin",
		opt...,
	)
	if err != nil {
		return HostHeatingMarginObservable{noop.Int64ObservableGauge{}}, err
	}
	return HostHeatingMarginObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostHeatingMarginObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HostHeatingMarginObservable) Name() string {
	return "hw.host.heating_margin"
}

// Unit returns the semantic convention unit of the instrument
func (HostHeatingMarginObservable) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (HostHeatingMarginObservable) Description() string {
	return "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (HostHeatingMarginObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostHeatingMarginObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostHeatingMarginObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// HostPower is an instrument used to record metric values conforming to the
// "hw.host.power" semantic conventions. It represents the instantaneous power
// consumed by the entire physical host in Watts (`hw.host.energy` is preferred).
type HostPower struct {
	metric.Int64Gauge
}

var newHostPowerOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)."),
	metric.WithUnit("W"),
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

	if len(opt) == 0 {
		opt = newHostPowerOpts
	} else {
		opt = append(opt, newHostPowerOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.host.power",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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

// HostPowerObservable is an instrument used to record metric values conforming
// to the "hw.host.power" semantic conventions. It represents the instantaneous
// power consumed by the entire physical host in Watts (`hw.host.energy` is
// preferred).
type HostPowerObservable struct {
	metric.Int64ObservableGauge
}

var newHostPowerObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)."),
	metric.WithUnit("W"),
}

// NewHostPowerObservable returns a new HostPowerObservable instrument.
func NewHostPowerObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HostPowerObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HostPowerObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHostPowerObservableOpts
	} else {
		opt = append(opt, newHostPowerObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.host.power",
		opt...,
	)
	if err != nil {
		return HostPowerObservable{noop.Int64ObservableGauge{}}, err
	}
	return HostPowerObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HostPowerObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HostPowerObservable) Name() string {
	return "hw.host.power"
}

// Unit returns the semantic convention unit of the instrument
func (HostPowerObservable) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (HostPowerObservable) Description() string {
	return "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (HostPowerObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (HostPowerObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (HostPowerObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskLimit is an instrument used to record metric values conforming to
// the "hw.logical_disk.limit" semantic conventions. It represents the size of
// the logical disk.
type LogicalDiskLimit struct {
	metric.Int64UpDownCounter
}

var newLogicalDiskLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Size of the logical disk."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newLogicalDiskLimitOpts
	} else {
		opt = append(opt, newLogicalDiskLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.logical_disk.limit",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m LogicalDiskLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// LogicalDiskLimitObservable is an instrument used to record metric values
// conforming to the "hw.logical_disk.limit" semantic conventions. It represents
// the size of the logical disk.
type LogicalDiskLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newLogicalDiskLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Size of the logical disk."),
	metric.WithUnit("By"),
}

// NewLogicalDiskLimitObservable returns a new LogicalDiskLimitObservable
// instrument.
func NewLogicalDiskLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (LogicalDiskLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newLogicalDiskLimitObservableOpts
	} else {
		opt = append(opt, newLogicalDiskLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.logical_disk.limit",
		opt...,
	)
	if err != nil {
		return LogicalDiskLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return LogicalDiskLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskLimitObservable) Name() string {
	return "hw.logical_disk.limit"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskLimitObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskLimitObservable) Description() string {
	return "Size of the logical disk."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (LogicalDiskLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskLimitObservable) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskUsage is an instrument used to record metric values conforming to
// the "hw.logical_disk.usage" semantic conventions. It represents the logical
// disk space usage.
type LogicalDiskUsage struct {
	metric.Int64UpDownCounter
}

var newLogicalDiskUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Logical disk space usage."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newLogicalDiskUsageOpts
	} else {
		opt = append(opt, newLogicalDiskUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.logical_disk.usage",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.logical_disk.state", string(logicalDiskState)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.logical_disk.state", string(logicalDiskState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m LogicalDiskUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// LogicalDiskUsageObservable is an instrument used to record metric values
// conforming to the "hw.logical_disk.usage" semantic conventions. It represents
// the logical disk space usage.
type LogicalDiskUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newLogicalDiskUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Logical disk space usage."),
	metric.WithUnit("By"),
}

// NewLogicalDiskUsageObservable returns a new LogicalDiskUsageObservable
// instrument.
func NewLogicalDiskUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (LogicalDiskUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newLogicalDiskUsageObservableOpts
	} else {
		opt = append(opt, newLogicalDiskUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.logical_disk.usage",
		opt...,
	)
	if err != nil {
		return LogicalDiskUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return LogicalDiskUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskUsageObservable) Name() string {
	return "hw.logical_disk.usage"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskUsageObservable) Description() string {
	return "Logical disk space usage."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (LogicalDiskUsageObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLogicalDiskState returns a required attribute for the
// "hw.logical_disk.state" semantic convention. It represents the state of the
// logical disk space usage.
func (LogicalDiskUsageObservable) AttrLogicalDiskState(val LogicalDiskStateAttr) attribute.KeyValue {
	return attribute.String("hw.logical_disk.state", string(val))
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskUsageObservable) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskUsageObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskUsageObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// LogicalDiskUtilization is an instrument used to record metric values
// conforming to the "hw.logical_disk.utilization" semantic conventions. It
// represents the logical disk space utilization as a fraction.
type LogicalDiskUtilization struct {
	metric.Int64Gauge
}

var newLogicalDiskUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Logical disk space utilization as a fraction."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newLogicalDiskUtilizationOpts
	} else {
		opt = append(opt, newLogicalDiskUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.logical_disk.utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.logical_disk.state", string(logicalDiskState)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.logical_disk.state", string(logicalDiskState)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m LogicalDiskUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// LogicalDiskUtilizationObservable is an instrument used to record metric values
// conforming to the "hw.logical_disk.utilization" semantic conventions. It
// represents the logical disk space utilization as a fraction.
type LogicalDiskUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newLogicalDiskUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Logical disk space utilization as a fraction."),
	metric.WithUnit("1"),
}

// NewLogicalDiskUtilizationObservable returns a new
// LogicalDiskUtilizationObservable instrument.
func NewLogicalDiskUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (LogicalDiskUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return LogicalDiskUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newLogicalDiskUtilizationObservableOpts
	} else {
		opt = append(opt, newLogicalDiskUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.logical_disk.utilization",
		opt...,
	)
	if err != nil {
		return LogicalDiskUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return LogicalDiskUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m LogicalDiskUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (LogicalDiskUtilizationObservable) Name() string {
	return "hw.logical_disk.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (LogicalDiskUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (LogicalDiskUtilizationObservable) Description() string {
	return "Logical disk space utilization as a fraction."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (LogicalDiskUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLogicalDiskState returns a required attribute for the
// "hw.logical_disk.state" semantic convention. It represents the state of the
// logical disk space usage.
func (LogicalDiskUtilizationObservable) AttrLogicalDiskState(val LogicalDiskStateAttr) attribute.KeyValue {
	return attribute.String("hw.logical_disk.state", string(val))
}

// AttrLogicalDiskRaidLevel returns an optional attribute for the
// "hw.logical_disk.raid_level" semantic convention. It represents the RAID Level
// of the logical disk.
func (LogicalDiskUtilizationObservable) AttrLogicalDiskRaidLevel(val string) attribute.KeyValue {
	return attribute.String("hw.logical_disk.raid_level", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (LogicalDiskUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (LogicalDiskUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// MemorySize is an instrument used to record metric values conforming to the
// "hw.memory.size" semantic conventions. It represents the size of the memory
// module.
type MemorySize struct {
	metric.Int64UpDownCounter
}

var newMemorySizeOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Size of the memory module."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newMemorySizeOpts
	} else {
		opt = append(opt, newMemorySizeOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.memory.size",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m MemorySize) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// MemorySizeObservable is an instrument used to record metric values conforming
// to the "hw.memory.size" semantic conventions. It represents the size of the
// memory module.
type MemorySizeObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newMemorySizeObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Size of the memory module."),
	metric.WithUnit("By"),
}

// NewMemorySizeObservable returns a new MemorySizeObservable instrument.
func NewMemorySizeObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (MemorySizeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemorySizeObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newMemorySizeObservableOpts
	} else {
		opt = append(opt, newMemorySizeObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.memory.size",
		opt...,
	)
	if err != nil {
		return MemorySizeObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return MemorySizeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemorySizeObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (MemorySizeObservable) Name() string {
	return "hw.memory.size"
}

// Unit returns the semantic convention unit of the instrument
func (MemorySizeObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemorySizeObservable) Description() string {
	return "Size of the memory module."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (MemorySizeObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrMemoryType returns an optional attribute for the "hw.memory.type" semantic
// convention. It represents the type of the memory module.
func (MemorySizeObservable) AttrMemoryType(val string) attribute.KeyValue {
	return attribute.String("hw.memory.type", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (MemorySizeObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (MemorySizeObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (MemorySizeObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (MemorySizeObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (MemorySizeObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkBandwidthLimit is an instrument used to record metric values conforming
// to the "hw.network.bandwidth.limit" semantic conventions. It represents the
// link speed.
type NetworkBandwidthLimit struct {
	metric.Int64UpDownCounter
}

var newNetworkBandwidthLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Link speed."),
	metric.WithUnit("By/s"),
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

	if len(opt) == 0 {
		opt = newNetworkBandwidthLimitOpts
	} else {
		opt = append(opt, newNetworkBandwidthLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.network.bandwidth.limit",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkBandwidthLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NetworkBandwidthLimitObservable is an instrument used to record metric values
// conforming to the "hw.network.bandwidth.limit" semantic conventions. It
// represents the link speed.
type NetworkBandwidthLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNetworkBandwidthLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Link speed."),
	metric.WithUnit("By/s"),
}

// NewNetworkBandwidthLimitObservable returns a new
// NetworkBandwidthLimitObservable instrument.
func NewNetworkBandwidthLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NetworkBandwidthLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkBandwidthLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkBandwidthLimitObservableOpts
	} else {
		opt = append(opt, newNetworkBandwidthLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.network.bandwidth.limit",
		opt...,
	)
	if err != nil {
		return NetworkBandwidthLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NetworkBandwidthLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkBandwidthLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkBandwidthLimitObservable) Name() string {
	return "hw.network.bandwidth.limit"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkBandwidthLimitObservable) Unit() string {
	return "By/s"
}

// Description returns the semantic convention description of the instrument
func (NetworkBandwidthLimitObservable) Description() string {
	return "Link speed."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (NetworkBandwidthLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkBandwidthLimitObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkBandwidthLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkBandwidthLimitObservable) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkBandwidthLimitObservable) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkBandwidthLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkBandwidthLimitObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkBandwidthLimitObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkBandwidthUtilization is an instrument used to record metric values
// conforming to the "hw.network.bandwidth.utilization" semantic conventions. It
// represents the utilization of the network bandwidth as a fraction.
type NetworkBandwidthUtilization struct {
	metric.Int64Gauge
}

var newNetworkBandwidthUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Utilization of the network bandwidth as a fraction."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newNetworkBandwidthUtilizationOpts
	} else {
		opt = append(opt, newNetworkBandwidthUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.network.bandwidth.utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m NetworkBandwidthUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// NetworkBandwidthUtilizationObservable is an instrument used to record metric
// values conforming to the "hw.network.bandwidth.utilization" semantic
// conventions. It represents the utilization of the network bandwidth as a
// fraction.
type NetworkBandwidthUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newNetworkBandwidthUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Utilization of the network bandwidth as a fraction."),
	metric.WithUnit("1"),
}

// NewNetworkBandwidthUtilizationObservable returns a new
// NetworkBandwidthUtilizationObservable instrument.
func NewNetworkBandwidthUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (NetworkBandwidthUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkBandwidthUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkBandwidthUtilizationObservableOpts
	} else {
		opt = append(opt, newNetworkBandwidthUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.network.bandwidth.utilization",
		opt...,
	)
	if err != nil {
		return NetworkBandwidthUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return NetworkBandwidthUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkBandwidthUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (NetworkBandwidthUtilizationObservable) Name() string {
	return "hw.network.bandwidth.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkBandwidthUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (NetworkBandwidthUtilizationObservable) Description() string {
	return "Utilization of the network bandwidth as a fraction."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (NetworkBandwidthUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkBandwidthUtilizationObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkBandwidthUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkBandwidthUtilizationObservable) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkBandwidthUtilizationObservable) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkBandwidthUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkBandwidthUtilizationObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkBandwidthUtilizationObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkIO is an instrument used to record metric values conforming to the
// "hw.network.io" semantic conventions. It represents the received and
// transmitted network traffic in bytes.
type NetworkIO struct {
	metric.Int64Counter
}

var newNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Received and transmitted network traffic in bytes."),
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
		"hw.network.io",
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
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

// NetworkIOObservable is an instrument used to record metric values conforming
// to the "hw.network.io" semantic conventions. It represents the received and
// transmitted network traffic in bytes.
type NetworkIOObservable struct {
	metric.Int64ObservableCounter
}

var newNetworkIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Received and transmitted network traffic in bytes."),
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
		"hw.network.io",
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
	return "hw.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetworkIOObservable) Description() string {
	return "Received and transmitted network traffic in bytes."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (NetworkIOObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrNetworkIODirection returns a required attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkIOObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkIOObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkIOObservable) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkIOObservable) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkIOObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkIOObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkIOObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkPackets is an instrument used to record metric values conforming to the
// "hw.network.packets" semantic conventions. It represents the received and
// transmitted network traffic in packets (or frames).
type NetworkPackets struct {
	metric.Int64Counter
}

var newNetworkPacketsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Received and transmitted network traffic in packets (or frames)."),
	metric.WithUnit("{packet}"),
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

	if len(opt) == 0 {
		opt = newNetworkPacketsOpts
	} else {
		opt = append(opt, newNetworkPacketsOpts...)
	}

	i, err := m.Int64Counter(
		"hw.network.packets",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
				attribute.String("network.io.direction", string(networkIoDirection)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkPackets) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NetworkPacketsObservable is an instrument used to record metric values
// conforming to the "hw.network.packets" semantic conventions. It represents the
// received and transmitted network traffic in packets (or frames).
type NetworkPacketsObservable struct {
	metric.Int64ObservableCounter
}

var newNetworkPacketsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Received and transmitted network traffic in packets (or frames)."),
	metric.WithUnit("{packet}"),
}

// NewNetworkPacketsObservable returns a new NetworkPacketsObservable instrument.
func NewNetworkPacketsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NetworkPacketsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkPacketsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkPacketsObservableOpts
	} else {
		opt = append(opt, newNetworkPacketsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.network.packets",
		opt...,
	)
	if err != nil {
		return NetworkPacketsObservable{noop.Int64ObservableCounter{}}, err
	}
	return NetworkPacketsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkPacketsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkPacketsObservable) Name() string {
	return "hw.network.packets"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkPacketsObservable) Unit() string {
	return "{packet}"
}

// Description returns the semantic convention description of the instrument
func (NetworkPacketsObservable) Description() string {
	return "Received and transmitted network traffic in packets (or frames)."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (NetworkPacketsObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrNetworkIODirection returns a required attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NetworkPacketsObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkPacketsObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkPacketsObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkPacketsObservable) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkPacketsObservable) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkPacketsObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkPacketsObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkPacketsObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// NetworkUp is an instrument used to record metric values conforming to the
// "hw.network.up" semantic conventions. It represents the link status: `1` (up)
// or `0` (down).
type NetworkUp struct {
	metric.Int64UpDownCounter
}

var newNetworkUpOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Link status: `1` (up) or `0` (down)."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newNetworkUpOpts
	} else {
		opt = append(opt, newNetworkUpOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.network.up",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NetworkUp) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NetworkUpObservable is an instrument used to record metric values conforming
// to the "hw.network.up" semantic conventions. It represents the link status:
// `1` (up) or `0` (down).
type NetworkUpObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNetworkUpObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Link status: `1` (up) or `0` (down)."),
	metric.WithUnit("1"),
}

// NewNetworkUpObservable returns a new NetworkUpObservable instrument.
func NewNetworkUpObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NetworkUpObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetworkUpObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNetworkUpObservableOpts
	} else {
		opt = append(opt, newNetworkUpObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.network.up",
		opt...,
	)
	if err != nil {
		return NetworkUpObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NetworkUpObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetworkUpObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NetworkUpObservable) Name() string {
	return "hw.network.up"
}

// Unit returns the semantic convention unit of the instrument
func (NetworkUpObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (NetworkUpObservable) Description() string {
	return "Link status: `1` (up) or `0` (down)."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (NetworkUpObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (NetworkUpObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (NetworkUpObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrNetworkLogicalAddresses returns an optional attribute for the
// "hw.network.logical_addresses" semantic convention. It represents the logical
// addresses of the adapter (e.g. IP address, or WWPN).
func (NetworkUpObservable) AttrNetworkLogicalAddresses(val ...string) attribute.KeyValue {
	return attribute.StringSlice("hw.network.logical_addresses", val)
}

// AttrNetworkPhysicalAddress returns an optional attribute for the
// "hw.network.physical_address" semantic convention. It represents the physical
// address of the adapter (e.g. MAC address, or WWNN).
func (NetworkUpObservable) AttrNetworkPhysicalAddress(val string) attribute.KeyValue {
	return attribute.String("hw.network.physical_address", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (NetworkUpObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (NetworkUpObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (NetworkUpObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PhysicalDiskEnduranceUtilization is an instrument used to record metric values
// conforming to the "hw.physical_disk.endurance_utilization" semantic
// conventions. It represents the endurance remaining for this SSD disk.
type PhysicalDiskEnduranceUtilization struct {
	metric.Int64Gauge
}

var newPhysicalDiskEnduranceUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Endurance remaining for this SSD disk."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newPhysicalDiskEnduranceUtilizationOpts
	} else {
		opt = append(opt, newPhysicalDiskEnduranceUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.physical_disk.endurance_utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.physical_disk.state", string(physicalDiskState)),
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
				attribute.String("hw.id", id),
				attribute.String("hw.physical_disk.state", string(physicalDiskState)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PhysicalDiskEnduranceUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// PhysicalDiskEnduranceUtilizationObservable is an instrument used to record
// metric values conforming to the "hw.physical_disk.endurance_utilization"
// semantic conventions. It represents the endurance remaining for this SSD disk.
type PhysicalDiskEnduranceUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newPhysicalDiskEnduranceUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Endurance remaining for this SSD disk."),
	metric.WithUnit("1"),
}

// NewPhysicalDiskEnduranceUtilizationObservable returns a new
// PhysicalDiskEnduranceUtilizationObservable instrument.
func NewPhysicalDiskEnduranceUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (PhysicalDiskEnduranceUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskEnduranceUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPhysicalDiskEnduranceUtilizationObservableOpts
	} else {
		opt = append(opt, newPhysicalDiskEnduranceUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.physical_disk.endurance_utilization",
		opt...,
	)
	if err != nil {
		return PhysicalDiskEnduranceUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return PhysicalDiskEnduranceUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskEnduranceUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskEnduranceUtilizationObservable) Name() string {
	return "hw.physical_disk.endurance_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskEnduranceUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskEnduranceUtilizationObservable) Description() string {
	return "Endurance remaining for this SSD disk."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PhysicalDiskEnduranceUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrPhysicalDiskState returns a required attribute for the
// "hw.physical_disk.state" semantic convention. It represents the state of the
// physical disk endurance utilization.
func (PhysicalDiskEnduranceUtilizationObservable) AttrPhysicalDiskState(val PhysicalDiskStateAttr) attribute.KeyValue {
	return attribute.String("hw.physical_disk.state", string(val))
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskEnduranceUtilizationObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskEnduranceUtilizationObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskEnduranceUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskEnduranceUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskEnduranceUtilizationObservable) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskEnduranceUtilizationObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskEnduranceUtilizationObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PhysicalDiskSize is an instrument used to record metric values conforming to
// the "hw.physical_disk.size" semantic conventions. It represents the size of
// the disk.
type PhysicalDiskSize struct {
	metric.Int64UpDownCounter
}

var newPhysicalDiskSizeOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Size of the disk."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newPhysicalDiskSizeOpts
	} else {
		opt = append(opt, newPhysicalDiskSizeOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.physical_disk.size",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PhysicalDiskSize) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// PhysicalDiskSizeObservable is an instrument used to record metric values
// conforming to the "hw.physical_disk.size" semantic conventions. It represents
// the size of the disk.
type PhysicalDiskSizeObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPhysicalDiskSizeObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Size of the disk."),
	metric.WithUnit("By"),
}

// NewPhysicalDiskSizeObservable returns a new PhysicalDiskSizeObservable
// instrument.
func NewPhysicalDiskSizeObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PhysicalDiskSizeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskSizeObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPhysicalDiskSizeObservableOpts
	} else {
		opt = append(opt, newPhysicalDiskSizeObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.physical_disk.size",
		opt...,
	)
	if err != nil {
		return PhysicalDiskSizeObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PhysicalDiskSizeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskSizeObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskSizeObservable) Name() string {
	return "hw.physical_disk.size"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskSizeObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskSizeObservable) Description() string {
	return "Size of the disk."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PhysicalDiskSizeObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskSizeObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskSizeObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskSizeObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskSizeObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskSizeObservable) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskSizeObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskSizeObservable) AttrVendor(val string) attribute.KeyValue {
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

var newPhysicalDiskSmartOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newPhysicalDiskSmartOpts
	} else {
		opt = append(opt, newPhysicalDiskSmartOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.physical_disk.smart",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PhysicalDiskSmart) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// PhysicalDiskSmartObservable is an instrument used to record metric values
// conforming to the "hw.physical_disk.smart" semantic conventions. It represents
// the value of the corresponding [S.M.A.R.T.] (Self-Monitoring, Analysis, and
// Reporting Technology) attribute.
//
// [S.M.A.R.T.]: https://wikipedia.org/wiki/S.M.A.R.T.
type PhysicalDiskSmartObservable struct {
	metric.Int64ObservableGauge
}

var newPhysicalDiskSmartObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute."),
	metric.WithUnit("1"),
}

// NewPhysicalDiskSmartObservable returns a new PhysicalDiskSmartObservable
// instrument.
func NewPhysicalDiskSmartObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (PhysicalDiskSmartObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PhysicalDiskSmartObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPhysicalDiskSmartObservableOpts
	} else {
		opt = append(opt, newPhysicalDiskSmartObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.physical_disk.smart",
		opt...,
	)
	if err != nil {
		return PhysicalDiskSmartObservable{noop.Int64ObservableGauge{}}, err
	}
	return PhysicalDiskSmartObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PhysicalDiskSmartObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PhysicalDiskSmartObservable) Name() string {
	return "hw.physical_disk.smart"
}

// Unit returns the semantic convention unit of the instrument
func (PhysicalDiskSmartObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PhysicalDiskSmartObservable) Description() string {
	return "Value of the corresponding [S.M.A.R.T.](https://wikipedia.org/wiki/S.M.A.R.T.) (Self-Monitoring, Analysis, and Reporting Technology) attribute."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PhysicalDiskSmartObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrFirmwareVersion returns an optional attribute for the
// "hw.firmware_version" semantic convention. It represents the firmware version
// of the hardware component.
func (PhysicalDiskSmartObservable) AttrFirmwareVersion(val string) attribute.KeyValue {
	return attribute.String("hw.firmware_version", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PhysicalDiskSmartObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PhysicalDiskSmartObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PhysicalDiskSmartObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrPhysicalDiskSmartAttribute returns an optional attribute for the
// "hw.physical_disk.smart_attribute" semantic convention. It represents the
// [S.M.A.R.T.] (Self-Monitoring, Analysis, and Reporting Technology) attribute
// of the physical disk.
//
// [S.M.A.R.T.]: https://wikipedia.org/wiki/S.M.A.R.T.
func (PhysicalDiskSmartObservable) AttrPhysicalDiskSmartAttribute(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.smart_attribute", val)
}

// AttrPhysicalDiskType returns an optional attribute for the
// "hw.physical_disk.type" semantic convention. It represents the type of the
// physical disk.
func (PhysicalDiskSmartObservable) AttrPhysicalDiskType(val string) attribute.KeyValue {
	return attribute.String("hw.physical_disk.type", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PhysicalDiskSmartObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PhysicalDiskSmartObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Power is an instrument used to record metric values conforming to the
// "hw.power" semantic conventions. It represents the instantaneous power
// consumed by the component.
type Power struct {
	metric.Int64Gauge
}

var newPowerOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Instantaneous power consumed by the component."),
	metric.WithUnit("W"),
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

	if len(opt) == 0 {
		opt = newPowerOpts
	} else {
		opt = append(opt, newPowerOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.power",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.type", string(hwType)),
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

// PowerObservable is an instrument used to record metric values conforming to
// the "hw.power" semantic conventions. It represents the instantaneous power
// consumed by the component.
type PowerObservable struct {
	metric.Int64ObservableGauge
}

var newPowerObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Instantaneous power consumed by the component."),
	metric.WithUnit("W"),
}

// NewPowerObservable returns a new PowerObservable instrument.
func NewPowerObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (PowerObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPowerObservableOpts
	} else {
		opt = append(opt, newPowerObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.power",
		opt...,
	)
	if err != nil {
		return PowerObservable{noop.Int64ObservableGauge{}}, err
	}
	return PowerObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PowerObservable) Name() string {
	return "hw.power"
}

// Unit returns the semantic convention unit of the instrument
func (PowerObservable) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (PowerObservable) Description() string {
	return "Instantaneous power consumed by the component."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PowerObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrType returns a required attribute for the "hw.type" semantic convention.
// It represents the type of the component.
func (PowerObservable) AttrType(val TypeAttr) attribute.KeyValue {
	return attribute.String("hw.type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// PowerSupplyLimit is an instrument used to record metric values conforming to
// the "hw.power_supply.limit" semantic conventions. It represents the maximum
// power output of the power supply.
type PowerSupplyLimit struct {
	metric.Int64UpDownCounter
}

var newPowerSupplyLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum power output of the power supply."),
	metric.WithUnit("W"),
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

	if len(opt) == 0 {
		opt = newPowerSupplyLimitOpts
	} else {
		opt = append(opt, newPowerSupplyLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.power_supply.limit",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PowerSupplyLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// PowerSupplyLimitObservable is an instrument used to record metric values
// conforming to the "hw.power_supply.limit" semantic conventions. It represents
// the maximum power output of the power supply.
type PowerSupplyLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPowerSupplyLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum power output of the power supply."),
	metric.WithUnit("W"),
}

// NewPowerSupplyLimitObservable returns a new PowerSupplyLimitObservable
// instrument.
func NewPowerSupplyLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PowerSupplyLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPowerSupplyLimitObservableOpts
	} else {
		opt = append(opt, newPowerSupplyLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.power_supply.limit",
		opt...,
	)
	if err != nil {
		return PowerSupplyLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PowerSupplyLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyLimitObservable) Name() string {
	return "hw.power_supply.limit"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyLimitObservable) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyLimitObservable) Description() string {
	return "Maximum power output of the power supply."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PowerSupplyLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (PowerSupplyLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyLimitObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyLimitObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyLimitObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PowerSupplyUsage is an instrument used to record metric values conforming to
// the "hw.power_supply.usage" semantic conventions. It represents the current
// power output of the power supply.
type PowerSupplyUsage struct {
	metric.Int64UpDownCounter
}

var newPowerSupplyUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Current power output of the power supply."),
	metric.WithUnit("W"),
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

	if len(opt) == 0 {
		opt = newPowerSupplyUsageOpts
	} else {
		opt = append(opt, newPowerSupplyUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.power_supply.usage",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m PowerSupplyUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// PowerSupplyUsageObservable is an instrument used to record metric values
// conforming to the "hw.power_supply.usage" semantic conventions. It represents
// the current power output of the power supply.
type PowerSupplyUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPowerSupplyUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Current power output of the power supply."),
	metric.WithUnit("W"),
}

// NewPowerSupplyUsageObservable returns a new PowerSupplyUsageObservable
// instrument.
func NewPowerSupplyUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PowerSupplyUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPowerSupplyUsageObservableOpts
	} else {
		opt = append(opt, newPowerSupplyUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.power_supply.usage",
		opt...,
	)
	if err != nil {
		return PowerSupplyUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PowerSupplyUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyUsageObservable) Name() string {
	return "hw.power_supply.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyUsageObservable) Unit() string {
	return "W"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyUsageObservable) Description() string {
	return "Current power output of the power supply."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PowerSupplyUsageObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyUsageObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyUsageObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyUsageObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyUsageObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyUsageObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// PowerSupplyUtilization is an instrument used to record metric values
// conforming to the "hw.power_supply.utilization" semantic conventions. It
// represents the utilization of the power supply as a fraction of its maximum
// output.
type PowerSupplyUtilization struct {
	metric.Int64Gauge
}

var newPowerSupplyUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Utilization of the power supply as a fraction of its maximum output."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newPowerSupplyUtilizationOpts
	} else {
		opt = append(opt, newPowerSupplyUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.power_supply.utilization",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m PowerSupplyUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// PowerSupplyUtilizationObservable is an instrument used to record metric values
// conforming to the "hw.power_supply.utilization" semantic conventions. It
// represents the utilization of the power supply as a fraction of its maximum
// output.
type PowerSupplyUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newPowerSupplyUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Utilization of the power supply as a fraction of its maximum output."),
	metric.WithUnit("1"),
}

// NewPowerSupplyUtilizationObservable returns a new
// PowerSupplyUtilizationObservable instrument.
func NewPowerSupplyUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (PowerSupplyUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PowerSupplyUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPowerSupplyUtilizationObservableOpts
	} else {
		opt = append(opt, newPowerSupplyUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.power_supply.utilization",
		opt...,
	)
	if err != nil {
		return PowerSupplyUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return PowerSupplyUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PowerSupplyUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PowerSupplyUtilizationObservable) Name() string {
	return "hw.power_supply.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PowerSupplyUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PowerSupplyUtilizationObservable) Description() string {
	return "Utilization of the power supply as a fraction of its maximum output."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (PowerSupplyUtilizationObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (PowerSupplyUtilizationObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (PowerSupplyUtilizationObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (PowerSupplyUtilizationObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (PowerSupplyUtilizationObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (PowerSupplyUtilizationObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Status is an instrument used to record metric values conforming to the
// "hw.status" semantic conventions. It represents the operational status: `1`
// (true) or `0` (false) for each of the possible states.
type Status struct {
	metric.Int64UpDownCounter
}

var newStatusOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Operational status: `1` (true) or `0` (false) for each of the possible states."),
	metric.WithUnit("1"),
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

	if len(opt) == 0 {
		opt = newStatusOpts
	} else {
		opt = append(opt, newStatusOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"hw.status",
		opt...,
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
			attribute.String("hw.state", string(state)),
			attribute.String("hw.type", string(hwType)),
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

// StatusObservable is an instrument used to record metric values conforming to
// the "hw.status" semantic conventions. It represents the operational status:
// `1` (true) or `0` (false) for each of the possible states.
type StatusObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatusObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Operational status: `1` (true) or `0` (false) for each of the possible states."),
	metric.WithUnit("1"),
}

// NewStatusObservable returns a new StatusObservable instrument.
func NewStatusObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatusObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatusObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatusObservableOpts
	} else {
		opt = append(opt, newStatusObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"hw.status",
		opt...,
	)
	if err != nil {
		return StatusObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatusObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatusObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatusObservable) Name() string {
	return "hw.status"
}

// Unit returns the semantic convention unit of the instrument
func (StatusObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (StatusObservable) Description() string {
	return "Operational status: `1` (true) or `0` (false) for each of the possible states."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (StatusObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrState returns a required attribute for the "hw.state" semantic convention.
// It represents the current state of the component.
func (StatusObservable) AttrState(val StateAttr) attribute.KeyValue {
	return attribute.String("hw.state", string(val))
}

// AttrType returns a required attribute for the "hw.type" semantic convention.
// It represents the type of the component.
func (StatusObservable) AttrType(val TypeAttr) attribute.KeyValue {
	return attribute.String("hw.type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (StatusObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (StatusObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// TapeDriveOperations is an instrument used to record metric values conforming
// to the "hw.tape_drive.operations" semantic conventions. It represents the
// operations performed by the tape drive.
type TapeDriveOperations struct {
	metric.Int64Counter
}

var newTapeDriveOperationsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Operations performed by the tape drive."),
	metric.WithUnit("{operation}"),
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

	if len(opt) == 0 {
		opt = newTapeDriveOperationsOpts
	} else {
		opt = append(opt, newTapeDriveOperationsOpts...)
	}

	i, err := m.Int64Counter(
		"hw.tape_drive.operations",
		opt...,
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m TapeDriveOperations) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// TapeDriveOperationsObservable is an instrument used to record metric values
// conforming to the "hw.tape_drive.operations" semantic conventions. It
// represents the operations performed by the tape drive.
type TapeDriveOperationsObservable struct {
	metric.Int64ObservableCounter
}

var newTapeDriveOperationsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Operations performed by the tape drive."),
	metric.WithUnit("{operation}"),
}

// NewTapeDriveOperationsObservable returns a new TapeDriveOperationsObservable
// instrument.
func NewTapeDriveOperationsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (TapeDriveOperationsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return TapeDriveOperationsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newTapeDriveOperationsObservableOpts
	} else {
		opt = append(opt, newTapeDriveOperationsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"hw.tape_drive.operations",
		opt...,
	)
	if err != nil {
		return TapeDriveOperationsObservable{noop.Int64ObservableCounter{}}, err
	}
	return TapeDriveOperationsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m TapeDriveOperationsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (TapeDriveOperationsObservable) Name() string {
	return "hw.tape_drive.operations"
}

// Unit returns the semantic convention unit of the instrument
func (TapeDriveOperationsObservable) Unit() string {
	return "{operation}"
}

// Description returns the semantic convention description of the instrument
func (TapeDriveOperationsObservable) Description() string {
	return "Operations performed by the tape drive."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (TapeDriveOperationsObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrModel returns an optional attribute for the "hw.model" semantic
// convention. It represents the descriptive model name of the hardware
// component.
func (TapeDriveOperationsObservable) AttrModel(val string) attribute.KeyValue {
	return attribute.String("hw.model", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (TapeDriveOperationsObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (TapeDriveOperationsObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSerialNumber returns an optional attribute for the "hw.serial_number"
// semantic convention. It represents the serial number of the hardware
// component.
func (TapeDriveOperationsObservable) AttrSerialNumber(val string) attribute.KeyValue {
	return attribute.String("hw.serial_number", val)
}

// AttrTapeDriveOperationType returns an optional attribute for the
// "hw.tape_drive.operation_type" semantic convention. It represents the type of
// tape drive operation.
func (TapeDriveOperationsObservable) AttrTapeDriveOperationType(val TapeDriveOperationTypeAttr) attribute.KeyValue {
	return attribute.String("hw.tape_drive.operation_type", string(val))
}

// AttrVendor returns an optional attribute for the "hw.vendor" semantic
// convention. It represents the vendor name of the hardware component.
func (TapeDriveOperationsObservable) AttrVendor(val string) attribute.KeyValue {
	return attribute.String("hw.vendor", val)
}

// Temperature is an instrument used to record metric values conforming to the
// "hw.temperature" semantic conventions. It represents the temperature in
// degrees Celsius.
type Temperature struct {
	metric.Int64Gauge
}

var newTemperatureOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Temperature in degrees Celsius."),
	metric.WithUnit("Cel"),
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

	if len(opt) == 0 {
		opt = newTemperatureOpts
	} else {
		opt = append(opt, newTemperatureOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.temperature",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m Temperature) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// TemperatureObservable is an instrument used to record metric values conforming
// to the "hw.temperature" semantic conventions. It represents the temperature in
// degrees Celsius.
type TemperatureObservable struct {
	metric.Int64ObservableGauge
}

var newTemperatureObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Temperature in degrees Celsius."),
	metric.WithUnit("Cel"),
}

// NewTemperatureObservable returns a new TemperatureObservable instrument.
func NewTemperatureObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (TemperatureObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return TemperatureObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newTemperatureObservableOpts
	} else {
		opt = append(opt, newTemperatureObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.temperature",
		opt...,
	)
	if err != nil {
		return TemperatureObservable{noop.Int64ObservableGauge{}}, err
	}
	return TemperatureObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m TemperatureObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (TemperatureObservable) Name() string {
	return "hw.temperature"
}

// Unit returns the semantic convention unit of the instrument
func (TemperatureObservable) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (TemperatureObservable) Description() string {
	return "Temperature in degrees Celsius."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (TemperatureObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (TemperatureObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (TemperatureObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (TemperatureObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// TemperatureLimit is an instrument used to record metric values conforming to
// the "hw.temperature.limit" semantic conventions. It represents the temperature
// limit in degrees Celsius.
type TemperatureLimit struct {
	metric.Int64Gauge
}

var newTemperatureLimitOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Temperature limit in degrees Celsius."),
	metric.WithUnit("Cel"),
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

	if len(opt) == 0 {
		opt = newTemperatureLimitOpts
	} else {
		opt = append(opt, newTemperatureLimitOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.temperature.limit",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m TemperatureLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// TemperatureLimitObservable is an instrument used to record metric values
// conforming to the "hw.temperature.limit" semantic conventions. It represents
// the temperature limit in degrees Celsius.
type TemperatureLimitObservable struct {
	metric.Int64ObservableGauge
}

var newTemperatureLimitObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Temperature limit in degrees Celsius."),
	metric.WithUnit("Cel"),
}

// NewTemperatureLimitObservable returns a new TemperatureLimitObservable
// instrument.
func NewTemperatureLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (TemperatureLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return TemperatureLimitObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newTemperatureLimitObservableOpts
	} else {
		opt = append(opt, newTemperatureLimitObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.temperature.limit",
		opt...,
	)
	if err != nil {
		return TemperatureLimitObservable{noop.Int64ObservableGauge{}}, err
	}
	return TemperatureLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m TemperatureLimitObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (TemperatureLimitObservable) Name() string {
	return "hw.temperature.limit"
}

// Unit returns the semantic convention unit of the instrument
func (TemperatureLimitObservable) Unit() string {
	return "Cel"
}

// Description returns the semantic convention description of the instrument
func (TemperatureLimitObservable) Description() string {
	return "Temperature limit in degrees Celsius."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (TemperatureLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (TemperatureLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (TemperatureLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (TemperatureLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (TemperatureLimitObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// Voltage is an instrument used to record metric values conforming to the
// "hw.voltage" semantic conventions. It represents the voltage measured by the
// sensor.
type Voltage struct {
	metric.Int64Gauge
}

var newVoltageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Voltage measured by the sensor."),
	metric.WithUnit("V"),
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

	if len(opt) == 0 {
		opt = newVoltageOpts
	} else {
		opt = append(opt, newVoltageOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.voltage",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m Voltage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// VoltageObservable is an instrument used to record metric values conforming to
// the "hw.voltage" semantic conventions. It represents the voltage measured by
// the sensor.
type VoltageObservable struct {
	metric.Int64ObservableGauge
}

var newVoltageObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Voltage measured by the sensor."),
	metric.WithUnit("V"),
}

// NewVoltageObservable returns a new VoltageObservable instrument.
func NewVoltageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (VoltageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return VoltageObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newVoltageObservableOpts
	} else {
		opt = append(opt, newVoltageObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.voltage",
		opt...,
	)
	if err != nil {
		return VoltageObservable{noop.Int64ObservableGauge{}}, err
	}
	return VoltageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m VoltageObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (VoltageObservable) Name() string {
	return "hw.voltage"
}

// Unit returns the semantic convention unit of the instrument
func (VoltageObservable) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (VoltageObservable) Description() string {
	return "Voltage measured by the sensor."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (VoltageObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (VoltageObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (VoltageObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (VoltageObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// VoltageLimit is an instrument used to record metric values conforming to the
// "hw.voltage.limit" semantic conventions. It represents the voltage limit in
// Volts.
type VoltageLimit struct {
	metric.Int64Gauge
}

var newVoltageLimitOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Voltage limit in Volts."),
	metric.WithUnit("V"),
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

	if len(opt) == 0 {
		opt = newVoltageLimitOpts
	} else {
		opt = append(opt, newVoltageLimitOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.voltage.limit",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m VoltageLimit) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// VoltageLimitObservable is an instrument used to record metric values
// conforming to the "hw.voltage.limit" semantic conventions. It represents the
// voltage limit in Volts.
type VoltageLimitObservable struct {
	metric.Int64ObservableGauge
}

var newVoltageLimitObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Voltage limit in Volts."),
	metric.WithUnit("V"),
}

// NewVoltageLimitObservable returns a new VoltageLimitObservable instrument.
func NewVoltageLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (VoltageLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return VoltageLimitObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newVoltageLimitObservableOpts
	} else {
		opt = append(opt, newVoltageLimitObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.voltage.limit",
		opt...,
	)
	if err != nil {
		return VoltageLimitObservable{noop.Int64ObservableGauge{}}, err
	}
	return VoltageLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m VoltageLimitObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (VoltageLimitObservable) Name() string {
	return "hw.voltage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (VoltageLimitObservable) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (VoltageLimitObservable) Description() string {
	return "Voltage limit in Volts."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (VoltageLimitObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrLimitType returns an optional attribute for the "hw.limit_type" semantic
// convention. It represents the type of limit for hardware components.
func (VoltageLimitObservable) AttrLimitType(val LimitTypeAttr) attribute.KeyValue {
	return attribute.String("hw.limit_type", string(val))
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (VoltageLimitObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (VoltageLimitObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (VoltageLimitObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}

// VoltageNominal is an instrument used to record metric values conforming to the
// "hw.voltage.nominal" semantic conventions. It represents the nominal
// (expected) voltage.
type VoltageNominal struct {
	metric.Int64Gauge
}

var newVoltageNominalOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Nominal (expected) voltage."),
	metric.WithUnit("V"),
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

	if len(opt) == 0 {
		opt = newVoltageNominalOpts
	} else {
		opt = append(opt, newVoltageNominalOpts...)
	}

	i, err := m.Int64Gauge(
		"hw.voltage.nominal",
		opt...,
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
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("hw.id", id),
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
				attribute.String("hw.id", id),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m VoltageNominal) RecordSet(ctx context.Context, val int64, set attribute.Set) {
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

// VoltageNominalObservable is an instrument used to record metric values
// conforming to the "hw.voltage.nominal" semantic conventions. It represents the
// nominal (expected) voltage.
type VoltageNominalObservable struct {
	metric.Int64ObservableGauge
}

var newVoltageNominalObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Nominal (expected) voltage."),
	metric.WithUnit("V"),
}

// NewVoltageNominalObservable returns a new VoltageNominalObservable instrument.
func NewVoltageNominalObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (VoltageNominalObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return VoltageNominalObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newVoltageNominalObservableOpts
	} else {
		opt = append(opt, newVoltageNominalObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"hw.voltage.nominal",
		opt...,
	)
	if err != nil {
		return VoltageNominalObservable{noop.Int64ObservableGauge{}}, err
	}
	return VoltageNominalObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m VoltageNominalObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (VoltageNominalObservable) Name() string {
	return "hw.voltage.nominal"
}

// Unit returns the semantic convention unit of the instrument
func (VoltageNominalObservable) Unit() string {
	return "V"
}

// Description returns the semantic convention description of the instrument
func (VoltageNominalObservable) Description() string {
	return "Nominal (expected) voltage."
}

// AttrID returns a required attribute for the "hw.id" semantic convention. It
// represents an identifier for the hardware component, unique within the
// monitored host.
func (VoltageNominalObservable) AttrID(val string) attribute.KeyValue {
	return attribute.String("hw.id", val)
}

// AttrName returns an optional attribute for the "hw.name" semantic convention.
// It represents an easily-recognizable name for the hardware component.
func (VoltageNominalObservable) AttrName(val string) attribute.KeyValue {
	return attribute.String("hw.name", val)
}

// AttrParent returns an optional attribute for the "hw.parent" semantic
// convention. It represents the unique identifier of the parent component
// (typically the `hw.id` attribute of the enclosure, or disk controller).
func (VoltageNominalObservable) AttrParent(val string) attribute.KeyValue {
	return attribute.String("hw.parent", val)
}

// AttrSensorLocation returns an optional attribute for the "hw.sensor_location"
// semantic convention. It represents the location of the sensor.
func (VoltageNominalObservable) AttrSensorLocation(val string) attribute.KeyValue {
	return attribute.String("hw.sensor_location", val)
}
