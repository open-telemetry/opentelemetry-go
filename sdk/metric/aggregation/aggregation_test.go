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

//go:build go1.17
// +build go1.17

package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type invalidOperation struct {
	operation
}

func TestAggregationErr(t *testing.T) {
	t.Run("DropOperation", func(t *testing.T) {
		agg := Aggregation{Operation: Drop{}}
		assert.NoError(t, agg.Err())
	})

	t.Run("SumOperation", func(t *testing.T) {
		agg := Aggregation{Operation: Sum{}}
		assert.NoError(t, agg.Err())
	})

	t.Run("LastValueOperation", func(t *testing.T) {
		agg := Aggregation{Operation: LastValue{}}
		assert.NoError(t, agg.Err())
	})

	t.Run("ExplicitBucketHistogramOperation", func(t *testing.T) {
		agg := Aggregation{Operation: ExplicitBucketHistogram{}}
		assert.NoError(t, agg.Err())

		agg = Aggregation{Operation: ExplicitBucketHistogram{
			Boundaries:   []float64{0},
			RecordMinMax: true,
		}}
		assert.NoError(t, agg.Err())

		agg = Aggregation{Operation: ExplicitBucketHistogram{
			Boundaries:   []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
			RecordMinMax: true,
		}}
		assert.NoError(t, agg.Err())
	})

	t.Run("UnsetOperation", func(t *testing.T) {
		agg := Aggregation{}
		assert.ErrorIs(t, agg.Err(), errAgg)
	})

	t.Run("UnknownOperation", func(t *testing.T) {
		agg := Aggregation{Operation: invalidOperation{}}
		assert.ErrorIs(t, agg.Err(), errAgg)
	})

	t.Run("NonmonotonicHistogramBoundaries", func(t *testing.T) {
		agg := Aggregation{Operation: ExplicitBucketHistogram{
			Boundaries: []float64{2, 1},
		}}
		assert.ErrorIs(t, agg.Err(), errAgg)
	})
}
