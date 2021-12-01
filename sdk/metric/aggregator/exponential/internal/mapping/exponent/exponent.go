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

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
)

// MinScale defines the point at which the exponential mapping
// function becomes useless for float64.  With scale -10, ignoring
// subnormal values, bucket indices range from -1 to 1.
const MinScale int32 = -10

// exponentMapping is used for negative scales, effectively a
// mapping of the base-2 logarithm of the exponent.
type exponentMapping struct {
	scale          int32
	underflowIndex int32
	overflowIndex  int32
}

func NewMapping(scale int32) (mapping.Mapping, error) {
	if scale > 0 {
		return nil, fmt.Errorf("exponent mapping requires scale <= 0")
	}
	if scale < MinScale {
		return nil, fmt.Errorf("scale too low")
	}
	return exponentMapping{
		scale:          scale,
		underflowIndex: (MinSubnormalExponent - 1) >> -scale,
		overflowIndex:  (MaxNormalExponent + 1) >> -scale,
	}, nil
}

func (e exponentMapping) MapToIndex(value float64) int32 {
	// Note: we can assume not a 0, Inf, or NaN, ignore sign bit.
	// Errors are impossible in this code path because negative
	// scale implies indexes have smaller magnitide than the
	// corresponding values, thus all values are mapped.

	// GetBase2 compensates for subnormal values.
	exp := GetBase2(value)

	// Note: bit-shifting does the right thing for negative
	// exponents, e.g., -1 >> 1 == -1.
	return exp >> -e.scale
}

func (e exponentMapping) LowerBoundary(index int32) (float64, error) {
	if index <= e.underflowIndex {
		return 0, mapping.ErrUnderflow
	}
	if index >= e.overflowIndex {
		return 0, mapping.ErrOverflow
	}

	unbiased := int64(index << -e.scale)

	var bits uint64

	if unbiased < int64(MinNormalExponent) {
		diff := MinNormalExponent - int32(unbiased)
		shift := SignificandWidth - diff
		bits = uint64(1) << shift
	} else {
		exponent := unbiased + ExponentBias
		bits = uint64(exponent << SignificandWidth)
	}

	return math.Float64frombits(bits), nil
}

func (e exponentMapping) Scale() int32 {
	return e.scale
}
