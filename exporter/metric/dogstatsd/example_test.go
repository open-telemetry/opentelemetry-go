package dogstatsd_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/dogstatsd"

	// "go.opentelemetry.io/otel/exporter/metric/dogstatsd"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func ExampleNew() {
	var wg sync.WaitGroup
	finish := make(chan struct{}, 1)

	// Create a server
	tmpfile, err := ioutil.TempFile("", "examplegram")
	if err != nil {
		log.Fatal("Could not create tempfile: ", err)
	}
	path := tmpfile.Name()
	_ = tmpfile.Close()
	_ = os.Remove(path)
	defer func() {
		err := os.Remove(path)
		if err != nil {
			log.Fatal("Could not remove tempfile: ", err)
		}
	}()

	laddr, err := net.ResolveUnixAddr("unixgram", tmpfile.Name())
	if err != nil {
		log.Fatal("Could not resolve address: ", tmpfile, ":", err)
	}

	wg.Add(1)
	listener, err := net.ListenUnixgram("unixgram", laddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	go func() {
		defer wg.Done()

		for {
			select {
			case <-finish:
				return
			default:
			}

			var buf [4096]byte
			n, _, err := listener.ReadFrom(buf[:])
			if err != nil {
				panic(fmt.Sprint("Read err: ", err))
			} else if n >= len(buf) {
				panic(fmt.Sprint("Read small buffer: ", n))
			} else {
				fmt.Print(string(buf[0:n]))
			}
		}
	}()

	// Create a meter
	selector := simple.NewWithExactMeasure()
	exporter, err := dogstatsd.New(dogstatsd.Config{
		URL: fmt.Sprint("unix://", path),
	})
	if err != nil {
		panic(fmt.Sprintln("Could not initialize dogstatsd exporter:", err))
	}
	batcher := ungrouped.New(selector, false)
	pusher := push.New(batcher, exporter, time.Hour)
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
	// a.counter:100|c|#key:value
}
