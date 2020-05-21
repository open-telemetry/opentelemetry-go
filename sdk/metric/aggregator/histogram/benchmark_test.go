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

package histogram_test

import (
	"context"
	"math/rand"
	"testing"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/test"
)

const inputRange = 1e6

func benchmarkHistogramSearchFloat64(b *testing.B, size int) {
	boundaries := make([]float64, size)

	for i := range boundaries {
		boundaries[i] = rand.Float64() * inputRange
	}

	values := make([]float64, b.N)
	for i := range values {
		values[i] = rand.Float64() * inputRange
	}
	desc := test.NewAggregatorTest(metric.ValueRecorderKind, metric.Float64NumberKind)
	agg := histogram.New(desc, boundaries)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = agg.Update(ctx, metric.NewFloat64Number(values[i]), desc)
	}
}

func BenchmarkHistogramSearchFloat64_1(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 1)
}
func BenchmarkHistogramSearchFloat64_8(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 8)
}
func BenchmarkHistogramSearchFloat64_16(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 16)
}
func BenchmarkHistogramSearchFloat64_32(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 32)
}
func BenchmarkHistogramSearchFloat64_64(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 64)
}
func BenchmarkHistogramSearchFloat64_128(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 128)
}
func BenchmarkHistogramSearchFloat64_256(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 256)
}
func BenchmarkHistogramSearchFloat64_512(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 512)
}
func BenchmarkHistogramSearchFloat64_1024(b *testing.B) {
	benchmarkHistogramSearchFloat64(b, 1024)
}

func benchmarkHistogramSearchInt64(b *testing.B, size int) {
	boundaries := make([]float64, size)

	for i := range boundaries {
		boundaries[i] = rand.Float64() * inputRange
	}

	values := make([]int64, b.N)
	for i := range values {
		values[i] = int64(rand.Float64() * inputRange)
	}
	desc := test.NewAggregatorTest(metric.ValueRecorderKind, metric.Int64NumberKind)
	agg := histogram.New(desc, boundaries)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = agg.Update(ctx, metric.NewInt64Number(values[i]), desc)
	}
}

func BenchmarkHistogramSearchInt64_1(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 1)
}
func BenchmarkHistogramSearchInt64_8(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 8)
}
func BenchmarkHistogramSearchInt64_16(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 16)
}
func BenchmarkHistogramSearchInt64_32(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 32)
}
func BenchmarkHistogramSearchInt64_64(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 64)
}
func BenchmarkHistogramSearchInt64_128(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 128)
}
func BenchmarkHistogramSearchInt64_256(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 256)
}
func BenchmarkHistogramSearchInt64_512(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 512)
}
func BenchmarkHistogramSearchInt64_1024(b *testing.B) {
	benchmarkHistogramSearchInt64(b, 1024)
}
