// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"errors"
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
	expoHistogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	noMinMax bool
	noSum    bool
}

func newExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
	noMinMax, noSum bool,
) *expoHistogramDataPoint[N] { // nolint:revive // we need this control flag
	return &expoHistogramDataPoint[N]{
		attrs:    attrs,
		noMinMax: noMinMax,
		noSum:    noSum,
		expoHistogramPointCounters: expoHistogramPointCounters[N]{
			posBuckets: expoBuckets{
				scale:   maxScale,
				maxSize: maxSize,
			},
			negBuckets: expoBuckets{
				scale:   maxScale,
				maxSize: maxSize,
			},
		},
	}
}

// hotColdExpoHistogramPoint a hot and cold exponential histogram points, used
// in cumulative aggregations.
type hotColdExpoHistogramPoint[N int64 | float64] struct {
	hcwg         hotColdWaitGroup
	hotColdPoint [2]expoHistogramPointCounters[N]

	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	noMinMax bool
	noSum    bool
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
		hotColdPoint: [2]expoHistogramPointCounters[N]{
			{
				posBuckets: expoBuckets{
					scale:   maxScale,
					maxSize: maxSize,
				},
				negBuckets: expoBuckets{
					scale:   maxScale,
					maxSize: maxSize,
				},
			},
			{
				posBuckets: expoBuckets{
					scale:   maxScale,
					maxSize: maxSize,
				},
				negBuckets: expoBuckets{
					scale:   maxScale,
					maxSize: maxSize,
				},
			},
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

func (e *expoHistogramPointCounters[N]) count() uint64 {
	count := e.zeroCount.Load()
	for i := range e.posBuckets.counts {
		count += e.posBuckets.counts[i]
	}
	for i := range e.negBuckets.counts {
		count += e.negBuckets.counts[i]
	}
	return count
}

// mergeIntoAndReset merges this set of histogram counter data into another,
// and resets the state of this set of counters. This is used by
// hotColdHistogramPoint to ensure that the cumulative counters continue to
// accumulate after being read.
func (p *expoHistogramPointCounters[N]) mergeIntoAndReset( // nolint:revive // Intentional internal control flag
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
		p.sum.reset()
	}
	p.posBuckets.mergeIntoAndReset(&into.posBuckets)
	p.negBuckets.mergeIntoAndReset(&into.negBuckets)
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
func scaleChange(bin, startBin int32, length, maxSize int) int32 {
	if length == 0 {
		// No need to rescale if there are no buckets.
		return 0
	}

	low := int(startBin)
	high := int(bin)
	if startBin >= bin {
		low = int(bin)
		high = int(startBin) + length - 1
	}

	var count int32
	for high-low >= maxSize {
		low >>= 1
		high >>= 1
		count++
		if count > expoMaxScale-expoMinScale {
			return count
		}
	}
	return count
}

// expoBuckets is a set of buckets in an exponential histogram.
type expoBuckets struct {
	scaleMux sync.Mutex
	scale    int32
	startBin int32
	counts   []uint64
	maxSize  int
}

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(absV float64) {
	b.scaleMux.Lock()
	defer b.scaleMux.Unlock()
	bin := b.getBin(absV)

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	if scaleDelta := scaleChange(bin, b.startBin, len(b.counts), b.maxSize); scaleDelta > 0 {
		if b.scale-scaleDelta < expoMinScale {
			// With a scale of -10 there is only two buckets for the whole range of float64 values.
			// This can only happen if there is a max size of 1.
			otel.Handle(errors.New("exponential histogram scale underflow"))
			return
		}
		// Downscale
		b.downscale(scaleDelta)

		bin = b.getBin(absV)
	}
	b.recordBucket(bin)
}

func (b *expoBuckets) recordBucket(bin int32) {
	b.resizeToInclude(bin)
	b.counts[bin-b.startBin]++
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

	if len(b.counts) <= 1 || delta < 1 {
		b.startBin >>= delta
		return
	}

	steps := int32(1) << delta
	offset := b.startBin % steps
	offset = (offset + steps) % steps // to make offset positive
	for i := 1; i < len(b.counts); i++ {
		idx := i + int(offset)
		if idx%int(steps) == 0 {
			b.counts[idx/int(steps)] = b.counts[i]
			continue
		}
		b.counts[idx/int(steps)] += b.counts[i]
	}

	lastIdx := (len(b.counts) - 1 + int(offset)) / int(steps)
	b.counts = b.counts[:lastIdx+1]
	b.startBin >>= delta
}

func (b *expoBuckets) endBin() int32 {
	return b.startBin + int32(len(b.counts)-1)
}

func (b *expoBuckets) resizeToInclude(bin int32) {
	if len(b.counts) == 0 {
		b.counts = []uint64{0}
		b.startBin = bin
		return
	}
	endBin := b.endBin()

	// if the new bin is inside the current range
	if bin >= b.startBin && bin <= endBin {
		return
	}
	// if the new bin is before the current start add spaces to the counts
	if bin < b.startBin {
		origLen := len(b.counts)
		newLength := int(endBin - bin + 1)
		shift := b.startBin - bin

		if newLength > cap(b.counts) {
			b.counts = append(b.counts, make([]uint64, newLength-len(b.counts))...)
		}

		copy(b.counts[shift:origLen+int(shift)], b.counts)
		b.counts = b.counts[:newLength]
		for i := 0; i < int(shift); i++ {
			b.counts[i] = 0
		}
		b.startBin = bin
		return
	}
	// if the new is after the end add spaces to the end
	if bin > endBin {
		if int(bin-b.startBin) < cap(b.counts) {
			b.counts = b.counts[:bin-b.startBin+1]
			for i := int(endBin + 1 - b.startBin); i < len(b.counts); i++ {
				b.counts[i] = 0
			}
			return
		}

		end := make([]uint64, int(bin-b.startBin)-len(b.counts)+1)
		b.counts = append(b.counts, end...)
	}
}

// mergeIntoAndReset merges this expoBuckets into another, and resets the state
// of the expoBuckets. This is used to ensure that the cumulative counters
// continue to accumulate after being read. It returns the scale change that
// was applied to the input buckets.
func (b *expoBuckets) mergeIntoAndReset(into *expoBuckets) {
	// we alraedy hold the lock for b.
	into.scaleMux.Lock()
	defer into.scaleMux.Unlock()
	// Rescale both to the same scale
	scaleDelta := into.scale - b.scale
	if scaleDelta > 0 {
		into.downscale(scaleDelta)
	} else if scaleDelta < 0 {
		b.downscale(-scaleDelta)
	}

	into.resizeToInclude(b.startBin)
	into.resizeToInclude(b.endBin())
	scaleDelta = scaleChange(into.endBin(), into.startBin, len(into.counts), b.maxSize)
	if scaleDelta > 0 {
		// Merging buckets required a scale change to the positive buckets to
		// fit within the max scale. Update scale and scale down the negative
		// buckets to match.
		b.downscale(scaleDelta)
		into.downscale(scaleDelta)
	}
	// At this point, into as been expanded to be a superset of b.
	// Now we finally increment buckets.
	startBinDelta := b.startBin - into.startBin
	for i := range b.counts {
		into.counts[i+int(startBinDelta)] += b.counts[i]
		b.counts[i] = 0
	}
	b.counts = b.counts[0:0]
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
		hPt := newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)
		hPt.res = e.newRes(attr)
		return hPt
	}).(*expoHistogramDataPoint[N])
	if !v.noMinMax {
		v.minMax.Update(value)
	}
	if !v.noSum {
		v.sum.add(value)
	}
	v.recordCount(value)
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
		hDPts[i].Count = val.count()
		hDPts[i].ZeroCount = val.zeroCount.Load()
		hDPts[i].ZeroThreshold = 0.0

		val.posBuckets.scaleMux.Lock()
		defer val.posBuckets.scaleMux.Unlock()
		val.negBuckets.scaleMux.Lock()
		defer val.negBuckets.scaleMux.Unlock()
		// Unify the positive and negative scales by downscaling the higher
		// scale to the lower one.
		scale := min(val.posBuckets.scale, val.negBuckets.scale)
		hDPts[i].Scale = scale
		if scaleDelta := val.posBuckets.scale - scale; scaleDelta > 0 {
			val.posBuckets.downscale(scaleDelta)
		}
		if scaleDelta := val.negBuckets.scale - scale; scaleDelta > 0 {
			val.negBuckets.downscale(scaleDelta)
		}

		hDPts[i].PositiveBucket.Offset = val.posBuckets.startBin
		hDPts[i].PositiveBucket.Counts = reset(
			hDPts[i].PositiveBucket.Counts,
			len(val.posBuckets.counts),
			len(val.posBuckets.counts),
		)
		copy(hDPts[i].PositiveBucket.Counts, val.posBuckets.counts)

		hDPts[i].NegativeBucket.Offset = val.negBuckets.startBin
		hDPts[i].NegativeBucket.Counts = reset(
			hDPts[i].NegativeBucket.Counts,
			len(val.negBuckets.counts),
			len(val.negBuckets.counts),
		)
		copy(hDPts[i].NegativeBucket.Counts, val.negBuckets.counts)

		if !e.noSum {
			hDPts[i].Sum = val.sum.load()
		}
		if !e.noMinMax {
			if val.minMax.set.Load() {
				hDPts[i].Min = metricdata.NewExtrema(val.minMax.minimum.Load())
				hDPts[i].Max = metricdata.NewExtrema(val.minMax.maximum.Load())
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
	v.hotColdPoint[hotIdx].recordCount(value)
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
		readIdx := val.hcwg.swapHotAndWait()
		newPt := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    val.attrs,
			StartTime:     e.start,
			Time:          t,
			Count:         val.hotColdPoint[readIdx].count(),
			ZeroCount:     val.hotColdPoint[readIdx].zeroCount.Load(),
			ZeroThreshold: 0.0,
		}

		val.hotColdPoint[readIdx].posBuckets.scaleMux.Lock()
		defer val.hotColdPoint[readIdx].posBuckets.scaleMux.Unlock()
		val.hotColdPoint[readIdx].negBuckets.scaleMux.Lock()
		defer val.hotColdPoint[readIdx].negBuckets.scaleMux.Unlock()
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

		newPt.PositiveBucket.Offset = val.hotColdPoint[readIdx].posBuckets.startBin
		newPt.PositiveBucket.Counts = reset(
			newPt.PositiveBucket.Counts,
			len(val.hotColdPoint[readIdx].posBuckets.counts),
			len(val.hotColdPoint[readIdx].posBuckets.counts),
		)
		copy(newPt.PositiveBucket.Counts, val.hotColdPoint[readIdx].posBuckets.counts)

		newPt.NegativeBucket.Offset = val.hotColdPoint[readIdx].negBuckets.startBin
		newPt.NegativeBucket.Counts = reset(
			newPt.NegativeBucket.Counts,
			len(val.hotColdPoint[readIdx].negBuckets.counts),
			len(val.hotColdPoint[readIdx].negBuckets.counts),
		)
		copy(newPt.NegativeBucket.Counts, val.hotColdPoint[readIdx].negBuckets.counts)

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
		val.hotColdPoint[readIdx].mergeIntoAndReset(&val.hotColdPoint[hotIdx], e.noMinMax, e.noSum)

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
