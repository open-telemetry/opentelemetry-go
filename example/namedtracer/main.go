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
	"time"

	"go.opentelemetry.io/api/distributedcontext"
	"go.opentelemetry.io/api/key"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/example/namedtracer/foo"
	sdktrace "go.opentelemetry.io/sdk/trace"
)

var (
	fooKey     = key.New("ex.com/foo")
	barKey     = key.New("ex.com/bar")
	anotherKey = key.New("ex.com/another")
)

// initTracer registers sdktrace as trace provider. It also registers exporter with
// sdktrace. In this example it is PrintExporter. Any default configuration such as
// default sampling should be done here.
func initTracer() {
	sdktrace.RegisterProvider()
	ssp := sdktrace.NewSimpleSpanProcessor(&sdktrace.PrintExporter{})
	sdktrace.RegisterSpanProcessor(ssp)
	sdktrace.ApplyConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()})
}

func main() {
	// initialize trace provider.
	initTracer()

	// Create a named tracer with package path as its name.
	tracer := trace.GlobalProvider().Tracer("example/namedtracer/main")
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

	// sleep 5 seconds to print span data
	time.Sleep(5 * time.Second)
}
