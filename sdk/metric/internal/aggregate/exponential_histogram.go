// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	expoMaxScale = 20
	expoMinScale = -10

	// These redefine the Math constants with a type, so the compiler won't coerce
	// them into an int on 32 bit platforms.
	maxInt64 int64 = math.MaxInt64
	minInt64 int64 = math.MinInt64
)

// expoHistogramDataPoint is a single data point in an exponential histogram.
type expoHistogramDataPoint[N int64 | float64] struct {
	attrs attribute.Set
	res   FilteredExemplarReservoir[N]

	min N
	max N
	sum N

	maxSize  int
	noMinMax bool
	noSum    bool

	scale int32

	posBuckets expoBuckets
	negBuckets expoBuckets
	zeroCount  uint64
}

func newExpoHistogramDataPoint[N int64 | float64](
	attrs attribute.Set,
	maxSize int,
	maxScale int32,
	noMinMax, noSum bool,
) *expoHistogramDataPoint[N] { // nolint:revive // we need this control flag
	f := math.MaxFloat64
	ma := N(f) // if N is int64, max will overflow to -9223372036854775808
	mi := N(-f)
	if N(maxInt64) > N(f) {
		ma = N(maxInt64)
		mi = N(minInt64)
	}
	return &expoHistogramDataPoint[N]{
		attrs:      attrs,
		min:        ma,
		max:        mi,
		maxSize:    maxSize,
		noMinMax:   noMinMax,
		noSum:      noSum,
		scale:      maxScale,
		posBuckets: newExpoBuckets(maxSize),
		negBuckets: newExpoBuckets(maxSize),
	}
}

