// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
)

const (
	queueSizeDefault    = 2048
	intervalDefault     = time.Second
	timeoutDefault      = 30 * time.Second
	maxBatchSizeDefault = 512
)

var _ Processor = (*BatchingProcessor)(nil)

// BatchingProcessor is an processor that asynchronously exports batches of log records.
type BatchingProcessor struct {
	exporter Exporter
	cfg      batcherConfig

	mu    sync.Mutex
	queue []Record

	batch []Record

	stop       chan exportRequest
	flush      chan exportRequest
	isShutdown atomic.Bool
	done       chan struct{}
}

type batcherConfig struct {
	queueSize    int
	interval     time.Duration
	timeout      time.Duration
	maxBatchSize int
}

type exportRequest struct {
	Context context.Context
	Result  chan error
}

// NewBatchingProcessor decorates the provided exporter
// so that the log records are batched before exporting.
func NewBatchingProcessor(exporter Exporter, opts ...BatchingOption) *BatchingProcessor {
	cfg := batcherConfig{
		queueSize:    queueSizeDefault,
		interval:     intervalDefault,
		timeout:      timeoutDefault,
		maxBatchSize: maxBatchSizeDefault,
	}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if v := os.Getenv("OTEL_BLRP_MAX_QUEUE_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_BLRP_MAX_QUEUE_SIZE value: %w", err))
		} else {
			cfg.queueSize = n
		}
	}
	if v := os.Getenv("OTEL_BSP_SCHEDULE_DELAY"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_BSP_SCHEDULE_DELAY value: %w", err))
		} else {
			cfg.interval = time.Duration(n) * time.Millisecond
		}
	}
	if v := os.Getenv("OTEL_BSP_EXPORT_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_BSP_EXPORT_TIMEOUT value: %w", err))
		} else {
			cfg.timeout = time.Duration(n) * time.Millisecond
		}
	}
	if v := os.Getenv("OTEL_BSP_MAX_EXPORT_BATCH_SIZE"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_BSP_MAX_EXPORT_BATCH_SIZE value: %w", err))
		} else {
			cfg.timeout = time.Duration(n) * time.Millisecond
		}
	}

	if cfg.queueSize <= 0 {
		otel.Handle(fmt.Errorf("batcher max queue size must be positive but was %v, setting default value", cfg.queueSize))
		cfg.queueSize = queueSizeDefault
	}
	if cfg.interval <= 0 {
		otel.Handle(fmt.Errorf("batcher interval must be positive but was %v, setting default value", cfg.interval))
		cfg.interval = intervalDefault
	}
	if cfg.timeout <= 0 {
		otel.Handle(fmt.Errorf("batcher timeout must be positive but was %v, setting default value", cfg.timeout))
		cfg.timeout = timeoutDefault
	}
	if cfg.maxBatchSize <= 0 {
		otel.Handle(fmt.Errorf("batcher max batch size must be positive but was %v, setting default value", cfg.maxBatchSize))
		cfg.maxBatchSize = maxBatchSizeDefault
	}

	b := &BatchingProcessor{
		exporter: exporter,
		cfg:      cfg,
		flush:    make(chan exportRequest),
		stop:     make(chan exportRequest),
		done:     make(chan struct{}),
		queue:    make([]Record, 0, cfg.queueSize),
		batch:    make([]Record, 0, cfg.maxBatchSize),
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

	if len(b.queue) == b.cfg.queueSize {
		// Queue is full.
		return nil
	}
	b.queue = append(b.queue, r)
	return nil
}

// Shutdown flushes queued log records and shuts down the decorated expoter.
func (b *BatchingProcessor) Shutdown(ctx context.Context) error {
	wasShutdown := b.isShutdown.Swap(true)
	if wasShutdown {
		return nil
	}

	req := exportRequest{
		Context: ctx,
		Result:  make(chan error, 1), // Heap allocation.
	}
	b.stop <- req
	return <-req.Result
}

// ForceFlush flushes queued log records and flushes the decorated expoter.
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
		return nil
	case b.flush <- req:
		return <-req.Result
	}
}

func (b *BatchingProcessor) run() {
	defer close(b.done)

	ticker := time.NewTicker(b.cfg.interval)
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
			ticker.Reset(b.cfg.interval)
		case req := <-b.stop:
			err := b.export(req.Context)
			err = errors.Join(err, b.exporter.Shutdown(req.Context))
			req.Result <- err
			return
		}
	}
}

func (b *BatchingProcessor) export(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.cfg.timeout) // 5 heap allocations.
	defer cancel()

	b.mu.Lock()
	if len(b.queue) == 0 {
		// Nothing to export
		b.mu.Unlock()
		return nil
	}
	if len(b.queue) > b.cfg.maxBatchSize {
		b.batch = append(b.batch, b.queue[:b.cfg.maxBatchSize]...)
		b.queue = b.queue[b.cfg.maxBatchSize:]
	} else {
		b.batch = append(b.batch, b.queue...)
		b.queue = b.queue[:0]
	}
	b.mu.Unlock()

	// Doing export outside of the log prevents deadlocks and improves efficiency.
	err := b.exporter.Export(ctx, b.batch)
	b.batch = b.batch[:0]
	return err
}

// BatchingOption applies a configuration to a Batcher.
type BatchingOption interface {
	apply(batcherConfig) batcherConfig
}

type batchingOptionFunc func(batcherConfig) batcherConfig

func (fn batchingOptionFunc) apply(c batcherConfig) batcherConfig {
	return fn(c)
}

// WithMaxQueueSize sets the maximum queue size used by the Batcher.
// After the size is reached log records are dropped.
//
// If the OTEL_BLRP_MAX_QUEUE_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BLRP_MAX_QUEUE_SIZE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 2048 will be used.
// The default value is also used when the provided value is not a positive value.
func WithMaxQueueSize(max int) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		cfg.queueSize = max
		return cfg
	})
}

// WithExportInterval sets the maximum duration between batched exports.
//
// If the OTEL_BSP_SCHEDULE_DELAY environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_SCHEDULE_DELAY will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 1s will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportInterval(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		cfg.interval = d
		return cfg
	})
}

// WithExportTimeout sets the duration after which a batched export is canceled.
//
// If the OTEL_BSP_EXPORT_TIMEOUT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_EXPORT_TIMEOUT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 30s will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportTimeout(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		cfg.timeout = d
		return cfg
	})
}

// WithExportMaxBatchSize sets the maximum batch size of every export.
//
// If the OTEL_BSP_MAX_EXPORT_BATCH_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_BSP_MAX_EXPORT_BATCH_SIZE will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, 512 will be used.
// The default value is also used when the provided value is not a positive value.
func WithExportMaxBatchSize(max int) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		cfg.maxBatchSize = max
		return cfg
	})
}
