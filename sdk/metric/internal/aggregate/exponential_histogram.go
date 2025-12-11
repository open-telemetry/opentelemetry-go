// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	expoMaxScale = 20
	expoMinScale = -10
)

// expoHistogramDataPoint is a single data point in an exponential histogram.
type expoHistogramDataPoint[N int64 | float64] struct {
	rescaleMux sync.Mutex
	expoHistogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]
}

func newExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
) *expoHistogramDataPoint[N] { // nolint:revive // we need this control flag
	return &expoHistogramDataPoint[N]{
		attrs:                      attrs,
		expoHistogramPointCounters: newExpoHistogramPointCounters[N](maxSize, maxScale),
	}
}

// hotColdExpoHistogramPoint a hot and cold exponential histogram points, used
// in cumulative aggregations.
type hotColdExpoHistogramPoint[N int64 | float64] struct {
	rescaleMux   sync.Mutex
	hcwg         hotColdWaitGroup
	hotColdPoint [2]expoHistogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	maxScale int32
}

func newHotColdExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
) *hotColdExpoHistogramPoint[N] { // nolint:revive // we need this control flag
	return &hotColdExpoHistogramPoint[N]{
		attrs:    attrs,
		maxScale: maxScale,
		hotColdPoint: [2]expoHistogramPointCounters[N]{
			newExpoHistogramPointCounters[N](maxSize, maxScale),
			newExpoHistogramPointCounters[N](maxSize, maxScale),
		},
	}
}

func (p *expoHistogramPointCounters[N]) tryFastRecord(v N, noMinMax, noSum bool) bool { // nolint:revive // we need this control flag
	absV := math.Abs(float64(v))
	if float64(absV) == 0.0 {
		p.zeroCount.Add(1)
		return true
	}
	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}
	if !bucket.tryFastRecord(absV) {
		return false
	}
	if !noMinMax {
		p.minMax.Update(v)
	}
	if !noSum {
		p.sum.add(v)
	}
	return true
}

// record adds a new measurement to the histogram. It will rescale the buckets if needed.
// The caller must hold the rescaleMux lock
func (p *expoHistogramPointCounters[N]) record(v N, noMinMax, noSum bool) { // nolint:revive // we need this control flag
	absV := math.Abs(float64(v))
	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}
	if !bucket.record(absV) {
		// We failed to record for an unrecoverable reason.
		return
	}
	if !noMinMax {
		p.minMax.Update(v)
	}
	if !noSum {
		p.sum.add(v)
	}
}

// expoHistogramPointCounters contains only the atomic counter data, and is
// used by both expoHistogramDataPoint and hotColdExpoHistogramPoint.
type expoHistogramPointCounters[N int64 | float64] struct {
	minMax    atomicMinMax[N]
	sum       atomicCounter[N]
	zeroCount atomic.Uint64

	posBuckets hotColdExpoBuckets
	negBuckets hotColdExpoBuckets
}

func newExpoHistogramPointCounters[N int64 | float64](
	maxSize int,
	maxScale int32) expoHistogramPointCounters[N] {
	return expoHistogramPointCounters[N]{
		posBuckets: newHotColdExpoBuckets(maxSize, maxScale),
		negBuckets: newHotColdExpoBuckets(maxSize, maxScale),
	}
}

// loadInto writes the values of the counters into the datapoint.
// It is safe to call concurrently, but callers need to use a hot/cold
// waitgroup to ensure consistent results.
func (e *expoHistogramPointCounters[N]) loadInto(into *metricdata.ExponentialHistogramDataPoint[N], noMinMax, noSum bool) {
	into.ZeroCount = e.zeroCount.Load()
	if !noSum {
		into.Sum = e.sum.load()
	}
	if !noMinMax && e.minMax.set.Load() {
		into.Min = metricdata.NewExtrema(e.minMax.minimum.Load())
		into.Max = metricdata.NewExtrema(e.minMax.maximum.Load())
	}
	into.Scale = e.posBuckets.unifyScale(&e.negBuckets)

	posCount, posOffset := e.posBuckets.loadCountsAndOffset(&into.PositiveBucket.Counts)
	into.PositiveBucket.Offset = posOffset

	negCount, negOffset := e.negBuckets.loadCountsAndOffset(&into.NegativeBucket.Counts)
	into.NegativeBucket.Offset = negOffset

	into.Count = posCount + negCount + into.ZeroCount

}

