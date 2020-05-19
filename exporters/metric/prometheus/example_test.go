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

// This test demonstrates that it is relatively difficult to setup a
// Prometheus export pipeline:
//
//   1. The default boundaries are difficult to pass, should be []float instead of []metric.Number
//   2. The push controller doesn't make sense b/c Prometheus is pull-bsaed
//
// TODO: Address these issues; add Resources to the test.

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
	counter := metric.Must(meter).NewInt64Counter(
		"a.counter",
		metric.WithDescription("Counts things"),
	)
	recorder := metric.Must(meter).NewInt64ValueRecorder(
		"a.valuerecorder",
		metric.WithDescription("Records values"),
	)

	counter.Add(ctx, 100, kv.String("key", "value"))
	recorder.Record(ctx, 100, kv.String("key", "value"))

	// Simulate a push
	meterImpl.Collect(ctx)
	err = exporter.Export(ctx, integrator.CheckpointSet())
	if err != nil {
		panic(err)
	}

	// GET the HTTP endpoint
	var input bytes.Buffer
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", &input)
	if err != nil {
		panic(err)
	}
	exporter.ServeHTTP(resp, req)
	data, err := ioutil.ReadAll(resp.Result().Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(data))

	// Output:
	// # HELP a_counter Counts things
	// # TYPE a_counter counter
	// a_counter{key="value"} 100
	// # HELP a_valuerecorder Records values
	// # TYPE a_valuerecorder histogram
	// a_valuerecorder_bucket{key="value",le="+Inf"} 1
	// a_valuerecorder_sum{key="value"} 100
	// a_valuerecorder_count{key="value"} 1
}
