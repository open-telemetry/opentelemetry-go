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

package exponential

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping/exponent"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping/logarithm"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
)

var (
	testDescriptor = metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
	intDescriptor  = metrictest.NewDescriptor("integer", sdkapi.HistogramInstrumentKind, number.Int64Kind)

	plusOne  = number.NewFloat64Number(1)
	minusOne = number.NewFloat64Number(-1)
)

type printableBucket struct {
	index int32
	count uint64
	lower float64
}

func (a *Aggregator) printBuckets(b *buckets) (r []printableBucket) {
	for i := uint32(0); i < b.Len(); i++ {
		lower, _ := a.state.mapping.LowerBoundary(b.Offset() + int32(i))
		r = append(r, printableBucket{
			index: b.Offset() + int32(i),
			count: b.At(i),
			lower: lower,
		})
	}
	return r
}

func getCounts(b *buckets) (r []uint64) {
	for i := uint32(0); i < b.Len(); i++ {
		r = append(r, b.At(i))
	}
	return r
}

func (s printableBucket) String() string {
	return fmt.Sprintf("%v=%v(%.2g)", s.index, s.count, s.lower)
}

func countNoError(t *testing.T, a *Aggregator) uint64 {
	cnt, err := a.Count()
	require.NoError(t, err)
	return cnt
}

func zeroCountNoError(t *testing.T, a *Aggregator) uint64 {
	cnt, err := a.ZeroCount()
	require.NoError(t, err)
	return cnt
}

func scaleNoError(t *testing.T, a *Aggregator) int32 {
	scale, err := a.Scale()
	require.NoError(t, err)
	return scale
}

func floatSumNoError(t *testing.T, a *Aggregator) float64 {
	sum, err := a.Sum()
	require.NoError(t, err)
	return sum.AsFloat64()
}

func intSumNoError(t *testing.T, a *Aggregator) int64 {
	sum, err := a.Sum()
	require.NoError(t, err)
	return sum.AsInt64()
}

// requireEqual is a helper used to require that two aggregators
// should have equal contents.  Because the backing array is cyclic,
// the two may are expected to have different underlying
// representations.
func requireEqual(t *testing.T, a, b *Aggregator) {
	aSum := floatSumNoError(t, a)
	bSum := floatSumNoError(t, b)
	if aSum == 0 || bSum == 0 {
		require.InDelta(t, aSum, bSum, 1e-6)
	} else {
		require.InEpsilon(t, aSum, bSum, 1e-6)
	}
	require.Equal(t, countNoError(t, a), countNoError(t, b))
	require.Equal(t, zeroCountNoError(t, a), zeroCountNoError(t, b))
	require.Equal(t, scaleNoError(t, a), scaleNoError(t, b))

	bstr := func(data *buckets) string {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintln("[@", data.Offset()))
		for i := uint32(0); i < data.Len(); i++ {
			sb.WriteString(fmt.Sprintln(data.At(i)))
		}
		sb.WriteString("]\n")
		return sb.String()
	}
	require.Equal(t, bstr(&a.state.positive), bstr(&b.state.positive), "positive %v %v", a.printBuckets(&a.state.positive), a.printBuckets(&b.state.positive))
	require.Equal(t, bstr(&a.state.negative), bstr(&b.state.negative), "negative %v %v", a.printBuckets(&a.state.negative), a.printBuckets(&b.state.negative))
}

// centerVal returns the midpoint of the histogram bucket with index
// `x`, used in tests to avoid rounding errors that happen near the
// bucket boundaries.
func centerVal(mapper mapping.Mapping, x int32) float64 {
	lb, err1 := mapper.LowerBoundary(x)
	ub, err2 := mapper.LowerBoundary(x + 1)
	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("unexpected errors: %v %v", err1, err2))
	}
	return (lb + ub) / 2
}

// Tests that the aggregation kind is correct.
func TestAggregationKind(t *testing.T) {
	agg := &New(1, &testDescriptor)[0]
	require.Equal(t, aggregation.ExponentialHistogramKind, agg.Aggregation().Kind())
}

