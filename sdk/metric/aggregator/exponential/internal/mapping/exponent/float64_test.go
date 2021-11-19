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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ldexp32(x float64, i int32) float64 {
	return math.Ldexp(x, int(i))
}

func TestFloat64Bits(t *testing.T) {
	assert.Equal(t, int32(-1022), MinNormalExponent)
	assert.Equal(t, int32(+1023), MaxNormalExponent)

	assert.Equal(t, int32(1022), GetBase2(0x1p+1022))

	assert.Equal(t, int32(-1022), GetBase2(0x1p-1022))

	// Subnormals below this point
	assert.Equal(t, int32(-1023), GetBase2(0x1p-1023))
	assert.Equal(t, int32(-1024), GetBase2(0x1p-1024))
	assert.Equal(t, int32(-1025), GetBase2(0x1p-1025))

	for i := 0; i <= SignificandWidth; i++ {
		assert.Equal(t, int32(-1022-i), GetBase2(0x1p-1022/float64(uint64(1)<<i)))
	}

	// This works b/c the raw significand is zero, so 64 leading zeros - 12 = 52
	zero := 0x1p-1022 / float64(uint64(1)<<53)
	assert.Equal(t, int32(-1022-53), GetBase2(zero))
	assert.NotEqual(t, int32(-1022-54), GetBase2(0x1p-1022/float64(uint64(1)<<54)))
}

func TestGetBase2(t *testing.T) {
	for x := int32(MinNormalExponent); x <= MaxNormalExponent; x++ {
		assert.Equal(t, x, GetBase2(ldexp32(1, x)))
		assert.Equal(t, x, GetBase2(ldexp32(-1, x)))
	}

	// Smallest exponent
	assert.Equal(t, MinSubnormalExponent, GetBase2(math.Float64frombits(1)))

	assert.Equal(t, MinNormalExponent, GetBase2(0x1p-1022))
	assert.Equal(t, SignedZeroSubnormalExponent, GetBase2(0x1p-1023))
	assert.Equal(t, SignedZeroSubnormalExponent-1, GetBase2(0x1p-1024))
	assert.Equal(t, SignedZeroSubnormalExponent-51, GetBase2(0x1p-1074))
}
