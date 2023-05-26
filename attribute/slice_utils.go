// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// hashKVs returns a new FNV-1a hash of kvs.
func hashKVs(kvs *[]KeyValue) fnv.Hash {
	if kvs == nil {
		return fnv.New()
	}

	h := fnv.New()
	for _, kv := range *kvs {
		h = hashKV(h, kv)
	}
	return h
}

// hashKV returns the FNV-1a hash of kv with h as the base.
func hashKV(h fnv.Hash, kv KeyValue) fnv.Hash {
	h = h.String(string(kv.Key))

	switch kv.Value.Type() {
	case BOOL:
		h = h.Bool(kv.Value.AsBool())
	case INT64:
		h = h.Int64(kv.Value.AsInt64())
	case FLOAT64:
		h = h.Float64(kv.Value.AsFloat64())
	case STRING:
		h = h.String(kv.Value.AsString())
	case BOOLSLICE:
		// Avoid allocating a new []bool with AsBoolSlice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Bool(rv.Index(i).Bool())
		}
	case INT64SLICE:
		// Avoid allocating a new []int64 with AsInt64Slice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Int64(rv.Index(i).Int())
		}
	case FLOAT64SLICE:
		// Avoid allocating a new []float64 with AsFloat64Slice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Float64(rv.Index(i).Float())
		}
	case STRINGSLICE:
		// Avoid allocating a new []string with AsStringSlice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.String(rv.Index(i).String())
		}
	default:
		// Logging is an alternative, but using the internal logger here
		// causes an import cycle so it is not done.
		v := kv.Value.AsInterface()
		msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
		panic(msg)
	}
	return h
}

// equalKVs returns if the sorted slices of []KeyValue a and b are equal.
func equalKVs(aPtr, bPtr *[]KeyValue) bool {
	a, b := *aPtr, *bPtr
	if len(a) != len(b) {
		return false
	}

	for i, aKV := range a {
		bKV := b[i]

		if aKV.Key != bKV.Key {
			return false
		}

		aVal, bVal := aKV.Value, bKV.Value
		if aVal.Type() != bVal.Type() {
			return false
		}

		switch aVal.Type() {
		case BOOL, INT64, FLOAT64:
			if aVal.numeric != bVal.numeric {
				return false
			}
		case STRING:
			if aVal.stringly != bVal.stringly {
				return false
			}
		case BOOLSLICE, INT64SLICE, FLOAT64SLICE, STRINGSLICE:
			if aVal.slice != bVal.slice {
				return false
			}
		default:
			// Logging is an alternative, but using the internal logger here
			// causes an import cycle so it is not done.
			v := aVal.AsInterface()
			msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
			panic(msg)
		}
	}

	return true
}
