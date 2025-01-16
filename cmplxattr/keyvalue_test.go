// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmplxattr_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/cmplxattr"
	"go.opentelemetry.io/otel/internal/global"
)

func TestKind(t *testing.T) {
	testCases := []struct {
		kind  cmplxattr.Kind
		str   string
		value int
	}{
		{cmplxattr.KindBool, "Bool", 1},
		{cmplxattr.KindBytes, "Bytes", 5},
		{cmplxattr.KindEmpty, "Empty", 0},
		{cmplxattr.KindFloat64, "Float64", 2},
		{cmplxattr.KindInt64, "Int64", 3},
		{cmplxattr.KindSlice, "Slice", 6},
		{cmplxattr.KindMap, "Map", 7},
		{cmplxattr.KindString, "String", 4},
	}
	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			assert.Equal(t, tc.value, int(tc.kind), "Kind value")
			assert.Equal(t, tc.str, tc.kind.String(), "Kind string")
		})
	}
}

func TestValueEqual(t *testing.T) {
	vals := []cmplxattr.Value{
		{},
		cmplxattr.Int64Value(1),
		cmplxattr.Int64Value(2),
		cmplxattr.Int64Value(-2),
		cmplxattr.Float64Value(3.5),
		cmplxattr.Float64Value(3.7),
		cmplxattr.BoolValue(true),
		cmplxattr.BoolValue(false),
		cmplxattr.StringValue("hi"),
		cmplxattr.StringValue("bye"),
		cmplxattr.BytesValue([]byte{1, 3, 5}),
		cmplxattr.SliceValue(cmplxattr.StringValue("foo")),
		cmplxattr.SliceValue(cmplxattr.IntValue(3), cmplxattr.StringValue("foo")),
		cmplxattr.MapValue(cmplxattr.Bool("b", true), cmplxattr.Int("i", 3)),
		cmplxattr.MapValue(
			cmplxattr.Slice("l", cmplxattr.IntValue(3), cmplxattr.StringValue("foo")),
			cmplxattr.Bytes("b", []byte{3, 5, 7}),
			cmplxattr.Empty("e"),
		),
	}
	for i, v1 := range vals {
		for j, v2 := range vals {
			assert.Equal(t, i == j, v1.Equal(v2), "%v.Equal(%v)", v1, v2)
		}
	}
}

func TestSortedValueEqual(t *testing.T) {
	testCases := []struct {
		value  cmplxattr.Value
		value2 cmplxattr.Value
	}{
		{
			value: cmplxattr.MapValue(
				cmplxattr.Slice("l", cmplxattr.IntValue(3), cmplxattr.StringValue("foo")),
				cmplxattr.Bytes("b", []byte{3, 5, 7}),
				cmplxattr.Empty("e"),
			),
			value2: cmplxattr.MapValue(
				cmplxattr.Bytes("b", []byte{3, 5, 7}),
				cmplxattr.Slice("l", cmplxattr.IntValue(3), cmplxattr.StringValue("foo")),
				cmplxattr.Empty("e"),
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.value.String(), func(t *testing.T) {
			assert.True(t, tc.value.Equal(tc.value2), "%v.Equal(%v)", tc.value, tc.value2)
		})
	}
}

func TestValueEmpty(t *testing.T) {
	v := cmplxattr.Value{}
	t.Run("Value.Empty", func(t *testing.T) {
		assert.True(t, v.Empty())
	})

	t.Run("Bytes", func(t *testing.T) {
		assert.Nil(t, cmplxattr.Bytes("b", nil).Value.AsBytes())
	})
	t.Run("Slice", func(t *testing.T) {
		assert.Nil(t, cmplxattr.Slice("s").Value.AsSlice())
	})
	t.Run("Map", func(t *testing.T) {
		assert.Nil(t, cmplxattr.Map("m").Value.AsMap())
	})
}

func TestEmptyGroupsPreserved(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		assert.Equal(t, []cmplxattr.KeyValue{
			cmplxattr.Int("a", 1),
			cmplxattr.Map("g1", cmplxattr.Map("g2")),
			cmplxattr.Map("g3", cmplxattr.Map("g4", cmplxattr.Int("b", 2))),
		}, cmplxattr.MapValue(
			cmplxattr.Int("a", 1),
			cmplxattr.Map("g1", cmplxattr.Map("g2")),
			cmplxattr.Map("g3", cmplxattr.Map("g4", cmplxattr.Int("b", 2))),
		).AsMap())
	})

	t.Run("Slice", func(t *testing.T) {
		assert.Equal(t, []cmplxattr.Value{{}}, cmplxattr.SliceValue(cmplxattr.Value{}).AsSlice())
	})
}

