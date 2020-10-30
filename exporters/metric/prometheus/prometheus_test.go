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
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestPrometheusExporter(t *testing.T) {
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{
			DefaultHistogramBoundaries: []float64{-0.5, 1},
		},
		pull.WithCachePeriod(0),
		pull.WithResource(resource.NewWithAttributes(label.String("R", "V"))),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")

	counter := otel.Must(meter).NewFloat64Counter("counter")
	valuerecorder := otel.Must(meter).NewFloat64ValueRecorder("valuerecorder")

	labels := []label.KeyValue{
		label.Key("A").String("B"),
		label.Key("C").String("D"),
	}
	ctx := context.Background()

	var expected []string

	counter.Add(ctx, 10, labels...)
	counter.Add(ctx, 5.3, labels...)

	expected = append(expected, `counter{A="B",C="D",R="V"} 15.3`)

	valuerecorder.Record(ctx, -0.6, labels...)
	valuerecorder.Record(ctx, -0.4, labels...)
	valuerecorder.Record(ctx, 0.6, labels...)
	valuerecorder.Record(ctx, 20, labels...)

	expected = append(expected, `valuerecorder_bucket{A="B",C="D",R="V",le="+Inf"} 4`)
	expected = append(expected, `valuerecorder_bucket{A="B",C="D",R="V",le="-0.5"} 1`)
	expected = append(expected, `valuerecorder_bucket{A="B",C="D",R="V",le="1"} 3`)
	expected = append(expected, `valuerecorder_count{A="B",C="D",R="V"} 4`)
	expected = append(expected, `valuerecorder_sum{A="B",C="D",R="V"} 19.6`)

	compareExport(t, exporter, expected)
	compareExport(t, exporter, expected)
}

func compareExport(t *testing.T, exporter *prometheus.Exporter, expected []string) {
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
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
		pull.WithCachePeriod(0),
	)
	require.NoError(t, err)

	meter := exporter.MeterProvider().Meter("test")

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

	counter := otel.Must(meter).NewInt64Counter(
		"a.counter",
		otel.WithDescription("Counts things"),
	)

	counter.Add(ctx, 100, label.String("key", "value"))

	require.Equal(t, `# HELP a_counter Counts things
# TYPE a_counter counter
a_counter{key="value"} 100
`, scrape())

	counter.Add(ctx, 100, label.String("key", "value"))

	require.Equal(t, `# HELP a_counter Counts things
# TYPE a_counter counter
a_counter{key="value"} 200
`, scrape())

}
