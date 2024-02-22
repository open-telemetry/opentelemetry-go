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

package log_test

import (
	"testing"

	"go.opentelemetry.io/otel/log"
)

// Store results in a file scope var to ensure compiler does not optimize the
// test away.
var (
	outV  log.Value
	outKV log.KeyValue

	outAny     any
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
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
	b.Run("AsAny", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			outAny = kv.Value.AsAny()
		}
	})
}
