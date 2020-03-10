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
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

func TestPrometheusExporter(t *testing.T) {
	exporter, err := prometheus.NewRawExporter(prometheus.Config{
		DefaultSummaryQuantiles: []float64{0.5, 0.9, 0.99},
	})
	if err != nil {
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	var expected []string
	checkpointSet := test.NewCheckpointSet(metric.NewDefaultLabelEncoder())

	counter := export.NewDescriptor(
		"counter", export.CounterKind, nil, "", "", core.Float64NumberKind, false)
	lastValue := export.NewDescriptor(
		"lastvalue", export.ObserverKind, nil, "", "", core.Float64NumberKind, false)
	measure := export.NewDescriptor(
		"measure", export.MeasureKind, nil, "", "", core.Float64NumberKind, false)

	labels := []core.KeyValue{
		key.New("A").String("B"),
		key.New("C").String("D"),
	}

	checkpointSet.AddCounter(counter, 15.3, labels...)
	expected = append(expected, `counter{A="B",C="D"} 15.3`)

	checkpointSet.AddLastValue(lastValue, 13.2, labels...)
	expected = append(expected, `lastvalue{A="B",C="D"} 13.2`)

	checkpointSet.AddMeasure(measure, 13, labels...)
	checkpointSet.AddMeasure(measure, 15, labels...)
	checkpointSet.AddMeasure(measure, 17, labels...)
	expected = append(expected, `measure{A="B",C="D",quantile="0.5"} 15`)
	expected = append(expected, `measure{A="B",C="D",quantile="0.9"} 17`)
	expected = append(expected, `measure{A="B",C="D",quantile="0.99"} 17`)
	expected = append(expected, `measure_sum{A="B",C="D"} 45`)
	expected = append(expected, `measure_count{A="B",C="D"} 3`)

	missingLabels := []core.KeyValue{
		key.New("A").String("E"),
		key.New("C").String(""),
	}

	checkpointSet.AddCounter(counter, 12, missingLabels...)
	expected = append(expected, `counter{A="E",C=""} 12`)

	checkpointSet.AddLastValue(lastValue, 32, missingLabels...)
	expected = append(expected, `lastvalue{A="E",C=""} 32`)

	checkpointSet.AddMeasure(measure, 19, missingLabels...)
	expected = append(expected, `measure{A="E",C="",quantile="0.5"} 19`)
	expected = append(expected, `measure{A="E",C="",quantile="0.9"} 19`)
	expected = append(expected, `measure{A="E",C="",quantile="0.99"} 19`)
	expected = append(expected, `measure_count{A="E",C=""} 1`)
	expected = append(expected, `measure_sum{A="E",C=""} 19`)

	compareExport(t, exporter, checkpointSet, expected)
}

func compareExport(t *testing.T, exporter *prometheus.Exporter, checkpointSet *test.CheckpointSet, expected []string) {
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
