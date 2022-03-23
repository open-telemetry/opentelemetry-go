// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exponent

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
)

type expectMapping struct {
	value float64
	index int32
}

// Tests that getBase2 returns the base-2 exponent as documented, unlike
// math.Frexp.
func TestGetBase2(t *testing.T) {
	require.Equal(t, int32(-1022), MinNormalExponent)
	require.Equal(t, int32(+1023), MaxNormalExponent)

	require.Equal(t, MaxNormalExponent, getBase2(0x1p+1023))
	require.Equal(t, int32(1022), getBase2(0x1p+1022))

	require.Equal(t, int32(0), getBase2(1))

	require.Equal(t, int32(-1021), getBase2(0x1p-1021))
	require.Equal(t, int32(-1022), getBase2(0x1p-1022))

	// Subnormals below this point
	require.Equal(t, int32(-1022), getBase2(0x1p-1023))
	require.Equal(t, int32(-1022), getBase2(0x1p-1024))
	require.Equal(t, int32(-1022), getBase2(0x1p-1025))
	require.Equal(t, int32(-1022), getBase2(0x1p-1074))
}

// Tests a few cases with scale=0.
func TestExponentMappingZero(t *testing.T) {
	m, err := NewMapping(0)
	require.NoError(t, err)

	require.Equal(t, int32(0), m.Scale())

	for _, pair := range []expectMapping{
		{math.MaxFloat64, MaxNormalExponent},
		{0x1p+1023, MaxNormalExponent},
		{0x1p-1022, MinNormalExponent},
		{math.SmallestNonzeroFloat64, MinNormalExponent},
		{4, 2},
		{3, 1},
		{2, 1},
		{1.5, 0},
		{1, 0},
		{0.75, -1},
		{0.5, -1},
		{0.25, -2},
	} {
		idx := m.MapToIndex(pair.value)

		require.Equal(t, pair.index, idx)
	}
}

// Tests a few cases with scale=MinScale.
func TestExponentMappingMinScale(t *testing.T) {
	m, err := NewMapping(MinScale)
	require.NoError(t, err)

	require.Equal(t, MinScale, m.Scale())

	for _, pair := range []expectMapping{
		{1, 0},
		{math.MaxFloat64 / 2, 0},
		{math.MaxFloat64, 0},
		{math.SmallestNonzeroFloat64, -1},
		{0.5, -1},
	} {
		t.Run(fmt.Sprint(pair.value), func(t *testing.T) {
			idx := m.MapToIndex(pair.value)

			require.Equal(t, pair.index, idx)
		})
	}
}

// Tests invalid scales.
func TestInvalidScale(t *testing.T) {
	m, err := NewMapping(1)
	require.Error(t, err)
	require.Nil(t, m)

	m, err = NewMapping(MinScale - 1)
	require.Error(t, err)
	require.Nil(t, m)
}

// Tests a few cases with scale=-1.
func TestExponentMappingNegOne(t *testing.T) {
	m, _ := NewMapping(-1)

	for _, pair := range []expectMapping{
		{16, 2},
		{15, 1},
		{9, 1},
		{8, 1},
		{5, 1},
		{4, 1},
		{3, 0},
		{2, 0},
		{1.5, 0},
		{1, 0},
		{0.75, -1},
		{0.5, -1},
		{0.25, -1},
		{0.20, -2},
		{0.13, -2},
		{0.125, -2},
		{0.10, -2},
		{0.0625, -2},
		{0.06, -3},
	} {
		idx := m.MapToIndex(pair.value)
		require.Equal(t, pair.index, idx, "value: %v", pair.value)
	}
}

