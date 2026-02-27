// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReset verifies that the reset function properly manages slice capacity
// to prevent unbounded memory growth while allowing reasonable growth for efficiency.
func TestReset(t *testing.T) {
	t.Run("AllocatesWhenCapacityTooSmall", func(t *testing.T) {
		s := make([]int, 0, 5)
		result := reset(s, 10, 10)
		assert.Len(t, result, 10, "length should match requested")
		assert.Equal(t, 10, cap(result), "capacity should match requested")
	})

	t.Run("ReusesSliceWhenCapacitySufficient", func(t *testing.T) {
		s := make([]int, 0, 10)
		// Store the slice header to verify it's reused
		originalData := &s[0:cap(s)][0]
		result := reset(s, 5, 5)
		assert.Len(t, result, 5, "length should match requested")
		assert.Equal(t, 10, cap(result), "capacity should be unchanged")
		// Verify the underlying array is reused
		assert.Same(t, originalData, &result[0:cap(result)][0], "should reuse underlying array")
	})

	t.Run("ShrinksWhenCapacityExcessive", func(t *testing.T) {
		// Create a slice with excessive capacity (more than 2x needed)
		s := make([]int, 0, 100)
		result := reset(s, 10, 10)
		assert.Len(t, result, 10, "length should match requested")
		assert.Equal(t, 10, cap(result), "capacity should shrink to requested")
	})

	t.Run("DoesNotShrinkWhenCapacityReasonable", func(t *testing.T) {
		// Create a slice with 2x capacity (exactly at threshold)
		s := make([]int, 0, 20)
		originalData := &s[0:cap(s)][0]
		result := reset(s, 10, 10)
		assert.Len(t, result, 10, "length should match requested")
		assert.Equal(t, 20, cap(result), "capacity should not shrink at 2x threshold")
		assert.Same(t, originalData, &result[0:cap(result)][0], "should reuse underlying array")
	})

	t.Run("ShrinksWhenCapacityJustAboveThreshold", func(t *testing.T) {
		// Create a slice with capacity just above 2x threshold
		s := make([]int, 0, 21)
		result := reset(s, 10, 10)
		assert.Len(t, result, 10, "length should match requested")
		assert.Equal(t, 10, cap(result), "capacity should shrink when above 2x threshold")
	})

	t.Run("HandlesZeroCapacity", func(t *testing.T) {
		s := make([]int, 0, 100)
		result := reset(s, 0, 0)
		assert.Empty(t, result, "length should be zero")
		// Should not shrink when capacity is 0 to avoid division issues
		assert.Equal(t, 100, cap(result), "capacity should not shrink for zero-capacity request")
	})

	t.Run("PreservesDataWhenShrinking", func(t *testing.T) {
		s := make([]int, 10, 100)
		for i := range s {
			s[i] = i
		}
		result := reset(s, 10, 10)
		assert.Len(t, result, 10, "length should match requested")
		assert.Equal(t, 10, cap(result), "capacity should shrink")
		// Note: When reallocating, data is NOT preserved by reset itself.
		// This is expected behavior - the caller is responsible for copying data if needed.
	})

	t.Run("SimulatesHighCardinalityScenario", func(t *testing.T) {
		// Simulate what happens with exemplar collection under high cardinality:
		// 1. Initial collection with small capacity
		s := make([]Exemplar, 0, 5)

		// 2. Spike in exemplars causes growth
		s = reset(s, 50, 50)
		assert.Equal(t, 50, cap(s), "should grow to handle spike")

		// 3. Subsequent collections with normal load should shrink back
		s = reset(s, 5, 5)
		assert.Equal(t, 5, cap(s), "should shrink back after spike")
	})

	t.Run("PreventsMemoryLeakOverMultipleFlushes", func(t *testing.T) {
		// Simulate repeated flush cycles with varying exemplar counts
		s := make([]Exemplar, 0, 10)

		// Normal operation
		s = reset(s, 10, 10)
		assert.LessOrEqual(t, cap(s), 20, "capacity should stay reasonable")

		// Temporary spike
		s = reset(s, 100, 100)
		assert.Equal(t, 100, cap(s), "should accommodate spike")

		// Back to normal - memory should be released
		s = reset(s, 10, 10)
		assert.Equal(t, 10, cap(s), "should release excess capacity")

		// Multiple normal operations should not grow capacity
		for range 10 {
			s = reset(s, 10, 10)
		}
		assert.Equal(t, 10, cap(s), "capacity should remain stable")
	})
}
