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
	"fmt"
	"sort"
	"sync"

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
)

var ErrNoSubtract = fmt.Errorf("histogram subtract not implemented")

// Note: This code uses a Mutex to govern access to the exclusive
// aggregator state.  This is in contrast to a lock-free approach
// (as in the Go prometheus client) that was reverted here:
// https://github.com/open-telemetry/opentelemetry-go/pull/669

type (
	State[N number.Any, Traits traits.Any[N]] struct {
		boundaries []float64

		lock         sync.Mutex
		bucketCounts []uint64
		sum          N
		count        uint64
	}

	Methods[N number.Any, Traits traits.Any[N], Storage State[N, Traits]] struct{}
)

var (
	_ aggregator.Methods[int64, State[int64, traits.Int64]]       = Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	_ aggregator.Methods[float64, State[float64, traits.Float64]] = Methods[float64, traits.Float64, State[float64, traits.Float64]]{}

	_ aggregation.Histogram = &State[int64, traits.Int64]{}
	_ aggregation.Histogram = &State[float64, traits.Float64]{}
)

func NewFloat64(boundries []float64, values ...float64) aggregation.Histogram {
	if len(boundries) < 1 {
		boundries = DefaultFloat64Boundaries
	}

	hist := &State[float64, traits.Float64]{
		boundaries:   boundries,
		bucketCounts: make([]uint64, len(boundries)+1),
	}
	methods := Methods[float64, traits.Float64, State[float64, traits.Float64]]{}
	for _, val := range values {
		methods.Update(hist, val)
	}
	return hist
}

func NewInt64(boundries []float64, values ...int64) aggregation.Histogram {
	if len(boundries) < 1 {
		boundries = DefaultFloat64Boundaries
	}

	hist := &State[int64, traits.Int64]{
		boundaries:   boundries,
		bucketCounts: make([]uint64, len(boundries)+1),
	}
	methods := Methods[int64, traits.Int64, State[int64, traits.Int64]]{}
	for _, val := range values {
		methods.Update(hist, val)
	}
	return hist
}

// DefaultBoundaries have been copied from prometheus.DefBuckets.
var DefaultFloat64Boundaries = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

// defaultInt64BoundaryMultiplier determines the default
// integer histogram boundaries.
const defaultInt64BoundaryMultiplier = 1e6

// DefaultInt64Boundaries applies a multiplier to the default
// float64 boundaries: [ 5K, 10K, 25K, ..., 2.5M, 5M, 10M ]
var DefaultInt64Boundaries = func(bounds []float64) (asint []float64) {
	for _, f := range bounds {
		asint = append(asint, defaultInt64BoundaryMultiplier*f)
	}
	return
}(DefaultFloat64Boundaries)

func NewConfig(bounds []float64) aggregator.HistogramConfig {
	// Boundaries MUST be ordered otherwise the histogram could not
	// be properly computed.
	sortedBoundaries := make([]float64, len(bounds))

	copy(sortedBoundaries, bounds)
	sort.Float64s(sortedBoundaries)
	return aggregator.HistogramConfig{ExplicitBoundaries: sortedBoundaries}
}

func (h *State[N, Traits]) Sum() number.Number {
	var traits Traits
	return traits.ToNumber(h.sum)
}

func (h *State[N, Traits]) Count() uint64 {
	return h.count
}

func (h *State[N, Traits]) Histogram() aggregation.Buckets {
	return aggregation.Buckets{
		Boundaries: h.boundaries,
		Counts:     h.bucketCounts,
	}
}

func (h *State[N, Traits]) Category() aggregation.Category {
	return aggregation.HistogramCategory
}

func (h *State[N, Traits]) clearState() {
	for i := range h.bucketCounts {
		h.bucketCounts[i] = 0
	}
	h.sum = 0
	h.count = 0
}

func (Methods[N, Traits, Storage]) Kind() aggregation.Kind {
	return aggregation.HistogramKind
}

func (Methods[N, Traits, Storage]) Init(state *State[N, Traits], cfg aggregator.Config) {
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

func (Methods[N, Traits, Storage]) SubtractSwap(value, operandToModify *State[N, Traits]) {
	operandToModify.sum = value.sum - operandToModify.sum
	operandToModify.count = value.count - operandToModify.count

	for i := range value.bucketCounts {
		operandToModify.bucketCounts[i] = value.bucketCounts[i] - operandToModify.bucketCounts[i]
	}
}
