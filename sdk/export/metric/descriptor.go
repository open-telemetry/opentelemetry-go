package metric

import (
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/sdk/resource"
)

// DescriptorConfig contains configuration for a metric Descriptor
type DescriptorConfig struct {
	// keys are the common keys related to all measurements.
	Keys []core.Key

	// description describes the metric in human readable terms.
	Description string

	// unit is the determinate quantity used as a standard of all
	// measurements for the metric.
	Unit unit.Unit

	// resource descrbes the entity for which this metric makes measurements.
	Resource resource.Resource
}

// DescriptorOption is the interface that applies the value to a DescriptorConfig option.
type DescriptorOption interface {
	// Apply sets the option value of a DescriptorConfig.
	Apply(*DescriptorConfig)
}

// WithKeys applies common label keys.
// Multiple `WithKeys` options accumulate.
func WithKeys(keys ...core.Key) DescriptorOption {
	return keysOption(keys)
}

type keysOption []core.Key

func (k keysOption) Apply(config *DescriptorConfig) {
	if config == nil {
		return
	}
	config.Keys = append(config.Keys, k...)
}

// WithDescription applies provided description.
func WithDescription(d string) DescriptorOption {
	return descriptionOption(d)
}

type descriptionOption string

func (d descriptionOption) Apply(config *DescriptorConfig) {
	if config == nil {
		return
	}
	config.Description = string(d)
}

// WithUnit applies provided unit.
func WithUnit(u unit.Unit) DescriptorOption {
	return unitOption(u)
}

type unitOption unit.Unit

func (u unitOption) Apply(config *DescriptorConfig) {
	if config == nil {
		return
	}
	config.Unit = unit.Unit(u)
}

// WithResource applies provided Resource.
func WithResource(r resource.Resource) DescriptorOption {
	return resourceOption(r)
}

type resourceOption resource.Resource

func (r resourceOption) Apply(config *DescriptorConfig) {
	if config == nil {
		return
	}
	config.Resource = resource.Resource(r)
}

// Descriptor describes a metric instrument to the exporter.
//
// Descriptors are created once per instrument and a pointer to the
// descriptor may be used to uniquely identify the instrument in an
// exporter.
type Descriptor struct {
	name        string
	metricKind  Kind
	keys        []core.Key
	description string
	unit        unit.Unit
	numberKind  core.NumberKind
	resource    resource.Resource
}

// NewDescriptor builds a new descriptor, for use by `Meter`
// implementations in constructing new metric instruments.
//
// Descriptors are created once per instrument and a pointer to the
// descriptor may be used to uniquely identify the instrument in an
// exporter.
func NewDescriptor(name string, metricKind Kind, numberKind core.NumberKind, opts ...DescriptorOption) *Descriptor {
	c := &DescriptorConfig{}
	for _, opt := range opts {
		opt.Apply(c)
	}

	return &Descriptor{
		name:        name,
		metricKind:  metricKind,
		keys:        c.Keys,
		description: c.Description,
		unit:        c.Unit,
		numberKind:  numberKind,
		resource:    c.Resource,
	}
}

// Name returns the metric instrument's name.
func (d *Descriptor) Name() string {
	return d.name
}

// MetricKind returns the kind of instrument: counter, measure, or
// observer.
func (d *Descriptor) MetricKind() Kind {
	return d.metricKind
}

// Keys returns the recommended keys included in the metric
// definition.  These keys may be used by a Batcher as a default set
// of grouping keys for the metric instrument.
func (d *Descriptor) Keys() []core.Key {
	return d.keys
}

// Description provides a human-readable description of the metric
// instrument.
func (d *Descriptor) Description() string {
	return d.description
}

// Unit describes the units of the metric instrument.  Unitless
// metrics return the empty string.
func (d *Descriptor) Unit() unit.Unit {
	return d.unit
}

// NumberKind returns whether this instrument is declared over int64
// or a float64 values.
func (d *Descriptor) NumberKind() core.NumberKind {
	return d.numberKind
}

// Resource returns the Resource describing the entity for whom the metric
// instrument measures.
func (d *Descriptor) Resource() resource.Resource {
	return d.resource
}
