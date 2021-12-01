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
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
)

type expectMapping struct {
	value float64
	index int32
}

type invalidMapping struct {
	scale int32
	index int32
}

const (
	testUnderflowIndex = math.MinInt32
	testOverflowIndex  = math.MaxInt32
)

func TestExponentMappingZero(t *testing.T) {
	m, err := NewMapping(0)
	require.NoError(t, err)

	require.Equal(t, int32(0), m.Scale())

	for _, pair := range []expectMapping{
		{math.MaxFloat64, 1023},
		{math.SmallestNonzeroFloat64, -1074},
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

func TestInvalidScale(t *testing.T) {
	m, err := NewMapping(1)
	require.Error(t, err)
	require.Nil(t, m)

	m, err = NewMapping(MinScale - 1)
	require.Error(t, err)
	require.Nil(t, m)
}

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
		{0x1.FFFFFFFFFFFFFp1023, testOverflowIndex},
		{0x1p1023, testOverflowIndex},
		{0x1p1019, testOverflowIndex},
		{0x1p1008, testOverflowIndex},
		{0x1p1007, 62},
		{0x1p1000, 62},
		{0x1p0992, 62},
		{0x1p0991, 61},

		// Min and subnormal values
		{0x1p-1074, testUnderflowIndex},
		{0x1p-1073, testUnderflowIndex},

		{0x1p-1072, -67}, // n.b. 67 * 2**4 == 1072
		{0x1p-1057, -67},
		{0x1p-1056, -66},
		{0x1p-1041, -66},
		{0x1p-1040, -65},
		{0x1p-1025, -65},
		{0x1p-1024, -64},
		{0x1p-1009, -64},
		{0x1p-1008, -63},
		{0x1p-0993, -63},
		{0x1p-0992, -62},
		{0x1p-0977, -62},
		{0x1p-0976, -61},
	} {
		t.Run(fmt.Sprintf("%x", pair.value), func(t *testing.T) {
			index := m.MapToIndex(pair.value)

			if pair.index != testUnderflowIndex && pair.index != testOverflowIndex {
				require.Equal(t, pair.index, index, "value: %#x", pair.value)
			}

			lb, err1 := m.LowerBoundary(index)
			ub, err2 := m.LowerBoundary(index + 1)

			if pair.index != testUnderflowIndex && pair.index != testOverflowIndex {
				require.NoError(t, err1)
				require.NoError(t, err2)

				require.NotEqual(t, 0., lb)
				require.NotEqual(t, 0., ub)
				require.LessOrEqual(t, lb, pair.value, fmt.Sprintf("value: %x index %v", pair.value, index))
				require.Greater(t, ub, pair.value, fmt.Sprintf("value: %x index %v", pair.value, index))
			} else {
				// The upper or lower boundary must see an error.
				require.True(t, err1 != nil || err2 != nil)
				var expectErr error
				if pair.index == testUnderflowIndex {
					expectErr = mapping.ErrUnderflow
				} else {
					expectErr = mapping.ErrOverflow
				}

				if err1 != nil {
					require.Equal(t, err1, expectErr)
				}
				if err2 != nil {
					require.Equal(t, err2, expectErr)
				}
			}
		})
	}
}

func TestExponentMappingInvalid(t *testing.T) {
	for _, pair := range []invalidMapping{
		{-4, 64},
		{-4, 65},
		{-4, 66},
		{-4, 1e6},

		{-4, -68},
		{-4, -69},
		{-4, -70},
		{-4, -1e6},

		{0, 1024},
		{0, 1025},
		{0, 1026},
		{0, 1e6},

		{0, -1075},
		{0, -1076},
		{0, -1077},
		{0, -1e6},

		{-1, 513},
		{-1, 514},
		{-1, 515},
		{-1, 1e6},

		{-1, -538},
		{-1, -539},
		{-1, -540},
		{-1, -1e6},
	} {
		t.Run(fmt.Sprintf("%v_%d", pair.scale, pair.index), func(t *testing.T) {
			m, err := NewMapping(pair.scale)
			require.NoError(t, err)
			_, err = m.LowerBoundary(pair.index)
			require.Error(t, err)
		})
	}
}
