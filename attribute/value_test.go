// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestValue(t *testing.T) {
	k := attribute.Key("test")
	for _, testcase := range []struct {
		name      string
		value     attribute.Value
		wantType  attribute.Type
		wantValue any
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
			name:      "Key.Int64() correctly returns negative keys's internal int64 value",
			value:     k.Int64(-42).Value,
			wantType:  attribute.INT64,
			wantValue: int64(-42),
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
		{
			name:      "Key.ByteSlice() correctly returns keys's internal []byte value",
			value:     k.ByteSlice([]byte("hello world")).Value,
			wantType:  attribute.BYTESLICE,
			wantValue: []byte("hello world"),
		},
		{
			name:      "Key.Slice() correctly returns keys's internal []Value value",
			value:     k.Slice([]attribute.Value{attribute.BoolValue(true), attribute.IntValue(42), attribute.StringValue("foo")}).Value,
			wantType:  attribute.SLICE,
			wantValue: []attribute.Value{attribute.BoolValue(true), attribute.IntValue(42), attribute.StringValue("foo")},
		},
		{
			name:      "empty value",
			value:     attribute.Value{},
			wantType:  attribute.EMPTY,
			wantValue: nil,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.value.Type() != testcase.wantType {
				t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
			}
			got := testcase.value.AsInterface()
			if diff := cmp.Diff(testcase.wantValue, got, cmp.AllowUnexported(attribute.Value{})); diff != "" {
				t.Errorf("+got, -want: %s", diff)
			}
		})
	}
}