// mergeInto merges this set of histogram counter data into another,
// and resets the state of this set of counters. This is used by
// hotColdHistogramPoint to ensure that the cumulative counters continue to
// accumulate after being read.
func (p *expoHistogramPointCounters[N]) mergeIntoAndReset( // nolint:revive // Intentional internal control flag
	into *expoHistogramPointCounters[N],
	noMinMax, noSum bool,
) {
	// Swap in 0 to reset the zero count.
	into.zeroCount.Add(p.zeroCount.Swap(0))
	// Do not reset min or max because cumulative min and max only ever grow
	// smaller or larger respectively.
	if !noMinMax && p.minMax.set.Load() {
		into.minMax.Update(p.minMax.minimum.Load())
		into.minMax.Update(p.minMax.maximum.Load())
	}
	if !noSum {
		into.sum.add(p.sum.load())
		p.sum.reset()
	}
	p.posBuckets.mergeIntoAndReset(&into.posBuckets)
	p.negBuckets.mergeIntoAndReset(&into.negBuckets)
}

type hotColdExpoBuckets struct {
	hcwg           hotColdWaitGroup
	hotColdBuckets [2]expoBuckets

	maxScale int32
}

func newHotColdExpoBuckets(maxSize int, maxScale int32) hotColdExpoBuckets {
	return hotColdExpoBuckets{
		hotColdBuckets: [2]expoBuckets{
			newExpoBuckets(maxSize, maxScale),
			newExpoBuckets(maxSize, maxScale),
		},
		maxScale: maxScale,
	}
}

// tryFastRecord is the fast-path for exponential histogram measurements. It
// succeeds if the value can be written without downscaling the buckets.
// If it fails, it returns false, and also expands the range of the buckets as
// far towards the required bin as possible to prevent the range from changing
// while we downscale.
func (b *hotColdExpoBuckets) tryFastRecord(v float64) bool {
	hotIdx := b.hcwg.start()
	defer b.hcwg.done(hotIdx)
	return b.hotColdBuckets[hotIdx].recordBucket(b.hotColdBuckets[hotIdx].getBin(v))
}

// record is the slow path, and is invoked when the tryFastRecord fails.
// It locks to prevent concurrent scale changes. It downscales buckets to fit
// the measurement, and then records it.
func (b *hotColdExpoBuckets) record(v float64) bool {
	// Hot may have been swapped while we were waiting for the lock.
	// We don't use p.hcwg.start() because we already hold the lock, and would
	// deadlock when waiting for writes to complete.
	hotIdx := b.hcwg.loadHot()
	hotBucket := &b.hotColdBuckets[hotIdx]

	// Try recording again in-case it was resized while we were waiting, and to
	// ensure the bucket range doesn't change.
	bin := hotBucket.getBin(v)
	if hotBucket.recordBucket(hotBucket.getBin(v)) {
		return true
	}

	hotBucket.startEndMux.Lock()
	defer hotBucket.startEndMux.Unlock()

	// Since recordBucket failed above, we know we need a scale change.
	scaleDelta := hotBucket.scaleChange(bin)
	if hotBucket.scale-scaleDelta < expoMinScale {
		// With a scale of -10 there is only two buckets for the whole range of float64 values.
		// This can only happen if there is a max size of 1.
		otel.Handle(errors.New("exponential histogram scale underflow"))
		return false
	}
	// Copy scale and min/max to cold
	coldIdx := (hotIdx + 1) % 2
	coldBucket := &b.hotColdBuckets[coldIdx]
	coldBucket.scale = hotBucket.scale
	startBin, endBin := hotBucket.startAndEnd.Load()
	coldBucket.startAndEnd.Store(startBin, endBin)
	// Downscale cold to the new scale
	coldBucket.downscale(scaleDelta)
	// Expand the cold prior to swapping to hot to ensure our measurement fits.
	bin = coldBucket.getBin(v)
	coldBucket.resizeToInclude(bin)

	coldBucket.startEndMux.Lock()
	defer coldBucket.startEndMux.Unlock()

	b.hcwg.swapHotAndWait()
	// Now that hot has become cold, downscale it, and merge it into the new hot buckets.
	hotBucket.downscale(scaleDelta)
	hotBucket.mergeIntoAndReset(coldBucket, b.maxScale)

	return coldBucket.recordBucket(bin)
}

