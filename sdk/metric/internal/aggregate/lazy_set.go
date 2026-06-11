// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math/bits"

	"go.opentelemetry.io/otel/attribute"
)

// lazyFilteredSet represents an attribute Set with a filter applied lazily.
// It is designed for performance-sensitive paths where filtering results
// should only be computed if needed.
type lazyFilteredSet struct {
	orig     attribute.Set
	mask     uint64
	distinct attribute.Distinct
	fallback *attribute.Set
	dropped  []attribute.KeyValue
}

// newLazyFilteredSet creates a new lazyFilteredSet.
// It evaluates the filter exactly once for each attribute in the set.
func newLazyFilteredSet(set attribute.Set, filter attribute.Filter) lazyFilteredSet {
	if filter == nil {
		return lazyFilteredSet{orig: set, distinct: set.Equivalent()}
	}
	if set.Len() == 0 {
		return lazyFilteredSet{orig: set, distinct: set.Equivalent()}
	}

	if set.Len() > 64 {
		filtered, dropped := set.Filter(filter)
		return lazyFilteredSet{orig: set, distinct: filtered.Equivalent(), fallback: &filtered, dropped: dropped}
	}

	distinct, mask := attribute.NewDistinctFiltered(set, filter)
	return lazyFilteredSet{orig: set, mask: mask, distinct: distinct}
}

// Distinct returns the hash of the filtered attributes.
func (s lazyFilteredSet) Distinct() attribute.Distinct {
	return s.distinct
}

// Filtered materializes and returns the filtered attribute set.
func (s lazyFilteredSet) Filtered() attribute.Set {
	if s.fallback != nil {
		return *s.fallback
	}
	if s.distinct == s.orig.Equivalent() {
		return s.orig
	}
	if s.mask == 0 {
		return attribute.NewSet()
	}

	count := bits.OnesCount64(s.mask)
	kvs := make([]attribute.KeyValue, 0, count)

	iter := s.orig.Iter()
	i := 0
	for iter.Next() {
		if s.mask&(1<<i) != 0 {
			kvs = append(kvs, iter.Attribute())
		}
		i++
	}

	return attribute.NewSet(kvs...)
}

// Dropped materializes and returns the attributes that were filtered out.
func (s lazyFilteredSet) Dropped() []attribute.KeyValue {
	if s.fallback != nil {
		return s.dropped
	}
	if s.distinct == s.orig.Equivalent() {
		return nil
	}
	if s.mask == 0 {
		return s.orig.ToSlice()
	}

	count := s.orig.Len() - bits.OnesCount64(s.mask)
	kvs := make([]attribute.KeyValue, 0, count)

	iter := s.orig.Iter()
	i := 0
	for iter.Next() {
		if s.mask&(1<<i) == 0 {
			kvs = append(kvs, iter.Attribute())
		}
		i++
	}
	return kvs
}
