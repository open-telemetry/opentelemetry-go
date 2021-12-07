package main

import (
	"fmt"
	"math"
	"math/big"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/internal/mapping/exponent"
)

func newBig() *big.Float {
	return (&big.Float{}).SetMode(big.ToNegativeInf).SetPrec(100)
}

func roundedBoundary(scale int32, index int64) float64 {
	one := big.NewFloat(1)
	f := newBig().SetMantExp(one, int(index))
	for i := scale; i > 0; i-- {
		f = newBig().Sqrt(f)
	}
	for i := scale; i < 0; i++ {
		f = newBig().Mul(f, f)
	}

	result, _ := f.Float64()
	return result
}

func calc() {
	for scale := int32(1); scale <= 20; scale++ {

		maxIndex := (int64(exponent.MaxNormalExponent+1) << scale) - 1

		if maxIndex > math.MaxInt32 {
			fmt.Println("index values are required to fit a signed 32 bit integer: ", maxIndex, "at scale", scale)
			return
		}

		maxBoundary := roundedBoundary(scale, maxIndex)
		if maxBoundary == math.Inf(+1) {
			panic("infinity detected")
		}

		ratio := math.MaxFloat64 / maxBoundary
		base := roundedBoundary(scale, 1)

		diff := math.Abs(ratio-base) / math.Abs(base)

		if diff > 1e-10 {
			panic(fmt.Sprint("expected MaxVal ratio out of range", diff))
		}

		scaleFactor := math.Ldexp(math.Log2E, int(scale))
		inverseFactor := math.Ldexp(math.Ln2, int(-scale))

		minIndex := int64(exponent.MinNormalExponent) << scale

		minBoundary := roundedBoundary(scale, minIndex)

		if minBoundary != 0x1p-1022 {
			panic("(unexpected logic error")
		}

		fmt.Printf(`	{
		scale:	       %d,
		minIndex:      %#x,
		maxIndex:      %#x,
		maxBoundary:   %x,
		scaleFactor:   %x,
		inverseFactor: %x,
	},
`, scale, minIndex, maxIndex, maxBoundary, scaleFactor, inverseFactor)
	}
}

func main() {
	fmt.Println(`// Copyright The OpenTelemetry Authors
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

// maxBoundaries is the exact lower boundary value corresponding with
// the largest bucket mapped by finite float64 values up to
// math.MaxFloat64.  This is offset by one, since the logarithm
// implementation is used with scale >= 1.
var prebuiltMappings = [MaxScale]logarithmMapping{`)
	calc()
	fmt.Println("}")
}
