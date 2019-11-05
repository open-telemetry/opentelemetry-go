// Copyright 2019, OpenTelemetry Authors
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

package ungrouped // import "go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type (
	Batcher struct {
		selector export.AggregationSelector
		batchMap batchMap
		stateful bool
	}

	batchKey struct {
		descriptor *export.Descriptor
		encoded    string
	}

	batchValue struct {
		aggregator export.Aggregator
		labels     []core.KeyValue
	}

	batchMap map[batchKey]batchValue
)

var _ export.Batcher = &Batcher{}
var _ export.Producer = batchMap{}

func New(selector export.AggregationSelector, stateful bool) *Batcher {
	return &Batcher{
		selector: selector,
		batchMap: batchMap{},
		stateful: stateful,
	}
}

func (b *Batcher) AggregatorFor(record export.Record) export.Aggregator {
	return b.selector.AggregatorFor(record)
}

func (b *Batcher) Process(_ context.Context, record export.Record, agg export.Aggregator) {
	desc := record.Descriptor()
	key := batchKey{
		descriptor: desc,
		encoded:    record.EncodedLabels(),
	}
	value, ok := b.batchMap[key]
	if !ok {
		b.batchMap[key] = batchValue{
			aggregator: agg,
			labels:     record.Labels(),
		}
	} else {
		value.aggregator.Merge(agg, desc)
	}
}

func (b *Batcher) ReadCheckpoint() export.Producer {
	checkpoint := b.batchMap
	if !b.stateful {
		b.batchMap = batchMap{}
	}
	return checkpoint
}

func (c batchMap) Foreach(f func(export.Aggregator, export.ProducedRecord)) {
	for key, value := range c {
		pr := export.ProducedRecord{
			Descriptor: key.descriptor,
			Labels:     value.labels,
		}
		f(value.aggregator, pr)
	}
}