// Tests insertion of [1, 2, 0.5].  The index of 1 (i.e., 0) becomes
// `indexBase`, the "2" goes to its right and the "0.5" goes in the
// last position of the backing array.  With 3 binary orders of
// magnitude and MaxSize=4, this must finish with scale=0 and offset=-1.
func TestAlternatingGrowth1(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	require.NoError(t, agg.Update(ctx, plusOne, &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(2), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0.5), &testDescriptor))

	require.Equal(t, int32(-1), agg.positive().Offset())
	require.Equal(t, int32(0), agg.scale())
	require.Equal(t, []uint64{1, 1, 1}, getCounts(agg.positive()))
}

// Tests insertion of [1, 1, 2, 0.5, 5, 0.25].  The test proceeds as
// above but then downscales once further to scale=-1, thus index -1
// holds range [0.25, 1.0), index 0 holds range [1.0, 4), index 1
// holds range [4, 16).
func TestAlternatingGrowth2(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	require.NoError(t, agg.Update(ctx, plusOne, &testDescriptor))
	require.NoError(t, agg.Update(ctx, plusOne, &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(2), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0.5), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(4), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0.25), &testDescriptor))

	require.Equal(t, int32(-1), agg.positive().Offset())
	require.Equal(t, int32(-1), agg.scale())
	require.Equal(t, []uint64{2, 3, 1}, getCounts(agg.positive()))
}

// Tests that every permutation of {1/2, 1, 2} with maxSize=2 results
// in the same scale=-1 histogram.
func TestScaleNegOneCentered(t *testing.T) {
	for j, order := range [][]float64{
		{1, 0.5, 2},
		{1, 2, 0.5},
		{2, 0.5, 1},
		{2, 1, 0.5},
		{0.5, 1, 2},
		{0.5, 2, 1},
	} {
		t.Run(fmt.Sprint(j), func(t *testing.T) {
			ctx := context.Background()
			agg := &New(1, &testDescriptor, WithMaxSize(2))[0]

			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor))
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor))

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// (1/2), counts[1] == 2 (the 1 and 2).
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor))
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(-1), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(1), agg.positive().At(0))
			require.Equal(t, uint64(2), agg.positive().At(1))
		})
	}
}

// Tests that every permutation of {1, 2, 4} with maxSize=2 results in
// the same scale=-1 histogram.
func TestScaleNegOnePositive(t *testing.T) {
	for j, order := range [][]float64{
		{1, 2, 4},
		{1, 4, 2},
		{2, 4, 1},
		{2, 1, 4},
		{4, 1, 2},
		{4, 2, 1},
	} {
		t.Run(fmt.Sprint(j), func(t *testing.T) {
			ctx := context.Background()
			agg := &New(1, &testDescriptor, WithMaxSize(2))[0]

			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor))
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor))

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// 1 and 2), counts[1] == 2 (the 4).
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor))
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(0), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(2), agg.positive().At(0))
			require.Equal(t, uint64(1), agg.positive().At(1))
		})
	}
}

// Tests that every permutation of {1, 1/2, 1/4} with maxSize=2
// results in the same scale=-1 histogram.
func TestScaleNegOneNegative(t *testing.T) {
	for j, order := range [][]float64{
		{1, 0.5, 0.25},
		{1, 0.25, 0.5},
		{0.5, 0.25, 1},
		{0.5, 1, 0.25},
		{0.25, 1, 0.5},
		{0.25, 0.5, 1},
	} {
		t.Run(fmt.Sprint(j), func(t *testing.T) {
			ctx := context.Background()
			agg := &New(1, &testDescriptor, WithMaxSize(2))[0]

			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor))
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor))

			// Enter order[2]: scale set to -1, expect counts[0] == 2 (the
			// 1/4 and 1/2, counts[1] == 2 (the 1).
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor))
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(-1), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(2), agg.positive().At(0))
			require.Equal(t, uint64(1), agg.positive().At(1))
		})
	}
}

