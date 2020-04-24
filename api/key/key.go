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

package key

import (
	"fmt"

	"go.opentelemetry.io/otel/api/core"
)

// New creates a new key with a passed name.
func New(name string) core.Key {
	return core.Key(name)
}

// Bool creates a new key-value pair with a passed name and a bool
// value.
func Bool(k string, v bool) core.KeyValue {
	return New(k).Bool(v)
}

// Int64 creates a new key-value pair with a passed name and an int64
// value.
func Int64(k string, v int64) core.KeyValue {
	return New(k).Int64(v)
}

// Uint64 creates a new key-value pair with a passed name and a uint64
// value.
func Uint64(k string, v uint64) core.KeyValue {
	return New(k).Uint64(v)
}

// Float64 creates a new key-value pair with a passed name and a float64
// value.
func Float64(k string, v float64) core.KeyValue {
	return New(k).Float64(v)
}

// Int32 creates a new key-value pair with a passed name and an int32
// value.
func Int32(k string, v int32) core.KeyValue {
	return New(k).Int32(v)
}

// Uint32 creates a new key-value pair with a passed name and a uint32
// value.
func Uint32(k string, v uint32) core.KeyValue {
	return New(k).Uint32(v)
}

// Float32 creates a new key-value pair with a passed name and a float32
// value.
func Float32(k string, v float32) core.KeyValue {
	return New(k).Float32(v)
}

// String creates a new key-value pair with a passed name and a string
// value.
func String(k, v string) core.KeyValue {
	return New(k).String(v)
}

// String creates a new key-value pair with a passed name and a string
// value.
func Stringer(k string, v fmt.Stringer) core.KeyValue {
	return New(k).String(v.String())
}

// Int creates a new key-value pair instance with a passed name and
// either an int32 or an int64 value, depending on whether the int
// type is 32 or 64 bits wide.
func Int(k string, v int) core.KeyValue {
	return New(k).Int(v)
}

// Uint creates a new key-value pair instance with a passed name and
// either an uint32 or an uint64 value, depending on whether the uint
// type is 32 or 64 bits wide.
func Uint(k string, v uint) core.KeyValue {
	return New(k).Uint(v)
}

// Infer creates a new key-value pair instance with a passed name and
// automatic type inference. This is slower, and not type-safe.
func Infer(k string, value interface{}) core.KeyValue {
	switch v := value.(type) {
	case bool:
		return Bool(k, v)
	case int:
		return Int(k, v)
	case uint:
		return Uint(k, v)
	case int32:
		return Int32(k, v)
	case int64:
		return Int64(k, v)
	case uint32:
		return Uint32(k, v)
	case uint64:
		return Uint64(k, v)
	case float32:
		return Float32(k, v)
	case float64:
		return Float64(k, v)
	case string:
		return String(k, v)
	case fmt.Stringer:
		return Stringer(k, v)
	default:
		return String(k, fmt.Sprint(v))
	}
}
