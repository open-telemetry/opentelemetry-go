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

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
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
		boundaries []float64
		kind       number.Kind
		state      state
	}

	// Config describes how the histogram is aggregated.
	Config struct {
		// ExplicitBoundaries support arbitrary bucketing schemes.  This
		// is the general case.
		ExplicitBoundaries []float64
	}

	// Option configures a histogram Config.
	Option interface {
		// Apply sets one or more Config fields.
		Apply(*Config)
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state struct {
		bucketCounts []uint64
		sum          number.Number
		count        uint64
	}
)

// WithExplicitBoundaries sets the ExplicitBoundaries configuration option of a Config.
func WithExplicitBoundaries(explicitBoundaries []float64) Option {
	return explicitBoundariesOption{explicitBoundaries}
}

type explicitBoundariesOption struct {
	boundaries []float64
}

func (o explicitBoundariesOption) Apply(config *Config) {
	config.ExplicitBoundaries = o.boundaries
}

// defaultExplicitBoundaries have been copied from prometheus.DefBuckets.
//
// Note we anticipate the use of a high-precision histogram sketch as
// the standard histogram aggregator for OTLP export.
// (https://github.com/open-telemetry/opentelemetry-specification/issues/982).
var defaultFloat64ExplicitBoundaries = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

// defaultInt64ExplicitBoundaryMultiplier determines the default
// integer histogram boundaries.
const defaultInt64ExplicitBoundaryMultiplier = 1e6

// defaultInt64ExplicitBoundaries applies a multiplier to the default
// float64 boundaries: [ 5K, 10K, 25K, ..., 2.5M, 5M, 10M ]
var defaultInt64ExplicitBoundaries = func(bounds []float64) (asint []float64) {
	for _, f := range bounds {
		asint = append(asint, defaultInt64ExplicitBoundaryMultiplier*f)
	}
	return
}(defaultFloat64ExplicitBoundaries)

var _ export.Aggregator = &Aggregator{}
var _ aggregation.Sum = &Aggregator{}
var _ aggregation.Count = &Aggregator{}
var _ aggregation.Histogram = &Aggregator{}

// New returns a new aggregator for computing Histograms.
//
// A Histogram observe events and counts them in pre-defined buckets.
// And also provides the total sum and count of all observations.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func New(cnt int, desc *metric.Descriptor, opts ...Option) []Aggregator {
	var cfg Config

	if desc.NumberKind() == number.Int64Kind {
		cfg.ExplicitBoundaries = defaultInt64ExplicitBoundaries
	} else {
		cfg.ExplicitBoundaries = defaultFloat64ExplicitBoundaries
	}

	for _, opt := range opts {
		opt.Apply(&cfg)
	}

	aggs := make([]Aggregator, cnt)

	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := make([]float64, len(cfg.ExplicitBoundaries))

	copy(sortedBoundaries, cfg.ExplicitBoundaries)
	sort.Float64s(sortedBoundaries)

	for i := range aggs {
		aggs[i] = Aggregator{
			kind:       desc.NumberKind(),
			boundaries: sortedBoundaries,
			state:      emptyState(sortedBoundaries),
		}
	}
	return aggs
}

// Aggregation returns an interface for reading the state of this aggregator.
func (c *Aggregator) Aggregation() aggregation.Aggregation {
	return c
}

// Kind returns aggregation.HistogramKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

// Sum returns the sum of all values in the checkpoint.
func (c *Aggregator) Sum() (number.Number, error) {
	return c.state.sum, nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator) Count() (uint64, error) {
	return c.state.count, nil
}

// Histogram returns the count of events in pre-determined buckets.
func (c *Aggregator) Histogram() (aggregation.Buckets, error) {
	return aggregation.Buckets{
		Boundaries: c.boundaries,
		Counts:     c.state.bucketCounts,
	}, nil
}

// SynchronizedMove saves the current state into oa and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator) SynchronizedMove(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)

	if oa != nil && o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.lock.Lock()
	if o != nil {
		o.state = c.state
	}
	c.state = emptyState(c.boundaries)
	c.lock.Unlock()

	return nil
}

func emptyState(boundaries []float64) state {
	return state{
		bucketCounts: make([]uint64, len(boundaries)+1),
	}
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator) Update(_ context.Context, number number.Number, desc *metric.Descriptor) error {
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

	c.state.count++
	c.state.sum.AddNumber(kind, number)
	c.state.bucketCounts[bucketID]++

	return nil
}

// Merge combines two histograms that have the same buckets into a single one.
func (c *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(c, oa)
	}

	c.state.sum.AddNumber(desc.NumberKind(), o.state.sum)
	c.state.count += o.state.count

	for i := 0; i < len(c.state.bucketCounts); i++ {
		c.state.bucketCounts[i] += o.state.bucketCounts[i]
	}
	return nil
}
