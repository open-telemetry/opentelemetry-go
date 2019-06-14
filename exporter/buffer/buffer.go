package buffer

import (
	"sync"
	"sync/atomic"

	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type (
	Buffer struct {
		observers []observer.Observer
		events    chan observer.Event
		dropped   uint64
		wait      sync.WaitGroup
		close     chan struct{}
	}
)

func NewBuffer(size int, observers ...observer.Observer) *Buffer {
	b := &Buffer{
		observers: observers,
		events:    make(chan observer.Event, size),
		close:     make(chan struct{}),
	}
	b.wait.Add(1)
	go b.run()
	return b
}

func (b *Buffer) Observe(data observer.Event) {
	select {
	case b.events <- data:
	default:
		atomic.AddUint64(&b.dropped, 1)
	}
}

func (b *Buffer) Close() {
	close(b.close)
	b.wait.Wait()
}

func (b *Buffer) run() {
	defer func() {
		_ = recover()
		b.wait.Done()
	}()

	for {
		select {
		case <-b.close:
			return
		case ev := <-b.events:
			// TODO: This has to ensure ordered arrival,
			// e.g., put into a heap and delay observations.
			for _, obs := range b.observers {
				obs.Observe(ev)
			}
		}
	}
}