// mergeInto merges the values of one hotColdExpoBuckets into another.
// The caller must already have exclusive access to b, and into can accept
// measurements concurrently with mergeInto.
func (b *hotColdExpoBuckets) mergeIntoAndReset(into *hotColdExpoBuckets) {
	// unifyScale is what is racing with writes
	b.unifyScale(into)
	bHotIdx := b.hcwg.loadHot()
	bBuckets := &b.hotColdBuckets[bHotIdx]
	intoHotIdx := into.hcwg.loadHot()
	intoBuckets := &into.hotColdBuckets[intoHotIdx]

	startBin, endBin := bBuckets.startAndEnd.Load()
	if startBin != endBin {
		intoBuckets.resizeToInclude(startBin)
		intoBuckets.resizeToInclude(endBin - 1)
	}
	scaleDelta := intoBuckets.scaleChange(endBin - 1)
	if scaleDelta > 0 {
		// Merging buckets required a scale change to the positive buckets to
		// fit within the max scale. Update scale and scale down the negative
		// buckets to match.
		b.downscale(scaleDelta, bHotIdx)
		into.downscale(scaleDelta, intoHotIdx)
	}
	b.hotColdBuckets[b.hcwg.loadHot()].mergeIntoAndReset(&into.hotColdBuckets[into.hcwg.loadHot()], b.maxScale)
}

// unifyScale downscales buckets as needed to make the scale of b and other
// the same. It returns the resulting scale. The caller must have exclusive
// access to both hotColdExpoBuckets.
func (b *hotColdExpoBuckets) unifyScale(other *hotColdExpoBuckets) int32 {
	bHotIdx := b.hcwg.loadHot()
	bScale := b.hotColdBuckets[bHotIdx].scale
	otherHotIdx := other.hcwg.loadHot()
	otherScale := other.hotColdBuckets[otherHotIdx].scale
	if bScale < otherScale {
		other.downscale(otherScale-bScale, otherHotIdx)
	} else if bScale > otherScale {
		b.downscale(bScale-otherScale, bHotIdx)
	}
	return min(bScale, otherScale)
}

// downscale force-downscales the bucket. It is assumed that the new scale is valid.
// The caller must hold the rescale mux.
func (b *hotColdExpoBuckets) downscale(delta int32, hotIdx uint64) {
	// Copy scale and min/max to cold
	coldIdx := (hotIdx + 1) % 2
	coldBucket := &b.hotColdBuckets[coldIdx]
	hotBucket := &b.hotColdBuckets[hotIdx]
	coldBucket.scale = hotBucket.scale
	startBin, endBin := hotBucket.startAndEnd.Load()

	coldBucket.startAndEnd.Store(startBin, endBin)
	// Downscale cold to the new scale
	coldBucket.downscale(delta)

	b.hcwg.swapHotAndWait()

	// Now that hot has become cold, downscale it, and merge it into the new hot buckets.
	hotBucket.downscale(delta)
	hotBucket.mergeIntoAndReset(coldBucket, b.maxScale)
}

// loadCountsAndOffset returns the buckets counts, the count, and the offset.
// It is not safe to call concurrently.
func (b *hotColdExpoBuckets) loadCountsAndOffset(buckets *[]uint64) (uint64, int32) {
	return b.hotColdBuckets[b.hcwg.loadHot()].loadCountsAndOffset(buckets)
}

