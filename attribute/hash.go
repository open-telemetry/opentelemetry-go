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

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// hashKVs returns a new FNV-1a hash of kvs.
func hashKVs(kvs []KeyValue) fnv.Hash {
	if kvs == nil {
		return fnv.New()
	}

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
		h = h.Bool(kv.Value.asBool())
	case INT64:
		h = h.Int64(kv.Value.asInt64())
	case FLOAT64:
		h = h.Float64(kv.Value.asFloat64())
	case STRING:
		h = h.String(kv.Value.asString())
	case BOOLSLICE:
		for _, v := range kv.Value.asBoolSlice() {
			h = h.Bool(v)
		}
	case INT64SLICE:
		for _, v := range kv.Value.asInt64Slice() {
			h = h.Int64(v)
		}
	case FLOAT64SLICE:
		for _, v := range kv.Value.asFloat64Slice() {
			h = h.Float64(v)
		}
	case STRINGSLICE:
		for _, v := range kv.Value.asStringSlice() {
			h = h.String(v)
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
