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

	t.Run("FilterSome", func(t *testing.T) {
		ls := attribute.NewLazyFilteredSet(s, func(kv attribute.KeyValue) bool {
			return string(kv.Key) == "k0" || string(kv.Key) == "k2"
		})
		expectedSet := attribute.NewSet(k0, k2)
		assert.Equal(t, expectedSet.Equivalent(), ls.Distinct())
		assert.Equal(t, expectedSet, ls.Filtered())
		assert.ElementsMatch(t, []attribute.KeyValue{k1}, ls.Dropped())
	})

	t.Run("EmptySet", func(t *testing.T) {
		ls := attribute.NewLazyFilteredSet(empty, func(_ attribute.KeyValue) bool { return true })
		assert.Equal(t, empty.Equivalent(), ls.Distinct())
		assert.Equal(t, empty, ls.Filtered())
		assert.Empty(t, ls.Dropped())
	})
}

func TestLazyFilteredSetFallback(t *testing.T) {
	var kvs []attribute.KeyValue
	for i := range 70 {
		kvs = append(kvs, attribute.Int(string(rune('a'+i%26))+string(rune('0'+i)), i))
	}
	s := attribute.NewSet(kvs...)

	filter := func(kv attribute.KeyValue) bool {
		// Keep even ones
		return kv.Value.AsInt64()%2 == 0
	}

	ls := attribute.NewLazyFilteredSet(s, filter)

	filtered, dropped := s.Filter(filter)

	assert.Equal(t, filtered.Equivalent(), ls.Distinct())
	assert.Equal(t, filtered, ls.Filtered())
	assert.ElementsMatch(t, dropped, ls.Dropped())
}

func TestLazyFilteredSetNonDeterministicFilter(t *testing.T) {
	k0 := attribute.String("k0", "v0")
	s := attribute.NewSet(k0)

	called := 0
	filter := func(_ attribute.KeyValue) bool {
		called++
		return called == 1 // True only on first call
	}

	ls := attribute.NewLazyFilteredSet(s, filter)

	filtered := ls.Filtered()
	dropped := ls.Dropped()

	assert.Equal(t, 1, called, "filter should be called exactly once per attribute")

	assert.Equal(t, attribute.NewSet(k0), filtered)
	assert.Empty(t, dropped)

	ls.Filtered()
	ls.Dropped()
	assert.Equal(t, 1, called, "filter should NOT be called again on materialization")
}

func TestLazyFilteredSetMediumSet(t *testing.T) {
	var kvs []attribute.KeyValue
	for i := range 20 {
		kvs = append(kvs, attribute.Int(fmt.Sprintf("k%d", i), i))
	}
	s := attribute.NewSet(kvs...)

	filter := func(kv attribute.KeyValue) bool {
		return kv.Value.AsInt64()%2 == 0
	}

	ls := attribute.NewLazyFilteredSet(s, filter)

	filtered, dropped := s.Filter(filter)

	assert.Equal(t, filtered.Equivalent(), ls.Distinct())
	assert.Equal(t, filtered, ls.Filtered())
	assert.ElementsMatch(t, dropped, ls.Dropped())
}
