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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/metric/controller/pull"
	"go.opentelemetry.io/otel/sdk/resource"
)

// This test demonstrates that it is relatively difficult to setup a
// Prometheus export pipeline:
//
//   1. The default boundaries are difficult to pass, should be []float instead of []otel.Number
//
// TODO: Address this issue.

func ExampleNewExportPipeline() {
	// Create a resource, with builtin attributes plus R=V.
	res, err := resource.New(
		context.Background(),
		resource.WithoutBuiltin(), // Test-only!
		resource.WithAttributes(label.String("R", "V")),
	)
	if err != nil {
		panic(err)
	}

	// Create a meter
	exporter, err := prometheus.NewExportPipeline(
		prometheus.Config{},
		pull.WithResource(res),
	)
	if err != nil {
		panic(err)
	}
	meter := exporter.MeterProvider().Meter("example")
	ctx := context.Background()

	// Use two instruments
	counter := otel.Must(meter).NewInt64Counter(
		"a.counter",
		otel.WithDescription("Counts things"),
	)
	recorder := otel.Must(meter).NewInt64ValueRecorder(
		"a.valuerecorder",
		otel.WithDescription("Records values"),
	)

	counter.Add(ctx, 100, label.String("key", "value"))
	recorder.Record(ctx, 100, label.String("key", "value"))

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
	// a_counter{R="V",key="value"} 100
	// # HELP a_valuerecorder Records values
	// # TYPE a_valuerecorder histogram
	// a_valuerecorder_bucket{R="V",key="value",le="+Inf"} 1
	// a_valuerecorder_sum{R="V",key="value"} 100
	// a_valuerecorder_count{R="V",key="value"} 1
}
