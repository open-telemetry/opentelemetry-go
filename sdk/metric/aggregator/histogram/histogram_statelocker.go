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

package histogram // import "go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"

import (
	"context"
	"sort"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/internal"
)

type (
	// AggregatorSL observe events and counts them in pre-determined buckets.
	// It also calculates the sum and count of all events.
	AggregatorSL struct {
		// This aggregator uses the StateLocker that enables a lock-free Update()
		// in exchange of a blocking and consistent Checkpoint(). Since Checkpoint()
		// is called by the sdk itself and it is not part of a hot path,
		// the user is not impacted by these blocking calls.
		//
		// The algorithm keeps two states. At every instance of time there exist one current state,
		// in which new updates are aggregated, and one checkpoint state, that represents the state
		// since the last Checkpoint(). These states are swapped when a `Checkpoint()` occur.

		// states needs to be aligned for 64-bit atomic operations.
		states     [2]stateSL
		lock       internal.StateLocker
		boundaries []core.Number
		kind       core.NumberKind
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	stateSL struct {
		// all fields have to be aligned for 64-bit atomic operations.
		buckets aggregator.Buckets
		count   core.Number
		sum     core.Number
	}
)

var _ export.Aggregator = &AggregatorSL{}
var _ aggregator.Sum = &AggregatorSL{}
var _ aggregator.Count = &AggregatorSL{}
var _ aggregator.Histogram = &AggregatorSL{}

// New returns a new measure aggregator for computing Histograms.
//
// A Histogram observe events and counts them in pre-defined buckets.
// And also provides the total sum and count of all observations.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func NewSL(desc *metric.Descriptor, boundaries []core.Number) *AggregatorSL {
	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := numbers{
		numbers: make([]core.Number, len(boundaries)),
		kind:    desc.NumberKind(),
	}

	copy(sortedBoundaries.numbers, boundaries)
	sort.Sort(&sortedBoundaries)
	boundaries = sortedBoundaries.numbers

	agg := AggregatorSL{
		kind:       desc.NumberKind(),
		boundaries: boundaries,
		states: [2]stateSL{
			{
				buckets: aggregator.Buckets{
					Boundaries: boundaries,
					Counts:     make([]core.Number, len(boundaries)+1),
				},
			},
			{
				buckets: aggregator.Buckets{
					Boundaries: boundaries,
					Counts:     make([]core.Number, len(boundaries)+1),
				},
			},
		},
	}
	return &agg
}

// Sum returns the sum of all values in the checkpoint.
func (c *AggregatorSL) Sum() (core.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint().sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *AggregatorSL) Count() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return int64(c.checkpoint().count), nil
}

// Histogram returns the count of events in pre-determined buckets.
func (c *AggregatorSL) Histogram() (aggregator.Buckets, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint().buckets, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *AggregatorSL) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	c.lock.SwapActiveState(c.resetCheckpoint)
}

// checkpoint returns the checkpoint state by inverting the lower bit of generationAndHotIdx.
func (c *AggregatorSL) checkpoint() *stateSL {
	return &c.states[c.lock.ColdIdx()]
}

func (c *AggregatorSL) resetCheckpoint() {
	checkpoint := c.checkpoint()

	checkpoint.count.SetUint64(0)
	checkpoint.sum.SetNumber(core.Number(0))
	checkpoint.buckets.Counts = make([]core.Number, len(checkpoint.buckets.Counts))
}

// Update adds the recorded measurement to the current data set.
func (c *AggregatorSL) Update(_ context.Context, number core.Number, desc *metric.Descriptor) error {
	kind := desc.NumberKind()

	cIdx := c.lock.Start()
	defer c.lock.End(cIdx)

	current := &c.states[cIdx]
	current.count.AddUint64Atomic(1)
	current.sum.AddNumberAtomic(kind, number)

	for i, boundary := range c.boundaries {
		if number.CompareNumber(kind, boundary) < 0 {
			current.buckets.Counts[i].AddUint64Atomic(1)
			return nil
		}
	}

	// Observed event is bigger than all defined boundaries.
	current.buckets.Counts[len(c.boundaries)].AddUint64Atomic(1)

	return nil
}

// Merge combines two histograms that have the same buckets into a single one.
func (c *AggregatorSL) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*AggregatorSL)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	// Lock() synchronize Merge() and Checkpoint() to make sure all operations of
	// Merge() is done to the same state.
	c.lock.Lock()
	defer c.lock.Unlock()

	current := c.checkpoint()
	// We assume that the aggregator being merged is not being updated nor checkpointed or this could be inconsistent.
	ocheckpoint := o.checkpoint()

	current.sum.AddNumber(desc.NumberKind(), ocheckpoint.sum)
	current.count.AddNumber(core.Uint64NumberKind, ocheckpoint.count)

	for i := 0; i < len(current.buckets.Counts); i++ {
		current.buckets.Counts[i].AddNumber(core.Uint64NumberKind, ocheckpoint.buckets.Counts[i])
	}
	return nil
}
