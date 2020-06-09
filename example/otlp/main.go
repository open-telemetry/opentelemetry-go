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
//
// Example application showcasing opentelemetry Go using the OTLP wire
// protocol

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (*otlp.Exporter, *push.Controller) {
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress("localhost:55680"),
	)
	handleErr(err, "Failed to create exporter: $v")

	traceProvider, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exp),
	)
	handleErr(err, "Failed to create trace provider: %v")

	pusher := push.New(
		simple.NewWithExactDistribution(),
		exp,
		push.WithStateful(true),
		push.WithPeriod(2*time.Second),
	)

	global.SetTraceProvider(traceProvider)
	global.SetMeterProvider(pusher.Provider())
	pusher.Start()

	return exp, pusher
}

func main() {
	exp, pusher := initProvider()
	defer func() { handleErr(exp.Stop(), "Failed to stop exporter") }()
	defer pusher.Stop() // pushes any last exports to the receiver

	tracer := global.Tracer("mage-sense")
	meter := global.Meter("mage-read")

	// labels represent additional descriptors that can be bound to a metric
	// observer or recorder. In this case they describe the location in
	// which a spell scroll is scribed.
	commonLabels := []kv.KeyValue{
		kv.String("work-room", "East Scriptorium"),
		kv.String("occupancy", "69,105"),
		kv.String("priority", "Ultra"),
	}

	// Observer metric example
	oneMetricCB := func(_ context.Context, result metric.Float64ObserverResult) {
		result.Observe(1, commonLabels...)
	}
	_ = metric.Must(meter).NewFloat64ValueObserver("scrying.glass.one", oneMetricCB,
		metric.WithDescription("A ValueObserver set to 1.0"),
	)

	// Recorder metric example
	valuerecorder := metric.Must(meter).
		NewFloat64ValueRecorder("scrying.glass.two").
		Bind(commonLabels...)
	defer valuerecorder.Unbind()

	// work begins
	ctx, span := tracer.Start(context.Background(), "Archmage-Overlord-Inspection")
	for i := 0; i < 10; i++ {
		_, innerSpan := tracer.Start(ctx, fmt.Sprintf("Minion-%d", i))
		log.Println("Minions hard at work, scribing...")
		valuerecorder.Record(ctx, float64(i)*1.5)

		<-time.After(time.Second)
		innerSpan.End()
	}

	span.End()
	<-time.After(time.Second)

	log.Println("Spell-scroll scribed!")
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
