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

type test struct{}

var (
	arrayVal    = []string{"one", "two"}
	arrayKeyVal = attribute.Array("array", arrayVal)

	boolVal    = true
	boolKeyVal = attribute.Bool("bool", boolVal)

	intVal    = int(1)
	intKeyVal = attribute.Int("int", intVal)

	int8Val    = int8(1)
	int8KeyVal = attribute.Int("int8", int(int8Val))

	int16Val    = int16(1)
	int16KeyVal = attribute.Int("int16", int(int16Val))

	int64Val    = int64(1)
	int64KeyVal = attribute.Int64("int64", int64Val)

	float64Val    = float64(1.0)
	float64KeyVal = attribute.Float64("float64", float64Val)

	stringVal    = "string"
	stringKeyVal = attribute.String("string", stringVal)

	bytesVal  = []byte("bytes")
	structVal = test{}
)

func BenchmarkArrayKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Array("array", arrayVal)
	}
}

func BenchmarkArrayKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("array", arrayVal)
	}
}

func BenchmarkBoolKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Bool("bool", boolVal)
	}
}

func BenchmarkBoolKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("bool", boolVal)
	}
}

func BenchmarkIntKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Int("int", intVal)
	}
}

func BenchmarkIntKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("int", intVal)
	}
}

func BenchmarkInt8KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("int8", int8Val)
	}
}

func BenchmarkInt16KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("int16", int16Val)
	}
}

func BenchmarkInt64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Int64("int64", int64Val)
	}
}

func BenchmarkInt64KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("int64", int64Val)
	}
}

func BenchmarkFloat64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Float64("float64", float64Val)
	}
}

func BenchmarkFloat64KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("float64", float64Val)
	}
}

func BenchmarkStringKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.String("string", stringVal)
	}
}

func BenchmarkStringKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("string", stringVal)
	}
}

func BenchmarkBytesKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("bytes", bytesVal)
	}
}

func BenchmarkStructKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = attribute.Any("struct", structVal)
	}
}

func BenchmarkEmitArray(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = arrayKeyVal.Value.Emit()
	}
}

func BenchmarkEmitBool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = boolKeyVal.Value.Emit()
	}
}

func BenchmarkEmitInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = intKeyVal.Value.Emit()
	}
}

func BenchmarkEmitInt8(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = int8KeyVal.Value.Emit()
	}
}

func BenchmarkEmitInt16(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = int16KeyVal.Value.Emit()
	}
}

func BenchmarkEmitInt64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = int64KeyVal.Value.Emit()
	}
}

func BenchmarkEmitFloat64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = float64KeyVal.Value.Emit()
	}
}

func BenchmarkEmitString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = stringKeyVal.Value.Emit()
	}
}
