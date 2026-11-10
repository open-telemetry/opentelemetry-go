// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	expoMaxScale   = 20
	expoMinScale   = -10
	defaultMaxSize = 160

	smallestNonZeroNormalFloat64 = 0x1p-1022
)

// expoHistogramDataPoint is a single data point in an exponential histogram.
// For maxSize > 2, record uses wg for lockless fast-path bucket increments and
// mu + wg.swapHotAndWait for slow-path slice growth and scale downscaling, while
// callers drain concurrent record calls via the outer aggregator hotColdWaitGroup
// before calling merge, reset, or uploadTo.
// For maxSize <= 2, record, merge, reset, and uploadTo are serialized by mu.
type expoHistogramDataPoint[N int64 | float64] struct {
	attrs         attribute.Set
	res           FilteredExemplarReservoir[N]
	dropExemplars bool

	minMax atomicMinMax[N]
	sum    atomicCounter[N]

	maxSize  int
	noMinMax bool
	noSum    bool

	mu sync.Mutex
	wg hotColdWaitGroup

	scale atomic.Int32

	posBuckets expoBuckets
	negBuckets expoBuckets
	zeroCount  atomic.Uint64
	startTime  time.Time
}

func newExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
	noMinMax, noSum bool,
	newRes func(attribute.Set) FilteredExemplarReservoir[N],
	startTime time.Time,
) *expoHistogramDataPoint[N] { // nolint:revive // we need this control flag
	dp := &expoHistogramDataPoint[N]{
		attrs:     attrs,
		maxSize:   maxSize,
		noMinMax:  noMinMax,
		noSum:     noSum,
		startTime: startTime,
	}
	dp.scale.Store(maxScale)
	if newRes != nil {
		dp.res = newRes(attrs)
		_, isDrop := dp.res.(*dropRes[N])
		dp.dropExemplars = isDrop
	}
	return dp
}

func (p *expoHistogramDataPoint[N]) reset(startTime time.Time, maxScale int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.scale.Store(maxScale)
	p.posBuckets.clear()
	p.negBuckets.clear()
	p.zeroCount.Store(0)
	p.sum.reset()
	p.minMax.reset()
	p.startTime = startTime
}

// record adds a new measurement to the histogram. It will rescale the buckets if needed.
func (p *expoHistogramDataPoint[N]) record(v N) {
	if p.maxSize <= 2 {
		p.recordLocked(v)
		return
	}

	if !p.noMinMax {
		p.minMax.Update(v)
	}
	if !p.noSum {
		p.sum.add(v)
	}

	absV := math.Abs(float64(v))
	if absV == 0.0 {
		p.zeroCount.Add(1)
		return
	}

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	idx := p.wg.start()
	if idx == 0 {
		ok := bucket.tryRecord(p.getBin(absV), 1, p.maxSize)
		p.wg.done(0)
		if ok {
			return
		}
	} else {
		p.wg.done(idx)
	}

	p.recordSlow(absV, bucket)
}

func (p *expoHistogramDataPoint[N]) recordSlow(absV float64, bucket *expoBuckets) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if bucket.tryRecord(p.getBin(absV), 1, p.maxSize) {
		return
	}

	_ = p.wg.swapHotAndWait()

	startBin, length := bucket.loadBounds()
	bin := p.getBin(absV)
	if scaleDelta := p.scaleChange(bin, startBin, int(length)); scaleDelta > 0 {
		p.scale.Add(-scaleDelta)
		p.posBuckets.downscale(scaleDelta)
		p.negBuckets.downscale(scaleDelta)
		bin = p.getBin(absV)
	}
	bucket.record(bin)

	_ = p.wg.swapHotAndWait()
}

