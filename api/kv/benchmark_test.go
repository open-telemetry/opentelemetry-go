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
	"context"
	"testing"

	"go.opentelemetry.io/otel/api/kv"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Note: The tests below load a real SDK to ensure the compiler isn't optimizing
// the test based on global analysis limited to the NoopSpan implementation.
func getSpan() trace.Span {
	_, sp := global.Tracer("Test").Start(context.Background(), "Span")
	tr, _ := sdktrace.NewProvider()
	global.SetTraceProvider(tr)
	return sp
}

func BenchmarkKeyInfer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kv.Infer("Attr", int(256))
	}
}

func BenchmarkMultiNoKeyInference(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttributes(kv.Int("Attr", 1))
	}
}

func BenchmarkMultiWithKeyInference(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttributes(kv.Infer("Attr", 1))
	}
}

func BenchmarkSingleWithKeyInferenceValue(b *testing.B) {
	sp := getSpan()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sp.SetAttribute("Attr", 1)
	}

	b.StopTimer()
}
