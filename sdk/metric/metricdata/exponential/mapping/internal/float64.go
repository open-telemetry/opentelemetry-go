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

package internal // import "go.opentelemetry.io/otel/sdk/metricdata/exponential/mapping/internal"

import "math"

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
	SignMask = (1 << (SignificandWidth + ExponentWidth))

	// MinNormalExponent is the minimum exponent of a normalized
	// floating point: -1022.
	MinNormalExponent int32 = -ExponentBias + 1

	// MaxNormalExponent is the maximum exponent of a normalized
	// floating point: 1023.
	MaxNormalExponent int32 = ExponentBias

	// MinValue is the smallest normal number.
	MinValue = 0x1p-1022

	// MaxValue is the largest normal number.
	MaxValue = math.MaxFloat64
)

// GetNormalBase2 extracts the normalized base-2 fractional exponent.
// Unlike Frexp(), this returns k for the equation f x 2**k where f is
// in the range [1, 2).  Note that this function is not called for
// subnormal numbers.
func GetNormalBase2(value float64) int32 {
	rawBits := math.Float64bits(value)
	rawExponent := (int64(rawBits) & ExponentMask) >> SignificandWidth
	return int32(rawExponent - ExponentBias)
}

// GetSignificand returns the 52 bit (unsigned) significand as a
// signed value.
func GetSignificand(value float64) int64 {
	return int64(math.Float64bits(value)) & SignificandMask
}