// Tests a variety of ascending sequences, calculated using known
// index ranges.  For example, with maxSize=3, using scale=0 and
// offset -5, add a sequence of numbers. Because the numbers have
// known range, we know the expected scale.
func TestAscendingSequence(t *testing.T) {
	for _, maxSize := range []int32{3, 4, 6, 9} {
		t.Run(fmt.Sprintf("maxSize=%d", maxSize), func(t *testing.T) {
			for offset := int32(-5); offset <= 5; offset++ {
				for _, initScale := range []int32{
					0, 4,
				} {
					testAscendingSequence(t, maxSize, offset, initScale)
				}
			}
		})
	}
}

func testAscendingSequence(t *testing.T, maxSize, offset, initScale int32) {
	for step := maxSize; step < 4*maxSize; step++ {
		ctx := context.Background()
		agg := &New(1, &testDescriptor, WithMaxSize(maxSize))[0]
		mapper, err := newMapping(initScale)
		require.NoError(t, err)

		minVal := centerVal(mapper, offset)
		maxVal := centerVal(mapper, offset+step)
		sum := 0.0

		for i := int32(0); i < maxSize; i++ {
			value := centerVal(mapper, offset+i)
			require.NoError(t, agg.Update(ctx, number.NewFloat64Number(value), &testDescriptor))
			sum += value
		}

		require.Equal(t, initScale, agg.scale())
		require.Equal(t, offset, agg.positive().Offset())

		require.NoError(t, agg.Update(ctx, number.NewFloat64Number(maxVal), &testDescriptor))
		sum += maxVal

		// The zeroth bucket is not empty.
		require.NotEqual(t, uint64(0), agg.positive().At(0))

		// The maximum-index filled bucket is at or
		// above the mid-point, (otherwise we
		// downscaled too much).
		maxFill := uint32(0)
		totalCount := uint64(0)

		for i := uint32(0); i < agg.positive().Len(); i++ {
			totalCount += agg.positive().At(i)
			if agg.positive().At(i) != 0 {
				maxFill = i
			}
		}
		require.GreaterOrEqual(t, maxFill, uint32(maxSize)/2)

		// Count is correct
		require.GreaterOrEqual(t, uint64(maxSize+1), totalCount)
		hcount, _ := agg.Count()
		require.GreaterOrEqual(t, uint64(maxSize+1), hcount)
		// Sum is correct
		hsum, _ := agg.Sum()
		require.GreaterOrEqual(t, sum, hsum.CoerceToFloat64(number.Float64Kind))

		// The offset is correct at the computed scale.
		mapper, err = newMapping(agg.scale())
		require.NoError(t, err)
		idx := mapper.MapToIndex(minVal)
		require.Equal(t, int32(idx), agg.positive().Offset())

		// The maximum range is correct at the computed scale.
		idx = mapper.MapToIndex(maxVal)
		require.Equal(t, int32(idx), agg.positive().Offset()+int32(agg.positive().Len())-1)
	}
}

