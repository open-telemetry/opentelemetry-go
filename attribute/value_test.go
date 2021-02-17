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

package attribute_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/attribute"
)

func TestValue(t *testing.T) {
	k := attribute.Key("test")
	bli := getBitlessInfo(42)
	for _, testcase := range []struct {
		name      string
		value     attribute.Value
		wantType  attribute.Type
		wantValue interface{}
	}{
		{
			name:      "Key.Bool() correctly returns keys's internal bool value",
			value:     k.Bool(true).Value,
			wantType:  attribute.BOOL,
			wantValue: true,
		},
		{
			name:      "Key.Array([]bool) correctly return key's internal bool values",
			value:     k.Array([]bool{true, false}).Value,
			wantType:  attribute.ARRAY,
			wantValue: [2]bool{true, false},
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  attribute.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  attribute.FLOAT64,
			wantValue: 42.1,
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  attribute.STRING,
			wantValue: "foo",
		},
		{
			name:      "Key.Int() correctly returns keys's internal signed integral value",
			value:     k.Int(bli.intValue).Value,
			wantType:  bli.signedType,
			wantValue: bli.signedValue,
		},
		{
			name:      "Key.Array([]int64) correctly returns keys's internal int64 values",
			value:     k.Array([]int64{42, 43}).Value,
			wantType:  attribute.ARRAY,
			wantValue: [2]int64{42, 43},
		},
		{
			name:      "Key.Array([]float64) correctly returns keys's internal float64 values",
			value:     k.Array([]float64{42, 43}).Value,
			wantType:  attribute.ARRAY,
			wantValue: [2]float64{42, 43},
		},
		{
			name:      "Key.Array([]string) correctly return key's internal string values",
			value:     k.Array([]string{"foo", "bar"}).Value,
			wantType:  attribute.ARRAY,
			wantValue: [2]string{"foo", "bar"},
		},
		{
			name:      "Key.Array([]int) correctly returns keys's internal signed integral values",
			value:     k.Array([]int{42, 43}).Value,
			wantType:  attribute.ARRAY,
			wantValue: [2]int{42, 43},
		},
		{
			name:      "Key.Array([][]int) refuses to create multi-dimensional array",
			value:     k.Array([][]int{{1, 2}, {3, 4}}).Value,
			wantType:  attribute.INVALID,
			wantValue: nil,
		},
	} {
		t.Logf("Running test case %s", testcase.name)
		if testcase.value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
		}
		if testcase.wantType == attribute.INVALID {
			continue
		}
		got := testcase.value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}

type bitlessInfo struct {
	intValue    int
	uintValue   uint
	signedType  attribute.Type
	signedValue interface{}
}

func getBitlessInfo(i int) bitlessInfo {
	return bitlessInfo{
		intValue:    i,
		uintValue:   uint(i),
		signedType:  attribute.INT64,
		signedValue: int64(i),
	}
}

func TestAsArrayValue(t *testing.T) {
	v := attribute.ArrayValue([]int{1, 2, 3}).AsArray()
	// Ensure the returned dynamic type is stable.
	if got, want := reflect.TypeOf(v).Kind(), reflect.Array; got != want {
		t.Errorf("AsArray() returned %T, want %T", got, want)
	}
}
