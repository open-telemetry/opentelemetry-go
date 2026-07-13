// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestLazyFilteredAttributes_CalledExactlyOnce(t *testing.T) {
	kvs := []attribute.KeyValue{
		attribute.String("a", "1"),
		attribute.String("b", "2"),
		attribute.String("c", "3"),
		attribute.String("d", "4"),
	}
	orig := attribute.NewSet(kvs...)

	var calls atomic.Int32
	filter := func(kv attribute.KeyValue) bool {
		calls.Add(1)
		return kv.Key == "a" || kv.Key == "c"
	}

	l := newLazyFilteredAttributes(orig, filter)
	assert.Equal(t, int32(len(kvs)), calls.Load(), "filter must be evaluated exactly once per attribute on creation")

	// Access all 3 stages repeatedly and concurrently across multiple goroutines
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				_ = l.Distinct()
				_ = l.Set()
				_ = l.Dropped()
			}
		})
	}
	wg.Wait()

	assert.Equal(t, int32(len(kvs)), calls.Load(), "subsequent stage access must never re-evaluate filter")
}

func TestLazyFilteredAttributes_NilFilter(t *testing.T) {
	kvs := []attribute.KeyValue{
		attribute.String("a", "1"),
		attribute.String("b", "2"),
	}
	orig := attribute.NewSet(kvs...)
	l := newLazyFilteredAttributes(orig, nil)

	assert.Equal(t, orig.Equivalent(), l.Distinct())
	assert.Equal(t, orig, l.Set())
	assert.Empty(t, l.Dropped())
}

func TestLazyFilteredAttributes_Correctness(t *testing.T) {
	t.Run("SubsetKept", func(t *testing.T) {
		orig := attribute.NewSet(
			attribute.String("keep1", "v1"),
			attribute.String("drop1", "v2"),
			attribute.String("keep2", "v3"),
		)
		filter := func(kv attribute.KeyValue) bool {
			return kv.Key == "keep1" || kv.Key == "keep2"
		}
		l := newLazyFilteredAttributes(orig, filter)

		wantSet := attribute.NewSet(
			attribute.String("keep1", "v1"),
			attribute.String("keep2", "v3"),
		)
		assert.Equal(t, wantSet.Equivalent(), l.Distinct())
		assert.Equal(t, wantSet, l.Set())
		assert.Equal(t, []attribute.KeyValue{attribute.String("drop1", "v2")}, l.Dropped())
	})

	t.Run("AllDropped", func(t *testing.T) {
		orig := attribute.NewSet(
			attribute.String("a", "1"),
			attribute.String("b", "2"),
		)
		filter := func(attribute.KeyValue) bool { return false }
		l := newLazyFilteredAttributes(orig, filter)

		wantSet := attribute.NewSet()
		assert.Equal(t, wantSet.Equivalent(), l.Distinct())
		assert.Equal(t, wantSet, l.Set())
		assert.ElementsMatch(t, orig.ToSlice(), l.Dropped())
	})

	t.Run("LargeSetMoreThan64", func(t *testing.T) {
		var kvs []attribute.KeyValue
		for i := range 75 {
			kvs = append(kvs, attribute.Int("k"+fmt.Sprint(i), i))
		}
		orig := attribute.NewSet(kvs...)
		filter := func(kv attribute.KeyValue) bool {
			return kv.Value.AsInt64()%2 == 0
		}
		l := newLazyFilteredAttributes(orig, filter)

		expectedKept, _ := orig.Filter(filter)
		assert.Equal(t, expectedKept.Equivalent(), l.Distinct())
		assert.Equal(t, expectedKept, l.Set())
		assert.Len(t, l.Dropped(), 75-expectedKept.Len())
	})

	t.Run("LargeSetMoreThan64_BigMaskNil", func(t *testing.T) {
		var kvs []attribute.KeyValue
		for i := range 75 {
			kvs = append(kvs, attribute.Int("k"+fmt.Sprint(i), i))
		}
		orig := attribute.NewSet(kvs...)
		filter := func(kv attribute.KeyValue) bool {
			// Only keep index < 64 (specifically only the first attribute k0) so bigMask is nil
			return kv.Key == "k0"
		}
		l := newLazyFilteredAttributes(orig, filter)

		expectedKept, _ := orig.Filter(filter)
		assert.Equal(t, expectedKept.Equivalent(), l.Distinct())
		assert.Equal(t, expectedKept, l.Set())
		assert.Len(t, l.Dropped(), 75-expectedKept.Len())
	})

	t.Run("LargeSetMoreThan64_AllKeptZeroAllocs", func(t *testing.T) {
		var kvs []attribute.KeyValue
		for i := range 75 {
			kvs = append(kvs, attribute.Int("k"+fmt.Sprint(i), i))
		}
		orig := attribute.NewSet(kvs...)
		filter := func(attribute.KeyValue) bool { return true }

		allocs := testing.AllocsPerRun(10, func() {
			_ = newLazyFilteredAttributes(orig, filter)
		})
		assert.Zero(t, allocs, "newLazyFilteredAttributes on >64 all-kept set must perform zero allocations")

		l := newLazyFilteredAttributes(orig, filter)
		assert.False(t, l.HasDroppedAttributes())
		assert.Equal(t, orig.Equivalent(), l.Distinct())
		assert.Equal(t, orig, l.Set())
		assert.Empty(t, l.Dropped())
	})

	t.Run("LargeSetMoreThan64_BackfillBigMask", func(t *testing.T) {
		var kvs []attribute.KeyValue
		for i := range 80 {
			kvs = append(kvs, attribute.Int("k"+fmt.Sprint(i), i))
		}
		orig := attribute.NewSet(kvs...)
		filter := func(kv attribute.KeyValue) bool {
			// Drop attribute at index 10 and index 66, keep everything else including 64, 65, 67, etc.
			return kv.Key != "k10" && kv.Key != "k66"
		}
		l := newLazyFilteredAttributes(orig, filter)
		assert.True(t, l.HasDroppedAttributes())

		expectedKept, _ := orig.Filter(filter)
		assert.Equal(t, expectedKept.Equivalent(), l.Distinct())
		assert.Equal(t, expectedKept, l.Set())
		assert.Len(t, l.Dropped(), 2)
	})

	t.Run("LargeSetMoreThan64_SuffixDropsAfter63", func(t *testing.T) {
		var kvs []attribute.KeyValue
		for i := range 80 {
			kvs = append(kvs, attribute.Int("k"+fmt.Sprint(i), i))
		}
		orig := attribute.NewSet(kvs...)
		filter := func(kv attribute.KeyValue) bool {
			// Keep all attributes up to k69, drop only the suffix (k70 to k79) after that.
			return kv.Value.AsInt64() < 70
		}
		l := newLazyFilteredAttributes(orig, filter)
		assert.True(t, l.HasDroppedAttributes())

		expectedKept, _ := orig.Filter(filter)
		assert.Equal(t, expectedKept.Equivalent(), l.Distinct())
		assert.Equal(t, expectedKept, l.Set())
		assert.Len(t, l.Dropped(), 10)
	})
}
