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
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

type expectMapping struct {
	value float64
	index int64
}

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
		require.Equal(t, pair.index, m.MapToIndex(pair.value), "value: %v", pair.value)
	}
}

func TestLogarithmError(t *testing.T) {
	// Note: avoid large-magnitude negative scales so the indices
	// used in the loop below are viable.
	for _, scale := range []int32{-1, 0, 1, 2, 3, 4, 10, 15} {
		t.Run(fmt.Sprint(scale), func(t *testing.T) {
			m, _ := NewMapping(scale)
			for _, index := range []int64{-100, -10, -1, 0, 1, 10, 100} {
				t.Run(fmt.Sprint(index), func(t *testing.T) {
					lowBoundary, err := m.LowerBoundary(index)
					require.NoError(t, err)
					mapped := m.MapToIndex(lowBoundary)

					// At or near the boundary should be off-by-one.
					require.LessOrEqual(t, index-1, mapped)
					require.GreaterOrEqual(t, index, mapped)

					// The values should be very close.
					require.InEpsilon(t, lowBoundary, roundedBoundary(scale, index), 1e-9)
				})
			}
		})
	}
}

func roundedBoundary(scale int32, index int64) float64 {
	one := big.NewFloat(1)
	f := (&big.Float{}).SetMantExp(one, int(index))
	for i := scale; i > 0; i-- {
		f = (&big.Float{}).Sqrt(f)
	}
	for i := scale; i < 0; i++ {
		f = (&big.Float{}).Mul(f, f)
	}

	result, _ := f.Float64()
	return result
}
