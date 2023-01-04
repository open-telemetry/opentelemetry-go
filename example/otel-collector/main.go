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

// Example using OTLP exporters + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	metricglobal "go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initOTEL(ctx context.Context) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display in observer backends
			semconv.ServiceNameKey.String("test-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// OpenTelemtry Collector is running on docker-compose in same network as this go app,
	// and exposes the OTLP receiver on port 4317.
	conn, err := grpc.DialContext(ctx, "otel-collector:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceShutdownFunc, err := initTracerProvider(ctx, conn, res)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize trace provider: %w", err)
	}
	metricShutdownFunc, err := initMetricProvider(ctx, conn, res)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metric provider: %w", err)
	}
	return func(ctx context.Context) error {
		if err := traceShutdownFunc(ctx); err != nil {
			return fmt.Errorf("failed to shutdown trace provider: %w", err)
		}
		if err := metricShutdownFunc(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metric provider: %w", err)
		}
		return nil
	}, nil
}

func initMetricProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (func(context.Context) error, error) {
	// Set up a metric exporter
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	metricProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(5*time.Second),
			),
		),
		sdkmetric.WithResource(res),
	)

	metricglobal.SetMeterProvider(metricProvider)
	return func(ctx context.Context) error {
		if err := metricProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown MeterProvider: %w", err)
		}
		return nil
	}, nil
}

func initTracerProvider(ctx context.Context, conn *grpc.ClientConn, res *resource.Resource) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown TracerProvider: %w", err)
		}
		return nil
	}, nil
}

func main() {
	log.Printf("Waiting for connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initOTEL(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown: %w", err)
		}
	}()

	tracer := otel.Tracer("test-tracer")
	counter, err := metricglobal.Meter("MetricCollectorExporter-Example").Int64Counter("work_duration")
	if err != nil {
		log.Fatal(err)
	}
	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}
	// work begins
	ctx, span := tracer.Start(
		ctx,
		"TraceCollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		startTime := time.Now()
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Duration(i) * time.Millisecond * 100)
		counter.Add(ctx, time.Since(startTime).Milliseconds(), commonAttrs...)
		iSpan.End()
	}

	log.Printf("Done!")
}
