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
	counter := meter.NewInt64Counter("a.counter", metric.WithKeys(key))
	labels := meter.Labels(key.String("value"))

	counter.Add(ctx, 100, labels)

	// Flush the exporter, close the pipe, and wait for the reader.
	pusher.Stop()
	writer.Close()
	wg.Wait()

	// Output:
	// a.counter:100|c|#key:value
}
