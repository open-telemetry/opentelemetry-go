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
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
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

func spans(n int) []tracesdk.ReadOnlySpan {
	now := time.Now()
	s := make(tracetest.SpanStubs, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("span %d", i)
		s[i] = tracetest.SpanStub{
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
	return s.Snapshots()
}

func benchmarkExportSpans(b *testing.B, o EndpointOption, i int) {
	ctx := context.Background()
	s := spans(i)
	exp, err := New(o)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if err := exp.ExportSpans(ctx, s); err != nil {
			b.Error(err)
		}
	}
}

func benchmarkCollector(b *testing.B, i int) {
	benchmarkExportSpans(b, withTestCollectorEndpoint(), i)
}

func benchmarkAgent(b *testing.B, i int) {
	benchmarkExportSpans(b, WithAgentEndpoint(), i)
}

func BenchmarkCollectorExportSpans1(b *testing.B)     { benchmarkCollector(b, 1) }
func BenchmarkCollectorExportSpans10(b *testing.B)    { benchmarkCollector(b, 10) }
func BenchmarkCollectorExportSpans100(b *testing.B)   { benchmarkCollector(b, 100) }
func BenchmarkCollectorExportSpans1000(b *testing.B)  { benchmarkCollector(b, 1000) }
func BenchmarkCollectorExportSpans10000(b *testing.B) { benchmarkCollector(b, 10000) }
func BenchmarkAgentExportSpans1(b *testing.B)         { benchmarkAgent(b, 1) }
func BenchmarkAgentExportSpans10(b *testing.B)        { benchmarkAgent(b, 10) }
func BenchmarkAgentExportSpans100(b *testing.B)       { benchmarkAgent(b, 100) }

/*
* BUG: These tests are not possible currently because the thrift payload size
* does not fit in a UDP packet with the default size (65000) and will return
* an error.

func BenchmarkAgentExportSpans1000(b *testing.B)      { benchmarkAgent(b, 1000) }
func BenchmarkAgentExportSpans10000(b *testing.B)     { benchmarkAgent(b, 10000) }

*/
