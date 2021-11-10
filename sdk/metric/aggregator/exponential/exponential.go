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
	"math/bits"
	"sync"

	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/logarithm"
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
		maxSize int32
		state   *state
	}

	// config describes how the histogram is aggregated.
	config struct {
		maxSize int32
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
		mapping   mapping.Mapping
	}

	buckets struct {
		backing    interface{} // nil, []uint8, []uint16, []uint32, or []uint64
		indexBase  int32       // value of backing[0] in [indexStart, indexEnd]
		indexStart int32
		indexEnd   int32
	}

	highLow struct {
		low  int64
		high int64
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

type maxSizeOption int32

func (o maxSizeOption) apply(config *config) {
	config.maxSize = int32(o)
}

// New returns `cnt` number of configured histogram aggregators for `desc`.
func New(cnt int, desc *sdkapi.Descriptor, opts ...Option) []Aggregator {
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
			state: &state{
				mapping: newMapping(DefaultNormalScale),
			},
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
	// return aggregation.ExponentialHistogramKind
	return aggregation.HistogramKind
}

// SynchronizedMove implements export.Aggregator.
func (a *Aggregator) SynchronizedMove(oa export.Aggregator, desc *sdkapi.Descriptor) error {
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
func (a *Aggregator) Update(_ context.Context, number number.Number, desc *sdkapi.Descriptor) error {
	value := number.CoerceToFloat64(desc.NumberKind())

	a.lock.Lock()
	defer a.lock.Unlock()

	if value == 0 {
		a.state.zeroCount++
	} else if value > 0 {
		a.update(&a.state.positive, value, 1)
	} else {
		a.update(&a.state.negative, -value, 1)
	}

	a.state.count++
	a.state.sum += value

	return nil
}

func int32min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func int64min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func int64max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
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
func (a *Aggregator) Scale() (int32, error) {
	return a.state.mapping.Scale(), nil
}

// ZeroCount implements aggregation.ExponentialHistogram.
func (a *Aggregator) ZeroCount() (uint64, error) {
	return a.state.zeroCount, nil
}

// Positive implements aggregation.ExponentialHistogram.
func (a *Aggregator) Positive() (aggregation.ExponentialBuckets, error) {
	return &a.state.positive, nil
}

// Negative implements aggregation.ExponentialHistogram.
func (a *Aggregator) Negative() (aggregation.ExponentialBuckets, error) {
	return &a.state.negative, nil
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
func (b *buckets) At(pos0 uint32) uint64 {
	pos := pos0
	bias := uint32(b.indexBase - b.indexStart)

	if pos < bias {
		pos += uint32(b.size())
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

func newMapping(scale int32) mapping.Mapping {
	if scale <= 0 {
		return exponent.NewMapping(scale)
	}
	return logarithm.NewMapping(scale)
}

// initialize enters the first value into a histogram and sets its
// ideal scale.
func (a *Aggregator) initialize(b *buckets, value float64) {
	firstScale := idealScale(value)

	a.state.mapping = newMapping(firstScale)

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
	exponent := mapping.GetExponent(value)

	scale := DefaultNormalScale
	if exponent > 0 {
		scale -= exponent
	} else {
		scale += exponent
	}
	return scale
}

func (a *Aggregator) downscale(change int32) {
	newScale := a.state.mapping.Scale() - change

	a.state.positive.downscale(change)
	a.state.negative.downscale(change)
	a.state.mapping = newMapping(newScale)
}

// changeScale computes the required change of scale.
//
// sizeReq = (high-low+1) is the minimum size needed to fit the new
// index at the current scale, i.e., the distance to the more-distant
// extreme inclusive bucket.  We have that:
//
//   sizeReq >= maxSize
//
// Compute the shift equal to the number of times sizeReq must be
// divided by two before sizeReq < maxSize.
//
// Note: this can be computed in a conservative way w/o use of a loop,
// e.g.,
//
//   shift := 64-bits.LeadingZeros64((high-low+1)/int64(a.maxSize))
//
// however this under-counts by 1 some of the time depending on
// alignment.
func (a *Aggregator) changeScale(hl highLow) int32 {
	var change int32
	for hl.high-hl.low >= int64(a.maxSize) {
		hl.high >>= 1
		hl.low >>= 1
		change++
	}
	return change
}

// size() reflects the allocated size of the array, not to be confused
// with Len() which is the range of non-zero values.
func (b *buckets) size() int32 {
	switch counts := b.backing.(type) {
	case []uint8:
		return int32(cap(counts))
	case []uint16:
		return int32(cap(counts))
	case []uint32:
		return int32(cap(counts))
	case []uint64:
		return int32(cap(counts))
	}
	return 0
}

// update increments the appropriate buckets for a given absolute
// value by the provided increment.
func (a *Aggregator) update(b *buckets, value float64, incr uint64) {
	// if there are zeros, we have not fixed the scale yet.
	if a.state.count == a.state.zeroCount {
		a.initialize(b, value)
		return
	}

	index := a.state.mapping.MapToIndex(value)

	hl, success := a.incrementIndexBy(b, index, incr)
	if success {
		return
	}

	a.downscale(a.changeScale(hl))

	index = a.state.mapping.MapToIndex(value)
	_, success = a.incrementIndexBy(b, index, incr)

	if !success {
		panic("downscale logic error")
	}
}

// increment determines if the index lies inside the current range
// [indexStart, indexEnd] and, if not, returns the minimum size (up to
// maxSize) will satisfy the new value.
func (a *Aggregator) incrementIndexBy(b *buckets, index int64, incr uint64) (highLow, bool) {
	if index < int64(b.indexStart) {
		if span := uint64(int64(b.indexEnd) - index); span >= uint64(a.maxSize) {
			return highLow{
				low:  index,
				high: int64(b.indexEnd),
			}, false // rescale needed
		} else if span >= uint64(b.size()) {
			a.grow(b, uint32(span+1))
		}
		b.indexStart = int32(index)
	} else if index > int64(b.indexEnd) {
		if span := uint64(index - int64(b.indexStart)); span >= uint64(a.maxSize) {
			return highLow{
				low:  int64(b.indexStart),
				high: index,
			}, false // rescale needed
		} else if span >= uint64(b.size()) {
			a.grow(b, uint32(span+1))
		}
		b.indexEnd = int32(index)
	}

	size := int32(b.size())
	bucketIndex := int32(index - int64(b.indexBase))
	if bucketIndex >= size {
		bucketIndex -= int32(b.size())
	} else if bucketIndex < 0 {
		bucketIndex += int32(b.size())
	}
	b.incrementBucket(bucketIndex, incr)
	return highLow{}, true
}

// grow resizes the backing array by doubling in size up to maxSize.
// this extends the array with a bunch of zeros and copies the
// existing counts to the same position.
func (a *Aggregator) grow(b *buckets, needed uint32) {
	size := b.size()
	bias := b.indexBase - b.indexStart
	diff := size - bias
	growTo := int32(1) << (32 - bits.LeadingZeros32(uint32(needed)))
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

// downscale first rotates, then collapses 2**`by`-to-1 buckets
func (b *buckets) downscale(by int32) {
	b.rotate()

	size := 1 + b.indexEnd - b.indexStart
	each := int64(1) << by
	inpos := int32(0)
	outpos := int32(0)
	for pos := b.indexStart; pos <= b.indexEnd; {
		mod := int64(pos) % each
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

// rotate shifts the backing array contents so that indexStart ==
// indexBase to simplify the downscale logic.
func (b *buckets) rotate() {
	bias := uint32(b.indexBase - b.indexStart)

	if bias == 0 {
		return
	}

	move := b.Len()
	size := uint32(b.size())
	endpoint := size - bias

	// Rotate the array so that indexBase == indexStart
	b.indexBase = b.indexStart

	switch counts := b.backing.(type) {
	case []uint8:
		for start := uint32(0); start < move; {
			for pos := endpoint; pos < size && start < move; pos++ {
				counts[start], counts[pos] = counts[pos], counts[start]
				start++
			}
		}
	case []uint16:
		for start := uint32(0); start < move; {
			for pos := endpoint; pos < size && start < move; pos++ {
				counts[start], counts[pos] = counts[pos], counts[start]
				start++
			}
		}
	case []uint32:
		for start := uint32(0); start < move; {
			for pos := endpoint; pos < size && start < move; pos++ {
				counts[start], counts[pos] = counts[pos], counts[start]
				start++
			}
		}
	case []uint64:
		for start := uint32(0); start < move; {
			for pos := endpoint; pos < size && start < move; pos++ {
				counts[start], counts[pos] = counts[pos], counts[start]
				start++
			}
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

// Merge combines two histograms that have the same buckets into a single one.
func (a *Aggregator) Merge(oa export.Aggregator, desc *sdkapi.Descriptor) error {
	o, _ := oa.(*Aggregator)
	if o == nil {
		return aggregator.NewInconsistentAggregatorError(a, oa)
	}

	a.state.sum += o.state.sum
	a.state.count += o.state.count
	a.state.zeroCount += o.state.zeroCount

	aScale, _ := a.Scale()
	oScale, _ := o.Scale()
	minScale := int32min(aScale, oScale)

	hlp := a.highLowAtScale(&a.state.positive, minScale)
	hlp = hlp.with(o.highLowAtScale(&o.state.positive, minScale))

	hln := a.highLowAtScale(&a.state.negative, minScale)
	hln = hln.with(o.highLowAtScale(&o.state.negative, minScale))

	minScale = int32min(
		minScale-a.changeScale(hlp),
		minScale-a.changeScale(hln),
	)

	aScale, _ = a.Scale()
	a.downscale(aScale - minScale)

	a.mergeBuckets(&a.state.positive, o, &o.state.positive, minScale)
	a.mergeBuckets(&a.state.negative, o, &o.state.negative, minScale)

	return nil
}

// mergeBuckets translates index values from another histogram into
// the corresponding buckets of this histogram.
func (a *Aggregator) mergeBuckets(mine *buckets, other *Aggregator, theirs *buckets, scale int32) {
	otherScale, _ := other.Scale()
	theirOffset := theirs.Offset()
	theirChange := otherScale - scale

	for i := uint32(0); i < theirs.Len(); i++ {
		_, success := a.incrementIndexBy(mine,
			int64(theirOffset+int32(i))>>theirChange, theirs.At(i))
		if !success {
			panic("incorrect merge scale")
		}
	}
}

func (a *Aggregator) highLowAtScale(b *buckets, scale int32) highLow {
	if b.Len() == 0 {
		return highLow{
			low:  0,
			high: -1,
		}
	}
	aScale, _ := a.Scale()
	shift := aScale - scale
	return highLow{
		low:  int64(b.indexStart >> shift),
		high: int64(b.indexEnd >> shift),
	}
}

func (h *highLow) with(o highLow) highLow {
	if o.empty() {
		return *h
	}
	if h.empty() {
		return o
	}
	return highLow{
		low:  int64min(h.low, o.low),
		high: int64max(h.high, o.high),
	}
}

func (h *highLow) empty() bool {
	return h.low >= h.high
}
