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
	"go.opentelemetry.io/otel/sdk/internal/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	expoMaxScale = 20
	expoMinScale = -10

	smallestNonZeroNormalFloat64 = 0x1p-1022
)

var errUnderflow = errors.New("exponential histogram underflow (exceeds maxSize at scale -10)")

// expoHistogramDataPoint is a single data point in an exponential histogram.
type expoHistogramDataPoint[N int64 | float64] struct {
	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	minMax atomicMinMax[N]
	sum    atomicCounter[N]

	maxSize  int
	noMinMax bool
	noSum    bool

	mu sync.Mutex

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
) *expoHistogramDataPoint[N] { // nolint:revive // we need this control flag
	dp := &expoHistogramDataPoint[N]{
		attrs:     attrs,
		maxSize:   maxSize,
		noMinMax:  noMinMax,
		noSum:     noSum,
		startTime: now(),
	}
	dp.scale.Store(maxScale)
	return dp
}

// record adds a new measurement to the histogram. It will rescale the buckets if needed.
func (p *expoHistogramDataPoint[N]) record(v N) {
	if !p.noMinMax {
		p.minMax.Update(v)
	}
	if !p.noSum {
		p.sum.add(v)
	}

	absV := math.Abs(float64(v))

	if float64(absV) == 0.0 {
		p.zeroCount.Add(1)
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	bin := p.getBin(absV)

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	if scaleDelta := p.scaleChange(bin, bucket.startBin, len(bucket.counts)); scaleDelta > 0 {
		currentScale := p.scale.Load()
		if currentScale-scaleDelta < expoMinScale {
			otel.Handle(errUnderflow)
			return
		}
		// Downscale
		p.scale.Add(-scaleDelta)
		p.posBuckets.downscale(scaleDelta)
		p.negBuckets.downscale(scaleDelta)

		bin = p.getBin(absV)
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

// expoBuckets is a set of buckets in an exponential histogram.
type expoBuckets struct {
	startBin int32
	counts   []atomic.Uint64
}

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(bin int32) {
	b.recordCount(bin, 1)
}

func (b *expoBuckets) recordCount(bin int32, count uint64) {
	if len(b.counts) == 0 {
		b.counts = make([]atomic.Uint64, 1)
		b.counts[0].Store(count)
		b.startBin = bin
		return
	}

	endBin := int(b.startBin) + len(b.counts) - 1

	// if the new bin is inside the current range
	if bin >= b.startBin && int(bin) <= endBin {
		b.counts[bin-b.startBin].Add(count)
		return
	}
	// if the new bin is before the current start add spaces to the counts
	if bin < b.startBin {
		origLen := len(b.counts)
		newLength := endBin - int(bin) + 1
		shift := b.startBin - bin

		if newLength > cap(b.counts) {
			b.counts = append(b.counts, make([]atomic.Uint64, newLength-len(b.counts))...)
		}

		b.counts = b.counts[:newLength]

		// Shift existing elements to the right. Go's copy() doesn't work for
		// structs like atomic.Uint64.
		for i := origLen - 1; i >= 0; i-- {
			b.counts[i+int(shift)].Store(b.counts[i].Load())
		}

		for i := 1; i < int(shift); i++ {
			b.counts[i].Store(0)
		}
		b.startBin = bin
		b.counts[0].Store(count)
		return
	}
	// if the new is after the end add spaces to the end
	if int(bin) > endBin {
		if int(bin-b.startBin) < cap(b.counts) {
			b.counts = b.counts[:bin-b.startBin+1]
			for i := endBin + 1 - int(b.startBin); i < len(b.counts); i++ {
				b.counts[i].Store(0)
			}
			b.counts[bin-b.startBin].Store(count)
			return
		}

		end := make([]atomic.Uint64, int(bin-b.startBin)-len(b.counts)+1)
		b.counts = append(b.counts, end...)
		b.counts[bin-b.startBin].Store(count)
	}
}

// downscale shrinks a bucket by a factor of 2*s. It will sum counts into the
// correct lower resolution bucket.
func (b *expoBuckets) downscale(delta int32) {
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
			b.counts[idx/int(steps)].Store(b.counts[i].Load())
			continue
		}
		b.counts[idx/int(steps)].Add(b.counts[i].Load())
	}

	lastIdx := (len(b.counts) - 1 + int(offset)) / int(steps)
	b.counts = b.counts[:lastIdx+1]
	b.startBin >>= delta
}

func (b *expoBuckets) merge(other *expoBuckets) {
	if len(other.counts) == 0 {
		return
	}
	for i := range other.counts {
		c := other.counts[i].Load()
		if c > 0 {
			b.recordCount(other.startBin+int32(i), c)
		}
	}
}

func (b *expoBuckets) count() uint64 {
	var total uint64
	for i := range b.counts {
		total += b.counts[i].Load()
	}
	return total
}

func (p *expoHistogramDataPoint[N]) merge(other *expoHistogramDataPoint[N]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	pStartScale := p.scale.Load()
	oStartScale := other.scale.Load()

	targetScale := min(pStartScale, oStartScale)
	pAlignShift := pStartScale - targetScale
	oAlignShift := oStartScale - targetScale

	var d int32
	if len(p.posBuckets.counts) > 0 && len(other.posBuckets.counts) > 0 {
		pMinBin := p.posBuckets.startBin >> pAlignShift
		oMinBin := other.posBuckets.startBin >> oAlignShift
		pMaxBin := (p.posBuckets.startBin + int32(len(p.posBuckets.counts)) - 1) >> pAlignShift         // nolint: gosec // length fits in int32
		oMaxBin := (other.posBuckets.startBin + int32(len(other.posBuckets.counts)) - 1) >> oAlignShift // nolint: gosec // length fits in int32

		minBin := min(pMinBin, oMinBin)
		maxBin := max(pMaxBin, oMaxBin)
		delta := p.scaleChange(maxBin, minBin, 1)
		if delta > d {
			d = delta
		}
	}
	if len(p.negBuckets.counts) > 0 && len(other.negBuckets.counts) > 0 {
		pMinBin := p.negBuckets.startBin >> pAlignShift
		oMinBin := other.negBuckets.startBin >> oAlignShift
		pMaxBin := (p.negBuckets.startBin + int32(len(p.negBuckets.counts)) - 1) >> pAlignShift         // nolint: gosec // length fits in int32
		oMaxBin := (other.negBuckets.startBin + int32(len(other.negBuckets.counts)) - 1) >> oAlignShift // nolint: gosec // length fits in int32

		minBin := min(pMinBin, oMinBin)
		maxBin := max(pMaxBin, oMaxBin)
		delta := p.scaleChange(maxBin, minBin, 1)
		if delta > d {
			d = delta
		}
	}

	pDownscale := pAlignShift + d
	if pDownscale > 0 {
		p.scale.Add(-pDownscale)
		p.posBuckets.downscale(pDownscale)
		p.negBuckets.downscale(pDownscale)
	}

	oDownscale := oAlignShift + d
	if oDownscale > 0 {
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
		hotColdValMap: [2]limitedSyncMap{
			{aggLimit: limit},
			{aggLimit: limit},
		},

		start: now(),
	}
}

func newCumulativeExpoHistogram[N int64 | float64](
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
		vp := newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)
		vp.res = e.newRes(attr)
		return vp
	}).(*expoHistogramDataPoint[N])

	v.record(value)
	v.res.Offer(ctx, value, droppedAttr)
}

