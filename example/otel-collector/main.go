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

// Example from otlp package: https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/open-telemetry/opentelemetry-collector/translator/conventions"
)

func main() {
	// If the OpenTelemetry Collector is running on a local cluster (minikube or microk8s),
	// it should be accessible through the NodePort service at the `localhost:30080` address.
	// Otherwise, replace `localhost` with the address of your cluster.
	// If you run the app inside k8s, then you can probably connect directly to the service through dns
	exp, err := otlp.NewExporter(otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
		otlp.WithGRPCDialOption(grpc.WithBlock()))
	if err != nil {
		log.Fatalf("Failed to create the collector exporter: %v", err)
	}
	defer func() {
		err := exp.Stop()
		if err != nil {
			log.Fatalf("Failed to stop the exporter: %v", err)
		}
	}()

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(resource.New(
			// the service name used to display traces in Jaeger
			kv.Key(conventions.AttributeServiceName).String("test-service"),
		)),
		sdktrace.WithSyncer(exp))
	if err != nil {
		log.Fatalf("error creating trace provider: %v\n", err)
	}

	tracer := tp.Tracer("test-tracer")

	// Then use the OpenTelemetry tracing library, like we normally would.
	ctx, span := tracer.Start(context.Background(), "CollectorExporter-Example")

	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		<-time.After(time.Second)
		iSpan.End()
	}

	span.End()
}
