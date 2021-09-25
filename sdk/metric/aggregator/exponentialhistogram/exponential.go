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

package exponential // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential"

import (
	"context"
	"math"
	"math/bits"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
)

// Note: This code uses a Mutex to govern access to the exclusive
// aggregator state.  For an example of a lock-free approach
// see https://github.com/open-telemetry/opentelemetry-go/pull/669.

const DefaultMaxSize = 320

// DefaultNormalScale is the default scale used for a number in the
// range [1, 2).
const DefaultNormalScale int32 = 30

type (
	// Aggregator observes events and counts them in
	// exponentially-spaced buckets.
	Aggregator struct {
		lock    sync.Mutex
		kind    number.Kind
		maxSize uint32
		state   *state
	}

	// config describes how the histogram is aggregated.
	config struct {
		maxSize uint32
	}

	// Option configures a histogram config.
	Option interface {
		// apply sets one or more config fields.
		apply(*config)
	}

	// state represents the state of a histogram, consisting of
	// the sum and counts for all observed values and
	// the less than equal bucket count for the pre-determined boundaries.
	state struct {
		sum       float64
		count     uint64
		zeroCount uint64
		positive  buckets
		negative  buckets
		mapping   LogarithmMapping
	}

	buckets struct {
		wrapped    interface{} // nil, []uint8, []uint16, []uint32, or []uint64
		indexBase  int32       // value of wrapped[0] in [indexStart, indexEnd]
		indexStart int32
		indexEnd   int32
	}
)

var _ export.Aggregator = &Aggregator{}
var _ aggregation.Sum = &Aggregator{}
var _ aggregation.Count = &Aggregator{}
var _ aggregation.ExponentialHistogram = &Aggregator{}

// WithMaxSize sets the maximimum number of buckets.
func WithMaxSize(size int32) Option {
	return maxSizeOption(size)
}

type maxSizeOption int

func (o maxSizeOption) apply(config *config) {
	config.maxSize = uint32(o)
}

func New(cnt int, desc *metric.Descriptor, opts ...Option) []Aggregator {
	cfg := config{
		maxSize: DefaultMaxSize,
	}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	aggs := make([]Aggregator, cnt)

	for i := range aggs {
		aggs[i] = Aggregator{
			kind:    desc.NumberKind(),
			maxSize: cfg.maxSize,
		}
		aggs[i].state = aggs[i].newState()
	}
	return aggs
}

// Aggregation returns an interface for reading the state of this aggregator.
func (a *Aggregator) Aggregation() aggregation.Aggregation {
	return a
}

// Kind returns aggregation.ExponentialHistogramKind.
func (c *Aggregator) Kind() aggregation.Kind {
	return aggregation.ExponentialHistogramKind
}

func (a *Aggregator) SynchronizedMove(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)

	if oa != nil && o == nil {
		return aggregator.NewInconsistentAggregatorError(a, oa)
	}

	if o != nil {
		// Swap case: This is the ordinary case for a
		// synchronous instrument, where the SDK allocates two
		// Aggregators and lock contention is anticipated.
		// Reset the target state before swapping it under the
		// lock below.
		o.clearState()
	}

	a.lock.Lock()
	if o != nil {
		a.state, o.state = o.state, a.state
	} else {
		// No swap case: This is the ordinary case for an
		// asynchronous instrument, where the SDK allocates a
		// single Aggregator and there is no anticipated lock
		// contention.
		a.clearState()
	}
	a.lock.Unlock()

	return nil
}

func (a *Aggregator) newState() *state {
	return &state{}
}

func (a *Aggregator) clearState() {
	a.state.positive.clearState()
	a.state.negative.clearState()
	a.state.sum = 0
	a.state.count = 0
	a.state.zeroCount = 0
}

func (b *buckets) clearState() {
	switch counts := b.wrapped.(type) {
	case []uint8:
		for i := range counts {
			counts[i] = 0
		}
	case []uint16:
		for i := range counts {
			counts[i] = 0
		}
	case []uint32:
		for i := range counts {
			counts[i] = 0
		}
	case []uint64:
		for i := range counts {
			counts[i] = 0
		}
	}
}

// Update adds the recorded measurement to the current data set.
func (a *Aggregator) Update(_ context.Context, number number.Number, desc *metric.Descriptor) error {
	value := number.CoerceToFloat64(desc.NumberKind())

	a.lock.Lock()
	defer a.lock.Unlock()

	if value == 0 {
		a.state.zeroCount++
	} else if value > 0 {
		a.update(&a.state.positive, value)
	} else {
		a.update(&a.state.negative, -value)
	}

	a.state.count++
	a.state.sum += value

	return nil
}

func (a *Aggregator) update(b *buckets, value float64) {
	// Are there any non-zero buckets yet?
	if a.state.count == a.state.zeroCount {
		a.initialize(b, value)
		return
	}

	index := a.state.mapping.MapToIndex(value)

	if index >= math.MinInt32 && index <= math.MaxInt32 {
		if a.increment(b, int32(index)) {
			return
		}
	}

	panic("NOT YET")
}

