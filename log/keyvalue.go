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

//go:generate stringer -type=Kind -trimprefix=Kind

package log // import "go.opentelemetry.io/otel/log"

// Kind is the kind of a [Value].
type Kind int

// Kind values.
const (
	KindEmpty Kind = iota
	KindBool
	KindFloat64
	KindInt64
	KindString
	KindBytes
	KindSlice
	KindMap
)

// A Value represents a structured log value.
type Value struct{} // TODO (#4914): implement.

// StringValue returns a new [Value] for a string.
func StringValue(v string) Value { return Value{} } // TODO (#4914): implement.

// IntValue returns a [Value] for an int.
func IntValue(v int) Value { return Value{} } // TODO (#4914): implement.

// Int64Value returns a [Value] for an int64.
func Int64Value(v int64) Value { return Value{} } // TODO (#4914): implement.

// Float64Value returns a [Value] for a float64.
func Float64Value(v float64) Value { return Value{} } // TODO (#4914): implement.

// BoolValue returns a [Value] for a bool.
func BoolValue(v bool) Value { //nolint:revive // Not a control flag.
	// TODO (#4914): implement.
	return Value{}
}

// BytesValue returns a [Value] for a byte slice. The passed slice must not be
// changed after it is passed.
func BytesValue(v []byte) Value { return Value{} } // TODO (#4914): implement.

// SliceValue returns a [Value] for a slice of [Value]. The passed slice must
// not be changed after it is passed.
func SliceValue(vs ...Value) Value { return Value{} } // TODO (#4914): implement.

// MapValue returns a new [Value] for a slice of key-value pairs. The passed
// slice must not be changed after it is passed.
func MapValue(kvs ...KeyValue) Value { return Value{} } // TODO (#4914): implement.

// AsAny returns the value held by v as an any.
func (v Value) AsAny() any { return nil } // TODO (#4914): implement

// AsString returns the value held by v as a string.
func (v Value) AsString() string { return "" } // TODO (#4914): implement

// AsInt64 returns the value held by v as an int64.
func (v Value) AsInt64() int64 { return 0 } // TODO (#4914): implement

// AsBool returns the value held by v as a bool.
func (v Value) AsBool() bool { return false } // TODO (#4914): implement

// AsFloat64 returns the value held by v as a float64.
func (v Value) AsFloat64() float64 { return 0 } // TODO (#4914): implement

// AsBytes returns the value held by v as a []byte.
func (v Value) AsBytes() []byte { return nil } // TODO (#4914): implement

// AsSlice returns the value held by v as a []Value.
func (v Value) AsSlice() []Value { return nil } // TODO (#4914): implement

// AsMap returns the value held by v as a []KeyValue.
func (v Value) AsMap() []KeyValue { return nil } // TODO (#4914): implement

// Kind returns the Kind of v.
func (v Value) Kind() Kind { return KindEmpty } // TODO (#4914): implement.

// Empty returns if v does not hold any value.
func (v Value) Empty() bool { return false } // TODO (#4914): implement

// Equal returns if v is equal to w.
func (v Value) Equal(w Value) bool { return false } // TODO (#4914): implement

// An KeyValue is a key-value pair used to represent a log attribute (a
// superset of [go.opentelemetry.io/otel/attribute.KeyValue]) and map item.
type KeyValue struct {
	Key   string
	Value Value
}

// Equal returns if a is equal to b.
func (a KeyValue) Equal(b KeyValue) bool { return false } // TODO (#4914): implement

// String returns an KeyValue for a string value.
func String(key, value string) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Int64 returns an KeyValue for an int64 value.
func Int64(key string, value int64) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Int returns an KeyValue for an int value.
func Int(key string, value int) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Float64 returns an KeyValue for a float64 value.
func Float64(key string, v float64) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Bool returns an KeyValue for a bool value.
func Bool(key string, v bool) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Bytes returns an KeyValue for a []byte value.
func Bytes(key string, v []byte) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Slice returns an KeyValue for a []Value value.
func Slice(key string, args ...Value) KeyValue { return KeyValue{} } // TODO (#4914): implement

// Map returns an KeyValue for a map value.
func Map(key string, args ...KeyValue) KeyValue { return KeyValue{} } // TODO (#4914): implement
