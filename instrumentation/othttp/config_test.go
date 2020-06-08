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

package othttp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func TestBasicFilter(t *testing.T) {
	rr := httptest.NewRecorder()

	var id uint64
	tracer := mocktrace.MockTracer{StartSpanID: &id}

	h := NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := io.WriteString(w, "hello world"); err != nil {
				t.Fatal(err)
			}
		}), "test_handler",
		WithTracer(&tracer),
		WithFilter(func(r *http.Request) bool {
			return false
		}),
	)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rr, r)
	if got, expected := rr.Result().StatusCode, http.StatusOK; got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	if got := rr.Header().Get("Traceparent"); got != "" {
		t.Fatal("expected empty trace header")
	}
	if got, expected := id, uint64(0); got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	d, err := ioutil.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := string(d), "hello world"; got != expected {
		t.Fatalf("got %q, expected %q", got, expected)
	}
}

func TestSpanNameFormatter(t *testing.T) {
	var testCases = []struct {
		name      string
		formatter func(s string, r *http.Request) string
		operation string
		expected  string
	}{
		{
			name:      "default handler formatter",
			formatter: defaultHandlerFormatter,
			operation: "test_operation",
			expected:  "test_operation",
		},
		{
			name:      "default transport formatter",
			formatter: defaultTransportFormatter,
			expected:  http.MethodGet,
		},
		{
			name: "custom formatter",
			formatter: func(s string, r *http.Request) string {
				return r.URL.Path
			},
			operation: "",
			expected:  "/hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			var id uint64
			var spanName string
			tracer := mocktrace.MockTracer{
				StartSpanID: &id,
				OnSpanStarted: func(span *mocktrace.MockSpan) {
					spanName = span.Name
				},
			}
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := io.WriteString(w, "hello world"); err != nil {
					t.Fatal(err)
				}
			})
			h := NewHandler(
				handler,
				tc.operation,
				WithTracer(&tracer),
				WithSpanNameFormatter(tc.formatter),
			)
			r, err := http.NewRequest(http.MethodGet, "http://localhost/hello", nil)
			if err != nil {
				t.Fatal(err)
			}
			h.ServeHTTP(rr, r)
			if got, expected := rr.Result().StatusCode, http.StatusOK; got != expected {
				t.Fatalf("got %d, expected %d", got, expected)
			}
			if got, expected := spanName, tc.expected; got != expected {
				t.Fatalf("got %q, expected %q", got, expected)
			}
		})
	}
}
