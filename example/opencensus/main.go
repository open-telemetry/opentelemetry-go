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
	"time"

	ocmetric "go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// instrumenttype differentiates between our gauge and view metrics.
	keyType = tag.MustNewKey("instrumenttype")
	// Counts the number of lines read in from standard input.
	countMeasure = stats.Int64("test_count", "A count of something", stats.UnitDimensionless)
	countView    = &view.View{
		Name:        "test_count",
		Measure:     countMeasure,
		Description: "A count of something",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyType},
	}
)

func main() {
	log.Println("Using OpenTelemetry stdout exporters.")
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(fmt.Errorf("error creating trace exporter: %w", err))
	}
	metricsExporter, err := stdoutmetric.New()
	if err != nil {
		log.Fatal(fmt.Errorf("error creating metric exporter: %w", err))
	}
	tracing(traceExporter)
	if err := monitoring(metricsExporter); err != nil {
		log.Fatal(err)
	}
}

// tracing demonstrates overriding the OpenCensus DefaultTracer to send spans
// to the OpenTelemetry exporter by calling OpenCensus APIs.
func tracing(otExporter sdktrace.SpanExporter) {
	ctx := context.Background()

	log.Println("Configuring OpenCensus.  Not Registering any OpenCensus exporters.")
	octrace.ApplyConfig(octrace.Config{DefaultSampler: octrace.AlwaysSample()})

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(otExporter))
	otel.SetTracerProvider(tp)

	log.Println("Installing the OpenCensus bridge to make OpenCensus libraries write spans using OpenTelemetry.")
	tracer := tp.Tracer("simple")
	octrace.DefaultTracer = opencensus.NewTracer(tracer)
	tp.ForceFlush(ctx)

	log.Println("Creating OpenCensus span, which should be printed out using the OpenTelemetry stdouttrace exporter.\n-- It should have no parent, since it is the first span.")
	ctx, outerOCSpan := octrace.StartSpan(ctx, "OpenCensusOuterSpan")
	outerOCSpan.End()
	tp.ForceFlush(ctx)

	log.Println("Creating OpenTelemetry span\n-- It should have the OpenCensus span as a parent, since the OpenCensus span was written with using OpenTelemetry APIs.")
	ctx, otspan := tracer.Start(ctx, "OpenTelemetrySpan")
	otspan.End()
	tp.ForceFlush(ctx)

	log.Println("Creating OpenCensus span, which should be printed out using the OpenTelemetry stdouttrace exporter.\n-- It should have the OpenTelemetry span as a parent, since it was written using OpenTelemetry APIs")
	_, innerOCSpan := octrace.StartSpan(ctx, "OpenCensusInnerSpan")
	innerOCSpan.End()
	tp.ForceFlush(ctx)
}

// monitoring demonstrates creating an IntervalReader using the OpenTelemetry
// exporter to send metrics to the exporter by using either an OpenCensus
// registry or an OpenCensus view.
func monitoring(otExporter metric.Exporter) error {
	ctx := context.Background()
	log.Println("Using the OpenCensus bridge to export OpenCensus and OpenTelemetry metrics to a single OpenTelemetry exporter.")

	// Register the exporter with an SDK via a periodic reader.
	provider := metric.NewMeterProvider(
		metric.WithResource(resource.Default()),
		metric.WithReader(metric.NewPeriodicReader(otExporter)),
		// Add the OpenCensus producer to the SDK. This causes metrics from
		// OpenCensus to be included in the batch of metrics sent to our exporter.
		metric.WithProducer(opencensus.NewProducer()),
	)

	log.Println("Emitting a 'foo' metric using OpenTelemetry APIs, which is emitted with an OpenTelemetry stdout exporter")
	meter := provider.Meter("github.com/open-telemetry/opentelemetry-go/example/opencensus")
	counter, err := meter.SyncFloat64().Counter("foo", instrument.WithDescription("a simple counter"))
	if err != nil {
		return fmt.Errorf("failed to add otel counter: %w", err)
	}
	counter.Add(ctx, 5, []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}...)

	log.Println("Emitting 'test_gauge' and 'test_count' metrics using OpenCensus APIs. These are printed out using the same OpenTelemetry stdoutmetric exporter.")

	log.Println("Registering a gauge metric using an OpenCensus registry.")
	r := ocmetric.NewRegistry()
	metricproducer.GlobalManager().AddProducer(r)
	gauge, err := r.AddInt64Gauge(
		"test_gauge",
		ocmetric.WithDescription("A gauge for testing"),
		ocmetric.WithConstLabel(map[metricdata.LabelKey]metricdata.LabelValue{
			{Key: keyType.Name()}: metricdata.NewLabelValue("gauge"),
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to add gauge: %w", err)
	}
	entry, err := gauge.GetEntry()
	if err != nil {
		return fmt.Errorf("failed to get gauge entry: %w", err)
	}

	log.Println("Registering a cumulative metric using an OpenCensus view.")
	if err := view.Register(countView); err != nil {
		return fmt.Errorf("failed to register views: %w", err)
	}
	ctx, err = tag.New(context.Background(), tag.Insert(keyType, "view"))
	if err != nil {
		return fmt.Errorf("failed to set tag: %w", err)
	}
	for i := int64(1); true; i++ {
		// update stats for our gauge
		entry.Set(i)
		// update stats for our view
		stats.Record(ctx, countMeasure.M(1))
		time.Sleep(time.Second)
	}
	return nil
}
