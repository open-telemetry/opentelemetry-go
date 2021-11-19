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

package logarithm

import (
	"fmt"
	"math"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping"
)

const (
	MinScale int32 = 1
	MaxScale int32 = 20

	// MinValue is the smallest normal floating point value.  This
	// limit is necessary because these mapping functions do not
	// perform correctly for subnormal values.  Subnormal values are
	// supported by the exponent mapper.
	MinValue = 0x1p-1022
)

// This implementation was copied from a Java prototype. See:
// https://github.com/newrelic-experimental/newrelic-sketch-java/blob/1ce245713603d61ba3a4510f6df930a5479cd3f6/src/main/java/com/newrelic/nrsketch/indexer/LogIndexer.java
// for the equations used here.
type logarithmMapping struct {
	scale int32

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
	scaleFactor float64

	// log(boundary) = index * log(base)
	// log(boundary) = index * log(2^(2^-scale))
	// log(boundary) = index * 2^-scale * log(2)
	// boundary = exp(index * inverseFactor)
	// where:
	// inverseFactor = 2^-scale * log(2)
	// = math.Ldexp(math.Ln2, -scale)
	inverseFactor float64

	overflowIndex  int32
	underflowIndex int32
}

var _ mapping.Mapping = &logarithmMapping{}

func NewMapping(scale int32) (mapping.Mapping, error) {
	// An assumption used in this code is that scale is > 0.  If
	// scale is <= 0 it's better to use the exponent mapping.
	if scale <= 0 || scale > 20 {
		// scale 20 can represent the entire float64 range
		// with a 31 bit index, and we don't handle larger
		// scales to simplify range tests in this package.
		return nil, fmt.Errorf("expect 0 < scale <= 20")
	}

	l := &logarithmMapping{
		scale:         scale,
		scaleFactor:   math.Ldexp(math.Log2E, int(scale)),
		inverseFactor: math.Ldexp(math.Ln2, int(-scale)),
	}

	maxIdx, _ := l.MapToIndex(math.MaxFloat64)
	minIdx, _ := l.MapToIndex(MinValue)

	for l.lowerBoundary(maxIdx) == math.Inf(+1) {
		maxIdx--
	}

	for l.lowerBoundary(minIdx) == 0 {
		minIdx++
	}
	l.overflowIndex = maxIdx + 1
	l.underflowIndex = minIdx - 1

	return l, nil
}

func (l *logarithmMapping) MapToIndex(value float64) (int32, error) {
	// Use Floor() to round toward -Inf.
	index := math.Floor(math.Log(value) * l.scaleFactor)
	// The restriction on scale <= 20 ensures that index
	// fits an int32 for all float64 values.
	return int32(index), nil
}

func (l *logarithmMapping) LowerBoundary(index int32) (float64, error) {
	if index <= l.underflowIndex {
		return 0, mapping.ErrUnderflow
	}
	if index >= l.overflowIndex {
		return 0, mapping.ErrOverflow
	}
	return l.lowerBoundary(index), nil
}

func (l *logarithmMapping) lowerBoundary(index int32) float64 {
	return math.Exp(float64(index) * l.inverseFactor)
}

func (l *logarithmMapping) Scale() int32 {
	return l.scale
}
