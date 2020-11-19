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

package main // import "go.opentelemetry.io/otel/bridge/opencensus/examples/simple"

import (
	"context"
	"log"

	octrace "go.opencensus.io/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()

	log.Println("Configuring OpenCensus.  Not Registering any OpenCensus exporters.")
	octrace.ApplyConfig(octrace.Config{DefaultSampler: octrace.AlwaysSample()})

	log.Println("Registering OpenTelemetry stdout exporter.")
	otExporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(otExporter))
	otel.SetTracerProvider(tp)

	log.Println("Installing the OpenCensus bridge to make OpenCensus libraries write spans using OpenTelemetry.")
	tracer := tp.Tracer("simple")
	octrace.DefaultTracer = opencensus.NewTracer(tracer)

	log.Println("Creating OpenCensus span, which should be printed out using the OpenTelemetry stdout exporter.\n-- It should have no parent, since it is the first span.")
	ctx, outerOCSpan := octrace.StartSpan(ctx, "OpenCensusOuterSpan")
	outerOCSpan.End()

	log.Println("Creating OpenTelemetry span\n-- It should have the OpenCensus span as a parent, since the OpenCensus span was written with using OpenTelemetry APIs.")
	ctx, otspan := tracer.Start(ctx, "OpenTelemetrySpan")
	otspan.End()

	log.Println("Creating OpenCensus span, which should be printed out using the OpenTelemetry stdout exporter.\n-- It should have the OpenTelemetry span as a parent, since it was written using OpenTelemetry APIs")
	_, innerOCSpan := octrace.StartSpan(ctx, "OpenCensusInnerSpan")
	innerOCSpan.End()
}
