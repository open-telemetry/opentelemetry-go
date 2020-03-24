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

package dogstatsd_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/dogstatsd"
)

func ExampleNew() {
	// Create a "server"
	wg := &sync.WaitGroup{}
	wg.Add(1)

	reader, writer := io.Pipe()

	go func() {
		defer wg.Done()

		for {
			var buf [4096]byte
			n, err := reader.Read(buf[:])
			if err == io.EOF {
				return
			} else if err != nil {
				log.Fatal("Read err: ", err)
			} else if n >= len(buf) {
				log.Fatal("Read small buffer: ", n)
			} else {
				fmt.Print(string(buf[0:n]))
			}
		}
	}()

	// Create a meter
	pusher, err := dogstatsd.NewExportPipeline(dogstatsd.Config{
		// The Writer field provides test support.
		Writer: writer,

		// In real code, use the URL field:
		//
		// URL: fmt.Sprint("unix://", path),
	}, time.Minute)
	if err != nil {
		log.Fatal("Could not initialize dogstatsd exporter:", err)
	}

	ctx := context.Background()

	key := key.New("key")

	// pusher implements the metric.MeterProvider interface:
	meter := pusher.Meter("example")

	// Create and update a single counter:
	counter := metric.Must(meter).NewInt64Counter("a.counter", metric.WithKeys(key))
	labels := meter.Labels(key.String("value"))

	counter.Add(ctx, 100, labels)

	// Flush the exporter, close the pipe, and wait for the reader.
	pusher.Stop()
	writer.Close()
	wg.Wait()

	// Output:
	// a.counter:100|c|#key:value
}
