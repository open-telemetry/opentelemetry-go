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
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	Integrator struct {
		kind export.ExporterKind

		export.AggregationSelector

		sequence     int64
		checkpointed int64

		state
	}

	stateKey struct {
		descriptor *metric.Descriptor
		distinct   label.Distinct
		resource   label.Distinct
	}

	stateValue struct {
		aggregator export.Aggregator
		labels     *label.Set
		resource   *resource.Resource
		updated    int64
		stateful   bool
		aggOwned   bool
	}

	state struct {
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex
		values map[stateKey]*stateValue
	}
)

var _ export.Integrator = &Integrator{}
var _ export.CheckpointSet = &state{}

func New(selector export.AggregationSelector, kind export.ExporterKind) *Integrator {
	return &Integrator{
		AggregationSelector: selector,
		kind:                kind,
		checkpointed:        -1,
		state: state{
			values: map[stateKey]*stateValue{},
		},
	}
}

func (b *Integrator) cloneCheckpoint(checkpointed export.Aggregator, desc *metric.Descriptor) (export.Aggregator, error) {
	agg := b.AggregatorFor(desc)
	if err := agg.Merge(checkpointed, desc); err != nil {
		return nil, err
	}
	return agg, nil
}

// cloneReplace is used to replace the current aggregation with another, leaving the
// checkpoint unchanged.  This is for use with asynchronous instruments when an overwrite
// occurs.
func (b *Integrator) cloneReplace(value *stateValue, replace export.Aggregator, desc *metric.Descriptor) error {
	if !desc.MetricKind().Asynchronous() {
		return fmt.Errorf("inconsistent integrator state")
	}
	value.aggregator.Checkpoint(context.Background(), desc)
	cstate, err := b.cloneCheckpoint(value.aggregator, desc)
	if err != nil {
		return err
	}
	value.aggregator = cstate

	return cstate.Merge(replace, desc)
}

func (b *Integrator) Process(record export.Record) error {
	desc := record.Descriptor()
	key := stateKey{
		descriptor: desc,
		distinct:   record.Labels().Equivalent(),
		resource:   record.Resource().Equivalent(),
	}
	stateful := b.kind.MemoryRequired(*desc)
	agg := record.Aggregator()

	// Check if there is an existing record.  If so, update it.
	if value, ok := b.state.values[key]; ok {
		// Advance the update sequence number:
		sameRound := b.sequence == value.updated
		value.updated = b.sequence

		// An existing record will be found when:
		// (a) stateful aggregation is required for an exporter
		if !sameRound {
			// This is the first record in the current checkpoint set.
			if !stateful && !value.aggOwned {
				// The first time through, refer to a checkpointed
				// aggregator.
				value.aggregator = agg
				return nil
			}
			// This can be synchronous or asynchronous.
			err := value.aggregator.Merge(agg, desc)
			return err
		}
		// (b) multiple accumulators (SDKs) are being used.
		// Another accumulator must have produced this.
		if desc.MetricKind().Asynchronous() && !stateful {
			// The last value across multiple accumulators is taken.
			value.aggregator = agg
			return nil
		}
		// Clone the (synchronous or stateful asynchronous) record.
		if !value.aggOwned {
			clone, err := b.cloneCheckpoint(value.aggregator, desc)
			if err != nil {
				return err
			}
			value.aggregator = clone
			value.aggOwned = true
		}
		// Synchronous case: Merge with the prior aggregation.
		if desc.MetricKind().Synchronous() {
			return value.aggregator.Merge(agg, desc)
		}
		// Asynchronous case: Replace the current value.
		return b.cloneReplace(value, agg, desc)
	}

	// There was no existing record.
	newValue := &stateValue{
		aggregator: agg,
		labels:     record.Labels(),
		resource:   record.Resource(),
		stateful:   stateful,
		updated:    b.sequence,
	}
	if stateful {
		var err error
		newValue.aggregator, err = b.cloneCheckpoint(agg, desc)
		newValue.aggOwned = true
		if err != nil {
			return err
		}
	}
	b.state.values[key] = newValue
	return nil
}

func (b *Integrator) CheckpointSet() export.CheckpointSet {
	if b.checkpointed != b.sequence {
		b.checkpointed = b.sequence
		for key, value := range b.state.values {
			if value.aggOwned {
				value.aggregator.Merge(value.aggregator, key.descriptor)
				value.aggregator.Checkpoint(context.Background(), key.descriptor)
			}
			if !value.stateful {
				value.aggOwned = false
			}
		}
	}
	return &b.state
}

func (b *Integrator) FinishedCollection() {
	b.sequence++
	// TODO For a DeltaExporter it's likely faster to create a new map and
	// copy only the cumulative instrument state.
	for key, value := range b.state.values {
		if !value.stateful {
			delete(b.state.values, key)
			continue
		}
	}
}

func (b *state) ForEach(_ export.ExporterKind, f func(export.Record) error) error {
	// @@@ Use the kind
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

func (b *stateValue) String() string {
	return fmt.Sprintf("%v %v %v %v", b.aggregator, b.updated, b.stateful, b.aggOwned)
}
