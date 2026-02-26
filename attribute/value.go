// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"

	attribute "go.opentelemetry.io/otel/attribute/internal"
)

//go:generate stringer -type=Type

// Type describes the type of the data Value holds.
type Type int // nolint: revive  // redefines builtin Type.

// Value represents the value part in key-value pairs.
type Value struct {
	vtype    Type
	numeric  uint64
	stringly string
	slice    any
}

const (
	// INVALID is used for a Value with no value set.
	INVALID Type = iota
	// BOOL is a boolean Type Value.
	BOOL
	// INT64 is a 64-bit signed integral Type Value.
	INT64
	// FLOAT64 is a 64-bit floating point Type Value.
	FLOAT64
	// STRING is a string Type Value.
	STRING
	// BOOLSLICE is a slice of booleans Type Value.
	BOOLSLICE
	// INT64SLICE is a slice of 64-bit signed integral numbers Type Value.
	INT64SLICE
	// FLOAT64SLICE is a slice of 64-bit floating point numbers Type Value.
	FLOAT64SLICE
	// STRINGSLICE is a slice of strings Type Value.
	STRINGSLICE
	// MAP is a map of string keys to Values Type Value.
	MAP
)

// BoolValue creates a BOOL Value.
func BoolValue(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: boolToRaw(v),
	}
}

// BoolSliceValue creates a BOOLSLICE Value.
func BoolSliceValue(v []bool) Value {
	return Value{vtype: BOOLSLICE, slice: attribute.BoolSliceValue(v)}
}

// IntValue creates an INT64 Value.
func IntValue(v int) Value {
	return Int64Value(int64(v))
}

// IntSliceValue creates an INTSLICE Value.
func IntSliceValue(v []int) Value {
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeFor[int64]()))
	for i, val := range v {
		cp.Elem().Index(i).SetInt(int64(val))
	}
	return Value{
		vtype: INT64SLICE,
		slice: cp.Elem().Interface(),
	}
}

// Int64Value creates an INT64 Value.
func Int64Value(v int64) Value {
	return Value{
		vtype:   INT64,
		numeric: int64ToRaw(v),
	}
}

// Int64SliceValue creates an INT64SLICE Value.
func Int64SliceValue(v []int64) Value {
	return Value{vtype: INT64SLICE, slice: attribute.Int64SliceValue(v)}
}

// Float64Value creates a FLOAT64 Value.
func Float64Value(v float64) Value {
	return Value{
		vtype:   FLOAT64,
		numeric: float64ToRaw(v),
	}
}

// Float64SliceValue creates a FLOAT64SLICE Value.
func Float64SliceValue(v []float64) Value {
	return Value{vtype: FLOAT64SLICE, slice: attribute.Float64SliceValue(v)}
}

// StringValue creates a STRING Value.
func StringValue(v string) Value {
	return Value{
		vtype:    STRING,
		stringly: v,
	}
}

// StringSliceValue creates a STRINGSLICE Value.
func StringSliceValue(v []string) Value {
	return Value{vtype: STRINGSLICE, slice: attribute.StringSliceValue(v)}
}

// MapValue creates a MAP Value.
//
// A nil map is treated the same as an empty map. Both will round-trip via
// AsMap as an empty (non-nil) map[string]Value. This is consistent with
// how nil slices are handled by the other slice-typed constructors.
func MapValue(v map[string]Value) Value {
	return Value{vtype: MAP, slice: mapValue(v)}
}

// mapValue converts a map into a sorted array of KeyValue stored as any.
func mapValue(v map[string]Value) any {
	n := len(v)
	if n == 0 {
		// Return a zero-length array to maintain type information.
		return reflect.New(reflect.ArrayOf(0, reflect.TypeFor[KeyValue]())).Elem().Interface()
	}

	keys := make([]string, 0, n)
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]KeyValue, n)
	for i, k := range keys {
		entries[i] = KeyValue{Key: Key(k), Value: v[k]}
	}

	cp := reflect.New(reflect.ArrayOf(n, reflect.TypeFor[KeyValue]())).Elem()
	reflect.Copy(cp, reflect.ValueOf(entries))
	return cp.Interface()
}

// Type returns a type of the Value.
func (v Value) Type() Type {
	return v.vtype
}

// AsBool returns the bool value. Make sure that the Value's type is
// BOOL.
func (v Value) AsBool() bool {
	return rawToBool(v.numeric)
}

// AsBoolSlice returns the []bool value. Make sure that the Value's type is
// BOOLSLICE.
func (v Value) AsBoolSlice() []bool {
	if v.vtype != BOOLSLICE {
		return nil
	}
	return v.asBoolSlice()
}

func (v Value) asBoolSlice() []bool {
	return attribute.AsBoolSlice(v.slice)
}

// AsInt64 returns the int64 value. Make sure that the Value's type is
// INT64.
func (v Value) AsInt64() int64 {
	return rawToInt64(v.numeric)
}

// AsInt64Slice returns the []int64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsInt64Slice() []int64 {
	if v.vtype != INT64SLICE {
		return nil
	}
	return v.asInt64Slice()
}

func (v Value) asInt64Slice() []int64 {
	return attribute.AsInt64Slice(v.slice)
}

// AsFloat64 returns the float64 value. Make sure that the Value's
// type is FLOAT64.
func (v Value) AsFloat64() float64 {
	return rawToFloat64(v.numeric)
}

