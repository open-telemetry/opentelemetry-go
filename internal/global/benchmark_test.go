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

package global_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/label"
	metricglobal "go.opentelemetry.io/otel/metric/global"
)

func BenchmarkGlobalInt64CounterAddNoSDK(b *testing.B) {
	// Compare with BenchmarkGlobalInt64CounterAddWithSDK() in
	// ../../sdk/metric/benchmark_test.go to see the overhead of the
	// global no-op system against a registered SDK.
	global.ResetForTest()
	ctx := context.Background()
	sdk := metricglobal.Meter("test")
	labs := []label.KeyValue{label.String("A", "B")}
	cnt := Must(sdk).NewInt64Counter("int64.counter")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cnt.Add(ctx, 1, labs...)
	}
}

func BenchmarkStartEndSpanNoSDK(b *testing.B) {
	// Compare with BenchmarkStartEndSpan() in
	// ../../sdk/trace/benchmark_test.go.
	global.ResetForTest()
	t := otel.Tracer("Benchmark StartEndSpan")
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := t.Start(ctx, "/foo")
		span.End()
	}
}
