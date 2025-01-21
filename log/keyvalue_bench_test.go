// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

// Store results in a file scope var to ensure compiler does not optimize the
// test away.
var (
	outV  log.Value
	outKV log.KeyValue

	outBool    bool
	outFloat64 float64
	outInt64   int64
	outMap     []log.KeyValue
	outSlice   []log.Value
	outStr     string
)

func BenchmarkBool(b *testing.B) {
	const k, v = "bool", true

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = log.BoolValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Bool(k, v)
		}
	})

	kv := log.Bool(k, v)
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
			outV = log.Float64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Float64(k, v)
		}
	})

	kv := log.Float64(k, v)
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
			outV = log.IntValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Int(k, v)
		}
	})

	kv := log.Int(k, v)
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
			outV = log.Int64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Int64(k, v)
		}
	})

	kv := log.Int64(k, v)
	b.Run("AsInt64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outInt64 = kv.Value.AsInt64()
		}
	})
}

func BenchmarkMap(b *testing.B) {
	const k = "map"
	v := []log.KeyValue{log.Bool("b", true), log.Int("i", 1)}

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = log.MapValue(v...)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Map(k, v...)
		}
	})

	kv := log.Map(k, v...)
	b.Run("AsMap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outMap = kv.Value.AsMap()
		}
	})
}

func BenchmarkSlice(b *testing.B) {
	const k = "slice"
	v := []log.Value{log.BoolValue(true), log.IntValue(1)}

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = log.SliceValue(v...)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.Slice(k, v...)
		}
	})

	kv := log.Slice(k, v...)
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
			outV = log.StringValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = log.String(k, v)
		}
	})

	kv := log.String(k, v)
	b.Run("AsString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStr = kv.Value.AsString()
		}
	})
}

func BenchmarkValueEqual(b *testing.B) {
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
			log.Empty("e"),
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

func BenchmarkKeyValueFromAttribute(b *testing.B) {
	testCases := []struct {
		desc string
		kv   attribute.KeyValue
	}{
		{
			desc: "Empty",
			kv:   attribute.KeyValue{},
		},
		{
			desc: "Bool",
			kv:   attribute.Bool("k", true),
		},
		{
			desc: "BoolSlice",
			kv:   attribute.BoolSlice("k", []bool{true, false}),
		},
		{
			desc: "Int64",
			kv:   attribute.Int64("k", 13),
		},
		{
			desc: "Int64Slice",
			kv:   attribute.Int64Slice("k", []int64{12, 34}),
		},
		{
			desc: "Float64",
			kv:   attribute.Float64("k", 3.14),
		},
		{
			desc: "Float64Slice",
			kv:   attribute.Float64Slice("k", []float64{3.14, 2.72}),
		},
		{
			desc: "String",
			kv:   attribute.String("k", "foo"),
		},
		{
			desc: "StringSlice",
			kv:   attribute.StringSlice("k", []string{"foo", "bar"}),
		},
	}
	for _, tc := range testCases {
		b.Run(tc.desc, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				outKV = log.KeyValueFromAttribute(tc.kv)
			}
		})
	}
}
