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

package structure // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/structure"

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/logarithm"
)

const (
	plusOne  = 1
	minusOne = -1
)

type printableBucket struct {
	index int32
	count uint64
	lower float64
}

func (h *Histogram[N]) printBuckets(b *Buckets) (r []printableBucket) {
	for i := uint32(0); i < b.Len(); i++ {
		lower, _ := h.mapping.LowerBoundary(b.Offset() + int32(i))
		r = append(r, printableBucket{
			index: b.Offset() + int32(i),
			count: b.At(i),
			lower: lower,
		})
	}
	return r
}

func getCounts(b *Buckets) (r []uint64) {
	for i := uint32(0); i < b.Len(); i++ {
		r = append(r, b.At(i))
	}
	return r
}

func (b printableBucket) String() string {
	return fmt.Sprintf("%v=%v(%.2g)", b.index, b.count, b.lower)
}

// requireEqual is a helper used to require that two aggregators
// should have equal contents.  Because the backing array is cyclic,
// the two may are expected to have different underlying
// representations.  This method is more useful than RequireEqualValues
// for debugging the internals, because it prints numeric boundaries.
func requireEqual(t *testing.T, a, b *Histogram[float64]) {
	aSum := a.Sum()
	bSum := b.Sum()
	if aSum == 0 || bSum == 0 {
		require.InDelta(t, aSum, bSum, 1e-6)
	} else {
		require.InEpsilon(t, aSum, bSum, 1e-6)
	}
	require.Equal(t, a.Count(), b.Count())
	require.Equal(t, a.ZeroCount(), b.ZeroCount())
	require.Equal(t, a.Scale(), b.Scale())

	bstr := func(data *Buckets) string {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintln("[@", data.Offset()))
		for i := uint32(0); i < data.Len(); i++ {
			sb.WriteString(fmt.Sprintln(data.At(i)))
		}
		sb.WriteString("]\n")
		return sb.String()
	}
	require.Equal(t, bstr(&a.positive), bstr(&b.positive), "positive %v %v", a.printBuckets(&a.positive), a.printBuckets(&b.positive))
	require.Equal(t, bstr(&a.negative), bstr(&b.negative), "negative %v %v", a.printBuckets(&a.negative), a.printBuckets(&b.negative))
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

// Tests insertion of [1, 2, 0.5].  The index of 1 (i.e., 0) becomes
// `indexBase`, the "2" goes to its right and the "0.5" goes in the
// last position of the backing array.  With 3 binary orders of
// magnitude and MaxSize=4, this must finish with scale=0 and offset=-1.
func TestAlternatingGrowth1(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(4)))
	agg.Update(1)
	agg.Update(2)
	agg.Update(0.5)

	require.Equal(t, int32(-1), agg.Positive().Offset())
	require.Equal(t, int32(0), agg.Scale())
	require.Equal(t, []uint64{1, 1, 1}, getCounts(agg.Positive()))
}

