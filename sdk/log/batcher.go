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
	"fmt"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
)

var _ Exporter = (*Batcher)(nil)

type batcherConfig struct {
	queueSize    int
	interval     time.Duration
	timeout      time.Duration
	maxBatchSize int
}

// Batcher is an exporter decorator
// that asynchronously exports batches of log records.
type Batcher struct {
	exporter Exporter
	cfg      batcherConfig
}

const (
	queueSizeDefault    = 2048
	intervalDefault     = time.Second
	timeoutDefault      = 30 * time.Second
	maxBatchSizeDefault = 512
)

// NewBatchingExporter decorates the provided exporter
// so that the log records are batched before exporting.
func NewBatchingExporter(exporter Exporter, opts ...BatchingOption) *Batcher {
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

	return &Batcher{exporter, cfg}
}

// Export batches provided log records.
func (b *Batcher) Export(ctx context.Context, records []*Record) error {
	return nil
}

// Shutdown flushes queued log records and shuts down the decorated expoter.
func (b *Batcher) Shutdown(ctx context.Context) error {
	return nil
}

// ForceFlush flushes queued log records and flushes the decorated expoter.
func (b *Batcher) ForceFlush(ctx context.Context) error {
	return nil
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
