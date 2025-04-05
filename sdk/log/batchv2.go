// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
)

type recordsBufferPool struct {
	pool sync.Pool
}

func newRecordsBufferPool(bufferSize int) *recordsBufferPool {
	return &recordsBufferPool{
		pool: sync.Pool{
			New: func() any {
				slice := make([]Record, bufferSize)
				return &slice
			},
		},
	}
}

func (p *recordsBufferPool) Get() *[]Record {
	return p.pool.Get().(*[]Record)
}

func (p *recordsBufferPool) Put(recordsBuffer *[]Record) {
	p.pool.Put(recordsBuffer)
}

type recordsChunk struct {
	ctx     context.Context
	buf     *[]Record
	records int

	respCh chan<- error
}

// Compile-time check BatchProcessor implements Processor.
var _ Processor = (*BatchProcessorV2)(nil)

// BatchProcessorV2 is a processor that exports batches of log records.
//
// Use [NewBatchProcessorV2] to create a BatchProcessor. An empty BatchProcessor
// is shut down by default, no records will be batched or exported.
type BatchProcessorV2 struct {
	// The BatchProcessor is designed to provide the highest throughput of
	// log records possible while being compatible with OpenTelemetry. The
	// entry point of log records is the OnEmit method. This method is designed
	// to receive records as fast as possible while still honoring shutdown
	// commands. All records received are enqueued to queue.
	//
	// In order to block OnEmit as little as possible, a separate "poll"
	// goroutine is spawned at the creation of a BatchProcessor. This
	// goroutine is responsible for batching the queue at regular polled
	// intervals, or when it is directly signaled to.
	//
	// To keep the polling goroutine from backing up, all batches it makes are
	// exported with a bufferedExporter. This exporter allows the poll
	// goroutine to enqueue an export payload that will be handled in a
	// separate goroutine dedicated to the export. This asynchronous behavior
	// allows the poll goroutine to maintain accurate interval polling.
	//
	//   ___BatchProcessor____     __Poll Goroutine__     __Export Goroutine__
	// ||                     || ||                  || ||                    ||
	// ||          ********** || ||                  || ||     **********     ||
	// || Records=>* OnEmit * || ||   | - ticker     || ||     * export *     ||
	// ||          ********** || ||   | - trigger    || ||     **********     ||
	// ||             ||      || ||   |              || ||         ||         ||
	// ||             ||      || ||   |              || ||         ||         ||
	// ||   __________\/___   || ||   |***********   || ||   ______/\_______  ||
	// ||  (____queue______)>=||=||===|*  batch  *===||=||=>[_export_buffer_] ||
	// ||                     || ||   |***********   || ||                    ||
	// ||_____________________|| ||__________________|| ||____________________||
	//
	//
	// The "release valve" in this processing is the record queue. This queue
	// is a ring buffer. It will overwrite the oldest records first when writes
	// to OnEmit are made faster than the queue can be flushed. If batches
	// cannot be flushed to the export buffer, the records will remain in the
	// queue.

	exporter Exporter

	// queue is the active queue of records that have not yet been exported.
	queue *queue
	// batchSize is the minimum number of records needed before an export is
	// triggered (unless the interval expires).
	batchSize int

	// pollTrigger triggers the poll goroutine to flush a batch from the queue.
	// This is sent to when it is known that the queue contains at least one
	// complete batch.
	//
	// When a send is made to the channel, the poll loop will be reset after
	// the flush. If there is still enough records in the queue for another
	// batch the reset of the poll loop will automatically re-trigger itself.
	// There is no need for the original sender to monitor and resend.
	pollTrigger chan struct{}
	// pollKill kills the poll goroutine. This is only expected to be closed
	// once by the Shutdown method.
	pollKill chan struct{}
	// pollDone signals the poll goroutine has completed.
	pollDone chan struct{}

	// stopped holds the stopped state of the BatchProcessor.
	stopped atomic.Bool

	queuedChunks       chan recordsChunk
	queuedChunksClosed atomic.Bool
	queuedChunksMu     sync.Mutex
	queuedChunksDone   chan struct{}
	recordsBufPool     *recordsBufferPool

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// NewBatchProcessorV2 decorates the provided exporter
// so that the log records are batched before exporting.
//
// All of the exporter's methods are called synchronously.
func NewBatchProcessorV2(exporter Exporter, opts ...BatchProcessorOption) *BatchProcessorV2 {
	cfg := newBatchConfig(opts)
	if exporter == nil {
		// Do not panic on nil export.
		exporter = defaultNoopExporter
	}
	// Order is important here. Wrap the timeoutExporter with the chunkExporter
	// to ensure each export completes in timeout (instead of all chunked
	// exports).
	exporter = newTimeoutExporter(exporter, cfg.expTimeout.Value)
	// Use a chunkExporter to ensure ForceFlush and Shutdown calls are batched
	// appropriately on export.
	// exporter = newChunkExporter(exporter, cfg.expMaxBatchSize.Value)

	b := &BatchProcessorV2{
		exporter: exporter,

		queue:       newQueue(cfg.maxQSize.Value),
		batchSize:   cfg.expMaxBatchSize.Value,
		pollTrigger: make(chan struct{}, 1),
		pollKill:    make(chan struct{}),

		queuedChunks:   make(chan recordsChunk, cfg.expBufferSize.Value),
		recordsBufPool: newRecordsBufferPool(cfg.expMaxBatchSize.Value),
	}
	b.pollDone = b.poll(cfg.expInterval.Value)
	b.queuedChunksDone = b.processQueuedChunks()
	return b
}

// poll spawns a goroutine to handle interval polling and batch exporting. The
// returned done chan is closed when the spawned goroutine completes.
func (b *BatchProcessorV2) poll(interval time.Duration) (done chan struct{}) {
	done = make(chan struct{})

	ticker := time.NewTicker(interval)
	go func() {
		defer close(done)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
			case <-b.pollTrigger:
				ticker.Reset(interval)
			case <-b.pollKill:
				return
			}

			if d := b.queue.Dropped(); d > 0 {
				global.Warn("dropped log records", "dropped", d)
			}

			enqueued, recordsInQueue := b.enqueueExport(nil)
			if !enqueued || recordsInQueue >= b.batchSize {
				// There is another full batch ready. Immediately trigger
				// another export attempt.
				select {
				case b.pollTrigger <- struct{}{}:
				default:
					// Another flush signal already received.
				}
			}
		}
	}()

	return done
}

