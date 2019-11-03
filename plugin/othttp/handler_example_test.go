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
package othttp_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/plugin/othttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func ExampleNewHandler() {
	/* curl -v -d "a painting" http://localhost:7777/hello/bob/ross
	...
	* upload completely sent off: 10 out of 10 bytes
	< HTTP/1.1 200 OK
	< Traceparent: 00-76ae040ee5753f38edf1c2bd9bd128bd-dd394138cfd7a3dc-01
	< Date: Fri, 04 Oct 2019 02:33:08 GMT
	< Content-Length: 45
	< Content-Type: text/plain; charset=utf-8
	<
	Hello, bob/ross!
	You sent me this:
	a painting
	*/

	// Write spans to stdout
	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}

	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)

	figureOutName := func(ctx context.Context, s string) (string, error) {
		pp := strings.SplitN(s, "/", 2)
		var err error
		switch pp[1] {
		case "":
			err = fmt.Errorf("expected /hello/:name in %q", s)
		default:
			trace.CurrentSpan(ctx).SetAttribute(core.Key("name").String(pp[1]))
		}
		return pp[1], err
	}

	var mux http.ServeMux
	mux.Handle("/hello/",
		othttp.WithRouteTag("/hello/:name", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				var name string
				// Wrap another function in it's own span
				if err := trace.CurrentSpan(ctx).Tracer().WithSpan(ctx, "figureOutName",
					func(ctx context.Context) error {
						var err error
						name, err = figureOutName(ctx, r.URL.Path[1:])
						return err
					}); err != nil {
					log.Println("error figuring out name: ", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				d, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Println("error reading body: ", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				n, err := io.WriteString(w, "Hello, "+name+"!\nYou sent me this:\n"+string(d))
				if err != nil {
					log.Printf("error writing reply after %d bytes: %s", n, err)
				}
			}),
		),
	)

	if err := http.ListenAndServe(":7777",
		othttp.NewHandler(&mux, "server",
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
		),
	); err != nil {
		log.Fatal(err)
	}
}
