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
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
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

func TestPrometheusExporter(t *testing.T) {
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{
			DefaultHistogramBoundaries: []float64{-0.5, 1},
		},
		controller.WithCollectPeriod(0),
		controller.WithResource(resource.NewSchemaless(attribute.String("R", "V"))),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")
	upDownCounter := metric.Must(meter).NewFloat64UpDownCounter("updowncounter")
	counter := metric.Must(meter).NewFloat64Counter("counter")
	valuerecorder := metric.Must(meter).NewFloat64ValueRecorder("valuerecorder")

	labels := []attribute.KeyValue{
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	}
	ctx := context.Background()

	var expected []expectedMetric

	counter.Add(ctx, 10, labels...)
	counter.Add(ctx, 5.3, labels...)

	expected = append(expected, expectCounter("counter", `counter{A="B",C="D",R="V"} 15.3`))

	_ = metric.Must(meter).NewInt64ValueObserver("intobserver", func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(1, labels...)
	})

	expected = append(expected, expectGauge("intobserver", `intobserver{A="B",C="D",R="V"} 1`))

	valuerecorder.Record(ctx, -0.6, labels...)
	valuerecorder.Record(ctx, -0.4, labels...)
	valuerecorder.Record(ctx, 0.6, labels...)
	valuerecorder.Record(ctx, 20, labels...)

	expected = append(expected, expectHistogram("valuerecorder",
		`valuerecorder_bucket{A="B",C="D",R="V",le="-0.5"} 1`,
		`valuerecorder_bucket{A="B",C="D",R="V",le="1"} 3`,
		`valuerecorder_bucket{A="B",C="D",R="V",le="+Inf"} 4`,
		`valuerecorder_sum{A="B",C="D",R="V"} 19.6`,
		`valuerecorder_count{A="B",C="D",R="V"} 4`,
	))

	upDownCounter.Add(ctx, 10, labels...)
	upDownCounter.Add(ctx, -3.2, labels...)

	expected = append(expected, expectGauge("updowncounter", `updowncounter{A="B",C="D",R="V"} 6.8`))

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
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
		controller.WithCollectPeriod(0),
		controller.WithResource(resource.Empty()),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")

	ctx := context.Background()

	counter := metric.Must(meter).NewInt64Counter(
		"a.counter",
		metric.WithDescription("Counts things"),
	)

	counter.Add(ctx, 100, attribute.String("key", "value"))

	compareExport(t, exporter, []expectedMetric{
		expectCounterWithHelp("a_counter", "Counts things", `a_counter{key="value"} 100`),
	})

	counter.Add(ctx, 100, attribute.String("key", "value"))

	compareExport(t, exporter, []expectedMetric{
		expectCounterWithHelp("a_counter", "Counts things", `a_counter{key="value"} 200`),
	})
}
