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

// hotColdExpoHistogramPoint a hot and cold exponential histogram points, used
// in cumulative aggregations.
type hotColdExpoHistogramPoint[N int64 | float64] struct {
	rescaleMux   sync.Mutex
	hcwg         hotColdWaitGroup
	hotColdPoint [2]expoHistogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	noMinMax bool
	noSum    bool
	maxScale int32
}

func newHotColdExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
	noMinMax, noSum bool,
) *hotColdExpoHistogramPoint[N] { // nolint:revive // we need this control flag
	return &hotColdExpoHistogramPoint[N]{
		attrs:    attrs,
		noMinMax: noMinMax,
		noSum:    noSum,
		maxScale: maxScale,
		hotColdPoint: [2]expoHistogramPointCounters[N]{
			newExpoHistogramPointCounters[N](maxSize, maxScale),
			newExpoHistogramPointCounters[N](maxSize, maxScale),
		},
	}
}

// expoHistogramPointCounters contains only the atomic counter data, and is
// used by both expoHistogramDataPoint and hotColdExpoHistogramPoint.
type expoHistogramPointCounters[N int64 | float64] struct {
	minMax    atomicMinMax[N]
	sum       atomicCounter[N]
	zeroCount atomic.Uint64

	posBuckets expoBuckets
	negBuckets expoBuckets
}

func newExpoHistogramPointCounters[N int64 | float64](maxSize int, maxScale int32) expoHistogramPointCounters[N] {
	return expoHistogramPointCounters[N]{
		posBuckets: expoBuckets{
			scale:  maxScale,
			counts: make([]atomic.Uint64, maxSize),
		},
		negBuckets: expoBuckets{
			scale:  maxScale,
			counts: make([]atomic.Uint64, maxSize),
		},
	}
}

func (e *expoHistogramPointCounters[N]) reset(maxScale int32) {
	e.sum.reset()
	e.zeroCount.Store(0)
	e.posBuckets.reset(maxScale)
	e.negBuckets.reset(maxScale)
}

func (e *expoHistogramPointCounters[N]) count() uint64 {
	count := e.zeroCount.Load()
	for i := range e.posBuckets.counts {
		count += e.posBuckets.counts[i].Load()
	}
	for i := range e.negBuckets.counts {
		count += e.negBuckets.counts[i].Load()
	}
	return count
}

// mergeInto merges this set of histogram counter data into another,
// and resets the state of this set of counters. This is used by
// hotColdHistogramPoint to ensure that the cumulative counters continue to
// accumulate after being read.
func (p *expoHistogramPointCounters[N]) mergeInto( // nolint:revive // Intentional internal control flag
	into *expoHistogramPointCounters[N],
	noMinMax, noSum bool,
) {
	if !noMinMax {
		// Do not reset min or max because cumulative min and max only ever grow
		// smaller or larger respectively.
		if p.minMax.set.Load() {
			into.minMax.Update(p.minMax.minimum.Load())
			into.minMax.Update(p.minMax.maximum.Load())
		}
	}
	if !noSum {
		into.sum.add(p.sum.load())
	}
	p.posBuckets.mergeInto(&into.posBuckets)
	p.negBuckets.mergeInto(&into.negBuckets)
}

// recordCount adds a new measurement to the histogram. It will rescale the buckets if needed.
func (p *expoHistogramPointCounters[N]) recordCount(v N) {
	absV := math.Abs(float64(v))

	if float64(absV) == 0.0 {
		p.zeroCount.Add(1)
		return
	}

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	bucket.record(absV)
}

// expoBuckets is a set of buckets in an exponential histogram.
type expoBuckets struct {
	scale       int32
	startAndEnd atomicLimitedRange
	counts      []atomic.Uint64
}

// getIdx returns the index into counts for the provided bin.
func (e *expoBuckets) getIdx(bin int32) int {
	newBin := int(bin) % len(e.counts)
	return (newBin + len(e.counts)) % len(e.counts)
}

func (e *expoBuckets) reset(maxScale int32) {
	e.scale = maxScale
	e.startAndEnd.Store(0, 0)
	for i := range e.counts {
		e.counts[i].Store(0)
	}
}

