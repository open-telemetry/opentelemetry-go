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

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/internal/env"
	"go.opentelemetry.io/otel/trace"
)

// Defaults for BatchLogRecordProcessorOptions.
const (
	DefaultMaxQueueSize       = 2048
	DefaultScheduleDelay      = 5000
	DefaultExportTimeout      = 30000
	DefaultMaxExportBatchSize = 512
)

// BatchLogRecordProcessorOption configures a BatchLogRecordProcessor.
type BatchLogRecordProcessorOption func(o *BatchLogRecordProcessorOptions)

// BatchLogRecordProcessorOptions is configuration settings for a
// BatchLogRecordProcessor.
type BatchLogRecordProcessorOptions struct {
	// MaxQueueSize is the maximum queue size to buffer logRecords for delayed processing. If the
	// queue gets full it drops the logRecords. Use BlockOnQueueFull to change this behavior.
	// The default value of MaxQueueSize is 2048.
	MaxQueueSize int

	// BatchTimeout is the maximum duration for constructing a batch. Processor
	// forcefully sends available logRecords when timeout is reached.
	// The default value of BatchTimeout is 5000 msec.
	BatchTimeout time.Duration

	// ExportTimeout specifies the maximum duration for exporting logRecords. If the timeout
	// is reached, the export will be cancelled.
	// The default value of ExportTimeout is 30000 msec.
	ExportTimeout time.Duration

	// MaxExportBatchSize is the maximum number of logRecords to process in a single batch.
	// If there are more than one batch worth of logRecords then it processes multiple batches
	// of logRecords one batch after the other without any delay.
	// The default value of MaxExportBatchSize is 512.
	MaxExportBatchSize int

	// BlockOnQueueFull blocks onEnd() and onStart() method if the queue is full
	// AND if BlockOnQueueFull is set to true.
	// Blocking option should be used carefully as it can severely affect the performance of an
	// application.
	BlockOnQueueFull bool
}

// batchLogRecordProcessor is a LogRecordProcessor that batches asynchronously-received
// logRecords and sends them to a trace.Exporter when complete.
type batchLogRecordProcessor struct {
	e LogRecordExporter
	o BatchLogRecordProcessorOptions

	queue   chan ReadOnlyLogRecord
	dropped uint32

	batch      []ReadOnlyLogRecord
	batchMutex sync.Mutex
	timer      *time.Timer
	stopWait   sync.WaitGroup
	stopOnce   sync.Once
	stopCh     chan struct{}
}

var _ LogRecordProcessor = (*batchLogRecordProcessor)(nil)

// NewBatchLogRecordProcessor creates a new LogRecordProcessor that will send completed
// span batches to the exporter with the supplied options.
//
// If the exporter is nil, the span processor will preform no action.
func NewBatchLogRecordProcessor(
	exporter LogRecordExporter, options ...BatchLogRecordProcessorOption,
) LogRecordProcessor {
	maxQueueSize := env.BatchSpanProcessorMaxQueueSize(DefaultMaxQueueSize)
	maxExportBatchSize := env.BatchSpanProcessorMaxExportBatchSize(DefaultMaxExportBatchSize)

	if maxExportBatchSize > maxQueueSize {
		if DefaultMaxExportBatchSize > maxQueueSize {
			maxExportBatchSize = maxQueueSize
		} else {
			maxExportBatchSize = DefaultMaxExportBatchSize
		}
	}

	o := BatchLogRecordProcessorOptions{
		BatchTimeout:       time.Duration(env.BatchSpanProcessorScheduleDelay(DefaultScheduleDelay)) * time.Millisecond,
		ExportTimeout:      time.Duration(env.BatchSpanProcessorExportTimeout(DefaultExportTimeout)) * time.Millisecond,
		MaxQueueSize:       maxQueueSize,
		MaxExportBatchSize: maxExportBatchSize,
	}
	for _, opt := range options {
		opt(&o)
	}
	bsp := &batchLogRecordProcessor{
		e:      exporter,
		o:      o,
		batch:  make([]ReadOnlyLogRecord, 0, o.MaxExportBatchSize),
		timer:  time.NewTimer(o.BatchTimeout),
		queue:  make(chan ReadOnlyLogRecord, o.MaxQueueSize),
		stopCh: make(chan struct{}),
	}

	bsp.stopWait.Add(1)
	go func() {
		defer bsp.stopWait.Done()
		bsp.processQueue()
		bsp.drainQueue()
	}()

	return bsp
}

