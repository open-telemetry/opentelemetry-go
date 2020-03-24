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

package testtrace_test

import (
	"context"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func BenchmarkExtractB3(b *testing.B) {
	testGroup := []struct {
		singleHeader bool
		name         string
		tests        []extractTest
	}{
		{
			singleHeader: false,
			name:         "multiple headers",
			tests:        extractMultipleHeaders,
		},
		{
			singleHeader: true,
			name:         "single headers",
			tests:        extractSingleHeader,
		},
		{
			singleHeader: false,
			name:         "invalid multiple headers",
			tests:        extractInvalidB3MultipleHeaders,
		},
		{
			singleHeader: true,
			name:         "invalid single headers",
			tests:        extractInvalidB3SingleHeader,
		},
	}

	for _, tg := range testGroup {
		propagator := trace.B3{SingleHeader: tg.singleHeader}
		for _, tt := range tg.tests {
			traceBenchmark(tg.name+"/"+tt.name, b, func(b *testing.B) {
				ctx := context.Background()
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				for h, v := range tt.headers {
					req.Header.Set(h, v)
				}
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = propagator.Extract(ctx, req.Header)
				}
			})
		}
	}
}

func BenchmarkInjectB3(b *testing.B) {
	var id uint64
	testGroup := []struct {
		singleHeader bool
		name         string
		tests        []injectTest
	}{
		{
			singleHeader: false,
			name:         "multiple headers",
			tests:        injectB3MultipleHeader,
		},
		{
			singleHeader: true,
			name:         "single headers",
			tests:        injectB3SingleleHeader,
		},
	}

	mockTracer := &mocktrace.MockTracer{
		Sampled:     false,
		StartSpanID: &id,
	}

	for _, tg := range testGroup {
		id = 0
		propagator := trace.B3{SingleHeader: tg.singleHeader}
		for _, tt := range tg.tests {
			traceBenchmark(tg.name+"/"+tt.name, b, func(b *testing.B) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				ctx := context.Background()
				if tt.parentSc.IsValid() {
					ctx = trace.ContextWithRemoteSpanContext(ctx, tt.parentSc)
				}
				ctx, _ = mockTracer.Start(ctx, "inject")
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					propagator.Inject(ctx, req.Header)
				}
			})
		}
	}
}

func traceBenchmark(name string, b *testing.B, fn func(*testing.B)) {
	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
	b.Run(name, func(b *testing.B) {
		b.ReportAllocs()
		fn(b)
	})
}