func (e *expoBuckets) loadCountsInto(into *[]uint64) {
	// TODO (#3047): Making copies for bounds and counts incurs a large
	// memory allocation footprint. Alternatives should be explored.
	start, end := e.startAndEnd.Load()
	length := int(end - start)
	counts := reset(*into, length, length)
	eIdx := start
	for i := range length {
		counts[i] = e.counts[e.getIdx(eIdx)].Load()
		eIdx++
	}
	*into = counts
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

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(absV float64) {
	bin := b.getBin(absV)
	fmt.Printf("record() Bin %v for val %v, with scale %v\n", bin, absV, b.scale)

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	if scaleDelta := b.scaleChange(bin); scaleDelta > 0 {
		if b.scale-scaleDelta < expoMinScale {
			// With a scale of -10 there is only two buckets for the whole range of float64 values.
			// This can only happen if there is a max size of 1.
			otel.Handle(errors.New("exponential histogram scale underflow"))
			return
		}
		// Downscale
		b.downscale(scaleDelta)

		bin = b.getBin(absV)
		fmt.Printf("record() After downscale Bin %v for val %v\n", bin, absV)
	} else {
		fmt.Printf("scale change was 0\n")
	}
	b.recordBucket(bin)
}

func (b *expoBuckets) recordBucket(bin int32) {
	b.resizeToInclude(bin)
	b.counts[b.getIdx(bin)].Add(1)
	startBin, endBin := b.startAndEnd.Load()
	fmt.Printf("incremented bin %v idx %v, startBin %v, endBin %v, counts %+v\n", bin, b.getIdx(bin), startBin, endBin, b.counts)
}

func (b *expoBuckets) validate() {
	startBin, endBin := b.startAndEnd.Load()
	if endBin-startBin > int32(len(b.counts)) {
		fmt.Printf("inconsistent start %v end %v len %v\n", startBin, endBin, len(b.counts))
	}
}

// downscale shrinks a bucket by a factor of 2*s. It will sum counts into the
// correct lower resolution bucket.
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
		b.validate()
		// Shift all elements left by the change in start position
		startShift := b.getIdx(startBin - newStartBin)
		b.counts = append(b.counts[startShift:], b.counts[:startShift]...)

		// Clear all elements that are outside of our start to end range
		for i := newEndBin; i < newStartBin+int32(len(b.counts)); i++ {
			b.counts[b.getIdx(i)].Store(0)
		}
		fmt.Printf("downscale SHORT downsize %v startBin %v, endBin %v, counts %v\n", delta, newStartBin, newEndBin, b.counts)
		return
	}

	steps := int32(1) << delta
	offset := startBin % steps
	offset = (offset + steps) % steps // to make offset positive
	newLen := (length-1+offset)/steps + 1
	newStartBin := startBin >> delta
	newEndBin := newStartBin + newLen
	startShift := b.getIdx(startBin - newStartBin)

	fmt.Printf("downscale START downsize %v startBin %v, newStartBin %v, shift %v, endBin %v, newEndBin %v, steps %v, offset %v, counts %v\n", delta, startBin, newStartBin, startShift, endBin, newEndBin, steps, offset, b.counts)
	for i := startBin + 1; i < endBin; i++ {
		newIdx := b.getIdx(int32(math.Floor(float64(i)/float64(steps))) + int32(startShift))
		if i%steps == 0 {
			b.counts[newIdx].Store(b.counts[b.getIdx(i)].Load())
			fmt.Printf("downscale SET idx %v counts %+v\n", b.getIdx(i/steps+int32(startShift)), b.counts)
			continue
		}
		b.counts[newIdx].Add(b.counts[b.getIdx(i)].Load())
		fmt.Printf("downscale ADD i %v i/steps %v idx %v, getIdx %v counts %+v\n", i, i/steps, i/steps+int32(startShift), b.getIdx(i/steps+int32(startShift)), b.counts)
	}
	fmt.Printf("downscale END Merge downsize %v startBin %v, endBin %v, counts %v\n", delta, startBin, endBin, b.counts)

	fmt.Printf("downscale LONG downsize %v startBin %v, endBin %v, counts %v\n", delta, newStartBin, newEndBin, b.counts)
	b.startAndEnd.Store(newStartBin, newEndBin)
	b.validate()
	// Shift all elements left by the change in start position
	fmt.Printf("downscale SHIFT startBin %v newStartBin %v shift %v\n", startBin, newStartBin, startShift)
	b.counts = append(b.counts[startShift:], b.counts[:startShift]...)

	fmt.Printf("downscale AfterShiFT downsize %v startBin %v, endBin %v, counts %v\n", delta, newStartBin, newEndBin, b.counts)
	// Clear all elements that are outside of our start to end range
	for i := newEndBin; i < newStartBin+int32(len(b.counts)); i++ {
		b.counts[b.getIdx(i)].Store(0)
	}
	fmt.Printf("downscale FINAL downsize %v startBin %v, endBin %v, counts %v\n", delta, newStartBin, newEndBin, b.counts)
}

