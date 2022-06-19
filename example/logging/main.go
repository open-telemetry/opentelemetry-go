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
