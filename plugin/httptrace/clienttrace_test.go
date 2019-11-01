// Copyright 2019, OpenTelemetry Authors
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
package httptrace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"go.opentelemetry.io/otel/sdk/export"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type testExporter struct {
	mu      sync.Mutex
	spanMap map[string][]*export.SpanData
}

func (t *testExporter) ExportSpan(ctx context.Context, s *export.SpanData) {
	t.mu.Lock()
	defer t.mu.Unlock()
	var spans []*export.SpanData
	var ok bool

	if spans, ok = t.spanMap[s.Name]; !ok {
		spans = []*export.SpanData{}
		t.spanMap[s.Name] = spans
	}
	spans = append(spans, s)
	t.spanMap[s.Name] = spans
}

var _ export.SpanSyncer = (*testExporter)(nil)

func TestClientTrace(t *testing.T) {
	var wg sync.WaitGroup

	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	tr := tp.GetTracer("httptrace/client")

	// Mock http server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}),
	)
	defer ts.Close()

	client := ts.Client()
	iterations := 50
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			err := tr.WithSpan(context.Background(), "test",
				func(ctx context.Context) error {
					req, _ := http.NewRequest("GET", ts.URL, nil)

					_, req = httptrace.W3C(ctx, req)

					res, err := client.Do(req)
					if err != nil {
						return err
					}
					res.Body.Close()

					return nil
				})

			if err != nil {
				panic("unexpected error in http request")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	exp.mu.Lock()
	testLen := []struct {
		name             string
		len              int
		ignoreExactMatch bool
	}{
		{
			name:             "go.opentelemetry.io/otel/plugin/httptrace/http.connect",
			len:              iterations,
			ignoreExactMatch: true,
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.getconn",
			len:  iterations,
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.receive",
			len:  iterations,
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.send",
			len:  iterations,
		},
		{
			name: "httptrace/client/test",
			len:  iterations,
		},
	}
	for _, tl := range testLen {
		want := tl.len
		spans, ok := exp.spanMap[tl.name]
		if !ok {
			t.Fatalf("no spans found with the name %s, %v", tl.name, exp.spanMap)
		}
		got := len(spans)
		if !tl.ignoreExactMatch {
			if got != want {
				t.Fatalf("got %d, want %d spans", got, want)
			}
		}
	}
	exp.mu.Unlock()
}
