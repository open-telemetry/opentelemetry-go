package exponential

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping/logarithm"
)

var (
	testDescriptor = metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
	intDescriptor  = metrictest.NewDescriptor("integer", sdkapi.HistogramInstrumentKind, number.Int64Kind)

	plusOne  = number.NewFloat64Number(1)
	minusOne = number.NewFloat64Number(-1)
)

type show struct {
	index int32
	count uint64
	lower float64
}

func (a *Aggregator) shows(b *buckets) (r []show) {
	for i := uint32(0); i < b.Len(); i++ {
		lower, _ := a.state.mapping.LowerBoundary(b.Offset() + int32(i))
		r = append(r, show{
			index: b.Offset() + int32(i),
			count: b.At(i),
			lower: lower,
		})
	}
	return r
}

func counts(b *buckets) (r []uint64) {
	for i := uint32(0); i < b.Len(); i++ {
		r = append(r, b.At(i))
	}
	return r
}

func (s show) String() string {
	return fmt.Sprintf("%v=%v(%.2g)", s.index, s.count, s.lower)
}

func (a *Aggregator) String() string {
	return fmt.Sprintf("%v %v\n%v\n%v", a.state.count, a.state.sum, a.shows(&a.state.positive), a.shows(&a.state.negative))
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
	require.Equal(t, bstr(&a.state.positive), bstr(&b.state.positive), "positive %v %v", a.shows(&a.state.positive), a.shows(&b.state.positive))
	require.Equal(t, bstr(&a.state.negative), bstr(&b.state.negative), "negative %v %v", a.shows(&a.state.negative), a.shows(&b.state.negative))
}

func centerVal(mapper mapping.Mapping, x int32) float64 {
	lb, err1 := mapper.LowerBoundary(x)
	ub, err2 := mapper.LowerBoundary(x + 1)
	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("unexpected errors: %v %v", err1, err2))
	}
	return (lb + ub) / 2
}

func TestAlternatingGrowth1(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	agg.Update(ctx, plusOne, &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(2), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(0.5), &testDescriptor)

	require.Equal(t, int32(-1), agg.positive().Offset())
	require.Equal(t, int32(0), agg.scale())
	require.Equal(t, []uint64{1, 1, 1}, counts(agg.positive()))
}

func TestAlternatingGrowth2(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	agg.Update(ctx, plusOne, &testDescriptor)
	agg.Update(ctx, plusOne, &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(2), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(0.5), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(4), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(0.25), &testDescriptor)

	require.Equal(t, int32(-1), agg.positive().Offset())
	require.Equal(t, int32(-1), agg.scale())
	require.Equal(t, []uint64{2, 3, 1}, counts(agg.positive()))
}

// tests that every permutation of {1/2, 1, 2} with maxSize=2 results
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
			agg := New(1, &testDescriptor, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor)

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// (1/2), counts[1] == 2 (the 1 and 2).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor)
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(-1), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(1), agg.positive().At(0))
			require.Equal(t, uint64(2), agg.positive().At(1))
		})
	}
}

// tests that every permutation of {1, 2, 4} with maxSize=2 results
// in the same scale=-1 histogram.
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
			agg := New(1, &testDescriptor, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor)

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// 1 and 2), counts[1] == 2 (the 4).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor)
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(0), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(2), agg.positive().At(0))
			require.Equal(t, uint64(1), agg.positive().At(1))
		})
	}
}

// tests that every permutation of {1, 1/2, 1/4} with maxSize=2 results
// in the same scale=-1 histogram.
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
			agg := New(1, &testDescriptor, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &testDescriptor)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &testDescriptor)

			// Enter order[2]: scale set to -1, expect counts[0] == 2 (the
			// 1/4 and 1/2, counts[1] == 2 (the 1).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &testDescriptor)
			require.Equal(t, int32(-1), agg.scale())
			require.Equal(t, int32(-1), agg.positive().Offset())
			require.Equal(t, uint32(2), agg.positive().Len())
			require.Equal(t, uint64(2), agg.positive().At(0))
			require.Equal(t, uint64(1), agg.positive().At(1))
		})
	}
}

func TestExhaustiveSmall(t *testing.T) {
	for _, maxSize := range []int32{3, 4, 5, 6, 7, 8, 9} {
		t.Run(fmt.Sprintf("maxSize=%d", maxSize), func(t *testing.T) {
			for offset := int32(-5); offset <= 5; offset++ {
				t.Run(fmt.Sprintf("offset=%d", offset), func(t *testing.T) {
					for _, initScale := range []int32{
						0, 1, 2, 3, 4,
					} {
						t.Run(fmt.Sprintf("initScale=%d", initScale), func(t *testing.T) {
							testExhaustive(t, maxSize, offset, initScale)
						})
					}
				})
			}
		})
	}
}

