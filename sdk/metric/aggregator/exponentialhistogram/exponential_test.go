package exponential

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
)

var (
	normalMapping = NewLogarithmMapping(30)
	oneAndAHalf   = centerVal(normalMapping, int32(normalMapping.MapToIndex(1.5)))
)

func centerVal(mapper LogarithmMapping, x int32) float64 {
	return (mapper.LowerBoundary(int64(x)) + mapper.LowerBoundary(int64(x)+1)) / 2
}

func TestInitialCondition(t *testing.T) {
	// Test with a 4-bucket-max exponential histogram
	ctx := context.Background()
	desc := metric.NewDescriptor("name", sdkapi.HistogramInstrumentKind, number.Float64Kind)
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

	mapper := NewLogarithmMapping(agg.Scale())

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

	// Expect 2 in each bucket.
	for i := uint32(0); i < 4; i++ {
		require.Equal(t, uint64(2), pos.At(i))
	}
}
