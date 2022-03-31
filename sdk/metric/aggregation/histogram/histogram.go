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

package histogram // import "go.opentelemetry.io/otel/sdk/metric/aggregation/histogram"

import (
	"fmt"
	"sort"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

var ErrNoSubtract = fmt.Errorf("histogram subtract not implemented")

// Note: This code uses a Mutex to govern access to the exclusive
// aggregator state.  This is in contrast to a lock-free approach
// (as in the Go prometheus client) that was reverted here:
// https://github.com/open-telemetry/opentelemetry-go/pull/669

type (
	Defaults interface {
		Boundaries() []float64
	}

	Int64Defaults   struct{}
	Float64Defaults struct{}

	State[N number.Any, Traits traits.Any[N]] struct {
		boundaries []float64

		lock         sync.Mutex
		bucketCounts []uint64
		sum          N
		count        uint64
	}

	Option interface {
		// apply sets one or more config fields.
		apply(*aggregation.HistogramConfig)
	}

	Methods[N number.Any, Traits traits.Any[N], Storage State[N, Traits]] struct{}
)

var (
	_ aggregation.Methods[int64, State[int64, traits.Int64]]       = Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	_ aggregation.Methods[float64, State[float64, traits.Float64]] = Methods[float64, traits.Float64, State[float64, traits.Float64]]{}

	_ aggregation.Histogram = &State[int64, traits.Int64]{}
	_ aggregation.Histogram = &State[float64, traits.Float64]{}
)

// WithExplicitBoundaries sets the ExplicitBoundaries configuration option of a config.
func WithExplicitBoundaries(explicitBoundaries []float64) Option {
	return explicitBoundariesOption{explicitBoundaries}
}

type explicitBoundariesOption struct {
	boundaries []float64
}

func (o explicitBoundariesOption) apply(config *aggregation.HistogramConfig) {
	config.ExplicitBoundaries = o.boundaries
}

// defaultExplicitBoundaries have been copied from prometheus.DefBuckets.
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

func NewConfig(def Defaults, opts ...Option) aggregation.HistogramConfig {
	cfg := aggregation.HistogramConfig{
		ExplicitBoundaries: def.Boundaries(),
	}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := make([]float64, len(cfg.ExplicitBoundaries))

	copy(sortedBoundaries, cfg.ExplicitBoundaries)
	sort.Float64s(sortedBoundaries)
	cfg.ExplicitBoundaries = sortedBoundaries
	return cfg
}

func (h *State[N, Traits]) Sum() number.Number {
	h.lock.Lock()
	defer h.lock.Unlock()
	var traits Traits
	return traits.ToNumber(h.sum)
}

func (h *State[N, Traits]) Count() uint64 {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.count
}

func (h *State[N, Traits]) Histogram() aggregation.Buckets {
	h.lock.Lock()
	defer h.lock.Unlock()
	return aggregation.Buckets{
		Boundaries: h.boundaries,
		Counts:     h.bucketCounts,
	}
}

func (h *State[N, Traits]) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

func (h *State[N, Traits]) clearState() {
	for i := range h.bucketCounts {
		h.bucketCounts[i] = 0
	}
	h.sum = 0
	h.count = 0
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], cfg aggregation.Config) {
	state.boundaries = cfg.Histogram.ExplicitBoundaries
	state.bucketCounts = make([]uint64, len(state.boundaries)+1)
}

func (Methods[N, Traits, Storage]) Reset(ptr *State[N, Traits]) {
	ptr.clearState()
}

func (Methods[N, Traits, Storage]) HasChange(ptr *State[N, Traits]) bool {
	return ptr.count == 0
}

func (Methods[N, Traits, Storage]) SynchronizedMove(resetSrc, dest *State[N, Traits]) {
	// Swap case: This is the ordinary case for a
	// synchronous instrument, where the SDK allocates two
	// Aggregators and lock contention is anticipated.
	// Reset the target state before swapping it under the
	// lock below.
	dest.clearState()

	resetSrc.lock.Lock()
	defer resetSrc.lock.Unlock()

	dest.sum, resetSrc.sum = resetSrc.sum, dest.sum
	dest.count, resetSrc.count = resetSrc.count, dest.count
	dest.bucketCounts, resetSrc.bucketCounts = resetSrc.bucketCounts, dest.bucketCounts
}

// Update adds the recorded measurement to the current data set.
func (Methods[N, Traits, Storage]) Update(state *State[N, Traits], number N) {
	asFloat := float64(number)

	bucketID := len(state.boundaries)
	for i, boundary := range state.boundaries {
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

	state.lock.Lock()
	defer state.lock.Unlock()

	state.count++
	state.sum += number
	state.bucketCounts[bucketID]++
}

// Merge combines two histograms that have the same buckets into a single one.
func (Methods[N, Traits, Storage]) Merge(to, from *State[N, Traits]) {
	to.sum += from.sum
	to.count += from.count

	for i := 0; i < len(to.bucketCounts); i++ {
		to.bucketCounts[i] += from.bucketCounts[i]
	}
}

func (Methods[N, Traits, Storage]) Aggregation(state *State[N, Traits]) aggregation.Aggregation {
	return state
}

func (Methods[N, Traits, Storage]) Storage(aggr aggregation.Aggregation) *State[N, Traits] {
	return aggr.(*State[N, Traits])
}

func (Methods[N, Traits, Storage]) SubtractSwap(valueToModify, operand *State[N, Traits]) {
	panic("@@@HERE")
}
