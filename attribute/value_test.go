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
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/attribute"
)

func TestValue(t *testing.T) {
	k := attribute.Key("test")
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
			name:      "Key.BoolSlice() correctly returns keys's internal []bool value",
			value:     k.BoolSlice([]bool{true, false, true}).Value,
			wantType:  attribute.BOOLSLICE,
			wantValue: []bool{true, false, true},
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  attribute.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Int64Slice() correctly returns keys's internal []int64 value",
			value:     k.Int64Slice([]int64{42, -3, 12}).Value,
			wantType:  attribute.INT64SLICE,
			wantValue: []int64{42, -3, 12},
		},
		{
			name:      "Key.Int() correctly returns keys's internal signed integral value",
			value:     k.Int(42).Value,
			wantType:  attribute.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.IntSlice() correctly returns keys's internal []int64 value",
			value:     k.IntSlice([]int{42, -3, 12}).Value,
			wantType:  attribute.INT64SLICE,
			wantValue: []int64{42, -3, 12},
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  attribute.FLOAT64,
			wantValue: 42.1,
		},
		{
			name:      "Key.Float64Slice() correctly returns keys's internal []float64 value",
			value:     k.Float64Slice([]float64{42, -3, 12}).Value,
			wantType:  attribute.FLOAT64SLICE,
			wantValue: []float64{42, -3, 12},
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  attribute.STRING,
			wantValue: "foo",
		},
		{
			name:      "Key.StringSlice() correctly returns keys's internal []string value",
			value:     k.StringSlice([]string{"forty-two", "negative three", "twelve"}).Value,
			wantType:  attribute.STRINGSLICE,
			wantValue: []string{"forty-two", "negative three", "twelve"},
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
