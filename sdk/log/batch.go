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
// A BatchingProcessor must be created with [NewBatchingProcessor].
type BatchingProcessor struct {
	exporter Exporter

	maxQueueSize       int
	exportInterval     time.Duration
	exportTimeout      time.Duration
	exportMaxBatchSize int

	mu    sync.Mutex
	queue []Record

	batch []Record

	stop       chan exportRequest
	flush      chan exportRequest
	isShutdown atomic.Bool
	done       chan struct{}
}

type exportRequest struct {
	Context context.Context
	Result  chan error
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
	b := &BatchingProcessor{
		exporter: exporter,

		maxQueueSize:       cfg.maxQSize.Value,
		exportInterval:     cfg.expInterval.Value,
		exportTimeout:      cfg.expTimeout.Value,
		exportMaxBatchSize: cfg.expMaxBatchSize.Value,

		flush: make(chan exportRequest),
		stop:  make(chan exportRequest, 1),
		done:  make(chan struct{}),
		queue: make([]Record, 0, cfg.maxQSize.Value),
		batch: make([]Record, 0, cfg.expMaxBatchSize.Value),
	}

	go b.run()
	return b
}

// OnEmit batches provided log record.
func (b *BatchingProcessor) OnEmit(ctx context.Context, r Record) error {
	if b.isShutdown.Load() {
		return nil
	}

	defer b.mu.Unlock()
	b.mu.Lock()

	if len(b.queue) == b.maxQueueSize {
		// Queue is full.
		// Drop the record.
		return nil
	}
	b.queue = append(b.queue, r)
	return nil
}

// Enabled returns true.
func (b *BatchingProcessor) Enabled(context.Context, Record) bool {
	return true
}

// Shutdown flushes queued log records and shuts down the decorated exporter.
func (b *BatchingProcessor) Shutdown(ctx context.Context) error {
	wasShutdown := b.isShutdown.Swap(true)
	if wasShutdown {
		return nil
	}

	req := exportRequest{
		Context: ctx,
		Result:  make(chan error, 1), // Heap allocation.
	}
	// Send to a buffered channel so that it eventually closes the exporting goroutine.
	// This line can be called only once so it will never be blocking operation.
	b.stop <- req

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-req.Result:
		return err
	}
}

// ForceFlush flushes queued log records and flushes the decorated exporter.
func (b *BatchingProcessor) ForceFlush(ctx context.Context) error {
	if b.isShutdown.Load() {
		return nil
	}

	req := exportRequest{
		Context: ctx,
		Result:  make(chan error, 1), // Heap allocation.
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.done:
		// The batching processor was concurrently shutdown.
		return nil
	case b.flush <- req:
		return <-req.Result
	}
}

func (b *BatchingProcessor) run() {
	defer close(b.done)

	ticker := time.NewTicker(b.exportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := b.export(context.Background())
			if err != nil {
				otel.Handle(err)
			}
		case req := <-b.flush:
			err := b.export(req.Context)
			err = errors.Join(err, b.exporter.ForceFlush(req.Context))
			req.Result <- err
			ticker.Reset(b.exportInterval)
		case req := <-b.stop:
			err := b.export(req.Context)
			err = errors.Join(err, b.exporter.Shutdown(req.Context))
			req.Result <- err
			return
		}
	}
}

func (b *BatchingProcessor) export(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.exportTimeout) // 5 heap allocations.
	defer cancel()

	b.mu.Lock()
	if len(b.queue) == 0 {
		// Nothing to export
		b.mu.Unlock()
		return nil
	}
	if len(b.queue) > b.exportMaxBatchSize {
		b.batch = append(b.batch, b.queue[:b.exportMaxBatchSize]...)
		b.queue = b.queue[b.exportMaxBatchSize:]
	} else {
		b.batch = append(b.batch, b.queue...)
		b.queue = b.queue[:0]
	}
	b.mu.Unlock()

	// Doing export outside of the lock prevents deadlocks and improves efficiency.
	err := b.exporter.Export(ctx, b.batch)
	b.batch = b.batch[:0]
	return err
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
