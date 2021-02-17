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

package label_test

import (
	"testing"

	"go.opentelemetry.io/otel/label"
)

type test struct{}

var (
	arrayVal    = []string{"one", "two"}
	arrayKeyVal = label.Array("array", arrayVal)

	boolVal    = true
	boolKeyVal = label.Bool("bool", boolVal)

	intVal    = int(1)
	intKeyVal = label.Int("int", intVal)

	int8Val    = int8(1)
	int8KeyVal = label.Int("int8", int(int8Val))

	int16Val    = int16(1)
	int16KeyVal = label.Int("int16", int(int16Val))

	int64Val    = int64(1)
	int64KeyVal = label.Int64("int64", int64Val)

	float64Val    = float64(1.0)
	float64KeyVal = label.Float64("float64", float64Val)

	stringVal    = "string"
	stringKeyVal = label.String("string", stringVal)

	bytesVal  = []byte("bytes")
	structVal = test{}
)

func BenchmarkArrayKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Array("array", arrayVal)
	}
}

func BenchmarkArrayKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("array", arrayVal)
	}
}

func BenchmarkBoolKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Bool("bool", boolVal)
	}
}

func BenchmarkBoolKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("bool", boolVal)
	}
}

func BenchmarkIntKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Int("int", intVal)
	}
}

func BenchmarkIntKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("int", intVal)
	}
}

func BenchmarkInt8KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("int8", int8Val)
	}
}

func BenchmarkInt16KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("int16", int16Val)
	}
}

func BenchmarkInt64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Int64("int64", int64Val)
	}
}

func BenchmarkInt64KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("int64", int64Val)
	}
}

func BenchmarkFloat64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Float64("float64", float64Val)
	}
}

func BenchmarkFloat64KeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("float64", float64Val)
	}
}

func BenchmarkStringKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.String("string", stringVal)
	}
}

func BenchmarkStringKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("string", stringVal)
	}
}

func BenchmarkBytesKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("bytes", bytesVal)
	}
}

func BenchmarkStructKeyAny(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = label.Any("struct", structVal)
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