// OnStart method does nothing.
func (bsp *batchLogRecordProcessor) OnEmit(parent context.Context, s ReadWriteLogRecord) {
	// Do not enqueue logRecords if we are just going to drop them.
	if bsp.e == nil {
		return
	}
	bsp.enqueue(s)
}

// Shutdown flushes the queue and waits until all logRecords are processed.
// It only executes once. Subsequent call does nothing.
func (bsp *batchLogRecordProcessor) Shutdown(ctx context.Context) error {
	var err error
	bsp.stopOnce.Do(
		func() {
			wait := make(chan struct{})
			go func() {
				close(bsp.stopCh)
				bsp.stopWait.Wait()
				if bsp.e != nil {
					if err := bsp.e.Shutdown(ctx); err != nil {
						otel.Handle(err)
					}
				}
				close(wait)
			}()
			// Wait until the wait group is done or the context is cancelled
			select {
			case <-wait:
			case <-ctx.Done():
				err = ctx.Err()
			}
		},
	)
	return err
}

type forceFlushLogRecord struct {
	ReadOnlyLogRecord
	flushed chan struct{}
}

func (f forceFlushLogRecord) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{TraceFlags: trace.FlagsSampled})
}

// ForceFlush exports all ended logRecords that have not yet been exported.
func (bsp *batchLogRecordProcessor) ForceFlush(ctx context.Context) error {
	var err error
	if bsp.e != nil {
		flushCh := make(chan struct{})
		if bsp.enqueueBlockOnQueueFull(ctx, forceFlushLogRecord{flushed: flushCh}) {
			select {
			case <-flushCh:
				// Processed any items in queue prior to ForceFlush being called
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		wait := make(chan error)
		go func() {
			wait <- bsp.exportLogRecords(ctx)
			close(wait)
		}()
		// Wait until the export is finished or the context is cancelled/timed out
		select {
		case err = <-wait:
		case <-ctx.Done():
			err = ctx.Err()
		}
	}
	return err
}

// WithMaxQueueSize returns a BatchLogRecordProcessorOption that configures the
// maximum queue size allowed for a BatchLogRecordProcessor.
func WithMaxQueueSize(size int) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.MaxQueueSize = size
	}
}

// WithMaxExportBatchSize returns a BatchLogRecordProcessorOption that configures
// the maximum export batch size allowed for a BatchLogRecordProcessor.
func WithMaxExportBatchSize(size int) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.MaxExportBatchSize = size
	}
}

// WithBatchTimeout returns a BatchLogRecordProcessorOption that configures the
// maximum delay allowed for a BatchLogRecordProcessor before it will export any
// held span (whether the queue is full or not).
func WithBatchTimeout(delay time.Duration) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.BatchTimeout = delay
	}
}

// WithExportTimeout returns a BatchLogRecordProcessorOption that configures the
// amount of time a BatchLogRecordProcessor waits for an exporter to export before
// abandoning the export.
func WithExportTimeout(timeout time.Duration) BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.ExportTimeout = timeout
	}
}

// WithBlocking returns a BatchLogRecordProcessorOption that configures a
// BatchLogRecordProcessor to wait for enqueue operations to succeed instead of
// dropping data when the queue is full.
func WithBlocking() BatchLogRecordProcessorOption {
	return func(o *BatchLogRecordProcessorOptions) {
		o.BlockOnQueueFull = true
	}
}

