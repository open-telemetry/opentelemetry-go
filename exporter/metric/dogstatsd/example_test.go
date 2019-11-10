package dogstatsd_test

import (
	"context"
	"fmt"
	"log"
	"net"
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
	// Create a server
	const addr = "127.0.0.1:18999"

	var wg sync.WaitGroup
	wg.Add(1)
	listener, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	finish := make(chan struct{}, 1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-finish:
				return
			default:
			}

			var buf [4096]byte
			if n, _, err := listener.ReadFrom(buf[:]); err != nil {
				panic(fmt.Sprint("Read err: ", err))
			} else if n >= len(buf) {
				panic(fmt.Sprint("Read small buffer: ", n))
			} else {
				fmt.Print(buf[0:n])
			}
		}
	}()

	// Create a meter
	selector := simple.NewWithExactMeasure()
	exporter, err := dogstatsd.New(dogstatsd.Config{
		URL: "udp://127.0.0.1:18899",
	})
	if err != nil {
		panic(fmt.Sprintln("Could not initialize dogstatsd exporter:", err))
	}
	batcher := ungrouped.New(selector, false)
	pusher := push.New(batcher, exporter, time.Second)
	pusher.Start()

	ctx := context.Background()

	key := key.New("key")
	meter := pusher.GetMeter("example")

	counter := meter.NewInt64Counter("a.counter", metric.WithKeys(key))
	labels := meter.Labels(key.String("value"))

	counter.Add(ctx, 100, labels)
	pusher.Stop()

	close(finish)
	wg.Wait()

	// Output:
	// X
}
