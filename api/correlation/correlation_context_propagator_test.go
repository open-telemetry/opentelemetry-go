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

package correlation_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/api/kv/value"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/propagation"
)

func TestExtractValidDistributedContextFromHTTPReq(t *testing.T) {
	props := propagation.New(propagation.WithExtractors(correlation.CorrelationContext{}))
	tests := []struct {
		name    string
		header  string
		wantKVs []kv.KeyValue
	}{
		{
			name:   "valid w3cHeader",
			header: "key1=val1,key2=val2",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2"),
			},
		},
		{
			name:   "valid w3cHeader with spaces",
			header: "key1 =   val1,  key2 =val2   ",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2"),
			},
		},
		{
			name:   "valid w3cHeader with properties",
			header: "key1=val1,key2=val2;prop=1",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2;prop=1"),
			},
		},
		{
			name:   "valid header with url-escaped comma",
			header: "key1=val1,key2=val2%2Cval3",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2,val3"),
			},
		},
		{
			name:   "valid header with an invalid header",
			header: "key1=val1,key2=val2,a,val3",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2"),
			},
		},
		{
			name:   "valid header with no value",
			header: "key1=,key2=val2",
			wantKVs: []kv.KeyValue{
				kv.Key("key1").String(""),
				kv.Key("key2").String("val2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("otcorrelations", tt.header)

			ctx := context.Background()
			ctx = propagation.ExtractHTTP(ctx, props, req.Header)
			gotCorCtx := correlation.MapFromContext(ctx)
			wantCorCtx := correlation.NewMap(correlation.MapUpdate{MultiKV: tt.wantKVs})
			if gotCorCtx.Len() != wantCorCtx.Len() {
				t.Errorf(
					"Got and Want CorCtx are not the same size %d != %d",
					gotCorCtx.Len(),
					wantCorCtx.Len(),
				)
			}
			totalDiff := ""
			wantCorCtx.Foreach(func(keyValue kv.KeyValue) bool {
				val, _ := gotCorCtx.Value(keyValue.Key)
				diff := cmp.Diff(keyValue, kv.KeyValue{Key: keyValue.Key, Value: val}, cmp.AllowUnexported(value.Value{}))
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
	props := propagation.New(propagation.WithExtractors(correlation.CorrelationContext{}))
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
			req.Header.Set("otcorrelations", tt.header)

			ctx := context.Background()
			ctx = propagation.ExtractHTTP(ctx, props, req.Header)
			gotCorCtx := correlation.MapFromContext(ctx)
			if gotCorCtx.Len() != 0 {
				t.Errorf("Got and Want CorCtx are not the same size %d != %d", gotCorCtx.Len(), 0)
			}
		})
	}
}

func TestInjectCorrelationContextToHTTPReq(t *testing.T) {
	propagator := correlation.CorrelationContext{}
	props := propagation.New(propagation.WithInjectors(propagator))
	tests := []struct {
		name         string
		kvs          []kv.KeyValue
		wantInHeader []string
		wantedLen    int
	}{
		{
			name: "two simple values",
			kvs: []kv.KeyValue{
				kv.Key("key1").String("val1"),
				kv.Key("key2").String("val2"),
			},
			wantInHeader: []string{"key1=val1", "key2=val2"},
		},
		{
			name: "two values with escaped chars",
			kvs: []kv.KeyValue{
				kv.Key("key1").String("val1,val2"),
				kv.Key("key2").String("val3=4"),
			},
			wantInHeader: []string{"key1=val1%2Cval2", "key2=val3%3D4"},
		},
		{
			name: "values of non-string types",
			kvs: []kv.KeyValue{
				kv.Key("key1").Bool(true),
				kv.Key("key2").Int(123),
				kv.Key("key3").Int64(123),
				kv.Key("key4").Int32(123),
				kv.Key("key5").Uint(123),
				kv.Key("key6").Uint32(123),
				kv.Key("key7").Uint64(123),
				kv.Key("key8").Float64(123.567),
				kv.Key("key9").Float32(123.567),
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
			ctx := correlation.ContextWithMap(context.Background(), correlation.NewMap(correlation.MapUpdate{MultiKV: tt.kvs}))
			propagation.InjectHTTP(ctx, props, req.Header)

			gotHeader := req.Header.Get("otcorrelations")
			wantedLen := len(strings.Join(tt.wantInHeader, ","))
			if wantedLen != len(gotHeader) {
				t.Errorf(
					"%s: Inject otcorrelations incorrect length %d != %d.", tt.name, tt.wantedLen, len(gotHeader),
				)
			}
			for _, inHeader := range tt.wantInHeader {
				if !strings.Contains(gotHeader, inHeader) {
					t.Errorf(
						"%s: Inject otcorrelations missing part of header: %s in %s", tt.name, inHeader, gotHeader,
					)
				}
			}
		})
	}
}

func TestTraceContextPropagator_GetAllKeys(t *testing.T) {
	var propagator correlation.CorrelationContext
	want := []string{"otcorrelations"}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
