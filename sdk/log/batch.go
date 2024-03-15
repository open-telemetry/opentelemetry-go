// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"time"
)

// Compile-time check BatchingProcessor implements Processor.
var _ Processor = (*BatchingProcessor)(nil)

// BatchingProcessor is an processor that asynchronously exports batches of log records.
type BatchingProcessor struct{}

type batcherConfig struct{}

// NewBatchingProcessor decorates the provided exporter
// so that the log records are batched before exporting.
//
// All of the exporter's methods are called from a single dedicated
// background goroutine. Therefore, the expoter does not need to
// be concurrent safe.
func NewBatchingProcessor(exporter Exporter, opts ...BatchingOption) *BatchingProcessor {
	// TODO (#5063): Implement.
	return nil
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

// Shutdown flushes queued log records and shuts down the decorated expoter.
func (b *BatchingProcessor) Shutdown(ctx context.Context) error {
	// TODO (#5063): Implement.
	return nil
}

// ForceFlush flushes queued log records and flushes the decorated expoter.
func (b *BatchingProcessor) ForceFlush(ctx context.Context) error {
	// TODO (#5063): Implement.
	return nil
}

// BatchingOption applies a configuration to a BatchingProcessor.
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
//
// By default, if an environment variable is not set, and this option is not
// passed, 2048 will be used.
// The default value is also used when the provided value is less than one.
func WithMaxQueueSize(max int) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		// TODO (#5063): Implement.
		return cfg
	})
}

// WithExportInterval sets the maximum duration between batched exports.
//
// If the OTEL_BSP_SCHEDULE_DELAY environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 1s will be used.
// The default value is also used when the provided value is less than one.
func WithExportInterval(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		// TODO (#5063): Implement.
		return cfg
	})
}

// WithExportTimeout sets the duration after which a batched export is canceled.
//
// If the OTEL_BSP_EXPORT_TIMEOUT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 30s will be used.
// The default value is also used when the provided value is less than one.
func WithExportTimeout(d time.Duration) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		// TODO (#5063): Implement.
		return cfg
	})
}

// WithExportMaxBatchSize sets the maximum batch size of every export.
// A batch will be split into multiple exports to not exceed this size.
//
// If the OTEL_BSP_MAX_EXPORT_BATCH_SIZE environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 512 will be used.
// The default value is also used when the provided value is less than one.
func WithExportMaxBatchSize(max int) BatchingOption {
	return batchingOptionFunc(func(cfg batcherConfig) batcherConfig {
		// TODO (#5063): Implement.
		return cfg
	})
}
