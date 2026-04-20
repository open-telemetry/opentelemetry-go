// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewHistogramReservoirSamplingCorrectness(t *testing.T) {
	sampleSize := 100

	u := rand.Uint32()
	seed := [32]byte{byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24)}
	t.Logf("rng seed: %x", seed)
	rng := rand.New(rand.NewChaCha8(seed))

	data := make([]float64, sampleSize)
	for i := range data {
		data[i] = rng.Float64()
	}
	// Sort to test position bias.
	slices.Sort(data)

	bounds := []float64{1000} // Large bound to put all in one bucket.

	// We run multiple times because reservoir size is 1 per bucket.
	runs := 20
	allLast := true

	for range runs {
		r := NewHistogramReservoir(bounds)
		for _, value := range data {
			r.Offer(t.Context(), staticTime, NewValue(value), nil)
		}

		var dest []Exemplar
		r.Collect(&dest)

		require.Len(t, dest, 1, "number of collected exemplars")
		if dest[0].Value.Float64() != float64(sampleSize) {
			allLast = false
			break
		}
	}

	// Check that not all collected exemplars are the last offered value,
	// ensuring no bias in our random sampling algorithm.
	assert.False(t, allLast, "expected at least one exemplar to not be the last offered value")
}
