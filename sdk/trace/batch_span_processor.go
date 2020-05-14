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
	"log"
	"sync"
	"sync/atomic"
	"time"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

const (
	defaultMaxQueueSize       = 2048
	defaultScheduledDelay     = 5000 * time.Millisecond
	defaultMaxExportBatchSize = 512
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

	// ScheduledDelayMillis is the delay interval in milliseconds between two consecutive
	// processing of batches.
	// The default value of ScheduledDelayMillis is 5000 msec.
	ScheduledDelayMillis time.Duration

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

	enqueueWait sync.WaitGroup
	stopWait    sync.WaitGroup
	stopOnce    sync.Once
	stopCh      chan struct{}
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
		ScheduledDelayMillis: defaultScheduledDelay,
		MaxQueueSize:         defaultMaxQueueSize,
		MaxExportBatchSize:   defaultMaxExportBatchSize,
	}
	for _, opt := range opts {
		opt(&o)
	}
	bsp := &BatchSpanProcessor{
		e: e,
		o: o,
	}

	bsp.queue = make(chan *export.SpanData, bsp.o.MaxQueueSize)

	bsp.stopCh = make(chan struct{})

	bsp.stopWait.Add(1)
	go func() {
		defer bsp.stopWait.Done()
		bsp.processQueue()
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

func WithScheduleDelayMillis(delay time.Duration) BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.ScheduledDelayMillis = delay
	}
}

func WithBlocking() BatchSpanProcessorOption {
	return func(o *BatchSpanProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// processQueue removes spans from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to ScheduledDelayMillis to form a batch.
func (bsp *BatchSpanProcessor) processQueue() {
	timer := time.NewTimer(bsp.o.ScheduledDelayMillis)
	defer timer.Stop()

	batch := make([]*export.SpanData, 0, bsp.o.MaxExportBatchSize)

	exportSpans := func() {
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(bsp.o.ScheduledDelayMillis)

		if len(batch) > 0 {
			bsp.e.ExportSpans(context.Background(), batch)
			batch = batch[:0]
		}
	}

loop:
	for {
		select {
		case <-bsp.stopCh:
			break loop
		case <-timer.C:
			exportSpans()
		case sd := <-bsp.queue:
			if sd.SpanContext.IsSampled() {
				batch = append(batch, sd)
				if len(batch) == bsp.o.MaxExportBatchSize {
					exportSpans()
				}
			}
		}
	}

	go func() {
		bsp.enqueueWait.Wait()
		close(bsp.queue)
	}()

	for {
		if !timer.Stop() {
			<-timer.C
		}
		const waitTimeout = 30 * time.Second
		timer.Reset(waitTimeout)

		select {
		case sd := <-bsp.queue:
			if sd == nil {
				exportSpans()
				return
			}

			if sd.SpanContext.IsSampled() {
				batch = append(batch, sd)
				if len(batch) == bsp.o.MaxExportBatchSize {
					exportSpans()
				}
			}
		case <-timer.C:
			log.Println("bsp.enqueueWait timeout")
			exportSpans()
			return
		}
	}
}

func (bsp *BatchSpanProcessor) enqueue(sd *export.SpanData) {
	bsp.enqueueWait.Add(1)

	select {
	case <-bsp.stopCh:
		bsp.enqueueWait.Done()
		return
	default:
	}

	if bsp.o.BlockOnQueueFull {
		bsp.queue <- sd
	} else {
		select {
		case bsp.queue <- sd:
		default:
			atomic.AddUint32(&bsp.dropped, 1)
		}
	}

	bsp.enqueueWait.Done()
}
