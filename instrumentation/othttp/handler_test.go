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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/api/trace"
	mockmeter "go.opentelemetry.io/otel/internal/metric"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func assertMetricLabels(t *testing.T, expectedLabels []kv.KeyValue, measurementBatches []mockmeter.Batch) {
	for _, batch := range measurementBatches {
		assert.ElementsMatch(t, expectedLabels, batch.Labels)
	}
}

func TestHandlerBasics(t *testing.T) {
	rr := httptest.NewRecorder()

	var id uint64
	tracer := mocktrace.MockTracer{StartSpanID: &id}
	meterimpl, meter := mockmeter.NewMeter()

	operation := "test_handler"

	h := NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := io.WriteString(w, "hello world"); err != nil {
				t.Fatal(err)
			}
		}), operation,
		WithTracer(&tracer),
		WithMeter(meter),
	)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", strings.NewReader("foo"))
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(rr, r)

	if len(meterimpl.MeasurementBatches) == 0 {
		t.Fatalf("got 0 recorded measurements, expected 1 or more")
	}

	labelsToVerify := []kv.KeyValue{
		standard.HTTPServerNameKey.String(operation),
		standard.HTTPSchemeHTTP,
		standard.HTTPHostKey.String(r.Host),
		standard.HTTPFlavorKey.String(fmt.Sprintf("1.%d", r.ProtoMinor)),
		standard.HTTPRequestContentLengthKey.Int64(3),
	}

	assertMetricLabels(t, labelsToVerify, meterimpl.MeasurementBatches)

	if got, expected := rr.Result().StatusCode, http.StatusOK; got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	if got := rr.Header().Get("Traceparent"); got == "" {
		t.Fatal("expected non empty trace header")
	}
	if got, expected := id, uint64(1); got != expected {
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

func TestHandlerNoWrite(t *testing.T) {
	rr := httptest.NewRecorder()

	var id uint64
	tracer := mocktrace.MockTracer{StartSpanID: &id}

	operation := "test_handler"
	var span trace.Span

	h := NewHandler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span = trace.SpanFromContext(r.Context())
		}), operation,
		WithTracer(&tracer),
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
	if got, expected := id, uint64(1); got != expected {
		t.Fatalf("got %d, expected %d", got, expected)
	}
	if mockSpan, ok := span.(*mocktrace.MockSpan); ok {
		if got, expected := mockSpan.Status, codes.OK; got != expected {
			t.Fatalf("got %q, expected %q", got, expected)
		}
	} else {
		t.Fatalf("Expected *moctrace.MockSpan, got %T", span)
	}
}
