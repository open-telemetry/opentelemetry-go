package correlation_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/propagation"
)

func TestExtractValidDistributedContextFromHTTPReq(t *testing.T) {
	props := propagation.New(propagation.WithExtractors(correlation.CorrelationContext{}))
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
			req.Header.Set("Correlation-Context", tt.header)

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
			ctx := correlation.ContextWithMap(context.Background(), correlation.NewMap(correlation.MapUpdate{MultiKV: tt.kvs}))
			propagation.InjectHTTP(ctx, props, req.Header)

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

func TestTraceContextPropagator_GetAllKeys(t *testing.T) {
	var propagator correlation.CorrelationContext
	want := []string{"Correlation-Context"}
	got := propagator.GetAllKeys()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("GetAllKeys: -got +want %s", diff)
	}
}
