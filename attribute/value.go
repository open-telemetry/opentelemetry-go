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
	"math"
	"slices"
	"strconv"
	"unsafe"
)

//go:generate stringer -type=Type

// Type describes the type of the data Value holds.
type Type int // nolint: revive  // redefines builtin Type.

// Value represents the value part in key-value pairs.
type Value struct {
	// num holds the value for of Type INT64, FLOAT64, and BOOL. It holds the
	// length for STRING, BOOLSLICE, INT64SLICE, FLOAT64SLICE, and STRINGSLICE.
	num uint64
	// any holds either the BOOL, INT64, FLOAT64, stringPtr, boolSlicePtr,
	// int64SlicePtr, float64SlicePtr, or stringSlicePtr. If Type is BOOL,
	// INT64, or FLOAT64 then the value of Value is in num as described above.
	// Otherwise, it contains the value wrapped in the appropriate type.
	any any
}

type (
	// stringPtr represents a value in Value.any for Type STRING Values.
	stringPtr       *byte
	boolSlicePtr    *bool
	int64SlicePtr   *int64
	float64SlicePtr *float64
	stringSlicePtr  *string
)

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

// BoolValue creates a BOOL Value.
func BoolValue(v bool) Value { //nolint:revive // Not a control flag.
	var n uint64
	if v {
		n = 1
	}
	return Value{num: n, any: BOOL}
}

// BoolSliceValue creates a BOOLSLICE Value.
func BoolSliceValue(v []bool) Value {
	cp := slices.Clone(v)
	return Value{
		num: uint64(len(cp)),
		any: boolSlicePtr(unsafe.SliceData(cp)),
	}
}

// IntValue creates an INT64 Value.
func IntValue(v int) Value {
	return Int64Value(int64(v))
}

// IntSliceValue creates an INTSLICE Value.
func IntSliceValue(v []int) Value {
	cp := make([]int64, len(v))
	for i := range cp {
		cp[i] = int64(v[i])
	}
	return Value{
		num: uint64(len(cp)),
		any: int64SlicePtr(unsafe.SliceData(cp)),
	}
}

// Int64Value creates an INT64 Value.
func Int64Value(v int64) Value {
	return Value{num: uint64(v), any: INT64}
}

// Int64SliceValue creates an INT64SLICE Value.
func Int64SliceValue(v []int64) Value {
	cp := slices.Clone(v)
	return Value{
		num: uint64(len(cp)),
		any: int64SlicePtr(unsafe.SliceData(cp)),
	}
}

// Float64Value creates a FLOAT64 Value.
func Float64Value(v float64) Value {
	return Value{num: math.Float64bits(v), any: FLOAT64}
}

// Float64SliceValue creates a FLOAT64SLICE Value.
func Float64SliceValue(v []float64) Value {
	cp := slices.Clone(v)
	return Value{
		num: uint64(len(cp)),
		any: float64SlicePtr(unsafe.SliceData(cp)),
	}
}

// StringValue creates a STRING Value.
func StringValue(v string) Value {
	return Value{
		num: uint64(len(v)),
		any: stringPtr(unsafe.StringData(v)),
	}
}

// StringSliceValue creates a STRINGSLICE Value.
func StringSliceValue(v []string) Value {
	cp := slices.Clone(v)
	return Value{
		num: uint64(len(cp)),
		any: stringSlicePtr(unsafe.SliceData(cp)),
	}
}

// Equal returns if v is equal to w.
func (v Value) Equal(w Value) bool {
	vType := v.Type()
	wType := w.Type()
	if vType != wType {
		return false
	}
	switch vType {
	case INT64, BOOL:
		return v.num == w.num
	case STRING:
		return v.asString() == w.asString()
	case FLOAT64:
		return v.asFloat64() == w.asFloat64()
	case BOOLSLICE:
		return slices.Equal(v.asBoolSlice(), w.asBoolSlice())
	case INT64SLICE:
		return slices.Equal(v.asInt64Slice(), w.asInt64Slice())
	case FLOAT64SLICE:
		return slices.Equal(v.asFloat64Slice(), w.asFloat64Slice())
	case STRINGSLICE:
		return slices.Equal(v.asStringSlice(), w.asStringSlice())
	case INVALID:
		return true
	default:
		return false
	}
}

// Type returns a type of the Value.
func (v Value) Type() Type {
	switch t := v.any.(type) {
	case Type:
		return t
	case stringPtr:
		return STRING
	case boolSlicePtr:
		return BOOLSLICE
	case int64SlicePtr:
		return INT64SLICE
	case float64SlicePtr:
		return FLOAT64SLICE
	case stringSlicePtr:
		return STRINGSLICE
	}
	return INVALID
}

