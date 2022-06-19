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

package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	instrumentationName = "instrumentationExample"
)

func main() {
	ctx := context.Background()
	shutdown, err := initExportPipeline(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// get global tracer and create span
	newCtx, span := otel.Tracer(instrumentationName).Start(ctx, "main func")
	defer span.End()
	for i := 0; i < 3; i++ {
		triggerOTLPToCollector(newCtx, i)
	}
}

func triggerOTLPToCollector(ctx context.Context, num int) {
	_, span := otel.Tracer(instrumentationName).Start(ctx, "triggerOTLPToCollector func")
	span.SetAttributes(attribute.Int("requestNum", num))

	defer span.End()
}

func initExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP trace exporter: %v", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("example-logging"),
		semconv.ServiceVersionKey.String("0.1.0"),
		attribute.String("environment", "test example logging"),
	)
}
