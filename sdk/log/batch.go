// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"container/ring"
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/internal/global"
)

const (
	dfltMaxQSize        = 2048
	dfltExpInterval     = time.Second
	dfltExpTimeout      = 30 * time.Second
	dfltExpMaxBatchSize = 512

	envarMaxQSize        = "OTEL_BLRP_MAX_QUEUE_SIZE"
	envarExpInterval     = "OTEL_BLRP_SCHEDULE_DELAY"
	envarExpTimeout      = "OTEL_BLRP_EXPORT_TIMEOUT"
	envarExpMaxBatchSize = "OTEL_BLRP_MAX_EXPORT_BATCH_SIZE"
)

var errOverflow = errors.New("export overflow")

// Compile-time check BatchingProcessor implements Processor.
var _ Processor = (*BatchingProcessor)(nil)

// BatchingProcessor is a processor that exports batches of log records.
type BatchingProcessor struct {
	exporter   Exporter
	exportCh   chan exportData
	exportDone chan struct{}

	// q is the active queue of records that have not yet been exported.
	q *queue
	// batchSize is the minimum number of Records needed before an export is
	// triggered (unless the interval expires).
	batchSize int

	batchComplete chan struct{}
	done          chan struct{}

	// ctx is the parent context for the BatchingProcessor asynchronous
	// operations. When this context is canceled exports are canceled and
	// polling is stopped.
	ctx    context.Context
	cancel context.CancelFunc

	// stopped holds the stopped state of the BatchingProcessor.
	stopped atomic.Bool
}

// NewBatchingProcessor decorates the provided exporter
// so that the log records are batched before exporting.
//
// All of the exporter's methods are called synchronously.
func NewBatchingProcessor(exporter Exporter, opts ...BatchingOption) *BatchingProcessor {
	if exporter == nil {
		// Do not panic on nil export.
		exporter = defaultNoopExporter
	}
	cfg := newBatchingConfig(opts)

	ctx, cancel := context.WithCancel(context.Background())

	// TODO: explore making the size of this configurable.
	expCh := make(chan exportData, 1)
	exporter = &timeoutExporter{
		Exporter: exporter,
		timeout:  cfg.expTimeout.Value,
	}

	b := &BatchingProcessor{
		exporter:   exporter,
		exportCh:   expCh,
		exportDone: exportSync(expCh, exporter),

		q:             newQueue(cfg.maxQSize.Value),
		batchSize:     cfg.expMaxBatchSize.Value,
		batchComplete: make(chan struct{}, 1),

		ctx:    ctx,
		cancel: cancel,
	}
	b.done = b.run()
	return b
}

func (b *BatchingProcessor) run() (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		defer func() { close(done) }()
		for {
			// TODO: add interval polling.
			select {
			case <-b.batchComplete:
			case <-b.ctx.Done():
				// Parent context done. Shutting down.
				return
			}

			// TODO: sync.Pool to hold these.
			buf := make([]Record, b.batchSize)

			qLen := b.q.TryFlush(buf, func(r []Record) bool {
				data := exportData{ctx: context.Background(), records: r}
				select {
				case b.exportCh <- data:
					return true
				default:
					return false
				}
			})

			if qLen >= b.batchSize {
				select {
				case b.batchComplete <- struct{}{}:
				default:
					// Another flush signal already received.
				}
			}
		}
	}()
	return done
}

// OnEmit batches provided log record.
func (b *BatchingProcessor) OnEmit(_ context.Context, r Record) error {
	if b.stopped.Load() {
		return nil
	}
	if n := b.q.Enqueue(r); n >= b.batchSize {
		select {
		case b.batchComplete <- struct{}{}:
		default:
			// Flush chan full. The full queue is scheduled to be flushed but
			// we are enqueueing faster than we can flush. Old data in the
			// queue is going to be dropped at this point.
		}
	}
	return nil
}

// Enabled returns if b is enabled.
func (b *BatchingProcessor) Enabled(context.Context, Record) bool {
	return !b.stopped.Load()
}

// Shutdown flushes queued log records and shuts down the decorated exporter.
func (b *BatchingProcessor) Shutdown(ctx context.Context) error {
	if b.stopped.Swap(true) {
		return nil
	}

	// Stop the run goroutine.
	b.cancel()
	select {
	case <-b.done:
	case <-ctx.Done():
		// Out of time.
		return errors.Join(ctx.Err(), b.exporter.Shutdown(ctx))
	}

	// Flush remaining queued.
	err := b.flushAll(ctx)

	// Stop the export goroutine, no more data will be sent.
	close(b.exportCh)
	select {
	case <-b.exportDone:
	case <-ctx.Done():
		err = errors.Join(err, ctx.Err())
	}
	return errors.Join(err, b.exporter.Shutdown(ctx))
}

// ForceFlush flushes queued log records and flushes the decorated exporter.
func (b *BatchingProcessor) ForceFlush(ctx context.Context) error {
	if b.stopped.Load() {
		return nil
	}
	err := b.flushAll(ctx)
	return errors.Join(err, b.exporter.ForceFlush(ctx))
}

