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

package trace

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	DefaultMaxQueueSize       = 2048
	DefaultBatchTimeout       = 5000 * time.Millisecond
	DefaultMaxExportBatchSize = 512
)

var (
	errNilExporter = errors.New("exporter is nil")
)

type BatchSpanProcessorOption func(o *BatchSpanProcessorOptions)

type BatchSpanProcessorOptions struct {
	// MaxQueueSize is the maximum queue size to buffer spans for delayed processing. If the
	// queue gets full it drops the spans. Use BlockOnQueueFull to change this behavior.
	// The default value of MaxQueueSize is 2048.
	MaxQueueSize int

	// BatchTimeout is the maximum duration for constructing a batch. Processor
	// forcefully sends available spans when timeout is reached.
	// The default value of BatchTimeout is 5000 msec.
	BatchTimeout time.Duration

	// MaxExportBatchSize is the maximum number of spans to process in a single batch.
	// If there are more than one batch worth of spans then it processes multiple batches
	// of spans one batch after the other without any delay.
	// The default value of MaxExportBatchSize is 512.
	MaxExportBatchSize int

	// BlockOnQueueFull blocks onEnd() and onStart() method if the queue is full
	// AND if BlockOnQueueFull is set to true.
	// Blocking option should be used carefully as it can severely affect the performance of an
	// application.
	BlockOnQueueFull bool
}

// BatchSpanProcessor implements SpanProcessor interfaces. It is used by
// exporters to receive export.SpanData asynchronously.
// Use BatchSpanProcessorOptions to change the behavior of the processor.
type BatchSpanProcessor struct {
	e export.SpanBatcher
	o BatchSpanProcessorOptions

	queue   chan *export.SpanData
	dropped uint32

	batch    []*export.SpanData
	timer    *time.Timer
	stopWait sync.WaitGroup
	stopOnce sync.Once
	stopCh   chan struct{}
}

var _ SpanProcessor = (*BatchSpanProcessor)(nil)

// NewBatchSpanProcessor creates a new instance of BatchSpanProcessor
// for a given export. It returns an error if exporter is nil.
// The newly created BatchSpanProcessor should then be registered with sdk
// using RegisterSpanProcessor.
func NewBatchSpanProcessor(e export.SpanBatcher, opts ...BatchSpanProcessorOption) (*BatchSpanProcessor, error) {
	if e == nil {
		return nil, errNilExporter
	}

	o := BatchSpanProcessorOptions{
		BatchTimeout:       DefaultBatchTimeout,
		MaxQueueSize:       DefaultMaxQueueSize,
		MaxExportBatchSize: DefaultMaxExportBatchSize,
	}
	for _, opt := range opts {
		opt(&o)
	}
	bsp := &BatchSpanProcessor{
		e:      e,
		o:      o,
		batch:  make([]*export.SpanData, 0, o.MaxExportBatchSize),
		timer:  time.NewTimer(o.BatchTimeout),
		queue:  make(chan *export.SpanData, o.MaxQueueSize),
		stopCh: make(chan struct{}),
	}

	bsp.stopWait.Add(1)
	go func() {
		defer bsp.stopWait.Done()
		bsp.processQueue()
		bsp.drainQueue()
	}()

	return bsp, nil
}

// OnStart method does nothing.
func (bsp *BatchSpanProcessor) OnStart(sd *export.SpanData) {
}

// OnEnd method enqueues export.SpanData for later processing.
func (bsp *BatchSpanProcessor) OnEnd(sd *export.SpanData) {
	bsp.enqueue(sd)
}

// Shutdown flushes the queue and waits until all spans are processed.
// It only executes once. Subsequent call does nothing.
func (bsp *BatchSpanProcessor) Shutdown() {
	bsp.stopOnce.Do(func() {
		close(bsp.stopCh)
		bsp.stopWait.Wait()
	})
}

func WithMaxQueueSize(size int) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.MaxQueueSize = size
	}
}

func WithMaxExportBatchSize(size int) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.MaxExportBatchSize = size
	}
}

func WithBatchTimeout(delay time.Duration) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.BatchTimeout = delay
	}
}

func WithBlocking() BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// exportSpans is a subroutine of processing and draining the queue.
func (bsp *BatchSpanProcessor) exportSpans() {
	bsp.timer.Reset(bsp.o.BatchTimeout)

	if len(bsp.batch) > 0 {
		bsp.e.ExportSpans(context.Background(), bsp.batch)
		bsp.batch = bsp.batch[:0]
	}
}

// processQueue removes spans from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to BatchTimeout to form a batch.
func (bsp *BatchSpanProcessor) processQueue() {
	defer bsp.timer.Stop()

	for {
		select {
		case <-bsp.stopCh:
			return
		case <-bsp.timer.C:
			bsp.exportSpans()
		case sd := <-bsp.queue:
			bsp.batch = append(bsp.batch, sd)
			if len(bsp.batch) == bsp.o.MaxExportBatchSize {
				if !bsp.timer.Stop() {
					<-bsp.timer.C
				}
				bsp.exportSpans()
			}
		}
	}
}

// drainQueue awaits the any caller that had added to bsp.stopWait
// to finish the enqueue, then exports the final batch.
func (bsp *BatchSpanProcessor) drainQueue() {
	for {
		select {
		case sd := <-bsp.queue:
			if sd == nil {
				bsp.exportSpans()
				return
			}

			bsp.batch = append(bsp.batch, sd)
			if len(bsp.batch) == bsp.o.MaxExportBatchSize {
				bsp.exportSpans()
			}
		default:
			close(bsp.queue)
		}
	}
}

func (bsp *BatchSpanProcessor) enqueue(sd *export.SpanData) {
	if !sd.SpanContext.IsSampled() {
		return
	}

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer func() {
		x := recover()
		switch err := x.(type) {
		case nil:
			return
		case runtime.Error:
			if err.Error() == "send on closed channel" {
				return
			}
		}
		panic(x)
	}()

	select {
	case <-bsp.stopCh:
		return
	default:
	}

	if bsp.o.BlockOnQueueFull {
		bsp.queue <- sd
		return
	}

	select {
	case bsp.queue <- sd:
	default:
		atomic.AddUint32(&bsp.dropped, 1)
	}
}
