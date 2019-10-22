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

	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/example/namedtracer/foo"
	"go.opentelemetry.io/exporter/trace/stdout"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

var (
	fooKey     = key.New("ex.com/foo")
	barKey     = key.New("ex.com/bar")
	anotherKey = key.New("ex.com/another")
)

var tp *sdktrace.Provider

// initTracer creates and registers trace provider instance.
func initTracer() {
	var err error
	exp, err := stdout.NewExporter(stdout.Options{})
	if err != nil {
		log.Panicf("failed to initialize stdout exporter %v\n", err)
		return
	}
	tp, err = sdktrace.NewProvider(sdktrace.WithSyncer(exp),
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}))
	if err != nil {
		log.Panicf("failed to initialize trace provider %v\n", err)
	}
	trace.SetGlobalProvider(tp)
}

func main() {
	// initialize trace provider.
	initTracer()

	// Create a named tracer with package path as its name.
	tracer := tp.GetTracer("example/namedtracer/main")
	ctx := context.Background()

	ctx = distributedcontext.NewContext(ctx,
		distributedcontext.Insert(fooKey.String("foo1")),
		distributedcontext.Insert(barKey.String("bar1")),
	)

	err := tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {

		trace.CurrentSpan(ctx).AddEvent(ctx, "Nice operation!", key.New("bogons").Int(100))

		trace.CurrentSpan(ctx).SetAttributes(anotherKey.String("yes"))

		return foo.SubOperation(ctx)
	})
	if err != nil {
		panic(err)
	}
}
