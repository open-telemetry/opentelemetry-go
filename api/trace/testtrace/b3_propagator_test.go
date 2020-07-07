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
)

func TestExtractB3(t *testing.T) {
	testGroup := []struct {
		name  string
		tests []extractTest
	}{
		{
			name:  "valid headers",
			tests: extractHeaders,
		},
		{
			name:  "invalid headers",
			tests: extractInvalidHeaders,
		},
	}

	for _, tg := range testGroup {
		propagator := trace.B3{}
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

type testSpan struct {
	trace.NoopSpan
	sc trace.SpanContext
}

func (s testSpan) SpanContext() trace.SpanContext {
	return s.sc
}

func TestInjectB3(t *testing.T) {
	testGroup := []struct {
		name  string
		tests []injectTest
	}{
		{
			name:  "valid headers",
			tests: injectHeader,
		},
		{
			name:  "invalid headers",
			tests: injectInvalidHeader,
		},
	}

	for _, tg := range testGroup {
		for _, tt := range tg.tests {
			propagator := trace.B3{InjectEncoding: tt.encoding}
			t.Run(tt.name, func(t *testing.T) {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				ctx := trace.ContextWithSpan(
					context.Background(),
					testSpan{sc: tt.sc},
				)
				propagator.Inject(ctx, req.Header)

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
	tests := []struct {
		name       string
		propagator trace.B3
		want       []string
	}{
		{
			name:       "no encoding specified",
			propagator: trace.B3{},
			want: []string{
				b3TraceID,
				b3SpanID,
				b3Sampled,
				b3Flags,
			},
		},
		{
			name:       "MultipleHeader encoding specified",
			propagator: trace.B3{InjectEncoding: trace.MultipleHeader},
			want: []string{
				b3TraceID,
				b3SpanID,
				b3Sampled,
				b3Flags,
			},
		},
		{
			name:       "SingleHeader encoding specified",
			propagator: trace.B3{InjectEncoding: trace.SingleHeader},
			want: []string{
				b3Context,
			},
		},
		{
			name:       "SingleHeader and MultipleHeader encoding specified",
			propagator: trace.B3{InjectEncoding: trace.SingleHeader | trace.MultipleHeader},
			want: []string{
				b3Context,
				b3TraceID,
				b3SpanID,
				b3Sampled,
				b3Flags,
			},
		},
	}

	for _, test := range tests {
		if diff := cmp.Diff(test.propagator.GetAllKeys(), test.want); diff != "" {
			t.Errorf("%s: GetAllKeys: -got +want %s", test.name, diff)
		}
	}
}
