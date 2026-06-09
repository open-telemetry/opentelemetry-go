// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFixedSizeReservoir(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))

	t.Run("Float64", ReservoirTest[float64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
}

func TestNewFixedSizeReservoirZeroSize(t *testing.T) {
	r := NewFixedSizeReservoir(0)
	require.NotNil(t, r)

	// Offer should be a no-op and not panic.
	r.Offer(t.Context(), staticTime, NewValue(float64(10)), nil)

	// Collect should leave dest empty.
	dest := []Exemplar{{}} // pre-filled sentinel
	r.Collect(&dest)
	assert.Empty(t, dest)
}

func TestNewFixedSizeReservoirSamplingCorrectness(t *testing.T) {
	intensity := 0.1
	sampleSize := 1000

	u := rand.Uint32()
	seed := [32]byte{byte(u), byte(u >> 8), byte(u >> 16), byte(u >> 24)}
	t.Logf("rng seed: %x", seed)
	rng := rand.New(rand.NewChaCha8(seed))

	data := make([]float64, sampleSize*1000)
	for i := range data {
		// Generate exponentially distributed data.
		data[i] = (-1.0 / intensity) * math.Log(rng.Float64())
	}
	// Sort to test position bias.
	slices.Sort(data)

	r := NewFixedSizeReservoir(sampleSize)
	for _, value := range data {
		r.Offer(t.Context(), staticTime, NewValue(value), nil)
	}

	var sum float64
	for i := range r.measurements {
		sum += r.measurements[i].Value.Float64()
	}
	mean := sum / float64(sampleSize)

	// Check the intensity/rate of the sampled distribution is preserved
	// ensuring no bias in our random sampling algorithm.
	assert.InDelta(t, 1/mean, intensity, 0.02) // Within 5σ.
}

func TestFixedSizeReservoirConcurrentSafe(t *testing.T) {
	t.Run("Int64", reservoirConcurrentSafeTest[int64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
	t.Run("Float64", reservoirConcurrentSafeTest[float64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
}

func TestFixedSizeReservoirSamplesAfterFilling(t *testing.T) {
	k := 1
	var sampledSecondItemCount int
	iterations := 10000
	for range iterations {
		r := NewFixedSizeReservoir(k)
		// Offer k items (1 item)
		r.Offer(t.Context(), staticTime, NewValue(float64(1)), nil)
		// Offer the k+1 item (2nd item)
		r.Offer(t.Context(), staticTime, NewValue(float64(2)), nil)

		var dest []Exemplar
		r.Collect(&dest)
		if len(dest) == 1 && dest[0].Value.Float64() == 2 {
			sampledSecondItemCount++
		}
	}
	rate := float64(sampledSecondItemCount) / float64(iterations)
	// For k=1, the probability of sampling the k+1 item is k/(k+1) = 1/2.
	// Expected rate is 0.5.
	// With 10000 iterations, the standard deviation of the count is:
	//   sqrt(10000 * 0.5 * 0.5) = 50.
	// 5 standard deviations is 250, which is 2.5% (0.025).
	assert.InDelta(t, 0.5, rate, 0.025, "should sample the second item with ~50% probability")
}
