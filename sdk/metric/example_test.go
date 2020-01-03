// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/stdout"
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

	key := key.New("key")
	meter := pusher.Meter("example")

	counter := meter.NewInt64Counter("a.counter", metric.WithKeys(key))
	labels := meter.Labels(key.String("value"))

	counter.Add(ctx, 100, labels)

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
