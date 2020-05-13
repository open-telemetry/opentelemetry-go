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

package value_test

import (
	"testing"
	"unsafe"

	"go.opentelemetry.io/otel/api/kv"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/kv/value"
)

func TestValue(t *testing.T) {
	k := kv.Key("test")
	bli := getBitlessInfo(42)
	for _, testcase := range []struct {
		name      string
		value     value.Value
		wantType  value.Type
		wantValue interface{}
	}{
		{
			name:      "Key.Bool() correctly returns keys's internal bool value",
			value:     k.Bool(true).Value,
			wantType:  value.BOOL,
			wantValue: true,
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  value.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Uint64() correctly returns keys's internal uint64 value",
			value:     k.Uint64(42).Value,
			wantType:  value.UINT64,
			wantValue: uint64(42),
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  value.FLOAT64,
			wantValue: 42.1,
		},
		{
			name:      "Key.Int32() correctly returns keys's internal int32 value",
			value:     k.Int32(42).Value,
			wantType:  value.INT32,
			wantValue: int32(42),
		},
		{
			name:      "Key.Uint32() correctly returns keys's internal uint32 value",
			value:     k.Uint32(42).Value,
			wantType:  value.UINT32,
			wantValue: uint32(42),
		},
		{
			name:      "Key.Float32() correctly returns keys's internal float32 value",
			value:     k.Float32(42.1).Value,
			wantType:  value.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  value.STRING,
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
	signedType    value.Type
	unsignedType  value.Type
	signedValue   interface{}
	unsignedValue interface{}
}

func getBitlessInfo(i int) bitlessInfo {
	if unsafe.Sizeof(i) == 4 {
		return bitlessInfo{
			intValue:      i,
			uintValue:     uint(i),
			signedType:    value.INT32,
			unsignedType:  value.UINT32,
			signedValue:   int32(i),
			unsignedValue: uint32(i),
		}
	}
	return bitlessInfo{
		intValue:      i,
		uintValue:     uint(i),
		signedType:    value.INT64,
		unsignedType:  value.UINT64,
		signedValue:   int64(i),
		unsignedValue: uint64(i),
	}
}
