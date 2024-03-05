// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// hashKVs returns a new FNV-1a hash of kvs.
func hashKVs(kvs []KeyValue) fnv.Hash {
	h := fnv.New()
	for _, kv := range kvs {
		h = hashKV(h, kv)
	}
	return h
}

// hashKV returns the FNV-1a hash of kv with h as the base.
func hashKV(h fnv.Hash, kv KeyValue) fnv.Hash {
	h = h.String(string(kv.Key))

	switch kv.Value.Type() {
	case BOOL:
		h = h.String("b")
		h = h.Uint64(kv.Value.numeric)
	case INT64:
		h = h.String("i")
		h = h.Uint64(kv.Value.numeric)
	case FLOAT64:
		h = h.String("f")
		// Assumes numeric stored with math.Float64bits.
		h = h.Uint64(kv.Value.numeric)
	case STRING:
		h = h.String("s")
		h = h.String(kv.Value.stringly)
	case BOOLSLICE:
		// Differentiate between bool and [1]bool.
		h = h.String("[]b")
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Bool(rv.Index(i).Bool())
		}
	case INT64SLICE:
		// Differentiate between int64 and [1]int64.
		h = h.String("[]i")
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Int64(rv.Index(i).Int())
		}
	case FLOAT64SLICE:
		// Differentiate between float64 and [1]float64.
		h = h.String("[]f")
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Float64(rv.Index(i).Float())
		}
	case STRINGSLICE:
		// Differentiate between string and [1]string.
		h = h.String("[]s")
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.String(rv.Index(i).String())
		}
	case INVALID:
	default:
		// Logging is an alternative, but using the internal logger here
		// causes an import cycle so it is not done.
		v := kv.Value.AsInterface()
		msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
		panic(msg)
	}
	return h
}
