package prometheus_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	integrator "go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func ExampleNewExportPipeline() {
	// Create a meter
	selector := simple.NewWithHistogramDistribution(nil)
	exporter, err := prometheus.NewRawExporter(prometheus.Config{})
	if err != nil {
		panic(err)
	}
	integrator := integrator.New(selector, true)
	meterImpl := sdk.NewAccumulator(integrator)
	meter := metric.WrapMeterImpl(meterImpl, "example")

	ctx := context.Background()

	// Use two instruments
	counter := metric.Must(meter).NewInt64Counter("a.counter")
	recorder := metric.Must(meter).NewInt64ValueRecorder("a.valuerecorder")

	counter.Add(ctx, 100, kv.String("key", "value"))
	recorder.Record(ctx, 100, kv.String("key", "value"))

	// Simulate a push
	meterImpl.Collect(ctx)
	exporter.Export(ctx, nil, integrator.CheckpointSet())

	// GET the HTTP endpoint
	var input bytes.Buffer
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", &input)
	if err != nil {
		panic(err)
	}
	exporter.ServeHTTP(resp, req)
	data, err := ioutil.ReadAll(resp.Result().Body)
	fmt.Print(string(data))

	// Output:
	// # HELP a_counter
	// # TYPE a_counter counter
	// a_counter{key="value"} 100
	// # HELP a_valuerecorder
	// # TYPE a_valuerecorder histogram
	// a_valuerecorder_bucket{key="value",le="+Inf"} 1
	// a_valuerecorder_sum{key="value"} 100
	// a_valuerecorder_count{key="value"} 1
}
