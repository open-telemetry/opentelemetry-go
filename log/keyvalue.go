// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate stringer -type=Kind -trimprefix=Kind

package log // import "go.opentelemetry.io/otel/log"

import (
	"bytes"
	"errors"
	"math"
	"slices"
	"unsafe"

	"go.opentelemetry.io/otel/internal/global"
)

// errKind is logged when a Value is decoded to an incompatible type.
var errKind = errors.New("invalid Kind")

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
type Value struct {
	// Ensure forward compatibility by explicitly making this not comparable.
	noCmp [0]func() //nolint: unused  // This is indeed used.

	// num holds the value for Int64, Float64, and Bool. It holds the length
	// for String, Bytes, Slice, Map.
	num uint64
	// any holds either the KindBool, KindInt64, KindFloat64, stringptr,
	// bytesptr, sliceptr, or mapptr. If KindBool, KindInt64, or KindFloat64
	// then the value of Value is in num as described above. Otherwise, it
	// contains the value wrapped in the appropriate type.
	any any
}

type (
	// sliceptr represents a value in Value.any for KindString Values.
	stringptr *byte
	// bytesptr represents a value in Value.any for KindBytes Values.
	bytesptr *byte
	// sliceptr represents a value in Value.any for KindSlice Values.
	sliceptr *Value
	// mapptr represents a value in Value.any for KindMap Values.
	mapptr *KeyValue
)

// StringValue returns a new [Value] for a string.
func StringValue(v string) Value {
	return Value{
		num: uint64(len(v)),
		any: stringptr(unsafe.StringData(v)),
	}
}

// IntValue returns a [Value] for an int.
func IntValue(v int) Value { return Int64Value(int64(v)) }

// Int64Value returns a [Value] for an int64.
func Int64Value(v int64) Value {
	return Value{num: uint64(v), any: KindInt64}
}

// Float64Value returns a [Value] for a float64.
func Float64Value(v float64) Value {
	return Value{num: math.Float64bits(v), any: KindFloat64}
}

// BoolValue returns a [Value] for a bool.
func BoolValue(v bool) Value { //nolint:revive // Not a control flag.
	var n uint64
	if v {
		n = 1
	}
	return Value{num: n, any: KindBool}
}

// BytesValue returns a [Value] for a byte slice. The passed slice must not be
// changed after it is passed.
func BytesValue(v []byte) Value {
	return Value{
		num: uint64(len(v)),
		any: bytesptr(unsafe.SliceData(v)),
	}
}

// SliceValue returns a [Value] for a slice of [Value]. The passed slice must
// not be changed after it is passed.
func SliceValue(vs ...Value) Value {
	return Value{
		num: uint64(len(vs)),
		any: sliceptr(unsafe.SliceData(vs)),
	}
}

// MapValue returns a new [Value] for a slice of key-value pairs. The passed
// slice must not be changed after it is passed.
func MapValue(kvs ...KeyValue) Value {
	return Value{
		num: uint64(len(kvs)),
		any: mapptr(unsafe.SliceData(kvs)),
	}
}

// AsString returns the value held by v as a string.
func (v Value) AsString() string {
	if sp, ok := v.any.(stringptr); ok {
		return unsafe.String(sp, v.num)
	}
	global.Error(errKind, "AsString", "Kind", v.Kind())
	return ""
}

// asString returns the value held by v as a string. It will panic if the Value
// is not KindString.
func (v Value) asString() string {
	return unsafe.String(v.any.(stringptr), v.num)
}

// AsInt64 returns the value held by v as an int64.
func (v Value) AsInt64() int64 {
	if v.Kind() != KindInt64 {
		global.Error(errKind, "AsInt64", "Kind", v.Kind())
		return 0
	}
	return v.asInt64()
}

// asInt64 returns the value held by v as an int64. If v is not of KindInt64,
// this will return garbage.
func (v Value) asInt64() int64 { return int64(v.num) }

// AsBool returns the value held by v as a bool.
func (v Value) AsBool() bool {
	if v.Kind() != KindBool {
		global.Error(errKind, "AsBool", "Kind", v.Kind())
		return false
	}
	return v.asBool()
}

// asBool returns the value held by v as a bool. If v is not of KindBool, this
// will return garbage.
func (v Value) asBool() bool { return v.num == 1 }

// AsFloat64 returns the value held by v as a float64.
func (v Value) AsFloat64() float64 {
	if v.Kind() != KindFloat64 {
		global.Error(errKind, "AsFloat64", "Kind", v.Kind())
		return 0
	}
	return v.asFloat64()
}

