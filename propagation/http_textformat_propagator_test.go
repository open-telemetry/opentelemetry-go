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
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/propagation"
)

var (
	traceID = core.TraceID{High: 0x4bf92f3577b34da6, Low: 0xa3ce929d0e0e4736}
	spanID  = uint64(0x00f067aa0ba902b7)
)

func TestExtractTraceContextFromHTTPReq(t *testing.T) {
	propagator := propagation.TextFormatPropagator()
	tests := []struct {
		name   string
		header string
		wantSc core.SpanContext
	}{
		{
			name:   "future version",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
		},
		{
			name:   "zero trace ID and span ID",
			header: "00-00000000000000000000000000000000-0000000000000000-01",
			wantSc: core.INVALID_SPAN_CONTEXT,
		},
		{
			name:   "valid header",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
		},
		{
			name:   "missing options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7",
			wantSc: core.INVALID_SPAN_CONTEXT,
		},
		{
			name:   "empty options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-",
			wantSc: core.INVALID_SPAN_CONTEXT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("traceparent", tt.header)
			e := propagator.Extractor(req)

			gotSc, _ := e.Extract()
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestInjectTraceContextToHTTPReq(t *testing.T) {
	propagator := propagation.TextFormatPropagator()
	tests := []struct {
		name       string
		sc         core.SpanContext
		wantHeader string
	}{
		{
			name: "valid spancontext, sampled",
			sc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		},
		{
			name: "valid spancontext, not sampled",
			sc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
		},
		{
			name:       "invalid spancontext",
			sc:         core.INVALID_SPAN_CONTEXT,
			wantHeader: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			i := propagator.Injector(req)
			i.Inject(tt.sc, nil)

			gotHeader := req.Header.Get("traceparent")
			if diff := cmp.Diff(gotHeader, tt.wantHeader); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}
