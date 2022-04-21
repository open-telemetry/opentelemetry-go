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
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (metric.MeterProvider, func()) {
	ctx := context.Background()

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	conn, err := grpc.DialContext(ctx, "localhost:30080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	handleErr(err, "failed to create gRPC connection to collector")

	client := otlpmetricgrpc.NewClient(otlpmetricgrpc.WithGRPCConn(conn))
	metricExporter, err := otlpmetric.New(ctx, client)
	handleErr(err, "failed to create the metric collector exporter")

	controller := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(2*time.Second),
	)

	err = controller.Start(ctx)
	handleErr(err, "failed to start metric controoler")

	return controller, func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(controller.Stop(ctx), "failed to stop meterProvider")
	}
}

func main() {
	ctx := context.Background()
	log.Printf("Waiting for connection...")

	meterProvider, shutdown := initProvider()
	defer shutdown()

	// TODO Bring Back Global package
	// meter := global.Meter("go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc_test")
	meter := meterProvider.Meter("test-metric")

	// Recorder metric example

	counter, err := meter.SyncFloat64().Counter("an_important_metric", instrument.WithDescription("Measures the cumulative epicness of the app"))
	if err != nil {
		log.Fatalf("Failed to create the instrument: %v", err)
	}

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []attribute.KeyValue{
		attribute.String("labelA", "chocolate"),
		attribute.String("labelB", "raspberry"),
		attribute.String("labelC", "vanilla"),
	}

	for i := 0; i < 10; i++ {
		counter.Add(ctx, 1.0, commonLabels...)
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Second)
	}

	log.Printf("Done!")
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
