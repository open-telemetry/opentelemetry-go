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
	"net/http"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

func main() {
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/open-telemetry/opentelemetry-go/example/prometheus")

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(exporter)

	attrs := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.SyncFloat64().Counter("foo", instrument.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, attrs...)

	gauge, err := meter.SyncFloat64().UpDownCounter("bar", instrument.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	gauge.Add(ctx, 100, attrs...)
	gauge.Add(ctx, -25, attrs...)

	// This is the equivalent of prometheus.NewHistogramVec
	histogram, err := meter.SyncFloat64().Histogram("baz", instrument.WithDescription("a very nice histogram"))
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 23, attrs...)
	histogram.Record(ctx, 7, attrs...)
	histogram.Record(ctx, 101, attrs...)
	histogram.Record(ctx, 105, attrs...)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics(exp *prometheus.Exporter) {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", exp)
	err := http.ListenAndServe(":2223", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