// Tests a simple case of merging [1, 2, 4, 8] with [1/2, 1/4, 1/8, 1/16].
func TestMergeSimpleEven(t *testing.T) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(4))
	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)   // 1, 2, 4, 8
		f2 := 1 / float64(int64(2)<<i) // 1/2, 1/4, 1/8, 1/16
		n1 := number.NewFloat64Number(f1)
		n2 := number.NewFloat64Number(f2)

		require.NoError(t, aggs[0].Update(ctx, n1, &testDescriptor))
		require.NoError(t, aggs[1].Update(ctx, n2, &testDescriptor))
		require.NoError(t, aggs[2].Update(ctx, n1, &testDescriptor))
		require.NoError(t, aggs[2].Update(ctx, n2, &testDescriptor))
	}
	require.Equal(t, int32(0), aggs[0].scale())
	require.Equal(t, int32(0), aggs[1].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	require.Equal(t, int32(0), aggs[0].positive().Offset())
	require.Equal(t, int32(-4), aggs[1].positive().Offset())
	require.Equal(t, int32(-2), aggs[2].positive().Offset())

	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(aggs[0].positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(aggs[1].positive()))
	require.Equal(t, []uint64{2, 2, 2, 2}, getCounts(aggs[2].positive()))

	require.NoError(t, aggs[0].Merge(&aggs[1], &testDescriptor))

	require.Equal(t, int32(-1), aggs[0].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	requireEqual(t, &aggs[0], &aggs[2])
}

// Tests a simple case of merging [1, 2, 4, 8] with [1, 1/2, 1/4, 1/8].
func TestMergeSimpleOdd(t *testing.T) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(4))
	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)
		f2 := 1 / float64(int64(1)<<i) // Diff from above test: 1 here vs 2 above.
		n1 := number.NewFloat64Number(f1)
		n2 := number.NewFloat64Number(f2)

		require.NoError(t, aggs[0].Update(ctx, n1, &testDescriptor))
		require.NoError(t, aggs[1].Update(ctx, n2, &testDescriptor))
		require.NoError(t, aggs[2].Update(ctx, n1, &testDescriptor))
		require.NoError(t, aggs[2].Update(ctx, n2, &testDescriptor))
	}

	require.Equal(t, uint64(4), aggs[0].state.count)
	require.Equal(t, uint64(4), aggs[1].state.count)
	require.Equal(t, uint64(8), aggs[2].state.count)

	require.Equal(t, int32(0), aggs[0].scale())
	require.Equal(t, int32(0), aggs[1].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	require.Equal(t, int32(0), aggs[0].positive().Offset())
	require.Equal(t, int32(-3), aggs[1].positive().Offset())
	require.Equal(t, int32(-2), aggs[2].positive().Offset())

	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(aggs[0].positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(aggs[1].positive()))
	require.Equal(t, []uint64{1, 2, 3, 2}, getCounts(aggs[2].positive()))

	require.NoError(t, aggs[0].Merge(&aggs[1], &testDescriptor))

	require.Equal(t, int32(-1), aggs[0].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	requireEqual(t, &aggs[0], &aggs[2])
}

// Tests a random data set, exhaustively partitioned in every way, ensuring that
// computing the aggregations and merging them produces the same result as computing
// a single aggregation.
func TestMergeExhaustive(t *testing.T) {
	const (
		factor = 1024.0
		repeat = 1
		count  = 16
	)

	means := []float64{
		0,
		factor,
	}

	stddevs := []float64{
		1,
		factor,
	}

	for _, mean := range means {
		t.Run(fmt.Sprint("mean=", mean), func(t *testing.T) {
			for _, stddev := range stddevs {
				t.Run(fmt.Sprint("stddev=", stddev), func(t *testing.T) {
					src := rand.NewSource(77777677777)
					rnd := rand.New(src)

					values := make([]float64, count)
					for i := range values {
						values[i] = mean + rnd.NormFloat64()*stddev
					}

					for part := 1; part < count; part++ {
						for _, size := range []int32{
							2,
							6,
							8,
							9,
							16,
						} {
							for _, incr := range []uint64{
								1,
								0x100,
								0x10000,
								0x100000000,
							} {
								testMergeExhaustive(t, values[0:part], values[part:count], size, incr)
							}
						}
					}
				})
			}
		})
	}
}

func testMergeExhaustive(t *testing.T, a, b []float64, size int32, incr uint64) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(size))

	aHist := &aggs[0]
	bHist := &aggs[1]
	cHist := &aggs[2]

	for _, av := range a {
		require.NoError(t, aHist.UpdateByIncr(ctx, number.NewFloat64Number(av), incr, &testDescriptor))
		require.NoError(t, cHist.UpdateByIncr(ctx, number.NewFloat64Number(av), incr, &testDescriptor))
	}
	for _, bv := range b {
		require.NoError(t, bHist.UpdateByIncr(ctx, number.NewFloat64Number(bv), incr, &testDescriptor))
		require.NoError(t, cHist.UpdateByIncr(ctx, number.NewFloat64Number(bv), incr, &testDescriptor))
	}

	require.NoError(t, aHist.Merge(bHist, &testDescriptor))

	// aHist and cHist should be equivalent
	requireEqual(t, cHist, aHist)
}

