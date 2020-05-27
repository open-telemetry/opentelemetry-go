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
	"context"
	"log"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
)

func ExampleNewExportPipeline() {
	// Create a meter
	pusher, err := stdout.NewExportPipeline(stdout.Config{
		PrettyPrint:    true,
		DoNotPrintTime: true,
	})
	if err != nil {
		log.Fatal("Could not initialize stdout exporter:", err)
	}
	defer pusher.Stop()

	ctx := context.Background()

	key := kv.Key("key")
	meter := pusher.Provider().Meter("example")

	// Create and update a single counter:
	counter := metric.Must(meter).NewInt64Counter("a.counter")
	labels := []kv.KeyValue{key.String("value")}

	counter.Add(ctx, 100, labels...)

	// Output:
	// {
	// 	"updates": [
	// 		{
	// 			"name": "a.counter{key=value}",
	// 			"sum": 100
	// 		}
	// 	]
	// }
}
