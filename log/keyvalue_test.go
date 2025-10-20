// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
)

func TestKind(t *testing.T) {
	testCases := []struct {
		kind  log.Kind
		str   string
		value int
	}{
		{log.KindBool, "Bool", 1},
		{log.KindBytes, "Bytes", 5},
		{log.KindEmpty, "Empty", 0},
		{log.KindFloat64, "Float64", 2},
		{log.KindInt64, "Int64", 3},
		{log.KindSlice, "Slice", 6},
		{log.KindMap, "Map", 7},
		{log.KindString, "String", 4},
	}
	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			assert.Equal(t, tc.value, int(tc.kind), "Kind value")
			assert.Equal(t, tc.str, tc.kind.String(), "Kind string")
		})
	}
}

func TestValueEqual(t *testing.T) {
	vals := []log.Value{
		{},
		log.Int64Value(1),
		log.Int64Value(2),
		log.Int64Value(-2),
		log.Float64Value(3.5),
		log.Float64Value(3.7),
		log.BoolValue(true),
		log.BoolValue(false),
		log.StringValue("hi"),
		log.StringValue("bye"),
		log.BytesValue([]byte{1, 3, 5}),
		log.SliceValue(log.StringValue("foo")),
		log.SliceValue(log.IntValue(3), log.StringValue("foo")),
		log.MapValue(log.Bool("b", true), log.Int("i", 3)),
		log.MapValue(
			log.Slice("l", log.IntValue(3), log.StringValue("foo")),
			log.Bytes("b", []byte{3, 5, 7}),
			log.Empty("e"),
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
		value  log.Value
		value2 log.Value
	}{
		{
			value: log.MapValue(
				log.Slice("l", log.IntValue(3), log.StringValue("foo")),
				log.Bytes("b", []byte{3, 5, 7}),
				log.Empty("e"),
			),
			value2: log.MapValue(
				log.Bytes("b", []byte{3, 5, 7}),
				log.Slice("l", log.IntValue(3), log.StringValue("foo")),
				log.Empty("e"),
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
	v := log.Value{}
	t.Run("Value.Empty", func(t *testing.T) {
		assert.True(t, v.Empty())
	})

	t.Run("Bytes", func(t *testing.T) {
		assert.Nil(t, log.Bytes("b", nil).Value.AsBytes())
	})
	t.Run("Slice", func(t *testing.T) {
		assert.Nil(t, log.Slice("s").Value.AsSlice())
	})
	t.Run("Map", func(t *testing.T) {
		assert.Nil(t, log.Map("m").Value.AsMap())
	})
}

func TestEmptyGroupsPreserved(t *testing.T) {
	t.Run("Map", func(t *testing.T) {
		assert.Equal(t, []log.KeyValue{
			log.Int("a", 1),
			log.Map("g1", log.Map("g2")),
			log.Map("g3", log.Map("g4", log.Int("b", 2))),
		}, log.MapValue(
			log.Int("a", 1),
			log.Map("g1", log.Map("g2")),
			log.Map("g3", log.Map("g4", log.Int("b", 2))),
		).AsMap())
	})

	t.Run("Slice", func(t *testing.T) {
		assert.Equal(t, []log.Value{{}}, log.SliceValue(log.Value{}).AsSlice())
	})
}

func TestBool(t *testing.T) {
	const key, val = "boolKey", true
	kv := log.Bool(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindBool
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
	kv := log.Float64(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindFloat64
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
	kv := log.Int(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindInt64
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
	kv := log.Int64(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindInt64
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
	kv := log.String(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindString
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
	kv := log.Bytes(key, val)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindBytes
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
	val := []log.Value{log.IntValue(3), log.StringValue("foo")}
	kv := log.Slice(key, val...)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindSlice
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
	val := []log.KeyValue{
		log.Slice("l", log.IntValue(3), log.StringValue("foo")),
		log.Bytes("b", []byte{3, 5, 7}),
	}
	kv := log.Map(key, val...)
	testKV(t, key, kv)

	v, k := kv.Value, log.KindMap
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
	kv := log.Empty(key)

	assert.Equal(t, key, kv.Key, "incorrect key")
	assert.True(t, kv.Value.Empty(), "value not empty")

	v, k := kv.Value, log.KindEmpty
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
		v    log.Value
		want string
	}{
		{log.Int64Value(-3), "-3"},
		{log.Float64Value(.15), "0.15"},
		{log.BoolValue(true), "true"},
		{log.StringValue("foo"), "foo"},
		{log.BytesValue([]byte{2, 4, 6}), "[2 4 6]"},
		{log.SliceValue(log.IntValue(3), log.StringValue("foo")), "[3 foo]"},
		{log.MapValue(log.Int("a", 1), log.Bool("b", true)), "[a:1 b:true]"},
		{log.Value{}, "<nil>"},
	} {
		got := test.v.String()
		assert.Equal(t, test.want, got)
	}
}

func TestValueFromAttribute(t *testing.T) {
	testCases := []struct {
		desc string
		v    attribute.Value
		want log.Value
	}{
		{
			desc: "Empty",
			v:    attribute.Value{},
			want: log.Value{},
		},
		{
			desc: "Bool",
			v:    attribute.BoolValue(true),
			want: log.BoolValue(true),
		},
		{
			desc: "BoolSlice",
			v:    attribute.BoolSliceValue([]bool{true, false}),
			want: log.SliceValue(log.BoolValue(true), log.BoolValue(false)),
		},
		{
			desc: "Int64",
			v:    attribute.Int64Value(13),
			want: log.Int64Value(13),
		},
		{
			desc: "Int64Slice",
			v:    attribute.Int64SliceValue([]int64{12, 34}),
			want: log.SliceValue(log.Int64Value(12), log.Int64Value(34)),
		},
		{
			desc: "Float64",
			v:    attribute.Float64Value(3.14),
			want: log.Float64Value(3.14),
		},
		{
			desc: "Float64Slice",
			v:    attribute.Float64SliceValue([]float64{3.14, 2.72}),
			want: log.SliceValue(log.Float64Value(3.14), log.Float64Value(2.72)),
		},
		{
			desc: "String",
			v:    attribute.StringValue("foo"),
			want: log.StringValue("foo"),
		},
		{
			desc: "StringSlice",
			v:    attribute.StringSliceValue([]string{"foo", "bar"}),
			want: log.SliceValue(log.StringValue("foo"), log.StringValue("bar")),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := log.ValueFromAttribute(tc.v)
			if !got.Equal(tc.want) {
				t.Errorf("got: %v; want:%v", got, tc.want)
			}
		})
	}
}

func TestKeyValueFromAttribute(t *testing.T) {
	testCases := []struct {
		desc string
		kv   attribute.KeyValue
		want log.KeyValue
	}{
		{
			desc: "Empty",
			kv:   attribute.KeyValue{},
			want: log.KeyValue{},
		},
		{
			desc: "Bool",
			kv:   attribute.Bool("k", true),
			want: log.Bool("k", true),
		},
		{
			desc: "BoolSlice",
			kv:   attribute.BoolSlice("k", []bool{true, false}),
			want: log.Slice("k", log.BoolValue(true), log.BoolValue(false)),
		},
		{
			desc: "Int64",
			kv:   attribute.Int64("k", 13),
			want: log.Int64("k", 13),
		},
		{
			desc: "Int64Slice",
			kv:   attribute.Int64Slice("k", []int64{12, 34}),
			want: log.Slice("k", log.Int64Value(12), log.Int64Value(34)),
		},
		{
			desc: "Float64",
			kv:   attribute.Float64("k", 3.14),
			want: log.Float64("k", 3.14),
		},
		{
			desc: "Float64Slice",
			kv:   attribute.Float64Slice("k", []float64{3.14, 2.72}),
			want: log.Slice("k", log.Float64Value(3.14), log.Float64Value(2.72)),
		},
		{
			desc: "String",
			kv:   attribute.String("k", "foo"),
			want: log.String("k", "foo"),
		},
		{
			desc: "StringSlice",
			kv:   attribute.StringSlice("k", []string{"foo", "bar"}),
			want: log.Slice("k", log.StringValue("foo"), log.StringValue("bar")),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := log.KeyValueFromAttribute(tc.kv)
			if !got.Equal(tc.want) {
				t.Errorf("got: %v; want:%v", got, tc.want)
			}
		})
	}
}

type logSink struct {
	logr.LogSink

	err           error
	msg           string
	keysAndValues []any
}

func (l *logSink) Error(err error, msg string, keysAndValues ...any) {
	l.err, l.msg, l.keysAndValues = err, msg, keysAndValues
	l.LogSink.Error(err, msg, keysAndValues...)
}

func testErrKind[T any](f func() T, msg string, k log.Kind) func(*testing.T) {
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

func testKV(t *testing.T, key string, kv log.KeyValue) {
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
		slice []log.Value
		m     []log.KeyValue
	)

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		b = log.Bool(key, true).Value.AsBool()
	}), "Bool.AsBool")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		f = log.Float64(key, 3.0).Value.AsFloat64()
	}), "Float.AsFloat64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		i = log.Int(key, 9).Value.AsInt64()
	}), "Int.AsInt64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		i = log.Int64(key, 8).Value.AsInt64()
	}), "Int64.AsInt64")

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		s = log.String(key, "value").Value.AsString()
	}), "String.AsString")

	byteVal := []byte{1, 3, 4}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		by = log.Bytes(key, byteVal).Value.AsBytes()
	}), "Byte.AsBytes")

	sliceVal := []log.Value{log.BoolValue(true), log.IntValue(32)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		slice = log.Slice(key, sliceVal...).Value.AsSlice()
	}), "Slice.AsSlice")

	mapVal := []log.KeyValue{log.Bool("b", true), log.Int("i", 32)}
	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		m = log.Map(key, mapVal...).Value.AsMap()
	}), "Map.AsMap")

	// Convince the linter these values are used.
	_, _, _, _, _, _, _ = i, f, b, by, s, slice, m
}
