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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

const meterName = "github.com/open-telemetry/opentelemetry-go/example/view"

func main() {
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader, allowing it to be used in WithReader.
	exporter := otelprom.New()

	// View to customize histogram buckets and rename a single histogram instrument.
	customBucketsView, err := view.New(
		// Match* to match instruments
		view.MatchInstrumentName("custom_histogram"),
		view.MatchInstrumentationScope(instrumentation.Scope{Name: meterName}),

		// With* to modify instruments
		view.WithSetAggregation(aggregation.ExplicitBucketHistogram{
			Boundaries: []float64{64, 128, 256, 512, 1024, 2048, 4096},
		}),
		view.WithRename("bar"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Default view to keep all instruments
	defaultView, err := view.New(view.MatchInstrumentName("*"))
	if err != nil {
		log.Fatal(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter, customBucketsView, defaultView))
	meter := provider.Meter(meterName)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(exporter.Collector)

	attrs := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}

	counter, err := meter.SyncFloat64().Counter("foo", instrument.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, attrs...)

	histogram, err := meter.SyncFloat64().Histogram("custom_histogram", instrument.WithDescription("a histogram with custom buckets and rename"))
	if err != nil {
		log.Fatal(err)
	}
	histogram.Record(ctx, 136, attrs...)
	histogram.Record(ctx, 64, attrs...)
	histogram.Record(ctx, 701, attrs...)
	histogram.Record(ctx, 830, attrs...)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics(collector prometheus.Collector) {
	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	if err != nil {
		fmt.Printf("error registering collector: %v", err)
		return
	}

	log.Printf("serving metrics at localhost:2222/metrics")
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err = http.ListenAndServe(":2222", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
