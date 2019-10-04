package httptrace_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/exporter/trace/stdout"
	"go.opentelemetry.io/plugin/httptrace"
	sdktrace "go.opentelemetry.io/sdk/trace"
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

	//import sdktrace "go.opentelemetry.io/sdk/trace"
	sdktrace.Register()

	// Write spans to stdout
	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}

	// Wrap stdout exporter with SimpleSpanProcessor and register the processor.
	ssp := sdktrace.NewSimpleSpanProcessor(exporter)
	sdktrace.RegisterSpanProcessor(ssp)

	// For the example, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})

	figureOutName := func(ctx context.Context, s string) (string, error) {
		pp := strings.SplitN(s, "/", 2)
		var err error
		switch pp[1] {
		case "":
			err = fmt.Errorf("expected /hello/:name in %q", s)
		default:
			span := trace.CurrentSpan(ctx)
			span.SetAttribute(
				core.KeyValue{Key: core.Key{Name: "name"},
					Value: core.Value{Type: core.STRING, String: pp[1]},
				},
			)
		}
		return pp[1], err
	}

	var mux http.ServeMux
	mux.Handle("/hello/",
		httptrace.WithRouteTag("/hello/:name", http.HandlerFunc(
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
		httptrace.NewHandler(&mux, "server",
			httptrace.WithMessageEvents(httptrace.EventRead, httptrace.EventWrite),
		),
	); err != nil {
		log.Fatal(err)
	}
}