// expoBuckets is a set of buckets in an exponential histogram.
type expoBuckets struct {
	scale       int32
	startEndMux sync.Mutex
	startAndEnd atomicLimitedRange
	counts      []atomic.Uint64
}

func newExpoBuckets(maxSize int, maxScale int32) expoBuckets {
	return expoBuckets{
		scale:  maxScale,
		counts: make([]atomic.Uint64, maxSize),
	}
}

// getIdx returns the index into counts for the provided bin.
func (e *expoBuckets) getIdx(bin int32) int {
	newBin := int(bin) % len(e.counts)
	return (newBin + len(e.counts)) % len(e.counts)
}

// loadCountsAndOffset returns the buckets counts, the count, and the offset.
// It is not safe to call concurrently.
func (e *expoBuckets) loadCountsAndOffset(into *[]uint64) (uint64, int32) {
	// TODO (#3047): Making copies for bounds and counts incurs a large
	// memory allocation footprint. Alternatives should be explored.
	start, end := e.startAndEnd.Load()
	length := int(end - start)
	counts := reset(*into, length, length)
	count := uint64(0)
	eIdx := start
	for i := range length {
		val := e.counts[e.getIdx(eIdx)].Load()
		counts[i] = val
		count += val
		eIdx++
	}
	*into = counts
	return count, start
}

// getBin returns the bin v should be recorded into.
func (p *expoBuckets) getBin(v float64) int32 {
	frac, expInt := math.Frexp(v)
	// 11-bit exponential.
	exp := int32(expInt) // nolint: gosec
	if p.scale <= 0 {
		// Because of the choice of fraction is always 1 power of two higher than we want.
		var correction int32 = 1
		if frac == .5 {
			// If v is an exact power of two the frac will be .5 and the exp
			// will be one higher than we want.
			correction = 2
		}
		return (exp - correction) >> (-p.scale)
	}
	return exp<<p.scale + int32(math.Log(frac)*scaleFactors[p.scale]) - 1
}

// scaleFactors are constants used in calculating the logarithm index. They are
// equivalent to 2^index/log(2).
var scaleFactors = [21]float64{
	math.Ldexp(math.Log2E, 0),
	math.Ldexp(math.Log2E, 1),
	math.Ldexp(math.Log2E, 2),
	math.Ldexp(math.Log2E, 3),
	math.Ldexp(math.Log2E, 4),
	math.Ldexp(math.Log2E, 5),
	math.Ldexp(math.Log2E, 6),
	math.Ldexp(math.Log2E, 7),
	math.Ldexp(math.Log2E, 8),
	math.Ldexp(math.Log2E, 9),
	math.Ldexp(math.Log2E, 10),
	math.Ldexp(math.Log2E, 11),
	math.Ldexp(math.Log2E, 12),
	math.Ldexp(math.Log2E, 13),
	math.Ldexp(math.Log2E, 14),
	math.Ldexp(math.Log2E, 15),
	math.Ldexp(math.Log2E, 16),
	math.Ldexp(math.Log2E, 17),
	math.Ldexp(math.Log2E, 18),
	math.Ldexp(math.Log2E, 19),
	math.Ldexp(math.Log2E, 20),
}

// scaleChange returns the magnitude of the scale change needed to fit bin in
// the bucket. If no scale change is needed 0 is returned.
func (b *expoBuckets) scaleChange(bin int32) int32 {
	startBin, endBin := b.startAndEnd.Load()
	if startBin == endBin {
		// No need to rescale if there are no buckets.
		return 0
	}

	lastBin := endBin - 1
	if bin < startBin {
		startBin = bin
	} else if bin > lastBin {
		lastBin = bin
	}

	var count int32
	for lastBin-startBin >= int32(len(b.counts)) {
		startBin >>= 1
		lastBin >>= 1
		count++
		if count > expoMaxScale-expoMinScale {
			return count
		}
	}
	return count
}