func (b *BatchProcessorV2) processQueuedChunks() (done chan struct{}) {
	done = make(chan struct{})

	go func() {
		defer close(done)
		for chunk := range b.queuedChunks {
			err := b.exporter.Export(chunk.ctx, (*chunk.buf)[:chunk.records])

			b.recordsBufPool.Put(chunk.buf)

			select {
			case chunk.respCh <- err:
			default:
				// e.respCh is nil or busy, default to otel.Handler.
				if err != nil {
					otel.Handle(err)
				}
			}
		}
	}()

	return done
}

func (b *BatchProcessorV2) enqueueExport(respCh chan<- error) (enqueued bool, queueLen int) {
	b.queuedChunksMu.Lock()
	defer b.queuedChunksMu.Unlock()

	if len(b.queuedChunks) == cap(b.queuedChunks) {
		if respCh != nil {
			respCh <- nil
		}
		return false, 0
	}

	buf := b.recordsBufPool.Get()

	queueLen, n := b.queue.Dequeue(*buf)

	data := recordsChunk{ctx: context.Background(), respCh: respCh, buf: buf, records: n}

	if b.queuedChunksClosed.Load() {
		if respCh != nil {
			respCh <- nil
		}
		return true, 0
	}

	// done in sync as available space is guaranteed by len check
	b.queuedChunks <- data

	return true, queueLen
}

// OnEmit batches provided log record.
func (b *BatchProcessorV2) OnEmit(_ context.Context, r *Record) error {
	if b.stopped.Load() || b.queue == nil {
		return nil
	}
	// The record is cloned so that changes done by subsequent processors
	// are not going to lead to a data race.
	if n := b.queue.Enqueue(r.Clone()); n >= b.batchSize {
		select {
		case b.pollTrigger <- struct{}{}:
		default:
			// Flush chan full. The poll goroutine will handle this by
			// re-sending any trigger until the queue has less than batchSize
			// records.
		}
	}
	return nil
}

// Shutdown flushes queued log records and shuts down the decorated exporter.
func (b *BatchProcessorV2) Shutdown(ctx context.Context) error {
	if b.stopped.Swap(true) || b.queue == nil {
		return nil
	}

	// Stop the poll goroutine.
	close(b.pollKill)
	select {
	case <-b.pollDone:
	case <-ctx.Done():
		// Out of time.
		return errors.Join(ctx.Err(), b.exporter.Shutdown(ctx))
	}

	err := b.flush(ctx)
	if err != nil {
		return errors.Join(err, b.exporter.Shutdown(ctx))
	}

	b.queuedChunksClosed.Store(true)
	b.queuedChunksMu.Lock()
	defer b.queuedChunksMu.Unlock()
	close(b.queuedChunks)
	select {
	case <-b.queuedChunksDone:
	case <-ctx.Done():
		// Out of time.
		return errors.Join(ctx.Err(), b.exporter.Shutdown(ctx))
	}

	return errors.Join(err, b.exporter.Shutdown(ctx))
}

// ForceFlush flushes queued log records and flushes the decorated exporter.
func (b *BatchProcessorV2) ForceFlush(ctx context.Context) error {
	if b.stopped.Load() || b.queue == nil {
		return nil
	}

	err := b.flush(ctx)

	return errors.Join(err, b.exporter.ForceFlush(ctx))
}

func (b *BatchProcessorV2) flush(ctx context.Context) error {
	var err error
	for {
		respCh := make(chan error, 1)
		enqueued, queueLen := b.enqueueExport(respCh)

		select {
		case respErr := <-respCh:
			if respErr != nil {
				err = errors.Join(respErr, errPartialFlush)
			}
		case <-ctx.Done():
			err = errors.Join(ctxErr(ctx), errPartialFlush)
		}

		if err != nil {
			break
		}

		if enqueued && queueLen == 0 {
			break
		}
	}

	return err
}
