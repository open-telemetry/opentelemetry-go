// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
	"math"
	"slices"
	"sync"
	"testing"

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

	data := make([]float64, sampleSize*1000)
	for i := range data {
		// Generate exponentially distributed data.
		data[i] = (-1.0 / intensity) * math.Log(random())
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

func TestRandomConcurrentSafe(t *testing.T) {
	const goRoutines = 10

	var wg sync.WaitGroup
	for n := 0; n < goRoutines; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = random()
		}()
	}

	wg.Wait()
}
