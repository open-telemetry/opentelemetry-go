// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"math"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewFixedSizeReservoir(t *testing.T) {
	t.Run("Int64", ReservoirTest[int64](func(n int) (ReservoirProvider, int) {
		provider := FixedSizeReservoirProvider(n)
		return provider, int(provider(attribute.NewSet()).(*FixedSizeReservoir).k)
	}))

	t.Run("Float64", ReservoirTest[float64](func(n int) (ReservoirProvider, int) {
		provider := FixedSizeReservoirProvider(n)
		return provider, int(provider(attribute.NewSet()).(*FixedSizeReservoir).k)
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

func TestFixedSizeReservoirConcurrentSafe(t *testing.T) {
	t.Run("Int64", reservoirConcurrentSafeTest[int64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
	t.Run("Float64", reservoirConcurrentSafeTest[float64](func(n int) (ReservoirProvider, int) {
		return FixedSizeReservoirProvider(n), n
	}))
}

func TestNextTrackerAtomics(t *testing.T) {
	capacity := uint32(10)
	nt := newNextTracker(capacity)
	nt.setCountAndNext(0, 11)
	count, next := nt.incrementCount()
	assert.Equal(t, uint32(0), count)
	assert.Equal(t, uint32(11), next)
	count, secondNext := nt.incrementCount()
	assert.Equal(t, uint32(1), count)
	assert.Equal(t, next, secondNext)
	nt.setCountAndNext(50, 100)
	count, next = nt.incrementCount()
	assert.Equal(t, uint32(50), count)
	assert.Equal(t, uint32(100), next)
}

func TestNewFixedSizeReservoirConcurrentSamplingCorrectness(t *testing.T) {
	sampleSize := 1000
	workers := 10
	itemsPerWorker := 10000
	totalItems := workers * itemsPerWorker

	// The first half of the data is positive, and the second half is negative.
	// This test is designed to ensure the reservoir doesn't stall early.
	data := make([]float64, totalItems)
	for i := range data {
		if i < totalItems/2 {
			data[i] = 1.0
		} else {
			data[i] = -1.0
		}
	}

	r := NewFixedSizeReservoir(sampleSize)

	var wg sync.WaitGroup
	var idx atomic.Uint32

	for range workers {
		wg.Go(func() {
			for {
				curr := idx.Add(1) - 1
				if curr >= uint32(len(data)) {
					break
				}
				r.Offer(t.Context(), staticTime, NewValue(data[curr]), nil)
			}
		})
	}
	wg.Wait()

	var negCount int
	for i := range r.measurements {
		if r.measurements[i].Value.Float64() < 0 {
			negCount++
		}
	}

	// We expect roughly half of the sampled items to be negative.
	// Allow a wide delta because of randomness and concurrency. If negCount is
	// zero, then we stalled before any negative observations were made.
	assert.Greater(t, negCount, 100, "Should have sampled items from the second half of the interval")
}
