// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func TestDefined(t *testing.T) {
	for _, testcase := range []struct {
		name string
		k    attribute.Key
		want bool
	}{
		{
			name: "Key.Defined() returns true when len(v.Name) != 0",
			k:    attribute.Key("foo"),
			want: true,
		},
		{
			name: "Key.Defined() returns false when len(v.Name) == 0",
			k:    attribute.Key(""),
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// func (k attribute.Key) Defined() bool {
			have := testcase.k.Defined()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestJSONValue(t *testing.T) {
	var kvs any = [2]attribute.KeyValue{
		attribute.String("A", "B"),
		attribute.Int64("C", 1),
	}

	data, err := json.Marshal(kvs)
	require.NoError(t, err)
	require.JSONEq(
		t,
		`[{"Key":"A","Value":{"Type":"STRING","Value":"B"}},{"Key":"C","Value":{"Type":"INT64","Value":1}}]`,
		string(data),
	)
}

func TestEmit(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    attribute.Value
		want string
	}{
		{
			name: `test Key.Emit() can emit a string representing self.BOOL`,
			v:    attribute.BoolValue(true),
			want: "true",
		},
		{
			name: `test Key.Emit() can emit a string representing self.BOOLSLICE`,
			v:    attribute.BoolSliceValue([]bool{true, false, true}),
			want: `[true false true]`,
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT64SLICE`,
			v:    attribute.Int64SliceValue([]int64{1, 42}),
			want: `[1,42]`,
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT64`,
			v:    attribute.Int64Value(42),
			want: "42",
		},
		{
			name: `test Key.Emit() can representing an int value`,
			v:    attribute.IntValue(7),
			want: "7",
		},
		{
			name: `test Key.Emit() can represent an []int value`,
			v:    attribute.IntSliceValue([]int{1, 2, 3}),
			want: `[1,2,3]`,
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT64SLICE`,
			v:    attribute.Float64SliceValue([]float64{1.0, 42.5}),
			want: `[1,42.5]`,
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT64`,
			v:    attribute.Float64Value(42.1),
			want: "42.1",
		},
		{
			name: `test Key.Emit() can emit a string representing self.STRING`,
			v:    attribute.StringValue("foo"),
			want: "foo",
		},
		{
			name: `test Key.Emit() can emit a string representing self.STRINGSLICE`,
			v:    attribute.StringSliceValue([]string{"foo", "bar"}),
			want: `["foo","bar"]`,
		},
		{
			name: `test Key.Emit() can emit a string representing self.MAP`,
			v: attribute.MapValue(map[string]attribute.Value{
				"key1": attribute.StringValue("value1"),
				"key2": attribute.Int64Value(42),
			}),
			want: `{"key1":"value1","key2":42}`,
		},
		{
			name: `test Key.Emit() can emit a string representing empty self.MAP`,
			v:    attribute.MapValue(nil),
			want: `{}`,
		},
		{
			name: `test Key.Emit() escapes quotes in MAP string values`,
			v: attribute.MapValue(map[string]attribute.Value{
				"key": attribute.StringValue(`hello "world"`),
			}),
			want: `{"key":"hello \"world\""}`,
		},
		{
			name: `test Key.Emit() escapes quotes in MAP keys`,
			v: attribute.MapValue(map[string]attribute.Value{
				`k"ey`: attribute.BoolValue(true),
			}),
			want: `{"k\"ey":true}`,
		},
		{
			name: `test Key.Emit() quotes NaN float64 in MAP`,
			v: attribute.MapValue(map[string]attribute.Value{
				"val": attribute.Float64Value(math.NaN()),
			}),
			want: `{"val":"NaN"}`,
		},
		{
			name: `test Key.Emit() quotes +Inf float64 in MAP`,
			v: attribute.MapValue(map[string]attribute.Value{
				"val": attribute.Float64Value(math.Inf(1)),
			}),
			want: `{"val":"Infinity"}`,
		},
		{
			name: `test Key.Emit() quotes -Inf float64 in MAP`,
			v: attribute.MapValue(map[string]attribute.Value{
				"val": attribute.Float64Value(math.Inf(-1)),
			}),
			want: `{"val":"-Infinity"}`,
		},
		{
			name: `test Key.Emit() emits null for INVALID values in MAP`,
			v: attribute.MapValue(map[string]attribute.Value{
				"a": {},
			}),
			want: `{"a":null}`,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// proto: func (v attribute.Value) Emit() string {
			have := testcase.v.Emit()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}
