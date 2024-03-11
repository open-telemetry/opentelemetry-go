// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	attributeCountLimitDefault       = 128
	attributeValueLengthLimitDefault = -1
)

// Compile-time check LoggerProvider implements log.LoggerProvider.
var _ log.LoggerProvider = (*LoggerProvider)(nil)

// LoggerProvider handles the creation and coordination of Loggers. All Loggers
// created by a LoggerProvider will be associated with the same Resource.
type LoggerProvider struct {
	embedded.LoggerProvider

	cfg providerConfig

	mu      sync.Mutex
	loggers map[instrumentation.Scope]*logger

	isShutdown atomic.Bool
}

type providerConfig struct {
	resource                  *resource.Resource
	processors                []Processor
	attributeCountLimit       int
	attributeValueLengthLimit int
}

// NewLoggerProvider returns a new and configured LoggerProvider.
//
// By default, the returned LoggerProvider is configured with the default
// Resource and no Processors. Processors cannot be added after a LoggerProvider is
// created. This means the returned MeterProvider, one created with no
// Processors, will perform no operations.
func NewLoggerProvider(opts ...LoggerProviderOption) *LoggerProvider {
	cfg := providerConfig{
		attributeCountLimit:       attributeCountLimitDefault,
		attributeValueLengthLimit: attributeValueLengthLimitDefault,
	}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if v := os.Getenv("OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT value: %w", err))
		} else {
			cfg.attributeCountLimit = n
		}
	}
	if v := os.Getenv("OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err != nil {
			otel.Handle(fmt.Errorf("invalid OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT value: %w", err))
		} else {
			cfg.attributeValueLengthLimit = n
		}
	}

	if cfg.resource == nil {
		cfg.resource = resource.Default()
	}

	return &LoggerProvider{
		cfg:     cfg,
		loggers: make(map[instrumentation.Scope]*logger),
	}
}

// Logger returns a new [log.Logger] with the provided name and configuration.
//
// This method can be called concurrently.
func (p *LoggerProvider) Logger(name string, opts ...log.LoggerOption) log.Logger {
	if p.isShutdown.Load() {
		return noop.NewLoggerProvider().Logger(name, opts...)
	}

	cfg := log.NewLoggerConfig(opts...)
	scope := instrumentation.Scope{
		Name:      name,
		Version:   cfg.InstrumentationVersion(),
		SchemaURL: cfg.SchemaURL(),
	}

	p.mu.Lock()
	l, ok := p.loggers[scope]
	if !ok {
		l = &logger{
			provider: p,
			scope:    scope,
		}
		p.loggers[scope] = l
	}
	p.mu.Unlock()

	return l
}

// Shutdown flushes queued log records and shuts down the decorated expoter.
func (p *LoggerProvider) Shutdown(ctx context.Context) error {
	wasShutdown := p.isShutdown.Swap(true)
	if wasShutdown {
		return nil
	}

	var err error
	for _, exporter := range p.cfg.processors {
		err = exporter.Shutdown(ctx)
	}
	return err
}

// ForceFlush flushes all exporters.
func (p *LoggerProvider) ForceFlush(ctx context.Context) error {
	if p.isShutdown.Load() {
		return nil
	}

	var err error
	for _, exporter := range p.cfg.processors {
		err = exporter.ForceFlush(ctx)
	}
	return err
}

// LoggerProviderOption applies a configuration option value to a LoggerProvider.
type LoggerProviderOption interface {
	apply(providerConfig) providerConfig
}

type loggerProviderOptionFunc func(providerConfig) providerConfig

func (fn loggerProviderOptionFunc) apply(c providerConfig) providerConfig {
	return fn(c)
}

// WithResource associates a Resource with a LoggerProvider. This Resource
// represents the entity producing telemetry and is associated with all Loggers
// the LoggerProvider will create.
//
// By default, if this Option is not used, the default Resource from the
// go.opentelemetry.io/otel/sdk/resource package will be used.
func WithResource(res *resource.Resource) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.resource = res
		return cfg
	})
}

// WithProcessor associates Processor with a LoggerProvider.
//
// By default, if this option is not used, the LoggerProvider will perform no
// operations; no data will be exported without a processor.
//
// Each WithProcessor creates a separate pipeline. Use custom decotarators
// for advanced scenarios such as enriching with attributes.
//
// Use NewBatchingProcessor to batch log records before they are exported.
// Use NewSimpleProcessor to synchronously export log records.
func WithProcessor(processor Processor) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.processors = append(cfg.processors, processor)
		return cfg
	})
}

// WithAttributeCountLimit sets the maximum allowed log record attribute count.
// Any attribute added to a log record once this limit is reached will be dropped.
//
// Setting this to zero means no attributes will be recorded.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit 128 will be used.
func WithAttributeCountLimit(limit int) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.attributeCountLimit = limit
		return cfg
	})
}

// AttributeValueLengthLimit sets the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value will be truncated to this length.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
// If both are set, OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT will take precedence.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit (-1) will be used.
func WithAttributeValueLengthLimit(limit int) LoggerProviderOption {
	return loggerProviderOptionFunc(func(cfg providerConfig) providerConfig {
		cfg.attributeValueLengthLimit = limit
		return cfg
	})
}
