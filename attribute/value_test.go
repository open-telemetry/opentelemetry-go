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
			name:      "empty value",
			value:     attribute.Value{},
			wantType:  attribute.EMPTY,
			wantValue: nil,
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

func TestMapValue(t *testing.T) {
	t.Run("TypeIsMAP", func(t *testing.T) {
		v := attribute.MapValue([]attribute.KeyValue{
			attribute.String("k", "v"),
		})
		assert.Equal(t, attribute.MAP, v.Type())
		assert.Equal(t, "MAP", v.Type().String())
	})

	t.Run("RoundTrip", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			attribute.String("host", "localhost"),
			attribute.Int("port", 8080),
			attribute.Bool("tls", true),
		}
		v := attribute.MapValue(kvs)
		got := v.AsMap()
		assert.Equal(t, kvs, got)
	})

	t.Run("EmptyMap", func(t *testing.T) {
		v := attribute.MapValue(nil)
		assert.Equal(t, attribute.MAP, v.Type())
		assert.Empty(t, v.AsMap())
	})

	t.Run("AsMapWrongType", func(t *testing.T) {
		v := attribute.StringValue("not a map")
		assert.Nil(t, v.AsMap())
	})

	t.Run("AsInterfaceReturnsSlice", func(t *testing.T) {
		kvs := []attribute.KeyValue{attribute.String("a", "b")}
		v := attribute.MapValue(kvs)
		got, ok := v.AsInterface().([]attribute.KeyValue)
		assert.True(t, ok, "AsInterface() should return []KeyValue for MAP type")
		assert.Equal(t, kvs, got)
	})
}

func TestMapValueDeduplication(t *testing.T) {
	t.Run("LastWriteWins", func(t *testing.T) {
		// "a" appears twice: last value ("second") must win.
		kvs := []attribute.KeyValue{
			attribute.String("a", "first"),
			attribute.String("b", "only"),
			attribute.String("a", "second"),
		}
		got := attribute.MapValue(kvs).AsMap()
		want := []attribute.KeyValue{
			attribute.String("a", "second"), // position preserved; value updated
			attribute.String("b", "only"),
		}
		assert.Equal(t, want, got)
	})

	t.Run("AllDuplicates", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			attribute.Int("x", 1),
			attribute.Int("x", 2),
			attribute.Int("x", 3),
		}
		got := attribute.MapValue(kvs).AsMap()
		want := []attribute.KeyValue{attribute.Int("x", 3)}
		assert.Equal(t, want, got)
	})

	t.Run("NoDuplicates", func(t *testing.T) {
		kvs := []attribute.KeyValue{
			attribute.String("p", "1"),
			attribute.String("q", "2"),
		}
		got := attribute.MapValue(kvs).AsMap()
		assert.Equal(t, kvs, got)
	})

	t.Run("OrderOfFirstOccurrencePreserved", func(t *testing.T) {
		// Keys: c, a, b — with "a" duplicated. Result should be [c, a, b].
		kvs := []attribute.KeyValue{
			attribute.String("c", "c-val"),
			attribute.String("a", "a-first"),
			attribute.String("b", "b-val"),
			attribute.String("a", "a-second"),
		}
		got := attribute.MapValue(kvs).AsMap()
		want := []attribute.KeyValue{
			attribute.String("c", "c-val"),
			attribute.String("a", "a-second"),
			attribute.String("b", "b-val"),
		}
		assert.Equal(t, want, got)
	})
}

func TestMapValueEmit(t *testing.T) {
	t.Run("SingleEntry", func(t *testing.T) {
		v := attribute.MapValue([]attribute.KeyValue{
			attribute.String("key", "value"),
		})
		emit := v.Emit()
		// JSON object — must contain the key and value.
		assert.Contains(t, emit, `"key"`)
		assert.Contains(t, emit, `"value"`)
	})

	t.Run("EmptyMap", func(t *testing.T) {
		v := attribute.MapValue(nil)
		assert.Equal(t, "{}", v.Emit())
	})

	t.Run("MultipleEntries", func(t *testing.T) {
		v := attribute.MapValue([]attribute.KeyValue{
			attribute.Int("count", 42),
			attribute.Bool("ok", true),
		})
		emit := v.Emit()
		assert.Contains(t, emit, `"count"`)
		assert.Contains(t, emit, `"ok"`)
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
}
