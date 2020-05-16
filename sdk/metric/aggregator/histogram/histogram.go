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
	"sync"

	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type (
	// Aggregator observe events and counts them in pre-determined buckets.
	// It also calculates the sum and count of all events.
	Aggregator struct {
		lock       sync.Mutex
		current    state
		checkpoint state
		boundaries []metric.Number
		kind       metric.NumberKind
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state struct {
		// all fields have to be aligned for 64-bit atomic operations.
		buckets aggregator.Buckets
		count   metric.Number
		sum     metric.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregator.Sum = &Aggregator{}
var _ aggregator.Count = &Aggregator{}
var _ aggregator.Histogram = &Aggregator{}

// New returns a new aggregator for computing Histograms.
//
// A Histogram observe events and counts them in pre-defined buckets.
// And also provides the total sum and count of all observations.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func New(desc *metric.Descriptor, boundaries []metric.Number) *Aggregator {
	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := numbers{
		numbers: make([]metric.Number, len(boundaries)),
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
				Counts:     make([]metric.Number, len(boundaries)+1),
			},
		},
	}
	return &agg
}

// Sum returns the sum of all values in the checkpoint.
func (c *Aggregator) Sum() (metric.Number, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return int64(c.checkpoint.count), nil
}

// Histogram returns the count of events in pre-determined buckets.
func (c *Aggregator) Histogram() (aggregator.Buckets, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.checkpoint.buckets, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, c.emptyState()
	c.lock.Unlock()
}

func (c *Aggregator) emptyState() state {
	return state{
		buckets: aggregator.Buckets{
			Boundaries: c.boundaries,
			Counts:     make([]metric.Number, len(c.boundaries)+1),
		},
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	kind := desc.NumberKind()

	bucketID := len(c.boundaries)
	for i, boundary := range c.boundaries {
		if number.CompareNumber(kind, boundary) < 0 {
			bucketID = i
			break
		}
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.current.count.AddInt64(1)
	c.current.sum.AddNumber(kind, number)
	c.current.buckets.Counts[bucketID].AddUint64(1)

	return nil
}

// Merge combines two histograms that have the same buckets into a single one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentMergeError(c, oa)
	}

	c.checkpoint.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)
	c.checkpoint.count.AddNumber(metric.Uint64NumberKind, o.checkpoint.count)

	for i := 0; i < len(c.checkpoint.buckets.Counts); i++ {
		c.checkpoint.buckets.Counts[i].AddNumber(metric.Uint64NumberKind, o.checkpoint.buckets.Counts[i])
	}
	return nil
}

// numbers is an auxiliary struct to order histogram bucket boundaries (slice of kv.Number)
type numbers struct {
	numbers []metric.Number
	kind    metric.NumberKind
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
