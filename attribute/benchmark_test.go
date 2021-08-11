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

	outBool    bool
	outInt64   int64
	outFloat64 float64
	outStr     string
)

func benchmarkAny(k string, v interface{}) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outKV = attribute.Any(k, v)
		}
	}
}

func benchmarkEmit(kv attribute.KeyValue) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outStr = kv.Value.Emit()
		}
	}
}

func benchmarkArray(k string, v interface{}) func(*testing.B) {
	a := attribute.Array(k, v)
	return func(b *testing.B) {
		b.Run("KeyValue", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				outKV = attribute.Array(k, v)
			}
		})
		b.Run("Emit", benchmarkEmit(a))
	}
}

func BenchmarkBool(b *testing.B) {
	k, v := "bool", true
	kv := attribute.Bool(k, v)
	array := []bool{true, false, true}

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
	b.Run("Any", benchmarkAny(k, v))
	b.Run("Emit", benchmarkEmit(kv))
	b.Run("Array", benchmarkArray(k, array))
}

func BenchmarkInt(b *testing.B) {
	k, v := "int", int(42)
	kv := attribute.Int(k, v)
	array := []int{42, -3, 12}

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
	b.Run("Any", benchmarkAny(k, v))
	b.Run("Emit", benchmarkEmit(kv))
	b.Run("Array", benchmarkArray(k, array))
}

func BenchmarkInt8(b *testing.B) {
	b.Run("Any", benchmarkAny("int8", int8(42)))
}

func BenchmarkInt16(b *testing.B) {
	b.Run("Any", benchmarkAny("int16", int16(42)))
}

func BenchmarkInt32(b *testing.B) {
	b.Run("Any", benchmarkAny("int32", int32(42)))
}

func BenchmarkInt64(b *testing.B) {
	k, v := "int64", int64(42)
	kv := attribute.Int64(k, v)
	array := []int64{42, -3, 12}

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
	b.Run("Any", benchmarkAny(k, v))
	b.Run("Emit", benchmarkEmit(kv))
	b.Run("Array", benchmarkArray(k, array))
}

func BenchmarkFloat64(b *testing.B) {
	k, v := "float64", float64(42)
	kv := attribute.Float64(k, v)
	array := []float64{42, -3, 12}

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
	b.Run("Any", benchmarkAny(k, v))
	b.Run("Emit", benchmarkEmit(kv))
	b.Run("Array", benchmarkArray(k, array))
}

func BenchmarkString(b *testing.B) {
	k, v := "string", "42"
	kv := attribute.String(k, v)
	array := []string{"forty-two", "negative three", "twelve"}

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
	b.Run("Any", benchmarkAny(k, v))
	b.Run("Emit", benchmarkEmit(kv))
	b.Run("Array", benchmarkArray(k, array))
}

func BenchmarkBytes(b *testing.B) {
	b.Run("Any", benchmarkAny("bytes", []byte("bytes")))
}

type test struct{}

func BenchmarkStruct(b *testing.B) {
	b.Run("Any", benchmarkAny("struct", test{}))
}
