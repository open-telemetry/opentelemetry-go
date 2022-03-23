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

package exponent // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"

import (
	"fmt"
	"math"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
)

const (
	// MinScale defines the point at which the exponential mapping
	// function becomes useless for float64.  With scale -10, ignoring
	// subnormal values, bucket indices range from -1 to 1.
	MinScale int32 = -10

	// MaxScale is the largest scale supported in this code.  Use
	// ../logarithm for larger scales.
	MaxScale int32 = 0
)

type exponentMapping struct {
	shift uint8 // equals negative scale
}

// exponentMapping is used for negative scales, effectively a
// mapping of the base-2 logarithm of the exponent.
var prebuiltMappings = [-MinScale + 1]exponentMapping{
	{10},
	{9},
	{8},
	{7},
	{6},
	{5},
	{4},
	{3},
	{2},
	{1},
	{0},
}

// NewMapping constructs an exponential mapping function, used for scales <= 0.
func NewMapping(scale int32) (mapping.Mapping, error) {
	if scale > MaxScale {
		return nil, fmt.Errorf("exponent mapping requires scale <= 0")
	}
	if scale < MinScale {
		return nil, fmt.Errorf("scale too low")
	}
	return &prebuiltMappings[scale-MinScale], nil
}

// MapToIndex implements mapping.Mapping.
func (e *exponentMapping) MapToIndex(value float64) int32 {
	// Note: we can assume not a 0, Inf, or NaN; positive sign bit.

	// Note: bit-shifting does the right thing for negative
	// exponents, e.g., -1 >> 1 == -1.
	return getBase2(value) >> e.shift
}

func (e *exponentMapping) minIndex() int32 {
	return int32(MinNormalExponent) >> e.shift
}

func (e *exponentMapping) maxIndex() int32 {
	return int32(MaxNormalExponent) >> e.shift
}

// LowerBoundary implements mapping.Mapping.
func (e *exponentMapping) LowerBoundary(index int32) (float64, error) {
	if min := e.minIndex(); index < min {
		return 0, mapping.ErrUnderflow
	}

	if max := e.maxIndex(); index > max {
		return 0, mapping.ErrOverflow
	}

	unbiased := int64(index << e.shift)

	// Note: although the mapping function rounds subnormal values
	// up to the smallest normal value, there are still buckets
	// that may be filled that start at subnormal values.  The
	// following code handles this correctly.  It's equivalent to and
	// faster than math.Ldexp(1, int(unbiased)).
	if unbiased < int64(MinNormalExponent) {
		subnormal := uint64(1 << SignificandWidth)
		for unbiased < int64(MinNormalExponent) {
			unbiased++
			subnormal >>= 1
		}
		return math.Float64frombits(subnormal), nil
	}
	exponent := unbiased + ExponentBias

	bits := uint64(exponent << SignificandWidth)
	return math.Float64frombits(bits), nil
}

// Scale implements mapping.Mapping.
func (e *exponentMapping) Scale() int32 {
	return -int32(e.shift)
}
