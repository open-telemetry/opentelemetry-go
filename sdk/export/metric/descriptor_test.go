// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
