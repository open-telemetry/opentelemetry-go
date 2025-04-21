// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/hw"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the type of error encountered by the component.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// StateAttr is an attribute conforming to the hw.state semantic conventions. It
// represents the current state of the component.
type StateAttr string

var (
	// StateOk is the ok.
	StateOk StateAttr = "ok"
	// StateDegraded is the degraded.
	StateDegraded StateAttr = "degraded"
	// StateFailed is the failed.
	StateFailed StateAttr = "failed"
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

// HwEnergy is an instrument used to record metric values conforming to the
// "hw.energy" semantic conventions. It represents the energy consumed by the
// component.
type Energy struct {
	inst metric.Int64Counter
}

// NewEnergy returns a new Energy instrument.
func NewEnergy(m metric.Meter) (Energy, error) {
	i, err := m.Int64Counter(
	    "hw.energy",
	    metric.WithDescription("Energy consumed by the component"),
	    metric.WithUnit("J"),
	)
	if err != nil {
	    return Energy{}, err
	}
	return Energy{i}, nil
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
	return "Energy consumed by the component"
}

// Add adds incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
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
	attrs ...EnergyAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)
}

func (m Energy) conv(in []EnergyAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.energyAttr()
	}
	return out
}

// EnergyAttr is an optional attribute for the Energy instrument.
type EnergyAttr interface {
    energyAttr() attribute.KeyValue
}

type energyAttr struct {
	kv attribute.KeyValue
}

func (a energyAttr) energyAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (Energy) Name(val string) EnergyAttr {
	return energyAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (Energy) Parent(val string) EnergyAttr {
	return energyAttr{kv: attribute.String("hw.parent", val)}
}

// HwErrors is an instrument used to record metric values conforming to the
// "hw.errors" semantic conventions. It represents the number of errors
// encountered by the component.
type Errors struct {
	inst metric.Int64Counter
}

// NewErrors returns a new Errors instrument.
func NewErrors(m metric.Meter) (Errors, error) {
	i, err := m.Int64Counter(
	    "hw.errors",
	    metric.WithDescription("Number of errors encountered by the component"),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return Errors{}, err
	}
	return Errors{i}, nil
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
	return "Number of errors encountered by the component"
}

// Add adds incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
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
	attrs ...ErrorsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)
}

func (m Errors) conv(in []ErrorsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.errorsAttr()
	}
	return out
}

// ErrorsAttr is an optional attribute for the Errors instrument.
type ErrorsAttr interface {
    errorsAttr() attribute.KeyValue
}

type errorsAttr struct {
	kv attribute.KeyValue
}

func (a errorsAttr) errorsAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the type of error encountered by the component.
func (Errors) ErrorType(val ErrorTypeAttr) ErrorsAttr {
	return errorsAttr{kv: attribute.String("error.type", string(val))}
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (Errors) Name(val string) ErrorsAttr {
	return errorsAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (Errors) Parent(val string) ErrorsAttr {
	return errorsAttr{kv: attribute.String("hw.parent", val)}
}

// HwHostAmbientTemperature is an instrument used to record metric values
// conforming to the "hw.host.ambient_temperature" semantic conventions. It
// represents the ambient (external) temperature of the physical host.
type HostAmbientTemperature struct {
	inst metric.Int64Gauge
}

// NewHostAmbientTemperature returns a new HostAmbientTemperature instrument.
func NewHostAmbientTemperature(m metric.Meter) (HostAmbientTemperature, error) {
	i, err := m.Int64Gauge(
	    "hw.host.ambient_temperature",
	    metric.WithDescription("Ambient (external) temperature of the physical host"),
	    metric.WithUnit("Cel"),
	)
	if err != nil {
	    return HostAmbientTemperature{}, err
	}
	return HostAmbientTemperature{i}, nil
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
	return "Ambient (external) temperature of the physical host"
}

// Record records incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostAmbientTemperature) Record(
    ctx context.Context,
    val int64,
	id string,
	attrs ...HostAmbientTemperatureAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
			)...,
		),
	)
}

func (m HostAmbientTemperature) conv(in []HostAmbientTemperatureAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.hostAmbientTemperatureAttr()
	}
	return out
}

// HostAmbientTemperatureAttr is an optional attribute for the
// HostAmbientTemperature instrument.
type HostAmbientTemperatureAttr interface {
    hostAmbientTemperatureAttr() attribute.KeyValue
}

