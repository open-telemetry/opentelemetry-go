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

	"go.opencensus.io/metric/metricdata"

	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/metric/metricproducer"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otmetricexport "go.opentelemetry.io/otel/sdk/export/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	// instrumenttype differentiates between our gauge and view metrics.
	keyType = tag.MustNewKey("instrumenttype")
	// Counts the number of lines read in from standard input
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
	metricsExporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		log.Fatal(fmt.Errorf("error creating metric exporter: %w", err))
	}
	tracing(traceExporter)
	monitoring(metricsExporter)
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
func monitoring(otExporter otmetricexport.Exporter) {
	log.Println("Using the OpenTelemetry stdoutmetric exporter to export OpenCensus metrics.  This allows routing telemetry from both OpenTelemetry and OpenCensus to a single exporter.")
	ocExporter := opencensus.NewMetricExporter(otExporter)
	intervalReader, err := metricexport.NewIntervalReader(&metricexport.Reader{}, ocExporter)
	if err != nil {
		log.Fatalf("Failed to create interval reader: %v\n", err)
	}
	intervalReader.ReportingInterval = 10 * time.Second
	log.Println("Emitting metrics using OpenCensus APIs.  These should be printed out using the OpenTelemetry stdoutmetric exporter.")
	err = intervalReader.Start()
	if err != nil {
		log.Fatalf("Failed to start interval reader: %v\n", err)
	}
	defer intervalReader.Stop()

	log.Println("Registering a gauge metric using an OpenCensus registry.")
	r := metric.NewRegistry()
	metricproducer.GlobalManager().AddProducer(r)
	gauge, err := r.AddInt64Gauge(
		"test_gauge",
		metric.WithDescription("A gauge for testing"),
		metric.WithConstLabel(map[metricdata.LabelKey]metricdata.LabelValue{
			{Key: keyType.Name()}: metricdata.NewLabelValue("gauge"),
		}),
	)
	if err != nil {
		log.Fatalf("Failed to add gauge: %v\n", err)
	}
	entry, err := gauge.GetEntry()
	if err != nil {
		log.Fatalf("Failed to get gauge entry: %v\n", err)
	}

	log.Println("Registering a cumulative metric using an OpenCensus view.")
	if err := view.Register(countView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}
	ctx, err := tag.New(context.Background(), tag.Insert(keyType, "view"))
	if err != nil {
		log.Fatalf("Failed to set tag: %v\n", err)
	}
	for i := int64(1); true; i++ {
		// update stats for our gauge
		entry.Set(i)
		// update stats for our view
		stats.Record(ctx, countMeasure.M(1))
		time.Sleep(time.Second)
	}
}
