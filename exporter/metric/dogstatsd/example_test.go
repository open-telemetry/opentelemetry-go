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
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

func newStdoutServer() (string, func(), func()) {
	wg := &sync.WaitGroup{}
	finishChan := make(chan struct{}, 1)

	tmpfile, err := ioutil.TempFile("", "examplegram")
	if err != nil {
		log.Fatal("Could not create tempfile: ", err)
	}
	path := tmpfile.Name()
	_ = tmpfile.Close()
	_ = os.Remove(path)

	laddr, err := net.ResolveUnixAddr("unixgram", tmpfile.Name())
	if err != nil {
		log.Fatal("Could not resolve address: ", tmpfile, ":", err)
	}

	wg.Add(1)
	listener, err := net.ListenUnixgram("unixgram", laddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		defer wg.Done()

		for {
			select {
			case <-finishChan:
				return
			default:
			}

			var buf [4096]byte
			n, _, err := listener.ReadFrom(buf[:])
			if err != nil {
				log.Fatal("Read err: ", err)
			} else if n >= len(buf) {
				log.Fatal("Read small buffer: ", n)
			} else {
				fmt.Print(string(buf[0:n]))
			}
		}
	}()

	return path, func() {
			close(finishChan)
			wg.Wait()
		}, func() {
			_ = listener.Close()
			err := os.Remove(path)
			if err != nil {
				log.Fatal("Could not remove tempfile: ", err)
			}
		}
}

func ExampleNew() {
	// Create a server
	path, waitFunc, finishFunc := newStdoutServer()
	defer finishFunc()

	// Create a meter
	selector := simple.NewWithExactMeasure()
	exporter, err := dogstatsd.New(dogstatsd.Config{
		URL: fmt.Sprint("unix://", path),
	})
	if err != nil {
		log.Fatal("Could not initialize dogstatsd exporter:", err)
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

	waitFunc()

	// Output:
	// a.counter:100|c|#key:value
}
