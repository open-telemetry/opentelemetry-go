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

package defaultkeys // import "go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
)

type (
	Batcher struct {
		selector export.AggregationSelector
		lencoder export.LabelEncoder
		stateful bool
		dki      dkiMap
		agg      aggMap
	}

	aggEntry struct {
		aggregator export.Aggregator
		descriptor *export.Descriptor
		labels     []core.KeyValue
	}

	dkiMap map[*export.Descriptor]map[core.Key]int
	aggMap map[string]aggEntry

	producer struct {
		aggMap   aggMap
		lencoder export.LabelEncoder
	}
)

var _ export.Batcher = &Batcher{}
var _ export.Producer = &producer{}

func New(selector export.AggregationSelector, lencoder export.LabelEncoder, stateful bool) *Batcher {
	return &Batcher{
		selector: selector,
		lencoder: lencoder,
		dki:      dkiMap{},
		agg:      aggMap{},
		stateful: stateful,
	}
}

func (b *Batcher) AggregatorFor(descriptor *export.Descriptor) export.Aggregator {
	return b.selector.AggregatorFor(descriptor)
}

func (b *Batcher) Process(_ context.Context, desc *export.Descriptor, labels []core.KeyValue, _ string, _ export.LabelEncoder, agg export.Aggregator) error {
	keys := desc.Keys()

	// Cache the mapping from Descriptor->Key->Index
	ki, ok := b.dki[desc]
	if !ok {
		ki = map[core.Key]int{}
		b.dki[desc] = ki

		for i, k := range keys {
			ki[k] = i
		}
	}

	// Compute the value list.  Note: Unspecified values become
	// empty strings.  TODO: pin this down, we have no appropriate
	// Value constructor.
	canon := make([]core.KeyValue, len(keys))

	for i, key := range keys {
		canon[i] = key.String("")
	}

	// Note also the possibility to speed this computation of
	// "encoded" via "canon" in the form of a (Descriptor,
	// LabelSet)->(Labels, Encoded) cache.
	for _, kv := range labels {
		pos, ok := ki[kv.Key]
		if !ok {
			continue
		}
		canon[pos].Value = kv.Value
	}

	// Compute an encoded lookup key.
	encoded := b.lencoder.EncodeLabels(canon)

	// Reduce dimensionality.
	rag, ok := b.agg[encoded]
	if !ok {
		b.agg[encoded] = aggEntry{
			aggregator: agg,
			labels:     canon,
			descriptor: desc,
		}
		return nil
	}
	return rag.aggregator.Merge(agg, desc)
}

func (b *Batcher) ReadCheckpoint() export.Producer {
	checkpoint := b.agg
	if !b.stateful {
		b.agg = aggMap{}
	}
	return &producer{
		aggMap:   checkpoint,
		lencoder: b.lencoder,
	}
}

func (p *producer) Foreach(f func(export.Record)) {
	for encoded, entry := range p.aggMap {
		f(export.NewRecord(entry.aggregator,
			entry.descriptor,
			entry.labels,
			p.lencoder,
			encoded,
		))
	}
}
