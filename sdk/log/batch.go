// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"slices"
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

// Compile-time check BatchingProcessor implements Processor.
var _ Processor = (*BatchingProcessor)(nil)

// BatchingProcessor is a processor that exports batches of log records.
type BatchingProcessor struct {
	exportCh   chan exportData
	exportDone chan struct{}
	exporter   Exporter

	batch *batch

	ctx      context.Context
	cancel   context.CancelFunc
	pollDone chan struct{}
	stopped  atomic.Bool
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

	expCh := make(chan exportData, 2)
	exp := chunker{
		Exporter: exporter,
		Size:     cfg.expMaxBatchSize.Value,
		Timeout:  cfg.expTimeout.Value,
	}

	ctx, cancel := context.WithCancel(context.Background())

	b := &BatchingProcessor{
		exporter:   exp,
		exportCh:   expCh,
		exportDone: exportSync(expCh, exp),
		batch:      newBatch(cfg.maxQSize.Value),
		ctx:        ctx,
		cancel:     cancel,
	}

	b.pollDone = b.poll(cfg.expInterval.Value)
	return b
}

var (
	errElapsed = errors.New("interval elapsed")

	// Used for testing.
	ctxWithDeadlineCause = context.WithDeadlineCause
)

func (b *BatchingProcessor) poll(interval time.Duration) (done chan struct{}) {
	done = make(chan struct{})

	start := time.Now()
	go func() {
		defer func() { close(b.pollDone) }()
		for {
			// Wait until an interval has passed or the parent context stopped.
			deadline := start.Add(interval)
			// Use a custom error for the deadline to distinguish between a
			// parent context deadline being reached and a child context one.
			ctx, cancel := ctxWithDeadlineCause(b.ctx, deadline, errElapsed)
			<-ctx.Done()
			cancel() // Release any resource held by ctx.

			if !errors.Is(context.Cause(ctx), errElapsed) {
				// Parent context done.
				return
			}

			var records []Record
			records, start = b.batch.FlushStale(start)
			b.enqueue(b.ctx, records, nil)
		}
	}()
	return done
}

// enqueue attempts to enqueue an export. If the exportCh is full, the export
// will be dropped and an error logged.
func (b *BatchingProcessor) enqueue(ctx context.Context, r []Record, rCh chan error) {
	select {
	case b.exportCh <- exportData{ctx, r, rCh}:
	default:
		// Export chan full. Do not block.
		err := errors.New("export overflow")
		global.Error(err, "dropping log batch", "dropped", len(r))
	}
}

// OnEmit batches provided log record.
func (b *BatchingProcessor) OnEmit(_ context.Context, r Record) error {
	if b.stopped.Load() {
		return nil
	}
	if flushed := b.batch.Append(r); flushed != nil {
		b.enqueue(b.ctx, flushed, nil)
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

	resp := make(chan error, 1)
	b.enqueue(ctx, b.batch.Flush(), resp)

	// Close poller.
	b.cancel()
	<-b.pollDone

	// Wait for response before closing exporter.
	var err error
	select {
	case err = <-resp:
		close(resp)
	case <-ctx.Done():
		// Out of time. Ignore flush response.
		close(b.exportCh)
		return errors.Join(ctx.Err(), b.exporter.Shutdown(ctx))
	}

	// Close exporter.
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
	resp := make(chan error, 1)
	defer func() { close(resp) }()
	b.enqueue(ctx, b.batch.Flush(), resp)

	var err error
	select {
	case err = <-resp:
	case <-ctx.Done():
		return ctx.Err()
	}
	return errors.Join(err, b.exporter.ForceFlush(ctx))
}

// batch holds a batch of logging records.
type batch struct {
	sync.Mutex

	data []Record
	// oldest is the timestamp of when the first Record in the batch was
	// received.
	oldest time.Time
}

func newBatch(n int) *batch {
	return &batch{data: make([]Record, 0, n)}
}

// Append adds r to the batch. If adding r fills the batch, the batch is
// flushed and its contents returned.
func (b *batch) Append(r Record) []Record {
	b.Lock()
	defer b.Unlock()

	b.data = append(b.data, r)
	if len(b.data) == 1 {
		b.oldest = time.Now()
	}
	if len(b.data) == cap(b.data) {
		return b.flush()
	}
	return nil
}

// Flush returns and clears the contents of the batch.
func (b *batch) Flush() []Record {
	b.Lock()
	defer b.Unlock()

	return b.flush()
}

// Flush flushes batch b if it is older than timestamp.
//
// If the batch is flushed, the flushed records along with the estimated start
// time (i.e. now) of the next batch is returned.
//
// If the batch is not flushed, because the batch is empty or it is younger
// than timestamp, the returned time will be the start time of the still held
// batch. This time will be the estimated start time (i.e. now) in the case
// where the batch is empty.
func (b *batch) FlushStale(timestamp time.Time) ([]Record, time.Time) {
	b.Lock()
	defer b.Unlock()

	if len(b.data) == 0 {
		return nil, time.Now()
	}

	if b.oldest.After(timestamp) {
		// Not stale.
		return nil, b.oldest
	}

	return b.flush(), time.Now()
}

func (b *batch) flush() []Record {
	clone := slices.Clone(b.data)
	b.data = b.data[:0]
	b.oldest = time.Time{}
	return clone
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
