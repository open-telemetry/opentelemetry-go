// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

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
		ls := attribute.NewLazyFilteredSet(s, nil)
		assert.Equal(t, s.Equivalent(), ls.Distinct())
		assert.Equal(t, s, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})

	t.Run("FilterAll", func(t *testing.T) {
		ls := attribute.NewLazyFilteredSet(s, func(_ attribute.KeyValue) bool { return true })
		assert.Equal(t, s.Equivalent(), ls.Distinct())
		assert.Equal(t, s, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})

	t.Run("FilterNone", func(t *testing.T) {
		ls := attribute.NewLazyFilteredSet(s, func(_ attribute.KeyValue) bool { return false })
		assert.Equal(t, empty.Equivalent(), ls.Distinct())
		assert.Equal(t, empty, ls.Filtered())
		assert.ElementsMatch(t, []attribute.KeyValue{k0, k1, k2}, ls.Dropped())
	})

	t.Run("EmptySet", func(t *testing.T) {
		ls := attribute.NewLazyFilteredSet(empty, func(_ attribute.KeyValue) bool { return true })
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
		// Sizes 1-10 use default filter (accept all) to guarantee full coverage of switch cases
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

			ls := attribute.NewLazyFilteredSet(s, fltr)

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
		{
			name: "SmallSet",
			size: 1,
		},
		{
			name: "LargeSet",
			size: 70,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var kvs []attribute.KeyValue
			for i := range tt.size {
				kvs = append(kvs, attribute.Int(fmt.Sprintf("k%d", i), i))
			}
			s := attribute.NewSet(kvs...)

			called := 0
			filter := func(_ attribute.KeyValue) bool {
				called++
				return called <= tt.size // True only on first pass
			}

			ls := attribute.NewLazyFilteredSet(s, filter)

			filtered := ls.Filtered()
			dropped := ls.Dropped()

			assert.Equal(t, tt.size, called, "filter should be called exactly once per attribute")

			assert.Equal(t, s, filtered)
			assert.Empty(t, dropped)

			ls.Filtered()
			ls.Dropped()
			assert.Equal(t, tt.size, called, "filter should NOT be called again on materialization")
		})
	}
}
