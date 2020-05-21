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

// Note: This code uses a Mutex to govern access to the exclusive
// aggregator state.  This is in contrast to a lock-free approach
// (as in the Go prometheus client) that was reverted here:
// https://github.com/open-telemetry/opentelemetry-go/pull/669

type (
	// Aggregator observe events and counts them in pre-determined buckets.
	// It also calculates the sum and count of all events.
	Aggregator struct {
		lock       sync.Mutex
		current    state
		checkpoint state
		boundaries []float64
		kind       metric.NumberKind
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state struct {
		bucketCounts []float64
		count        metric.Number
		sum          metric.Number
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
func New(desc *metric.Descriptor, boundaries []float64) *Aggregator {
	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := make([]float64, len(boundaries))

	copy(sortedBoundaries, boundaries)
	sort.Float64s(sortedBoundaries)

	return &Aggregator{
		kind:       desc.NumberKind(),
		boundaries: sortedBoundaries,
		current:    emptyState(sortedBoundaries),
		checkpoint: emptyState(sortedBoundaries),
	}
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
	return aggregator.Buckets{
		Boundaries: c.boundaries,
		Counts:     c.checkpoint.bucketCounts,
	}, nil
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(ctx context.Context, desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, emptyState(c.boundaries)
	c.lock.Unlock()
}

func emptyState(boundaries []float64) state {
	return state{
		bucketCounts: make([]float64, len(boundaries)+1),
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number metric.Number, desc *metric.Descriptor) error {
	kind := desc.NumberKind()
	asFloat := number.CoerceToFloat64(kind)

	bucketID := len(c.boundaries)
	for i, boundary := range c.boundaries {
		if asFloat < boundary {
			bucketID = i
			break
		}
	}
	// Note: Binary-search was compared using the benchmarks. The following
	// code is equivalent to the linear search above:
	//
	//     bucketID := sort.Search(len(c.boundaries), func(i int) bool {
	//         return asFloat < c.boundaries[i]
	//     })
	//
	// The binary search wins for very large boundary sets, but
	// the linear search performs better up through arrays between
	// 256 and 512 elements, which is a relatively large histogram, so we
	// continue to prefer linear search.

	c.lock.Lock()
	defer c.lock.Unlock()

	c.current.count.AddInt64(1)
	c.current.sum.AddNumber(kind, number)
	c.current.bucketCounts[bucketID]++

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

	for i := 0; i < len(c.checkpoint.bucketCounts); i++ {
		c.checkpoint.bucketCounts[i] += o.checkpoint.bucketCounts[i]
	}
	return nil
}
