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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func TestTransportBasics(t *testing.T) {
	var id uint64
	tracer := mocktrace.MockTracer{StartSpanID: &id}
	content := []byte("Hello, world!")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := propagation.ExtractHTTP(r.Context(), global.Propagators(), r.Header)
		span := trace.RemoteSpanContextFromContext(ctx)
		tgtID, err := trace.SpanIDFromHex(fmt.Sprintf("%016x", id))
		if err != nil {
			t.Fatalf("Error converting id to SpanID: %s", err.Error())
		}
		if span.SpanID != tgtID {
			t.Fatalf("testing remote SpanID: got %s, expected %s", span.SpanID, tgtID)
		}
		if _, err := w.Write(content); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	r, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	tr := NewTransport(
		http.DefaultTransport,
		WithTracer(&tracer),
	)

	c := http.Client{Transport: tr}
	res, err := c.Do(r)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(body, content) {
		t.Fatalf("unexpected content: got %s, expected %s", body, content)
	}
}
