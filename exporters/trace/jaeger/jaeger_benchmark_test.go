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

package jaeger

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/instrumentation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	traceID     trace.TraceID
	spanID      trace.SpanID
	spanContext trace.SpanContext

	instrLibName = "benchmark.tests"
)

func init() {
	var err error
	traceID, err = trace.TraceIDFromHex("0102030405060708090a0b0c0d0e0f10")
	if err != nil {
		panic(err)
	}
	spanID, err = trace.SpanIDFromHex("0102030405060708")
	if err != nil {
		panic(err)
	}
	spanContext = trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
	})
}

func spans(n int) []*tracesdk.SpanSnapshot {
	now := time.Now()
	s := make([]*tracesdk.SpanSnapshot, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("span %d", i)
		s[i] = &tracesdk.SpanSnapshot{
			SpanContext: spanContext,
			Name:        name,
			StartTime:   now,
			EndTime:     now,
			SpanKind:    trace.SpanKindClient,
			InstrumentationLibrary: instrumentation.Library{
				Name: instrLibName,
			},
		}
	}
	return s
}

func benchmarkExportSpans(b *testing.B, i int) {
	ctx := context.Background()
	s := spans(i)
	exp, err := NewRawExporter(
		withTestCollectorEndpoint(),
		WithBatchMaxCount(i+1),
		WithBufferMaxCount(i+1),
	)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if err := exp.ExportSpans(ctx, s); err != nil {
			b.Error(err)
		}
		exp.bundler.Flush()
	}
}

func BenchmarkExportSpans1(b *testing.B)     { benchmarkExportSpans(b, 1) }
func BenchmarkExportSpans10(b *testing.B)    { benchmarkExportSpans(b, 10) }
func BenchmarkExportSpans100(b *testing.B)   { benchmarkExportSpans(b, 100) }
func BenchmarkExportSpans1000(b *testing.B)  { benchmarkExportSpans(b, 1000) }
func BenchmarkExportSpans10000(b *testing.B) { benchmarkExportSpans(b, 10000) }
