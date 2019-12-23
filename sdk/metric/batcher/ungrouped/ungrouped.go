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
		labels     export.Labels
	}

	batchMap map[batchKey]batchValue
)

var _ export.Batcher = &Batcher{}
var _ export.CheckpointSet = batchMap{}

func New(selector export.AggregationSelector, stateful bool) *Batcher {
	return &Batcher{
		selector: selector,
		batchMap: batchMap{},
		stateful: stateful,
	}
}

func (b *Batcher) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	return b.selector.AggregatorFor(descriptor)
}

func (b *Batcher) Process(_ context.Context, record export.Record) error {
	desc := record.Descriptor()
	key := batchKey{
		descriptor: desc,
		encoded:    record.Labels().Encoded(),
	}
	agg := record.Aggregator()
	value, ok := b.batchMap[key]
	if ok {
		// Note: The call to Merge here combines only
		// identical records.  It is required even for a
		// stateless Batcher because such identical records
		// may arise in the Meter implementation due to race
		// conditions.
		return value.aggregator.Merge(agg, desc)
	}
	// If this Batcher is stateful, create a copy of the
	// Aggregator for long-term storage.  Otherwise the
	// Meter implementation will checkpoint the aggregator
	// again, overwriting the long-lived state.
	if b.stateful {
		tmp := agg
		// Note: the call to AggregatorFor() followed by Merge
		// is effectively a Clone() operation.
		agg = b.AggregatorFor(desc)
		if err := agg.Merge(tmp, desc); err != nil {
			return err
		}
	}
	b.batchMap[key] = batchValue{
		aggregator: agg,
		labels:     record.Labels(),
	}
	return nil
}

func (b *Batcher) CheckpointSet() export.CheckpointSet {
	return b.batchMap
}

func (b *Batcher) FinishedCollection() {
	if !b.stateful {
		b.batchMap = batchMap{}
	}
}

func (c batchMap) ForEach(f func(export.Record)) {
	for key, value := range c {
		f(export.NewRecord(
			key.descriptor,
			value.labels,
			value.aggregator,
		))
	}
}
