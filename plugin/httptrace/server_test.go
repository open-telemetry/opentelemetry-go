package httptrace_test

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"go.opentelemetry.io/exporter/trace/jaeger"
	"go.opentelemetry.io/plugin/httptrace"
	"go.opentelemetry.io/sdk/trace"
)

func ExampleNewHandler() {
	trace.Register()

	// Create Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "trace-demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Wrap Jaeger exporter with SimpleSpanProcessor and register the processor.
	ssp := trace.NewSimpleSpanProcessor(exporter)
	trace.RegisterSpanProcessor(ssp)

	// For the example, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	helloHandler := func(w http.ResponseWriter, r *http.Request) {

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error reading body: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		n, err := io.WriteString(w, "Hello, world!\n"+string(d)+"\n")
		if err != nil {
			log.Printf("error writing reply after %d bytes: %s", n, err)
		}

	}

	http.Handle("/hello", httptrace.NewHandler(http.HandlerFunc(helloHandler), "hello"))
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal(err)
	}
}
