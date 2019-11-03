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

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
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
	address := ts.Listener.Addr()

	client := ts.Client()
	err := tr.WithSpan(context.Background(), "test",
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", ts.URL, nil)
			_, req = httptrace.W3C(ctx, req)

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err.Error())
			}
			_ = res.Body.Close()

			return nil
		})
	if err != nil {
		panic("unexpected error in http request: " + err.Error())
	}

	testLen := []struct {
		name       string
		attributes []core.KeyValue
	}{
		{
			name:       "go.opentelemetry.io/otel/plugin/httptrace/http.connect",
			attributes: []core.KeyValue{key.String("http.remote", address.String())},
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.getconn",
			attributes: []core.KeyValue{
				key.String("http.remote", address.String()),
				key.String("http.host", address.String()),
			},
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.receive",
		},
		{
			name: "go.opentelemetry.io/otel/plugin/httptrace/http.send",
		},
		{
			name: "httptrace/client/test",
		},
	}
	for _, tl := range testLen {
		spans, ok := exp.spanMap[tl.name]
		if !ok {
			t.Fatalf("no spans found with the name %s, %v", tl.name, exp.spanMap)
		}

		if len(spans) != 1 {
			t.Fatalf("Expected exactly one span for %s but found %d", tl.name, len(spans))
		}
		span := spans[0]

		actualAttrs := make(map[core.Key]string)
		for _, attr := range span.Attributes {
			actualAttrs[attr.Key] = attr.Value.Emit()
		}

		expectedAttrs := make(map[core.Key]string)
		for _, attr := range tl.attributes {
			expectedAttrs[attr.Key] = attr.Value.Emit()
		}

		if tl.name == "go.opentelemetry.io/otel/plugin/httptrace/http.getconn" {
			local := key.New("http.local")
			// http.local attribute is not deterministic, just make sure it exists for `getconn`.
			if _, ok := actualAttrs[local]; ok {
				delete(actualAttrs, local)
			} else {
				t.Fatalf("[span %s] is missing attribute %v", tl.name, local)
			}
		}

		if diff := cmp.Diff(actualAttrs, expectedAttrs); diff != "" {
			t.Fatalf("[span %s] Attributes are different: %v", tl.name, diff)
		}
	}
}
