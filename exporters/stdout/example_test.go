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

package stdout_test

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

var (
	tracer = global.TracerProvider().Tracer(
		instrumentationName,
		otel.WithInstrumentationVersion(instrumentationVersion),
	)

	meter = global.MeterProvider().Meter(
		instrumentationName,
		otel.WithInstrumentationVersion(instrumentationVersion),
	)

	loopCounter = otel.Must(meter).NewInt64Counter("function.loops")
	paramValue  = otel.Must(meter).NewInt64ValueRecorder("function.param")

	nameKey = label.Key("function.name")
)

func add(ctx context.Context, x, y int64) int64 {
	nameKV := nameKey.String("add")

	var span otel.Span
	ctx, span = tracer.Start(ctx, "Addition")
	defer span.End()

	loopCounter.Add(ctx, 1, nameKV)
	paramValue.Record(ctx, x, nameKV)
	paramValue.Record(ctx, y, nameKV)

	return x + y
}

func multiply(ctx context.Context, x, y int64) int64 {
	nameKV := nameKey.String("multiply")

	var span otel.Span
	ctx, span = tracer.Start(ctx, "Multiplication")
	defer span.End()

	loopCounter.Add(ctx, 1, nameKV)
	paramValue.Record(ctx, x, nameKV)
	paramValue.Record(ctx, y, nameKV)

	return x * y
}

func Example() {
	exportOpts := []stdout.Option{
		stdout.WithQuantiles([]float64{0.5}),
		stdout.WithPrettyPrint(),
	}
	// Registers both a trace and meter Provider globally.
	pusher, err := stdout.InstallNewPipeline(exportOpts, nil)
	if err != nil {
		log.Fatal("Could not initialize stdout exporter:", err)
	}
	defer pusher.Stop()

	ctx := context.Background()
	log.Println("the answer is", add(ctx, multiply(ctx, multiply(ctx, 2, 2), 10), 2))
}