// asFloat64 returns the value held by v as a float64. If v is not of
// KindFloat64, this will return garbage.
func (v Value) asFloat64() float64 { return math.Float64frombits(v.num) }

// AsBytes returns the value held by v as a []byte.
func (v Value) AsBytes() []byte {
	if sp, ok := v.any.(bytesptr); ok {
		return unsafe.Slice((*byte)(sp), v.num)
	}
	global.Error(errKind, "AsBytes", "Kind", v.Kind())
	return nil
}

// asBytes returns the value held by v as a []byte. It will panic if the Value
// is not KindBytes.
func (v Value) asBytes() []byte {
	return unsafe.Slice((*byte)(v.any.(bytesptr)), v.num)
}

// AsSlice returns the value held by v as a []Value.
func (v Value) AsSlice() []Value {
	if sp, ok := v.any.(sliceptr); ok {
		return unsafe.Slice((*Value)(sp), v.num)
	}
	global.Error(errKind, "AsSlice", "Kind", v.Kind())
	return nil
}

// asSlice returns the value held by v as a []Value. It will panic if the Value
// is not KindSlice.
func (v Value) asSlice() []Value {
	return unsafe.Slice((*Value)(v.any.(sliceptr)), v.num)
}

// AsMap returns the value held by v as a []KeyValue.
func (v Value) AsMap() []KeyValue {
	if sp, ok := v.any.(mapptr); ok {
		return unsafe.Slice((*KeyValue)(sp), v.num)
	}
	global.Error(errKind, "AsMap", "Kind", v.Kind())
	return nil
}

// asMap returns the value held by v as a []KeyValue. It will panic if the
// Value is not KindMap.
func (v Value) asMap() []KeyValue {
	return unsafe.Slice((*KeyValue)(v.any.(mapptr)), v.num)
}

// Kind returns the Kind of v.
func (v Value) Kind() Kind {
	switch x := v.any.(type) {
	case Kind:
		return x
	case stringptr:
		return KindString
	case bytesptr:
		return KindBytes
	case sliceptr:
		return KindSlice
	case mapptr:
		return KindMap
	default:
		return KindEmpty
	}
}

// Empty returns if v does not hold any value.
func (v Value) Empty() bool { return v.Kind() == KindEmpty }

// Equal returns if v is equal to w.
func (v Value) Equal(w Value) bool {
	k1 := v.Kind()
	k2 := w.Kind()
	if k1 != k2 {
		return false
	}
	switch k1 {
	case KindInt64, KindBool:
		return v.num == w.num
	case KindString:
		return v.asString() == w.asString()
	case KindFloat64:
		return v.asFloat64() == w.asFloat64()
	case KindSlice:
		return slices.EqualFunc(v.asSlice(), w.asSlice(), Value.Equal)
	case KindMap:
		return slices.EqualFunc(v.asMap(), w.asMap(), KeyValue.Equal)
	case KindBytes:
		return bytes.Equal(v.asBytes(), w.asBytes())
	case KindEmpty:
		return true
	default:
		global.Error(errKind, "Equal", "Kind", k1)
		return false
	}
}

// An KeyValue is a key-value pair used to represent a log attribute (a
// superset of [go.opentelemetry.io/otel/attribute.KeyValue]) and map item.
type KeyValue struct {
	Key   string
	Value Value
}

// Equal returns if a is equal to b.
func (a KeyValue) Equal(b KeyValue) bool {
	return a.Key == b.Key && a.Value.Equal(b.Value)
}

// String returns an KeyValue for a string value.
func String(key, value string) KeyValue {
	return KeyValue{key, StringValue(value)}
}

// Int64 returns an KeyValue for an int64 value.
func Int64(key string, value int64) KeyValue {
	return KeyValue{key, Int64Value(value)}
}

// Int returns an KeyValue for an int value.
func Int(key string, value int) KeyValue {
	return KeyValue{key, IntValue(value)}
}

// Float64 returns an KeyValue for a float64 value.
func Float64(key string, value float64) KeyValue {
	return KeyValue{key, Float64Value(value)}
}

// Bool returns an KeyValue for a bool value.
func Bool(key string, value bool) KeyValue {
	return KeyValue{key, BoolValue(value)}
}

// Bytes returns an KeyValue for a []byte value.
func Bytes(key string, value []byte) KeyValue {
	return KeyValue{key, BytesValue(value)}
}

// Slice returns an KeyValue for a []Value value.
func Slice(key string, value ...Value) KeyValue {
	return KeyValue{key, SliceValue(value...)}
}

// Map returns an KeyValue for a map value.
func Map(key string, value ...KeyValue) KeyValue {
	return KeyValue{key, MapValue(value...)}
}
