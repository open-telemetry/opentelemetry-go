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
	"strconv"

	"go.opentelemetry.io/otel/internal"
)

//go:generate stringer -type=Type

// Type describes the type of the data Value holds.
type Type int // nolint: revive  // redefines builtin Type.

// Value represents the value part in key-value pairs.
type Value struct {
	vtype    Type
	numeric  uint64
	stringly string
	slice    interface{}
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
)

func sliceValue[T bool | int64 | float64 | string](v []T) any {
	var zero T
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero)))
	copy(cp.Elem().Slice(0, len(v)).Interface().([]T), v)
	return cp.Elem().Interface()
}

func asSlice[T bool | int64 | float64 | string](v any) []T {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	var zero T
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(zero))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]T)
}

// BoolValue creates a BOOL Value.
func BoolValue(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: internal.BoolToRaw(v),
	}
}

// BoolSliceValue creates a BOOLSLICE Value.
func BoolSliceValue(v []bool) Value {
	return Value{vtype: BOOLSLICE, slice: sliceValue(v)}
}

// IntValue creates an INT64 Value.
func IntValue(v int) Value {
	return Int64Value(int64(v))
}

// IntSliceValue creates an INTSLICE Value.
func IntSliceValue(v []int) Value {
	var int64Val int64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(int64Val)))
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
		numeric: internal.Int64ToRaw(v),
	}
}

// Int64SliceValue creates an INT64SLICE Value.
func Int64SliceValue(v []int64) Value {
	return Value{vtype: INT64SLICE, slice: sliceValue(v)}
}

// Float64Value creates a FLOAT64 Value.
func Float64Value(v float64) Value {
	return Value{
		vtype:   FLOAT64,
		numeric: internal.Float64ToRaw(v),
	}
}

// Float64SliceValue creates a FLOAT64SLICE Value.
func Float64SliceValue(v []float64) Value {
	return Value{vtype: FLOAT64SLICE, slice: sliceValue(v)}
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
	return Value{vtype: STRINGSLICE, slice: sliceValue(v)}
}

// Type returns a type of the Value.
func (v Value) Type() Type {
	return v.vtype
}

// AsBool returns the bool value. Make sure that the Value's type is
// BOOL.
func (v Value) AsBool() bool {
	return internal.RawToBool(v.numeric)
}

// AsBoolSlice returns the []bool value. Make sure that the Value's type is
// BOOLSLICE.
func (v Value) AsBoolSlice() []bool {
	rv := reflect.ValueOf(v.slice)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	correctLen := rv.Len()
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(false))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]bool)
}

// AsInt64 returns the int64 value. Make sure that the Value's type is
// INT64.
func (v Value) AsInt64() int64 {
	return internal.RawToInt64(v.numeric)
}

// AsInt64Slice returns the []int64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsInt64Slice() []int64 {
	rv := reflect.ValueOf(v.slice)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	correctLen := rv.Len()
	var int64var int64
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(int64var))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]int64)
}

// AsFloat64 returns the float64 value. Make sure that the Value's
// type is FLOAT64.
func (v Value) AsFloat64() float64 {
	return internal.RawToFloat64(v.numeric)
}

// AsFloat64Slice returns the []float64 value. Make sure that the Value's type is
// FLOAT64SLICE.
func (v Value) AsFloat64Slice() []float64 {
	rv := reflect.ValueOf(v.slice)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	correctLen := rv.Len()
	var float64var float64
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(float64var))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]float64)
}

// AsString returns the string value. Make sure that the Value's type
// is STRING.
func (v Value) AsString() string {
	return v.stringly
}

// AsStringSlice returns the []string value. Make sure that the Value's type is
// STRINGSLICE.
func (v Value) AsStringSlice() []string {
	rv := reflect.ValueOf(v.slice)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	correctLen := rv.Len()
	var stringStr string
	correctType := reflect.ArrayOf(correctLen, reflect.TypeOf(stringStr))
	cpy := reflect.New(correctType)
	_ = reflect.Copy(cpy.Elem(), rv)
	return cpy.Elem().Slice(0, correctLen).Interface().([]string)
}

type unknownValueType struct{}

// AsInterface returns Value's data as interface{}.
func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case BOOL:
		return v.AsBool()
	case BOOLSLICE:
		return v.AsBoolSlice()
	case INT64:
		return v.AsInt64()
	case INT64SLICE:
		return v.AsInt64Slice()
	case FLOAT64:
		return v.AsFloat64()
	case FLOAT64SLICE:
		return v.AsFloat64Slice()
	case STRING:
		return v.stringly
	case STRINGSLICE:
		return v.AsStringSlice()
	}
	return unknownValueType{}
}

// Emit returns a string representation of Value's data.
func (v Value) Emit() string {
	switch v.Type() {
	case BOOLSLICE:
		return fmt.Sprint(*(v.slice.(*[]bool)))
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case INT64SLICE:
		return fmt.Sprint(*(v.slice.(*[]int64)))
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case FLOAT64SLICE:
		return fmt.Sprint(*(v.slice.(*[]float64)))
	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())
	case STRINGSLICE:
		return fmt.Sprint(*(v.slice.(*[]string)))
	case STRING:
		return v.stringly
	default:
		return "unknown"
	}
}

// MarshalJSON returns the JSON encoding of the Value.
func (v Value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value interface{}
	}
	jsonVal.Type = v.Type().String()
	jsonVal.Value = v.AsInterface()
	return json.Marshal(jsonVal)
}
