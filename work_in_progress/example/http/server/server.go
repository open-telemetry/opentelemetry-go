package main

import (
	"io"
	"net/http"

	"github.com/lightstep/opentelemetry-golang-prototype/api/core"
	"github.com/lightstep/opentelemetry-golang-prototype/api/log"
	"github.com/lightstep/opentelemetry-golang-prototype/api/tag"
	"github.com/lightstep/opentelemetry-golang-prototype/api/trace"
	"github.com/lightstep/opentelemetry-golang-prototype/plugin/httptrace"

	_ "github.com/lightstep/opentelemetry-golang-prototype/exporter/loader"
)

var (
	tracer = trace.GlobalTracer().
		WithService("server").
		WithComponent("main").
		WithResources(
			tag.New("whatevs").String("nooooo"),
		)
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		attrs, tags, spanCtx := httptrace.Extract(req)

		req = req.WithContext(tag.WithMap(req.Context(), tag.NewMap(core.KeyValue{}, tags, core.Mutator{}, nil)))

		ctx, span := tracer.Start(
			req.Context(),
			"hello",
			trace.WithAttributes(attrs...),
			trace.ChildOf(spanCtx),
		)
		defer span.Finish()

		log.Log(ctx, "handling this...")

		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		panic(err)
	}
}