// Tests the logic to switch between uint8, uint16, uint32, and
// uint64.  Test is based on the UpdateByIncr code path.
func TestOverflowBits(t *testing.T) {
	ctx := context.Background()

	for _, limit := range []uint64{
		0x100,
		0x10000,
		0x100000000,
	} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			aggs := New(3, &testDescriptor)

			aHist := &aggs[0]
			bHist := &aggs[1]
			cHist := &aggs[2]

			if limit <= 0x10000 {
				for i := uint64(0); i < limit; i++ {
					require.NoError(t, aHist.Update(ctx, plusOne, &testDescriptor))
					require.NoError(t, aHist.Update(ctx, minusOne, &testDescriptor))

					cnt, _ := aHist.Count()
					require.Equal(t, 2*(i+1), cnt)
				}
			} else {
				require.NoError(t, aHist.UpdateByIncr(ctx, plusOne, limit/2, &testDescriptor))
				require.NoError(t, aHist.UpdateByIncr(ctx, plusOne, limit/2, &testDescriptor))
				require.NoError(t, aHist.UpdateByIncr(ctx, minusOne, limit/2, &testDescriptor))
				require.NoError(t, aHist.UpdateByIncr(ctx, minusOne, limit/2, &testDescriptor))
			}
			require.NoError(t, bHist.UpdateByIncr(ctx, plusOne, limit-1, &testDescriptor))
			require.NoError(t, bHist.Update(ctx, plusOne, &testDescriptor))
			require.NoError(t, bHist.UpdateByIncr(ctx, minusOne, limit-1, &testDescriptor))
			require.NoError(t, bHist.Update(ctx, minusOne, &testDescriptor))
			require.NoError(t, cHist.UpdateByIncr(ctx, plusOne, limit, &testDescriptor))
			require.NoError(t, cHist.UpdateByIncr(ctx, minusOne, limit, &testDescriptor))

			aCnt, _ := aHist.Count()
			require.Equal(t, 2*limit, aCnt)
			sum, _ := aHist.Sum()
			require.Equal(t, float64(0), sum.AsFloat64())

			aPos, _ := aHist.Positive()
			require.Equal(t, uint32(1), aPos.Len())
			require.Equal(t, limit, aPos.At(0))

			aNeg, _ := aHist.Negative()
			require.Equal(t, uint32(1), aNeg.Len())
			require.Equal(t, limit, aNeg.At(0))

			requireEqual(t, cHist, aHist)
			requireEqual(t, bHist, aHist)
		})
	}
}