// Tests insertion of [1, 1, 2, 0.5, 5, 0.25].  The test proceeds as
// above but then downscales once further to scale=-1, thus index -1
// holds range [0.25, 1.0), index 0 holds range [1.0, 4), index 1
// holds range [4, 16).
func TestAlternatingGrowth2(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(4)))
	agg.Update(1)
	agg.Update(1)
	agg.Update(2)
	agg.Update(0.5)
	agg.Update(4)
	agg.Update(0.25)

	require.Equal(t, int32(-1), agg.Positive().Offset())
	require.Equal(t, int32(-1), agg.Scale())
	require.Equal(t, []uint64{2, 3, 1}, getCounts(agg.Positive()))
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
			agg := NewFloat64(NewConfig(WithMaxSize(2)), order...)

			// After three updates: scale set to -1, expect counts[0] == 1 (the
			// (1/2), counts[1] == 2 (the 1 and 2).

			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(-1), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(1), agg.Positive().At(0))
			require.Equal(t, uint64(2), agg.Positive().At(1))
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
			agg := NewFloat64(NewConfig(WithMaxSize(2)), order...)

			// After three updates: scale set to -1, expect counts[0] == 1 (the
			// 1 and 2), counts[1] == 2 (the 4).
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(0), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(2), agg.Positive().At(0))
			require.Equal(t, uint64(1), agg.Positive().At(1))
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
			agg := NewFloat64(NewConfig(WithMaxSize(2)), order...)

			// After 3 updates: scale set to -1, expect counts[0] == 2 (the
			// 1/4 and 1/2, counts[1] == 2 (the 1).
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(-1), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(2), agg.Positive().At(0))
			require.Equal(t, uint64(1), agg.Positive().At(1))
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
		agg := NewFloat64(NewConfig(WithMaxSize(maxSize)))
		mapper, err := newMapping(initScale)
		require.NoError(t, err)

		minVal := centerVal(mapper, offset)
		maxVal := centerVal(mapper, offset+step)
		sum := 0.0

		for i := int32(0); i < maxSize; i++ {
			value := centerVal(mapper, offset+i)
			agg.Update(value)
			sum += value
		}

		require.Equal(t, initScale, agg.Scale())
		require.Equal(t, offset, agg.Positive().Offset())

		agg.Update(maxVal)
		sum += maxVal

		// The zeroth bucket is not empty.
		require.NotEqual(t, uint64(0), agg.Positive().At(0))

		// The maximum-index filled bucket is at or
		// above the mid-point, (otherwise we
		// downscaled too much).
		maxFill := uint32(0)
		totalCount := uint64(0)

		for i := uint32(0); i < agg.Positive().Len(); i++ {
			totalCount += agg.Positive().At(i)
			if agg.Positive().At(i) != 0 {
				maxFill = i
			}
		}
		require.GreaterOrEqual(t, maxFill, uint32(maxSize)/2)

		// Count is correct
		require.GreaterOrEqual(t, uint64(maxSize+1), totalCount)
		require.GreaterOrEqual(t, uint64(maxSize+1), agg.Count())
		// Sum is correct
		require.GreaterOrEqual(t, sum, agg.Sum())

		// The offset is correct at the computed scale.
		mapper, err = newMapping(agg.Scale())
		require.NoError(t, err)
		idx := mapper.MapToIndex(minVal)
		require.Equal(t, int32(idx), agg.Positive().Offset())

		// The maximum range is correct at the computed scale.
		idx = mapper.MapToIndex(maxVal)
		require.Equal(t, int32(idx), agg.Positive().Offset()+int32(agg.Positive().Len())-1)
	}
}

// Tests a simple case of merging [1, 2, 4, 8] with [1/2, 1/4, 1/8, 1/16].
func TestMergeSimpleEven(t *testing.T) {
	agg0 := NewFloat64(NewConfig(WithMaxSize(4)))
	agg1 := NewFloat64(NewConfig(WithMaxSize(4)))
	agg2 := NewFloat64(NewConfig(WithMaxSize(4)))

	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)   // 1, 2, 4, 8
		f2 := 1 / float64(int64(2)<<i) // 1/2, 1/4, 1/8, 1/16

		agg0.Update(f1)
		agg1.Update(f2)
		agg2.Update(f1)
		agg2.Update(f2)
	}
	require.Equal(t, int32(0), agg0.Scale())
	require.Equal(t, int32(0), agg1.Scale())
	require.Equal(t, int32(-1), agg2.Scale())

	require.Equal(t, int32(0), agg0.Positive().Offset())
	require.Equal(t, int32(-4), agg1.Positive().Offset())
	require.Equal(t, int32(-2), agg2.Positive().Offset())

	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(agg0.Positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(agg1.Positive()))
	require.Equal(t, []uint64{2, 2, 2, 2}, getCounts(agg2.Positive()))

	agg0.MergeFrom(agg1)

	require.Equal(t, int32(-1), agg0.Scale())
	require.Equal(t, int32(-1), agg2.Scale())

	requireEqual(t, agg0, agg2)
}

