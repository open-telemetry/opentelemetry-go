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
	"fmt"
	"log"
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	prombridge "go.opentelemetry.io/otel/bridge/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
)

// promHistogram is a histogram defined using the prometheus client library.
var promHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "random_numbers",
	Help:    "A histogram of normally distributed random numbers.",
	Buckets: prometheus.LinearBuckets(-3, .1, 61),
})

func main() {
	// Create an OpenTelemetry exporter.  Use the stdout exporter for our
	// example, but you could use OTLP or another exporter instead.
	otelExporter, err := stdoutmetric.New()
	if err != nil {
		log.Fatal(fmt.Errorf("error creating metric exporter: %w", err))
	}
	// Construct a reader which periodically exports to our exporter.
	reader := metric.NewPeriodicReader(otelExporter)
	// Register the Prometheus metric Producer to add metrics from the
	// Prometheus DefaultGatherer to the output.
	reader.RegisterProducer(prombridge.NewMetricProducer())
	// Create an OTel MeterProvider that periodically reads from Prometheus,
	// but don't use it to create any meters or instruments. We will use
	// Prometheus instruments in this example instead.
	metric.NewMeterProvider(metric.WithReader(reader))
	// Make observations using our Prometheus histogram.
	for {
		promHistogram.Observe(rand.NormFloat64())
	}
	// We should see our histogram in the output, as well as the metrics
	// registered by prometheus by default, which includes go runtime metrics
	// and process metrics.
}