// recordBucket returns true if the bucket was incremented, or false if a downscale is required to
func (b *expoBuckets) recordBucket(bin int32) bool {
	startBin, endBin := b.startAndEnd.Load()
	if bin >= startBin && bin < endBin {
		b.counts[b.getIdx(bin)].Add(1)
		return true
	}
	return false
}

// downscale shrinks a bucket by a factor of 2*s. It will sum counts into the
// correct lower resolution bucket. downscale is not concurrent safe.
func (b *expoBuckets) downscale(delta int32) {
	b.scale -= delta
	// Example
	// delta = 2
	// Original offset: -6
	// Counts: [ 3,  1,  2,  3,  4,  5, 6, 7, 8, 9, 10]
	// bins:    -6  -5, -4, -3, -2, -1, 0, 1, 2, 3, 4
	// new bins:-2, -2, -1, -1, -1, -1, 0, 0, 0, 0, 1
	// new Offset: -2
	// new Counts: [4, 14, 30, 10]

	startBin, endBin := b.startAndEnd.Load()
	length := endBin - startBin
	if length <= 1 || delta < 1 {
		newStartBin := startBin >> delta
		newEndBin := newStartBin + length
		b.startAndEnd.Store(newStartBin, newEndBin)
		// Shift all elements left by the change in start position
		startShift := b.getIdx(startBin - newStartBin)
		b.counts = append(b.counts[startShift:], b.counts[:startShift]...)

		// Clear all elements that are outside of our start to end range
		for i := newEndBin; i < newStartBin+int32(len(b.counts)); i++ {
			b.counts[b.getIdx(i)].Store(0)
		}
		return
	}

	steps := int32(1) << delta
	offset := startBin % steps
	offset = (offset + steps) % steps // to make offset positive
	newLen := (length-1+offset)/steps + 1
	newStartBin := startBin >> delta
	newEndBin := newStartBin + newLen
	startShift := b.getIdx(startBin - newStartBin)

	for i := startBin + 1; i < endBin; i++ {
		newIdx := b.getIdx(int32(math.Floor(float64(i)/float64(steps))) + int32(startShift))
		if i%steps == 0 {
			b.counts[newIdx].Store(b.counts[b.getIdx(i)].Load())
			continue
		}
		b.counts[newIdx].Add(b.counts[b.getIdx(i)].Load())
	}
	b.startAndEnd.Store(newStartBin, newEndBin)
	// Shift all elements left by the change in start position
	b.counts = append(b.counts[startShift:], b.counts[:startShift]...)

	// Clear all elements that are outside of our start to end range
	for i := newEndBin; i < newStartBin+int32(len(b.counts)); i++ {
		b.counts[b.getIdx(i)].Store(0)
	}
}

// resizeToInclude force-expands the range of b to include the bin.
// resizeToInclude is not safe to call concurrently.
func (b *expoBuckets) resizeToInclude(bin int32) {
	b.startEndMux.Lock()
	defer b.startEndMux.Unlock()
	startBin, endBin := b.startAndEnd.Load()
	if startBin == endBin {
		startBin = bin
		endBin = bin + 1
		b.startAndEnd.Store(startBin, endBin)
	} else if bin < startBin {
		b.startAndEnd.Store(bin, endBin)
	} else if bin >= endBin {
		b.startAndEnd.Store(startBin, bin+1)
	}
}

// mergeInto merges this expoBuckets into another, and resets the state
// of the expoBuckets. This is used to ensure that the cumulative counters
// continue to accumulate after being read.
// mergeInto requires that scales are equal.
func (b *expoBuckets) mergeIntoAndReset(into *expoBuckets, maxScale int32) {
	if b.scale != into.scale {
		panic("scales not equal")
	}
	intoStartBin, intoEndBin := into.startAndEnd.Load()
	bStartBin, bEndBin := b.startAndEnd.Load()
	b.startAndEnd.Store(0, 0)
	if bStartBin != bEndBin && (intoStartBin > bStartBin || intoEndBin < bEndBin) {
		panic(fmt.Sprintf("into is not a superset of b.  intoStartBin %v, bStartBin %v, intoEndBin %v, bEndBin %v", intoStartBin, bStartBin, intoEndBin, bEndBin))
	}
	for i := bStartBin; i < bEndBin; i++ {
		// Swap in 0 to reset
		val := b.counts[b.getIdx(int32(i))].Swap(0)
		into.counts[into.getIdx(int32(i))].Add(val)
	}
	b.scale = maxScale
}