// AsFloat64Slice returns the []float64 value. Make sure that the Value's type is
// FLOAT64SLICE.
func (v Value) AsFloat64Slice() []float64 {
	if v.vtype != FLOAT64SLICE {
		return nil
	}
	return v.asFloat64Slice()
}

func (v Value) asFloat64Slice() []float64 {
	return attribute.AsFloat64Slice(v.slice)
}

// AsString returns the string value. Make sure that the Value's type
// is STRING.
func (v Value) AsString() string {
	return v.stringly
}

// AsStringSlice returns the []string value. Make sure that the Value's type is
// STRINGSLICE.
func (v Value) AsStringSlice() []string {
	if v.vtype != STRINGSLICE {
		return nil
	}
	return v.asStringSlice()
}

func (v Value) asStringSlice() []string {
	return attribute.AsStringSlice(v.slice)
}

// AsMap returns the map[string]Value value. Make sure that the Value's type is
// MAP.
func (v Value) AsMap() map[string]Value {
	if v.vtype != MAP {
		return nil
	}
	return v.asMap()
}

func (v Value) asMap() map[string]Value {
	entries := v.asMapKeyValues()
	if entries == nil {
		return nil
	}
	result := make(map[string]Value, len(entries))
	for _, kv := range entries {
		result[string(kv.Key)] = kv.Value
	}
	return result
}

func (v Value) asMapKeyValues() []KeyValue {
	rv := reflect.ValueOf(v.slice)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	n := rv.Len()
	if n == 0 {
		return []KeyValue{}
	}
	cpy := make([]KeyValue, n)
	reflect.Copy(reflect.ValueOf(cpy), rv)
	return cpy
}

type unknownValueType struct{}

// AsInterface returns Value's data as any.
func (v Value) AsInterface() any {
	switch v.Type() {
	case BOOL:
		return v.AsBool()
	case BOOLSLICE:
		return v.asBoolSlice()
	case INT64:
		return v.AsInt64()
	case INT64SLICE:
		return v.asInt64Slice()
	case FLOAT64:
		return v.AsFloat64()
	case FLOAT64SLICE:
		return v.asFloat64Slice()
	case STRING:
		return v.stringly
	case STRINGSLICE:
		return v.asStringSlice()
	case MAP:
		return v.asMap()
	}
	return unknownValueType{}
}

// encodeSlice encodes a slice to a JSON string using json.Encoder.
// SetEscapeHTML(false) is used to avoid escaping characters like '<', '>', '&'.
// The trailing newline added by Encode is trimmed before returning.
func encodeSlice(v any) (string, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return "", err
	}
	out := b.Bytes()
	// json.Encoder.Encode always appends a newline, trim it.
	if len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}
	return string(out), nil
}

// encodeString encodes a string to a quoted JSON string.
// SetEscapeHTML(false) is used to avoid escaping HTML characters.
func encodeString(s string) []byte {
	// json.Marshal for a plain string never errors.
	b, _ := json.Marshal(s)
	return b
}

// Emit returns a string representation of Value's data.
func (v Value) Emit() string {
	switch v.Type() {
	case BOOLSLICE:
		return fmt.Sprint(v.asBoolSlice())

	case BOOL:
		return strconv.FormatBool(v.AsBool())

	case INT64SLICE:
		// json.Encoder errors here only if the slice contains unsupported types,
		// which cannot happen for []int64, so the error path is a safety fallback.
		s, err := encodeSlice(v.asInt64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asInt64Slice())
		}
		return s

	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)

	case FLOAT64SLICE:
		// json.Encoder returns an error if the slice contains NaN or Inf values
		// since these are not valid JSON numbers. The fallback uses fmt.Sprintf
		// which renders them as "NaN", "+Inf", "-Inf" etc.
		s, err := encodeSlice(v.asFloat64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asFloat64Slice())
		}
		return s

	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())

	case STRINGSLICE:
		s, err := encodeSlice(v.asStringSlice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asStringSlice())
		}
		return s

	case STRING:
		return v.stringly

	case MAP:
		entries := v.asMapKeyValues()
		if len(entries) == 0 {
			return "{}"
		}

		var b bytes.Buffer
		b.WriteByte('{') // bytes.Buffer.WriteByte never errors.
		for i, kv := range entries {
			if i > 0 {
				_ = b.WriteByte(',')
			}

			// Encode the key as a quoted JSON string.
			_,_=b.Write(encodeString(string(kv.Key)))
			_=b.WriteByte(':')

			// Encode the value based on its type.
			switch kv.Value.Type() {
			case INVALID:
				_,_=b.WriteString("null")

			case STRING:
				_,_=b.Write(encodeString(kv.Value.stringly))

			case FLOAT64:
				// NaN and Inf are not valid JSON — represent them as strings
				// to avoid silent encoding errors or panics.
				f := kv.Value.AsFloat64()
				switch {
				case math.IsNaN(f):
					_,_=b.WriteString(`"NaN"`)
				case math.IsInf(f, 1):
					_,_=b.WriteString(`"Infinity"`)
				case math.IsInf(f, -1):
					_,_=b.WriteString(`"-Infinity"`)
				default:
					_,_=b.WriteString(kv.Value.Emit())
				}

			default:
				_,_=b.WriteString(kv.Value.Emit())
			}
		}
		_=b.WriteByte('}')
		return b.String()

	default:
		return "unknown"
	}
}

// MarshalJSON returns the JSON encoding of the Value.
func (v Value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value any
	}
	jsonVal.Type = v.Type().String()
	jsonVal.Value = v.AsInterface()
	return json.Marshal(jsonVal)
}
