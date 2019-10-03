package httptrace_test

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/exporter/trace/stdout"
	"go.opentelemetry.io/plugin/httptrace"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

func ExampleNewHandler() {
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

	importantFunc := func(ctx context.Context, s string) error {
		span := trace.CurrentSpan(ctx)
		// do stuff with the span as needed
		_ = span
		return nil
	}

	var mux http.ServeMux
	mux.HandleFunc("/hello",
		func(w http.ResponseWriter, r *http.Request) {
			d, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error reading body: ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			s := string(d)
			if s == "" {
				s = "unknown"
			}
			ctx := r.Context()

			// Wrap another function in it's own span
			if err := trace.CurrentSpan(ctx).Tracer().WithSpan(ctx, "importantFunc",
				func(ctx context.Context) error {
					return importantFunc(ctx, s)
				}); err != nil {
				log.Println("error doing important thing: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			n, err := io.WriteString(w, "Hello, "+s+"!\n")
			if err != nil {
				log.Printf("error writing reply after %d bytes: %s", n, err)
			}
		})

	if err := http.ListenAndServe(":7777",
		httptrace.NewHandler(&mux, "server",
			httptrace.WithMessageEvents(httptrace.EventRead, httptrace.EventWrite),
		),
	); err != nil {
		log.Fatal(err)
	}
}