func (p *expoHistogramDataPoint[N]) recordLocked(v N) {
	p.mu.Lock()
	defer p.mu.Unlock()

	absV := math.Abs(float64(v))

	if absV == 0.0 {
		p.zeroCount.Add(1)
		if !p.noMinMax {
			p.minMax.Update(v)
		}
		if !p.noSum {
			p.sum.add(v)
		}
		return
	}

	bin := p.getBin(absV)

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	startBin, length := bucket.loadBounds()
	if scaleDelta := p.scaleChange(bin, startBin, int(length)); scaleDelta > 0 {
		currentScale := p.scale.Load()
		if currentScale-scaleDelta < expoMinScale {
			// With a scale of -10 there is only two buckets for the whole range of float64 values.
			// This can only happen if there is a max size of 1.
			otel.Handle(errors.New("exponential histogram scale underflow"))
			return
		}
		// Downscale
		p.scale.Add(-scaleDelta)
		p.posBuckets.downscale(scaleDelta)
		p.negBuckets.downscale(scaleDelta)

		bin = p.getBin(absV)
	}

	if !p.noMinMax {
		p.minMax.Update(v)
	}
	if !p.noSum {
		p.sum.add(v)
	}
	bucket.record(bin)
}

// getBin returns the bin v should be recorded into.
func (p *expoHistogramDataPoint[N]) getBin(v float64) int32 {
	frac, expInt := math.Frexp(v)
	// 11-bit exponential.
	exp := int32(expInt) // nolint: gosec
	scale := p.scale.Load()
	if scale <= 0 {
		// Because of the choice of fraction is always 1 power of two higher than we want.
		var correction int32 = 1
		if frac == .5 {
			// If v is an exact power of two the frac will be .5 and the exp
			// will be one higher than we want.
			correction = 2
		}
		return (exp - correction) >> (-scale)
	}
	return exp<<scale + int32(math.Log(frac)*scaleFactors[scale]) - 1
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
func (p *expoHistogramDataPoint[N]) scaleChange(bin, startBin int32, length int) int32 {
	if length == 0 {
		// No need to rescale if there are no buckets.
		return 0
	}

	low := int64(startBin)
	high := int64(bin)
	if startBin >= bin {
		low = int64(bin)
		high = int64(startBin) + int64(length) - 1
	}

	var count int32
	for high-low >= int64(p.maxSize) {
		low >>= 1
		high >>= 1
		count++
		if count > expoMaxScale-expoMinScale {
			return count
		}
	}
	return count
}

func (p *expoHistogramDataPoint[N]) count() uint64 {
	return p.posBuckets.count() + p.negBuckets.count() + p.zeroCount.Load()
}

const minBucketCapacity = 1

// expoBuckets is a circular buffer of bucket counts in an exponential histogram.
//
// The logical range is [startBin, startBin + length - 1], packed into bounds.
// Physical storage is counts []atomic.Uint64 with len(counts) always a power of 2 (or nil when unallocated).
// Indexing into counts uses bitwise AND mask: int(bin) & (len(counts)-1).
// Invariant: every slot in counts outside the logical range holds zero.
type expoBuckets struct {
	// bounds packs startBin (upper 32 bits) and length (lower 32 bits).
	bounds atomic.Uint64
	counts []atomic.Uint64
}

func packBounds(startBin, length int32) uint64 {
	//nolint:gosec // intentional bit-packing of signed int32 values
	return (uint64(uint32(startBin)) << 32) | uint64(uint32(length))
}

func (b *expoBuckets) loadBounds() (startBin, length int32) {
	v := b.bounds.Load()
	//nolint:gosec // intentional bit-unpacking into signed int32 values
	return int32(uint32(v >> 32)), int32(uint32(v))
}

func (b *expoBuckets) storeBounds(startBin, length int32) {
	b.bounds.Store(packBounds(startBin, length))
}

func (b *expoBuckets) swapBounds(oldStart, oldLen, newStart, newLen int32) bool {
	return b.bounds.CompareAndSwap(packBounds(oldStart, oldLen), packBounds(newStart, newLen))
}

// index returns the slot in counts for bin. It requires len(b.counts) > 0.
func (b *expoBuckets) index(bin int32) int {
	return int(bin) & (len(b.counts) - 1)
}

// bin returns a pointer to the bucket count for the given logical bin.
func (b *expoBuckets) bin(bin int32) *atomic.Uint64 {
	return &b.counts[b.index(bin)]
}

func nextPowerOf2(n int32) int {
	if n <= minBucketCapacity {
		return minBucketCapacity
	}
	v := uint32(n - 1)
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	return int(v + 1)
}

func (b *expoBuckets) clear() {
	startBin, length := b.loadBounds()
	for i := range length {
		b.bin(startBin + i).Store(0)
	}
	b.storeBounds(0, 0)
}

// tryRecord attempts to record count into bin without reallocating counts or rescaling.
// It returns true if the measurement was recorded (either within the current bounds or
// by atomically expanding bounds within the existing counts capacity and maxSize),
// and false if counts must be grown or the histogram must be downscaled.
func (b *expoBuckets) tryRecord(bin int32, count uint64, maxSize int) bool {
	if count == 0 {
		return true
	}
	if len(b.counts) == 0 {
		return false
	}
	for {
		startBin, length := b.loadBounds()
		if length == 0 {
			if maxSize < 1 {
				return false
			}
			if b.swapBounds(startBin, 0, bin, 1) {
				b.bin(bin).Add(count)
				return true
			}
			continue
		}

		endBin := startBin + length - 1
		if bin >= startBin && bin <= endBin {
			b.bin(bin).Add(count)
			return true
		}

		newStartBin := min(bin, startBin)
		newEndBin := max(bin, endBin)
		newLength := int64(newEndBin) - int64(newStartBin) + 1
		if newLength > int64(len(b.counts)) || newLength > int64(maxSize) {
			return false
		}

		//nolint:gosec // newLength <= maxSize fits in int32
		if b.swapBounds(startBin, length, newStartBin, int32(newLength)) {
			b.bin(bin).Add(count)
			return true
		}
	}
}

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(bin int32) {
	b.recordCount(bin, 1)
}

func (b *expoBuckets) recordCount(bin int32, count uint64) {
	if count == 0 {
		return
	}
	startBin, length := b.loadBounds()
	if length == 0 {
		if len(b.counts) == 0 {
			b.counts = make([]atomic.Uint64, minBucketCapacity)
		}
		b.storeBounds(bin, 1)
		b.bin(bin).Store(count)
		return
	}

	endBin := startBin + length - 1
	if bin >= startBin && bin <= endBin {
		b.bin(bin).Add(count)
		return
	}

	newStartBin := min(bin, startBin)
	newEndBin := max(bin, endBin)
	newLength := newEndBin - newStartBin + 1

	if int(newLength) > len(b.counts) {
		newCap := nextPowerOf2(newLength)
		newCounts := make([]atomic.Uint64, newCap)
		mask := newCap - 1
		for i := range length {
			oldBin := startBin + i
			if val := b.bin(oldBin).Load(); val > 0 {
				newCounts[int(oldBin)&mask].Store(val)
			}
		}
		b.counts = newCounts
	}

	b.storeBounds(newStartBin, newLength)
	b.bin(bin).Add(count)
}

func (b *expoBuckets) uploadTo(dest *metricdata.ExponentialBucket) {
	startBin, length := b.loadBounds()
	dest.Offset = startBin
	dest.Counts = reset(dest.Counts, int(length), int(length))
	for i := range length {
		dest.Counts[i] = b.bin(startBin + i).Load()
	}
}

func (b *expoBuckets) merge(other *expoBuckets) {
	otherStart, otherLen := other.loadBounds()
	if otherLen == 0 {
		return
	}
	for i := range otherLen {
		if c := other.bin(otherStart + i).Load(); c > 0 {
			b.recordCount(otherStart+i, c)
		}
	}
}

func (p *expoHistogramDataPoint[N]) mergeScaleDelta(b, other *expoBuckets, bShift, oShift int32) int32 {
	bStart, bLen := b.loadBounds()
	oStart, oLen := other.loadBounds()
	if bLen == 0 || oLen == 0 {
		return 0
	}
	bEnd := bStart + bLen - 1
	oEnd := oStart + oLen - 1
	minBin := min(bStart>>bShift, oStart>>oShift)
	maxBin := max(bEnd>>bShift, oEnd>>oShift)
	return p.scaleChange(maxBin, minBin, 1)
}

func (p *expoHistogramDataPoint[N]) merge(other *expoHistogramDataPoint[N]) {
	// Caller must ensure other is idle (no concurrent record/merge), as merge reads and mutates other.
	// Note: scale underflow below expoMinScale cannot happen during merge because maxSize <= 2 is routed to
	// lockedCumulativeExpoHistogram, and maxSize > 2 guarantees scale changes remain bounded within [expoMinScale, expoMaxScale].
	p.mu.Lock()
	defer p.mu.Unlock()

	pStartScale := p.scale.Load()
	oStartScale := other.scale.Load()

	targetScale := min(pStartScale, oStartScale)
	pAlignShift := pStartScale - targetScale
	oAlignShift := oStartScale - targetScale

	d := max(
		p.mergeScaleDelta(&p.posBuckets, &other.posBuckets, pAlignShift, oAlignShift),
		p.mergeScaleDelta(&p.negBuckets, &other.negBuckets, pAlignShift, oAlignShift),
	)

	pDownscale := pAlignShift + d
	if pDownscale > 0 {
		p.scale.Add(-pDownscale)
		p.posBuckets.downscale(pDownscale)
		p.negBuckets.downscale(pDownscale)
	}

	oDownscale := oAlignShift + d
	if oDownscale > 0 {
		other.scale.Add(-oDownscale)
		other.posBuckets.downscale(oDownscale)
		other.negBuckets.downscale(oDownscale)
	}

	p.posBuckets.merge(&other.posBuckets)
	p.negBuckets.merge(&other.negBuckets)

	if !p.noSum {
		p.sum.add(other.sum.load())
	}
	if !p.noMinMax {
		if other.minMax.set.Load() {
			p.minMax.Update(other.minMax.minimum.Load())
			p.minMax.Update(other.minMax.maximum.Load())
		}
	}
	p.zeroCount.Add(other.zeroCount.Load())
}

func (p *expoHistogramDataPoint[N]) uploadTo(
	dest *metricdata.ExponentialHistogramDataPoint[N],
	startTime, endTime time.Time,
) {
	p.mu.Lock()
	defer p.mu.Unlock()

	dest.Attributes = p.attrs
	dest.StartTime = startTime
	dest.Time = endTime
	dest.Count = p.count()
	dest.Scale = p.scale.Load()
	dest.ZeroCount = p.zeroCount.Load()
	dest.ZeroThreshold = 0.0

	p.posBuckets.uploadTo(&dest.PositiveBucket)
	p.negBuckets.uploadTo(&dest.NegativeBucket)

	if !p.noSum {
		dest.Sum = p.sum.load()
	} else {
		var zero N
		dest.Sum = zero
	}
	if !p.noMinMax && p.minMax.set.Load() {
		dest.Min = metricdata.NewExtrema(p.minMax.minimum.Load())
		dest.Max = metricdata.NewExtrema(p.minMax.maximum.Load())
	} else {
		dest.Min = metricdata.Extrema[N]{}
		dest.Max = metricdata.Extrema[N]{}
	}
}

// downscale shrinks a bucket by a factor of 2^delta. It will sum counts into the
// correct lower resolution bucket.
func (b *expoBuckets) downscale(delta int32) {
	startBin, length := b.loadBounds()
	if length == 0 || delta < 1 {
		return
	}

	endBin := startBin + length - 1
	newStartBin := startBin >> delta
	newEndBin := endBin >> delta
	newLength := newEndBin - newStartBin + 1

	// stackScratch is sized to defaultMaxSize/2 + 1 (81), the maximum possible
	// downscaled span for defaultMaxSize (160) when delta >= 1, avoiding heap
	// allocations under the default configuration.
	var stackScratch [defaultMaxSize/2 + 1]uint64
	var scratch []uint64
	if int(newLength) <= len(stackScratch) {
		scratch = stackScratch[:newLength]
	} else {
		scratch = make([]uint64, newLength)
	}

	for i := range length {
		slot := b.bin(startBin + i)
		if val := slot.Load(); val > 0 {
			newBin := (startBin + i) >> delta
			scratch[newBin-newStartBin] += val
			slot.Store(0)
		}
	}

	b.storeBounds(newStartBin, newLength)
	for i := range newLength {
		if val := scratch[i]; val > 0 {
			b.bin(newStartBin + i).Store(val)
		}
	}
}

func (b *expoBuckets) count() uint64 {
	startBin, length := b.loadBounds()
	var total uint64
	for i := range length {
		total += b.bin(startBin + i).Load()
	}
	return total
}

// newDeltaExpoHistogram returns an Aggregator that summarizes a set of
// measurements as a delta exponential histogram.
func newDeltaExpoHistogram[N int64 | float64](
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
		hotColdValMap: [2]limitedSyncMap[*expoHistogramDataPoint[N]]{
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
	hotColdValMap [2]limitedSyncMap[*expoHistogramDataPoint[N]]

	start time.Time
}

func (e *deltaExpoHistogram[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	hotIdx := e.hcwg.start()
	defer e.hcwg.done(hotIdx)

	v := e.hotColdValMap[hotIdx].LoadOrStoreAttr(lazy, func(attr attribute.Set) *expoHistogramDataPoint[N] {
		return newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum, e.newRes, now())
	})

	v.record(value)
	if !v.dropExemplars {
		v.res.Offer(ctx, value, lazy)
	}
}

func (e *deltaExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, // nolint:gocritic // dest is an interface pointer used to avoid allocations
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.DeltaTemporality

	coldIdx := e.hcwg.swapHotAndWait()

	n := e.hotColdValMap[coldIdx].Len()
	hDPts := reset(h.DataPoints, n, n)

	var i int
	e.hotColdValMap[coldIdx].Range(func(_, value any) bool {
		val := value.(*expoHistogramDataPoint[N])

		val.uploadTo(&hDPts[i], e.start, t)

		if val.res != nil {
			collectExemplars(&hDPts[i].Exemplars, val.res.Collect)
		} else {
			hDPts[i].Exemplars = hDPts[i].Exemplars[:0]
		}

		i++
		return true
	})

	e.start = t
	e.hotColdValMap[coldIdx].Clear()

	h.DataPoints = hDPts
	*dest = h
	return n
}

// newCumulativeExpoHistogram returns an Aggregator that summarizes a set of
// measurements as a cumulative exponential histogram.
// It returns the measure and collect functions directly, hiding the underlying
// implementation (double-buffered or locked) from the caller.
func newCumulativeExpoHistogram[N int64 | float64](
	maxSize, maxScale int32,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) (fltrMeasure[N], func(*metricdata.Aggregation) int) {
	// For maxSize <= 2, use lockedCumulativeExpoHistogram to avoid scale underflow issues
	// during merge while reserving doubleBufferedCumulativeExpoHistogram for high-performance
	// concurrent aggregation at larger max sizes.
	if maxSize <= 2 {
		h := &lockedCumulativeExpoHistogram[N]{
			noSum:    noSum,
			noMinMax: noMinMax,
			maxSize:  int(maxSize),
			maxScale: maxScale,
			newRes:   r,
			values:   limitedSyncMap[*expoHistogramDataPoint[N]]{aggLimit: limit},
			start:    now(),
		}
		return h.measure, h.collect
	}

	h := &doubleBufferedCumulativeExpoHistogram[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		maxSize:  int(maxSize),
		maxScale: maxScale,
		newRes:   r,
		values:   limitedSyncMap[*doubleBufferedCumulativePoint[N]]{aggLimit: limit},
		start:    now(),
	}
	return h.measure, h.collect
}

// doubleBufferedCumulativePoint is the high-performance double-buffered
// implementation of a single cumulative series. It uses two active points
// (hot/cold) and one persistent cumulative point.
type doubleBufferedCumulativePoint[N int64 | float64] struct {
	wg         hotColdWaitGroup
	points     [2]*expoHistogramDataPoint[N]
	cumulative *expoHistogramDataPoint[N]

	maxScale int32
}

func (cp *doubleBufferedCumulativePoint[N]) measure(ctx context.Context, value N, lazy lazyFilteredAttributes) {
	hotIdx := cp.wg.start()
	defer cp.wg.done(hotIdx)
	v := cp.points[hotIdx]
	v.record(value)
	if cp.cumulative.res != nil && !cp.cumulative.dropExemplars {
		// Note: Measurements occurring after swapHotAndWait can offer an exemplar to
		// cp.cumulative.res before collectExemplars runs, which is an accepted tradeoff.
		cp.cumulative.res.Offer(ctx, value, lazy)
	}
}

func (cp *doubleBufferedCumulativePoint[N]) collect(
	dp *metricdata.ExponentialHistogramDataPoint[N],
	t, startTime time.Time,
) {
	val := cp.cumulative
	// swapHotAndWait makes sure all active writes to coldIdx are done before merge runs,
	// allowing merge to mutate delta without needing delta.mu.
	coldIdx := cp.wg.swapHotAndWait()
	delta := cp.points[coldIdx]

	val.merge(delta)

	// Reset the inactive cold point for reuse without allocation.
	delta.reset(t, cp.maxScale)

	val.uploadTo(dp, startTime, t)

	if val.res != nil && !val.dropExemplars {
		collectExemplars(&dp.Exemplars, val.res.Collect)
	} else {
		dp.Exemplars = dp.Exemplars[:0]
	}
}

// doubleBufferedCumulativeExpoHistogram summarizes a set of measurements as a cumulative
// exponential histogram using double-buffering for high performance.
type doubleBufferedCumulativeExpoHistogram[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int32

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	values limitedSyncMap[*doubleBufferedCumulativePoint[N]]
	start  time.Time
}

func (e *doubleBufferedCumulativeExpoHistogram[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	vp := e.values.LoadOrStoreAttr(lazy, func(attr attribute.Set) *doubleBufferedCumulativePoint[N] {
		n := now()
		cp := &doubleBufferedCumulativePoint[N]{
			maxScale: e.maxScale,
		}
		cp.cumulative = newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum, e.newRes, n)
		cp.points[0] = newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum, nil, n)
		cp.points[1] = newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum, nil, n)
		return cp
	})

	vp.measure(ctx, value, lazy)
}

func (e *doubleBufferedCumulativeExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, // nolint:gocritic // dest is an interface pointer used to avoid allocations
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	n := e.values.Len()
	hDPts := reset(h.DataPoints, n, n)

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	e.values.Range(func(_, value any) bool {
		cp := value.(*doubleBufferedCumulativePoint[N])

		if i >= len(hDPts) {
			hDPts = append(hDPts, metricdata.ExponentialHistogramDataPoint[N]{})
		}
		dp := &hDPts[i]

		startTime := e.start
		if perSeriesStartTimeEnabled {
			startTime = cp.cumulative.startTime
		}

		cp.collect(dp, t, startTime)
		i++
		return true
	})

	// TODO (#3006): This will use an unbounded amount of memory if there
	// are unbounded number of attribute sets being aggregated. Attribute
	// sets that become "stale" need to be forgotten so this will not
	// overload the system.

	h.DataPoints = hDPts[:i]
	*dest = h
	return i
}

// lockedCumulativeExpoHistogram summarizes a set of measurements as a cumulative
// exponential histogram using a simple locked implementation. This is used when
// maxSize <= 2 to avoid scale underflow issues during merge.
type lockedCumulativeExpoHistogram[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int32

	newRes func(attribute.Set) FilteredExemplarReservoir[N]
	values limitedSyncMap[*expoHistogramDataPoint[N]]
	start  time.Time
}

