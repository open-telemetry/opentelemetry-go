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
	"bytes"
	"context"
	"fmt"

	"go.opentelemetry.io/otel/api/kv"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/stdout"
)

func ExampleNew() {
	buf := bytes.Buffer{}
	_, pusher, err := stdout.NewExportPipeline([]stdout.Option{
		// Defaults to STDOUT.
		stdout.WithWriter(&buf),
		stdout.WithPrettyPrint(),
		stdout.WithoutTimestamps(), // This makes the output deterministic
	}, nil)
	if err != nil {
		panic(fmt.Sprintln("Could not initialize stdout exporter:", err))
	}

	meter := metric.Must(pusher.Provider().Meter("example"))
	counter := meter.NewInt64Counter("a.counter")
	counter.Add(context.Background(), 100, kv.String("key", "value"))

	// Flush everything
	pusher.Stop()

	fmt.Println(buf.String())
	// Output:
	// [
	// 	{
	// 		"Name": "a.counter{instrumentation.name=example,key=value}",
	// 		"Sum": 100
	// 	}
	// ]
}
