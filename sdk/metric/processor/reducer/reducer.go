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
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type (
	Processor struct {
		next           export.Processor
		filterSelector LabelFilterSelector
	}

	LabelFilterSelector interface {
		LabelFilterFor(*metric.Descriptor) label.Filter
	}
)

var _ export.Processor = &Processor{}

func New(filterSelector LabelFilterSelector, next export.Processor) *Processor {
	return &Processor{
		next:           next,
		filterSelector: filterSelector,
	}
}

func (p *Processor) Process(accum export.Accumulation) error {
	// Note: the removed labels are returned and ignored here.
	// Conceivably these inputs could be useful to a sampler.
	reduced, _ := accum.Labels().Filter(
		p.filterSelector.LabelFilterFor(
			accum.Descriptor(),
		),
	)
	return p.next.Process(
		export.NewAccumulation(
			accum.Descriptor(),
			&reduced,
			accum.Resource(),
			accum.Aggregator(),
		),
	)
}

func (p *Processor) AggregatorFor(desc *metric.Descriptor, aggs ...*export.Aggregator) {
	p.next.AggregatorFor(desc, aggs...)
}
