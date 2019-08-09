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
	"go.opentelemetry.io/api/tag"
	"go.opentelemetry.io/propagation"
)

func TestExtractSpanContextFromB3Headers(t *testing.T) {
	traceID := core.TraceID{High: 0x463ac35c9f6413ad, Low: 0x48485a3953bb6124}
	spanID := uint64(0x0020000000000001)
	spanIDStr := "0020000000000001"
	propagator := propagation.HttpB3Propagator()
	tests := []struct {
		name    string
		makeReq func() *http.Request
		wantSc  core.SpanContext
	}{
		{
			name: "128-bit trace ID + 64-bit span ID; sampled=1",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "463ac35c9f6413ad48485a3953bb6124")
				req.Header.Set(propagation.B3SpanIDHeader, "0020000000000001")
				req.Header.Set(propagation.B3SampledHeader, "1")
				return req
			},
			wantSc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
		},
		{
			name: "short trace ID + short span ID; sampled=1",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "000102")
				req.Header.Set(propagation.B3SpanIDHeader, "000102")
				req.Header.Set(propagation.B3SampledHeader, "1")
				return req
			},
			wantSc: core.SpanContext{
				TraceID:      core.TraceID{High: 0x0, Low: 0x102},
				SpanID:       0x102,
				TraceOptions: core.TraceOptionSampled,
			},
		},
		{
			name: "64-bit trace ID + 64-bit span ID; sampled=0",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "0020000000000001")
				req.Header.Set(propagation.B3SpanIDHeader, spanIDStr)
				req.Header.Set(propagation.B3SampledHeader, "0")
				return req
			},
			wantSc: core.SpanContext{
				TraceID: core.TraceID{High: 0x0, Low: 0x0020000000000001},
				SpanID:  spanID,
			},
		},
		{
			name: "128-bit trace ID + 64-bit span ID; no sampled header",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "463ac35c9f6413ad48485a3953bb6124")
				req.Header.Set(propagation.B3SpanIDHeader, spanIDStr)
				return req
			},
			wantSc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
		{
			name: "invalid trace ID + 64-bit span ID; no sampling header",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "")
				req.Header.Set(propagation.B3SpanIDHeader, "0020000000000001")
				return req
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name: "128-bit trace ID; invalid span ID; no sampling header",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "463ac35c9f6413ad48485a3953bb6124")
				return req
			},
			wantSc: core.EmptySpanContext(),
		},
		{
			name: "128-bit trace ID + 64-bit span ID; sampled=true",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "463ac35c9f6413ad48485a3953bb6124")
				req.Header.Set(propagation.B3SpanIDHeader, spanIDStr)
				req.Header.Set(propagation.B3SampledHeader, "true")
				return req
			},
			wantSc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
		},
		{
			name: "128-bit trace ID + 64-bit span ID; sampled=false",
			makeReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "http://example.com", nil)
				req.Header.Set(propagation.B3TraceIDHeader, "463ac35c9f6413ad48485a3953bb6124")
				req.Header.Set(propagation.B3SpanIDHeader, spanIDStr)
				req.Header.Set(propagation.B3SampledHeader, "false")
				return req
			},
			wantSc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.makeReq()
			f := propagator.CarrierExtractor(req)
			gotSc, _ := f.Extract()
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractSpanContextFromInvalidCarrierUsingB3Propagator(t *testing.T) {
	propagator := propagation.HttpB3Propagator()
	tests := []struct {
		name   string
		wantSc core.SpanContext
	}{
		{
			name:   "invalid carrier",
			wantSc: core.EmptySpanContext(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := propagator.CarrierExtractor("")
			gotSc, _ := f.Extract()
			if diff := cmp.Diff(gotSc, tt.wantSc); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestInjectSpanContextAsB3HeadersToHTTPReq(t *testing.T) {
	traceID := core.TraceID{High: 0x463ac35c9f6413ad, Low: 0x48485a3953bb6124}
	spanID := uint64(0x0020000000000001)
	propagator := propagation.HttpB3Propagator()
	tests := []struct {
		name        string
		sc          core.SpanContext
		wantHeaders map[string]string
	}{
		{
			name: "valid traceID, header ID, sampled=1",
			sc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
			wantHeaders: map[string]string{
				"X-B3-TraceId": "463ac35c9f6413ad48485a3953bb6124",
				"X-B3-SpanId":  "0020000000000001",
				"X-B3-Sampled": "1",
			},
		},
		{
			name: "valid traceID, header ID, sampled=0",
			sc: core.SpanContext{
				TraceID: traceID,
				SpanID:  spanID,
			},
			wantHeaders: map[string]string{
				"X-B3-TraceId": "463ac35c9f6413ad48485a3953bb6124",
				"X-B3-SpanId":  "0020000000000001",
				"X-B3-Sampled": "0",
			},
		},
		{
			name: "invalid spancontext",
			sc:   core.EmptySpanContext(),
			wantHeaders: map[string]string{
				"X-B3-TraceId": "",
				"X-B3-SpanId":  "",
				"X-B3-Sampled": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			f := propagator.CarrierInjector(req)
			f.Inject(tt.sc, tag.NewEmptyMap())

			for h, wantValue := range tt.wantHeaders {
				gotValue := req.Header.Get(h)
				if diff := cmp.Diff(gotValue, wantValue); diff != "" {
					t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
				}
			}
		})
	}
}

func TestInjectSpanContextAsB3HeadersToInvalidCarrier(t *testing.T) {
	propagator := propagation.HttpB3Propagator()
	tests := []struct {
		name string
		sc   core.SpanContext
	}{
		{
			name: "valid spancontext to invalid carrier does nothing.",
			sc: core.SpanContext{
				TraceID:      traceID,
				SpanID:       spanID,
				TraceOptions: core.TraceOptionSampled,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := propagator.CarrierInjector("")
			i.Inject(tt.sc, tag.NewEmptyMap())
		})
	}
}
