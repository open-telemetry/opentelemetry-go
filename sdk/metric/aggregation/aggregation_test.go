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

package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregationErr(t *testing.T) {
	t.Run("DropOperation", func(t *testing.T) {
		assert.NoError(t, Drop{}.Err())
	})

	t.Run("SumOperation", func(t *testing.T) {
		assert.NoError(t, Sum{}.Err())
	})

	t.Run("LastValueOperation", func(t *testing.T) {
		assert.NoError(t, LastValue{}.Err())
	})

	t.Run("ExplicitBucketHistogramOperation", func(t *testing.T) {
		assert.NoError(t, ExplicitBucketHistogram{}.Err())

		assert.NoError(t, ExplicitBucketHistogram{
			Boundaries: []float64{0},
			NoMinMax:   true,
		}.Err())

		assert.NoError(t, ExplicitBucketHistogram{
			Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
		}.Err())
	})

	t.Run("NonmonotonicHistogramBoundaries", func(t *testing.T) {
		assert.ErrorIs(t, ExplicitBucketHistogram{
			Boundaries: []float64{2, 1},
		}.Err(), errAgg)

		assert.ErrorIs(t, ExplicitBucketHistogram{
			Boundaries: []float64{0, 1, 2, 1, 3, 4},
		}.Err(), errAgg)
	})
}

func TestExplicitBucketHistogramDeepCopy(t *testing.T) {
	const orig = 0.0
	b := []float64{orig}
	h := ExplicitBucketHistogram{Boundaries: b}
	cpH := h.Copy().(ExplicitBucketHistogram)
	b[0] = orig + 1
	assert.Equal(t, orig, cpH.Boundaries[0], "changing the underlying slice data should not affect the copy")
}