// Tests the use of number.Int64Kind as opposed to floating point. The
// aggregator internal state is identical except for the Sum, which is
// maintained as a `number.Number`.
func TestIntegerAggregation(t *testing.T) {
	ctx := context.Background()
	aggs := New(2, &intDescriptor, WithMaxSize(256))
	agg := &aggs[0]
	alt := &aggs[1]

	expect := int64(0)
	for i := int64(1); i < 256; i++ {
		expect += i
		require.NoError(t, agg.Update(ctx, number.NewInt64Number(i), &intDescriptor))
		require.NoError(t, alt.Update(ctx, number.NewInt64Number(i), &intDescriptor))
	}

	require.Equal(t, expect, intSumNoError(t, agg))
	require.Equal(t, uint64(255), countNoError(t, agg))

	// Scale should be 5.  Here's why.  The upper power-of-two is
	// 256 == 2**8.  We expect the exponential base = 2**(2**-5)
	// raised to the 256th power to equal 256:
	//
	//   2**((2**-5)*256)
	// = 2**((2**-5)*(2**8))
	// = 2**(2**3)
	// = 2**8
	require.Equal(t, int32(5), scaleNoError(t, agg))

	expect0 := func(b aggregation.ExponentialBuckets) {
		require.Equal(t, uint32(0), b.Len())
	}
	expect256 := func(b aggregation.ExponentialBuckets, factor int) {
		require.Equal(t, uint32(256), b.Len())
		require.Equal(t, int32(0), b.Offset())
		// Bucket 254 has 6 elements, bucket 255 has 5
		// bucket 253 has 5, ...
		for i := uint32(0); i < 256; i++ {
			require.LessOrEqual(t, b.At(i), uint64(6*factor))
		}
	}

	pos, err := agg.Positive()
	require.NoError(t, err)
	expect256(pos, 1)

	neg, err := agg.Negative()
	require.NoError(t, err)
	expect0(neg)

	// Merge!
	require.NoError(t, agg.Merge(alt, &intDescriptor))

	pos, err = agg.Positive()
	require.NoError(t, err)
	expect256(pos, 2)

	require.Equal(t, 2*expect, intSumNoError(t, agg))

	// Reset!  Repeat with negative.
	require.NoError(t, agg.SynchronizedMove(nil, &intDescriptor))
	require.NoError(t, alt.SynchronizedMove(nil, &intDescriptor))

	expect = int64(0)
	for i := int64(1); i < 256; i++ {
		expect -= i
		require.NoError(t, agg.Update(ctx, number.NewInt64Number(-i), &intDescriptor))
		require.NoError(t, alt.Update(ctx, number.NewInt64Number(-i), &intDescriptor))
	}

	require.Equal(t, expect, intSumNoError(t, agg))
	require.Equal(t, uint64(255), countNoError(t, agg))

	neg, err = agg.Negative()
	require.NoError(t, err)
	expect256(neg, 1)

	pos, err = agg.Positive()
	require.NoError(t, err)
	expect0(pos)

	// Merge!
	require.NoError(t, agg.Merge(alt, &intDescriptor))

	neg, err = agg.Negative()
	require.NoError(t, err)
	expect256(neg, 2)

	require.Equal(t, 2*expect, intSumNoError(t, agg))

	// Scale should not change after filling in the negative range.
	require.Equal(t, int32(5), scaleNoError(t, agg))
}

// Tests the reset code path via SynchronizedMove.
func TestReset(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(256))[0]

	for _, incr := range []uint64{
		1,
		0x100,
		0x10000,
		0x100000000,

		// Another 32-bit increment tests the 64-bit reset path.
		0x200000000,
	} {
		t.Run(fmt.Sprint(incr), func(t *testing.T) {
			require.NoError(t, agg.SynchronizedMove(nil, &testDescriptor))

			// Note that scale is zero b/c no values
			require.Equal(t, int32(0), scaleNoError(t, agg))

			expect := 0.0
			for i := int64(1); i < 256; i++ {
				expect += float64(i) * float64(incr)
				err := agg.UpdateByIncr(ctx, number.NewFloat64Number(float64(i)), incr, &testDescriptor)
				require.NoError(t, err)
			}

			require.Equal(t, expect, floatSumNoError(t, agg))
			require.Equal(t, uint64(255)*incr, countNoError(t, agg))

			// See TestIntegerAggregation about why scale is 5.
			require.Equal(t, int32(5), scaleNoError(t, agg))

			pos, err := agg.Positive()
			require.NoError(t, err)

			require.Equal(t, uint32(256), pos.Len())
			require.Equal(t, int32(0), pos.Offset())
			// Bucket 254 has 6 elements, bucket 255 has 5
			// bucket 253 has 5, ...
			for i := uint32(0); i < 256; i++ {
				require.LessOrEqual(t, pos.At(i), uint64(6)*incr)
			}
		})
	}

}

