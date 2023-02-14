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

package attribute

import (
	"reflect"
	"testing"
)

func TestBoolSliceValue(t *testing.T) {
	type args struct {
		v []bool
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "TestBoolSliceValue", args: args{v: []bool{true, false, true}}, want: [3]bool{true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolSliceValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolSliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64SliceValue(t *testing.T) {
	type args struct {
		v []int64
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{name: "TestInt64SliceValue", args: args{[]int64{1, 2}}, want: [2]int64{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64SliceValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64SliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64SliceValue(t *testing.T) {
	type args struct {
		v []float64
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{name: "TestFloat64SliceValue", args: args{v: []float64{1, 2.3}}, want: [2]float64{1, 2.3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64SliceValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64SliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceValue(t *testing.T) {
	type args struct {
		v []string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{name: "TestStringSliceValue", args: args{[]string{"123456", "2"}}, want: [2]string{"123456", "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceValue(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsBoolSlice(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{name: "TestAsBoolSlice", args: args{[2]bool{true, false}}, want: []bool{true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsBoolSlice(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsSliceBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsInt64Slice(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{name: "TestAsInt64Slice", args: args{[2]int64{1, 3}}, want: []int64{1, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsInt64Slice(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsInt64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsFloat64Slice(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{name: "TestAsFloat64Slice", args: args{[2]float64{1.2, 3.1}}, want: []float64{1.2, 3.1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsFloat64Slice(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsFloat64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsStringSlice(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "TestAsStringSlice", args: args{[2]string{"1234", "12"}}, want: []string{"1234", "12"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsStringSlice(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
