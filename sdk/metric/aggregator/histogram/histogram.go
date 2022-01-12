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
	"sort"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

// Note: This code uses a Mutex to govern access to the exclusive
// aggregator state.  This is in contrast to a lock-free approach
// (as in the Go prometheus client) that was reverted here:
// https://github.com/open-telemetry/opentelemetry-go/pull/669

type (
	Defaults interface {
		Boundaries() []float64
	}

	Int64Defaults struct{}
	Float64Defaults struct{}
	
	// Aggregator observe events and counts them in pre-determined buckets.
	// It also calculates the sum and count of all events.
	Aggregator[N number.Any, Traits traits.Any[N]] struct {
		lock       sync.Mutex
		boundaries []float64
		state      *state[N, Traits]
	}

	// config describes how the histogram is aggregated.
	Config struct {
		// explicitBoundaries support arbitrary bucketing schemes.  This
		// is the general case.
		explicitBoundaries []float64
	}

	// Option configures a histogram config.
	Option interface {
		// apply sets one or more config fields.
		apply(*Config)
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state[N number.Any, Traits traits.Any[N]] struct {
		bucketCounts []uint64
		sum          N
		count        uint64
		traits       Traits
	}
)

// WithExplicitBoundaries sets the ExplicitBoundaries configuration option of a config.
func WithExplicitBoundaries(explicitBoundaries []float64) Option {
	return explicitBoundariesOption{explicitBoundaries}
}

type explicitBoundariesOption struct {
	boundaries []float64
}

func (o explicitBoundariesOption) apply(config *Config) {
	config.explicitBoundaries = o.boundaries
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

func (Int64Defaults) Boundaries() []float64 {
	return defaultInt64ExplicitBoundaries
}

func (Float64Defaults) Boundaries() []float64 {
	return defaultFloat64ExplicitBoundaries
}

type Int64Histogram = Aggregator[int64, traits.Int64]
type Float64Histogram = Aggregator[float64, traits.Float64]

var _ aggregator.Aggregator[int64, Int64Histogram, Config] = &Int64Histogram{}
var _ aggregator.Aggregator[float64, Float64Histogram, Config] = &Float64Histogram{}

func NewConfig(def Defaults, opts ...Option) Config {
	cfg := Config{
		explicitBoundaries: def.Boundaries(),
	}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := make([]float64, len(cfg.explicitBoundaries))

	copy(sortedBoundaries, cfg.explicitBoundaries)
	sort.Float64s(sortedBoundaries)
	cfg.explicitBoundaries = sortedBoundaries
	return cfg
}

// New returns a new aggregator for computing Histograms.
//
// A Histogram observe events and counts them in pre-defined buckets.
// And also provides the total sum and count of all observations.
//
// Note that this aggregator maintains each value using independent
// atomic operations, which introduces the possibility that
// checkpoints are inconsistent.
func (a *Aggregator[N, Traits]) Init(cfg Config) {
	a.boundaries = cfg.explicitBoundaries
	a.state = a.newState()
}

// Sum returns the sum of all values in the checkpoint.
func (c *Aggregator[N, Traits]) Sum() (number.Number, error) {
	var traits Traits
	return traits.ToNumber(c.state.sum), nil
}

// Count returns the number of values in the checkpoint.
func (c *Aggregator[N, Traits]) Count() (uint64, error) {
	return c.state.count, nil
}

// Histogram returns the count of events in pre-determined buckets.
func (c *Aggregator[N, Traits]) Histogram() (aggregation.Buckets, error) {
	return aggregation.Buckets{
		Boundaries: c.boundaries,
		Counts:     c.state.bucketCounts,
	}, nil
}

// SynchronizedMove saves the current state into oa and resets the current state to
// the empty set.  Since no locks are taken, there is a chance that
// the independent Sum, Count and Bucket Count are not consistent with each
// other.
func (c *Aggregator[N, Traits]) SynchronizedMove(o *Aggregator[N, Traits]) {
	if o != nil {
		// Swap case: This is the ordinary case for a
		// synchronous instrument, where the SDK allocates two
		// Aggregators and lock contention is anticipated.
		// Reset the target state before swapping it under the
		// lock below.
		o.clearState()
	}

	c.lock.Lock()
	if o != nil {
		c.state, o.state = o.state, c.state
	} else {
		// No swap case: This is the ordinary case for an
		// asynchronous instrument, where the SDK allocates a
		// single Aggregator and there is no anticipated lock
		// contention.
		c.clearState()
	}
	c.lock.Unlock()
}

func (c *Aggregator[N, Traits]) newState() *state[N, Traits] {
	return &state[N, Traits]{
		bucketCounts: make([]uint64, len(c.boundaries)+1),
	}
}

func (c *Aggregator[N, Traits]) clearState() {
	for i := range c.state.bucketCounts {
		c.state.bucketCounts[i] = 0
	}
	c.state.sum = 0
	c.state.count = 0
}

// Update adds the recorded measurement to the current data set.
func (c *Aggregator[N, Traits]) Update(number N) {
	asFloat := float64(number)

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
	c.state.sum += number
	c.state.bucketCounts[bucketID]++
}

// Merge combines two histograms that have the same buckets into a single one.
func (c *Aggregator[N, Traits]) Merge(o *Aggregator[N, Traits]) {
	c.state.sum += o.state.sum
	c.state.count += o.state.count

	for i := 0; i < len(c.state.bucketCounts); i++ {
		c.state.bucketCounts[i] += o.state.bucketCounts[i]
	}
}
