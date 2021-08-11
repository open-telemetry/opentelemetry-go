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
	"encoding/json"
	"fmt"
	"reflect"
)

// KeyValue holds a key and value pair.
type KeyValue struct {
	Key   Key
	Value Value
}

// Valid returns if kv is a valid OpenTelemetry attribute.
func (kv KeyValue) Valid() bool {
	return kv.Key != "" && kv.Value.Type() != INVALID
}

// Bool creates a KeyValue with a BOOL Value type.
func Bool(k string, v bool) KeyValue {
	return Key(k).Bool(v)
}

// BoolSlice creates a KeyValue with a BOOLSLICE Value type.
func BoolSlice(k string, v []bool) KeyValue {
	return Key(k).BoolSlice(v)
}

// Int creates a KeyValue with an INT64 Value type.
func Int(k string, v int) KeyValue {
	return Key(k).Int(v)
}

// IntSlice creates a KeyValue with an INT64SLICE Value type.
func IntSlice(k string, v []int) KeyValue {
	return Key(k).IntSlice(v)
}

// Int64 creates a KeyValue with an INT64 Value type.
func Int64(k string, v int64) KeyValue {
	return Key(k).Int64(v)
}

// Int64Slice creates a KeyValue with an INT64SLICE Value type.
func Int64Slice(k string, v []int64) KeyValue {
	return Key(k).Int64Slice(v)
}

// Float64 creates a KeyValue with a FLOAT64 Value type.
func Float64(k string, v float64) KeyValue {
	return Key(k).Float64(v)
}

// Float64Slice creates a KeyValue with a FLOAT64SLICE Value type.
func Float64Slice(k string, v []float64) KeyValue {
	return Key(k).Float64Slice(v)
}

// String creates a KeyValue with a STRING Value type.
func String(k, v string) KeyValue {
	return Key(k).String(v)
}

// StringSlice creates a KeyValue with a STRINGSLICE Value type.
func StringSlice(k string, v []string) KeyValue {
	return Key(k).StringSlice(v)
}

// Stringer creates a new key-value pair with a passed name and a string
// value generated by the passed Stringer interface.
func Stringer(k string, v fmt.Stringer) KeyValue {
	return Key(k).String(v.String())
}

// Array creates a new key-value pair with a passed name and a array.
// Only arrays of primitive type are supported.
//
// Deprecated: Use the typed *Slice functions instead.
func Array(k string, v interface{}) KeyValue {
	return Key(k).Array(v)
}

// Any creates a new key-value pair instance with a passed name and
// automatic type inference. This is slower, and not type-safe.
func Any(k string, value interface{}) KeyValue {
	if value == nil {
		return String(k, "<nil>")
	}

	if stringer, ok := value.(fmt.Stringer); ok {
		return String(k, stringer.String())
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Array:
		return Array(k, value)
	case reflect.Slice:
		switch reflect.TypeOf(value).Elem().Kind() {
		case reflect.Bool:
			return BoolSlice(k, rv.Interface().([]bool))
		case reflect.Int:
			return IntSlice(k, rv.Interface().([]int))
		case reflect.Int64:
			return Int64Slice(k, rv.Interface().([]int64))
		case reflect.Float64:
			return Float64Slice(k, rv.Interface().([]float64))
		case reflect.String:
			return StringSlice(k, rv.Interface().([]string))
		default:
			return Array(k, value)
		}
	case reflect.Bool:
		return Bool(k, rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int64(k, rv.Int())
	case reflect.Float64:
		return Float64(k, rv.Float())
	case reflect.String:
		return String(k, rv.String())
	}
	if b, err := json.Marshal(value); b != nil && err == nil {
		return String(k, string(b))
	}
	return String(k, fmt.Sprint(value))
}
