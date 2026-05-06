// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"math/bits"

	"go.opentelemetry.io/otel/attribute/internal/xxhash"
)

// LazyFilteredSet represents an attribute Set with a filter applied lazily.
// It is designed for performance-sensitive paths where filtering results
// should only be computed if needed.
type LazyFilteredSet struct {
	orig     Set
	mask     uint64
	distinct Distinct
	fallback *Set
	dropped  []KeyValue
}

// NewLazyFilteredSet creates a new LazyFilteredSet.
// It evaluates the filter exactly once for each attribute in the set.
func NewLazyFilteredSet(set Set, filter Filter) LazyFilteredSet {
	if filter == nil {
		return LazyFilteredSet{orig: set, distinct: set.Equivalent()}
	}
	if set.Len() == 0 {
		return LazyFilteredSet{orig: set, distinct: Distinct{hash: emptySet.hash}}
	}

	if set.Len() > 64 {
		filtered, dropped := set.Filter(filter)
		return LazyFilteredSet{orig: set, distinct: filtered.Equivalent(), fallback: &filtered, dropped: dropped}
	}

	h := xxhash.New()
	var mask uint64
	iter := set.Iter()
	i := 0
	hasAttributes := false
	for iter.Next() {
		kv := iter.Attribute()
		if filter(kv) {
			h = hashKV(h, kv)
			mask |= 1 << i
			hasAttributes = true
		}
		i++
	}

	var distinct Distinct
	if !hasAttributes {
		distinct = Distinct{hash: emptySet.hash}
	} else {
		distinct = Distinct{hash: h.Sum64()}
	}

	return LazyFilteredSet{orig: set, mask: mask, distinct: distinct}
}

// Distinct returns the hash of the filtered attributes.
func (s LazyFilteredSet) Distinct() Distinct {
	return s.distinct
}

// Filtered materializes and returns the filtered attribute set.
func (s LazyFilteredSet) Filtered() Set {
	if s.fallback != nil {
		return *s.fallback
	}
	if s.distinct == s.orig.Equivalent() {
		return s.orig
	}
	if s.mask == 0 {
		return emptySet
	}

	count := bits.OnesCount64(s.mask)
	kvs := make([]KeyValue, 0, count)

	switch d := s.orig.data.(type) {
	case [1]KeyValue:
		if s.mask&1 != 0 {
			kvs = append(kvs, d[0])
		}
	case [2]KeyValue:
		for i := range 2 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [3]KeyValue:
		for i := range 3 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [4]KeyValue:
		for i := range 4 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [5]KeyValue:
		for i := range 5 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [6]KeyValue:
		for i := range 6 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [7]KeyValue:
		for i := range 7 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [8]KeyValue:
		for i := range 8 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [9]KeyValue:
		for i := range 9 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	case [10]KeyValue:
		for i := range 10 {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, d[i])
			}
		}
	default:
		iter := s.orig.Iter()
		i := 0
		for iter.Next() {
			if s.mask&(1<<i) != 0 {
				kvs = append(kvs, iter.Attribute())
			}
			i++
		}
	}

	sSet := Set{
		hash: s.distinct.hash,
		data: computeDataFixed(kvs),
	}
	if sSet.data == nil {
		sSet.data = computeDataReflect(kvs)
	}
	return sSet
}

// Dropped materializes and returns the attributes that were filtered out.
func (s LazyFilteredSet) Dropped() []KeyValue {
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
	kvs := make([]KeyValue, 0, count)

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
