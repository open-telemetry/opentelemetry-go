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

package metric_test

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/api/kv"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
)

func ExampleNew() {
	pusher, err := stdout.NewExportPipeline(stdout.Config{
		PrettyPrint:    true,
		DoNotPrintTime: true, // This makes the output deterministic
	})
	if err != nil {
		panic(fmt.Sprintln("Could not initialize stdout exporter:", err))
	}
	defer pusher.Stop()

	ctx := context.Background()

	key := kv.Key("key")
	meter := pusher.Provider().Meter("example")

	counter := metric.Must(meter).NewInt64Counter("a.counter")

	counter.Add(ctx, 100, key.String("value"))

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
