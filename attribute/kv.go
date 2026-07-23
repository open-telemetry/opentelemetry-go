// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute

import (
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// KeyValue holds a key and value pair.
type KeyValue struct {
	Key   Key
	Value Value
}

// String returns a string representation of kv using the
// [OpenTelemetry Attribute representation for non-OTLP] rules.
//
// The returned string is a JSON object containing a single key-value pair.
//
// The returned string is meant for debugging;
// the string representation is not stable.
//
// [OpenTelemetry Attribute representation for non-OTLP]: https://opentelemetry.io/docs/specs/otel/common/#attribute-representation-for-non-otlp
func (kv KeyValue) String() string {
	const jsonObjectSyntaxLen = len(`{"":}`)

	var b strings.Builder
	n := len(kv.Key) + jsonObjectSyntaxLen
	switch kv.Value.Type() {
	case BOOL:
		if kv.Value.AsBool() {
			n += len("true")
		} else {
			n += len("false")
		}
	case BOOLSLICE:
		n += jsonArrayBracketsLen
		if l := reflect.ValueOf(kv.Value.slice).Len(); l > 0 {
			n += l*boolArrayElemMaxLen + (l-1)*commaLen
		}
	case INT64:
		var buf [int64ArrayElemMaxLen]byte
		n += len(strconv.AppendInt(buf[:0], kv.Value.AsInt64(), 10))
	case INT64SLICE:
		n += jsonArrayBracketsLen
		if l := reflect.ValueOf(kv.Value.slice).Len(); l > 0 {
			n += l*int64ArrayElemMaxLen + (l-1)*commaLen
		}
	case FLOAT64:
		val := kv.Value.AsFloat64()
		switch {
		case math.IsNaN(val):
			n += len(`"NaN"`)
		case math.IsInf(val, 1):
			n += len(`"Infinity"`)
		case math.IsInf(val, -1):
			n += len(`"-Infinity"`)
		default:
			var buf [float64ArrayElemMaxLen]byte
			n += len(strconv.AppendFloat(buf[:0], val, 'g', -1, 64))
		}
	case FLOAT64SLICE:
		n += jsonArrayBracketsLen
		if l := reflect.ValueOf(kv.Value.slice).Len(); l > 0 {
			n += l*float64ArrayElemMaxLen + (l-1)*commaLen
		}
	case STRING:
		n += len(kv.Value.stringly) + quotesLen
	case STRINGSLICE:
		n += jsonArrayBracketsLen
		if l := reflect.ValueOf(kv.Value.slice).Len(); l > 0 {
			n += l*smallObjectLen + (l-1)*commaLen
		}
	case BYTESLICE:
		n += base64.StdEncoding.EncodedLen(len(kv.Value.stringly)) + quotesLen
	case SLICE:
		n += jsonArrayBracketsLen
		if l := reflect.ValueOf(kv.Value.slice).Len(); l > 0 {
			n += l*smallObjectLen + (l-1)*commaLen
		}
	case EMPTY:
		n += len("null")
	default:
		n += len(`"unknown"`)
	}
	b.Grow(n)

	_ = b.WriteByte('{')
	appendJSONString(&b, string(kv.Key))
	_ = b.WriteByte(':')
	appendJSONValue(&b, kv.Value)
	_ = b.WriteByte('}')
	return b.String()
}

// Valid reports whether kv is a valid OpenTelemetry attribute.
func (kv KeyValue) Valid() bool {
	return kv.Key.Defined()
}

// Bool returns a [KeyValue] for a bool value.
func Bool(k string, v bool) KeyValue {
	return Key(k).Bool(v)
}

// BoolSlice returns a [KeyValue] for a []bool value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func BoolSlice(k string, v []bool) KeyValue {
	return Key(k).BoolSlice(v)
}

// Int returns a [KeyValue] for an int value.
//
// It is provided as a convenience for [Int64].
func Int(k string, v int) KeyValue {
	return Key(k).Int(v)
}

// IntSlice returns a [KeyValue] for a []int value.
//
// It is provided as a convenience for [Int64Slice].
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func IntSlice(k string, v []int) KeyValue {
	return Key(k).IntSlice(v)
}

// Int64 returns a [KeyValue] for an int64 value.
func Int64(k string, v int64) KeyValue {
	return Key(k).Int64(v)
}

// Int64Slice returns a [KeyValue] for a []int64 value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func Int64Slice(k string, v []int64) KeyValue {
	return Key(k).Int64Slice(v)
}

// Float64 returns a [KeyValue] for a float64 value.
func Float64(k string, v float64) KeyValue {
	return Key(k).Float64(v)
}

// Float64Slice returns a [KeyValue] for a []float64 value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func Float64Slice(k string, v []float64) KeyValue {
	return Key(k).Float64Slice(v)
}

// String returns a [KeyValue] for a string value.
func String(k, v string) KeyValue {
	return Key(k).String(v)
}

// StringSlice returns a [KeyValue] for a []string value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func StringSlice(k string, v []string) KeyValue {
	return Key(k).StringSlice(v)
}

// ByteSlice returns a [KeyValue] for a []byte value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func ByteSlice(k string, v []byte) KeyValue {
	return Key(k).ByteSlice(v)
}

// Slice returns a [KeyValue] for a []Value value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
func Slice(k string, v ...Value) KeyValue {
	return Key(k).Slice(v...)
}

// Map returns a [KeyValue] for a []KeyValue value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// Users should avoid providing duplicate keys; many receivers handle maps
// containing duplicate keys unpredictably.
//
// The order of v is not preserved.
func Map(k string, v ...KeyValue) KeyValue {
	return Key(k).Map(v...)
}

// Stringer creates a new key-value pair with a passed name and a string
// value generated by the passed Stringer interface.
func Stringer(k string, v fmt.Stringer) KeyValue {
	return Key(k).String(v.String())
}
