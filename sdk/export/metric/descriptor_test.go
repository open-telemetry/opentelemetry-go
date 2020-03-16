package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestDescriptorNoOpts(t *testing.T) {
	tests := []Descriptor{
		{}, // Empty values for everything
		{
			name:       "empty value for options",
			metricKind: CounterKind,
			numberKind: core.Float64NumberKind,
		},
	}

	for _, test := range tests {
		got := NewDescriptor(test.name, test.metricKind, test.numberKind)
		assert.Equal(t, test, *got)
	}
}

func TestDescriptorConfiguration(t *testing.T) {
	tests := []Descriptor{
		{}, // Empty values for everything
		{
			name:       "empty value for options",
			metricKind: CounterKind,
			numberKind: core.Float64NumberKind,
		},
		{
			name:        "with options",
			metricKind:  CounterKind,
			numberKind:  core.Float64NumberKind,
			description: "test description",
			unit:        unit.Bytes,
			keys:        []core.Key{"keys key"},
			resource:    *(resource.New(core.Key("resource key").Bool(true))),
		},
	}

	for _, test := range tests {
		opts := []DescriptorOption{
			WithKeys(test.keys...),
			WithDescription(test.description),
			WithUnit(test.unit),
			WithResource(test.resource),
		}

		got := NewDescriptor(test.name, test.metricKind, test.numberKind, opts...)
		assert.Equal(t, test, *got)
	}
}
