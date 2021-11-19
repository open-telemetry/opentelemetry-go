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
	"math/bits"
)

const (
	// SignificandWidth is the size of an IEEE 754 double-precision
	// floating-point significand.
	SignificandWidth = 52
	// ExponentWidth is the size of an IEEE 754 double-precision
	// floating-point exponent.
	ExponentWidth = 11

	// SignificandMask is the mask for the significand of an IEEE 754
	// double-precision floating-point value: 0xFFFFFFFFFFFFF.
	SignificandMask = 1<<SignificandWidth - 1

	// ExponentBias is the exponent bias specified for encoding
	// the IEEE 754 double-precision floating point exponent: 1023.
	ExponentBias = 1<<(ExponentWidth-1) - 1

	// ExponentMask are set to 1 for the bits of an IEEE 754
	// floating point exponent: 0x7FF0000000000000.
	ExponentMask = ((1 << ExponentWidth) - 1) << SignificandWidth

	// SignMask selects the sign bit of an IEEE 754 floating point
	// number.
	SignMask = (1 << 63)

	// MinNormalExponent is the minimum exponent of a normalized
	// floating point.
	MinNormalExponent int32 = -ExponentBias + 1

	// MaxNormalExponent is the maximum exponent of a normalized
	// floating point.
	MaxNormalExponent int32 = ExponentBias

	// SignedZeroSubnormalExponent is the exponent value after
	// removing bias for signed zero and subnormal values.
	SignedZeroSubnormalExponent int32 = -ExponentBias

	// InfAndNaNExponent is the exponent value after removing bias
	// for Inf and NaN values.
	InfAndNaNExponent int32 = ExponentBias + 1

	// Smallest positive subnormal exponent:
	MinSubnormalExponent int32 = MinNormalExponent - SignificandWidth
)

// GetExponent extracts the normalized base-2 fractional exponent.
// Let the value be represented as `1.significand x 2**exponent`,
// this returns `exponent`.  Not defined for 0, Inf, or NaN values.
//
// Note! THIS RETURNS A DIFFERENT RESULT THAN math.Frexp(), which
// does not handle subnormal values.
func GetExponent(value float64) int32 {
	rawBits := math.Float64bits(value)
	rawExponent := (int64(rawBits) & ExponentMask) >> SignificandWidth
	rawSignificand := rawBits & SignificandMask
	if rawExponent == 0 {
		// Handle subnormal values: rawSignificand cannot be zero
		// unless value is zero.
		rawExponent -= int64(bits.LeadingZeros64(rawSignificand) - 12)
	}
	return int32(rawExponent - ExponentBias)
}
