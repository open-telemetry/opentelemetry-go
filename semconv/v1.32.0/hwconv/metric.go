// Code generated from semantic convention specification. DO NOT EDIT.

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
			metric.WithDescription("Energy consumed by the component"),
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
	return "Energy consumed by the component"
}

// Add adds incr to the existing count.
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
			metric.WithDescription("Number of errors encountered by the component"),
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
	return "Number of errors encountered by the component"
}

// Add adds incr to the existing count.
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
			metric.WithDescription("Ambient (external) temperature of the physical host"),
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
	return "Ambient (external) temperature of the physical host"
}

// Record records val to the current distribution.
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
			metric.WithDescription("Total energy consumed by the entire physical host, in joules"),
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
	return "Total energy consumed by the entire physical host, in joules"
}

// Add adds incr to the existing count.
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
			metric.WithDescription("By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors"),
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
	return "By how many degrees Celsius the temperature of the physical host can be increased, before reaching a warning threshold on one of the internal sensors"
}

// Record records val to the current distribution.
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
			metric.WithDescription("Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)"),
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
	return "Instantaneous power consumed by the entire physical host in Watts (`hw.host.energy` is preferred)"
}

// Record records val to the current distribution.
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
			metric.WithDescription("Instantaneous power consumed by the component"),
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
	return "Instantaneous power consumed by the component"
}

// Record records val to the current distribution.
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
			metric.WithDescription("Operational status: `1` (true) or `0` (false) for each of the possible states"),
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
	return "Operational status: `1` (true) or `0` (false) for each of the possible states"
}

// Add adds incr to the existing count.
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