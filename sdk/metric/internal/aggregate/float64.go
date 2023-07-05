// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import "math"

const (
	// significandWidth is the size of an IEEE 754 double-precision
	// floating-point significand.
	significandWidth = 52
	// SignificandMask is the mask for the significand of an IEEE 754
	// double-precision floating-point value: 0xFFFFFFFFFFFFF.
	significandMask = 1<<significandWidth - 1
	// exponentWidth is the size of an IEEE 754 double-precision
	// floating-point exponent.
	exponentWidth = 11
	// exponentBias is the exponent bias specified for encoding
	// the IEEE 754 double-precision floating point exponent: 1023.
	exponentBias = 1<<(exponentWidth-1) - 1
	// exponentMask are set to 1 for the bits of an IEEE 754
	// floating point exponent: 0x7FF0000000000000.
	exponentMask = ((1 << exponentWidth) - 1) << significandWidth
)

// getNormalBase2 extracts the normalized base-2 fractional exponent.
// Unlike Frexp(), this returns k for the equation f x 2**k where f is
// in the range [1, 2).  Note that this function is not called for
// subnormal numbers.
func getNormalBase2(value float64) int {
	rawBits := math.Float64bits(value)
	rawExponent := (int64(rawBits) & exponentMask) >> significandWidth
	return int(rawExponent - exponentBias)
}

// getSignificand returns the 52 bit (unsigned) significand as a
// signed value.
func getSignificand(value float64) int {
	return int(int64(math.Float64bits(value)) & significandMask)
}
