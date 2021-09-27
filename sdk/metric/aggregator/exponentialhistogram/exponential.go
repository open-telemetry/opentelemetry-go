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

// DefaultMaxSize is the default number of buckets.
//
// 256 is a good choice
// 320 is a historical choice
//
// The OpenHistogram representation of the Prometheus default explicit
// histogram boundaries (spanning 0.005 to 10) yields 320 base-10
// 90-per-decade log-linear buckets.
//
// NrSketch uses this default.
const DefaultMaxSize = 320

// DefaultNormalScale is the default scale used for a number in the
// range [1, 2).  This is chosen to ensure that indices are
// approximately in the range [-2**30, 2**30].
const DefaultNormalScale int32 = 30

type (
	// Aggregator observes events and counts them in
	// exponentially-spaced buckets.  It is configured with a
	// maximum scale factor which determines resolution.  Scale is
	// automatically adjusted to accomodate the range of input
	// data.
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
		backing    interface{} // nil, []uint8, []uint16, []uint32, or []uint64
		indexBase  int32       // value of backing[0] in [indexStart, indexEnd]
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

// New returns `cnt` number of configured histogram aggregators for `desc`.
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
			state:   &state{},
		}
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

// SynchronizedMove implements export.Aggregator.
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

// Count implements aggregation.Sum.
func (a *Aggregator) Count() (uint64, error) {
	return a.state.count, nil
}

// Sum implements aggregation.Sum.
func (a *Aggregator) Sum() (number.Number, error) {
	if a.kind == number.Int64Kind {
		return number.NewInt64Number(int64(a.state.sum)), nil
	}
	return number.NewFloat64Number(a.state.sum), nil
}

// Scale implements aggregation.ExponentialHistogram.
func (a *Aggregator) Scale() int32 {
	return a.state.mapping.scale
}

// ZeroCount implements aggregation.ExponentialHistogram.
func (a *Aggregator) ZeroCount() uint64 {
	return a.state.zeroCount
}

// Positive implements aggregation.ExponentialHistogram.
func (a *Aggregator) Positive() aggregation.ExponentialBuckets {
	return &a.state.positive
}

// Negatiev implements aggregation.ExponentialHistogram.
func (a *Aggregator) Negative() aggregation.ExponentialBuckets {
	return &a.state.negative
}

// Offset implements aggregation.ExponentialBucket.
func (b *buckets) Offset() int32 {
	return b.indexStart
}

// Len implements aggregation.ExponentialBucket.
func (b *buckets) Len() uint32 {
	if b.backing == nil {
		return 0
	}
	return uint32(b.indexEnd - b.indexStart + 1)
}

// At returns the count of the bucket at a position in the logical
// array of counts.
func (b *buckets) At(pos uint32) uint64 {
	bias := uint32(b.indexBase - b.indexStart)

	if pos < bias {
		pos += b.Len()
	}
	pos -= bias

	switch counts := b.backing.(type) {
	case []uint8:
		return uint64(counts[pos])
	case []uint16:
		return uint64(counts[pos])
	case []uint32:
		return uint64(counts[pos])
	case []uint64:
		return counts[pos]
	default:
		panic("At() with size() == 0")
	}
}

// clearState resets a histogram to the empty state without changing
// its scale or backing array.
func (a *Aggregator) clearState() {
	a.state.positive.clearState()
	a.state.negative.clearState()
	a.state.sum = 0
	a.state.count = 0
	a.state.zeroCount = 0
}