// initialize enters the first value into a histogram and sets its
// initial scale.
func (a *Aggregator) initialize(b *buckets, value float64) {
	exponent := getExponent(value)
	firstScale := DefaultNormalScale - exponent

	a.state.mapping = NewLogarithmMapping(firstScale)

	index := a.state.mapping.MapToIndex(value)

	b.wrapped = []uint8{1}
	b.indexStart = int32(index)
	b.indexEnd = int32(index)
	b.indexBase = b.indexStart
}

// Merge combines two histograms that have the same buckets into a single one.
func (a *Aggregator) Merge(oa export.Aggregator, desc *metric.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(a, oa)
	}

	// a.state.sum.AddNumber(desc.NumberKind(), o.state.sum)
	// a.state.count += o.state.count

	// for i := 0; i < len(a.state.bucketCounts); i++ {
	// 	a.state.bucketCounts[i] += o.state.bucketCounts[i]
	// }
	return nil
}

func (a *Aggregator) Count() (uint64, error) {
	return a.state.count, nil
}

func (a *Aggregator) Sum() (number.Number, error) {
	if a.kind == number.Int64Kind {
		return number.NewInt64Number(int64(a.state.sum)), nil
	}
	return number.NewFloat64Number(a.state.sum), nil
}

func (a *Aggregator) Scale() int32 {
	return a.state.mapping.scale
}

func (a *Aggregator) ZeroCount() uint64 {
	return a.state.zeroCount
}

func (a *Aggregator) Positive() aggregation.ExponentialBuckets {
	return &a.state.positive
}

func (a *Aggregator) Negative() aggregation.ExponentialBuckets {
	return &a.state.negative
}

func (b *buckets) Offset() int32 {
	return b.indexStart
}

func (b *buckets) Len() uint32 {
	if b.wrapped == nil {
		return 0
	}
	return uint32(b.indexEnd - b.indexStart + 1)
}

func (b *buckets) size() uint32 {
	switch counts := b.wrapped.(type) {
	case []uint8:
		return uint32(cap(counts))
	case []uint16:
		return uint32(cap(counts))
	case []uint32:
		return uint32(cap(counts))
	case []uint64:
		return uint32(cap(counts))
	}
	return 0
}

// At returns the count of the bucket at a position in the logical
// array of counts.
func (b *buckets) At(pos uint32) uint64 {
	area := b.Len()
	bias := uint32(b.indexBase - b.indexStart)
	diff := area - bias

	if pos < bias {
		pos += diff
	} else {
		pos -= bias
	}

	switch counts := b.wrapped.(type) {
	case []uint8:
		return uint64(counts[pos])
	case []uint16:
		return uint64(counts[pos])
	case []uint32:
		return uint64(counts[pos])
	case []uint64:
		return counts[pos]
	default:
		panic("impossible case")
	}
}

// grow resizes the wrapped array by doubling in size up to maxSize.
// this extends the array with a bunch of zeros and copies the
// existing counts to the same position.
func (a *Aggregator) grow(b *buckets, needed uint32) {
	size := b.size()
	bias := uint32(b.indexBase - b.indexStart)
	diff := size - bias
	growTo := uint32(1) << (32 - bits.LeadingZeros32(uint32(needed)))
	if growTo > a.maxSize {
		growTo = a.maxSize
	}
	part := growTo - bias
	switch counts := b.wrapped.(type) {
	case []uint8:
		tmp := make([]uint8, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.wrapped = tmp
	case []uint16:
		tmp := make([]uint16, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.wrapped = tmp
	case []uint32:
		tmp := make([]uint32, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.wrapped = tmp
	case []uint64:
		tmp := make([]uint64, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.wrapped = tmp
	default:
		panic("impossible case")
	}
}

// increment determines if the index lies inside the current range
// [indexStart, indexEnd] and if not whether growing the array up to
// maxSize will satisfy the new value.
func (a *Aggregator) increment(b *buckets, index int32) bool {
	space := b.size()

	if index < b.indexStart {
		if span := uint32(b.indexEnd - index); span >= a.maxSize {
			return false // rescale needed
		} else if span >= space {
			a.grow(b, span+1)
		}
		b.indexStart = index
	} else if index > b.indexEnd {
		if span := uint32(index - b.indexStart); span >= a.maxSize {
			return false // rescale needed
		} else if span >= space {
			a.grow(b, span+1)
		}
		b.indexEnd = index
	}

	l := b.size()
	i := index - b.indexBase
	if i >= int32(l) {
		i -= int32(l)
	} else if i < 0 {
		i += int32(l)
	}

	for {
		switch counts := b.wrapped.(type) {
		case []uint8:
			if counts[i] < 0xff {
				counts[i]++
				return true
			}
			tmp := make([]uint16, len(counts))
			for i := range counts {
				tmp[i] = uint16(counts[i])
			}
			b.wrapped = tmp
			continue
		case []uint16:
			if counts[i] < 0xffff {
				counts[i]++
				return true
			}
			tmp := make([]uint32, len(counts))
			for i := range counts {
				tmp[i] = uint32(counts[i])
			}
			b.wrapped = tmp
			continue
		case []uint32:
			if counts[i] < 0xffffffff {
				counts[i]++
				return true
			}
			tmp := make([]uint64, len(counts))
			for i := range counts {
				tmp[i] = uint64(counts[i])
			}
			b.wrapped = tmp
			continue
		case []uint64:
			counts[i]++
			return true
		default:
			panic("impossible case")
		}
	}
}
