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

package propagation_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
	"go.opentelemetry.io/otel/propagation"
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
		propagator := propagation.HTTPB3Propagator{tg.singleHeader}
		for _, tt := range tg.tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				for h, v := range tt.headers {
					req.Header.Set(h, v)
				}

				ctx := context.Background()
				gotSc, _ := propagator.Extract(ctx, req.Header)
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
		propagator := propagation.HTTPB3Propagator{tg.singleHeader}
		for _, tt := range tg.tests {
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				ctx := context.Background()
				if tt.parentSc.IsValid() {
					ctx, _ = mockTracer.Start(ctx, "inject", trace.ChildOf(tt.parentSc))
				} else {
					ctx, _ = mockTracer.Start(ctx, "inject")
				}
				propagator.Inject(ctx, req.Header)

				for h, v := range tt.wantHeaders {
					got, want := req.Header.Get(h), v
					if diff := cmp.Diff(got, want); diff != "" {
						t.Errorf("%s: %s, header=%s: -got +want %s", tg.name, tt.name, h, diff)
					}
				}
				wantOk := false
				for _, h := range tt.doNotWantHeaders {
					v, gotOk := req.Header[h]
					if diff := cmp.Diff(gotOk, wantOk); diff != "" {
						t.Errorf("%s: %s, header=%s: -got +want %s, value=%s", tg.name, tt.name, h, diff, v)
					}

				}
			})
		}
	}
}

func TestHTTPB3Propagator_GetAllKeys(t *testing.T) {
	propagator := propagation.HTTPB3Propagator{false}
	want := []string{
		propagation.B3TraceIDHeader,
		propagation.B3SpanIDHeader,
		propagation.B3SampledHeader,
	}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}

func TestHTTPB3PropagatorWithSingleHeader_GetAllKeys(t *testing.T) {
	propagator := propagation.HTTPB3Propagator{true}
	want := []string{
		propagation.B3SingleHeader,
	}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
