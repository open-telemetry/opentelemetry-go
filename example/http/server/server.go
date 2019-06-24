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

package main

import (
	"io"
	"net/http"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/log"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/api/trace"
	"github.com/open-telemetry/opentelemetry-go/plugin/httptrace"

	_ "github.com/open-telemetry/opentelemetry-go/exporter/loader"
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