// Tests a few cases with scale=-4.
func TestExponentMappingNegFour(t *testing.T) {
	m, err := NewMapping(-4)
	require.NoError(t, err)
	require.Equal(t, int32(-4), m.Scale())

	for _, pair := range []expectMapping{
		{float64(0x1), 0},
		{float64(0x10), 0},
		{float64(0x100), 0},
		{float64(0x1000), 0},
		{float64(0x10000), 1}, // Base == 2**16
		{float64(0x100000), 1},
		{float64(0x1000000), 1},
		{float64(0x10000000), 1},
		{float64(0x100000000), 2}, // == 2**32
		{float64(0x1000000000), 2},
		{float64(0x10000000000), 2},
		{float64(0x100000000000), 2},
		{float64(0x1000000000000), 3}, // 2**48
		{float64(0x10000000000000), 3},
		{float64(0x100000000000000), 3},
		{float64(0x1000000000000000), 3},
		{float64(0x10000000000000000), 4}, // 2**64
		{float64(0x100000000000000000), 4},
		{float64(0x1000000000000000000), 4},
		{float64(0x10000000000000000000), 4},
		{float64(0x100000000000000000000), 5},

		{1 / float64(0x1), 0},
		{1 / float64(0x10), -1},
		{1 / float64(0x100), -1},
		{1 / float64(0x1000), -1},
		{1 / float64(0x10000), -1}, // 2**-16
		{1 / float64(0x100000), -2},
		{1 / float64(0x1000000), -2},
		{1 / float64(0x10000000), -2},
		{1 / float64(0x100000000), -2}, // 2**-32
		{1 / float64(0x1000000000), -3},
		{1 / float64(0x10000000000), -3},
		{1 / float64(0x100000000000), -3},
		{1 / float64(0x1000000000000), -3}, // 2**-48
		{1 / float64(0x10000000000000), -4},
		{1 / float64(0x100000000000000), -4},
		{1 / float64(0x1000000000000000), -4},
		{1 / float64(0x10000000000000000), -4}, // 2**-64
		{1 / float64(0x100000000000000000), -5},

		// Max values
		{0x1.FFFFFFFFFFFFFp1023, 63},
		{0x1p1023, 63},
		{0x1p1019, 63},
		{0x1p1008, 63},
		{0x1p1007, 62},
		{0x1p1000, 62},
		{0x1p0992, 62},
		{0x1p0991, 61},

		// Min and subnormal values
		{0x1p-1074, -64},
		{0x1p-1073, -64},
		{0x1p-1072, -64},
		{0x1p-1057, -64},
		{0x1p-1056, -64},
		{0x1p-1041, -64},
		{0x1p-1040, -64},
		{0x1p-1025, -64},
		{0x1p-1024, -64},
		{0x1p-1023, -64},
		{0x1p-1022, -64},
		{0x1p-1009, -64},
		{0x1p-1008, -63},
		{0x1p-0993, -63},
		{0x1p-0992, -62},
		{0x1p-0977, -62},
		{0x1p-0976, -61},
	} {
		t.Run(fmt.Sprintf("%x", pair.value), func(t *testing.T) {
			index := m.MapToIndex(pair.value)

			require.Equal(t, pair.index, index, "value: %#x", pair.value)
		})
	}
}

// roundedBoundary computes the correct boundary rounded to a float64
// using math/big.  Note that this function uses a Square() where the
// one in ../logarithm uses a SquareRoot().
func roundedBoundary(scale, index int32) float64 {
	one := big.NewFloat(1)
	f := (&big.Float{}).SetMantExp(one, int(index))
	for i := scale; i < 0; i++ {
		f = (&big.Float{}).Mul(f, f)
	}

	result, _ := f.Float64()
	return result
}

// TestExponentIndexMax ensures that for every valid scale, MaxFloat
// maps into the correct maximum index.  Also tests that the reverse
// lookup does not produce infinity and the following index produces
// an overflow error.
func TestExponentIndexMax(t *testing.T) {
	for scale := MinScale; scale <= MaxScale; scale++ {
		m, err := NewMapping(scale)
		require.NoError(t, err)

		index := m.MapToIndex(MaxValue)

		// Correct max index is one less than the first index
		// that overflows math.MaxFloat64, i.e., one less than
		// the index of +Inf.
		maxIndex := (int32(MaxNormalExponent+1) >> -scale) - 1
		require.Equal(t, index, int32(maxIndex))

		// The index maps to a finite boundary.
		bound, err := m.LowerBoundary(index)
		require.NoError(t, err)

		require.Equal(t, bound, roundedBoundary(scale, maxIndex))

		// One larger index will overflow.
		_, err = m.LowerBoundary(index + 1)
		require.Equal(t, err, mapping.ErrOverflow)
	}
}

// TestExponentIndexMin ensures that for every valid scale, the
// smallest normal number and all smaller numbers map to the correct
// index, which is that of the smallest normal number.
//
// Tests that the lower boundary of the smallest bucket is correct,
// even when that number is subnormal.
func TestExponentIndexMin(t *testing.T) {
	for scale := MinScale; scale <= MaxScale; scale++ {
		m, err := NewMapping(scale)
		require.NoError(t, err)

		minIndex := m.MapToIndex(MinValue)

		boundary, err := m.LowerBoundary(minIndex)
		require.NoError(t, err)

		correctMinIndex := int64(MinNormalExponent) >> -scale
		require.Greater(t, correctMinIndex, int64(math.MinInt32))
		require.Equal(t, int32(correctMinIndex), minIndex)

		correctBoundary := roundedBoundary(scale, int32(correctMinIndex))

		require.Equal(t, correctBoundary, boundary)
		require.Greater(t, roundedBoundary(scale, int32(correctMinIndex+1)), boundary)

		// Subnormal values map to the min index:
		require.Equal(t, m.MapToIndex(MinValue/2), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(MinValue/3), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(MinValue/100), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(0x1p-1050), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(0x1p-1073), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(0x1.1p-1073), int32(correctMinIndex))
		require.Equal(t, m.MapToIndex(0x1p-1074), int32(correctMinIndex))

		// One smaller index will underflow.
		_, err = m.LowerBoundary(minIndex - 1)
		require.Equal(t, err, mapping.ErrUnderflow)
	}
}
