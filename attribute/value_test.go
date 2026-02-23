// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"reflect"
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
	} {
		t.Logf("Running test case %s", testcase.name)
		if testcase.value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
		}
		if testcase.wantType == attribute.INVALID {
			continue
		}
		got := testcase.value.AsInterface()

		// Use a cmp.Comparer for attribute.Value to handle unexported fields.
		opt := cmp.Comparer(func(a, b attribute.Value) bool {
			return a.Type() == b.Type() && reflect.DeepEqual(a.AsInterface(), b.AsInterface())
		})
		if diff := cmp.Diff(testcase.wantValue, got, opt); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
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

func TestMapValue(t *testing.T) {
	// Test basic map
	m := map[string]attribute.Value{
		"string": attribute.StringValue("test"),
		"int":    attribute.Int64Value(123),
		"float":  attribute.Float64Value(3.14),
		"bool":   attribute.BoolValue(true),
	}

	kv := attribute.Map("test", m)
	assert.Equal(t, attribute.MAP, kv.Value.Type())

	result := kv.Value.AsMap()
	assert.Equal(t, m, result)

	// Test nested map
	nested := map[string]attribute.Value{
		"outer": attribute.MapValue(map[string]attribute.Value{
			"inner": attribute.StringValue("nested value"),
		}),
	}

	kvNested := attribute.Map("nested", nested)
	assert.Equal(t, attribute.MAP, kvNested.Value.Type())

	resultNested := kvNested.Value.AsMap()
	assert.Equal(t, nested, resultNested)

	// Verify nested value can be extracted
	outerMap := resultNested["outer"].AsMap()
	assert.NotNil(t, outerMap)
	assert.Equal(t, "nested value", outerMap["inner"].AsString())

	// Test empty map
	emptyMap := map[string]attribute.Value{}
	kvEmpty := attribute.Map("empty", emptyMap)
	assert.Equal(t, attribute.MAP, kvEmpty.Value.Type())
	assert.Equal(t, emptyMap, kvEmpty.Value.AsMap())

	// Test AsInterface returns the map
	iface := kv.Value.AsInterface()
	mapIface, ok := iface.(map[string]attribute.Value)
	assert.True(t, ok)
	assert.Equal(t, m, mapIface)
}

func TestMapValue_DeepNesting(t *testing.T) {
	// Test deeply nested maps to verify recursive map<string, AnyValue> support
	// as per OpenTelemetry spec
	deepMap := map[string]attribute.Value{
		"level1": attribute.MapValue(map[string]attribute.Value{
			"level2": attribute.MapValue(map[string]attribute.Value{
				"level3": attribute.MapValue(map[string]attribute.Value{
					"string": attribute.StringValue("deep value"),
					"int":    attribute.Int64Value(123),
					"bool":   attribute.BoolValue(true),
					"slice":  attribute.StringSliceValue([]string{"a", "b"}),
				}),
				"sibling": attribute.StringValue("level2 value"),
			}),
		}),
		"topString": attribute.StringValue("top level"),
	}

	kv := attribute.Map("config", deepMap)
	assert.Equal(t, attribute.MAP, kv.Value.Type())

	// Navigate through nested structure
	result := kv.Value.AsMap()
	assert.NotNil(t, result)

	level1 := result["level1"].AsMap()
	assert.NotNil(t, level1)

	level2 := level1["level2"].AsMap()
	assert.NotNil(t, level2)

	level3 := level2["level3"].AsMap()
	assert.NotNil(t, level3)

	// Verify leaf values
	assert.Equal(t, "deep value", level3["string"].AsString())
	assert.Equal(t, int64(123), level3["int"].AsInt64())
	assert.Equal(t, true, level3["bool"].AsBool())
	assert.Equal(t, []string{"a", "b"}, level3["slice"].AsStringSlice())

	// Verify sibling at level2
	assert.Equal(t, "level2 value", level2["sibling"].AsString())

	// Verify top level
	assert.Equal(t, "top level", result["topString"].AsString())
}

func TestMapValue_NilMap(t *testing.T) {
	// nil map round-trips as an empty (non-nil) map,
	// consistent with how nil slices are handled.
	v := attribute.MapValue(nil)
	assert.Equal(t, attribute.MAP, v.Type())
	got := v.AsMap()
	assert.NotNil(t, got, "nil map should round-trip as non-nil empty map")
	assert.Len(t, got, 0)
}

func TestAsMap_WrongType(t *testing.T) {
	// AsMap on a non-MAP Value must return nil without panicking.
	assert.Nil(t, attribute.StringValue("hello").AsMap())
	assert.Nil(t, attribute.Int64Value(42).AsMap())
	assert.Nil(t, attribute.BoolValue(true).AsMap())
	assert.Nil(t, attribute.Float64Value(3.14).AsMap())
	assert.Nil(t, attribute.BoolSliceValue([]bool{true}).AsMap())
}

func BenchmarkMapValue(b *testing.B) {
	m := map[string]attribute.Value{
		"string":  attribute.StringValue("value"),
		"int":     attribute.Int64Value(42),
		"float":   attribute.Float64Value(3.14),
		"bool":    attribute.BoolValue(true),
		"strings": attribute.StringSliceValue([]string{"a", "b", "c"}),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = attribute.MapValue(m)
	}
}

func BenchmarkAsMap(b *testing.B) {
	v := attribute.MapValue(map[string]attribute.Value{
		"string":  attribute.StringValue("value"),
		"int":     attribute.Int64Value(42),
		"float":   attribute.Float64Value(3.14),
		"bool":    attribute.BoolValue(true),
		"strings": attribute.StringSliceValue([]string{"a", "b", "c"}),
	})
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = v.AsMap()
	}
}