// Tests the move aspect of SynchronizedMove.
func TestMove(t *testing.T) {
	ctx := context.Background()
	aggs := New(2, &testDescriptor, WithMaxSize(256))
	agg := &aggs[0]
	cpy := &aggs[1]

	expect := 0.0
	for i := int64(1); i < 256; i++ {
		expect += float64(i)
		require.NoError(t, agg.Update(ctx, number.NewFloat64Number(float64(i)), &testDescriptor))
		require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0), &testDescriptor))
	}

	require.NoError(t, agg.SynchronizedMove(cpy, &testDescriptor))

	// agg was reset
	require.Equal(t, 0.0, floatSumNoError(t, agg))
	require.Equal(t, uint64(0), countNoError(t, agg))
	require.Equal(t, uint64(0), zeroCountNoError(t, agg))
	require.Equal(t, int32(0), scaleNoError(t, agg))

	// cpy is as expected
	require.Equal(t, expect, floatSumNoError(t, cpy))
	require.Equal(t, uint64(255*2), countNoError(t, cpy))
	require.Equal(t, uint64(255), zeroCountNoError(t, cpy))

	// See TestIntegerAggregation about why scale is 5,
	// max bucket count is 6, and so on.
	require.Equal(t, int32(5), scaleNoError(t, cpy))

	pos, err := cpy.Positive()
	require.NoError(t, err)

	require.Equal(t, uint32(256), pos.Len())
	require.Equal(t, int32(0), pos.Offset())
	for i := uint32(0); i < 256; i++ {
		require.LessOrEqual(t, pos.At(i), uint64(6))
	}
}

// Tests with maxSize=2 that very large numbers (but not the full
// range) yield scales -7 and -8.
func TestVeryLargeNumbers(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(2))[0]

	expectBalanced := func(c uint64) {
		pos, _ := agg.Positive()
		require.Equal(t, uint32(2), pos.Len())
		require.Equal(t, int32(-1), pos.Offset())
		require.Equal(t, c, pos.At(0))
		require.Equal(t, c, pos.At(1))
	}

	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p-100), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p+100), &testDescriptor))

	require.InEpsilon(t, 0x1p100, floatSumNoError(t, agg), 1e-5)
	require.Equal(t, uint64(2), countNoError(t, agg))
	require.Equal(t, int32(-7), scaleNoError(t, agg))

	expectBalanced(1)

	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p-128), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p+127), &testDescriptor))

	require.InEpsilon(t, 0x1p127, floatSumNoError(t, agg), 1e-5)
	require.Equal(t, uint64(4), countNoError(t, agg))
	require.Equal(t, int32(-7), scaleNoError(t, agg))

	expectBalanced(2)

	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p-129), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(0x1p+255), &testDescriptor))

	require.InEpsilon(t, 0x1p255, floatSumNoError(t, agg), 1e-5)
	require.Equal(t, uint64(6), countNoError(t, agg))
	require.Equal(t, int32(-8), scaleNoError(t, agg))

	expectBalanced(3)
}

// Tests that WithRangeLimit() performs as specified, restricting the
// histogram scale to best cover the specified range.
func TestFixedLimits(t *testing.T) {
	ctx := context.Background()

	for _, test := range []struct {
		min, max float64
		maxSize  int32
	}{
		{0.25, 4, 8},
		{0.25, 4, 7},
		{0.5, 8, 8},
		{0.5, 8, 7},

		{0.001, 60000, 256},
		{1e200, 1e300, 256},

		{0.1, 1, 256},
		{0.1, 1, 5},
		{1, 10, 5},
		{1, 1e300, 5},
		{1e-300, 1e300, 5},
	} {
		var expectScale int32
		sizeAtScale := func(s int32) int32 {
			m, err := newMapping(s)
			require.NoError(t, err)
			return m.MapToIndex(test.max) - m.MapToIndex(test.min) + 1
		}

		// Find the ideal scale exhaustively: choose the largest scale
		// value for which the size test below holds.
		for s := exponent.MinScale; s <= logarithm.MaxScale; s++ {
			sz := sizeAtScale(s)
			if sz <= test.maxSize && sz > test.maxSize/2 {
				// Note that we do not break the loop here because
				// in case of odd maximum size, it's possible for
				// the next larger scale to work.  See the comment
				// about the same issue in changeScale().
				expectScale = s
			}
		}

		agg := &New(1, &testDescriptor, WithRangeLimit(test.min, test.max), WithMaxSize(test.maxSize))[0]

		// Scale should be set correctly before any updates
		scale, _ := agg.Scale()
		require.Equal(t, expectScale, scale)

		// Update: average value
		require.NoError(t, agg.Update(ctx, number.NewFloat64Number((test.min+test.max)/2), &testDescriptor))

		scale, _ = agg.Scale()
		require.Equal(t, int32(expectScale), scale)

		// Update: min and max values
		require.NoError(t, agg.Update(ctx, number.NewFloat64Number(test.min), &testDescriptor))
		require.NoError(t, agg.Update(ctx, number.NewFloat64Number(test.max), &testDescriptor))

		scale, _ = agg.Scale()
		require.Equal(t, int32(expectScale), scale)

		// No negatives
		neg, _ := agg.Negative()
		require.Equal(t, uint32(0), neg.Len())

		// Positive count 3, expected size
		pos, _ := agg.Positive()
		require.Equal(t, uint32(sizeAtScale(expectScale)), pos.Len())
		cnt, _ := agg.Count()
		require.Equal(t, uint64(3), cnt)

		// Error case
		require.Equal(t, mapping.ErrUnderflow, agg.Update(ctx, number.NewFloat64Number(test.min*0.99), &testDescriptor))
		require.Equal(t, mapping.ErrOverflow, agg.Update(ctx, number.NewFloat64Number(test.max*1.01), &testDescriptor))

		// Make sure count didn't change
		cnt, _ = agg.Count()
		require.Equal(t, uint64(3), cnt)
	}
}