func (e *deltaExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, // nolint:gocritic // dest is an interface pointer used to avoid allocations
) int {
	t := now()

	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.DeltaTemporality

	coldIdx := e.hcwg.swapHotAndWait()

	n := e.hotColdValMap[coldIdx].Len()
	hDPts := reset(h.DataPoints, 0, n)

	var i int
	e.hotColdValMap[coldIdx].Range(func(_, value any) bool {
		val := value.(*expoHistogramDataPoint[N])

		dPt := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    val.attrs,
			StartTime:     e.start,
			Time:          t,
			Count:         val.count(),
			Scale:         val.scale.Load(),
			ZeroCount:     val.zeroCount.Load(),
			ZeroThreshold: 0.0,
		}

		dPt.PositiveBucket.Offset = val.posBuckets.startBin
		if i < len(h.DataPoints) && len(h.DataPoints[i].PositiveBucket.Counts) >= len(val.posBuckets.counts) {
			dPt.PositiveBucket.Counts = reset(
				h.DataPoints[i].PositiveBucket.Counts,
				len(val.posBuckets.counts),
				len(val.posBuckets.counts),
			)
		} else {
			dPt.PositiveBucket.Counts = make([]uint64, len(val.posBuckets.counts))
		}
		for j := range val.posBuckets.counts {
			dPt.PositiveBucket.Counts[j] = val.posBuckets.counts[j].Load()
		}

		dPt.NegativeBucket.Offset = val.negBuckets.startBin
		if i < len(h.DataPoints) && len(h.DataPoints[i].NegativeBucket.Counts) >= len(val.negBuckets.counts) {
			dPt.NegativeBucket.Counts = reset(
				h.DataPoints[i].NegativeBucket.Counts,
				len(val.negBuckets.counts),
				len(val.negBuckets.counts),
			)
		} else {
			dPt.NegativeBucket.Counts = make([]uint64, len(val.negBuckets.counts))
		}
		for j := range val.negBuckets.counts {
			dPt.NegativeBucket.Counts[j] = val.negBuckets.counts[j].Load()
		}

		if !e.noSum {
			dPt.Sum = val.sum.load()
		}
		if !e.noMinMax {
			if val.minMax.set.Load() {
				dPt.Min = metricdata.NewExtrema(val.minMax.minimum.Load())
				dPt.Max = metricdata.NewExtrema(val.minMax.maximum.Load())
			}
		}

		collectExemplars(&dPt.Exemplars, val.res.Collect)

		hDPts = append(hDPts, dPt)
		i++
		return true
	})

	e.start = t
	e.hotColdValMap[coldIdx].Clear()

	h.DataPoints = hDPts
	*dest = h
	return n
}