func (e *lockedCumulativeExpoHistogram[N]) measure(
	ctx context.Context,
	value N,
	lazy lazyFilteredAttributes,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	vp := e.values.LoadOrStoreAttr(lazy, func(attr attribute.Set) *expoHistogramDataPoint[N] {
		return newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum, e.newRes, now())
	})
	vp.record(value)
	if vp.res != nil && !vp.dropExemplars {
		vp.res.Offer(ctx, value, lazy)
	}
}

func (e *lockedCumulativeExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, // nolint:gocritic // dest is an interface pointer used to avoid allocations
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	n := e.values.Len()
	hDPts := reset(h.DataPoints, n, n)

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	e.values.Range(func(_, value any) bool {
		val := value.(*expoHistogramDataPoint[N])

		if i >= len(hDPts) {
			hDPts = append(hDPts, metricdata.ExponentialHistogramDataPoint[N]{})
		}
		dp := &hDPts[i]

		startTime := e.start
		if perSeriesStartTimeEnabled {
			startTime = val.startTime
		}

		val.uploadTo(dp, startTime, t)

		if val.res != nil && !val.dropExemplars {
			collectExemplars(&dp.Exemplars, val.res.Collect)
		} else {
			dp.Exemplars = dp.Exemplars[:0]
		}

		i++
		return true
	})

	// TODO (#3006): This will use an unbounded amount of memory if there
	// are unbounded number of attribute sets being aggregated. Attribute
	// sets that become "stale" need to be forgotten so this will not
	// overload the system.

	h.DataPoints = hDPts[:i]
	*dest = h
	return i
}