// Tests the largest and smallest finite numbers with below-minimum
// size.  Expect a size=MinSize histogram with MinScale.
func TestFullRange(t *testing.T) {
	ctx := context.Background()
	aggs := New(1, &testDescriptor, WithMaxSize(1))
	agg := &aggs[0]

	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(math.MaxFloat64), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(1), &testDescriptor))
	require.NoError(t, agg.Update(ctx, number.NewFloat64Number(math.SmallestNonzeroFloat64), &testDescriptor))

	require.Equal(t, logarithm.MaxValue, floatSumNoError(t, agg))
	require.Equal(t, uint64(3), countNoError(t, agg))

	require.Equal(t, exponent.MinScale, scaleNoError(t, agg))

	pos, err := agg.Positive()
	require.NoError(t, err)

	require.Equal(t, uint32(MinSize), pos.Len())
	require.Equal(t, int32(-1), pos.Offset())
	require.Equal(t, pos.At(0), uint64(1))
	require.Equal(t, pos.At(1), uint64(2))
}

// Tests the inconsistent aggregator checks.
func TestInconsistentAggregator(t *testing.T) {
	agg := &New(1, &testDescriptor)[0]
	wrong := &sum.New(1)[0]

	err := agg.SynchronizedMove(wrong, &testDescriptor)
	require.Error(t, err)
	require.True(t, errors.Is(err, aggregation.ErrInconsistentType))

	err = agg.Merge(wrong, &testDescriptor)
	require.Error(t, err)
	require.True(t, errors.Is(err, aggregation.ErrInconsistentType))
}

// Benchmarks the Update() function for values in the range [1,2)
func BenchmarkLinear(b *testing.B) {
	ctx := context.Background()
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := &New(1, &testDescriptor, WithMaxSize(1024))[0]
	for i := 0; i < b.N; i++ {
		x := 2 - rnd.Float64()
		_ = agg.Update(ctx, number.NewFloat64Number(x), &testDescriptor)
	}
}

// Benchmarks the Update() function for values in the range [1,2) with fixed scale
func BenchmarkLinearFixed(b *testing.B) {
	ctx := context.Background()
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := &New(1, &testDescriptor, WithMaxSize(1024), WithRangeLimit(1, 1.999999))[0]
	for i := 0; i < b.N; i++ {
		x := 2 - rnd.Float64()
		_ = agg.Update(ctx, number.NewFloat64Number(x), &testDescriptor)
	}
}

// Benchmarks the Update() function for values in the range (0, MaxValue]
func BenchmarkExponential(b *testing.B) {
	ctx := context.Background()
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := &New(1, &testDescriptor, WithMaxSize(1024))[0]
	for i := 0; i < b.N; i++ {
		x := rnd.ExpFloat64()
		_ = agg.Update(ctx, number.NewFloat64Number(x), &testDescriptor)
	}
}
