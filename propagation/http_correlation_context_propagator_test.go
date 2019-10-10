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

	"go.opentelemetry.io/api/key"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/trace"
	mocktrace "go.opentelemetry.io/internal/trace"
	"go.opentelemetry.io/propagation"
)

func TestExtractValidDistributedContextFromHTTPReq(t *testing.T) {
	trace.SetGlobalTracer(&mocktrace.MockTracer{})
	propagator := propagation.HttpCorrelationContextPropagator()
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
			gotSc := propagator.Extract(ctx, req.Header)
			if diff := cmp.Diff(gotSc, tt.wantKVs); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestExtractInvalidDistributedContextFromHTTPReq(t *testing.T) {
	trace.SetGlobalTracer(&mocktrace.MockTracer{})
	propagator := propagation.HttpCorrelationContextPropagator()
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
			gotKvs := propagator.Extract(ctx, req.Header)
			if diff := cmp.Diff(len(gotKvs), 0); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestInjectCorrelationContextToHTTPReq(t *testing.T) {
	propagator := propagation.HttpCorrelationContextPropagator()
	tests := []struct {
		name       string
		kvs        []core.KeyValue
		wantHeader string
	}{
		{
			name: "two simple values",
			kvs: []core.KeyValue{
				key.New("key1").String("val1"),
				key.New("key2").String("val2"),
			},
			wantHeader: "key1=val1,key2=val2",
		},
		{
			name: "two values with escaped chars",
			kvs: []core.KeyValue{
				key.New("key1").String("val1,val2"),
				key.New("key2").String("val3=4"),
			},
			wantHeader: "key1=val1%2Cval2,key2=val3%3D4",
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
				key.New("key10").Bytes([]byte{0x68, 0x69}),
			},
			wantHeader: "key1=true,key2=123,key3=123,key4=123,key5=123,key6=123,key7=123,key8=123.567,key9=123.56700134277344,key10=hi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			propagator.Inject(tt.kvs, req.Header)

			gotHeader := req.Header.Get("Correlation-Context")
			if diff := cmp.Diff(gotHeader, tt.wantHeader); diff != "" {
				t.Errorf("Extract Tracecontext: %s: -got +want %s", tt.name, diff)
			}
		})
	}
}

func TestHttpCorrelationContextPropagator_GetAllKeys(t *testing.T) {
	propagator := propagation.HttpCorrelationContextPropagator()
	want := []string{"Correlation-Context"}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
