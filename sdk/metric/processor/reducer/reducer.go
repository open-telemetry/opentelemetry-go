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
	"go.opentelemetry.io/otel/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type (
	// Processor implements "dimensionality reduction" by
	// filtering keys from export label sets.
	Processor struct {
		export.Checkpointer
		filterSelector LabelFilterSelector
	}

	// LabelFilterSelector is the interface used to configure a
	// specific Filter to an instrument.
	LabelFilterSelector interface {
		LabelFilterFor(descriptor *metric.Descriptor) attribute.Filter
	}
)

var _ export.Processor = &Processor{}
var _ export.Checkpointer = &Processor{}

// New returns a dimensionality-reducing Processor that passes data to
// the next stage in an export pipeline.
func New(filterSelector LabelFilterSelector, ckpter export.Checkpointer) *Processor {
	return &Processor{
		Checkpointer:   ckpter,
		filterSelector: filterSelector,
	}
}

// Process implements export.Processor.
func (p *Processor) Process(accum export.Accumulation) error {
	// Note: the removed labels are returned and ignored here.
	// Conceivably these inputs could be useful to a sampler.
	reduced, _ := accum.Labels().Filter(
		p.filterSelector.LabelFilterFor(
			accum.Descriptor(),
		),
	)
	return p.Checkpointer.Process(
		export.NewAccumulation(
			accum.Descriptor(),
			&reduced,
			accum.Resource(),
			accum.Aggregator(),
		),
	)
}
