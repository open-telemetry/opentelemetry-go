// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
	"math"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedSize(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(n int) (Reservoir, int) {
		return FixedSize(n), n
	}))

	t.Run("Float64", ReservoirTest[float64](func(n int) (Reservoir, int) {
		return FixedSize(n), n
	}))
}

func TestFixedSizeSamplingCorrectness(t *testing.T) {
	intensity := 0.1
	sampleSize := 1000

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	data := make([]float64, sampleSize*1000)
	for i := range data {
		// Generate exponentially distributed data.
		data[i] = (-1.0 / intensity) * math.Log(rng.Float64())
	}
	// Sort to test position bias.
	slices.Sort(data)

	r := FixedSize(sampleSize)
	for _, value := range data {
		r.Offer(context.Background(), staticTime, NewValue(value), nil)
	}

	var sum float64
	for _, m := range r.(*randRes).store {
		sum += m.Value.Float64()
	}
	mean := sum / float64(sampleSize)

	// Check the intensity/rate of the sampled distribution is preserved
	// ensuring no bias in our random sampling algorithm.
	assert.InDelta(t, 1/mean, intensity, 0.02) // Within 5σ.
}
