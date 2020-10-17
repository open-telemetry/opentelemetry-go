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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/example/namedtracer/foo"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/global"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	fooKey     = label.Key("ex.com/foo")
	barKey     = label.Key("ex.com/bar")
	anotherKey = label.Key("ex.com/another")
)

var tp *sdktrace.TracerProvider

// initTracer creates and registers trace provider instance.
func initTracer() func() {
	var err error
	exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Panicf("failed to initialize stdout exporter %v\n", err)
		return nil
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithConfig(
			sdktrace.Config{
				DefaultSampler: sdktrace.AlwaysSample(),
			},
		),
		sdktrace.WithSpanProcessor(bsp),
	)
	global.SetTracerProvider(tp)
	return bsp.Shutdown
}

func main() {
	// initialize trace provider.
	shutdown := initTracer()
	defer shutdown()

	// Create a named tracer with package path as its name.
	tracer := tp.Tracer("example/namedtracer/main")
	ctx := context.Background()
	ctx = otel.ContextWithBaggageValues(ctx, fooKey.String("foo1"), barKey.String("bar1"))

	var span otel.Span
	ctx, span = tracer.Start(ctx, "operation")
	defer span.End()
	span.AddEvent("Nice operation!", otel.WithAttributes(label.Int("bogons", 100)))
	span.SetAttributes(anotherKey.String("yes"))
	if err := foo.SubOperation(ctx); err != nil {
		panic(err)
	}
}
