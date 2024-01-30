// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestKind(t *testing.T) {
	testCases := []struct {
		kind  Kind
		str   string
		value int
	}{
		{KindBool, "Bool", 1},
		{KindBytes, "Bytes", 5},
		{KindEmpty, "Empty", 0},
		{KindFloat64, "Float64", 2},
		{KindInt64, "Int64", 3},
		{KindList, "List", 6},
		{KindMap, "Map", 7},
		{KindString, "String", 4},
	}
	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			assert.Equal(t, tc.value, int(tc.kind))
			assert.Equal(t, tc.str, tc.kind.String())
		})
	}
}

func TestValueEqual(t *testing.T) {
	vals := []Value{
		{},
		Int64Value(1),
		Int64Value(2),
		Float64Value(3.5),
		Float64Value(3.7),
		BoolValue(true),
		BoolValue(false),
		StringValue("hi"),
		BytesValue([]byte{1, 3, 5}),
		ListValue(IntValue(3), StringValue("foo")),
		MapValue(Bool("b", true), Int("i", 3)),
		MapValue(List("l", IntValue(3), StringValue("foo")), Bytes("b", []byte{3, 5, 7})),
	}
	for i, v1 := range vals {
		for j, v2 := range vals {
			got := v1.Equal(v2)
			want := i == j
			if got != want {
				t.Errorf("%v.Equal(%v): got %t, want %t", v1, v2, got, want)
			}
		}
	}
}

func TestValueString(t *testing.T) {
	for _, test := range []struct {
		v    Value
		want string
	}{
		{Int64Value(-3), "-3"},
		{Float64Value(.15), "0.15"},
		{BoolValue(true), "true"},
		{StringValue("foo"), "foo"},
		{BytesValue([]byte{2, 4, 6}), "[2 4 6]"},
		{ListValue(IntValue(3), StringValue("foo")), "[3 foo]"},
		{MapValue(Int("a", 1), Bool("b", true)), "[a=1 b=true]"},
		{Value{}, "<nil>"},
	} {
		got := test.v.String()
		assert.Equal(t, test.want, got)
	}
}

func TestValueNoAlloc(t *testing.T) {
	// Assign values just to make sure the compiler doesn't optimize away the statements.
	var (
		i  int64
		f  float64
		b  bool
		by []byte
		s  string
	)
	bytes := []byte{1, 3, 4}
	a := int(testing.AllocsPerRun(5, func() {
		i = Int64Value(1).AsInt64()
		f = Float64Value(1).AsFloat64()
		b = BoolValue(true).AsBool()
		by = BytesValue(bytes).AsBytes()
		s = StringValue("foo").AsString()
	}))
	assert.Zero(t, a)
	_ = i
	_ = f
	_ = b
	_ = by
	_ = s
}

func TestKeyValueNoAlloc(t *testing.T) {
	// Assign values just to make sure the compiler doesn't optimize away the statements.
	var (
		i  int64
		f  float64
		b  bool
		by []byte
		s  string
	)
	bytes := []byte{1, 3, 4}
	a := int(testing.AllocsPerRun(5, func() {
		i = Int64("key", 1).Value.AsInt64()
		f = Float64("key", 1).Value.AsFloat64()
		b = Bool("key", true).Value.AsBool()
		by = Bytes("key", bytes).Value.AsBytes()
		s = String("key", "foo").Value.AsString()
	}))
	assert.Zero(t, a)
	_ = i
	_ = f
	_ = b
	_ = by
	_ = s
}

func TestValueAny(t *testing.T) {
	for _, test := range []struct {
		want any
		in   Value
	}{
		{"s", StringValue("s")},
		{true, BoolValue(true)},
		{int64(4), IntValue(4)},
		{int64(11), Int64Value(11)},
		{1.5, Float64Value(1.5)},
		{[]byte{1, 2, 3}, BytesValue([]byte{1, 2, 3})},
		{[]Value{IntValue(3)}, ListValue(IntValue(3))},
		{[]KeyValue{Int("i", 3)}, MapValue(Int("i", 3))},
		{nil, Value{}},
	} {
		got := test.in.AsAny()
		assert.Equal(t, test.want, got)
	}
}

func TestEmptyMap(t *testing.T) {
	g := Map("g")
	got := g.Value.AsMap()
	assert.Nil(t, got)
}

func TestEmptyList(t *testing.T) {
	l := ListValue()
	got := l.AsList()
	assert.Nil(t, got)
}

func TestMapValueWithEmptyMaps(t *testing.T) {
	// Preserve empty groups.
	g := MapValue(
		Int("a", 1),
		Map("g1", Map("g2")),
		Map("g3", Map("g4", Int("b", 2))))
	got := g.AsMap()
	want := []KeyValue{Int("a", 1), Map("g1", Map("g2")), Map("g3", Map("g4", Int("b", 2)))}
	assert.Equal(t, want, got)
}

func TestListValueWithEmptyValues(t *testing.T) {
	// Preserve empty values.
	l := ListValue(Value{})
	got := l.AsList()
	want := []Value{{}}
	assert.Equal(t, want, got)
}

// A Value with "unsafe" strings is significantly faster:
// safe:  1785 ns/op, 0 allocs
// unsafe: 690 ns/op, 0 allocs

// Run this with and without -tags unsafe_kvs to compare.
func BenchmarkUnsafeStrings(b *testing.B) {
	b.ReportAllocs()
	dst := make([]Value, 100)
	src := make([]Value, len(dst))
	b.Logf("Value size = %d", unsafe.Sizeof(Value{}))
	for i := range src {
		src[i] = StringValue(fmt.Sprintf("string#%d", i))
	}
	b.ResetTimer()
	var d string
	for i := 0; i < b.N; i++ {
		copy(dst, src)
		for _, a := range dst {
			d = a.AsString()
		}
	}
	_ = d
}