func (b *BatchingProcessor) flushAll(ctx context.Context) error {
	records := b.q.Flush()
	n := len(records)
	if n == 0 {
		return nil
	}

	bSize := b.batchSize
	batches := 1 + ((n - 1) / bSize)
	resp := make(chan error, batches)

	for i, j := 0, min(bSize, n); i < n; i, j = i+bSize, min(j+bSize, n) {
		select {
		case b.exportCh <- exportData{ctx, records[i:j], resp}:
		case <-ctx.Done():
			global.Error(errOverflow, "dropping records", "count", len(records[i:]))
			return ctx.Err()
		}
	}

	var err error
	for i := 0; i < batches; i++ {
		select {
		case e := <-resp:
			errors.Join(err, e)
		case <-ctx.Done():
			return errors.Join(err, ctx.Err())
		}
	}
	return err
}

// queue holds a queue of logging records.
type queue struct {
	sync.Mutex

	cap, len    int
	read, write *ring.Ring
}

func newQueue(size int) *queue {
	r := ring.New(size)
	return &queue{
		cap:   size,
		read:  r,
		write: r,
	}
}

// Enqueue adds r to the queue. The queue size, including the addition of r, is
// returned.
func (q *queue) Enqueue(r Record) int {
	q.Lock()
	defer q.Unlock()

	q.write.Value = r
	q.write = q.write.Next()

	q.len++
	if q.len > q.cap {
		// Overflow. Advance read to be the new "oldest".
		q.len = q.cap
		q.read = q.read.Next()
	}
	return q.len
}

// TryFlush attempts to flush up to len(buf) Records into buf. It returns the
// number of Records flushed (0 <= n <= len(p)).
func (q *queue) TryFlush(buf []Record, flush func([]Record) bool) int {
	q.Lock()
	defer q.Unlock()

	origRead := q.read

	n := min(len(buf), q.len)
	for i := 0; i < n; i++ {
		buf[i] = q.read.Value.(Record)
		q.read = q.read.Next()
	}

	if flush(buf[:n]) {
		q.len -= n
	} else {
		q.read = origRead
	}
	return q.len
}

func (q *queue) Flush() []Record {
	q.Lock()
	defer q.Unlock()

	out := make([]Record, q.len)
	for i := range out {
		out[i] = q.read.Value.(Record)
		q.read = q.read.Next()
	}
	q.len = 0

	return out
}

type batchingConfig struct {
	maxQSize        setting[int]
	expInterval     setting[time.Duration]
	expTimeout      setting[time.Duration]
	expMaxBatchSize setting[int]
}

func newBatchingConfig(options []BatchingOption) batchingConfig {
	var c batchingConfig
	for _, o := range options {
		c = o.apply(c)
	}

	c.maxQSize = c.maxQSize.Resolve(
		clearLessThanOne[int](),
		getenv[int](envarMaxQSize),
		clearLessThanOne[int](),
		fallback[int](dfltMaxQSize),
	)
	c.expInterval = c.expInterval.Resolve(
		clearLessThanOne[time.Duration](),
		getenv[time.Duration](envarExpInterval),
		clearLessThanOne[time.Duration](),
		fallback[time.Duration](dfltExpInterval),
	)
	c.expTimeout = c.expTimeout.Resolve(
		clearLessThanOne[time.Duration](),
		getenv[time.Duration](envarExpTimeout),
		clearLessThanOne[time.Duration](),
		fallback[time.Duration](dfltExpTimeout),
	)
	c.expMaxBatchSize = c.expMaxBatchSize.Resolve(
		clearLessThanOne[int](),
		getenv[int](envarExpMaxBatchSize),
		clearLessThanOne[int](),
		fallback[int](dfltExpMaxBatchSize),
	)

	return c
}

// BatchingOption applies a configuration to a [BatchingProcessor].
type BatchingOption interface {
	apply(batchingConfig) batchingConfig
}

type batchingOptionFunc func(batchingConfig) batchingConfig

func (fn batchingOptionFunc) apply(c batchingConfig) batchingConfig {
	return fn(c)
}

// WithMaxQueueSize sets the maximum queue size used by the Batcher.
// After the size is reached log records are dropped.
//
// If the OTEL_BLRP_MAX_QUEUE_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 2048 will be used.
// The default value is also used when the provided value is less than one.
func WithMaxQueueSize(size int) BatchingOption {
	return batchingOptionFunc(func(cfg batchingConfig) batchingConfig {
		cfg.maxQSize = newSetting(size)
		return cfg
	})
}

// WithExportInterval sets the maximum duration between batched exports.
//
// If the OTEL_BLRP_SCHEDULE_DELAY environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 1s will be used.
// The default value is also used when the provided value is less than one.
func WithExportInterval(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batchingConfig) batchingConfig {
		cfg.expInterval = newSetting(d)
		return cfg
	})
}

// WithExportTimeout sets the duration after which a batched export is canceled.
//
// If the OTEL_BLRP_EXPORT_TIMEOUT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 30s will be used.
// The default value is also used when the provided value is less than one.
func WithExportTimeout(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batchingConfig) batchingConfig {
		cfg.expTimeout = newSetting(d)
		return cfg
	})
}

// WithExportMaxBatchSize sets the maximum batch size of every export.
// A batch will be split into multiple exports to not exceed this size.
//
// If the OTEL_BLRP_MAX_EXPORT_BATCH_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 512 will be used.
// The default value is also used when the provided value is less than one.
func WithExportMaxBatchSize(size int) BatchingOption {
	return batchingOptionFunc(func(cfg batchingConfig) batchingConfig {
		cfg.expMaxBatchSize = newSetting(size)
		return cfg
	})
}
