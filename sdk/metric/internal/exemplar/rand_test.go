// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixedSize(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(n int) (Reservoir[int64], int) {
		return FixedSize[int64](n), n
	}))

	t.Run("Float64", ReservoirTest[float64](func(n int) (Reservoir[float64], int) {
		return FixedSize[float64](n), n
	}))
}

func TestFixedSizeSamplingCorrectness(t *testing.T) {
	intensity := 0.1
	sampleSize := 1000

	data := make([]float64, sampleSize*1000)
	for i := range data {
		// Generate exponentially distributed data.
		data[i] = (-1.0 / intensity) * math.Log(random())
	}
	// Sort to test position bias.
	slices.Sort(data)

	r := FixedSize[float64](sampleSize)
	for _, value := range data {
		r.Offer(context.Background(), staticTime, value, nil)
	}

	var sum float64
	for _, m := range r.(*randRes[float64]).store {
		sum += m.Value
	}
	mean := sum / float64(sampleSize)

	// Check the intensity/rate of the sampled distribution is preserved
	// ensuring no bias in our random sampling algorithm.
	assert.InDelta(t, 1/mean, intensity, 0.02) // Within 5σ.
}