// AsBool returns the bool value. Make sure that the Value's type is
// BOOL.
func (v Value) AsBool() bool {
	if v.Type() != BOOL {
		return false
	}
	return v.asBool()
}

// asBool returns the value held by v as a bool. If v is not of KindBool, this
// will return garbage.
func (v Value) asBool() bool { return v.num == 1 }

// AsBoolSlice returns the []bool value. Make sure that the Value's type is
// BOOLSLICE.
func (v Value) AsBoolSlice() []bool {
	if sp, ok := v.any.(boolSlicePtr); ok {
		return slices.Clone(unsafe.Slice((*bool)(sp), v.num))
	}
	return nil
}

func (v Value) asBoolSlice() []bool {
	return unsafe.Slice((*bool)(v.any.(boolSlicePtr)), v.num)
}

// AsInt64 returns the int64 value. Make sure that the Value's type is
// INT64.
func (v Value) AsInt64() int64 {
	if v.Type() != INT64 {
		return 0
	}
	return v.asInt64()
}

// asInt64 returns the value held by v as an int64. If v is not of KindInt64,
// this will return garbage.
func (v Value) asInt64() int64 { return int64(v.num) }

// AsInt64Slice returns the []int64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsInt64Slice() []int64 {
	if sp, ok := v.any.(int64SlicePtr); ok {
		return slices.Clone(unsafe.Slice((*int64)(sp), v.num))
	}
	return nil
}

func (v Value) asInt64Slice() []int64 {
	return unsafe.Slice((*int64)(v.any.(int64SlicePtr)), v.num)
}

// AsFloat64 returns the float64 value. Make sure that the Value's
// type is FLOAT64.
func (v Value) AsFloat64() float64 {
	if v.Type() != FLOAT64 {
		return 0
	}
	return v.asFloat64()
}

// asFloat64 returns the value held by v as a float64. If v is not of
// KindFloat64, this will return garbage.
func (v Value) asFloat64() float64 { return math.Float64frombits(v.num) }

// AsFloat64Slice returns the []float64 value. Make sure that the Value's type is
// FLOAT64SLICE.
func (v Value) AsFloat64Slice() []float64 {
	if sp, ok := v.any.(float64SlicePtr); ok {
		return slices.Clone(unsafe.Slice((*float64)(sp), v.num))
	}
	return nil
}

func (v Value) asFloat64Slice() []float64 {
	return unsafe.Slice((*float64)(v.any.(float64SlicePtr)), v.num)
}

// AsString returns the string value. Make sure that the Value's type
// is STRING.
func (v Value) AsString() string {
	if sp, ok := v.any.(stringPtr); ok {
		return unsafe.String(sp, v.num)
	}
	return ""
}

// asString returns the value held by v as a string. It will panic if the Value
// is not KindString.
func (v Value) asString() string {
	return unsafe.String(v.any.(stringPtr), v.num)
}

// AsStringSlice returns the []string value. Make sure that the Value's type is
// STRINGSLICE.
func (v Value) AsStringSlice() []string {
	if sp, ok := v.any.(stringSlicePtr); ok {
		return slices.Clone(unsafe.Slice((*string)(sp), v.num))
	}
	return nil
}

func (v Value) asStringSlice() []string {
	return unsafe.Slice((*string)(v.any.(stringSlicePtr)), v.num)
}

type unknownValueType struct{}

// AsInterface returns Value's data as interface{}.
func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case BOOL:
		return v.asBool()
	case BOOLSLICE:
		return slices.Clone(v.asBoolSlice())
	case INT64:
		return v.asInt64()
	case INT64SLICE:
		return slices.Clone(v.asInt64Slice())
	case FLOAT64:
		return v.asFloat64()
	case FLOAT64SLICE:
		return slices.Clone(v.asFloat64Slice())
	case STRING:
		return v.asString()
	case STRINGSLICE:
		return slices.Clone(v.asStringSlice())
	}
	return unknownValueType{}
}

// Emit returns a string representation of Value's data.
func (v Value) Emit() string {
	switch v.Type() {
	case BOOLSLICE:
		return fmt.Sprint(v.asBoolSlice())
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case INT64SLICE:
		return fmt.Sprint(v.asInt64Slice())
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case FLOAT64SLICE:
		return fmt.Sprint(v.asFloat64Slice())
	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())
	case STRINGSLICE:
		return fmt.Sprint(v.asStringSlice())
	case STRING:
		return v.asString()
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
