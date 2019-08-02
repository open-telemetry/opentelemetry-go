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

package buffer

import (
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/experimental/streaming/exporter/observer"
)

type Buffer struct {
	observers []observer.Observer
	events    chan observer.Event
	dropped   uint64
	wait      sync.WaitGroup
	close     chan struct{}
}

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
