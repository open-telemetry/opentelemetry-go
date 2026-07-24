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
	// A single bucket (-Inf, 100] so all offered measurements fall into it.
	bounds := []float64{100}
	r := NewHistogramReservoir(bounds)

	const (
		N = 10    // Items offered per run
		M = 20000 // Number of runs
	)

	var counts [N]int
	var dest []Exemplar

	for range M {
		for j := range N {
			r.Offer(t.Context(), staticTime, NewValue(int64(j)), nil)
		}
		r.Collect(&dest)
		require.Len(t, dest, 1)
		pos := int(dest[0].Value.Int64())
		require.GreaterOrEqual(t, pos, 0)
		require.Less(t, pos, N)
		counts[pos]++
		dest = dest[:0]
	}

	// Chi-square goodness-of-fit test against the discrete uniform distribution.
	// Expected frequency for each position under uniform sampling: E = M / N.
	expected := float64(M) / float64(N)
	var chiSquare float64
	for _, count := range counts {
		diff := float64(count) - expected
		chiSquare += (diff * diff) / expected
	}

	// With N=10 categories (df = N-1 = 9), the critical value of Chi-square at
	// alpha = 0.00006 is 35.0. This prevents test flakiness while rejecting
	// non-uniform sampling distributions (such as symmetric U-shaped or skewed
	// distributions).
	assert.Less(
		t,
		chiSquare,
		35.0,
		"Chi-square goodness-of-fit statistic %v exceeds critical value 35.0 (df=%d) for uniform distribution",
		chiSquare,
		N-1,
	)
}
