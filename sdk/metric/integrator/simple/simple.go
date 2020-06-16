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

		// lock protects the remaining fields, synchronizes computing the
		// checkpoint state on the fly.
		lock sync.Mutex

		// updated indicates the last sequence number when this value had
		// Process() called by an accumulator.
		updated int64

		// stateful indicates that the last-value of the aggregation (since
		// process start time) is being maintained.
		stateful bool

		current    export.Aggregator // refers to single-accumulator checkpoint or delta.
		delta      export.Aggregator // owned if multi accumulator else nil.
		cumulative export.Aggregator // owned if stateful else nil.
	}

	state struct {
		// RWMutex implements locking for the `CheckpointSet` interface.
		sync.RWMutex
		sequence   int64
		processSeq int64

		foreachLock sync.Mutex
		foreachSeq  int64

		processStart  time.Time
		intervalStart time.Time
		intervalEnd   time.Time
		values        map[stateKey]*stateValue
	}
)

var _ export.Integrator = &Integrator{}
var _ export.CheckpointSet = &state{}

func New(aselector export.AggregationSelector, eselector export.ExportKindSelector) *Integrator {
	now := time.Now()
	return &Integrator{
		AggregationSelector: aselector,
		ExportKindSelector:  eselector,
		state: state{
			values:        map[stateKey]*stateValue{},
			processStart:  now,
			intervalStart: now,
			processSeq:    -1,
			foreachSeq:    -1,
		},
	}
}

func (b *Integrator) Process(accum export.Accumulation) error {
	if b.processSeq != b.sequence {
		// This is the beginning of the end of the interval.
		b.processSeq = b.sequence
		b.state.intervalEnd = time.Now()
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
		stateful := b.ExportKindFor(desc).MemoryRequired(desc.MetricKind())

		newValue := &stateValue{
			labels:   accum.Labels(),
			resource: accum.Resource(),
			updated:  b.state.sequence,
			stateful: stateful,
			current:  agg,
		}
		if stateful {
			b.AggregatorFor(desc, &newValue.cumulative)
		}
		b.state.values[key] = newValue
		return nil
	}

	value.lock.Lock()
	defer value.lock.Unlock()

	// Advance the update sequence number:
	sameRound := b.state.sequence == value.updated
	value.updated = b.state.sequence

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

		if value.stateful && value.delta == nil {
			b.AggregationSelector.AggregatorFor(desc, &value.delta)
		}

		return nil
	}
	if value.delta == nil {
		// Allocate the delta.
		b.AggregationSelector.AggregatorFor(desc, &value.delta)
	}
	if value.current != value.delta {
		// Merging two values, first copy the singleton.
		value.current.SynchronizedCopy(value.delta, desc)
		value.current = value.delta
	}
	return value.delta.Merge(agg, desc)
}

func (b *Integrator) CheckpointSet() export.CheckpointSet {
	return &b.state
}

func (b *Integrator) FinishedCollection() {
	b.foreachLock.Lock()
	needForeach := b.foreachSeq != b.sequence
	b.foreachLock.Unlock()

	if needForeach {
		// Note: Users are not expected to skip ForEach, but
		// if they do they may lose errors here.
		_ = b.ForEach(nil, func(_ export.Record) error { return nil })
	}

	b.state.sequence++
	b.state.intervalStart = b.state.intervalEnd
	b.state.intervalEnd = time.Time{}
}

func (b *state) ForEach(exporter export.ExportKindSelector, f func(export.Record) error) error {
	b.foreachLock.Lock()
	firstTime := b.foreachSeq != b.sequence
	b.foreachSeq = b.sequence
	b.foreachLock.Unlock()

	for key, value := range b.values {
		mkind := key.descriptor.MetricKind()
		value.lock.Lock()

		if !value.stateful && value.updated != b.sequence {
			delete(b.values, key)
			continue
		}

		if firstTime && value.stateful {
			// TODO: SumObserver monotonicity should be
			// tested here.  This is a case where memory
			// is required even for a cumulative exporter.

			var err error
			if mkind.PrecomputedSum() {
				// We need to compute a delta.  We have the last value.
				if subt, ok := value.current.(export.Subtractor); ok {
					err = subt.Subtract(value.cumulative, value.delta, key.descriptor)
				} else {
					err = aggregation.ErrNoSubtraction
				}
			} else {
				err = value.cumulative.Merge(value.current, key.descriptor)
			}
			if err != nil {
				value.lock.Unlock()
				return err
			}
		}

		var agg aggregation.Aggregation
		var start time.Time

		switch exporter.ExportKindFor(key.descriptor) {
		case export.PassThroughExporter:
			// No state is required, pass through the checkpointed value.
			agg = value.current

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
				agg = value.cumulative
			} else {
				agg = value.current
			}
			start = b.processStart

		case export.DeltaExporter:

			if mkind.PrecomputedSum() {
				agg = value.delta
			} else {
				agg = value.cumulative
			}
			start = b.intervalStart
		}

		value.lock.Unlock()

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

func (b *stateValue) String() string {
	return fmt.Sprintf("%v %v %v", b.current, b.updated, b.stateful)
}
