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

package kv_test

import (
	"testing"
	"unsafe"

	"go.opentelemetry.io/otel/api/kv"

	"github.com/google/go-cmp/cmp"
)

func TestValue(t *testing.T) {
	k := kv.Key("test")
	bli := getBitlessInfo(42)
	for _, testcase := range []struct {
		name      string
		value     kv.Value
		wantType  kv.Type
		wantValue interface{}
	}{
		{
			name:      "Key.Bool() correctly returns keys's internal bool value",
			value:     k.Bool(true).Value,
			wantType:  kv.BOOL,
			wantValue: true,
		},
		{
			name:      "Key.Array([]bool) correctly return key's internal bool values",
			value:     k.Array([]bool{true, false}).Value,
			wantType:  kv.ARRAY,
			wantValue: []bool{true, false},
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  kv.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Uint64() correctly returns keys's internal uint64 value",
			value:     k.Uint64(42).Value,
			wantType:  kv.UINT64,
			wantValue: uint64(42),
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  kv.FLOAT64,
			wantValue: 42.1,
		},
		{
			name:      "Key.Int32() correctly returns keys's internal int32 value",
			value:     k.Int32(42).Value,
			wantType:  kv.INT32,
			wantValue: int32(42),
		},
		{
			name:      "Key.Uint32() correctly returns keys's internal uint32 value",
			value:     k.Uint32(42).Value,
			wantType:  kv.UINT32,
			wantValue: uint32(42),
		},
		{
			name:      "Key.Float32() correctly returns keys's internal float32 value",
			value:     k.Float32(42.1).Value,
			wantType:  kv.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  kv.STRING,
			wantValue: "foo",
		},
		{
			name:      "Key.Int() correctly returns keys's internal signed integral value",
			value:     k.Int(bli.intValue).Value,
			wantType:  bli.signedType,
			wantValue: bli.signedValue,
		},
		{
			name:      "Key.Uint() correctly returns keys's internal unsigned integral value",
			value:     k.Uint(bli.uintValue).Value,
			wantType:  bli.unsignedType,
			wantValue: bli.unsignedValue,
		},
		{
			name:      "Key.Array([]int64) correctly returns keys's internal int64 values",
			value:     k.Array([]int64{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []int64{42, 43},
		},
		{
			name:      "KeyArray([]uint64) correctly returns keys's internal uint64 values",
			value:     k.Array([]uint64{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []uint64{42, 43},
		},
		{
			name:      "Key.Array([]float64) correctly returns keys's internal float64 values",
			value:     k.Array([]float64{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []float64{42, 43},
		},
		{
			name:      "Key.Array([]int32) correctly returns keys's internal int32 values",
			value:     k.Array([]int32{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []int32{42, 43},
		},
		{
			name:      "Key.Array([]uint32) correctly returns keys's internal uint32 values",
			value:     k.Array([]uint32{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []uint32{42, 43},
		},
		{
			name:      "Key.Array([]float32) correctly returns keys's internal float32 values",
			value:     k.Array([]float32{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []float32{42, 43},
		},
		{
			name:      "Key.Array([]string) correctly return key's internal string values",
			value:     k.Array([]string{"foo", "bar"}).Value,
			wantType:  kv.ARRAY,
			wantValue: []string{"foo", "bar"},
		},
		{
			name:      "Key.Array([]int) correctly returns keys's internal signed integral values",
			value:     k.Array([]int{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []int{42, 43},
		},
		{
			name:      "Key.Array([]uint) correctly returns keys's internal unsigned integral values",
			value:     k.Array([]uint{42, 43}).Value,
			wantType:  kv.ARRAY,
			wantValue: []uint{42, 43},
		},
		{
			name:      "Key.Array([][]int) correctly return key's multi dimensional array",
			value:     k.Array([][]int{{1, 2}, {3, 4}}).Value,
			wantType:  kv.ARRAY,
			wantValue: [][]int{{1, 2}, {3, 4}},
		},
	} {
		t.Logf("Running test case %s", testcase.name)
		if testcase.value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
		}
		got := testcase.value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}

type bitlessInfo struct {
	intValue      int
	uintValue     uint
	signedType    kv.Type
	unsignedType  kv.Type
	signedValue   interface{}
	unsignedValue interface{}
}

func getBitlessInfo(i int) bitlessInfo {
	if unsafe.Sizeof(i) == 4 {
		return bitlessInfo{
			intValue:      i,
			uintValue:     uint(i),
			signedType:    kv.INT32,
			unsignedType:  kv.UINT32,
			signedValue:   int32(i),
			unsignedValue: uint32(i),
		}
	}
	return bitlessInfo{
		intValue:      i,
		uintValue:     uint(i),
		signedType:    kv.INT64,
		unsignedType:  kv.UINT64,
		signedValue:   int64(i),
		unsignedValue: uint64(i),
	}
}
