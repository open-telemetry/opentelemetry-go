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
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/logarithm"
)

var (
	normalMapping  = logarithm.NewMapping(30)
	oneAndAHalf    = centerVal(normalMapping, int32(normalMapping.MapToIndex(1.5)))
	testDescriptor = metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
)

// TEST SUPPORT

func requireEqual(t *testing.T, a, b *Aggregator) {
	require.InEpsilon(t, a.state.sum, b.state.sum, 1e-10)
	require.Equal(t, a.state.count, b.state.count)
	require.Equal(t, a.state.zeroCount, b.state.zeroCount)
	require.Equal(t, a.state.mapping.Scale(), b.state.mapping.Scale())

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
	lb, err1 := mapper.LowerBoundary(int64(x))
	ub, err2 := mapper.LowerBoundary(int64(x) + 1)
	if err1 != nil || err2 != nil {
		panic(fmt.Sprintf("unexpected errors: %v %v", err1, err2))
	}
	return (lb + ub) / 2
}

// tests a simple case of 8 counts entered into a maxSize=4 histogram,
// causing a single downscale and no rotation.
func TestSimpleSize4(t *testing.T) {
	// Test with a 4-bucket-max exponential histogram
	ctx := context.Background()
	agg := New(1, &testDescriptor, WithMaxSize(4))[0]
	pos := agg.positive()
	neg := agg.negative()

	// Add a zero
	agg.Update(ctx, 0, &testDescriptor)

	require.Equal(t, uint64(1), agg.zeroCount())
	require.Equal(t, uint32(0), pos.Len())
	require.Equal(t, uint32(0), neg.Len())

	// Add a oneAndAHalf, which is in the normal range => DefaultNormalScale.
	agg.Update(ctx, number.NewFloat64Number(oneAndAHalf), &testDescriptor)

	// Test count=2
	cnt, err := agg.Count()
	require.Equal(t, uint64(2), cnt)
	require.NoError(t, err)

	// Test sum=oneAndAHalf
	sum, err := agg.Sum()
	require.Equal(t, number.NewFloat64Number(oneAndAHalf), sum)
	require.NoError(t, err)

	// Test a single positive bucket with count 1 at DefaultNormalScale.
	require.Equal(t, uint32(1), pos.Len())
	require.Equal(t, uint32(0), neg.Len())
	require.Equal(t, uint64(1), pos.At(0))
	require.Equal(t, int32(DefaultNormalScale), agg.scale())

	mapper := logarithm.NewMapping(agg.scale())

	// Check that the initial count maps to Offset().
	offset := mapper.MapToIndex(oneAndAHalf)
	require.Equal(t, int32(offset), pos.Offset())

	// Add 3 more values in each of the subsequent buckets.
	for i := int32(1); i < 4; i++ {
		value := centerVal(mapper, int32(offset)+i)
		agg.Update(ctx, number.NewFloat64Number(value), &testDescriptor)
	}

	require.Equal(t, uint32(4), pos.Len())

	for i := uint32(0); i < 4; i++ {
		require.Equal(t, uint64(1), pos.At(i))
	}

	// Add the next value!
	for i := int32(4); i < 8; i++ {
		value := centerVal(mapper, int32(offset)+i)
		agg.Update(ctx, number.NewFloat64Number(value), &testDescriptor)
	}

	// Expect 2 in each bucket
	for i := uint32(0); i < 4; i++ {
		require.Equal(t, uint64(2), pos.At(i))
	}
}

func TestAlternatingGrowth1(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	agg.Update(ctx, number.NewFloat64Number(1), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(2), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(0.5), &testDescriptor)

	require.Equal(t, int32(-1), agg.positive().Offset())
	require.Equal(t, int32(0), agg.scale())
	require.Equal(t, []uint64{1, 1, 1}, counts(agg.positive()))
}

func TestAlternatingGrowth2(t *testing.T) {
	ctx := context.Background()
	agg := &New(1, &testDescriptor, WithMaxSize(4))[0]
	agg.Update(ctx, number.NewFloat64Number(1), &testDescriptor)
	agg.Update(ctx, number.NewFloat64Number(1), &testDescriptor)
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
			require.Equal(t, int32(mapper.MapToIndex(minVal)), agg.positive().Offset())

			// The maximum range is correct at the computed scale.
			require.Equal(t, int32(mapper.MapToIndex(maxVal)), agg.positive().Offset()+int32(agg.positive().Len())-1)
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
		repeat = 64
		count  = 64
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
										count / 4,
										count / 2,
										count,
									} {
										t.Run(fmt.Sprint("size=", size), func(t *testing.T) {
											testMergeExhaustive(t, values[0:part], values[part:count], size)
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

func testMergeExhaustive(t *testing.T, a, b []float64, size int32) {
	ctx := context.Background()
	aggs := New(3, &testDescriptor, WithMaxSize(size))

	aHist := &aggs[0]
	bHist := &aggs[1]
	cHist := &aggs[2]

	for _, av := range a {
		aHist.Update(ctx, number.NewFloat64Number(av), &testDescriptor)
		cHist.Update(ctx, number.NewFloat64Number(av), &testDescriptor)
	}
	for _, bv := range b {
		bHist.Update(ctx, number.NewFloat64Number(bv), &testDescriptor)
		cHist.Update(ctx, number.NewFloat64Number(bv), &testDescriptor)
	}
	aHist.Merge(bHist, &testDescriptor)

	// aHist and cHist should be equivalent
	requireEqual(t, cHist, aHist)
}
