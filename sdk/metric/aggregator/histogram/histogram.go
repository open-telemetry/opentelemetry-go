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
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
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
		self         *Aggregator
		bucketCounts []float64
		count        metric.Number
		sum          metric.Number
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregation.Sum = &state{}
var _ aggregation.Count = &state{}
var _ aggregation.Histogram = &state{}

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

	agg := &Aggregator{
		kind:       desc.NumberKind(),
		boundaries: sortedBoundaries,
	}
	agg.current = agg.emptyState(sortedBoundaries)
	agg.checkpoint = agg.emptyState(sortedBoundaries)
	return agg
}

// Kind returns aggregation.HistogramKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

// Checkpoint saves the current state and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator) Checkpoint(desc *metric.Descriptor) {
	c.lock.Lock()
	c.checkpoint, c.current = c.current, c.emptyState(c.boundaries)
	c.lock.Unlock()
}

func (c *Aggregator) emptyState(boundaries []float64) state {
	// TODO: It is possible to avoid allocating new arrays on each
	// checkpoint by re-using the existing slices.
	return state{
		self:         c,
		bucketCounts: make([]float64, len(boundaries)+1),
	}
}

func (c *Aggregator) Swap() {
	c.checkpoint, c.current = c.current, c.checkpoint
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

	c.current.sum.AddNumber(desc.NumberKind(), o.checkpoint.sum)
	c.current.count.AddNumber(metric.Uint64NumberKind, o.checkpoint.count)

	for i := 0; i < len(c.current.bucketCounts); i++ {
		c.current.bucketCounts[i] += o.checkpoint.bucketCounts[i]
	}
	return nil
}

func (c *Aggregator) CheckpointedValue() aggregation.Aggregation {
	return &c.checkpoint
}

func (c *Aggregator) AccumulatedValue() aggregation.Aggregation {
	return &c.current
}

// Kind returns aggregation.HistogramKind.
func (s *state) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

// Sum returns the sum of all values in the checkpoint.
func (s *state) Sum() (metric.Number, error) {
	return s.sum, nil
}

// Count returns the number of values in the checkpoint.
func (s *state) Count() (int64, error) {
	return int64(s.count), nil
}

// Histogram returns the count of events in pre-determined buckets.
func (s *state) Histogram() (aggregation.Buckets, error) {
	return aggregation.Buckets{
		Boundaries: s.self.boundaries,
		Counts:     s.bucketCounts,
	}, nil
}
