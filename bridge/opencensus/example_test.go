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

package opencensus_test

import (
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/sdk/metric"
)

func ExampleNewMetricProducer() {
	// Create the OpenCensus Metric bridge.
	bridge := opencensus.NewMetricProducer()
	// Add the bridge as a producer to your reader.
	// If using a push exporter, such as OTLP exporter,
	// use metric.NewPeriodicReader with metric.WithProducer option.
	// If using a pull exporter which acts as a reader, such as prometheus exporter,
	// use a dedicated option like prometheus.WithProducer.
	reader := metric.NewManualReader(metric.WithProducer(bridge))
	// Add the reader to your MeterProvider.
	_ = metric.NewMeterProvider(metric.WithReader(reader))
}
