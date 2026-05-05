// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
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

func TestHistogramReservoirTimeUnbiased(t *testing.T) {
	bounds := []float64{10}
	r := NewHistogramReservoir(bounds)

	const (
		N = 100   // Items per run
		M = 20000 // Number of runs
	)

	var sum float64
	var dest []Exemplar

	for range M {
		for j := 1; j <= N; j++ {
			val := 10.0 * float64(j) / float64(N)
			r.Offer(t.Context(), staticTime, NewValue(val), nil)
		}
		r.Collect(&dest)
		require.Len(t, dest, 1)
		sum += dest[0].Value.Float64()
		dest = dest[:0]
	}

	mean := sum / float64(M)
	expectedMean := 5.0 * float64(N+1) / float64(N)

	// Standard deviation of the discrete uniform distribution is approx 2.88.
	// Standard error of the mean for M=20000 is 2.88 / sqrt(20000) approx 0.02.
	// A delta of 0.1 is approx 5 standard errors, which makes flakiness extremely unlikely.
	assert.InDelta(t, expectedMean, mean, 0.1)
}