// clearState zeros the backing array.
func (b *buckets) clearState() {
	switch counts := b.backing.(type) {
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

// initialize enters the first value into a histogram and sets its
// ideal scale.
func (a *Aggregator) initialize(b *buckets, value float64) {
	firstScale := idealScale(value)

	a.state.mapping = NewLogarithmMapping(firstScale)

	index := a.state.mapping.MapToIndex(value)

	if b.backing == nil {
		b.backing = []uint8{1}
	} else {
		b.incrementBucket(0, 1)
	}
	b.indexStart = int32(index)
	b.indexEnd = int32(index)
	b.indexBase = b.indexStart
}

// idealScale computes the best scale that results in a valid index.
// the default scale is ideal for normalized range [1,2) and a
// non-zero exponent degrades scale in either direction from zero.
func idealScale(value float64) int32 {
	exponent := getExponent(value)

	scale := DefaultNormalScale
	if exponent > 0 {
		scale -= exponent
	} else {
		scale += exponent
	}
	return scale
}

// size() reflects the allocated size of the array, not to be confused
// with Len() which is the range of non-zero values.
func (b *buckets) size() uint32 {
	switch counts := b.backing.(type) {
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

func (a *Aggregator) update(b *buckets, value float64) {
	// Are there any non-zero buckets yet?
	if a.state.count == a.state.zeroCount {
		a.initialize(b, value)
		return
	}

	index := a.state.mapping.MapToIndex(value)

	var span uint32
	if index >= math.MinInt32 && index <= math.MaxInt32 {
		var success bool
		if span, success = a.incrementIndex(b, int32(index)); success {
			return
		}
	}

	// two reasons for this, both call for change of scale:
	// (1) index does not fit a 32-bit value
	// (2) index is outside the maxSize range relative to current extrema.
	down := (span + a.maxSize - 1) / a.maxSize
	shift := int32(31 - bits.LeadingZeros32(down))
	newScale := a.state.mapping.scale - shift

	if ideal := idealScale(value); ideal < newScale {
		newScale = ideal
	}

	change := a.state.mapping.scale - newScale

	a.state.positive.downscale(change)
	a.state.negative.downscale(change)
	a.state.mapping = NewLogarithmMapping(newScale)

	index = a.state.mapping.MapToIndex(value)
	span, success := a.incrementIndex(b, int32(index))

	if !success {
		panic("downscale logic error")
	}
}

// increment determines if the index lies inside the current range
// [indexStart, indexEnd] and if not whether growing the array up to
// maxSize will satisfy the new value.
func (a *Aggregator) incrementIndex(b *buckets, index int32) (uint32, bool) {
	space := b.size()

	if index < b.indexStart {
		if span := uint32(b.indexEnd - index); span >= a.maxSize {
			return span + 1, false // rescale needed
		} else if span >= space {
			a.grow(b, span+1)
		}
		b.indexStart = index
	} else if index > b.indexEnd {
		if span := uint32(index - b.indexStart); span >= a.maxSize {
			return span + 1, false // rescale needed
		} else if span >= space {
			a.grow(b, span+1)
		}
		b.indexEnd = index
	}

	size := int32(b.size())
	bucketIndex := index - b.indexBase
	if bucketIndex >= size {
		bucketIndex -= size
	} else if bucketIndex < 0 {
		bucketIndex += size
	}
	b.incrementBucket(bucketIndex, 1)
	return 0, true
}

// grow resizes the backing array by doubling in size up to maxSize.
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
	switch counts := b.backing.(type) {
	case []uint8:
		tmp := make([]uint8, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.backing = tmp
	case []uint16:
		tmp := make([]uint16, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.backing = tmp
	case []uint32:
		tmp := make([]uint32, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.backing = tmp
	case []uint64:
		tmp := make([]uint64, growTo)
		copy(tmp[part:], counts[diff:])
		copy(tmp[0:diff], counts[0:diff])
		b.backing = tmp
	default:
		panic("grow() with size() == 0")
	}
}

func (b *buckets) downscale(by int32) {
	b.rotate()

	size := 1 + b.indexEnd - b.indexStart
	each := int32(1) << by
	inpos := int32(0)
	outpos := int32(0)
	for pos := b.indexStart; pos <= b.indexEnd; {
		base := pos
		mod := base % each
		if mod < 0 {
			mod += each
		}
		for i := mod; i < each && inpos < size; i++ {
			b.relocateBucket(outpos, inpos)
			inpos++
			pos++
		}
		outpos++
	}

	b.indexStart >>= by
	b.indexEnd >>= by
	b.indexBase = b.indexStart
}

func (b *buckets) rotate() {
	bias := uint32(b.indexBase - b.indexStart)

	if bias == 0 {
		return
	}

	post := b.Len() - bias

	// Rotate the array so that indexBase == indexStart
	b.indexBase = b.indexStart
	switch counts := b.backing.(type) {
	case []uint8:
		for off := uint32(0); off < post; {
			copy(counts[off:off+bias], counts[post:])
			off += bias
		}
	case []uint16:
		for off := uint32(0); off < post; {
			copy(counts[off:off+bias], counts[post:])
			off += bias
		}
	case []uint32:
		for off := uint32(0); off < post; {
			copy(counts[off:off+bias], counts[post:])
			off += bias
		}
	case []uint64:
		for off := uint32(0); off < post; {
			copy(counts[off:off+bias], counts[post:])
			off += bias
		}
	}
}

// relocateBucket adds the count in counts[src] to counts[dest] and
// resets count[src] to zero.
func (b *buckets) relocateBucket(dest, src int32) {
	if dest == src {
		return
	}
	switch counts := b.backing.(type) {
	case []uint8:
		tmp := counts[src]
		counts[src] = 0
		b.incrementBucket(dest, uint64(tmp))
	case []uint16:
		tmp := counts[src]
		counts[src] = 0
		b.incrementBucket(dest, uint64(tmp))
	case []uint32:
		tmp := counts[src]
		counts[src] = 0
		b.incrementBucket(dest, uint64(tmp))
	case []uint64:
		tmp := counts[src]
		counts[src] = 0
		b.incrementBucket(dest, uint64(tmp))
	}
}

// incrementBucket increments the backing array index by `incr`.
func (b *buckets) incrementBucket(bucketIndex int32, incr uint64) {
	for {
		switch counts := b.backing.(type) {
		case []uint8:
			if uint64(counts[bucketIndex])+incr < 0x100 {
				counts[bucketIndex] += uint8(incr)
				return
			}
			tmp := make([]uint16, len(counts))
			for i := range counts {
				tmp[i] = uint16(counts[i])
			}
			b.backing = tmp
			continue
		case []uint16:
			if uint64(counts[bucketIndex])+incr < 0x10000 {
				counts[bucketIndex] += uint16(incr)
				return
			}
			tmp := make([]uint32, len(counts))
			for i := range counts {
				tmp[i] = uint32(counts[i])
			}
			b.backing = tmp
			continue
		case []uint32:
			if uint64(counts[bucketIndex])+incr < 0x100000000 {
				counts[bucketIndex] += uint32(incr)
				return
			}
			tmp := make([]uint64, len(counts))
			for i := range counts {
				tmp[i] = uint64(counts[i])
			}
			b.backing = tmp
			continue
		case []uint64:
			counts[bucketIndex] += incr
			return
		default:
			panic("increment with nil slice")
		}
	}
}
