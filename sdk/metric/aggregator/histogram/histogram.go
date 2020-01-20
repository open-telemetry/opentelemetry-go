// Copyright 2020, OpenTelemetry Authors
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
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator observe events and counts them in pre-determined buckets.
	// It also calculates the sum and count of all events.
	Aggregator struct {
		// state needs to be aligned for 64-bit atomic operations.
		current state
		// checkpoint needs to be aligned for 64-bit atomic operations.
		checkpoint state
		boundaries []core.Number
		kind       core.NumberKind
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state struct {
		// all fields have to be aligned for 64-bit atomic operations.
		buckets aggregator.Buckets
		count   core.Number
		sum     core.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.Sum = &Aggregator{}
var _ aggregator.Count = &Aggregator{}
var _ aggregator.Histogram = &Aggregator{}

// New returns a new measure aggregator for computing Histograms.
//
// A Histogram observe events and counts them in pre-defined buckets.
// And also provides the total sum and count of all observations.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func New(desc *export.Descriptor, boundaries []core.Number) *Aggregator {
	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := numbers{
		numbers: make([]core.Number, len(boundaries)),
		kind:    desc.NumberKind(),
	}

	copy(sortedBoundaries.numbers, boundaries)
	sort.Sort(&sortedBoundaries)
	boundaries = sortedBoundaries.numbers

	agg := Aggregator{
		kind:       desc.NumberKind(),
		boundaries: boundaries,
		current: state{
			buckets: aggregator.Buckets{
				Boundaries: boundaries,
				Counts:     make([]core.Number, len(boundaries)+1),
			},
		},
		checkpoint: state{
			buckets: aggregator.Buckets{
				Boundaries: boundaries,
				Counts:     make([]core.Number, len(boundaries)+1),
			},
		},
	}
	return &agg
}

// Sum returns the sum of all values in the checkpoint.
func (c *Aggregator) Sum() (core.Number, error) {
	return c.checkpoint.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	return int64(c.checkpoint.count.AsUint64()), nil
}

// Histogram returns the count of events in pre-determined buckets.
func (c *Aggregator) Histogram() (aggregator.Buckets, error) {
	return c.checkpoint.buckets, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *export.Descriptor) {
	// N.B. There is no atomic operation that can update all three
	// values at once without a memory allocation.
	//
	// This aggregator is intended to trade this correctness for
	// speed.
	//
	// Therefore, atomically swap fields independently, knowing
	// that individually the three parts of this aggregation could
	// be spread across multiple collections in rare cases.

	c.checkpoint.count.SetUint64(c.current.count.SwapUint64Atomic(0))
	c.checkpoint.sum = c.current.sum.SwapNumberAtomic(core.Number(0))

	for i := 0; i < len(c.checkpoint.buckets.Counts); i++ {
		c.checkpoint.buckets.Counts[i].SetUint64(c.current.buckets.Counts[i].SwapUint64Atomic(0))
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number core.Number, desc *export.Descriptor) error {
	kind := desc.NumberKind()

	c.current.count.AddUint64Atomic(1)
	c.current.sum.AddNumberAtomic(kind, number)

	for i, boundary := range c.boundaries {
		if number.CompareNumber(kind, boundary) < 0 {
			c.current.buckets.Counts[i].AddUint64Atomic(1)
			return nil
		}
	}

	// Observed event is bigger than all defined boundaries.
	c.current.buckets.Counts[len(c.boundaries)].AddUint64Atomic(1)
	return nil
}

// Merge combines two data sets into one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *export.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)
	c.checkpoint.count.AddNumber(core.Uint64NumberKind, o.checkpoint.count)

	for i := 0; i < len(c.current.buckets.Counts); i++ {
		c.checkpoint.buckets.Counts[i].AddNumber(core.Uint64NumberKind, o.checkpoint.buckets.Counts[i])
	}
	return nil
}

// numbers is an auxiliary struct to order histogram bucket boundaries (slice of core.Number)
type numbers struct {
	numbers []core.Number
	kind    core.NumberKind
}

var _ sort.Interface = (*numbers)(nil)

func (n *numbers) Len() int {
	return len(n.numbers)
}

func (n *numbers) Less(i, j int) bool {
	return -1 == n.numbers[i].CompareNumber(n.kind, n.numbers[j])
}

func (n *numbers) Swap(i, j int) {
	n.numbers[i], n.numbers[j] = n.numbers[j], n.numbers[i]
}
