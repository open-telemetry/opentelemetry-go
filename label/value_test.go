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

package label_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/label"
)

func TestValue(t *testing.T) {
	k := label.Key("test")
	bli := getBitlessInfo(42)
	for _, testcase := range []struct {
		name      string
		value     label.Value
		wantType  label.Type
		wantValue interface{}
	}{
		{
			name:      "Key.Bool() correctly returns keys's internal bool value",
			value:     k.Bool(true).Value,
			wantType:  label.BOOL,
			wantValue: true,
		},
		{
			name:      "Key.Array([]bool) correctly return key's internal bool values",
			value:     k.Array([]bool{true, false}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]bool{true, false},
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  label.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Uint64() correctly returns keys's internal uint64 value",
			value:     k.Uint64(42).Value,
			wantType:  label.UINT64,
			wantValue: uint64(42),
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  label.FLOAT64,
			wantValue: 42.1,
		},
		{
			name:      "Key.Int32() correctly returns keys's internal int32 value",
			value:     k.Int32(42).Value,
			wantType:  label.INT32,
			wantValue: int32(42),
		},
		{
			name:      "Key.Uint32() correctly returns keys's internal uint32 value",
			value:     k.Uint32(42).Value,
			wantType:  label.UINT32,
			wantValue: uint32(42),
		},
		{
			name:      "Key.Float32() correctly returns keys's internal float32 value",
			value:     k.Float32(42.1).Value,
			wantType:  label.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  label.STRING,
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
			wantType:  label.ARRAY,
			wantValue: [2]int64{42, 43},
		},
		{
			name:      "KeyArray([]uint64) correctly returns keys's internal uint64 values",
			value:     k.Array([]uint64{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]uint64{42, 43},
		},
		{
			name:      "Key.Array([]float64) correctly returns keys's internal float64 values",
			value:     k.Array([]float64{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]float64{42, 43},
		},
		{
			name:      "Key.Array([]int32) correctly returns keys's internal int32 values",
			value:     k.Array([]int32{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]int32{42, 43},
		},
		{
			name:      "Key.Array([]uint32) correctly returns keys's internal uint32 values",
			value:     k.Array([]uint32{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]uint32{42, 43},
		},
		{
			name:      "Key.Array([]float32) correctly returns keys's internal float32 values",
			value:     k.Array([]float32{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]float32{42, 43},
		},
		{
			name:      "Key.Array([]string) correctly return key's internal string values",
			value:     k.Array([]string{"foo", "bar"}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]string{"foo", "bar"},
		},
		{
			name:      "Key.Array([]int) correctly returns keys's internal signed integral values",
			value:     k.Array([]int{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]int{42, 43},
		},
		{
			name:      "Key.Array([]uint) correctly returns keys's internal unsigned integral values",
			value:     k.Array([]uint{42, 43}).Value,
			wantType:  label.ARRAY,
			wantValue: [2]uint{42, 43},
		},
		{
			name:      "Key.Array([][]int) refuses to create multi-dimensional array",
			value:     k.Array([][]int{{1, 2}, {3, 4}}).Value,
			wantType:  label.INVALID,
			wantValue: nil,
		},
	} {
		t.Logf("Running test case %s", testcase.name)
		if testcase.value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
		}
		if testcase.wantType == label.INVALID {
			continue
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
	signedType    label.Type
	unsignedType  label.Type
	signedValue   interface{}
	unsignedValue interface{}
}

func getBitlessInfo(i int) bitlessInfo {
	if unsafe.Sizeof(i) == 4 {
		return bitlessInfo{
			intValue:      i,
			uintValue:     uint(i),
			signedType:    label.INT32,
			unsignedType:  label.UINT32,
			signedValue:   int32(i),
			unsignedValue: uint32(i),
		}
	}
	return bitlessInfo{
		intValue:      i,
		uintValue:     uint(i),
		signedType:    label.INT64,
		unsignedType:  label.UINT64,
		signedValue:   int64(i),
		unsignedValue: uint64(i),
	}
}

func TestAsArrayValue(t *testing.T) {
	v := label.ArrayValue([]uint{1, 2, 3}).AsArray()
	// Ensure the returned dynamic type is stable.
	if got, want := reflect.TypeOf(v).Kind(), reflect.Array; got != want {
		t.Errorf("AsArray() returned %T, want %T", got, want)
	}
}