// newDeltaExponentialHistogram returns an Aggregator that summarizes a set of
// measurements as a delta exponential histogram. Each histogram is scoped by
// attributes and the aggregation cycle the measurements were made in.
func newDeltaExponentialHistogram[N int64 | float64](
	maxSize, maxScale int32,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *deltaExpoHistogram[N] {
	return &deltaExpoHistogram[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		maxSize:  int(maxSize),
		maxScale: maxScale,

		newRes: r,
		hotColdValMap: [2]limitedSyncMap{
			{aggLimit: limit},
			{aggLimit: limit},
		},

		start: now(),
	}
}

// deltaExpoHistogram summarizes a set of measurements as an histogram with exponentially
// defined buckets.
type deltaExpoHistogram[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int32

	newRes        func(attribute.Set) FilteredExemplarReservoir[N]
	hcwg          hotColdWaitGroup
	hotColdValMap [2]limitedSyncMap

	start time.Time
}

func (e *deltaExpoHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	hotIdx := e.hcwg.start()
	defer e.hcwg.done(hotIdx)
	v := e.hotColdValMap[hotIdx].LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		hPt := newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale)
		hPt.res = e.newRes(attr)
		return hPt
	}).(*expoHistogramDataPoint[N])
	if !v.tryFastRecord(value, e.noMinMax, e.noSum) {
		v.rescaleMux.Lock()
		v.record(value, e.noMinMax, e.noSum)
		v.rescaleMux.Unlock()
	}
	v.res.Offer(ctx, value, droppedAttr)
}

func (e *deltaExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.DeltaTemporality

	// delta always clears values on collection
	readIdx := e.hcwg.swapHotAndWait()

	// The len will not change while we iterate over values, since we waited
	// for all writes to finish to the cold values and len.
	n := e.hotColdValMap[readIdx].Len()
	hDPts := reset(h.DataPoints, n, n)

	var i int
	e.hotColdValMap[readIdx].Range(func(_, value any) bool {
		val := value.(*expoHistogramDataPoint[N])
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = e.start
		hDPts[i].Time = t
		hDPts[i].ZeroThreshold = 0.0

		val.rescaleMux.Lock()
		defer val.rescaleMux.Unlock()
		val.loadInto(&hDPts[i], e.noMinMax, e.noSum)
		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
		return true
	})
	// Unused attribute sets do not report.
	e.hotColdValMap[readIdx].Clear()

	e.start = t
	h.DataPoints = hDPts
	*dest = h
	return n
}

// newCumulativeExponentialHistogram returns an Aggregator that summarizes a
// set of measurements as a cumulative exponential histogram. Each histogram is
// scoped by attributes and the aggregation cycle the measurements were made
// in.
func newCumulativeExponentialHistogram[N int64 | float64](
	maxSize, maxScale int32,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *cumulativeExpoHistogram[N] {
	return &cumulativeExpoHistogram[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		maxSize:  int(maxSize),
		maxScale: maxScale,

		newRes: r,
		values: limitedSyncMap{aggLimit: limit},

		start: now(),
	}
}

// cumulativeExpoHistogram summarizes a set of measurements as an cumulative
// histogram with exponentially defined buckets.
type cumulativeExpoHistogram[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int32

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	values limitedSyncMap

	start time.Time
}

func (e *cumulativeExpoHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	v := e.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		hPt := newHotColdExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale)
		hPt.res = e.newRes(attr)
		return hPt
	}).(*hotColdExpoHistogramPoint[N])

	hotIdx := v.hcwg.start()
	fastRecordSuccess := v.hotColdPoint[hotIdx].tryFastRecord(value, e.noMinMax, e.noSum)
	v.hcwg.done(hotIdx)
	if !fastRecordSuccess {
		v.rescaleMux.Lock()
		// we hold the lock, so no need to use start/end from hcwg
		v.hotColdPoint[v.hcwg.loadHot()].record(value, e.noMinMax, e.noSum)
		v.rescaleMux.Unlock()
	}
	v.res.Offer(ctx, value, droppedAttr)
}

