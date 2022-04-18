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

package prometheus_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

type expectedMetric struct {
	kind   string
	name   string
	help   string
	values []string
}

func (e *expectedMetric) lines() []string {
	ret := []string{
		fmt.Sprintf("# HELP %s %s", e.name, e.help),
		fmt.Sprintf("# TYPE %s %s", e.name, e.kind),
	}

	ret = append(ret, e.values...)

	return ret
}

func expectCounterWithHelp(name, help, value string) expectedMetric {
	return expectedMetric{
		kind:   "counter",
		name:   name,
		help:   help,
		values: []string{value},
	}
}

func expectCounter(name, value string) expectedMetric {
	return expectCounterWithHelp(name, "", value)
}

func expectGauge(name, value string) expectedMetric {
	return expectedMetric{
		kind:   "gauge",
		name:   name,
		values: []string{value},
	}
}

func expectHistogram(name string, values ...string) expectedMetric {
	return expectedMetric{
		kind:   "histogram",
		name:   name,
		values: values,
	}
}

func newPipeline(config prometheus.Config, options ...controller.Option) (*prometheus.Exporter, error) {
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		options...,
	)
	return prometheus.New(config, c)
}

func TestPrometheusExporter(t *testing.T) {
	exporter, err := newPipeline(
		prometheus.Config{
			DefaultHistogramBoundaries: []float64{-0.5, 1},
		},
		controller.WithCollectPeriod(0),
		controller.WithResource(resource.NewSchemaless(attribute.String("R", "V"))),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")
	upDownCounter, err := meter.SyncFloat64().UpDownCounter("updowncounter")
	require.NoError(t, err)
	counter, err := meter.SyncFloat64().Counter("counter")
	require.NoError(t, err)
	histogram, err := meter.SyncFloat64().Histogram("histogram")
	require.NoError(t, err)

	attrs := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}
	ctx := context.Background()

	var expected []expectedMetric

	counter.Add(ctx, 10, attrs...)
	counter.Add(ctx, 5.3, attrs...)

	expected = append(expected, expectCounter("counter", `counter{A="B",C="D",R="V"} 15.3`))

	gaugeObserver, err := meter.AsyncInt64().Gauge("intgaugeobserver")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{gaugeObserver}, func(ctx context.Context) {
		gaugeObserver.Observe(ctx, 1, attrs...)
	})
	require.NoError(t, err)

	expected = append(expected, expectGauge("intgaugeobserver", `intgaugeobserver{A="B",C="D",R="V"} 1`))

	histogram.Record(ctx, -0.6, attrs...)
	histogram.Record(ctx, -0.4, attrs...)
	histogram.Record(ctx, 0.6, attrs...)
	histogram.Record(ctx, 20, attrs...)

	expected = append(expected, expectHistogram("histogram",
		`histogram_bucket{A="B",C="D",R="V",le="-0.5"} 1`,
		`histogram_bucket{A="B",C="D",R="V",le="1"} 3`,
		`histogram_bucket{A="B",C="D",R="V",le="+Inf"} 4`,
		`histogram_sum{A="B",C="D",R="V"} 19.6`,
		`histogram_count{A="B",C="D",R="V"} 4`,
	))

	upDownCounter.Add(ctx, 10, attrs...)
	upDownCounter.Add(ctx, -3.2, attrs...)

	expected = append(expected, expectGauge("updowncounter", `updowncounter{A="B",C="D",R="V"} 6.8`))

	counterObserver, err := meter.AsyncFloat64().Counter("floatcounterobserver")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		counterObserver.Observe(ctx, 7.7, attrs...)
	})
	require.NoError(t, err)

	expected = append(expected, expectCounter("floatcounterobserver", `floatcounterobserver{A="B",C="D",R="V"} 7.7`))

	upDownCounterObserver, err := meter.AsyncFloat64().UpDownCounter("floatupdowncounterobserver")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{upDownCounterObserver}, func(ctx context.Context) {
		upDownCounterObserver.Observe(ctx, -7.7, attrs...)
	})
	require.NoError(t, err)

	expected = append(expected, expectGauge("floatupdowncounterobserver", `floatupdowncounterobserver{A="B",C="D",R="V"} -7.7`))

	compareExport(t, exporter, expected)
	compareExport(t, exporter, expected)
}

func compareExport(t *testing.T, exporter *prometheus.Exporter, expected []expectedMetric) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	exporter.ServeHTTP(rec, req)

	output := rec.Body.String()
	lines := strings.Split(output, "\n")

	expectedLines := []string{""}
	for _, v := range expected {
		expectedLines = append(expectedLines, v.lines()...)
	}

	sort.Strings(lines)
	sort.Strings(expectedLines)

	require.Equal(t, expectedLines, lines)
}

func TestPrometheusStatefulness(t *testing.T) {
	// Create a meter
	exporter, err := newPipeline(
		prometheus.Config{},
		controller.WithCollectPeriod(0),
		controller.WithResource(resource.Empty()),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")

	ctx := context.Background()

	counter, err := meter.SyncInt64().Counter("a.counter", instrument.WithDescription("Counts things"))
	require.NoError(t, err)

	counter.Add(ctx, 100, attribute.String("key", "value"))

	compareExport(t, exporter, []expectedMetric{
		expectCounterWithHelp("a_counter", "Counts things", `a_counter{key="value"} 100`),
	})

	counter.Add(ctx, 100, attribute.String("key", "value"))

	compareExport(t, exporter, []expectedMetric{
		expectCounterWithHelp("a_counter", "Counts things", `a_counter{key="value"} 200`),
	})
}
