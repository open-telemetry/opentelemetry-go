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

package logarithm

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"
)

type expectMapping struct {
	value float64
	index int32
}

// Tests an invalid scale.
func TestInvalidScale(t *testing.T) {
	_, err := NewMapping(-1)
	require.Error(t, err)
}

// Tests a few values are mapped correctly at scale 1, where the
// exponentiation factor is SquareRoot(2).
func TestLogarithmMapping(t *testing.T) {
	// Scale 1 means 1 division between every power of two, having
	// a factor sqrt(2) times the lower boundary.
	m, err := NewMapping(+1)
	require.NoError(t, err)
	require.Equal(t, int32(+1), m.Scale())

	// Note: Do not test exact boundaries, with the exception of
	// 1, because we expect errors in that case (e.g.,
	// MapToIndex(8) returns 5, an off-by-one.  See the following
	// test.
	for _, pair := range []expectMapping{
		{15, 7},
		{9, 6},
		{7, 5},
		{5, 4},
		{3, 3},
		{2.5, 2},
		{1.5, 1},
		{1.2, 0},
		{1, 0},
		{0.75, -1},
		{0.55, -2},
		{0.45, -3},
	} {
		idx := m.MapToIndex(pair.value)
		require.Equal(t, pair.index, idx, "value: %v", pair.value)
	}
}

// Tests the mapping function for correctness-within-epsilon for a few
// scales and index values.
func TestLogarithmBoundary(t *testing.T) {
	for _, scale := range []int32{1, 2, 3, 4, 10, 15} {
		t.Run(fmt.Sprint(scale), func(t *testing.T) {
			m, _ := NewMapping(scale)
			for _, index := range []int32{-100, -10, -1, 0, 1, 10, 100} {
				t.Run(fmt.Sprint(index), func(t *testing.T) {
					lowBoundary, err := m.LowerBoundary(index)
					require.NoError(t, err)
					mapped := m.MapToIndex(lowBoundary)

					// At or near the boundary expected to be off-by-one sometimes.
					require.LessOrEqual(t, index-1, mapped)
					require.GreaterOrEqual(t, index, mapped)

					// The values should be very close.
					require.InEpsilon(t, lowBoundary, roundedBoundary(scale, index), 1e-9)
				})
			}
		})
	}
}

// roundedBoundary computes the correct boundary rounded to a float64
// using math/big.  Note that this function uses a SquareRoot() where the
// one in ../exponent uses a Square().
func roundedBoundary(scale, index int32) float64 {
	one := big.NewFloat(1)
	f := (&big.Float{}).SetMantExp(one, int(index))
	for i := scale; i > 0; i-- {
		f = (&big.Float{}).Sqrt(f)
	}

	result, _ := f.Float64()
	return result
}

// TestLogarithmIndexMax ensures that for every valid scale, MaxFloat
// maps into the correct maximum index.  Also tests that the reverse
// lookup does not produce infinity and the following index produces
// an overflow error.
func TestLogarithmIndexMax(t *testing.T) {
	for scale := MinScale; scale <= MaxScale; scale++ {
		m, err := NewMapping(scale)
		require.NoError(t, err)

		index := m.MapToIndex(MaxValue)

		// Correct max index is one less than the first index
		// that overflows math.MaxFloat64, i.e., one less than
		// the index of +Inf.
		maxIndex64 := (int64(exponent.MaxNormalExponent+1) << scale) - 1
		require.Less(t, maxIndex64, int64(math.MaxInt32))
		require.Equal(t, index, int32(maxIndex64))

		// The index maps to a finite boundary near MaxFloat.
		bound, err := m.LowerBoundary(index)
		require.NoError(t, err)

		base, _ := m.LowerBoundary(1)

		require.Less(t, bound, MaxValue)

		// The expected ratio equals the base factor.
		require.InEpsilon(t, (MaxValue-bound)/bound, base-1, 1e-6)

		// One larger index will overflow.
		_, err = m.LowerBoundary(index + 1)
		require.Equal(t, err, mapping.ErrOverflow)
	}
}

// TestLogarithmIndexMin ensures that for every valid scale, Non-zero numbers.
func TestLogarithmIndexMin(t *testing.T) {
	for scale := MinScale; scale <= MaxScale; scale++ {
		m, err := NewMapping(scale)
		require.NoError(t, err)

		minIndex := m.MapToIndex(MinValue)

		mapped, err := m.LowerBoundary(minIndex)
		require.NoError(t, err)

		correctMinIndex := int64(exponent.MinNormalExponent) << scale
		require.Greater(t, correctMinIndex, int64(math.MinInt32))

		correctMapped := roundedBoundary(scale, int32(correctMinIndex))
		require.Equal(t, correctMapped, MinValue)
		require.InEpsilon(t, mapped, MinValue, 1e-6)

		require.Equal(t, minIndex, int32(correctMinIndex))

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
		maxIndex64 := (int64(exponent.MaxNormalExponent+1) << scale) - 1
		require.Less(t, maxIndex64, int64(math.MaxInt32))
		require.Equal(t, index, int32(maxIndex64))

		// The index maps to a finite boundary near MaxFloat.
		bound, err := m.LowerBoundary(index)
		require.NoError(t, err)

		base, _ := m.LowerBoundary(1)

		require.Less(t, bound, MaxValue)

		// The expected ratio equals the base factor.
		require.InEpsilon(t, (MaxValue-bound)/bound, base-1, 1e-6)

		// One larger index will overflow.
		_, err = m.LowerBoundary(index + 1)
		require.Equal(t, err, mapping.ErrOverflow)
	}
}

// TestExponentIndexMin ensures that for every valid scale, the
// smallest normal number and all smaller numbers map to the correct
// index, which is that of the smallest normal number.
func TestExponentIndexMin(t *testing.T) {
	for scale := MinScale; scale <= MaxScale; scale++ {
		m, err := NewMapping(scale)
		require.NoError(t, err)

		minIndex := m.MapToIndex(MinValue)

		mapped, err := m.LowerBoundary(minIndex)
		require.NoError(t, err)

		correctMinIndex := int64(exponent.MinNormalExponent) << scale
		require.Greater(t, correctMinIndex, int64(math.MinInt32))

		correctMapped := roundedBoundary(scale, int32(correctMinIndex))
		require.Equal(t, correctMapped, MinValue)
		require.InEpsilon(t, mapped, MinValue, 1e-6)

		require.Equal(t, minIndex, int32(correctMinIndex))

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
