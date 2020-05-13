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
package httptrace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	nhtrace "net/http/httptrace"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/plugin/httptrace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
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

func TestHTTPRequestWithClientTrace(t *testing.T) {
	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	tr := tp.Tracer("httptrace/client")

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

	getSpan := func(name string) *export.SpanData {
		spans, ok := exp.spanMap[name]
		if !ok {
			t.Fatalf("no spans found with the name %s, %v", name, exp.spanMap)
		}

		if len(spans) != 1 {
			t.Fatalf("Expected exactly one span for %s but found %d", name, len(spans))
		}

		return spans[0]
	}

	testLen := []struct {
		name       string
		attributes []kv.KeyValue
		parent     string
	}{
		{
			name:       "http.connect",
			attributes: []kv.KeyValue{kv.String("http.remote", address.String())},
			parent:     "http.getconn",
		},
		{
			name: "http.getconn",
			attributes: []kv.KeyValue{
				kv.String("http.remote", address.String()),
				kv.String("http.host", address.String()),
			},
			parent: "test",
		},
		{
			name:   "http.receive",
			parent: "test",
		},
		{
			name:   "http.send",
			parent: "test",
		},
		{
			name: "test",
		},
	}
	for _, tl := range testLen {
		span := getSpan(tl.name)

		if tl.parent != "" {
			parentSpan := getSpan(tl.parent)

			if span.ParentSpanID != parentSpan.SpanContext.SpanID {
				t.Fatalf("[span %s] does not have expected parent span %s", tl.name, tl.parent)
			}
		}

		actualAttrs := make(map[kv.Key]string)
		for _, attr := range span.Attributes {
			actualAttrs[attr.Key] = attr.Value.Emit()
		}

		expectedAttrs := make(map[kv.Key]string)
		for _, attr := range tl.attributes {
			expectedAttrs[attr.Key] = attr.Value.Emit()
		}

		if tl.name == "http.getconn" {
			local := kv.Key("http.local")
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

func TestConcurrentConnectionStart(t *testing.T) {
	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	ct := httptrace.NewClientTrace(context.Background())

	tts := []struct {
		name string
		run  func()
	}{
		{
			name: "Open1Close1Open2Close2",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectDone("tcp", "[::1]:3000", nil)
			},
		},
		{
			name: "Open2Close2Open1Close1",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectDone("tcp", "[::1]:3000", nil)
				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
			},
		},
		{
			name: "Open1Open2Close1Close2",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
				ct.ConnectDone("tcp", "[::1]:3000", nil)
			},
		},
		{
			name: "Open1Open2Close2Close1",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectDone("tcp", "[::1]:3000", nil)
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
			},
		},
		{
			name: "Open2Open1Close1Close2",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
				ct.ConnectDone("tcp", "[::1]:3000", nil)
			},
		},
		{
			name: "Open2Open1Close2Close1",
			run: func() {
				exp.spanMap = make(map[string][]*export.SpanData)

				ct.ConnectStart("tcp", "[::1]:3000")
				ct.ConnectStart("tcp", "127.0.0.1:3000")
				ct.ConnectDone("tcp", "[::1]:3000", nil)
				ct.ConnectDone("tcp", "127.0.0.1:3000", nil)
			},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			tt.run()
			spans := exp.spanMap["http.connect"]

			if l := len(spans); l != 2 {
				t.Fatalf("Expected 2 'http.connect' traces but found %d", l)
			}

			remotes := make(map[string]struct{})
			for _, span := range spans {
				if l := len(span.Attributes); l != 1 {
					t.Fatalf("Expected 1 attribute on each span but found %d", l)
				}

				attr := span.Attributes[0]
				if attr.Key != "http.remote" {
					t.Fatalf("Expected attribute to be 'http.remote' but found %s", attr.Key)
				}
				remotes[attr.Value.Emit()] = struct{}{}
			}

			if l := len(remotes); l != 2 {
				t.Fatalf("Expected 2 different 'http.remote' but found %d", l)
			}

			for _, remote := range []string{"127.0.0.1:3000", "[::1]:3000"} {
				if _, ok := remotes[remote]; !ok {
					t.Fatalf("Missing remote %s", remote)
				}
			}
		})
	}
}

func TestEndBeforeStartCreatesSpan(t *testing.T) {
	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	tr := tp.Tracer("httptrace/client")
	ctx, span := tr.Start(context.Background(), "test")
	defer span.End()

	ct := httptrace.NewClientTrace(ctx)

	ct.DNSDone(nhtrace.DNSDoneInfo{})
	ct.DNSStart(nhtrace.DNSStartInfo{Host: "example.com"})

	getSpan := func(name string) *export.SpanData {
		spans, ok := exp.spanMap[name]
		if !ok {
			t.Fatalf("no spans found with the name %s, %v", name, exp.spanMap)
		}

		if len(spans) != 1 {
			t.Fatalf("Expected exactly one span for %s but found %d", name, len(spans))
		}

		return spans[0]
	}

	getSpan("http.dns")

}
