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

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log_test

import (
	golog "log"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		),
	}
	for i, v1 := range vals {
		for j, v2 := range vals {
			assert.Equal(t, i == j, v1.Equal(v2), "%v.Equal(%v)", v1, v2)
		}
	}
}

func TestEmpty(t *testing.T) {
	v := log.Value{}
	t.Run("Value.Empty", func(t *testing.T) {
		assert.True(t, v.Empty())
	})
	t.Run("Value.AsAny", func(t *testing.T) {
		assert.Nil(t, v.AsAny())
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
	const key, val = "key", true
	kv := log.Bool(key, val)
	testKV(t, key, val, kv)

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
	const key, val = "key", 3.0
	kv := log.Float64(key, val)
	testKV(t, key, val, kv)

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
	const key, val = "key", 1
	kv := log.Int(key, val)
	testKV[int64](t, key, val, kv)

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
	const key, val = "key", 1
	kv := log.Int64(key, val)
	testKV[int64](t, key, val, kv)

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
	const key, val = "key", "test string value"
	kv := log.String(key, val)
	testKV(t, key, val, kv)

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
	const key = "key"
	val := []byte{3, 2, 1}
	kv := log.Bytes(key, val)
	testKV(t, key, val, kv)

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
	const key = "key"
	val := []log.Value{log.IntValue(3), log.StringValue("foo")}
	kv := log.Slice(key, val...)
	testKV(t, key, val, kv)

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
	const key = "key"
	val := []log.KeyValue{
		log.Slice("l", log.IntValue(3), log.StringValue("foo")),
		log.Bytes("b", []byte{3, 5, 7}),
	}
	kv := log.Map(key, val...)
	testKV(t, key, val, kv)

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

type logSink struct {
	logr.LogSink

	err           error
	msg           string
	keysAndValues []interface{}
}

func (l *logSink) Error(err error, msg string, keysAndValues ...interface{}) {
	l.err, l.msg, l.keysAndValues = err, msg, keysAndValues
	l.LogSink.Error(err, msg, keysAndValues)
}

var stdLogger = stdr.New(golog.New(os.Stderr, "", golog.LstdFlags|golog.Lshortfile))

func testErrKind[T any](f func() T, msg string, k log.Kind) func(*testing.T) {
	return func(t *testing.T) {
		l := &logSink{LogSink: testr.New(t).GetSink()}
		global.SetLogger(logr.New(l))
		t.Cleanup(func() { global.SetLogger(stdLogger) })

		assert.Zero(t, f())

		assert.ErrorContains(t, l.err, "invalid Kind")
		assert.Equal(t, msg, l.msg)
		require.Len(t, l.keysAndValues, 2, "logged attributes")
		assert.Equal(t, l.keysAndValues[1], k)
	}
}

func testKV[T any](t *testing.T, key string, val T, kv log.KeyValue) {
	t.Helper()

	assert.Equal(t, key, kv.Key, "incorrect key")
	assert.False(t, kv.Value.Empty(), "value empty")
	assert.Equal(t, kv.Value.AsAny(), T(val), "AsAny wrong value")
}
