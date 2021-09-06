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
type Type int

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
	// ARRAY is an array Type Value used to store 1-dimensional slices or
	// arrays of bool, int, int32, int64, uint, uint32, uint64, float,
	// float32, float64, or string types.
	//
	// Deprecated: Use slice types instead.
	ARRAY
)

// BoolValue creates a BOOL Value.
func BoolValue(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: internal.BoolToRaw(v),
	}
}

// BoolSliceValue creates a BOOLSLICE Value.
func BoolSliceValue(v []bool) Value {
	return sliceValue(v, BOOLSLICE)
}

// IntValue creates an INT64 Value.
func IntValue(v int) Value {
	return Int64Value(int64(v))
}

// IntSliceValue creates an INTSLICE Value.
func IntSliceValue(v []int) Value {
	return sliceValue(v, INT64SLICE)
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
	return sliceValue(v, INT64SLICE)
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
	return sliceValue(v, FLOAT64SLICE)
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
	return sliceValue(v, STRINGSLICE)
}

func sliceValue(v interface{}, vtype Type) Value {
	// get array type regardless of dimensions
	typ := reflect.TypeOf(v).Elem()
	kind := typ.Kind()
	switch kind {
	case reflect.Bool, reflect.Int, reflect.Int64,
		reflect.Float64, reflect.String:
		val := reflect.ValueOf(v)
		length := val.Len()
		frozen := reflect.Indirect(reflect.New(reflect.ArrayOf(length, typ)))
		reflect.Copy(frozen, val)
		return Value{
			vtype: vtype,
			slice: frozen.Interface(),
		}
	default:
		return Value{vtype: INVALID}
	}
}

// ArrayValue creates an ARRAY value from an array or slice.
// Only arrays or slices of bool, int, int64, float, float64, or string types are allowed.
// Specifically, arrays  and slices can not contain other arrays, slices, structs, or non-standard
// types. If the passed value is not an array or slice of these types an
// INVALID value is returned.
//
// Deprecated: Use the typed *SliceValue functions instead.
func ArrayValue(v interface{}) Value {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice:
		return sliceValue(v, ARRAY)
	}
	return Value{vtype: INVALID}
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
	if v.vtype != BOOLSLICE {
		return nil
	}
	r := []bool{}
	s := reflect.ValueOf(v.slice)
	for i := 0; i < s.Len(); i++ {
		r = append(r, s.Index(i).Bool())
	}
	return r
}

// AsInt64 returns the int64 value. Make sure that the Value's type is
// INT64.
func (v Value) AsInt64() int64 {
	return internal.RawToInt64(v.numeric)
}

// AsInt64Slice returns the []int64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsInt64Slice() []int64 {
	if v.vtype != INT64SLICE {
		return nil
	}
	r := []int64{}
	s := reflect.ValueOf(v.slice)
	for i := 0; i < s.Len(); i++ {
		r = append(r, s.Index(i).Int())
	}
	return r
}

// AsFloat64 returns the float64 value. Make sure that the Value's
// type is FLOAT64.
func (v Value) AsFloat64() float64 {
	return internal.RawToFloat64(v.numeric)
}

// AsFloat64Slice returns the []float64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsFloat64Slice() []float64 {
	if v.vtype != FLOAT64SLICE {
		return nil
	}
	r := []float64{}
	s := reflect.ValueOf(v.slice)
	for i := 0; i < s.Len(); i++ {
		r = append(r, s.Index(i).Float())
	}
	return r
}

// AsString returns the string value. Make sure that the Value's type
// is STRING.
func (v Value) AsString() string {
	return v.stringly
}

// AsStringSlice returns the []string value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsStringSlice() []string {
	if v.vtype != STRINGSLICE {
		return nil
	}
	r := []string{}
	s := reflect.ValueOf(v.slice)
	for i := 0; i < s.Len(); i++ {
		r = append(r, s.Index(i).String())
	}
	return r
}

// AsArray returns the array Value as an interface{}.
//
// Deprecated: Use the typed As*Slice functions instead.
func (v Value) AsArray() interface{} {
	return v.slice
}

type unknownValueType struct{}

// AsInterface returns Value's data as interface{}.
func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case ARRAY:
		return v.AsArray()
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
	case ARRAY, BOOLSLICE, INT64SLICE, FLOAT64SLICE, STRINGSLICE:
		return fmt.Sprint(v.slice)
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())
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
