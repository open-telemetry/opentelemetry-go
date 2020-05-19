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
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	exportTest "go.opentelemetry.io/otel/exporters/metric/test"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	controllerTest "go.opentelemetry.io/otel/sdk/metric/controller/test"
)

func TestPrometheusExporter(t *testing.T) {
	exporter, err := prometheus.NewRawExporter(prometheus.Config{
		DefaultSummaryQuantiles: []float64{0.5, 0.9, 0.99},
	})
	require.NoError(t, err)

	var expected []string
	checkpointSet := exportTest.NewCheckpointSet(nil)

	counter := metric.NewDescriptor(
		"counter", metric.CounterKind, metric.Float64NumberKind)
	lastValue := metric.NewDescriptor(
		"lastvalue", metric.ValueObserverKind, metric.Float64NumberKind)
	valuerecorder := metric.NewDescriptor(
		"valuerecorder", metric.ValueRecorderKind, metric.Float64NumberKind)
	histogramValueRecorder := metric.NewDescriptor(
		"histogram_valuerecorder", metric.ValueRecorderKind, metric.Float64NumberKind)

	labels := []kv.KeyValue{
		kv.Key("A").String("B"),
		kv.Key("C").String("D"),
	}

	checkpointSet.AddCounter(&counter, 15.3, labels...)
	expected = append(expected, `counter{A="B",C="D"} 15.3`)

	checkpointSet.AddLastValue(&lastValue, 13.2, labels...)
	expected = append(expected, `lastvalue{A="B",C="D"} 13.2`)

	checkpointSet.AddValueRecorder(&valuerecorder, 13, labels...)
	checkpointSet.AddValueRecorder(&valuerecorder, 15, labels...)
	checkpointSet.AddValueRecorder(&valuerecorder, 17, labels...)
	expected = append(expected, `valuerecorder{A="B",C="D",quantile="0.5"} 15`)
	expected = append(expected, `valuerecorder{A="B",C="D",quantile="0.9"} 17`)
	expected = append(expected, `valuerecorder{A="B",C="D",quantile="0.99"} 17`)
	expected = append(expected, `valuerecorder_sum{A="B",C="D"} 45`)
	expected = append(expected, `valuerecorder_count{A="B",C="D"} 3`)

	boundaries := []metric.Number{metric.NewFloat64Number(-0.5), metric.NewFloat64Number(1)}
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, -0.6, labels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, -0.4, labels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, 0.6, labels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, 20, labels...)

	expected = append(expected, `histogram_valuerecorder_bucket{A="B",C="D",le="+Inf"} 4`)
	expected = append(expected, `histogram_valuerecorder_bucket{A="B",C="D",le="-0.5"} 1`)
	expected = append(expected, `histogram_valuerecorder_bucket{A="B",C="D",le="1"} 3`)
	expected = append(expected, `histogram_valuerecorder_count{A="B",C="D"} 4`)
	expected = append(expected, `histogram_valuerecorder_sum{A="B",C="D"} 19.6`)

	missingLabels := []kv.KeyValue{
		kv.Key("A").String("E"),
		kv.Key("C").String(""),
	}

	checkpointSet.AddCounter(&counter, 12, missingLabels...)
	expected = append(expected, `counter{A="E",C=""} 12`)

	checkpointSet.AddLastValue(&lastValue, 32, missingLabels...)
	expected = append(expected, `lastvalue{A="E",C=""} 32`)

	checkpointSet.AddValueRecorder(&valuerecorder, 19, missingLabels...)
	expected = append(expected, `valuerecorder{A="E",C="",quantile="0.5"} 19`)
	expected = append(expected, `valuerecorder{A="E",C="",quantile="0.9"} 19`)
	expected = append(expected, `valuerecorder{A="E",C="",quantile="0.99"} 19`)
	expected = append(expected, `valuerecorder_count{A="E",C=""} 1`)
	expected = append(expected, `valuerecorder_sum{A="E",C=""} 19`)

	boundaries = []metric.Number{metric.NewFloat64Number(0), metric.NewFloat64Number(1)}
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, -0.6, missingLabels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, -0.4, missingLabels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, -0.1, missingLabels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, 15, missingLabels...)
	checkpointSet.AddHistogramValueRecorder(&histogramValueRecorder, boundaries, 15, missingLabels...)

	expected = append(expected, `histogram_valuerecorder_bucket{A="E",C="",le="+Inf"} 5`)
	expected = append(expected, `histogram_valuerecorder_bucket{A="E",C="",le="0"} 3`)
	expected = append(expected, `histogram_valuerecorder_bucket{A="E",C="",le="1"} 3`)
	expected = append(expected, `histogram_valuerecorder_count{A="E",C=""} 5`)
	expected = append(expected, `histogram_valuerecorder_sum{A="E",C=""} 28.9`)

	compareExport(t, exporter, checkpointSet, expected)
}

func compareExport(t *testing.T, exporter *prometheus.Exporter, checkpointSet *exportTest.CheckpointSet, expected []string) {
	err := exporter.Export(context.Background(), checkpointSet)
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

func TestPrometheusStatefulness(t *testing.T) {
	// Create a meter
	controller, exporter, err := prometheus.NewExportPipeline(prometheus.Config{}, push.WithPeriod(time.Minute))
	require.NoError(t, err)

	meter := controller.Provider().Meter("test")
	mock := controllerTest.NewMockClock()
	controller.SetClock(mock)
	controller.Start()

	// GET the HTTP endpoint
	scrape := func() string {
		var input bytes.Buffer
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", &input)
		require.NoError(t, err)

		exporter.ServeHTTP(resp, req)
		data, err := ioutil.ReadAll(resp.Result().Body)
		require.NoError(t, err)

		return string(data)
	}

	ctx := context.Background()

	counter := metric.Must(meter).NewInt64Counter(
		"a.counter",
		metric.WithDescription("Counts things"),
	)

	counter.Add(ctx, 100, kv.String("key", "value"))

	// Trigger a push
	mock.Add(time.Minute)
	runtime.Gosched()

	require.Equal(t, `# HELP a_counter Counts things
# TYPE a_counter counter
a_counter{key="value"} 100
`, scrape())

	counter.Add(ctx, 100, kv.String("key", "value"))

	// Again, now expect cumulative count
	mock.Add(time.Minute)
	runtime.Gosched()

	require.Equal(t, `# HELP a_counter Counts things
# TYPE a_counter counter
a_counter{key="value"} 200
`, scrape())

}
