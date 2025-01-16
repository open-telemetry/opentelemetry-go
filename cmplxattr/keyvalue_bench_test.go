// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package cmplxattr_test

import (
	"testing"

	"go.opentelemetry.io/otel/cmplxattr"
)

// Store results in a file scope var to ensure compiler does not optimize the
// test away.
var (
	outV  cmplxattr.Value
	outKV cmplxattr.KeyValue

	outBool    bool
	outFloat64 float64
	outInt64   int64
	outMap     []cmplxattr.KeyValue
	outSlice   []cmplxattr.Value
	outStr     string
)

func BenchmarkBool(b *testing.B) {
	const k, v = "bool", true

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.BoolValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Bool(k, v)
		}
	})

	kv := cmplxattr.Bool(k, v)
	b.Run("AsBool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outBool = kv.Value.AsBool()
		}
	})
}

func BenchmarkFloat64(b *testing.B) {
	const k, v = "float64", 3.0

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.Float64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Float64(k, v)
		}
	})

	kv := cmplxattr.Float64(k, v)
	b.Run("AsFloat64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outFloat64 = kv.Value.AsFloat64()
		}
	})
}

func BenchmarkInt(b *testing.B) {
	const k, v = "int", 32

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.IntValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Int(k, v)
		}
	})

	kv := cmplxattr.Int(k, v)
	b.Run("AsInt64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outInt64 = kv.Value.AsInt64()
		}
	})
}

func BenchmarkInt64(b *testing.B) {
	const k, v = "int64", int64(32)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.Int64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Int64(k, v)
		}
	})

	kv := cmplxattr.Int64(k, v)
	b.Run("AsInt64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outInt64 = kv.Value.AsInt64()
		}
	})
}

func BenchmarkMap(b *testing.B) {
	const k = "map"
	v := []cmplxattr.KeyValue{cmplxattr.Bool("b", true), cmplxattr.Int("i", 1)}

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.MapValue(v...)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Map(k, v...)
		}
	})

	kv := cmplxattr.Map(k, v...)
	b.Run("AsMap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outMap = kv.Value.AsMap()
		}
	})
}

func BenchmarkSlice(b *testing.B) {
	const k = "slice"
	v := []cmplxattr.Value{cmplxattr.BoolValue(true), cmplxattr.IntValue(1)}

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.SliceValue(v...)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.Slice(k, v...)
		}
	})

	kv := cmplxattr.Slice(k, v...)
	b.Run("AsSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outSlice = kv.Value.AsSlice()
		}
	})
}

func BenchmarkString(b *testing.B) {
	const k, v = "str", "value"

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = cmplxattr.StringValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = cmplxattr.String(k, v)
		}
	})

	kv := cmplxattr.String(k, v)
	b.Run("AsString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStr = kv.Value.AsString()
		}
	})
}

func BenchmarkValueEqual(b *testing.B) {
	vals := []cmplxattr.Value{
		{},
		cmplxattr.Int64Value(1),
		cmplxattr.Int64Value(2),
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
	for _, v1 := range vals {
		for _, v2 := range vals {
			b.Run(v1.String()+" with "+v2.String(), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = v1.Equal(v2)
				}
			})
		}
	}
}
