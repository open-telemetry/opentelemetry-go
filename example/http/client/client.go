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
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/plugin/httptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}

func main() {
	initTracer()

	client := http.DefaultClient
	ctx := distributedcontext.NewContext(context.Background(),
		key.String("username", "donuts"),
	)

	var body []byte

	tr := global.TraceProvider().Tracer("example/client")
	err := tr.WithSpan(ctx, "say hello",
		func(ctx context.Context) error {
			req, _ := http.NewRequest("GET", "http://localhost:7777/hello", nil)

			ctx, req = httptrace.W3C(ctx, req)
			httptrace.Inject(ctx, req)

			fmt.Printf("Sending request...\n")
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			body, err = ioutil.ReadAll(res.Body)
			res.Body.Close()
			trace.CurrentSpan(ctx).SetStatus(codes.OK)

			return err
		})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Received: %s\n\n\n", body)
	fmt.Printf("Waiting for few seconds to export spans ...\n\n")
	time.Sleep(10 * time.Second)
	fmt.Printf("Inspect traces on stdout")
}
