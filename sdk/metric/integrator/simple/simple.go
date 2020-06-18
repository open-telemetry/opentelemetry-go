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
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
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
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex
		values map[batchKey]batchValue

		// Note: the timestamp logic currently assumes all
		// exports are deltas.

		intervalStart time.Time
		intervalEnd   time.Time

		// startedCollection and finishedCollection are the
		// number of StartCollection() and FinishCollection()
		// calls, used to ensure that the sequence of starts
		// and finishes are correctly balanced.

		startedCollection  int64
		finishedCollection int64
	}
)

var _ export.Integrator = &Integrator{}
var _ export.CheckpointSet = &batch{}
var ErrInconsistentState = fmt.Errorf("inconsistent integrator state")

func New(selector export.AggregationSelector, stateful bool) *Integrator {
	return &Integrator{
		AggregationSelector: selector,
		stateful:            stateful,
		batch: batch{
			values:        map[batchKey]batchValue{},
			intervalStart: time.Now(),
		},
	}
}

func (b *Integrator) Process(accumulation export.Accumulation) error {
	if b.startedCollection != b.finishedCollection+1 {
		return ErrInconsistentState
	}

	desc := accumulation.Descriptor()
	key := batchKey{
		descriptor: desc,
		distinct:   accumulation.Labels().Equivalent(),
		resource:   accumulation.Resource().Equivalent(),
	}
	agg := accumulation.Aggregator()
	value, ok := b.batch.values[key]
	if ok {
		// Note: The call to Merge here combines only
		// identical accumulations.  It is required even for a
		// stateless Integrator because such identical accumulations
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
		b.AggregatorFor(desc, &agg)
		if err := agg.Merge(tmp, desc); err != nil {
			return err
		}
	}
	b.batch.values[key] = batchValue{
		aggregator: agg,
		labels:     accumulation.Labels(),
		resource:   accumulation.Resource(),
	}
	return nil
}

func (b *Integrator) CheckpointSet() export.CheckpointSet {
	return &b.batch
}

func (b *Integrator) StartCollection() {
	if b.startedCollection != 0 {
		b.intervalStart = b.intervalEnd
	}
	b.startedCollection++
	if !b.stateful {
		b.batch.values = map[batchKey]batchValue{}
	}
}

func (b *Integrator) FinishCollection() error {
	b.finishedCollection++
	b.intervalEnd = time.Now()
	if b.startedCollection != b.finishedCollection {
		return ErrInconsistentState
	}
	return nil
}

func (b *batch) ForEach(f func(export.Record) error) error {
	if b.startedCollection != b.finishedCollection {
		return ErrInconsistentState
	}

	for key, value := range b.values {
		if err := f(export.NewRecord(
			key.descriptor,
			value.labels,
			value.resource,
			value.aggregator.Aggregation(),
			b.intervalStart,
			b.intervalEnd,
		)); err != nil && !errors.Is(err, aggregation.ErrNoData) {
			return err
		}
	}
	return nil
}
