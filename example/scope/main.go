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
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	metricstdout "go.opentelemetry.io/otel/exporter/metric/stdout"
	tracestdout "go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	namespace = "ex.com"
)

var (
	environmentKey1 = key.New("environment1")
	environmentKey2 = key.New("environment2")

	resourceKey1 = key.New("resource1")
	resourceKey2 = key.New("resource2")

	attrKey1 = key.New("attribute1")
	attrKey2 = key.New("attribute2")

	// Note: metrics are allocated statically.  They use the
	// global scope's namespace when it is initialized.
	counter1 = metric.NewFloat64Counter(
		"counter1",
		metric.WithKeys(attrKey1, attrKey2),
	)
	gauge1 = metric.NewFloat64Gauge(
		"gauge1",
		metric.WithKeys(attrKey1, attrKey2),
	)
)

// start sets the global scope with the configured tracer, meter, and resources.
func start() func() {
	tracer := initTracer()
	meter := initMeter()

	telemetry := scope.
		WithTracerSDK(tracer).
		WithMeterSDK(meter.Meter()).
		WithNamespace(namespace).
		AddResources(
			environmentKey1.String("ENV1"),
			environmentKey2.String("ENV2"),
		)

	global.SetScope(telemetry)

	return func() {
		meter.Stop()
	}
}

func main() {
	defer start()()

	// Start with no telemetry state
	ctx := context.Background()

	// Add scoped resources.  These are on top of the global resources.
	ctx = scope.Current(ctx).AddResources(
		resourceKey1.String("res1"),
		resourceKey2.String("res2"),
	).InContext(ctx)

	// Now consider four ways to add "attrKey1" and "attrKey2" attributes
	// to a pair of metric events.

	////////////////////////////////////////////////////////////
	// 1 As a batch, labels passed at the call site

	// Using the Meter() from a scope ensures that scope's
	// resources are attached.
	scope.Current(ctx).Meter().RecordBatch(ctx, []core.KeyValue{
		attrKey1.String("val1"),
		attrKey2.String("val2"),
	},
		counter1.Measurement(1),
		gauge1.Measurement(2),
	)
	////////////////////////////////////////////////////////////
	// 2 Individual events, labels passed a the call site

	// The batch could be written as two events:
	counter1.Add(ctx, 1, attrKey1.String("val1"), attrKey2.String("val2"))
	gauge1.Set(ctx, 2, attrKey1.String("val1"), attrKey2.String("val2"))

	////////////////////////////////////////////////////////////
	// 3 By placing the labels in the current resource scope

	// Instead of repeating the two attributes above, and where
	// LabelSets are currently specified, use scope to introduce local resources:
	if true {
		ctx := scope.Current(ctx).AddResources(
			attrKey1.String("val1"),
			attrKey2.String("val2"),
		).InContext(ctx)

		// Now the "LabelSet" is part of the resource scope.
		counter1.Add(ctx, 1)
		gauge1.Set(ctx, 2)
	}

	////////////////////////////////////////////////////////////
	// 4 By starting a span with corresponding attributes, which
	// enter the current resource scope.

	// Creating a new span updates the scope with the span
	// attributes as resources.
	ctx, span := trace.Start(
		ctx,
		"a_span",
		trace.WithAttributes(
			attrKey1.String("val1"),
			attrKey2.String("val2"),
		),
	)
	defer span.End()

	// These metric events automatically have the current scope's resources.
	counter1.Add(ctx, 1)
	gauge1.Set(ctx, 2)
}

// initMeter configures the tracing SDK.
func initTracer() trace.TracerSDK {
	var err error
	exp, err := tracestdout.NewExporter(tracestdout.Options{PrettyPrint: false})
	if err != nil {
		log.Panicf("failed to initialize trace stdout exporter %v", err)
		return nil

	}
	tri, err := sdktrace.NewTracer(sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		log.Panicf("failed to initialize trace provider %v", err)
	}
	return tri
}

// initMeter configures the metrics SDK.
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
