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

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/stdout"
)

const (
	instrumentationName    = "github.com/instrumentron"
	instrumentationVersion = "v0.1.0"
)

var (
	tracer = global.TraceProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
	)

	meter = global.MeterProvider().Meter(
		instrumentationName,
		metric.WithInstrumentationVersion(instrumentationVersion),
	)

	loopCounter = metric.Must(meter).NewInt64Counter("function.loops")
	paramValue  = metric.Must(meter).NewFloat64ValueRecorder("function.param")

	nameKey = kv.Key("function.name")
)

func myFunction(ctx context.Context, values ...float64) error {
	nameKV := nameKey.String("myFunction")
	boundCount := loopCounter.Bind(nameKV)
	boundValue := paramValue.Bind(nameKV)
	for _, value := range values {
		boundCount.Add(ctx, 1)
		boundValue.Record(ctx, value)
	}
	return nil
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

	err = tracer.WithSpan(
		context.Background(),
		"myFunction/call",
		func(ctx context.Context) error {
			err := tracer.WithSpan(
				ctx,
				"internal/call",
				func(ctx context.Context) error { return myFunction(ctx, 200, 100, 5000, 600) },
			)
			if err != nil {
				return err
			}
			return myFunction(ctx, 100, 200, 500, 800)
		},
	)
	if err != nil {
		log.Fatal("Failed to call myFunction", err)
	}
}