func TestEquivalence(t *testing.T) {
	pairs := [][2]attribute.KeyValue{
		{
			attribute.Bool("Bool", true),
			attribute.Bool("Bool", true),
		},
		{
			attribute.BoolSlice("BoolSlice", []bool{true, false, true}),
			attribute.BoolSlice("BoolSlice", []bool{true, false, true}),
		},
		{
			attribute.Int("Int", 34),
			attribute.Int("Int", 34),
		},
		{
			attribute.IntSlice("IntSlice", []int{312, 1, -2}),
			attribute.IntSlice("IntSlice", []int{312, 1, -2}),
		},
		{
			attribute.Int64("Int64", 98),
			attribute.Int64("Int64", 98),
		},
		{
			attribute.Int64Slice("Int64Slice", []int64{12, 1298, -219, 2}),
			attribute.Int64Slice("Int64Slice", []int64{12, 1298, -219, 2}),
		},
		{
			attribute.Float64("Float64", 19.09),
			attribute.Float64("Float64", 19.09),
		},
		{
			attribute.Float64Slice("Float64Slice", []float64{12398.1, -37.1713873737, 3}),
			attribute.Float64Slice("Float64Slice", []float64{12398.1, -37.1713873737, 3}),
		},
		{
			attribute.String("String", "string value"),
			attribute.String("String", "string value"),
		},
		{
			attribute.StringSlice("StringSlice", []string{"one", "two", "three"}),
			attribute.StringSlice("StringSlice", []string{"one", "two", "three"}),
		},
		{
			attribute.ByteSlice("ByteSlice", []byte("one")),
			attribute.ByteSlice("ByteSlice", []byte("one")),
		},
		{
			attribute.Slice("Slice", []attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(42),
				attribute.SliceValue([]attribute.Value{attribute.StringValue("nested")}),
			}),
			attribute.Slice("Slice", []attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(42),
				attribute.SliceValue([]attribute.Value{attribute.StringValue("nested")}),
			}),
		},
		{
			attribute.KeyValue{Key: "Empty"},
			attribute.KeyValue{Key: "Empty"},
		},
	}

	t.Run("Distinct", func(t *testing.T) {
		for _, p := range pairs {
			s0, s1 := attribute.NewSet(p[0]), attribute.NewSet(p[1])
			m := map[attribute.Distinct]struct{}{s0.Equivalent(): {}}
			_, ok := m[s1.Equivalent()]
			assert.Truef(
				t,
				ok,
				"Distinct comparison of %s type: not equivalent: %s != %s",
				p[0].Value.Type(),
				s0.Encoded(attribute.DefaultEncoder()),
				s1.Encoded(attribute.DefaultEncoder()),
			)
		}
	})

	t.Run("Equality operator", func(t *testing.T) {
		// Maintain backwards compatibility.
		for _, p := range pairs {
			if p[0] != p[1] {
				t.Errorf("Expected %v to be equal to %v", p[0], p[1])
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		// Maintain backwards compatibility.
		for _, p := range pairs {
			s0, s1 := attribute.NewSet(p[0]), attribute.NewSet(p[1])
			m := map[attribute.Set]struct{}{s0: {}}
			_, ok := m[s1]
			assert.Truef(
				t,
				ok,
				"Set comparison of %s type: not equivalent: %s != %s",
				p[0].Value.Type(),
				s0.Encoded(attribute.DefaultEncoder()),
				s1.Encoded(attribute.DefaultEncoder()),
			)
		}
	})
}

func TestNotEquivalence(t *testing.T) {
	pairs := [][2]attribute.KeyValue{
		{
			attribute.Int("Key", 0),
			attribute.Bool("Key", false),
		},
		{
			attribute.Bool("Bool", true),
			attribute.Bool("Bool", false),
		},
		{
			attribute.BoolSlice("BoolSlice", []bool{true, false, true}),
			attribute.BoolSlice("BoolSlice", []bool{true, true, true}),
		},
		{
			attribute.Int("Int", 34),
			attribute.Int("Int", 32),
		},
		{
			attribute.IntSlice("IntSlice", []int{312, 1, -2}),
			attribute.IntSlice("IntSlice", []int{312, 2, -2}),
		},
		{
			attribute.Int64("Int64", 98),
			attribute.Int64("Int64", 97),
		},
		{
			attribute.Int64Slice("Int64Slice", []int64{12, 1298, -219, 2}),
			attribute.Int64Slice("Int64Slice", []int64{12, 1298, -219, 1}),
		},
		{
			attribute.Float64("Float64", 19.09),
			attribute.Float64("Float64", 22.09),
		},
		{
			attribute.ByteSlice("ByteSlice", []byte("bytes value")),
			attribute.ByteSlice("ByteSlice", []byte("another value")),
		},
		{
			attribute.Float64Slice("Float64Slice", []float64{12398.1, -37.1713873737, 3}),
			attribute.Float64Slice("Float64Slice", []float64{12398.1, -37.1713873737, 5}),
		},
		{
			attribute.String("String", "string value"),
			attribute.String("String", "another value"),
		},
		{
			attribute.StringSlice("StringSlice", []string{"one", "two", "three"}),
			attribute.StringSlice("StringSlice", []string{"one", "two"}),
		},
		{
			attribute.Slice("Slice", []attribute.Value{attribute.BoolValue(true), attribute.IntValue(42)}),
			attribute.Slice("Slice", []attribute.Value{attribute.BoolValue(true), attribute.IntValue(43)}),
		},
		{
			attribute.KeyValue{Key: "Empty"},
			attribute.String("Empty", ""),
		},
	}

	t.Run("Distinct", func(t *testing.T) {
		for _, p := range pairs {
			s0, s1 := attribute.NewSet(p[0]), attribute.NewSet(p[1])
			m := map[attribute.Distinct]struct{}{s0.Equivalent(): {}}
			_, ok := m[s1.Equivalent()]
			assert.Falsef(
				t,
				ok,
				"Distinct comparison of %s type: equivalent: %s == %s",
				p[0].Value.Type(),
				s0.Encoded(attribute.DefaultEncoder()),
				s1.Encoded(attribute.DefaultEncoder()),
			)
		}
	})

	t.Run("Equality operator", func(t *testing.T) {
		// Maintain backwards compatibility.
		for _, p := range pairs {
			if p[0] == p[1] {
				t.Errorf("Expected %v to not be equal to %v", p[0], p[1])
			}
		}
	})

	t.Run("Set", func(t *testing.T) {
		// Maintain backwards compatibility.
		for _, p := range pairs {
			s0, s1 := attribute.NewSet(p[0]), attribute.NewSet(p[1])
			m := map[attribute.Set]struct{}{s0: {}}
			_, ok := m[s1]
			assert.Falsef(
				t,
				ok,
				"Set comparison of %s type: equivalent: %s == %s",
				p[0].Value.Type(),
				s0.Encoded(attribute.DefaultEncoder()),
				s1.Encoded(attribute.DefaultEncoder()),
			)
		}
	})
}

func TestAsSlice(t *testing.T) {
	bs1 := []bool{true, false, true}
	kv := attribute.BoolSlice("BoolSlice", bs1)
	bs2 := kv.Value.AsBoolSlice()
	assert.Equal(t, bs1, bs2)

	i64s1 := []int64{12, 1298, -219, 2}
	kv = attribute.Int64Slice("Int64Slice", i64s1)
	i64s2 := kv.Value.AsInt64Slice()
	assert.Equal(t, i64s1, i64s2)

	is1 := []int{12, 1298, -219, 2}
	kv = attribute.IntSlice("IntSlice", is1)
	i64s2 = kv.Value.AsInt64Slice()
	assert.Equal(t, i64s1, i64s2)

	fs1 := []float64{12398.1, -37.1713873737, 3}
	kv = attribute.Float64Slice("Float64Slice", fs1)
	fs2 := kv.Value.AsFloat64Slice()
	assert.Equal(t, fs1, fs2)

	ss1 := []string{"one", "two", "three"}
	kv = attribute.StringSlice("StringSlice", ss1)
	ss2 := kv.Value.AsStringSlice()
	assert.Equal(t, ss1, ss2)

	b1 := []byte("one")
	kv = attribute.ByteSlice("ByteSlice", b1)
	b2 := kv.Value.AsByteSlice()
	assert.Equal(t, b1, b2)

	for _, tc := range []struct {
		name string
		in   []attribute.Value
	}{
		{
			name: "empty",
			in:   []attribute.Value{},
		},
		{
			name: "len1",
			in:   []attribute.Value{attribute.BoolValue(true)},
		},
		{
			name: "len2",
			in:   []attribute.Value{attribute.BoolValue(true), attribute.IntValue(42)},
		},
		{
			name: "len3",
			in:   []attribute.Value{attribute.BoolValue(true), attribute.IntValue(42), attribute.StringValue("test")},
		},
		{
			name: "len4",
			in: []attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(42),
				attribute.StringValue("test"),
				attribute.Float64Value(1.25),
			},
		},
		{
			name: "len5",
			in: []attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(42),
				attribute.StringValue("test"),
				attribute.Float64Value(1.25),
				attribute.ByteSliceValue([]byte("bin")),
			},
		},
		{
			name: "reflect path",
			in: []attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(42),
				attribute.StringValue("test"),
				attribute.Float64Value(1.25),
				attribute.ByteSliceValue([]byte("bin")),
				attribute.SliceValue([]attribute.Value{attribute.BoolValue(false)}),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			kv = attribute.Slice("Slice", tc.in)
			assert.Equal(t, tc.in, kv.Value.AsSlice())
		})
	}
}

