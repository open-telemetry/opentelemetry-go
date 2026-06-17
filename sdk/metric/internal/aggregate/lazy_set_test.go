// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestLazyFilteredSet(t *testing.T) {
	k0 := attribute.String("k0", "v0")
	k1 := attribute.String("k1", "v1")
	k2 := attribute.String("k2", "v2")
	s := attribute.NewSet(k0, k1, k2)
	empty := attribute.NewSet()

	t.Run("NilFilter", func(t *testing.T) {
		ls := newLazyFilteredSet(s, nil)
		assert.Equal(t, s.Equivalent(), ls.Distinct())
		assert.Equal(t, s, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})

	t.Run("FilterAll", func(t *testing.T) {
		ls := newLazyFilteredSet(s, func(_ attribute.KeyValue) bool { return true })
		assert.Equal(t, s.Equivalent(), ls.Distinct())
		assert.Equal(t, s, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})

	t.Run("FilterNone", func(t *testing.T) {
		ls := newLazyFilteredSet(s, func(_ attribute.KeyValue) bool { return false })
		assert.Equal(t, empty.Equivalent(), ls.Distinct())
		assert.Equal(t, empty, ls.Filtered())
		assert.ElementsMatch(t, []attribute.KeyValue{k0, k1, k2}, ls.Dropped())
	})

	t.Run("EmptySet", func(t *testing.T) {
		ls := newLazyFilteredSet(empty, func(_ attribute.KeyValue) bool { return true })
		assert.Equal(t, empty.Equivalent(), ls.Distinct())
		assert.Equal(t, empty, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})
}

func TestLazyFilteredSetVariousSizes(t *testing.T) {
	testCases := []struct {
		name   string
		size   int
		filter func(attribute.KeyValue) bool
	}{
		// Sizes 1-10 use the default filter (accept even values) to exercise the <=64 mask-based path across small set sizes
		{name: "Size1", size: 1},
		{name: "Size2", size: 2},
		{name: "Size3", size: 3},
		{name: "Size4", size: 4},
		{name: "Size5", size: 5},
		{name: "Size6", size: 6},
		{name: "Size7", size: 7},
		{name: "Size8", size: 8},
		{name: "Size9", size: 9},
		{name: "Size10", size: 10},
		// Specific boundary tests with a realistic filter
		{
			name:   "SmallSetFiltered",
			size:   3,
			filter: func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 },
		},
		{
			name:   "MediumSetFiltered",
			size:   25,
			filter: func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 },
		},
		{
			name:   "LargeSetFiltered",
			size:   70,
			filter: func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 },
		},
		{
			name:   "Size64Boundary",
			size:   64,
			filter: func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 },
		},
		{
			name:   "Size65Boundary",
			size:   65,
			filter: func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 },
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var kvs []attribute.KeyValue
			for i := range tt.size {
				kvs = append(kvs, attribute.Int(fmt.Sprintf("k%d", i), i))
			}
			s := attribute.NewSet(kvs...)

			fltr := tt.filter
			if fltr == nil {
				fltr = func(kv attribute.KeyValue) bool { return kv.Value.AsInt64()%2 == 0 }
			}

			ls := newLazyFilteredSet(s, fltr)

			filtered, dropped := s.Filter(fltr)

			assert.Equal(t, filtered.Equivalent(), ls.Distinct())
			assert.Equal(t, filtered, ls.Filtered())
			assert.ElementsMatch(t, dropped, ls.Dropped())
		})
	}
}

func TestLazyFilteredSetInconsistentFilter(t *testing.T) {
	testCases := []struct {
		name string
		size int
	}{
		{name: "SmallSet", size: 5},
		{name: "LargeSet", size: 65},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var kvs []attribute.KeyValue
			for i := range tc.size {
				kvs = append(kvs, attribute.Int(fmt.Sprintf("k%d", i), i))
			}
			s := attribute.NewSet(kvs...)

			called := 0
			// Filter that accepts even values
			filter := func(kv attribute.KeyValue) bool {
				called++
				return kv.Value.AsInt64()%2 == 0
			}

			ls := newLazyFilteredSet(s, filter)

			filtered := ls.Filtered()
			dropped := ls.Dropped()

			assert.Equal(t, tc.size, called, "filter should be called exactly once per attribute")

			// Verify filtered and dropped sets are consistent
			expectedFilteredKVs := []attribute.KeyValue{}
			expectedDroppedKVs := []attribute.KeyValue{}
			// NewSet sorts the KVs, so we need to match the sorted order for expectedDroppedKVs
			// actually, expectedFiltered will be sorted by NewSet, but expectedDroppedKVs we build manually.
			// Set.Iter() returns in sorted order.
			iter := s.Iter()
			for iter.Next() {
				kv := iter.Attribute()
				if kv.Value.AsInt64()%2 == 0 {
					expectedFilteredKVs = append(expectedFilteredKVs, kv)
				} else {
					expectedDroppedKVs = append(expectedDroppedKVs, kv)
				}
			}
			expectedFiltered := attribute.NewSet(expectedFilteredKVs...)

			assert.Equal(t, expectedFiltered, filtered)
			assert.Equal(t, expectedDroppedKVs, dropped)

			// Call again to ensure no more evaluations
			ls.Filtered()
			ls.Dropped()
			assert.Equal(t, tc.size, called, "filter should NOT be called again on materialization")
		})
	}
}