func testExhaustive(t *testing.T, maxSize, offset, initScale int32) {
	for step := maxSize; step < 4*maxSize; step++ {
		t.Run(fmt.Sprintf("step=%d", step), func(t *testing.T) {
			ctx := context.Background()
			agg := New(1, &testDescriptor, WithMaxSize(maxSize))[0]
			mapper := newMapping(initScale)

			minVal := centerVal(mapper, offset)
			maxVal := centerVal(mapper, offset+step)
			sum := 0.0

			for i := int32(0); i < maxSize; i++ {
				value := centerVal(mapper, offset+i)
				agg.Update(ctx, number.NewFloat64Number(value), &testDescriptor)
				sum += value
			}

			require.Equal(t, initScale, agg.scale())
			require.Equal(t, offset, agg.positive().Offset())

			agg.Update(ctx, number.NewFloat64Number(maxVal), &testDescriptor)
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
			mapper = newMapping(agg.scale())
			idx, err := mapper.MapToIndex(minVal)
			require.NoError(t, err)
			require.Equal(t, int32(idx), agg.positive().Offset())

			// The maximum range is correct at the computed scale.
			idx, err = mapper.MapToIndex(maxVal)
			require.NoError(t, err)
			require.Equal(t, int32(idx), agg.positive().Offset()+int32(agg.positive().Len())-1)
		})
	}
}

