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

package simple // import "go.opentelemetry.io/otel/sdk/metric/integrator/simple"

import (
	"context"
	"errors"
	"sync"

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Integrator struct {
		export.AggregationSelector
		stateful bool
		batch
	}

	batchKey struct {
		descriptor *metric.Descriptor
		distinct   label.Distinct
		resource   label.Distinct
	}

	batchValue struct {
		aggregator export.Aggregator
		labels     *label.Set
		resource   *resource.Resource
	}

	batch struct {
		// RWMutex implements locking for the `CheckpoingSet` interface.
		sync.RWMutex
		values map[batchKey]batchValue
	}
)

var _ export.Integrator = &Integrator{}
var _ export.CheckpointSet = &batch{}

func New(selector export.AggregationSelector, stateful bool) *Integrator {
	return &Integrator{
		AggregationSelector: selector,
		stateful:            stateful,
		batch: batch{
			values: map[batchKey]batchValue{},
		},
	}
}

func (b *Integrator) Process(_ context.Context, record export.Record) error {
	desc := record.Descriptor()
	key := batchKey{
		descriptor: desc,
		distinct:   record.Labels().Equivalent(),
		resource:   record.Resource().Equivalent(),
	}
	agg := record.Aggregator()
	value, ok := b.batch.values[key]
	if ok {
		// Note: The call to Merge here combines only
		// identical records.  It is required even for a
		// stateless Integrator because such identical records
		// may arise in the Meter implementation due to race
		// conditions.
		return value.aggregator.Merge(agg, desc)
	}
	// If this integrator is stateful, create a copy of the
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
	b.batch.values[key] = batchValue{
		aggregator: agg,
		labels:     record.Labels(),
		resource:   record.Resource(),
	}
	return nil
}

func (b *Integrator) CheckpointSet() export.CheckpointSet {
	return &b.batch
}

func (b *Integrator) FinishedCollection() {
	if !b.stateful {
		b.batch.values = map[batchKey]batchValue{}
	}
}

func (b *batch) ForEach(f func(export.Record) error) error {
	for key, value := range b.values {
		if err := f(export.NewRecord(
			key.descriptor,
			value.labels,
			value.resource,
			value.aggregator,
		)); err != nil && !errors.Is(err, aggregator.ErrNoData) {
			return err
		}
	}
	return nil
}
