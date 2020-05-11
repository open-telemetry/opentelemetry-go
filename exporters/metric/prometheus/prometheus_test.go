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
	"log"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/metric/test"
)

func TestPrometheusExporter(t *testing.T) {
	exporter, err := prometheus.NewRawExporter(prometheus.Config{
		DefaultSummaryQuantiles: []float64{0.5, 0.9, 0.99},
	})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	var expected []string
	checkpointSet := test.NewCheckpointSet()

	counter := metric.NewDescriptor(
		"counter", metric.CounterKind, metric.Float64NumberKind)
	lastValue := metric.NewDescriptor(
		"lastvalue", metric.ObserverKind, metric.Float64NumberKind)
	measure := metric.NewDescriptor(
		"measure", metric.MeasureKind, metric.Float64NumberKind)
	histogramMeasure := metric.NewDescriptor(
		"histogram_measure", metric.MeasureKind, metric.Float64NumberKind)

	labels := []core.KeyValue{
		key.New("A").String("B"),
		key.New("C").String("D"),
	}

	checkpointSet.AddCounter(&counter, 15.3, labels...)
	expected = append(expected, `counter{A="B",C="D"} 15.3`)

	checkpointSet.AddLastValue(&lastValue, 13.2, labels...)
	expected = append(expected, `lastvalue{A="B",C="D"} 13.2`)

	checkpointSet.AddMeasure(&measure, 13, labels...)
	checkpointSet.AddMeasure(&measure, 15, labels...)
	checkpointSet.AddMeasure(&measure, 17, labels...)
	expected = append(expected, `measure{A="B",C="D",quantile="0.5"} 15`)
	expected = append(expected, `measure{A="B",C="D",quantile="0.9"} 17`)
	expected = append(expected, `measure{A="B",C="D",quantile="0.99"} 17`)
	expected = append(expected, `measure_sum{A="B",C="D"} 45`)
	expected = append(expected, `measure_count{A="B",C="D"} 3`)

	boundaries := []metric.Number{metric.NewFloat64Number(-0.5), metric.NewFloat64Number(1)}
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, -0.6, labels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, -0.4, labels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, 0.6, labels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, 20, labels...)

	expected = append(expected, `histogram_measure_bucket{A="B",C="D",le="+Inf"} 4`)
	expected = append(expected, `histogram_measure_bucket{A="B",C="D",le="-0.5"} 1`)
	expected = append(expected, `histogram_measure_bucket{A="B",C="D",le="1"} 3`)
	expected = append(expected, `histogram_measure_count{A="B",C="D"} 4`)
	expected = append(expected, `histogram_measure_sum{A="B",C="D"} 19.6`)

	missingLabels := []core.KeyValue{
		key.New("A").String("E"),
		key.New("C").String(""),
	}

	checkpointSet.AddCounter(&counter, 12, missingLabels...)
	expected = append(expected, `counter{A="E",C=""} 12`)

	checkpointSet.AddLastValue(&lastValue, 32, missingLabels...)
	expected = append(expected, `lastvalue{A="E",C=""} 32`)

	checkpointSet.AddMeasure(&measure, 19, missingLabels...)
	expected = append(expected, `measure{A="E",C="",quantile="0.5"} 19`)
	expected = append(expected, `measure{A="E",C="",quantile="0.9"} 19`)
	expected = append(expected, `measure{A="E",C="",quantile="0.99"} 19`)
	expected = append(expected, `measure_count{A="E",C=""} 1`)
	expected = append(expected, `measure_sum{A="E",C=""} 19`)

	boundaries = []metric.Number{metric.NewFloat64Number(0), metric.NewFloat64Number(1)}
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, -0.6, missingLabels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, -0.4, missingLabels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, -0.1, missingLabels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, 15, missingLabels...)
	checkpointSet.AddHistogramMeasure(&histogramMeasure, boundaries, 15, missingLabels...)

	expected = append(expected, `histogram_measure_bucket{A="E",C="",le="+Inf"} 5`)
	expected = append(expected, `histogram_measure_bucket{A="E",C="",le="0"} 3`)
	expected = append(expected, `histogram_measure_bucket{A="E",C="",le="1"} 3`)
	expected = append(expected, `histogram_measure_count{A="E",C=""} 5`)
	expected = append(expected, `histogram_measure_sum{A="E",C=""} 28.9`)

	compareExport(t, exporter, checkpointSet, expected)
}

func compareExport(t *testing.T, exporter *prometheus.Exporter, checkpointSet *test.CheckpointSet, expected []string) {
	err := exporter.Export(context.Background(), nil, checkpointSet)
	require.Nil(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	exporter.ServeHTTP(rec, req)

	output := rec.Body.String()
	lines := strings.Split(output, "\n")
	var metricsOnly []string
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && line != "" {
			metricsOnly = append(metricsOnly, line)
		}
	}

	sort.Strings(metricsOnly)
	sort.Strings(expected)

	require.Equal(t, strings.Join(expected, "\n"), strings.Join(metricsOnly, "\n"))
}