// Tests a simple case of merging [1, 2, 4, 8] with [1, 1/2, 1/4, 1/8].
func TestMergeSimpleOdd(t *testing.T) {
	agg0 := NewFloat64(NewConfig(WithMaxSize(4)))
	agg1 := NewFloat64(NewConfig(WithMaxSize(4)))
	agg2 := NewFloat64(NewConfig(WithMaxSize(4)))

	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)
		f2 := 1 / float64(int64(1)<<i) // Diff from above test: 1 here vs 2 above.

		agg0.Update(f1)
		agg1.Update(f2)
		agg2.Update(f1)
		agg2.Update(f2)
	}

	require.Equal(t, uint64(4), agg0.Count())
	require.Equal(t, uint64(4), agg1.Count())
	require.Equal(t, uint64(8), agg2.Count())

	require.Equal(t, int32(0), agg0.Scale())
	require.Equal(t, int32(0), agg1.Scale())
	require.Equal(t, int32(-1), agg2.Scale())

	require.Equal(t, int32(0), agg0.Positive().Offset())
	require.Equal(t, int32(-3), agg1.Positive().Offset())
	require.Equal(t, int32(-2), agg2.Positive().Offset())

	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(agg0.Positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, getCounts(agg1.Positive()))
	require.Equal(t, []uint64{1, 2, 3, 2}, getCounts(agg2.Positive()))

	agg0.MergeFrom(agg1)

	require.Equal(t, int32(-1), agg0.Scale())
	require.Equal(t, int32(-1), agg2.Scale())

	requireEqual(t, agg0, agg2)
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
	aHist := NewFloat64(NewConfig(WithMaxSize(size)))
	bHist := NewFloat64(NewConfig(WithMaxSize(size)))
	cHist := NewFloat64(NewConfig(WithMaxSize(size)))

	for _, av := range a {
		aHist.UpdateByIncr(av, incr)
		cHist.UpdateByIncr(av, incr)
	}
	for _, bv := range b {
		bHist.UpdateByIncr(bv, incr)
		cHist.UpdateByIncr(bv, incr)
	}

	aHist.MergeFrom(bHist)

	// aHist and cHist should be equivalent
	requireEqual(t, cHist, aHist)
}

// Tests the logic to switch between uint8, uint16, uint32, and
// uint64.  Test is based on the UpdateByIncr code path.
func TestOverflowBits(t *testing.T) {
	for _, limit := range []uint64{
		0x100,
		0x10000,
		0x100000000,
	} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			aHist := NewFloat64(NewConfig())
			bHist := NewFloat64(NewConfig())
			cHist := NewFloat64(NewConfig())

			if limit <= 0x10000 {
				for i := uint64(0); i < limit; i++ {
					aHist.Update(plusOne)
					aHist.Update(minusOne)

					require.Equal(t, 2*(i+1), aHist.Count())
				}
			} else {
				aHist.UpdateByIncr(plusOne, limit/2)
				aHist.UpdateByIncr(plusOne, limit/2)
				aHist.UpdateByIncr(minusOne, limit/2)
				aHist.UpdateByIncr(minusOne, limit/2)
			}
			bHist.UpdateByIncr(plusOne, limit-1)
			bHist.Update(plusOne)
			bHist.UpdateByIncr(minusOne, limit-1)
			bHist.Update(minusOne)
			cHist.UpdateByIncr(plusOne, limit)
			cHist.UpdateByIncr(minusOne, limit)

			require.Equal(t, 2*limit, aHist.Count())
			require.Equal(t, float64(0), aHist.Sum())

			aPos := aHist.Positive()
			require.Equal(t, uint32(1), aPos.Len())
			require.Equal(t, limit, aPos.At(0))

			aNeg := aHist.Negative()
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
	agg := NewInt64(NewConfig(WithMaxSize(256)))
	alt := NewInt64(NewConfig(WithMaxSize(256)))

	expect := int64(0)
	for i := int64(1); i < 256; i++ {
		expect += i
		agg.Update(i)
		alt.Update(i)
	}

	require.Equal(t, expect, agg.Sum())
	require.Equal(t, uint64(255), agg.Count())

	// Scale should be 5.  Here's why.  The upper power-of-two is
	// 256 == 2**8.  We expect the exponential base = 2**(2**-5)
	// raised to the 256th power to equal 256:
	//
	//   2**((2**-5)*256)
	// = 2**((2**-5)*(2**8))
	// = 2**(2**3)
	// = 2**8
	require.Equal(t, int32(5), agg.Scale())

	expect0 := func(b *Buckets) {
		require.Equal(t, uint32(0), b.Len())
	}
	expect256 := func(b *Buckets, factor int) {
		require.Equal(t, uint32(256), b.Len())
		require.Equal(t, int32(0), b.Offset())
		// Bucket 254 has 6 elements, bucket 255 has 5
		// bucket 253 has 5, ...
		for i := uint32(0); i < 256; i++ {
			require.LessOrEqual(t, b.At(i), uint64(6*factor))
		}
	}

	expect256(agg.Positive(), 1)
	expect0(agg.Negative())

	// Merge!
	agg.MergeFrom(alt)

	expect256(agg.Positive(), 2)

	require.Equal(t, 2*expect, agg.Sum())

	// Reset!  Repeat with negative.
	agg.Clear()
	alt.Clear()

	expect = int64(0)
	for i := int64(1); i < 256; i++ {
		expect -= i
		agg.Update(-i)
		alt.Update(-i)
	}

	require.Equal(t, expect, agg.Sum())
	require.Equal(t, uint64(255), agg.Count())

	expect256(agg.Negative(), 1)
	expect0(agg.Positive())

	// Merge!
	agg.MergeFrom(alt)

	expect256(agg.Negative(), 2)

	require.Equal(t, 2*expect, agg.Sum())

	// Scale should not change after filling in the negative range.
	require.Equal(t, int32(5), agg.Scale())
}