type hostAmbientTemperatureAttr struct {
	kv attribute.KeyValue
}

func (a hostAmbientTemperatureAttr) hostAmbientTemperatureAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (HostAmbientTemperature) Name(val string) HostAmbientTemperatureAttr {
	return hostAmbientTemperatureAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (HostAmbientTemperature) Parent(val string) HostAmbientTemperatureAttr {
	return hostAmbientTemperatureAttr{kv: attribute.String("hw.parent", val)}
}

// HwHostEnergy is an instrument used to record metric values conforming to the
// "hw.host.energy" semantic conventions. It represents the total energy consumed
// by the entire physical host, in joules.
type HostEnergy struct {
	inst metric.Int64Counter
}

// NewHostEnergy returns a new HostEnergy instrument.
func NewHostEnergy(m metric.Meter) (HostEnergy, error) {
	i, err := m.Int64Counter(
	    "hw.host.energy",
	    metric.WithDescription("Total energy consumed by the entire physical host, in joules"),
	    metric.WithUnit("J"),
	)
	if err != nil {
	    return HostEnergy{}, err
	}
	return HostEnergy{i}, nil
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
	return "Total energy consumed by the entire physical host, in joules"
}

// Add adds incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostEnergy) Add(
    ctx context.Context,
    incr int64,
	id string,
	attrs ...HostEnergyAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
			)...,
		),
	)
}

func (m HostEnergy) conv(in []HostEnergyAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.hostEnergyAttr()
	}
	return out
}

// HostEnergyAttr is an optional attribute for the HostEnergy instrument.
type HostEnergyAttr interface {
    hostEnergyAttr() attribute.KeyValue
}

type hostEnergyAttr struct {
	kv attribute.KeyValue
}

func (a hostEnergyAttr) hostEnergyAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (HostEnergy) Name(val string) HostEnergyAttr {
	return hostEnergyAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (HostEnergy) Parent(val string) HostEnergyAttr {
	return hostEnergyAttr{kv: attribute.String("hw.parent", val)}
}

// HwHostHeatingMargin is an instrument used to record metric values conforming
// to the "hw.host.heating_margin" semantic conventions. It represents the by how
// many degrees Celsius the temperature of the physical host can be increased,
// before reaching a warning threshold on one of the internal sensors.
type HostHeatingMargin struct {
	inst metric.Int64Gauge
}

// NewHostHeatingMargin returns a new HostHeatingMargin instrument.
func NewHostHeatingMargin(m metric.Meter) (HostHeatingMargin, error) {
	i, err := m.Int64Gauge(
	    "hw.host.heating_margin",
	    metric.WithDescription("By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors"),
	    metric.WithUnit("Cel"),
	)
	if err != nil {
	    return HostHeatingMargin{}, err
	}
	return HostHeatingMargin{i}, nil
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
	return "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors"
}

// Record records incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostHeatingMargin) Record(
    ctx context.Context,
    val int64,
	id string,
	attrs ...HostHeatingMarginAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
			)...,
		),
	)
}

func (m HostHeatingMargin) conv(in []HostHeatingMarginAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.hostHeatingMarginAttr()
	}
	return out
}

// HostHeatingMarginAttr is an optional attribute for the HostHeatingMargin
// instrument.
type HostHeatingMarginAttr interface {
    hostHeatingMarginAttr() attribute.KeyValue
}

type hostHeatingMarginAttr struct {
	kv attribute.KeyValue
}

func (a hostHeatingMarginAttr) hostHeatingMarginAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (HostHeatingMargin) Name(val string) HostHeatingMarginAttr {
	return hostHeatingMarginAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (HostHeatingMargin) Parent(val string) HostHeatingMarginAttr {
	return hostHeatingMarginAttr{kv: attribute.String("hw.parent", val)}
}

// HwHostPower is an instrument used to record metric values conforming to the
// "hw.host.power" semantic conventions. It represents the instantaneous power
// consumed by the entire physical host in Watts (`hw.host.energy` is preferred).
type HostPower struct {
	inst metric.Int64Gauge
}

// NewHostPower returns a new HostPower instrument.
func NewHostPower(m metric.Meter) (HostPower, error) {
	i, err := m.Int64Gauge(
	    "hw.host.power",
	    metric.WithDescription("Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)"),
	    metric.WithUnit("W"),
	)
	if err != nil {
	    return HostPower{}, err
	}
	return HostPower{i}, nil
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
	return "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)"
}

// Record records incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// All additional attrs passed are included in the recorded value.
func (m HostPower) Record(
    ctx context.Context,
    val int64,
	id string,
	attrs ...HostPowerAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
			)...,
		),
	)
}

