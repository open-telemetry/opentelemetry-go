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

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/resource"
)

// @@@ Need a systematic test of all the instruments / aggregators.

type (
	Integrator struct {
		kind export.ExporterKind

		export.AggregationSelector

		state
	}

	stateKey struct {
		descriptor *metric.Descriptor
		distinct   label.Distinct
		resource   label.Distinct
	}

	stateValue struct {
		// labels corresponds to the stateKey.distinct field.
		labels *label.Set

		// resource corresponds to the stateKey.resource field.
		resource *resource.Resource

		// lock protects the remaining fields, synchronizes computing the
		// checkpoint state on the fly.
		lock sync.Mutex

		aggregator export.Aggregator

		// updated indicates the last sequence number when this value had
		// Process() called by an accumulator.
		updated      int64
		checkpointed int64

		// stateful indicates that the last-value of the aggregation (since
		// process start time) is being maintained.
		stateful bool

		// aggOwned is always true for stateful aggregators.  aggOWned may also be
		// true for stateless synchronous aggregators, either the set of labels is
		// reduced (by in-process aggregation) or multiple accumulators are used.
		//
		// When aggOwned is true, the current accumulated value is held in the
		// aggregator's current register.  If the aggregator is stateless, the
		// last stateful value is held in the checkpoint register.
		aggOwned bool
	}

	state struct {
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex
		sequence int64
		values   map[stateKey]*stateValue
	}
)

var _ export.Integrator = &Integrator{}
var _ export.CheckpointSet = &state{}

func New(selector export.AggregationSelector, kind export.ExporterKind) *Integrator {
	return &Integrator{
		AggregationSelector: selector,
		kind:                kind,
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
	value.aggregator.Checkpoint(desc)
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
		value.lock.Lock()
		defer value.lock.Unlock()

		// Advance the update sequence number:
		sameRound := b.state.sequence == value.updated
		value.updated = b.state.sequence

		// An existing record will be found when:
		// (a) stateful aggregation is required for an exporter
		if !sameRound {
			if stateful {
				// The prior stateful value is in the checkpoint register, and
				// the last accumulator value is in the current register.
				value.aggregator.Swap()
				value.aggregator.Checkpoint(desc)
			}
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
		aggregator:   agg,
		labels:       record.Labels(),
		resource:     record.Resource(),
		stateful:     stateful,
		checkpointed: -1,
		updated:      b.state.sequence,
		aggOwned:     false,
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
	return &b.state
}

func (b *Integrator) FinishedCollection() {
	b.state.sequence++
}

func (b *state) ForEach(_ export.ExporterKind, f func(export.Record) error) error {
	for key, value := range b.values {
		value.lock.Lock()

		if !value.stateful && value.updated != b.sequence {
			delete(b.values, key)
			continue
		}

		if value.checkpointed != b.sequence {
			value.checkpointed = b.sequence
			if value.stateful {
				// Accumulated value in current; last value in checkpoint.
				value.aggregator.Swap()

				// Last value in current, accumulated value in checkpoint:
				// add into current.
				err := value.aggregator.Merge(value.aggregator, key.descriptor)
				if err != nil {
					return err
				}

				// Place up-to-date value in checkpoint, accumulated value in current.
				value.aggregator.Swap()
			} else {
				// In this case, we'll the accumulated value (otherwise
				// we'd have state).
				if value.aggOwned {
					value.aggregator.Checkpoint(key.descriptor)
					value.aggOwned = false
				}
			}
		}

		value.lock.Unlock()

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
