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
		export.ExportKindSelector
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

		// updated indicates the last sequence number when this value had
		// Process() called by an accumulator.
		updated int64

		// stateful indicates that a cumulative aggregation is
		// being maintained, taken from the process start time.
		stateful bool

		current    export.Aggregator // refers to single-accumulator checkpoint or delta.
		delta      export.Aggregator // owned if multi accumulator else nil.
		cumulative export.Aggregator // owned if stateful else nil.
	}

	state struct {
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex
		values map[stateKey]*stateValue

		// Note: the timestamp logic currently assumes all
		// exports are deltas.

		processStart  time.Time
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
var _ export.CheckpointSet = &state{}
var ErrInconsistentState = fmt.Errorf("inconsistent integrator state")

// New returns a basic Integrator using the provided
// AggregationSelector to select Aggregators.  The ExportKindSelector
// is consulted to determine the kind(s) of exporter that will consume
// data, so that this Integrator can prepare to compute Delta or
// Cumulative Aggregations as needed.
func New(aselector export.AggregationSelector, eselector export.ExportKindSelector) *Integrator {
	now := time.Now()
	return &Integrator{
		AggregationSelector: aselector,
		ExportKindSelector:  eselector,
		state: state{
			values:        map[stateKey]*stateValue{},
			processStart:  now,
			intervalStart: now,
		},
	}
}

// Process implements export.Integrator.
func (b *Integrator) Process(accum export.Accumulation) error {
	if b.startedCollection != b.finishedCollection+1 {
		return ErrInconsistentState
	}
	desc := accum.Descriptor()
	key := stateKey{
		descriptor: desc,
		distinct:   accum.Labels().Equivalent(),
		resource:   accum.Resource().Equivalent(),
	}
	agg := accum.Aggregator()

	// Check if there is an existing record.
	value, ok := b.state.values[key]
	if !ok {
		stateful := b.ExportKindFor(desc, agg.Aggregation().Kind()).MemoryRequired(desc.MetricKind())

		newValue := &stateValue{
			labels:   accum.Labels(),
			resource: accum.Resource(),
			updated:  b.state.finishedCollection,
			stateful: stateful,
			current:  agg,
		}
		if stateful {
			// If stateful, allocate a cumulative aggregator.
			b.AggregatorFor(desc, &newValue.cumulative)

			if desc.MetricKind().PrecomputedSum() {
				// If we need to compute deltas, allocate another aggregator.
				b.AggregatorFor(desc, &newValue.delta)
			}
		}
		b.state.values[key] = newValue
		return nil
	}

	// Advance the update sequence number:
	sameRound := b.state.finishedCollection == value.updated
	value.updated = b.state.finishedCollection

	// An existing record will be found when:
	// (a) stateful aggregation is required for an exporter
	// (b) multiple accumulators (SDKs) are being used.
	// Another accumulator must have produced this.

	if !sameRound {
		// This is the first time through in a new round.
		value.current = agg
		return nil
	}
	if desc.MetricKind().Asynchronous() {
		// The last value across multiple accumulators is taken.
		value.current = agg
		return nil
	}
	if value.delta == nil {
		// Merging values: may need to allocate the delta aggregator.
		b.AggregationSelector.AggregatorFor(desc, &value.delta)
	}
	if value.current != value.delta {
		// Merging two values, first copy the singleton.
		err := value.current.SynchronizedCopy(value.delta, desc)
		if err != nil {
			return err
		}
		value.current = value.delta
	}
	return value.delta.Merge(agg, desc)
}

// CheckpointSet returns the associated CheckpointSet.  Use the
// CheckpointSet Locker interface to synchronize access to this
// object.  The CheckpointSet.ForEach() method cannot be called
// concurrently with Process().
func (b *Integrator) CheckpointSet() export.CheckpointSet {
	return &b.state
}

// StartCollection signals to the Integrator one or more Accumulators
// will begin calling Process() calls during collection.
func (b *Integrator) StartCollection() {
	if b.startedCollection != 0 {
		b.intervalStart = b.intervalEnd
	}
	b.startedCollection++
}

// FinishCollection signals to the Integrator that a complete round of
// collection is finished and that ForEach will be called to access
// the CheckpointSet.
func (b *Integrator) FinishCollection() error {
	b.intervalEnd = time.Now()
	if b.startedCollection != b.finishedCollection+1 {
		return ErrInconsistentState
	}
	defer func() { b.finishedCollection++ }()

	for key, value := range b.values {
		mkind := key.descriptor.MetricKind()

		if !value.stateful {
			if value.updated != b.finishedCollection {
				delete(b.values, key)
			}
			continue
		}

		var err error
		if mkind.PrecomputedSum() {
			// We need to compute a delta.  We have the prior cumulative value.
			if subt, ok := value.current.(export.Subtractor); ok {
				err = subt.Subtract(value.cumulative, value.delta, key.descriptor)

				if err == nil {
					err = value.current.SynchronizedCopy(value.cumulative, key.descriptor)
				}
			} else {
				err = aggregation.ErrNoSubtraction
			}
		} else {
			err = value.cumulative.Merge(value.current, key.descriptor)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ForEach iterates through the CheckpointSet, for passing an
// export.Record with the appropriate Cumulative or Delta aggregation
// to an exporter.
func (b *state) ForEach(exporter export.ExportKindSelector, f func(export.Record) error) error {
	if b.startedCollection != b.finishedCollection {
		return ErrInconsistentState
	}
	for key, value := range b.values {
		mkind := key.descriptor.MetricKind()

		var agg aggregation.Aggregation
		var start time.Time

		switch exporter.ExportKindFor(key.descriptor, value.current.Aggregation().Kind()) {
		case export.PassThroughExporter:
			// No state is required, pass through the checkpointed value.
			agg = value.current.Aggregation()

			if mkind.PrecomputedSum() {
				start = b.processStart
			} else {
				start = b.intervalStart
			}

		case export.CumulativeExporter:
			// If stateful, the sum has been computed.  If stateless, the
			// input was already cumulative.  Either way, use the checkpointed
			// value:
			if value.stateful {
				agg = value.cumulative.Aggregation()
			} else {
				agg = value.current.Aggregation()
			}
			start = b.processStart

		case export.DeltaExporter:
			// Precomputed sums are a special case.
			if mkind.PrecomputedSum() {
				agg = value.delta.Aggregation()
			} else {
				agg = value.current.Aggregation()
			}
			start = b.intervalStart
		}

		if err := f(export.NewRecord(
			key.descriptor,
			value.labels,
			value.resource,
			agg,
			start,
			b.intervalEnd,
		)); err != nil && !errors.Is(err, aggregation.ErrNoData) {
			return err
		}
	}
	return nil
}