// exportLogRecords is a subroutine of processing and draining the queue.
func (bsp *batchLogRecordProcessor) exportLogRecords(ctx context.Context) error {
	bsp.timer.Reset(bsp.o.BatchTimeout)

	bsp.batchMutex.Lock()
	defer bsp.batchMutex.Unlock()

	if bsp.o.ExportTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, bsp.o.ExportTimeout)
		defer cancel()
	}

	if l := len(bsp.batch); l > 0 {
		global.Debug("exporting log records", "count", len(bsp.batch), "total_dropped", atomic.LoadUint32(&bsp.dropped))
		err := bsp.e.ExportLogRecords(ctx, bsp.batch)

		// A new batch is always created after exporting, even if the batch failed to be exported.
		//
		// It is up to the exporter to implement any type of retry logic if a batch is failing
		// to be exported, since it is specific to the protocol and backend being sent to.
		bsp.batch = bsp.batch[:0]

		if err != nil {
			return err
		}
	}
	return nil
}

// processQueue removes logRecords from the `queue` channel until processor
// is shut down. It calls the exporter in batches of up to MaxExportBatchSize
// waiting up to BatchTimeout to form a batch.
func (bsp *batchLogRecordProcessor) processQueue() {
	defer bsp.timer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case <-bsp.stopCh:
			return
		case <-bsp.timer.C:
			if err := bsp.exportLogRecords(ctx); err != nil {
				otel.Handle(err)
			}
		case sd := <-bsp.queue:
			if ffs, ok := sd.(forceFlushLogRecord); ok {
				close(ffs.flushed)
				continue
			}
			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) >= bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()
			if shouldExport {
				if !bsp.timer.Stop() {
					<-bsp.timer.C
				}
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
			}
		}
	}
}

// drainQueue awaits the any caller that had added to bsp.stopWait
// to finish the enqueue, then exports the final batch.
func (bsp *batchLogRecordProcessor) drainQueue() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case sd := <-bsp.queue:
			if sd == nil {
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
				return
			}

			bsp.batchMutex.Lock()
			bsp.batch = append(bsp.batch, sd)
			shouldExport := len(bsp.batch) == bsp.o.MaxExportBatchSize
			bsp.batchMutex.Unlock()

			if shouldExport {
				if err := bsp.exportLogRecords(ctx); err != nil {
					otel.Handle(err)
				}
			}
		default:
			close(bsp.queue)
		}
	}
}

func (bsp *batchLogRecordProcessor) enqueue(sd ReadOnlyLogRecord) {
	ctx := context.TODO()
	if bsp.o.BlockOnQueueFull {
		bsp.enqueueBlockOnQueueFull(ctx, sd)
	} else {
		bsp.enqueueDrop(ctx, sd)
	}
}

func recoverSendOnClosedChan() {
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
}

func (bsp *batchLogRecordProcessor) enqueueBlockOnQueueFull(ctx context.Context, sd ReadOnlyLogRecord) bool {
	if !sd.SpanContext().IsSampled() {
		return false
	}

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	case <-ctx.Done():
		return false
	}
}

func (bsp *batchLogRecordProcessor) enqueueDrop(ctx context.Context, sd ReadOnlyLogRecord) bool {
	if !sd.SpanContext().IsSampled() {
		return false
	}

	// This ensures the bsp.queue<- below does not panic as the
	// processor shuts down.
	defer recoverSendOnClosedChan()

	select {
	case <-bsp.stopCh:
		return false
	default:
	}

	select {
	case bsp.queue <- sd:
		return true
	default:
		atomic.AddUint32(&bsp.dropped, 1)
	}
	return false
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (bsp *batchLogRecordProcessor) MarshalLog() interface{} {
	return struct {
		Type         string
		SpanExporter LogRecordExporter
		Config       BatchLogRecordProcessorOptions
	}{
		Type:         "BatchLogRecordProcessor",
		SpanExporter: bsp.e,
		Config:       bsp.o,
	}
}
