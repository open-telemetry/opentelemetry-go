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
	normalMapping = newLogarithmMapping(30)
	oneAndAHalf   = centerVal(normalMapping, int32(normalMapping.MapToIndex(1.5)))
)

func centerVal(mapper logarithmMapping, x int32) float64 {
	return (mapper.LowerBoundary(int64(x)) + mapper.LowerBoundary(int64(x)+1)) / 2
}

// tests a simple case of 8 counts entered into a maxSize=4 histogram,
// causing a single downscale and no rotation.
func TestSimpleSize4(t *testing.T) {
	// Test with a 4-bucket-max exponential histogram
	ctx := context.Background()
	desc := metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
	agg := New(1, &desc, WithMaxSize(4))[0]
	pos := agg.Positive()
	neg := agg.Negative()

	// Add a zero
	agg.Update(ctx, 0, &desc)

	require.Equal(t, uint64(1), agg.ZeroCount())
	require.Equal(t, uint32(0), pos.Len())
	require.Equal(t, uint32(0), neg.Len())

	// Add a oneAndAHalf, which is in the normal range => DefaultNormalScale.
	agg.Update(ctx, number.NewFloat64Number(oneAndAHalf), &desc)

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
		agg.Update(ctx, number.NewFloat64Number(value), &desc)
	}

	require.Equal(t, uint32(4), pos.Len())

	for i := uint32(0); i < 4; i++ {
		require.Equal(t, uint64(1), pos.At(i))
	}

	// Add the next value!
	for i := int32(4); i < 8; i++ {
		value := centerVal(mapper, int32(offset)+i)
		agg.Update(ctx, number.NewFloat64Number(value), &desc)
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
			desc := metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
			agg := New(1, &desc, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &desc)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &desc)

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// (1/2), counts[1] == 2 (the 1 and 2).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &desc)
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
			desc := metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
			agg := New(1, &desc, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &desc)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &desc)

			// Enter order[2]: scale set to -1, expect counts[0] == 1 (the
			// 1 and 2), counts[1] == 2 (the 4).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &desc)
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
			desc := metrictest.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
			agg := New(1, &desc, WithMaxSize(2))[0]

			agg.Update(ctx, number.NewFloat64Number(order[0]), &desc)
			agg.Update(ctx, number.NewFloat64Number(order[1]), &desc)

			// Enter order[2]: scale set to -1, expect counts[0] == 2 (the
			// 1/4 and 1/2, counts[1] == 2 (the 1).
			agg.Update(ctx, number.NewFloat64Number(order[2]), &desc)
			require.Equal(t, int32(-1), agg.Scale())
			require.Equal(t, int32(-1), agg.Positive().Offset())
			require.Equal(t, uint32(2), agg.Positive().Len())
			require.Equal(t, uint64(2), agg.Positive().At(0))
			require.Equal(t, uint64(1), agg.Positive().At(1))
		})
	}
}