// Tests the reset code path via MoveInto.
func TestReset(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(256)))

	for _, incr := range []uint64{
		1,
		0x100,
		0x10000,
		0x100000000,

		// Another 32-bit increment tests the 64-bit reset path.
		0x200000000,
	} {
		t.Run(fmt.Sprint(incr), func(t *testing.T) {
			agg.Clear()

			// Note that scale is zero b/c no values
			require.Equal(t, int32(0), agg.Scale())

			expect := 0.0
			for i := int64(1); i < 256; i++ {
				expect += float64(i) * float64(incr)
				agg.UpdateByIncr(float64(i), incr)
			}

			require.Equal(t, expect, agg.Sum())
			require.Equal(t, uint64(255)*incr, agg.Count())

			// See TestIntegerAggregation about why scale is 5.
			require.Equal(t, int32(5), agg.Scale())

			pos := agg.Positive()

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

// Tests the swap operation.
func TestMoveInto(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(256)))
	cpy := NewFloat64(NewConfig(WithMaxSize(256)))

	expect := 0.0
	for i := int64(1); i < 256; i++ {
		expect += float64(i)
		agg.Update(float64(i))
		agg.Update(0)
	}

	agg.Swap(cpy)

	// agg was reset
	require.Equal(t, 0.0, agg.Sum())
	require.Equal(t, uint64(0), agg.Count())
	require.Equal(t, uint64(0), agg.ZeroCount())
	require.Equal(t, int32(0), agg.Scale())

	// cpy is as expected
	require.Equal(t, expect, cpy.Sum())
	require.Equal(t, uint64(255*2), cpy.Count())
	require.Equal(t, uint64(255), cpy.ZeroCount())

	// See TestIntegerAggregation about why scale is 5,
	// max bucket count is 6, and so on.
	require.Equal(t, int32(5), cpy.Scale())

	pos := cpy.Positive()

	require.Equal(t, uint32(256), pos.Len())
	require.Equal(t, int32(0), pos.Offset())
	for i := uint32(0); i < 256; i++ {
		require.LessOrEqual(t, pos.At(i), uint64(6))
	}
}

// Tests with maxSize=2 that very large numbers (but not the full
// range) yield scales -7 and -8.
func TestVeryLargeNumbers(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(2)))

	expectBalanced := func(c uint64) {
		pos := agg.Positive()
		require.Equal(t, uint32(2), pos.Len())
		require.Equal(t, int32(-1), pos.Offset())
		require.Equal(t, c, pos.At(0))
		require.Equal(t, c, pos.At(1))
	}

	agg.Update(0x1p-100)
	agg.Update(0x1p+100)

	require.InEpsilon(t, 0x1p100, agg.Sum(), 1e-5)
	require.Equal(t, uint64(2), agg.Count())
	require.Equal(t, int32(-7), agg.Scale())

	expectBalanced(1)

	agg.Update(0x1p-128)
	agg.Update(0x1p+127)

	require.InEpsilon(t, 0x1p127, agg.Sum(), 1e-5)
	require.Equal(t, uint64(4), agg.Count())
	require.Equal(t, int32(-7), agg.Scale())

	expectBalanced(2)

	agg.Update(0x1p-129)
	agg.Update(0x1p+255)

	require.InEpsilon(t, 0x1p255, agg.Sum(), 1e-5)
	require.Equal(t, uint64(6), agg.Count())
	require.Equal(t, int32(-8), agg.Scale())

	expectBalanced(3)
}

// Tests the largest and smallest finite numbers with below-minimum
// size.  Expect a size=MinSize histogram with MinScale.
func TestFullRange(t *testing.T) {
	agg := NewFloat64(NewConfig(WithMaxSize(1)))

	agg.Update(math.MaxFloat64)
	agg.Update(1)
	agg.Update(math.SmallestNonzeroFloat64)

	require.Equal(t, logarithm.MaxValue, agg.Sum())
	require.Equal(t, uint64(3), agg.Count())

	require.Equal(t, exponent.MinScale, agg.Scale())

	pos := agg.Positive()

	require.Equal(t, uint32(MinSize), pos.Len())
	require.Equal(t, int32(-1), pos.Offset())
	require.Equal(t, pos.At(0), uint64(1))
	require.Equal(t, pos.At(1), uint64(2))
}