type cumulativePoint[N int64 | float64] struct {
	wg         hotColdWaitGroup
	points     [2]*expoHistogramDataPoint[N]
	cumulative *expoHistogramDataPoint[N]

	tracker atomicUnderflowTracker
}

// cumulativeExpoHistogram summarizes a set of measurements as an histogram with exponentially
// defined buckets.
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

	vp := e.values.LoadOrStoreAttr(fltrAttr, func(attr attribute.Set) any {
		cp := &cumulativePoint[N]{}
		cp.tracker.maxSize = e.maxSize
		cp.cumulative = newExpoHistogramDataPoint[N](attr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)

		cp.points[0] = &expoHistogramDataPoint[N]{
			attrs:    attr,
			maxSize:  e.maxSize,
			noMinMax: e.noMinMax,
			noSum:    e.noSum,
		}
		cp.points[0].scale.Store(e.maxScale)
		cp.points[0].res = e.newRes(attr)

		cp.points[1] = &expoHistogramDataPoint[N]{
			attrs:    attr,
			maxSize:  e.maxSize,
			noMinMax: e.noMinMax,
			noSum:    e.noSum,
		}
		cp.points[1].scale.Store(e.maxScale)
		cp.points[1].res = e.newRes(attr)

		return cp
	}).(*cumulativePoint[N])

	if !vp.tracker.checkAndRecord(float64(value)) {
		otel.Handle(errUnderflow)
		return
	}

	hotIdx := vp.wg.start()
	v := vp.points[hotIdx]
	v.record(value)
	v.res.Offer(ctx, value, droppedAttr)
	vp.wg.done(hotIdx)
}

func (e *cumulativeExpoHistogram[N]) collect(
	dest *metricdata.Aggregation, // nolint:gocritic // dest is an interface pointer used to avoid allocations
) int {
	t := now()

	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	n := e.values.Len()
	hDPts := reset(h.DataPoints, 0, n)

	perSeriesStartTimeEnabled := x.PerSeriesStartTimestamps.Enabled()

	var i int
	e.values.Range(func(_, value any) bool {
		cp := value.(*cumulativePoint[N])

		coldIdx := cp.wg.swapHotAndWait()
		delta := cp.points[coldIdx]

		cp.cumulative.merge(delta)

		// Replace the cold delta point and its reservoir
		cp.points[coldIdx] = &expoHistogramDataPoint[N]{
			attrs:    delta.attrs,
			maxSize:  e.maxSize,
			noMinMax: e.noMinMax,
			noSum:    e.noSum,
		}
		cp.points[coldIdx].scale.Store(e.maxScale)
		cp.points[coldIdx].res = e.newRes(delta.attrs)

		val := cp.cumulative
		val.mu.Lock()

		startTime := e.start
		if perSeriesStartTimeEnabled {
			startTime = val.startTime
		}

		dPt := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    val.attrs,
			StartTime:     startTime,
			Time:          t,
			Count:         val.count(),
			Scale:         val.scale.Load(),
			ZeroCount:     val.zeroCount.Load(),
			ZeroThreshold: 0.0,
		}

		dPt.PositiveBucket.Offset = val.posBuckets.startBin
		if i < len(h.DataPoints) && len(h.DataPoints[i].PositiveBucket.Counts) >= len(val.posBuckets.counts) {
			dPt.PositiveBucket.Counts = reset(
				h.DataPoints[i].PositiveBucket.Counts,
				len(val.posBuckets.counts),
				len(val.posBuckets.counts),
			)
		} else {
			dPt.PositiveBucket.Counts = make([]uint64, len(val.posBuckets.counts))
		}
		for j := range val.posBuckets.counts {
			dPt.PositiveBucket.Counts[j] = val.posBuckets.counts[j].Load()
		}

		dPt.NegativeBucket.Offset = val.negBuckets.startBin
		if i < len(h.DataPoints) && len(h.DataPoints[i].NegativeBucket.Counts) >= len(val.negBuckets.counts) {
			dPt.NegativeBucket.Counts = reset(
				h.DataPoints[i].NegativeBucket.Counts,
				len(val.negBuckets.counts),
				len(val.negBuckets.counts),
			)
		} else {
			dPt.NegativeBucket.Counts = make([]uint64, len(val.negBuckets.counts))
		}
		for j := range val.negBuckets.counts {
			dPt.NegativeBucket.Counts[j] = val.negBuckets.counts[j].Load()
		}

		if !e.noSum {
			dPt.Sum = val.sum.load()
		}
		if !e.noMinMax {
			if val.minMax.set.Load() {
				dPt.Min = metricdata.NewExtrema(val.minMax.minimum.Load())
				dPt.Max = metricdata.NewExtrema(val.minMax.maximum.Load())
			}
		}

		if delta.res != nil {
			// Extract exemplars from the delta collector since they represent the current reporting cycle
			collectExemplars(&dPt.Exemplars, delta.res.Collect)
		}

		val.mu.Unlock()

		hDPts = append(hDPts, dPt)
		i++
		return true
	})

	h.DataPoints = hDPts
	*dest = h
	return n
}
