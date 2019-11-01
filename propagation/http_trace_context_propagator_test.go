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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	dctx "go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	mocktrace "go.opentelemetry.io/otel/internal/trace"
	"go.opentelemetry.io/otel/propagation"
)

var (
	traceID = mustTraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID  = mustSpanIDFromHex("00f067aa0ba902b7")
)

func mustTraceIDFromHex(s string) (t core.TraceID) {
	t, _ = core.TraceIDFromHex(s)
	return
}

func mustSpanIDFromHex(s string) (t core.SpanID) {
	t, _ = core.SpanIDFromHex(s)
	return
}

func TestExtractValidTraceContextFromHTTPReq(t *testing.T) {
	var propagator propagation.HTTPTraceContextPropagator
	tests := []struct {
		name   string
		header string
		wantSc core.SpanContext
	}{
		{
			name:   "valid w3cHeader",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			wantSc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name:   "valid w3cHeader and sampled",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name:   "future version",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name:   "future options with sampled bit set",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name:   "future options with sampled bit cleared",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-08",
			wantSc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name:   "future additional data",
			header: "02-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09-XYZxsf09",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name:   "valid b3Header ending in dash",
			header: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01-",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
		{
			name:   "future valid b3Header ending in dash",
			header: "01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-09-",
			wantSc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("traceparent", tt.header)

			ctx := context.Background()
			gotSc, _ := propagator.Extract(ctx, req.Header)
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractInvalidTraceContextFromHTTPReq(t *testing.T) {
	var propagator propagation.HTTPTraceContextPropagator
	wantSc := core.EmptySpanContext()
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
			gotSc, _ := propagator.Extract(ctx, req.Header)
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
	var propagator propagation.HTTPTraceContextPropagator
	tests := []struct {
		name       string
		sc         core.SpanContext
		wantHeader string
	}{
		{
			name: "valid spancontext, sampled",
			sc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: core.TraceFlagsSampled,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000001-01",
		},
		{
			name: "valid spancontext, not sampled",
			sc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000002-00",
		},
		{
			name: "valid spancontext, with unsupported bit set in traceflags",
			sc: core.SpanContext{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: 0xff,
			},
			wantHeader: "00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000003-01",
		},
		{
			name:       "invalid spancontext",
			sc:         core.EmptySpanContext(),
			wantHeader: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			ctx := context.Background()
			if tt.sc.IsValid() {
				ctx, _ = mockTracer.Start(ctx, "inject", trace.ChildOf(tt.sc))
			}
			propagator.Inject(ctx, req.Header)

			gotHeader := req.Header.Get("traceparent")
			if diff := cmp.Diff(gotHeader, tt.wantHeader); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractValidDistributedContextFromHTTPReq(t *testing.T) {
	propagator := propagation.HTTPTraceContextPropagator{}
	tests := []struct {
		name    string
		header  string
		wantKVs []core.KeyValue
	}{
		{
			name:   "valid w3cHeader",
			header: "key1=val1,key2=val2",
			wantKVs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2"),
			},
		},
		{
			name:   "valid w3cHeader with spaces",
			header: "key1 =   val1,  key2 =val2   ",
			wantKVs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2"),
			},
		},
		{
			name:   "valid w3cHeader with properties",
			header: "key1=val1,key2=val2;prop=1",
			wantKVs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2;prop=1"),
			},
		},
		{
			name:   "valid header with url-escaped comma",
			header: "key1=val1,key2=val2%2Cval3",
			wantKVs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2,val3"),
			},
		},
		{
			name:   "valid header with an invalid header",
			header: "key1=val1,key2=val2,a,val3",
			wantKVs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2"),
			},
		},
		{
			name:   "valid header with no value",
			header: "key1=,key2=val2",
			wantKVs: []core.KeyValue{
				key.New("key1").String(""),
				key.New("key2").String("val2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("Correlation-Context", tt.header)

			ctx := context.Background()
			_, gotCorCtx := propagator.Extract(ctx, req.Header)
			wantCorCtx := dctx.NewMap(dctx.MapUpdate{MultiKV: tt.wantKVs})
			if gotCorCtx.Len() != wantCorCtx.Len() {
				t.Errorf(
					"Got and Want CorCtx are not the same size %d != %d",
					gotCorCtx.Len(),
					wantCorCtx.Len(),
				)
			}
			totalDiff := ""
			wantCorCtx.Foreach(func(kv core.KeyValue) bool {
				val, _ := gotCorCtx.Value(kv.Key)
				diff := cmp.Diff(kv, core.KeyValue{Key: kv.Key, Value: val}, cmp.AllowUnexported(core.Value{}))
				if diff != "" {
					totalDiff += diff + "\n"
				}
				return true
			})
			if totalDiff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, totalDiff)
			}
		})
	}
}

func TestExtractInvalidDistributedContextFromHTTPReq(t *testing.T) {
	propagator := propagation.HTTPTraceContextPropagator{}
	tests := []struct {
		name   string
		header string
	}{
		{
			name:   "no key values",
			header: "header1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("Correlation-Context", tt.header)

			ctx := context.Background()
			_, gotCorCtx := propagator.Extract(ctx, req.Header)
			if gotCorCtx.Len() != 0 {
				t.Errorf("Got and Want CorCtx are not the same size %d != %d", gotCorCtx.Len(), 0)
			}
		})
	}
}

func TestInjectCorrelationContextToHTTPReq(t *testing.T) {
	propagator := propagation.HTTPTraceContextPropagator{}
	tests := []struct {
		name         string
		kvs          []core.KeyValue
		wantInHeader []string
		wantedLen    int
	}{
		{
			name: "two simple values",
			kvs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2"),
			},
			wantInHeader: []string{"key1=val1", "key2=val2"},
		},
		{
			name: "two values with escaped chars",
			kvs: []core.KeyValue{
				key.New("key1").String("val1,val2"),
				key.New("key2").String("val3=4"),
			},
			wantInHeader: []string{"key1=val1%2Cval2", "key2=val3%3D4"},
		},
		{
			name: "values of non-string types",
			kvs: []core.KeyValue{
				key.New("key1").Bool(true),
				key.New("key2").Int(123),
				key.New("key3").Int64(123),
				key.New("key4").Int32(123),
				key.New("key5").Uint(123),
				key.New("key6").Uint32(123),
				key.New("key7").Uint64(123),
				key.New("key8").Float64(123.567),
				key.New("key9").Float32(123.567),
			},
			wantInHeader: []string{
				"key1=true",
				"key2=123",
				"key3=123",
				"key4=123",
				"key5=123",
				"key6=123",
				"key7=123",
				"key8=123.567",
				"key9=123.567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			ctx := dctx.WithMap(context.Background(), dctx.NewMap(dctx.MapUpdate{MultiKV: tt.kvs}))
			propagator.Inject(ctx, req.Header)

			gotHeader := req.Header.Get("Correlation-Context")
			wantedLen := len(strings.Join(tt.wantInHeader, ","))
			if wantedLen != len(gotHeader) {
				t.Errorf(
					"%s: Inject Correlation-Context incorrect length %d != %d.", tt.name, tt.wantedLen, len(gotHeader),
				)
			}
			for _, inHeader := range tt.wantInHeader {
				if !strings.Contains(gotHeader, inHeader) {
					t.Errorf(
						"%s: Inject Correlation-Context missing part of header: %s in %s", tt.name, inHeader, gotHeader,
					)
				}
			}
		})
	}
}

func TestHTTPTraceContextPropagator_GetAllKeys(t *testing.T) {
	var propagator propagation.HTTPTraceContextPropagator
	want := []string{"Traceparent", "Correlation-Context"}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