// TestAggregatorMinMax verifies the min and max values.
func TestAggregatorMinMax(t *testing.T) {
	h1 := NewFloat64(NewConfig(), 1, 3, 5, 7, 9)
	require.Equal(t, 1.0, h1.Min())
	require.Equal(t, 9.0, h1.Max())

	h2 := NewFloat64(NewConfig(), -1, -3, -5, -7, -9)
	require.Equal(t, -9.0, h2.Min())
	require.Equal(t, -1.0, h2.Max())
}

// TestAggregatorCopySwap tests both Copy and Swap.
func TestAggregatorCopySwap(t *testing.T) {
	h1 := NewFloat64(NewConfig(), 1, 3, 5, 7, 9, -1, -3, -5)
	h2 := NewFloat64(NewConfig(), 5, 4, 3, 2)
	h3 := NewFloat64(NewConfig())

	h1.Swap(h2)
	h2.CopyInto(h3)

	requireEqual(t, h2, h3)
}

// Benchmarks the Update() function for values in the range [1,2)
func BenchmarkLinear(b *testing.B) {
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := NewFloat64(NewConfig(WithMaxSize(1024)))
	for i := 0; i < b.N; i++ {
		x := 2 - rnd.Float64()
		agg.Update(x)
	}
}

// Benchmarks the Update() function for values in the range (0, MaxValue]
func BenchmarkExponential(b *testing.B) {
	src := rand.NewSource(77777677777)
	rnd := rand.New(src)
	agg := NewFloat64(NewConfig(WithMaxSize(1024)))
	for i := 0; i < b.N; i++ {
		x := rnd.ExpFloat64()
		agg.Update(x)
	}
}

func benchmarkMapping(b *testing.B, name string, mapper mapping.Mapping) {
	b.Run(fmt.Sprintf("mapping_%s", name), func(b *testing.B) {
		src := rand.New(rand.NewSource(54979))

		for i := 0; i < b.N; i++ {
			_ = mapper.MapToIndex(1 + src.Float64())
		}
	})
}

func benchmarkBoundary(b *testing.B, name string, mapper mapping.Mapping) {
	b.Run(fmt.Sprintf("boundary_%s", name), func(b *testing.B) {
		src := rand.New(rand.NewSource(54979))

		for i := 0; i < b.N; i++ {
			_, _ = mapper.LowerBoundary(int32(src.Int63()))
		}
	})
}

// An earlier draft of this benchmark included a lookup-table based
// implementation:
// https://github.com/open-telemetry/opentelemetry-go-contrib/pull/1353
// That mapping function uses O(2^scale) extra space and falls
// somewhere between the exponent and logarithm methods compared here.
// In the test, lookuptable was 40% faster than logarithm, which did
// not justify the significant extra complexity.

// Benchmarks the MapToIndex function.
func BenchmarkMapping(b *testing.B) {
	em, _ := exponent.NewMapping(-1)
	lm, _ := logarithm.NewMapping(1)
	benchmarkMapping(b, "exponent", em)
	benchmarkMapping(b, "logarithm", lm)
}

// Benchmarks the LowerBoundary function.
func BenchmarkReverseMapping(b *testing.B) {
	em, _ := exponent.NewMapping(-1)
	lm, _ := logarithm.NewMapping(1)
	benchmarkBoundary(b, "exponent", em)
	benchmarkBoundary(b, "logarithm", lm)
}

// Statistical test: how biased are the exact power-of-two boundaries?
func TestBoundaryStatistics(t *testing.T) {
	for scale := logarithm.MinScale; scale <= logarithm.MaxScale; scale++ {

		m, _ := logarithm.NewMapping(scale)

		var above, below int

		total := exponent.MaxNormalExponent - exponent.MinNormalExponent + 1
		for exp := exponent.MinNormalExponent; exp <= exponent.MaxNormalExponent; exp++ {
			value := math.Ldexp(1, int(exp))

			index := m.MapToIndex(value)

			bound, err := m.LowerBoundary(index)
			require.NoError(t, err)

			if bound == value {
			} else if bound < value {
				above++
			} else {
				below++
			}
		}

		// The sample results here not guaranteed.  Test that this is approximately unbiased.
		// (Results on dev machine: 1059 above, 963 below, 24 equal, total = 2046.)
		require.InEpsilon(t, 0.5, float64(above)/float64(total), 0.05)
		require.InEpsilon(t, 0.5, float64(below)/float64(total), 0.06)
	}
}
