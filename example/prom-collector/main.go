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
	"log"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

func initMeter() {
	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(attribute.String("R", "V")),
	)
	if err != nil {
		log.Fatal("could not initialize resource:", err)
	}

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint("localhost:30080"),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)
	otlpExporter, err := otlp.NewExporter(ctx, driver)

	if err != nil {
		log.Fatal("could not initialize OTLP:", err)
	}

	cont := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries([]float64{
					0.001, 0.01, 0.1, 1, 10, 100, 1000,
				}),
			),
			otlpExporter, // otlpExporter is an ExportKindSelector
			processor.WithMemory(true),
		),
		controller.WithResource(res),
		controller.WithExporter(otlpExporter),
	)

	if err := cont.Start(context.Background()); err != nil {
		log.Fatal("could not start controller:", err)
	}

	promExporter, err := prometheus.NewExporter(prometheus.Config{}, cont)
	if err != nil {
		log.Fatal("could not initialize prometheus:", err)
	}
	http.HandleFunc("/", promExporter.ServeHTTP)
	go func() {
		log.Fatal(http.ListenAndServe(":17000", nil))
	}()

	global.SetMeterProvider(cont.MeterProvider())

	log.Println("Prometheus server running on :17000")
	log.Println("Exporting OTLP to :30080")
}

func main() {
	initMeter()

	labels := []attribute.KeyValue{
		attribute.String("label1", "value1"),
	}

	meter := global.Meter("ex.com/prom-collector")
	_ = metric.Must(meter).NewFloat64ValueObserver(
		"randval",
		func(_ context.Context, result metric.Float64ObserverResult) {
			result.Observe(
				rand.Float64(),
				labels...,
			)
		},
		metric.WithDescription("A random value"),
	)

	temperature := metric.Must(meter).NewFloat64ValueRecorder("temperature")
	interrupts := metric.Must(meter).NewInt64Counter("interrupts")

	ctx := context.Background()

	log.Println("Example is running, please visit :17000")

	for {
		temperature.Record(ctx, 100+10*rand.NormFloat64(), labels...)
		interrupts.Add(ctx, int64(rand.Intn(100)), labels...)

		time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	}
}