func TestValueString(t *testing.T) {
	for _, tc := range []struct {
		name string
		v    attribute.Value
		want string
	}{
		{
			name: "bool",
			v:    attribute.BoolValue(true),
			want: "true",
		},
		{
			name: "bool false",
			v:    attribute.BoolValue(false),
			want: "false",
		},
		{
			name: "bool slice len1 fast path",
			v:    attribute.BoolSliceValue([]bool{false}),
			want: `[false]`,
		},
		{
			name: "bool slice len2 fast path",
			v:    attribute.BoolSliceValue([]bool{true, false}),
			want: `[true,false]`,
		},
		{
			name: "empty bool slice",
			v:    attribute.BoolSliceValue(nil),
			want: "[]",
		},
		{
			name: "empty bool slice literal",
			v:    attribute.BoolSliceValue([]bool{}),
			want: "[]",
		},
		{
			name: "bool slice",
			v:    attribute.BoolSliceValue([]bool{true, false, true}),
			want: `[true,false,true]`,
		},
		{
			name: "bool slice reflect path",
			v:    attribute.BoolSliceValue([]bool{false, true, false, true}),
			want: `[false,true,false,true]`,
		},
		{
			name: "int64",
			v:    attribute.Int64Value(-42),
			want: "-42",
		},
		{
			name: "int",
			v:    attribute.IntValue(7),
			want: "7",
		},
		{
			name: "int64 slice len1 fast path",
			v:    attribute.Int64SliceValue([]int64{-1}),
			want: `[-1]`,
		},
		{
			name: "int64 slice len2 fast path",
			v:    attribute.Int64SliceValue([]int64{1, -2}),
			want: `[1,-2]`,
		},
		{
			name: "empty int slice",
			v:    attribute.IntSliceValue(nil),
			want: "[]",
		},
		{
			name: "empty int slice literal",
			v:    attribute.IntSliceValue([]int{}),
			want: "[]",
		},
		{
			name: "empty int64 slice literal",
			v:    attribute.Int64SliceValue([]int64{}),
			want: "[]",
		},
		{
			name: "int slice",
			v:    attribute.IntSliceValue([]int{1, -2, 3}),
			want: `[1,-2,3]`,
		},
		{
			name: "int64 slice reflect path",
			v:    attribute.Int64SliceValue([]int64{1, -2, 3, -4}),
			want: `[1,-2,3,-4]`,
		},
		{
			name: "float64",
			v:    attribute.Float64Value(1.23e10),
			want: "1.23e+10",
		},
		{
			name: "float64 negative zero",
			v:    attribute.Float64Value(math.Copysign(0, -1)),
			want: "-0",
		},
		{
			name: "float64 NaN",
			v:    attribute.Float64Value(math.NaN()),
			want: "NaN",
		},
		{
			name: "float64 +Inf",
			v:    attribute.Float64Value(math.Inf(1)),
			want: "Infinity",
		},
		{
			name: "float64 -Inf",
			v:    attribute.Float64Value(math.Inf(-1)),
			want: "-Infinity",
		},
		{
			name: "empty float64 slice",
			v:    attribute.Float64SliceValue(nil),
			want: "[]",
		},
		{
			name: "empty float64 slice literal",
			v:    attribute.Float64SliceValue([]float64{}),
			want: "[]",
		},
		{
			name: "float64 slice len1 fast path",
			v:    attribute.Float64SliceValue([]float64{math.Inf(-1)}),
			want: `["-Infinity"]`,
		},
		{
			name: "float64 slice len3 fast path",
			v:    attribute.Float64SliceValue([]float64{1.25, math.Copysign(0, -1), 2.5}),
			want: `[1.25,-0,2.5]`,
		},
		{
			name: "float64 slice",
			v: attribute.Float64SliceValue([]float64{
				1,
				math.NaN(),
				math.Inf(1),
				math.Inf(-1),
				math.Copysign(0, -1),
			}),
			want: `[1,"NaN","Infinity","-Infinity",-0]`,
		},
		{
			name: "float64 slice fast path",
			v: attribute.Float64SliceValue([]float64{
				math.NaN(),
				math.Inf(1),
			}),
			want: `["NaN","Infinity"]`,
		},
		{
			name: "string",
			v:    attribute.StringValue(`hello "world"`),
			want: `hello "world"`,
		},
		{
			name: "empty string",
			v:    attribute.StringValue(""),
			want: "",
		},
		{
			name: "empty string slice",
			v:    attribute.StringSliceValue(nil),
			want: "[]",
		},
		{
			name: "empty string slice literal",
			v:    attribute.StringSliceValue([]string{}),
			want: "[]",
		},
		{
			name: "string slice len1 fast path",
			v:    attribute.StringSliceValue([]string{""}),
			want: `[""]`,
		},
		{
			name: "string slice len3 fast path",
			v:    attribute.StringSliceValue([]string{"snowman ☃", "left\u2028right", "left\u2029right"}),
			want: `["snowman ☃","left\u2028right","left\u2029right"]`,
		},
		{
			name: "string slice",
			v: attribute.StringSliceValue([]string{
				`hello "world"`,
				"line\nbreak",
				string([]byte{0xff, 'a'}),
				"\u2028",
			}),
			want: `["hello \"world\"","line\nbreak","\ufffda","\u2028"]`,
		},
		{
			name: "string slice fast path escapes",
			v: attribute.StringSliceValue([]string{
				"tab\treturn\rformfeed\fbackslash\\quote\"backspace\b",
				string([]byte{0x01}) + "\u2029",
			}),
			want: `["tab\treturn\rformfeed\fbackslash\\quote\"backspace\b","\u0001\u2029"]`,
		},
		{
			name: "string slice leaves HTML characters unescaped",
			v:    attribute.StringSliceValue([]string{"<tag>&"}),
			want: `["<tag>&"]`,
		},
		{
			name: "string slice replaces invalid utf8 after copied prefix",
			v:    attribute.StringSliceValue([]string{string([]byte{'a', 0xff, 'b'})}),
			want: `["a\ufffdb"]`,
		},
		{
			name: "byte slice",
			v:    attribute.ByteSliceValue([]byte("hello world")),
			want: "aGVsbG8gd29ybGQ=",
		},
		{
			name: "empty byte slice",
			v:    attribute.ByteSliceValue(nil),
			want: "",
		},
		{
			name: "empty slice",
			v:    attribute.SliceValue(nil),
			want: "[]",
		},
		{
			name: "slice len5 fast path",
			v: attribute.SliceValue([]attribute.Value{
				attribute.BoolValue(true),
				attribute.IntValue(7),
				attribute.Float64Value(math.Copysign(0, -1)),
				attribute.StringValue(`hello "world"`),
				attribute.ByteSliceValue([]byte("bin")),
			}),
			want: `[true,7,-0,"hello \"world\"","Ymlu"]`,
		},
		{
			name: "slice len1 fast path",
			v:    attribute.SliceValue([]attribute.Value{attribute.BoolValue(false)}),
			want: `[false]`,
		},
		{
			name: "slice len2 fast path",
			v: attribute.SliceValue([]attribute.Value{
				attribute.IntValue(7),
				attribute.StringValue(`hello "world"`),
			}),
			want: `[7,"hello \"world\""]`,
		},
		{
			name: "slice len3 fast path",
			v: attribute.SliceValue([]attribute.Value{
				attribute.Float64Value(1.25),
				attribute.Float64Value(math.Inf(1)),
				attribute.Float64Value(math.Inf(-1)),
			}),
			want: `[1.25,"Infinity","-Infinity"]`,
		},
		{
			name: "slice",
			v: attribute.SliceValue([]attribute.Value{
				attribute.StringValue("hello \"world\""),
				attribute.Float64Value(math.NaN()),
				attribute.ByteSliceValue([]byte("bin")),
				attribute.SliceValue([]attribute.Value{attribute.BoolValue(true), {}}),
			}),
			want: `["hello \"world\"","NaN","Ymlu",[true,null]]`,
		},
		{
			name: "slice reflect path nested slice values",
			v: attribute.SliceValue([]attribute.Value{
				attribute.BoolSliceValue([]bool{}),
				attribute.BoolSliceValue([]bool{true}),
				attribute.BoolSliceValue([]bool{true, false}),
				attribute.BoolSliceValue([]bool{true, false, true}),
				attribute.BoolSliceValue([]bool{false, true, false, true}),
				attribute.Int64SliceValue([]int64{}),
				attribute.Int64SliceValue([]int64{-1}),
				attribute.Int64SliceValue([]int64{1, -2}),
				attribute.Int64SliceValue([]int64{1, -2, 3}),
				attribute.Int64SliceValue([]int64{1, -2, 3, -4}),
				attribute.Float64SliceValue([]float64{}),
				attribute.Float64SliceValue([]float64{math.Inf(-1)}),
				attribute.Float64SliceValue([]float64{math.NaN(), math.Inf(1)}),
				attribute.Float64SliceValue([]float64{1.25, math.Copysign(0, -1), 2.5}),
				attribute.Float64SliceValue([]float64{1, math.NaN(), math.Inf(1), math.Inf(-1)}),
				attribute.StringSliceValue([]string{}),
				attribute.StringSliceValue([]string{""}),
				attribute.StringSliceValue([]string{`hello "world"`, "line\nbreak"}),
				attribute.StringSliceValue([]string{"snowman ☃", "left\u2028right", "left\u2029right"}),
				attribute.StringSliceValue([]string{
					"tab\treturn\rformfeed\fbackslash\\quote\"backspace\b",
					string([]byte{0x01}) + "\u2029",
					"<tag>&",
					string([]byte{'a', 0xff, 'b'}),
				}),
				attribute.SliceValue([]attribute.Value{}),
				attribute.SliceValue([]attribute.Value{attribute.BoolValue(true)}),
				attribute.SliceValue([]attribute.Value{attribute.BoolValue(true), attribute.IntValue(2)}),
				attribute.SliceValue([]attribute.Value{attribute.BoolValue(true), attribute.IntValue(2), attribute.StringValue("x")}),
				attribute.SliceValue([]attribute.Value{
					attribute.BoolValue(true),
					attribute.IntValue(2),
					attribute.StringValue("x"),
					attribute.Float64Value(math.Inf(1)),
				}),
				attribute.SliceValue([]attribute.Value{
					attribute.BoolValue(true),
					attribute.IntValue(2),
					attribute.StringValue("x"),
					attribute.Float64Value(math.Inf(1)),
					attribute.ByteSliceValue([]byte("bin")),
				}),
				attribute.SliceValue([]attribute.Value{
					attribute.BoolValue(true),
					attribute.IntValue(2),
					attribute.StringValue("x"),
					attribute.Float64Value(math.Inf(1)),
					attribute.ByteSliceValue([]byte("bin")),
					{},
				}),
			}),
			want: `[[],[true],[true,false],[true,false,true],[false,true,false,true],[]` +
				`,[-1],[1,-2],[1,-2,3],[1,-2,3,-4],[]` +
				`,["-Infinity"],["NaN","Infinity"],[1.25,-0,2.5],[1,"NaN","Infinity","-Infinity"],[]` +
				`,[""],["hello \"world\"","line\nbreak"],["snowman ☃","left\u2028right","left\u2029right"]` +
				`,["tab\treturn\rformfeed\fbackslash\\quote\"backspace\b","\u0001\u2029","<tag>&","a\ufffdb"]` +
				`,[],[true],[true,2],[true,2,"x"],[true,2,"x","Infinity"],[true,2,"x","Infinity","Ymlu"],[true,2,"x","Infinity","Ymlu",null]]`,
		},
		{
			name: "empty",
			v:    attribute.Value{},
			want: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.v.String())
		})
	}
}
