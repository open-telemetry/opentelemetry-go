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

package otlp_test

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/internal/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func ExampleNew() {
	ctx := context.Background()

	metricsDriver := otlpgrpc.NewDriver(
	// Configure metrics driver here
	)
	driver := otlp.NewSplitDriver(otlp.WithMetricDriver(metricsDriver))
	exporter, err := otlp.New(ctx, driver) // Configure as needed.
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}
	defer func() {
		err := exporter.Shutdown(ctx)
		if err != nil {
			log.Fatalf("failed to stop exporter: %v", err)
		}
	}()

	c := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exporter),
		controller.WithExporter(exporter))

	global.SetMeterProvider(c.MeterProvider())
}