func (b *expoBuckets) resizeToInclude(bin int32) {
	startBin, endBin := b.startAndEnd.Load()
	if startBin == endBin {
		startBin = bin
		endBin = bin + 1
		b.startAndEnd.Store(startBin, endBin)
		b.validate()
		fmt.Printf("resizeToInclude AAAAA bin %v start %v , end %v\n", bin, startBin, endBin)
	} else if bin < startBin {
		b.startAndEnd.Store(bin, endBin)
		b.validate()
		fmt.Printf("resizeToInclude CCCCC bin %v start %v , end %v\n", bin, startBin, endBin)
	} else if bin >= endBin {
		b.startAndEnd.Store(startBin, bin+1)
		b.validate()
		fmt.Printf("resizeToInclude DDDDDD bin %v start %v , end %v\n", bin, startBin, endBin)
	} else {
		fmt.Printf("resizeToInclude BBBBBB bin %v start %v , end %v\n", bin, startBin, endBin)
	}
}

// mergeInto merges this expoBuckets into another, and resets the state
// of the expoBuckets. This is used to ensure that the cumulative counters
// continue to accumulate after being read. It returns the scale change that
// was applied to the input buckets.
func (b *expoBuckets) mergeInto(into *expoBuckets) {
	// Rescale both to the same scale
	scaleDelta := into.scale - b.scale
	if scaleDelta > 0 {
		into.downscale(scaleDelta)
	} else if scaleDelta < 0 {
		b.downscale(-scaleDelta)
	}

	startBin, endBin := b.startAndEnd.Load()
	into.resizeToInclude(startBin)
	into.resizeToInclude(endBin - 1)
	scaleDelta = into.scaleChange(endBin - 1)
	if scaleDelta > 0 {
		// Merging buckets required a scale change to the positive buckets to
		// fit within the max scale. Update scale and scale down the negative
		// buckets to match.
		b.downscale(scaleDelta)
		into.downscale(scaleDelta)
	}
	// At this point, into as been expanded to be a superset of b.
	// Now we finally increment buckets.
	bStartBin, _ := b.startAndEnd.Load()
	intoStartBin, _ := b.startAndEnd.Load()
	startBinDelta := bStartBin - intoStartBin
	for i := range b.counts {
		into.counts[i+int(startBinDelta)].Add(b.counts[i].Load())
	}
	b.validate()
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
		hPt := newHotColdExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)
		hPt.res = e.newRes(attr)
		return hPt
	}).(*hotColdExpoHistogramPoint[N])
	ptHotIdx := v.hcwg.start()
	defer v.hcwg.done(ptHotIdx)
	if !v.noMinMax {
		v.hotColdPoint[ptHotIdx].minMax.Update(value)
	}
	if !v.noSum {
		v.hotColdPoint[ptHotIdx].sum.add(value)
	}
	v.rescaleMux.Lock()
	v.hotColdPoint[ptHotIdx].recordCount(value)
	v.rescaleMux.Unlock()
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
		val := value.(*hotColdExpoHistogramPoint[N])
		val.rescaleMux.Lock()
		defer val.rescaleMux.Unlock()
		readIdx := val.hcwg.swapHotAndWait()
		defer val.hotColdPoint[readIdx].reset(e.maxScale)
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = e.start
		hDPts[i].Time = t
		hDPts[i].Count = val.hotColdPoint[readIdx].count()
		hDPts[i].ZeroCount = val.hotColdPoint[readIdx].zeroCount.Load()
		hDPts[i].ZeroThreshold = 0.0

		// Unify the positive and negative scales by downscaling the higher
		// scale to the lower one.
		scale := min(val.hotColdPoint[readIdx].posBuckets.scale, val.hotColdPoint[readIdx].negBuckets.scale)
		hDPts[i].Scale = scale
		if scaleDelta := val.hotColdPoint[readIdx].posBuckets.scale - scale; scaleDelta > 0 {
			val.hotColdPoint[readIdx].posBuckets.downscale(scaleDelta)
		}
		if scaleDelta := val.hotColdPoint[readIdx].negBuckets.scale - scale; scaleDelta > 0 {
			val.hotColdPoint[readIdx].negBuckets.downscale(scaleDelta)
		}

		offset, _ := val.hotColdPoint[readIdx].posBuckets.startAndEnd.Load()
		hDPts[i].PositiveBucket.Offset = offset
		val.hotColdPoint[readIdx].posBuckets.loadCountsInto(&hDPts[i].PositiveBucket.Counts)

		offset, _ = val.hotColdPoint[readIdx].negBuckets.startAndEnd.Load()
		hDPts[i].NegativeBucket.Offset = offset
		val.hotColdPoint[readIdx].negBuckets.loadCountsInto(&hDPts[i].NegativeBucket.Counts)

		if !e.noSum {
			hDPts[i].Sum = val.hotColdPoint[readIdx].sum.load()
		}
		if !e.noMinMax {
			if val.hotColdPoint[readIdx].minMax.set.Load() {
				hDPts[i].Min = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.minimum.Load())
				hDPts[i].Max = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.maximum.Load())
			}
		}

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
		hPt := newHotColdExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)
		hPt.res = e.newRes(attr)
		return hPt
	}).(*hotColdExpoHistogramPoint[N])
	hotIdx := v.hcwg.start()
	defer v.hcwg.done(hotIdx)
	if !v.noMinMax {
		v.hotColdPoint[hotIdx].minMax.Update(value)
	}
	if !v.noSum {
		v.hotColdPoint[hotIdx].sum.add(value)
	}
	v.rescaleMux.Lock()
	v.hotColdPoint[hotIdx].recordCount(value)
	v.rescaleMux.Unlock()
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
		val.rescaleMux.Lock()
		defer val.rescaleMux.Unlock()
		readIdx := val.hcwg.swapHotAndWait()
		newPt := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    val.attrs,
			StartTime:     e.start,
			Time:          t,
			Count:         val.hotColdPoint[readIdx].count(),
			ZeroCount:     val.hotColdPoint[readIdx].zeroCount.Load(),
			ZeroThreshold: 0.0,
		}

		// Unify the positive and negative scales by downscaling the higher
		// scale to the lower one.
		scale := min(val.hotColdPoint[readIdx].posBuckets.scale, val.hotColdPoint[readIdx].negBuckets.scale)
		newPt.Scale = scale
		if scaleDelta := val.hotColdPoint[readIdx].posBuckets.scale - scale; scaleDelta > 0 {
			val.hotColdPoint[readIdx].posBuckets.downscale(scaleDelta)
		}
		if scaleDelta := val.hotColdPoint[readIdx].negBuckets.scale - scale; scaleDelta > 0 {
			val.hotColdPoint[readIdx].negBuckets.downscale(scaleDelta)
		}

		offset, _ := val.hotColdPoint[readIdx].posBuckets.startAndEnd.Load()
		newPt.PositiveBucket.Offset = offset
		val.hotColdPoint[readIdx].posBuckets.loadCountsInto(&newPt.PositiveBucket.Counts)

		offset, _ = val.hotColdPoint[readIdx].negBuckets.startAndEnd.Load()
		newPt.NegativeBucket.Offset = offset
		val.hotColdPoint[readIdx].negBuckets.loadCountsInto(&newPt.NegativeBucket.Counts)

		if !e.noSum {
			newPt.Sum = val.hotColdPoint[readIdx].sum.load()
		}
		if !e.noMinMax {
			if val.hotColdPoint[readIdx].minMax.set.Load() {
				newPt.Min = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.minimum.Load())
				newPt.Max = metricdata.NewExtrema(val.hotColdPoint[readIdx].minMax.maximum.Load())
			}
		}
		// Once we've read the point, merge it back into the hot histogram
		// point since it is cumulative.
		hotIdx := (readIdx + 1) % 2
		val.hotColdPoint[readIdx].mergeInto(&val.hotColdPoint[hotIdx], e.noMinMax, e.noSum)
		val.hotColdPoint[readIdx].reset(val.maxScale)

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