func TestBool(t *testing.T) {
	const key, val = "boolKey", true
	kv := cmplxattr.Bool(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindBool
	t.Run("AsBool", func(t *testing.T) {
		assert.Equal(t, val, kv.Value.AsBool(), "AsBool")
	})
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestFloat64(t *testing.T) {
	const key, val = "float64Key", 3.0
	kv := cmplxattr.Float64(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindFloat64
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", func(t *testing.T) {
		assert.Equal(t, val, v.AsFloat64(), "AsFloat64")
	})
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestInt(t *testing.T) {
	const key, val = "intKey", 1
	kv := cmplxattr.Int(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindInt64
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", func(t *testing.T) {
		assert.Equal(t, int64(val), v.AsInt64(), "AsInt64")
	})
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestInt64(t *testing.T) {
	const key, val = "int64Key", 1
	kv := cmplxattr.Int64(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindInt64
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", func(t *testing.T) {
		assert.Equal(t, int64(val), v.AsInt64(), "AsInt64")
	})
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestString(t *testing.T) {
	const key, val = "stringKey", "test string value"
	kv := cmplxattr.String(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindString
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", func(t *testing.T) {
		assert.Equal(t, val, v.AsString(), "AsString")
	})
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestBytes(t *testing.T) {
	const key = "bytesKey"
	val := []byte{3, 2, 1}
	kv := cmplxattr.Bytes(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindBytes
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", func(t *testing.T) {
		assert.Equal(t, val, v.AsBytes(), "AsBytes")
	})
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestSlice(t *testing.T) {
	const key = "sliceKey"
	val := []cmplxattr.Value{cmplxattr.IntValue(3), cmplxattr.StringValue("foo")}
	kv := cmplxattr.Slice(key, val...)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindSlice
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", func(t *testing.T) {
		assert.Equal(t, val, v.AsSlice(), "AsSlice")
	})
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestMap(t *testing.T) {
	const key = "mapKey"
	val := []cmplxattr.KeyValue{
		cmplxattr.Slice("l", cmplxattr.IntValue(3), cmplxattr.StringValue("foo")),
		cmplxattr.Bytes("b", []byte{3, 5, 7}),
	}
	kv := cmplxattr.Map(key, val...)
	testKV(t, key, kv)

	v, k := kv.Value, cmplxattr.KindMap
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", func(t *testing.T) {
		assert.Equal(t, val, v.AsMap(), "AsMap")
	})
}

func TestEmpty(t *testing.T) {
	const key = "key"
	kv := cmplxattr.Empty(key)

	assert.Equal(t, key, kv.Key, "incorrect key")
	assert.True(t, kv.Value.Empty(), "value not empty")

	v, k := kv.Value, cmplxattr.KindEmpty
	t.Run("AsBool", testErrKind(v.AsBool, "AsBool", k))
	t.Run("AsFloat64", testErrKind(v.AsFloat64, "AsFloat64", k))
	t.Run("AsInt64", testErrKind(v.AsInt64, "AsInt64", k))
	t.Run("AsString", testErrKind(v.AsString, "AsString", k))
	t.Run("AsBytes", testErrKind(v.AsBytes, "AsBytes", k))
	t.Run("AsSlice", testErrKind(v.AsSlice, "AsSlice", k))
	t.Run("AsMap", testErrKind(v.AsMap, "AsMap", k))
}

