// Copyright The OpenTelemetry Authors
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

package reducer // import "go.opentelemetry.io/otel/sdk/metric/processor/reducer"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

type (
	// Processor implements "dimensionality reduction" by
	// filtering keys from export attribute sets.
	Processor struct {
		export.Checkpointer
		filterSelector AttributeFilterSelector
	}

	// AttributeFilterSelector selects an attribute filter based on the
	// instrument described by the descriptor.
	AttributeFilterSelector interface {
		AttributeFilterFor(descriptor *sdkapi.Descriptor) attribute.Filter
	}
)

var _ export.Processor = &Processor{}
var _ export.Checkpointer = &Processor{}

// New returns a dimensionality-reducing Processor that passes data to the
// next stage in an export pipeline.
func New(filterSelector AttributeFilterSelector, ckpter export.Checkpointer) *Processor {
	return &Processor{
		Checkpointer:   ckpter,
		filterSelector: filterSelector,
	}
}

// Process implements export.Processor.
func (p *Processor) Process(accum export.Accumulation) error {
	// Note: the removed attributes are returned and ignored here.
	// Conceivably these inputs could be useful to a sampler.
	reduced, _ := accum.Attributes().Filter(
		p.filterSelector.AttributeFilterFor(
			accum.Descriptor(),
		),
	)
	return p.Checkpointer.Process(
		export.NewAccumulation(
			accum.Descriptor(),
			&reduced,
			accum.Aggregator(),
		),
	)
}