func (e *cumulativeExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	// Values are being concurrently written while we iterate, so only use the
	// current length for capacity.
	hDPts := reset(h.DataPoints, 0, e.values.Len())

	var i int
	e.values.Range(func(_, value any) bool {
		val := value.(*hotColdExpoHistogramPoint[N])
		newPt := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    val.attrs,
			StartTime:     e.start,
			Time:          t,
			ZeroThreshold: 0.0,
		}
		val.rescaleMux.Lock()
		defer val.rescaleMux.Unlock()
		// Prevent buckets from changing start or end ranges during collection
		// so mergeInto below succeeds.
		hotIdx := val.hcwg.loadHot()
		val.hotColdPoint[hotIdx].posBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].posBuckets.hcwg.loadHot()].startEndMux.Lock()
		defer val.hotColdPoint[hotIdx].posBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].posBuckets.hcwg.loadHot()].startEndMux.Unlock()
		val.hotColdPoint[hotIdx].negBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].negBuckets.hcwg.loadHot()].startEndMux.Lock()
		defer val.hotColdPoint[hotIdx].negBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].negBuckets.hcwg.loadHot()].startEndMux.Unlock()

		// Set the range of the cold point to the hot range to ensure we don't
		// accept measurements that would have resulted in an underflow, and
		// would block mergeInto below.
		coldIdx := (hotIdx + 1) % 2
		start, end := val.hotColdPoint[hotIdx].posBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].posBuckets.hcwg.loadHot()].startAndEnd.Load()
		val.hotColdPoint[coldIdx].posBuckets.hotColdBuckets[val.hotColdPoint[coldIdx].posBuckets.hcwg.loadHot()].startAndEnd.Store(start, end)
		start, end = val.hotColdPoint[hotIdx].negBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].negBuckets.hcwg.loadHot()].startAndEnd.Load()
		val.hotColdPoint[coldIdx].negBuckets.hotColdBuckets[val.hotColdPoint[coldIdx].negBuckets.hcwg.loadHot()].startAndEnd.Store(start, end)

		// Set the scale to the minimum of the pos and negative bucket scale.
		posScale := val.hotColdPoint[hotIdx].posBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].posBuckets.hcwg.loadHot()].scale
		negScale := val.hotColdPoint[hotIdx].negBuckets.hotColdBuckets[val.hotColdPoint[hotIdx].negBuckets.hcwg.loadHot()].scale
		newScale := min(posScale, negScale)
		val.hotColdPoint[coldIdx].posBuckets.hotColdBuckets[val.hotColdPoint[coldIdx].posBuckets.hcwg.loadHot()].scale = newScale
		val.hotColdPoint[coldIdx].negBuckets.hotColdBuckets[val.hotColdPoint[coldIdx].negBuckets.hcwg.loadHot()].scale = newScale

		// Swap so we can read from hot
		readIdx := val.hcwg.swapHotAndWait()
		readPt := &val.hotColdPoint[readIdx]
		readPt.loadInto(&newPt, e.noMinMax, e.noSum)
		// Once we've read the point, merge it back into the now-hot histogram
		// point since it is cumulative.
		readPt.mergeIntoAndReset(&val.hotColdPoint[(readIdx+1)%2], e.noMinMax, e.noSum)

		collectExemplars(&newPt.Exemplars, val.res.Collect)
		hDPts = append(hDPts, newPt)

		i++
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
		return true
	})

	h.DataPoints = hDPts
	*dest = h
	return i
}
