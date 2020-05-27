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

package test

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
)

type mapkey struct {
	desc     *metric.Descriptor
	distinct label.Distinct
}

type CheckpointSet struct {
	sync.RWMutex
	records  map[mapkey]export.Record
	updates  []export.Record
	resource *resource.Resource
}

// NewCheckpointSet returns a test CheckpointSet that new records could be added.
// Records are grouped by their encoded labels.
func NewCheckpointSet(resource *resource.Resource) *CheckpointSet {
	return &CheckpointSet{
		records:  make(map[mapkey]export.Record),
		resource: resource,
	}
}

func (p *CheckpointSet) Reset() {
	p.records = make(map[mapkey]export.Record)
	p.updates = nil
}

// Add a new descriptor to a Checkpoint.
//
// If there is an existing record with the same descriptor and labels,
// the stored aggregator will be returned and should be merged.
func (p *CheckpointSet) Add(desc *metric.Descriptor, newAgg export.Aggregator, labels ...kv.KeyValue) (agg export.Aggregator, added bool) {
	elabels := label.NewSet(labels...)

	key := mapkey{
		desc:     desc,
		distinct: elabels.Equivalent(),
	}
	if record, ok := p.records[key]; ok {
		return record.Aggregator(), false
	}

	rec := export.NewRecord(desc, &elabels, p.resource, newAgg)
	p.updates = append(p.updates, rec)
	p.records[key] = rec
	return newAgg, true
}

func createNumber(desc *metric.Descriptor, v float64) metric.Number {
	if desc.NumberKind() == metric.Float64NumberKind {
		return metric.NewFloat64Number(v)
	}
	return metric.NewInt64Number(int64(v))
}

func (p *CheckpointSet) AddLastValue(desc *metric.Descriptor, v float64, labels ...kv.KeyValue) {
	p.updateAggregator(desc, lastvalue.New(), v, labels...)
}

func (p *CheckpointSet) AddCounter(desc *metric.Descriptor, v float64, labels ...kv.KeyValue) {
	p.updateAggregator(desc, sum.New(), v, labels...)
}

func (p *CheckpointSet) AddValueRecorder(desc *metric.Descriptor, v float64, labels ...kv.KeyValue) {
	p.updateAggregator(desc, array.New(), v, labels...)
}

func (p *CheckpointSet) AddHistogramValueRecorder(desc *metric.Descriptor, boundaries []float64, v float64, labels ...kv.KeyValue) {
	p.updateAggregator(desc, histogram.New(desc, boundaries), v, labels...)
}

func (p *CheckpointSet) updateAggregator(desc *metric.Descriptor, newAgg export.Aggregator, v float64, labels ...kv.KeyValue) {
	ctx := context.Background()
	// Updates and checkpoint the new aggregator
	_ = newAgg.Update(ctx, createNumber(desc, v), desc)
	newAgg.Checkpoint(ctx, desc)

	// Try to add this aggregator to the CheckpointSet
	agg, added := p.Add(desc, newAgg, labels...)
	if !added {
		// An aggregator already exist for this descriptor and label set, we should merge them.
		_ = agg.Merge(newAgg, desc)
	}
}

func (p *CheckpointSet) ForEach(f func(export.Record) error) error {
	for _, r := range p.updates {
		if err := f(r); err != nil && !errors.Is(err, aggregator.ErrNoData) {
			return err
		}
	}
	return nil
}
