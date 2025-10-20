// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// Store results in a file scope var to ensure compiler does not optimize the
// test away.
var (
	outV  attribute.Value
	outKV attribute.KeyValue

	outBool         bool
	outBoolSlice    []bool
	outInt64        int64
	outInt64Slice   []int64
	outFloat64      float64
	outFloat64Slice []float64
	outStr          string
	outStrSlice     []string
)

func benchmarkEmit(kv attribute.KeyValue) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStr = kv.Value.Emit()
		}
	}
}

func BenchmarkBool(b *testing.B) {
	k, v := "bool", true
	kv := attribute.Bool(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.BoolValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Bool(k, v)
		}
	})
	b.Run("AsBool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outBool = kv.Value.AsBool()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkBoolSlice(b *testing.B) {
	k, v := "bool slice", []bool{true, false, true}
	kv := attribute.BoolSlice(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.BoolSliceValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.BoolSlice(k, v)
		}
	})
	b.Run("AsBoolSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outBoolSlice = kv.Value.AsBoolSlice()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkInt(b *testing.B) {
	k, v := "int", int(42)
	kv := attribute.Int(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.IntValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Int(k, v)
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkIntSlice(b *testing.B) {
	k, v := "int slice", []int{42, -3, 12}
	kv := attribute.IntSlice(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.IntSliceValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.IntSlice(k, v)
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkInt64(b *testing.B) {
	k, v := "int64", int64(42)
	kv := attribute.Int64(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.Int64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Int64(k, v)
		}
	})
	b.Run("AsInt64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outInt64 = kv.Value.AsInt64()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkInt64Slice(b *testing.B) {
	k, v := "int64 slice", []int64{42, -3, 12}
	kv := attribute.Int64Slice(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.Int64SliceValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Int64Slice(k, v)
		}
	})
	b.Run("AsInt64Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outInt64Slice = kv.Value.AsInt64Slice()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkFloat64(b *testing.B) {
	k, v := "float64", float64(42)
	kv := attribute.Float64(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.Float64Value(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Float64(k, v)
		}
	})
	b.Run("AsFloat64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outFloat64 = kv.Value.AsFloat64()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkFloat64Slice(b *testing.B) {
	k, v := "float64 slice", []float64{42, -3, 12}
	kv := attribute.Float64Slice(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.Float64SliceValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Float64Slice(k, v)
		}
	})
	b.Run("AsFloat64Slice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outFloat64Slice = kv.Value.AsFloat64Slice()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkString(b *testing.B) {
	k, v := "string", "42"
	kv := attribute.String(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.StringValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.String(k, v)
		}
	})
	b.Run("AsString", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStr = kv.Value.AsString()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkStringSlice(b *testing.B) {
	k, v := "float64 slice", []string{"forty-two", "negative three", "twelve"}
	kv := attribute.StringSlice(k, v)

	b.Run("Value", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outV = attribute.StringSliceValue(v)
		}
	})
	b.Run("KeyValue", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.StringSlice(k, v)
		}
	})
	b.Run("AsStringSlice", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStrSlice = kv.Value.AsStringSlice()
		}
	})
	b.Run("Emit", benchmarkEmit(kv))
}

func BenchmarkSetEquals(b *testing.B) {
	b.Run("Empty", func(b *testing.B) {
		benchmarkSetEquals(b, attribute.EmptySet())
	})
	b.Run("1 string attribute", func(b *testing.B) {
		set := attribute.NewSet(attribute.String("string", "42"))
		benchmarkSetEquals(b, &set)
	})
	b.Run("10 string attributes", func(b *testing.B) {
		set := attribute.NewSet(
			attribute.String("a", "42"),
			attribute.String("b", "42"),
			attribute.String("c", "42"),
			attribute.String("d", "42"),
			attribute.String("e", "42"),
			attribute.String("f", "42"),
			attribute.String("g", "42"),
			attribute.String("h", "42"),
			attribute.String("i", "42"),
			attribute.String("j", "42"),
		)
		benchmarkSetEquals(b, &set)
	})
	b.Run("1 int attribute", func(b *testing.B) {
		set := attribute.NewSet(attribute.Int("string", 42))
		benchmarkSetEquals(b, &set)
	})
	b.Run("10 int attributes", func(b *testing.B) {
		set := attribute.NewSet(
			attribute.Int("a", 42),
			attribute.Int("b", 42),
			attribute.Int("c", 42),
			attribute.Int("d", 42),
			attribute.Int("e", 42),
			attribute.Int("f", 42),
			attribute.Int("g", 42),
			attribute.Int("h", 42),
			attribute.Int("i", 42),
			attribute.Int("j", 42),
		)
		benchmarkSetEquals(b, &set)
	})
}

func benchmarkSetEquals(b *testing.B, set *attribute.Set) {
	b.ResetTimer()
	for range b.N {
		if !set.Equals(set) {
			b.Fatal("not equal")
		}
	}
}

// BenchmarkEquivalentMapAccess measures how expensive it is to use
// Equivalent() as a map key. This is on the hot path for making synchronous
// measurements on the metrics API/SDK. It will likely be on the hot path for
// the trace and logs API/SDK in the future.
func BenchmarkEquivalentMapAccess(b *testing.B) {
	b.Run("Empty", func(b *testing.B) {
		benchmarkEquivalentMapAccess(b, attribute.EmptySet())
	})
	b.Run("1 string attribute", func(b *testing.B) {
		set := attribute.NewSet(attribute.String("string", "42"))
		benchmarkEquivalentMapAccess(b, &set)
	})
	b.Run("10 string attributes", func(b *testing.B) {
		set := attribute.NewSet(
			attribute.String("a", "42"),
			attribute.String("b", "42"),
			attribute.String("c", "42"),
			attribute.String("d", "42"),
			attribute.String("e", "42"),
			attribute.String("f", "42"),
			attribute.String("g", "42"),
			attribute.String("h", "42"),
			attribute.String("i", "42"),
			attribute.String("j", "42"),
		)
		benchmarkEquivalentMapAccess(b, &set)
	})
	b.Run("1 int attribute", func(b *testing.B) {
		set := attribute.NewSet(attribute.Int("string", 42))
		benchmarkEquivalentMapAccess(b, &set)
	})
	b.Run("10 int attributes", func(b *testing.B) {
		set := attribute.NewSet(
			attribute.Int("a", 42),
			attribute.Int("b", 42),
			attribute.Int("c", 42),
			attribute.Int("d", 42),
			attribute.Int("e", 42),
			attribute.Int("f", 42),
			attribute.Int("g", 42),
			attribute.Int("h", 42),
			attribute.Int("i", 42),
			attribute.Int("j", 42),
		)
		benchmarkEquivalentMapAccess(b, &set)
	})
}

func benchmarkEquivalentMapAccess(b *testing.B, set *attribute.Set) {
	values := map[attribute.Distinct]int{}
	b.ResetTimer()
	for range b.N {
		values[set.Equivalent()]++
	}
}
