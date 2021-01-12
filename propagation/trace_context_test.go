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

package propagation_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestExtractValidTraceContextFromHTTPReq(t *testing.T) {
	prop := propagation.TraceContext{}
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
			ctx = prop.Extract(ctx, req.Header)
			gotSc := trace.RemoteSpanContextFromContext(ctx)
			if diff := cmp.Diff(gotSc, tt.wantSc, cmp.AllowUnexported(trace.TraceState{})); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractInvalidTraceContextFromHTTPReq(t *testing.T) {
	wantSc := trace.SpanContext{}
	prop := propagation.TraceContext{}
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
			ctx = prop.Extract(ctx, req.Header)
			gotSc := trace.RemoteSpanContextFromContext(ctx)
			if diff := cmp.Diff(gotSc, wantSc, cmp.AllowUnexported(trace.TraceState{})); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestInjectTraceContextToHTTPReq(t *testing.T) {
	mockTracer := oteltest.DefaultTracer()
	prop := propagation.TraceContext{}
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
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000002-01",
		},
		{
			name: "valid spancontext, not sampled",
			sc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000003-00",
		},
		{
			name: "valid spancontext, with unsupported bit set in traceflags",
			sc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: 0xff,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000004-01",
		},
		{
			name:       "invalid spancontext",
			sc:         trace.SpanContext{},
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
			prop.Inject(ctx, req.Header)

			gotHeader := req.Header.Get("traceparent")
			if diff := cmp.Diff(gotHeader, tt.wantHeader); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestTraceContextPropagator_GetAllKeys(t *testing.T) {
	var propagator propagation.TraceContext
	want := []string{"traceparent", "tracestate"}
	got := propagator.Fields()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}

func TestTraceStatePropagation(t *testing.T) {
	prop := propagation.TraceContext{}
	stateHeader := "tracestate"
	parentHeader := "traceparent"
	state, err := trace.TraceStateFromKeyValues(label.String("key1", "value1"), label.String("key2", "value2"))
	if err != nil {
		t.Fatalf("Unable to construct expected TraceState: %s", err.Error())
	}

	tests := []struct {
		name    string
		headers map[string]string
		valid   bool
		wantSc  trace.SpanContext
	}{
		{
			name: "valid parent and state",
			headers: map[string]string{
				parentHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
				stateHeader:  "key1=value1,key2=value2",
			},
			valid: true,
			wantSc: trace.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceState: state,
			},
		},
		{
			name: "valid parent, invalid state",
			headers: map[string]string{
				parentHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
				stateHeader:  "key1=value1,invalid$@#=invalid",
			},
			valid: false,
			wantSc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name: "valid parent, malformed state",
			headers: map[string]string{
				parentHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
				stateHeader:  "key1=value1,invalid",
			},
			valid: false,
			wantSc: trace.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inReq, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
			for hk, hv := range tt.headers {
				inReq.Header.Add(hk, hv)
			}

			ctx := prop.Extract(context.Background(), inReq.Header)
			if diff := cmp.Diff(
				trace.RemoteSpanContextFromContext(ctx),
				tt.wantSc,
				cmp.AllowUnexported(label.Value{}),
				cmp.AllowUnexported(trace.TraceState{}),
			); diff != "" {
				t.Errorf("Extracted tracestate: -got +want %s", diff)
			}

			if tt.valid {
				mockTracer := oteltest.DefaultTracer()
				ctx, _ = mockTracer.Start(ctx, "inject")
				outReq, _ := http.NewRequest(http.MethodGet, "http://www.example.com", nil)
				prop.Inject(ctx, outReq.Header)

				if diff := cmp.Diff(outReq.Header.Get(stateHeader), tt.headers[stateHeader]); diff != "" {
					t.Errorf("Propagated tracestate: -got +want %s", diff)
				}
			}
		})
	}
}
