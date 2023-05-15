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

package metric_test

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func Example() {
	// This reader is used as a stand-in for a reader that will actually export
	// data. See exporters in the go.opentelemetry.io/otel/exporters package
	// for more information.
	reader := metric.NewManualReader()

	// See the go.opentelemetry.io/otel/sdk/resource package for more
	// information about how to create and use Resources.
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-service"),
		semconv.ServiceVersion("v0.1.0"),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(reader),
	)
	otel.SetMeterProvider(meterProvider)
	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()
	// The MeterProvider is configured and registered globally. You can now run
	// your code instrumented with the OpenTelemetry API that uses the global
	// MeterProvider without having to pass this MeterProvider instance. Or,
	// you can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.
	//
	// See the go.opentelemetry.io/otel/metric package for more information
	// about the metric API.
}