func TestMergeSimpleEven(t *testing.T) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(4))
	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)   // 1, 2, 4, 8
		f2 := 1 / float64(int64(2)<<i) // 1/2, 1/4, 1/8, 1/16
		n1 := number.NewFloat64Number(f1)
		n2 := number.NewFloat64Number(f2)

		aggs[0].Update(ctx, n1, &testDescriptor)
		aggs[1].Update(ctx, n2, &testDescriptor)
		aggs[2].Update(ctx, n1, &testDescriptor)
		aggs[2].Update(ctx, n2, &testDescriptor)
	}
	require.Equal(t, int32(0), aggs[0].scale())
	require.Equal(t, int32(0), aggs[1].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	require.Equal(t, int32(0), aggs[0].positive().Offset())
	require.Equal(t, int32(-4), aggs[1].positive().Offset())
	require.Equal(t, int32(-2), aggs[2].positive().Offset())

	require.Equal(t, []uint64{1, 1, 1, 1}, counts(aggs[0].positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, counts(aggs[1].positive()))
	require.Equal(t, []uint64{2, 2, 2, 2}, counts(aggs[2].positive()))

	require.NoError(t, aggs[0].Merge(&aggs[1], &testDescriptor))

	require.Equal(t, int32(-1), aggs[0].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	requireEqual(t, &aggs[0], &aggs[2])
}

func TestMergeSimpleOdd(t *testing.T) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(4))
	for i := 0; i < 4; i++ {
		f1 := float64(int64(1) << i)
		f2 := 1 / float64(int64(1)<<i) // Diff from above test: 1 here vs 2 above.
		n1 := number.NewFloat64Number(f1)
		n2 := number.NewFloat64Number(f2)

		aggs[0].Update(ctx, n1, &testDescriptor)
		aggs[1].Update(ctx, n2, &testDescriptor)
		aggs[2].Update(ctx, n1, &testDescriptor)
		aggs[2].Update(ctx, n2, &testDescriptor)
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

	require.Equal(t, []uint64{1, 1, 1, 1}, counts(aggs[0].positive()))
	require.Equal(t, []uint64{1, 1, 1, 1}, counts(aggs[1].positive()))
	require.Equal(t, []uint64{1, 2, 3, 2}, counts(aggs[2].positive()))

	require.NoError(t, aggs[0].Merge(&aggs[1], &testDescriptor))

	require.Equal(t, int32(-1), aggs[0].scale())
	require.Equal(t, int32(-1), aggs[2].scale())

	requireEqual(t, &aggs[0], &aggs[2])
}

func TestMergeExhaustive(t *testing.T) {
	const (
		factor = 1024.0
		repeat = 16
		count  = 32
	)

	means := []float64{
		0,
		1,
		factor - 1,
		factor,
		factor + 1,
		factor*factor - factor,
		factor * factor,
		factor*factor + factor,
	}

	stddevs := []float64{
		1,
		factor,
		factor * factor,
	}

	for _, mean := range means {
		t.Run(fmt.Sprint("mean=", mean), func(t *testing.T) {
			for _, stddev := range stddevs {
				t.Run(fmt.Sprint("stddev=", stddev), func(t *testing.T) {
					for r := 0; r < repeat; r++ {
						src := rand.NewSource(77777677777)
						rnd := rand.New(src)

						t.Run(fmt.Sprint("repeat=", r), func(t *testing.T) {
							values := make([]float64, count)
							for i := range values {
								values[i] = mean + rnd.NormFloat64()*stddev
							}

							for part := 1; part < count; part++ {
								t.Run(fmt.Sprint("part=", part), func(t *testing.T) {
									for _, size := range []int32{
										2,
										3,
										4,
										6,
										9,
										12,
										16,
									} {
										t.Run(fmt.Sprint("size=", size), func(t *testing.T) {
											for _, incr := range []uint64{
												1,
												17,
												0x100,
												0x10000,
												0x100000000,
											} {
												t.Run(fmt.Sprintf("incr=%x", incr), func(t *testing.T) {
													testMergeExhaustive(t, values[0:part], values[part:count], size, incr)
												})
											}
										})
									}
								})
							}
						})
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
		aHist.UpdateByIncr(ctx, number.NewFloat64Number(av), incr, &testDescriptor)
		cHist.UpdateByIncr(ctx, number.NewFloat64Number(av), incr, &testDescriptor)
	}
	for _, bv := range b {
		bHist.UpdateByIncr(ctx, number.NewFloat64Number(bv), incr, &testDescriptor)
		cHist.UpdateByIncr(ctx, number.NewFloat64Number(bv), incr, &testDescriptor)
	}

	aHist.Merge(bHist, &testDescriptor)

	// aHist and cHist should be equivalent
	requireEqual(t, cHist, aHist)
}

func TestOverflow8bits(t *testing.T) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor)

	aHist := &aggs[0]
	bHist := &aggs[1]
	cHist := &aggs[2]

	for i := 0; i < 256; i++ {
		aHist.Update(ctx, plusOne, &testDescriptor)
	}
	bHist.UpdateByIncr(ctx, plusOne, 255, &testDescriptor)
	bHist.Update(ctx, plusOne, &testDescriptor)
	cHist.UpdateByIncr(ctx, plusOne, 256, &testDescriptor)

	requireEqual(t, cHist, aHist)
	requireEqual(t, bHist, aHist)
}

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
					aHist.Update(ctx, plusOne, &testDescriptor)
					aHist.Update(ctx, minusOne, &testDescriptor)

					cnt, _ := aHist.Count()
					require.Equal(t, 2*(i+1), cnt)
				}
			} else {
				aHist.UpdateByIncr(ctx, plusOne, limit/2, &testDescriptor)
				aHist.UpdateByIncr(ctx, plusOne, limit/2, &testDescriptor)
				aHist.UpdateByIncr(ctx, minusOne, limit/2, &testDescriptor)
				aHist.UpdateByIncr(ctx, minusOne, limit/2, &testDescriptor)
			}
			bHist.UpdateByIncr(ctx, plusOne, limit-1, &testDescriptor)
			bHist.Update(ctx, plusOne, &testDescriptor)
			bHist.UpdateByIncr(ctx, minusOne, limit-1, &testDescriptor)
			bHist.Update(ctx, minusOne, &testDescriptor)
			cHist.UpdateByIncr(ctx, plusOne, limit, &testDescriptor)
			cHist.UpdateByIncr(ctx, minusOne, limit, &testDescriptor)

			aCnt, _ := aHist.Count()
			require.Equal(t, 2*limit, aCnt)
			sum, _ := aHist.Sum()
			require.Equal(t, float64(0), sum.CoerceToFloat64(number.Float64Kind))

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

func TestIntegerAggregation(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &intDescriptor, WithMaxSize(256))[0]

	expect := int64(0)
	for i := int64(1); i < 256; i++ {
		expect += i
		agg.Update(ctx, number.NewInt64Number(i), &intDescriptor)
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
	expect256 := func(b aggregation.ExponentialBuckets) {
		require.Equal(t, uint32(256), b.Len())
		require.Equal(t, int32(0), b.Offset())
		// Bucket 254 has 6 elements, bucket 255 has 5
		// bucket 253 has 5, ...
		for i := uint32(0); i < 256; i++ {
			require.LessOrEqual(t, b.At(i), uint64(6))
		}
	}

	pos, err := agg.Positive()
	require.NoError(t, err)
	expect256(pos)

	neg, err := agg.Negative()
	require.NoError(t, err)
	expect0(neg)

	// Reset!
	agg.SynchronizedMove(nil, &intDescriptor)

	expect = int64(0)
	for i := int64(1); i < 256; i++ {
		expect -= i
		agg.Update(ctx, number.NewInt64Number(-i), &intDescriptor)
	}

	require.Equal(t, expect, intSumNoError(t, agg))
	require.Equal(t, uint64(255), countNoError(t, agg))

	neg, err = agg.Negative()
	require.NoError(t, err)
	expect256(neg)

	pos, err = agg.Positive()
	require.NoError(t, err)
	expect0(pos)

	// Scale should not change after filling in the negative range.
	require.Equal(t, int32(5), scaleNoError(t, agg))
}

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
			agg.SynchronizedMove(nil, &testDescriptor)

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
	require.Equal(t, logarithm.MaxScale, scaleNoError(t, agg))

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

// func TestFixedLimits(t *testing.T) {
// 	const min = 0.001
// 	const max = 60000
// 	for scale := int32(0); scale < 10; scale++ {
// 		m := newMapping(scale)
// 		fmt.Println("required size at scale", scale, "is", m.MapToIndex(max)-m.MapToIndex(min))
// 	}
// }
