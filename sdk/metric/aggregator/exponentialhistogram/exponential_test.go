package exponential

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

var (
	normalMapping  = newLogarithmMapping(30)
	oneAndAHalf    = centerVal(normalMapping, int32(normalMapping.MapToIndex(1.5)))
	testDescriptor = metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
)

func centerVal(mapper Mapping, x int32) float64 {
	return (mapper.LowerBoundary(int64(x)) + mapper.LowerBoundary(int64(x)+1)) / 2
}

// tests a simple case of 8 counts entered into a maxSize=4 histogram,
// causing a single downscale and no rotation.
func TestSimpleSize4(t *testing.T) {
	// Test with a 4-bucket-max exponential histogram
	ctx := context.Background()
	agg := New(1, &testDescriptor, WithMaxSize(4))[0]
	pos := agg.Positive()
	neg := agg.Negative()

	// Add a zero
	agg.Update(ctx, 0, &testDescriptor)

	require.Equal(t, uint64(1), agg.ZeroCount())
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
	require.Equal(t, int32(DefaultNormalScale), agg.Scale())

	mapper := newLogarithmMapping(agg.Scale())

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
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(-1), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(1), agg.Positive().At(0))
			require.Equal(t, uint64(2), agg.Positive().At(1))
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
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(0), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(2), agg.Positive().At(0))
			require.Equal(t, uint64(1), agg.Positive().At(1))
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
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(-1), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(2), agg.Positive().At(0))
			require.Equal(t, uint64(1), agg.Positive().At(1))
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

			require.Equal(t, initScale, agg.Scale())
			require.Equal(t, offset, agg.Positive().Offset())

			agg.Update(ctx, number.NewFloat64Number(maxVal), &testDescriptor)
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
			hcount, _ := agg.Count()
			require.GreaterOrEqual(t, uint64(maxSize+1), hcount)
			// Sum is correct
			hsum, _ := agg.Sum()
			require.GreaterOrEqual(t, sum, hsum.CoerceToFloat64(number.Float64Kind))

			// The offset is correct at the computed scale.
			mapper = newMapping(agg.Scale())
			require.Equal(t, int32(mapper.MapToIndex(minVal)), agg.Positive().Offset())

			// The maximum range is correct at the computed scale.
			require.Equal(t, int32(mapper.MapToIndex(maxVal)), agg.Positive().Offset()+int32(agg.Positive().Len())-1)
		})
	}
}
