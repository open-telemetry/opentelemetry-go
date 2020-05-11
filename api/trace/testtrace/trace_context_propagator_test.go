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

var (
	traceID = mustTraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID  = mustSpanIDFromHex("00f067aa0ba902b7")
)

func mustTraceIDFromHex(s string) (t trace.ID) {
	t, _ = trace.IDFromHex(s)
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	t, _ = trace.SpanIDFromHex(s)
	return
}

func TestExtractValidTraceContextFromHTTPReq(t *testing.T) {
	props := propagation.New(propagation.WithExtractors(trace.TraceContext{}))
	tests := []struct {
		name   string
		header string
		wantSc trace.SpanContext
	}{
		{
			name:   "valid w3cHeader",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			wantSc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name:   "valid w3cHeader and sampled",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			name:   "future version",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			name:   "future options with sampled bit set",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			name:   "future options with sampled bit cleared",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-08",
			wantSc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name:   "future additional data",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09-XYZxsf09",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			name:   "valid b3Header ending in dash",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01-",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
		{
			name:   "future valid b3Header ending in dash",
			header: "01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09-",
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("traceparent", tt.header)

			ctx := context.Background()
			ctx = propagation.ExtractHTTP(ctx, props, req.Header)
			gotSc := trace.RemoteSpanContextFromContext(ctx)
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractInvalidTraceContextFromHTTPReq(t *testing.T) {
	wantSc := trace.EmptySpanContext()
	props := propagation.New(propagation.WithExtractors(trace.TraceContext{}))
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "wrong version length",
			header: "0000-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "wrong trace ID length",
			header: "00-ab00000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "wrong span ID length",
			header: "00-ab000000000000000000000000000000-cd0000000000000000-01",
		},
		{
			name:   "wrong trace flag length",
			header: "00-ab000000000000000000000000000000-cd00000000000000-0100",
		},
		{
			name:   "bogus version",
			header: "qw-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "bogus trace ID",
			header: "00-qw000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "bogus span ID",
			header: "00-ab000000000000000000000000000000-qw00000000000000-01",
		},
		{
			name:   "bogus trace flag",
			header: "00-ab000000000000000000000000000000-cd00000000000000-qw",
		},
		{
			name:   "upper case version",
			header: "A0-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "upper case trace ID",
			header: "00-AB000000000000000000000000000000-cd00000000000000-01",
		},
		{
			name:   "upper case span ID",
			header: "00-ab000000000000000000000000000000-CD00000000000000-01",
		},
		{
			name:   "upper case trace flag",
			header: "00-ab000000000000000000000000000000-cd00000000000000-A1",
		},
		{
			name:   "zero trace ID and span ID",
			header: "00-00000000000000000000000000000000-0000000000000000-01",
		},
		{
			name:   "trace-flag unused bits set",
			header: "00-ab000000000000000000000000000000-cd00000000000000-09",
		},
		{
			name:   "missing options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7",
		},
		{
			name:   "empty options",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("traceparent", tt.header)

			ctx := context.Background()
			ctx = propagation.ExtractHTTP(ctx, props, req.Header)
			gotSc := trace.RemoteSpanContextFromContext(ctx)
			if diff := cmp.Diff(gotSc, wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestInjectTraceContextToHTTPReq(t *testing.T) {
	var id uint64
	mockTracer := &mocktrace.MockTracer{
		Sampled:     false,
		StartSpanID: &id,
	}
	props := propagation.New(propagation.WithInjectors(trace.TraceContext{}))
	tests := []struct {
		name       string
		sc         trace.SpanContext
		wantHeader string
	}{
		{
			name: "valid spancontext, sampled",
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: trace.FlagsSampled,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000001-01",
		},
		{
			name: "valid spancontext, not sampled",
			sc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000002-00",
		},
		{
			name: "valid spancontext, with unsupported bit set in traceflags",
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: 0xff,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000003-01",
		},
		{
			name:       "invalid spancontext",
			sc:         trace.EmptySpanContext(),
			wantHeader: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			ctx := context.Background()
			if tt.sc.IsValid() {
				ctx = trace.ContextWithRemoteSpanContext(ctx, tt.sc)
				ctx, _ = mockTracer.Start(ctx, "inject")
			}
			propagation.InjectHTTP(ctx, props, req.Header)

			gotHeader := req.Header.Get("traceparent")
			if diff := cmp.Diff(gotHeader, tt.wantHeader); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestTraceContextPropagator_GetAllKeys(t *testing.T) {
	var propagator trace.TraceContext
	want := []string{"traceparent", "tracestate"}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}

func TestTraceStatePropagation(t *testing.T) {
	props := propagation.New(propagation.WithInjectors(trace.TraceContext{}), propagation.WithExtractors(trace.TraceContext{}))
	want := "opaquevalue"
	headerName := "tracestate"

	inReq, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	inReq.Header.Add(headerName, want)
	ctx := propagation.ExtractHTTP(context.Background(), props, inReq.Header)

	outReq, _ := http.NewRequest(http.MethodGet, "http://www.example.com", nil)
	propagation.InjectHTTP(ctx, props, outReq.Header)

	if diff := cmp.Diff(outReq.Header.Get(headerName), want); diff != "" {
		t.Errorf("Propagate tracestate: -got +want %s", diff)
	}
}
