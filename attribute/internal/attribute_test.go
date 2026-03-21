// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute

import (
	"reflect"
	"testing"
)

var wrapFloat64SliceValue = func(v any) any {
	if vi, ok := v.([]float64); ok {
		return SliceValue(vi)
	}
	return nil
}

var wrapInt64SliceValue = func(v any) any {
	if vi, ok := v.([]int64); ok {
		return SliceValue(vi)
	}
	return nil
}

var wrapBoolSliceValue = func(v any) any {
	if vi, ok := v.([]bool); ok {
		return SliceValue(vi)
	}
	return nil
}

var wrapStringSliceValue = func(v any) any {
	if vi, ok := v.([]string); ok {
		return SliceValue(vi)
	}
	return nil
}

var wrapBytesValue = func(v any) any {
	if vi, ok := v.([]byte); ok {
		return BytesValue(vi)
	}
	return nil
}

var (
	wrapAsBoolSlice    = func(v any) any { return AsSlice[bool](v) }
	wrapAsInt64Slice   = func(v any) any { return AsSlice[int64](v) }
	wrapAsFloat64Slice = func(v any) any { return AsSlice[float64](v) }
	wrapAsStringSlice  = func(v any) any { return AsSlice[string](v) }
)

func TestSliceValue(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want any
		fn   func(any) any
	}{
		{
			name: "Float64SliceValue() two items",
			args: args{v: []float64{1, 2.3}}, want: [2]float64{1, 2.3}, fn: wrapFloat64SliceValue,
		},
		{
			name: "Int64SliceValue() two items",
			args: args{[]int64{1, 2}}, want: [2]int64{1, 2}, fn: wrapInt64SliceValue,
		},
		{
			name: "BoolSliceValue() two items",
			args: args{v: []bool{true, false}}, want: [2]bool{true, false}, fn: wrapBoolSliceValue,
		},
		{
			name: "StringSliceValue() two items",
			args: args{[]string{"123", "2"}}, want: [2]string{"123", "2"}, fn: wrapStringSliceValue,
		},
		{
			name: "BytesValue() two items",
			args: args{v: []byte{1, 2}}, want: [2]byte{1, 2}, fn: wrapBytesValue,
		},
		{
			name: "AsBoolSlice() two items",
			args: args{[2]bool{true, false}}, want: []bool{true, false}, fn: wrapAsBoolSlice,
		},
		{
			name: "AsInt64Slice() two items",
			args: args{[2]int64{1, 3}}, want: []int64{1, 3}, fn: wrapAsInt64Slice,
		},
		{
			name: "AsFloat64Slice() two items",
			args: args{[2]float64{1.2, 3.1}}, want: []float64{1.2, 3.1}, fn: wrapAsFloat64Slice,
		},
		{
			name: "AsStringSlice() two items",
			args: args{[2]string{"1234", "12"}}, want: []string{"1234", "12"}, fn: wrapAsStringSlice,
		},
		{
			name: "AsBytes() two items",
			args: args{[2]byte{1, 2}}, want: []byte{1, 2}, fn: wrapAsBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsSliceMismatchedType(t *testing.T) {
	tests := []struct {
		name string
		fn   func() any
	}{
		{name: "bool from int64 array", fn: func() any { return AsSlice[bool]([2]int64{1, 2}) }},
		{name: "int64 from float64 array", fn: func() any { return AsSlice[int64]([2]float64{1, 2}) }},
		{name: "float64 from string array", fn: func() any { return AsSlice[float64]([2]string{"1", "2"}) }},
		{name: "string from bool array", fn: func() any { return AsSlice[string]([2]bool{true, false}) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn()
			rv := reflect.ValueOf(got)
			if !rv.IsNil() {
				t.Fatalf("got %v, want nil", got)
			}
		})
	}
}

// sync is a global used to ensure the benchmark are not optimized away.
var sync any

func BenchmarkBoolSliceValue(b *testing.B) {
	for _, bench := range []struct {
		name string
		s    []bool
	}{
		{name: "Len2", s: []bool{true, false}},
		{name: "Len8", s: []bool{true, false, true, false, true, false, true, false}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sync = SliceValue(bench.s)
			}
		})
	}
}

func BenchmarkInt64SliceValue(b *testing.B) {
	for _, bench := range []struct {
		name string
		s    []int64
	}{
		{name: "Len2", s: []int64{1, 2}},
		{name: "Len8", s: []int64{1, 2, 3, 4, 5, 6, 7, 8}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sync = SliceValue(bench.s)
			}
		})
	}
}

func BenchmarkFloat64SliceValue(b *testing.B) {
	for _, bench := range []struct {
		name string
		s    []float64
	}{
		{name: "Len2", s: []float64{1.2, 3.4}},
		{name: "Len8", s: []float64{1.2, 3.4, 5.6, 7.8, 9.1, 2.3, 4.5, 6.7}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sync = SliceValue(bench.s)
			}
		})
	}
}

func BenchmarkStringSliceValue(b *testing.B) {
	for _, bench := range []struct {
		name string
		s    []string
	}{
		{name: "Len2", s: []string{"a", "b"}},
		{name: "Len8", s: []string{"a", "b", "c", "d", "e", "f", "g", "h"}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sync = SliceValue(bench.s)
			}
		})
	}
}

func BenchmarkAsFloat64Slice(b *testing.B) {
	for _, bench := range []struct {
		name string
		in   any
	}{
		{name: "Len2", in: [2]float64{1, 2.3}},
		{name: "Len8", in: [8]float64{1, 2.3, 3.4, 4.5, 5.6, 6.7, 7.8, 8.9}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sync = AsSlice[float64](bench.in)
			}
		})
	}
}

func BenchmarkBytesValue(b *testing.B) {
	b.ReportAllocs()
	bs := []byte("foo")

	for b.Loop() {
		BytesValue(bs)
	}
}

func BenchmarkAsBytes(b *testing.B) {
	b.ReportAllocs()
	bs := [2]byte{1, 2}

	for b.Loop() {
		AsBytes(bs)
	}
}
