// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
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
			name: "Key.Map() correctly returns keys's internal map[string]Value value",
			value: k.Map(map[string]attribute.Value{
				"key1": attribute.StringValue("value1"),
				"key2": attribute.Int64Value(42),
			}).Value,
			wantType: attribute.MAP,
			wantValue: map[string]attribute.Value{
				"key1": attribute.StringValue("value1"),
				"key2": attribute.Int64Value(42),
			},
		},
		{
			name:      "MapValue with nil map",
			value:     attribute.MapValue(nil),
			wantType:  attribute.MAP,
			wantValue: map[string]attribute.Value{},
		},
		{
			name:      "MapValue with empty map",
			value:     attribute.MapValue(map[string]attribute.Value{}),
			wantType:  attribute.MAP,
			wantValue: map[string]attribute.Value{},
		},
		{
			name:     "MapValue with single item",
			value:    attribute.MapValue(map[string]attribute.Value{"key": attribute.StringValue("val")}),
			wantType: attribute.MAP,
			wantValue: map[string]attribute.Value{
				"key": attribute.StringValue("val"),
			},
		},
		{
			name: "MapValue with multiple items",
			value: attribute.MapValue(map[string]attribute.Value{
				"string": attribute.StringValue("test"),
				"int":    attribute.Int64Value(123),
				"float":  attribute.Float64Value(3.14),
				"bool":   attribute.BoolValue(true),
			}),
			wantType: attribute.MAP,
			wantValue: map[string]attribute.Value{
				"string": attribute.StringValue("test"),
				"int":    attribute.Int64Value(123),
				"float":  attribute.Float64Value(3.14),
				"bool":   attribute.BoolValue(true),
			},
		},
		{
			name: "MapValue with nested maps",
			value: attribute.MapValue(map[string]attribute.Value{
				"outer": attribute.MapValue(map[string]attribute.Value{
					"inner": attribute.StringValue("nested value"),
				}),
				"top": attribute.StringValue("top level"),
			}),
			wantType: attribute.MAP,
			wantValue: map[string]attribute.Value{
				"outer": attribute.MapValue(map[string]attribute.Value{
					"inner": attribute.StringValue("nested value"),
				}),
				"top": attribute.StringValue("top level"),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.value.Type() != testcase.wantType {
				t.Errorf(
					"wrong value type, got %#v, expected %#v",
					testcase.value.Type(),
					testcase.wantType,
				)
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
			attribute.Map("Map", map[string]attribute.Value{
				"key1": attribute.StringValue("value1"),
				"key2": attribute.Int64Value(42),
			}),
			attribute.Map("Map", map[string]attribute.Value{
				"key1": attribute.StringValue("value1"),
				"key2": attribute.Int64Value(42),
			}),
		},
		{
			attribute.Map("NestedMap", map[string]attribute.Value{
				"outer": attribute.MapValue(map[string]attribute.Value{
					"inner": attribute.StringValue("nested"),
				}),
				"top": attribute.BoolValue(true),
			}),
			attribute.Map("NestedMap", map[string]attribute.Value{
				"outer": attribute.MapValue(map[string]attribute.Value{
					"inner": attribute.StringValue("nested"),
				}),
				"top": attribute.BoolValue(true),
			}),
		},
	}

	t.Run("Distinct", func(t *testing.T) {
		for _, p := range pairs {
			s0, s1 := attribute.NewSet(p[0]), attribute.NewSet(p[1])
			m := map[attribute.Distinct]struct{}{s0.Equivalent(): {}}
			_, ok := m[s1.Equivalent()]
			assert.Truef(t, ok, "Distinct comparison of %s type: not equivalent", p[0].Value.Type())
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

	m1 := map[string]attribute.Value{
		"key1": attribute.StringValue("value1"),
		"key2": attribute.Int64Value(42),
		"key3": attribute.BoolValue(true),
	}
	kv = attribute.Map("Map", m1)
	m2 := kv.Value.AsMap()
	assert.Equal(t, m1, m2)
}
