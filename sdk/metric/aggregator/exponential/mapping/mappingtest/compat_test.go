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

package histogram_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/logarithm"
)

func testCompatibility(t *testing.T, histoScale int32) {
	src := rand.New(rand.NewSource(54979))
	t.Run(fmt.Sprintf("compat_%d", histoScale), func(t *testing.T) {
		const trials = 1e5

		ltm := exponent.NewMapping(histoScale)
		lgm := logarithm.NewMapping(histoScale)

		for i := 0; i < trials; i++ {
			// Generate a random normalized number.
			v := mapping.Scalb(
				1+src.Float64(),
				mapping.MinNormalExponent+int32(src.Intn(int(1+mapping.MaxNormalExponent-mapping.MinNormalExponent))))

			lti := ltm.MapToIndex(v)
			lgi := lgm.MapToIndex(v)

			assert.Equal(t, lti, lgi)

			ltb, err1 := ltm.LowerBoundary(lti)
			lgb, err2 := lgm.LowerBoundary(lti)

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			assert.InEpsilon(t, ltb, lgb, 0.000001)
		}
	})
}

func TestCompatExponentMapping(t *testing.T) {
	for scale := int32(0); scale >= -4; scale-- {
		testCompatibility(t, scale)
	}
}
