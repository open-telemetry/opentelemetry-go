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

// Command zipkin is an example program that creates spans
// and uploads to openzipkin collector.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(url string) func() {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.NewRawExporter(
		url,
		zipkin.WithLogger(logger),
		zipkin.WithSDKOptions(sdktrace.WithSampler(sdktrace.AlwaysSample())),
	)
	if err != nil {
		log.Fatal(err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String("zipkin-test"),
		)),
	)
	otel.SetTracerProvider(tp)

	return func() {
		_ = tp.Shutdown(context.Background())
	}
}

func main() {
	url := flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	shutdown := initTracer(*url)
	defer shutdown()

	ctx := context.Background()

	tr := otel.GetTracerProvider().Tracer("component-main")
	ctx, span := tr.Start(ctx, "foo", trace.WithSpanKind(trace.SpanKindServer))
	<-time.After(6 * time.Millisecond)
	bar(ctx)
	<-time.After(6 * time.Millisecond)
	span.End()
}

func bar(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	<-time.After(6 * time.Millisecond)
	span.End()
}
