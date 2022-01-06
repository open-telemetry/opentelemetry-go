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

package logarithm // import "go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping/logarithm"

import (
	"fmt"
	"math"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
)

const (
	// MinScale ensures that the ../exponent mapper is used for
	// zero and negative scale values.  Do not use the logarithm
	// mapper for scales <= 0.
	MinScale int32 = 1

	// MaxScale is selected as the largest scale that is possible
	// in current code, considering there are 10 bits of base-2
	// exponent combined with scale-bits of range.  At this scale,
	// the growth factor is 0.0000661%.
	//
	// Scales larger than 20 complicate the logic in cmd/prebuild,
	// because math/big overflows when exponent is math.MaxInt32
	// (== the index of math.MaxFloat64 at scale=21),
	//
	// At scale=20, index values are in the interval [-0x3fe00000,
	// 0x3fffffff], having 31 bits of information.  This is
	// sensible given that the OTLP exponential histogram data
	// point uses a signed 32 bit integer for indices.
	MaxScale int32 = 20

	// MaxValue is the largest normal number.
	MaxValue = math.MaxFloat64

	// MinValue is the smallest normal number.
	MinValue = 0x1p-1022
)

// logarithmMapping contains the constants used to implement the
// exponential mapping function for a particular scale > 0.  Note that
// these structs are compiled in using code generated by the
// ./cmd/prebuild package, this way no allocations are required as the
// aggregators switch between mapping functions and the two mapping
// functions are kept separate.
//
// Note that some of these fields could be calculated easily at
// runtime, but they are compiled in to avoid those operations at
// runtime (e.g., calls to math.Ldexp(math.Log2E, scale) for every
// measurement).
type logarithmMapping struct {
	// scale is between MinScale and MaxScale
	scale int32

	// minIndex is the index of MinValue
	minIndex int32
	// maxIndex is the index of MaxValue
	maxIndex int32

	// maxBoundary is the correct LowerBoundary() of maxIndex.
	// Note that this cannot be easily computed using the ordinary
	// LowerBoundary() equations because of overflows.
	maxBoundary float64

	// scaleFactor is used and computed as follows:
	// index = log(value) / log(base)
	// = log(value) / log(2^(2^-scale))
	// = log(value) / (2^-scale * log(2))
	// = log(value) * (1/log(2) * 2^scale)
	// = log(value) * scaleFactor
	// where:
	// scaleFactor = (1/log(2) * 2^scale)
	// = math.Log2E * math.Exp2(scale)
	// = math.Ldexp(math.Log2E, scale)
	// Because multiplication is faster than division, we define scaleFactor as a multiplier.
	// This implementation was copied from a Java prototype. See:
	// https://github.com/newrelic-experimental/newrelic-sketch-java/blob/1ce245713603d61ba3a4510f6df930a5479cd3f6/src/main/java/com/newrelic/nrsketch/indexer/LogIndexer.java
	// for the equations used here.
	scaleFactor float64

	// log(boundary) = index * log(base)
	// log(boundary) = index * log(2^(2^-scale))
	// log(boundary) = index * 2^-scale * log(2)
	// boundary = exp(index * inverseFactor)
	// where:
	// inverseFactor = 2^-scale * log(2)
	// = math.Ldexp(math.Ln2, -scale)
	inverseFactor float64
}

var _ mapping.Mapping = &logarithmMapping{}

// NewMapping constructs a logarithm mapping function, used for scales > 0.
func NewMapping(scale int32) (mapping.Mapping, error) {
	// An assumption used in this code is that scale is > 0.  If
	// scale is <= 0 it's better to use the exponent mapping.
	if scale < MinScale || scale > MaxScale {
		// scale 20 can represent the entire float64 range
		// with a 30 bit index, and we don't handle larger
		// scales to simplify range tests in this package.
		return nil, fmt.Errorf("scale out of bounds")
	}
	return &prebuiltMappings[scale-MinScale], nil
}

// MapToIndex implements mapping.Mapping.
func (l *logarithmMapping) MapToIndex(value float64) int32 {
	// Note: we can assume not a 0, Inf, or NaN; positive sign bit.
	if value >= l.maxBoundary {
		return l.maxIndex
	}
	if value <= MinValue {
		return l.minIndex
	}
	// Use Floor() to round toward 0.
	return int32(math.Floor(math.Log(value) * l.scaleFactor))
}

// LowerBoundary implements mapping.Mapping.
func (l *logarithmMapping) LowerBoundary(index int32) (float64, error) {
	if index >= l.maxIndex {
		if index == l.maxIndex {
			// Note: the formula below behaves poorly
			// near the boundary, will return +Inf instead
			// of the correct last bucket lower boundary.
			// This implementation hard-codes the maximum
			// boundary for this reason.
			return l.maxBoundary, nil
		}
		return 0, mapping.ErrOverflow
	}
	if index <= l.minIndex {
		if index == l.minIndex {
			return MinValue, nil
		}
		return 0, mapping.ErrUnderflow
	}
	return math.Exp(float64(index) * l.inverseFactor), nil
}

// Scale implements mapping.Mapping.
func (l *logarithmMapping) Scale() int32 {
	return l.scale
}
