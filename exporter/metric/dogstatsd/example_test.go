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
	"go.opentelemetry.io/otel/exporter/metric/dogstatsd"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
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
	selector := simple.NewWithExactMeasure()
	exporter, err := dogstatsd.New(dogstatsd.Config{
		// The Writer field provides test support.
		Writer: writer,

		// In real code, use the URL field:
		//
		// URL: fmt.Sprint("unix://", path),
	})
	if err != nil {
		log.Fatal("Could not initialize dogstatsd exporter:", err)
	}
	// The ungrouped batcher ensures that the export sees the full
	// set of labels as dogstatsd tags.
	batcher := ungrouped.New(selector, false)

	// The pusher automatically recognizes that the exporter
	// implements the LabelEncoder interface, which ensures the
	// export encoding for labels is encoded in the LabelSet.
	pusher := push.New(batcher, exporter, time.Hour)
	pusher.Start()

	ctx := context.Background()

	key := key.New("key")

	// pusher implements the metric.MeterProvider interface:
	meter := pusher.GetMeter("example")

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
