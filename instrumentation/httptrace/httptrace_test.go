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

package httptrace_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/instrumentation/httptrace"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestRoundtrip(t *testing.T) {
	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	tr := tp.Tracer("httptrace/client")

	var expectedAttrs map[kv.Key]string
	expectedCorrs := map[kv.Key]string{kv.Key("foo"): "bar"}

	// Mock http server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attrs, corrs, span := httptrace.Extract(r.Context(), r)

			actualAttrs := make(map[kv.Key]string)
			for _, attr := range attrs {
				if attr.Key == standard.NetPeerPortKey {
					// Peer port will be non-deterministic
					continue
				}
				actualAttrs[attr.Key] = attr.Value.Emit()
			}

			if diff := cmp.Diff(actualAttrs, expectedAttrs); diff != "" {
				t.Fatalf("[TestRoundtrip] Attributes are different: %v", diff)
			}

			actualCorrs := make(map[kv.Key]string)
			for _, corr := range corrs {
				actualCorrs[corr.Key] = corr.Value.Emit()
			}

			if diff := cmp.Diff(actualCorrs, expectedCorrs); diff != "" {
				t.Fatalf("[TestRoundtrip] Correlations are different: %v", diff)
			}

			if !span.IsValid() {
				t.Fatalf("[TestRoundtrip] Invalid span extracted: %v", span)
			}

			_, err := w.Write([]byte("OK"))
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer ts.Close()

	address := ts.Listener.Addr()
	hp := strings.Split(address.String(), ":")
	expectedAttrs = map[kv.Key]string{
		standard.HTTPFlavorKey:               "1.1",
		standard.HTTPHostKey:                 address.String(),
		standard.HTTPMethodKey:               "GET",
		standard.HTTPSchemeKey:               "http",
		standard.HTTPTargetKey:               "/",
		standard.HTTPUserAgentKey:            "Go-http-client/1.1",
		standard.HTTPRequestContentLengthKey: "3",
		standard.NetHostIPKey:                hp[0],
		standard.NetHostPortKey:              hp[1],
		standard.NetPeerIPKey:                "127.0.0.1",
		standard.NetTransportKey:             "IP.TCP",
	}

	client := ts.Client()
	err := tr.WithSpan(context.Background(), "test",
		func(ctx context.Context) error {
			ctx = correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{SingleKV: kv.Key("foo").String("bar")}))
			req, _ := http.NewRequest("GET", ts.URL, strings.NewReader("foo"))
			httptrace.Inject(ctx, req)

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err.Error())
			}
			_ = res.Body.Close()

			return nil
		})
	if err != nil {
		panic("unexpected error in http request: " + err.Error())
	}
}

func TestSpecifyPropagators(t *testing.T) {
	exp := &testExporter{
		spanMap: make(map[string][]*export.SpanData),
	}
	tp, _ := sdktrace.NewProvider(sdktrace.WithSyncer(exp), sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	global.SetTraceProvider(tp)

	tr := tp.Tracer("httptrace/client")

	expectedCorrs := map[kv.Key]string{kv.Key("foo"): "bar"}

	// Mock http server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, corrs, span := httptrace.Extract(r.Context(), r, httptrace.WithPropagators(propagation.New(propagation.WithExtractors(correlation.DefaultHTTPPropagator()))))

			actualCorrs := make(map[kv.Key]string)
			for _, corr := range corrs {
				actualCorrs[corr.Key] = corr.Value.Emit()
			}

			if diff := cmp.Diff(actualCorrs, expectedCorrs); diff != "" {
				t.Fatalf("[TestRoundtrip] Correlations are different: %v", diff)
			}

			if span.IsValid() {
				t.Fatalf("[TestRoundtrip] valid span extracted, expected none: %v", span)
			}

			_, err := w.Write([]byte("OK"))
			if err != nil {
				t.Fatal(err)
			}
		}),
	)
	defer ts.Close()

	client := ts.Client()
	err := tr.WithSpan(context.Background(), "test",
		func(ctx context.Context) error {
			ctx = correlation.ContextWithMap(ctx, correlation.NewMap(correlation.MapUpdate{SingleKV: kv.Key("foo").String("bar")}))
			req, _ := http.NewRequest("GET", ts.URL, nil)
			httptrace.Inject(ctx, req, httptrace.WithPropagators(propagation.New(propagation.WithInjectors(correlation.DefaultHTTPPropagator()))))

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %s", err.Error())
			}
			_ = res.Body.Close()

			return nil
		})
	if err != nil {
		panic("unexpected error in http request: " + err.Error())
	}
}
