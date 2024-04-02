// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"container/ring"
	"context"
	"sync"
	"time"
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

// Compile-time check BatchingProcessor implements Processor.
var _ Processor = (*BatchingProcessor)(nil)

// BatchingProcessor is a processor that exports batches of log records.
type BatchingProcessor struct {
	exporter Exporter

	maxQueueSize       int
	exportInterval     time.Duration
	exportTimeout      time.Duration
	exportMaxBatchSize int
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
	return &BatchingProcessor{
		exporter: exporter,

		maxQueueSize:       cfg.maxQSize.Value,
		exportInterval:     cfg.expInterval.Value,
		exportTimeout:      cfg.expTimeout.Value,
		exportMaxBatchSize: cfg.expMaxBatchSize.Value,
	}
}

// OnEmit batches provided log record.
func (b *BatchingProcessor) OnEmit(ctx context.Context, r Record) error {
	// TODO (#5063): Implement.
	return nil
}

// Enabled returns true.
func (b *BatchingProcessor) Enabled(context.Context, Record) bool {
	return true
}

// Shutdown flushes queued log records and shuts down the decorated exporter.
func (b *BatchingProcessor) Shutdown(ctx context.Context) error {
	// TODO (#5063): Implement.
	return nil
}

// ForceFlush flushes queued log records and flushes the decorated exporter.
func (b *BatchingProcessor) ForceFlush(ctx context.Context) error {
	// TODO (#5063): Implement.
	return nil
}

// queue holds a queue of logging records.
//
// When the queue becomes full, the oldest records in the queue are
// overwritten.
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
//
// If enqueueing r will exceed the capacity of q, the oldest Record held in q
// will be dropped and r retained.
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

// TryFlush attempts to flush up to len(buf) Records. The available Records
// will be assigned into buf and passed to flush. If flush fails, returning
// false, the Records will not be removed from the queue. If flush succeeds,
// returning true, the flushed Records are removed from the queue. The number
// of Records remaining in the queue are returned.
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

// Flush returns all the Records held in the queue and resets it to be
// empty.
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
