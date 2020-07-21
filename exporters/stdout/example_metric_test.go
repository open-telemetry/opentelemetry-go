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

package stdout_test

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/stdout"
)

func ExampleExport() {
	// Default output is STDOUT.
	buf := &bytes.Buffer{}
	exportOpts := []stdout.Option{
		stdout.WithPrettyPrint(),
		// Used in testing to make output predictable.
		stdout.WithoutTimestamps(),
		stdout.WithWriter(buf),
	}
	_, pusher, err := stdout.NewExportPipeline(exportOpts, nil)
	if err != nil {
		log.Fatal("Could not initialize stdout exporter:", err)
	}

	meter := pusher.Provider().Meter(
		"github.com/instrumentron",
		metric.WithInstrumentationVersion("v0.1.0"),
	)

	// Create a counter.
	counter, err := meter.NewInt64Counter("a.counter")
	if err != nil {
		log.Fatal("Could not initialize a.counter:", err)
	}

	// Update the counter.
	ctx := context.Background()
	counter.Add(ctx, 100, kv.String("key", "value"))

	// Flush everything.
	pusher.Stop()

	fmt.Println(buf.String())
	// Output:
	// {
	// 	"updates": [
	// 		{
	// 			"name": "a.counter{instrumentation.name=github.com/instrumentron,instrumentation.version=v0.1.0,key=value}",
	// 			"sum": 100
	// 		}
	// 	]
	// }
}
