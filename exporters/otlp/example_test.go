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

package otlp_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Example_insecure() {
	ctx := context.Background()
	exp, err := otlp.NewExporter(ctx, otlp.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create the collector exporter: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("test-tracer")

	// Then use the OpenTelemetry tracing library, like we normally would.
	ctx, span := tracer.Start(ctx, "CollectorExporter-Example")
	defer span.End()

	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		<-time.After(6 * time.Millisecond)
		iSpan.End()
	}
}

func Example_withTLS() {
	// Please take at look at https://pkg.go.dev/google.golang.org/grpc/credentials#TransportCredentials
	// for ways on how to initialize gRPC TransportCredentials.
	creds, err := credentials.NewClientTLSFromFile("my-cert.pem", "")
	if err != nil {
		log.Fatalf("failed to create gRPC client TLS credentials: %v", err)
	}

	ctx := context.Background()
	exp, err := otlp.NewExporter(ctx, otlp.WithTLSCredentials(creds))
	if err != nil {
		log.Fatalf("failed to create the collector exporter: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithBatcher(
			exp,
			// add following two options to ensure flush
			sdktrace.WithBatchTimeout(5),
			sdktrace.WithMaxExportBatchSize(10),
		),
	)
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("test-tracer")

	// Then use the OpenTelemetry tracing library, like we normally would.
	ctx, span := tracer.Start(ctx, "Securely-Talking-To-Collector-Span")
	defer span.End()

	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		<-time.After(6 * time.Millisecond)
		iSpan.End()
	}
}