func TestValueString(t *testing.T) {
	for _, test := range []struct {
		v    cmplxattr.Value
		want string
	}{
		{cmplxattr.Int64Value(-3), "-3"},
		{cmplxattr.Float64Value(.15), "0.15"},
		{cmplxattr.BoolValue(true), "true"},
		{cmplxattr.StringValue("foo"), "foo"},
		{cmplxattr.BytesValue([]byte{2, 4, 6}), "[2 4 6]"},
		{cmplxattr.SliceValue(cmplxattr.IntValue(3), cmplxattr.StringValue("foo")), "[3 foo]"},
		{cmplxattr.MapValue(cmplxattr.Int("a", 1), cmplxattr.Bool("b", true)), "[a:1 b:true]"},
		{cmplxattr.Value{}, "<nil>"},
	} {
		got := test.v.String()
		assert.Equal(t, test.want, got)
	}
}

type logSink struct {
	logr.LogSink

	err           error
	msg           string
	keysAndValues []interface{}
}

func (l *logSink) Error(err error, msg string, keysAndValues ...interface{}) {
	l.err, l.msg, l.keysAndValues = err, msg, keysAndValues
	l.LogSink.Error(err, msg, keysAndValues...)
}

func testErrKind[T any](f func() T, msg string, k cmplxattr.Kind) func(*testing.T) {
	return func(t *testing.T) {
		t.Cleanup(func(l logr.Logger) func() {
			return func() { global.SetLogger(l) }
		}(global.GetLogger()))

		l := &logSink{LogSink: testr.New(t).GetSink()}
		global.SetLogger(logr.New(l))

		assert.Zero(t, f())

		assert.ErrorContains(t, l.err, "invalid Kind")
		assert.Equal(t, msg, l.msg)
		require.Len(t, l.keysAndValues, 2, "logged attributes")
		assert.Equal(t, l.keysAndValues[1], k)
	}
}

func testKV(t *testing.T, key string, kv cmplxattr.KeyValue) {
	t.Helper()

	assert.Equal(t, key, kv.Key, "incorrect key")
	assert.False(t, kv.Value.Empty(), "value empty")
}

func TestAllocationLimits(t *testing.T) {
	const (
		runs = 5
		key  = "key"
	)

	// Assign testing results to external scope so the compiler doesn't
	// optimize away the testing statements.
	var (
		i     int64
		f     float64
		b     bool
		by    []byte
		s     string
		slice []cmplxattr.Value
		m     []cmplxattr.KeyValue
	)

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		b = cmplxattr.Bool(key, true).Value.AsBool()
	}), "Bool.AsBool")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		f = cmplxattr.Float64(key, 3.0).Value.AsFloat64()
	}), "Float.AsFloat64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		i = cmplxattr.Int(key, 9).Value.AsInt64()
	}), "Int.AsInt64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		i = cmplxattr.Int64(key, 8).Value.AsInt64()
	}), "Int64.AsInt64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		s = cmplxattr.String(key, "value").Value.AsString()
	}), "String.AsString")

	byteVal := []byte{1, 3, 4}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		by = cmplxattr.Bytes(key, byteVal).Value.AsBytes()
	}), "Byte.AsBytes")

	sliceVal := []cmplxattr.Value{cmplxattr.BoolValue(true), cmplxattr.IntValue(32)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		slice = cmplxattr.Slice(key, sliceVal...).Value.AsSlice()
	}), "Slice.AsSlice")

	mapVal := []cmplxattr.KeyValue{cmplxattr.Bool("b", true), cmplxattr.Int("i", 32)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		m = cmplxattr.Map(key, mapVal...).Value.AsMap()
	}), "Map.AsMap")

	// Convince the linter these values are used.
	_, _, _, _, _, _, _ = i, f, b, by, s, slice, m
}
