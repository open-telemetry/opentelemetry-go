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

package mapping

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalb(t *testing.T) {
	assert.Equal(t, 2.0, Scalb(1, 1))
	assert.Equal(t, 0.5, Scalb(1, -1))
	assert.Equal(t, -2.0, Scalb(-1, 1))
	assert.Equal(t, -0.5, Scalb(-1, -1))

	assert.Equal(t, 0x1p-1000, Scalb(1, -1000))
	assert.Equal(t, 0x1p+1000, Scalb(1, +1000))
}

func TestFloat64Bits(t *testing.T) {
	assert.Equal(t, int32(-1022), MinNormalExponent)
	assert.Equal(t, int32(+1023), MaxNormalExponent)

	assert.Equal(t, int32(1022), GetExponent(0x1p+1022))

	assert.Equal(t, int32(-1022), GetExponent(0x1p-1022))

	// Subnormals below this point
	assert.Equal(t, int32(-1023), GetExponent(0x1p-1023))
	assert.Equal(t, int32(-1024), GetExponent(0x1p-1024))
	assert.Equal(t, int32(-1025), GetExponent(0x1p-1025))

	for i := 0; i <= SignificandWidth; i++ {
		assert.Equal(t, int32(-1022-i), GetExponent(0x1p-1022/float64(uint64(1)<<i)))
	}

	// This works b/c the raw significand is zero, so 64 leading zeros - 12 = 52
	zero := 0x1p-1022 / float64(uint64(1)<<53)
	assert.Equal(t, int32(-1022-53), GetExponent(zero))
	assert.NotEqual(t, int32(-1022-54), GetExponent(0x1p-1022/float64(uint64(1)<<54)))
}

func TestGetExponent(t *testing.T) {
	for x := int32(MinNormalExponent); x <= MaxNormalExponent; x++ {
		assert.Equal(t, x, GetExponent(Scalb(1, x)))
		assert.Equal(t, x, GetExponent(Scalb(-1, x)))
	}

	// GetExponent and Scalb work for the special exponents (good or bad)
	assert.Equal(t, InfAndNaNExponent, GetExponent(Scalb(1, InfAndNaNExponent)))
	assert.Equal(t, InfAndNaNExponent, GetExponent(Scalb(-1, InfAndNaNExponent)))

	// Check that scalb overflows
	assert.NotEqual(t, InfAndNaNExponent+1, GetExponent(Scalb(1, InfAndNaNExponent+1)))
	assert.NotEqual(t, SignedZeroSubnormalExponent-1, GetExponent(Scalb(1, SignedZeroSubnormalExponent-1)))

	// Smallest exponent
	assert.Equal(t, MinSubnormalExponent, GetExponent(math.Float64frombits(1)))

	assert.Equal(t, MinNormalExponent, GetExponent(0x1p-1022))
	assert.Equal(t, SignedZeroSubnormalExponent, GetExponent(0x1p-1023))
	assert.Equal(t, SignedZeroSubnormalExponent-1, GetExponent(0x1p-1024))
	assert.Equal(t, SignedZeroSubnormalExponent-51, GetExponent(0x1p-1074))
}
