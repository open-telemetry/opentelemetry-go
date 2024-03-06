// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregationErr(t *testing.T) {
	t.Run("DropOperation", func(t *testing.T) {
		assert.NoError(t, AggregationDrop{}.err())
	})

	t.Run("SumOperation", func(t *testing.T) {
		assert.NoError(t, AggregationSum{}.err())
	})

	t.Run("LastValueOperation", func(t *testing.T) {
		assert.NoError(t, AggregationLastValue{}.err())
	})

	t.Run("ExplicitBucketHistogramOperation", func(t *testing.T) {
		assert.NoError(t, AggregationExplicitBucketHistogram{}.err())

		assert.NoError(t, AggregationExplicitBucketHistogram{
			Boundaries: []float64{0},
			NoMinMax:   true,
		}.err())

		assert.NoError(t, AggregationExplicitBucketHistogram{
			Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
		}.err())
	})

	t.Run("NonmonotonicHistogramBoundaries", func(t *testing.T) {
		assert.ErrorIs(t, AggregationExplicitBucketHistogram{
			Boundaries: []float64{2, 1},
		}.err(), errAgg)

		assert.ErrorIs(t, AggregationExplicitBucketHistogram{
			Boundaries: []float64{0, 1, 2, 1, 3, 4},
		}.err(), errAgg)
	})

	t.Run("ExponentialHistogramOperation", func(t *testing.T) {
		assert.NoError(t, AggregationBase2ExponentialHistogram{
			MaxSize:  160,
			MaxScale: 20,
		}.err())

		assert.NoError(t, AggregationBase2ExponentialHistogram{
			MaxSize:  1,
			NoMinMax: true,
		}.err())

		assert.NoError(t, AggregationBase2ExponentialHistogram{
			MaxSize:  1024,
			MaxScale: -3,
		}.err())
	})

	t.Run("InvalidExponentialHistogramOperation", func(t *testing.T) {
		// MazSize must be greater than 0
		assert.ErrorIs(t, AggregationBase2ExponentialHistogram{}.err(), errAgg)

		// MaxScale Must be <=20
		assert.ErrorIs(t, AggregationBase2ExponentialHistogram{
			MaxSize:  1,
			MaxScale: 30,
		}.err(), errAgg)
	})
}

func TestExplicitBucketHistogramDeepCopy(t *testing.T) {
	const orig = 0.0
	b := []float64{orig}
	h := AggregationExplicitBucketHistogram{Boundaries: b}
	cpH := h.copy().(AggregationExplicitBucketHistogram)
	b[0] = orig + 1
	assert.Equal(t, orig, cpH.Boundaries[0], "changing the underlying slice data should not affect the copy")
}
