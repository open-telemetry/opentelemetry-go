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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"errors"
	"math"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

const (
	expoMaxScale = 20
	expoMinScale = -10

	smallestNonZeroNormalFloat64 = 0x1p-1022

	// These redefine the Math constants with a type, so the compiler won't coerce
	// them into an int on 32 bit platforms.
	maxInt64 int64 = math.MaxInt64
	minInt64 int64 = math.MinInt64
)

// expoHistogramValues summarizes a set of measurements as an histValues with
// explicitly defined buckets.
type expoHistogramValues[N int64 | float64] struct {
	noSum    bool
	noMinMax bool
	maxSize  int
	maxScale int

	values   map[attribute.Set]*expoHistogramDataPoint[N]
	valuesMu sync.Mutex
}

func newExpoHistValues[N int64 | float64](maxSize, maxScale int, noMinMax, noSum bool) *expoHistogramValues[N] {
	return &expoHistogramValues[N]{
		noSum:    noSum,
		noMinMax: noMinMax,
		maxSize:  maxSize,
		maxScale: maxScale,

		values: make(map[attribute.Set]*expoHistogramDataPoint[N]),
	}
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (e *expoHistogramValues[N]) Aggregate(value N, attr attribute.Set) {
	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	v, ok := e.values[attr]
	if !ok {
		v = newExpoHistogramDataPoint[N](e.maxSize, e.maxScale, e.noMinMax, e.noSum)
		e.values[attr] = v
	}
	v.record(value)
}

// expoHistogramDataPoint is a single bucket in an exponential histogram.
type expoHistogramDataPoint[N int64 | float64] struct {
	count uint64
	min   N
	max   N
	sum   N

	maxSize  int
	noMinMax bool
	noSum    bool

	scale int

	posBuckets expoBuckets
	negBuckets expoBuckets
	zeroCount  uint64
}

func newExpoHistogramDataPoint[N int64 | float64](maxSize, maxScale int, noMinMax, noSum bool) *expoHistogramDataPoint[N] {
	f := math.MaxFloat64
	max := N(f) // if N is int64, max will overflow to -9223372036854775808
	min := N(-f)
	if N(maxInt64) > N(f) {
		max = N(maxInt64)
		min = N(minInt64)
	}
	return &expoHistogramDataPoint[N]{
		min:      max,
		max:      min,
		maxSize:  maxSize,
		noMinMax: noMinMax,
		noSum:    noSum,
		scale:    maxScale,
	}
}

// record adds a new measurement to the histogram. It will rescale the buckets if needed.
func (p *expoHistogramDataPoint[N]) record(v N) {
	p.count++

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

	bin := getBin(absV, p.scale)

	bucket := &p.posBuckets
	if v < 0 {
		bucket = &p.negBuckets
	}

	// If the new bin would make the counts larger than maxScale, we need to
	// downscale current measurements.
	if scaleDelta := scaleChange(bin, bucket.startBin, len(bucket.counts), p.maxSize); scaleDelta > 0 {
		if p.scale-scaleDelta < expoMinScale {
			// With a scale of -10 there is only two buckets for the whole range of float64 values.
			// This can only happen if there is a max size of 1.
			otel.Handle(errors.New("exponential histogram scale underflow"))
			return
		}
		//Downscale
		p.scale -= scaleDelta
		p.posBuckets.downscale(scaleDelta)
		p.negBuckets.downscale(scaleDelta)

		bin = getBin(absV, p.scale)
	}

	bucket.record(bin)
}

// getBin returns the bin of the bucket that the value v should be recorded
// into at the given scale.
func getBin(v float64, scale int) int {
	frac, exp := math.Frexp(v)
	if scale <= 0 {
		// Because of the choice of fraction is always 1 power of two higher than we want.
		correction := 1
		if frac == .5 {
			// If v is an exact power of two the frac will be .5 and the exp
			// will be one higher than we want.
			correction = 2
		}
		return (exp - correction) >> (-scale)
	}
	return exp<<scale + int(math.Log(frac)*scaleFactors[scale]) - 1
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

// scaleChange returns the magnitude of the scale change needed to fit bin in the bucket.
func scaleChange(bin, startBin, length, maxSize int) int {
	if length == 0 {
		// No need to rescale if there are no buckets.
		return 0
	}

	low := startBin
	high := bin
	if startBin >= bin {
		low = bin
		high = startBin + length - 1
	}

	count := 0
	for high-low >= maxSize {
		low = low >> 1
		high = high >> 1
		count++
		if count > expoMaxScale-expoMinScale {
			return count
		}
	}
	return count
}

// expoBuckets is a single bucket in an exponential histogram.
type expoBuckets struct {
	startBin int
	counts   []uint64
}

// record increments the count for the given bin, and expands the buckets if needed.
// Size changes must be done before calling this function.
func (b *expoBuckets) record(bin int) {
	if len(b.counts) == 0 {
		b.counts = []uint64{1}
		b.startBin = bin
		return
	}

	endBin := b.startBin + len(b.counts) - 1

	// if the new bin is inside the current range
	if bin >= b.startBin && bin <= endBin {
		b.counts[bin-b.startBin]++
		return
	}
	// if the new bin is before the current start add spaces to the counts
	if bin < b.startBin {
		origLen := len(b.counts)
		newLength := endBin - bin + 1
		shift := b.startBin - bin

		if newLength > cap(b.counts) {
			b.counts = append(b.counts, make([]uint64, newLength-len(b.counts))...)
		}

		copy(b.counts[shift:origLen+shift], b.counts[:])
		b.counts = b.counts[:newLength]
		for i := 1; i < shift; i++ {
			b.counts[i] = 0
		}
		b.startBin = bin
		b.counts[0] = 1
		return
	}
	// if the new is after the end add spaces to the end
	if bin > endBin {
		if bin-b.startBin < cap(b.counts) {
			b.counts = b.counts[:bin-b.startBin+1]
			for i := endBin + 1 - b.startBin; i < len(b.counts); i++ {
				b.counts[i] = 0
			}
			b.counts[bin-b.startBin] = 1
			return
		}

		end := make([]uint64, bin-b.startBin-len(b.counts)+1)
		b.counts = append(b.counts, end...)
		b.counts[bin-b.startBin] = 1
	}
}

// downscale shrinks a bucket by a factor of 2*s. It will sum counts into the
// correct lower resolution bucket.
func (b *expoBuckets) downscale(delta int) {
	// Example
	// delta = 2
	// Original offset: -6
	// Counts: [ 3,  1,  2,  3,  4,  5, 6, 7, 8, 9, 10]
	// bins:    -6  -5, -4, -3, -2, -1, 0, 1, 2, 3, 4
	// new bins:-2, -2, -1, -1, -1, -1, 0, 0, 0, 0, 1
	// new Offset: -2
	// new Counts: [4, 14, 30, 10]

	if len(b.counts) <= 1 || delta < 1 {
		b.startBin = b.startBin >> delta
		return
	}

	steps := 1 << delta
	offset := b.startBin % steps
	offset = (offset + steps) % steps // to make offset positive
	for i := 1; i < len(b.counts); i++ {
		idx := i + offset
		if idx%steps == 0 {
			b.counts[idx/steps] = b.counts[i]
			continue
		}
		b.counts[idx/steps] += b.counts[i]
	}

	lastIdx := (len(b.counts) - 1 + offset) / steps
	b.counts = b.counts[:lastIdx+1]
	b.startBin = b.startBin >> delta
}

// newDeltaExponentialHistogram returns an Aggregator that summarizes a set of
// measurements as an exponential histogram. Each histogram is scoped by attributes
// and the aggregation cycle the measurements were made in.
//
// Each aggregation cycle is treated independently. When the returned
// Aggregator's Aggregations method is called it will reset all histogram
// counts to zero.
func newDeltaExponentialHistogram[N int64 | float64](cfg aggregation.Base2ExponentialHistogram, noSum bool) aggregator[N] {

	return &deltaExponentialHistogram[N]{
		expoHistogramValues: newExpoHistValues[N](
			int(cfg.MaxSize),
			int(cfg.MaxScale),
			cfg.NoMinMax,
			noSum,
		),
		start: now(),
	}
}

// deltaExponentialHistogram summarizes a set of measurements made in a single
// aggregation cycle as an Exponential histogram with explicitly defined buckets.
type deltaExponentialHistogram[N int64 | float64] struct {
	*expoHistogramValues[N]

	start time.Time
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (e *deltaExponentialHistogram[N]) Aggregation() metricdata.Aggregation {
	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	if len(e.values) == 0 {
		return nil
	}

	t := now()
	h := metricdata.ExponentialHistogram[N]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  make([]metricdata.ExponentialHistogramDataPoint[N], 0, len(e.values)),
	}
	for a, b := range e.values {
		ehdp := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    a,
			StartTime:     e.start,
			Time:          t,
			Count:         b.count,
			Scale:         int32(b.scale),
			ZeroCount:     b.zeroCount,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: int32(b.posBuckets.startBin),
				Counts: make([]uint64, len(b.posBuckets.counts)),
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: int32(b.negBuckets.startBin),
				Counts: make([]uint64, len(b.negBuckets.counts)),
			},
		}
		copy(ehdp.PositiveBucket.Counts, b.posBuckets.counts)
		copy(ehdp.NegativeBucket.Counts, b.negBuckets.counts)

		if !e.noMinMax {
			ehdp.Min = metricdata.NewExtrema(b.min)
			ehdp.Max = metricdata.NewExtrema(b.max)
		}
		if !e.noSum {
			ehdp.Sum = b.sum
		}
		h.DataPoints = append(h.DataPoints, ehdp)

		delete(e.values, a)
	}
	e.start = t
	return h
}

