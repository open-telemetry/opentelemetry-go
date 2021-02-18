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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/baggage"
	"go.opentelemetry.io/otel/propagation"
)

func TestExtractValidBaggageFromHTTPReq(t *testing.T) {
	prop := propagation.TextMapPropagator(propagation.Baggage{})
	tests := []struct {
		name    string
		header  string
		wantKVs []attribute.KeyValue
	}{
		{
			name:   "valid w3cHeader",
			header: "key1=val1,key2=val2",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
		},
		{
			name:   "valid w3cHeader with spaces",
			header: "key1 =   val1,  key2 =val2   ",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
		},
		{
			name:   "valid w3cHeader with properties",
			header: "key1=val1,key2=val2;prop=1",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2;prop=1"),
			},
		},
		{
			name:   "valid header with url-escaped comma",
			header: "key1=val1,key2=val2%2Cval3",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2,val3"),
			},
		},
		{
			name:   "valid header with an invalid header",
			header: "key1=val1,key2=val2,a,val3",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
		},
		{
			name:   "valid header with no value",
			header: "key1=,key2=val2",
			wantKVs: []attribute.KeyValue{
				attribute.String("key1", ""),
				attribute.String("key2", "val2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("baggage", tt.header)

			ctx := context.Background()
			ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
			gotBaggage := baggage.MapFromContext(ctx)
			wantBaggage := baggage.NewMap(baggage.MapUpdate{MultiKV: tt.wantKVs})
			if gotBaggage.Len() != wantBaggage.Len() {
				t.Errorf(
					"Got and Want Baggage are not the same size %d != %d",
					gotBaggage.Len(),
					wantBaggage.Len(),
				)
			}
			totalDiff := ""
			wantBaggage.Foreach(func(keyValue attribute.KeyValue) bool {
				val, _ := gotBaggage.Value(keyValue.Key)
				diff := cmp.Diff(keyValue, attribute.KeyValue{Key: keyValue.Key, Value: val}, cmp.AllowUnexported(attribute.Value{}))
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
	prop := propagation.TextMapPropagator(propagation.Baggage{})
	tests := []struct {
		name   string
		header string
		hasKVs []attribute.KeyValue
	}{
		{
			name:   "no key values",
			header: "header1",
		},
		{
			name:   "invalid header with existing context",
			header: "header2",
			hasKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
		},
		{
			name:   "empty header value",
			header: "",
			hasKVs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("baggage", tt.header)

			ctx := baggage.NewContext(context.Background(), tt.hasKVs...)
			wantBaggage := baggage.MapFromContext(ctx)
			ctx = prop.Extract(ctx, propagation.HeaderCarrier(req.Header))
			gotBaggage := baggage.MapFromContext(ctx)
			if gotBaggage.Len() != wantBaggage.Len() {
				t.Errorf(
					"Got and Want Baggage are not the same size %d != %d",
					gotBaggage.Len(),
					wantBaggage.Len(),
				)
			}
			totalDiff := ""
			wantBaggage.Foreach(func(keyValue attribute.KeyValue) bool {
				val, _ := gotBaggage.Value(keyValue.Key)
				diff := cmp.Diff(keyValue, attribute.KeyValue{Key: keyValue.Key, Value: val}, cmp.AllowUnexported(attribute.Value{}))
				if diff != "" {
					totalDiff += diff + "\n"
				}
				return true
			})
		})
	}
}

func TestInjectBaggageToHTTPReq(t *testing.T) {
	propagator := propagation.Baggage{}
	tests := []struct {
		name         string
		kvs          []attribute.KeyValue
		wantInHeader []string
		wantedLen    int
	}{
		{
			name: "two simple values",
			kvs: []attribute.KeyValue{
				attribute.String("key1", "val1"),
				attribute.String("key2", "val2"),
			},
			wantInHeader: []string{"key1=val1", "key2=val2"},
		},
		{
			name: "two values with escaped chars",
			kvs: []attribute.KeyValue{
				attribute.String("key1", "val1,val2"),
				attribute.String("key2", "val3=4"),
			},
			wantInHeader: []string{"key1=val1%2Cval2", "key2=val3%3D4"},
		},
		{
			name: "values of non-string types",
			kvs: []attribute.KeyValue{
				attribute.Bool("key1", true),
				attribute.Int("key2", 123),
				attribute.Int64("key3", 123),
				attribute.Float64("key4", 123.567),
			},
			wantInHeader: []string{
				"key1=true",
				"key2=123",
				"key3=123",
				"key4=123.567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://example.com", nil)
			ctx := baggage.ContextWithMap(context.Background(), baggage.NewMap(baggage.MapUpdate{MultiKV: tt.kvs}))
			propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

			gotHeader := req.Header.Get("baggage")
			wantedLen := len(strings.Join(tt.wantInHeader, ","))
			if wantedLen != len(gotHeader) {
				t.Errorf(
					"%s: Inject baggage incorrect length %d != %d.", tt.name, tt.wantedLen, len(gotHeader),
				)
			}
			for _, inHeader := range tt.wantInHeader {
				if !strings.Contains(gotHeader, inHeader) {
					t.Errorf(
						"%s: Inject baggage missing part of header: %s in %s", tt.name, inHeader, gotHeader,
					)
				}
			}
		})
	}
}

func TestBaggagePropagatorGetAllKeys(t *testing.T) {
	var propagator propagation.Baggage
	want := []string{"baggage"}
	got := propagator.Fields()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
