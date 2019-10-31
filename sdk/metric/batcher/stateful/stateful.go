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

package stateful

import (
	"context"
	"strings"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/sdk/export"
)

type (
	Batcher struct {
		dki      dkiMap
		agg      aggMap
		selector export.MetricAggregationSelector
	}

	aggEntry struct {
		aggregator export.MetricAggregator
		descriptor *export.Descriptor

		// NOTE: When only a single exporter is in use,
		// there's a potential to avoid encoding the labels
		// twice, since this class has to encode them once.
		labels []core.Value
	}

	dkiMap map[*export.Descriptor]map[core.Key]int
	aggMap map[string]aggEntry
)

var _ export.MetricBatcher = &Batcher{}
var _ export.MetricProducer = aggMap{}

func New(selector export.MetricAggregationSelector) *Batcher {
	return &Batcher{
		selector: selector,
		dki:      dkiMap{},
		agg:      aggMap{},
	}
}

func (b *Batcher) AggregatorFor(record export.MetricRecord) export.MetricAggregator {
	return b.selector.AggregatorFor(record)
}

func (b *Batcher) Process(_ context.Context, record export.MetricRecord, agg export.MetricAggregator) {
	desc := record.Descriptor()
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
	// empty strings.  TODO: pin this down.
	canon := make([]core.Value, len(keys))

	for i := 0; i < len(keys); i++ {
		canon[i] = core.String("")
	}

	for _, kv := range record.Labels() {
		pos, ok := ki[kv.Key]
		if !ok {
			continue
		}
		canon[pos] = kv.Value
	}

	// Compute an encoded lookup key.
	//
	// Note the opportunity to use an export-specific
	// representation here, then avoid recomputing it in the
	// exporter.  For example, depending on the exporter, we could
	// use an OpenMetrics representation, a statsd representation,
	// etc.  This only benefits a single exporter, of course.
	//
	// Note also the possibility to speed this computation of
	// "encoded" from "canon" in the form of a (Descriptor,
	// LabelSet)->Encoded cache.
	var sb strings.Builder
	for i := 0; i < len(keys); i++ {
		sb.WriteString(string(keys[i]))
		sb.WriteRune('=')
		sb.WriteString(canon[i].Emit())

		if i < len(keys)-1 {
			sb.WriteRune(',')
		}
	}

	encoded := sb.String()

	// Reduce dimensionality.
	rag, ok := b.agg[encoded]
	if !ok {
		b.agg[encoded] = aggEntry{
			aggregator: agg,
			labels:     canon,
			descriptor: record.Descriptor(),
		}
	} else {
		rag.aggregator.Merge(agg, record.Descriptor())
	}
}

func (b *Batcher) Reset() {
	b.agg = aggMap{}
}

func (b *Batcher) ReadCheckpoint() export.MetricProducer {
	checkpoint := b.agg
	return checkpoint
}

func (c aggMap) Foreach(f func(export.MetricAggregator, *export.Descriptor, []core.Value)) {
	for _, entry := range c {
		f(entry.aggregator, entry.descriptor, entry.labels)
	}
}