// newCumulativeExponentialHistogram returns an Aggregator that summarizes a set of
// measurements as an exponential histogram. Each histogram is scoped by attributes.
//
// Each aggregation cycle builds from the previous, the histogram counts are
// the bucketed counts of all values aggregated since the returned Aggregator
// was created.
func newCumulativeExponentialHistogram[N int64 | float64](cfg aggregation.Base2ExponentialHistogram, noSum bool) aggregator[N] {
	return &cumulativeExponentialHistogram[N]{
		expoHistogramValues: newExpoHistValues[N](
			int(cfg.MaxSize),
			int(cfg.MaxScale),
			cfg.NoMinMax,
			noSum,
		),
		start: now(),
	}
}

// cumulativeExponentialHistogram summarizes a set of measurements made in a single
// aggregation cycle as an Exponential histogram with explicitly defined buckets.
type cumulativeExponentialHistogram[N int64 | float64] struct {
	*expoHistogramValues[N]

	start time.Time
}

// Aggregate records the measurement, scoped by attr, and aggregates it
// into an aggregation.
func (e *cumulativeExponentialHistogram[N]) Aggregation() metricdata.Aggregation {
	e.valuesMu.Lock()
	defer e.valuesMu.Unlock()

	if len(e.values) == 0 {
		return nil
	}
	t := now()
	h := metricdata.ExponentialHistogram[N]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  make([]metricdata.ExponentialHistogramDataPoint[N], 0, len(e.values)),
	}
	for a, b := range e.values {
		ehdp := metricdata.ExponentialHistogramDataPoint[N]{
			Attributes:    a,
			StartTime:     e.start,
			Time:          t,
			Count:         b.count,
			Scale:         int32(b.scale),
			ZeroCount:     b.zeroCount,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: int32(b.posBuckets.startBin),
				Counts: make([]uint64, len(b.posBuckets.counts)),
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: int32(b.negBuckets.startBin),
				Counts: make([]uint64, len(b.negBuckets.counts)),
			},
		}
		copy(ehdp.PositiveBucket.Counts, b.posBuckets.counts)
		copy(ehdp.NegativeBucket.Counts, b.negBuckets.counts)

		if !e.noMinMax {
			ehdp.Min = metricdata.NewExtrema(b.min)
			ehdp.Max = metricdata.NewExtrema(b.max)
		}
		if !e.noSum {
			ehdp.Sum = b.sum
		}
		h.DataPoints = append(h.DataPoints, ehdp)
		// TODO (#3006): This will use an unbounded amount of memory if there
		// are unbounded number of attribute sets being aggregated. Attribute
		// sets that become "stale" need to be forgotten so this will not
		// overload the system.
	}

	return h
}
