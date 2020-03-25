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

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
)

func TestExtractB3(t *testing.T) {
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
		props := propagation.New(propagation.WithExtractors(propagator))

		for _, tt := range tg.tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				for h, v := range tt.headers {
					req.Header.Set(h, v)
				}

				ctx := context.Background()
				ctx = propagation.ExtractHTTP(ctx, props, req.Header)
				gotSc := trace.RemoteSpanContextFromContext(ctx)
				if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
					t.Errorf("%s: %s: -got +want %s", tg.name, tt.name, diff)
				}
			})
		}
	}
}

func TestInjectB3(t *testing.T) {
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
		props := propagation.New(propagation.WithInjectors(propagator))
		for _, tt := range tg.tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				ctx := context.Background()
				if tt.parentSc.IsValid() {
					ctx = trace.ContextWithRemoteSpanContext(ctx, tt.parentSc)
				}
				ctx, _ = mockTracer.Start(ctx, "inject")
				propagation.InjectHTTP(ctx, props, req.Header)

				for h, v := range tt.wantHeaders {
					got, want := req.Header.Get(h), v
					if diff := cmp.Diff(got, want); diff != "" {
						t.Errorf("%s: %s, header=%s: -got +want %s", tg.name, tt.name, h, diff)
					}
				}
				for _, h := range tt.doNotWantHeaders {
					v, gotOk := req.Header[h]
					if diff := cmp.Diff(gotOk, false); diff != "" {
						t.Errorf("%s: %s, header=%s: -got +want %s, value=%s", tg.name, tt.name, h, diff, v)
					}
				}
			})
		}
	}
}

func TestB3Propagator_GetAllKeys(t *testing.T) {
	propagator := trace.B3{SingleHeader: false}
	want := []string{
		trace.B3TraceIDHeader,
		trace.B3SpanIDHeader,
		trace.B3SampledHeader,
	}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}

func TestB3PropagatorWithSingleHeader_GetAllKeys(t *testing.T) {
	propagator := trace.B3{SingleHeader: true}
	want := []string{
		trace.B3SingleHeader,
	}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
