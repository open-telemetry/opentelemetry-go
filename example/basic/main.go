// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/context/scope"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	metricstdout "go.opentelemetry.io/otel/exporter/metric/stdout"
	tracestdout "go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	fooKey     = key.New("ex.com/foo")
	barKey     = key.New("ex.com/bar")
	lemonsKey  = key.New("ex.com/lemons")
	anotherKey = key.New("ex.com/another")

	// Note that metric instruments are declared globally.  They
	// are initialized when the global scope is set.
	exGauge = metric.NewFloat64Gauge("gauge.one",
		metric.WithKeys(fooKey, barKey, lemonsKey),
		metric.WithDescription("A gauge set to 1.0"),
	)

	exMeasure = metric.NewFloat64Measure("measure.two")
)

func initTracer() trace.TracerSDK {
	var err error
	exp, err := tracestdout.NewExporter(tracestdout.Options{PrettyPrint: false})
	if err != nil {
		log.Panicf("failed to initialize trace stdout exporter %v", err)
	}
	tr, err := sdktrace.NewTracer(sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		log.Panicf("failed to initialize trace provider %v", err)
	}
	return tr
}

func initMeter() *push.Controller {
	pusher, err := metricstdout.NewExportPipeline(metricstdout.Config{
		Quantiles:   []float64{0.5, 0.9, 0.99},
		PrettyPrint: false,
	})
	if err != nil {
		log.Panicf("failed to initialize metric stdout exporter %v", err)
	}
	return pusher
}

func initTelemetry() func() {
	tracer := initTracer()
	pusher := initMeter()
	global.SetScope(
		scope.WithTracerSDK(tracer).
			WithMeterSDK(pusher.Meter()).
			WithNamespace("example").
			AddResources(
				key.String("process1", "value1"),
				key.String("process2", "value2"),
			),
	)
	return pusher.Stop
}

func main() {
	defer initTelemetry()()

	// Use the global scope, provide a namespace & resources, get a base context.
	ctx := global.Scope().
		WithNamespace("ex.com/basic").
		AddResources(
			lemonsKey.Int(10),
			key.String("A", "1"),
			key.String("B", "2"),
			key.String("C", "3"),
		).
		InContext(context.Background())

	// Setup a distributed context
	ctx = distributedcontext.NewContext(
		ctx,
		fooKey.String("foo1"),
		barKey.String("bar1"),
	)

	// Binding in this context gets the process-wide labels and
	// the scoped labels entered above automatically.  One new
	// label is added at the call site for each bound instrument.
	gauge := exGauge.Bind(ctx, key.Float64("D", 1.3))
	defer gauge.Unbind()

	measure := exMeasure.Bind(ctx, key.Bool("E", false))
	defer measure.Unbind()

	// Using the static method `trace.WithSpan` here, it uses
	// the current scope's tracer this inherits the resource
	// scope.
	err := trace.WithSpan(ctx, "operation", func(ctx context.Context) error {
		span := trace.SpanFromContext(ctx)

		span.AddEvent(ctx, "Nice operation!", key.New("bogons").Int(100))

		span.SetAttributes(anotherKey.String("yes"))

		gauge.Set(ctx, 1)

		metric.RecordBatch(
			ctx,
			[]core.KeyValue{
				anotherKey.String("xyz"),
			},
			exGauge.Measurement(1.0),
			exMeasure.Measurement(2.0),
		)

		return trace.WithSpan(
			ctx,
			"Sub operation...",
			func(ctx context.Context) error {
				span := trace.SpanFromContext(ctx)
				span.SetAttributes(lemonsKey.String("five"))

				span.AddEvent(ctx, "Sub span event")

				measure.Record(ctx, 1.3)

				return nil
			},
		)
	})
	if err != nil {
		panic(err)
	}
}
