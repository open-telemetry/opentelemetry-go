// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHist(t *testing.T) {
	bounds := []float64{0, 100}
	t.Run("Int64", ReservoirTest[int64](func(int) (ReservoirProvider, int) {
		return HistogramReservoirProvider(bounds), len(bounds)
	}))

	t.Run("Float64", ReservoirTest[float64](func(int) (ReservoirProvider, int) {
		return HistogramReservoirProvider(bounds), len(bounds)
	}))
}

func TestHistogramReservoirConcurrentSafe(t *testing.T) {
	bounds := []float64{0, 100}
	t.Run("Int64", reservoirConcurrentSafeTest[int64](func(int) (ReservoirProvider, int) {
		return HistogramReservoirProvider(bounds), len(bounds)
	}))
	t.Run("Float64", reservoirConcurrentSafeTest[float64](func(int) (ReservoirProvider, int) {
		return HistogramReservoirProvider(bounds), len(bounds)
	}))
}

func TestHistogramReservoirSamplingCorrectness(t *testing.T) {
	bounds := []float64{10}
	// Bucket 0: < 10
	// Bucket 1: >= 10

	// We will offer 5 measurements to each bucket and repeat the experiment 1000 times.
	// We expect each of the 5 measurements to be sampled approximately 20% of the time.

	counts0 := make(map[int]int)
	counts1 := make(map[int]int)

	experiments := 1000
	offersPerBucket := 5

	for range experiments {
		r := NewHistogramReservoir(bounds)

		// Interleave offers to test that sharded locks work and state is independent
		for i := 1; i <= offersPerBucket; i++ {
			r.Offer(t.Context(), staticTime, NewValue[float64](float64(i)), nil)    // falls in bucket 0
			r.Offer(t.Context(), staticTime, NewValue[float64](float64(i+10)), nil) // falls in bucket 1
		}

		var dest []Exemplar
		r.Collect(&dest)

		for _, ex := range dest {
			v := int(ex.Value.Float64())
			if v <= 5 {
				counts0[v]++
			} else {
				counts1[v-10]++
			}
		}
	}

	// Assert that each item was sampled at least some times (not 0).
	// And that it is close to the expected 200 counts.
	expectedCount := experiments / offersPerBucket
	delta := float64(experiments) * 0.05 // 5% margin

	for i := 1; i <= offersPerBucket; i++ {
		assert.InDelta(t, float64(expectedCount), float64(counts0[i]), delta, "Bucket 0 item %d", i)
		assert.InDelta(t, float64(expectedCount), float64(counts1[i]), delta, "Bucket 1 item %d", i)
	}
}
