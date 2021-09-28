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

package exponential

import (
	"fmt"
	"math"
)

const (
	// MantissaWidth is the size of an IEEE 754 double-precision
	// floating-point mantissa.
	MantissaWidth = 52
	// ExponentWidth is the size of an IEEE 754 double-precision
	// floating-point exponent.
	ExponentWidth = 11

	// MantissaMask is the mask for the mantissa of an IEEE 754
	// double-precision floating-point value.
	MantissaMask = 1<<MantissaWidth - 1

	// ExponentBias is the exponent bias specified for encoding
	// the IEEE 754 double-precision floating point exponent.
	ExponentBias = 1<<(ExponentWidth-1) - 1

	// ExponentMask are set to 1 for the bits of an IEEE 754
	// floating point exponent (as distinct from the mantissa and
	// sign).
	ExponentMask = ((1 << ExponentWidth) - 1) << MantissaWidth

	// SignMask selects the sign bit of an IEEE 754 floating point
	// number.
	SignMask = (1 << 63)
)

type (
	// Mapping is the interface of a mapper.
	Mapping interface {
		MapToIndex(value float64) int64
		LowerBoundary(index int64) float64
		Scale() int32
	}

	// exponentMapping is used for negative scales, effectively a
	// mapping of the base-2 logarithm of the exponent.
	exponentMapping struct {
		scale int32
	}
)

func newMapping(scale int32) Mapping {
	fmt.Println("new scale", scale)
	if scale < 0 {
		return exponentMapping{scale}
	}
	return newLogarithmMapping(scale)
}

func (e exponentMapping) MapToIndex(value float64) int64 {
	return int64(getExponent(value) >> -e.scale)
}

func (e exponentMapping) LowerBoundary(index int64) float64 {
	exponent := int64(index<<-e.scale) + ExponentBias

	return math.Float64frombits(uint64(exponent << MantissaWidth))
}

func (e exponentMapping) Scale() int32 {
	return e.scale
}

// java.lang.Math.scalb(float f, int scaleFactor) returns f x
// 2**scaleFactor, rounded as if performed by a single correctly
// rounded floating-point multiply to a member of the double value set.
func scalb(f float64, sf int32) float64 {
	if f == 0 {
		return 0
	}
	valueBits := math.Float64bits(f)

	signBit := valueBits & SignMask
	mantissa := MantissaMask & valueBits

	exponent := (int64(valueBits) & ExponentMask) >> MantissaWidth
	exponent += int64(sf)

	return math.Float64frombits(signBit | uint64(exponent<<MantissaWidth) | mantissa)
}

func getExponent(value float64) int32 {
	valueBits := math.Float64bits(value)
	exponent := ((int64(valueBits) & ExponentMask) >> MantissaWidth) - ExponentBias
	return int32(exponent)
}
