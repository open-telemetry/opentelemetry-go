// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute/internal/xxhash"
)

// Type identifiers. These identifiers are hashed before the value of the
// corresponding type. This is done to distinguish values that are hashed with
// the same value representation (e.g. `int64(1)` and `true`, []int64{0} and
// int64(0)).
//
// These are all 8 byte length strings converted to a uint64 representation. A
// uint64 is used instead of the string directly as an optimization, it avoids
// the for loop in [xxhash] which adds minor overhead.
const (
	boolID         uint64 = 7953749933313450591 // "_boolean" (little endian)
	int64ID        uint64 = 7592915492740740150 // "64_bit_i" (little endian)
	float64ID      uint64 = 7376742710626956342 // "64_bit_f" (little endian)
	stringID       uint64 = 6874584755375207263 // "_string_" (little endian)
	boolSliceID    uint64 = 6875993255270243167 // "_[]bool_" (little endian)
	int64SliceID   uint64 = 3762322556277578591 // "_[]int64" (little endian)
	float64SliceID uint64 = 7308324551835016539 // "[]double" (little endian)
	stringSliceID  uint64 = 7453010373645655387 // "[]string" (little endian)
	byteSliceID    uint64 = 6874028470941080415 // "_[]byte_" (little endian)
	emptyID        uint64 = 7305809155345288421 // "__empty_" (little endian)
)

// hashKVs returns a new xxHash64 hash of kvs.
func hashKVs(kvs []KeyValue) uint64 {
	h := xxhash.New()
	for _, kv := range kvs {
		h = hashKV(h, kv)
	}
	return h.Sum64()
}

// hashKV returns the xxHash64 hash of kv with h as the base.
func hashKV(h xxhash.Hash, kv KeyValue) xxhash.Hash {
	h = h.String(string(kv.Key))

	switch kv.Value.Type() {
	case BOOL:
		h = h.Uint64(boolID)
		h = h.Uint64(kv.Value.numeric)
	case INT64:
		h = h.Uint64(int64ID)
		h = h.Uint64(kv.Value.numeric)
	case FLOAT64:
		h = h.Uint64(float64ID)
		// Assumes numeric stored with math.Float64bits.
		h = h.Uint64(kv.Value.numeric)
	case STRING:
		h = h.Uint64(stringID)
		h = h.String(kv.Value.stringly)
	case BOOLSLICE:
		h = h.Uint64(boolSliceID)
		for _, v := range kv.Value.asBoolSlice() {
			h = h.Bool(v)
		}
	case INT64SLICE:
		h = h.Uint64(int64SliceID)
		for _, v := range kv.Value.asInt64Slice() {
			h = h.Int64(v)
		}
	case FLOAT64SLICE:
		h = h.Uint64(float64SliceID)
		for _, v := range kv.Value.asFloat64Slice() {
			h = h.Float64(v)
		}
	case STRINGSLICE:
		h = h.Uint64(stringSliceID)
		for _, v := range kv.Value.asStringSlice() {
			h = h.String(v)
		}
	case BYTESLICE:
		h = h.Uint64(byteSliceID)
		h = h.String(kv.Value.stringly)
	case EMPTY:
		h = h.Uint64(emptyID)
	default:
		// Logging is an alternative, but using the internal logger here
		// causes an import cycle so it is not done.
		v := kv.Value.AsInterface()
		msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
		panic(msg)
	}
	return h
}
