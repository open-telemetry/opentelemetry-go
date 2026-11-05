// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
)

// rawSet mirrors attribute.Set. It is used only to read the immutable backing
// array without calling Set.Get on the per-measurement path.
type rawSet struct {
	hash uint64
	data any
}

// SetMayHaveComposite reports false only when set is known to hold no SLICE or
// MAP value, so [SetDedup] cannot change it.
func SetMayHaveComposite(set *attribute.Set) bool {
	switch kvs := (*rawSet)(unsafe.Pointer(set)).data.(type) { //nolint:gosec // Read-only mirror of attribute.Set.
	case [0]attribute.KeyValue:
		return false
	case [1]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [2]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [3]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [4]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [5]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [6]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [7]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [8]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [9]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [10]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [11]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [12]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [13]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [14]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [15]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [16]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [17]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [18]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [19]attribute.KeyValue:
		return anyComposite(kvs[:])
	case [20]attribute.KeyValue:
		return anyComposite(kvs[:])
	default:
		// Rarer large sets keep the full SetDedup path.
		return true
	}
}

func anyComposite(kvs []attribute.KeyValue) bool {
	for i := range kvs {
		if t := kvs[i].Value.Type(); t == attribute.SLICE || t == attribute.MAP {
			return true
		}
	}
	return false
}
