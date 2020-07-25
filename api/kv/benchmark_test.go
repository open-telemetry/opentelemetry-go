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

package kv_test

import (
	"testing"

	"go.opentelemetry.io/otel/api/kv"
)

type test struct{}

var (
	arrayVal    = []string{"one", "two"}
	arrayKeyVal = kv.Array("array", arrayVal)

	boolVal    = true
	boolKeyVal = kv.Bool("bool", boolVal)

	intVal    = int(1)
	intKeyVal = kv.Int("int", intVal)

	int8Val    = int8(1)
	int8KeyVal = kv.Int("int8", int(int8Val))

	int16Val    = int16(1)
	int16KeyVal = kv.Int("int16", int(int16Val))

	int32Val    = int32(1)
	int32KeyVal = kv.Int32("int32", int32Val)

	int64Val    = int64(1)
	int64KeyVal = kv.Int64("int64", int64Val)

	uintVal    = uint(1)
	uintKeyVal = kv.Uint("uint", uintVal)

	uint8Val    = uint8(1)
	uint8KeyVal = kv.Uint("uint8", uint(uint8Val))

	uint16Val    = uint16(1)
	uint16KeyVal = kv.Uint("uint16", uint(uint16Val))

	uint32Val    = uint32(1)
	uint32KeyVal = kv.Uint32("uint32", uint32Val)

	uint64Val    = uint64(1)
	uint64KeyVal = kv.Uint64("uint64", uint64Val)

	float32Val    = float32(1.0)
	float32KeyVal = kv.Float32("float32", float32Val)

	float64Val    = float64(1.0)
	float64KeyVal = kv.Float64("float64", float64Val)

	stringVal    = "string"
	stringKeyVal = kv.String("string", stringVal)

	bytesVal  = []byte("bytes")
	structVal = test{}
)

func BenchmarkArrayKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Array("array", arrayVal)
	}
}

func BenchmarkArrayKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("array", arrayVal)
	}
}

func BenchmarkBoolKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Bool("bool", boolVal)
	}
}

func BenchmarkBoolKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("bool", boolVal)
	}
}

func BenchmarkIntKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Int("int", intVal)
	}
}

func BenchmarkIntKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("int", intVal)
	}
}

func BenchmarkInt8KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("int8", int8Val)
	}
}

func BenchmarkInt16KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("int16", int16Val)
	}
}

func BenchmarkInt32Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Int32("int32", int32Val)
	}
}

func BenchmarkInt32KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("int32", int32Val)
	}
}

func BenchmarkInt64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Int64("int64", int64Val)
	}
}

func BenchmarkInt64KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("int64", int64Val)
	}
}

func BenchmarkUintKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Uint("uint", uintVal)
	}
}

func BenchmarkUintKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("uint", uintVal)
	}
}

func BenchmarkUint8KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("uint8", uint8Val)
	}
}

func BenchmarkUint16KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("uint16", uint16Val)
	}
}

func BenchmarkUint32Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Uint32("uint32", uint32Val)
	}
}

func BenchmarkUint32KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("uint32", uint32Val)
	}
}

func BenchmarkUint64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Uint64("uint64", uint64Val)
	}
}

func BenchmarkUint64KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("uint64", uint64Val)
	}
}

func BenchmarkFloat32Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Float32("float32", float32Val)
	}
}

func BenchmarkFloat32KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("float32", float32Val)
	}
}

func BenchmarkFloat64Key(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Float64("float64", float64Val)
	}
}

func BenchmarkFloat64KeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("float64", float64Val)
	}
}

func BenchmarkStringKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.String("string", stringVal)
	}
}

func BenchmarkStringKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("string", stringVal)
	}
}

func BenchmarkBytesKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("bytes", bytesVal)
		// TODO struct json
	}
}

func BenchmarkStructKeyInfer(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = kv.Infer("struct", structVal)
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

func BenchmarkEmitInt32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = int32KeyVal.Value.Emit()
	}
}

func BenchmarkEmitInt64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = int64KeyVal.Value.Emit()
	}
}

func BenchmarkEmitUint(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uintKeyVal.Value.Emit()
	}
}

func BenchmarkEmitUint8(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uint8KeyVal.Value.Emit()
	}
}

func BenchmarkEmitUint16(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uint16KeyVal.Value.Emit()
	}
}

func BenchmarkEmitUint32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uint32KeyVal.Value.Emit()
	}
}

func BenchmarkEmitUint64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = uint64KeyVal.Value.Emit()
	}
}

func BenchmarkEmitFloat32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = float32KeyVal.Value.Emit()
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
