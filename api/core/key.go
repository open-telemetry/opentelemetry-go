// Copyright 2019, OpenTelemetry Authors
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

package core

//go:generate stringer -type=ValueType

import (
	"fmt"
	"unsafe"
)

type Key string

type KeyValue struct {
	Key   Key
	Value Value
}

type ValueType int

type Value struct {
	vtype    ValueType
	numeric  uint64
	stringly string
	// TODO Lazy value type?
}

const (
	INVALID ValueType = iota
	BOOL
	INT32
	INT64
	UINT32
	UINT64
	FLOAT32
	FLOAT64
	STRING
)

func Bool(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: boolToRaw(v),
	}
}

func Int64(v int64) Value {
	return Value{
		vtype:   INT64,
		numeric: int64ToRaw(v),
	}
}

func Uint64(v uint64) Value {
	return Value{
		vtype:   UINT64,
		numeric: uint64ToRaw(v),
	}
}

func Float64(v float64) Value {
	return Value{
		vtype:   FLOAT64,
		numeric: float64ToRaw(v),
	}
}

func Int32(v int32) Value {
	return Value{
		vtype:   INT32,
		numeric: int32ToRaw(v),
	}
}

func Uint32(v uint32) Value {
	return Value{
		vtype:   UINT32,
		numeric: uint32ToRaw(v),
	}
}

func Float32(v float32) Value {
	return Value{
		vtype:   FLOAT32,
		numeric: float32ToRaw(v),
	}
}

func String(v string) Value {
	return Value{
		vtype:    STRING,
		stringly: v,
	}
}

func Int(v int) Value {
	if unsafe.Sizeof(v) == 4 {
		return Int32(int32(v))
	}
	return Int64(int64(v))
}

func Uint(v uint) Value {
	if unsafe.Sizeof(v) == 4 {
		return Uint32(uint32(v))
	}
	return Uint64(uint64(v))
}

func (k Key) Bool(v bool) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Bool(v),
	}
}

func (k Key) Int64(v int64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int64(v),
	}
}

func (k Key) Uint64(v uint64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Uint64(v),
	}
}

func (k Key) Float64(v float64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float64(v),
	}
}

func (k Key) Int32(v int32) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int32(v),
	}
}

func (k Key) Uint32(v uint32) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Uint32(v),
	}
}

func (k Key) Float32(v float32) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float32(v),
	}
}

func (k Key) String(v string) KeyValue {
	return KeyValue{
		Key:   k,
		Value: String(v),
	}
}

func (k Key) Int(v int) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int(v),
	}
}

func (k Key) Uint(v uint) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Uint(v),
	}
}

func (k Key) Defined() bool {
	return len(k) != 0
}

func (v *Value) Type() ValueType {
	return v.vtype
}

func (v *Value) AsBool() bool {
	return rawToBool(v.numeric)
}

func (v *Value) AsInt32() int32 {
	return rawToInt32(v.numeric)
}

func (v *Value) AsInt64() int64 {
	return rawToInt64(v.numeric)
}

func (v *Value) AsUint32() uint32 {
	return rawToUint32(v.numeric)
}

func (v *Value) AsUint64() uint64 {
	return rawToUint64(v.numeric)
}

func (v *Value) AsFloat32() float32 {
	return rawToFloat32(v.numeric)
}

func (v *Value) AsFloat64() float64 {
	return rawToFloat64(v.numeric)
}

func (v *Value) AsString() string {
	return v.stringly
}

type unknownValueType struct{}

func (v *Value) AsInterface() interface{} {
	switch v.Type() {
	case BOOL:
		return v.AsBool()
	case INT32:
		return v.AsInt32()
	case INT64:
		return v.AsInt64()
	case UINT32:
		return v.AsUint32()
	case UINT64:
		return v.AsUint64()
	case FLOAT32:
		return v.AsFloat32()
	case FLOAT64:
		return v.AsFloat64()
	case STRING:
		return v.stringly
	}
	return unknownValueType{}
}

func (v *Value) Emit() string {
	if v.Type() == STRING {
		return v.stringly
	}
	i := v.AsInterface()
	if _, ok := i.(unknownValueType); ok {
		return "unknown"
	}
	return fmt.Sprint(i)
}
