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

// Example using the OTLP exporter + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/standard"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (*otlp.Exporter, *push.Controller) {

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` address. Otherwise, replace `localhost` with the
	// address of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:30080"),
		otlp.WithGRPCDialOption(grpc.WithBlock()), // useful for testing
	)
	handleErr(err, "failed to create exporter")

	traceProvider, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithResource(resource.New(
			// the service name used to display traces in backends
			kv.Key(standard.ServiceNameKey).String("test-service"),
		)),
		sdktrace.WithSyncer(exp),
	)
	handleErr(err, "failed to create trace provider")

	pusher := push.New(
		simple.NewWithExactDistribution(),
		exp,
		push.WithPeriod(2*time.Second),
	)

	global.SetTraceProvider(traceProvider)
	global.SetMeterProvider(pusher.Provider())
	pusher.Start()

	return exp, pusher
}

func main() {
	log.Printf("Waiting for connection...")

	exp, pusher := initProvider()
	defer func() { handleErr(exp.Stop(), "failed to stop exporter") }()
	defer pusher.Stop() // pushes any last exports to the receiver

	tracer := global.Tracer("test-tracer")
	meter := global.Meter("test-meter")

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []kv.KeyValue{
		kv.String("labelA", "chocolate"),
		kv.String("labelB", "raspberry"),
		kv.String("labelC", "vanilla"),
	}

	// Recorder metric example
	valuerecorder := metric.Must(meter).
		NewFloat64Counter(
			"an_important_metric",
			metric.WithDescription("Measures the cumulative epicness of the app"),
		).Bind(commonLabels...)
	defer valuerecorder.Unbind()

	// work begins
	ctx, span := tracer.Start(context.Background(), "CollectorExporter-Example")
	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)
		valuerecorder.Add(ctx, 1.0)

		<-time.After(time.Second)
		iSpan.End()
	}

	log.Printf("Done!")
	span.End()
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