func (m HostPower) conv(in []HostPowerAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.hostPowerAttr()
	}
	return out
}

// HostPowerAttr is an optional attribute for the HostPower instrument.
type HostPowerAttr interface {
    hostPowerAttr() attribute.KeyValue
}

type hostPowerAttr struct {
	kv attribute.KeyValue
}

func (a hostPowerAttr) hostPowerAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (HostPower) Name(val string) HostPowerAttr {
	return hostPowerAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (HostPower) Parent(val string) HostPowerAttr {
	return hostPowerAttr{kv: attribute.String("hw.parent", val)}
}

// HwPower is an instrument used to record metric values conforming to the
// "hw.power" semantic conventions. It represents the instantaneous power
// consumed by the component.
type Power struct {
	inst metric.Int64Gauge
}

// NewPower returns a new Power instrument.
func NewPower(m metric.Meter) (Power, error) {
	i, err := m.Int64Gauge(
	    "hw.power",
	    metric.WithDescription("Instantaneous power consumed by the component"),
	    metric.WithUnit("W"),
	)
	if err != nil {
	    return Power{}, err
	}
	return Power{i}, nil
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
	return "Instantaneous power consumed by the component"
}

// Record records incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
func (m Power) Record(
    ctx context.Context,
    val int64,
	id string,
	hwType TypeAttr,
	attrs ...PowerAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)
}

func (m Power) conv(in []PowerAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.powerAttr()
	}
	return out
}

// PowerAttr is an optional attribute for the Power instrument.
type PowerAttr interface {
    powerAttr() attribute.KeyValue
}

type powerAttr struct {
	kv attribute.KeyValue
}

func (a powerAttr) powerAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (Power) Name(val string) PowerAttr {
	return powerAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (Power) Parent(val string) PowerAttr {
	return powerAttr{kv: attribute.String("hw.parent", val)}
}

// HwStatus is an instrument used to record metric values conforming to the
// "hw.status" semantic conventions. It represents the operational status: `1`
// (true) or `0` (false) for each of the possible states.
type Status struct {
	inst metric.Int64UpDownCounter
}

// NewStatus returns a new Status instrument.
func NewStatus(m metric.Meter) (Status, error) {
	i, err := m.Int64UpDownCounter(
	    "hw.status",
	    metric.WithDescription("Operational status: `1` (true) or `0` (false) for each of the possible states"),
	    metric.WithUnit("1"),
	)
	if err != nil {
	    return Status{}, err
	}
	return Status{i}, nil
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
	return "Operational status: `1` (true) or `0` (false) for each of the possible states"
}

// Add adds incr to the existing count.
//
// The hwId is the an identifier for the hardware component, unique within the
// monitored host
//
// The hwState is the the current state of the component
//
// The hwType is the type of the component
//
// All additional attrs passed are included in the recorded value.
func (m Status) Add(
    ctx context.Context,
    incr int64,
	id string,
	state StateAttr,
	hwType TypeAttr,
	attrs ...StatusAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("hw.id", id),
				attribute.String("hw.state", string(state)),
				attribute.String("hw.type", string(hwType)),
			)...,
		),
	)
}

func (m Status) conv(in []StatusAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.statusAttr()
	}
	return out
}

// StatusAttr is an optional attribute for the Status instrument.
type StatusAttr interface {
    statusAttr() attribute.KeyValue
}

type statusAttr struct {
	kv attribute.KeyValue
}

func (a statusAttr) statusAttr() attribute.KeyValue {
    return a.kv
}

// Name returns an optional attribute for the "hw.name" semantic convention. It
// represents an easily-recognizable name for the hardware component.
func (Status) Name(val string) StatusAttr {
	return statusAttr{kv: attribute.String("hw.name", val)}
}

// Parent returns an optional attribute for the "hw.parent" semantic convention.
// It represents the unique identifier of the parent component (typically the
// `hw.id` attribute of the enclosure, or disk controller).
func (Status) Parent(val string) StatusAttr {
	return statusAttr{kv: attribute.String("hw.parent", val)}
}