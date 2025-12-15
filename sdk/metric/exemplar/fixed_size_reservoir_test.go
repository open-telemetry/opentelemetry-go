// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFixedSizeReservoir(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))

	t.Run("Float64", ReservoirTest[float64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
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

func TestNextTrackerAtomics(t *testing.T) {
	capacity := 10
	nt := newNextTracker(capacity)
	nt.setCountAndNext(0, 11)
	count, next := nt.incrementCount()
	assert.Equal(t, uint64(0), count)
	assert.Equal(t, uint64(11), next)
	count, secondNext := nt.incrementCount()
	assert.Equal(t, uint64(1), count)
	assert.Equal(t, next, secondNext)
	nt.setCountAndNext(50, 100)
	count, next = nt.incrementCount()
	assert.Equal(t, uint64(50), count)
	assert.Equal(t, uint64(100), next)
}