// record adds a new measurement to the histogram. It will rescale the buckets if needed.
func (p *expoHistogramDataPoint[N]) record(v N) {
	if !p.noMinMax {
		if v < p.min {
			p.min = v
		}
		if v > p.max {
			p.max = v
		}
	}
	if !p.noSum {
		p.sum += v
	}

	absV := math.Abs(float64(v))

	if float64(absV) == 0.0 {
		p.zeroCount++
		return
	}

	bin := p.getBin(absV)

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	if scaleDelta := p.scaleChange(bin, bucket.startBin, bucket.endBin); scaleDelta > 0 {
		if p.scale-scaleDelta < expoMinScale {
			// With a scale of -10 there is only two buckets for the whole range of float64 values.
			// This can only happen if there is a max size of 1.
			otel.Handle(errors.New("exponential histogram scale underflow"))
			return
		}
		// Downscale
		p.scale -= scaleDelta
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
func (p *expoHistogramDataPoint[N]) scaleChange(bin, startBin, endBin int32) int32 {
	if startBin == endBin {
		// No need to rescale if there are no buckets.
		return 0
	}

	low := int(startBin)
	high := int(bin)
	if startBin >= bin {
		low = int(bin)
		high = int(endBin) - 1
	}

	var count int32
	for high-low >= p.maxSize {
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
	return p.posBuckets.count() + p.negBuckets.count() + p.zeroCount
}

func newExpoBuckets(maxSize int) expoBuckets {
	return expoBuckets{
		counts: make([]uint64, maxSize),
	}
}

// expoBuckets is a set of buckets in an exponential histogram.
type expoBuckets struct {
	// startBin is the inclusive start of the range.
	startBin int32
	// endBin is the exclusive end of the range.
	endBin int32
	// counts is a circular slice of size maxSize, where bin i is stored at
	// position i % maxSize.
	counts []uint64
}

func (e *expoBuckets) getIdx(bin int) int {
	newBin := bin % len(e.counts)
	// ensure the index is positive.
	return (newBin + len(e.counts)) % len(e.counts)
}

// loadCountsAndOffset returns the buckets counts, the count, and the offset.
// It is not safe to call concurrently.
func (e *expoBuckets) loadCountsAndOffset(into *[]uint64) int32 {
	// TODO (#3047): Making copies for bounds and counts incurs a large
	// memory allocation footprint. Alternatives should be explored.
	length := int(e.endBin - e.startBin)
	counts := reset(*into, length, length)
	count := uint64(0)
	eIdx := int(e.startBin)
	for i := range length {
		val := e.counts[e.getIdx(eIdx)]
		counts[i] = val
		count += val
		eIdx++
	}
	*into = counts
	return e.startBin
}

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(bin int32) {
	// we already downscaled, so we are guaranteed that we can fit within the
	// counts.
	b.counts[b.getIdx(int(bin))]++
	switch {
	case b.startBin == b.endBin:
		b.endBin = bin + 1
		b.startBin = bin
	case bin < b.startBin:
		b.startBin = bin
	case bin >= b.endBin:
		b.endBin = bin + 1
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

	oldStartBin := b.startBin
	oldEndBin := b.endBin
	oldLength := b.endBin - b.startBin
	if oldLength <= 1 || delta < 1 {
		b.startBin >>= delta
		b.endBin += b.startBin - oldStartBin
		// Shift all elements left by the change in start position
		startShift := b.getIdx(int(oldStartBin - b.startBin))
		b.counts = append(b.counts[startShift:], b.counts[:startShift]...)
		return
	}

	steps := int32(1) << delta

	offset := oldStartBin % steps
	offset = (offset + steps) % steps // to make offset positive
	newLen := (oldLength-1+offset)/steps + 1
	b.startBin >>= delta
	b.endBin = b.startBin + newLen
	startShift := b.getIdx(int(oldStartBin - b.startBin))

	for i := int(oldStartBin) + 1; i < int(oldEndBin); i++ {
		newIdx := b.getIdx(int(math.Floor(float64(i)/float64(steps))) + startShift)
		if i%int(steps) == 0 {
			b.counts[newIdx] = b.counts[b.getIdx(i)]
			continue
		}
		b.counts[newIdx] += b.counts[b.getIdx(i)]
	}
	// Shift all elements left by the change in start position
	b.counts = append(b.counts[startShift:], b.counts[:startShift]...)

	// Clear all elements that are outside of our start to end range
	for i := int(b.endBin); i < int(b.startBin)+len(b.counts); i++ {
		b.counts[b.getIdx(i)] = 0
	}
}

func (b *expoBuckets) count() uint64 {
	var total uint64
	for _, count := range b.counts {
		total += count
	}
	return total
}

// newExponentialHistogram returns an Aggregator that summarizes a set of
// measurements as an exponential histogram. Each histogram is scoped by attributes
// and the aggregation cycle the measurements were made in.
func newExponentialHistogram[N int64 | float64](
	maxSize, maxScale int32,
	noMinMax, noSum bool,
	limit int,
	r func(attribute.Set) FilteredExemplarReservoir[N],
) *expoHistogram[N] {
	return &expoHistogram[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		maxSize:  int(maxSize),
		maxScale: maxScale,

		newRes: r,
		limit:  newLimiter[expoHistogramDataPoint[N]](limit),
		values: make(map[attribute.Distinct]*expoHistogramDataPoint[N]),

		start: now(),
	}
}

// expoHistogram summarizes a set of measurements as an histogram with exponentially
// defined buckets.
type expoHistogram[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int32

	newRes   func(attribute.Set) FilteredExemplarReservoir[N]
	limit    limiter[expoHistogramDataPoint[N]]
	values   map[attribute.Distinct]*expoHistogramDataPoint[N]
	valuesMu sync.Mutex

	start time.Time
}

func (e *expoHistogram[N]) measure(
	ctx context.Context,
	value N,
	fltrAttr attribute.Set,
	droppedAttr []attribute.KeyValue,
) {
	// Ignore NaN and infinity.
	if math.IsInf(float64(value), 0) || math.IsNaN(float64(value)) {
		return
	}

	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	v, ok := e.values[fltrAttr.Equivalent()]
	if !ok {
		fltrAttr = e.limit.Attributes(fltrAttr, e.values)
		// If we overflowed, make sure we add to the existing overflow series
		// if it already exists.
		v, ok = e.values[fltrAttr.Equivalent()]
		if !ok {
			v = newExpoHistogramDataPoint[N](fltrAttr, e.maxSize, e.maxScale, e.noMinMax, e.noSum)
			v.res = e.newRes(fltrAttr)

			e.values[fltrAttr.Equivalent()] = v
		}
	}
	v.record(value)
	v.res.Offer(ctx, value, droppedAttr)
}

func (e *expoHistogram[N]) delta(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.DeltaTemporality

	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	n := len(e.values)
	hDPts := reset(h.DataPoints, n, n)

	var i int
	for _, val := range e.values {
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = e.start
		hDPts[i].Time = t
		hDPts[i].Count = val.count()
		hDPts[i].Scale = val.scale
		hDPts[i].ZeroCount = val.zeroCount
		hDPts[i].ZeroThreshold = 0.0

		offset := val.posBuckets.loadCountsAndOffset(&hDPts[i].PositiveBucket.Counts)
		hDPts[i].PositiveBucket.Offset = offset

		offset = val.negBuckets.loadCountsAndOffset(&hDPts[i].NegativeBucket.Counts)
		hDPts[i].NegativeBucket.Offset = offset

		if !e.noSum {
			hDPts[i].Sum = val.sum
		}
		if !e.noMinMax {
			hDPts[i].Min = metricdata.NewExtrema(val.min)
			hDPts[i].Max = metricdata.NewExtrema(val.max)
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
	}
	// Unused attribute sets do not report.
	clear(e.values)

	e.start = t
	h.DataPoints = hDPts
	*dest = h
	return n
}

func (e *expoHistogram[N]) cumulative(
	dest *metricdata.Aggregation, //nolint:gocritic // The pointer is needed for the ComputeAggregation interface
) int {
	t := now()

	// If *dest is not a metricdata.ExponentialHistogram, memory reuse is missed.
	// In that case, use the zero-value h and hope for better alignment next cycle.
	h, _ := (*dest).(metricdata.ExponentialHistogram[N])
	h.Temporality = metricdata.CumulativeTemporality

	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	n := len(e.values)
	hDPts := reset(h.DataPoints, n, n)

	var i int
	for _, val := range e.values {
		hDPts[i].Attributes = val.attrs
		hDPts[i].StartTime = e.start
		hDPts[i].Time = t
		hDPts[i].Count = val.count()
		hDPts[i].Scale = val.scale
		hDPts[i].ZeroCount = val.zeroCount
		hDPts[i].ZeroThreshold = 0.0

		offset := val.posBuckets.loadCountsAndOffset(&hDPts[i].PositiveBucket.Counts)
		hDPts[i].PositiveBucket.Offset = offset

		offset = val.negBuckets.loadCountsAndOffset(&hDPts[i].NegativeBucket.Counts)
		hDPts[i].NegativeBucket.Offset = offset

		if !e.noSum {
			hDPts[i].Sum = val.sum
		}
		if !e.noMinMax {
			hDPts[i].Min = metricdata.NewExtrema(val.min)
			hDPts[i].Max = metricdata.NewExtrema(val.max)
		}

		collectExemplars(&hDPts[i].Exemplars, val.res.Collect)

		i++
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}

	h.DataPoints = hDPts
	*dest = h
	return n
}
